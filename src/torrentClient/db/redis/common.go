package redis

import (
	"fmt"

	"torrentClient/parser/env"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

const (
	fileSliceIndexes       = "slices"
)

type manager struct {
	conn	*redis.Client
}

var Manager manager

func (m *manager) GetSliceIndexesKey(fileName string) string {
	return fmt.Sprintf("%s:%s", fileSliceIndexes, fileName)
}

func (m *manager) AddSliceIndexForFile(fileName string, sliceByteIdx ...int64) {
	for _, idx := range sliceByteIdx {
		_, err := m.conn.SAdd(m.GetSliceIndexesKey(fileName), idx).Result()
		if err != nil {
			logrus.Errorf("Error AddSliceIndexForFile: %v", err)
		} else {
			logrus.Debugf("Added slices for file %v: %v", fileName, sliceByteIdx)
		}
	}
}

func (m *manager) DeleteSliceIndexesSet(fileName string) {
	if _, err := m.conn.Del(m.GetSliceIndexesKey(fileName)).Result(); err != nil {
		logrus.Errorf("Error deleting key: %v", err)
	}
}

func (m *manager) InitConnection() {
	m.conn = redis.NewClient(&redis.Options{
		Addr: env.GetParser().GetRedisDbAddr(),
		Password: env.GetParser().GetRedisDbPasswd(),
		DB: 0,
	})

	if err := m.conn.Ping().Err(); err != nil {
		logrus.Fatalf("Error pinging redis: %v", err)
	}
}

func (m *manager) CloseConnection() {
	if err := m.conn.Close(); err != nil {
		logrus.Errorf("Error closing redis conn: %v", err)
	}
}



