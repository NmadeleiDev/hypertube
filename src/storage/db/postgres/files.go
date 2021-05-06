package postgres

import (
	"fmt"
	"time"

	"hypertube_storage/model"

	"github.com/sirupsen/logrus"
)

func (d *manager) GetFileInfoById(id string) (info model.LoadInfo, err error)  {
	query := `
SELECT coalesce(file_name, ''), file_length, in_progress, is_loaded FROM %s WHERE file_id LIKE $1`

	err = d.conn.QueryRow(fmt.Sprintf(query, d.LoadedFilesTablePath()), id).Scan(
		&info.VideoFile.Name, &info.VideoFile.Length, &info.InProgress, &info.IsLoaded)
	return info, err
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

func (d *manager) GetInProgressFileIds() (result []string) {
	query := `
SELECT file_id FROM %s WHERE in_progress=true`

	rows, err := d.conn.Query(fmt.Sprintf(query, d.LoadedFilesTablePath()))
	if err != nil {
		logrus.Errorf("Error getting in progress files: %v", err)
		return nil
	}

	result = make([]string, 0, 10)

	for rows.Next() {
		var dest string
		if err := rows.Scan(&dest); err != nil {
			logrus.Errorf("Scan error: %v", err)
			continue
		}
		result = append(result, dest)
	}

	return result
}

func (d *manager) GetFileIdsWithWatchedUnder(under time.Time) (ids []string, names []string) {
	query := `
SELECT file_id, file_name FROM %s WHERE last_watched < $1`

	rows, err := d.conn.Query(fmt.Sprintf(query, d.LoadedFilesTablePath()), under)
	if err != nil {
		logrus.Errorf("Error getting in progress files: %v", err)
		return nil, nil
	}

	ids = make([]string, 0, 10)
	names = make([]string, 0, 10)

	for rows.Next() {
		var dest1 string
		var dest2 string
		if err := rows.Scan(&dest1, &dest2); err != nil {
			logrus.Errorf("Scan error: %v", err)
			continue
		}
		ids = append(ids, dest1)
		names = append(names, dest2)
	}

	return ids, names
}

func (d *manager) UpdateLastWatchedDate(fileId string) {
	query := `UPDATE %s SET last_watched = now()::timestamp WHERE file_id=$1`

	if _, err := d.conn.Exec(fmt.Sprintf(query, d.LoadedFilesTablePath()), fileId); err != nil {
		logrus.Errorf("Error deleteing loaded file record: %v", err)
	}
}
