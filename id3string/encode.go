package id3string

import "unicode/utf16"

func EncodeASCIIWithNullTerminator(val string) []byte {
	return append([]byte(val), '\x00')
}

// EncodeRunesWithNullTerminator adds all the extra bytes that id3v2.3 expects.
// In the case where the information is only ascii, (enc: 0), simply add a null terminator
// In the case of UTF-16, (enc: 1), add a BOM, then the string, then two null terminators (unicode null).
func EncodeRunesWithNullTerminator(enc byte, val []rune) []byte {
	runes := EncodeRunes(enc, val)
	switch enc {
	case 0:
		return append(runes, '\x00')
	case 1:
		return append(runes, '\x00', '\x00')
	default:
		panic("unknown encoding for id3v2.3")
	}
}

func EncodeRunes(enc byte, val []rune) []byte {
	switch enc {
	case 0:
		return []byte(string(val))
	case 1:
		encoded := utf16.Encode(val)
		bytes := []byte{}
		for _, c := range encoded {
			bytes = append(bytes, byte(c>>8), byte(c&0xFF))
		}
		bom := []byte{'\xFE', '\xFF'}
		return append(bom, bytes...)
	default:
		panic("unknown encoding for id3v2.3")
	}
}

// IsASCII returns true if it's only ascii; otherwise assume UTF-16 since UTF-8 is not used.
func IsASCII(in []rune) bool {
	for _, c := range in {
		if c > 127 {
			return false
		}
	}
	return true
}
func IsASCIIBytes(in []byte) bool {
	for _, c := range in {
		if c > 127 {
			return false
		}
	}
	return true
}

func Equal(a, b []rune) bool {
	if len(a) != len(b) {
		return false
	}
	for i, c := range a {
		if c != b[i] {
			return false
		}
	}
	return true
}

func EqualBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, c := range a {
		if c != b[i] {
			return false
		}
	}
	return true
}
