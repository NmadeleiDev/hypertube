package filesReader

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"

	"hypertube_storage/dao"
	"hypertube_storage/db"
	"hypertube_storage/parser/env"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
)

const readMaxLimit = 1e6
const readMinLimit = 1e5

var filesDir = env.GetParser().GetFilesDir()

type fileReader struct {
}

func (f *fileReader) HasNullBytes(src []byte) bool {
	if src == nil {
		return false
	}
	return bytes.IndexByte(src, 0) >= 0
}

func (f *fileReader) HasNotNullBytes(src []byte) bool {
	for _, v := range src {
		if v != 0 {
			return true
		}
	}
	return false
}

func (f *fileReader) IsPartWritten(fileName string, part []byte, start int64) bool {
	if part == nil {
		return false
	}

	slices := db.GetLoadedStateDb().GetSliceIndexesForFile(fileName)
	//logrus.Debugf("Got file %v slices: %v", fileName, slices)
	partStart := start
	partEnd := start + int64(len(part))
	nSlices := len(slices)
	for i, _ := range slices {
		if slices[i] == 0 {
			continue
		}
		prevSliceVal := 0
		if i > 0 {
			prevSliceVal = int(slices[i-1])
		}
		if slices[i] > partStart && slices[i] < partEnd {
			left := prevSliceVal - int(start)
			right := slices[i] - start
			if left < 0 {
				left = 0
			}

			if i == 0 && !f.HasNotNullBytes(part[:right]) {
				return false
			} else if i != nSlices - 1 && !f.HasNotNullBytes(part[left:right]) {
				return false
			} else if i == nSlices - 1 && !f.HasNotNullBytes(part[left:]) {
				return false
			}
		}
	}
	return true
}

func (f *fileReader) WaitForFilePart(ctx context.Context, fileName string, start int64, expectedLen int64) ([]byte, int64, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logrus.Errorf("Error init watcher: %v", err)
		return nil, 0, err
	}
	defer watcher.Close()

	type resultStruct struct {
		Data	[]byte
		TotalLen	int64
		Err		error
	}

	done := make(chan resultStruct)
	defer close(done)

	watchCtx, watchCancel := context.WithCancel(ctx)

	defer watchCancel()

	go func(ctx context.Context, resultsChan chan <- resultStruct) {
		for {
			select {
			case <- ctx.Done():
				logrus.Debugf("Exiting file watch due to context DONE")
				return
			case event, ok := <-watcher.Events:
				if !ok {
					resultsChan <- resultStruct{Data: nil, Err: fmt.Errorf("ok is false during watching event")}
				}

				buf, totalLen, err := f.GetFileInRange(fileName, start, expectedLen)
				if err != nil {
					logrus.Debugf("File load err in watch, keep watching: %v", err)
					//resultsChan <- resultStruct{Data: nil, Err: fmt.Errorf("read file error: %v", err)}
				} else if buf != nil && f.IsPartWritten(fileName, buf, start) {
					logrus.Debugf("Sending read res with len = %v", totalLen)
					resultsChan <- resultStruct{Data: buf, TotalLen: totalLen, Err: nil}
					return
				} else {
					logrus.Debugf("Read buf not written or nil (%v %v)", f.IsPartWritten(fileName, buf, start), buf == nil)
					if buf != nil {
						logrus.Debugf("Buf: {%v; %v}", buf[:100], buf[len(buf) - 100:])
					}
				}

				logrus.Debugf("event: %v", event)
				if event.Op & fsnotify.Write == fsnotify.Write {
					logrus.Debugf("modified file: %v", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					resultsChan <- resultStruct{Data: nil, Err: fmt.Errorf("ok is false during watching errors")}
					return
				}
				logrus.Errorf("error: %v", err)
			}
		}
	}(watchCtx, done)

	err = watcher.Add(path.Join(env.GetParser().GetFilesDir(), fileName))
	if err != nil {
		logrus.Errorf("add file %v to watch err: %v", fileName, err)
	}

	select {
	case result := <-done:
		return result.Data, result.TotalLen, result.Err
	case <-ctx.Done():
		logrus.Debugf("Exiting file watch wait due to context DONE")
		return nil, 0, fmt.Errorf("context finished")
	}
}

func (f *fileReader) GetFileInRange(fileName string, start int64, expectedLen int64) (result []byte, totalLength int64, err error) {
	file, err := os.Open(path.Join(filesDir, fileName))
	if err != nil {
		logrus.Errorf("Error open file: %v", err)
		return nil, 0, err
	}

	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		logrus.Errorf("Error getting file info: %v", err)
		return nil, 0, err
	}

	end := start + readMaxLimit
	minEnd := start + readMinLimit

	if start > info.Size() {
		logrus.Debugf("start byte %v exceeds file lenght (%d)", start, info.Size())
		return nil, info.Size(), nil
	}
	if minEnd > expectedLen {
		minEnd = expectedLen
		end = expectedLen
	}

	if minEnd > info.Size() {
		logrus.Debugf("min_end (%v) > info.Size() (%v), returning nil", minEnd, info.Size())
		return nil, info.Size(), fmt.Errorf("no enoght bytes written yet")
	} else if end > info.Size() {
		end = info.Size()
	}

	totalLength = end - start

	logrus.Debugf("Reading %v bytes from %v", totalLength, fileName)

	buf := make([]byte, totalLength)
	n, err := file.ReadAt(buf, start)
	if err != nil {
		logrus.Errorf("Error reading file to buf: %v", err)
		return nil, info.Size(), err
	}

	return buf[:n], info.Size(), nil
}

func (f *fileReader) RemoveFile(fileName string) bool {
	if err := os.Remove(path.Join(filesDir, fileName)); err != nil {
		logrus.Errorf("Error deleting file %v: %v", fileName, err)
		return false
	}
	return true
}

func GetManager() dao.FileReader {
	return &fileReader{}
}



