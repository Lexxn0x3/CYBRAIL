package config

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// configExists checks if the configuration file already exists.
func ConfigExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func CheckLogFilesInConfig(logDir, configPath string) bool {
	// Load the configuration
	config, err := LoadConfig(configPath)
	if err != nil {
		log.Errorf("Failed to load config: %v", err)
		return false
	}

	// Retrieve all log files in the directory
	logFiles := GetAllLogFiles(logDir)

	// Create a map of student log files from the config for quick lookup
	studentLogMap := make(map[string]bool)
	for _, student := range config.Students {
		// Add each log path to the map
		studentLogMap[student.Logs.BrowserLog] = true
		studentLogMap[student.Logs.ClientLog] = true
		studentLogMap[student.Logs.RuntimeLog] = true
		studentLogMap[student.Logs.TypeIntervalLog] = true
	}

	// Check if each log file in the directory has a corresponding student config
	for _, logFile := range logFiles {
		if !studentLogMap[logFile] {
			return false
		}
	}

	return true
}

// getAllLogFiles retrieves all log files from the given directory.
func GetAllLogFiles(dir string) []string {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	var logFiles []string
	for _, file := range files {
		if !file.IsDir() && (strings.HasSuffix(file.Name(), ".log") || strings.HasSuffix(file.Name(), "TypeIntervall.json")) {
			logFiles = append(logFiles, filepath.Join(dir, file.Name()))
		}
	}
	return logFiles
}

// groupLogsByBaseName groups log files by their base name.
func GroupLogsByBaseName(logFiles []string) map[string]LogPaths {
	groups := make(map[string]LogPaths)

	for _, logFile := range logFiles {
		baseName := getBaseName(logFile)
		if strings.Contains(logFile, "_Browser.log") {
			groups[baseName] = LogPaths{BrowserLog: logFile}
		} else if strings.Contains(logFile, "_Client.log") {
			logs := groups[baseName]
			logs.ClientLog = logFile
			groups[baseName] = logs
		} else if strings.Contains(logFile, "_Runtime.log") {
			logs := groups[baseName]
			logs.RuntimeLog = logFile
			groups[baseName] = logs
		} else if strings.Contains(logFile, "_TypeIntervall.json") {
			logs := groups[baseName]
			logs.TypeIntervalLog = logFile
			groups[baseName] = logs
		}
	}

	return groups
}

// getBaseName extracts the base name of the log file (common prefix before _Browser/_Client/_Runtime).
func getBaseName(logFile string) string {
	base := filepath.Base(logFile)
	base = strings.TrimSuffix(base, filepath.Ext(base))
	parts := strings.Split(base, "_")
	return strings.Join(parts[:len(parts)-1], "_")
}
