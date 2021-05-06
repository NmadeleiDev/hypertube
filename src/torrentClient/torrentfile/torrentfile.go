package torrentfile

import (
	"context"
	"crypto/rand"
	"fmt"

	"torrentClient/db"
	"torrentClient/fsWriter"
	"torrentClient/loadMaster"
	"torrentClient/p2p"
	"torrentClient/parser/env"

	"github.com/sirupsen/logrus"
)

type torrentsManager struct {
}

func (t *TorrentFile) DownloadToFile() error {
	downloadCtx, downloadCancel := context.WithCancel(context.Background())
	defer downloadCancel()

	fileId := t.GetFileId()
	db.GetFilesManagerDb().SetInProgressStatusForRecord(fileId, true)
	defer db.GetFilesManagerDb().SetInProgressStatusForRecord(fileId, false)

	loadEntry, ok := loadMaster.GetMaster().AddLoadEntry(fileId, downloadCancel, len(t.PieceHashes))
	if !ok {
		logrus.Debugf("Failed to add loading entry (propably, file is already in progress)")
		return fmt.Errorf("failed to add loading entry")
	}

	t.InitMyPeerIDAndPort()

	peersPoolObj := PeersPool{}
	peersPoolObj.Init(t)
	defer peersPoolObj.DestroyPool()
	poolCtx, poolCancel := context.WithCancel(downloadCtx)
	defer poolCancel()
	go peersPoolObj.StartRefreshing(poolCtx)

	priorityManager := LoadPriority{torrentFile: t}

	torrent := p2p.TorrentMeta{
		ActivatedClientsChan:     peersPoolObj.ClientMaker.InitializedPeersChan,
		DeadPeersChan: 			peersPoolObj.ClientMaker.DeadPeersChan,
		PieceLoadPriorityUpdates: priorityManager.StartPriorityUpdating(downloadCtx),
		PeerID:                   t.Download.MyPeerId,
		InfoHash:                 t.InfoHash,
		PieceHashes:              t.PieceHashes,
		PieceLength:              t.PieceLength,
		Length:                   t.Length,
		Name:                     t.Name,
		FileId:                   fileId,
		ResultsChan:              make(chan p2p.LoadedPiece, 100),
		LoadStats:                loadEntry,
	}

	t.CreateFileBoundariesMapping()

	db.GetFilesManagerDb().PreparePlaceForFile(torrent.FileId)
	//defer db.GetFilesManagerDb().RemoveFilePartsPlace(torrent.FileId)
	//logrus.Infof("Prepared table for parts, starting download")
	videoFile := t.getHeaviestFile()
	db.GetFilesManagerDb().SetFileNameForRecord(fileId, videoFile.EncodeFileName())

	go t.WaitForDataAndWriteToDisk(downloadCtx, torrent.ResultsChan)

	if err := torrent.Download(downloadCtx); err != nil {
		return fmt.Errorf("file download error: %v", err)
	}
	db.GetFilesManagerDb().SetLoadedStatusForRecord(fileId, true)
	logrus.Infof("Download for %v completed!", fileId)
	return nil
}

func (t *TorrentFile) PrepareDownload() (string, int64) {
	videoFile := t.getHeaviestFile()
	subtitlesFiles := t.getSubtitlesFiles()
	for _, subFile := range subtitlesFiles {
		fsWriter.GetWriter().CreateEmptyFile(subFile.EncodeFileName())
	}
	fsWriter.GetWriter().CreateEmptyFile(videoFile.EncodeFileName())
	db.GetFilesManagerDb().SetVideoFileNameAndLengthForRecord(t.GetFileId(), videoFile.EncodeFileName(), int64(videoFile.Length))
	return videoFile.EncodeFileName(), int64(videoFile.Length)
}

func (t *TorrentFile) GetVideoFileName() string {
	videoFile := t.getHeaviestFile()
	return videoFile.EncodeFileName()
}

func (t *TorrentFile) GetVideoFileLength() int {
	videoFile := t.getHeaviestFile()
	return videoFile.Length
}

func (t *TorrentFile) InitMyPeerIDAndPort() {
	var peerID [20]byte

	_, err := rand.Read(peerID[:])
	if err != nil{
		logrus.Errorf("read rand error: %v", err)
		copy(peerID[:], []byte("you suck")[:20])
	}
	t.Download.MyPeerId = peerID
	t.Download.MyPeerPort = env.GetParser().GetTorrentPeerPort()
}

func (t *TorrentFile) getHeaviestFile() bencodeTorrentFile {
	allFiles := t.GetFiles()

	if len(allFiles) == 1 {
		return allFiles[0]
	}

	longest := allFiles[0]

	for _, file := range allFiles {
		if file.Length > longest.Length {
			longest = file
		}
	}

	return longest
}

func (t *TorrentFile) getSubtitlesFiles() []bencodeTorrentFile {
	allFiles := t.GetFiles()
	res := make([]bencodeTorrentFile, len(allFiles))
	for _, file := range allFiles {
		if file.Extension() == "srt" {
			res = append(res, file)
		}
	}

	return res
}

func (t *TorrentFile) WaitForDataAndWriteToDisk(ctx context.Context, dataParts chan p2p.LoadedPiece) {
	fileBoundariesMapping := t.GetFileBoundariesMapping()

	for {
		select {
		case <- ctx.Done():
			logrus.Debugf("Got DONE in ctx in WaitForDataAndWriteToDisk, exiting!")
			close(dataParts)
			return
		case loaded := <- dataParts:
			logrus.Debugf("Got loaded part: start=%v, len=%v", loaded.StartByte, loaded.Len)
			for _, file := range fileBoundariesMapping {
				if loaded.StartByte > file.End || loaded.StartByte + loaded.Len < file.Start {
					//logrus.Debugf("Skipping '%v' write due to (%v, %v); (%v, %v)", file.FileName, loaded.StartByte > file.End, loaded.StartByte + loaded.Len < file.Start, file.Start, file.End)
					continue
				}

				sliceStart := file.Start - loaded.StartByte
				sliceEnd := loaded.Len
				fileEndBias := loaded.StartByte + loaded.Len - file.End
				if fileEndBias > 0 {
					sliceEnd -= fileEndBias
				}

				if sliceStart < 0 {
					sliceStart = 0
				}
				if sliceEnd < 0 {
					sliceEnd = loaded.Len
				} else if sliceEnd < sliceStart {
					logrus.Warnf("sliceEnd < sliceStart! %v %v; start=%v, len=%v; file: %v", sliceEnd, sliceStart, loaded.StartByte, loaded.Len, file)
				}

				offset := loaded.StartByte - file.Start
				if offset <= 0 {
					offset = 0
				}

				//writeTask := fsWriter.WriteTask{Data: loaded.Data[sliceStart:sliceEnd], Offset: offset, FileName: file.FileName}
				//logrus.Debugf("Write task: name=%v, offset=%v, slice=(%v:%v)", writeTask.FileName, writeTask.Offset, sliceStart, sliceEnd)
				data := loaded.Data[sliceStart:sliceEnd]
				fsWriter.GetWriter().AddToWriteQue(file.FileName, data, offset)
			}
		}
	}
}

func (t *TorrentFile) SaveLoadedPiecesToFS() error {
	start := 0

	loadChan := make(chan []byte, 100)
	writePartsChan := make(chan p2p.LoadedPiece, 100)

	go db.GetFilesManagerDb().LoadPartsForFile(t.GetFileId(), loadChan)

	loadCtx := context.TODO()

	go t.WaitForDataAndWriteToDisk(loadCtx, writePartsChan)

	for {
		loadedData := <- loadChan
		if loadedData == nil {
			logrus.Infof("Loaded nil, exiting")
			break
		}
		logrus.Debugf("Got %v bytes to save. Start=%v", len(loadedData), start)

		writePartsChan <- p2p.LoadedPiece{StartByte: int64(start), Len: int64(len(loadedData)), Data: loadedData}
		start += len(loadedData)
	}

	loadCtx.Done()
	return nil
}

func (t *TorrentFile) CreateFileBoundariesMapping() {
	files := t.GetFiles()
	t.SetFileBoundariesMapping(make([]FileBoundaries, len(files)))

	fileStart := 0
	for i, file := range files {
		t.mu.Lock()
		t.FileBoundariesMapping[i].FileName = file.EncodeFileName()
		t.FileBoundariesMapping[i].Index = i
		t.FileBoundariesMapping[i].Start = int64(fileStart)
		t.FileBoundariesMapping[i].End = int64(fileStart + file.Length)
		t.mu.Unlock()
		fileStart += file.Length
	}
	logrus.Infof("Calculated files borders: %v", t.GetFileBoundariesMapping())
}
