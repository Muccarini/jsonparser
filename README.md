# jsonparser

## Disclaimer

⚠️ **This is an educational project and my first code written in Go.** No extensive testing has been conducted, so usage in production is discouraged. This project was inspired by @antirez's [video](https://www.youtube.com/watch?v=EP75QgSC73I).

## About

This is an educational project created to learn and explore Go programming concepts including:
- Memory management: escape analysis and zero-allocation programming techniques
- Slice manipulation and optimization
- Performance benchmarking

## Performance Comparison

Benchmarks comparing this library against [buger/jsonparser](https://github.com/buger/jsonparser) and Go's standard `encoding/json`.

| Operation | Library | Time (ns/op) | Memory (B/op) | Allocs/op |
|-----------|---------|--------------|---------------|-----------|
| Simple String | Muccarini | 36.00 | 16 | 1 |
| | Buger | 57.05 | 16 | 1 |
| | Standard | 25,786 | 11,088 | 296 |
| Nested String | Muccarini | 1,386 | 48 | 1 |
| | Buger | 1,555 | 48 | 1 |
| | Standard | 25,675 | 11,088 | 296 |
| Integer | Muccarini | 546.5 | 0 | 0 |
| | Buger | 685.2 | 0 | 0 |
| | Standard | 29,029 | 11,088 | 296 |
| Float | Muccarini | 1,035 | 0 | 0 |
| | Buger | 991.0 | 0 | 0 |
| | Standard | 28,884 | 11,088 | 296 |
| Boolean | Muccarini | 503.9 | 0 | 0 |
| | Buger | 656.7 | 0 | 0 |
| | Standard | 28,437 | 11,088 | 296 |
| Array Access | Muccarini | 1,512 | 5 | 1 |
| | Buger | 1,622 | 5 | 1 |
| | Standard | 28,455 | 11,088 | 296 |
| Many Fields | Muccarini | 2,262 | 8 | 1 |
| | Buger | 2,436 | 8 | 1 |
| | Standard | 28,346 | 11,088 | 296 |

## Key Results

- **Speed**: Competitive with buger/jsonparser, 17-717x faster than standard library
- **Memory**: Zero allocations for primitive types

## TODO

- [ ] Array iteration functionality
- [ ] Object serialization support
- [ ] Comprehensive error handling

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
    score, _ := jsonparser.GetFloat64(json, 64, "score")

    // Array and nested access
    tag, _ := jsonparser.GetString(json, "tags", "1")
    city, _ := jsonparser.GetString(json, "profile", "city")

    fmt.Printf("Name: %s, Age: %d, Active: %t\n", name, age, active)
    fmt.Printf("Score: %.1f, Tag: %s, City: %s\n", score, tag, city)
}
```

## Installation

```bash
go get github.com/muccarini/jsonparser
```
