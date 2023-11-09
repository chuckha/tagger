package frames

import "testing"

func TestCommentsEncoding(t *testing.T) {
	t.Run("marshal is inverse of unmarshal", func(t *testing.T) {
		testcases := []struct {
			name  string
			input *Comment
		}{
			{
				name: "valid comment",
				input: &Comment{
					TextEncoding:            0,
					Language:                "eng",
					ShortContentDescription: "short",
					ActualText:              "actual",
				},
			},
		}

		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				b, err := tt.input.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				c := &Comment{}
				if err := c.UnmarshalBinary(b); err != nil {
					t.Fatal(err)
				}
				if !c.Equal(tt.input) {
					t.Fatalf("\nexpected: %v\n     got: %v", tt.input, c)
				}
			})
		}
	})
}
