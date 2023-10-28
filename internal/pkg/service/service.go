package service

import (
	"bufio"
	"io"

	"github.com/pvarentsov/powtcp/internal/pkg/lib/message"
)

// Opts - options to create new cache instance
type ServerOpts struct {
	Cache        Cache
	Logger       Logger
	ErrorChecker ErrorChecker
}

// NewServer - create new server-side service
func NewServer(opts ServerOpts) *Server {
	return &Server{
		cache:        opts.Cache,
		logger:       opts.Logger,
		errorChecker: opts.ErrorChecker,
	}
}

// Server - server-side service
type Server struct {
	cache        Cache
	logger       Logger
	errorChecker ErrorChecker
}

// HandleMessages - handle client messages
func (s *Server) HandleMessages(clientID string, rw io.ReadWriter) {
	const op = "service.Server.HandleMessages"
	msgReader := bufio.NewReader(rw)

	for {
		rawMsg, err := msgReader.ReadString(message.DelimiterMessage)
		if err != nil {
			s.logger.Error(err.Error(), "op", op, "clientID", clientID)

			if s.errorChecker.IsTimeout(err) {
				s.responseError(clientID, ErrTimeoutExceeded, rw)
				return
			}

			s.responseError(clientID, ErrInternalError, rw)
			return
		}

		msg, err := message.ParseMessage(rawMsg)
		if err != nil {
			s.logger.Error(err.Error(), "op", op, "clientID", clientID)
			s.responseError(clientID, ErrIncorrectMessageFormat, rw)
			return
		}

		switch msg.Command {
		case message.CommandRequestPuzzle:
			s.responsePuzzle(clientID, msg.Payload, rw)
		case message.CommandRequestResource:
			s.responseResource(clientID, msg.Payload, rw)
		default:
			s.responseError(clientID, ErrIncorrectMessageFormat, rw)
		}
	}
}

func (s *Server) responsePuzzle(clientID string, payload string, w io.Writer) {
	const op = "service.Server.responsePuzzle"

	msg := message.Message{
		Command: message.CommandResponsePuzzle,
	}

	if _, err := w.Write(msg.Bytes()); err != nil {
		s.logger.Error(err.Error(), "op", op, "clientID", clientID)
	}
}

func (s *Server) responseResource(clientID string, payload string, w io.Writer) {
	const op = "service.Server.responseResource"

	msg := message.Message{
		Command: message.CommandResponseResource,
	}

	if _, err := w.Write(msg.Bytes()); err != nil {
		s.logger.Error(err.Error(), "op", op, "clientID", clientID)
	}
}

func (s *Server) responseError(clientID string, handleErr error, w io.Writer) {
	const op = "service.Server.responseError"

	if _, err := w.Write(errorMessage(handleErr).Bytes()); err != nil {
		s.logger.Error(err.Error(), "op", op, "clientID", clientID)
	}
}
