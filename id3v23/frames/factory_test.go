package frames

import "testing"

func TestFrames_SetTextInformationFrame(t *testing.T) {
	t.Run("empty frames", func(t *testing.T) {
		frames := make(Frames, 0)
		frames.SetTextInformationFrame("TRCK", "1/10")
		if len(frames) != 1 {
			t.Errorf("expected 1 frame, got %d", len(frames))
		}
		if frames[0].Header.ID != "TRCK" {
			t.Errorf("expected frame id TRCK, got %q", frames[0].Header.ID)
		}
		if frames[0].Body.(*TextInformation).Information != "1/10" {
			t.Errorf("expected frame body 1/10, got %q", frames[0].Body.String())
		}
	})

	t.Run("existing frame", func(t *testing.T) {
		frames := Frames{
			&Frame{
				Header: &FrameHeader{ID: "TRCK"},
				Body:   &TextInformation{Information: "1/10"},
			},
		}
		frames.SetTextInformationFrame("TRCK", "2/10")
		if len(frames) != 1 {
			t.Errorf("expected 1 frame, got %d", len(frames))
		}
		if frames[0].Header.ID != "TRCK" {
			t.Errorf("expected frame id TRCK, got %q", frames[0].Header.ID)
		}
		if frames[0].Body.(*TextInformation).Information != "2/10" {
			t.Errorf("expected frame body 2/10, got %q", frames[0].Body.String())
		}
	})

	t.Run("multiple frames", func(t *testing.T) {
		frames := Frames{
			&Frame{
				Header: &FrameHeader{ID: "TPE1"},
				Body:   &TextInformation{Information: "Artist"},
			},
			&Frame{
				Header: &FrameHeader{ID: "TRCK"},
				Body:   &TextInformation{Information: "1/10"},
			},
		}
		frames.SetTextInformationFrame("TRCK", "2/10")
		if len(frames) != 2 {
			t.Errorf("expected 2 frames, got %d", len(frames))
		}
		if frames[1].Header.ID != "TRCK" {
			t.Errorf("expected frame id TRCK, got %q", frames[0].Header.ID)
		}
		if frames[1].Body.(*TextInformation).Information != "2/10" {
			t.Errorf("expected frame body 2/10, got %q", frames[0].Body.String())
		}
	})
}
