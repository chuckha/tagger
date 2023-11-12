package frames

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/chuckha/tagger/id3string"

	"gitlab.com/tozd/go/errors"
)

// AttachedPicture are all of the attached picture frames.
// AttachedPicture frames have an ID of APIC.
type AttachedPicture struct {
	TextEncoding byte
	MIMEType     string
	PictureType  byte
	Description  []rune
	PictureData  []byte
}

// weird bug, we get image/jpeg0x03ffd8 (03 is the picture type then ffd8 starts the JFIF)
// This means that the mp3 is missing a 0x00 after the MIME type and the description is omitted entirely.

func (a *AttachedPicture) UnmarshalBinary(data []byte) error {
	a.TextEncoding = data[0]
	ptr := 1

	a.MIMEType = id3string.ExtractNullTerminatedASCII(data[ptr:])
	// if the mime type is too long, re-parse it and assume the MIME-type is not 0 terminated and the description is missing
	// not sure why the MP3s i tested had this malformatting...
	if len(a.MIMEType) > 10 {
		// find the first byte less than the lowest ascii value
		for i := 0; i < len(a.MIMEType); i++ {
			if a.MIMEType[i] < 0x20 {
				a.MIMEType = string(data[ptr : ptr+i])
				ptr += i
				break
			}
		}
		a.PictureType = data[ptr]
		ptr++
		a.Description = []rune{}
	} else {
		// otherwise we have a normal layout
		ptr += len(a.MIMEType) + 1
		a.PictureType = data[ptr]
		ptr++
		desc, n := id3string.ExtractNullTerminatedValueWithEncoding(a.TextEncoding, data[ptr:])
		a.Description = desc
		ptr += len(a.Description) + n
	}
	a.PictureData = data[ptr:]
	return nil
}

func (a *AttachedPicture) UnmarshalJSON(data []byte) error {
	return errors.New("not implemented for APIC")
	// but actually here read in a file reference and use that as the picture data
	if err := json.Unmarshal(data, a); err != nil {
		return errors.WithStack(err)
	}
	if !id3string.IsASCII(a.Description) {
		a.TextEncoding = 1
	}
	return nil
}

func (a *AttachedPicture) ExtractPicture() {
	f, err := os.Create("test.jpg")
	if err != nil {
		panic(err)
	}
	if _, err := f.Write(a.PictureData); err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
}

func (a *AttachedPicture) String() string {
	return fmt.Sprintf("enc: %x; mime: %q; type: %q; desc: %q", a.TextEncoding, a.MIMEType, PictureTypes[a.PictureType], a.Description)
}

func (a *AttachedPicture) MarshalBinary() ([]byte, error) {
	out := []byte{a.TextEncoding}
	out = append(out, id3string.EncodeASCIIWithNullTerminator(a.MIMEType)...)
	out = append(out, a.PictureType)
	out = append(out, id3string.EncodeRunesWithNullTerminator(a.TextEncoding, a.Description)...)
	out = append(out, a.PictureData...)
	return out, nil
}

func (a *AttachedPicture) Equal(a2 *AttachedPicture) bool {
	return a.TextEncoding == a2.TextEncoding &&
		a.MIMEType == a2.MIMEType &&
		a.PictureType == a2.PictureType &&
		id3string.Equal(a.Description, a2.Description) &&
		id3string.EqualBytes(a.PictureData, a2.PictureData)
}

var PictureTypes = map[byte]string{
	0x00: "Other",
	0x01: "32x32 pixels 'file icon' (PNG only)",
	0x02: "Other file icon",
	0x03: "Cover (front)",
	0x04: "Cover (back)",
	0x05: "Leaflet page",
	0x06: "Media (e.g. lable side of CD)",
	0x07: "Lead artist/lead performer/soloist",
	0x08: "Artist/performer",
	0x09: "Conductor",
	0x0A: "Band/Orchestra",
	0x0B: "Composer",
	0x0C: "Lyricist/text writer",
	0x0D: "Recording Location",
	0x0E: "During recording",
	0x0F: "During performance",
	0x10: "Movie/video screen capture",
	0x11: "A bright coloured fish",
	0x12: "Illustration",
	0x13: "Band/artist logotype",
	0x14: "Publisher/Studio logotype",
}
