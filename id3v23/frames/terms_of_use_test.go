package frames

import "testing"

func TestTermsOfUseEncoding(t *testing.T) {
	t.Run("marshal is inverse of unmarshal", func(t *testing.T) {
		testcases := []struct {
			name  string
			input *TermsOfUse
		}{
			{
				name: "valid terms of use",
				input: &TermsOfUse{
					TextEncoding: 0,
					Language:     "eng",
					Text:         "text",
				},
			},
		}

		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				b, err := tt.input.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				tu := &TermsOfUse{}
				if err := tu.UnmarshalBinary(b); err != nil {
					t.Fatal(err)
				}
				if !tu.Equal(tt.input) {
					t.Fatalf("\nexpected: %v\n     got: %v", tt.input, tu)
				}
			})
		}
	})
}
