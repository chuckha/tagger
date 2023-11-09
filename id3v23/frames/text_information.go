package frames

import (
	"fmt"
	"tagger/id3string"
	"unicode/utf8"
)

// TextInformation are all of the text frames.
// Text frames have IDs of T000-TZZZ excluding TXXX.
type TextInformation struct {
	TextEncoding byte
	Information  string
}

func isUnicode(in string) bool {
	for len(in) > 0 {
		_, size := utf8.DecodeRuneInString(in)
		if size > 1 {
			return true
		}
		in = in[size:]
	}
	return false
}

func NewTextInformation(info string) *TextInformation {
	ti := &TextInformation{Information: info}
	if isUnicode(info) {
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
