package torrentfile

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"torrentClient/db"
	"torrentClient/fsWriter"
	"torrentClient/p2p"
	"torrentClient/parser/env"

	"github.com/jackpal/bencode-go"
	"github.com/sirupsen/logrus"
)

type torrentsManager struct {
}

var activeDownloads = make(map[string]bool, 100)

func (t *torrentsManager) ReadTorrentFileFromFS(path string) (TorrentFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return TorrentFile{}, err
	}
	defer file.Close()

	return t.ParseReaderToTorrent(file)
}

func (t *torrentsManager) ReadTorrentFileFromBytes(body io.Reader) (TorrentFile, error) {
	return t.ParseReaderToTorrent(body)
}

func (t *torrentsManager) ParseReaderToTorrent(body io.Reader) (TorrentFile, error) {
	btoSingle := bencodeTorrentSingleFile{}
	btoMultiFile := bencodeTorrentMultiFiles{}

	readBody, err := io.ReadAll(body)
	if err != nil {
		logrus.Errorf("error readall body: %v", err)
		return TorrentFile{}, err
	}
	err = bencode.Unmarshal(bytes.NewBuffer(readBody), &btoSingle)
	if err != nil {
		return TorrentFile{}, err
	}
	err = bencode.Unmarshal(bytes.NewBuffer(readBody), &btoMultiFile)
	if err != nil {
		return TorrentFile{}, err
	}
	logrus.Infof("Parsed torrent!")

	var result TorrentFile

	if btoSingle.Info.Length == 0 {
		result, err = btoMultiFile.toTorrentFile()
	} else {
		result, err = btoSingle.toTorrentFile()
	}
	if err != nil {
		logrus.Errorf("Error creating torret from bto: %v", err)
	}

	logrus.Infof("Bto info: name='%v'; len=%v; files = %v; pieces = (%v)",
		result.Name, result.Length, result.Files,
		[]string{hex.EncodeToString(result.PieceHashes[0][:]),
			hex.EncodeToString(result.PieceHashes[1][:])})
	return result, nil
}

func GetManager() TorrentFilesManager {
	return &torrentsManager{}
}

func (t *TorrentFile) DownloadToFile() error {
	if active, exists := activeDownloads[t.SysInfo.FileId]; active && exists {
		return fmt.Errorf("file %v is already loading", t.SysInfo.FileId)
	} else {
		activeDownloads[t.SysInfo.FileId] = true
		defer delete(activeDownloads, t.SysInfo.FileId)
	}
	downloadCtx, downloadCancel := context.WithCancel(context.Background())
	defer downloadCancel()

	t.InitMyPeerIDAndPort()

	peersPoolObj := PeersPool{}
	peersPoolObj.InitPool()
	defer peersPoolObj.DestroyPool()
	peersPoolObj.SetTorrent(t)
	poolCtx, poolCancel := context.WithCancel(downloadCtx)
	defer poolCancel()
	go peersPoolObj.StartRefreshing(poolCtx)

	priorityManager := LoadPriority{torrentFile: t}

	torrent := p2p.TorrentMeta{
		ActiveClientsChan: peersPoolObj.ActiveClientsChan,
		PieceLoadPriorityUpdates: priorityManager.StartPriorityUpdating(downloadCtx),
		PeerID:      t.Download.MyPeerId,
		InfoHash:    t.InfoHash,
		PieceHashes: t.PieceHashes,
		PieceLength: t.PieceLength,
		Length:      t.Length,
		Name:        t.Name,
		FileId: 	 t.SysInfo.FileId,
		ResultsChan: make(chan p2p.LoadedPiece, 100),
	}

	t.CreateFileBoundariesMapping()

	db.GetFilesManagerDb().PreparePlaceForFile(torrent.FileId)
	//defer db.GetFilesManagerDb().RemoveFilePartsPlace(torrent.FileId)
	//logrus.Infof("Prepared table for parts, starting download")
	videoFile := t.getHeaviestFile()
	db.GetFilesManagerDb().SetFileNameForRecord(t.SysInfo.FileId, videoFile.EncodeFileName())

	go t.WaitForDataAndWriteToDisk(downloadCtx, torrent.ResultsChan)

	if err := torrent.Download(downloadCtx); err != nil {
		return fmt.Errorf("file download error: %v", err)
	}

	logrus.Infof("Download for %v completed!", t.SysInfo.FileId)
	return nil
}

func (t *TorrentFile) PrepareFile() (string, int64) {
	videoFile := t.getHeaviestFile()
	fsWriter.GetWriter().CreateEmptyFile(videoFile.EncodeFileName())
	db.GetFilesManagerDb().SetFileNameAndLengthForRecord(t.SysInfo.FileId, videoFile.EncodeFileName(), int64(videoFile.Length))
	return videoFile.EncodeFileName(), int64(videoFile.Length)
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
	if len(t.Files) == 1 {
		return t.Files[0]
	}

	longest := t.Files[0]

	for _, file := range t.Files {
		if file.Length > longest.Length {
			longest = file
		}
	}

	return longest
}

func (t *TorrentFile) WaitForDataAndWriteToDisk(ctx context.Context, dataParts chan p2p.LoadedPiece) {
	for {
		select {
		case <- ctx.Done():
			logrus.Debugf("Got DONE in ctx in WaitForDataAndWriteToDisk, exiting!")
			close(dataParts)
			return
		case loaded := <- dataParts:
			logrus.Debugf("Got loaded part: start=%v, len=%v", loaded.StartByte, loaded.Len)
			for _, file := range t.FileBoundariesMapping {
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

				writeTask := fsWriter.WriteTask{Data: loaded.Data[sliceStart:sliceEnd], Offset: offset, FileName: file.FileName}
				logrus.Debugf("Write task: name=%v, offset=%v, slice=(%v:%v)", writeTask.FileName, writeTask.Offset, sliceStart, sliceEnd)
				fsWriter.GetWriter().DataChan <- writeTask
			}
		}
	}
}

func (i *bencodeInfoSingleFile) hash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i)
	if err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil
}

func (t *TorrentFile) SaveLoadedPiecesToFS() error {
	start := 0

	loadChan := make(chan []byte, 100)
	writePartsChan := make(chan p2p.LoadedPiece, 100)

	go db.GetFilesManagerDb().LoadPartsForFile(t.SysInfo.FileId, loadChan)

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

func (i *bencodeInfoMultiFiles) hash() ([20]byte, error) {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i)
	if err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buf.Bytes())
	return h, nil
}

func (i *bencodeInfoSingleFile) splitPieceHashes() ([][20]byte, error) {
	hashLen := 20 // Length of SHA-1 hash
	buf := []byte(i.Pieces)
	if len(buf)%hashLen != 0 {
		err := fmt.Errorf("Received malformed pieces of length %d", len(buf))
		return nil, err
	}
	numHashes := len(buf) / hashLen
	hashes := make([][20]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
}

func (i *bencodeInfoMultiFiles) splitPieceHashes() ([][20]byte, error) {
	hashLen := 20 // Length of SHA-1 hash
	buf := []byte(i.Pieces)
	if len(buf)%hashLen != 0 {
		err := fmt.Errorf("Received malformed pieces of length %d", len(buf))
		return nil, err
	}
	numHashes := len(buf) / hashLen
	hashes := make([][20]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
}

func (t *TorrentFile) CreateFileBoundariesMapping() {
	t.FileBoundariesMapping = make([]FileBoundaries, len(t.Files))
	fileStart := 0
	for i, file := range t.Files {
		t.FileBoundariesMapping[i].FileName = file.EncodeFileName()
		t.FileBoundariesMapping[i].Index = i
		t.FileBoundariesMapping[i].Start = int64(fileStart)
		t.FileBoundariesMapping[i].End = int64(fileStart + file.Length)
		fileStart += file.Length
	}
	logrus.Infof("Calculated files borders: %v", t.FileBoundariesMapping)
}

func (bto *bencodeTorrentSingleFile) toTorrentFile() (TorrentFile, error) {
	infoHash, err := bto.Info.hash()
	if err != nil {
		return TorrentFile{}, err
	}
	pieceHashes, err := bto.Info.splitPieceHashes()
	if err != nil {
		return TorrentFile{}, err
	}
	t := TorrentFile{
		Announce:    bto.Announce,
		AnnounceList: UnfoldArray(bto.AnnounceList),
		InfoHash:    infoHash,
		PieceHashes: pieceHashes,
		PieceLength: bto.Info.PieceLength,
		Length:      bto.Info.Length,
		Name:        bto.Info.Name,
	}
	return t, nil
}

func (bto *bencodeTorrentMultiFiles) toTorrentFile() (TorrentFile, error) {
	infoHash, err := bto.Info.hash()
	if err != nil {
		return TorrentFile{}, err
	}
	pieceHashes, err := bto.Info.splitPieceHashes()
	if err != nil {
		return TorrentFile{}, err
	}
	t := TorrentFile{
		Announce:     bto.Announce,
		AnnounceList: UnfoldArray(bto.AnnounceList),
		InfoHash:     infoHash,
		PieceHashes:  pieceHashes,
		PieceLength:  bto.Info.PieceLength,
		Length:       bto.SumFilesLength(),
		Files:        bto.Info.Files,
		Name:         bto.Info.Name,
		SysInfo:      SystemInfo{},
		Download:     DownloadUtils{},
	}
	return t, nil
}
