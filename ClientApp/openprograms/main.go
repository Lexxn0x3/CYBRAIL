package openprograms

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/shirou/gopsutil/process"
)

type ProcessInfo struct {
	PID         int32  `json:"pid"`
	Name        string `json:"name"`
	CommandLine string `json:"cmdline"`
}

func RunProcessWatch(exit bool) {
	var logFileName string
	if exit {
		logFileName = "process_log_exit.json"
	} else {
		logFileName = "process_log.json"
	}
	// Create or open the log file for appending
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening/creating log file:", err)
		return
	}
	defer logFile.Close()

	// Get the current timestamp
	timestamp := time.Now().Format(time.RFC3339)

	// Write the timestamp to the log file in JSON format
	_, err = logFile.WriteString(fmt.Sprintf("{\n  \"timestamp\": \"%s\",\n  \"processes\": [\n", timestamp))
	if err != nil {
		fmt.Println("Error writing to log file:", err)
		return
	}

	// Get a list of all processes
	processes, err := process.Processes()
	if err != nil {
		fmt.Println("Error retrieving processes:", err)
		return
	}

	// Loop through the processes and create a slice of ProcessInfo structs
	var processList []ProcessInfo
	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			name = "unknown"
		}
		cmdline, err := p.Cmdline()
		if err != nil {
			cmdline = "unknown"
		}

		processInfo := ProcessInfo{
			PID:         p.Pid,
			Name:        name,
			CommandLine: cmdline,
		}
		processList = append(processList, processInfo)
	}

	// Marshal the process list to JSON
	processesJSON, err := json.MarshalIndent(processList, "  ", "    ")
	if err != nil {
		fmt.Println("Error marshaling process list to JSON:", err)
		return
	}

	// Write the processes JSON to the log file
	_, err = logFile.Write(processesJSON)
	if err != nil {
		fmt.Println("Error writing process list to log file:", err)
		return
	}

	// Close the JSON object
	_, err = logFile.WriteString("\n  ]\n}\n")
	if err != nil {
		fmt.Println("Error closing JSON log file:", err)
		return
	}

	fmt.Println("Process list logged in JSON format to process_log.json")
}
