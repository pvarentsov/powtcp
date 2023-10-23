package cache

// Logger - logger interface
type Logger interface {
	Debug(msg string, args ...any)
}
