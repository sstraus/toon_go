package toon

import (
	"errors"
	"strings"
	"testing"
)

func TestEncodeError_ErrorAndUnwrap(t *testing.T) {
	cause := errors.New("root-cause")
	e := &EncodeError{Message: "encode failed", Value: "x", Cause: cause}

	msg := e.Error()
	if !strings.Contains(msg, "encode failed") || !strings.Contains(msg, "x") {
		t.Fatalf("unexpected Error() output: %s", msg)
	}
	if e.Unwrap() != cause {
		t.Fatalf("unexpected Unwrap(): %v", e.Unwrap())
	}
	if !errors.Is(e, cause) {
		t.Fatalf("errors.Is failed")
	}
}

func TestDecodeError_Formatting(t *testing.T) {
	e := &DecodeError{
		Message: "parse error",
		Line:    2,
		Column:  5,
		Token:   "tok",
		Context: "line content",
	}

	msg := e.Error()
	if !strings.Contains(msg, "parse error at line 2, column 5") {
		t.Fatalf("missing line/column in message: %s", msg)
	}
	if !strings.Contains(msg, "(token: 'tok')") {
		t.Fatalf("missing token in message: %s", msg)
	}
	if !strings.Contains(msg, "Context:\nline content") {
		t.Fatalf("missing context: %s", msg)
	}
}

func TestDecodeError_Error_AllFormats(t *testing.T) {
	tests := []struct {
		name     string
		err      *DecodeError
		contains []string
		notContains []string
	}{
		{
			name: "message only",
			err: &DecodeError{
				Message: "simple error",
			},
			contains: []string{"simple error"},
			notContains: []string{"at line", "at column", "token:", "Context:"},
		},
		{
			name: "with line and column",
			err: &DecodeError{
				Message: "position error",
				Line:    10,
				Column:  25,
			},
			contains: []string{"position error at line 10, column 25"},
		},
		{
			name: "with line only",
			err: &DecodeError{
				Message: "line error",
				Line:    5,
			},
			contains: []string{"line error at line 5"},
			notContains: []string{"column"},
		},
		{
			name: "with column only",
			err: &DecodeError{
				Message: "column error",
				Column:  15,
			},
			contains: []string{"column error at column 15"},
			notContains: []string{"at line"},
		},
		{
			name: "with token",
			err: &DecodeError{
				Message: "token error",
				Token:   "invalid_token",
			},
			contains: []string{"token error", "(token: 'invalid_token')"},
		},
		{
			name: "with context",
			err: &DecodeError{
				Message: "context error",
				Context: "surrounding\ncode\nlines",
			},
			contains: []string{"context error", "Context:", "surrounding\ncode\nlines"},
		},
		{
			name: "all fields",
			err: &DecodeError{
				Message: "complete error",
				Line:    20,
				Column:  30,
				Token:   "bad_token",
				Context: "error context here",
			},
			contains: []string{
				"complete error at line 20, column 30",
				"(token: 'bad_token')",
				"Context:",
				"error context here",
			},
		},
		{
			name: "empty token",
			err: &DecodeError{
				Message: "no token",
				Token:   "",
			},
			contains: []string{"no token"},
			notContains: []string{"token:"},
		},
		{
			name: "empty context",
			err: &DecodeError{
				Message: "no context",
				Context: "",
			},
			contains: []string{"no context"},
			notContains: []string{"Context:"},
		},
		{
			name: "zero line and column",
			err: &DecodeError{
				Message: "no position",
				Line:    0,
				Column:  0,
			},
			contains: []string{"no position"},
			notContains: []string{"at line", "at column"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tt.err.Error()
			for _, s := range tt.contains {
				if !strings.Contains(msg, s) {
					t.Errorf("expected message to contain %q, got: %s", s, msg)
				}
			}
			for _, s := range tt.notContains {
				if strings.Contains(msg, s) {
					t.Errorf("expected message not to contain %q, got: %s", s, msg)
				}
			}
		})
	}
}

func TestDecodeError_Unwrap(t *testing.T) {
	tests := []struct {
		name      string
		err       *DecodeError
		wantCause error
	}{
		{
			name: "with cause",
			err: &DecodeError{
				Message: "wrapped error",
				Cause:   errors.New("underlying cause"),
			},
			wantCause: errors.New("underlying cause"),
		},
		{
			name: "without cause",
			err: &DecodeError{
				Message: "no cause",
				Cause:   nil,
			},
			wantCause: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Unwrap()
			if tt.wantCause == nil {
				if got != nil {
					t.Errorf("Unwrap() = %v, want nil", got)
				}
			} else {
				if got == nil {
					t.Errorf("Unwrap() = nil, want non-nil")
				} else if got.Error() != tt.wantCause.Error() {
					t.Errorf("Unwrap() = %v, want %v", got, tt.wantCause)
				}
			}
		})
	}
}

func TestDecodeError_ErrorsIs(t *testing.T) {
	rootCause := errors.New("root cause")
	decErr := &DecodeError{
		Message: "decode failed",
		Cause:   rootCause,
	}

	if !errors.Is(decErr, rootCause) {
		t.Error("errors.Is should find the root cause")
	}

	otherErr := errors.New("other error")
	if errors.Is(decErr, otherErr) {
		t.Error("errors.Is should not match unrelated error")
	}
}

func TestEncodeError_Error_AllFormats(t *testing.T) {
	tests := []struct {
		name     string
		err      *EncodeError
		contains []string
		notContains []string
	}{
		{
			name: "message only",
			err: &EncodeError{
				Message: "encode failed",
			},
			contains: []string{"encode failed"},
		},
		{
			name: "with value",
			err: &EncodeError{
				Message: "invalid value",
				Value:   "test_value",
			},
			contains: []string{"invalid value:", "test_value"},
		},
		{
			name: "with cause",
			err: &EncodeError{
				Message: "encoding error",
				Cause:   errors.New("underlying issue"),
			},
			contains: []string{"encoding error:", "underlying issue"},
		},
		{
			name: "value takes precedence over cause",
			err: &EncodeError{
				Message: "both present",
				Value:   "some_value",
				Cause:   errors.New("some cause"),
			},
			contains: []string{"both present:", "some_value"},
			notContains: []string{"some cause"},
		},
		{
			name: "nil value",
			err: &EncodeError{
				Message: "nil value test",
				Value:   nil,
			},
			contains: []string{"nil value test"},
		},
		{
			name: "complex value type",
			err: &EncodeError{
				Message: "complex type",
				Value:   map[string]int{"key": 123},
			},
			contains: []string{"complex type:", "key", "123"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tt.err.Error()
			for _, s := range tt.contains {
				if !strings.Contains(msg, s) {
					t.Errorf("expected message to contain %q, got: %s", s, msg)
				}
			}
			for _, s := range tt.notContains {
				if strings.Contains(msg, s) {
					t.Errorf("expected message not to contain %q, got: %s", s, msg)
				}
			}
		})
	}
}

func TestEncodeError_Unwrap(t *testing.T) {
	tests := []struct {
		name      string
		err       *EncodeError
		wantCause error
	}{
		{
			name: "with cause",
			err: &EncodeError{
				Message: "wrapped error",
				Cause:   errors.New("root error"),
			},
			wantCause: errors.New("root error"),
		},
		{
			name: "without cause",
			err: &EncodeError{
				Message: "no cause",
				Cause:   nil,
			},
			wantCause: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Unwrap()
			if tt.wantCause == nil {
				if got != nil {
					t.Errorf("Unwrap() = %v, want nil", got)
				}
			} else {
				if got == nil {
					t.Errorf("Unwrap() = nil, want non-nil")
				} else if got.Error() != tt.wantCause.Error() {
					t.Errorf("Unwrap() = %v, want %v", got, tt.wantCause)
				}
			}
		})
	}
}

func TestEncodeError_ErrorsIs(t *testing.T) {
	rootCause := errors.New("root cause")
	encErr := &EncodeError{
		Message: "encode failed",
		Cause:   rootCause,
	}

	if !errors.Is(encErr, rootCause) {
		t.Error("errors.Is should find the root cause")
	}

	otherErr := errors.New("other error")
	if errors.Is(encErr, otherErr) {
		t.Error("errors.Is should not match unrelated error")
	}
}

func TestDecodeError_EdgeCases(t *testing.T) {
	tests := []struct {
		name string
		err  *DecodeError
	}{
		{
			name: "empty message",
			err:  &DecodeError{Message: ""},
		},
		{
			name: "very long message",
			err: &DecodeError{
				Message: strings.Repeat("a", 1000),
			},
		},
		{
			name: "special characters in token",
			err: &DecodeError{
				Message: "special chars",
				Token:   "{}[]\n\t\"'",
			},
		},
		{
			name: "multiline context",
			err: &DecodeError{
				Message: "multiline",
				Context: "line1\nline2\nline3\nline4\nline5",
			},
		},
		{
			name: "negative line number",
			err: &DecodeError{
				Message: "negative line",
				Line:    -1,
			},
		},
		{
			name: "negative column number",
			err: &DecodeError{
				Message: "negative column",
				Column:  -1,
			},
		},
		{
			name: "large line and column numbers",
			err: &DecodeError{
				Message: "large numbers",
				Line:    999999,
				Column:  888888,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			msg := tt.err.Error()
			if msg == "" && tt.err.Message != "" {
				t.Error("Error() returned empty string for non-empty message")
			}
		})
	}
}

func TestEncodeError_EdgeCases(t *testing.T) {
	tests := []struct {
		name string
		err  *EncodeError
	}{
		{
			name: "empty message",
			err:  &EncodeError{Message: ""},
		},
		{
			name: "very long message",
			err: &EncodeError{
				Message: strings.Repeat("x", 1000),
			},
		},
		{
			name: "nil everything",
			err: &EncodeError{
				Message: "nil test",
				Value:   nil,
				Cause:   nil,
			},
		},
		{
			name: "zero value",
			err: &EncodeError{
				Message: "zero",
				Value:   0,
			},
		},
		{
			name: "false value",
			err: &EncodeError{
				Message: "false",
				Value:   false,
			},
		},
		{
			name: "empty string value",
			err: &EncodeError{
				Message: "empty string",
				Value:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic
			msg := tt.err.Error()
			if msg == "" && tt.err.Message != "" {
				t.Error("Error() returned empty string for non-empty message")
			}
		})
	}
}