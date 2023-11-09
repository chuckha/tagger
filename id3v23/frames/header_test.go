package frames

import (
	"bytes"
	"reflect"
	"testing"
)

func TestFrameHeader_MarshalBinary(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		testcases := []struct {
			name     string
			input    *FrameHeader
			expected []byte
		}{
			{
				name: "valid header",
				input: &FrameHeader{
					ID:                       "TIT2",
					Size:                     1000,
					PreserveTagOnAlteration:  true,
					PreserveFileOnAlteration: true,
					ReadOnly:                 true,
					Compressed:               true,
					Encrypted:                true,
					ContainsGroupingIdentity: true,
				},
				expected: []byte{84, 73, 84, 50, 0, 0, 3, 232, 224, 28},
			},
		}

		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				b, err := tt.input.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(b, tt.expected) {
					t.Fatalf("\nexpected: %v\n     got: %v", tt.expected, b)
				}
			})
		}
	})

	t.Run("marshal is inverse of unmarshal", func(t *testing.T) {
		testcases := []struct {
			name  string
			input *FrameHeader
		}{
			{
				name: "valid header",
				input: &FrameHeader{
					ID:                       "TIT2",
					Size:                     1030,
					PreserveTagOnAlteration:  true,
					PreserveFileOnAlteration: true,
					ReadOnly:                 true,
					Compressed:               true,
					Encrypted:                true,
					ContainsGroupingIdentity: true,
				},
			},
		}

		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				b, err := tt.input.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}

				fh := &FrameHeader{}
				err = fh.UnmarshalBinary(b)
				if err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(tt.input, fh) {
					t.Fatalf("not equal\na: %v\nb: %v", tt.input, fh)
				}
			})
		}
	})
}
