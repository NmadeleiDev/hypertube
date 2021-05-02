package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type FileRangeDescription struct {
	Start	int64
	End		int64
}

func (f *FileRangeDescription) ParseHeader(header string) error {
	parts := strings.Split(header, "=")
	if len(parts) != 2 || parts[0] != "bytes" {
		return fmt.Errorf("range header is incorrect: '%v'", header)
	}
	limits := strings.Split(parts[1], "-")
	if len(limits) != 2 {
		logrus.Errorf("Error parsing range header! Len != 2: %v", header)
		return fmt.Errorf("len after split != 2")
	}

	if start, err := strconv.ParseInt(limits[0], 10, 64); err != nil {
		logrus.Errorf("Error parsing start in range header! err: %v; header: %v", err, header)
		f.Start = -1
	} else {
		f.Start = start
	}

	if end, err := strconv.ParseInt(limits[1], 10, 64); err != nil {
		logrus.Debugf("Failed to parse end in range header. Err: %v; Header: %v", err, header)
	} else  {
		f.End = -1
		f.End = end
	}

	return nil
}
