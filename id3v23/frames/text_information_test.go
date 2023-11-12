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
					Information:  []rune("text"),
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

	t.Run("Marshal does not write null bytes", func(t *testing.T) {
		ti := &TextInformation{
			TextEncoding: 1,
			Information:  []rune("食べる"),
		}
		b, err := ti.MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}
		for _, v := range b {
			if v == 0x00 {
				t.Fatal("null byte found")
			}
		}
	})

	t.Run("Unmarshal test...", func(t *testing.T) {
		ti := &TextInformation{}
		if err := ti.UnmarshalJSON([]byte(`{"Information":"しろくまカフェ"}`)); err != nil {
			t.Fatal(err)
		}
		if ti.TextEncoding != 1 {
			t.Fatalf("expected text encoding to be 1, got %d", ti.TextEncoding)
		}
		if string(ti.Information) != "しろくまカフェ" {
			t.Fatalf("expected information to be しろくまカフェ, got %s", string(ti.Information))
		}
	})
}
