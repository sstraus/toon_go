# TOON Go Examples

This directory contains comprehensive examples demonstrating the TOON Go library features using the v1.1.0 API with functional options and string convenience functions.

## Running Examples

Each example is a standalone Go program. To run any example:

```bash
cd examples/<example-name>
go run main.go
```

Or from the repository root:

```bash
go run ./examples/<example-name>
```

## Available Examples

### 1. Basic (`basic/`)

**What it demonstrates:**
- Simple primitives encoding
- Objects and nested objects
- All three array formats (inline, tabular, list)
- Custom encoding options
- Complex nested structures
- Token efficiency vs JSON

**Best for:** Getting started with TOON

```bash
go run ./examples/basic
```

### 2. Advanced Encoding (`advanced-encoding/`)

**What it demonstrates:**
- Inline arrays with primitives
- Tabular arrays with uniform objects
- List arrays with mixed/nested content
- Custom delimiters (tabs, pipes)
- Custom indentation (4 spaces)
- Length markers with prefixes
- Complex microservices architecture
- Combining multiple functional options

**Best for:** Learning all encoding capabilities and options

```bash
go run ./examples/advanced-encoding
```

**Key features shown:**
- `WithDelimiter("\t")` - Tab-separated values
- `WithIndent(4)` - Custom indentation
- `WithLengthMarker("#")` - Array length prefixes
- Multiple options combined

### 3. Decoding (`decoding/`)

**What it demonstrates:**
- Decoding simple objects
- Decoding nested structures
- Decoding all array formats
- Using `UnmarshalFromString()` convenience function
- Decoding with `WithStrictDecoding(false)`
- Round-trip encoding/decoding
- Handling different data types
- Decoding to `interface{}`

**Best for:** Understanding TOON parsing and decoding options

```bash
go run ./examples/decoding
```

**Key features shown:**
- `toon.UnmarshalFromString()` - Direct string decoding
- `toon.WithStrictDecoding(false)` - Flexible parsing
- Type preservation (int, float, bool, null, string)

### 4. OrderedMap (`orderedmap/`)

**What it demonstrates:**
- Creating OrderedMaps
- Preserving key insertion order
- Nested OrderedMaps
- OrderedMap methods (Get, Set, Delete, Keys, Len)
- Comparison with regular maps
- Use cases: HTTP headers, build pipelines, form fields, documentation structure

**Best for:** Learning when and how to preserve key ordering

```bash
go run ./examples/orderedmap
```

**Key concepts:**
- Regular maps → alphabetically sorted keys
- OrderedMaps → insertion order preserved
- Perfect for configuration files, APIs, and sequential data

### 5. Configuration Files (`configuration/`)

**What it demonstrates:**
- Reading application configuration from TOON
- Generating configuration files
- Environment-specific configurations (dev, staging, prod)
- Feature flags configuration
- Microservices configuration
- Real-world config structure (database, cache, logging, security)

**Best for:** Using TOON as a config file format

```bash
go run ./examples/configuration
```

**Use cases shown:**
- Web application configuration
- Payment gateway settings
- Microservices deployment configs
- Feature flag management

### 6. API Responses (`api-response/`)

**What it demonstrates:**
- User profile API responses
- Paginated list responses
- Error responses with validation details
- Search results
- Analytics dashboards
- Webhook payloads
- GraphQL-style nested responses

**Best for:** Using TOON for LLM-friendly API communication

```bash
go run ./examples/api-response
```

**Benefits highlighted:**
- 30-60% token savings vs JSON
- More readable for LLMs
- Perfect for AI agents and LLM-powered applications

## Example Comparison

| Example | Focus | Complexity | API Features Shown |
|---------|-------|------------|-------------------|
| basic | Getting started | ⭐ Simple | MarshalToString, basic usage |
| advanced-encoding | All encoding options | ⭐⭐ Moderate | All functional options, multiple formats |
| decoding | Parsing & decoding | ⭐⭐ Moderate | UnmarshalFromString, strict mode |
| orderedmap | Key ordering | ⭐⭐ Moderate | OrderedMap methods, use cases |
| configuration | Config files | ⭐⭐⭐ Advanced | Real-world structures, environments |
| api-response | API communication | ⭐⭐⭐ Advanced | Complex nested data, LLM use cases |

## New in v1.1.0

All examples use the new v1.1.0 API featuring:

### String Convenience Functions
```go
// Instead of bytes.Buffer
result, err := toon.MarshalToString(data)
err := toon.UnmarshalFromString(input, &result)
```

### Functional Options Pattern
```go
// Instead of struct pointers
result, _ := toon.MarshalToString(data,
    toon.WithIndent(4),
    toon.WithDelimiter("\t"),
    toon.WithLengthMarker("#"),
)

err := toon.UnmarshalFromString(input, &result,
    toon.WithStrictDecoding(false),
    toon.WithExpandPaths("safe"),
)
```

## Common Patterns

### Encoding with Options
```go
data := map[string]interface{}{
    "values": []interface{}{1, 2, 3},
}

result, _ := toon.MarshalToString(data,
    toon.WithIndent(4),
    toon.WithDelimiter("|"),
)
```

### Decoding with Flexibility
```go
input := `name: Test
value: 123`

var config map[string]interface{}
err := toon.UnmarshalFromString(input, &config,
    toon.WithStrictDecoding(false),
)
```

### Using OrderedMap
```go
om := toon.NewOrderedMap()
om.Set("first", 1)
om.Set("second", 2)
om.Set("third", 3)

result, _ := toon.MarshalToString(om)
// Keys appear in insertion order
```

## Learning Path

1. Start with **basic** to understand core concepts
2. Explore **advanced-encoding** for all encoding features
3. Learn **decoding** for parsing capabilities
4. Use **orderedmap** when key order matters
5. See **configuration** for config file use cases
6. Check **api-response** for LLM/API applications

## Additional Resources

- [Main README](../README.md) - Library overview and installation
- [API Documentation](https://pkg.go.dev/github.com/sstraus/toon_go/toon) - Complete API reference
- [TOON Specification](https://github.com/toon-format/spec) - Format specification
- [CHANGELOG](../CHANGELOG.md) - Version history and migration guides

## Token Efficiency

TOON achieves 30-60% token reduction compared to JSON:

**JSON Example (123 tokens):**
```json
{
  "users": [
    {"name": "Alice", "age": 30},
    {"name": "Bob", "age": 25}
  ]
}
```

**TOON Example (52 tokens - 58% reduction):**
```
users[2]{age,name}:
  30,Alice
  25,Bob
```

This makes TOON ideal for:
- LLM-based applications
- AI agent communication
- Token-constrained APIs
- Configuration files for AI systems

## Contributing Examples

Have a great example? Consider contributing!

1. Create a new directory under `examples/`
2. Add a clear, focused `main.go`
3. Update this README with your example
4. Submit a pull request

Examples should:
- Demonstrate a specific feature or use case
- Use the v1.1.0 API (functional options, string functions)
- Include clear output and comments
- Be self-contained and runnable
