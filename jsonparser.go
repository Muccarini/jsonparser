package jsonparser

// IMPORT
import (
	"bytes"
	"fmt"
)

// API

func Get(json []byte, fields ...string) (string, error) {
	return get(json, fields...)
}

//FUNC

func get(json []byte, fields ...string) (string, error) {

	// Simple validations
	if len(json) == 0 {
		return "", fmt.Errorf("json is empty")
	}

	if len(fields) == 0 {
		return "", fmt.Errorf("no fields specified")
	}

	pos := 0
	depth := 0
	isValue := false
	fieldFoundIndex := 0

	// Initialize depth based on JSON start 0 or -1 based on whether it starts with {
	for pos < len(json) && isWhitespace(json[pos]) {
		pos++
	}

	if pos < len(json) && json[pos] == '{' {
		depth = 0
		pos++
	} else {
		depth = -1
	}

	for pos < len(json) {

		switch json[pos] {
		case '"':
			// "candidate_field" must not:
			// 1) be \"escaped\"
			// 2) be longer/shorter than the field we are looking for.
			// 3) have different depth that the field we are looking for.
			if depth == fieldFoundIndex &&
				json[pos-1] != '\\' &&
				pos+len(fields[fieldFoundIndex]) <= len(json)-1 && json[pos+len(fields[fieldFoundIndex])] == '"' {
					// yeyy our "candidate_field" is in the right position and has the right depth
					if bytes.Equal(json[pos+1:pos+len(fields[fieldFoundIndex])], []byte(fields[fieldFoundIndex])) {
						// we found a match for the current field in the path
						fieldFoundIndex++
						if fieldFoundIndex == len(fields) {
							// We found the full path, now extract the value
							getNextToken(json, pos+len(fields[fieldFoundIndex-1])+1)
							colonPos := findColon(json, pos+len(fields[fieldFoundIndex-1])+1)
						}
					}
				}
			}
		case '{':
			depth++
			isValue = false
			pos++
		case '[':
			// array
			depth++
			isValue = false
			pos++
		case '}':
			depth--
			isValue = false
			pos++
		case ',':
			isValue = false
		case ':':
			isValue = true
		}
	}
	return "", fmt.Errorf("field path not found: %v", fields)
}

// Helper functions

func getNextToken(json []byte, pos int) (byte, int, error) {
	for pos < len(json) {
		if !isWhitespace(json[pos]) {
			return json[pos], pos, nil
		}
		pos++
	}
	return 0, -1, fmt.Errorf("no more tokens")
}

func isWhitespace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

// checked
func findFieldEnd(json []byte, start int) int {
	for i := start; i < len(json); i++ {
		if json[i] == '"' {
			return i
		}
	}
	return -1
}

// checked
func findColon(json []byte, start int) int {
	for i := start; i < len(json); i++ {
		if json[i] == ':' {
			return i
		}
		if !isWhitespace(json[i]) {
			return -1 // Found non-whitespace before colon
		}
	}
	return -1
}

func extractValue(json []byte, start int) ([]byte, error) {
	// Skip whitespace
	pos := start
	for pos < len(json) && isWhitespace(json[pos]) {
		pos++
	}

	if pos >= len(json) {
		return nil, fmt.Errorf("no value found")
	}

	char := json[pos]

	// Handle different value types
	switch char {
	case '"':
		// String value
		return extractString(json, pos)
	case '{':
		// Object value
		return extractObject(json, pos)
	case '[':
		// Array value
		return extractArray(json, pos)
	case 't', 'f':
		// Boolean value
		return extractBoolean(json, pos)
	case 'n':
		// Null value
		return extractNull(json, pos)
	default:
		// Number value
		if (char >= '0' && char <= '9') || char == '-' {
			return extractNumber(json, pos)
		}
	}

	return nil, fmt.Errorf("invalid value at position %d", pos)
}

func extractString(json []byte, start int) ([]byte, error) {
	pos := start + 1 // Skip opening quote
	escaped := false

	for pos < len(json) {
		if escaped {
			escaped = false
		} else if json[pos] == '\\' {
			escaped = true
		} else if json[pos] == '"' {
			// Found closing quote, return the correct slice
			return json[start : pos+1], nil
		}
		pos++
	}

	return nil, fmt.Errorf("unterminated string")
}

func extractNumber(json []byte, start int) ([]byte, error) {
	pos := start

	// Skip optional minus
	if pos < len(json) && json[pos] == '-' {
		pos++
	}

	// Must have at least one digit
	if pos >= len(json) || (json[pos] < '0' || json[pos] > '9') {
		return nil, fmt.Errorf("invalid number")
	}

	// Skip digits
	for pos < len(json) && json[pos] >= '0' && json[pos] <= '9' {
		pos++
	}

	// Optional decimal part
	if pos < len(json) && json[pos] == '.' {
		pos++
		for pos < len(json) && json[pos] >= '0' && json[pos] <= '9' {
			pos++
		}
	}

	// Optional exponent
	if pos < len(json) && (json[pos] == 'e' || json[pos] == 'E') {
		pos++
		if pos < len(json) && (json[pos] == '+' || json[pos] == '-') {
			pos++
		}
		for pos < len(json) && json[pos] >= '0' && json[pos] <= '9' {
			pos++
		}
	}

	return json[start:pos], nil
}

func extractBoolean(json []byte, start int) ([]byte, error) {
	if start+4 <= len(json) && bytes.Equal(json[start:start+4], []byte("true")) {
		return json[start : start+4], nil
	}
	if start+5 <= len(json) && bytes.Equal(json[start:start+5], []byte("false")) {
		return json[start : start+5], nil
	}
	return nil, fmt.Errorf("invalid boolean")
}

func extractNull(json []byte, start int) ([]byte, error) {
	if start+4 <= len(json) && bytes.Equal(json[start:start+4], []byte("null")) {
		return json[start : start+4], nil
	}
	return nil, fmt.Errorf("invalid null")
}

func parseObject(json []byte, pos int, depth int) ([]byte, error) {
	inString := false
	escaped := false

	for pos < len(json) {
		char := json[pos]

		if inString {
			if escaped {
				escaped = false
			} else if char == '\\' {
				escaped = true
			} else if char == '"' {
				inString = false
			}
		} else {
			if char == '"' {
				inString = true
			} else if char == '{' {
				depth++
			} else if char == '}' {
				depth--
				if depth == 0 {
					return json[start : pos+1], nil
				}
			}
		}
		pos++
	}

	return nil, fmt.Errorf("unterminated object")
}

func extractArray(json []byte, start int) ([]byte, error) {
	depth := 0
	pos := start
	inString := false
	escaped := false

	for pos < len(json) {
		char := json[pos]

		if inString {
			if escaped {
				escaped = false
			} else if char == '\\' {
				escaped = true
			} else if char == '"' {
				inString = false
			}
		} else {
			if char == '"' {
				inString = true
			} else if char == '[' {
				depth++
			} else if char == ']' {
				depth--
				if depth == 0 {
					return json[start : pos+1], nil
				}
			}
		}
		pos++
	}

	return nil, fmt.Errorf("unterminated array")
}
