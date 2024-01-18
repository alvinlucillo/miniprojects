package processor

import "fmt"

type Logger interface {
	// Log(args ...interface{})
	Println(args ...interface{}) // Write robot messages (e.g., to stdout)
}

type StdLogger struct{}

func (s *StdLogger) Println(args ...interface{}) {
	fmt.Println(args...)
}
