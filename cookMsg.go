package mbox

import (
	"bufio"
	"fmt"
	"mime"
	"strings"
)

// Remove all but vital headers, // TODO multipart, HTML, attachments
func CookMsg(text string) (string, error) {
	var report strings.Builder

	// headers we want to display
	headerKeys := strings.Fields("Date From To Subject" +
		" Content-Type Content-Transfer-Encoding")
	headers := make(map[string]string)
	for _, k := range headerKeys {
		headers[k] = "-"
	}

	// parse line by line
	inHead := true // start in headings
	var ct ContentType
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := scanner.Text()

		if inHead {
			if len(line) > 0 {
				if pos := strings.IndexRune(line, ':'); pos > 0 {
					key := strings.TrimSpace(line[:pos])
					if _, ok := headers[key]; ok {
						headers[key] = decode(strings.TrimSpace(line[pos+1:]))
					}
				}
			} else { // empty line marks end of headers
				for _, k := range headerKeys {
					report.WriteString(fmt.Sprintf("%-7s %s", k, headers[k]))
					report.WriteRune('\n')
					if k == "Subject" {
						report.WriteRune('\n')
					}
				}
				inHead = false
				ct = parseContentType(headers["Content-Type"])
				report.WriteRune('\n')
			}

		} else { // in body
			switch ct.basicType {
			case "text/plain":
				fallthrough // temp
			case "text/html":
				fallthrough // temp
			case "multipart/alternative":
				fallthrough // temp
			case "multipart/mixed":
				fallthrough // temp
			default:
				report.WriteString(line)
				report.WriteRune('\n')
			}
		}
	}

	err := scanner.Err()
	return report.String(), err
}

type ContentType struct {
	basicType, name, fileName, mType, method, format, charset, extra,
	boundary string //, reportType string
}

// get `multipart/mixed` from `multipart/mixed; boundary="abcd"`
func parseContentType(header string) ContentType {
	var ct ContentType
	k := "content-type" // Content-Type: type; key=value; key=value  ...
	inQuotes := false
	header = strings.TrimSpace(header) + "\n"

	var sb strings.Builder
	for _, r := range header {
		if r == '"' {
			inQuotes = !inQuotes
		} else {
			if inQuotes && r != '\n' {
				sb.WriteRune(r)
			} else {
				switch r {
				case '=':
					k = strings.TrimSpace(strings.ToLower(sb.String()))
					sb.Reset()
				case ';', '\n':
					switch k {
					case "":
						// ignore successive ;; or ;\n
					case "content-type":
						ct.basicType = sb.String()
					case "name":
						ct.name = sb.String()
					case "filename":
						ct.fileName = sb.String()
					case "method":
						ct.method = sb.String()
					case "type":
						ct.mType = sb.String()
					case "format":
						ct.format = sb.String()
					case "charset":
						ct.charset = sb.String()
					case "boundary":
						ct.boundary = sb.String()
					// case report-type:
					//	ct.reportType = sb.String()
					default:
						ct.extra = ct.extra + " " + k + "=" + sb.String()
					}
					sb.Reset()
					k = ""
					//fmt.Printf("Rune# %d: is ; or EOL. k is '%s'\n", n, k)
					//fmt.Printf("ct is %+v\n", ct)
				default:
					sb.WriteRune(r)
				}
			} // inQuotes
		} // not "
	} // next rune

	return ct
}

// Decode UTF-8 encoded headers typically starting =?UTF-8? or =?utf-8?
func decode(s string) string {
	if strings.Contains(s, "=?") {
		decoder := new(mime.WordDecoder)
		decoded, err := decoder.DecodeHeader(s)
		if err == nil {
			s = decoded
		}
	}
	return s
}
