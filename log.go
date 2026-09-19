package mbox

import (
	"log"
)

func (p Parser) log(format string, v ...interface{}) {
	if p.logLevel > 0 {
		log.Printf("mbox "+format, v...)
	}
}
