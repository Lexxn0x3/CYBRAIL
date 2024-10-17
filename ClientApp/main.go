package main

import (
	"ClientApp/keylogger"
	"ClientApp/openprograms"
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"syscall"
	"time"
)

func main() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	exitChannel := make(chan bool, 1)

	go func() {
		sig := <-signalChannel
		fmt.Printf("Received signal: %v\n", sig)
		openprograms.RunProcessWatch(true)
		exitChannel <- true
	}()

	go openprograms.RunProcessWatch(false)

	keylogger.RunKeywatch(exitChannel)

	err := uploadLogs()
	if err != nil {
		log.Printf("Failed to upload logs: %v", err)
	}
}

func uploadLogs() error {
	// Server URL
	url := "https://cybrail-api.mattzi.de/upload-logs"

	// Create a buffer to hold the multipart form data
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Rename "typing_intervals.json" to "typeIntervall.json"
	err := os.Rename("typing_intervals.json", "typeIntervall.json")
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Failed to rename typing_intervals.json: %v", err)
	}

	// List of log files to upload
	jsonFiles := []string{
		"process_log.json",
		"process_log_exit.json",
		"TypeIntervall.json",
	}

	// Add the SafeExamBrowser logs
	sebLogDir := filepath.Join(os.Getenv("LOCALAPPDATA"), "SafeExamBrowser", "Logs")
	sebLogs, _ := filepath.Glob(filepath.Join(sebLogDir, "*"))

	// Map to group files by timestamp
	timestampToFiles := make(map[string][]string)

	// Regex to extract timestamp
	re := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}_\d{2}h\d{2}m\d{2}s)_.*\.log$`)

	for _, filePath := range sebLogs {
		filename := filepath.Base(filePath)
		matches := re.FindStringSubmatch(filename)
		if matches != nil {
			timestamp := matches[1]
			timestampToFiles[timestamp] = append(timestampToFiles[timestamp], filePath)
		} else {
			// Skip files that do not match the expected pattern
			continue
		}
	}

	// Find the latest timestamp
	var latestTimestamp time.Time
	var latestTimestampStr string
	for timestampStr := range timestampToFiles {
		// Parse timestamp string into time.Time
		t, err := time.Parse("2006-01-02_15h04m05s", timestampStr)
		if err != nil {
			log.Printf("Failed to parse timestamp (%s): %v", timestampStr, err)
			continue
		}
		if t.After(latestTimestamp) {
			latestTimestamp = t
			latestTimestampStr = timestampStr
		}
	}

	// Add files with the latest timestamp to logFiles
	logFiles := []string{}
	if files, ok := timestampToFiles[latestTimestampStr]; ok {
		logFiles = append(logFiles, files...)
	} else {
		log.Printf("No log files found for the latest timestamp")
	}

	// Prefix JSON files with the latest timestamp and add to logFiles
	for _, jsonFile := range jsonFiles {
		// Check if file exists
		if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
			log.Printf("File does not exist: %s", jsonFile)
			continue
		}

		// New filename with timestamp prefix
		newFilename := latestTimestampStr + "_" + jsonFile

		// Copy and rename the file to a temporary location
		tempFilePath := filepath.Join(os.TempDir(), newFilename)
		sourceFile, err := os.Open(jsonFile)
		if err != nil {
			log.Printf("Failed to open file (%s): %v", jsonFile, err)
			continue
		}
		defer sourceFile.Close()

		destFile, err := os.Create(tempFilePath)
		if err != nil {
			log.Printf("Failed to create temp file (%s): %v", tempFilePath, err)
			continue
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, sourceFile)
		if err != nil {
			log.Printf("Failed to copy file data to temp file: %v", err)
			continue
		}

		logFiles = append(logFiles, tempFilePath)
	}

	// Add files to the multipart form
	for _, filePath := range logFiles {
		file, err := os.Open(filePath)
		if err != nil {
			log.Printf("Failed to open file (%s): %v", filePath, err)
			continue // Skip files that can't be opened
		}
		defer file.Close()

		part, err := writer.CreateFormFile("files", filepath.Base(filePath))
		if err != nil {
			log.Printf("Failed to create form file: %v", err)
			continue
		}

		_, err = io.Copy(part, file)
		if err != nil {
			log.Printf("Failed to copy file data: %v", err)
			continue
		}
	}

	// Close the writer to finalize the multipart form
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close multipart writer: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set the content type to multipart/form-data with the boundary
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Set basic auth credentials
	req.SetBasicAuth("uploads", "c4K8vaD!6&iYN9G") // Replace with your credentials

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server returned status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
