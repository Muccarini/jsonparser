package jsonparser

import (
	"encoding/json"
	"testing"

	buger "github.com/buger/jsonparser"
)

// Test data for benchmarks
var testJSON = []byte(`{
	"user": {
		"id": 12345,
		"name": "John Doe",
		"email": "john.doe@example.com",
		"age": 30,
		"isActive": true,
		"profile": {
			"bio": "Software developer with 5+ years experience",
			"location": {
				"city": "San Francisco",
				"state": "CA",
				"country": "USA",
				"coordinates": {
					"latitude": 37.7749,
					"longitude": -122.4194
				}
			},
			"preferences": {
				"theme": "dark",
				"notifications": true,
				"language": "en-US"
			}
		}
	},
	"posts": [ //Get(json, "posts").Where("tags", "go")
		{
			"id": 1,
			"title": "Getting Started with Go",
			"content": "Go is a powerful programming language for building scalable applications. It offers excellent performance and simplicity.",
			"author": "John Doe",
			"publishDate": "2024-01-15T10:30:00Z",
			"tags": ["go", "programming", "tutorial"],
			"metadata": {
				"views": 1520,
				"likes": 89,
				"comments": 12,
				"featured": true
			}
		},
		{
			"id": 2,
			"title": "JSON Parsing Best Practices",
			"content": "When working with JSON data in applications, it's important to consider performance and memory usage.",
			"author": "Jane Smith",
			"publishDate": "2024-02-01T14:45:00Z",
			"tags": ["json", "parsing", "best-practices"],
			"metadata": {
				"views": 2340,
				"likes": 156,
				"comments": 23,
				"featured": false
			}
		}
	],
	"categories": ["Technology", "Programming", "Web Development", "Software Engineering"],
	"settings": {
		"version": "1.2.3",
		"apiEndpoint": "https://api.example.com/v1",
		"timeout": 5000,
		"retryAttempts": 3,
		"features": {
			"enableCaching": true,
			"enableLogging": false,
			"enableMetrics": true,
			"debugMode": null
		}
	},
	"stats": {
		"totalUsers": 15420,
		"activeUsers": 8934,
		"averageSessionTime": 285.5,
		"conversionRate": 0.034,
		"lastUpdated": "2024-03-15T09:00:00Z"
	},
	"emptyArray": [],
	"emptyObject": {},
	"nullValue": null,
	"booleanTrue": true,
	"booleanFalse": false,
	"largeNumber": 9223372036854775807,
	"floatingPoint": 3.14159265359,
	"scientificNotation": 1.23e-4,
	"unicodeString": "Hello 世界 🌍",
	"escapedCharacters": "Line 1\nLine 2\tTabbed\r\nWindows line ending",
	"specialChars": "Quote: \"Hello\", Backslash: \\, Forward slash: /"
}`)

// Large JSON for stress testing
var largeJSON []byte

func init() {
	// Create a larger JSON by repeating the test structure
	largeJSONStr := `{"users": [`
	for i := 0; i < 1000; i++ {
		if i > 0 {
			largeJSONStr += ","
		}
		largeJSONStr += `{
			"id": ` + string(rune('0'+i%10)) + `,
			"name": "User ` + string(rune('0'+i%10)) + `",
			"profile": {
				"location": {
					"city": "City ` + string(rune('0'+i%10)) + `",
					"coordinates": {"lat": 37.7749, "lng": -122.4194}
				}
			},
			"posts": [
				{"id": 1, "title": "Post 1", "content": "Content here"},
				{"id": 2, "title": "Post 2", "content": "More content"}
			]
		}`
	}
	largeJSONStr += `]}`
	largeJSON = []byte(largeJSONStr)
}

// Benchmark: Simple field extraction (root level)
func BenchmarkGet_SimpleField(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Get(testJSON, "booleanTrue")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark: Nested field extraction (2 levels deep)
func BenchmarkGet_NestedField(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Get(testJSON, "user", "name")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark: Deep nested field extraction (4 levels deep)
func BenchmarkGet_DeepNestedField(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Get(testJSON, "user", "profile", "location", "city")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark: Field at end of JSON
func BenchmarkGet_FieldAtEnd(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Get(testJSON, "specialChars")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark: Non-existent field (worst case)
func BenchmarkGet_NonExistentField(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Get(testJSON, "nonExistentField")
		if err == nil {
			b.Fatal("Expected error for non-existent field")
		}
	}
}

// Benchmark: Large JSON with deep nesting
func BenchmarkGet_LargeJSON(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Get(largeJSON, "users")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Memory allocation benchmarks
func BenchmarkGet_SimpleField_Memory(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Get(testJSON, "booleanTrue")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGet_NestedField_Memory(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Get(testJSON, "user", "name")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGet_DeepNestedField_Memory(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Get(testJSON, "user", "profile", "location", "city")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Comparison with standard library encoding/json
func BenchmarkStandardJSON_Unmarshal(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result map[string]interface{}
		err := json.Unmarshal(testJSON, &result)
		if err != nil {
			b.Fatal(err)
		}
		// Access nested field
		if user, ok := result["user"].(map[string]interface{}); ok {
			_ = user["name"]
		}
	}
}

func BenchmarkStandardJSON_Unmarshal_Memory(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result map[string]interface{}
		err := json.Unmarshal(testJSON, &result)
		if err != nil {
			b.Fatal(err)
		}
		// Access nested field
		if user, ok := result["user"].(map[string]interface{}); ok {
			_ = user["name"]
		}
	}
}

// Benchmark different value types
func BenchmarkGet_StringValue(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Get(testJSON, "user", "name")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGet_NumberValue(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Get(testJSON, "user", "id")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGet_BooleanValue(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Get(testJSON, "booleanTrue")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGet_ObjectValue(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Get(testJSON, "user", "profile")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGet_ArrayValue(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Get(testJSON, "categories")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Comparison with buger/jsonparser
func BenchmarkBuger_SimpleField(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := buger.GetString(testJSON, "booleanTrue")
		if err != nil {
			// booleanTrue is not a string, try as boolean
			_, err = buger.GetBoolean(testJSON, "booleanTrue")
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkBuger_NestedField(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := buger.GetString(testJSON, "user", "name")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBuger_DeepNestedField(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := buger.GetString(testJSON, "user", "profile", "location", "city")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBuger_NumberValue(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := buger.GetInt(testJSON, "user", "id")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBuger_BooleanValue(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := buger.GetBoolean(testJSON, "booleanTrue")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBuger_ArrayValue(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _, err := buger.Get(testJSON, "categories")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBuger_ObjectValue(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _, err := buger.Get(testJSON, "user", "profile")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBuger_NonExistentField(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := buger.GetString(testJSON, "nonExistentField")
		if err == nil {
			b.Fatal("Expected error for non-existent field")
		}
	}
}

func BenchmarkBuger_LargeJSON(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _, err := buger.Get(largeJSON, "users")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Memory allocation benchmarks for buger
func BenchmarkBuger_NestedField_Memory(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := buger.GetString(testJSON, "user", "name")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBuger_DeepNestedField_Memory(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := buger.GetString(testJSON, "user", "profile", "location", "city")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBuger_NumberValue_Memory(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := buger.GetInt(testJSON, "user", "id")
		if err != nil {
			b.Fatal(err)
		}
	}
}
