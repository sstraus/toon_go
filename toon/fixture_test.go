package toon

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

// FixtureTest represents a single test case from a fixture file
type FixtureTest struct {
	Name        string                 `json:"name"`
	Input       interface{}            `json:"input"`
	Expected    interface{}            `json:"expected"`
	ShouldError bool                   `json:"shouldError"`
	Options     map[string]interface{} `json:"options"`
	SpecSection string                 `json:"specSection"`
	Note        string                 `json:"note"`
}

// Fixture represents the structure of a fixture file
type Fixture struct {
	Version     string        `json:"version"`
	Category    string        `json:"category"`
	Description string        `json:"description"`
	Tests       []FixtureTest `json:"tests"`
}

// loadFixture reads and parses a JSON fixture file
func loadFixture(filename string) (*Fixture, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read fixture file: %w", err)
	}

	// First parse into OrderedMap to preserve key order in input objects
	var fixtureMap OrderedMap
	if err := json.Unmarshal(data, &fixtureMap); err != nil {
		return nil, fmt.Errorf("failed to parse fixture JSON: %w", err)
	}

	// Convert to Fixture struct
	fixture := &Fixture{}
	if v, ok := fixtureMap.Get("version"); ok {
		if s, ok := v.(string); ok {
			fixture.Version = s
		}
	}
	if v, ok := fixtureMap.Get("category"); ok {
		if s, ok := v.(string); ok {
			fixture.Category = s
		}
	}
	if v, ok := fixtureMap.Get("description"); ok {
		if s, ok := v.(string); ok {
			fixture.Description = s
		}
	}
	if v, ok := fixtureMap.Get("tests"); ok {
		if tests, ok := v.([]interface{}); ok {
			fixture.Tests = make([]FixtureTest, len(tests))
			for i, t := range tests {
				if testMap, ok := t.(OrderedMap); ok {
					fixture.Tests[i] = fixtureTestFromOrderedMap(testMap)
				} else if testMap, ok := t.(map[string]interface{}); ok {
					fixture.Tests[i] = fixtureTestFromMap(testMap)
				}
			}
		}
	}

	return fixture, nil
}

// fixtureTestFromOrderedMap converts an OrderedMap to a FixtureTest
func fixtureTestFromOrderedMap(m OrderedMap) FixtureTest {
	test := FixtureTest{}
	if v, ok := m.Get("name"); ok {
		if s, ok := v.(string); ok {
			test.Name = s
		}
	}
	if v, ok := m.Get("input"); ok {
		test.Input = v
	}
	if v, ok := m.Get("expected"); ok {
		test.Expected = v
	}
	if v, ok := m.Get("shouldError"); ok {
		if b, ok := v.(bool); ok {
			test.ShouldError = b
		}
	}
	if v, ok := m.Get("options"); ok {
		if opts, ok := v.(map[string]interface{}); ok {
			test.Options = opts
		} else if opts, ok := v.(OrderedMap); ok {
			test.Options = opts.Values()
		}
	}
	if v, ok := m.Get("specSection"); ok {
		if s, ok := v.(string); ok {
			test.SpecSection = s
		}
	}
	if v, ok := m.Get("note"); ok {
		if s, ok := v.(string); ok {
			test.Note = s
		}
	}
	return test
}

// fixtureTestFromMap converts a regular map to a FixtureTest (fallback)
func fixtureTestFromMap(m map[string]interface{}) FixtureTest {
	test := FixtureTest{}
	if v, ok := m["name"]; ok {
		if s, ok := v.(string); ok {
			test.Name = s
		}
	}
	if v, ok := m["input"]; ok {
		test.Input = v
	}
	if v, ok := m["expected"]; ok {
		test.Expected = v
	}
	if v, ok := m["shouldError"]; ok {
		if b, ok := v.(bool); ok {
			test.ShouldError = b
		}
	}
	if v, ok := m["options"]; ok {
		if opts, ok := v.(map[string]interface{}); ok {
			test.Options = opts
		}
	}
	if v, ok := m["specSection"]; ok {
		if s, ok := v.(string); ok {
			test.SpecSection = s
		}
	}
	if v, ok := m["note"]; ok {
		if s, ok := v.(string); ok {
			test.Note = s
		}
	}
	return test
}

// fixtureOptionsToEncodeOptions converts fixture options to EncodeOptions
func fixtureOptionsToEncodeOptions(opts map[string]interface{}) *EncodeOptions {
	if opts == nil {
		return nil
	}

	encOpts := &EncodeOptions{}

	if delimiter, ok := opts["delimiter"].(string); ok {
		encOpts.Delimiter = delimiter
	}

	if lengthMarker, ok := opts["length_marker"].(string); ok {
		encOpts.LengthMarker = lengthMarker
	} else if lengthMarker, ok := opts["lengthMarker"].(string); ok {
		encOpts.LengthMarker = lengthMarker
	}

	if indent, ok := opts["indent"].(float64); ok {
		encOpts.Indent = int(indent)
	} else if indentSize, ok := opts["indent_size"].(float64); ok {
		encOpts.Indent = int(indentSize)
	} else if indentSize, ok := opts["indentSize"].(float64); ok {
		encOpts.Indent = int(indentSize)
	}

	// Handle keyFolding option: "safe" enables flattening, "off" disables it
	if keyFolding, ok := opts["keyFolding"].(string); ok {
		encOpts.FlattenPaths = (keyFolding == "safe")
	}
	
	if flattenPaths, ok := opts["flattenPaths"].(bool); ok {
		encOpts.FlattenPaths = flattenPaths
	}
	
	if flattenDepth, ok := opts["flattenDepth"].(float64); ok {
		encOpts.FlattenDepth = int(flattenDepth)
	} else if flattenDepth, ok := opts["flattenDepth"].(int); ok {
		encOpts.FlattenDepth = flattenDepth
	} else {
		// Not set - use -1 as sentinel for "not specified"
		encOpts.FlattenDepth = -1
	}
	
	if strict, ok := opts["strict"].(bool); ok {
		encOpts.Strict = strict
	}

	return encOpts
}

// fixtureOptionsToDecodeOptions converts fixture options to DecodeOptions
func fixtureOptionsToDecodeOptions(opts map[string]interface{}) *DecodeOptions {
	if opts == nil {
		return nil
	}

	decOpts := &DecodeOptions{}

	if strict, ok := opts["strict"].(bool); ok {
		decOpts.Strict = strict
	}

	if indentSize, ok := opts["indent_size"].(float64); ok {
		decOpts.IndentSize = int(indentSize)
	} else if indentSize, ok := opts["indentSize"].(float64); ok {
		decOpts.IndentSize = int(indentSize)
	}

	if expandPaths, ok := opts["expandPaths"].(string); ok {
		decOpts.ExpandPaths = expandPaths
	}

	// Keys field is KeyMode, not []string - use default StringKeys
	decOpts.Keys = StringKeys

	return decOpts
}

// deepEqual compares two values for semantic equality, handling type conversions
func deepEqual(a, b interface{}) bool {
	// Handle nil
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Direct comparison
	if reflect.DeepEqual(a, b) {
		return true
	}

	// Handle numeric type conversions
	aInt, aIsInt := toInt64(a)
	bInt, bIsInt := toInt64(b)
	if aIsInt && bIsInt {
		return aInt == bInt
	}

	aFloat, aIsFloat := toFloat64(a)
	bFloat, bIsFloat := toFloat64(b)
	if aIsFloat && bIsFloat {
		return aFloat == bFloat
	}

	// Handle arrays
	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)

	if aVal.Kind() == reflect.Slice && bVal.Kind() == reflect.Slice {
		if aVal.Len() != bVal.Len() {
			return false
		}
		for i := 0; i < aVal.Len(); i++ {
			if !deepEqual(aVal.Index(i).Interface(), bVal.Index(i).Interface()) {
				return false
			}
		}
		return true
	}

	// Handle maps
	if aVal.Kind() == reflect.Map && bVal.Kind() == reflect.Map {
		if aVal.Len() != bVal.Len() {
			return false
		}
		
		aKeys := aVal.MapKeys()
		for _, key := range aKeys {
			aElem := aVal.MapIndex(key)
			bElem := bVal.MapIndex(key)
			
			if !bElem.IsValid() {
				return false
			}
			
			if !deepEqual(aElem.Interface(), bElem.Interface()) {
				return false
			}
		}
		return true
	}

	return false
}

// normalizeValue converts fixture values to comparable Go types
func normalizeValue(v interface{}) interface{} {
	if v == nil {
		return nil
	}

	// Handle OrderedMap - convert to regular map for comparison
	if orderedMap, ok := v.(OrderedMap); ok {
		result := make(map[string]interface{})
		for _, key := range orderedMap.Keys() {
			if val, ok := orderedMap.Get(key); ok {
				result[key] = normalizeValue(val)
			}
		}
		return result
	}
	if orderedMapPtr, ok := v.(*OrderedMap); ok {
		result := make(map[string]interface{})
		for _, key := range orderedMapPtr.Keys() {
			if val, ok := orderedMapPtr.Get(key); ok {
				result[key] = normalizeValue(val)
			}
		}
		return result
	}

	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Slice:
		result := make([]interface{}, val.Len())
		for i := 0; i < val.Len(); i++ {
			result[i] = normalizeValue(val.Index(i).Interface())
		}
		return result

	case reflect.Map:
		result := make(map[string]interface{})
		for _, key := range val.MapKeys() {
			keyStr := key.String()
			result[keyStr] = normalizeValue(val.MapIndex(key).Interface())
		}
		return result

	default:
		return v
	}
}