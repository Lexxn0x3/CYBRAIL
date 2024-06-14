package parser

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type RuntimeLog struct {
	Application struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		Build   string `json:"build"`
	} `json:"application"`
	System struct {
		StartTime    string `json:"start_time"`
		OS           string `json:"os"`
		Computer     string `json:"computer"`
		Model        string `json:"model"`
		Manufacturer string `json:"manufacturer"`
		RuntimeID    string `json:"runtime_id"`
	} `json:"system"`
	Logs []struct {
		Timestamp string `json:"timestamp"`
		Thread    string `json:"thread"`
		Level     string `json:"level"`
		Message   string `json:"message"`
	} `json:"logs"`
}

type RuntimeLogParser struct{}

func (p *RuntimeLogParser) Parse(logData []byte) ([]byte, error) {
	var runtimeLog RuntimeLog

	lines := strings.Split(string(logData), "\n")
	headerParsed := false
	var logs []struct {
		Timestamp string `json:"timestamp"`
		Thread    string `json:"thread"`
		Level     string `json:"level"`
		Message   string `json:"message"`
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		line = removeBOMAndNonASCII(line)
		if line == "" {
			continue
		}

		if !headerParsed {
			if strings.HasPrefix(line, "/*") || strings.HasPrefix(line, "#") {
				parseHeaderInfo(&runtimeLog, lines)
				headerParsed = true
			}
		}

		if headerParsed && !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "/*") {
			log, err := parseLogEntry(line)
			if err != nil {
				return nil, err
			}
			logs = append(logs, log)
		}
	}

	runtimeLog.Logs = logs

	logJSON, err := json.MarshalIndent(runtimeLog, "", "  ")
	if err != nil {
		return nil, err
	}

	// Output the JSON for demonstration purposes
	//fmt.Printf("%s\n", logJSON)

	return logJSON, err
}

func parseHeaderInfo(runtimeLog *RuntimeLog, lines []string) {
	appInfoRegex := regexp.MustCompile(`^\/\* Safe Exam Browser, Version ([^\(]+) \(([^)]+)\), Build ([^\n]+)$`)
	systemInfoRegex := regexp.MustCompile(`^# Application started at (.*?)$`)
	osInfoRegex := regexp.MustCompile(`^# Running on (.*?), (.*?)$`)
	computerInfoRegex := regexp.MustCompile(`^# Computer '(.*?)' is a (.*?) manufactured by (.*?)$`)
	runtimeIDRegex := regexp.MustCompile(`^# Runtime-ID: (.*?)$`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		line = removeBOMAndNonASCII(line)
		if line == "" || strings.HasPrefix(line, "/*") || strings.HasPrefix(line, "#") {
			if matches := appInfoRegex.FindStringSubmatch(line); matches != nil {
				runtimeLog.Application.Name = "Safe Exam Browser"
				runtimeLog.Application.Version = strings.TrimSpace(matches[1])
				runtimeLog.Application.Build = strings.TrimSpace(matches[3])
			}
			if matches := systemInfoRegex.FindStringSubmatch(line); matches != nil {
				runtimeLog.System.StartTime = strings.TrimSpace(matches[1])
			}
			if matches := osInfoRegex.FindStringSubmatch(line); matches != nil {
				runtimeLog.System.OS = strings.TrimSpace(matches[1])
				runtimeLog.System.Computer = strings.TrimSpace(matches[2])
			}
			if matches := computerInfoRegex.FindStringSubmatch(line); matches != nil {
				runtimeLog.System.Computer = strings.TrimSpace(matches[1])
				runtimeLog.System.Model = strings.TrimSpace(matches[2])
				runtimeLog.System.Manufacturer = strings.TrimSpace(matches[3])
			}
			if matches := runtimeIDRegex.FindStringSubmatch(line); matches != nil {
				runtimeLog.System.RuntimeID = strings.TrimSpace(matches[1])
			}
		}
	}
}

func parseLogEntry(line string) (struct {
	Timestamp string `json:"timestamp"`
	Thread    string `json:"thread"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}, error) {
	// Extract log entry info
	re := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}) \[(\d+)\] - (INFO|DEBUG|ERROR|WARNING): (.*)`)
	matches := re.FindStringSubmatch(line)

	if matches == nil {
		return struct {
			Timestamp string `json:"timestamp"`
			Thread    string `json:"thread"`
			Level     string `json:"level"`
			Message   string `json:"message"`
		}{}, fmt.Errorf("log line does not match expected format: %s", line)
	}

	timestamp := matches[1]
	thread := matches[2]
	level := matches[3]
	message := matches[4]

	return struct {
		Timestamp string `json:"timestamp"`
		Thread    string `json:"thread"`
		Level     string `json:"level"`
		Message   string `json:"message"`
	}{
		Timestamp: timestamp,
		Thread:    thread,
		Level:     level,
		Message:   message,
	}, nil
}
