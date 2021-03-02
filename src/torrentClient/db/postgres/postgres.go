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

func (d *manager) GetTorrentFileForByFileId(fileId string) ([]byte, string, bool) {
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
		constraint table_name_pk
			primary key,
	start bigint default 0 not null,
	size bigint default 0 not null,
	data bytea not null
)`

	if _, err := d.conn.Exec(fmt.Sprintf(query, d.PartsTablePathForFile(fileId))); err != nil {
		logrus.Errorf("Error creating table for parts: %v", err)
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

func (d *manager) RemoveFilePartsPlace(fileId string) {
	query := `
DROP TABLE %s`

	if _, err := d.conn.Exec(fmt.Sprintf(query, d.PartsTablePathForFile(fileId))); err != nil {
		logrus.Errorf("Error dropping table of parts: %v", err)
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
