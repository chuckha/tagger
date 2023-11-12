package frames

import "testing"

func TestUnsynchronizedLyricsEncoding(t *testing.T) {
	t.Run("marshal is inverse of unmarshal", func(t *testing.T) {
		testcases := []struct {
			name  string
			input *UnsynchronizedLyrics
		}{
			{
				name: "valid unsynchronized lyrics",
				input: &UnsynchronizedLyrics{
					TextEncoding:      0,
					Language:          "eng",
					ContentDescriptor: []rune("content"),
					Lyrics:            "lyrics",
				},
			},
		}

		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				b, err := tt.input.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				ul := &UnsynchronizedLyrics{}
				if err := ul.UnmarshalBinary(b); err != nil {
					t.Fatal(err)
				}
				if !ul.Equal(tt.input) {
					t.Fatalf("\nexpected: %v\n     got: %v", tt.input, ul)
				}
			})
		}
	})
}
