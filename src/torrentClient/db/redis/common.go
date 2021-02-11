package redis

import (
	"fmt"

	"torrent_client/parser/env"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

const (
	activeDownloadsKey       = "downloads"
	loadedPartsChannelPrefix = "loaded"
	loadedPartRecordPrefix	= "complete"
)

type manager struct {
	conn	*redis.Client
}

var Manager manager

func (m *manager) GetActiveDownloads() []string {
	res, err := m.conn.SMembers(activeDownloadsKey).Result()
	if err != nil {
		logrus.Errorf("Error getting active downloads: %v", err)
	}
	return res
}

func (m *manager) CheckIfFileIsActiveLoading(file string) bool {
	res, err := m.conn.SIsMember(activeDownloadsKey, file).Result()
	if err != nil {
		logrus.Errorf("Error checking active downloads: %v", err)
	}
	return res
}

func (m *manager) AddFileIdToActiveDownloads(id string) {
	_, err := m.conn.SAdd(activeDownloadsKey, id).Result()
	if err != nil {
		logrus.Errorf("Error checking active downloads: %v", err)
	}
}

func (m *manager) AnnounceLoadedPart(fileId, partId string, start, size int64) {
	if _, err := m.conn.Publish(fmt.Sprintf("%s:%s", loadedPartsChannelPrefix, fileId), fmt.Sprintf("id=%s&start=%d&size=%d", partId, start, size)).Result(); err != nil {
		logrus.Errorf("Error publish message: %v", err )
	}
}

func (m *manager) SaveLoadedPartInfo(fileId, partId string, start, size int64) {
	data := map[string]interface{}{
		"start": start,
		"size": size,
	}
	if _, err := m.conn.HMSet(fmt.Sprintf("%s:%s:%s", loadedPartRecordPrefix, fileId, partId), data).Result(); err != nil {
		logrus.Errorf("Error saving loaded part record: %v", err)
	}
}

func (m *manager) CleanLoadingLogsForFile(file string) {
	cur := uint64(0)
	//keys := make([]string, 0, 100)

	for {
		var res []string
		var err error

		res, cur, err = m.conn.Scan(cur, fmt.Sprintf("%s:%s:*", loadedPartRecordPrefix, file), 10).Result()
		if err != nil {
			logrus.Errorf("Error saving loaded part record: %v", err)
		}

		m.conn.Del(res...)
	}
}

func (m *manager) DeleteFileFromActiveDownloads(file string) {
	if _, err := m.conn.SRem(activeDownloadsKey, file).Result(); err != nil {
		logrus.Errorf("Error removing id from active downloads: %v", err)
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



