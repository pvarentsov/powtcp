package service

import "time"

// Cache - cache interface
type Cache interface {
	Add(k, v string, exp time.Time)
	Get(k string) (v string, ok bool)
}

// Logger - logger interface
type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

// ErrorChecker - error checker interface
type ErrorChecker interface {
	IsTimeout(err error) bool
}
