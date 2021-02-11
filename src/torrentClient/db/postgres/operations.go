package postgres

import (
	"fmt"
	"math"
	"os"

	"github.com/sirupsen/logrus"
)

func (d *manager) SaveFilePartsToFile(dest *os.File, fileId string)  {

	query := `
SELECT data FROM %s ORDER BY id INC`

	rows, err := d.conn.Query(fmt.Sprintf(query, d.PartsTablePathForFile(fileId)))
	if err != nil {
		logrus.Errorf("Error getting file parts: %v", err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		cont := make([]byte, 0, int64(math.Pow(2, 14)))

		if err := rows.Scan(&cont); err != nil {
			logrus.Errorf("Error scanning file part: %v", err)
			continue
		}

		if _, err := dest.Write(cont); err != nil {
			logrus.Errorf("Error writing part to file: %v", err)
		}
	}
}
