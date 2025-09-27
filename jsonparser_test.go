package jsonparser

import (
	"fmt"
	"testing"

	buger "github.com/buger/jsonparser"
)

var testJson = []byte(`{
	"email": "This is a long email to settings performance",
	"age": 25,
	"settings": true,
	"floatingPoint": 12.34,
	"nullValue": null,
	"profile": {
		"email": "luca@example.com",
		"settings": {
			"theme": "dark theme is a bit longer to test performance",
			"notifications": true
		}
	},
	"hobbies": ["coding", "reading", "gaming"]
}`)

var largeTestJson = []byte(`{
	"user": "Luca",
	"age": 25,
	"active": true,
	"profile": {
		"email": "luca@example.com",
		"firstName": "Luca",
		"lastName": "Muccarini",
		"address": {
			"street": "123 Main St",
			"city": "Rome",
			"country": "Italy",
			"zipCode": "00100"
		},
		"settings": {
			"theme": "dark",
			"notifications": true,
			"privacy": {
				"showEmail": false,
				"showPhone": true
			}
		}
	},
	"hobbies": ["coding", "reading", "gaming", "traveling", "photography"],
	"friends": [
		{"name": "Marco", "age": 28},
		{"name": "Sara", "age": 26},
		{"name": "Giovanni", "age": 30}
	],
	"metadata": {
		"created": "2024-01-01",
		"updated": "2024-12-01",
		"version": 1
	}
}`)

func Benchmark_GetString(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	var value string
	var err error

	for i := 0; i < b.N; i++ {
		value, err = GetString(testJson, "user")
		if err != nil {
			b.Error(err)
		}
	}

	fmt.Printf("Value:%v | %T\n", value, value)
}

func Benchmark_GetString_Nested(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	var value string
	var err error

	for i := 0; i < b.N; i++ {
		value, err = GetString(testJson, "profile", "settings", "theme")
		if err != nil {
			b.Error(err)
		}
	}

	fmt.Printf("Value:%v | %T\n", value, value)
}

func Benchmark_GetInt(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	var value int
	var err error

	for i := 0; i < b.N; i++ {
		value, err = GetInt(testJson, "age")
		if err != nil {
			b.Error(err)
		}
	}

	fmt.Printf("Value:%v | %T\n", value, value)
}

func Benchmark_GetInt64(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	var value int64
	var err error

	for i := 0; i < b.N; i++ {
		value, err = GetInt64(testJson, "age")
		if err != nil {
			b.Error(err)
		}
	}

	fmt.Printf("Value:%v | %T\n", value, value)
}

func Benchmark_GetBool(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	var value bool
	var err error

	for i := 0; i < b.N; i++ {
		value, err = GetBool(testJson, "active")
		if err != nil {
			b.Error(err)
		}
	}

	fmt.Printf("Value:%v | %T\n", value, value)
}

func Benchmark_GetFloat64(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	var value float64
	var err error

	for i := 0; i < b.N; i++ {
		value, err = GetFloat64(testJson, 64, "floatingPoint")
		if err != nil {
			b.Error(err)
		}
	}

	fmt.Printf("Value:%v | %T\n", value, value)
}

func Benchmark_GetNull(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	var value interface{}
	var err error

	for i := 0; i < b.N; i++ {
		value, err = GetString(testJson, "nullValue")
		if err != nil {
			b.Error(err)
		}
	}

	fmt.Printf("Value:%v | %T\n", value, value)
}

func Benchmark_GetBurger(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	var valueRes int
	var err error
	var valueSlice []string

	for i := 0; i < b.N; i++ {
		valueSlice = make([]string, 0, 3)
		valueRes, err = buger.ArrayEach(largeTestJson, func(value []byte, dataType buger.ValueType, offset int, err error) {
			valueSlice = append(valueSlice, string(value))
		}, "hobbies")
		if err != nil {
			b.Error(err)
		}
	}

	_ = valueRes

	fmt.Printf("Value:%v | %T\n", valueSlice, valueSlice)
}

func Benchmark_GetBurger_Nested(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	var value string
	var err error

	for i := 0; i < b.N; i++ {
		value, err = buger.GetString(testJson, "profile", "settings", "theme")
		if err != nil {
			b.Error(err)
		}
	}

	fmt.Printf("Value:%v | %T\n", value, value)
}
