package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollapseDuplicates(t *testing.T) {
	fixedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		logs          []*LogEntry
		from          time.Time
		to            time.Time
		expected      []*LogEntry
		expectedError string
	}{
		{
			name: "неправильный диапазон",
			logs: []*LogEntry{
				{message: "test", level: INFO, timestamp: fixedTime},
			},
			from:          fixedTime.Add(1 * time.Hour),
			to:            fixedTime,
			expectedError: "from >= to!",
		},
		{
			name:     "пустые логи",
			logs:     []*LogEntry{},
			from:     fixedTime,
			to:       fixedTime.Add(1 * time.Hour),
			expected: []*LogEntry{},
		},
		{
			name: "нет логов в диапазоне",
			logs: []*LogEntry{
				{message: "test", level: INFO, timestamp: fixedTime.Add(-2 * time.Hour)},
				{message: "test", level: INFO, timestamp: fixedTime.Add(-1 * time.Hour)},
			},
			from:     fixedTime,
			to:       fixedTime.Add(1 * time.Hour),
			expected: []*LogEntry{},
		},
		{
			name: "меньше трёх логов с одинаковым интервалом",
			logs: []*LogEntry{
				{message: "test", level: INFO, timestamp: fixedTime},
				{message: "test", level: INFO, timestamp: fixedTime.Add(10 * time.Second)},
			},
			from: fixedTime,
			to:   fixedTime.Add(1 * time.Hour),
			expected: []*LogEntry{
				{message: "test", level: INFO, timestamp: fixedTime},
				{message: "test", level: INFO, timestamp: fixedTime.Add(10 * time.Second)},
			},
		},
		{
			name: "одинаковые интервалы",
			logs: []*LogEntry{
				{message: "test", level: INFO, timestamp: fixedTime},
				{message: "test", level: INFO, timestamp: fixedTime.Add(10 * time.Second)},
				{message: "test", level: INFO, timestamp: fixedTime.Add(20 * time.Second)},
				{message: "test", level: INFO, timestamp: fixedTime.Add(30 * time.Second)},
			},
			from: fixedTime,
			to:   fixedTime.Add(1 * time.Hour),
			expected: []*LogEntry{
				{message: "test", level: INFO, timestamp: fixedTime},
				{message: "test", level: INFO, timestamp: fixedTime.Add(30 * time.Second)},
			},
		},
		{
			name: "разные интервалы",
			logs: []*LogEntry{
				{message: "test", level: INFO, timestamp: fixedTime},
				{message: "test", level: INFO, timestamp: fixedTime.Add(10 * time.Second)},
				{message: "test", level: INFO, timestamp: fixedTime.Add(40 * time.Second)},
				{message: "test", level: INFO, timestamp: fixedTime.Add(50 * time.Second)},
			},
			from: fixedTime,
			to:   fixedTime.Add(1 * time.Hour),
			expected: []*LogEntry{
				{message: "test", level: INFO, timestamp: fixedTime},
				{message: "test", level: INFO, timestamp: fixedTime.Add(10 * time.Second)},
				{message: "test", level: INFO, timestamp: fixedTime.Add(40 * time.Second)},
				{message: "test", level: INFO, timestamp: fixedTime.Add(50 * time.Second)},
			},
		},
		{
			name: "смещанные логи",
			logs: []*LogEntry{
				{message: "test1", level: INFO, timestamp: fixedTime},
				{message: "test1", level: INFO, timestamp: fixedTime.Add(10 * time.Second)},
				{message: "test1", level: INFO, timestamp: fixedTime.Add(20 * time.Second)},
				{message: "test2", level: INFO, timestamp: fixedTime.Add(5 * time.Second)},
				{message: "test2", level: INFO, timestamp: fixedTime.Add(15 * time.Second)},
			},
			from: fixedTime,
			to:   fixedTime.Add(1 * time.Hour),
			expected: []*LogEntry{
				{message: "test1", level: INFO, timestamp: fixedTime},
				{message: "test1", level: INFO, timestamp: fixedTime.Add(20 * time.Second)},
				{message: "test2", level: INFO, timestamp: fixedTime.Add(5 * time.Second)},
				{message: "test2", level: INFO, timestamp: fixedTime.Add(15 * time.Second)},
			},
		},
		{
			name: "одинаковые сообщения, но разные уровни",
			logs: []*LogEntry{
				{message: "test", level: INFO, timestamp: fixedTime},
				{message: "test", level: INFO, timestamp: fixedTime.Add(10 * time.Second)},
				{message: "test", level: INFO, timestamp: fixedTime.Add(20 * time.Second)},
				{message: "test", level: ERROR, timestamp: fixedTime.Add(10 * time.Second)},
				{message: "test", level: ERROR, timestamp: fixedTime.Add(20 * time.Second)},
			},
			from: fixedTime,
			to:   fixedTime.Add(1 * time.Hour),
			expected: []*LogEntry{
				{message: "test", level: INFO, timestamp: fixedTime},
				{message: "test", level: INFO, timestamp: fixedTime.Add(20 * time.Second)},
				{message: "test", level: ERROR, timestamp: fixedTime.Add(10 * time.Second)},
				{message: "test", level: ERROR, timestamp: fixedTime.Add(20 * time.Second)},
			},
		},
		{
			name: "логи с нарушенной хронологией",
			logs: []*LogEntry{
				{message: "test", level: INFO, timestamp: fixedTime.Add(20 * time.Second)},
				{message: "test", level: INFO, timestamp: fixedTime.Add(10 * time.Second)},
				{message: "test", level: INFO, timestamp: fixedTime},
			},
			from:          fixedTime,
			to:            fixedTime.Add(1 * time.Hour),
			expectedError: "inconsistent timestamps",
		},
		{
			name: "границы диапазона и логов совпадают",
			logs: []*LogEntry{
				{message: "test", level: INFO, timestamp: fixedTime},
				{message: "test", level: INFO, timestamp: fixedTime.Add(30 * time.Minute)},
				{message: "test", level: INFO, timestamp: fixedTime.Add(1 * time.Hour)},
			},
			from: fixedTime,
			to:   fixedTime.Add(1 * time.Hour),
			expected: []*LogEntry{
				{message: "test", level: INFO, timestamp: fixedTime},
				{message: "test", level: INFO, timestamp: fixedTime.Add(1 * time.Hour)},
			},
		},
		{
			name: "интервалы отличаются, но меньше, чем на дельту",
			logs: []*LogEntry{
				{message: "test", level: INFO, timestamp: fixedTime.Add(10 * time.Second)},
				{message: "test", level: INFO, timestamp: fixedTime.Add(20 * time.Second)},
				{message: "test", level: INFO, timestamp: fixedTime.Add(34 * time.Second)},
			},
			from: fixedTime,
			to:   fixedTime.Add(1 * time.Hour),
			expected: []*LogEntry{
				{message: "test", level: INFO, timestamp: fixedTime.Add(10 * time.Second)},
				{message: "test", level: INFO, timestamp: fixedTime.Add(34 * time.Second)},
			},
		},
		{
			name: "интервалы отличаются больше, чем на дельту",
			logs: []*LogEntry{
				{message: "test", level: INFO, timestamp: fixedTime.Add(10 * time.Second)},
				{message: "test", level: INFO, timestamp: fixedTime.Add(20 * time.Second)},
				{message: "test", level: INFO, timestamp: fixedTime.Add(36 * time.Second)},
			},
			from: fixedTime,
			to:   fixedTime.Add(1 * time.Hour),
			expected: []*LogEntry{
				{message: "test", level: INFO, timestamp: fixedTime.Add(10 * time.Second)},
				{message: "test", level: INFO, timestamp: fixedTime.Add(20 * time.Second)},
				{message: "test", level: INFO, timestamp: fixedTime.Add(36 * time.Second)},
			},
		},
		{
			name: "одинаковые ts у логов",
			logs: []*LogEntry{
				{message: "test", level: INFO, timestamp: fixedTime.Add(10 * time.Second)},
				{message: "test", level: INFO, timestamp: fixedTime.Add(10 * time.Second)},
				{message: "test", level: INFO, timestamp: fixedTime.Add(10 * time.Second)},
			},
			from: fixedTime,
			to:   fixedTime.Add(1 * time.Hour),
			expected: []*LogEntry{
				{message: "test", level: INFO, timestamp: fixedTime.Add(10 * time.Second)},
			},
		},
		{
			name: "одинковые логи, один нарушает порядок",
			logs: []*LogEntry{
				{message: "test", level: INFO, timestamp: fixedTime.Add(10 * time.Second)},
				{message: "test", level: INFO, timestamp: fixedTime.Add(11 * time.Second)},
				{message: "test", level: INFO, timestamp: fixedTime.Add(10 * time.Second)},
				{message: "test", level: INFO, timestamp: fixedTime.Add(10 * time.Second)},
			},
			from:          fixedTime,
			to:            fixedTime.Add(1 * time.Hour),
			expectedError: "inconsistent timestamps",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Logger{entries: tt.logs}
			result, err := l.CollapseDuplicates(tt.from, tt.to)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result, "exp: %+v\nres: %+v", tt.expected, result)
		})
	}
}

func TestInRange(t *testing.T) {
	fixedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		ts       time.Time
		from     time.Time
		to       time.Time
		expected bool
	}{
		{
			name:     "одинаковые границы",
			ts:       fixedTime,
			from:     fixedTime,
			to:       fixedTime.Add(1 * time.Hour),
			expected: true,
		},
		{
			name:     "ts равен верхней границе",
			ts:       fixedTime.Add(1 * time.Hour),
			from:     fixedTime,
			to:       fixedTime.Add(1 * time.Hour),
			expected: true,
		},
		{
			name:     "ts равен нижней границе",
			ts:       fixedTime,
			from:     fixedTime,
			to:       fixedTime.Add(1 * time.Hour),
			expected: true,
		},
		{
			name:     "ts в диапазоне",
			ts:       fixedTime.Add(30 * time.Minute),
			from:     fixedTime,
			to:       fixedTime.Add(1 * time.Hour),
			expected: true,
		},
		{
			name:     "ts мешьне нижней границы",
			ts:       fixedTime.Add(-1 * time.Hour),
			from:     fixedTime,
			to:       fixedTime.Add(1 * time.Hour),
			expected: false,
		},
		{
			name:     "ts больше верхней границы",
			ts:       fixedTime.Add(2 * time.Hour),
			from:     fixedTime,
			to:       fixedTime.Add(1 * time.Hour),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := inRange(tt.ts, tt.from, tt.to)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHasPeriodicity(t *testing.T) {
	tests := []struct {
		name      string
		intervals []time.Duration
		expected  bool
	}{
		{
			name:      "пустой список",
			intervals: []time.Duration{},
			expected:  true,
		},
		{
			name:      "один интервал",
			intervals: []time.Duration{10 * time.Second},
			expected:  true,
		},
		{
			name: "одинаковый интервал",
			intervals: []time.Duration{
				10 * time.Second,
				10 * time.Second,
				10 * time.Second,
			},
			expected: true,
		},
		{
			name: "интервалы отличаются, но меньше, чем на дельту",
			intervals: []time.Duration{
				10 * time.Second,
				12 * time.Second,
				11 * time.Second,
			},
			expected: true,
		},
		{
			name: "интервалы отличаются больше, чем на дельту",
			intervals: []time.Duration{
				10 * time.Second,
				16 * time.Second,
				10 * time.Second,
				16 * time.Second,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasPeriodicity(tt.intervals)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetIntervals(t *testing.T) {
	fixedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		logs          []*LogEntry
		expected      []time.Duration
		expectedError string
	}{
		{
			name: "правильные интервалы",
			logs: []*LogEntry{
				{timestamp: fixedTime},
				{timestamp: fixedTime.Add(10 * time.Second)},
				{timestamp: fixedTime.Add(20 * time.Second)},
			},
			expected: []time.Duration{
				10 * time.Second,
				10 * time.Second,
			},
		},
		{
			name: "наружение хронологии",
			logs: []*LogEntry{
				{timestamp: fixedTime.Add(20 * time.Second)},
				{timestamp: fixedTime.Add(10 * time.Second)},
			},
			expectedError: "inconsistent timestamps",
		},
		{
			name:     "пустой список",
			logs:     []*LogEntry{},
			expected: []time.Duration{},
		},
		{
			name: "один лог",
			logs: []*LogEntry{
				{timestamp: fixedTime},
			},
			expected: []time.Duration{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			intervals, err := getIntervals(tt.logs)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, intervals)
		})
	}
}
