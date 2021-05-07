package postgres

import (
	"fmt"

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

func (d *manager) UpdateLastWatchedDate(fileId string) {
	query := `UPDATE %s SET last_watched = now()::timestamp WHERE file_id=$1`

	if _, err := d.conn.Exec(fmt.Sprintf(query, d.LoadedFilesTablePath()), fileId); err != nil {
		logrus.Errorf("Error deleteing loaded file record: %v", err)
	}
}
