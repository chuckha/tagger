package tags

import (
	"strings"
	"testing"

	"github.com/chuckha/tagger"
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

		err = tag.ApplyConfig(&tagger.Config{
			Frames: map[string]frames.FrameBody{
				"TIT1": &frames.TextInformation{Information: "test99"},
				"TPOS": &frames.TextInformation{Information: "1/1"},
				"TPE1": &frames.TextInformation{Information: strings.Repeat("a", 10000)},
			},
		})
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

		err = tag.ApplyConfig(&tagger.Config{
			Frames: map[string]frames.FrameBody{
				"TIT1": &frames.TextInformation{Information: "t"},
			},
		})
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
			Body:   &frames.TextInformation{Information: strings.Repeat("a", 10000)},
		}
		tag := createTag(t, bigFrame)
		out, err := tag.MarshalBinary()
		if err != nil {
			t.Error(err)
		}
		tag.Header.Size = len(out) - 10

		err = tag.ApplyConfig(&tagger.Config{
			Frames: map[string]frames.FrameBody{
				"TIT1": &frames.TextInformation{Information: "t"},
				"TPE1": &frames.TextInformation{Information: "smaller"},
			},
		})
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
}

func createTag(t *testing.T, fs ...*frames.Frame) *ID3v2 {
	tag := NewID3v2()
	tag.Header.FileIdentifier = []byte("ID3")
	tag.Header.MajorVersion = 3
	tag.Frames.ApplyFrame(&frames.Frame{
		Header: &frames.FrameHeader{ID: "TIT2"},
		Body:   &frames.TextInformation{Information: "test2"},
	})
	tag.Frames.ApplyFrame(&frames.Frame{
		Header: &frames.FrameHeader{ID: "TIT1"},
		Body:   &frames.TextInformation{Information: "test1"},
	})
	tag.Frames.ApplyFrame(&frames.Frame{
		Header: &frames.FrameHeader{ID: "TIT3"},
		Body:   &frames.TextInformation{Information: "test3"},
	})
	for _, f := range fs {
		if err := tag.Frames.ApplyFrame(f); err != nil {
			t.Error(f)
		}
	}
	return tag
}
