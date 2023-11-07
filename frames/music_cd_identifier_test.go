package frames

import "testing"

func TestMusicCDIdentifierEncoding(t *testing.T) {
	t.Run("marshal is inverse of unmarshal", func(t *testing.T) {
		testcases := []struct {
			name  string
			input *MusicCDIdentifier
		}{
			{
				name: "valid music cd identifier",
				input: &MusicCDIdentifier{
					TableOfContents: []byte{0x01, 0x02, 0x03, 0x04},
				},
			},
		}

		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				b, err := tt.input.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				m := &MusicCDIdentifier{}
				if err := m.UnmarshalBinary(b); err != nil {
					t.Fatal(err)
				}
				if !m.Equal(tt.input) {
					t.Fatalf("\nexpected: %v\n     got: %v", tt.input, m)
				}
			})
		}
	})

}
