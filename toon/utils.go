package toon

import (
	"math"
	"reflect"
)

// isPrimitive checks if a value is a primitive type (nil, bool, number, or string).
func isPrimitive(v Value) bool {
	if v == nil {
		return true
	}

	switch v.(type) {
	case bool, int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64, string:
		return true
	default:
		return false
	}
}

// isMap checks if a value is a map.
func isMap(v Value) bool {
	if v == nil {
		return false
	}
	// Check for OrderedMap
	if _, ok := v.(OrderedMap); ok {
		return true
	}
	if _, ok := v.(*OrderedMap); ok {
		return true
	}
	rv := reflect.ValueOf(v)
	return rv.Kind() == reflect.Map
}

// isList checks if a value is a list (slice or array).
func isList(v Value) bool {
	if v == nil {
		return false
	}
	rv := reflect.ValueOf(v)
	k := rv.Kind()
	return k == reflect.Slice || k == reflect.Array
}

// allPrimitives checks if all elements in a slice are primitives.
func allPrimitives(v Value) bool {
	if !isList(v) {
		return false
	}

	rv := reflect.ValueOf(v)
	for i := 0; i < rv.Len(); i++ {
		if !isPrimitive(rv.Index(i).Interface()) {
			return false
		}
	}
	return true
}

// allMaps checks if all elements in a slice are maps.
func allMaps(v Value) bool {
	if !isList(v) {
		return false
	}

	rv := reflect.ValueOf(v)
	for i := 0; i < rv.Len(); i++ {
		if !isMap(rv.Index(i).Interface()) {
			return false
		}
	}
	return true
}

// sameKeys checks if all maps in a slice have the same keys.
func sameKeys(v Value) bool {
	if !isList(v) {
		return false
	}

	rv := reflect.ValueOf(v)
	if rv.Len() == 0 {
		return true
	}

	// Get keys from first map
	first := rv.Index(0).Interface()
	if !isMap(first) {
		return false
	}

	firstKeys := getMapKeys(first)
	if len(firstKeys) == 0 {
		return true
	}

	// Check all other maps have the same keys
	for i := 1; i < rv.Len(); i++ {
		item := rv.Index(i).Interface()
		if !isMap(item) {
			return false
		}

		keys := getMapKeys(item)
		if !sameStringSlice(firstKeys, keys) {
			return false
		}
	}

	return true
}

// getMapKeys returns sorted string keys from a map.
func getMapKeys(v Value) []string {
	// Handle OrderedMap
	if orderedMap, ok := v.(OrderedMap); ok {
		keys := make([]string, len(orderedMap.Keys()))
		copy(keys, orderedMap.Keys())
		sortStrings(keys)
		return keys
	}
	if orderedMapPtr, ok := v.(*OrderedMap); ok {
		keys := make([]string, len(orderedMapPtr.Keys()))
		copy(keys, orderedMapPtr.Keys())
		sortStrings(keys)
		return keys
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Map {
		return nil
	}

	keys := make([]string, 0, rv.Len())
	for _, k := range rv.MapKeys() {
		keys = append(keys, k.String())
	}

	// Sort for consistent comparison
	sortStrings(keys)
	return keys
}

// sameStringSlice checks if two string slices contain the same elements (order-independent).
func sameStringSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	// Create map for O(1) lookup
	m := make(map[string]bool, len(a))
	for _, s := range a {
		m[s] = true
	}

	for _, s := range b {
		if !m[s] {
			return false
		}
	}

	return true
}

// sortStrings sorts a slice of strings in place (simple bubble sort for small slices).
func sortStrings(s []string) {
	for i := 0; i < len(s); i++ {
		for j := i + 1; j < len(s); j++ {
			if s[i] > s[j] {
				s[i], s[j] = s[j], s[i]
			}
		}
	}
}

// sortKeysWithArraysFirst sorts keys alphabetically but puts array fields first in list items
func sortKeysWithArraysFirst(keys []string, rv reflect.Value) {
	// Separate array keys from non-array keys
	arrayKeys := []string{}
	otherKeys := []string{}

	for _, k := range keys {
		mapKey := reflect.ValueOf(k)
		val := rv.MapIndex(mapKey).Interface()
		if isList(val) {
			arrayKeys = append(arrayKeys, k)
		} else {
			otherKeys = append(otherKeys, k)
		}
	}

	// Sort each group
	sortStrings(arrayKeys)
	sortStrings(otherKeys)

	// Rebuild keys with arrays first
	copy(keys, arrayKeys)
	copy(keys[len(arrayKeys):], otherKeys)
}

// normalize normalizes a value for encoding, converting to JSON-compatible types.
func normalize(v Value) Value {
	if v == nil {
		return nil
	}

	switch val := v.(type) {
	case bool, string:
		return val

	case int, int8, int16, int32, int64:
		return reflect.ValueOf(val).Int()

	case uint, uint8, uint16, uint32, uint64:
		return int64(reflect.ValueOf(val).Uint())

	case float32:
		return normalizeFloat(float64(val))

	case float64:
		return normalizeFloat(val)

	case []interface{}:
		return normalizeSlice(val)

	case map[string]interface{}:
		return normalizeMap(val)

	case OrderedMap:
		return normalizeOrderedMap(&val)

	case *OrderedMap:
		return normalizeOrderedMap(val)

	default:
		return normalizeReflection(v)
	}
}

// normalizeSlice normalizes a slice of values.
func normalizeSlice(slice []interface{}) Value {
	result := make([]Value, len(slice))
	for i, item := range slice {
		result[i] = normalize(item)
	}
	return result
}

// normalizeMap normalizes a map[string]interface{}.
func normalizeMap(m map[string]interface{}) Value {
	result := make(map[string]Value, len(m))
	for k, item := range m {
		result[k] = normalize(item)
	}
	return result
}

// normalizeOrderedMap normalizes an OrderedMap.
func normalizeOrderedMap(om *OrderedMap) Value {
	result := NewOrderedMap()
	for _, k := range om.Keys() {
		if v, ok := om.Get(k); ok {
			result.Set(k, normalize(v))
		}
	}
	return *result
}

// normalizeReflection handles normalization using reflection for non-standard types.
func normalizeReflection(v Value) Value {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		return normalizeReflectSlice(rv)

	case reflect.Map:
		return normalizeReflectMap(rv)

	default:
		// Unsupported type, return nil
		return nil
	}
}

// normalizeReflectSlice normalizes a slice using reflection.
func normalizeReflectSlice(rv reflect.Value) Value {
	length := rv.Len()
	result := make([]Value, length)
	for i := 0; i < length; i++ {
		result[i] = normalize(rv.Index(i).Interface())
	}
	return result
}

// normalizeReflectMap normalizes a map using reflection.
func normalizeReflectMap(rv reflect.Value) Value {
	result := make(map[string]Value)
	for _, k := range rv.MapKeys() {
		key := k.String()
		result[key] = normalize(rv.MapIndex(k).Interface())
	}
	return result
}

// normalizeFloat handles special float values.
func normalizeFloat(f float64) Value {
	// Handle negative zero - normalize to 0
	if f == 0 && math.Signbit(f) {
		return int64(0)
	}

	// Handle NaN and Infinity
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return nil
	}

	// Check if it's a whole number
	if f == math.Trunc(f) && f >= math.MinInt64 && f <= math.MaxInt64 {
		return int64(f)
	}

	return f
}

// toFloat64 converts a numeric value to float64.
func toFloat64(v Value) (float64, bool) {
	if v == nil {
		return 0, false
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Float32, reflect.Float64:
		return rv.Float(), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(rv.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(rv.Uint()), true
	default:
		return 0, false
	}
}

// toInt64 converts a numeric value to int64.
func toInt64(v Value) (int64, bool) {
	switch val := v.(type) {
	case int, int8, int16, int32, int64:
		return toInt64FromSigned(val), true
	case uint, uint8, uint16, uint32, uint64:
		return toInt64FromUnsigned(val)
	case float32, float64:
		return toInt64FromFloat(val)
	default:
		return 0, false
	}
}

// toInt64FromSigned converts signed integers to int64.
func toInt64FromSigned(v Value) int64 {
	rv := reflect.ValueOf(v)
	return rv.Int()
}

// toInt64FromUnsigned converts unsigned integers to int64.
func toInt64FromUnsigned(v Value) (int64, bool) {
	rv := reflect.ValueOf(v)
	uval := rv.Uint()
	if uval <= math.MaxInt64 {
		return int64(uval), true
	}
	return 0, false
}

// toInt64FromFloat converts floats to int64 if they are whole numbers.
func toInt64FromFloat(v Value) (int64, bool) {
	rv := reflect.ValueOf(v)
	fval := rv.Float()
	ival := int64(fval)
	if fval == float64(ival) {
		return ival, true
	}
	return 0, false
}
