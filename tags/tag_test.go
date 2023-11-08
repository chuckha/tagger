package tags

import (
	"fmt"
	"testing"
)

func TestID3v2_MarshalBinary(t *testing.T) {
	f, err := NewID3v2FromFile("testdata/some_tag.mp3")
	if err != nil {
		t.Error(err)
	}
	f.Frames.SetTextInformationFrame("TIT3", "title")
	b, err := f.MarshalBinary()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(f)
	fmt.Println(b)
}
