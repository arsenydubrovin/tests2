package main

import "fmt"

func main() {
	logger := NewLogger()

	entry1, _ := NewLogEntry("info log message")
	logger.AddEntry(entry1)

	entry2, _ := NewLogEntry("warning log message")
	entry2.SetLevel(WARNING)
	logger.AddEntry(entry2)

	entry3, _ := NewLogEntry("error log message")
	entry3.SetLevel(ERROR)
	logger.AddEntry(entry3)

	filtered := logger.GetEntries(WARNING, ERROR)
	for _, entry := range filtered {
		fmt.Println(entry.String())
	}
}
