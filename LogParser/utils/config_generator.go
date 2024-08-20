package utils

import (
	"LogParser/models"
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// GenerateConfig checks for the existence of a config file and generates one if needed.
func GenerateConfig(configPath string, dataDir string) {
	if promptUserForConfigCreation() {
		logFiles := getAllLogFiles(dataDir)
		groups := groupLogsByBaseName(logFiles)

		students := []models.Student{}
		for baseName, logs := range groups {
			id, name := promptUserForStudentInfo(baseName)
			student := models.Student{
				ID:   id,
				Name: name,
				Logs: logs,
			}
			students = append(students, student)
		}

		cfg := models.Config{Students: students}
		saveConfig(configPath, cfg)
		fmt.Println("Configuration saved.")
	} else {
   return 
	}
}

// configExists checks if the configuration file already exists.
func ConfigExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// promptUserForConfigCreation asks the user if they want to create a new configuration file.
func promptUserForConfigCreation() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to create a configuration? (yes/no): ")
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))
	return answer == "yes" || answer == "y"
}

// getAllLogFiles retrieves all log files from the given directory.
func getAllLogFiles(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalf("Failed to read directory: %v", err)
	}

	var logFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".log") {
			logFiles = append(logFiles, filepath.Join(dir, file.Name()))
		}
	}
	return logFiles
}

// groupLogsByBaseName groups log files by their base name.
func groupLogsByBaseName(logFiles []string) map[string]models.LogPaths {
	groups := make(map[string]models.LogPaths)

	for _, logFile := range logFiles {
		baseName := getBaseName(logFile)
		if strings.Contains(logFile, "_Browser.log") {
			groups[baseName] = models.LogPaths{BrowserLog: logFile}
		} else if strings.Contains(logFile, "_Client.log") {
			logs := groups[baseName]
			logs.ClientLog = logFile
			groups[baseName] = logs
		} else if strings.Contains(logFile, "_Runtime.log") {
			logs := groups[baseName]
			logs.RuntimeLog = logFile
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

// promptUserForStudentInfo asks the user for the student ID and name.
func promptUserForStudentInfo(baseName string) (string, string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter the student ID for logs related to '%s': ", baseName)
	id, _ := reader.ReadString('\n')
	fmt.Printf("Enter the student name for logs related to '%s': ", baseName)
	name, _ := reader.ReadString('\n')
	return strings.TrimSpace(id), strings.TrimSpace(name)
}

// saveConfig saves the configuration to a file.
func saveConfig(path string, cfg models.Config) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("Failed to create config file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(cfg)
	if err != nil {
		log.Fatalf("Failed to write config file: %v", err)
	}
}

