package agi

import (
	"fmt"
	"strconv"
	"strings"
)

// EscapeString escapes a string for use in AGI commands
func EscapeString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}

// UnescapeString unescapes a string from AGI responses
func UnescapeString(s string) string {
	s = strings.ReplaceAll(s, "\\\"", "\"")
	s = strings.ReplaceAll(s, "\\\\", "\\")
	return s
}

// ParseAGIResult parses an AGI result string
func ParseAGIResult(s string) (int, string, error) {
	if !strings.HasPrefix(s, "200") {
		return 0, "", fmt.Errorf("invalid AGI response: %s", s)
	}

	s = strings.TrimPrefix(s, "200 ")
	if !strings.HasPrefix(s, "result=") {
		return 0, "", fmt.Errorf("invalid AGI response format: %s", s)
	}

	s = strings.TrimPrefix(s, "result=")
	parts := strings.SplitN(s, " ", 2)

	result, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", fmt.Errorf("failed to parse result code: %v", err)
	}

	var data string
	if len(parts) > 1 {
		data = strings.TrimSpace(parts[1])
		if strings.HasPrefix(data, "(") && strings.HasSuffix(data, ")") {
			data = data[1 : len(data)-1]
		}
	}

	return result, data, nil
}

// ParseAGIEnv parses AGI environment variables
func ParseAGIEnv(input string) map[string]string {
	env := make(map[string]string)
	lines := strings.Split(input, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "agi_") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				env[key] = value
			}
		}
	}

	return env
}

// FormatDateTime formats a date/time string for Asterisk
func FormatDateTime(format string) string {
	// Common format substitutions
	replacements := map[string]string{
		"YYYY": "2006",
		"YY":   "06",
		"MM":   "01",
		"DD":   "02",
		"hh":   "15",
		"mm":   "04",
		"ss":   "05",
	}

	result := format
	for from, to := range replacements {
		result = strings.ReplaceAll(result, from, to)
	}
	return result
}

// SplitCommand splits an AGI command into its components
func SplitCommand(cmd string) []string {
	var parts []string
	var current strings.Builder
	inQuotes := false
	escaped := false

	for _, c := range cmd {
		if escaped {
			current.WriteRune(c)
			escaped = false
			continue
		}

		switch c {
		case '\\':
			escaped = true
		case '"':
			inQuotes = !inQuotes
		case ' ':
			if inQuotes {
				current.WriteRune(c)
			} else if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(c)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

// JoinCommand joins command parts with proper escaping
func JoinCommand(parts []string) string {
	escaped := make([]string, len(parts))
	for i, part := range parts {
		if strings.ContainsAny(part, " \"") {
			escaped[i] = fmt.Sprintf("\"%s\"", EscapeString(part))
		} else {
			escaped[i] = part
		}
	}
	return strings.Join(escaped, " ")
}
