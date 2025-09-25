package jsonparser

import (
	"testing"

	buger "github.com/buger/jsonparser"
)

var testJson = []byte(`{
	"user": "Muccarini",
	"age": 25,
	"active": true,
	"profile": {
		"email": "luca@example.com",
		"settings": {
			"theme": "dark",
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

func BenchmarkGet_User(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := GetString(testJson, "user")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkGetBurger_User(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := buger.GetString(testJson, "user")
		if err != nil {
			b.Error(err)
		}
	}
}
