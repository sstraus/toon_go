# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-11-20
### Added
- Initial stable release of the TOON encoder and decoder for Go.
- Support for primitive values, objects, and all TOON array formats (inline, tabular, list).
- Custom encoding and decoding options for indentation, delimiters, length markers, and flattening.
- Comprehensive fixture-based test suite with coverage tooling.
- Public `Version` constant for programmatic access to the library version.

### Tooling
- Continuous integration workflow running `go test` and collecting coverage.
