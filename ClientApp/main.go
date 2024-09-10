package main

import (
	"ClientApp/keylogger"
	"fmt"
)

func main() {
	// Start the keylogger with the exit sequence enabled
	jsonData, err := keylogger.StartKeylogger(true)
	if err != nil {
		panic(err)
	}

	// Output JSON
	fmt.Println("Session data in JSON format:")
	fmt.Println(string(jsonData))

	// Save to file
	err = keylogger.SaveTypingSessionToFile(jsonData, "typing_intervals.json")
	if err != nil {
		panic(err)
	}
}
