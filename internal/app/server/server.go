package server

import (
	"context"
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
		service:  opts.Service,
	}

	server.shutdownWg.Add(1)
	go server.acceptConnections(ctx)

	return server, nil
}

// Opts - options to run server
type Opts struct {
	Address string
	Logger  Logger
	Service Service
}

// Sever - tcp server
type Server struct {
	listener net.Listener
	logger   Logger
	service  Service

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
	case <-time.After(2000 * time.Millisecond):
		s.logger.Debug("shutdown server by timeout", "op", op)
		return
	}
}

func (s *Server) acceptConnections(ctx context.Context) {
	const op = "server.acceptConnections"
	defer s.shutdownWg.Done()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.isShutingDown.Load() {
				s.logger.Debug("server closed", "op", op)
				return
			}

			s.logger.Error(err.Error(), "op", op)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	const op = "server.handleConnection"
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(60000 * time.Millisecond))

	if s.isShutingDown.Load() {
		s.logger.Error("server closed", "op", op)
		return
	}

	s.service.HandleMessages(conn.LocalAddr().String(), conn)
}
