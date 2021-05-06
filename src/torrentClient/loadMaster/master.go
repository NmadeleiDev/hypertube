package loadMaster

import (
	"context"
	"fmt"
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
	mu                 sync.Mutex
	ExecutionCtxCancel context.CancelFunc
	NumOfActivePeers   int
	totalPieces        int
	ProcessedPieces    map[int]bool
	
	DonePieces []int
	InProgressPieces []int
}

type LoadStat struct {
	NumOfActivePeers	int
	TotalPieces		int
	DonePieces []int
	InProgressPieces []int
}

func (m *LoadsMaster) Init() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.loads = make(map[string]*LoadEntry, 100)
}

func (m *LoadsMaster) AddLoadEntry(fileId string, ctxCancel context.CancelFunc, totalPieces int) (*LoadEntry, bool) {
	m.mu.Lock()
	if _, exists := m.loads[fileId]; exists {
		return nil, false
	}
	m.mu.Unlock()

	loadEntry := &LoadEntry{
		ExecutionCtxCancel: ctxCancel,
		totalPieces:        totalPieces,
		ProcessedPieces:    make(map[int]bool, totalPieces),
	}
	logrus.Debugf("Added load entry for %v", fileId)

	m.mu.Lock()
	m.loads[fileId] = loadEntry
	m.mu.Unlock()

	return loadEntry, true
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

func (m *LoadsMaster) GetStatsForEntry(fileId string) (result LoadStat, ok bool) {
	logrus.Debugf("Getting stats for %v", fileId)
	m.mu.Lock()
	entry, exists := m.loads[fileId]
	m.mu.Unlock()

	if exists {
		logrus.Debugf("fuck1 stats for %v", fileId)
		totalPieces := entry.TotalPieces()
		logrus.Debugf("fuck2 stats for %v", fileId)
		nDone := entry.GetLoadedPieces()
		logrus.Debugf("fuck3 stats for %v", fileId)
		inProgress := entry.GetInProgressPieces()

		result, ok = LoadStat{
			NumOfActivePeers: entry.NumOfActivePeers, 
			DonePieces:       nDone,
			InProgressPieces: inProgress,
			TotalPieces:      totalPieces}, true
	}
	return result, ok
}

func (l *LoadEntry) GetLoadedPercent() int {
	done := float64(l.CountDone())
	total := float64(l.TotalPieces())

	return int(done / total) * 100
}

func (l *LoadEntry) GetLoadedPieces() (res []int) {
	total := l.TotalPieces()

	l.mu.Lock()
	defer l.mu.Unlock()

	res = make([]int, 0, total)
	for k, v := range l.ProcessedPieces {
		if v {
			res = append(res, k)
		}
	}
	return res
}

func (l *LoadEntry) GetInProgressPieces() (res []int) {
	total := l.TotalPieces()

	l.mu.Lock()
	defer l.mu.Unlock()

	res = make([]int, 0, total)
	for k, v := range l.ProcessedPieces {
		if !v {
			res = append(res, k)
		}
	}
	return res
}

func (l *LoadEntry) GetNumOfActivePeers() (res int) {
	l.mu.Lock()
	res = l.NumOfActivePeers
	l.mu.Unlock()
	return res
}

func (l *LoadEntry) CountDone() (count int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, v := range l.ProcessedPieces {
		if v {
			count++
		}
	}
	return count
}

func (l *LoadEntry) SetDone(idx int) (err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	ready, exists := l.ProcessedPieces[idx]
	if !exists {
		err = fmt.Errorf("tryied to set idx=%v as done, but it is not processed yet", idx)
	} else if ready {
		err = fmt.Errorf("tryied to set idx=%v as done, but it is already set as done", idx)
	} else {
		l.ProcessedPieces[idx] = true
		err = nil
	}
	return err
}

func (l *LoadEntry) ForceSetDone(idx int) {
	l.mu.Lock()
	l.ProcessedPieces[idx] = true
	l.mu.Unlock()
}

func (l *LoadEntry) AddProcessed(idx int) (err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	_, exists := l.ProcessedPieces[idx]
	if exists {
		err = fmt.Errorf("tryied to add idx=%v to process, but it is already there", idx)
	} else {
		l.ProcessedPieces[idx] = false
		err = nil
	}
	return err
}

func (l *LoadEntry) DeleteProcessed(idx int) (err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	_, exists := l.ProcessedPieces[idx]
	if !exists {
		err = fmt.Errorf("tryied to delete idx=%v from process, but it is not there: %v", idx, l.ProcessedPieces)
	} else {
		delete(l.ProcessedPieces, idx)
		logrus.Debugf("Deleted processed piece idx=%v", idx)
		err = nil
	}
	return err
}

func (l *LoadEntry) TotalPieces() int {
	l.mu.Lock()
	val := l.totalPieces
	l.mu.Unlock()

	return val
}

func (l *LoadEntry) SetTotalPieces(val int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.totalPieces = val
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

