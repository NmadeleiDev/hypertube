package env

import (
	"sync"

	"torrent_client/parser/env/impl"
)

var syncOnce sync.Once
var parser Parser

type Parser interface {
	GetRedisDbAddr() string
	GetRedisDbPasswd() string
	IsDevMode() bool
	GetFilesDir() string
	GetPostgresDbDsn() string
}

func GetParser() Parser {
	syncOnce.Do(func() {
		parser = &impl.Parser{}
	})
	return parser
}
