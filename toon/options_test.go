package toon

import (
	"testing"
)

// TestValidateEncodeOptions tests the validateEncodeOptions function
func TestValidateEncodeOptions(t *testing.T) {
	tests := []struct {
		name    string
		opts    *EncodeOptions
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil options",
			opts:    nil,
			wantErr: false,
		},
		{
			name:    "valid default indent",
			opts:    &EncodeOptions{Indent: 2},
			wantErr: false,
		},
		{
			name:    "valid zero indent",
			opts:    &EncodeOptions{Indent: 0},
			wantErr: false,
		},
		{
			name:    "valid large indent",
			opts:    &EncodeOptions{Indent: 8},
			wantErr: false,
		},
		{
			name:    "valid extremely large indent",
			opts:    &EncodeOptions{Indent: 100},
			wantErr: false,
		},
		{
			name:    "invalid negative indent",
			opts:    &EncodeOptions{Indent: -1},
			wantErr: true,
			errMsg:  "indent must be non-negative",
		},
		{
			name:    "invalid large negative indent",
			opts:    &EncodeOptions{Indent: -100},
			wantErr: true,
			errMsg:  "indent must be non-negative",
		},
		{
			name:    "valid comma delimiter",
			opts:    &EncodeOptions{Delimiter: comma},
			wantErr: false,
		},
		{
			name:    "valid tab delimiter",
			opts:    &EncodeOptions{Delimiter: tab},
			wantErr: false,
		},
		{
			name:    "valid pipe delimiter",
			opts:    &EncodeOptions{Delimiter: pipe},
			wantErr: false,
		},
		{
			name:    "empty delimiter",
			opts:    &EncodeOptions{Delimiter: ""},
			wantErr: false,
		},
		{
			name:    "invalid multi-char delimiter",
			opts:    &EncodeOptions{Delimiter: ";;"},
			wantErr: true,
			errMsg:  "invalid delimiter",
		},
		{
			name:    "invalid semicolon delimiter",
			opts:    &EncodeOptions{Delimiter: ";"},
			wantErr: true,
			errMsg:  "invalid delimiter",
		},
		{
			name:    "invalid space delimiter",
			opts:    &EncodeOptions{Delimiter: " "},
			wantErr: true,
			errMsg:  "invalid delimiter",
		},
		{
			name:    "invalid colon delimiter",
			opts:    &EncodeOptions{Delimiter: ":"},
			wantErr: true,
			errMsg:  "invalid delimiter",
		},
		{
			name:    "invalid newline delimiter",
			opts:    &EncodeOptions{Delimiter: "\n"},
			wantErr: true,
			errMsg:  "invalid delimiter",
		},
		{
			name: "valid options combination",
			opts: &EncodeOptions{
				Indent:       4,
				Delimiter:    comma,
				LengthMarker: "#",
			},
			wantErr: false,
		},
		{
			name: "valid with flatten paths",
			opts: &EncodeOptions{
				Indent:       2,
				FlattenPaths: true,
				FlattenDepth: 3,
			},
			wantErr: false,
		},
		{
			name: "invalid indent with valid delimiter",
			opts: &EncodeOptions{
				Indent:    -5,
				Delimiter: comma,
			},
			wantErr: true,
			errMsg:  "indent must be non-negative",
		},
		{
			name: "valid indent with invalid delimiter",
			opts: &EncodeOptions{
				Indent:    4,
				Delimiter: "!",
			},
			wantErr: true,
			errMsg:  "invalid delimiter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEncodeOptions(tt.opts)
			if tt.wantErr {
				if err == nil {
					t.Errorf("validateEncodeOptions() expected error containing %q, got nil", tt.errMsg)
					return
				}
				if tt.errMsg != "" {
					errStr := err.Error()
					if len(errStr) == 0 || errStr[:len(tt.errMsg)] != tt.errMsg {
						// Check if error message contains the expected substring
						found := false
						for i := 0; i <= len(errStr)-len(tt.errMsg); i++ {
							if errStr[i:i+len(tt.errMsg)] == tt.errMsg {
								found = true
								break
							}
						}
						if !found {
							t.Errorf("validateEncodeOptions() error = %q, want error containing %q", errStr, tt.errMsg)
						}
					}
				}
			} else {
				if err != nil {
					t.Errorf("validateEncodeOptions() unexpected error = %v", err)
				}
			}
		})
	}
}

// TestValidateDecodeOptions tests the validateDecodeOptions function
func TestValidateDecodeOptions(t *testing.T) {
	tests := []struct {
		name    string
		opts    *DecodeOptions
		wantErr bool
		errMsg  string
	}{
		{
			name:    "nil options",
			opts:    nil,
			wantErr: false,
		},
		{
			name:    "valid default indent size",
			opts:    &DecodeOptions{IndentSize: 2},
			wantErr: false,
		},
		{
			name:    "valid indent size 1",
			opts:    &DecodeOptions{IndentSize: 1},
			wantErr: false,
		},
		{
			name:    "valid indent size 4",
			opts:    &DecodeOptions{IndentSize: 4},
			wantErr: false,
		},
		{
			name:    "valid large indent size",
			opts:    &DecodeOptions{IndentSize: 8},
			wantErr: false,
		},
		{
			name:    "valid extremely large indent size",
			opts:    &DecodeOptions{IndentSize: 100},
			wantErr: false,
		},
		{
			name:    "invalid zero indent size",
			opts:    &DecodeOptions{IndentSize: 0},
			wantErr: true,
			errMsg:  "indent_size must be positive",
		},
		{
			name:    "invalid negative indent size",
			opts:    &DecodeOptions{IndentSize: -1},
			wantErr: true,
			errMsg:  "indent_size must be positive",
		},
		{
			name:    "invalid large negative indent size",
			opts:    &DecodeOptions{IndentSize: -100},
			wantErr: true,
			errMsg:  "indent_size must be positive",
		},
		{
			name: "valid with strict keys",
			opts: &DecodeOptions{
				IndentSize: 2,
				Strict:     true,
			},
			wantErr: false,
		},
		{
			name: "valid with non-strict keys",
			opts: &DecodeOptions{
				IndentSize: 2,
				Strict:     false,
			},
			wantErr: false,
		},
		{
			name: "valid with expand paths off",
			opts: &DecodeOptions{
				IndentSize:  2,
				ExpandPaths: "off",
			},
			wantErr: false,
		},
		{
			name: "valid with expand paths safe",
			opts: &DecodeOptions{
				IndentSize:  2,
				ExpandPaths: "safe",
			},
			wantErr: false,
		},
		{
			name: "valid with empty expand paths",
			opts: &DecodeOptions{
				IndentSize:  2,
				ExpandPaths: "",
			},
			wantErr: false,
		},
		{
			name: "valid combination of options",
			opts: &DecodeOptions{
				Keys:        StringKeys,
				Strict:      true,
				IndentSize:  4,
				ExpandPaths: "safe",
			},
			wantErr: false,
		},
		{
			name: "valid with all flags",
			opts: &DecodeOptions{
				Keys:        StringKeys,
				Strict:      true,
				IndentSize:  2,
				ExpandPaths: "off",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDecodeOptions(tt.opts)
			if tt.wantErr {
				if err == nil {
					t.Errorf("validateDecodeOptions() expected error containing %q, got nil", tt.errMsg)
					return
				}
				if tt.errMsg != "" {
					errStr := err.Error()
					if len(errStr) == 0 || errStr[:len(tt.errMsg)] != tt.errMsg {
						// Check if error message contains the expected substring
						found := false
						for i := 0; i <= len(errStr)-len(tt.errMsg); i++ {
							if errStr[i:i+len(tt.errMsg)] == tt.errMsg {
								found = true
								break
							}
						}
						if !found {
							t.Errorf("validateDecodeOptions() error = %q, want error containing %q", errStr, tt.errMsg)
						}
					}
				}
			} else {
				if err != nil {
					t.Errorf("validateDecodeOptions() unexpected error = %v", err)
				}
			}
		})
	}
}

// TestIsValidDelimiter tests the isValidDelimiter function
func TestIsValidDelimiter(t *testing.T) {
	tests := []struct {
		name      string
		delimiter string
		want      bool
	}{
		{
			name:      "valid comma",
			delimiter: comma,
			want:      true,
		},
		{
			name:      "valid tab",
			delimiter: tab,
			want:      true,
		},
		{
			name:      "valid pipe",
			delimiter: pipe,
			want:      true,
		},
		{
			name:      "invalid empty string",
			delimiter: "",
			want:      false,
		},
		{
			name:      "invalid space",
			delimiter: " ",
			want:      false,
		},
		{
			name:      "invalid semicolon",
			delimiter: ";",
			want:      false,
		},
		{
			name:      "invalid colon",
			delimiter: ":",
			want:      false,
		},
		{
			name:      "invalid multi-char semicolons",
			delimiter: ";;",
			want:      false,
		},
		{
			name:      "invalid multi-char commas",
			delimiter: ",,",
			want:      false,
		},
		{
			name:      "invalid newline",
			delimiter: "\n",
			want:      false,
		},
		{
			name:      "invalid carriage return",
			delimiter: "\r",
			want:      false,
		},
		{
			name:      "invalid dash",
			delimiter: "-",
			want:      false,
		},
		{
			name:      "invalid underscore",
			delimiter: "_",
			want:      false,
		},
		{
			name:      "invalid period",
			delimiter: ".",
			want:      false,
		},
		{
			name:      "invalid slash",
			delimiter: "/",
			want:      false,
		},
		{
			name:      "invalid backslash",
			delimiter: "\\",
			want:      false,
		},
		{
			name:      "invalid exclamation",
			delimiter: "!",
			want:      false,
		},
		{
			name:      "invalid question mark",
			delimiter: "?",
			want:      false,
		},
		{
			name:      "invalid at sign",
			delimiter: "@",
			want:      false,
		},
		{
			name:      "invalid hash",
			delimiter: "#",
			want:      false,
		},
		{
			name:      "invalid dollar",
			delimiter: "$",
			want:      false,
		},
		{
			name:      "invalid percent",
			delimiter: "%",
			want:      false,
		},
		{
			name:      "invalid caret",
			delimiter: "^",
			want:      false,
		},
		{
			name:      "invalid ampersand",
			delimiter: "&",
			want:      false,
		},
		{
			name:      "invalid asterisk",
			delimiter: "*",
			want:      false,
		},
		{
			name:      "invalid plus",
			delimiter: "+",
			want:      false,
		},
		{
			name:      "invalid equals",
			delimiter: "=",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidDelimiter(tt.delimiter)
			if got != tt.want {
				t.Errorf("isValidDelimiter(%q) = %v, want %v", tt.delimiter, got, tt.want)
			}
		})
	}
}

// TestGetEncodeOptions tests the getEncodeOptions function with defaults
func TestGetEncodeOptions(t *testing.T) {
	tests := []struct {
		name string
		opts *EncodeOptions
		want *EncodeOptions
	}{
		{
			name: "nil options apply defaults",
			opts: nil,
			want: &EncodeOptions{
				Indent:       defaultIndent,
				Delimiter:    defaultDelimiter,
				FlattenDepth: 0,
			},
		},
		{
			name: "empty options apply defaults",
			opts: &EncodeOptions{},
			want: &EncodeOptions{
				Indent:       defaultIndent,
				Delimiter:    defaultDelimiter,
				FlattenDepth: 0,
			},
		},
		{
			name: "custom indent preserved",
			opts: &EncodeOptions{Indent: 4},
			want: &EncodeOptions{
				Indent:       4,
				Delimiter:    defaultDelimiter,
				FlattenDepth: 0,
			},
		},
		{
			name: "custom delimiter preserved",
			opts: &EncodeOptions{Delimiter: tab},
			want: &EncodeOptions{
				Indent:       defaultIndent,
				Delimiter:    tab,
				FlattenDepth: 0,
			},
		},
		{
			name: "flatten paths with default depth",
			opts: &EncodeOptions{
				FlattenPaths: true,
				FlattenDepth: -1,
			},
			want: &EncodeOptions{
				Indent:       defaultIndent,
				Delimiter:    defaultDelimiter,
				FlattenPaths: true,
				FlattenDepth: 9999,
			},
		},
		{
			name: "flatten paths with explicit depth",
			opts: &EncodeOptions{
				FlattenPaths: true,
				FlattenDepth: 3,
			},
			want: &EncodeOptions{
				Indent:       defaultIndent,
				Delimiter:    defaultDelimiter,
				FlattenPaths: true,
				FlattenDepth: 3,
			},
		},
		{
			name: "flatten paths disabled",
			opts: &EncodeOptions{
				FlattenPaths: false,
				FlattenDepth: 0,
			},
			want: &EncodeOptions{
				Indent:       defaultIndent,
				Delimiter:    defaultDelimiter,
				FlattenPaths: false,
				FlattenDepth: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getEncodeOptions(tt.opts)
			if got.Indent != tt.want.Indent {
				t.Errorf("getEncodeOptions().Indent = %v, want %v", got.Indent, tt.want.Indent)
			}
			if got.Delimiter != tt.want.Delimiter {
				t.Errorf("getEncodeOptions().Delimiter = %v, want %v", got.Delimiter, tt.want.Delimiter)
			}
			if got.FlattenPaths != tt.want.FlattenPaths {
				t.Errorf("getEncodeOptions().FlattenPaths = %v, want %v", got.FlattenPaths, tt.want.FlattenPaths)
			}
			if got.FlattenDepth != tt.want.FlattenDepth {
				t.Errorf("getEncodeOptions().FlattenDepth = %v, want %v", got.FlattenDepth, tt.want.FlattenDepth)
			}
		})
	}
}

// TestGetDecodeOptions tests the getDecodeOptions function with defaults
func TestGetDecodeOptions(t *testing.T) {
	tests := []struct {
		name string
		opts *DecodeOptions
		want *DecodeOptions
	}{
		{
			name: "nil options apply defaults",
			opts: nil,
			want: &DecodeOptions{
				Keys:        StringKeys,
				Strict:      true,
				IndentSize:  defaultIndent,
				ExpandPaths: "off",
			},
		},
		{
			name: "empty options apply defaults",
			opts: &DecodeOptions{},
			want: &DecodeOptions{
				Keys:        StringKeys,
				Strict:      false,
				IndentSize:  defaultIndent,
				ExpandPaths: "off",
			},
		},
		{
			name: "custom indent size preserved",
			opts: &DecodeOptions{IndentSize: 4},
			want: &DecodeOptions{
				Keys:        StringKeys,
				Strict:      false,
				IndentSize:  4,
				ExpandPaths: "off",
			},
		},
		{
			name: "strict mode preserved",
			opts: &DecodeOptions{
				IndentSize: 2,
				Strict:     true,
			},
			want: &DecodeOptions{
				Keys:        StringKeys,
				Strict:      true,
				IndentSize:  2,
				ExpandPaths: "off",
			},
		},
		{
			name: "expand paths safe",
			opts: &DecodeOptions{
				IndentSize:  2,
				ExpandPaths: "safe",
			},
			want: &DecodeOptions{
				Keys:        StringKeys,
				Strict:      false,
				IndentSize:  2,
				ExpandPaths: "safe",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getDecodeOptions(tt.opts)
			if got.Keys != tt.want.Keys {
				t.Errorf("getDecodeOptions().Keys = %v, want %v", got.Keys, tt.want.Keys)
			}
			if got.Strict != tt.want.Strict {
				t.Errorf("getDecodeOptions().Strict = %v, want %v", got.Strict, tt.want.Strict)
			}
			if got.IndentSize != tt.want.IndentSize {
				t.Errorf("getDecodeOptions().IndentSize = %v, want %v", got.IndentSize, tt.want.IndentSize)
			}
			if got.ExpandPaths != tt.want.ExpandPaths {
				t.Errorf("getDecodeOptions().ExpandPaths = %v, want %v", got.ExpandPaths, tt.want.ExpandPaths)
			}
		})
	}
}
