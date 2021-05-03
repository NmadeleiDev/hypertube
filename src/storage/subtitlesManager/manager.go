package subtitlesManager

import (
	"io"

	subtitles "github.com/asticode/go-astisub"
	"github.com/sirupsen/logrus"
)

type SubtitlesManager struct {
}

func GetManager() *SubtitlesManager {
	return &SubtitlesManager{}
}

func (m *SubtitlesManager) ConvertSrtToVtt(src io.Reader, dest io.Writer) error {
	srt, err := subtitles.ReadFromSRT(src)
	if err != nil {
		logrus.Errorf("Failed to read srt file: %v", err)
		return err
	}
	if err := srt.WriteToWebVTT(dest); err != nil {
		logrus.Errorf("Error converting to vtt: %v", err)
		return err
	}
	return nil
}
