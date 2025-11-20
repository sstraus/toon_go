package toon

import (
	"bytes"
)

// marshalToBytes is a test helper that wraps the new Marshal API for backward compatibility.
// This allows existing tests to continue using the old []byte-based API.
func marshalToBytes(v interface{}, opts *EncodeOptions) ([]byte, error) {
	var buf bytes.Buffer
	err := Marshal(v, &buf, opts)
	return buf.Bytes(), err
}

// unmarshalFromBytes is a test helper that wraps the new Unmarshal API for backward compatibility.
// This allows existing tests to continue using the old []byte-based API.
func unmarshalFromBytes(data []byte, v interface{}, opts *DecodeOptions) error {
	return Unmarshal(bytes.NewReader(data), v, opts)
}
