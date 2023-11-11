package id3string

import (
	"fmt"
	"unicode/utf8"
)

func EncodeString(enc byte, s string) []byte {
	switch enc {
	case 0, 3:
		return append([]byte(s), '\x00')
	case 1, 2:
		return append([]byte(s), '\x00', '\x00')
	default:
		panic(fmt.Sprintf("unhandled text encoding: %08b", enc))
	}
}

func TextEncoding(s string) byte {
	if IsUnicode(s) {
		return 1
	}
	return 0
}

func IsUnicode(in string) bool {
	for len(in) > 0 {
		r, size := utf8.DecodeRuneInString(in)
		if r == utf8.RuneError {
			// This means we have something that's kind of unicode but is a bit broken.
			// So treat it as unicode
			return true
		}
		if size > 1 {
			return true
		}
		in = in[size:]
	}
	return false
}
