package service

import "time"

// PuzzleCache - puzzle cache interface
type PuzzleCache interface {
	AddWithExp(k string, v struct{}, exp time.Time)
	Get(k string) (v struct{}, ok bool)
	Delete(k string)
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

// ServerConfig - server config interface
type ServerConfig interface {
	PuzzleTTL() time.Duration
	PuzzleZeroBits() int
}

// ClientConfig - client config interface
type ClientConfig interface {
	PuzzleComputeMaxAttempts() int
}
