package main

type Level int

const (
	INFO Level = iota
	WARNING
	ERROR
)

func (l Level) String() string {
	str := "INFO"
	switch l {
	case WARNING:
		str = "WARNING"
	case ERROR:
		str = "ERROR"
	}
	return str
}
