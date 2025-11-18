package toon

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMarshalSimplePrimitive(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "null",
			input:    nil,
			expected: "null",
		},
		{
			name:     "true",
			input:    true,
			expected: "true",
		},
		{
			name:     "false",
			input:    false,
			expected: "false",
		},
		{
			name:     "integer",
			input:    42,
			expected: "42",
		},
		{
			name:     "float",
			input:    3.14,
			expected: "3.14",
		},
		{
			name:     "string",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "string with space",
			input:    "hello world",
			expected: "hello world", // Root primitives don't need quoting
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Marshal(tt.input, nil)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}
			if string(result) != tt.expected {
				t.Errorf("Marshal() = %q, want %q", string(result), tt.expected)
			}
		})
	}
}

func TestMarshalSimpleObject(t *testing.T) {
	input := map[string]interface{}{
		"name": "Alice",
		"age":  30,
	}

	result, err := Marshal(input, nil)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	expected := "age: 30\nname: Alice"
	if string(result) != expected {
		t.Errorf("Marshal() = %q, want %q", string(result), expected)
	}
}

func TestMarshalInlineArray(t *testing.T) {
	input := map[string]interface{}{
		"tags": []interface{}{"go", "toon"},
	}

	result, err := Marshal(input, nil)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	expected := "tags[2]: go,toon"
	if string(result) != expected {
		t.Errorf("Marshal() = %q, want %q", string(result), expected)
	}
}

func TestMarshalNestedObject(t *testing.T) {
	input := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "Bob",
		},
	}

	result, err := Marshal(input, nil)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	expected := "user:\n  name: Bob"
	if string(result) != expected {
		t.Errorf("Marshal() = %q, want %q", string(result), expected)
	}
}

func TestMarshalTabularArray(t *testing.T) {
	input := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"name": "Alice", "age": 30},
			map[string]interface{}{"name": "Bob", "age": 25},
		},
	}

	result, err := Marshal(input, nil)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	// Keys should be sorted: age, name
	expected := "users[2]{age,name}:\n  30,Alice\n  25,Bob"
	if string(result) != expected {
		t.Errorf("Marshal() = %q, want %q", string(result), expected)
	}
}

func TestMarshalWithOptions(t *testing.T) {
	input := map[string]interface{}{
		"tags": []interface{}{"a", "b", "c"},
	}

	opts := &EncodeOptions{
		Delimiter:    "\t",
		LengthMarker: "#",
	}

	result, err := Marshal(input, opts)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	expected := "tags[#3\t]: a\tb\tc"
	if string(result) != expected {
		t.Errorf("Marshal() = %q, want %q", string(result), expected)
	}
}

func TestMarshalEmptyArray(t *testing.T) {
	input := map[string]interface{}{
		"items": []interface{}{},
	}

	result, err := Marshal(input, nil)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	expected := "items[0]:"
	if string(result) != expected {
		t.Errorf("Marshal() = %q, want %q", string(result), expected)
	}
}
// TestEncodeFixtures runs all official TOON specification encode tests
func TestEncodeFixtures(t *testing.T) {
	fixtureDir := "../testdata/fixtures/encode"
	
	entries, err := os.ReadDir(fixtureDir)
	if err != nil {
		t.Fatalf("Failed to read fixture directory: %v", err)
	}

	totalTests := 0
	passedTests := 0

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".json" {
			fixturePath := filepath.Join(fixtureDir, entry.Name())
			
			fixture, err := loadFixture(fixturePath)
			if err != nil {
				t.Errorf("Failed to load fixture %s: %v", entry.Name(), err)
				continue
			}

			t.Run(entry.Name(), func(t *testing.T) {
				for _, test := range fixture.Tests {
					totalTests++
					
					t.Run(test.Name, func(t *testing.T) {
						// Convert fixture options to encode options
						opts := fixtureOptionsToEncodeOptions(test.Options)

						// Encode the input
						encoded, err := Marshal(test.Input, opts)
						if err != nil {
							t.Errorf("Marshal() error = %v", err)
							return
						}

						encodedStr := string(encoded)
						expectedStr := test.Expected.(string)

						// For semantic comparison, decode both and compare structures
						var encodedResult interface{}
						var expectedResult interface{}

						decodeErr1 := Unmarshal([]byte(encodedStr), &encodedResult, nil)
						decodeErr2 := Unmarshal([]byte(expectedStr), &expectedResult, nil)

						if decodeErr1 == nil && decodeErr2 == nil {
							// Compare decoded structures for semantic equivalence
							encodedNorm := normalizeValue(encodedResult)
							expectedNorm := normalizeValue(expectedResult)

							if deepEqual(encodedNorm, expectedNorm) {
								passedTests++
								return
							}
						}

						// If semantic comparison fails or decode errors, do string comparison
						if encodedStr != expectedStr {
							t.Errorf("Marshal() mismatch\nSpec Section: %s\nInput: %#v\nExpected: %q\nGot: %q\nNote: %s",
								test.SpecSection, test.Input, expectedStr, encodedStr, test.Note)
						} else {
							passedTests++
						}
					})
				}
			})
		}
	}

	t.Logf("Encode Fixtures: %d/%d tests passed", passedTests, totalTests)
}

// TestAssignResultMapTarget tests assigning decoded result to *map[string]interface{}
func TestAssignResultMapTarget(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantErr     bool
		errContains string
		validate    func(t *testing.T, result map[string]interface{})
	}{
		{
			name:    "simple map",
			input:   "name: Alice\nage: 30",
			wantErr: false,
			validate: func(t *testing.T, result map[string]interface{}) {
				if result["name"] != "Alice" {
					t.Errorf("name = %v, want Alice", result["name"])
				}
				if result["age"] != int64(30) {
					t.Errorf("age = %v, want 30", result["age"])
				}
			},
		},
		{
			name:    "nested map",
			input:   "user:\n  name: Bob\n  age: 25",
			wantErr: false,
			validate: func(t *testing.T, result map[string]interface{}) {
				// assignResult converts map[string]Value to map[string]interface{}
				// so nested maps stay as map[string]Value (which is interface{})
				user, ok := result["user"].(map[string]Value)
				if !ok {
					t.Fatalf("user is not a map[string]Value, got %T", result["user"])
				}
				if user["name"] != "Bob" {
					t.Errorf("user.name = %v, want Bob", user["name"])
				}
			},
		},
		{
			name:    "empty map",
			input:   "",
			wantErr: false,
			validate: func(t *testing.T, result map[string]interface{}) {
				if len(result) != 0 {
					t.Errorf("expected empty map, got %v", result)
				}
			},
		},
		{
			name:    "map with array value",
			input:   "tags[2]: go,toon",
			wantErr: false,
			validate: func(t *testing.T, result map[string]interface{}) {
				tags, ok := result["tags"].([]Value)
				if !ok {
					t.Fatalf("tags is not []Value, got %T", result["tags"])
				}
				if len(tags) != 2 {
					t.Errorf("len(tags) = %d, want 2", len(tags))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result map[string]interface{}
			err := Unmarshal([]byte(tt.input), &result, nil)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errContains != "" && err.Error() != tt.errContains {
					t.Errorf("error = %q, want to contain %q", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

// TestAssignResultArrayTarget tests assigning decoded result to *[]interface{}
// Note: TOON format always decodes to maps at the root level, arrays are values within maps
// So we test the error case where trying to assign a map to an array target
func TestAssignResultArrayTarget(t *testing.T) {
	t.Run("map result to array target should error", func(t *testing.T) {
		input := "name: Alice"
		var result []interface{}
		err := Unmarshal([]byte(input), &result, nil)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "cannot assign non-array to array target" {
			t.Errorf("error = %q, want 'cannot assign non-array to array target'", err.Error())
		}
	})

	t.Run("map with array to array target should error", func(t *testing.T) {
		input := "tags[2]: go,toon"
		var result []interface{}
		err := Unmarshal([]byte(input), &result, nil)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "cannot assign non-array to array target" {
			t.Errorf("error = %q, want 'cannot assign non-array to array target'", err.Error())
		}
	})
}

// TestAssignResultInterfaceTarget tests assigning decoded result to *interface{}
func TestAssignResultInterfaceTarget(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		validate func(t *testing.T, result interface{})
	}{
		{
			name:    "map to interface{}",
			input:   "name: Alice\nage: 30",
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				m, ok := result.(map[string]Value)
				if !ok {
					t.Fatalf("result is not a map[string]Value, got %T", result)
				}
				if m["name"] != "Alice" {
					t.Errorf("name = %v, want Alice", m["name"])
				}
			},
		},
		{
			name:    "array to interface{}",
			input:   "tags[2]: go,toon",
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				m, ok := result.(map[string]Value)
				if !ok {
					t.Fatalf("result is not a map[string]Value, got %T", result)
				}
				arr, ok := m["tags"].([]Value)
				if !ok {
					t.Fatalf("tags is not []Value, got %T", m["tags"])
				}
				if len(arr) != 2 {
					t.Errorf("tags length = %d, want 2", len(arr))
				}
			},
		},
		{
			name:    "primitive to interface{}",
			input:   "42",
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result != int64(42) {
					t.Errorf("result = %v, want 42", result)
				}
			},
		},
		{
			name:    "string to interface{}",
			input:   "hello",
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result != "hello" {
					t.Errorf("result = %v, want hello", result)
				}
			},
		},
		{
			name:    "boolean to interface{}",
			input:   "true",
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result != true {
					t.Errorf("result = %v, want true", result)
				}
			},
		},
		{
			name:    "null to interface{}",
			input:   "null",
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				if result != nil {
					t.Errorf("result = %v, want nil", result)
				}
			},
		},
		{
			name:    "nested structure to interface{}",
			input:   "user:\n  profile:\n    name: Charlie",
			wantErr: false,
			validate: func(t *testing.T, result interface{}) {
				m, ok := result.(map[string]Value)
				if !ok {
					t.Fatalf("result is not a map, got %T", result)
				}
				user, ok := m["user"].(map[string]Value)
				if !ok {
					t.Fatal("user is not a map")
				}
				profile, ok := user["profile"].(map[string]Value)
				if !ok {
					t.Fatal("profile is not a map")
				}
				if profile["name"] != "Charlie" {
					t.Errorf("profile.name = %v, want Charlie", profile["name"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result interface{}
			err := Unmarshal([]byte(tt.input), &result, nil)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

// TestAssignResultUnsupportedTarget tests error cases for unsupported target types
func TestAssignResultUnsupportedTarget(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		target      interface{}
		errContains string
	}{
		{
			name:        "non-pointer target",
			input:       "name: Alice",
			target:      map[string]interface{}{},
			errContains: "unsupported target type",
		},
		{
			name:        "string target",
			input:       "name: Alice",
			target:      new(string),
			errContains: "unsupported target type",
		},
		{
			name:        "int target",
			input:       "42",
			target:      new(int),
			errContains: "unsupported target type",
		},
		{
			name:        "bool target",
			input:       "true",
			target:      new(bool),
			errContains: "unsupported target type",
		},
		{
			name:        "struct target",
			input:       "name: Alice",
			target:      &struct{ Name string }{},
			errContains: "unsupported target type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Unmarshal([]byte(tt.input), tt.target, nil)

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if tt.errContains != "" && err.Error() != tt.errContains {
				t.Errorf("error = %q, want to contain %q", err.Error(), tt.errContains)
			}
		})
	}
}

// TestAssignResultComplexStructures tests complex nested structures
func TestAssignResultComplexStructures(t *testing.T) {
	input := `users[2]{age,name}:
  30,Alice
  25,Bob
metadata:
  version: 1.0
  tags[3]: go,toon,format`

	t.Run("complex to map", func(t *testing.T) {
		var result map[string]interface{}
		err := Unmarshal([]byte(input), &result, &DecodeOptions{Strict: false})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Check users array - assignResult converts to []interface{} from []Value
		users, ok := result["users"].([]Value)
		if !ok {
			t.Fatalf("users is not []Value, got %T", result["users"])
		}
		if len(users) != 2 {
			t.Errorf("users length = %d, want 2", len(users))
		}

		// Check metadata - nested maps stay as map[string]Value
		metadata, ok := result["metadata"].(map[string]Value)
		if !ok {
			t.Fatalf("metadata is not map[string]Value, got %T", result["metadata"])
		}
		// version is parsed as float64 1.0
		if metadata["version"] != 1.0 {
			t.Errorf("version = %v (type %T), want 1.0", metadata["version"], metadata["version"])
		}
	})

	t.Run("complex to interface{}", func(t *testing.T) {
		var result interface{}
		err := Unmarshal([]byte(input), &result, &DecodeOptions{Strict: false})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		m, ok := result.(map[string]Value)
		if !ok {
			t.Fatalf("result is not map[string]Value, got %T", result)
		}

		if _, ok := m["users"]; !ok {
			t.Error("users key not found")
		}
		if _, ok := m["metadata"]; !ok {
			t.Error("metadata key not found")
		}
	})
}

// TestAssignResultEmptyAndNil tests edge cases with empty and nil values
func TestAssignResultEmptyAndNil(t *testing.T) {
	t.Run("empty input to map", func(t *testing.T) {
		var result map[string]interface{}
		err := Unmarshal([]byte(""), &result, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 0 {
			t.Errorf("expected empty map, got %v", result)
		}
	})

	t.Run("empty input to interface{}", func(t *testing.T) {
		var result interface{}
		err := Unmarshal([]byte(""), &result, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		m, ok := result.(map[string]Value)
		if !ok {
			t.Fatalf("result is not map, got %T", result)
		}
		if len(m) != 0 {
			t.Errorf("expected empty map, got %v", m)
		}
	})

	t.Run("null value to interface{}", func(t *testing.T) {
		var result interface{}
		err := Unmarshal([]byte("null"), &result, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result != nil {
			t.Errorf("expected nil, got %v", result)
		}
	})
}

// TestAssignResultDirectArrayAssignment tests direct assignment of arrays
// This tests the successful branch of *[]interface{} case in assignResult
func TestAssignResultDirectArrayAssignment(t *testing.T) {
	t.Run("assign array value to array target", func(t *testing.T) {
		// Create an array value directly
		arrayValue := []Value{"item1", "item2", int64(3)}
		
		// Test direct assignment
		var target []interface{}
		err := assignResult(arrayValue, &target)
		
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		
		if len(target) != 3 {
			t.Errorf("expected length 3, got %d", len(target))
		}
		
		if target[0] != "item1" {
			t.Errorf("target[0] = %v, want item1", target[0])
		}
		if target[1] != "item2" {
			t.Errorf("target[1] = %v, want item2", target[1])
		}
		if target[2] != int64(3) {
			t.Errorf("target[2] = %v, want 3", target[2])
		}
	})

	t.Run("assign empty array to array target", func(t *testing.T) {
		arrayValue := []Value{}
		
		var target []interface{}
		err := assignResult(arrayValue, &target)
		
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		
		if len(target) != 0 {
			t.Errorf("expected empty array, got length %d", len(target))
		}
	})

	t.Run("assign nested array to array target", func(t *testing.T) {
		// Array with nested structures
		nestedArray := []Value{
			map[string]Value{"key": "value"},
			[]Value{"nested", "array"},
			int64(42),
		}
		
		var target []interface{}
		err := assignResult(nestedArray, &target)
		
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		
		if len(target) != 3 {
			t.Errorf("expected length 3, got %d", len(target))
		}
		
		// Check nested map
		nestedMap, ok := target[0].(map[string]Value)
		if !ok {
			t.Errorf("target[0] is not map[string]Value, got %T", target[0])
		} else if nestedMap["key"] != "value" {
			t.Errorf("nestedMap[key] = %v, want value", nestedMap["key"])
		}
		
		// Check nested array
		nestedArr, ok := target[1].([]Value)
		if !ok {
			t.Errorf("target[1] is not []Value, got %T", target[1])
		} else if len(nestedArr) != 2 {
			t.Errorf("nested array length = %d, want 2", len(nestedArr))
		}
	})
}

// TestAssignResultDirectMapAssignment tests direct assignment of maps
// This ensures we test the map conversion loop thoroughly
func TestAssignResultDirectMapAssignment(t *testing.T) {
	t.Run("assign map value to map target", func(t *testing.T) {
		// Create a map value directly
		mapValue := map[string]Value{
			"string": "value",
			"number": int64(42),
			"bool":   true,
			"null":   nil,
		}
		
		// Test direct assignment
		var target map[string]interface{}
		err := assignResult(mapValue, &target)
		
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		
		if len(target) != 4 {
			t.Errorf("expected length 4, got %d", len(target))
		}
		
		if target["string"] != "value" {
			t.Errorf("target[string] = %v, want value", target["string"])
		}
		if target["number"] != int64(42) {
			t.Errorf("target[number] = %v, want 42", target["number"])
		}
		if target["bool"] != true {
			t.Errorf("target[bool] = %v, want true", target["bool"])
		}
		if target["null"] != nil {
			t.Errorf("target[null] = %v, want nil", target["null"])
		}
	})

	t.Run("assign empty map to map target", func(t *testing.T) {
		mapValue := map[string]Value{}
		
		var target map[string]interface{}
		err := assignResult(mapValue, &target)
		
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		
		if len(target) != 0 {
			t.Errorf("expected empty map, got length %d", len(target))
		}
	})

	t.Run("assign nested map to map target", func(t *testing.T) {
		// Map with nested structures
		nestedMap := map[string]Value{
			"nested_map":   map[string]Value{"inner": "value"},
			"nested_array": []Value{"a", "b", "c"},
			"primitive":    "simple",
		}
		
		var target map[string]interface{}
		err := assignResult(nestedMap, &target)
		
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		
		if len(target) != 3 {
			t.Errorf("expected length 3, got %d", len(target))
		}
		
		// Check nested map
		innerMap, ok := target["nested_map"].(map[string]Value)
		if !ok {
			t.Errorf("nested_map is not map[string]Value, got %T", target["nested_map"])
		} else if innerMap["inner"] != "value" {
			t.Errorf("innerMap[inner] = %v, want value", innerMap["inner"])
		}
		
		// Check nested array
		innerArr, ok := target["nested_array"].([]Value)
		if !ok {
			t.Errorf("nested_array is not []Value, got %T", target["nested_array"])
		} else if len(innerArr) != 3 {
			t.Errorf("inner array length = %d, want 3", len(innerArr))
		}
	})
}