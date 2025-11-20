# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2025-11-20
### Changed
- **BREAKING**: `Marshal` now uses variadic functional options instead of struct pointer
  - Old: `Marshal(v, w, &EncodeOptions{Indent: 4})`
  - New: `Marshal(v, w, WithIndent(4))`
  - Backward compatibility maintained in tests via helper functions
- **BREAKING**: `Unmarshal` now uses variadic functional options instead of struct pointer
  - Old: `Unmarshal(r, v, &DecodeOptions{Strict: false})`
  - New: `Unmarshal(r, v, WithStrictDecoding(false))`
- File organization: Renamed implementation files with `encode_*` and `decode_*` prefixes for better organization
- Updated package documentation with comprehensive examples of new API

### Added
- `MarshalToString()` convenience function for direct string encoding
- `UnmarshalFromString()` convenience function for direct string decoding
- Functional option functions for cleaner API:
  - **Encoding**: `WithIndent()`, `WithDelimiter()`, `WithLengthMarker()`, `WithFlattenPaths()`, `WithFlattenDepth()`, `WithStrict()`
  - **Decoding**: `WithKeyMode()`, `WithStrictDecoding()`, `WithIndentSize()`, `WithExpandPaths()`
- Enhanced `doc.go` with detailed usage examples and functional options guide
- Created dedicated `api.go` file separating public API from implementation

### Documentation
- Updated README.md with functional options examples
- Updated all code examples to use new API style
- Added comprehensive functional options reference
- Updated examples/basic/main.go to demonstrate new API

## [1.0.0] - 2025-11-20
### Added
- Initial stable release of the TOON encoder and decoder for Go.
- Support for primitive values, objects, and all TOON array formats (inline, tabular, list).
- Custom encoding and decoding options for indentation, delimiters, length markers, and flattening.
- Comprehensive fixture-based test suite with coverage tooling.
- Public `Version` constant for programmatic access to the library version.

### Tooling
- Continuous integration workflow running `go test` and collecting coverage.
