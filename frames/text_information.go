package frames

import (
	"fmt"
	"tagger/id3string"
)

// TextInformation are all of the text frames.
// Text frames have IDs of T000-TZZZ excluding TXXX.
type TextInformation struct {
	TextEncoding byte
	Information  string
}

func (t *TextInformation) UnmarshalBinary(data []byte) error {
	t.TextEncoding = data[0]
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
