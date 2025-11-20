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
//	Marshal(v interface{}, w io.Writer, opts ...EncodeOption) error
//	Unmarshal(r io.Reader, v interface{}, opts ...DecodeOption) error
//	MarshalToString(v interface{}, opts ...EncodeOption) (string, error)
//	UnmarshalFromString(s string, v interface{}, opts ...DecodeOption) error
//
// Additional exported types:
//
//	OrderedMap - Preserves key insertion order
//	EncodeOptions - Encoding configuration struct (for advanced use)
//	DecodeOptions - Decoding configuration struct (for advanced use)
//	EncodeOption - Functional option for encoding
//	DecodeOption - Functional option for decoding
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
//	// Using io.Writer
//	var buf bytes.Buffer
//	err := toon.Marshal(data, &buf)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(buf.String())
//
//	// Or use string convenience function
//	result, err := toon.MarshalToString(data)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(result)
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
//	// Using io.Reader
//	var result map[string]interface{}
//	err := toon.Unmarshal(strings.NewReader(input), &result)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("%+v\n", result)
//
//	// Or use string convenience function
//	var result2 map[string]interface{}
//	err = toon.UnmarshalFromString(input, &result2)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("%+v\n", result2)
//	// Output: map[age:25 name:Bob active:true]
//
// # Functional Options
//
// Configure encoding and decoding using functional options:
//
//	// Encoding with options
//	result, err := toon.MarshalToString(data,
//	    toon.WithIndent(4),           // 4 spaces instead of 2
//	    toon.WithDelimiter("\t"),     // Tab delimiter for arrays
//	    toon.WithLengthMarker("#"),   // Prefix array length with #
//	)
//
//	// Decoding with options
//	err := toon.UnmarshalFromString(input, &result,
//	    toon.WithStrictDecoding(false),  // Disable strict mode
//	    toon.WithExpandPaths("safe"),    // Expand dotted keys
//	)
//
// Available encoding options:
//
//	WithIndent(n)            - Set indentation size in spaces (default: 2)
//	WithDelimiter(s)         - Set array delimiter: "," | "\t" | "|" (default: ",")
//	WithLengthMarker(s)      - Set length marker prefix (default: "")
//	WithFlattenPaths(bool)   - Enable path flattening (default: false)
//	WithFlattenDepth(n)      - Limit flattening depth (default: 0 = unlimited)
//	WithStrict(bool)         - Enable strict collision detection (default: false)
//
// Available decoding options:
//
//	WithStrictDecoding(bool) - Enable strict validation (default: true)
//	WithIndentSize(n)        - Expected indent size (default: 2)
//	WithExpandPaths(mode)    - Expand dotted keys: "off" | "safe" (default: "off")
//	WithKeyMode(mode)        - Key decoding mode (default: StringKeys)
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
//	result, _ := toon.MarshalToString(om)
//	// Keys will be encoded in insertion order: first, second, third
//
// # Error Handling
//
// The package returns detailed error types for encoding and decoding failures:
//
//	result, err := toon.MarshalToString(data)
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
//   - api.go - Public API entry points
//   - encode_*.go - Encoding logic for objects, arrays, and primitives
//   - decode_*.go - Decoding logic with structural and token parsers
//   - orderedmap.go - Ordered map implementation
//   - types.go, errors.go, options.go - Public type definitions
//
// The implementation achieves:
//
//   - 82.2% test coverage with 1,088 tests
//   - Average cyclomatic complexity of 5.39
//   - 340/340 official TOON specification fixtures passing
//   - File organization with encode_* and decode_* prefixes for clarity
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
