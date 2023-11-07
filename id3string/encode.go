package id3string

import "fmt"

func EncodeString(enc byte, s string) []byte {
	switch enc {
	case 0, 2, 3:
		return append([]byte(s), '\x00')
	case 1:
		return append([]byte(s), '\x00', '\x00')
	default:
		panic(fmt.Sprintf("unhandled text encoding: %08b", enc))
	}
}
