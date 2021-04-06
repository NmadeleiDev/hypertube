package impl

import (
	"fmt"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

type Parser struct {
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

func (p *Parser) IsDevMode() bool {
	return os.Getenv("DEV_MODE") == "on"
}

func (p *Parser) GetFilesDir() string {
	return os.Getenv("FILES_DIR")
}

func (p *Parser) GetTorrentPeerPort() uint16 {
	port, err := strconv.ParseInt(os.Getenv("TORRENT_PEER_PORT"), 10, 16)
	if err != nil {
		logrus.Errorf("Error parsing peer port: %v; src: %v", err, port )
		return 6881
	} else {
		return uint16(port)
	}
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
