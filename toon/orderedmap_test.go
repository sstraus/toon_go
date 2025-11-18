package toon

import (
	"encoding/json"
	"sort"
	"strings"
	"testing"
)

// TestPairKeyValue tests the Key() and Value() methods
func TestPairKeyValue(t *testing.T) {
	pair := &Pair{
		key:   "testKey",
		value: "testValue",
	}

	if got := pair.Key(); got != "testKey" {
		t.Errorf("Key() = %v, want %v", got, "testKey")
	}

	if got := pair.Value(); got != "testValue" {
		t.Errorf("Value() = %v, want %v", got, "testValue")
	}
}

// TestByPairSorting tests Len, Swap, and Less methods
func TestByPairSorting(t *testing.T) {
	pairs := []*Pair{
		{key: "z", value: 3},
		{key: "a", value: 1},
		{key: "m", value: 2},
	}

	byPair := ByPair{
		Pairs: pairs,
		LessFunc: func(a *Pair, b *Pair) bool {
			return a.key < b.key
		},
	}

	// Test Len
	if got := byPair.Len(); got != 3 {
		t.Errorf("Len() = %v, want %v", got, 3)
	}

	// Test Swap
	byPair.Swap(0, 2)
	if byPair.Pairs[0].key != "m" || byPair.Pairs[2].key != "z" {
		t.Errorf("Swap failed: got [%v, %v, %v]", byPair.Pairs[0].key, byPair.Pairs[1].key, byPair.Pairs[2].key)
	}

	// Test Less
	if !byPair.Less(1, 0) { // "a" < "m"
		t.Error("Less(1, 0) should be true")
	}
	if byPair.Less(0, 1) { // "m" < "a"
		t.Error("Less(0, 1) should be false")
	}

	// Test full sort
	sort.Sort(byPair)
	expected := []string{"a", "m", "z"}
	for i, pair := range byPair.Pairs {
		if pair.key != expected[i] {
			t.Errorf("After sort, index %d: got %v, want %v", i, pair.key, expected[i])
		}
	}
}

// TestSetEscapeHTML tests the SetEscapeHTML functionality
func TestSetEscapeHTML(t *testing.T) {
	om := NewOrderedMap()
	om.Set("html", "<script>alert('xss')</script>")

	// Test with escapeHTML = true (default)
	data, err := om.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}
	dataStr := string(data)
	// Remove newlines that json.Encoder adds
	dataStr = strings.ReplaceAll(dataStr, "\n", "")
	if !strings.Contains(dataStr, "\\u003c") && !strings.Contains(dataStr, "&lt;") {
		t.Errorf("Expected HTML to be escaped by default, got: %s", dataStr)
	}

	// Test with escapeHTML = false
	om2 := NewOrderedMap()
	om2.SetEscapeHTML(false)
	om2.Set("html", "<script>alert('xss')</script>")
	data, err = om2.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}
	dataStr = string(data)
	dataStr = strings.ReplaceAll(dataStr, "\n", "")
	if !strings.Contains(dataStr, "<script>") {
		t.Errorf("Expected HTML to not be escaped when SetEscapeHTML(false), got: %s", dataStr)
	}
}

// TestDelete tests the Delete method
func TestDelete(t *testing.T) {
	om := NewOrderedMap()
	om.Set("key1", "value1")
	om.Set("key2", "value2")
	om.Set("key3", "value3")

	// Test deleting existing key
	om.Delete("key2")
	if _, exists := om.Get("key2"); exists {
		t.Error("key2 should not exist after Delete")
	}
	if om.Len() != 2 {
		t.Errorf("Len() = %v, want 2", om.Len())
	}
	expectedKeys := []string{"key1", "key3"}
	keys := om.Keys()
	for i, key := range keys {
		if key != expectedKeys[i] {
			t.Errorf("Keys()[%d] = %v, want %v", i, key, expectedKeys[i])
		}
	}

	// Test deleting non-existent key (should not panic)
	om.Delete("nonexistent")
	if om.Len() != 2 {
		t.Errorf("Len() should still be 2 after deleting non-existent key")
	}

	// Test deleting first key
	om.Delete("key1")
	if om.Len() != 1 {
		t.Errorf("Len() = %v, want 1", om.Len())
	}
	if om.Keys()[0] != "key3" {
		t.Errorf("Keys()[0] = %v, want key3", om.Keys()[0])
	}

	// Test deleting last key
	om.Delete("key3")
	if om.Len() != 0 {
		t.Errorf("Len() = %v, want 0", om.Len())
	}
}

// TestSortKeys tests the SortKeys method
func TestSortKeys(t *testing.T) {
	om := NewOrderedMap()
	om.Set("zebra", 1)
	om.Set("apple", 2)
	om.Set("mango", 3)

	om.SortKeys(func(keys []string) {
		sort.Strings(keys)
	})

	expected := []string{"apple", "mango", "zebra"}
	keys := om.Keys()
	for i, key := range keys {
		if key != expected[i] {
			t.Errorf("Keys()[%d] = %v, want %v", i, key, expected[i])
		}
	}
}

// TestSort tests the Sort method
func TestSort(t *testing.T) {
	om := NewOrderedMap()
	om.Set("c", 3)
	om.Set("a", 1)
	om.Set("b", 2)

	// Sort by key
	om.Sort(func(a *Pair, b *Pair) bool {
		return a.key < b.key
	})

	expected := []string{"a", "b", "c"}
	keys := om.Keys()
	for i, key := range keys {
		if key != expected[i] {
			t.Errorf("After Sort by key, Keys()[%d] = %v, want %v", i, key, expected[i])
		}
	}

	// Sort by value
	om.Sort(func(a *Pair, b *Pair) bool {
		return a.value.(int) > b.value.(int) // descending
	})

	expectedDesc := []string{"c", "b", "a"}
	keys = om.Keys()
	for i, key := range keys {
		if key != expectedDesc[i] {
			t.Errorf("After Sort by value desc, Keys()[%d] = %v, want %v", i, key, expectedDesc[i])
		}
	}
}

// TestMarshalJSON tests the MarshalJSON method
func TestMarshalJSON(t *testing.T) {
	om := NewOrderedMap()
	om.Set("first", 1)
	om.Set("second", "two")
	om.Set("third", true)

	data, err := om.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Unmarshaling result failed: %v", err)
	}

	// Verify values
	if result["first"] != float64(1) {
		t.Errorf("first = %v, want 1", result["first"])
	}
	if result["second"] != "two" {
		t.Errorf("second = %v, want two", result["second"])
	}
	if result["third"] != true {
		t.Errorf("third = %v, want true", result["third"])
	}

	// Verify order is preserved (check string representation)
	str := string(data)
	firstPos := strings.Index(str, "first")
	secondPos := strings.Index(str, "second")
	thirdPos := strings.Index(str, "third")
	if firstPos > secondPos || secondPos > thirdPos {
		t.Error("Key order not preserved in JSON output")
	}
}

// TestMarshalJSONNested tests MarshalJSON with nested structures
func TestMarshalJSONNested(t *testing.T) {
	nested := NewOrderedMap()
	nested.Set("inner", "value")

	om := NewOrderedMap()
	om.Set("outer", nested)
	om.Set("number", 42)

	data, err := om.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Unmarshaling result failed: %v", err)
	}

	outer, ok := result["outer"].(map[string]interface{})
	if !ok {
		t.Fatal("outer is not a map")
	}
	if outer["inner"] != "value" {
		t.Errorf("inner = %v, want value", outer["inner"])
	}
}

// TestUnmarshalJSONDuplicateKeys tests decodeOrderedMap with duplicate keys
func TestUnmarshalJSONDuplicateKeys(t *testing.T) {
	// JSON with duplicate keys - last one should win and be at the end
	jsonData := []byte(`{"a": 1, "b": 2, "a": 3}`)

	om := NewOrderedMap()
	err := om.UnmarshalJSON(jsonData)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	// Check value (should be the last occurrence)
	val, exists := om.Get("a")
	if !exists {
		t.Fatal("key 'a' should exist")
	}
	if val != float64(3) {
		t.Errorf("a = %v, want 3", val)
	}

	// Check that 'a' is at the end of keys
	keys := om.Keys()
	if len(keys) != 2 {
		t.Errorf("Len() = %v, want 2", len(keys))
	}
	if keys[1] != "a" {
		t.Errorf("Last key should be 'a', got %v", keys[1])
	}
}

// TestUnmarshalJSONNestedObjects tests decodeOrderedMap with nested objects
func TestUnmarshalJSONNestedObjects(t *testing.T) {
	jsonData := []byte(`{
		"outer": {
			"inner": {
				"deep": "value"
			}
		}
	}`)

	om := NewOrderedMap()
	err := om.UnmarshalJSON(jsonData)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	outer, exists := om.Get("outer")
	if !exists {
		t.Fatal("outer should exist")
	}

	outerMap, ok := outer.(OrderedMap)
	if !ok {
		t.Fatalf("outer should be OrderedMap, got %T", outer)
	}

	inner, exists := outerMap.Get("inner")
	if !exists {
		t.Fatal("inner should exist")
	}

	innerMap, ok := inner.(OrderedMap)
	if !ok {
		t.Fatalf("inner should be OrderedMap, got %T", inner)
	}

	deep, exists := innerMap.Get("deep")
	if !exists {
		t.Fatal("deep should exist")
	}

	if deep != "value" {
		t.Errorf("deep = %v, want value", deep)
	}
}

// TestUnmarshalJSONNestedArrays tests decodeSlice functionality
func TestUnmarshalJSONNestedArrays(t *testing.T) {
	jsonData := []byte(`{
		"array": [
			{"key": "value1"},
			{"key": "value2"}
		]
	}`)

	om := NewOrderedMap()
	err := om.UnmarshalJSON(jsonData)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	array, exists := om.Get("array")
	if !exists {
		t.Fatal("array should exist")
	}

	slice, ok := array.([]interface{})
	if !ok {
		t.Fatalf("array should be []interface{}, got %T", array)
	}

	if len(slice) != 2 {
		t.Fatalf("array length = %v, want 2", len(slice))
	}

	// Check first element
	first, ok := slice[0].(OrderedMap)
	if !ok {
		t.Fatalf("slice[0] should be OrderedMap, got %T", slice[0])
	}
	val, _ := first.Get("key")
	if val != "value1" {
		t.Errorf("slice[0].key = %v, want value1", val)
	}

	// Check second element
	second, ok := slice[1].(OrderedMap)
	if !ok {
		t.Fatalf("slice[1] should be OrderedMap, got %T", slice[1])
	}
	val, _ = second.Get("key")
	if val != "value2" {
		t.Errorf("slice[1].key = %v, want value2", val)
	}
}

// TestDecodeSliceNestedArrays tests nested array handling
func TestDecodeSliceNestedArrays(t *testing.T) {
	jsonData := []byte(`{
		"matrix": [
			[1, 2],
			[3, 4]
		]
	}`)

	om := NewOrderedMap()
	err := om.UnmarshalJSON(jsonData)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	matrix, exists := om.Get("matrix")
	if !exists {
		t.Fatal("matrix should exist")
	}

	outer, ok := matrix.([]interface{})
	if !ok {
		t.Fatalf("matrix should be []interface{}, got %T", matrix)
	}

	if len(outer) != 2 {
		t.Fatalf("outer array length = %v, want 2", len(outer))
	}
}

// TestDecodeSliceMixedTypes tests array with mixed object and array types
func TestDecodeSliceMixedTypes(t *testing.T) {
	jsonData := []byte(`{
		"mixed": [
			{"type": "object"},
			[1, 2, 3],
			"string"
		]
	}`)

	om := NewOrderedMap()
	err := om.UnmarshalJSON(jsonData)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	mixed, exists := om.Get("mixed")
	if !exists {
		t.Fatal("mixed should exist")
	}

	slice, ok := mixed.([]interface{})
	if !ok {
		t.Fatalf("mixed should be []interface{}, got %T", mixed)
	}

	if len(slice) != 3 {
		t.Fatalf("slice length = %v, want 3", len(slice))
	}

	// Check object
	if _, ok := slice[0].(OrderedMap); !ok {
		t.Errorf("slice[0] should be OrderedMap, got %T", slice[0])
	}

	// Check array
	if _, ok := slice[1].([]interface{}); !ok {
		t.Errorf("slice[1] should be []interface{}, got %T", slice[1])
	}

	// Check string
	if slice[2] != "string" {
		t.Errorf("slice[2] = %v, want string", slice[2])
	}
}

// TestDecodeOrderedMapEdgeCases tests edge cases in decodeOrderedMap
func TestDecodeOrderedMapEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
	}{
		{
			name:    "empty object",
			json:    `{}`,
			wantErr: false,
		},
		{
			name:    "object with null values",
			json:    `{"key": null}`,
			wantErr: false,
		},
		{
			name:    "object with empty array",
			json:    `{"arr": []}`,
			wantErr: false,
		},
		{
			name:    "object with empty nested object",
			json:    `{"nested": {}}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := NewOrderedMap()
			err := om.UnmarshalJSON([]byte(tt.json))
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestDecodeSliceEdgeCases tests edge cases in decodeSlice
func TestDecodeSliceEdgeCases(t *testing.T) {
	jsonData := []byte(`{
		"empty": [],
		"single": [1],
		"nested_empty": [[]],
		"nested_objects": [{"a": 1}, {"b": 2}]
	}`)

	om := NewOrderedMap()
	err := om.UnmarshalJSON(jsonData)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	// Check empty array
	empty, _ := om.Get("empty")
	if emptySlice, ok := empty.([]interface{}); !ok || len(emptySlice) != 0 {
		t.Errorf("empty should be empty slice")
	}

	// Check single element
	single, _ := om.Get("single")
	if singleSlice, ok := single.([]interface{}); !ok || len(singleSlice) != 1 {
		t.Errorf("single should have 1 element")
	}

	// Check nested empty
	nestedEmpty, _ := om.Get("nested_empty")
	if nestedSlice, ok := nestedEmpty.([]interface{}); !ok || len(nestedSlice) != 1 {
		t.Errorf("nested_empty should have 1 element")
	}
}

// TestValues tests the Values() method
func TestValues(t *testing.T) {
	om := NewOrderedMap()
	om.Set("key1", "value1")
	om.Set("key2", 42)
	om.Set("key3", true)

	values := om.Values()
	if len(values) != 3 {
		t.Errorf("Values() length = %v, want 3", len(values))
	}
	if values["key1"] != "value1" {
		t.Errorf("values[key1] = %v, want value1", values["key1"])
	}
	if values["key2"] != 42 {
		t.Errorf("values[key2] = %v, want 42", values["key2"])
	}
	if values["key3"] != true {
		t.Errorf("values[key3] = %v, want true", values["key3"])
	}
}

// TestUnmarshalJSONErrors tests error handling in UnmarshalJSON
func TestUnmarshalJSONErrors(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		wantErr bool
	}{
		{
			name:    "invalid json",
			json:    `{invalid}`,
			wantErr: true,
		},
		{
			name:    "unclosed object",
			json:    `{"key": "value"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := NewOrderedMap()
			err := om.UnmarshalJSON([]byte(tt.json))
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestDecodeOrderedMapWithExistingOrderedMap tests decoding when value is already an OrderedMap
func TestDecodeOrderedMapWithExistingOrderedMap(t *testing.T) {
	// First create an OrderedMap with nested OrderedMap
	inner := NewOrderedMap()
	inner.Set("inner_key", "inner_value")
	
	om := NewOrderedMap()
	om.Set("outer", *inner)
	
	// Now marshal and unmarshal to trigger the OrderedMap path
	data, err := om.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}
	
	om2 := NewOrderedMap()
	err = om2.UnmarshalJSON(data)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	
	outer, exists := om2.Get("outer")
	if !exists {
		t.Fatal("outer should exist")
	}
	
	outerMap, ok := outer.(OrderedMap)
	if !ok {
		t.Fatalf("outer should be OrderedMap, got %T", outer)
	}
	
	val, exists := outerMap.Get("inner_key")
	if !exists || val != "inner_value" {
		t.Errorf("inner_key = %v, want inner_value", val)
	}
}

// TestDecodeSliceWithExistingOrderedMap tests decodeSlice when slice element is already an OrderedMap
func TestDecodeSliceWithExistingOrderedMap(t *testing.T) {
	inner := NewOrderedMap()
	inner.Set("key", "value")
	
	om := NewOrderedMap()
	om.Set("array", []interface{}{*inner})
	
	data, err := om.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}
	
	om2 := NewOrderedMap()
	err = om2.UnmarshalJSON(data)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	
	array, exists := om2.Get("array")
	if !exists {
		t.Fatal("array should exist")
	}
	
	slice, ok := array.([]interface{})
	if !ok {
		t.Fatalf("array should be []interface{}, got %T", array)
	}
	
	if len(slice) != 1 {
		t.Fatalf("slice length = %v, want 1", len(slice))
	}
	
	elem, ok := slice[0].(OrderedMap)
	if !ok {
		t.Fatalf("slice[0] should be OrderedMap, got %T", slice[0])
	}
	
	val, _ := elem.Get("key")
	if val != "value" {
		t.Errorf("key = %v, want value", val)
	}
}

// TestDecodeSliceOutOfBoundsObject tests decodeSlice when object index is out of bounds
func TestDecodeSliceOutOfBoundsObject(t *testing.T) {
	jsonData := []byte(`{
		"arr": [{"a": 1}, {"b": 2}, {"c": 3}]
	}`)
	
	om := NewOrderedMap()
	// Pre-populate with shorter slice
	om.Set("arr", []interface{}{})
	
	err := om.UnmarshalJSON(jsonData)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	
	arr, _ := om.Get("arr")
	slice, ok := arr.([]interface{})
	if !ok {
		t.Fatal("arr should be []interface{}")
	}
	
	// Should have processed all 3 objects even though initial slice was empty
	if len(slice) != 3 {
		t.Errorf("slice length = %v, want 3", len(slice))
	}
}

// TestDecodeSliceOutOfBoundsArray tests decodeSlice when array index is out of bounds
func TestDecodeSliceOutOfBoundsArray(t *testing.T) {
	jsonData := []byte(`{
		"arr": [[1], [2], [3]]
	}`)
	
	om := NewOrderedMap()
	om.Set("arr", []interface{}{})
	
	err := om.UnmarshalJSON(jsonData)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	
	arr, _ := om.Get("arr")
	slice, ok := arr.([]interface{})
	if !ok {
		t.Fatal("arr should be []interface{}")
	}
	
	if len(slice) != 3 {
		t.Errorf("slice length = %v, want 3", len(slice))
	}
}

// TestDecodeSliceNonMapNonSliceObject tests when slice element is neither map nor slice but object encountered
func TestDecodeSliceNonMapNonSliceObject(t *testing.T) {
	jsonData := []byte(`{
		"arr": [42, {"key": "value"}]
	}`)
	
	om := NewOrderedMap()
	om.Set("arr", []interface{}{42, "string"})
	
	err := om.UnmarshalJSON(jsonData)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	
	arr, _ := om.Get("arr")
	slice, ok := arr.([]interface{})
	if !ok {
		t.Fatal("arr should be []interface{}")
	}
	
	if len(slice) != 2 {
		t.Errorf("slice length = %v, want 2", len(slice))
	}
}

// TestDecodeSliceNonSliceArray tests when slice element is not a slice but array encountered
func TestDecodeSliceNonSliceArray(t *testing.T) {
	jsonData := []byte(`{
		"arr": ["string", [1, 2]]
	}`)
	
	om := NewOrderedMap()
	om.Set("arr", []interface{}{"string", 42})
	
	err := om.UnmarshalJSON(jsonData)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	
	arr, _ := om.Get("arr")
	slice, ok := arr.([]interface{})
	if !ok {
		t.Fatal("arr should be []interface{}")
	}
	
	if len(slice) != 2 {
		t.Errorf("slice length = %v, want 2", len(slice))
	}
}

// TestDecodeOrderedMapNonMapNonOrderedMap tests when value is neither map nor OrderedMap but object encountered
func TestDecodeOrderedMapNonMapNonOrderedMap(t *testing.T) {
	jsonData := []byte(`{
		"key": {"nested": "value"}
	}`)
	
	om := NewOrderedMap()
	om.Set("key", "string_value")
	
	err := om.UnmarshalJSON(jsonData)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	
	key, exists := om.Get("key")
	if !exists {
		t.Fatal("key should exist")
	}
	
	// Should have been replaced with OrderedMap
	keyMap, ok := key.(OrderedMap)
	if !ok {
		t.Fatalf("key should be OrderedMap, got %T", key)
	}
	
	nested, _ := keyMap.Get("nested")
	if nested != "value" {
		t.Errorf("nested = %v, want value", nested)
	}
}

// TestDecodeOrderedMapNonSliceArray tests when value is not a slice but array encountered
func TestDecodeOrderedMapNonSliceArray(t *testing.T) {
	jsonData := []byte(`{
		"arr": [1, 2, 3]
	}`)
	
	om := NewOrderedMap()
	om.Set("arr", "not_a_slice")
	
	err := om.UnmarshalJSON(jsonData)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	
	arr, exists := om.Get("arr")
	if !exists {
		t.Fatal("arr should exist")
	}
	
	// Should still be the original value since we can't convert
	slice, ok := arr.([]interface{})
	if !ok {
		t.Fatalf("arr should be []interface{}, got %T", arr)
	}
	
	if len(slice) != 3 {
		t.Errorf("slice length = %v, want 3", len(slice))
	}
}

// TestMarshalJSONEmptyMap tests marshaling an empty OrderedMap
func TestMarshalJSONEmptyMap(t *testing.T) {
	om := NewOrderedMap()
	
	data, err := om.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}
	
	// Remove newlines
	str := strings.ReplaceAll(string(data), "\n", "")
	if str != "{}" {
		t.Errorf("MarshalJSON() = %v, want {}", str)
	}
}

// TestOrderedMapEscapeHTMLPropagation tests that escapeHTML is propagated to nested maps
func TestOrderedMapEscapeHTMLPropagation(t *testing.T) {
	jsonData := []byte(`{
		"nested": {
			"html": "<script>"
		}
	}`)

	om := NewOrderedMap()
	om.SetEscapeHTML(false)
	err := om.UnmarshalJSON(jsonData)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	nested, _ := om.Get("nested")
	nestedMap, ok := nested.(OrderedMap)
	if !ok {
		t.Fatal("nested should be OrderedMap")
	}

	// Marshal the nested map
	data, err := nestedMap.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	dataStr := string(data)
	dataStr = strings.ReplaceAll(dataStr, "\n", "")
	// Should not escape HTML
	if !strings.Contains(dataStr, "<script>") {
		t.Errorf("Expected HTML to not be escaped in nested map, got: %s", dataStr)
	}
}
// TestDecodeSliceInBoundsNonMapObject tests when slice has element that's not map/OrderedMap but object encountered
func TestDecodeSliceInBoundsNonMapObject(t *testing.T) {
	jsonData := []byte(`{
		"arr": [{"key": "value"}]
	}`)
	
	om := NewOrderedMap()
	// Pre-populate with non-map value
	om.Set("arr", []interface{}{42})
	
	err := om.UnmarshalJSON(jsonData)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	
	arr, _ := om.Get("arr")
	slice, ok := arr.([]interface{})
	if !ok {
		t.Fatal("arr should be []interface{}")
	}
	
	if len(slice) != 1 {
		t.Errorf("slice length = %v, want 1", len(slice))
	}
}

// TestDecodeSliceInBoundsNonSliceArray tests when slice has element that's not slice but array encountered
func TestDecodeSliceInBoundsNonSliceArray(t *testing.T) {
	jsonData := []byte(`{
		"arr": [[1, 2]]
	}`)
	
	om := NewOrderedMap()
	// Pre-populate with non-slice value
	om.Set("arr", []interface{}{"string"})
	
	err := om.UnmarshalJSON(jsonData)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	
	arr, _ := om.Get("arr")
	slice, ok := arr.([]interface{})
	if !ok {
		t.Fatal("arr should be []interface{}")
	}
	
	if len(slice) != 1 {
		t.Errorf("slice length = %v, want 1", len(slice))
	}
}

// TestDecodeSliceWithMapInSlice tests slice containing map[string]interface{}
func TestDecodeSliceWithMapInSlice(t *testing.T) {
	jsonData := []byte(`{
		"arr": [{"nested": "value"}]
	}`)
	
	om := NewOrderedMap()
	// Pre-populate with map[string]interface{}
	innerMap := make(map[string]interface{})
	innerMap["old"] = "data"
	om.Set("arr", []interface{}{innerMap})
	
	err := om.UnmarshalJSON(jsonData)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	
	arr, _ := om.Get("arr")
	slice, ok := arr.([]interface{})
	if !ok {
		t.Fatal("arr should be []interface{}")
	}
	
	if len(slice) != 1 {
		t.Errorf("slice length = %v, want 1", len(slice))
	}
	
	elem, ok := slice[0].(OrderedMap)
	if !ok {
		t.Fatalf("slice[0] should be OrderedMap, got %T", slice[0])
	}
	
	val, _ := elem.Get("nested")
	if val != "value" {
		t.Errorf("nested = %v, want value", val)
	}
}

// TestDecodeSliceWithSliceInSlice tests slice containing []interface{}
func TestDecodeSliceWithSliceInSlice(t *testing.T) {
	jsonData := []byte(`{
		"matrix": [[1, 2, 3]]
	}`)
	
	om := NewOrderedMap()
	// Pre-populate with existing slice
	inner := []interface{}{99}
	om.Set("matrix", []interface{}{inner})
	
	err := om.UnmarshalJSON(jsonData)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	
	matrix, _ := om.Get("matrix")
	outer, ok := matrix.([]interface{})
	if !ok {
		t.Fatal("matrix should be []interface{}")
	}
	
	if len(outer) != 1 {
		t.Errorf("outer length = %v, want 1", len(outer))
	}
}

// TestUnmarshalJSONErrorInToken tests error handling when decoder fails
func TestUnmarshalJSONErrorInToken(t *testing.T) {
	// Malformed JSON that will cause decoder errors
	tests := []struct {
		name string
		json string
	}{
		{
			name: "truncated after key",
			json: `{"key":`,
		},
		{
			name: "invalid nested object",
			json: `{"key": {invalid}}`,
		},
		{
			name: "invalid nested array",
			json: `{"arr": [invalid]}`,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			om := NewOrderedMap()
			err := om.UnmarshalJSON([]byte(tt.json))
			if err == nil {
				t.Error("Expected error for malformed JSON")
			}
		})
	}
}

// TestMarshalJSONWithError tests error handling in MarshalJSON
func TestMarshalJSONWithError(t *testing.T) {
	om := NewOrderedMap()
	// Add a value that can't be marshaled (channels can't be marshaled)
	ch := make(chan int)
	om.Set("invalid", ch)
	
	_, err := om.MarshalJSON()
	if err == nil {
		t.Error("Expected error when marshaling unmarshalable value")
	}
}