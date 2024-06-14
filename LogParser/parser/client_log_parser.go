package parser

import (
	"encoding/json"
	"regexp"
	"strings"
)

// ClientLog represents the structure of a client log.
type ClientLog struct {
	Sessions []struct {
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
		Logs      []struct {
			Timestamp string `json:"timestamp"`
			ThreadID  string `json:"thread_id"`
			Level     string `json:"level"`
			Message   string `json:"message"`
			Details   string `json:"details,omitempty"`
		} `json:"logs"`
	} `json:"sessions"`
}

// ClientLogParser is a parser for client logs.
type ClientLogParser struct{}

// Parse processes the plain text log data and converts it to JSON objects.
func (p *ClientLogParser) Parse(logData []byte) ([]byte, error) {
	logLines := strings.Split(string(logData), "\n")
	var sessions []struct {
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
		Logs      []struct {
			Timestamp string `json:"timestamp"`
			ThreadID  string `json:"thread_id"`
			Level     string `json:"level"`
			Message   string `json:"message"`
			Details   string `json:"details,omitempty"`
		} `json:"logs"`
	}
	var currentSession *struct {
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
		Logs      []struct {
			Timestamp string `json:"timestamp"`
			ThreadID  string `json:"thread_id"`
			Level     string `json:"level"`
			Message   string `json:"message"`
			Details   string `json:"details,omitempty"`
		} `json:"logs"`
	}

	sessionStartRe := regexp.MustCompile(`(?m)^# New client instance started at (.+)$`)
	sessionEndRe := regexp.MustCompile(`^# Client instance terminated at (.+)$`)
	logEntryRe := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) \[(\d+)\] - (\w+): (.*)$`)
	detailRe := regexp.MustCompile(`^   (.+)$`)

	for _, line := range logLines {
		line = strings.TrimSpace(line)
		line = removeBOMAndNonASCII(strings.TrimSpace(line))

		if matches := sessionStartRe.FindStringSubmatch(line); matches != nil {
			startTime := matches[1]
			currentSession = &struct {
				StartTime string `json:"start_time"`
				EndTime   string `json:"end_time"`
				Logs      []struct {
					Timestamp string `json:"timestamp"`
					ThreadID  string `json:"thread_id"`
					Level     string `json:"level"`
					Message   string `json:"message"`
					Details   string `json:"details,omitempty"`
				} `json:"logs"`
			}{
				StartTime: startTime,
			}
		} else if matches := sessionEndRe.FindStringSubmatch(line); matches != nil {
			if currentSession != nil {
				endTime := matches[1]
				currentSession.EndTime = endTime
				sessions = append(sessions, *currentSession)
				currentSession = nil
			}
		} else if currentSession != nil {
			if logEntryMatches := logEntryRe.FindStringSubmatch(line); logEntryMatches != nil {
				timestamp := logEntryMatches[1]
				threadID := logEntryMatches[2]
				level := logEntryMatches[3]
				message := logEntryMatches[4]
				details := ""

				for i := len(logLines) - 1; i >= 0; i-- {
					if detailMatches := detailRe.FindStringSubmatch(logLines[i]); detailMatches != nil {
						details = detailMatches[1] + "\n" + details
						logLines = logLines[:i]
					} else {
						break
					}
				}

				logEntry := struct {
					Timestamp string `json:"timestamp"`
					ThreadID  string `json:"thread_id"`
					Level     string `json:"level"`
					Message   string `json:"message"`
					Details   string `json:"details,omitempty"`
				}{
					Timestamp: timestamp,
					ThreadID:  threadID,
					Level:     level,
					Message:   message,
					Details:   strings.TrimSpace(details),
				}

				currentSession.Logs = append(currentSession.Logs, logEntry)
			}
		}
	}

	clientLog := ClientLog{Sessions: sessions}
	logJSON, err := json.MarshalIndent(clientLog, "", "  ")
	if err != nil {
		return nil, err
	}

	// Process logs (in JSON format)
	//fmt.Printf("%s\n", logJSON)
	return logJSON, err
}
