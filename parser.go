package mbox

import (
	"io"
	"log"
)

type Parser struct {
	logLevel int
}

func New() *Parser {
	p := Parser{}
	return &p
}

func (p *Parser) LogLevel(level int) *Parser {
	p.logLevel = level
	return p
}

// Wraps IndexFile function in a method of Parser so we can manage logging
func (p Parser) IndexFile(fileName string) (MsgList, error) {

	// Turn off logging for indexing process if not requested
	logTarget := log.Writer()
	if p.logLevel == 0 {
		log.SetOutput(io.Discard)
	}

	msgs, err := IndexFile(fileName)

	// Reinstate whatever logging caller had previously
	if p.logLevel == 0 {
		log.SetOutput(logTarget)
	}

	return msgs, err
}
