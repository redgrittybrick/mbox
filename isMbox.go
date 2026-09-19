package mbox

import (
	"io"
	"os"
)

// returns true if specified file has content superficially like TB Mbox format
func IsMbox(path string) bool {
	r, err := os.Open(path)
	if err != nil {
		return false
	}

	defer r.Close()

	// get first few bytes of file
	var header [5]byte
	_, err = io.ReadFull(r, header[:])
	if err != nil {
		return false
	}

	// Thunderbird Mbox files seem to begin with "From "
	hs := string(header[:])
	return hs == "From "
}
