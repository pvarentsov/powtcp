package service

import (
	"errors"

	"github.com/pvarentsov/powtcp/internal/pkg/lib/message"
)

// Errors
var (
	ErrIncorrectMessageFormat = errors.New("incorrect message format")
	ErrTimeoutExceeded        = errors.New("timeout exceeded")
	ErrUnknownCommand         = errors.New("unknown command")
	ErrInternalError          = errors.New("internal error")
)

func errorMessage(err error) message.Message {
	return message.Message{
		Command: message.CommandError,
		Payload: err.Error(),
	}
}
