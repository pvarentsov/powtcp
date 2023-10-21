package log

import (
	"log/slog"
	"os"
)

type Opts struct {
	Debug  bool
	Json   bool
	Writer *os.File
}

func New(opts Opts) (logger *slog.Logger) {
	writer := opts.Writer
	if writer == nil {
		writer = os.Stderr
	}

	level := new(slog.LevelVar)
	if opts.Debug {
		level.Set(slog.LevelDebug)
	}

	handlerOpts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if opts.Json {
		handler = slog.NewJSONHandler(writer, handlerOpts)
	} else {
		handler = slog.NewTextHandler(writer, handlerOpts)
	}

	return slog.New(handler)
}
