package frames

import "fmt"

// TermsOfUse have the ID USER
type TermsOfUse struct {
	TextEncoding byte
	Language     string
	Text         string
}

func (t *TermsOfUse) UnmarshalBinary(data []byte) error {
	t.TextEncoding = data[0]
	ptr := 1
	t.Language = string(data[ptr : ptr+3])
	ptr += 3
	t.Text = string(data[ptr:])
	return nil
}

func (t *TermsOfUse) String() string {
	return fmt.Sprintf("enc: %x; lang: %q; text: %q", t.TextEncoding, t.Language, t.Text)
}

func (t *TermsOfUse) MarshalBinary() ([]byte, error) {
	out := []byte{t.TextEncoding}
	out = append(out, []byte(t.Language)...)
	out = append(out, []byte(t.Text)...)
	return out, nil
}

func (t *TermsOfUse) Equal(t2 *TermsOfUse) bool {
	return t.TextEncoding == t2.TextEncoding &&
		t.Language == t2.Language &&
		t.Text == t2.Text
}
