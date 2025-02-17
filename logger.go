package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type LogEntry struct {
	message   string
	level     Level
	timestamp time.Time
}

func NewLogEntry(message string) (*LogEntry, error) {
	if message == "" {
		return nil, errors.New("empty message")
	}
	return &LogEntry{
		message: message,
		level:   INFO,
	}, nil
}

func (l *LogEntry) SetLevel(level Level) {
	l.level = level
}

func (l *LogEntry) SetTimestamp(timestamp time.Time) {
	l.timestamp = timestamp
}

func (l *LogEntry) String() string {
	return fmt.Sprintf("[%s] %s: %s",
		l.timestamp.Format(time.RFC3339),
		l.level.String(),
		l.message)
}

type Logger struct {
	entries []*LogEntry
}

func NewLogger() *Logger {
	return &Logger{
		entries: make([]*LogEntry, 0),
	}
}

func (l *Logger) AddEntry(entry *LogEntry) {
	if entry.timestamp.IsZero() {
		entry.SetTimestamp(time.Now())
	}
	l.entries = append(l.entries, entry)
}

func (l *Logger) GetEntries(level ...Level) []*LogEntry {
	if len(level) == 0 {
		return l.entries
	}

	levelMap := make(map[Level]bool)
	for _, lvl := range level {
		levelMap[lvl] = true
	}

	filtered := make([]*LogEntry, 0)
	for _, entry := range l.entries {
		if levelMap[entry.level] {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func (l *Logger) MarshalJSON() ([]byte, error) {
	type jsonEntry struct {
		Message   string    `json:"message"`
		Level     string    `json:"level"`
		Timestamp time.Time `json:"timestamp"`
	}
	entries := make([]jsonEntry, len(l.entries))
	for i, entry := range l.entries {
		entries[i] = jsonEntry{
			Message:   entry.message,
			Level:     entry.level.String(),
			Timestamp: entry.timestamp,
		}
	}
	return json.Marshal(entries)
}
