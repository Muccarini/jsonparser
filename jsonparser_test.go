package jsonparser

import (
	"testing"
)

func TestGet(t *testing.T) {
	data := []byte(`{"name":"John","age":30,"address":{"city":"New York","zip":"10001"},"hobbies":["reading","swimming"]}`)
	
	tests := []struct {
		name        string
		keys        []string
		expected    string
		expectedType ValueType
		shouldError bool
	}{
		{"get name", []string{"name"}, `"John"`, String, false},
		{"get age", []string{"age"}, "30", Number, false},
		{"get nested city", []string{"address", "city"}, `"New York"`, String, false},
		{"get nested zip", []string{"address", "zip"}, `"10001"`, String, false},
		{"get array element", []string{"hobbies", "0"}, `"reading"`, String, false},
		{"get array element 2", []string{"hobbies", "1"}, `"swimming"`, String, false},
		{"get nonexistent key", []string{"nonexistent"}, "", NotExist, true},
		{"get out of bounds array", []string{"hobbies", "5"}, "", NotExist, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, dataType, _, err := Get(data, tt.keys...)
			
			if tt.shouldError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			
			if string(value) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, string(value))
			}
			
			if dataType != tt.expectedType {
				t.Errorf("expected type %v, got %v", tt.expectedType, dataType)
			}
		})
	}
}

func TestGetString(t *testing.T) {
	data := []byte(`{"name":"John Doe","empty":"","with_quotes":"He said \"hello\""}`)
	
	tests := []struct {
		name        string
		keys        []string
		expected    string
		shouldError bool
	}{
		{"simple string", []string{"name"}, "John Doe", false},
		{"empty string", []string{"empty"}, "", false},
		{"string with quotes", []string{"with_quotes"}, `He said "hello"`, false},
		{"nonexistent key", []string{"missing"}, "", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetString(data, tt.keys...)
			
			if tt.shouldError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestGetInt(t *testing.T) {
	data := []byte(`{"positive":42,"negative":-17,"zero":0,"large":9223372036854775807}`)
	
	tests := []struct {
		name        string
		keys        []string
		expected    int64
		shouldError bool
	}{
		{"positive number", []string{"positive"}, 42, false},
		{"negative number", []string{"negative"}, -17, false},
		{"zero", []string{"zero"}, 0, false},
		{"large number", []string{"large"}, 9223372036854775807, false},
		{"nonexistent key", []string{"missing"}, 0, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetInt(data, tt.keys...)
			
			if tt.shouldError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestGetFloat(t *testing.T) {
	data := []byte(`{"pi":3.14159,"negative":-2.5,"scientific":1.23e-4,"zero":0.0}`)
	
	tests := []struct {
		name        string
		keys        []string
		expected    float64
		shouldError bool
	}{
		{"decimal number", []string{"pi"}, 3.14159, false},
		{"negative float", []string{"negative"}, -2.5, false},
		{"scientific notation", []string{"scientific"}, 1.23e-4, false},
		{"zero float", []string{"zero"}, 0.0, false},
		{"nonexistent key", []string{"missing"}, 0, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetFloat(data, tt.keys...)
			
			if tt.shouldError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			
			if result != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestGetBoolean(t *testing.T) {
	data := []byte(`{"true_val":true,"false_val":false}`)
	
	tests := []struct {
		name        string
		keys        []string
		expected    bool
		shouldError bool
	}{
		{"true value", []string{"true_val"}, true, false},
		{"false value", []string{"false_val"}, false, false},
		{"nonexistent key", []string{"missing"}, false, true},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetBoolean(data, tt.keys...)
			
			if tt.shouldError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			
			if result != tt.expected {
				t.Errorf("expected %t, got %t", tt.expected, result)
			}
		})
	}
}

func TestArrayEach(t *testing.T) {
	data := []byte(`{"numbers":[1,2,3],"strings":["a","b","c"],"mixed":[1,"hello",true,null]}`)
	
	t.Run("numbers array", func(t *testing.T) {
		var results []string
		var types []ValueType
		
		ArrayEach(data, func(value []byte, dataType ValueType, offset int, err error) {
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			results = append(results, string(value))
			types = append(types, dataType)
		}, "numbers")
		
		expected := []string{"1", "2", "3"}
		expectedTypes := []ValueType{Number, Number, Number}
		
		if len(results) != len(expected) {
			t.Errorf("expected %d results, got %d", len(expected), len(results))
			return
		}
		
		for i, result := range results {
			if result != expected[i] {
				t.Errorf("expected %q at index %d, got %q", expected[i], i, result)
			}
			if types[i] != expectedTypes[i] {
				t.Errorf("expected type %v at index %d, got %v", expectedTypes[i], i, types[i])
			}
		}
	})
	
	t.Run("mixed array", func(t *testing.T) {
		var results []string
		var types []ValueType
		
		ArrayEach(data, func(value []byte, dataType ValueType, offset int, err error) {
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			results = append(results, string(value))
			types = append(types, dataType)
		}, "mixed")
		
		expected := []string{"1", `"hello"`, "true", "null"}
		expectedTypes := []ValueType{Number, String, Boolean, Null}
		
		if len(results) != len(expected) {
			t.Errorf("expected %d results, got %d", len(expected), len(results))
			return
		}
		
		for i, result := range results {
			if result != expected[i] {
				t.Errorf("expected %q at index %d, got %q", expected[i], i, result)
			}
			if types[i] != expectedTypes[i] {
				t.Errorf("expected type %v at index %d, got %v", expectedTypes[i], i, types[i])
			}
		}
	})
}

func TestObjectEach(t *testing.T) {
	data := []byte(`{"user":{"name":"John","age":30,"active":true}}`)
	
	t.Run("object iteration", func(t *testing.T) {
		var keys []string
		var values []string
		var types []ValueType
		
		err := ObjectEach(data, func(key []byte, value []byte, dataType ValueType, offset int) error {
			keys = append(keys, string(key))
			values = append(values, string(value))
			types = append(types, dataType)
			return nil
		}, "user")
		
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		
		expectedKeys := []string{"name", "age", "active"}
		expectedValues := []string{`"John"`, "30", "true"}
		expectedTypes := []ValueType{String, Number, Boolean}
		
		if len(keys) != len(expectedKeys) {
			t.Errorf("expected %d keys, got %d", len(expectedKeys), len(keys))
			return
		}
		
		for i, key := range keys {
			if key != expectedKeys[i] {
				t.Errorf("expected key %q at index %d, got %q", expectedKeys[i], i, key)
			}
			if values[i] != expectedValues[i] {
				t.Errorf("expected value %q at index %d, got %q", expectedValues[i], i, values[i])
			}
			if types[i] != expectedTypes[i] {
				t.Errorf("expected type %v at index %d, got %v", expectedTypes[i], i, types[i])
			}
		}
	})
}

func TestComplexJSON(t *testing.T) {
	complexJSON := []byte(`{
		"users": [
			{
				"id": 1,
				"name": "Alice Johnson",
				"email": "alice@example.com",
				"address": {
					"street": "123 Main St",
					"city": "Boston",
					"coordinates": {
						"lat": 42.3601,
						"lng": -71.0589
					}
				},
				"active": true,
				"roles": ["admin", "user"]
			},
			{
				"id": 2,
				"name": "Bob Smith",
				"email": "bob@example.com",
				"address": {
					"street": "456 Oak Ave",
					"city": "Seattle",
					"coordinates": {
						"lat": 47.6062,
						"lng": -122.3321
					}
				},
				"active": false,
				"roles": ["user"]
			}
		],
		"metadata": {
			"total": 2,
			"version": "1.0.0"
		}
	}`)
	
	tests := []struct {
		name        string
		keys        []string
		expected    string
		expectedType ValueType
	}{
		{"first user name", []string{"users", "0", "name"}, `"Alice Johnson"`, String},
		{"first user lat", []string{"users", "0", "address", "coordinates", "lat"}, "42.3601", Number},
		{"second user active", []string{"users", "1", "active"}, "false", Boolean},
		{"first user first role", []string{"users", "0", "roles", "0"}, `"admin"`, String},
		{"metadata total", []string{"metadata", "total"}, "2", Number},
		{"metadata version", []string{"metadata", "version"}, `"1.0.0"`, String},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, dataType, _, err := Get(complexJSON, tt.keys...)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			
			if string(value) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, string(value))
			}
			
			if dataType != tt.expectedType {
				t.Errorf("expected type %v, got %v", tt.expectedType, dataType)
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	t.Run("empty json", func(t *testing.T) {
		_, _, _, err := Get([]byte(""), "key")
		if err == nil {
			t.Errorf("expected error for empty JSON")
		}
	})
	
	t.Run("invalid json", func(t *testing.T) {
		_, _, _, err := Get([]byte("{invalid"), "key")
		if err == nil {
			t.Errorf("expected error for invalid JSON")
		}
	})
	
	t.Run("null value", func(t *testing.T) {
		data := []byte(`{"value":null}`)
		value, dataType, _, err := Get(data, "value")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}
		
		if string(value) != "null" {
			t.Errorf("expected null, got %q", string(value))
		}
		
		if dataType != Null {
			t.Errorf("expected Null type, got %v", dataType)
		}
	})
	
	t.Run("empty object", func(t *testing.T) {
		data := []byte(`{}`)
		_, _, _, err := Get(data, "key")
		if err == nil {
			t.Errorf("expected error for missing key in empty object")
		}
	})
	
	t.Run("empty array", func(t *testing.T) {
		data := []byte(`{"arr":[]}`)
		_, _, _, err := Get(data, "arr", "0")
		if err == nil {
			t.Errorf("expected error for accessing element in empty array")
		}
	})
}

// Benchmark tests to ensure zero allocation performance
func BenchmarkGet(b *testing.B) {
	data := []byte(`{"users":[{"id":1,"name":"Alice","email":"alice@example.com","address":{"city":"Boston","coordinates":{"lat":42.3601,"lng":-71.0589}}}]}`)
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		_, _, _, _ = Get(data, "users", "0", "address", "coordinates", "lat")
	}
}

func BenchmarkGetString(b *testing.B) {
	data := []byte(`{"users":[{"id":1,"name":"Alice Johnson","email":"alice@example.com"}]}`)
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		_, _ = GetString(data, "users", "0", "name")
	}
}

func BenchmarkArrayEach(b *testing.B) {
	data := []byte(`{"numbers":[1,2,3,4,5,6,7,8,9,10]}`)
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		ArrayEach(data, func(value []byte, dataType ValueType, offset int, err error) {
			// Do nothing, just iterate
		}, "numbers")
	}
}

func BenchmarkObjectEach(b *testing.B) {
	data := []byte(`{"user":{"id":1,"name":"Alice","email":"alice@example.com","city":"Boston","active":true}}`)
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		ObjectEach(data, func(key []byte, value []byte, dataType ValueType, offset int) error {
			return nil
		}, "user")
	}
}