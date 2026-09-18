package mbox

import (
	//	"bufio"
	"fmt"
	//	"log"
	"os"
	// "strconv"
	// "strings"
)

// -------------------------------------------------------------
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

// -------------------------------------------------------------
// data in form we want to write to indexes

const fileIndexName = "message-files.csv"

type fileInfo struct {
	path string
	size int
}

type fileInfos map[rune]fileInfo

/*
func (fis *fileInfos) Read() {
	m := *fis
	f, err := os.Open(fileIndexName)
	if err != nil {
		log.Printf("Can't read file index %q because %v", fileIndexName, err)
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		items := strings.Split(line, "|")
		if len(items) == 3 {
			keyStr := items[0]
			key := rune(keyStr[0])
			path := items[1]
			size, err := strconv.Atoi(items[2])
			check(err, "Atoi("+items[2]+")")
			m[key] = fileInfo{path: path, size: size}
			log.Printf("Read File Info %c : %s (%d)", key, m[key].path, m[key].size)
		}
	}
	err = scanner.Err()
	if err != nil {
		log.Printf("Error reading file index %q : %v", fileIndexName, err)
	}

}
*/

func (fis fileInfos) Write() {
	f, err := os.Create(fileIndexName)
	check(err, "Creating "+fileIndexName)
	defer f.Close()
	for k, fi := range fis {
		fmt.Fprintf(f, "%d|%s|%d\n", k, fi.path, fi.size)
	}
}

// ...............................................................

/*
type messageInfo struct {
	fileRef             rune
	from, date, subject string
	size                int
}

func (mi messageInfo) String() string {
	items := []string{string(mi.fileRef), mi.from, mi.date, mi.subject, strconv.Itoa(mi.size)}
	return strings.Join(items, "|")
}

type messageInfos []messageInfo

func (mis *messageInfos) Read() {

}

func (mis messageInfos) Write() {

}
*/
