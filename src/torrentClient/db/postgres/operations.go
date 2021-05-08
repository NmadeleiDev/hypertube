package postgres

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	defaultPieceSize = 1e4
)

func (d *manager) GetTorrentOrMagnetForByFileId(fileId string) ([]byte, string, bool) {
	query := `
SELECT coalesce(torrent_file, ''), coalesce(magnet_link, '') FROM %s 
WHERE file_id LIKE $1 AND (torrent_file is not null OR magnet_link is not null)`

	var torrentFile []byte
	var magnetLink string

	if err := d.conn.QueryRow(fmt.Sprintf(query, d.LoadedFilesTablePath()), fileId).Scan(&torrentFile, &magnetLink); err != nil {
		logrus.Errorf("Error getting torrent file: %v", err)
		return nil, "", false
	}

	if len(torrentFile) == 0 && len(magnetLink) == 0 {
		return torrentFile, magnetLink, false
	}

	return torrentFile, magnetLink, true
}

func (d *manager) GetFileIdsWithWatchedUnder(under time.Time) (ids []string) {
	query := `
SELECT file_id FROM %s WHERE last_watched < $1`

	rows, err := d.conn.Query(fmt.Sprintf(query, d.LoadedFilesTablePath()), under)
	if err != nil {
		logrus.Errorf("Error getting in progress files: %v", err)
		return nil
	}

	ids = make([]string, 0, 10)

	for rows.Next() {
		var dest1 string
		if err := rows.Scan(&dest1); err != nil {
			logrus.Errorf("Scan error: %v", err)
			continue
		}
		ids = append(ids, dest1)
	}

	return ids
}

func (d *manager) DeleteLoadedFileInfo(id string) error  {
	query := `
DELETE FROM %s WHERE file_id LIKE $1`

	if _, err := d.conn.Exec(fmt.Sprintf(query, d.LoadedFilesTablePath()), id); err != nil {
		logrus.Errorf("Error deleteing loaded file record: %v", err)
		return err
	}
	return nil
}

func (d *manager) PreparePlaceForFile(fileId string) {
	query := ` 
CREATE TABLE %s
(
	id bigint not null
		constraint %s_pk
			primary key,
	start bigint default 0 not null,
	size bigint default 0 not null,
	data bytea not null
)`


	if _, err := d.conn.Exec(fmt.Sprintf(query, d.PartsTablePathForFile(fileId), d.PartsTableNameForFile(fileId))); err != nil {
		logrus.Errorf("Error creating table for parts: %v", err)
	} else {
		logrus.Infof("Created table: %v", d.PartsTablePathForFile(fileId))
	}
}

func (d *manager) SaveFilePart(fileId string, part []byte, start, size, index int64) {
	query := `
INSERT INTO %s (id, start, size, data) VALUES ($1, $2, $3, $4)
`

	if _, err := d.conn.Exec(fmt.Sprintf(query, d.PartsTablePathForFile(fileId)), index, start, size, part); err != nil {
		logrus.Errorf("Error inserting file part: %v", err)
	}
}

func (d *manager) SetFileNameForRecord(fileId, name string) {
	query := `
UPDATE %s SET file_name=$1 WHERE file_id=$2`

	if _, err := d.conn.Exec(fmt.Sprintf(query, d.LoadedFilesTablePath()), name, fileId); err != nil {
		logrus.Errorf("Error saving file name: %v", err)
	}
}

func (d *manager) SetFileLengthForRecord(fileId string, length int64) {
	query := `
UPDATE %s SET file_length=$1 WHERE file_id=$2`

	if _, err := d.conn.Exec(fmt.Sprintf(query, d.LoadedFilesTablePath()), length, fileId); err != nil {
		logrus.Errorf("Error saving file len: %v", err)
	}
}

func (d *manager) SetVideoFileNameAndLengthForRecord(fileId, fileName string, length int64) {
	query := `
UPDATE %s SET file_name=$1, file_length=$2 WHERE file_id=$3`

	if _, err := d.conn.Exec(fmt.Sprintf(query, d.LoadedFilesTablePath()), fileName, length, fileId); err != nil {
		logrus.Errorf("Error saving file name and len: %v", err)
	}
}

func (d *manager) SetSrtFileNameAndLengthForRecord(fileId, fileName string, length int64) {
	query := `
UPDATE %s SET srt_file_name=$1, srt_file_length=$2 WHERE file_id=$3`

	if _, err := d.conn.Exec(fmt.Sprintf(query, d.LoadedFilesTablePath()), fileName, length, fileId); err != nil {
		logrus.Errorf("Error saving srt file name and len: %v", err)
	}
}

func (d *manager) SetInProgressStatusForRecord(fileId string, status bool) {
	query := `
UPDATE %s SET in_progress=$1 WHERE file_id=$2`

	if _, err := d.conn.Exec(fmt.Sprintf(query, d.LoadedFilesTablePath()), status, fileId); err != nil {
		logrus.Errorf("Error saving file in progess status: %v", err)
	}
}

func (d *manager) SetLoadedStatusForRecord(fileId string, status bool) {
	query := `
UPDATE %s SET is_loaded=$1 WHERE file_id=$2`

	if _, err := d.conn.Exec(fmt.Sprintf(query, d.LoadedFilesTablePath()), status, fileId); err != nil {
		logrus.Errorf("Error saving file loaded status: %v", err)
	}
}

func (d *manager) GetFileStatus(fileId string) (inProgress bool, isLoaded bool, ok bool) {
	query := `
SELECT in_progress, is_loaded FROM %s WHERE file_id=$1`

	if err := d.conn.QueryRow(fmt.Sprintf(query, d.LoadedFilesTablePath()), fileId).Scan(&inProgress, &isLoaded); err != nil {
		logrus.Errorf("Error getting file loaded status: %v", err)
		return false, false, false
	}
	return inProgress, isLoaded, true
}

func (d *manager) GetInProgressFileIds() (fileIds []string, ok bool) {
	query := `
SELECT file_id FROM %s WHERE in_progress=true`

	rows, err := d.conn.Query(fmt.Sprintf(query, d.LoadedFilesTablePath()))
	if err != nil {
		logrus.Errorf("Error getting in progress file ids: %v", err)
		return nil, false
	}

	fileIds = make([]string, 0, 10)

	for rows.Next() {
		var cont string

		if err := rows.Scan(&cont); err != nil {
			logrus.Errorf("Error scanning file id: %v", err)
			continue
		}

		fileIds = append(fileIds, cont)
	}
	return fileIds, true
}

func (d *manager) GetLoadedIndexesForFile(fileId string) []int {

	result := make([]int, 0, 400)

	query := fmt.Sprintf(`
SELECT id FROM %s`, d.PartsTablePathForFile(fileId))

	rows, err := d.conn.Query(query)
	if err != nil {
		logrus.Errorf("Error getting loaded indexes from db: %v", err)
		return nil
	}

	for rows.Next() {
		var idx int
		if err := rows.Scan(&idx); err != nil {
			logrus.Errorf("Error scan idx: %v", err)
		}
		result = append(result, idx)
	}

	return result
}

func (d *manager) LoadPartsForFile(fileId string, writeChan chan []byte) {

	defer func() {
		writeChan <- nil
	}()

	query := fmt.Sprintf(`
SELECT data FROM %s ORDER BY start`, d.PartsTablePathForFile(fileId))

	rows, err := d.conn.Query(query)
	if err != nil {
		logrus.Errorf("Error getting loaded data for %v from db: %v", fileId, err)
		return
	}

	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			logrus.Errorf("Error scan data: %v", err)
		}
		writeChan <- data
	}

	return
}

func (d *manager) GetPartDataByIdx(fileId string, idx int) ([]byte, int64, int64, bool) {
	query := fmt.Sprintf(`
SELECT data, start, size FROM %s WHERE id = $1`, d.PartsTablePathForFile(fileId))

	data := make([]byte, 0, defaultPieceSize)
	var start, size int64

	if err := d.conn.QueryRow(query, idx).Scan(&data, &start, &size); err != nil {
		logrus.Errorf("Error getting data by idx: %v", err)
		return nil, 0, 0, false
	}
	return data, start, size, true
}

func (d *manager) DropDataPartByIdx(fileId string, idx int) bool {
	query := fmt.Sprintf(`
DROP %s WHERE id = $1`, d.PartsTablePathForFile(fileId))

	if _, err := d.conn.Exec(query, idx); err != nil {
		logrus.Errorf("Error dropping data by idx: %v", err)
		return false
	}
	return true
}

func (d *manager) RemoveFilePartsPlace(fileId string) {
	query := `
DROP TABLE %s`

	if _, err := d.conn.Exec(fmt.Sprintf(query, d.PartsTablePathForFile(fileId))); err != nil {
		logrus.Errorf("Error dropping table of parts: %v", err)
	} else {
		logrus.Infof("Dropped table: %v", d.PartsTablePathForFile(fileId))
	}
}
