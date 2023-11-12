package tags

import (
	"strings"
	"testing"

	"github.com/chuckha/tagger/id3v23/frames"
)

func TestID3v2_MarshalBinary(t *testing.T) {
	t.Run("the new tag is bigger than the original tag", func(t *testing.T) {
		tag := createTag(t)
		out, err := tag.MarshalBinary()
		if err != nil {
			t.Error(err)
		}
		tag.Header.Size = len(out) - 10

		err = tag.ApplyFrames(map[string]frames.FrameBody{
			"TIT1": &frames.TextInformation{Information: []rune("test99")},
			"TPOS": &frames.TextInformation{Information: []rune("1/1")},
			"TPE1": &frames.TextInformation{Information: []rune(strings.Repeat("a", 10000))}})
		if err != nil {
			t.Error(err)
		}
		out2, err := tag.MarshalBinary()
		if err != nil {
			t.Error(err)
		}
		if len(out2) < len(out) {
			t.Errorf("expected %d to be greater than %d", len(out2), len(out))
		}
	})

	t.Run("the new tag is smaller than the original tag, but only a bit", func(t *testing.T) {
		tag := createTag(t)
		out, err := tag.MarshalBinary()
		if err != nil {
			t.Error(err)
		}
		tag.Header.Size = len(out) - 10

		err = tag.ApplyFrames(map[string]frames.FrameBody{
			"TIT1": &frames.TextInformation{Information: []rune("t")}})
		if err != nil {
			t.Error(err)
		}
		out2, err := tag.MarshalBinary()
		if err != nil {
			t.Error(err)
		}
		if len(out2) != len(out) {
			t.Errorf("expected %d to be equal to %d", len(out2), len(out))
		}
	})

	t.Run("the new tag is significantly smaller than the original tag", func(t *testing.T) {
		bigFrame := &frames.Frame{
			Header: &frames.FrameHeader{ID: "TPE1"},
			Body:   &frames.TextInformation{Information: []rune(strings.Repeat("a", 10000))},
		}
		tag := createTag(t, bigFrame)
		out, err := tag.MarshalBinary()
		if err != nil {
			t.Error(err)
		}
		tag.Header.Size = len(out) - 10

		err = tag.ApplyFrames(map[string]frames.FrameBody{
			"TIT1": &frames.TextInformation{Information: []rune("t")},
			"TPE1": &frames.TextInformation{Information: []rune("smaller")}})
		if err != nil {
			t.Error(err)
		}
		out2, err := tag.MarshalBinary()
		if err != nil {
			t.Error(err)
		}
		if len(out2) > len(out) {
			t.Errorf("expected %d to be less than %d", len(out2), len(out))
		}
	})

	t.Run("supports japanese in text tags", func(t *testing.T) {
		japaneseFrame := &frames.Frame{
			Header: &frames.FrameHeader{ID: "TPE1"},
			Body:   frames.NewTextInformation("日本語"),
		}
		tag := createTag(t, japaneseFrame)
		out, err := tag.MarshalBinary()
		if err != nil {
			t.Error(err)
		}
		nt := NewID3v2()
		if err := nt.UnmarshalBinary(out); err != nil {
			t.Error(err)
		}
		for _, f := range *nt.Frames {
			if f.Header.ID == "TPE1" {
				f.Body.(*frames.TextInformation).Information = []rune("日本語")
			}
		}
	})
}

func createTag(t *testing.T, fs ...*frames.Frame) *ID3v2 {
	tag := NewID3v2()
	tag.Header.FileIdentifier = []byte("ID3")
	tag.Header.MajorVersion = 3
	tag.Frames.ApplyFrame(&frames.Frame{
		Header: &frames.FrameHeader{ID: "TIT2"},
		Body:   &frames.TextInformation{Information: []rune("test2")},
	})
	tag.Frames.ApplyFrame(&frames.Frame{
		Header: &frames.FrameHeader{ID: "TIT1"},
		Body:   &frames.TextInformation{Information: []rune("test1")},
	})
	tag.Frames.ApplyFrame(&frames.Frame{
		Header: &frames.FrameHeader{ID: "TIT3"},
		Body:   &frames.TextInformation{Information: []rune("test3")},
	})
	for _, f := range fs {
		if err := tag.Frames.ApplyFrame(f); err != nil {
			t.Error(f)
		}
	}
	return tag
}
