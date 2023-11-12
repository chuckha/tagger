package frames

import (
	"fmt"

	"github.com/chuckha/tagger/id3string"
)

// GeneralEncapsulationObject have the ID GEOB
type GeneralEncapsulationObject struct {
	TextEncoding       byte
	MIMEType           string
	Filename           []rune
	ContentDescription []rune
	EncapsulatedObject []byte
}

func (g *GeneralEncapsulationObject) UnmarshalBinary(data []byte) error {
	g.TextEncoding = data[0]
	ptr := 1
	g.MIMEType = id3string.ExtractNullTerminatedASCII(data[ptr:])
	ptr += len(g.MIMEType) + 1
	filename, n := id3string.ExtractNullTerminatedValueWithEncoding(g.TextEncoding, data[ptr:])
	g.Filename = filename
	ptr += len(g.Filename) + n
	contentDescription, n := id3string.ExtractNullTerminatedValueWithEncoding(g.TextEncoding, data[ptr:])
	g.ContentDescription = contentDescription
	ptr += len(g.ContentDescription) + n
	g.EncapsulatedObject = data[ptr:]
	return nil
}

func (g *GeneralEncapsulationObject) String() string {
	return fmt.Sprintf("mime: %q; filename: %q; contentdesc: %q; encapsulatedobject: %q", g.MIMEType, g.Filename, g.ContentDescription, g.EncapsulatedObject)
}

func (g *GeneralEncapsulationObject) MarshalBinary() ([]byte, error) {
	out := []byte{g.TextEncoding}
	out = append(out, id3string.EncodeASCIIWithNullTerminator(g.MIMEType)...)
	out = append(out, id3string.EncodeRunesWithNullTerminator(g.TextEncoding, g.Filename)...)
	out = append(out, id3string.EncodeRunesWithNullTerminator(g.TextEncoding, g.ContentDescription)...)
	out = append(out, g.EncapsulatedObject...)
	return out, nil
}

func (g *GeneralEncapsulationObject) Equal(g2 *GeneralEncapsulationObject) bool {
	return g.TextEncoding == g2.TextEncoding &&
		g.MIMEType == g2.MIMEType &&
		id3string.Equal(g.Filename, g2.Filename) &&
		id3string.Equal(g.ContentDescription, g2.ContentDescription) &&
		string(g.EncapsulatedObject) == string(g2.EncapsulatedObject)
}
