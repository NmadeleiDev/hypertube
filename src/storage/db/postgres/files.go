package postgres

import (
	"fmt"
)

func (d *manager) GetFilePathById(id string) (path string, err error)  {
	query := `
SELECT file_name FROM %s WHERE file_id LIKE $1 AND file_name IS NOT NULL`

	err = d.conn.QueryRow(fmt.Sprintf(query, d.LoadedFilesTablePath()), id).Scan(&path)
	return path, err
}
