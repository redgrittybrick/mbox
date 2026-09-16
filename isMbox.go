package mbox

import (
	"io"
	"os"
)

// Mbox files seem to begin with "From "
func IsMbox(path string) bool {
	r, err := os.Open(path)
	if err != nil {
		return false
	}

	defer r.Close()

	var header [5]byte
	_, err = io.ReadFull(r, header[:])
	if err != nil {
		return false
	}

	hs := string(header[:])
	//log.Printf("hs=%q %x", hs,hs)
	return hs == "From "
}
