package benchmarks

import (
	"encoding/json"
	"os"
	"testing"

	buger "github.com/buger/jsonparser"
	"github.com/muccarini/jsonparser"
)

var comparisonJson []byte

func init() {
	// Load the sample JSON file for benchmarking comparison
	var err error
	comparisonJson, err = os.ReadFile("../sample_primitives.json")
	if err != nil {
		panic("Failed to load sample_primitives.json: " + err.Error())
	}
}

// Benchmark GetString operations
func BenchmarkString_GetString_Mucca(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := jsonparser.GetString(comparisonJson, "stringValue")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkString_Get_Mucca(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var result string
		var err error
		if _, err = jsonparser.Get(&result, comparisonJson, "stringValue"); err != nil {
			b.Error(err)
		}

	}
}

func BenchmarkString_GetString_Buger(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := buger.GetString(comparisonJson, "stringValue")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkString_GetString_Std(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var result map[string]interface{}
		if err := json.Unmarshal(comparisonJson, &result); err != nil {
			b.Error(err)
		}
		_ = result["stringValue"].(string)
	}
}

// Benchmark nested string operations
func BenchmarkNestedString_GetString_Mucca(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := jsonparser.GetString(comparisonJson, "nested", "level2", "level3", "extremelyDeepString")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkNestedString_Get_Mucca(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var result string
		var err error
		if _, err = jsonparser.Get(&result, comparisonJson, "nested", "level2", "level3", "extremelyDeepString"); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkNestedString_GetString_Buger(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := buger.GetString(comparisonJson, "nested", "level2", "level3", "extremelyDeepString")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkNestedString_GetString_Std(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var result map[string]interface{}
		if err := json.Unmarshal(comparisonJson, &result); err != nil {
			b.Error(err)
		}
		nested := result["nested"].(map[string]interface{})
		level2 := nested["level2"].(map[string]interface{})
		level3 := level2["level3"].(map[string]interface{})
		_ = level3["extremelyDeepString"].(string)
	}
}

// Benchmark GetInt operations
func BenchmarkInt_GetInt_Mucca(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := jsonparser.GetInt(comparisonJson, "intPositive")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkInt_Get_Mucca(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var result int64
		var err error
		if _, err = jsonparser.Get(&result, comparisonJson, "intPositive"); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkInt_GetInt_Buger(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := buger.GetInt(comparisonJson, "intPositive")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkInt_GetInt_Std(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var result map[string]interface{}
		if err := json.Unmarshal(comparisonJson, &result); err != nil {
			b.Error(err)
		}
		_ = int(result["intPositive"].(float64))
	}
}

// Benchmark GetFloat operations
func BenchmarkFloat_GetFloat_Mucca(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := jsonparser.GetFloat64(comparisonJson, "floatPositive")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkFloat_Get_Mucca(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var result float64
		var err error
		if _, err = jsonparser.Get(&result, comparisonJson, "floatPositive"); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkFloat_GetFloat_Buger(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := buger.GetFloat(comparisonJson, "floatPositive")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkFloat_GetFloat_Std(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var result map[string]interface{}
		if err := json.Unmarshal(comparisonJson, &result); err != nil {
			b.Error(err)
		}
		_ = result["floatPositive"].(float64)
	}
}

// Benchmark GetBool operations
func BenchmarkBool_GetBool_Mucca(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := jsonparser.GetBool(comparisonJson, "boolTrue")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkBoolean_Get_Mucca(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result bool
		var err error
		if _, err = jsonparser.Get(&result, comparisonJson, "boolTrue"); err != nil {
			b.Error(err)
		}

	}
}

func BenchmarkBool_GetBool_Buger(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := buger.GetBoolean(comparisonJson, "boolTrue")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkBool_GetBool_Std(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var result map[string]interface{}
		if err := json.Unmarshal(comparisonJson, &result); err != nil {
			b.Error(err)
		}
		_ = result["boolTrue"].(bool)
	}
}

// Benchmark array access operations
func BenchmarkArrayAccess_GetString_Mucca(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := jsonparser.GetString(comparisonJson, "arrayOfStrings", "2")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkArrayAccess_Get_Mucca(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result string
		var err error
		if _, err = jsonparser.Get(&result, comparisonJson, "arrayOfStrings", "2"); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkArrayAccess_GetString_Buger(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := buger.GetString(comparisonJson, "arrayOfStrings", "[2]")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkArrayAccess_GetString_Std(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var result map[string]interface{}
		if err := json.Unmarshal(comparisonJson, &result); err != nil {
			b.Error(err)
		}
		arr := result["arrayOfStrings"].([]interface{})
		_ = arr[2].(string)
	}
}

// Benchmark large object field access
func BenchmarkManyFields_GetString_Mucca(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := jsonparser.GetString(comparisonJson, "performance", "manyFields", "field15")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkManyFields_Get_Mucca(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result string
		var err error
		if _, err = jsonparser.Get(&result, comparisonJson, "performance", "manyFields", "field15"); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkManyFields_GetString_Buger(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := buger.GetString(comparisonJson, "performance", "manyFields", "field15")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkManyFields_GetString_Std(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var result map[string]interface{}
		if err := json.Unmarshal(comparisonJson, &result); err != nil {
			b.Error(err)
		}
		perf := result["performance"].(map[string]interface{})
		manyFields := perf["manyFields"].(map[string]interface{})
		_ = manyFields["field15"].(string)
	}
}

// Benchmark array iteration operations
func BenchmarkArrayIteration_ForeachArrayElement_Mucca(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := jsonparser.Foreach(comparisonJson, func(valueSlice []byte, index int) {
			_ = valueSlice
		}, "arrayOfStrings")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkArrayIteration_ArrayEach_Buger(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := buger.ArrayEach(comparisonJson, func(value []byte, dataType buger.ValueType, offset int, err error) {
			_ = value
		}, "arrayOfStrings")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkArrayIteration_ArrayEach_Std(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var result map[string]interface{}
		if err := json.Unmarshal(comparisonJson, &result); err != nil {
			b.Error(err)
		}
		arr := result["arrayOfStrings"].([]interface{})
		for _, v := range arr {
			_ = v.(string)
		}
	}
}
