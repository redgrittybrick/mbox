package mbox

import (
	"io"
	"strings"
	"github.com/emersion/go-message"
	_ "github.com/emersion/go-message/charset"
)

func Outline(msg string) (string, string) {
	var sbo, sbt strings.Builder

	// remove the "From - date" line as that is nonstandard addition by T'bird
	p := strings.IndexRune(msg, '\x0a') // LF

	r := strings.NewReader(msg[p+1:])
	m, err := message.Read(r)
	if message.IsUnknownCharset(err) { // This error is not fatal
		sbo.WriteString("Unknown encoding: " + err.Error() + "\n")
	} else if err != nil {
		sbo.WriteString("Error: " + err.Error() + "\n")
		return sbo.String(), sbt.String()
	}

	if mr := m.MultipartReader(); mr != nil { // This is a multipart message
		sbo.WriteString("== This is a multipart message containing:\n")
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				sbo.WriteString("Error: " + err.Error() + "\n")
			}

			t, _, _ := p.Header.ContentType()
			sbo.WriteString("== A part with type " + t + "\n")

			switch t {
			case "text/plain", "text/html":
				/*
				b, _ := io.ReadAll(p.Body)
				sbt.WriteString("================ text part ================\n")
				sbt.WriteString(string(b))
				sbt.WriteString("-------------------------------------------\n")
				*/
			case "multipart/alternative":
				// BUG - doing io.ReadAll(p.Body) here steals from the library parse
				
				b, _ := io.ReadAll(p.Body)
				o,t := Outline(string(b))
				sbo.WriteString(o)
				sbt.WriteString(t)
				
			}

		}
	} else {
		t, _, _ := m.Header.ContentType()
		sbo.WriteString("== This is a non-multipart message with type " + t + "\n")
		switch t {
		case "text/plain", "text/html":
			// BUG
			/*
			b, _ := io.ReadAll(m.Body)
			sbt.WriteString("================ text part ================\n")
			sbt.WriteString(string(b))
			sbt.WriteString("-------------------------------------------\n")
			*/
		}
	}
	return sbo.String(), sbt.String()
}
