package main

import (
	"LogParser/parser"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"

	"LogParser/checkstrategies" // Adjust this path to match your actual module path
	"LogParser/config"
	"LogParser/utils" // Adjust this path to match your actual module path
)

func main() {

	// Load the template engine
	engine := html.New("./views", ".html")
	engine.AddFunc("contains", strings.Contains)

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error { return c.Render("index", nil) })
	app.Post("/run-overview", runOverview)
	app.Post("/save-config", saveConfig)
	app.Get("/config", listModuleConfigs)
	app.Get("/edit-module-config", editModuleConfig)
	app.Get("/logfilelist", logfilelist)
	app.Post("/save-module-config", saveModuleConfig)

	log.Fatal(app.Listen(":3000"))
}

func listModuleConfigs(c *fiber.Ctx) error {
	modulesDir := "./Modules"

	files, err := os.ReadDir(modulesDir)
	if err != nil {
		return c.Status(500).SendString("Failed to read modules directory")
	}

	var configFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".config") {
			configFiles = append(configFiles, file.Name())
		}
	}

	return c.Render("listconfigs", fiber.Map{
		"ConfigFiles": configFiles,
	})
}

func editModuleConfig(c *fiber.Ctx) error {
	configFile := c.Query("file")
	modulesDir := "./Modules"
	filePath := filepath.Join(modulesDir, configFile)

	// Load the .config file
	file, err := os.Open(filePath)
	if err != nil {
		return c.Status(500).SendString("Failed to open config file")
	}
	defer file.Close()

	var config map[string]map[string]interface{}
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString("Failed to decode config file")
	}

	return c.Render("editconfig", fiber.Map{
		"Config":   config,
		"FileName": configFile,
	})
}

func saveModuleConfig(c *fiber.Ctx) error {
	modulesDir := "./Modules"
	configFile := c.FormValue("FileName")
	filePath := filepath.Join(modulesDir, configFile)
	fmt.Println(filePath, c.FormValue("FileName"))

	// Load the existing configuration
	file, err := os.Open(filePath)
	if err != nil {
		return c.Status(500).SendString("Failed to open config file")
	}
	defer file.Close()

	var config map[string]map[string]interface{}
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).SendString("Failed to decode config file")
	}

	// Update the config with form values
	for key := range config {
		config[key]["value"] = c.FormValue(key)
		if config[key]["type"] == "list" {
			config[key]["value"] = strings.Split(c.FormValue(key), ",")
		}
	}

	// Save the updated config back to the file
	file, err = os.Create(filePath)
	if err != nil {
		return c.Status(500).SendString("Failed to save config file")
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(&config)
	if err != nil {
		return c.Status(500).SendString("Failed to encode config file")
	}

	return c.Redirect("/edit-module-config?file=" + configFile)
}

func saveConfig(c *fiber.Ctx) error {
	baseName := c.FormValue("baseName")
	id := c.FormValue("id")
	name := c.FormValue("name")

	// Load the existing configuration
	baseDir := "./data"
	configPath := filepath.Join(baseDir, "config.json")
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return c.Status(500).SendString("Failed to load config")
	}

	// Find or add the entry
	found := false
	for i, student := range cfg.Students {
		if student.Logs.BrowserLog == baseName {
			cfg.Students[i].ID = id
			cfg.Students[i].Name = name
			found = true
			break
		}
	}

	newBaseName := filepath.Join(baseDir, baseName)

	if !found {
		cfg.Students = append(cfg.Students, config.Student{
			ID:   id,
			Name: name,
			Logs: config.LogPaths{BrowserLog: newBaseName + "_Browser.log", ClientLog: newBaseName + "_Client.log", RuntimeLog: newBaseName + "_Runtime.log"},
		})
	}

	// Save the updated configuration
	err = config.SaveConfig(configPath, cfg)
	if err != nil {
		return c.Status(500).SendString("Failed to save config")
	}

	return c.Redirect("/")
}

func runOverview(c *fiber.Ctx) error {
	baseDir := "./data"
	configPath := filepath.Join(baseDir, "config.json")
	existingConfig, _ := config.LoadConfig(configPath)

	if config.CheckLogFilesInConfig(baseDir, configPath) {
		results := make(chan map[string]interface{}, len(existingConfig.Students))
		for _, student := range existingConfig.Students {
			go func(student config.Student) {
				fmt.Printf("Processing Student: %s (%v)\n", student.Name, student.ID)
				result := process(baseDir, student.Logs)
				results <- map[string]interface{}{
					"student": student.Name,
					"status":  result["overallStatus"],
					"details": result["overallDetails"],
				}
			}(student)
		}

		finalResults := make([]map[string]interface{}, 0, len(existingConfig.Students))
		for i := 0; i < len(existingConfig.Students); i++ {
			finalResults = append(finalResults, <-results)
		}

		// Close the results channel
		close(results)

		// Render the results
		return c.Render("clioutput", fiber.Map{
			"Output": finalResults,
		})
	}

	return c.Redirect("/logfilelist")
}

func logfilelist(c *fiber.Ctx) error {
	baseDir := "./data"
	configPath := filepath.Join(baseDir, "config.json")
	existingConfig, _ := config.LoadConfig(configPath)

	logFiles := config.GetAllLogFiles(baseDir)
	logGroups := config.GroupLogsByBaseName(logFiles)

	existingEntries := make(map[string]config.Student)
	newEntries := make(map[string]config.LogPaths)

	for baseName, logs := range logGroups {
		found := false
		for _, student := range existingConfig.Students {
			if student.Logs.BrowserLog == logs.BrowserLog {
				existingEntries[baseName] = student
				found = true
				break
			}
		}
		if !found {
			newEntries[baseName] = logs
		}
	}

	return c.Render("logfilelist", fiber.Map{
		"LogGroups":       logGroups,
		"ExistingEntries": existingEntries,
		"NewEntries":      newEntries,
	})
}

func process(baseDir string, logPaths config.LogPaths) map[string]interface{} {
	// Derive the JSON paths from the log paths
	jsonPaths := config.DeriveJSONPaths(baseDir, logPaths)

	// Read and process the logs
	browserResult := processLogFile(logPaths.BrowserLog, jsonPaths.BrowserLogJSON, &parser.BrowserLogParser{})
	clientResult := processLogFile(logPaths.ClientLog, jsonPaths.ClientLogJSON, &parser.ClientLogParser{})
	runtimeResult := processLogFile(logPaths.RuntimeLog, jsonPaths.RuntimeLogJSON, &parser.RuntimeLogParser{})

	// Execute the check modules and get the result
	checkResults := runCheckModules(jsonPaths)

	// Combine results and determine overall status
	overallStatus := checkResults["overallStatus"].(string)
	overallDetails := []string{}

	if !browserResult || !clientResult || !runtimeResult {
		overallStatus = "error"
		overallDetails = append(overallDetails, "Error in log processing")
	}

	// Append log details if the status is cheating
	if overallStatus == "cheating" {
		for moduleName, result := range checkResults["results"].(map[string]interface{}) {
			if resultMap, ok := result.(map[string]interface{}); ok && resultMap["status"] == "cheating" {
				overallDetails = append(overallDetails, fmt.Sprintf("Module: %s", moduleName))
				if errors, ok := resultMap["errors"].([]interface{}); ok {
					for _, err := range errors {
						if errMap, ok := err.(map[string]interface{}); ok {
							overallDetails = append(overallDetails, fmt.Sprintf("%s: %s", errMap["message"], errMap["details"]))
						}
					}
				}
			}
		}
	}

	return map[string]interface{}{
		"overallStatus":  overallStatus,
		"results":        checkResults["results"],
		"overallDetails": overallDetails,
	}
}

// processLogFile reads, processes, and writes the log data to a JSON file.
func processLogFile(logPath, jsonPath string, curparser parser.LogParser) bool {
	logData, err := os.ReadFile(logPath)
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

	err = os.WriteFile(jsonPath, processedData, 0644)
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

func runCheckModules(jsonPaths config.JSONPaths) map[string]interface{} {
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
	results := make(map[string]interface{})

	for _, module := range modules {
		// Check if the log file path is set and not empty
		if module.LogPath == "" {
			results[module.Name] = map[string]interface{}{
				"status": "error",
				"errors": []map[string]string{
					{
						"message": "LogPath is not set or is empty.",
						"details": "The path to the log file is missing. Please ensure that the correct log file path is provided.",
					},
				},
			}
			if utils.CompareStatus("error", overallStatus) > 0 {
				overallStatus = "error"
			}
			continue
		}

		// Check if the log file exists and was processed
		if !utils.FileExists(module.LogPath) {
			results[module.Name] = map[string]interface{}{
				"status": "error",
				"errors": []map[string]string{
					{
						"message": "Required log file not found or not processed.",
						"details": fmt.Sprintf("The log file at path %s could not be found or was not processed.", module.LogPath),
					},
				},
			}
			if utils.CompareStatus("error", overallStatus) > 0 {
				overallStatus = "error"
			}
			continue
		}

		// Execute the module strategy
		output, err := module.Strategy.Execute(module.LogPath)
		result := utils.ParseModuleOutput(output, err)
		results[module.Name] = result

		// Update the overall status based on this module's result
		if status, ok := result["status"].(string); ok && utils.CompareStatus(status, overallStatus) > 0 {
			overallStatus = status
		}
	}

	// Return the overall status and detailed results
	return map[string]interface{}{
		"overallStatus": overallStatus,
		"results":       results,
	}
}
