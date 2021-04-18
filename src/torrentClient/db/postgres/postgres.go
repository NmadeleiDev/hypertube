package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var Manager manager

type manager struct {
	conn *sqlx.DB

	schemaName       string
}

const filePartsTablePrefix = "part_"

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

func (d *manager) SaveFileNameForReadyFile(fileId, name string) {
	query := `
UPDATE %s SET file_name=$1 WHERE file_id=$2`

	if _, err := d.conn.Exec(fmt.Sprintf(query, d.LoadedFilesTablePath()), name, fileId); err != nil {
		logrus.Errorf("Error saving file name: %v", err)
	}
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

	tableName := filePartsTablePrefix + fileId

	if _, err := d.conn.Exec(fmt.Sprintf(query, d.PartsTablePathForFile(tableName), tableName)); err != nil {
		logrus.Errorf("Error creating table for parts: %v", err)
	} else {
		logrus.Infof("Created table: %v", tableName)
	}
}

func (d *manager) SaveFilePart(fileId string, part []byte, start, size, index int64) {
	query := `
INSERT INTO %s (id, start, size, data) VALUES ($1, $2, $3, $4)
`

	tableName := filePartsTablePrefix + fileId

	if _, err := d.conn.Exec(fmt.Sprintf(query, d.PartsTablePathForFile(tableName)), index, start, size, part); err != nil {
		logrus.Errorf("Error inserting file part: %v", err)
	}
}

func (d *manager) GetLoadedIndexesForFile(fileId string) []int {
	tableName := filePartsTablePrefix + fileId

	result := make([]int, 0, 400)

	query := fmt.Sprintf(`
SELECT id FROM %s`, d.PartsTablePathForFile(tableName))

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
	tableName := filePartsTablePrefix + fileId

	query := fmt.Sprintf(`
SELECT data FROM %s ORDER BY start`, d.PartsTablePathForFile(tableName))

	rows, err := d.conn.Query(query)
	if err != nil {
		logrus.Errorf("Error getting loaded indexes from db: %v", err)
		return
	}

	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			logrus.Errorf("Error scan idx: %v", err)
		}
		writeChan <- data
	}

	return
}

func (d *manager) RemoveFilePartsPlace(fileId string) {
	query := `
DROP TABLE %s`

	tableName := filePartsTablePrefix + fileId

	if _, err := d.conn.Exec(fmt.Sprintf(query, d.PartsTablePathForFile(tableName))); err != nil {
		logrus.Errorf("Error dropping table of parts: %v", err)
	} else {
		logrus.Infof("Dropped table: %v", tableName)
	}
}

func (d *manager) InitTables() {
	query := `create schema if not exists %v`

	if _, err := d.conn.Exec(fmt.Sprintf(query, d.schemaName));err != nil {
		logrus.Fatalf("Error creating schema: %v", err)
	}

}

func (d *manager) InitConnection(connStr string) {
	db := manager{
		schemaName:       "file_parts",
	}
	conn, err := sqlx.Open("postgres", connStr)
	if err != nil {
		logrus.Fatalf("Error connecting to database: ", err)
	}
	if err := conn.Ping(); err != nil {
		logrus.Fatalf("Error pinging db: %v", err)
	}
	logrus.Debugf("Connected to %v; db conf: %v", connStr, db)
	db.conn = conn
	Manager = db
}

func (d *manager) CloseConnection() {
	_ = d.conn.Close()
}
