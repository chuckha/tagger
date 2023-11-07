package frames

import "testing"

func TestUserDefinedTextInformationEncoding(t *testing.T) {
	t.Run("marshal is inverse of unmarshal", func(t *testing.T) {
		testcases := []struct {
			name  string
			input *UserDefinedTextInformation
		}{
			{
				name: "valid user defined text information",
				input: &UserDefinedTextInformation{
					TextEncoding: 0,
					Description:  "text",
					Value:        "value",
				},
			},
		}

		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				b, err := tt.input.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				udti := &UserDefinedTextInformation{}
				if err := udti.UnmarshalBinary(b); err != nil {
					t.Fatal(err)
				}
				if !udti.Equal(tt.input) {
					t.Fatalf("\nexpected: %v\n     got: %v", tt.input, udti)
				}
			})
		}
	})
}
