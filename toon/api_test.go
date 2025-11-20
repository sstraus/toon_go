package toon

import (
	"bytes"
	"strings"
	"testing"
)

// TestMarshalToString tests the MarshalToString convenience function.
func TestMarshalToString(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		opts     []EncodeOption
		expected string
		wantErr  bool
	}{
		{
			name: "simple map",
			input: map[string]interface{}{
				"name": "Alice",
				"age":  30,
			},
			expected: "age: 30\nname: Alice",
		},
		{
			name: "with custom indent",
			input: map[string]interface{}{
				"nested": map[string]interface{}{
					"value": 42,
				},
			},
			opts:     []EncodeOption{WithIndent(4)},
			expected: "nested:\n    value: 42",
		},
		{
			name: "with custom delimiter",
			input: map[string]interface{}{
				"values": []interface{}{1, 2, 3},
			},
			opts:     []EncodeOption{WithDelimiter("|")},
			expected: "values[3|]: 1|2|3",
		},
		{
			name: "ordered map",
			input: func() *OrderedMap {
				om := NewOrderedMap()
				om.Set("first", 1)
				om.Set("second", 2)
				return om
			}(),
			expected: "first: 1\nsecond: 2",
		},
		{
			name:     "nil value",
			input:    nil,
			expected: "null",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "\"\"",
		},
		{
			name:     "number",
			input:    42,
			expected: "42",
		},
		{
			name:     "boolean",
			input:    true,
			expected: "true",
		},
		{
			name: "array of primitives",
			input: map[string]interface{}{
				"items": []interface{}{"a", "b", "c"},
			},
			expected: "items[3]: a,b,c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MarshalToString(tt.input, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("MarshalToString() =\n%s\nexpected:\n%s", result, tt.expected)
			}
		})
	}
}

// TestUnmarshalFromString tests the UnmarshalFromString convenience function.
func TestUnmarshalFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		opts     []DecodeOption
		expected interface{}
		wantErr  bool
	}{
		{
			name:  "simple object",
			input: "name: Alice\nage: 30",
			expected: map[string]interface{}{
				"name": "Alice",
				"age":  float64(30),
			},
		},
		{
			name:  "nested object",
			input: "person:\n  name: Bob\n  age: 25",
			expected: map[string]interface{}{
				"person": map[string]interface{}{
					"name": "Bob",
					"age":  float64(25),
				},
			},
		},
		{
			name:  "inline array",
			input: "values[3]: 1,2,3",
			expected: map[string]interface{}{
				"values": []interface{}{float64(1), float64(2), float64(3)},
			},
		},
		{
			name:     "null value",
			input:    "null",
			expected: nil,
		},
		{
			name:     "boolean true",
			input:    "true",
			expected: true,
		},
		{
			name:     "boolean false",
			input:    "false",
			expected: false,
		},
		{
			name:     "number",
			input:    "42",
			expected: float64(42),
		},
		{
			name:     "string",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "quoted string",
			input:    "\"hello world\"",
			expected: "hello world",
		},
		{
			name: "list array",
			input: `items[2]:
  - name: Item 1
  - name: Item 2`,
			expected: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{"name": "Item 1"},
					map[string]interface{}{"name": "Item 2"},
				},
			},
		},
		{
			name: "multiple colons",
			input: "key: value: extra",
			expected: map[string]interface{}{
				"key": "value: extra",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result interface{}
			err := UnmarshalFromString(tt.input, &result, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !deepEqual(result, tt.expected) {
				t.Errorf("UnmarshalFromString() =\n%#v\nexpected:\n%#v", result, tt.expected)
			}
		})
	}
}

// TestMarshalUnmarshalRoundTrip tests round-trip conversion.
func TestMarshalUnmarshalRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{
			name: "simple object",
			input: map[string]interface{}{
				"name":   "Charlie",
				"age":    35,
				"active": true,
			},
		},
		{
			name: "nested object",
			input: map[string]interface{}{
				"user": map[string]interface{}{
					"name": "Dave",
					"profile": map[string]interface{}{
						"bio": "Developer",
					},
				},
			},
		},
		{
			name: "with arrays",
			input: map[string]interface{}{
				"tags":   []interface{}{"go", "toon", "test"},
				"scores": []interface{}{95, 87, 92},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to string
			encoded, err := MarshalToString(tt.input)
			if err != nil {
				t.Fatalf("MarshalToString() error = %v", err)
			}

			// Unmarshal back
			var decoded interface{}
			err = UnmarshalFromString(encoded, &decoded)
			if err != nil {
				t.Fatalf("UnmarshalFromString() error = %v", err)
			}

			// Normalize and compare
			normalized := normalize(tt.input)
			if !deepEqual(decoded, normalized) {
				t.Errorf("Round-trip failed:\nOriginal: %#v\nDecoded:  %#v", normalized, decoded)
			}
		})
	}
}

// TestMarshalVsMarshalToString ensures Marshal and MarshalToString produce the same output.
func TestMarshalVsMarshalToString(t *testing.T) {
	input := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
		"key3": []interface{}{1, 2, 3},
	}

	// Marshal to buffer
	var buf bytes.Buffer
	err := Marshal(input, &buf)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	bufResult := buf.String()

	// MarshalToString
	strResult, err := MarshalToString(input)
	if err != nil {
		t.Fatalf("MarshalToString() error = %v", err)
	}

	if bufResult != strResult {
		t.Errorf("Marshal and MarshalToString produce different output:\nMarshal:         %s\nMarshalToString: %s", bufResult, strResult)
	}
}

// TestUnmarshalVsUnmarshalFromString ensures Unmarshal and UnmarshalFromString produce the same output.
func TestUnmarshalVsUnmarshalFromString(t *testing.T) {
	input := "key1: value1\nkey2: 42\nkey3[3]: 1,2,3"

	// Unmarshal from reader
	var result1 interface{}
	err := Unmarshal(strings.NewReader(input), &result1)
	if err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	// UnmarshalFromString
	var result2 interface{}
	err = UnmarshalFromString(input, &result2)
	if err != nil {
		t.Fatalf("UnmarshalFromString() error = %v", err)
	}

	if !deepEqual(result1, result2) {
		t.Errorf("Unmarshal and UnmarshalFromString produce different output:\nUnmarshal:           %#v\nUnmarshalFromString: %#v", result1, result2)
	}
}

// TestEncodeOptionsCoverage tests functional options for encoding.
func TestEncodeOptionsCoverage(t *testing.T) {
	input := map[string]interface{}{
		"items": []interface{}{1, 2, 3},
	}

	tests := []struct {
		name string
		opts []EncodeOption
	}{
		{
			name: "WithLengthMarker",
			opts: []EncodeOption{WithLengthMarker("#")},
		},
		{
			name: "WithFlattenPaths",
			opts: []EncodeOption{WithFlattenPaths(true)},
		},
		{
			name: "WithFlattenDepth",
			opts: []EncodeOption{WithFlattenDepth(2)},
		},
		{
			name: "WithStrict",
			opts: []EncodeOption{WithStrict(true)},
		},
		{
			name: "Multiple options",
			opts: []EncodeOption{
				WithIndent(4),
				WithDelimiter("|"),
				WithLengthMarker("COUNT:"),
				WithFlattenPaths(false),
				WithFlattenDepth(3),
				WithStrict(false),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := MarshalToString(input, tt.opts...)
			if err != nil {
				t.Errorf("MarshalToString() with %s error = %v", tt.name, err)
			}
		})
	}
}

// TestDecodeOptionsCoverage tests functional options for decoding.
func TestDecodeOptionsCoverage(t *testing.T) {
	input := "key: value"

	tests := []struct {
		name string
		opts []DecodeOption
	}{
		{
			name: "WithKeyMode",
			opts: []DecodeOption{WithKeyMode(StringKeys)},
		},
		{
			name: "WithStrictDecoding",
			opts: []DecodeOption{WithStrictDecoding(false)},
		},
		{
			name: "WithIndentSize",
			opts: []DecodeOption{WithIndentSize(4)},
		},
		{
			name: "WithExpandPaths",
			opts: []DecodeOption{WithExpandPaths("safe")},
		},
		{
			name: "Multiple options",
			opts: []DecodeOption{
				WithKeyMode(StringKeys),
				WithStrictDecoding(true),
				WithIndentSize(2),
				WithExpandPaths("off"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result interface{}
			err := UnmarshalFromString(input, &result, tt.opts...)
			if err != nil {
				t.Errorf("UnmarshalFromString() with %s error = %v", tt.name, err)
			}
		})
	}
}

