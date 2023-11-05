package service

import (
	"bufio"
	"io"
	"time"

	"github.com/pvarentsov/powtcp/internal/pkg/lib/hashcash"
	"github.com/pvarentsov/powtcp/internal/pkg/lib/message"
)

// Opts - options to create new cache instance
type ServerOpts struct {
	Logger       Logger
	Config       ServerConfig
	PuzzleCache  PuzzleCache
	ErrorChecker ErrorChecker
}

// NewServer - create new server-side service
func NewServer(opts ServerOpts) *Server {
	return &Server{
		logger:       opts.Logger,
		config:       opts.Config,
		puzzleCache:  opts.PuzzleCache,
		errorChecker: opts.ErrorChecker,
	}
}

// Server - server-side service
type Server struct {
	logger       Logger
	config       ServerConfig
	puzzleCache  PuzzleCache
	errorChecker ErrorChecker
}

// HandleMessages - handle client messages
func (s *Server) HandleMessages(clientID string, rw io.ReadWriter) {
	const op = "service.Server.HandleMessages"

	s.logger.Info("connected new client", "clientID", clientID)

	for {
		rawMsg, err := bufio.NewReader(rw).ReadString(message.DelimiterMessage)
		if err != nil {
			clientErr := ErrInternalError
			if s.errorChecker.IsTimeout(err) {
				clientErr = ErrTimeoutExceeded
				s.logger.Info(clientErr.Error(), "clientID", clientID)
			} else {
				s.logger.Error(err.Error(), "op", op, "clientID", clientID)
			}
			s.writeError(clientID, clientErr, rw)

			return
		}

		msg, err := message.ParseMessage(rawMsg)
		if err != nil {
			s.logger.Info(ErrIncorrectMessageFormat.Error(), "clientID", clientID, "message", rawMsg)
			s.writeError(clientID, ErrIncorrectMessageFormat, rw)
			return
		}

		switch msg.Command {
		case message.CommandRequestPuzzle:
			s.responsePuzzle(clientID, msg.Payload, rw)
		case message.CommandRequestResource:
			s.responseResource(clientID, msg.Payload, rw)
			return
		default:
			s.writeError(clientID, ErrIncorrectMessageFormat, rw)
			return
		}
	}
}

func (s *Server) responsePuzzle(clientID string, payload string, w io.Writer) {
	const op = "service.Server.responsePuzzle"

	s.logger.Info("requested new puzzle", "clientID", clientID)

	hashcash, err := hashcash.New(s.config.PuzzleZeroBits(), clientID)
	if err != nil {
		s.logger.Error(err.Error(), "op", op, "clientID", clientID)
		s.writeError(clientID, ErrInternalError, w)
		return
	}

	exp := time.Now().Add(s.config.PuzzleTTL())
	s.puzzleCache.AddWithExp(hashcash.Key(), struct{}{}, exp)

	msg := message.Message{
		Command: message.CommandResponsePuzzle,
		Payload: string(hashcash.Header()),
	}

	s.writeMsg(clientID, msg, w)
	s.logger.Info("puzzle sent", "clientID", clientID, "puzzle", msg.Payload)
}

func (s *Server) responseResource(clientID string, payload string, w io.Writer) {
	const op = "service.Server.responseResource"

	s.logger.Info("requested resource", "clientID", clientID, "solution", payload)

	hashcash, err := hashcash.ParseHeader(payload)
	if err != nil {
		s.logger.Info(ErrHashcashHeaderNotCorrect.Error(), "clientID", clientID, "header", payload)
		s.writeError(clientID, ErrHashcashHeaderNotCorrect, w)
		return
	}

	if _, ok := s.puzzleCache.Get(hashcash.Key()); !ok {
		s.logger.Info(ErrHashcashHeaderNotFound.Error(), "clientID", clientID, "header", payload)
		s.writeError(clientID, ErrHashcashHeaderNotFound, w)
		return
	}
	if !hashcash.EqualResource(clientID) {
		s.logger.Info(ErrHashcashHeaderNotFound.Error(), "clientID", clientID, "header", payload)
		s.writeError(clientID, ErrHashcashHeaderNotFound, w)
		return
	}
	if !hashcash.IsActual(s.config.PuzzleTTL()) {
		s.logger.Info(ErrHashcashExpirationExceeded.Error(), "clientID", clientID, "header", payload)
		s.writeError(clientID, ErrHashcashExpirationExceeded, w)
		return
	}

	isHashCorrect, err := hashcash.Header().IsHashCorrect(hashcash.Bits())
	if err != nil {
		s.logger.Error(err.Error(), "op", op, "clientID", clientID)
		s.writeError(clientID, ErrInternalError, w)
		return
	}
	if !isHashCorrect {
		s.logger.Info(ErrHashcashHeaderNotCorrect.Error(), "clientID", clientID, "header", payload)
		s.writeError(clientID, ErrHashcashHeaderNotCorrect, w)
		return
	}

	msg := message.Message{
		Command: message.CommandResponseResource,
		Payload: "resource",
	}

	s.writeMsg(clientID, msg, w)
	s.puzzleCache.Delete(hashcash.Key())
	s.logger.Info("resource sent", "clientID", clientID, "resource", msg.Payload)
}

func (s *Server) writeMsg(clientID string, msg message.Message, w io.Writer) {
	const op = "service.Server.writeMsg"

	if _, err := w.Write(msg.Bytes()); err != nil {
		s.logger.Error(err.Error(), "op", op, "clientID", clientID)
	}
}

func (s *Server) writeError(clientID string, handleErr error, w io.Writer) {
	const op = "service.Server.writeError"

	if _, err := w.Write(errorMessage(handleErr).Bytes()); err != nil {
		s.logger.Error(err.Error(), "op", op, "clientID", clientID)
	}
}
