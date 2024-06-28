package main

import (
	"LogParser/parser"
	"io/ioutil"
	"log"
	"path/filepath"

	"LogParser/checkstrategies" // Adjust this path to match your actual module path
	"LogParser/utils"           // Adjust this path to match your actual module path
)

func main() {
	// Define paths to log files
	baseDir := "./data"
	browserLogPath := filepath.Join(baseDir, "2024-05-09_17h53m25s_Browser.log")
	clientLogPath := filepath.Join(baseDir, "2024-05-09_17h53m25s_Client.log")
	runtimeLogPath := filepath.Join(baseDir, "2024-05-09_17h53m25s_Runtime.log")
	runtimeLogsJSONPath := filepath.Join(baseDir, "runtime_logs.json")
	clientLogsJSONPath := filepath.Join(baseDir, "client_logs.json")

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
	process, err := logProcessor.Process(browserLogData)
	if err != nil {
		log.Fatalf("Failed to process browser logs: %v", err)
	}
	err = ioutil.WriteFile(filepath.Join(baseDir, "browser_logs.json"), process, 0644)
	if err != nil {
		log.Fatalf("Failed to write browser logs JSON file: %v", err)
	}

	// Process Client Logs
	clientParser := &parser.ClientLogParser{}
	logProcessor.SetParser(clientParser)
	process, err = logProcessor.Process(clientLogData)
	if err != nil {
		log.Fatalf("Failed to process Client logs: %v", err)
	}
	err = ioutil.WriteFile(filepath.Join(baseDir, "client_logs.json"), process, 0644)
	if err != nil {
		log.Fatalf("Failed to write client logs JSON file: %v", err)
	}

	// Process Runtime Logs
	runtimeParser := &parser.RuntimeLogParser{}
	logProcessor.SetParser(runtimeParser)
	process, err = logProcessor.Process(runtimeLogData)
	if err != nil {
		log.Fatalf("Failed to process Runtime logs: %v", err)
	}
	err = ioutil.WriteFile(runtimeLogsJSONPath, process, 0644)
	if err != nil {
		log.Fatalf("Failed to write runtime logs JSON file: %v", err)
	}

	// Define check modules with appropriate log paths
	checkModules := []checkstrategies.CheckModule{
		{"RuntimeLogCheck", &checkstrategies.PythonScriptStrategy{ScriptPath: "Z:\\A_Projekte\\CYBRAIL\\Modules\\module.py"}, runtimeLogsJSONPath},
		{"IntegrityCheck", &checkstrategies.PythonScriptStrategy{ScriptPath: "Z:\\A_Projekte\\CYBRAIL\\Modules\\integrity_check.py"}, runtimeLogsJSONPath},
		{"DisplayCheck", &checkstrategies.PythonScriptStrategy{ScriptPath: "Z:\\A_Projekte\\CYBRAIL\\Modules\\display_check.py"}, runtimeLogsJSONPath},
		{"NetworkConfigCheck", &checkstrategies.PythonScriptStrategy{ScriptPath: "Z:\\A_Projekte\\CYBRAIL\\Modules\\network_config_check.py"}, clientLogsJSONPath},
		{"FrequentReinitialization", &checkstrategies.PythonScriptStrategy{ScriptPath: "Z:\\A_Projekte\\CYBRAIL\\Modules\\frequent_reinitialization.py"}, runtimeLogsJSONPath},
		{"UnusualShutdownCheck", &checkstrategies.PythonScriptStrategy{ScriptPath: "Z:\\A_Projekte\\CYBRAIL\\Modules\\unusual_shutdown_check.py"}, clientLogsJSONPath},
	}

	overallStatus := "success"
	results := make(map[string]interface{})

	// Execute check modules
	for _, module := range checkModules {
		output, err := module.Strategy.Execute(module.LogPath)
		result := utils.ParseModuleOutput(output, err)
		results[module.Name] = result

		// Update overallStatus based on the worst status encountered
		if status, ok := result["status"].(string); ok && utils.CompareStatus(status, overallStatus) > 0 {
			overallStatus = status
		}
	}

	// Print the overall status and individual check statuses
	utils.PrintOverallStatus(overallStatus, results)
}
