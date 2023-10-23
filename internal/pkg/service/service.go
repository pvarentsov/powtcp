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
				rw.Write(errorMessage(ErrTimeoutExceeded).Bytes())
				return
			}

			rw.Write(errorMessage(ErrInternalError).Bytes())
			return
		}

		msg, err := message.ParseMessage(rawMsg)
		if err != nil {
			s.logger.Error(err.Error(), "op", op, "clientID", clientID)
			rw.Write(errorMessage(ErrIncorrectMessageFormat).Bytes())
			return
		}

		rw.Write(msg.Bytes())
	}
}
