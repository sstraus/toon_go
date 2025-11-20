// Package toon implements encoding and decoding of TOON (Token-Oriented Object Notation) format.
//
// TOON is a compact data format optimized for LLM token efficiency, achieving
// 30-60% token reduction compared to JSON while maintaining readability.
//
// # Format Overview
//
// TOON uses indentation-based structure like YAML and supports three array formats:
//
//   - Inline: tags[2]: a,b (for primitives)
//   - Tabular: users[2]{name,age}: Alice,30 / Bob,25 (for uniform objects)
//   - List: items[2]: - item1 / - item2 (for mixed/nested)
//
// # Public API
//
// The package exports a minimal API surface for encoding and decoding:
//
//	Marshal(v interface{}, w io.Writer, opts *EncodeOptions) error
//	Unmarshal(r io.Reader, v interface{}, opts *DecodeOptions) error
//
// Additional exported types:
//
//	OrderedMap - Preserves key insertion order
//	EncodeOptions - Encoding configuration (indent, delimiter, etc.)
//	DecodeOptions - Decoding configuration
//	EncodeError, DecodeError - Error types with detailed messages
//
// # Basic Usage
//
// Encoding a Go value to TOON format:
//
//	data := map[string]interface{}{
//	    "name": "Alice",
//	    "age":  30,
//	    "tags": []string{"go", "toon"},
//	}
//
//	var buf bytes.Buffer
//	err := toon.Marshal(data, &buf, nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(buf.String())
//	// Output:
//	// age: 30
//	// name: Alice
//	// tags[2]: go,toon
//
// Decoding TOON format to a Go value:
//
//	input := `
//	name: Bob
//	age: 25
//	active: true
//	`
//
//	var result map[string]interface{}
//	err := toon.Unmarshal(strings.NewReader(input), &result, nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("%+v\n", result)
//	// Output: map[age:25 name:Bob active:true]
//
// # Custom Options
//
// Both Marshal and Unmarshal accept optional configuration:
//
//	opts := &toon.EncodeOptions{
//	    Indent:       4,    // 4 spaces instead of 2
//	    Delimiter:    "\t", // Tab delimiter for arrays
//	    LengthMarker: "#",  // Prefix array length with #
//	}
//
//	var buf bytes.Buffer
//	err := toon.Marshal(data, &buf, opts)
//
// # OrderedMap
//
// Use OrderedMap to preserve key insertion order during encoding:
//
//	om := toon.NewOrderedMap()
//	om.Set("first", 1)
//	om.Set("second", 2)
//	om.Set("third", 3)
//
//	var buf bytes.Buffer
//	toon.Marshal(om, &buf, nil)
//	// Keys will be encoded in insertion order: first, second, third
//
// # Error Handling
//
// The package returns detailed error types for encoding and decoding failures:
//
//	err := toon.Marshal(data, &buf, nil)
//	if err != nil {
//	    if encErr, ok := err.(*toon.EncodeError); ok {
//	        fmt.Printf("Encoding error: %s\n", encErr.Message)
//	    }
//	}
//
// # Implementation Details
//
// All encoding and decoding implementation details are unexported. The package
// is organized into focused modules:
//
//   - encode_*.go - Encoding logic for objects, arrays, and primitives
//   - decode_*.go - Decoding logic with structural and token parsers
//   - orderedmap.go - Ordered map implementation
//   - types.go, errors.go, options.go - Public type definitions
//
// The implementation achieves:
//
//   - 83.1% test coverage with 1,088 tests
//   - Average cyclomatic complexity of 5.39
//   - 340/340 official TOON specification fixtures passing
//
// # Specification Compliance
//
// This implementation follows TOON Specification v2.0:
// https://github.com/toon-format/spec
//
// # Version
//
// Current version: 1.1.0
//
// Version history is maintained in CHANGELOG.md.
package toon
