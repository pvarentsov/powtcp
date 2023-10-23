package server

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// Listen - listen tcp connections
func Listen(ctx context.Context, opts Opts) (server *Server, err error) {
	listener, err := net.Listen("tcp", opts.Address)
	if err != nil {
		return server, err
	}

	server = &Server{
		listener: listener,
		logger:   opts.Logger,
	}

	server.shutdownWg.Add(1)
	go server.acceptConnections(ctx)

	return server, nil
}

// Opts - options to run server
type Opts struct {
	Address string
	Logger  Logger
}

// Sever - tcp server
type Server struct {
	listener net.Listener
	logger   Logger

	shutdownWg    sync.WaitGroup
	isShutingDown atomic.Bool
}

// Shutdown - shutdown server gracefully
func (s *Server) Shutdown() {
	const op = "server.Shutdown"

	s.isShutingDown.Store(true)
	s.listener.Close()

	done := make(chan struct{})
	go func() {
		s.shutdownWg.Wait()
		close(done)
	}()

	select {
	case <-done:
		s.logger.Debug("shutdown server gracefully", "op", op)
		return
	case <-time.After(2 * time.Second):
		s.logger.Debug("shutdown server by timeout", "op", op)
		return
	}
}

func (s *Server) acceptConnections(ctx context.Context) {
	const op = "server.acceptConnections"
	defer s.shutdownWg.Done()

	for {
		select {
		case <-ctx.Done():
			s.logger.Debug("context canceled", "op", op)
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				if s.isShutingDown.Load() {
					s.logger.Debug("server closed", "op", op)
					return
				}

				s.logger.Warn(err.Error(), "op", op)
				continue
			}

			go s.handleConnection(conn)
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	const op = "server.handleConnection"
	defer conn.Close()

	if s.isShutingDown.Load() {
		s.logger.Error("server closed", "op", op)
		return
	}

	_, err := conn.Write([]byte("Hi from server!"))
	if err != nil {
		s.logger.Error(fmt.Sprintf("sent message: %v", err), "op", op)
	}
}
