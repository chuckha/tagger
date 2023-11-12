package frames

import "testing"

func TestGeneralEncapsulationObjectEncoding(t *testing.T) {
	t.Run("marshal is inverse of unmarshal", func(t *testing.T) {
		testcases := []struct {
			name  string
			input *GeneralEncapsulationObject
		}{
			{
				name: "valid comment",
				input: &GeneralEncapsulationObject{
					TextEncoding:       0,
					MIMEType:           "text/plain",
					Filename:           []rune("filename.txt"),
					ContentDescription: []rune("content description"),
					EncapsulatedObject: []byte("encapsulated object"),
				},
			},
		}

		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				b, err := tt.input.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				g := &GeneralEncapsulationObject{}
				if err := g.UnmarshalBinary(b); err != nil {
					t.Fatal(err)
				}
				if !g.Equal(tt.input) {
					t.Fatalf("\nexpected: %v\n     got: %v", tt.input, g)
				}
			})
		}
	})
}
