package service

import (
	"errors"

	"github.com/pvarentsov/powtcp/internal/pkg/lib/message"
)

// Errors
var (
	ErrIncorrectMessageFormat     = errors.New("incorrect message format")
	ErrTimeoutExceeded            = errors.New("timeout exceeded")
	ErrUnknownCommand             = errors.New("unknown command")
	ErrHashcashHeaderNotFound     = errors.New("hashcah header not found")
	ErrHashcashHeaderNotCorrect   = errors.New("hashcah header not correct")
	ErrHashcashExpirationExceeded = errors.New("hashcah expiration exceeded")
	ErrInternalError              = errors.New("internal error")
	ErrResponseCommandNotcorrect  = errors.New("response command is not correct")
)

func errorMessage(err error) message.Message {
	return message.Message{
		Command: message.CommandError,
		Payload: err.Error(),
	}
}
