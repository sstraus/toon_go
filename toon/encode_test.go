package toon

import (
	"strings"
	"testing"
)

func TestFlattenObjectCollisionStrict(t *testing.T) {
	obj := map[string]Value{
		"a":   map[string]Value{"b": int64(1)},
		"a.b": int64(2),
	}
	opts := &EncodeOptions{Strict: true}

	_, err := flattenObject(obj, "", 0, opts)
	if err == nil {
		t.Fatalf("expected collision error in strict mode")
	}
	if _, ok := err.(*EncodeError); !ok {
		t.Fatalf("expected *EncodeError, got %T", err)
	}
	if !strings.Contains(err.Error(), "key collision") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestEscapeAndEncodeStringRoundTrip(t *testing.T) {
	orig := "hello\tback\\slash\"newline\nend"
	esc := escapeString(orig)
	un := unescapeString(esc)
	if un != orig {
		t.Fatalf("unescape mismatch:\nwant: %q\ngot:  %q", orig, un)
	}

	enc := encodeString(orig, ",")
	if !strings.HasPrefix(enc, "\"") || !strings.HasSuffix(enc, "\"") {
		t.Fatalf("expected quoted string, got %q", enc)
	}
	inner := enc[1 : len(enc)-1]
	if unescapeString(inner) != orig {
		t.Fatalf("encodeString produced incorrect escaping")
	}
}
