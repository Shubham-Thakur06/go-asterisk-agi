package agi

import (
	"fmt"
	"strings"
)

// ChannelStatus gets the status of the current channel
func (s *AgiSession) ChannelStatus() (int, error) {
	resp, err := s.execute("CHANNEL STATUS")
	if err != nil {
		return 0, err
	}
	return resp.Result, nil
}

// Execute executes a dialplan application
func (s *AgiSession) Execute(application string, options ...string) error {
	cmd := fmt.Sprintf("EXEC %s", application)
	if len(options) > 0 {
		cmd += fmt.Sprintf(" \"%s\"", strings.Join(options, ","))
	}
	_, err := s.execute(cmd)
	return err
}

// GetOption streams a file and gets a digit
func (s *AgiSession) GetOption(filename string, escapeDigits string, timeout int) (string, error) {
	cmd := fmt.Sprintf("GET OPTION %s \"%s\" %d", filename, escapeDigits, timeout)
	resp, err := s.execute(cmd)
	if err != nil {
		return "", err
	}
	if resp.Result == -1 {
		return "", nil // Timeout
	}
	return string(resp.Result), nil
}

// SayNumber says a number
func (s *AgiSession) SayNumber(number int, escapeDigits string) (string, error) {
	cmd := fmt.Sprintf("SAY NUMBER %d \"%s\"", number, escapeDigits)
	resp, err := s.execute(cmd)
	if err != nil {
		return "", err
	}
	if resp.Result == -1 {
		return "", nil // Timeout or hangup
	}
	return string(resp.Result), nil
}

// SayDigits says digits
func (s *AgiSession) SayDigits(digits string, escapeDigits string) (string, error) {
	cmd := fmt.Sprintf("SAY DIGITS %s \"%s\"", digits, escapeDigits)
	resp, err := s.execute(cmd)
	if err != nil {
		return "", err
	}
	if resp.Result == -1 {
		return "", nil
	}
	return string(resp.Result), nil
}

// SayDateTime says a date/time
func (s *AgiSession) SayDateTime(timestamp int64, escapeDigits string, format string, timezone string) (string, error) {
	cmd := fmt.Sprintf("SAY DATETIME %d \"%s\" %s %s", timestamp, escapeDigits, format, timezone)
	resp, err := s.execute(cmd)
	if err != nil {
		return "", err
	}
	if resp.Result == -1 {
		return "", nil
	}
	return string(resp.Result), nil
}

// DatabaseGet gets a value from the Asterisk database
func (s *AgiSession) DatabaseGet(family, key string) (string, error) {
	cmd := fmt.Sprintf("DATABASE GET %s %s", family, key)
	resp, err := s.execute(cmd)
	if err != nil {
		return "", err
	}
	if resp.Result == 0 {
		return "", nil
	}
	return resp.Data, nil
}

// DatabasePut puts a value into the Asterisk database
func (s *AgiSession) DatabasePut(family, key, value string) error {
	cmd := fmt.Sprintf("DATABASE PUT %s %s \"%s\"", family, key, value)
	_, err := s.execute(cmd)
	return err
}

// DatabaseDel deletes a key from the Asterisk database
func (s *AgiSession) DatabaseDel(family, key string) error {
	cmd := fmt.Sprintf("DATABASE DEL %s %s", family, key)
	_, err := s.execute(cmd)
	return err
}

// Verbose sends a message to the Asterisk verbose log
func (s *AgiSession) Verbose(message string, level int) error {
	cmd := fmt.Sprintf("VERBOSE \"%s\" %d", message, level)
	_, err := s.execute(cmd)
	return err
}

// RecordFile records audio to a file
func (s *AgiSession) RecordFile(filename, format, escapeDigits string, timeout, offset, beep int, silence int) (string, error) {
	cmd := fmt.Sprintf("RECORD FILE %s %s \"%s\" %d %d %d %d",
		filename, format, escapeDigits, timeout, offset, beep, silence)
	resp, err := s.execute(cmd)
	if err != nil {
		return "", err
	}
	if resp.Result == -1 {
		return "", nil
	}
	return string(resp.Result), nil
}

// SendText sends text to channels that support it
func (s *AgiSession) SendText(text string) error {
	cmd := fmt.Sprintf("SEND TEXT \"%s\"", text)
	_, err := s.execute(cmd)
	return err
}

// SendImage sends an image to channels that support it
func (s *AgiSession) SendImage(image string) error {
	cmd := fmt.Sprintf("SEND IMAGE %s", image)
	_, err := s.execute(cmd)
	return err
}

// SetMusic enables/disables music on hold
func (s *AgiSession) SetMusic(on bool, class string) error {
	onStr := "OFF"
	if on {
		onStr = "ON"
	}
	cmd := fmt.Sprintf("SET MUSIC %s %s", onStr, class)
	_, err := s.execute(cmd)
	return err
}

// SetCallerID sets the caller ID
func (s *AgiSession) SetCallerID(number string) error {
	cmd := fmt.Sprintf("SET CALLERID %s", number)
	_, err := s.execute(cmd)
	return err
}

// SetContext sets the context for the channel
func (s *AgiSession) SetContext(context string) error {
	cmd := fmt.Sprintf("SET CONTEXT %s", context)
	_, err := s.execute(cmd)
	return err
}

// SetExtension sets the extension for the channel
func (s *AgiSession) SetExtension(extension string) error {
	cmd := fmt.Sprintf("SET EXTENSION %s", extension)
	_, err := s.execute(cmd)
	return err
}

// SetPriority sets the priority for the channel
func (s *AgiSession) SetPriority(priority int) error {
	cmd := fmt.Sprintf("SET PRIORITY %d", priority)
	_, err := s.execute(cmd)
	return err
}

// Noop does nothing (but can be used for debugging)
func (s *AgiSession) Noop() error {
	_, err := s.execute("NOOP")
	return err
}

// ReceiveChar receives a character from channels that support it
func (s *AgiSession) ReceiveChar(timeout int) (string, error) {
	cmd := fmt.Sprintf("RECEIVE CHAR %d", timeout)
	resp, err := s.execute(cmd)
	if err != nil {
		return "", err
	}
	if resp.Result == -1 {
		return "", nil
	}
	return string(resp.Result), nil
}

// ReceiveText receives text from channels that support it
func (s *AgiSession) ReceiveText(timeout int) (string, error) {
	cmd := fmt.Sprintf("RECEIVE TEXT %d", timeout)
	resp, err := s.execute(cmd)
	if err != nil {
		return "", err
	}
	return resp.Data, nil
}

// AGIResultCode represents the possible result codes from AGI commands
type AGIResultCode int

const (
	Success AGIResultCode = iota
	Hangup
	Error
	Timeout
)

// GetResultCode parses an AGI response result into a ResultCode
func GetResultCode(result int) AGIResultCode {
	switch {
	case result == 0:
		return Success
	case result == -1:
		return Hangup
	case result > 0:
		return Success
	default:
		return Error
	}
}
