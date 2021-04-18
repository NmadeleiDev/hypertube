package filesReader

import (
	"fmt"
	"io"
	"os"
	"path"

	"hypertube_storage/dao"
	"hypertube_storage/model"
	"hypertube_storage/parser/env"

	"github.com/sirupsen/logrus"
)

const readLimit = 10e6
var filesDir = env.GetParser().GetFilesDir()

type fileReader struct {

}

// в данный момент учитывает только начало range, возвращает всегда
func (f *fileReader) GetFileInRange(fileName string, description *model.FileRangeDescription) (result io.Reader, totalLength int64, err error) {
	file, err := os.Open(path.Join(filesDir, fileName))
	if err != nil {
		logrus.Errorf("Error open file: %v", err)
		return nil, 0, err
	}

	//defer file.Close()

	info, err := file.Stat()
	if err != nil {
		logrus.Errorf("Error getting file info: %v", err)
		return nil, 0, err
	}

	start := int64(0)
	if description.Start > 0 {
		start = description.Start
	} else {
		description.Start = start
	}

	if start + readLimit > info.Size() {
		description.End = info.Size()
	} else {
		description.End = start + readLimit
	}

	if start > info.Size() {
		return nil, info.Size(), fmt.Errorf("start byte exceeds file lenght (%d)", info.Size())
	}

	totalLength = description.End - start
	return io.NewSectionReader(file, start, description.End - start), info.Size(), err
}

func GetManager() dao.FileReader {
	return &fileReader{}
}



