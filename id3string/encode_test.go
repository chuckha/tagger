package id3string

import "testing"

func TestIsUnicode(t *testing.T) {
	tests := []struct {
		in   []rune
		want bool
	}{
		{[]rune(""), false},
		{[]rune("a"), false},
		{[]rune("\x00"), false},
		{[]rune("\x01"), false},
		{[]rune("\x7f"), false},
		{[]rune("\x80"), true},
		{[]rune("\xff"), true},
		{[]rune("\u0100"), true},
		{[]rune("\u07ff"), true},
		{[]rune("\u0800"), true},
		{[]rune("\uffff"), true},
		{[]rune("\U00010000"), true},
		{[]rune("\U0010ffff"), true},
		{[]rune("日本語"), true},
	}
	for _, tt := range tests {
		t.Run(string(tt.in), func(t *testing.T) {
			if got := !IsASCII(tt.in); got != tt.want {
				t.Errorf("IsUnicode(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}
