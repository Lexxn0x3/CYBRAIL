package keylogger

import (
	"fmt"
)

func RunKeywatch(termination chan bool) {
	// Start the keylogger with the exit sequence enabled
	jsonData, err := StartKeylogger(false, termination)
	if err != nil {
		panic(err)
	}

	// Output JSON
	fmt.Println("Session data in JSON format:")
	fmt.Println(string(jsonData))

	// Save to file
	err = SaveTypingSessionToFile(jsonData, "typing_intervals.json")
	if err != nil {
		panic(err)
	}
}
