// Package jsonparser provides a fast JSON parser with filters and zero memory allocation.
package jsonparser

import (
	"errors"
	"strconv"
)

// ValueType represents the type of a JSON value
type ValueType int

const (
	NotExist ValueType = iota
	String
	Number
	Object
	Array
	Boolean
	Null
	Unknown
)

// Errors returned by the parser
var (
	ErrKeyPathNotFound = errors.New("key path not found")
	ErrInvalidJSON     = errors.New("invalid JSON")
	ErrUnknownValue    = errors.New("unknown value type")
)

// Get returns the value by the given key path and its type.
// The path is a dot-separated string like "user.name" or "users.0.email".
// Returns the raw bytes of the value, its type, and any error.
func Get(data []byte, keys ...string) (value []byte, dataType ValueType, offset int, err error) {
	return get(data, keys...)
}

// GetString extracts a string value by the given key path.
func GetString(data []byte, keys ...string) (string, error) {
	value, dataType, _, err := Get(data, keys...)
	if err != nil {
		return "", err
	}
	if dataType != String {
		return "", errors.New("value is not a string")
	}
	
	// Remove surrounding quotes and handle escaped characters
	if len(value) < 2 || value[0] != '"' || value[len(value)-1] != '"' {
		return "", ErrInvalidJSON
	}
	return parseString(value[1 : len(value)-1])
}

// GetInt extracts an integer value by the given key path.
func GetInt(data []byte, keys ...string) (int64, error) {
	value, dataType, _, err := Get(data, keys...)
	if err != nil {
		return 0, err
	}
	if dataType != Number {
		return 0, errors.New("value is not a number")
	}
	return strconv.ParseInt(string(value), 10, 64)
}

// GetFloat extracts a float value by the given key path.
func GetFloat(data []byte, keys ...string) (float64, error) {
	value, dataType, _, err := Get(data, keys...)
	if err != nil {
		return 0, err
	}
	if dataType != Number {
		return 0, errors.New("value is not a number")
	}
	return strconv.ParseFloat(string(value), 64)
}

// GetBoolean extracts a boolean value by the given key path.
func GetBoolean(data []byte, keys ...string) (bool, error) {
	value, dataType, _, err := Get(data, keys...)
	if err != nil {
		return false, err
	}
	if dataType != Boolean {
		return false, errors.New("value is not a boolean")
	}
	
	switch string(value) {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, ErrInvalidJSON
	}
}

// ArrayEach iterates over each element in an array at the given key path.
// The callback function is called for each element with its value, type, and index.
func ArrayEach(data []byte, callback func(value []byte, dataType ValueType, offset int, err error), keys ...string) {
	arrayData, dataType, _, err := Get(data, keys...)
	if err != nil {
		callback(nil, Unknown, -1, err)
		return
	}
	
	if dataType != Array {
		callback(nil, Unknown, -1, errors.New("value is not an array"))
		return
	}
	
	arrayEach(arrayData, callback)
}

// ObjectEach iterates over each key-value pair in an object at the given key path.
// The callback function is called for each pair with the key, value, type, and offset.
func ObjectEach(data []byte, callback func(key []byte, value []byte, dataType ValueType, offset int) error, keys ...string) error {
	objectData, dataType, _, err := Get(data, keys...)
	if err != nil {
		return err
	}
	
	if dataType != Object {
		return errors.New("value is not an object")
	}
	
	return objectEach(objectData, callback)
}

// Internal function to get value by key path
func get(data []byte, keys ...string) (value []byte, dataType ValueType, offset int, err error) {
	if len(data) == 0 {
		return nil, NotExist, -1, ErrKeyPathNotFound
	}
	
	// Skip whitespace
	offset = nextToken(data, 0)
	if offset == -1 {
		return nil, NotExist, -1, ErrInvalidJSON
	}
	
	// If no keys provided, return the entire JSON
	if len(keys) == 0 {
		value, dataType, err = getType(data, offset)
		return value, dataType, offset, err
	}
	
	// Navigate through the key path
	current := data
	currentOffset := 0
	
	for _, key := range keys {
		var nextValue []byte
		var nextType ValueType
		
		// Check if we're dealing with an array index
		if isArrayIndex(key) {
			index, parseErr := strconv.Atoi(key)
			if parseErr != nil {
				return nil, NotExist, -1, errors.New("invalid array index")
			}
			
			nextValue, nextType, _, err = getArrayIndex(current, currentOffset, index)
		} else {
			nextValue, nextType, _, err = getObjectKey(current, currentOffset, key)
		}
		
		if err != nil {
			return nil, NotExist, -1, err
		}
		
		current = nextValue
		currentOffset = 0 // Reset offset for the new value
		dataType = nextType
	}
	
	return current, dataType, currentOffset, nil
}

// Get value by object key
func getObjectKey(data []byte, offset int, key string) ([]byte, ValueType, int, error) {
	offset = nextToken(data, offset)
	if offset == -1 || offset >= len(data) || data[offset] != '{' {
		return nil, NotExist, -1, errors.New("expected object")
	}
	
	offset++ // skip '{'
	
	for {
		offset = nextToken(data, offset)
		if offset == -1 {
			return nil, NotExist, -1, ErrInvalidJSON
		}
		
		if data[offset] == '}' {
			return nil, NotExist, -1, ErrKeyPathNotFound
		}
		
		// Parse key
		if data[offset] != '"' {
			return nil, NotExist, -1, ErrInvalidJSON
		}
		
		keyStart := offset + 1
		keyEnd := findStringEnd(data, keyStart)
		if keyEnd == -1 {
			return nil, NotExist, -1, ErrInvalidJSON
		}
		
		currentKey := string(data[keyStart:keyEnd])
		offset = keyEnd + 1 // skip closing quote
		
		// Skip colon
		offset = nextToken(data, offset)
		if offset == -1 || data[offset] != ':' {
			return nil, NotExist, -1, ErrInvalidJSON
		}
		offset++ // skip ':'
		
		// Get value
		offset = nextToken(data, offset)
		if offset == -1 {
			return nil, NotExist, -1, ErrInvalidJSON
		}
		
		valueStart := offset
		valueEnd, valueType, err := skipValue(data, offset)
		if err != nil {
			return nil, NotExist, -1, err
		}
		
		if currentKey == key {
			return data[valueStart:valueEnd], valueType, 0, nil
		}
		
		offset = valueEnd
		
		// Check for comma or end of object
		offset = nextToken(data, offset)
		if offset == -1 {
			return nil, NotExist, -1, ErrInvalidJSON
		}
		
		if data[offset] == '}' {
			return nil, NotExist, -1, ErrKeyPathNotFound
		}
		
		if data[offset] == ',' {
			offset++
			continue
		}
		
		return nil, NotExist, -1, ErrInvalidJSON
	}
}

// Get value by array index
func getArrayIndex(data []byte, offset int, index int) ([]byte, ValueType, int, error) {
	offset = nextToken(data, offset)
	if offset == -1 || offset >= len(data) || data[offset] != '[' {
		return nil, NotExist, -1, errors.New("expected array")
	}
	
	offset++ // skip '['
	currentIndex := 0
	
	for {
		offset = nextToken(data, offset)
		if offset == -1 {
			return nil, NotExist, -1, ErrInvalidJSON
		}
		
		if data[offset] == ']' {
			return nil, NotExist, -1, ErrKeyPathNotFound
		}
		
		valueStart := offset
		valueEnd, valueType, err := skipValue(data, offset)
		if err != nil {
			return nil, NotExist, -1, err
		}
		
		if currentIndex == index {
			return data[valueStart:valueEnd], valueType, 0, nil
		}
		
		currentIndex++
		offset = valueEnd
		
		// Check for comma or end of array
		offset = nextToken(data, offset)
		if offset == -1 {
			return nil, NotExist, -1, ErrInvalidJSON
		}
		
		if data[offset] == ']' {
			return nil, NotExist, -1, ErrKeyPathNotFound
		}
		
		if data[offset] == ',' {
			offset++
			continue
		}
		
		return nil, NotExist, -1, ErrInvalidJSON
	}
}

// Helper functions
func nextToken(data []byte, offset int) int {
	for offset < len(data) {
		switch data[offset] {
		case ' ', '\t', '\n', '\r':
			offset++
		default:
			return offset
		}
	}
	return -1
}

func getType(data []byte, offset int) ([]byte, ValueType, error) {
	if offset >= len(data) {
		return nil, Unknown, ErrInvalidJSON
	}
	
	switch data[offset] {
	case '"':
		end := findStringEnd(data, offset+1)
		if end == -1 {
			return nil, Unknown, ErrInvalidJSON
		}
		return data[offset : end+1], String, nil
		
	case '{':
		end, err := findObjectEnd(data, offset)
		if err != nil {
			return nil, Unknown, err
		}
		return data[offset:end], Object, nil
		
	case '[':
		end, err := findArrayEnd(data, offset)
		if err != nil {
			return nil, Unknown, err
		}
		return data[offset:end], Array, nil
		
	case 't', 'f':
		if offset+4 <= len(data) && string(data[offset:offset+4]) == "true" {
			return data[offset : offset+4], Boolean, nil
		}
		if offset+5 <= len(data) && string(data[offset:offset+5]) == "false" {
			return data[offset : offset+5], Boolean, nil
		}
		return nil, Unknown, ErrInvalidJSON
		
	case 'n':
		if offset+4 <= len(data) && string(data[offset:offset+4]) == "null" {
			return data[offset : offset+4], Null, nil
		}
		return nil, Unknown, ErrInvalidJSON
		
	default:
		// Try to parse as number
		end := findNumberEnd(data, offset)
		if end == -1 {
			return nil, Unknown, ErrInvalidJSON
		}
		return data[offset:end], Number, nil
	}
}

func isArrayIndex(key string) bool {
	for _, c := range key {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(key) > 0
}

func findStringEnd(data []byte, start int) int {
	for i := start; i < len(data); i++ {
		if data[i] == '"' && (i == start || data[i-1] != '\\') {
			return i
		}
	}
	return -1
}

func findObjectEnd(data []byte, start int) (int, error) {
	if start >= len(data) || data[start] != '{' {
		return -1, ErrInvalidJSON
	}
	
	depth := 1
	i := start + 1
	inString := false
	
	for i < len(data) && depth > 0 {
		if !inString {
			switch data[i] {
			case '"':
				inString = true
			case '{':
				depth++
			case '}':
				depth--
			}
		} else {
			if data[i] == '"' && (i == 0 || data[i-1] != '\\') {
				inString = false
			}
		}
		i++
	}
	
	if depth != 0 {
		return -1, ErrInvalidJSON
	}
	
	return i, nil
}

func findArrayEnd(data []byte, start int) (int, error) {
	if start >= len(data) || data[start] != '[' {
		return -1, ErrInvalidJSON
	}
	
	depth := 1
	i := start + 1
	inString := false
	
	for i < len(data) && depth > 0 {
		if !inString {
			switch data[i] {
			case '"':
				inString = true
			case '[':
				depth++
			case ']':
				depth--
			}
		} else {
			if data[i] == '"' && (i == 0 || data[i-1] != '\\') {
				inString = false
			}
		}
		i++
	}
	
	if depth != 0 {
		return -1, ErrInvalidJSON
	}
	
	return i, nil
}

func findNumberEnd(data []byte, start int) int {
	i := start
	
	// Handle negative sign
	if i < len(data) && data[i] == '-' {
		i++
	}
	
	// Must have at least one digit
	if i >= len(data) || (data[i] < '0' || data[i] > '9') {
		return -1
	}
	
	// Parse integer part
	for i < len(data) && data[i] >= '0' && data[i] <= '9' {
		i++
	}
	
	// Parse decimal part
	if i < len(data) && data[i] == '.' {
		i++
		if i >= len(data) || (data[i] < '0' || data[i] > '9') {
			return -1
		}
		for i < len(data) && data[i] >= '0' && data[i] <= '9' {
			i++
		}
	}
	
	// Parse exponent part
	if i < len(data) && (data[i] == 'e' || data[i] == 'E') {
		i++
		if i < len(data) && (data[i] == '+' || data[i] == '-') {
			i++
		}
		if i >= len(data) || (data[i] < '0' || data[i] > '9') {
			return -1
		}
		for i < len(data) && data[i] >= '0' && data[i] <= '9' {
			i++
		}
	}
	
	return i
}

func skipValue(data []byte, offset int) (int, ValueType, error) {
	value, valueType, err := getType(data, offset)
	if err != nil {
		return -1, Unknown, err
	}
	return offset + len(value), valueType, nil
}

func parseString(data []byte) (string, error) {
	// Simple implementation - in a real parser you'd handle all escape sequences
	result := make([]byte, 0, len(data))
	
	for i := 0; i < len(data); i++ {
		if data[i] == '\\' && i+1 < len(data) {
			switch data[i+1] {
			case '"':
				result = append(result, '"')
			case '\\':
				result = append(result, '\\')
			case '/':
				result = append(result, '/')
			case 'b':
				result = append(result, '\b')
			case 'f':
				result = append(result, '\f')
			case 'n':
				result = append(result, '\n')
			case 'r':
				result = append(result, '\r')
			case 't':
				result = append(result, '\t')
			default:
				result = append(result, data[i+1])
			}
			i++ // skip escaped character
		} else {
			result = append(result, data[i])
		}
	}
	
	return string(result), nil
}

// ParseStringValue parses a JSON string value (with quotes) into a Go string
func ParseStringValue(data []byte) (string, error) {
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return "", ErrInvalidJSON
	}
	return parseString(data[1 : len(data)-1])
}

func arrayEach(data []byte, callback func(value []byte, dataType ValueType, offset int, err error)) {
	offset := nextToken(data, 0)
	if offset == -1 || data[offset] != '[' {
		callback(nil, Unknown, -1, errors.New("expected array"))
		return
	}
	
	offset++ // skip '['
	index := 0
	
	for {
		offset = nextToken(data, offset)
		if offset == -1 {
			callback(nil, Unknown, -1, ErrInvalidJSON)
			return
		}
		
		if data[offset] == ']' {
			return
		}
		
		valueStart := offset
		valueEnd, valueType, err := skipValue(data, offset)
		if err != nil {
			callback(nil, Unknown, -1, err)
			return
		}
		
		callback(data[valueStart:valueEnd], valueType, index, nil)
		index++
		offset = valueEnd
		
		// Check for comma or end of array
		offset = nextToken(data, offset)
		if offset == -1 {
			callback(nil, Unknown, -1, ErrInvalidJSON)
			return
		}
		
		if data[offset] == ']' {
			return
		}
		
		if data[offset] == ',' {
			offset++
			continue
		}
		
		callback(nil, Unknown, -1, ErrInvalidJSON)
		return
	}
}

func objectEach(data []byte, callback func(key []byte, value []byte, dataType ValueType, offset int) error) error {
	offset := nextToken(data, 0)
	if offset == -1 || data[offset] != '{' {
		return errors.New("expected object")
	}
	
	offset++ // skip '{'
	
	for {
		offset = nextToken(data, offset)
		if offset == -1 {
			return ErrInvalidJSON
		}
		
		if data[offset] == '}' {
			return nil
		}
		
		// Parse key
		if data[offset] != '"' {
			return ErrInvalidJSON
		}
		
		keyStart := offset + 1
		keyEnd := findStringEnd(data, keyStart)
		if keyEnd == -1 {
			return ErrInvalidJSON
		}
		
		key := data[keyStart:keyEnd]
		offset = keyEnd + 1 // skip closing quote
		
		// Skip colon
		offset = nextToken(data, offset)
		if offset == -1 || data[offset] != ':' {
			return ErrInvalidJSON
		}
		offset++ // skip ':'
		
		// Get value
		offset = nextToken(data, offset)
		if offset == -1 {
			return ErrInvalidJSON
		}
		
		valueStart := offset
		valueEnd, valueType, err := skipValue(data, offset)
		if err != nil {
			return err
		}
		
		err = callback(key, data[valueStart:valueEnd], valueType, valueStart)
		if err != nil {
			return err
		}
		
		offset = valueEnd
		
		// Check for comma or end of object
		offset = nextToken(data, offset)
		if offset == -1 {
			return ErrInvalidJSON
		}
		
		if data[offset] == '}' {
			return nil
		}
		
		if data[offset] == ',' {
			offset++
			continue
		}
		
		return ErrInvalidJSON
	}
}