/*
The dates found in email headers vary a lot in format and are often
malformed

Some real examples

	Date: 07 Feb 18 12:03:02                    // E.ON 2 digit year and no offset
	Date: 10 Mar 2017 20:37:12 +0100
	Date: 9 Mar 2017 16:38:21 +0000             // day can be single digit
	Date: 04 May 2017 11:36:05 +0200            // or have leading zero
	Date: Tue, 30 Jan 2018 22:54:19 +0000       // can have dayname prefix
	Date: Tue, 30 Jan 2018 07:44:40 +0000 (GMT) // and timezone in a (comment)
	Date: Mon, 23 Apr 2018 12:51:21 -0400       // B.G. offset can be negative

Some found examples for which mail.ParseDate() returns an error

	Date: Fri, 0 Jan 2011 12:20:30 +0100                   // invalid day of month
	Received: ... ; 6 May 2004 16:14:42 UT                 // should be UTC
	Received: ... ; Fri, 26 Mar 2004 15:30:19 +0000 (GMT)  // (GMT) is comment
*/
package mbox

import (
	"log"
	"net/mail"
	"unicode"
)


// Tries to get msg date from "Date:" or "Received:" headers or initial "From "
func parseDate(h header) string {
	cause := "no date found"

	for _, d := range []string{h.date, h.receivedDate, h.fromDate} {
		if len(d) > 0 {
			t, err := mail.ParseDate(fixDate(d))
			if err == nil {
				return t.UTC().Format("2006-01-02 15:04")
			}
			cause = err.Error()
		}
	}

	log.Printf("dates '%s' '%s' '%s' %s\n",
		h.date, h.receivedDate, h.fromDate, cause)
	return "unknown"
}

// fixes some broken date formats
func fixDate(date string) string {
	l := len(date)

	if l > 5 && date[l-5:] == "(GMT)" { // () marks comment
		date = date[:l-5] + "GMT"
	}

	if l > 3 && date[l-3:] == " UT" { // should be UTC
		date = date + "C"
	}

	switch {
	case matches(date, "D Ull DD DD:DD:DD"): // 7 May 17 23:17:54
		return date + " +0000"
	case matches(date, "DD Ull DD DD:DD:DD"): // 17 May 17 23:17:54
		return date + " +0000"
	case matches(date, "D Ull DDDD DD:DD:DD"): // 7 May 2017 23:17:54
		return date + " +0000"
	case matches(date, "DD Ull DDDD DD:DD:DD"): // 17 May 2017 23:17:54
		return date + " +0000"
	default:
		return date
	}
}

// check if string patches a fixed pattern (not an R.E.)
func matches(s, pattern string) bool {
	if len(s) != len(pattern) {
		return false
	}
	runes := []rune(s)
	for i, p := range pattern {
		r := runes[i]
		switch p {
		case 'L':
			if !unicode.IsLetter(r) {
				return false
			}
		case 'D':
			if !unicode.IsDigit(r) {
				return false
			}
		case '^':
			if !unicode.IsControl(r) {
				return false
			}
		case 'l':
			if !unicode.IsLower(r) {
				return false
			}
		case 'U':
			if !unicode.IsUpper(r) {
				return false
			}
		case 'P':
			if !unicode.IsPrint(r) {
				return false
			}
		case 'M':
			if !unicode.IsMark(r) {
				return false
			}
		case 'N':
			if !unicode.IsNumber(r) {
				return false
			}
		case ' ':
			if !unicode.IsSpace(r) {
				return false
			}
		case '.':
			if !unicode.IsPunct(r) {
				return false
			}
		case 'S':
			if !unicode.IsSymbol(r) {
				return false
			}
		case 'G':
			if !unicode.IsGraphic(r) {
				return false
			}
		default:
			if r != p {
				return false
			}
		}
	}
	return true
}
