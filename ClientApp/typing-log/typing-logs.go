
package main

import (
    "encoding/json"
    "fmt"
    "os"
    "time"

    "github.com/eiannone/keyboard"
)

// KeyEventIntervalThreshold defines how long of a pause qualifies as a new event
const KeyEventIntervalThreshold = 2.0 // 2 seconds

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

func main() {
    // Start capturing keyboard events
    if err := keyboard.Open(); err != nil {
        panic(err)
    }
    defer keyboard.Close()

    fmt.Println("Start typing... Press ESC to quit.")

    var session TypingSession
    var prevTime time.Time

    for {
        char, key, err := keyboard.GetKey()
        if err != nil {
            panic(err)
        }

        if key == keyboard.KeyEsc {
            break
        }

        currentTime := time.Now()

        // If this is not the first keypress, calculate the interval
        if !prevTime.IsZero() {
            interval := currentTime.Sub(prevTime).Seconds()
            isNewEvent := interval > KeyEventIntervalThreshold // Detect if it's a new event
            session.AddInterval(interval, isNewEvent)
        }

        prevTime = currentTime

        fmt.Printf("You pressed: %q at %v\n", char, currentTime)
    }

    // Convert session to JSON
    jsonOutput, err := session.ToJSON()
    if err != nil {
        panic(err)
    }

    fmt.Println("Session data in JSON format (with multiple events):")
    fmt.Println(string(jsonOutput))

    // Optionally, write JSON to a file
    if err := os.WriteFile("typing_intervals.json", jsonOutput, 0644); err != nil {
        panic(err)
    }
}

