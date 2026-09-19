package mbox

import (
	"log"
)

// Writes a log message if log level > 0 (see parser.LogLevel(n))
func (p Parser) log(format string, v ...interface{}) {
	if p.logLevel > 0 {
		log.Printf("mbox "+format, v...)
	}
}
