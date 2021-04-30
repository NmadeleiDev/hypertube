package p2p

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

const defaultLen = 1e4

func TestSorter(t *testing.T) {
	rand.Seed(time.Now().Unix())

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	nPiecesToTest := 15
	pieces := make([]*pieceWork, nPiecesToTest)
	for i, _ := range pieces {
		pieces[i] = &pieceWork{}
		pieces[i].index = i
		copy(pieces[i].hash[:], fmt.Sprint(rand.Int()))
		pieces[i].length = defaultLen
	}
	sorter := prioritySorter{Pieces: pieces}
	topChan := sorter.InitSorter(ctx)

	updates := make(chan int64, 100)

	sorter.PriorityUpdates = updates

	go func() {
		ticker := time.NewTicker(time.Second * 5)
		for {
			select {
			case <- ctx.Done():
				return
			case <- ticker.C:
				newTop := int64(rand.Intn(nPiecesToTest))
				t.Logf("New top=%v", newTop)
				updates <- newTop
			}
		}
	}()

	for piece := range topChan {
		t.Logf("\tGot piece idx=%v", piece.index)
		time.Sleep(time.Second * time.Duration(rand.Intn(5)))
	}
}
