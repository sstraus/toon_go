package toon

import (
	"math"
	"reflect"
	"testing"
)

// TestToFloat64 tests the toFloat64 function with all numeric types
func TestToFloat64(t *testing.T) {
	tests := []struct {
		name     string
		input    Value
		expected float64
		ok       bool
	}{
		// float64 cases
		{"float64 zero", float64(0), 0, true},
		{"float64 positive", float64(42.5), 42.5, true},
		{"float64 negative", float64(-42.5), -42.5, true},
		{"float64 max", math.MaxFloat64, math.MaxFloat64, true},
		{"float64 small", float64(0.0000001), 0.0000001, true},

		// float32 cases
		{"float32 zero", float32(0), 0, true},
		{"float32 positive", float32(42.5), 42.5, true},
		{"float32 negative", float32(-42.5), -42.5, true},
		{"float32 max", float32(math.MaxFloat32), float64(float32(math.MaxFloat32)), true},

		// int cases
		{"int zero", int(0), 0, true},
		{"int positive", int(42), 42, true},
		{"int negative", int(-42), -42, true},
		{"int max", int(math.MaxInt32), float64(math.MaxInt32), true},

		// int8 cases
		{"int8 zero", int8(0), 0, true},
		{"int8 max", int8(127), 127, true},
		{"int8 min", int8(-128), -128, true},

		// int16 cases
		{"int16 zero", int16(0), 0, true},
		{"int16 max", int16(32767), 32767, true},
		{"int16 min", int16(-32768), -32768, true},

		// int32 cases
		{"int32 zero", int32(0), 0, true},
		{"int32 max", int32(math.MaxInt32), float64(math.MaxInt32), true},
		{"int32 min", int32(math.MinInt32), float64(math.MinInt32), true},

		// int64 cases
		{"int64 zero", int64(0), 0, true},
		{"int64 max", int64(math.MaxInt64), float64(math.MaxInt64), true},
		{"int64 min", int64(math.MinInt64), float64(math.MinInt64), true},

		// uint cases
		{"uint zero", uint(0), 0, true},
		{"uint positive", uint(42), 42, true},
		{"uint max", uint(math.MaxUint32), float64(math.MaxUint32), true},

		// uint8 cases
		{"uint8 zero", uint8(0), 0, true},
		{"uint8 max", uint8(255), 255, true},

		// uint16 cases
		{"uint16 zero", uint16(0), 0, true},
		{"uint16 max", uint16(65535), 65535, true},

		// uint32 cases
		{"uint32 zero", uint32(0), 0, true},
		{"uint32 max", uint32(math.MaxUint32), float64(math.MaxUint32), true},

		// uint64 cases
		{"uint64 zero", uint64(0), 0, true},
		{"uint64 max", uint64(math.MaxUint64), float64(math.MaxUint64), true},

		// Invalid types
		{"string", "42.5", 0, false},
		{"bool", true, 0, false},
		{"nil", nil, 0, false},
		{"slice", []int{1, 2, 3}, 0, false},
		{"map", map[string]int{"a": 1}, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := toFloat64(tt.input)
			if ok != tt.ok {
				t.Errorf("toFloat64(%v) ok = %v, want %v", tt.input, ok, tt.ok)
			}
			if ok && result != tt.expected {
				t.Errorf("toFloat64(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestToInt64 tests the toInt64 function with all numeric types
func TestToInt64(t *testing.T) {
	tests := []struct {
		name     string
		input    Value
		expected int64
		ok       bool
	}{
		// int cases
		{"int zero", int(0), 0, true},
		{"int positive", int(42), 42, true},
		{"int negative", int(-42), -42, true},

		// int8 cases
		{"int8 zero", int8(0), 0, true},
		{"int8 max", int8(127), 127, true},
		{"int8 min", int8(-128), -128, true},

		// int16 cases
		{"int16 zero", int16(0), 0, true},
		{"int16 max", int16(32767), 32767, true},
		{"int16 min", int16(-32768), -32768, true},

		// int32 cases
		{"int32 zero", int32(0), 0, true},
		{"int32 max", int32(math.MaxInt32), int64(math.MaxInt32), true},
		{"int32 min", int32(math.MinInt32), int64(math.MinInt32), true},

		// int64 cases
		{"int64 zero", int64(0), 0, true},
		{"int64 max", int64(math.MaxInt64), math.MaxInt64, true},
		{"int64 min", int64(math.MinInt64), math.MinInt64, true},

		// uint cases
		{"uint zero", uint(0), 0, true},
		{"uint positive", uint(42), 42, true},

		// uint8 cases
		{"uint8 zero", uint8(0), 0, true},
		{"uint8 max", uint8(255), 255, true},

		// uint16 cases
		{"uint16 zero", uint16(0), 0, true},
		{"uint16 max", uint16(65535), 65535, true},

		// uint32 cases
		{"uint32 zero", uint32(0), 0, true},
		{"uint32 max", uint32(math.MaxUint32), int64(math.MaxUint32), true},

		// uint64 cases
		{"uint64 zero", uint64(0), 0, true},
		{"uint64 small", uint64(42), 42, true},
		{"uint64 max valid", uint64(math.MaxInt64), math.MaxInt64, true},
		{"uint64 overflow", uint64(math.MaxUint64), 0, false},

		// float32 cases
		{"float32 whole number", float32(42.0), 42, true},
		{"float32 zero", float32(0.0), 0, true},
		{"float32 negative", float32(-42.0), -42, true},
		{"float32 fractional", float32(42.5), 0, false},

		// float64 cases
		{"float64 whole number", float64(42.0), 42, true},
		{"float64 zero", float64(0.0), 0, true},
		{"float64 negative", float64(-42.0), -42, true},
		{"float64 fractional", float64(42.5), 0, false},
		{"float64 large whole", float64(1e15), int64(1e15), true},

		// Invalid types
		{"string", "42", 0, false},
		{"bool", true, 0, false},
		{"nil", nil, 0, false},
		{"slice", []int{1, 2, 3}, 0, false},
		{"map", map[string]int{"a": 1}, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := toInt64(tt.input)
			if ok != tt.ok {
				t.Errorf("toInt64(%v) ok = %v, want %v", tt.input, ok, tt.ok)
			}
			if ok && result != tt.expected {
				t.Errorf("toInt64(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestNormalize tests the normalize function with various data types
func TestNormalize(t *testing.T) {
	tests := []struct {
		name     string
		input    Value
		expected Value
	}{
		// Primitives
		{"nil", nil, nil},
		{"bool true", true, true},
		{"bool false", false, false},
		{"string", "hello", "hello"},
		{"string empty", "", ""},

		// Integers
		{"int", int(42), int64(42)},
		{"int8", int8(42), int64(42)},
		{"int16", int16(42), int64(42)},
		{"int32", int32(42), int64(42)},
		{"int64", int64(42), int64(42)},
		{"int negative", int(-42), int64(-42)},

		// Unsigned integers
		{"uint", uint(42), int64(42)},
		{"uint8", uint8(42), int64(42)},
		{"uint16", uint16(42), int64(42)},
		{"uint32", uint32(42), int64(42)},
		{"uint64", uint64(42), int64(42)},

		// Floats
		{"float32 whole", float32(42.0), int64(42)},
		{"float32 fractional", float32(42.5), float64(42.5)},
		{"float64 whole", float64(42.0), int64(42)},
		{"float64 fractional", float64(42.5), 42.5},
		{"float64 negative zero", math.Copysign(0, -1), int64(0)},
		{"float64 NaN", math.NaN(), nil},
		{"float64 +Inf", math.Inf(1), nil},
		{"float64 -Inf", math.Inf(-1), nil},

		// Slices
		{"empty slice", []interface{}{}, []Value{}},
		{"slice of ints", []interface{}{1, 2, 3}, []Value{int64(1), int64(2), int64(3)}},
		{"slice of mixed", []interface{}{1, "hello", true}, []Value{int64(1), "hello", true}},
		{"slice nested", []interface{}{[]interface{}{1, 2}, []interface{}{3, 4}},
			[]Value{[]Value{int64(1), int64(2)}, []Value{int64(3), int64(4)}}},

		// Maps
		{"empty map", map[string]interface{}{}, map[string]Value{}},
		{"simple map", map[string]interface{}{"a": 1, "b": "hello"},
			map[string]Value{"a": int64(1), "b": "hello"}},
		{"nested map", map[string]interface{}{"a": map[string]interface{}{"b": 1}},
			map[string]Value{"a": map[string]Value{"b": int64(1)}}},
		{"map with array", map[string]interface{}{"items": []interface{}{1, 2, 3}},
			map[string]Value{"items": []Value{int64(1), int64(2), int64(3)}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalize(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("normalize(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestNormalizeOrderedMap tests normalize with OrderedMap types
func TestNormalizeOrderedMap(t *testing.T) {
	// Test OrderedMap value
	om := NewOrderedMap()
	om.Set("a", 1)
	om.Set("b", "hello")
	om.Set("c", []interface{}{1, 2, 3})

	result := normalize(*om)
	resultOM, ok := result.(OrderedMap)
	if !ok {
		t.Fatalf("normalize(OrderedMap) should return OrderedMap, got %T", result)
	}

	// Check normalized values
	if v, ok := resultOM.Get("a"); !ok || v != int64(1) {
		t.Errorf("normalized OrderedMap 'a' = %v, want int64(1)", v)
	}
	if v, ok := resultOM.Get("b"); !ok || v != "hello" {
		t.Errorf("normalized OrderedMap 'b' = %v, want 'hello'", v)
	}
	if v, ok := resultOM.Get("c"); ok {
		if arr, ok := v.([]Value); ok {
			expected := []Value{int64(1), int64(2), int64(3)}
			if !reflect.DeepEqual(arr, expected) {
				t.Errorf("normalized OrderedMap 'c' = %v, want %v", arr, expected)
			}
		} else {
			t.Errorf("normalized OrderedMap 'c' should be []Value, got %T", v)
		}
	} else {
		t.Error("normalized OrderedMap should have key 'c'")
	}

	// Test OrderedMap pointer
	resultPtr := normalize(om)
	resultOMPtr, ok := resultPtr.(OrderedMap)
	if !ok {
		t.Fatalf("normalize(*OrderedMap) should return OrderedMap, got %T", resultPtr)
	}
	if v, ok := resultOMPtr.Get("a"); !ok || v != int64(1) {
		t.Errorf("normalized *OrderedMap 'a' = %v, want int64(1)", v)
	}
}

// TestNormalizeReflection tests normalize with reflection-based types
func TestNormalizeReflection(t *testing.T) {
	// Test generic slice types
	intSlice := []int{1, 2, 3}
	result := normalize(intSlice)
	expected := []Value{int64(1), int64(2), int64(3)}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("normalize([]int) = %v, want %v", result, expected)
	}

	// Test generic map types
	stringMap := map[string]int{"a": 1, "b": 2}
	result = normalize(stringMap)
	expectedMap := map[string]Value{"a": int64(1), "b": int64(2)}
	if !reflect.DeepEqual(result, expectedMap) {
		t.Errorf("normalize(map[string]int) = %v, want %v", result, expectedMap)
	}

	// Test unsupported type
	type customStruct struct {
		Field string
	}
	result = normalize(customStruct{Field: "test"})
	if result != nil {
		t.Errorf("normalize(customStruct) = %v, want nil", result)
	}
}

// TestSortKeysWithArraysFirst tests the sortKeysWithArraysFirst function
func TestSortKeysWithArraysFirst(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected []string
	}{
		{
			name:     "empty map",
			input:    map[string]interface{}{},
			expected: []string{},
		},
		{
			name: "only arrays",
			input: map[string]interface{}{
				"items": []int{1, 2, 3},
				"data":  []string{"a", "b"},
			},
			expected: []string{"data", "items"},
		},
		{
			name: "only non-arrays",
			input: map[string]interface{}{
				"name": "John",
				"age":  30,
			},
			expected: []string{"age", "name"},
		},
		{
			name: "mixed with arrays first",
			input: map[string]interface{}{
				"name":  "John",
				"items": []int{1, 2, 3},
				"age":   30,
				"tags":  []string{"a", "b"},
			},
			expected: []string{"items", "tags", "age", "name"},
		},
		{
			name: "single array",
			input: map[string]interface{}{
				"data": []int{1},
			},
			expected: []string{"data"},
		},
		{
			name: "single non-array",
			input: map[string]interface{}{
				"name": "test",
			},
			expected: []string{"name"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a map and get its keys
			keys := []string{}
			for k := range tt.input {
				keys = append(keys, k)
			}

			// Sort using the function
			rv := reflect.ValueOf(tt.input)
			sortKeysWithArraysFirst(keys, rv)

			// Compare results
			if !reflect.DeepEqual(keys, tt.expected) {
				t.Errorf("sortKeysWithArraysFirst() = %v, want %v", keys, tt.expected)
			}
		})
	}
}

// TestSameKeys tests the sameKeys function
func TestSameKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    Value
		expected bool
	}{
		{"not a list", "string", false},
		{"empty list", []interface{}{}, true},
		{"single map", []interface{}{map[string]interface{}{"a": 1}}, true},
		{"same keys", []interface{}{
			map[string]interface{}{"a": 1, "b": 2},
			map[string]interface{}{"a": 3, "b": 4},
		}, true},
		{"different keys", []interface{}{
			map[string]interface{}{"a": 1, "b": 2},
			map[string]interface{}{"a": 3, "c": 4},
		}, false},
		{"different key count", []interface{}{
			map[string]interface{}{"a": 1},
			map[string]interface{}{"a": 1, "b": 2},
		}, false},
		{"not all maps", []interface{}{
			map[string]interface{}{"a": 1},
			"string",
		}, false},
		{"first not map", []interface{}{"string"}, false},
		{"empty maps", []interface{}{
			map[string]interface{}{},
			map[string]interface{}{},
		}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sameKeys(tt.input)
			if result != tt.expected {
				t.Errorf("sameKeys(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestSameKeysWithOrderedMap tests sameKeys with OrderedMap
func TestSameKeysWithOrderedMap(t *testing.T) {
	om1 := NewOrderedMap()
	om1.Set("a", 1)
	om1.Set("b", 2)

	om2 := NewOrderedMap()
	om2.Set("b", 3)
	om2.Set("a", 4)

	// Same keys, different order
	input := []interface{}{om1, om2}
	if !sameKeys(input) {
		t.Error("sameKeys with OrderedMaps should return true for same keys")
	}

	om3 := NewOrderedMap()
	om3.Set("a", 1)
	om3.Set("c", 2)

	// Different keys
	input = []interface{}{om1, om3}
	if sameKeys(input) {
		t.Error("sameKeys with OrderedMaps should return false for different keys")
	}
}

// TestGetMapKeys tests the getMapKeys function
func TestGetMapKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    Value
		expected []string
	}{
		{"nil", nil, nil},
		{"not a map", "string", nil},
		{"empty map", map[string]interface{}{}, []string{}},
		{"simple map", map[string]interface{}{"c": 1, "a": 2, "b": 3}, []string{"a", "b", "c"}},
		{"single key", map[string]interface{}{"key": "value"}, []string{"key"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getMapKeys(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("getMapKeys(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestGetMapKeysOrderedMap tests getMapKeys with OrderedMap
func TestGetMapKeysOrderedMap(t *testing.T) {
	// Test with OrderedMap value
	om := NewOrderedMap()
	om.Set("z", 1)
	om.Set("a", 2)
	om.Set("m", 3)

	result := getMapKeys(*om)
	expected := []string{"a", "m", "z"} // Should be sorted
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("getMapKeys(OrderedMap) = %v, want %v", result, expected)
	}

	// Test with OrderedMap pointer
	result = getMapKeys(om)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("getMapKeys(*OrderedMap) = %v, want %v", result, expected)
	}
}

// TestIsMap tests the isMap function with various types
func TestIsMap(t *testing.T) {
	tests := []struct {
		name     string
		input    Value
		expected bool
	}{
		{"nil", nil, false},
		{"string", "test", false},
		{"int", 42, false},
		{"bool", true, false},
		{"slice", []int{1, 2, 3}, false},
		{"map[string]interface{}", map[string]interface{}{"a": 1}, true},
		{"map[string]int", map[string]int{"a": 1}, true},
		{"map[string]string", map[string]string{"a": "b"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isMap(tt.input)
			if result != tt.expected {
				t.Errorf("isMap(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIsMapOrderedMap tests isMap with OrderedMap types
func TestIsMapOrderedMap(t *testing.T) {
	om := NewOrderedMap()
	om.Set("a", 1)

	// Test OrderedMap value
	if !isMap(*om) {
		t.Error("isMap(OrderedMap) should return true")
	}

	// Test OrderedMap pointer
	if !isMap(om) {
		t.Error("isMap(*OrderedMap) should return true")
	}
}

// TestIsList tests the isList function with various types
func TestIsList(t *testing.T) {
	tests := []struct {
		name     string
		input    Value
		expected bool
	}{
		{"nil", nil, false},
		{"string", "test", false},
		{"int", 42, false},
		{"bool", true, false},
		{"map", map[string]interface{}{"a": 1}, false},
		{"[]int", []int{1, 2, 3}, true},
		{"[]string", []string{"a", "b"}, true},
		{"[]interface{}", []interface{}{1, "a", true}, true},
		{"empty slice", []int{}, true},
		{"array", [3]int{1, 2, 3}, true},
		{"empty array", [0]int{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isList(tt.input)
			if result != tt.expected {
				t.Errorf("isList(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}