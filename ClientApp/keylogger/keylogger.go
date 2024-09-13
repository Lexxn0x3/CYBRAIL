package keylogger

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kindlyfire/go-keylogger"
)

// KeyEventIntervalThreshold defines how long of a pause qualifies as a new event
const KeyEventIntervalThreshold = 2.0 // 2 seconds
const ExitSequence = "exit"           // The sequence of keys that will trigger the program to exit

// Event represents a set of keypress intervals
type Event struct {
	Intervals []float64 `json:"intervals"`
}

// TypingSession holds multiple events
type TypingSession struct {
	Events []Event `json:"events"`
}

// AddInterval adds a new interval to the latest event or creates a new event if needed
func (ts *TypingSession) AddInterval(interval float64, isNewEvent bool) {
	if isNewEvent || len(ts.Events) == 0 {
		ts.Events = append(ts.Events, Event{}) // Start a new event
	}
	ts.Events[len(ts.Events)-1].Intervals = append(ts.Events[len(ts.Events)-1].Intervals, interval)
}

// ToJSON converts the TypingSession to JSON format
func (ts *TypingSession) ToJSON() ([]byte, error) {
	return json.MarshalIndent(ts, "", "    ")
}

// StartKeylogger starts the keylogger and returns the typing session data in JSON format.
// exitSequenceEnabled allows control over whether the "exit" sequence will end the session.
func StartKeylogger(exitSequenceEnabled bool, termination chan bool) ([]byte, error) {
	// Create a new keylogger
	kl := keylogger.NewKeylogger()

	fmt.Println("Start typing... Press 'exit' to quit if exit sequence is enabled.")

	var session TypingSession
	var prevTime time.Time
	var inputBuffer strings.Builder

loop:
	for {
		select {
		case <-termination:
			// Exit the loop
			break loop
		default:
			// Capture the key event
			key := kl.GetKey()

			// Skip empty key events
			if key.Empty {
				time.Sleep(10 * time.Millisecond)
				continue
			}

			currentTime := time.Now()

			// If this is not the first keypress, calculate the interval
			if !prevTime.IsZero() {
				interval := currentTime.Sub(prevTime).Seconds()
				isNewEvent := interval > KeyEventIntervalThreshold // Detect if it's a new event
				session.AddInterval(interval, isNewEvent)
			}

			prevTime = currentTime

			// Output the key pressed in a more readable format
			fmt.Printf("You pressed: %s at %v (Keycode: %d)\n", string(key.Rune), currentTime, key.Keycode)

			// Add the key to the input buffer for exit sequence detection
			if exitSequenceEnabled {
				inputBuffer.WriteRune(key.Rune)
				if strings.Contains(inputBuffer.String(), ExitSequence) {
					fmt.Println("Exit sequence detected. Exiting...")
					break
				}
			}
		}
	}

	// Convert session to JSON
	return session.ToJSON()
}

// SaveTypingSessionToFile saves the typing session to a file
func SaveTypingSessionToFile(jsonData []byte, filename string) error {
	return os.WriteFile(filename, jsonData, 0644)
}
