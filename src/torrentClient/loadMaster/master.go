package loadMaster

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

type LoadsMaster struct {
	mu sync.Mutex
	loads		map[string]*LoadEntry
}

var master LoadsMaster

func GetMaster() *LoadsMaster {
	return &master
}

type LoadEntry struct {
	mu sync.Mutex
	ExecutionCtxCancel	context.CancelFunc
	NumOfActivePeers	int
	TotalPieces		int
	DonePieces		int
}

func (m *LoadsMaster) Init() {
	m.loads = make(map[string]*LoadEntry, 100)
}

func (m *LoadsMaster) AddLoadEntry(fileId string, entry *LoadEntry) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.loads[fileId]; exists {
		return false
	}
	logrus.Debugf("Added load entry for %v", fileId)
	m.loads[fileId] = entry

	return true
}

func (m *LoadsMaster) StopLoad(fileId string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if entry, exists := m.loads[fileId]; exists {
		entry.ExecutionCtxCancel()
		delete(m.loads, fileId)
		return true
	} else {
		return false
	}
}

func (m *LoadsMaster) GetStatsForEntry(fileId string) (LoadEntry, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if entry, exists := m.loads[fileId]; exists {
		logrus.Debugf("Returning stats for %v: %v", fileId, entry)
		return LoadEntry{NumOfActivePeers: entry.NumOfActivePeers, DonePieces: entry.DonePieces, TotalPieces: entry.TotalPieces}, true
	} else {
		return LoadEntry{}, false
	}
}

func (l *LoadEntry) GetLoadedPercent() int {
	l.mu.Lock()
	defer l.mu.Unlock()

	return int((float64(l.DonePieces) / float64(l.TotalPieces)) * 100)
}

func (l *LoadEntry) CountDone() int {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.DonePieces
}

func (l *LoadEntry) IncrDone() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.DonePieces ++
}

func (l *LoadEntry) SetTotalPieces(val int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.TotalPieces = val
}


func (l *LoadEntry) IncrActivePeers() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.NumOfActivePeers ++
}

func (l *LoadEntry) DecrActivePeers() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.NumOfActivePeers --
}

