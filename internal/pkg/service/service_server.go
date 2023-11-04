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
	PuzzleCache  PuzzleCache
	ErrorChecker ErrorChecker
}

// NewServer - create new server-side service
func NewServer(opts ServerOpts) *Server {
	return &Server{
		logger:       opts.Logger,
		puzzleCache:  opts.PuzzleCache,
		errorChecker: opts.ErrorChecker,
	}
}

// Server - server-side service
type Server struct {
	logger       Logger
	puzzleCache  PuzzleCache
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

			clientErr := ErrInternalError
			if s.errorChecker.IsTimeout(err) {
				clientErr = ErrTimeoutExceeded
			}
			s.writeError(clientID, clientErr, rw)

			return
		}

		msg, err := message.ParseMessage(rawMsg)
		if err != nil {
			s.logger.Error(err.Error(), "op", op, "clientID", clientID)
			s.writeError(clientID, ErrIncorrectMessageFormat, rw)
			return
		}

		switch msg.Command {
		case message.CommandRequestPuzzle:
			s.responsePuzzle(clientID, msg.Payload, rw)
		case message.CommandRequestResource:
			s.responseResource(clientID, msg.Payload, rw)
		default:
			s.writeError(clientID, ErrIncorrectMessageFormat, rw)
		}
	}
}

func (s *Server) responsePuzzle(clientID string, payload string, w io.Writer) {
	const op = "service.Server.responsePuzzle"

	hashcash, err := hashcash.New(4, clientID)
	if err != nil {
		s.logger.Error(err.Error(), "op", op, "clientID", clientID)
		s.writeError(clientID, ErrIncorrectMessageFormat, w)
		return
	}

	exp := time.Now().Add(120 * time.Second)
	s.puzzleCache.AddWithExp(hashcash.Key(), struct{}{}, exp)

	msg := message.Message{
		Command: message.CommandResponsePuzzle,
		Payload: string(hashcash.Header()),
	}

	s.writeMsg(clientID, msg, w)
}

func (s *Server) responseResource(clientID string, payload string, w io.Writer) {
	const op = "service.Server.responseResource"

	hashcash, err := hashcash.ParseHeader(payload)
	if err != nil {
		s.logger.Error(err.Error(), "op", op, "clientID", clientID)
		s.writeError(clientID, ErrIncorrectMessageFormat, w)
		return
	}

	if _, ok := s.puzzleCache.Get(hashcash.Key()); !ok {
		s.writeError(clientID, ErrHashcashHeaderNotFound, w)
		return
	}
	if !hashcash.EqualResource(clientID) {
		s.writeError(clientID, ErrHashcashHeaderNotFound, w)
		return
	}

	isHashCorrect, err := hashcash.Header().IsHashCorrect(hashcash.Bits())
	if err != nil {
		s.logger.Error(err.Error(), "op", op, "clientID", clientID)
		s.writeError(clientID, ErrInternalError, w)
		return
	}
	if !isHashCorrect {
		s.writeError(clientID, ErrHashcashHeaderNotCorrect, w)
		return
	}

	msg := message.Message{
		Command: message.CommandResponseResource,
		Payload: "resource",
	}

	s.writeMsg(clientID, msg, w)
	s.puzzleCache.Delete(hashcash.Key())
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
