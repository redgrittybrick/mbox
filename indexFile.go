package mbox

import (
	"bufio"
	"fmt"
	"log"
	"mime"
	"os"
	"strings"
)

// ---------------------------------------------------------------
// Keeps track of which part of a message contains line just read
type State int

const (
	inHeader State = iota
	inMulti
	inTail
)

// ---------------------------------------------------------------
// Items needed to make a list of messages
type MsgInfo struct {
	date, to, from, subject string // from message headers
	linePos                 int    // line on which message starts
	startByte, endByte      int64  // byte offset in file
	size                    int    // bytes in message
}

// returns printable form of "Subject:" header data
func (m *MsgInfo) Subject() string {
	subject := m.subject

	// Decode UTF-8 encoded headers typically starting =?UTF-8? or =?utf-8?
	if strings.Contains(subject, "=?") {
		dec := new(mime.WordDecoder)
		header, err := dec.DecodeHeader(subject)
		if err == nil {
			subject = header
		}
	}
	return subject
}

func (m *MsgInfo) Date() string {
	return m.date
}

func (m *MsgInfo) To() string {
	return m.to
}

// returns readable for of "From:" header
func (m *MsgInfo) From() string {
	from := m.from

	// remove superfluous quotes from first part of "From:" header value
	if len(from) > 2 && from[0] == '"' {
		p := strings.IndexRune(from[1:], '"')
		if p > -1 {
			from = from[1:p+1] + from[p+2:]
		}
	}

	// and do any UTF-8 etc decoding indicated
	if len(from) > 2 && from[:2] == "=?" {
		dec := new(mime.WordDecoder)
		header, err := dec.DecodeHeader(from)
		if err == nil {
			from = header
		}
	}
	return from
}

func (m *MsgInfo) Start() int64 {
	return m.startByte
}

func (m *MsgInfo) End() int64 {
	return m.endByte
}

func (m *MsgInfo) Size() int {
	return m.size
}

type MsgList []MsgInfo

func (l *MsgList) String() string {
	var sb strings.Builder
	sb.WriteString(`
Date             | From                 | Lines   | Subject
-----------------|----------------------|---------|---------
`)
	for _, m := range *l {
		sb.WriteString(fmt.Sprintf("%16.16s | %-20.20s |%8d | %-30.30s\n",
			m.date, m.From(), m.size, m.Subject()))
	}
	return sb.String()
}

// ---------------------------------------------------------------
// Returns a list of messages in a specified mbox file
func IndexFile(fileName string) (MsgList, error) {
	var list MsgList

	file, err := os.Open(fileName)
	if err != nil {
		return list, err
	}
	defer file.Close()

	state := inHeader
	h := header{}
	h.Reset()
	lineNo := 0
	var startByte, endByte int64 = 0, 0
	key := ""
	boundary := "" // only for MIME/Multipart messages

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lineNo++

		// first line may be a "From - " line with a date we might need
		if lineNo == 1 && len(line) > 7 && line[:7] == "From - " {
			h.fromDate = fromDate(line)
		}

		switch state {
		case inHeader:
			if len(line) > 0 && line[0] != ' ' && line[0] != '\t' { // not a continuation
				if i := strings.IndexRune(line, ':'); i > -1 { // poss key:value
					key = line[:i]
					//log.Printf("--key--%s--", key)
					switch key {
					case "From", "FROM":
						h.from = strings.TrimSpace(line[5:])
					case "To":
						h.to = strings.TrimSpace(line[3:])
					case "Subject":
						h.subject = strings.TrimSpace(line[8:])
					case "Date":
						h.date = strings.TrimSpace(line[5:])
					case "Received":
						h.received = strings.TrimSpace(line[10:])
						d := getRcvdDate(h.received)
						if len(d) > 0 {
							h.receivedDate = d
						}
					case "Content-Type":
						h.contentType = strings.TrimSpace(line[13:])
					}
					log.Printf("l:%-6d k:%-20s line:%s\n", lineNo, key, line)
				}
			}

			if len(line) > 0 && (line[0] == ' ' || line[0] == '\t') { // Continuation?
				//extra := strings.TrimSpace(line)
				extra := line
				switch key {
				case "From", "FROM":
					h.from = h.from + extra
				case "To":
					h.to = h.to + extra
				case "Subject":
					h.subject = h.subject + extra
					log.Printf(">>>>>>>>Subject>>%s<<", h.subject)
				case "Date":
					h.date = h.date + extra
				case "Received":
					h.received = h.received + extra
					d := getRcvdDate(h.received)
					if len(d) > 0 {
						h.receivedDate = d
					}
				case "Content-Type":
					h.contentType = h.contentType + extra
				}
				log.Printf("l:%-6d k+%-20s line:%s\n", lineNo, key, line)
			}

			// end of headers is marked by an empty line
			if len(line) == 0 {
				if len(h.contentType) > 10 && h.contentType[:10] == "multipart/" {
					state = inMulti
					p := strings.LastIndex(h.contentType, "boundary=")
					if p > 0 {
						b := h.contentType[p+9:]
						if len(b) > 1 && b[0] == '"' && b[len(b)-1] == '"' {
							b = b[1 : len(b)-1] // remove quotes from `boundary="something"`
						}
						boundary = "--" + b + "--"
						log.Printf("multipart boundary = '%s'\n", boundary)
					}
				} else {
					state = inTail
				}
			}

		case inMulti: // no longer in header part of a message
			// look for marker of start of next message
			// Use multipart marker + "--" when available
			// because at least one email had line == "From " in text/plain body part!
			if line == boundary {
				log.Printf("=== End of multipart message === '%s'\n", line)
				state = inTail // now look for next "From " line
			}

		case inTail: // look for "From " indicating start of next messge
			if line == "From " || (len(line) > 7 && line[:7] == "From - ") {
				list = append(list, MsgInfo{
					date:      parseDate(h),
					to:        h.to,
					from:      h.from,
					subject:   h.subject,
					linePos:   lineNo, // TODO: this is end of msg not start.
					startByte: startByte,
					endByte:   endByte - 1, // Not sure why I am out by 1 - CRLF issue?
					size:      int(endByte - startByte),
				})
				//printLine(h, bcount)
				h.Reset()
				boundary = ""
				state = inHeader
				startByte = endByte // + 1 // See above endByte issue
				log.Printf("=== Start of next message === '%s'\n", line)
				if len(line) > 7 && line[:7] == "From - " {
					h.fromDate = fromDate(line) // pseudo header, last resort for a msg date
				}
			}
		} // state switch

		endByte += int64(len(line) + 2) // assume CR LF line separator
	} // next line

	// include final message from file (if any)
	if len(h.to) > 0 { // TODO avoid code repeat
		list = append(list, MsgInfo{
			date:      parseDate(h),
			to:        h.to,
			from:      h.from,
			subject:   h.subject,
			linePos:   lineNo, // TODO: this is end of msg not start.
			startByte: startByte,
			endByte:   endByte - 1, // Not sure why I am out by 1 - CRLF issue?
			size:      int(endByte - startByte),
		})
	}

	if err := scanner.Err(); err != nil {
		return list, err
	}
	log.Printf("%d lines processed in %s\n", lineNo, fileName)
	log.Printf("%d messages found\n", len(list))
	return list, nil
}

// extract any date info from "Received" header
func getRcvdDate(rcvd string) string {
	/*
		RFC2822 defines the Received header as
			received        =       "Received:" name-val-list ";" date-time CRLF
		Examples from real mail, without any date:
			Received: from mail.interdns.co.uk [83.170.124.82]
			Received: from unknown ([10.237.48.10])
			Received: from unknown (HELO smtp1.interdns.co.uk) (192.168.2.150)
			Received: from unknown (HELO ?10.0.0.185?) (Alice.Wilson@ian1@81.149.177.123)
		Examples from real mail, with a date:
			Received: (qmail 29245 invoked from network); 9 Mar 2017 15:40:06 -0000
			Received: from mail.interdns.co.uk [83.170.124.82]
			 for <Briars@localhost> (single-drop); Mon, 13 Mar 2017 21:55:11 +0000 (GMT)
			Received: from DHV ([192.168.55.8]) by interviewermail.djsresearch.com with Microsoft SMTPSVC(8.5.9600.16384);
			 Wed, 15 Mar 2017 06:00:57 +0000
			Received: from mail.cixhosting.co.uk [83.170.124.82]
			 for <ian@localhost> (single-drop); Wed, 01 Jan 2020 13:10:33 +0000 (WET)
			Received: from [10.0.0.253] (unknown [10.0.0.253])
			 Wed, 29 May 2024 08:32:16 +0100 (+01)
		Notes:
			Continuations ("folding" in RFC) have an intial tab (HT) character.
			Braces (...) are defined in the RFC as comments.
			The Panasonic BLC-101 camera fails to put a semicolon ";" before the date.
	*/
	if len(rcvd) > 1 {
		//p := strings.IndexRune(rcvd, ';')
		p := strings.LastIndex(rcvd, ";")
		if p > -1 && p < len(rcvd) { // occurs and not at EOL
			return strings.TrimSpace(rcvd[p+1:])
		}
	}
	return ""
}

// converts date found on "From - " line to RFC2822 format
func fromDate(s string) string {
	// expect "From - Wed Feb 03 17:49:17 2016"
	// want   "Wed, 03 Feb 2016 17:49:17 +0000"
	if matches(s, "Ulll - Ull Ull DD DD:DD:DD DDDD") {
		f := strings.Fields(s)
		return fmt.Sprintf("%s, %s %s %s %s +0000", f[2], f[4], f[3], f[6], f[5])
	}
	return s
}
