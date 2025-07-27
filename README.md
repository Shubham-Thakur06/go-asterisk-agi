# go-asterisk-agi

A robust and feature-complete Asterisk Gateway Interface (AGI) library for Go applications. This library provides a high-performance, type-safe way to interact with Asterisk PBX through AGI scripts.

## Features

- Full AGI command support with type-safe interfaces
- FastAGI server implementation with connection pooling
- MRCP (Media Resource Control Protocol) support for speech recognition and synthesis
- Context-aware operations with proper resource cleanup
- Thread-safe command execution
- Automatic error handling and recovery
- Comprehensive logging and debugging capabilities
- Environment variable management
- String escaping and command parsing utilities
- Extensive test coverage

## Installation

```bash
go get github.com/Shubham-Thakur06/go-asterisk-agi
```

## Quick Start

### Regular AGI Script

```go
package main

import (
    "log"
    agi "github.com/Shubham-Thakur06/go-asterisk-agi"
)

func main() {
    // Create new AGI session
    session, err := agi.New()
    if err != nil {
        log.Fatal(err)
    }
    defer session.Close()

    // Enable debug mode
    session.SetDebug(true)

    // Answer the call
    if err := session.Answer(); err != nil {
        log.Fatal(err)
    }

    // Play welcome message
    if err := session.StreamFile("welcome", "0123456789*#"); err != nil {
        log.Fatal(err)
    }

    // Get user input
    digits, err := session.GetData("enter-account", 5000, 4)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("User entered: %s", digits)
}
```

### FastAGI Server

```go
package main

import (
    "context"
    "log"
    agi "github.com/Shubham-Thakur06/go-asterisk-agi"
)

func main() {
    handler := agi.HandlerFunc(func(ctx context.Context, s *agi.Session) error {
        s.SetDebug(true)
        
        if err := s.Answer(); err != nil {
            return err
        }
        
        if err := s.StreamFile("welcome", "0123456789*#"); err != nil {
            return err
        }
        
        return s.Hangup()
    })

    server, err := agi.NewFastAGIServer(":4573", handler)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("FastAGI server listening on :4573")
    if err := server.Serve(); err != nil {
        log.Fatal(err)
    }
}
```

## Core Features

### Session Management

- `New()` - Create a new AGI session
- `NewWithContext(ctx)` - Create a session with context
- `Close()` - Clean up resources
- `SetDebug(enabled)` - Enable/disable debug logging
- `SetTimeout(duration)` - Set operation timeout

### Basic Channel Operations

- `Answer()` - Answer channel
- `Hangup()` - Hangup channel
- `ChannelStatus()` - Get channel status
- `Execute(app, ...options)` - Execute Asterisk application

### Variable Management

- `GetVariable(name)` - Get channel variable
- `SetVariable(name, value)` - Set channel variable
- `GetEnv(key)` - Get AGI environment variable

### Audio Operations

- `StreamFile(filename, digits)` - Play audio file
- `WaitForDigit(timeout)` - Wait for DTMF input
- `GetData(filename, timeout, maxDigits)` - Get user input
- `SayNumber(num, digits)` - Say number
- `SayDigits(digits, escape)` - Say digits
- `SayDateTime(timestamp, escape, format, timezone)` - Say date/time

### Database Operations

- `DatabaseGet(family, key)` - Get from AstDB
- `DatabasePut(family, key, value)` - Put into AstDB
- `DatabaseDel(family, key)` - Delete from AstDB

### MRCP Speech Support

```go
// Speech Recognition
recog := &agi.MRCPRecog{
    Grammar:   "grammar.xml",
    Timeout:   5000,
    Options:   map[string]string{"confidence": "0.7"},
    ResultVar: "recognition_result",
}
if err := session.SpeechRecognize(recog); err != nil {
    log.Fatal(err)
}

// Speech Synthesis
synth := &agi.MRCPSynth{
    Text:    "Welcome to the system",
    Options: map[string]string{"voice": "male"},
}
if err := session.SpeechSpeak(synth); err != nil {
    log.Fatal(err)
}
```

### Utility Functions

- `EscapeString(s)` - Escape string for AGI commands
- `UnescapeString(s)` - Unescape AGI response strings
- `ParseAGIResult(s)` - Parse AGI result string
- `ParseAGIEnv(input)` - Parse AGI environment
- `FormatDateTime(format)` - Format date/time string
- `SplitCommand(cmd)` - Split AGI command into parts
- `JoinCommand(parts)` - Join command parts with escaping

## Asterisk Configuration

### extensions.conf Example

```asterisk
[default]
; Regular AGI script
exten => 1000,1,AGI(script.go)

; FastAGI server
exten => 1001,1,AGI(agi://localhost:4573)
```

### Debugging

Enable debug mode to see all AGI commands and responses:

```go
session.SetDebug(true)
```

Debug output will show:
```
AGI Command: ANSWER
AGI Response: 200 result=0
AGI Command: STREAM FILE welcome "0123456789*#"
AGI Response: 200 result=0
```

## Error Handling

The library uses `github.com/pkg/errors` for better error handling:

```go
if err := session.Answer(); err != nil {
    switch {
    case strings.Contains(err.Error(), "channel unavailable"):
        // Handle channel unavailable
    case strings.Contains(err.Error(), "timeout"):
        // Handle timeout
    default:
        // Handle other errors
    }
}
```

## Thread Safety

All AGI operations are thread-safe. The library handles concurrent access to the AGI session using mutexes.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Asterisk Team for their excellent documentation
- Go community for best practices and patterns