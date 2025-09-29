package jsonparser

import (
	"bytes"
	"fmt"
	"strconv"
)

var (
	ERROR_INVALID_JSON       = fmt.Errorf("invalid JSON")
	ERROR_FIELD_NOT_FOUND    = fmt.Errorf("field not found")
	ERROR_ARGUMENTS          = fmt.Errorf("invalid arguments")
	ERROR_COLON_NOT_FOUND    = fmt.Errorf("no colon found")
	ERROR_INVALID_INTEGER    = fmt.Errorf("invalid integer")
	ERROR_INVALID_FLOAT      = fmt.Errorf("invalid float")
	ERROR_INVALID_BOOLEAN    = fmt.Errorf("invalid boolean")
	ERROR_INVALID_STRING     = fmt.Errorf("invalid string")
	ERROR_INVALID_NULL       = fmt.Errorf("invalid null")
	ERROR_UNTERMINATED_ARRAY = fmt.Errorf("unterminated array")
)

type irange struct {
	start int
	end   int
}

// API

func GetString(json []byte, fields ...string) (string, error) {
	if len(json) == 0 {
		return "", ERROR_INVALID_JSON
	}

	if len(fields) == 0 {
		return "", ERROR_ARGUMENTS
	}

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

	slice := json

	for _, field := range fields {
		if isNumericField(field) {
			intField, err := strconv.Atoi(field)
			if err != nil {
				return nil, err
			}
			// field is an integer, search on this depth level the nth element of the array
			// find the position next to comma ', {value}' and extract the value
			valuePos, err := findArrayValue(slice, 0, intField)
			if err != nil {
				return nil, err
			}

			slice = slice[valuePos:]
		} else {
			valuePos, err := findFieldValuePos(slice, 0, field)
			if err != nil {
				return nil, err
			}

			slice = slice[valuePos:] // starts from the value
		}
	}

	value, err := extractValue(slice, 0)
	if err != nil {
		return nil, ERROR_FIELD_NOT_FOUND
	}

	return value, nil
}

func isWhitespace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

// isNumericField checks if a field string represents a number without allocating
func isNumericField(s string) bool {
	if len(s) == 0 {
		return false
	}

	for _, c := range []byte(s) {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// return the position of the next colon starting from pos
func nextColon(json []byte, pos int) (int, error) {
	for pos < len(json) {
		if json[pos] == ':' {
			return pos, nil
		}
		pos++
	}
	return -1, ERROR_COLON_NOT_FOUND
}

func findArrayValue(json []byte, pos int, elementIndex int) (int, error) {
	if len(json) == 0 {
		return -1, ERROR_INVALID_JSON
	}

	if json[pos] != '[' {
		return -1, ERROR_INVALID_JSON
	}

	pos++
	isValue := false
	index := 0

	if elementIndex == 0 {
		for pos < len(json) && isWhitespace(json[pos]) {
			pos++
		}
		return pos, nil
	}

	for pos < len(json) {
		switch json[pos] {
		case ',':
			if !isValue {
				index++
				if elementIndex == index {
					pos++ //skip comma and whitespace
					for pos < len(json) && isWhitespace(json[pos]) {
						pos++
					}
					return pos, nil
				}
			}
			pos++
		case '"':
			if json[pos-1] != '\\' {
				isValue = !isValue
			}
			pos++
		case '{':
			if json[pos-1] != '\\' {
				posRes, err := skipObject(json, pos)
				if err != nil {
					return -1, err
				}
				pos = posRes
				continue
			}
			pos++
		case '[':
			if json[pos-1] != '\\' {
				posRes, err := skipMatrix(json, pos)
				if err != nil {
					return -1, err
				}
				pos = posRes
				continue
			}
			pos++
		default:
			pos++
		}
	}

	return -1, ERROR_FIELD_NOT_FOUND
}

// findFieldValuePos returns the position of the colon after the field name at the specified depth
func findFieldValuePos(json []byte, pos int, field string) (int, error) {
	if len(json) == 0 {
		return -1, ERROR_INVALID_JSON
	}

	for pos < len(json) && isWhitespace(json[pos]) {
		pos++
	}

	//check if is an object
	if json[pos] == '{' {
		pos++
	}

	// We are looking for the field at the same relative depth we called this function
	depth := 0
	isValue := false

	for pos < len(json) {
		switch json[pos] {
		case '"':
			if depth == 0 &&
				!isValue &&
				pos > 0 && json[pos-1] != '\\' {

				candidate := irange{start: pos + 1, end: pos + 1 + len(field)}

				if json[candidate.end] == '"' &&
					bytes.Equal(json[candidate.start:candidate.end], []byte(field)) {

					pos, err := nextColon(json, candidate.end)
					if err != nil {
						return -1, err
					}

					pos++ //skip colon and whitespace
					for pos < len(json) && isWhitespace(json[pos]) {
						pos++
					}

					return pos, nil
				}
			}
			pos++
		case '{':
			if pos > 0 && json[pos-1] != '\\' {
				isValue = false
				posRes, err := skipObject(json, pos)
				if err != nil {
					return -1, err
				}
				pos = posRes
				continue
			}
			pos++
		case '[':
			if pos > 0 && json[pos-1] != '\\' {
				isValue = false
				posRes, err := skipMatrix(json, pos)
				if err != nil {
					return -1, err
				}
				pos = posRes
				continue
			}
			pos++
		case ':':
			if pos > 0 && json[pos-1] != '\\' {
				isValue = true
			}
			pos++
		case ',':
			if pos > 0 && json[pos-1] != '\\' {
				isValue = false
			}
			pos++
		default:
			pos++
		}
	}

	return -1, ERROR_FIELD_NOT_FOUND
}

func skipObject(json []byte, pos int) (int, error) {
	if json[pos] != '{' {
		return -1, ERROR_INVALID_JSON
	}

	depth := 1
	pos++

	for pos < len(json) {
		switch json[pos] {
		case '{':
			if json[pos-1] != '\\' {
				depth++
			}
		case '}':
			if json[pos-1] != '\\' {
				depth--
				if depth == 0 {
					return pos + 1, nil
				}
			}
		}
		pos++
	}

	return -1, ERROR_INVALID_JSON
}

func skipMatrix(json []byte, pos int) (int, error) {
	if json[pos] != '[' {
		return -1, ERROR_INVALID_JSON
	}

	count := 1
	pos++

	for pos < len(json) {
		switch json[pos] {
		case '[':
			if json[pos-1] != '\\' {
				count++
			}
		case ']':
			if json[pos-1] != '\\' {
				count--
				if count == 0 {
					return pos + 1, nil
				}
			}
		}
		pos++
	}

	return -1, ERROR_INVALID_JSON
}

func extractValue(json []byte, pos int) ([]byte, error) {
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
		slice, err := extractObject(json, pos)
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
			if pos > 0 && json[pos-1] != '\\' {
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
	return nil, ERROR_INVALID_BOOLEAN
}

func extractNull(json []byte, start int) ([]byte, error) {
	if start+4 <= len(json) && bytes.Equal(json[start:start+4], []byte("null")) {
		return json[start : start+4], nil
	}
	return nil, ERROR_INVALID_NULL
}

func extractObject(json []byte, pos int) ([]byte, error) {

	depth := 0
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

	return nil, ERROR_UNTERMINATED_ARRAY
}
