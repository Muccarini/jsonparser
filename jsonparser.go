package jsonparser

import (
	"bytes"
	"fmt"
	"strconv"
)

var (
	ERROR_INVALID_JSON         = fmt.Errorf("invalid JSON")
	ERROR_FIELD_NOT_FOUND      = fmt.Errorf("field not found")
	ERROR_ARGUMENTS            = fmt.Errorf("invalid arguments")
	ERROR_INVALID_INTEGER      = fmt.Errorf("invalid integer")
	ERROR_INVALID_FLOAT        = fmt.Errorf("invalid float")
	ERROR_INVALID_BOOLEAN      = fmt.Errorf("invalid boolean")
	ERROR_INVALID_STRING       = fmt.Errorf("invalid string")
	ERROR_INVALID_NULL         = fmt.Errorf("invalid null")
	ERROR_COLON_NOT_FOUND      = fmt.Errorf("no colon found")
	ERROR_NEXT_TOKEN_NOT_FOUND = fmt.Errorf("no more tokens found")
)

type irange struct {
	start int
	end   int
}

// API

func GetString(json []byte, fields ...string) (string, error) {
	res, err := get(json, fields...)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func GetBool(json []byte, fields ...string) (bool, error) {
	res, err := get(json, fields...)
	if err != nil {
		return false, err
	}

	switch string(res) {
	case "true":
		return true, nil
	case "false":
		return false, nil
	}

	return false, ERROR_INVALID_BOOLEAN
}

func GetInt(json []byte, fields ...string) (int, error) {
	res, err := get(json, fields...)
	if err != nil {
		return -1, err
	}

	resInt, err := strconv.Atoi(string(res))
	if err != nil {
		return -1, err
	}

	return resInt, nil
}

func GetInt64(json []byte, fields ...string) (int64, error) {
	res, err := get(json, fields...)
	if err != nil {
		return -1, err
	}

	resInt64, err := strconv.ParseInt(string(res), 10, 64)
	if err != nil {
		return -1, err
	}

	return resInt64, nil
}

func GetFloat64(json []byte, bitSize int, fields ...string) (float64, error) {
	res, err := get(json, fields...)
	if err != nil {
		return -1, err
	}

	resFloat, err := strconv.ParseFloat(string(res), bitSize)
	if err != nil {
		return -1, err
	}

	return resFloat, nil
}

// INTERNAL

func get(json []byte, fields ...string) ([]byte, error) {
	// Simple validations
	if len(json) == 0 {
		return nil, ERROR_INVALID_JSON
	}

	if len(fields) == 0 {
		return nil, ERROR_ARGUMENTS
	}

	pos := 0
	depth := 0
	fieldFoundIndex := 0
	isValue := false
	isArrayValue := false

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
			candidate := irange{start: pos + 1, end: pos + 1 + len(fields[fieldFoundIndex])}

			if depth == fieldFoundIndex &&
				json[pos-1] != '\\' &&
				!isValue &&
				candidate.end+1 <= len(json) &&
				json[candidate.end] == '"' {

				if bytes.Equal(json[candidate.start:candidate.end], []byte(fields[fieldFoundIndex])) {
					fieldFoundIndex++
					if fieldFoundIndex == len(fields) {
						// We found the full path, now extract the value
						colonPos, err := nextColon(json, candidate.end+1)
						if err != nil {
							return nil, err
						}

						slice, err := extractValue(json, colonPos+1, depth)
						if err != nil {
							return nil, err
						}

						return slice, nil
					}
				}
			}
			pos++
		case '{':
			if json[pos-1] != '\\' {
				isValue = false
				depth++
			}
			pos++
		case '}':
			if json[pos-1] != '\\' {
				isValue = false
				depth--
			}
			pos++
		case '[':
			if json[pos-1] != '\\' {
				isArrayValue = true
				isValue = false
			}
			pos++
		case ']':
			if json[pos-1] != '\\' {
				isArrayValue = false
				isValue = false
			}
			pos++
		case ':':
			isValue = true
			pos++
		case ',':
			isValue = isArrayValue
			pos++
		default:
			pos++
		}
	}

	return nil, ERROR_FIELD_NOT_FOUND
}

func extractValue(json []byte, pos int, depth int) ([]byte, error) {
	for pos < len(json) && isWhitespace(json[pos]) {
		pos++
	}

	switch json[pos] {
	case 't', 'f':
		slice, err := extractBoolean(json, pos)
		if err != nil {
			return nil, err
		}
		return slice, nil

	case 'n':
		slice, err := extractNull(json, pos)
		if err != nil {
			return nil, err
		}
		return slice, nil

	case '{':
		slice, err := extractObject(json, pos, depth+1)
		if err != nil {
			return nil, err
		}
		return slice, nil

	case '[':
		slice, err := extractArray(json, pos)
		if err != nil {
			return nil, err
		}
		return slice, nil

	case '"':
		slice, err := extractString(json, pos)
		if err != nil {
			return nil, err
		}
		return slice, nil

	default:
		slice, err := extractNumber(json, pos)
		if err != nil {
			return nil, err
		}
		return slice, nil
	}
}

// return the first position of a non-whitespace character starting from pos + 1
func nextToken(json []byte, pos int) (int, error) {
	for pos < len(json) {
		pos++
		if !isWhitespace(json[pos]) {
			return pos, nil
		}
	}
	return -1, ERROR_INVALID_JSON
}

func isWhitespace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

// return the position of the next colon starting from pos
func nextColon(json []byte, pos int) (int, error) {
	for pos < len(json) {
		if json[pos] == ':' {
			return pos, nil
		}
		pos++
	}
	return -1, fmt.Errorf("no colon found")
}

func findArrayElementPos(json []byte, pos int, field int, depth int) (int, error) {
	for pos < len(json) && isWhitespace(json[pos]) {
		pos++
	}

	isValue := false
	index := 0
	arrayDepth := depth

	for pos < len(json) {
		switch json[pos] {
		case ',':
			if !isValue && arrayDepth == depth {
				index++
				if field == index {
					for pos < len(json) && isWhitespace(json[pos]) {
						pos++
					}
				}
			}
			return pos, nil
		case '"':
			if json[pos-1] != '\\' {
				isValue = !isValue
			}
			pos++
		case '{':
			if json[pos-1] != '\\' {
				depth++
			}
			pos++
		case '}':
			if json[pos-1] != '\\' {
				depth--
			}
			pos++
		case '[':
			if json[pos-1] != '\\' {
				depth++
			}
			pos++
		case ']':
			if json[pos-1] != '\\' {
				depth--
			}
			pos++
		case ':':
			pos++
		default:
			pos++
		}
	}

	return -1, ERROR_FIELD_NOT_FOUND
}

func findFieldSlice(json []byte, pos int, field string, fieldIndex int, depth int) ([]byte, error) {
	if len(json) == 0 {
		return nil, ERROR_INVALID_JSON
	}

	for pos < len(json) && isWhitespace(json[pos]) {
		pos++
	}

	isValue := false
	isArrayValue := false

	for pos < len(json) {
		switch json[pos] {
		case '"':
			candidate := irange{start: pos + 1, end: pos + 1 + len(field)}
			if depth == fieldIndex &&
				json[pos-1] != '\\' &&
				!isValue &&
				candidate.end+1 <= len(json) &&
				json[candidate.end] == '"' {

				if bytes.Equal(json[candidate.start:candidate.end], []byte(field)) {
					return json[candidate.start:candidate.end], nil
				}
			}
			pos++
		case '{':
			if json[pos-1] != '\\' {
				isValue = false
				depth++
			}
			pos++
		case '}':
			if json[pos-1] != '\\' {
				isValue = false
				depth--
			}
			pos++
		case '[':
			if json[pos-1] != '\\' {
				isArrayValue = true
				isValue = false
			}
			pos++
		case ']':
			if json[pos-1] != '\\' {
				isArrayValue = false
				isValue = false
			}
			pos++
		case ':':
			isValue = true
			pos++
		case ',':
			isValue = isArrayValue
			pos++
		default:
			pos++
		}
	}

	return nil, ERROR_FIELD_NOT_FOUND
}

// checked
func extractString(json []byte, pos int) ([]byte, error) {
	// Skip opening quote
	var start int
	if json[pos] == '"' {
		start = pos
		pos++
	} else {
		start = pos
	}

	// Find closing quote
	for pos < len(json) {
		if json[pos] == '"' {
			if json[pos-1] != '\\' {
				// Returning the value without the quotes.
				return json[start+1 : pos], nil
			}
		}
		pos++
	}

	return nil, ERROR_INVALID_JSON
}

func extractNumber(json []byte, pos int) ([]byte, error) {
	start := pos

	for pos < len(json) {
		switch json[pos] {
		case '-', '+', '.', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'e', 'E':
			pos++
			continue
		default:
			return json[start:pos], nil
		}
	}

	return nil, ERROR_INVALID_JSON
}

func extractBoolean(json []byte, pos int) ([]byte, error) {
	if pos+4 <= len(json) && bytes.Equal(json[pos:pos+4], []byte("true")) {
		return json[pos : pos+4], nil
	}
	if pos+5 <= len(json) && bytes.Equal(json[pos:pos+5], []byte("false")) {
		return json[pos : pos+5], nil
	}
	return nil, fmt.Errorf("invalid boolean")
}

func extractNull(json []byte, start int) ([]byte, error) {
	if start+4 <= len(json) && bytes.Equal(json[start:start+4], []byte("null")) {
		return json[start : start+4], nil
	}
	return nil, ERROR_INVALID_NULL
}

func extractObject(json []byte, pos int, depth int) ([]byte, error) {
	isValue := false
	escaped := false
	isArray := false

	switch json[pos] {
	case '{':
		if json[pos-1] != '\\' && !isValue {
			// Start of an object
			depth++
		}
	case '}':
		if json[pos-1] != '\\' && !isValue {
			depth--
		}
		pos++
	case '[':
		if json[pos-1] != '\\' && !isValue {
			isArray = true
		}
		pos++
	case ']':
		if json[pos-1] != '\\' && !isValue {
			isArray = false
		}
		pos++
	case ':':
		if json[pos-1] != '\\' {
			isValue = true
		}
		pos++
	case ',':
		isValue = isArray
		pos++
	}

	for pos < len(json) {

		if isValue {
			if escaped {
				escaped = false
			} else if json[pos] == '\\' {
				escaped = true
			} else if json[pos] == '"' {
				isValue = false
			}
		} else {
			if json[pos] == '"' {
				isValue = true
			} else if json[pos] == '{' {
				depth++
			} else if json[pos] == '}' {
				depth--
				if depth == 0 {
					return json[pos : pos+1], nil
				}
			}
		}
		pos++
	}

	return nil, ERROR_INVALID_JSON
}

func extractArray(json []byte, start int) ([]byte, error) {
	count := 0
	pos := start
	isValue := false
	escaped := false

	for pos < len(json) {
		char := json[pos]

		if isValue {
			if escaped {
				escaped = false
			} else if char == '\\' {
				escaped = true
			} else if char == '"' {
				isValue = false
			}
		} else {
			if char == '"' {
				isValue = true
			} else if char == '[' {
				count++
			} else if char == ']' {
				count--
				if count == 0 {
					return json[start : pos+1], nil
				}
			}
		}
		pos++
	}

	return nil, fmt.Errorf("unterminated array")
}
