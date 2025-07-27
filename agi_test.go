package agi

import (
	"bufio"
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockIO struct {
	reader *strings.Reader
	writer *bytes.Buffer
}

func newMockIO(input string) *mockIO {
	return &mockIO{
		reader: strings.NewReader(input),
		writer: &bytes.Buffer{},
	}
}

func TestNewSession(t *testing.T) {
	// Test successful session creation
	t.Run("successful creation", func(t *testing.T) {
		input := "agi_network: yes\nagi_network_script: test.agi\n\n"
		mock := newMockIO(input)

		session := &AgiSession{
			reader:    bufio.NewReader(mock.reader),
			writer:    mock.writer,
			env:       make(map[string]string),
			variables: make(map[string]string),
			debugMode: false,
			timeout:   30 * time.Second,
		}

		err := session.readEnvironment()
		require.NoError(t, err)

		assert.Equal(t, "yes", session.env["agi_network"])
		assert.Equal(t, "test.agi", session.env["agi_network_script"])
	})

	// Test invalid environment format
	t.Run("invalid environment", func(t *testing.T) {
		input := "invalid line\n\n"
		mock := newMockIO(input)

		session := &AgiSession{
			reader:    bufio.NewReader(mock.reader),
			writer:    mock.writer,
			env:       make(map[string]string),
			variables: make(map[string]string),
			debugMode: false,
			timeout:   30 * time.Second,
		}

		err := session.readEnvironment()
		require.Error(t, err)
	})
}

func TestAGICommands(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		response string
		wantErr  bool
	}{
		{
			name:     "answer success",
			command:  "ANSWER\n",
			response: "200 result=0\n",
			wantErr:  false,
		},
		{
			name:     "hangup success",
			command:  "HANGUP\n",
			response: "200 result=1\n",
			wantErr:  false,
		},
		{
			name:     "get variable success",
			command:  "GET VARIABLE test\n",
			response: "200 result=1 (test)\n",
			wantErr:  false,
		},
		{
			name:     "invalid response",
			command:  "INVALID\n",
			response: "500 invalid command\n",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := newMockIO(tt.response)

			session := &AgiSession{
				reader:    bufio.NewReader(mock.reader),
				writer:    mock.writer,
				env:       make(map[string]string),
				variables: make(map[string]string),
				debugMode: false,
				timeout:   30 * time.Second,
			}

			resp, err := session.execute(tt.command[:len(tt.command)-1])

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, resp)

			// Verify command was sent correctly
			assert.Equal(t, tt.command, mock.writer.String())
		})
	}
}

func TestFastAGIServer(t *testing.T) {
	t.Run("server lifecycle", func(t *testing.T) {
		handler := HandlerFunc(func(ctx context.Context, s *AgiSession) error {
			return s.Answer()
		})

		server, err := NewFastAGIServer(":0", handler)
		require.NoError(t, err)

		// Start server in goroutine
		go func() {
			err := server.Serve()
			require.NoError(t, err)
		}()

		// Allow server to start
		time.Sleep(100 * time.Millisecond)

		// Stop server
		err = server.Stop()
		require.NoError(t, err)
	})
}

func TestResponseParsing(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *AgiResponse
		wantErr bool
	}{
		{
			name:  "simple response",
			input: "200 result=1",
			want: &AgiResponse{
				Status: 1,
				Result: 1,
				Raw:    "200 result=1",
			},
			wantErr: false,
		},
		{
			name:  "response with data",
			input: "200 result=1 (data)",
			want: &AgiResponse{
				Status: 1,
				Result: 1,
				Data:   "(data)",
				Raw:    "200 result=1 (data)",
			},
			wantErr: false,
		},
		{
			name:    "invalid response",
			input:   "invalid",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseResponse(tt.input)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want.Status, got.Status)
			assert.Equal(t, tt.want.Result, got.Result)
			assert.Equal(t, tt.want.Data, got.Data)
			assert.Equal(t, tt.want.Raw, got.Raw)
		})
	}
}
