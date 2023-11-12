package frames

import (
	"encoding/json"
	"fmt"

	"github.com/chuckha/tagger/id3string"

	"gitlab.com/tozd/go/errors"
)

// UserDefinedTextInformation
type UserDefinedTextInformation struct {
	TextEncoding byte
	Description  []rune
	Value        []rune
}

func (u *UserDefinedTextInformation) UnmarshalBinary(data []byte) error {
	u.TextEncoding = data[0]
	ptr := 1
	desc, n := id3string.ExtractNullTerminatedValueWithEncoding(u.TextEncoding, data[ptr:])
	u.Description = desc
	ptr += len(u.Description) + n
	u.Value = []rune(string(data[ptr:]))
	return nil
}

func (u *UserDefinedTextInformation) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, u); err != nil {
		return errors.WithStack(err)
	}
	if !id3string.IsASCII(u.Description) || !id3string.IsASCII(u.Value) {
		u.TextEncoding = 1
	}
	return nil
}

func (u *UserDefinedTextInformation) String() string {
	return fmt.Sprintf("enc: %x; desc: %q; value: %q", u.TextEncoding, u.Description, u.Value)
}

func (u *UserDefinedTextInformation) MarshalBinary() ([]byte, error) {
	out := []byte{u.TextEncoding}
	out = append(out, id3string.EncodeRunesWithNullTerminator(u.TextEncoding, u.Description)...)
	out = append(out, id3string.EncodeRunes(u.TextEncoding, u.Value)...)
	return out, nil
}

func (u *UserDefinedTextInformation) Equal(u2 *UserDefinedTextInformation) bool {
	return u.TextEncoding == u2.TextEncoding &&
		id3string.Equal(u.Description, u2.Description) &&
		id3string.Equal(u.Value, u2.Value)
}
