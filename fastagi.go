package agi

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

// FastAGIServer represents a FastAGI server
type FastAGIServer struct {
	listener   net.Listener
	handler    Handler
	wg         sync.WaitGroup
	mu         sync.Mutex
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// Handler is the interface that must be implemented to handle FastAGI requests
type Handler interface {
	Handle(ctx context.Context, s *AgiSession) error
}

// HandlerFunc is an adapter to allow the use of ordinary functions as AGI handlers
type HandlerFunc func(ctx context.Context, s *AgiSession) error

// Handle calls f(ctx, s)
func (f HandlerFunc) Handle(ctx context.Context, s *AgiSession) error {
	return f(ctx, s)
}

// NewFastAGIServer creates a new FastAGI server
func NewFastAGIServer(address string, handler Handler) (*FastAGIServer, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &FastAGIServer{
		listener:   listener,
		handler:    handler,
		ctx:        ctx,
		cancelFunc: cancel,
	}, nil
}

// Serve starts serving FastAGI requests
func (s *FastAGIServer) Serve() error {
	defer s.cancelFunc()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.ctx.Done():
				return nil
			default:
				return fmt.Errorf("failed to accept connection: %v", err)
			}
		}

		s.wg.Add(1)
		go s.handleConnection(conn)
	}
}

// Stop stops the FastAGI server
func (s *FastAGIServer) Stop() error {
	s.cancelFunc()
	err := s.listener.Close()
	s.wg.Wait()
	return err
}

func (s *FastAGIServer) handleConnection(conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()

	// Set connection timeout
	conn.SetDeadline(time.Now().Add(30 * time.Second))

	session := &AgiSession{
		reader:    bufio.NewReader(conn),
		writer:    conn,
		env:       make(map[string]string),
		variables: make(map[string]string),
		debugMode: false,
		timeout:   30 * time.Second,
	}

	// Read environment
	if err := session.readEnvironment(); err != nil {
		fmt.Printf("Failed to read environment: %v\n", err)
		return
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(s.ctx, session.timeout)
	defer cancel()

	// Handle the request
	if err := s.handler.Handle(ctx, session); err != nil {
		fmt.Printf("Handler error: %v\n", err)
	}
}

// Example usage:
/*
func main() {
    handler := HandlerFunc(func(ctx context.Context, s *Session) error {
        if err := s.Answer(); err != nil {
            return err
        }

        if err := s.StreamFile("welcome", "0123456789*#"); err != nil {
            return err
        }

        return s.Hangup()
    })

    server, err := NewFastAGIServer(":4573", handler)
    if err != nil {
        log.Fatal(err)
    }

    if err := server.Serve(); err != nil {
        log.Fatal(err)
    }
}
*/
