package toon

import "testing"

func TestExpandDottedKeyCreatesNested(t *testing.T) {
	sp := newStructuralParser("", &DecodeOptions{Strict: true})
	target := map[string]Value{}

	if err := sp.expandDottedKey("a.b.c", int64(1), target); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	a, ok := target["a"].(map[string]Value)
	if !ok {
		t.Fatalf("expected 'a' to be a map, got: %T", target["a"])
	}
	b, ok := a["b"].(map[string]Value)
	if !ok {
		t.Fatalf("expected 'a.b' to be a map, got: %T", a["b"])
	}
	if b["c"] != int64(1) {
		t.Fatalf("unexpected value at a.b.c: %v", b["c"])
	}
}

func TestExpandDottedKeyConflictStrict(t *testing.T) {
	sp := newStructuralParser("", &DecodeOptions{Strict: true})
	target := map[string]Value{"a": int64(5)}

	err := sp.expandDottedKey("a.b", int64(1), target)
	if err == nil {
		t.Fatalf("expected error due to path conflict in strict mode")
	}
	if _, ok := err.(*DecodeError); !ok {
		t.Fatalf("expected DecodeError, got: %T", err)
	}
}

func TestExpandDottedKeyNonStrictOverwrites(t *testing.T) {
	sp := newStructuralParser("", &DecodeOptions{Strict: false})
	target := map[string]Value{"a": int64(5)}

	if err := sp.expandDottedKey("a.b", int64(1), target); err != nil {
		t.Fatalf("unexpected error in non-strict mode: %v", err)
	}

	a, ok := target["a"].(map[string]Value)
	if !ok {
		t.Fatalf("expected 'a' to be a map after non-strict overwrite, got: %T", target["a"])
	}
	if a["b"] != int64(1) {
		t.Fatalf("unexpected value at a.b after overwrite: %v", a["b"])
	}
}
