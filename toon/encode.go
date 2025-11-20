package toon

import "fmt"

// encode encodes a normalized value to TOON format string.
func encode(v Value, opts *EncodeOptions) (string, error) {
	w := newWriter(opts.Indent)

	if err := encodeValue(w, "", v, 0, opts); err != nil {
		return "", err
	}

	return w.String(), nil
}

// encodeValue encodes a value with an optional key.
func encodeValue(w *writer, key string, v Value, depth int, opts *EncodeOptions) error {
	if v == nil {
		return encodeValueNil(w, key, depth)
	}

	if isPrimitive(v) {
		return encodeValuePrimitive(w, key, v, depth, opts)
	}

	if isMap(v) {
		return encodeValueMap(w, key, v, depth, opts)
	}

	if isList(v) {
		return encodeArray(w, key, v, depth, opts)
	}

	return &EncodeError{
		Message: "unsupported type",
		Value:   v,
	}
}

// encodeValueNil encodes a nil value.
func encodeValueNil(w *writer, key string, depth int) error {
	if key != "" {
		w.push(key+colon+space+nullLiteral, depth)
	} else {
		w.push(nullLiteral, depth)
	}
	return nil
}

// encodeValuePrimitive encodes a primitive value.
func encodeValuePrimitive(w *writer, key string, v Value, depth int, opts *EncodeOptions) error {
	encoded, err := encodePrimitive(v, opts.Delimiter)
	if err != nil {
		return err
	}

	if key != "" {
		w.push(key+colon+space+encoded, depth)
	} else {
		if depth == 0 && w.Len() == 0 {
			// Root primitive without indentation
			w.pushRaw(encoded)
		} else {
			w.push(encoded, depth)
		}
	}
	return nil
}

// encodeValueMap encodes a map value with optional flattening.
func encodeValueMap(w *writer, key string, v Value, depth int, opts *EncodeOptions) error {
	mapValue, err := convertToMapValue(v)
	if err != nil {
		return err
	}

	// Apply flattening if enabled (only at root level, depth==0)
	if opts.FlattenPaths && depth == 0 {
		return encodeFlattenedMap(w, key, mapValue, depth, opts)
	}

	return encodeObject(w, key, v, depth, opts)
}

// convertToMapValue converts various map types to map[string]Value.
func convertToMapValue(v Value) (map[string]Value, error) {
	if orderedMap, ok := v.(OrderedMap); ok {
		mapValue := make(map[string]Value)
		for k, val := range orderedMap.Values() {
			mapValue[k] = val
		}
		return mapValue, nil
	}

	if orderedMapPtr, ok := v.(*OrderedMap); ok {
		mapValue := make(map[string]Value)
		for k, val := range orderedMapPtr.Values() {
			mapValue[k] = val
		}
		return mapValue, nil
	}

	if m, ok := v.(map[string]Value); ok {
		return m, nil
	}

	return nil, &EncodeError{
		Message: "unsupported map type",
		Value:   v,
	}
}

// encodeFlattenedMap flattens and encodes a map.
func encodeFlattenedMap(w *writer, key string, mapValue map[string]Value, depth int, opts *EncodeOptions) error {
	flattened, err := flattenObject(mapValue, "", 0, opts)
	if err != nil {
		return fmt.Errorf("flatten failed: %w", err)
	}

	// Create new options with flattening disabled for nested encoding
	nestedOpts := *opts
	nestedOpts.FlattenPaths = false
	return encodeObject(w, key, flattened, depth, &nestedOpts)
}
