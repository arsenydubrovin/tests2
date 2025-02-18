package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)


func TestNewLogEntry(t *testing.T) {
	tests := []struct {
		message string
		wantErr bool
	}{
		{
			message: "test message",
			wantErr: false,
		},
		{
			message: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		entry, err := NewLogEntry(tt.message)
		if tt.wantErr {
			assert.Error(t, err, "expected error for message: %q", tt.message)
		} else {
			assert.NoError(t, err, "unexpected error for message: %q", tt.message)
			assert.Equal(t, tt.message, entry.message, "message mismatch")
		}
	}
}

func TestLoggerMarshalJSON(t *testing.T) {
	fixedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		entries []*LogEntry
		want    string
	}{
		{
			entries: []*LogEntry{},
			want:    "[]",
		},
		{
			entries: []*LogEntry{{
				message:   "test",
				level:     INFO,
				timestamp: fixedTime,
			}},
			want: `[{"message":"test","level":"INFO","timestamp":"2024-01-01T12:00:00Z"}]`,
		},
	}

	for _, tt := range tests {
		logger := &Logger{entries: tt.entries}
		got, err := logger.MarshalJSON()
		assert.NoError(t, err, "unexpected error")
		assert.Equal(t, tt.want, string(got), "json mismatch")
	}
}

func TestLoggerIntegrationAddAndFilterEntries(t *testing.T) {
	logger := NewLogger()

	entry1, _ := NewLogEntry("info log message")
	entry1.SetLevel(INFO)
	logger.AddEntry(entry1)

	entry2, _ := NewLogEntry("warning log message")
	entry2.SetLevel(WARNING)
	logger.AddEntry(entry2)

	entry3, _ := NewLogEntry("error log message")
	entry3.SetLevel(ERROR)
	logger.AddEntry(entry3)

	filtered := logger.GetEntries(WARNING, ERROR)

	assert.Len(t, filtered, 2, "should be WARNING and ERROR only")
	assert.Equal(t, "warning log message", filtered[0].message)
	assert.Equal(t, "error log message", filtered[1].message)
}

func TestLoggerAttestation(t *testing.T) {
	logger := NewLogger()

	fixedDate := time.Date(2024, 2, 17, 15, 0, 0, 0, time.UTC)

	entry1, _ := NewLogEntry("info message")
	entry1.SetLevel(INFO)
	entry1.SetTimestamp(fixedDate)
	logger.AddEntry(entry1)

	entry2, _ := NewLogEntry("warning message")
	entry2.SetLevel(WARNING)
	entry2.SetTimestamp(fixedDate.Add(time.Hour))
	logger.AddEntry(entry2)

	entry3, _ := NewLogEntry("error message")
	entry3.SetLevel(ERROR)
	entry3.SetTimestamp(fixedDate.Add(time.Hour * 2))
	logger.AddEntry(entry3)

	filtered := logger.GetEntries(WARNING, ERROR)

	assert.Len(t, filtered, 2, "expected only WARNING and ERROR")
	assert.Equal(t, "warning message", filtered[0].message, "should be WARNING")
	assert.Equal(t, "error message", filtered[1].message, "should be ERROR")

	data, err := json.Marshal(logger)

	assert.NoError(t, err, "JSON serialization error")

	expected := `[
		{"message":"info message","level":"INFO","timestamp":"2024-02-17T15:00:00Z"},
		{"message":"warning message","level":"WARNING","timestamp":"2024-02-17T16:00:00Z"},
		{"message":"error message","level":"ERROR","timestamp":"2024-02-17T17:00:00Z"}
	]`

	assert.JSONEq(t, expected, string(data), "JSON doen't match")
}
