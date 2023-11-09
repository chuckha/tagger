package frames

import "testing"

func TestPrivateDataEncoding(t *testing.T) {
	t.Run("marshal is inverse of unmarshal", func(t *testing.T) {
		testcases := []struct {
			name  string
			input *PrivateData
		}{
			{
				name: "valid private data",
				input: &PrivateData{
					OwnerIdentifier: "owner",
					Data:            "data",
				},
			},
		}

		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				b, err := tt.input.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				p := &PrivateData{}
				if err := p.UnmarshalBinary(b); err != nil {
					t.Fatal(err)
				}
				if !p.Equal(tt.input) {
					t.Fatalf("\nexpected: %v\n     got: %v", tt.input, p)
				}
			})
		}
	})
}
