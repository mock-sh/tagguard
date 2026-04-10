// Package tags provides utilities for parsing Go struct tag strings.
package tags

import "strings"

// Tag represents a single key-value pair in a struct tag.
// For example, in `json:"name,omitempty"`:
//
//	Key     = "json"
//	Value   = "name"
//	Options = ["omitempty"]
type Tag struct {
	Key     string
	Value   string
	Options []string
	Raw     string // full raw value e.g. "name,omitempty"
}

// Parse parses a raw struct tag string (without backticks) into individual tags.
// It handles the standard Go struct tag format: key:"value,opt1,opt2"
func Parse(raw string) []Tag {
	var result []Tag
	s := raw

	for s != "" {
		// Skip leading spaces
		s = strings.TrimLeft(s, " \t")
		if s == "" {
			break
		}

		// Read key (everything before the colon)
		colonIdx := strings.Index(s, ":")
		if colonIdx < 0 {
			break
		}
		key := s[:colonIdx]
		s = s[colonIdx+1:]

		// Value must start with a double-quote
		if len(s) == 0 || s[0] != '"' {
			break
		}

		// Find closing quote (handles escaped quotes)
		value, rest, ok := readQuoted(s)
		if !ok {
			break
		}
		s = rest

		parts := strings.SplitN(value, ",", -1)
		mainValue := parts[0]
		var options []string
		if len(parts) > 1 {
			options = parts[1:]
		}

		result = append(result, Tag{
			Key:     key,
			Value:   mainValue,
			Options: options,
			Raw:     value,
		})
	}

	return result
}

// readQuoted reads a double-quoted string from the start of s.
// Returns the unquoted string content, the remainder of s, and whether it succeeded.
func readQuoted(s string) (string, string, bool) {
	if len(s) == 0 || s[0] != '"' {
		return "", s, false
	}
	i := 1
	for i < len(s) {
		if s[i] == '\\' {
			i += 2
			continue
		}
		if s[i] == '"' {
			return s[1:i], s[i+1:], true
		}
		i++
	}
	return "", s, false
}
