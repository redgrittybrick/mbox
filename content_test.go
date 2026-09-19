package mbox

import (
	"strings"
	"testing"
)

// Tests the parsing by parseContentType()  of the `Content-Type` header
func TestContent(t *testing.T) {
	for _, example := range strings.Split(examples, "\n") {
		if len(example) > 0 {
			e := strings.Split(example, "|")
			if len(e) < 9 {
				t.Errorf(`testdata has %d elements, 9 needed`, len(e))
			} else {
				ct := parseContentType(e[0])
				compar(t, ct, "basicType", ct.basicType, e[1], e[0])
				compar(t, ct, "name", ct.name, e[2], e[0])
				compar(t, ct, "filename", ct.fileName, e[3], e[0])
				compar(t, ct, "mtype", ct.mType, e[4], e[0])
				compar(t, ct, "method", ct.method, e[5], e[0])
				compar(t, ct, "format", ct.format, e[6], e[0])
				compar(t, ct, "charset", ct.charset, e[7], e[0])
				compar(t, ct, "boundary", ct.boundary, e[8], e[0])
				compar(t, ct, "extra", ct.extra, e[9], e[0])
			}
		}
	}
}

func compar(t *testing.T, ct ContentType, name, have, want, data string) {
	if have != want {
		t.Errorf(`%s is '%s' not '%s' in %#v from %s`, name, have, want, ct, data)
	}
}

// header data|basicType|name|fileName|mType|method|format|charset|boundary|extra
const examples = `
application/octet-stream; name=1016606.csv|application/octet-stream|1016606.csv|||||||
application/pdf;|application/pdf||||||||
application/pdf; name="ACE Form 6 - Customer Keyholder Form -|application/pdf|ACE Form 6 - Customer Keyholder Form -||||||||
audio/mp3; name=+441138231761-0523-104046-359.mp3|audio/mp3|+441138231761-0523-104046-359.mp3||||||||
image/gif;|image/gif||||||||
image/png; filename="logo.png"|image/png||logo.png|||||||
image/png; name="7d6f6DDp1MzO6Rc6.png"|image/png|7d6f6DDp1MzO6Rc6.png|||||||
message/delivery-status|message/delivery-status||||||||
message/global|message/global||||||||
multipart/alternative;|multipart/alternative||||||||
multipart/mixed; boundary="----------=_1712752802-570927-17"|multipart/mixed|||||||----------=_1712752802-570927-17|
multipart/related;|multipart/related|
multipart/related; type="multipart/alternative";|multipart/related|||multipart/alternative|||||
multipart/related; type="text/html";|multipart/related|||text/html|||||
`

/*
multipart/related;
multipart/related; type="multipart/alternative";
multipart/related; type="text/html";
multipart/report; report-type=delivery-status;
text/calendar; charset="utf-8"; method=REQUEST
text/html
text/html;
text/html; charset="iso-8859-1"
text/html; charset=iso-8859-1
text/html; charset="us-ascii"
text/html; charset=us-ascii
text/html; charset = "utf-8"
text/html; charset="utf-8"
text/html; charset=utf-8
text/html; charset="UTF-8"
text/html; charset=UTF-8
text/html;charset=UTF-8
text/plain
text/plain;
text/plain; charset=ISO-2022-JP
text/plain; charset="iso-8859-1"
text/plain; charset=iso-8859-1
text/plain; charset=US-ASCII;
text/plain; charset=US-ASCII; format=flowed
text/plain; charset="utf-8"; format="fixed"
text/plain; name="SpamAssassinReport.txt"
text/rfc822-headers`
*/
