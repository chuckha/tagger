package id3string

import "testing"

func TestIsUnicode(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"", false},
		{"a", false},
		{"\x00", false},
		{"\x01", false},
		{"\x7f", false},
		{"\x80", true},
		{"\xff", true},
		{"\u0100", true},
		{"\u07ff", true},
		{"\u0800", true},
		{"\uffff", true},
		{"\U00010000", true},
		{"\U0010ffff", true},
		{"日本語", true},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := IsUnicode(tt.in); got != tt.want {
				t.Errorf("IsUnicode(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}
