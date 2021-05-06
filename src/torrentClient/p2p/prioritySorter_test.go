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
	topChan, returnPieces := sorter.InitSorter(ctx)

	updates := make(chan int, 100)

	sorter.PriorityUpdates = updates

	go func() {
		timer := time.NewTimer(time.Second * 5)
		for {
			select {
			case <- ctx.Done():
				return
			case <- timer.C:
				newTop := rand.Intn(nPiecesToTest)
				updates <- newTop
				t.Logf("sent new top=%v", newTop)
				timer.Reset(time.Second * time.Duration(4 + rand.Intn(5)))
			}
		}
	}()

	gotten := make(map[int]bool, nPiecesToTest)

	for piece := range topChan {
		t.Logf("Got piece idx=%v", piece.index)

		if isTrue, exists := gotten[piece.index]; isTrue && exists {
			t.Fatalf("FUCK: %v\n; %v\n; %v\n", piece.index, gotten, sorter.Pieces)
		}
		gotten[piece.index] = true
		time.Sleep(time.Second * time.Duration(1 + rand.Intn(5)))

		if rand.Intn(10) > 5 {
			t.Logf("Returning piece idx=%v", piece.index)
			returnPieces <- piece
			gotten[piece.index] = false
		}

		//t.Logf("Gotten: %v", gotten)
		printTrue(gotten, "here")
		sorter.PrintPieces()
	}
}

func printTrue(src map[int]bool, name string) {
	fmt.Printf("%v pieces: [", name)
	for idx, piece := range src {
		if piece {
			fmt.Printf("%v, ", idx)
		}
	}
	fmt.Print("]\n")
}
