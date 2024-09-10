package main

import (
    "encoding/json"
    "fmt"
    "os"
    "time"

    "github.com/eiannone/keyboard"
)

// KeyEvent represents a key press with a timestamp
type KeyEvent struct {
    Key       rune    `json:"key"`
    Timestamp string  `json:"timestamp"`
}

// TypingSession holds the sequence of key events
type TypingSession struct {
    Events []KeyEvent `json:"events"`
}

// AddEvent adds a new KeyEvent to the TypingSession
func (ts *TypingSession) AddEvent(key rune) {
    ts.Events = append(ts.Events, KeyEvent{
        Key:       key,
        Timestamp: time.Now().Format(time.RFC3339Nano),
    })
}

// CalculateIntervals calculates the intervals between consecutive key events
func (ts *TypingSession) CalculateIntervals() []time.Duration {
    intervals := make([]time.Duration, 0)
    for i := 1; i < len(ts.Events); i++ {
        start, _ := time.Parse(time.RFC3339Nano, ts.Events[i-1].Timestamp)
        end, _ := time.Parse(time.RFC3339Nano, ts.Events[i].Timestamp)
        interval := end.Sub(start)
        intervals = append(intervals, interval)
    }
    return intervals
}

// ToJSON converts TypingSession to JSON format
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

    for {
        char, key, err := keyboard.GetKey()
        if err != nil {
            panic(err)
        }

        if key == keyboard.KeyEsc {
            break
        }

        session.AddEvent(char)
        fmt.Printf("You pressed: %q at %v\n", char, time.Now())
    }

    jsonOutput, err := session.ToJSON()
    if err != nil {
        panic(err)
    }

    fmt.Println("Session data in JSON format:")
    fmt.Println(string(jsonOutput))

    // Optionally, write JSON to a file
    if err := os.WriteFile("typing_session.json", jsonOutput, 0644); err != nil {
        panic(err)
    }
}

