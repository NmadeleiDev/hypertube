package impl

import (
	"fmt"
	"os"
)

type Parser struct {
}

func (p *Parser) GetFilesDir() string {
	return os.Getenv("FILES_DIR")
}

func (p *Parser) GetRedisDbAddr() string {
	return fmt.Sprintf(
		"%v:%v",
		os.Getenv("REDIS_HOST"),
		os.Getenv("REDIS_PORT"))
}

func (p *Parser) GetRedisDbPasswd() string {
	return os.Getenv("REDIS_PASSWORD")
}


func (p *Parser) GetPostgresDbDsn() string {
	return fmt.Sprintf(
		"host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))
}

func (p *Parser) GetLoaderServiceHost() string {
	return os.Getenv("LOADER_SERVICE_ADDR")
}

func (p *Parser) IsDevMode() bool {
	return os.Getenv("DEV_MODE") == "on"
}

