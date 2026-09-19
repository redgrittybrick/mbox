package mbox

import (
	"fmt"
	"io"
	"os"
)

// reads mbox file skipping to specific message
func GetMsg(path string, startByte, endByte int64) (string, error) {

	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.Seek(startByte, io.SeekStart)
	if err != nil {
		return "", err
	}

	l := int(endByte - startByte + 1)
	buf := make([]byte, l)
	n, err := f.Read(buf)

	if n < l {
		return string(buf), fmt.Errorf("Only %d bytes read, expected %d", n, l)
	}

	return string(buf), nil
}
