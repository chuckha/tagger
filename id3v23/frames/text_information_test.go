package frames

import (
	"testing"
)

func TestTextInformationEncoding(t *testing.T) {
	t.Run("marshal is inverse of unmarshal", func(t *testing.T) {
		testcases := []struct {
			name  string
			input *TextInformation
		}{
			{
				name: "valid text information",
				input: &TextInformation{
					TextEncoding: 0,
					Information:  "text",
				},
			},
		}

		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				b, err := tt.input.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				ti := &TextInformation{}
				if err := ti.UnmarshalBinary(b); err != nil {
					t.Fatal(err)
				}
				if !ti.Equal(tt.input) {
					t.Fatalf("\nexpected: %v\n     got: %v", tt.input, ti)
				}
			})
		}
	})

	t.Run("UnmarshalJSON and Unicode", func(t *testing.T) {
		testcases := []struct {
			name             string
			input            string
			expectedEncoding byte
		}{
			{
				name:             "unicode text information",
				input:            `{"Information":"日本語"}`,
				expectedEncoding: 1,
			},
			{
				name:             "ascii text information",
				input:            `{"Information":"ascii"}`,
				expectedEncoding: 0,
			},
		}

		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				ti := &TextInformation{}
				if err := ti.UnmarshalJSON([]byte(tt.input)); err != nil {
					t.Fatal(err)
				}
				if ti.TextEncoding != tt.expectedEncoding {
					t.Fatalf("expected text encoding to be %d, got %d", tt.expectedEncoding, ti.TextEncoding)
				}
			})
		}
	})
}
