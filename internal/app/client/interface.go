package client

import (
	"io"
)

// Config - config interface
type Config interface {
	ServerAddress() string
}

// Logger - logger interface
type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}

// Service - clisnt service to get sever resource
type Service interface {
	RequestResource(clientID string, rw io.ReadWriter) (resource string, err error)
}
