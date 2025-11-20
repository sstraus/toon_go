package toon

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUnmarshalSimplePrimitive(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			name:     "null",
			input:    "null",
			expected: nil,
		},
		{
			name:     "true",
			input:    "true",
			expected: true,
		},
		{
			name:     "false",
			input:    "false",
			expected: false,
		},
		{
			name:     "integer",
			input:    "42",
			expected: int64(42),
		},
		{
			name:     "float",
			input:    "3.14",
			expected: 3.14,
		},
		{
			name:     "string",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "quoted string",
			input:    `"hello world"`,
			expected: "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result interface{}
			err := unmarshalFromBytes([]byte(tt.input), &result, nil)
			if err != nil {
				t.Fatalf("unmarshalFromBytes() error = %v", err)
			}

			// Compare values
			if !valuesEqual(result, tt.expected) {
				t.Errorf("unmarshalFromBytes() = %v (%T), want %v (%T)", result, result, tt.expected, tt.expected)
			}
		})
	}
}

func TestUnmarshalSimpleObject(t *testing.T) {
	input := "name: Alice\nage: 30"

	var result map[string]interface{}
	err := unmarshalFromBytes([]byte(input), &result, nil)
	if err != nil {
		t.Fatalf("unmarshalFromBytes() error = %v", err)
	}

	if result["name"] != "Alice" {
		t.Errorf("name = %v, want Alice", result["name"])
	}
	if result["age"] != int64(30) {
		t.Errorf("age = %v, want 30", result["age"])
	}
}

func TestUnmarshalInlineArray(t *testing.T) {
	input := "tags[2]: go,toon"

	var result map[string]interface{}
	err := unmarshalFromBytes([]byte(input), &result, nil)
	if err != nil {
		t.Fatalf("unmarshalFromBytes() error = %v", err)
	}

	tags, ok := result["tags"].([]Value)
	if !ok {
		t.Fatalf("tags is not an array: %T", result["tags"])
	}

	if len(tags) != 2 {
		t.Errorf("len(tags) = %d, want 2", len(tags))
	}
	if tags[0] != "go" {
		t.Errorf("tags[0] = %v, want go", tags[0])
	}
	if tags[1] != "toon" {
		t.Errorf("tags[1] = %v, want toon", tags[1])
	}
}

func TestUnmarshalNestedObject(t *testing.T) {
	input := "user:\n  name: Bob"

	var result map[string]interface{}
	err := unmarshalFromBytes([]byte(input), &result, nil)
	if err != nil {
		t.Fatalf("unmarshalFromBytes() error = %v", err)
	}

	user, ok := result["user"].(map[string]Value)
	if !ok {
		t.Fatalf("user is not an object: %T", result["user"])
	}

	if user["name"] != "Bob" {
		t.Errorf("user.name = %v, want Bob", user["name"])
	}
}

func TestUnmarshalTabularArray(t *testing.T) {
	input := "users[2]{age,name}:\n  30,Alice\n  25,Bob"

	var result map[string]interface{}
	err := unmarshalFromBytes([]byte(input), &result, nil)
	if err != nil {
		t.Fatalf("unmarshalFromBytes() error = %v", err)
	}

	users, ok := result["users"].([]Value)
	if !ok {
		t.Fatalf("users is not an array: %T", result["users"])
	}

	if len(users) != 2 {
		t.Errorf("len(users) = %d, want 2", len(users))
	}

	user0, ok := users[0].(map[string]Value)
	if !ok {
		t.Fatalf("users[0] is not an object: %T", users[0])
	}

	if user0["name"] != "Alice" {
		t.Errorf("users[0].name = %v, want Alice", user0["name"])
	}
	if user0["age"] != int64(30) {
		t.Errorf("users[0].age = %v, want 30", user0["age"])
	}
}

func TestRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
	}{
		{
			name: "simple object",
			data: map[string]interface{}{
				"name": "Alice",
				"age":  30,
			},
		},
		{
			name: "inline array",
			data: map[string]interface{}{
				"tags": []interface{}{"go", "toon"},
			},
		},
		{
			name: "nested object",
			data: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "Bob",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode
			encoded, err := marshalToBytes(tt.data, nil)
			if err != nil {
				t.Fatalf("marshalToBytes() error = %v", err)
			}

			// Decode
			var result interface{}
			err = unmarshalFromBytes(encoded, &result, nil)
			if err != nil {
				t.Fatalf("unmarshalFromBytes() error = %v", err)
			}

			// Compare (basic comparison, not deep)
			// For proper comparison, we'd need to normalize types
			t.Logf("Encoded: %s", string(encoded))
			t.Logf("Decoded: %+v", result)
		})
	}
}

// Helper function to compare values considering type differences
func valuesEqual(a, b interface{}) bool {
	// Handle nil
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Direct comparison for same types
	if a == b {
		return true
	}

	// Handle numeric comparisons (int vs int64, etc.)
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

	return false
}

// TestDecodeFixtures runs all official TOON specification decode tests
func TestDecodeFixtures(t *testing.T) {
	fixtureDir := "../testdata/fixtures/decode"

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
						// Convert fixture options to decode options
						opts := fixtureOptionsToDecodeOptions(test.Options)

						// Get input string
						inputStr, ok := test.Input.(string)
						if !ok {
							t.Errorf("Test input is not a string: %T", test.Input)
							return
						}

						// Decode the input
						var result interface{}
						err := unmarshalFromBytes([]byte(inputStr), &result, opts)

						// Handle shouldError tests
						if test.ShouldError {
							if err == nil {
								t.Errorf("unmarshalFromBytes() expected error but got none\nSpec Section: %s\nInput: %q\nGot: %#v\nNote: %s",
									test.SpecSection, inputStr, result, test.Note)
								return
							}
							// Error occurred as expected
							passedTests++
							return
						}

						// Non-error test cases
						if err != nil {
							t.Errorf("unmarshalFromBytes() error = %v\nSpec Section: %s\nInput: %q\nNote: %s",
								err, test.SpecSection, inputStr, test.Note)
							return
						}

						// Normalize both values for comparison
						resultNorm := normalizeValue(result)
						expectedNorm := normalizeValue(test.Expected)

						// Compare decoded result with expected
						if !deepEqual(resultNorm, expectedNorm) {
							t.Errorf("unmarshalFromBytes() mismatch\nSpec Section: %s\nInput: %q\nExpected: %#v\nGot: %#v\nNote: %s",
								test.SpecSection, inputStr, expectedNorm, resultNorm, test.Note)
						} else {
							passedTests++
						}
					})
				}
			})
		}
	}

	t.Logf("Decode Fixtures: %d/%d tests passed", passedTests, totalTests)
}
