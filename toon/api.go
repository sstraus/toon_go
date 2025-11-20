package toon

import (
	"io"
	"io/ioutil"
)

// Version is the current version of the TOON library.
const Version = "1.1.0"

// Marshal encodes a Go value to TOON format and writes it to w.
//
// The value v can be any JSON-compatible type: nil, bool, int, int64, float64,
// string, []interface{}, or map[string]interface{}.
//
// Options can be nil to use defaults, or specify custom encoding behavior.
//
// Example:
//
//	data := map[string]interface{}{
//		"tags": []interface{}{"go", "toon"},
//	}
//	var buf bytes.Buffer
//	err := toon.Marshal(data, &buf, nil)
//	// buf contains: tags[2]: go,toon
func Marshal(v interface{}, w io.Writer, opts *EncodeOptions) error {
	// Validate options
	opts = getEncodeOptions(opts)
	if err := validateEncodeOptions(opts); err != nil {
		return err
	}

	// Normalize the value
	normalized := normalize(v)

	// Encode
	result, err := encode(normalized, opts)
	if err != nil {
		return err
	}

	// Write to io.Writer
	_, err = w.Write([]byte(result))
	return err
}

// Unmarshal decodes TOON format data from r into a Go value.
//
// The value pointed to by v will be populated with the decoded data.
// v must be a pointer to a value of the appropriate type.
//
// Options can be nil to use defaults, or specify custom decoding behavior.
//
// Example:
//
//	input := strings.NewReader("name: Alice\nage: 30")
//	var result map[string]interface{}
//	err := toon.Unmarshal(input, &result, nil)
//	// result: map[string]interface{}{"name": "Alice", "age": 30}
func Unmarshal(r io.Reader, v interface{}, opts *DecodeOptions) error {
	// Read from io.Reader
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	// Validate options
	opts = getDecodeOptions(opts)
	if err := validateDecodeOptions(opts); err != nil {
		return err
	}

	// Decode
	result, err := decode(string(data), opts)
	if err != nil {
		return err
	}

	// Assign result to v
	return assignResult(result, v)
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
