package jsonparser_test

import (
	"os"
	"testing"

	"github.com/muccarini/jsonparser"
	"github.com/stretchr/testify/assert"
)

var primitivesTestJson []byte

func init() {
	// Load the sample JSON file for testing
	var err error
	primitivesTestJson, err = os.ReadFile("sample_primitives.json")
	if err != nil {
		panic("Failed to load sample_primitives.json: " + err.Error())
	}
}

func TestGetString_Primitives(t *testing.T) {
	tests := []struct {
		name     string
		fields   []string
		expected string
	}{
		{"simple string", []string{"stringValue"}, "Hello, World!"},
		{"empty string", []string{"emptyString"}, ""},
		{"unicode string", []string{"unicodeString"}, "Hello 世界 🌍 Testing unicode characters"},
		{"string with escapes", []string{"stringWithEscapes"}, "Line 1\\nLine 2\\tTab\\r\\nCarriage return and quotes: \\\"Hello\\\""},
		{"nested string", []string{"nested", "deepString"}, "Nested string value"},
		{"deep nested string", []string{"nested", "level2", "level3", "extremelyDeepString"}, "Extremely deep string for performance testing"},
		{"array element", []string{"arrayOfStrings", "1"}, "second"},
		{"many fields", []string{"performance", "manyFields", "field5"}, "value5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := jsonparser.GetString(primitivesTestJson, tt.fields...)
			assert.NoError(t, err, "Error getting %v", tt.fields)
			assert.Equal(t, tt.expected, result, "%s should equal %s", tt.fields, tt.expected)
		})
	}
}

func TestGetBool_Primitives(t *testing.T) {
	tests := []struct {
		name     string
		fields   []string
		expected bool
	}{
		{"bool true", []string{"boolTrue"}, true},
		{"bool false", []string{"boolFalse"}, false},
		{"nested bool", []string{"nested", "deepBool"}, true},
		{"deep nested bool", []string{"nested", "level2", "veryDeepBool"}, false},
		{"array element true", []string{"arrayOfBools", "0"}, true},
		{"array element false", []string{"arrayOfBools", "1"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := jsonparser.GetBool(primitivesTestJson, tt.fields...)
			assert.NoError(t, err, "Error getting %v", tt.fields)
			assert.Equal(t, tt.expected, result, "%s should equal %t", tt.fields, tt.expected)
		})
	}
}

func TestGetInt_Primitives(t *testing.T) {
	tests := []struct {
		name     string
		fields   []string
		expected int
	}{
		{"int zero", []string{"intZero"}, 0},
		{"int positive", []string{"intPositive"}, 42},
		{"int negative", []string{"intNegative"}, -123},
		{"int large", []string{"intLarge"}, 2147483647},
		{"int minimum", []string{"intMinimum"}, -2147483648},
		{"nested int", []string{"nested", "deepInt"}, 999},
		{"deep nested int", []string{"nested", "level2", "level3", "extremelyDeepInt"}, 555},
		{"array element", []string{"arrayOfInts", "2"}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := jsonparser.GetInt(primitivesTestJson, tt.fields...)
			assert.NoError(t, err, "Error getting %v", tt.fields)
			assert.Equal(t, tt.expected, result, "%s should equal %d", tt.fields, tt.expected)
		})
	}
}

func TestGetInt64_Primitives(t *testing.T) {
	tests := []struct {
		name     string
		fields   []string
		expected int64
	}{
		{"int64 large", []string{"int64Large"}, 9223372036854775807},
		{"int64 minimum", []string{"int64Minimum"}, -9223372036854775808},
		{"nested int64", []string{"nested", "level2", "level3", "extremelyDeepInt"}, 555},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := jsonparser.GetInt64(primitivesTestJson, tt.fields...)
			assert.NoError(t, err, "Error getting %v", tt.fields)
			assert.Equal(t, tt.expected, result, "%s should equal %d", tt.fields, tt.expected)
		})
	}
}

func TestGetFloat64_Primitives(t *testing.T) {
	tests := []struct {
		name     string
		fields   []string
		expected float64
		bitSize  int
	}{
		{"float zero", []string{"floatZero"}, 0.0, 64},
		{"float positive", []string{"floatPositive"}, 3.14159, 64},
		{"float negative", []string{"floatNegative"}, -2.71828, 64},
		{"float scientific", []string{"floatScientific"}, 1.23e-4, 64},
		{"nested float", []string{"nested", "deepFloat"}, 123.456, 64},
		{"deep nested float", []string{"nested", "level2", "level3", "extremelyDeepFloat"}, 111.222, 64},
		{"array element", []string{"arrayOfFloats", "1"}, 2.2, 64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := jsonparser.GetFloat64(primitivesTestJson, tt.fields...)
			assert.NoError(t, err, "Error getting %v", tt.fields)
			assert.Equal(t, tt.expected, result, "%s should equal %f", tt.fields, tt.expected)
		})
	}
}

func TestGetString_ErrorCases(t *testing.T) {
	tests := []struct {
		name   string
		fields []string
	}{
		{"non-existent field", []string{"nonExistent"}},
		{"wrong path", []string{"nested", "wrongField"}},
		{"array out of bounds", []string{"arrayOfStrings", "100"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := jsonparser.GetString(primitivesTestJson, tt.fields...)
			assert.Error(t, err, "GetString() expected error for %v, but got none", tt.fields)
		})
	}
}
