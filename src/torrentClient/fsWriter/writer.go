package fsWriter

import (
	"os"

	"torrentClient/parser/env"
	"path"

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

	for {
		data := <- w.DataChan

		w.WriteDataToFile(data.FileName, data.Data, data.Offset)
	}
}

func (w *FsWriter) WriteDataToFile(fileName string, data []byte, offset int64) error {
	file, err := os.OpenFile(path.Join(env.GetParser().GetFilesDir(), fileName), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logrus.Errorf("Error open file: %v", err)
		return err
	}

	defer file.Close()

	if _, err := file.WriteAt(data, offset); err != nil {
		logrus.Error("Error writing to file: %v", err)
	}
	logrus.Debugf("Wrote %v bytes to file %v starting from %v", len(data), fileName, offset)

	return nil
}
