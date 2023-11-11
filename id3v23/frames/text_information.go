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
	Information  string
}

func NewTextInformation(info string) *TextInformation {
	ti := &TextInformation{Information: info}
	if id3string.IsUnicode(info) {
		ti.TextEncoding = 1
	}
	return ti
}

func (t *TextInformation) UnmarshalBinary(data []byte) error {
	t.TextEncoding = data[0]
	// this extracts the string that is either null terminated; double null terminated; or all the bytes.
	info, _ := id3string.ExtractStringFromEncoding(t.TextEncoding, data[1:])
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
	t.Information = in.Information
	if id3string.IsUnicode(t.Information) {
		t.TextEncoding = 1
	}
	return nil
}

func (t *TextInformation) MarshalBinary() ([]byte, error) {
	return append([]byte{t.TextEncoding}, id3string.EncodeString(t.TextEncoding, t.Information)...), nil
}

func (t *TextInformation) String() string {
	return fmt.Sprintf("enc: %x; info: %q", t.TextEncoding, t.Information)
}

func (t *TextInformation) Equal(t2 *TextInformation) bool {
	return t.TextEncoding == t2.TextEncoding &&
		t.Information == t2.Information
}
