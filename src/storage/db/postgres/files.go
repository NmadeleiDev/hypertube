package postgres

import (
	"fmt"
)

func (d *manager) GetFileInfoById(id string) (path string, inProgress, isLoaded bool, fLen int64, err error)  {
	query := `
SELECT coalesce(file_name, ''), in_progress, is_loaded, file_length FROM %s WHERE file_id LIKE $1`

	err = d.conn.QueryRow(fmt.Sprintf(query, d.LoadedFilesTablePath()), id).Scan(&path, &inProgress, &isLoaded, &fLen)
	return path, inProgress, isLoaded, fLen, err
}
