package logger

import (
	"fmt"
	"runtime"
)

type Logger interface {
	Debug(message string)
	Debugf(format string, a ...any)
	Info(message string)
	Infof(format string, a ...any)
	Warn(message string)
	Warnf(format string, a ...any)
	Error(message string)
	Errorf(format string, a ...any)
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

func (s *nullOutputLogger) Infof(format string, a ...any) {
	// noop
}

func (s *nullOutputLogger) Warn(_ string) {
	// noop
}

func (s *nullOutputLogger) Warnf(format string, a ...any) {
	// noop
}

func (s *nullOutputLogger) Error(_ string) {
	// noop
}

func (s *nullOutputLogger) Errorf(format string, a ...any) {
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

func (s *stdOutLogger) Infof(format string, a ...any) {
	s.Info(fmt.Sprintf(format, a...))
}

func (s *stdOutLogger) Warn(message string) {
	fmt.Println(message)
}

func (s *stdOutLogger) Warnf(format string, a ...any) {
	s.Warn(fmt.Sprintf(format, a...))
}

func (s *stdOutLogger) Error(message string) {
	fmt.Println(message)
}

func (s *stdOutLogger) Errorf(format string, a ...any) {
	s.Error(fmt.Sprintf(format, a...))
}
