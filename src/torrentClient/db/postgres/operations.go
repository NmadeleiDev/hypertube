package postgres

import (
	"fmt"
	"math"
	"os"

	"github.com/sirupsen/logrus"
)

func (d *manager) SaveFilePartsToFile(dest *os.File, fileId string, start int, length int) error {

	logrus.Debugf("Loading file from db. FileId=%v, start=%v, length=%v", fileId, start, length)

	tableName := d.PartsTablePathForFile(filePartsTablePrefix + fileId)
	query := `
SELECT id, start, size, data FROM ` + tableName + `
WHERE start
    BETWEEN $1 - ` + tableName + `.size AND $2
ORDER BY id`

	rows, err := d.conn.Query(query, start, start + length)
	if err != nil {
		logrus.Errorf("Error getting file parts: %v", err)
		return fmt.Errorf("query error: %v", err)
	}

	defer rows.Close()

	firstWriteStart := -1

	for rows.Next() {
		var idx int
		var fStart int
		var size int
		dataCont := make([]byte, 0, int64(math.Pow(2, 14)))

		if err := rows.Scan(&idx, &fStart, &size, &dataCont); err != nil {
			logrus.Errorf("Error scanning file part: %v", err)
			continue
		}

		logrus.Debugf("Part data: id=%v, start=%v, size=%v, data_len=%v", idx, fStart, size, len(dataCont))

		if firstWriteStart < 0 {
			firstWriteStart = fStart
		}

		writeStart := 0
		if fStart < start {
			writeStart = start
		}
		writeEnd := size
		if length < writeEnd - firstWriteStart {
			writeEnd = length % size
		}

		logrus.Debugf("Writing part ot fs. buf_len=%v, writeStart=%v, writeEnd=%v, offset=%v",
			len(dataCont), writeStart, writeEnd, fStart - firstWriteStart)

		if _, err := dest.WriteAt(dataCont[writeStart:writeEnd], int64(fStart - firstWriteStart)); err != nil {
			logrus.Errorf("Error writing part to file: %v", err)
			return fmt.Errorf("write error: %v", err)
		}
	}
	return nil
}
