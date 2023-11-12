package frames

import "testing"

func TestUserDefinedURLEncoding(t *testing.T) {
	t.Run("marshal is inverse of unmarshal", func(t *testing.T) {
		testcases := []struct {
			name  string
			input *UserDefinedURL
		}{
			{
				name: "valid user defined url",
				input: &UserDefinedURL{
					TextEncoding: 0,
					Description:  []rune("description"),
					URL:          "url",
				},
			},
		}

		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				b, err := tt.input.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				udurl := &UserDefinedURL{}
				if err := udurl.UnmarshalBinary(b); err != nil {
					t.Fatal(err)
				}
				if !udurl.Equal(tt.input) {
					t.Fatalf("\nexpected: %v\n     got: %v", tt.input, udurl)
				}
			})
		}
	})
}
