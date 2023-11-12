package frames

import (
	"fmt"

	"github.com/chuckha/tagger/id3string"
)

// UnsynchronizedLyrics have an ID of USLT.
type UnsynchronizedLyrics struct {
	TextEncoding      byte
	Language          string
	ContentDescriptor []rune
	Lyrics            string
}

func (u *UnsynchronizedLyrics) UnmarshalBinary(data []byte) error {
	u.TextEncoding = data[0]
	ptr := 1
	u.Language = string(data[ptr : ptr+3])
	ptr += 3
	contentDesc, n := id3string.ExtractNullTerminatedValueWithEncoding(u.TextEncoding, data[ptr:])
	u.ContentDescriptor = contentDesc
	ptr += len(u.ContentDescriptor) + n
	u.Lyrics = string(data[ptr:])
	return nil
}

func (u *UnsynchronizedLyrics) String() string {
	return fmt.Sprintf("enc: %x; lang: %q; desc: %q; lyrics: %q", u.TextEncoding, u.Language, u.ContentDescriptor, u.Lyrics)
}

func (u *UnsynchronizedLyrics) MarshalBinary() ([]byte, error) {
	out := []byte{u.TextEncoding}
	out = append(out, []byte(u.Language)...)
	out = append(out, id3string.EncodeRunesWithNullTerminator(u.TextEncoding, u.ContentDescriptor)...)
	out = append(out, []byte(u.Lyrics)...)
	return out, nil
}

func (u *UnsynchronizedLyrics) Equal(u2 *UnsynchronizedLyrics) bool {
	return u.TextEncoding == u2.TextEncoding &&
		u.Language == u2.Language &&
		id3string.Equal(u.ContentDescriptor, u2.ContentDescriptor) &&
		u.Lyrics == u2.Lyrics
}
