package postgres

import (
	"fmt"
	"math"
	"os"

	"github.com/sirupsen/logrus"
)

func (d *manager) SaveFilePartsToFile(dest *os.File, fileId string, fileStart int, fileLength int) error {

	logrus.Debugf("Loading file from db. FileId=%v, fileStart=%v, fileLength=%v", fileId, fileStart, fileLength)

	tableName := d.PartsTablePathForFile(tableNamePrefix + fileId)
	query := `
SELECT id, start, size, data FROM ` + tableName + `
WHERE start
    BETWEEN $1 - ` + tableName + `.size AND $2
ORDER BY id`

	rows, err := d.conn.Query(query, fileStart, fileStart+fileLength)
	if err != nil {
		logrus.Errorf("Error getting file parts: %v", err)
		return fmt.Errorf("query error: %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var idx int
		var partStart int
		var partSize int
		dataCont := make([]byte, 0, int64(math.Pow(2, 14)))

		if err := rows.Scan(&idx, &partStart, &partSize, &dataCont); err != nil {
			logrus.Errorf("Error scanning file part: %v", err)
			continue
		}

		logrus.Debugf("Part data: id=%v, fileStart=%v, partSize=%v, data_len=%v", idx, partStart, partSize, len(dataCont))

		bufSliceStart := 0
		if partStart < fileStart {
			bufSliceStart = partStart - fileStart
		}
		bufSliceEnd := partSize
		if fileLength < bufSliceEnd + bufSliceStart {
			bufSliceEnd = bufSliceEnd + bufSliceStart - fileLength
		}

		writeOffset := partStart - fileStart
		if writeOffset < 0 {
			writeOffset = 0
		}

		logrus.Debugf("Writing part to fs. buf_len_to_write=%v, bufSliceStart=%v, bufSliceEnd=%v, offset=%v",
			len(dataCont[bufSliceStart:bufSliceEnd]), bufSliceStart, bufSliceEnd, writeOffset)

		if _, err := dest.WriteAt(dataCont[bufSliceStart:bufSliceEnd], int64(writeOffset)); err != nil {
			logrus.Errorf("Error writing part to file: %v", err)
			return fmt.Errorf("write error: %v", err)
		}
	}
	return nil
}
