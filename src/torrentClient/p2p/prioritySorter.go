package p2p

import (
	"context"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

type prioritySorter struct {
	mu sync.Mutex
	Pieces	[]*pieceWork
	PriorityUpdates	<- chan int

	topPiece *pieceWork
	topPieceIdx	int
}

func (s *prioritySorter) InitSorter(ctx context.Context) (topPriorityPieceChan, returnedPiecesChan chan *pieceWork) {
	topPriorityPieceChan = make(chan *pieceWork)
	returnedPiecesChan = make(chan *pieceWork, 50)

	s.UpdateTopPieceIndex(0)
	s.RecalculateTopPiece()
	
	go func() {
		for {
			select {
			case <- ctx.Done():
				close(topPriorityPieceChan)
				close(returnedPiecesChan)
				return
			case newTopIdx := <- s.PriorityUpdates:
				logrus.Debugf("Got priority update in InitSorter: %v", newTopIdx)
				s.UpdateTopPieceIndex(newTopIdx)
				s.RecalculateTopPiece()
				logrus.Infof("Priority que after priority update new=%v:", newTopIdx)
				s.PrintPieces()
			case returnedPiece := <- returnedPiecesChan: // нам вернули часть, которую не получилось скачать
				s.InsertPieceByIdx(returnedPiece)
				s.RecalculateTopPiece()
				logrus.Infof("Priority que after recieving returned piece idx=%v:", returnedPiece.index)
				s.PrintPieces()
			default:
				topPiece := s.GetTopPiece()

				if topPiece == nil {
					continue
				}
				topPriorityPieceChan <- topPiece
				fmt.Printf("Passed top piece idx=%v\n", topPiece.index)
				s.DeletePieceByIdx(topPiece)
				s.RecalculateTopPiece()
				logrus.Infof("Priority que after giving piece idx=%v:", topPiece.index)
				s.PrintPieces()
			}
		}
	}()

	return topPriorityPieceChan, returnedPiecesChan
}


func (s *prioritySorter) UpdateTopPieceIndex(newTopIdx int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.topPieceIdx = newTopIdx
}

func (s *prioritySorter) RecalculateTopPiece() {
	s.mu.Lock()
	pieces := s.Pieces
	topIdx := s.topPieceIdx
	s.mu.Unlock()

	topPiece, _ := s.findClosestPiece(pieces, topIdx)

	s.mu.Lock()
	s.topPiece = topPiece
	s.mu.Unlock()
}

func (s *prioritySorter) GetTopPiece() *pieceWork {
	piece := &pieceWork{}

	s.mu.Lock()
	if s.topPiece == nil {
		piece = nil
	} else {
		*piece = *s.topPiece
	}
	s.mu.Unlock()
	return piece
}

func (s *prioritySorter) InsertPieceByIdx(returnedPiece *pieceWork) {
	s.mu.Lock()
	defer s.mu.Unlock()

	pLen := len(s.Pieces)
	if pLen == 0 {
		s.Pieces = append(s.Pieces, returnedPiece)
	} else {
		insertIndex := pLen
		for i, piece := range s.Pieces { // важно вставить часть на свое место
			if piece.index > returnedPiece.index {
				insertIndex = i
				break
			}
		}
		if insertIndex > pLen - 1 {
			s.Pieces = append(s.Pieces, returnedPiece)
		} else {
			s.Pieces = append(s.Pieces[:insertIndex + 1], s.Pieces[insertIndex:]...)
			s.Pieces[insertIndex] = returnedPiece
		}
		//
		//concatLeft := append(s.Pieces[:insertIndex], returnedPiece)
		//s.Pieces = append(concatLeft, s.Pieces[insertIndex:]...)
	}
}

func (s *prioritySorter) DeletePieceByIdx(pieceToDelete *pieceWork) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	pLen := len(s.Pieces)
	pieceIndex := -1
	for i, piece := range s.Pieces {
		if piece.index == pieceToDelete.index {
			pieceIndex = i
			break
		}
	}
	if pieceIndex < 0 {
		logrus.Errorf("Piece to delete (%v) not found in pieces: %v", pieceToDelete.index, s.Pieces)
		return
	}
	if pieceIndex == pLen - 1 {
		s.Pieces = s.Pieces[:pieceIndex]
	} else {
		s.Pieces = append(s.Pieces[:pieceIndex], s.Pieces[pieceIndex + 1:]...)
	}
}

func (s *prioritySorter) findClosestPiece(pieces []*pieceWork, ideal int) (piece *pieceWork, distance int) {
	numPieces := len(pieces)
	if numPieces == 0 {
		return nil, 0
	}
	if numPieces == 1 {
		return pieces[0], int(pieces[0].index) - ideal
	}

	center := numPieces / 2
	hereFoundPiece := pieces[center]
	hereFoundDistance := int(hereFoundPiece.index) - ideal

	var thereFoundPiece *pieceWork
	var thereFoundDistance int

	if int(hereFoundPiece.index) == ideal {
		return hereFoundPiece, 0
	} else if int(hereFoundPiece.index) > ideal || center + 1 >= numPieces {
		thereFoundPiece, thereFoundDistance = s.findClosestPiece(pieces[:center], ideal)
	} else {
		thereFoundPiece, thereFoundDistance = s.findClosestPiece(pieces[center + 1:], ideal)
	}

	if thereFoundDistance == 0 {
		return thereFoundPiece, thereFoundDistance
	} else if thereFoundDistance > 0 {
		if hereFoundDistance > 0 && hereFoundDistance < thereFoundDistance {
			return hereFoundPiece, hereFoundDistance
		} else {
			return thereFoundPiece, thereFoundDistance
		}
	} else {
		if hereFoundDistance < 0 && thereFoundDistance > hereFoundDistance {
			return thereFoundPiece, thereFoundDistance
		} else {
			return hereFoundPiece, hereFoundDistance
		}
	}
}

func (s *prioritySorter) PrintPieces() {
	fmt.Print("sorter pieces: [")
	s.mu.Lock()
	for _, piece := range s.Pieces {
		fmt.Printf("%v, ", piece.index)
	}
	s.mu.Unlock()
	fmt.Print("]\n")
}