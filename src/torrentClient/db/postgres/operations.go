package postgres

import (
	"fmt"

	"github.com/sirupsen/logrus"
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

func (d *manager) RemoveFilePartsPlace(fileId string) {
	query := `
DROP TABLE %s`

	if _, err := d.conn.Exec(fmt.Sprintf(query, d.PartsTablePathForFile(fileId))); err != nil {
		logrus.Errorf("Error dropping table of parts: %v", err)
	} else {
		logrus.Infof("Dropped table: %v", d.PartsTablePathForFile(fileId))
	}
}
