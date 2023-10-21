package server

import (
	"context"
	"log/slog"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

func Listen(ctx context.Context, opts Opts) (server *Server, err error) {
	listener, err := net.Listen("tcp", opts.Address)
	if err != nil {
		return server, err
	}

	server = &Server{
		listener: listener,
		logger:   opts.Logger,
	}

	server.wg.Add(1)
	go server.acceptConnections(ctx)

	return server, nil
}

type Opts struct {
	Address string
	Logger  *slog.Logger
}

type Server struct {
	listener net.Listener
	wg       sync.WaitGroup
	logger   *slog.Logger

	isShutingDown atomic.Bool
}

func (s *Server) Shutdown() {
	op := "server.Shutdown"

	s.isShutingDown.Store(true)
	s.listener.Close()

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		s.logger.Error("server gracefully shutdown", "op", op)
		return
	case <-time.After(2 * time.Second):
		s.logger.Error("timed out waiting for server shutdown", "op", op)
		return
	}
}

func (s *Server) acceptConnections(ctx context.Context) {
	op := "server.acceptConnections"
	defer s.wg.Done()

	for {
		select {
		case <-ctx.Done():
			s.logger.Error("server closed", "op", op)
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				if s.isShutingDown.Load() {
					s.logger.Error("server closed", "op", op)
					return
				}

				s.logger.Error(err.Error(), "op", op)
				continue
			}

			go s.handleConnection(conn)
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	op := "server.handleConnection"
	defer conn.Close()

	if s.isShutingDown.Load() {
		s.logger.Error("server closed", "op", op)
		return
	}

	s.logger.Info("handle connection", "op", op)
	conn.Write([]byte("Hi from server!"))
}
