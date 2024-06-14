package parser

import "fmt"

// LogParser is the common interface for all log parsers
type LogParser interface {
	Parse(logData []byte) ([]byte, error)
}

// LogProcessor uses the LogParser interface to process logs
type LogProcessor struct {
	parser LogParser
}

func (lp *LogProcessor) SetParser(parser LogParser) {
	lp.parser = parser
}

func (lp *LogProcessor) Process(logData []byte) ([]byte, error) {
	if lp.parser == nil {
		return nil, fmt.Errorf("no parser set")
	}
	return lp.parser.Parse(logData)
}
