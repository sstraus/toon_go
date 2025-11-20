package toon

import (
	"bytes"
)

// marshalToBytes is a test helper that wraps the new Marshal API for backward compatibility.
// This allows existing tests to continue using the old []byte-based API with *EncodeOptions.
func marshalToBytes(v interface{}, opts *EncodeOptions) ([]byte, error) {
	var buf bytes.Buffer
	// Convert *EncodeOptions to functional options
	funcOpts := encodeOptionsToFunctional(opts)
	err := Marshal(v, &buf, funcOpts...)
	return buf.Bytes(), err
}

// unmarshalFromBytes is a test helper that wraps the new Unmarshal API for backward compatibility.
// This allows existing tests to continue using the old []byte-based API with *DecodeOptions.
func unmarshalFromBytes(data []byte, v interface{}, opts *DecodeOptions) error {
	// Convert *DecodeOptions to functional options
	funcOpts := decodeOptionsToFunctional(opts)
	return Unmarshal(bytes.NewReader(data), v, funcOpts...)
}

// encodeOptionsToFunctional converts *EncodeOptions to functional options.
func encodeOptionsToFunctional(opts *EncodeOptions) []EncodeOption {
	// Normalize options through getEncodeOptions to fill in defaults
	normalizedOpts := getEncodeOptions(opts)

	// Create a single functional option that applies all fields from the normalized struct
	return []EncodeOption{
		func(o *EncodeOptions) {
			*o = *normalizedOpts
		},
	}
}

// decodeOptionsToFunctional converts *DecodeOptions to functional options.
func decodeOptionsToFunctional(opts *DecodeOptions) []DecodeOption {
	// Normalize options through getDecodeOptions to fill in defaults
	normalizedOpts := getDecodeOptions(opts)

	// Create a single functional option that applies all fields from the normalized struct
	return []DecodeOption{
		func(o *DecodeOptions) {
			*o = *normalizedOpts
		},
	}
}
