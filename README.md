# TOON Go Library

A Go implementation of TOON (Token-Oriented Object Notation) format encoder and decoder, compliant with TOON Specification v2.0.

✅ **Encoder & Decoder Complete** ✅

Currently implemented:
- ✅ Complete encoder with all three array formats (inline, tabular, list)
- ✅ Complete decoder with structural parsing
- ✅ Primitive types (nil, bool, int, float, string)
- ✅ Objects (maps) with nested support
- ✅ Arrays with format auto-detection
- ✅ Custom delimiters and length markers
- ✅ String quoting and escaping
- ✅ Round-trip encoding/decoding
- ✅ Official specification fixtures (340/340 passing)
- ✅ Comprehensive test suite (196 decode + 144 encode)
- ✅ 65.4% code coverage

## About TOON

TOON is a compact data format optimized for LLM token efficiency, achieving **30-60% token reduction** compared to JSON while maintaining readability.

### Key Features

- **Token Efficient**: 30-60% fewer tokens than JSON
- **Human Readable**: Indentation-based structure like YAML
- **Three Array Formats**: 
  - Inline: `tags[2]: a,b` (for primitives)
  - Tabular: `users[2]{name,age}: Alice,30 / Bob,25` (for uniform objects)
  - List: `items[2]: - item1 / - item2` (for mixed/nested)

## Installation

```bash
go get github.com/sstraus/toon_go/toon
```

## Usage

### Encoding & Decoding

```go
package main

import (
    "fmt"
    "github.com/sstraus/toon_go/toon"
)

func main() {
    // Simple object
    data := map[string]interface{}{
        "name": "Alice",
        "age":  30,
    }
    
    encoded, err := toon.Marshal(data, nil)
    if err != nil {
        panic(err)
    }
    
    fmt.Println(string(encoded))
    // Output:
    // age: 30
    // name: Alice
}
```

### Arrays

```go
// Inline array (primitives)
data := map[string]interface{}{
    "tags": []interface{}{"go", "toon", "llm"},
}
encoded, _ := toon.Marshal(data, nil)
fmt.Println(string(encoded))
// Output: tags[3]: go,toon,llm

// Tabular array (uniform objects)
users := map[string]interface{}{
    "users": []interface{}{
        map[string]interface{}{"name": "Alice", "age": 30},
        map[string]interface{}{"name": "Bob", "age": 25},
    },
}
encoded, _ := toon.Marshal(users, nil)
fmt.Println(string(encoded))
// Output:
// users[2]{age,name}:
//   30,Alice
//   25,Bob
```

### Custom Options

```go
opts := &toon.EncodeOptions{
    Indent:       4,           // 4 spaces instead of 2
    Delimiter:    "\t",        // Tab delimiter
    LengthMarker: "#",         // Prefix length with #
}

data := map[string]interface{}{
    "values": []interface{}{1, 2, 3},
}

encoded, _ := toon.Marshal(data, opts)
fmt.Println(string(encoded))
// Output: values[#3	]: 1	2	3
```

## Project Structure

```
toon/
├── toon.go              # Public API (Marshal/Unmarshal)
├── encode.go            # Encoder implementation
├── decode.go            # Decoder implementation
├── structural_parser.go # Structural/indentation-based parser
├── parser.go            # Token parser
├── types.go             # Type definitions
├── constants.go         # Format constants
├── errors.go            # Error types
├── primitives.go        # Primitive encoding/decoding
├── strings.go           # String utilities
├── arrays.go            # Array format logic
├── objects.go           # Object encoding
├── orderedmap.go        # Ordered map implementation
├── writer.go            # Output writer
├── options.go           # Option types
├── utils.go             # Utilities
├── toon_test.go         # Encoder tests
└── decode_test.go       # Decoder tests
```

## Testing

The implementation includes comprehensive test coverage:

```bash
# Run all tests
go test

# Run with verbose output
go test -v

# Check coverage
go test -cover
```

**Test Status:**
- ✅ 340/340 specification fixtures passing
- ✅ 196 decode fixtures
- ✅ 144 encode fixtures
- ✅ 65.4% code coverage

## Specification

This implementation follows [TOON Specification v2.0](https://github.com/toon-format/spec).

## License

MIT

## Author

Created by [Stefano Straus](https://github.com/sstraus) in 2025.