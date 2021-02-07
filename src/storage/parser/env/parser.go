package env

import (
	"sync"

	"hypertube_storage/parser/env/impl"
)

var syncOnce sync.Once
var parser Parser

type Parser interface {
	GetRedisDbAddr() string
	GetPostgresDbDsn() string
	IsDevMode() bool
	GetFilesDir() string
}

func GetParser() Parser {
	syncOnce.Do(func() {
		parser = &impl.Parser{}
	})
	return parser
}
