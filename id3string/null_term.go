package id3string

import (
	"bytes"
	"fmt"
)

func ExtractStringFromEncoding(enc byte, data []byte) (string, int) {
	switch enc {
	case 0:
		return ExtractNullTerminated(data), 1
	case 1:
		return ExtractNullTerminatedUnicode(data), 2
	case 2:
		fmt.Println("warning: id3v2.4 (2) encoding used")
		return ExtractNullTerminated(data), 1
	case 3:
		fmt.Println("warning: id3v2.4 (3) encoding used")
		return ExtractNullTerminated(data), 1
	default:
		panic(fmt.Sprintf("unhandled text encoding: %08b", enc))
	}
}

// ExtractNullTerminated is to be used when only a single null terminator ends a string.
func ExtractNullTerminated(b []byte) string {
	n := bytes.IndexByte(b, 0)
	if n == -1 {
		return string(b)
	}
	return string(b[:n])
}

// ExtractNullTerminatedUnicode looks for two null bytes to end a string.
func ExtractNullTerminatedUnicode(b []byte) string {
	n := bytes.Index(b, []byte{0, 0})
	if n == -1 {
		return string(b)
	}
	return string(b[:n])
}
