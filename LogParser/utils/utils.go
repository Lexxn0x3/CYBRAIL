package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ANSI color codes
const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Reset  = "\033[0m"
)

// ParseModuleOutput parses the output of a module and handles errors
func ParseModuleOutput(output string, err error) map[string]interface{} {
	if err != nil {
		return map[string]interface{}{
			"status": "error",
			"errors": []map[string]string{
				{
					"message": "Execution failed",
					"details": err.Error(),
				},
			},
		}
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		return map[string]interface{}{
			"status": "error",
			"errors": []map[string]string{
				{
					"message": "Failed to parse output",
					"details": err.Error(),
				},
			},
		}
	}

	return result
}

// PrintOverallStatus prints the overall status and individual check statuses with appropriate formatting
func PrintOverallStatus(overallStatus string, results map[string]interface{}) {
	PrintStatus("Overall Status", overallStatus, "")
	for name, result := range results {
		resultMap := result.(map[string]interface{})
		status := resultMap["status"].(string)
		PrintStatus(fmt.Sprintf("Check Module: %s", name), status, "|-- ")
		if errors, ok := resultMap["errors"].([]interface{}); ok {
			for i, err := range errors {
				errMap := err.(map[string]interface{})
				isLast := i == len(errors)-1
				PrintIndented(fmt.Sprintf("Error: %s", errMap["message"].(string)), "|    ", isLast)
				PrintIndented(fmt.Sprintf("Details: %s", errMap["details"].(string)), "|    |  ", isLast)
			}
		}
	}
}

// CompareStatus compares two statuses and returns an integer indicating their relative order
func CompareStatus(status1, status2 string) int {
	statusOrder := map[string]int{
		"success":  0,
		"error":    1,
		"cheating": 2,
	}

	return statusOrder[status1] - statusOrder[status2]
}

// PrintStatus prints the status with the appropriate color
func PrintStatus(label, status, indent string) {
	color := Reset
	switch status {
	case "success":
		color = Green
	case "cheating":
		color = Red
	case "error":
		color = Yellow
	}
	fmt.Printf("%s%s%s: %s%s\n", indent, label, color, status, Reset)
}

// PrintIndented prints the indented message with structure
func PrintIndented(message, indent string, isLast bool) {
	lines := strings.Split(message, "\n")
	for i, line := range lines {
		if isLast && i == len(lines)-1 {
			fmt.Printf("%s    |-- %s\n", indent[:len(indent)-4], line) // Remove the last "|  " and replace with "   "
		} else {
			fmt.Printf("%s|-- %s\n", indent, line)
		}
	}
}

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
