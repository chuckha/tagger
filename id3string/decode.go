package id3string

import "unicode/utf8"

// DecodeUTF8 is used when we have a JSON representation of a Frame.
func DecodeUTF8(in string) []rune {
	out := []rune{}
	for len(in) > 0 {
		r, size := utf8.DecodeRune([]byte(in))
		in = in[size:]
		out = append(out, r)
	}
	return out
}
