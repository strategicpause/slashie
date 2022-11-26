package logger

import (
	"fmt"
	"runtime"
)

type Logger interface {
	Debug(message string)
	Debugf(format string, a ...any)
	Info(message string)
	Warn(message string)
	Error(message string)
}

type nullOutputLogger struct{}

func NewNullOutputLogger() Logger {
	return &nullOutputLogger{}
}

func (s *nullOutputLogger) Debugf(_ string, _ ...any) {
	// noop
}

func (s *nullOutputLogger) Debug(_ string) {
	// noop
}

func (s *nullOutputLogger) Info(_ string) {
	// noop
}

func (s *nullOutputLogger) Warn(_ string) {
	// noop
}

func (s *nullOutputLogger) Error(_ string) {
	// noop
}

type stdOutLogger struct{}

func NewStdOutLogger() Logger {
	return &stdOutLogger{}
}

func (s *stdOutLogger) Debugf(format string, a ...any) {
	_, file, line, _ := runtime.Caller(1)
	s.debug(file, line, fmt.Sprintf(format, a...))
}

func (s *stdOutLogger) Debug(message string) {
	_, file, line, _ := runtime.Caller(1)
	s.debug(file, line, message)
}

func (s *stdOutLogger) debug(file string, line int, message string) {
	fmt.Printf("%s:%d %s\n", file, line, message)
}

func (s *stdOutLogger) Info(message string) {
	fmt.Println(message)
}

func (s *stdOutLogger) Warn(message string) {
	fmt.Println(message)
}

func (s *stdOutLogger) Error(message string) {
	fmt.Println(message)
}
