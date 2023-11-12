package frames

import (
	"encoding/json"
	"fmt"

	"github.com/chuckha/tagger/id3string"

	"gitlab.com/tozd/go/errors"
)

// UserDefinedURL have the identifier WXXX.
type UserDefinedURL struct {
	TextEncoding byte
	Description  []rune
	// URL is always ascii
	URL string
}

func (u *UserDefinedURL) UnmarshalBinary(data []byte) error {
	u.TextEncoding = data[0]
	info, n := id3string.ExtractNullTerminatedValueWithEncoding(u.TextEncoding, data[1:])
	u.Description = info
	u.URL = string(data[len(info)+1+n:])
	return nil
}

func (u *UserDefinedURL) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, u); err != nil {
		return errors.WithStack(err)
	}
	if !id3string.IsASCII(u.Description) {
		u.TextEncoding = 1
	}
	return nil
}

func (u *UserDefinedURL) String() string {
	return fmt.Sprintf("enc: %x; desc: %q; url: %q", u.TextEncoding, u.Description, u.URL)
}

func (u *UserDefinedURL) MarshalBinary() ([]byte, error) {
	out := []byte{u.TextEncoding}
	out = append(out, id3string.EncodeRunesWithNullTerminator(u.TextEncoding, u.Description)...)
	out = append(out, []byte(u.URL)...)
	return out, nil
}

func (u *UserDefinedURL) Equal(u2 *UserDefinedURL) bool {
	return u.TextEncoding == u2.TextEncoding &&
		id3string.Equal(u.Description, u2.Description) &&
		u.URL == u2.URL
}
