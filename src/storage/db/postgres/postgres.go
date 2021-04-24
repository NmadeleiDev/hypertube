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
	loadedFilesTable string
}

func (d *manager) InitTables() {
	query := `create schema if not exists %v`

	if _, err := d.conn.Exec(fmt.Sprintf(query, d.schemaName));err != nil {
		logrus.Fatalf("Error creating schema: %v", err)
	}

	query = `create table if not exists %s
(
    id           serial                                      not null
        constraint %s_pk
            primary key,
    file_id      varchar(64)                                 not null,
    torrent_file bytea,
    magnet_link  varchar(2048) default NULL::character varying,
    file_name    varchar(256)  default NULL::character varying,
	file_length  bigint			not null default 0,
    comment      varchar(256)  default ''::character varying not null,
    error        varchar(1024) default ''::character varying not null,
    in_progress  boolean       default false                 not null,
    is_loaded    boolean       default false                 not null
)`

	if _, err := d.conn.Exec(fmt.Sprintf(query, d.LoadedFilesTablePath(), d.loadedFilesTable)); err != nil {
		logrus.Fatalf("Error creating table %v: %v", d.LoadedFilesTablePath(), err)
	}

}

func (d *manager) InitConnection(connStr string) {
	db := manager{
		schemaName:       "hypertube",
		loadedFilesTable: "loaded_files",
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
