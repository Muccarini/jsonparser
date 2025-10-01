# jsonparser

## Disclaimer

⚠️ **This is an educational project and my first code written in Go.** No extensive testing has been conducted, so usage in production is discouraged. This project was inspired by Salvatore Sanfilippo's [video](https://www.youtube.com/watch?v=EP75QgSC73I).

## About

This is an educational project created to learn and explore Go programming concepts including:
- Memory management: escape analysis and zero-allocation programming techniques
- Slice manipulation and optimization
- Performance benchmarking

## Performance Comparison

Benchmarks comparing this library against [buger/jsonparser](https://github.com/buger/jsonparser) and Go's standard `encoding/json`.

| Operation | Library | Time (ns/op) | Memory (B/op) | Allocs/op |
|-----------|---------|--------------|---------------|-----------|
| Simple String | Muccarini | 42.81 | 16 | 1 |
| | Buger | 57.79 | 16 | 1 |
| | Standard | 25,836 | 11,088 | 296 |
| Nested String | Muccarini | 1,352 | 48 | 1 |
| | Buger | 1,556 | 48 | 1 |
| | Standard | 25,477 | 11,088 | 296 |
| Integer | Muccarini | 529.2 | 0 | 0 |
| | Buger | 670.1 | 0 | 0 |
| | Standard | 25,538 | 11,088 | 296 |
| Float | Muccarini | 793.4 | 0 | 0 |
| | Buger | 938.4 | 0 | 0 |
| | Standard | 25,685 | 11,088 | 296 |
| Boolean | Muccarini | 455.7 | 0 | 0 |
| | Buger | 600.2 | 0 | 0 |
| | Standard | 25,731 | 11,088 | 296 |
| Array Access | Muccarini | 1,289 | 5 | 1 |
| | Buger | 1,525 | 5 | 1 |
| | Standard | 25,441 | 11,088 | 296 |
| Array Iteration | Muccarini | 1,289 | 0 | 0 |
| | Buger | 1,514 | 0 | 0 |
| | Standard | 25,750 | 11,088 | 296 |
| Many Fields | Muccarini | 1,972 | 8 | 1 |
| | Buger | 2,281 | 8 | 1 |
| | Standard | 25,829 | 11,088 | 296 |

## Key Results

- **Speed vs buger/jsonparser**: 13-26% faster across all operations
- **Speed vs standard library**: 1,200-5,540% faster (13-56x) across all operations
- **Memory**: Zero allocations for primitive types (Integer, Float, Boolean, Array Iteration)

## TODO

- [ ] Array iteration functionality
- [ ] Object serialization support

## Usage

```go
package main

import (
    "fmt"
    "github.com/muccarini/jsonparser"
)

func main() {
    json := []byte(`{
        "name": "John",
        "age": 30,
        "active": true,
        "score": 85.5,
        "tags": ["dev", "golang"],
        "profile": {"city": "NYC"}
    }`)

    // Basic extraction
    name, _ := jsonparser.GetString(json, "name")
    age, _ := jsonparser.GetInt(json, "age")
    active, _ := jsonparser.GetBool(json, "active")
    score, _ := jsonparser.GetFloat64(json, "score")

    // Array and nested access
    tag, _ := jsonparser.GetString(json, "tags", "1")
    city, _ := jsonparser.GetString(json, "profile", "city")

    fmt.Printf("Name: %s, Age: %d, Active: %t\n", name, age, active)
    fmt.Printf("Score: %.1f, Tag: %s, City: %s\n", score, tag, city)

    // Array iteration
    fmt.Print("Tags: ")
    jsonparser.ForeachArrayElement(json, func(value []byte, index int) {
        fmt.Printf("[%d]=%s ", index, string(value))
    }, "tags")
    fmt.Println()
}
```

## Installation

```bash
go get github.com/muccarini/jsonparser
```
