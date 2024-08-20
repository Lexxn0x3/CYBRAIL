package models

type Student struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    Logs LogPaths `json:"logs"`
}

