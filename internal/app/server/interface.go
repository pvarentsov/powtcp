package server

import (
	"io"
	"time"
)

// Config - config interface
type Config interface {
	Address() string
	ShutdownTimeout() time.Duration
	ConnectionTimeout() time.Duration
}

// Logger - logger interface
type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}

// Service - server service to handle client messages
type Service interface {
	HandleMessages(clientID string, rw io.ReadWriter)
}
