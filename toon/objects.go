package toon

import (
	"fmt"
	"reflect"
	"strings"
)

// encodeObject encodes a map to TOON format.
func encodeObject(w *writer, key string, v Value, depth int, opts *EncodeOptions) error {
	// Handle OrderedMap - preserve key order
	var keys []string
	var rv reflect.Value

	if orderedMap, ok := v.(OrderedMap); ok {
		keys = orderedMap.Keys()
		rv = reflect.ValueOf(orderedMap.Values())
	} else if orderedMapPtr, ok := v.(*OrderedMap); ok {
		keys = orderedMapPtr.Keys()
		rv = reflect.ValueOf(orderedMapPtr.Values())
	} else {
		// Regular map
		rv = reflect.ValueOf(v)
		if rv.Kind() != reflect.Map {
			return &EncodeError{Message: "not a map", Value: v}
		}

		// Get and sort keys
		keys = make([]string, 0, rv.Len())
		for _, k := range rv.MapKeys() {
			keys = append(keys, k.String())
		}
		sortStrings(keys)
	}

	// Handle empty object
	if len(keys) == 0 {
		if key != "" {
			w.push(key+colon, depth)
		}
		return nil
	}

	// If this is a keyed object, write the key first
	if key != "" {
		w.push(key+colon, depth)
		depth++
	}

	// Encode each key-value pair
	for _, k := range keys {
		var mapValue interface{}

		// Get value based on type
		if rv.Kind() == reflect.Map {
			mapKey := reflect.ValueOf(k)
			mapValue = rv.MapIndex(mapKey).Interface()
		} else {
			// For OrderedMap, rv is already the values map
			mapValue = rv.Interface().(map[string]interface{})[k]
		}

		encodedKey := encodeKey(k)

		if err := encodeValue(w, encodedKey, mapValue, depth, opts); err != nil {
			return err
		}
	}

	return nil
}

// flattenObject flattens nested maps into dotted key notation.
// Example: {"a":{"b":1}} becomes {"a.b":1}
func flattenObject(obj map[string]Value, currentPath string, depth int, opts *EncodeOptions) (map[string]Value, error) {
	result := make(map[string]Value)

	// First pass: check if any potential flattened paths would collide with literal keys
	literalKeys := make(map[string]bool)
	for key := range obj {
		literalKeys[key] = true
	}

	// Check for collisions and handle accordingly
	if hasCollision := checkFlattenCollisions(obj, literalKeys); hasCollision {
		if opts.Strict {
			return nil, &EncodeError{
				Message: "key collision: flattened path would conflict with existing literal key",
				Value:   obj,
			}
		}
		return obj, nil
	}

	for key, value := range obj {
		fullPath := buildFullPath(currentPath, key)
		segmentCount := countPathSegments(currentPath)

		if shouldFlatten := shouldFlattenValue(value, key, segmentCount, opts); shouldFlatten {
			isEmpty, nestedMap := isEmptyNestedMap(value)
			if isEmpty {
				// Empty object - add at current path
				if err := addToResultWithCollisionCheck(result, fullPath, value, opts); err != nil {
					return nil, err
				}
				continue
			}

			// Recursively flatten nested maps up to depth limit
			nested, err := flattenObject(nestedMap, fullPath, depth+1, opts)
			if err != nil {
				return nil, err
			}
			// Merge nested results
			for k, v := range nested {
				if err := addToResultWithCollisionCheck(result, k, v, opts); err != nil {
					return nil, err
				}
			}
		} else {
			// Stop flattening - add the value at current path
			if err := addToResultWithCollisionCheck(result, fullPath, value, opts); err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

// convertToNestedMap converts various map types to map[string]Value.
func convertToNestedMap(value Value) (map[string]Value, bool) {
	if orderedMap, ok := value.(OrderedMap); ok {
		nestedMap := make(map[string]Value)
		for k, v := range orderedMap.Values() {
			nestedMap[k] = v
		}
		return nestedMap, true
	}
	if orderedMapPtr, ok := value.(*OrderedMap); ok {
		nestedMap := make(map[string]Value)
		for k, v := range orderedMapPtr.Values() {
			nestedMap[k] = v
		}
		return nestedMap, true
	}
	if m, ok := value.(map[string]Value); ok {
		return m, true
	}
	return nil, false
}

// collectPotentialPaths recursively collects potential flattened paths from a nested map.
func collectPotentialPaths(m map[string]Value, prefix string, paths map[string]bool) {
	for k, v := range m {
		var path string
		if prefix == "" {
			path = k
		} else {
			path = prefix + "." + k
		}
		paths[path] = true
		if isMap(v) {
			if nm, ok := convertToNestedMap(v); ok {
				collectPotentialPaths(nm, path, paths)
			}
		}
	}
}

// checkFlattenCollisions checks if any potential flattened paths would collide with literal keys.
func checkFlattenCollisions(obj map[string]Value, literalKeys map[string]bool) bool {
	for key, value := range obj {
		if !isMap(value) {
			continue
		}
		nestedMap, ok := convertToNestedMap(value)
		if !ok {
			continue
		}
		// Collect potential paths from this nested structure
		tempPaths := make(map[string]bool)
		collectPotentialPaths(nestedMap, key, tempPaths)

		// Check if any of these paths collide with literal keys
		for path := range tempPaths {
			if literalKeys[path] {
				return true
			}
		}
	}
	return false
}

// containsQuotedKeys checks if a nested structure contains any keys that need quoting.
func containsQuotedKeys(m map[string]Value) bool {
	for k, v := range m {
		if !safeKey(k) {
			return true
		}
		if isMap(v) {
			if nm, ok := convertToNestedMap(v); ok && containsQuotedKeys(nm) {
				return true
			}
		}
	}
	return false
}

// countPathSegments counts the number of segments in a path.
func countPathSegments(currentPath string) int {
	if currentPath == "" {
		return 1
	}
	return len(strings.Split(currentPath, ".")) + 1
}

// buildFullPath builds the full path from current path and key.
func buildFullPath(currentPath, key string) string {
	if currentPath == "" {
		return key
	}
	return currentPath + "." + key
}

// shouldFlattenValue determines if a value should be flattened.
func shouldFlattenValue(value Value, key string, segmentCount int, opts *EncodeOptions) bool {
	if !isMap(value) {
		return false
	}
	if !safeKey(key) {
		return false
	}
	if opts.FlattenDepth <= 0 || segmentCount >= opts.FlattenDepth {
		return false
	}
	// Check if the nested structure contains any keys that need quoting
	nestedMap, ok := convertToNestedMap(value)
	if !ok {
		return false
	}
	return !containsQuotedKeys(nestedMap)
}

// addToResultWithCollisionCheck adds a value to result with collision checking.
func addToResultWithCollisionCheck(result map[string]Value, fullPath string, value Value, opts *EncodeOptions) error {
	if existing, exists := result[fullPath]; exists {
		if opts.Strict {
			return &EncodeError{
				Message: fmt.Sprintf("key collision: %q", fullPath),
				Value:   existing,
			}
		}
		// Non-strict: last value wins
	}
	result[fullPath] = value
	return nil
}

// isEmptyNestedMap checks if a nested map is empty.
func isEmptyNestedMap(value Value) (bool, map[string]Value) {
	if orderedMap, ok := value.(OrderedMap); ok {
		nestedMap := make(map[string]Value)
		for k, v := range orderedMap.Values() {
			nestedMap[k] = v
		}
		return orderedMap.Len() == 0, nestedMap
	}
	if orderedMapPtr, ok := value.(*OrderedMap); ok {
		nestedMap := make(map[string]Value)
		for k, v := range orderedMapPtr.Values() {
			nestedMap[k] = v
		}
		return orderedMapPtr.Len() == 0, nestedMap
	}
	rv := reflect.ValueOf(value)
	isEmpty := rv.Len() == 0
	nestedMap := value.(map[string]Value)
	return isEmpty, nestedMap
}
