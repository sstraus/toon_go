package toon

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
)

// Version is the current version of the TOON library.
const Version = "1.1.0"

// Marshal encodes a Go value to TOON format and writes it to w.
//
// The value v can be any JSON-compatible type: nil, bool, int, int64, float64,
// string, []interface{}, or map[string]interface{}.
//
// Options are configured using functional options.
//
// Example:
//
//	data := map[string]interface{}{
//		"tags": []interface{}{"go", "toon"},
//	}
//	var buf bytes.Buffer
//	err := toon.Marshal(data, &buf)
//	// buf contains: tags[2]: go,toon
//
// With options:
//
//	err := toon.Marshal(data, &buf, WithIndent(4), WithDelimiter("\t"))
func Marshal(v interface{}, w io.Writer, opts ...EncodeOption) error {
	// Apply functional options
	encOpts := applyEncodeOptions(opts...)

	// Validate options
	if err := validateEncodeOptions(encOpts); err != nil {
		return err
	}

	// Normalize the value
	normalized := normalize(v)

	// Encode
	result, err := encode(normalized, encOpts)
	if err != nil {
		return err
	}

	// Write to io.Writer
	_, err = w.Write([]byte(result))
	return err
}

// MarshalToString encodes a Go value to TOON format and returns it as a string.
//
// This is a convenience function that wraps Marshal.
//
// Example:
//
//	data := map[string]interface{}{
//		"name": "Alice",
//		"age":  30,
//	}
//	result, err := toon.MarshalToString(data)
//	// result: "age: 30\nname: Alice"
//
// With options:
//
//	result, err := toon.MarshalToString(data, WithIndent(4), WithDelimiter("\t"))
func MarshalToString(v interface{}, opts ...EncodeOption) (string, error) {
	var buf bytes.Buffer
	err := Marshal(v, &buf, opts...)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Unmarshal decodes TOON format data from r into a Go value.
//
// The value pointed to by v will be populated with the decoded data.
// v must be a pointer to a value of the appropriate type.
//
// Options are configured using functional options.
//
// Example:
//
//	input := strings.NewReader("name: Alice\nage: 30")
//	var result map[string]interface{}
//	err := toon.Unmarshal(input, &result)
//	// result: map[string]interface{}{"name": "Alice", "age": 30}
//
// With options:
//
//	err := toon.Unmarshal(input, &result, WithStrictDecoding(false))
func Unmarshal(r io.Reader, v interface{}, opts ...DecodeOption) error {
	// Read from io.Reader
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	// Apply functional options
	decOpts := applyDecodeOptions(opts...)

	// Validate options
	if err := validateDecodeOptions(decOpts); err != nil {
		return err
	}

	// Decode
	result, err := decode(string(data), decOpts)
	if err != nil {
		return err
	}

	// Assign result to v
	return assignResult(result, v)
}

// UnmarshalFromString decodes TOON format data from a string into a Go value.
//
// This is a convenience function that wraps Unmarshal.
//
// Example:
//
//	input := "name: Alice\nage: 30"
//	var result map[string]interface{}
//	err := toon.UnmarshalFromString(input, &result)
//	// result: map[string]interface{}{"name": "Alice", "age": 30}
//
// With options:
//
//	err := toon.UnmarshalFromString(input, &result, WithStrictDecoding(false))
func UnmarshalFromString(s string, v interface{}, opts ...DecodeOption) error {
	return Unmarshal(strings.NewReader(s), v, opts...)
}

// applyEncodeOptions applies functional options to create EncodeOptions.
func applyEncodeOptions(opts ...EncodeOption) *EncodeOptions {
	// Start with defaults
	encOpts := getEncodeOptions(nil)

	// Apply functional options
	for _, opt := range opts {
		opt(encOpts)
	}

	return encOpts
}

// applyDecodeOptions applies functional options to create DecodeOptions.
func applyDecodeOptions(opts ...DecodeOption) *DecodeOptions {
	// Start with defaults
	decOpts := getDecodeOptions(nil)

	// Apply functional options
	for _, opt := range opts {
		opt(decOpts)
	}

	return decOpts
}

// assignResult assigns the decoded result to the target variable.
func assignResult(result Value, v interface{}) error {
	switch target := v.(type) {
	case *map[string]interface{}:
		if m, ok := result.(map[string]Value); ok {
			// Convert map[string]Value to map[string]interface{}
			converted := make(map[string]interface{}, len(m))
			for k, val := range m {
				converted[k] = val
			}
			*target = converted
		} else {
			return &DecodeError{Message: "cannot assign non-map to map target"}
		}
	case *[]interface{}:
		if arr, ok := result.([]Value); ok {
			// Convert []Value to []interface{}
			converted := make([]interface{}, len(arr))
			for i, val := range arr {
				converted[i] = val
			}
			*target = converted
		} else {
			return &DecodeError{Message: "cannot assign non-array to array target"}
		}
	case *interface{}:
		*target = result
	default:
		return &DecodeError{Message: "unsupported target type"}
	}
	return nil
}
