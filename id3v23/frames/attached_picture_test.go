package frames

import "testing"

func TestAttachedPictureEncoding(t *testing.T) {
	t.Run("marshal is inverse of unmarshal", func(t *testing.T) {
		testcases := []struct {
			name  string
			input *AttachedPicture
		}{
			{
				name: "valid attached picture",
				input: &AttachedPicture{
					TextEncoding: 0,
					MIMEType:     "image/png",
					PictureType:  0,
					Description:  "description",
					PictureData:  []byte("picture data"),
				},
			},
		}

		for _, tt := range testcases {
			t.Run(tt.name, func(t *testing.T) {
				b, err := tt.input.MarshalBinary()
				if err != nil {
					t.Fatal(err)
				}
				ap := &AttachedPicture{}
				if err := ap.UnmarshalBinary(b); err != nil {
					t.Fatal(err)
				}
				if !ap.Equal(tt.input) {
					t.Fatalf("\nexpected: %v\n     got: %v", tt.input, ap)
				}
			})
		}
	})
}

func TestAttachedPictureUnmarshal(t *testing.T) {
	testcases := []struct {
		name     string
		input    []byte
		expected *AttachedPicture
	}{
		{
			name: "valid attached picture",
			input: []byte{
				0x00,                                                       // Text encoding
				0x69, 0x6d, 0x61, 0x67, 0x65, 0x2F, 0x70, 0x6E, 0x67, 0x00, // MIME type
				0x00,                                                                   // Picture type
				0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x00, // Description
				0x64, 0x61, 0x74, 0x61, // Picture data
			},
			expected: &AttachedPicture{
				TextEncoding: 0,
				MIMEType:     "image/png",
				PictureType:  0,
				Description:  "description",
				PictureData:  []byte("data"),
			},
		},
		{
			name: "weird bug with 0x03 and no description field",
			input: []byte{
				0x00,                                                 // Text encoding
				0x69, 0x6d, 0x61, 0x67, 0x65, 0x2F, 0x70, 0x6e, 0x67, // MIME type
				0x03,                   // Picture type
				0x64, 0x61, 0x74, 0x61, // Picture data
			},
			expected: &AttachedPicture{
				TextEncoding: 0,
				MIMEType:     "image/png",
				PictureType:  3,
				Description:  "",
				PictureData:  []byte("data"),
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			ap := &AttachedPicture{}
			if err := ap.UnmarshalBinary(tt.input); err != nil {
				t.Fatal(err)
			}
			if !ap.Equal(tt.expected) {
				t.Fatalf("\nexpected: %v\n     got: %v", tt.expected, ap)
			}
		})
	}
}
