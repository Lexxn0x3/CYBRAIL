package parser

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// BrowserLog represents the structure of a log entry.
type BrowserLog struct {
	Timestamp        string `json:"timestamp"`
	LogLevel         string `json:"log_level"`
	Component        string `json:"component"`
	LineNumber       int    `json:"line_number"`
	Message          string `json:"message"`
	Source           string `json:"source,omitempty"`
	SourceLineNumber int    `json:"source_line_number,omitempty"`
}

// BrowserLogParser is a parser for browser logs.
type BrowserLogParser struct{}

// Parse processes the plain text log data and converts it to JSON objects.
func (p *BrowserLogParser) Parse(logData []byte) ([]byte, error) {
	logLines := strings.Split(string(logData), "\n")
	var logs []BrowserLog

	for _, line := range logLines {
		if line == "" {
			continue
		}

		line = removeBOMAndNonASCII(line)

		log, err := parseLogLine(line)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	logJSON, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		return nil, err
	}

	// Process logs (in JSON format)
	//fmt.Printf("%s\n", logJSON)
	return logJSON, err
}

// parseLogLine parses a single log line and returns a BrowserLog object.
func parseLogLine(line string) (BrowserLog, error) {
	// Regex to match the log entry format with optional source and source line number
	re := regexp.MustCompile(`\[(\d{4}/\d{6}\.\d{3}):(INFO|WARNING|ERROR|DEBUG):([^\(]+)\((\d+)\)\] "?(.*?)"?(?:, source: (.*?)(?: \((\d+)\))?)?$`)
	matches := re.FindStringSubmatch(line)

	if matches == nil {
		return BrowserLog{}, fmt.Errorf("log line does not match expected format: %s", line)
	}

	// Parse the matched groups
	timestamp := matches[1]
	logLevel := matches[2]
	component := matches[3]
	lineNumber, _ := strconv.Atoi(matches[4])
	message := matches[5]
	source := ""
	sourceLineNumber := 0
	if len(matches) > 6 {
		source = matches[6]
	}
	if len(matches) > 7 && matches[7] != "" {
		sourceLineNumber, _ = strconv.Atoi(matches[7])
	}

	log := BrowserLog{
		Timestamp:        timestamp,
		LogLevel:         logLevel,
		Component:        component,
		LineNumber:       lineNumber,
		Message:          message,
		Source:           source,
		SourceLineNumber: sourceLineNumber,
	}

	return log, nil
}
