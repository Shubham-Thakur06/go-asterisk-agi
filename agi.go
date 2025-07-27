package agi

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// AgiSession represents an AGI session
type AgiSession struct {
	reader     *bufio.Reader
	writer     io.Writer
	env        map[string]string
	mutex      sync.Mutex
	variables  map[string]string
	debugMode  bool
	timeout    time.Duration
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// AgiResponse represents an AGI response
type AgiResponse struct {
	Status  int
	Result  int
	Data    string
	Raw     string
	EndPos  int
	Digits  string
	Timeout bool
}

// NewAgiSession creates a new AGI session
func NewAgiSession() (*AgiSession, error) {
	return NewWithContext(context.Background())
}

// NewWithContext creates a new AGISession with a context
func NewWithContext(ctx context.Context) (*AgiSession, error) {
	ctx, cancel := context.WithCancel(ctx)

	s := &AgiSession{
		reader:     bufio.NewReader(os.Stdin),
		writer:     os.Stdout,
		env:        make(map[string]string),
		variables:  make(map[string]string),
		debugMode:  false,
		timeout:    30 * time.Second,
		ctx:        ctx,
		cancelFunc: cancel,
	}

	if err := s.readEnvironment(); err != nil {
		cancel()
		return nil, errors.Wrap(err, "failed to read AGI environment")
	}

	return s, nil
}

// readEnvironment reads the AGI environment variables
func (s *AgiSession) readEnvironment() error {
	for {
		line, err := s.reader.ReadString('\n')
		if err != nil {
			return errors.Wrap(err, "failed to read environment line")
		}

		line = strings.TrimSpace(line)
		if line == "" {
			break
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return errors.Errorf("invalid environment line: %s", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		s.env[key] = value
	}

	return nil
}

// channel functions

// GetVariable gets a channel variable
func (s *AgiSession) GetVariable(name string) (string, error) {
	resp, err := s.execute(fmt.Sprintf("GET VARIABLE %s", name))
	if err != nil {
		return "", err
	}

	if resp.Status != 1 {
		return "", nil
	}

	return resp.Data, nil
}

// SetVariable sets a channel variable
func (s *AgiSession) SetVariable(name, value string) error {
	_, err := s.execute(fmt.Sprintf("SET VARIABLE %s \"%s\"", name, value))
	return err
}

// Answer answers the channel
func (s *AgiSession) Answer() error {
	_, err := s.execute("ANSWER")
	return err
}

// Hangup hangs up the channel
func (s *AgiSession) Hangup() error {
	_, err := s.execute("HANGUP")
	return err
}

// StreamFile plays a sound file
func (s *AgiSession) StreamFile(filename string, escapeDigits string) error {
	_, err := s.execute(fmt.Sprintf("STREAM FILE %s \"%s\"", filename, escapeDigits))
	return err
}

// WaitForDigit waits for a DTMF digit
func (s *AgiSession) WaitForDigit(timeout int) (string, error) {
	resp, err := s.execute(fmt.Sprintf("WAIT FOR DIGIT %d", timeout))
	if err != nil {
		return "", err
	}

	if resp.Status == 0 {
		return "", nil
	}

	return string(resp.Result), nil
}

// GetData gets data from the user
func (s *AgiSession) GetData(filename string, timeout, maxDigits int) (string, error) {
	resp, err := s.execute(fmt.Sprintf("GET DATA %s %d %d", filename, timeout, maxDigits))
	if err != nil {
		return "", err
	}

	return resp.Data, nil
}

// execute sends a command to Asterisk and waits for the response
func (s *AgiSession) execute(command string) (*AgiResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.debugMode {
		fmt.Fprintf(os.Stderr, "AGI Command: %s\n", command)
	}

	if _, err := fmt.Fprintf(s.writer, "%s\n", command); err != nil {
		return nil, errors.Wrap(err, "failed to send command")
	}

	line, err := s.reader.ReadString('\n')
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response")
	}

	if s.debugMode {
		fmt.Fprintf(os.Stderr, "AGI Response: %s", line)
	}

	return parseResponse(line)
}

// parseResponse parses an AGI response
func parseResponse(line string) (*AgiResponse, error) {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "200") {
		return nil, errors.Errorf("invalid response: %s", line)
	}

	parts := strings.SplitN(line[4:], " ", 2)
	if len(parts) != 2 {
		return nil, errors.Errorf("invalid response format: %s", line)
	}

	result, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse result")
	}

	resp := &AgiResponse{
		Status: 1,
		Result: result,
		Raw:    line,
	}

	if len(parts) > 1 {
		resp.Data = strings.TrimSpace(parts[1])
	}

	return resp, nil
}

// Close closes the AGI session
func (s *AgiSession) Close() error {
	s.cancelFunc()
	return nil
}

// SetDebug enables or disables debug mode
func (s *AgiSession) SetDebug(enabled bool) {
	s.debugMode = enabled
}

// SetTimeout sets the timeout for AGI operations
func (s *AgiSession) SetTimeout(timeout time.Duration) {
	s.timeout = timeout
}

// GetEnv gets an environment variable
func (s *AgiSession) GetEnv(key string) string {
	return s.env[key]
}
