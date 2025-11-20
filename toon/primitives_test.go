package toon

import (
	"math"
	"testing"
)

// TestEncodePrimitive_Nil tests encoding nil values
func TestEncodePrimitive_Nil(t *testing.T) {
	result, err := encodePrimitive(nil, ",")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "null" {
		t.Errorf("expected 'null', got %q", result)
	}
}

// TestEncodePrimitive_Bool tests encoding boolean values
func TestEncodePrimitive_Bool(t *testing.T) {
	tests := []struct {
		name     string
		value    bool
		expected string
	}{
		{"true", true, "true"},
		{"false", false, "false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encodePrimitive(tt.value, ",")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestEncodePrimitive_SignedIntegers tests encoding signed integer types
func TestEncodePrimitive_SignedIntegers(t *testing.T) {
	tests := []struct {
		name     string
		value    Value
		expected string
	}{
		{"int zero", int(0), "0"},
		{"int positive", int(42), "42"},
		{"int negative", int(-42), "-42"},
		{"int max", int(math.MaxInt32), "2147483647"},
		{"int8 zero", int8(0), "0"},
		{"int8 positive", int8(127), "127"},
		{"int8 negative", int8(-128), "-128"},
		{"int16 zero", int16(0), "0"},
		{"int16 positive", int16(32767), "32767"},
		{"int16 negative", int16(-32768), "-32768"},
		{"int32 zero", int32(0), "0"},
		{"int32 positive", int32(2147483647), "2147483647"},
		{"int32 negative", int32(-2147483648), "-2147483648"},
		{"int64 zero", int64(0), "0"},
		{"int64 positive", int64(9223372036854775807), "9223372036854775807"},
		{"int64 negative", int64(-9223372036854775808), "-9223372036854775808"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encodePrimitive(tt.value, ",")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestEncodePrimitive_UnsignedIntegers tests encoding unsigned integer types
func TestEncodePrimitive_UnsignedIntegers(t *testing.T) {
	tests := []struct {
		name     string
		value    Value
		expected string
	}{
		{"uint zero", uint(0), "0"},
		{"uint positive", uint(42), "42"},
		{"uint max", uint(4294967295), "4294967295"},
		{"uint8 zero", uint8(0), "0"},
		{"uint8 positive", uint8(255), "255"},
		{"uint16 zero", uint16(0), "0"},
		{"uint16 positive", uint16(65535), "65535"},
		{"uint32 zero", uint32(0), "0"},
		{"uint32 positive", uint32(4294967295), "4294967295"},
		{"uint64 zero", uint64(0), "0"},
		{"uint64 small", uint64(42), "42"},
		{"uint64 large", uint64(18446744073709551615), "18446744073709551615"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encodePrimitive(tt.value, ",")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestEncodePrimitive_FloatingPoint tests encoding floating point types
func TestEncodePrimitive_FloatingPoint(t *testing.T) {
	tests := []struct {
		name     string
		value    Value
		expected string
	}{
		{"float32 zero", float32(0.0), "0"},
		{"float32 positive", float32(3.14), "3.140000104904175"},
		{"float32 negative", float32(-3.14), "-3.140000104904175"},
		{"float32 whole number", float32(42.0), "42"},
		{"float64 zero", float64(0.0), "0"},
		{"float64 positive", float64(3.14159), "3.14159"},
		{"float64 negative", float64(-3.14159), "-3.14159"},
		{"float64 whole number", float64(42.0), "42"},
		{"float64 small", float64(0.000001), "0.000001"},
		{"float64 large", float64(123456789.123456), "123456789.123456"},
		{"float64 very small", float64(1e-10), "0.0000000001"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encodePrimitive(tt.value, ",")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestEncodePrimitive_SpecialFloats tests encoding special floating point values
func TestEncodePrimitive_SpecialFloats(t *testing.T) {
	tests := []struct {
		name     string
		value    Value
		expected string
	}{
		{"float64 NaN", math.NaN(), "null"},
		{"float64 +Inf", math.Inf(1), "null"},
		{"float64 -Inf", math.Inf(-1), "null"},
		{"float64 negative zero", math.Copysign(0.0, -1), "0"},
		{"float32 NaN", float32(math.NaN()), "null"},
		{"float32 +Inf", float32(math.Inf(1)), "null"},
		{"float32 -Inf", float32(math.Inf(-1)), "null"},
		{"float32 negative zero", float32(math.Copysign(0.0, -1)), "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encodePrimitive(tt.value, ",")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestEncodePrimitive_Strings tests encoding string values
func TestEncodePrimitive_Strings(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		delimiter string
		expected  string
	}{
		{"empty string", "", ",", `""`},
		{"simple string", "hello", ",", "hello"},
		{"string with space", "hello world", ",", "hello world"},
		{"leading space", " hello", ",", `" hello"`},
		{"trailing space", "hello ", ",", `"hello "`},
		{"string with comma (active delimiter)", "hello,world", ",", `"hello,world"`},
		{"string with tab (inactive delimiter)", "hello\tworld", ",", `"hello\tworld"`},
		{"string with pipe (inactive delimiter)", "hello|world", ",", "hello|world"},
		{"string with tab as active delimiter", "hello\tworld", "\t", `"hello\tworld"`},
		{"string with pipe as active delimiter", "hello|world", "|", `"hello|world"`},
		{"string with quotes", `hello"world`, ",", `"hello\"world"`},
		{"string with newline", "hello\nworld", ",", `"hello\nworld"`},
		{"string with tab", "hello\tworld", ",", `"hello\tworld"`},
		{"string with carriage return", "hello\rworld", ",", `"hello\rworld"`},
		{"string with backslash", `hello\world`, ",", `"hello\\world"`},
		{"string with multiple escapes", "hello\n\t\"world\\", ",", `"hello\n\t\"world\\"`},
		{"string that looks like true", "true", ",", `"true"`},
		{"string that looks like false", "false", ",", `"false"`},
		{"string that looks like null", "null", ",", `"null"`},
		{"string that looks like integer", "42", ",", `"42"`},
		{"string that looks like float", "3.14", ",", `"3.14"`},
		{"string that looks like negative", "-42", ",", `"-42"`},
		{"string with leading zero", "007", ",", `"007"`},
		{"string with colon", "key:value", ",", `"key:value"`},
		{"string with brackets", "array[0]", ",", `"array[0]"`},
		{"string with braces", "{object}", ",", `"{object}"`},
		{"string with parens", "(value)", ",", `"(value)"`},
		{"string starting with hyphen", "-value", ",", `"-value"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encodePrimitive(tt.value, tt.delimiter)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestEncodePrimitive_UnsupportedTypes tests error handling for unsupported types
func TestEncodePrimitive_UnsupportedTypes(t *testing.T) {
	tests := []struct {
		name  string
		value Value
	}{
		{"map", map[string]interface{}{"key": "value"}},
		{"slice", []int{1, 2, 3}},
		{"struct", struct{ Name string }{Name: "test"}},
		{"array", [3]int{1, 2, 3}},
		{"pointer", new(int)},
		{"channel", make(chan int)},
		{"function", func() {}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := encodePrimitive(tt.value, ",")
			if err == nil {
				t.Fatal("expected error for unsupported type, got nil")
			}
			if encErr, ok := err.(*EncodeError); !ok {
				t.Errorf("expected *EncodeError, got %T", err)
			} else if encErr.Message != "unsupported primitive type" {
				t.Errorf("expected 'unsupported primitive type', got %q", encErr.Message)
			}
		})
	}
}

// TestEncodeFloat_EdgeCases tests edge cases in float encoding
func TestEncodeFloat_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		value    Value
		expected string
		wantErr  bool
	}{
		{"float64 max", math.MaxFloat64, "179769313486231570000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", false},
		{"float64 smallest positive", math.SmallestNonzeroFloat64, "0.000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005", false},
		{"float64 1.0", 1.0, "1", false},
		{"float64 -1.0", -1.0, "-1", false},
		{"float32 max", float32(math.MaxFloat32), "340282346638528860000000000000000000000", false},
		{"whole number within int64 range", float64(9223372036854775807), "9223372036854775807", false},
		{"whole number at max int64", float64(math.MaxInt64), "9223372036854775807", false},
		{"invalid type (string)", "not a float", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encodeFloat(tt.value)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestEncodeString_AllDelimiters tests string encoding with all valid delimiters
func TestEncodeString_AllDelimiters(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		delimiter string
		expected  string
	}{
		{"comma delimiter - contains comma", "a,b", ",", `"a,b"`},
		{"comma delimiter - no comma", "ab", ",", "ab"},
		{"tab delimiter - contains tab", "a\tb", "\t", `"a\tb"`},
		{"tab delimiter - no tab", "ab", "\t", "ab"},
		{"pipe delimiter - contains pipe", "a|b", "|", `"a|b"`},
		{"pipe delimiter - no pipe", "ab", "|", "ab"},
		{"value with inactive delimiters", "a,b|c", "\t", "a,b|c"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := encodeString(tt.value, tt.delimiter)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestNeedsQuoting_Comprehensive tests comprehensive quoting scenarios
func TestNeedsQuoting_Comprehensive(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		delimiter string
		expected  bool
	}{
		{"empty string", "", ",", true},
		{"simple word", "hello", ",", false},
		{"leading space", " hello", ",", true},
		{"trailing space", "hello ", ",", true},
		{"true literal", "true", ",", true},
		{"false literal", "false", ",", true},
		{"null literal", "null", ",", true},
		{"looks like integer", "42", ",", true},
		{"looks like float", "3.14", ",", true},
		{"looks like negative", "-42", ",", true},
		{"leading zero number", "007", ",", true},
		{"zero itself", "0", ",", true},
		{"contains colon", "key:value", ",", true},
		{"contains open bracket", "array[", ",", true},
		{"contains close bracket", "array]", ",", true},
		{"contains open brace", "{json", ",", true},
		{"contains close brace", "json}", ",", true},
		{"contains open paren", "(value", ",", true},
		{"contains close paren", "value)", ",", true},
		{"contains quote", `hello"world`, ",", true},
		{"contains backslash", `hello\world`, ",", true},
		{"contains newline", "hello\nworld", ",", true},
		{"contains tab", "hello\tworld", ",", true},
		{"contains carriage return", "hello\rworld", ",", true},
		{"contains active delimiter comma", "a,b", ",", true},
		{"contains inactive delimiter pipe", "a|b", ",", false},
		{"starts with hyphen", "-value", ",", true},
		{"hyphen in middle", "my-value", ",", false},
		{"safe alphanumeric", "hello123", ",", false},
		{"safe with underscore", "hello_world", ",", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := needsQuoting(tt.value, tt.delimiter)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestLooksLikeNumber tests number detection in strings
func TestLooksLikeNumber(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected bool
	}{
		{"integer", "42", true},
		{"negative integer", "-42", true},
		{"float", "3.14", true},
		{"negative float", "-3.14", true},
		{"scientific notation", "1e10", true},
		{"leading zero", "007", true},
		{"zero", "0", true},
		{"just decimal", ".5", true},
		{"not a number", "hello", false},
		{"number with text", "42hello", false},
		{"empty", "", false},
		{"just minus", "-", false},
		{"just dot", ".", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := looksLikeNumber(tt.value)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestEscapeString tests string escaping
func TestEscapeString(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{"no escaping", "hello", "hello"},
		{"backslash", `hello\world`, `hello\\world`},
		{"quote", `hello"world`, `hello\"world`},
		{"tab", "hello\tworld", `hello\tworld`},
		{"newline", "hello\nworld", `hello\nworld`},
		{"carriage return", "hello\rworld", `hello\rworld`},
		{"multiple escapes", `a\b"c` + "\nd\te", `a\\b\"c\nd\te`},
		{"backslash then quote", `\"`, `\\\"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := escapeString(tt.value)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestEncodeKey tests key encoding with strict rules
func TestEncodeKey(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected string
	}{
		{"simple key", "key", "key"},
		{"uppercase start", "Key", "Key"},
		{"underscore start", "_key", "_key"},
		{"with dot", "key.subkey", "key.subkey"},
		{"with numbers", "key123", "key123"},
		{"empty key", "", `""`},
		{"number start", "123key", `"123key"`},
		{"with space", "my key", `"my key"`},
		{"with hyphen", "my-key", `"my-key"`},
		{"with special chars", "key:value", `"key:value"`},
		{"lowercase start", "key", "key"},
		{"mixed case", "myKey", "myKey"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := encodeKey(tt.key)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestSafeKey tests key safety validation
func TestSafeKey(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected bool
	}{
		{"simple lowercase", "key", true},
		{"simple uppercase", "Key", true},
		{"underscore start", "_key", true},
		{"with dot", "key.subkey", true},
		{"with numbers", "key123", true},
		{"mixed case with dot and number", "myKey.sub_1", true},
		{"empty", "", false},
		{"number start", "1key", false},
		{"hyphen start", "-key", false},
		{"with space", "my key", false},
		{"with hyphen", "my-key", false},
		{"with colon", "key:value", false},
		{"with special char", "key@value", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := safeKey(tt.key)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
