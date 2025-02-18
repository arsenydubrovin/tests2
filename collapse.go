package main

import (
	"errors"
	"reflect"
	"time"
)

func (l *Logger) CollapseDuplicates(from, to time.Time) ([]*LogEntry, error) {
	if len(l.entries) == 0 {
		return []*LogEntry{}, nil
	}

	if from.After(to) {
		return nil, errors.New("from >= to!")
	}

	logsInRange := make([]*LogEntry, 0)
	for _, log := range l.entries {
		if inRange(log.timestamp, from, to) {
			logsInRange = append(logsInRange, log)
		}
	}

	type logKey struct {
		message string
		level   Level
	}
	groups := make(map[logKey][]*LogEntry)

	for _, log := range logsInRange {
		key := logKey{log.message, log.level}
		groups[key] = append(groups[key], log)
	}

	result := make([]*LogEntry, 0)

	for _, logs := range groups {
		if len(logs) < 3 {
			result = append(result, logs...)
			continue
		}

		intervals, err := getIntervals(logs)
		if err != nil {
			return nil, err
		}

		if !hasPeriodicity(intervals) {
			result = append(result, logs...)
			continue
		}

		if reflect.DeepEqual(logs[0], logs[len(logs)-1]) {
			result = append(result, logs[0])
			continue
		}

		result = append(result, logs[0], logs[len(logs)-1])
	}

	return result, nil
}

func inRange(ts, from, to time.Time) bool {
	return (ts.Equal(from) || ts.After(from)) && (ts.Equal(to) || ts.Before(to))
}

const delta = 5 * time.Second

func hasPeriodicity(intervals []time.Duration) bool {
	for i := 0; i < len(intervals)-1; i++ {
		diff := intervals[i] - intervals[i+1]
		if diff < 0 {
			diff = -diff
		}
		if diff > delta {
			return false
		}
	}
	return true
}

func getIntervals(logs []*LogEntry) ([]time.Duration, error) {
	if len(logs) == 0 {
		return []time.Duration{}, nil
	}
	intervals := make([]time.Duration, len(logs)-1)
	for i := 0; i < len(logs)-1; i++ {
		interval := logs[i+1].timestamp.Sub(logs[i].timestamp)
		if interval < 0 {
			return nil, errors.New("inconsistent timestamps")
		}
		intervals[i] = interval
	}
	return intervals, nil
}
