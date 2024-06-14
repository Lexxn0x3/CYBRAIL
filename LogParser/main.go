package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"LogParser/parser"
)

func main() {
	// Define paths to log files
	baseDir := "./data"
	browserLogPath := filepath.Join(baseDir, "2024-05-09_17h53m25s_Browser.log")
	clientLogPath := filepath.Join(baseDir, "2024-05-09_17h53m25s_Client.log")
	runtimeLogPath := filepath.Join(baseDir, "2024-05-09_17h53m25s_Runtime.log")

	// Read Browser Log Data
	browserLogData, err := ioutil.ReadFile(browserLogPath)
	if err != nil {
		log.Fatalf("Failed to read browser log file: %v", err)
	}

	// Read Client Log Data
	clientLogData, err := ioutil.ReadFile(clientLogPath)
	if err != nil {
		log.Fatalf("Failed to read client log file: %v", err)
	}

	// Read Runtime Log Data
	runtimeLogData, err := ioutil.ReadFile(runtimeLogPath)
	if err != nil {
		log.Fatalf("Failed to read runtime log file: %v", err)
	}

	// Process Browser Logs
	browserParser := &parser.BrowserLogParser{}
	logProcessor := &parser.LogProcessor{}
	logProcessor.SetParser(browserParser)
	fmt.Println("Processing Browser Logs:")
	process, err := logProcessor.Process(browserLogData)
	if err != nil {
		log.Fatalf("Failed to process browser logs: %v", err)
	}
	err = ioutil.WriteFile(filepath.Join(baseDir, "browser_logs.json"), process, 0644)
	if err != nil {
		log.Fatalf("Failed to write browser logs JSON file: %v", err)
	}
	fmt.Println("Finsihed brower logs")

	// Process Client Logs
	clientParser := &parser.ClientLogParser{}
	logProcessor.SetParser(clientParser)
	fmt.Println("Processing Client Logs:")
	process, err = logProcessor.Process(clientLogData)
	if err != nil {
		log.Fatalf("Failed to process Client logs: %v", err)
	}
	err = ioutil.WriteFile(filepath.Join(baseDir, "client_logs.json"), process, 0644)
	if err != nil {
		log.Fatalf("Failed to write client logs JSON file: %v", err)
	}
	fmt.Println("Finsihed client logs")

	// Process Runtime Logs
	runtimeParser := &parser.RuntimeLogParser{}
	logProcessor.SetParser(runtimeParser)
	fmt.Println("Processing Runtime Logs:")
	process, err = logProcessor.Process(runtimeLogData)
	if err != nil {
		log.Fatalf("Failed to process Runtime logs: %v", err)
	}
	err = ioutil.WriteFile(filepath.Join(baseDir, "Runtime_logs.json"), process, 0644)
	if err != nil {
		log.Fatalf("Failed to write client logs JSON file: %v", err)
	}
	fmt.Println("Finsihed Runtime logs")
}
