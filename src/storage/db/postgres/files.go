package postgres

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

func (d *manager) GetFileInfoById(id string) (path string, inProgress, isLoaded bool, fLen int64, err error)  {
	query := `
SELECT coalesce(file_name, ''), in_progress, is_loaded, file_length FROM %s WHERE file_id LIKE $1`

	err = d.conn.QueryRow(fmt.Sprintf(query, d.LoadedFilesTablePath()), id).Scan(&path, &inProgress, &isLoaded, &fLen)
	return path, inProgress, isLoaded, fLen, err
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

func (d *manager) GetFileIdsWithLoadedDateUnder(under time.Time) (result []string) {
	query := `
SELECT file_id FROM %s WHERE loaded_date < $1`

	rows, err := d.conn.Query(fmt.Sprintf(query, d.LoadedFilesTablePath()), under)
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
