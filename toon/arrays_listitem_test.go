package toon

import (
	"testing"
)

// TestEncodeListItemComplexNestedArrays tests complex nested array structures
// to improve coverage of encodeListItem function (lines 366-384)
func TestEncodeListItemComplexNestedArrays(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name: "array of arrays with 3+ levels of nesting",
			input: map[string]interface{}{
				"data": []interface{}{
					[]interface{}{
						[]interface{}{1, 2},
						[]interface{}{3, 4},
					},
					[]interface{}{
						[]interface{}{5, 6},
					},
				},
			},
			expected: "data[2]:\n  - [2]:\n    - [2]: 1,2\n    - [2]: 3,4\n  - [1]:\n    - [2]: 5,6",
		},
		{
			name: "array containing objects with nested arrays",
			input: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{
						"id":   1,
						"tags": []interface{}{"a", "b"},
					},
					map[string]interface{}{
						"id":   2,
						"tags": []interface{}{"c", "d", "e"},
					},
				},
			},
			expected: "items[2]:\n  - tags[2]: a,b\n    id: 1\n  - tags[3]: c,d,e\n    id: 2",
		},
		{
			name: "deeply nested mixed structures (arrays with maps with arrays)",
			input: map[string]interface{}{
				"complex": []interface{}{
					map[string]interface{}{
						"level1": []interface{}{
							map[string]interface{}{
								"level2": []interface{}{"deep1", "deep2"},
							},
						},
					},
				},
			},
			expected: "complex[1]:\n  - level1[1]:\n    - level2[2]: deep1,deep2",
		},
		{
			name: "array with empty nested arrays at various levels",
			input: map[string]interface{}{
				"nested": []interface{}{
					[]interface{}{},
					[]interface{}{[]interface{}{}},
					[]interface{}{[]interface{}{}, []interface{}{}},
				},
			},
			expected: "nested[3]:\n  - [0]:\n  - [1]:\n    - [0]:\n  - [2]:\n    - [0]:\n    - [0]:",
		},
		{
			name: "list items with nested arrays containing primitives vs complex objects",
			input: map[string]interface{}{
				"mixed": []interface{}{
					[]interface{}{1, 2, 3},
					[]interface{}{
						map[string]interface{}{"a": 1},
						map[string]interface{}{"b": 2},
					},
				},
			},
			expected: "mixed[2]:\n  - [3]: 1,2,3\n  - [2]:\n    - a: 1\n    - b: 2",
		},
		{
			name: "object in list with array field on first key (lines 375-379)",
			input: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{
						"arr": []interface{}{1, 2},
						"id":  1,
					},
				},
			},
			expected: "items[1]:\n  - arr[2]: 1,2\n    id: 1",
		},
		{
			name: "object in list with complex nested object on first key (lines 380-386)",
			input: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{
						"nested": map[string]interface{}{
							"inner": "value",
						},
						"id": 1,
					},
				},
			},
			expected: "items[1]:\n  - id: 1\n    nested:\n      inner: value",
		},
		{
			name: "object in list with array field on subsequent key (lines 403-408)",
			input: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{
						"id":   1,
						"tags": []interface{}{"a", "b"},
					},
				},
			},
			expected: "items[1]:\n  - tags[2]: a,b\n    id: 1",
		},
		{
			name: "object in list with nested object on subsequent key (lines 410-415)",
			input: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{
						"id": 1,
						"meta": map[string]interface{}{
							"created": "2024-01-01",
						},
					},
				},
			},
			expected: "items[1]:\n  - id: 1\n    meta:\n      created: 2024-01-01",
		},
		{
			name: "4 levels of array nesting",
			input: map[string]interface{}{
				"deep": []interface{}{
					[]interface{}{
						[]interface{}{
							[]interface{}{1, 2},
						},
					},
				},
			},
			expected: "deep[1]:\n  - [1]:\n    - [1]:\n      - [2]: 1,2",
		},
		{
			name: "mixed array-object nesting with multiple levels",
			input: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"items": []interface{}{
							map[string]interface{}{
								"values": []interface{}{1, 2, 3},
							},
						},
					},
				},
			},
			expected: "data[1]:\n  - items[1]:\n    - values[3]: 1,2,3",
		},
		{
			name: "array with nil values in nested context",
			input: map[string]interface{}{
				"items": []interface{}{
					[]interface{}{nil, "text", nil},
					[]interface{}{nil},
				},
			},
			expected: "items[2]:\n  - [3]: null,text,null\n  - [1]: null",
		},
		{
			name: "empty nested arrays mixed with populated ones",
			input: map[string]interface{}{
				"data": []interface{}{
					[]interface{}{},
					[]interface{}{1},
					[]interface{}{},
					[]interface{}{2, 3},
				},
			},
			expected: "data[4]:\n  - [0]:\n  - [1]: 1\n  - [0]:\n  - [2]: 2,3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Marshal(tt.input, nil)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}
			got := string(result)
			if got != tt.expected {
				t.Errorf("Marshal() output mismatch:\nGot:\n%s\n\nWant:\n%s", got, tt.expected)
			}
		})
	}
}

// TestEncodeListItemDelimiterHandling tests delimiter handling in nested contexts
func TestEncodeListItemDelimiterHandling(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		delimiter string
		expected  string
	}{
		{
			name: "tab delimiter with nested arrays",
			input: map[string]interface{}{
				"data": []interface{}{
					[]interface{}{"a", "b"},
					[]interface{}{"c", "d"},
				},
			},
			delimiter: "\t",
			expected:  "data[2\t]:\n  - [2\t]: a\tb\n  - [2\t]: c\td",
		},
		{
			name: "pipe delimiter with complex nesting",
			input: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{
						"tags": []interface{}{"x", "y", "z"},
					},
				},
			},
			delimiter: "|",
			expected:  "items[1|]:\n  - tags[3|]: x|y|z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &EncodeOptions{
				Delimiter: tt.delimiter,
			}
			result, err := Marshal(tt.input, opts)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}
			got := string(result)
			if got != tt.expected {
				t.Errorf("Marshal() output mismatch:\nGot:\n%s\n\nWant:\n%s", got, tt.expected)
			}
		})
	}
}

// TestEncodeListItemNullEncoding tests how unsupported types encode as null
func TestEncodeListItemNullEncoding(t *testing.T) {
	// Note: Go's reflection treats channels and functions as nil when marshaled
	// These test cases verify the actual behavior - they encode as null rather than error
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name: "channel type in nested array encodes as null",
			input: map[string]interface{}{
				"items": []interface{}{
					[]interface{}{
						make(chan int), // Encoded as null
					},
				},
			},
			expected: "items[1]:\n  - [1]: null",
		},
		{
			name: "function type in list item encodes as null",
			input: map[string]interface{}{
				"items": []interface{}{
					func() {},
				},
			},
			expected: "items[1]: null",
		},
		{
			name: "channel in deeply nested structure encodes as null",
			input: map[string]interface{}{
				"data": []interface{}{
					[]interface{}{
						[]interface{}{
							make(chan string),
						},
					},
				},
			},
			expected: "data[1]:\n  - [1]:\n    - [1]: null",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Marshal(tt.input, nil)
			if err != nil {
				t.Fatalf("Marshal() unexpected error = %v", err)
			}
			got := string(result)
			if got != tt.expected {
				t.Errorf("Marshal() output mismatch:\nGot:\n%s\n\nWant:\n%s", got, tt.expected)
			}
		})
	}
}

// TestEncodeListItemEdgeCases tests edge cases for list item encoding
func TestEncodeListItemEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name: "single element nested array",
			input: map[string]interface{}{
				"data": []interface{}{
					[]interface{}{
						[]interface{}{42},
					},
				},
			},
			expected: "data[1]:\n  - [1]:\n    - [1]: 42",
		},
		{
			name: "alternating empty and filled nested arrays",
			input: map[string]interface{}{
				"items": []interface{}{
					[]interface{}{},
					[]interface{}{1},
					[]interface{}{},
					[]interface{}{2},
					[]interface{}{},
				},
			},
			expected: "items[5]:\n  - [0]:\n  - [1]: 1\n  - [0]:\n  - [1]: 2\n  - [0]:",
		},
		{
			name: "object with all array fields in list",
			input: map[string]interface{}{
				"items": []interface{}{
					map[string]interface{}{
						"a": []interface{}{1},
						"b": []interface{}{2, 3},
						"c": []interface{}{4, 5, 6},
					},
				},
			},
			expected: "items[1]:\n  - a[1]: 1\n      b[2]: 2,3\n      c[3]: 4,5,6",
		},
		{
			name: "nested array with boolean and string primitives",
			input: map[string]interface{}{
				"flags": []interface{}{
					[]interface{}{true, false, "yes", "no"},
				},
			},
			expected: "flags[1]:\n  - [4]: true,false,yes,no",
		},
		{
			name: "array with mixed number types in nested context",
			input: map[string]interface{}{
				"numbers": []interface{}{
					[]interface{}{int64(1), 2.5, int32(3)},
				},
			},
			expected: "numbers[1]:\n  - [3]: 1,2.5,3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Marshal(tt.input, nil)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}
			got := string(result)
			if got != tt.expected {
				t.Errorf("Marshal() output mismatch:\nGot:\n%s\n\nWant:\n%s", got, tt.expected)
			}
		})
	}
}

// TestEncodeListItemRootLevelArrays tests root-level complex array scenarios
func TestEncodeListItemRootLevelArrays(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name: "root array of nested arrays",
			input: []interface{}{
				[]interface{}{1, 2},
				[]interface{}{
					[]interface{}{3, 4},
					[]interface{}{5, 6},
				},
			},
			expected: "[2]:\n  - [2]: 1,2\n  - [2]:\n    - [2]: 3,4\n    - [2]: 5,6",
		},
		{
			name: "root array with mixed primitives and nested arrays",
			input: []interface{}{
				"text",
				[]interface{}{1, 2},
				42,
				[]interface{}{
					[]interface{}{"a", "b"},
				},
			},
			expected: "[4]:\n  - text\n  - [2]: 1,2\n  - 42\n  - [1]:\n    - [2]: a,b",
		},
		{
			name: "root array with objects containing nested arrays",
			input: []interface{}{
				map[string]interface{}{
					"data": []interface{}{
						[]interface{}{1, 2},
					},
				},
			},
			expected: "[1]:\n  - data[1]:\n    - [2]: 1,2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Marshal(tt.input, nil)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}
			got := string(result)
			if got != tt.expected {
				t.Errorf("Marshal() output mismatch:\nGot:\n%s\n\nWant:\n%s", got, tt.expected)
			}
		})
	}
}

// TestEncodeListItemWithCustomOptions tests list items with various encoding options
func TestEncodeListItemWithCustomOptions(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		opts     *EncodeOptions
		expected string
	}{
		{
			name: "nested arrays with custom indent",
			input: map[string]interface{}{
				"data": []interface{}{
					[]interface{}{
						[]interface{}{1, 2},
					},
				},
			},
			opts: &EncodeOptions{
				Indent: 4,
			},
			expected: "data[1]:\n    - [1]:\n        - [2]: 1,2",
		},
		{
			name: "nested arrays with length marker",
			input: map[string]interface{}{
				"items": []interface{}{
					[]interface{}{1, 2},
					[]interface{}{3, 4, 5},
				},
			},
			opts: &EncodeOptions{
				LengthMarker: "#",
			},
			expected: "items[#2]:\n  - [#2]: 1,2\n  - [#3]: 3,4,5",
		},
		{
			name: "complex nesting with all custom options",
			input: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"tags": []interface{}{"a", "b"},
					},
				},
			},
			opts: &EncodeOptions{
				Indent:       3,
				Delimiter:    "|",
				LengthMarker: "#",
			},
			expected: "data[#1|]:\n   - tags[#2|]: a|b",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Marshal(tt.input, tt.opts)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}
			got := string(result)
			if got != tt.expected {
				t.Errorf("Marshal() output mismatch:\nGot:\n%s\n\nWant:\n%s", got, tt.expected)
			}
		})
	}
}
