package id3string

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unicode/utf16"
)

func ExtractValueWithEncoding(enc byte, data []byte) ([]rune, int) {
	switch enc {
	case 0:
		return []rune(string(data)), 0
	case 1:
		return ExtractUnicode(data), 2 // consume 2 BOM bytes
	default:
		panic(fmt.Sprintf("unhandled text encoding: %08b", enc))
	}
}

func ExtractNullTerminatedValueWithEncoding(enc byte, data []byte) ([]rune, int) {
	switch enc {
	case 0:
		return ExtractNullTerminated(data), 1
	case 1:
		// 4 are the BOM and the unicode null terminator
		return ExtractUnicodeNullTerminated(data), 4
	default:
		panic(fmt.Sprintf("unhandled text encoding: %08b", enc))
	}
}

func ExtractNullTerminatedASCII(b []byte) string {
	n := bytes.IndexByte(b, 0)
	if n == -1 {
		return string(b)
	}
	return string(b[:n])
}

// ExtractNullTerminated is to be used when only a single null terminator ends a string.
func ExtractNullTerminated(b []byte) []rune {
	n := bytes.IndexByte(b, 0)
	if n == -1 {
		return []rune(string(b))
	}
	return []rune(string(b[:n]))
}

// ExtractNullTerminatedUnicode gets the two BOM bytes, looks for a unicode null, then extracts the middle.
func ExtractUnicodeNullTerminated(b []byte) []rune {
	// TODO: use the bom to determine if it's big or little endian; for now assume big endian
	_ = b[0:2]
	n := bytes.Index(b[2:], []byte{0, 0})
	if n == -1 {
		return []rune(string(b))
	}
	return bytesToRunes(b[2:n])
}

func ExtractUnicode(b []byte) []rune {
	// TODO: use the bom to determine if it's big or little endian; for now assume big endian
	_ = b[0:2]
	return bytesToRunes(b[2:])
}

func bytesToRunes(b []byte) []rune {
	// Check if byte slice length is even
	if len(b)%2 != 0 {
		panic("byte slice has odd length, cannot be valid UTF-16")
	}

	// Convert []byte to []uint16
	uints := make([]uint16, 0, len(b)/2)
	for i := 0; i < len(b); i += 2 {
		// For big endian, use binary.BigEndian.Uint16
		// For little endian, use binary.LittleEndian.Uint16
		uints = append(uints, binary.BigEndian.Uint16(b[i:i+2]))
	}

	// Decode []uint16 to []rune
	return utf16.Decode(uints)
}
