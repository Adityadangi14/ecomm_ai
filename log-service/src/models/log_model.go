package models

import (
	"encoding/json"
	"time"
)

type LogEntry struct {
	Time    time.Time       `json:"time"`
	Level   string          `json:"level"`
	Message string          `json:"msg"`
	Source  *SlogSource     `json:"source,omitempty"`
	Attrs   json.RawMessage `json:"-"` // optional catch-all if you want dynamic fields
}

type SlogSource struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}
