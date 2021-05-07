package fsWriter

import (
	"os"
	"path"

	"torrentClient/db"
	"torrentClient/parser/env"

	"github.com/sirupsen/logrus"
)

type FsWriter struct {
	DataChan	chan WriteTask
}

type WriteTask struct {
	FileName	string
	Data	[]byte
	Offset	int64
}

var fsWriter FsWriter

func GetWriter() *FsWriter {
	return &fsWriter
}

func (w *FsWriter) StartWaitingForData() {
	w.DataChan = make(chan WriteTask, 1000)

	for data := range w.DataChan {
		w.WriteDataToFile(data.FileName, data.Data, data.Offset)
	}
}

func (w *FsWriter) AddToWriteQue(fileName string, data []byte, offset int64) {
	task := WriteTask{Data: data, Offset: offset, FileName: fileName}

	w.DataChan <- task
}

func (w *FsWriter) CreateEmptyFile(fileName string) bool {
	if fileExists(fileName) {
		return true
	}

	if  f, err := os.Create(path.Join(env.GetParser().GetFilesDir(), fileName)); err != nil {
		logrus.Errorf("Error creating file %v: %v", fileName, err)
		return false
	} else {
		f.Close()
		return true
	}
}

func (w *FsWriter) WriteDataToFile(fileName string, data []byte, offset int64) error {
	file, err := os.OpenFile(path.Join(env.GetParser().GetFilesDir(), fileName), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logrus.Errorf("Error open file: %v", err)
		return err
	}

	defer file.Close()
	dataLen := len(data)

	db.GetLoadedStateDb().AddSliceIndexForFile(fileName, offset, offset + int64(dataLen))

	if _, err := file.WriteAt(data, offset); err != nil {
		logrus.Error("Error writing to file: %v", err)
	}

	logrus.Debugf("Wrote %v bytes to file %v starting from %v", dataLen, fileName, offset)
	return nil
}

func (w *FsWriter) RemoveFile(fileName string) bool {
	if err := os.Remove(path.Join(env.GetParser().GetFilesDir(), fileName)); err != nil {
		logrus.Errorf("Error deleting file %v: %v", fileName, err)
		return false
	}
	return true
}
