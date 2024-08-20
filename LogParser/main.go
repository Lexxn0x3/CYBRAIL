package main

import (
	"LogParser/parser"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"LogParser/checkstrategies" // Adjust this path to match your actual module path
	"LogParser/models"
	"LogParser/utils" // Adjust this path to match your actual module path
)

func main() {
	// Define paths to log files
	baseDir := "./data"
  configPath := filepath.Join(baseDir, "config.json")

  utils.GenerateConfig(configPath,baseDir)

  if !utils.ConfigExists(configPath) {
	  fmt.Println("No configuration created. Exiting.")
		  return
	}

	// Load the configuration
	config, err := models.LoadConfig(configPath)
	if err != nil {
		fmt.Println("Failed to load config:", err)
		return
	}

	// Example: Process each student's log files
	for _, student := range config.Students {
    fmt.Printf("Processing Student: %s (%v)\n", student.Name, student.ID)
		process(baseDir,student.Logs)
    fmt.Println()
	}
}

func process(baseDir string, logPaths models.LogPaths) {
	// Derive the JSON paths from the log paths
	jsonPaths := models.DeriveJSONPaths(baseDir, logPaths)

	// Read and process the logs
	processLogFile(logPaths.BrowserLog, jsonPaths.BrowserLogJSON, &parser.BrowserLogParser{})
	processLogFile(logPaths.ClientLog, jsonPaths.ClientLogJSON, &parser.ClientLogParser{})
	processLogFile(logPaths.RuntimeLog, jsonPaths.RuntimeLogJSON, &parser.RuntimeLogParser{})

	// Execute the check modules
	runCheckModules(jsonPaths)
}

// processLogFile reads, processes, and writes the log data to a JSON file.
func processLogFile(logPath, jsonPath string, curparser parser.LogParser) bool {
	logData, err := ioutil.ReadFile(logPath)
	if err != nil {
		log.Printf("Failed to read log file (%s): %v", logPath, err)
		return false
	}

	logProcessor := &parser.LogProcessor{}
	logProcessor.SetParser(curparser)

	processedData, err := logProcessor.Process(logData)
	if err != nil {
		log.Printf("Failed to process log file (%s): %v", logPath, err)
		return false
	}

	err = ioutil.WriteFile(jsonPath, processedData, 0644)
	if err != nil {
		log.Printf("Failed to write JSON file (%s): %v", jsonPath, err)
		return false
	}

	return true
}

func getExecutableDir() string {
	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	return filepath.Dir(execPath)
}

func runCheckModules(jsonPaths models.JSONPaths) {
  execDir := getExecutableDir()

  modules := []checkstrategies.CheckModule{
		{"RuntimeLogCheck", &checkstrategies.PythonScriptStrategy{ScriptPath: filepath.Join(execDir, "Modules/module.py")}, jsonPaths.RuntimeLogJSON},
		{"IntegrityCheck", &checkstrategies.PythonScriptStrategy{ScriptPath: filepath.Join(execDir, "Modules/integrity_check.py")}, jsonPaths.RuntimeLogJSON},
		{"DisplayCheck", &checkstrategies.PythonScriptStrategy{ScriptPath: filepath.Join(execDir, "Modules/display_check.py")}, jsonPaths.RuntimeLogJSON},
		{"NetworkConfigCheck", &checkstrategies.PythonScriptStrategy{ScriptPath: filepath.Join(execDir, "Modules/network_config_check.py")}, jsonPaths.ClientLogJSON},
		{"FrequentReinitialization", &checkstrategies.PythonScriptStrategy{ScriptPath: filepath.Join(execDir, "Modules/frequent_reinitialization.py")}, jsonPaths.RuntimeLogJSON},
		{"UnusualShutdownCheck", &checkstrategies.PythonScriptStrategy{ScriptPath: filepath.Join(execDir, "Modules/unusual_shutdown_check.py")}, jsonPaths.ClientLogJSON},
	}

	overallStatus := "success"
  overallDetails := []string{}  // Collect details of all errors
  results := make(map[string]interface{})

  for _, module := range modules {
    // Check if the log file path is set and not empty
    if module.LogPath == "" {
        errorMessage := fmt.Sprintf("LogPath is not set or is empty for module: %s", module.Name)
        log.Println(errorMessage)
        results[module.Name] = map[string]interface{}{
            "status": "error",
            "errors": []map[string]string{
                {
                    "message": "LogPath is not set or is empty.",
                    "details": "The path to the log file is missing. Please ensure that the correct log file path is provided.",
                },
            },
        }
        overallStatus = "error"
        overallDetails = append(overallDetails, errorMessage)
        continue
    }

    // Check if the log file exists and was processed
    if !utils.FileExists(module.LogPath) {
        errorMessage := fmt.Sprintf("Required log file not found or not processed: %s", module.LogPath)
        log.Println(errorMessage)
        results[module.Name] = map[string]interface{}{
            "status": "error",
            "errors": []map[string]string{
                {
                    "message": "Required log file not found or not processed.",
                    "details": fmt.Sprintf("The log file at path %s could not be found or was not processed.", module.LogPath),
                },
            },
        }
        overallStatus = "error"
        overallDetails = append(overallDetails, errorMessage)
        continue
    }

    // Execute the module strategy
    output, err := module.Strategy.Execute(module.LogPath)
    result := utils.ParseModuleOutput(output, err)
    results[module.Name] = result

    if status, ok := result["status"].(string); ok && utils.CompareStatus(status, overallStatus) > 0 {
        overallStatus = status
    }
}
  // Print the overall status and results
  utils.PrintOverallStatus(overallStatus, results)
}
