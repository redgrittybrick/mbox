package mbox

import (
	"fmt"
	"io"
	"log"
	"os"
)

// reads mbox file skipping to specific message
func GetMsg(path string, startByte, endByte int64) (string, error) {
	log.Printf("Seeking %d %d in '%s'", startByte, endByte, path)

	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	o, err := f.Seek(startByte, io.SeekStart)
	if err != nil {
		return "", err
	}
	if o != endByte {
		/*
			return "",
				fmt.Errorf("File '%s' Seek %d reached %d instead", path, endByte, o)
		*/
	}

	l := int(endByte - startByte + 1)
	buf := make([]byte, l)
	n, err := f.Read(buf)

	if n < l {
		return string(buf), fmt.Errorf("Only %d bytes read, expected %d", n, l)
	}

	return string(buf), nil
}
