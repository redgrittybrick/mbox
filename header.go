package mbox

import (
	"fmt"
	"os"
)

// extracted data - also see type MsgInfo in indexFile.go

type header struct {
	to,
	from,
	date, // date sent by sender's Mail User Agent
	subject,
	received,
	receivedDate, // date received by Mail Transfer Agent at ISP
	fromDate, // date collected by our Mail User Agent (TBird)
	contentType string
}

func (h *header) Reset() {
	h.to = ""
	h.from = ""
	h.date = ""
	h.subject = ""
	h.received = ""
	h.receivedDate = ""
	h.fromDate = ""
	h.contentType = ""
}

// data in form we want to write to indexes

const fileIndexName = "message-files.csv"

type fileInfo struct {
	path string
	size int
}

type fileInfos map[rune]fileInfo

func (fis fileInfos) Write() {
	f, err := os.Create(fileIndexName)
	check(err, "Creating "+fileIndexName)
	defer f.Close()
	for k, fi := range fis {
		fmt.Fprintf(f, "%d|%s|%d\n", k, fi.path, fi.size)
	}
}
