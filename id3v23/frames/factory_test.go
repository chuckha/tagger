package frames

import (
	"fmt"
	"sort"
	"testing"
)

func TestFrames_FramesSorting(t *testing.T) {
	frames := Frames{
		&Frame{Header: &FrameHeader{ID: "APIC"}},
		&Frame{Header: &FrameHeader{ID: "TIT2"}},
		&Frame{Header: &FrameHeader{ID: "TALB"}},
		&Frame{Header: &FrameHeader{ID: "TPE1"}},
		&Frame{Header: &FrameHeader{ID: "TPE2"}},
	}

	sort.Sort(&frames)

	if frames[len(frames)-1].Header.ID != "APIC" {
		for _, h := range frames {
			fmt.Println(h.Header.ID)
		}
		t.Fatal("did not sort correctly")
	}
}
