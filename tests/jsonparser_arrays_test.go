package jsonparser_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/muccarini/jsonparser"
)

var arrayTestJson []byte

func init() {
	// Load the sample_arrays.json file for testing
	var err error
	arrayTestJson, err = os.ReadFile("sample_arrays.json")
	if err != nil {
		panic("Failed to load sample_arrays.json: " + err.Error())
	}
}

// Test empty array
func TestEmptyArray(t *testing.T) {
	result, err := jsonparser.GetString(arrayTestJson, "empty_array")
	assert.NoError(t, err)
	assert.Equal(t, "[]", result)
}

// Test string array elements
func TestStringArrayElements(t *testing.T) {
	tests := []struct {
		index    string
		expected string
	}{
		{"0", "apple"},
		{"1", "banana"},
		{"2", "cherry"},
		{"3", "date"},
	}

	for _, test := range tests {
		result, err := jsonparser.GetString(arrayTestJson, "string_array", test.index)
		assert.NoError(t, err, "Error getting string_array[%s]", test.index)
		assert.Equal(t, test.expected, result, "string_array[%s] should equal %s", test.index, test.expected)
	}
}

// Test number array elements
func TestNumberArrayElements(t *testing.T) {
	tests := []struct {
		index    string
		expected int
	}{
		{"0", 1},
		{"1", 2},
		{"2", 3},
		{"3", 4},
		{"4", 5},
		{"5", 42},
		{"6", 100},
	}

	for _, test := range tests {
		result, err := jsonparser.GetInt(arrayTestJson, "number_array", test.index)
		assert.NoError(t, err, "Error getting number_array[%s]", test.index)
		assert.Equal(t, test.expected, result, "number_array[%s] should equal %d", test.index, test.expected)
	}
}

// Test float array elements
func TestFloatArrayElements(t *testing.T) {
	tests := []struct {
		index    string
		expected float64
	}{
		{"0", 1.1},
		{"1", 2.5},
		{"2", 3.14},
		{"3", 4.0},
		{"4", 5.999},
	}

	for _, test := range tests {
		result, err := jsonparser.GetFloat64(arrayTestJson, "float_array", test.index)
		assert.NoError(t, err, "Error getting float_array[%s]", test.index)
		assert.Equal(t, test.expected, result, "float_array[%s] should equal %f", test.index, test.expected)
	}
}

// Test boolean array elements
func TestBooleanArrayElements(t *testing.T) {
	tests := []struct {
		index    string
		expected bool
	}{
		{"0", true},
		{"1", false},
		{"2", true},
		{"3", true},
		{"4", false},
	}

	for _, test := range tests {
		result, err := jsonparser.GetBool(arrayTestJson, "boolean_array", test.index)
		assert.NoError(t, err, "Error getting boolean_array[%s]", test.index)
		assert.Equal(t, test.expected, result, "boolean_array[%s] should equal %t", test.index, test.expected)
	}
}

// Test mixed array elements
func TestMixedArrayElements(t *testing.T) {
	// Test integer element
	intResult, err := jsonparser.GetInt(arrayTestJson, "mixed_array", "0")
	assert.NoError(t, err, "Error getting mixed_array[0]")
	assert.Equal(t, 1, intResult, "mixed_array[0] should equal 1")

	// Test string element
	stringResult, err := jsonparser.GetString(arrayTestJson, "mixed_array", "1")
	assert.NoError(t, err, "Error getting mixed_array[1]")
	assert.Equal(t, "hello", stringResult, "mixed_array[1] should equal hello")

	// Test boolean element
	boolResult, err := jsonparser.GetBool(arrayTestJson, "mixed_array", "2")
	assert.NoError(t, err, "Error getting mixed_array[2]")
	assert.Equal(t, true, boolResult, "mixed_array[2] should equal true")

	// Test float element
	floatResult, err := jsonparser.GetFloat64(arrayTestJson, "mixed_array", "3")
	assert.NoError(t, err, "Error getting mixed_array[3]")
	assert.Equal(t, 3.14, floatResult, "mixed_array[3] should equal 3.14")

	// Test null element (as string, should return empty or null representation)
	nullResult, err := jsonparser.GetString(arrayTestJson, "mixed_array", "4")
	assert.NoError(t, err, "Error getting mixed_array[4]")
	assert.Equal(t, "null", nullResult, "mixed_array[4] should equal null")
}

// Test null array elements
func TestNullArrayElements(t *testing.T) {
	for i := 0; i < 3; i++ {
		result, err := jsonparser.GetString(arrayTestJson, "null_array", string(rune('0'+i)))
		assert.NoError(t, err, "Error getting null_array[%d]", i)
		assert.Equal(t, "null", result, "null_array[%d] should equal null", i)
	}
}

// Test nested arrays
func TestNestedArrays(t *testing.T) {
	// Test first nested array [1, 2, 3]
	result, err := jsonparser.GetInt(arrayTestJson, "nested_arrays", "0", "1")
	assert.NoError(t, err, "Error getting nested_arrays[0][1]")
	assert.Equal(t, 2, result, "nested_arrays[0][1] should equal 2")

	// Test second nested array ["a", "b", "c"]
	result2, err := jsonparser.GetString(arrayTestJson, "nested_arrays", "1", "2")
	assert.NoError(t, err, "Error getting nested_arrays[1][2]")
	assert.Equal(t, "c", result2, "nested_arrays[1][2] should equal c")

	// Test boolean nested array
	result3, err := jsonparser.GetBool(arrayTestJson, "nested_arrays", "2", "0")
	assert.NoError(t, err, "Error getting nested_arrays[2][0]")
	assert.Equal(t, true, result3, "nested_arrays[2][0] should equal true")

	// Test deeply nested array [[1, 2], [3, 4]]
	result4, err := jsonparser.GetInt(arrayTestJson, "nested_arrays", "3", "1", "0")
	assert.NoError(t, err, "Error getting nested_arrays[3][1][0]")
	assert.Equal(t, 3, result4, "nested_arrays[3][1][0] should equal 3")
}

// Test array of objects
func TestArrayOfObjects(t *testing.T) {
	//Test first object properties
	id, err := jsonparser.GetInt(arrayTestJson, "array_of_objects", "0", "id")
	assert.NoError(t, err, "Error getting array_of_objects[0].id")
	assert.Equal(t, 1, id, "array_of_objects[0].id should equal 1")

	name, err := jsonparser.GetString(arrayTestJson, "array_of_objects", "0", "name")
	assert.NoError(t, err, "Error getting array_of_objects[0].name")
	assert.Equal(t, "John Doe", name, "array_of_objects[0].name should equal John Doe")

	email, err := jsonparser.GetString(arrayTestJson, "array_of_objects", "0", "email")
	assert.NoError(t, err, "Error getting array_of_objects[0].email")
	assert.Equal(t, "john@example.com", email, "array_of_objects[0].email should equal john@example.com")

	active, err := jsonparser.GetBool(arrayTestJson, "array_of_objects", "0", "active")
	assert.NoError(t, err, "Error getting array_of_objects[0].active")
	assert.Equal(t, true, active, "array_of_objects[0].active should equal true")

	// Test second object
	name2, err := jsonparser.GetString(arrayTestJson, "array_of_objects", "1", "name")
	assert.NoError(t, err, "Error getting array_of_objects[1].name")
	assert.Equal(t, "Jane Smith", name2, "array_of_objects[1].name should equal Jane Smith")

	active2, err := jsonparser.GetBool(arrayTestJson, "array_of_objects", "1", "active")
	assert.NoError(t, err, "Error getting array_of_objects[1].active")
	assert.Equal(t, false, active2, "array_of_objects[1].active should equal false")
}

// Test objects with arrays
func TestObjectsWithArrays(t *testing.T) {
	// Test Product A tags
	tag, err := jsonparser.GetString(arrayTestJson, "objects_with_arrays", "0", "tags", "1")
	assert.NoError(t, err, "Error getting objects_with_arrays[0].tags[1]")
	assert.Equal(t, "mobile", tag, "objects_with_arrays[0].tags[1] should equal mobile")

	// Test Product A prices
	price, err := jsonparser.GetFloat64(arrayTestJson, "objects_with_arrays", "0", "prices", "0")
	assert.NoError(t, err, "Error getting objects_with_arrays[0].prices[0]")
	assert.Equal(t, 299.99, price, "objects_with_arrays[0].prices[0] should equal 299.99")

	// Test nested features colors
	color, err := jsonparser.GetString(arrayTestJson, "objects_with_arrays", "0", "features", "colors", "2")
	assert.NoError(t, err, "Error getting objects_with_arrays[0].features.colors[2]")
	assert.Equal(t, "red", color, "objects_with_arrays[0].features.colors[2] should equal red")

	// Test nested features storage
	storage, err := jsonparser.GetInt(arrayTestJson, "objects_with_arrays", "0", "features", "storage", "1")
	assert.NoError(t, err, "Error getting objects_with_arrays[0].features.storage[1]")
	assert.Equal(t, 128, storage, "objects_with_arrays[0].features.storage[1] should equal 128")
}

// Test deeply nested structure
func TestDeeplyNested(t *testing.T) {
	result, err := jsonparser.GetString(arrayTestJson, "deeply_nested", "0", "level1", "0", "level2", "0", "level3", "1", "data")
	assert.NoError(t, err, "Error getting deeply nested value")
	assert.Equal(t, "deep value 2", result, "deeply nested value should equal 'deep value 2'")
}

// Test large array
func TestLargeArray(t *testing.T) {
	tests := []struct {
		index    string
		expected string
	}{
		{"0", "item_1"},
		{"9", "item_10"},
		{"19", "item_20"},
	}

	for _, test := range tests {
		result, err := jsonparser.GetString(arrayTestJson, "large_array", test.index)
		assert.NoError(t, err, "Error getting large_array[%s]", test.index)
		assert.Equal(t, test.expected, result, "large_array[%s] should equal %s", test.index, test.expected)
	}
}

// Test unicode array
func TestUnicodeArray(t *testing.T) {
	tests := []struct {
		index    string
		expected string
	}{
		{"0", "café"},
		{"1", "naïve"},
		{"2", "résumé"},
		{"3", "🌟"},
		{"4", "🚀"},
		{"5", "💡"},
	}

	for _, test := range tests {
		result, err := jsonparser.GetString(arrayTestJson, "unicode_array", test.index)
		assert.NoError(t, err, "Error getting unicode_array[%s]", test.index)
		assert.Equal(t, test.expected, result, "unicode_array[%s] should equal %s", test.index, test.expected)
	}
}

// Test special characters array
func TestSpecialCharactersArray(t *testing.T) {
	tests := []struct {
		index    string
		expected string
	}{
		{"0", "line\\nbreak"},
		{"1", "tab\\there"},
		{"2", "quote\\\"test"},
		{"3", "backslash\\\\test"},
	}

	for _, test := range tests {
		result, err := jsonparser.GetString(arrayTestJson, "special_characters", test.index)
		assert.NoError(t, err, "Error getting special_characters[%s]", test.index)
		assert.Equal(t, test.expected, result, "special_characters[%s] should equal %s", test.index, test.expected)
	}
}

// Test single element arrays
func TestSingleElementArrays(t *testing.T) {
	// Test single string
	result, err := jsonparser.GetString(arrayTestJson, "single_element_arrays", "one_string", "0")
	assert.NoError(t, err, "Error getting single string")
	assert.Equal(t, "single", result, "single string should equal 'single'")

	// Test single number
	numberResult, err := jsonparser.GetInt(arrayTestJson, "single_element_arrays", "one_number", "0")
	assert.NoError(t, err, "Error getting single number")
	assert.Equal(t, 42, numberResult, "single number should equal 42")

	// Test single boolean
	boolResult, err := jsonparser.GetBool(arrayTestJson, "single_element_arrays", "one_boolean", "0")
	assert.NoError(t, err, "Error getting single boolean")
	assert.Equal(t, true, boolResult, "single boolean should equal true")

	// Test single null
	nullResult, err := jsonparser.GetString(arrayTestJson, "single_element_arrays", "one_null", "0")
	assert.NoError(t, err, "Error getting single null")
	assert.Equal(t, "null", nullResult, "single null should equal 'null'")

	// Test single object
	objResult, err := jsonparser.GetString(arrayTestJson, "single_element_arrays", "one_object", "0", "key")
	assert.NoError(t, err, "Error getting single object")
	assert.Equal(t, "value", objResult, "single object value should equal 'value'")

	// Test single nested array
	nestedResult, err := jsonparser.GetInt(arrayTestJson, "single_element_arrays", "one_array", "0", "1")
	assert.NoError(t, err, "Error getting single nested array")
	assert.Equal(t, 2, nestedResult, "single nested array element should equal 2")
}

// Test arrays with duplicates
func TestArraysWithDuplicates(t *testing.T) {
	// Test duplicate strings
	result1, err := jsonparser.GetString(arrayTestJson, "arrays_with_duplicates", "duplicate_strings", "0")
	assert.NoError(t, err, "Error getting duplicate_strings[0]")
	result2, err := jsonparser.GetString(arrayTestJson, "arrays_with_duplicates", "duplicate_strings", "2")
	assert.NoError(t, err, "Error getting duplicate_strings[2]")
	assert.Equal(t, "apple", result1, "duplicate_strings[0] should equal 'apple'")
	assert.Equal(t, "apple", result2, "duplicate_strings[2] should equal 'apple'")

	// Test duplicate numbers
	num1, err := jsonparser.GetInt(arrayTestJson, "arrays_with_duplicates", "duplicate_numbers", "0")
	assert.NoError(t, err, "Error getting duplicate_numbers[0]")
	num2, err := jsonparser.GetInt(arrayTestJson, "arrays_with_duplicates", "duplicate_numbers", "3")
	assert.NoError(t, err, "Error getting duplicate_numbers[3]")
	assert.Equal(t, 1, num1, "duplicate_numbers[0] should equal 1")
	assert.Equal(t, 1, num2, "duplicate_numbers[3] should equal 1")
}

// Test edge cases
func TestEdgeCases(t *testing.T) {
	// Test array with empty strings
	result, err := jsonparser.GetString(arrayTestJson, "edge_cases", "array_with_empty_strings", "0")
	assert.NoError(t, err, "Error getting empty string")
	assert.Equal(t, "", result, "empty string should be empty")

	// Test array with zeros
	zero, err := jsonparser.GetInt(arrayTestJson, "edge_cases", "array_with_zeros", "0")
	assert.NoError(t, err, "Error getting zero")
	assert.Equal(t, 0, zero, "zero should equal 0")

	// Test negative numbers
	negative, err := jsonparser.GetInt(arrayTestJson, "edge_cases", "array_with_negative_numbers", "0")
	assert.NoError(t, err, "Error getting negative number")
	assert.Equal(t, -1, negative, "negative number should equal -1")

	// Test very long string
	longString, err := jsonparser.GetString(arrayTestJson, "edge_cases", "array_with_very_long_string", "0")
	assert.NoError(t, err, "Error getting long string")
	expectedStart := "This is a very long string"
	assert.GreaterOrEqual(t, len(longString), len(expectedStart), "long string should be long enough")
	assert.Equal(t, expectedStart, longString[:len(expectedStart)], "long string should start with expected text")

	// Test scientific notation
	scientific, err := jsonparser.GetFloat64(arrayTestJson, "edge_cases", "array_with_scientific_notation", "0")
	assert.NoError(t, err, "Error getting scientific notation")
	assert.Equal(t, 1e10, scientific, "scientific notation should equal 1e10")
}

// Benchmark tests for array operations
func BenchmarkStringArrayAccess(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := jsonparser.GetString(arrayTestJson, "string_array", "1")
		assert.NoError(b, err)
	}
}

func BenchmarkNestedArrayAccess(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := jsonparser.GetInt(arrayTestJson, "nested_arrays", "0", "1")
		assert.NoError(b, err)
	}
}

func BenchmarkDeeplyNestedAccess(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := jsonparser.GetString(arrayTestJson, "deeply_nested", "0", "level1", "0", "level2", "0", "level3", "1", "data")
		assert.NoError(b, err)
	}
}

func BenchmarkArrayOfObjectsAccess(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := jsonparser.GetString(arrayTestJson, "array_of_objects", "0", "name")
		assert.NoError(b, err)
	}
}

// ===== Tests for ForeachArrayElement =====
func TestForeachArrayElement_StringArray(t *testing.T) {
	expected := []string{"apple", "banana", "cherry", "date"}
	results := []string{}

	err := jsonparser.ForeachArrayElement(arrayTestJson, func(valueSlice []byte, index int) {
		results = append(results, string(valueSlice))
	}, "string_array")

	assert.NoError(t, err)
	assert.Equal(t, len(expected), len(results), "Should have correct number of elements")
	for i, exp := range expected {
		assert.Equal(t, exp, results[i], "Element at index %d should match", i)
	}
}

// Test ForeachArrayElement with number array
func TestForeachArrayElement_NumberArray(t *testing.T) {
	expected := []string{"1", "2", "3", "4", "5", "42", "100"}
	results := []string{}

	err := jsonparser.ForeachArrayElement(arrayTestJson, func(valueSlice []byte, index int) {
		results = append(results, string(valueSlice))
	}, "number_array")

	assert.NoError(t, err)
	assert.Equal(t, len(expected), len(results), "Should have correct number of elements")
	for i, exp := range expected {
		assert.Equal(t, exp, results[i], "Element at index %d should match", i)
	}
}

// Test ForeachArrayElement with float array
func TestForeachArrayElement_FloatArray(t *testing.T) {
	expected := []string{"1.1", "2.5", "3.14", "4.0", "5.999"}
	results := []string{}

	err := jsonparser.ForeachArrayElement(arrayTestJson, func(valueSlice []byte, index int) {
		results = append(results, string(valueSlice))
	}, "float_array")

	assert.NoError(t, err)
	assert.Equal(t, len(expected), len(results), "Should have correct number of elements")
	for i, exp := range expected {
		assert.Equal(t, exp, results[i], "Element at index %d should match", i)
	}
}

// Test ForeachArrayElement with boolean array
func TestForeachArrayElement_BooleanArray(t *testing.T) {
	expected := []string{"true", "false", "true", "true", "false"}
	results := []string{}

	err := jsonparser.ForeachArrayElement(arrayTestJson, func(valueSlice []byte, index int) {
		results = append(results, string(valueSlice))
	}, "boolean_array")

	assert.NoError(t, err)
	assert.Equal(t, len(expected), len(results), "Should have correct number of elements")
	for i, exp := range expected {
		assert.Equal(t, exp, results[i], "Element at index %d should match", i)
	}
}

// Test ForeachArrayElement with mixed array
func TestForeachArrayElement_MixedArray(t *testing.T) {
	expected := []string{"1", "hello", "true", "3.14", "null", "false", "world"}
	results := []string{}

	err := jsonparser.ForeachArrayElement(arrayTestJson, func(valueSlice []byte, index int) {
		results = append(results, string(valueSlice))
	}, "mixed_array")

	assert.NoError(t, err)
	assert.Equal(t, len(expected), len(results), "Should have correct number of elements")
	for i, exp := range expected {
		assert.Equal(t, exp, results[i], "Element at index %d should match", i)
	}
}

// Test ForeachArrayElement with empty array
func TestForeachArrayElement_EmptyArray(t *testing.T) {
	callCount := 0

	err := jsonparser.ForeachArrayElement(arrayTestJson, func(valueSlice []byte, index int) {
		callCount++
	}, "empty_array")

	assert.NoError(t, err)
	assert.Equal(t, 0, callCount, "Callback should not be called for empty array")
}

// Test ForeachArrayElement with null array
func TestForeachArrayElement_NullArray(t *testing.T) {
	expected := []string{"null", "null", "null"}
	results := []string{}

	err := jsonparser.ForeachArrayElement(arrayTestJson, func(valueSlice []byte, index int) {
		results = append(results, string(valueSlice))
	}, "null_array")

	assert.NoError(t, err)
	assert.Equal(t, len(expected), len(results), "Should have correct number of elements")
	for i, exp := range expected {
		assert.Equal(t, exp, results[i], "Element at index %d should match", i)
	}
}

// Test ForeachArrayElement with nested arrays - iterating outer array
func TestForeachArrayElement_NestedArrays(t *testing.T) {
	results := []string{}

	err := jsonparser.ForeachArrayElement(arrayTestJson, func(valueSlice []byte, index int) {
		// Remove whitespace for comparison
		normalized := string(valueSlice)
		results = append(results, normalized)
	}, "nested_arrays")

	assert.NoError(t, err)
	assert.Equal(t, 4, len(results), "Should have 4 nested arrays")
	// Verify each result is an array
	for i, result := range results {
		assert.True(t, len(result) > 0 && result[0] == '[', "Element at index %d should be an array", i)
	}
}

// Test ForeachArrayElement with nested arrays - iterating inner array
func TestForeachArrayElement_NestedArraysInner(t *testing.T) {
	expected := []string{"1", "2", "3"}
	results := []string{}

	err := jsonparser.ForeachArrayElement(arrayTestJson, func(valueSlice []byte, index int) {
		results = append(results, string(valueSlice))
	}, "nested_arrays", "0")

	assert.NoError(t, err)
	assert.Equal(t, len(expected), len(results), "Should have correct number of elements")
	for i, exp := range expected {
		assert.Equal(t, exp, results[i], "Element at index %d should match", i)
	}
}

// Test ForeachArrayElement with array of objects
func TestForeachArrayElement_ArrayOfObjects(t *testing.T) {
	results := []string{}
	expectedCount := 3

	err := jsonparser.ForeachArrayElement(arrayTestJson, func(valueSlice []byte, index int) {
		results = append(results, string(valueSlice))
	}, "array_of_objects")

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, len(results), "Should have 3 objects")

	// Each result should be a JSON object
	for i, result := range results {
		assert.Contains(t, result, "id", "Object at index %d should contain 'id'", i)
		assert.Contains(t, result, "name", "Object at index %d should contain 'name'", i)
		assert.Contains(t, result, "email", "Object at index %d should contain 'email'", i)
	}
}

// Test ForeachArrayElement with objects containing arrays
func TestForeachArrayElement_ObjectsWithArrays(t *testing.T) {
	results := []string{}
	expectedCount := 2

	err := jsonparser.ForeachArrayElement(arrayTestJson, func(valueSlice []byte, index int) {
		results = append(results, string(valueSlice))
	}, "objects_with_arrays")

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, len(results), "Should have 2 objects")
}

// Test ForeachArrayElement accessing nested array within object
func TestForeachArrayElement_NestedArrayInObject(t *testing.T) {
	expected := []string{"electronics", "mobile", "smartphone"}
	results := []string{}

	err := jsonparser.ForeachArrayElement(arrayTestJson, func(valueSlice []byte, index int) {
		results = append(results, string(valueSlice))
	}, "objects_with_arrays", "0", "tags")

	assert.NoError(t, err)
	assert.Equal(t, len(expected), len(results), "Should have correct number of tags")
	for i, exp := range expected {
		assert.Equal(t, exp, results[i], "Tag at index %d should match", i)
	}
}

// Benchmark ForeachArrayElement
func BenchmarkForeachArrayElement_StringArray(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = jsonparser.ForeachArrayElement(arrayTestJson, func(valueSlice []byte, index int) {
			// Do minimal work
			_ = valueSlice
		}, "string_array")
	}
}

func BenchmarkForeachArrayElement_LargeArray(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = jsonparser.ForeachArrayElement(arrayTestJson, func(valueSlice []byte, index int) {
			// Do minimal work
			_ = valueSlice
		}, "large_array")
	}
}

func BenchmarkForeachArrayElement_ArrayOfObjects(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = jsonparser.ForeachArrayElement(arrayTestJson, func(valueSlice []byte, index int) {
			// Do minimal work
			_ = valueSlice
		}, "array_of_objects")
	}
}
