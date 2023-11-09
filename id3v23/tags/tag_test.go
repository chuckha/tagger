package tags

import (
	"testing"
)

func TestID3v2_MarshalBinary(t *testing.T) {
	tag, err := NewID3v2FromFile("testdata/some_tag.mp3")
	if err != nil {
		t.Error(err)
	}
	tag.Frames.SetTextInformationFrame("TIT3", "title")
	b, err := tag.MarshalBinary()
	if err != nil {
		t.Error(err)
	}
	if len(b) != tag.Header.Size+10 {
		t.Errorf("expected %d bytes, got %d", tag.Header.Size+10, len(b))
	}
}
