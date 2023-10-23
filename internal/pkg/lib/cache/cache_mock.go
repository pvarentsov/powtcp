package cache

type mockLogger struct {
	cancelSignalHandled bool
}

func (l *mockLogger) Debug(msg string, args ...any) {
	l.cancelSignalHandled = msg == "context canceled"
}
