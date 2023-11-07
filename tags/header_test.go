package tags

import (
	"fmt"
	"tagger/id3math"
	"testing"
)

func TestHeader_UnmarshalBinary(t *testing.T) {
	t.Run("valid headers", func(t *testing.T) {
		tests := []struct {
			name     string
			data     []byte
			expected *Header
		}{
			{
				name: "valid header",
				data: []byte{'I', 'D', '3', 4, 0, 0, 0, 0, 0, 0},
				expected: &Header{
					FileIdentifier:    []byte{'I', 'D', '3'},
					MajorVersion:      4,
					Revision:          0,
					Unsynchronisation: false,
					ExtendedHeader:    false,
					Experimental:      false,
					Size:              0,
				},
			},
			{
				name: "valid header with flags",
				data: []byte{'I', 'D', '3', 4, 0, 0b10100000, 0, 0, 0, 0},
				expected: &Header{
					FileIdentifier:    []byte{'I', 'D', '3'},
					MajorVersion:      4,
					Revision:          0,
					Unsynchronisation: true,
					ExtendedHeader:    false,
					Experimental:      true,
					Size:              0,
				},
			},
			{
				name: "valid header with flags and size",
				data: []byte{'I', 'D', '3', 4, 0, 0b11100000, 0, 0, 0, 1},
				expected: &Header{
					FileIdentifier:    []byte{'I', 'D', '3'},
					MajorVersion:      4,
					Revision:          0,
					Unsynchronisation: true,
					ExtendedHeader:    true,
					Experimental:      true,
					Size:              1,
				},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				h := &Header{}
				if err := h.UnmarshalBinary(tt.data); err != nil {
					t.Errorf("did not expect an error")
				}
				if !h.Equal(tt.expected) {
					t.Log("flags", h.Unsynchronisation, h.ExtendedHeader, h.Experimental)
					t.Log("flags", tt.expected.Unsynchronisation, tt.expected.ExtendedHeader, tt.expected.Experimental)
				}
			})
		}
	})
}

func TestHeader_MarshalBinary(t *testing.T) {
	t.Run("valid headers", func(t *testing.T) {
		tests := []struct {
			name     string
			input    *Header
			expected []byte
		}{
			{
				name: "valid header",
				input: &Header{
					FileIdentifier:    []byte{'I', 'D', '3'},
					MajorVersion:      4,
					Revision:          0,
					Unsynchronisation: false,
					ExtendedHeader:    false,
					Experimental:      false,
					Size:              5030,
				},
				expected: []byte{'I', 'D', '3', 4, 0, 0, 0, 0, 39, 166},
			},
			{
				name: "valid header with flags",
				input: &Header{
					FileIdentifier:    []byte{'I', 'D', '3'},
					MajorVersion:      4,
					Revision:          0,
					Unsynchronisation: false,
					ExtendedHeader:    false,
					Experimental:      false,
					Size:              5030,
				},
				expected: []byte{'I', 'D', '3', 4, 0, 0, 0, 0, 39, 166},
			},
			{
				name: "valid header with flags and size",
				input: &Header{
					FileIdentifier:    []byte{'I', 'D', '3'},
					MajorVersion:      4,
					Revision:          0,
					Unsynchronisation: false,
					ExtendedHeader:    true,
					Experimental:      false,
					Size:              5030,
				},
				expected: []byte{'I', 'D', '3', 4, 0, 64, 0, 0, 39, 166},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				b, err := tt.input.MarshalBinary()
				if err != nil {
					t.Errorf("did not expect an error")
				}
				if !equal(b, tt.expected) {
					fmt.Println(b)
					fmt.Println(tt.expected)
					t.Log("flags", tt.input.Unsynchronisation, tt.input.ExtendedHeader, tt.input.Experimental)
					t.Log("flags", tt.expected[3]&FlagUnsynchronisation, tt.expected[3]&FlagExtendedHeader, tt.expected[3]&FlagExperimental)
					t.Log("sizes", tt.input.Size, id3math.SyncSafeToInt(tt.expected[6:10]))
					t.Fatal("expected equal")
				}
			})
		}
	})
}

func equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if b[i] != v {
			return false
		}
	}
	return true
}
