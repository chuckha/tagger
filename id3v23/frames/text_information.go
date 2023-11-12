package frames

import (
	"encoding/json"
	"fmt"

	"github.com/chuckha/tagger/id3string"

	"gitlab.com/tozd/go/errors"
)

// TextInformation are all of the text frames.
// Text frames have IDs of T000-TZZZ excluding TXXX.
type TextInformation struct {
	TextEncoding byte `json:"-"`
	Information  []rune
}

func NewTextInformation(info string) *TextInformation {
	// encode to utf-16
	val := []rune(info)
	// check if any value is outside of ascii
	ti := &TextInformation{Information: val}
	if !id3string.IsASCII(val) {
		ti.TextEncoding = 1
	}
	return ti
}

func (t *TextInformation) UnmarshalBinary(data []byte) error {
	t.TextEncoding = data[0]
	// this extracts the string that is either null terminated; double null terminated; or all the bytes.
	info, _ := id3string.ExtractValueWithEncoding(t.TextEncoding, data[1:])
	t.Information = info
	return nil
}

func (t *TextInformation) UnmarshalJSON(data []byte) error {
	var in struct {
		Information string
	}
	if err := json.Unmarshal(data, &in); err != nil {
		return errors.WithStack(err)
	}
	t.Information = []rune(in.Information)
	if !id3string.IsASCII(t.Information) {
		t.TextEncoding = 1
	}
	return nil
}

func (t *TextInformation) MarshalBinary() ([]byte, error) {
	return append([]byte{t.TextEncoding}, id3string.EncodeRunes(t.TextEncoding, t.Information)...), nil
}

func (t *TextInformation) String() string {
	return fmt.Sprintf("enc: %x; info: %v", t.TextEncoding, t.Information)
}

func (t *TextInformation) Equal(t2 *TextInformation) bool {
	return t.TextEncoding == t2.TextEncoding &&
		id3string.Equal(t.Information, t2.Information)
}
