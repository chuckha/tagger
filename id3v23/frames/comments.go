package frames

import (
	"encoding/json"
	"fmt"

	"github.com/chuckha/tagger/id3string"

	"gitlab.com/tozd/go/errors"
)

// Comment are all of the comment frames.
// Comment frames have an ID of COMM.
type Comment struct {
	TextEncoding            byte
	Language                string
	ShortContentDescription []rune
	ActualText              []rune
}

func (c *Comment) UnmarshalBinary(data []byte) error {
	ptr := 0
	c.TextEncoding = data[0]
	ptr++
	c.Language = string(data[1:4])
	ptr += 3
	desc, n := id3string.ExtractNullTerminatedValueWithEncoding(c.TextEncoding, data[ptr:])
	c.ShortContentDescription = desc
	ptr += len(desc) + n
	at, n := id3string.ExtractNullTerminatedValueWithEncoding(c.TextEncoding, data[ptr:])
	c.ActualText = at
	ptr += n
	// TODO: check length maybe?
	return nil
}

func (c *Comment) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, c); err != nil {
		return errors.WithStack(err)
	}
	if !id3string.IsASCII(c.ShortContentDescription) || !id3string.IsASCII(c.ActualText) {
		c.TextEncoding = 1
	}
	return nil
}

func (c *Comment) MarshalBinary() ([]byte, error) {
	out := []byte{c.TextEncoding}
	out = append(out, []byte(c.Language)...)
	out = append(out, id3string.EncodeRunesWithNullTerminator(c.TextEncoding, c.ShortContentDescription)...)
	out = append(out, id3string.EncodeRunes(c.TextEncoding, c.ActualText)...)
	return out, nil
}

func (c *Comment) String() string {
	return fmt.Sprintf("enc: %x; lang: %q; short: %q; text: %q", c.TextEncoding, c.Language, c.ShortContentDescription, c.ActualText)
}

func (c *Comment) Equal(c2 *Comment) bool {
	return c.TextEncoding == c2.TextEncoding &&
		c.Language == c2.Language &&
		id3string.Equal(c.ShortContentDescription, c2.ShortContentDescription) &&
		id3string.Equal(c.ActualText, c2.ActualText)
}
