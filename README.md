# jsonparser

A fast JSON parser with filters and zero memory allocation for Go.

## Features

- **Zero Memory Allocation**: Core parsing operations don't allocate memory
- **Fast Performance**: Direct byte slice manipulation without copying
- **JSONPath-like Navigation**: Access nested values using dot notation
- **Type Safety**: Strongly typed value extraction functions
- **Array & Object Iteration**: Efficient iteration over arrays and objects
- **Error Handling**: Comprehensive error reporting for invalid JSON or missing keys

## Performance

Benchmark results on AMD EPYC 7763:

```
BenchmarkGet-4          	 1803128	       664.3 ns/op	       0 B/op	       0 allocs/op
BenchmarkArrayEach-4    	 5781552	       218.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkObjectEach-4   	 5352842	       225.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkGetString-4    	 4034606	       293.7 ns/op	      32 B/op	       2 allocs/op
```

The core `Get` operations achieve zero allocations. Only string extraction allocates memory for the result string.

## Installation

```bash
go get github.com/Muccarini/jsonparser
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/Muccarini/jsonparser"
)

func main() {
    jsonData := []byte(`{"user":{"name":"Alice","age":30,"roles":["admin","user"]}}`)
    
    // Extract values by path
    name, _ := jsonparser.GetString(jsonData, "user", "name")
    age, _ := jsonparser.GetInt(jsonData, "user", "age")
    firstRole, _ := jsonparser.GetString(jsonData, "user", "roles", "0")
    
    fmt.Printf("Name: %s, Age: %d, First Role: %s\n", name, age, firstRole)
    // Output: Name: Alice, Age: 30, First Role: admin
}
```

## API Reference

### Core Functions

#### Get
```go
func Get(data []byte, keys ...string) (value []byte, dataType ValueType, offset int, err error)
```
Returns the raw bytes of a value, its type, and offset. This is the zero-allocation foundation function.

#### Typed Getters
```go
func GetString(data []byte, keys ...string) (string, error)
func GetInt(data []byte, keys ...string) (int64, error)
func GetFloat(data []byte, keys ...string) (float64, error)
func GetBoolean(data []byte, keys ...string) (bool, error)
```
Extract values with automatic type conversion and validation.

#### Iteration Functions
```go
func ArrayEach(data []byte, callback func(value []byte, dataType ValueType, offset int, err error), keys ...string)
func ObjectEach(data []byte, callback func(key []byte, value []byte, dataType ValueType, offset int) error, keys ...string) error
```
Efficiently iterate over arrays and objects without loading all elements into memory.

### Path Navigation

Use keys to navigate nested structures:

```go
// Object navigation
jsonparser.GetString(data, "user", "address", "city")

// Array navigation (use string indices)
jsonparser.GetString(data, "users", "0", "name")

// Mixed navigation
jsonparser.GetFloat(data, "users", "0", "coordinates", "lat")
```

### Value Types

```go
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
```

## Examples

### Basic Usage

```go
jsonData := []byte(`{
    "name": "Alice Johnson",
    "age": 30,
    "active": true,
    "balance": 1234.56,
    "metadata": null
}`)

name, _ := jsonparser.GetString(jsonData, "name")        // "Alice Johnson"
age, _ := jsonparser.GetInt(jsonData, "age")             // 30
active, _ := jsonparser.GetBoolean(jsonData, "active")   // true
balance, _ := jsonparser.GetFloat(jsonData, "balance")   // 1234.56

// Check type without parsing
_, valueType, _, _ := jsonparser.Get(jsonData, "metadata")
if valueType == jsonparser.Null {
    fmt.Println("Metadata is null")
}
```

### Array Handling

```go
jsonData := []byte(`{"users":[{"name":"Alice"},{"name":"Bob"}]}`)

// Access specific array elements
firstName, _ := jsonparser.GetString(jsonData, "users", "0", "name") // "Alice"
secondName, _ := jsonparser.GetString(jsonData, "users", "1", "name") // "Bob"

// Iterate over array
jsonparser.ArrayEach(jsonData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
    if dataType == jsonparser.Object {
        name, _ := jsonparser.GetString(value, "name")
        fmt.Printf("User: %s\n", name)
    }
}, "users")
```

### Object Iteration

```go
jsonData := []byte(`{"user":{"name":"Alice","age":30,"city":"Boston"}}`)

jsonparser.ObjectEach(jsonData, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
    fmt.Printf("Key: %s, Value: %s, Type: %v\n", string(key), string(value), dataType)
    return nil
}, "user")
```

### Error Handling

```go
jsonData := []byte(`{"user":{"name":"Alice"}}`)

// Missing key
_, err := jsonparser.GetString(jsonData, "user", "age")
if err == jsonparser.ErrKeyPathNotFound {
    fmt.Println("Age not found")
}

// Type mismatch
_, err = jsonparser.GetInt(jsonData, "user", "name")
if err != nil {
    fmt.Println("Name is not a number")
}

// Array out of bounds
_, err = jsonparser.GetString(jsonData, "users", "10", "name")
if err != nil {
    fmt.Println("Array index out of bounds")
}
```

### Complex JSON Navigation

```go
complexJSON := []byte(`{
    "users": [
        {
            "id": 1,
            "profile": {
                "name": "Alice",
                "address": {
                    "coordinates": {
                        "lat": 42.3601,
                        "lng": -71.0589
                    }
                }
            }
        }
    ]
}`)

// Navigate deep nested structure
lat, _ := jsonparser.GetFloat(complexJSON, "users", "0", "profile", "address", "coordinates", "lat")
fmt.Printf("Latitude: %f\n", lat) // 42.360100
```

## Design Principles

### Zero Allocation
The parser works directly with the input byte slice, returning sub-slices that reference the original data. No intermediate objects or copies are created during parsing.

### Streaming Friendly
Since the parser doesn't build a complete AST, it can work with partial JSON data and is suitable for streaming applications.

### Error Recovery
The parser provides detailed error information to help identify issues with JSON structure or incorrect key paths.

## Limitations

- String values require allocation when extracted (unavoidable in Go)
- The parser assumes valid UTF-8 input
- Escape sequence handling in strings is basic (covers common cases)
- No support for JSON Schema validation

## Contributing

Contributions are welcome! Please ensure:

1. All tests pass: `go test -v`
2. Benchmarks maintain zero allocations: `go test -bench=. -benchmem`
3. Code follows Go conventions

## License

MIT License - see [LICENSE](LICENSE) for details.
