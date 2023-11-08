package tags

import (
	"fmt"
	"os"
	"strings"
	"tagger/frames"
	"text/tabwriter"

	"github.com/pkg/errors"
)

const (
	AdditionalPaddingSize = 2048
)

type ID3v2 struct {
	Header  *Header
	Frames  frames.Frames
	Padding []byte
}

func NewID3v2() *ID3v2 {
	return &ID3v2{
		Header:  &Header{},
		Frames:  make(frames.Frames, 0),
		Padding: make([]byte, 0),
	}
}

func NewID3v2FromFile(file string) (*ID3v2, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	tag := NewID3v2()
	headerBytes := make([]byte, 10)
	if _, err := f.Read(headerBytes); err != nil {
		return nil, errors.WithStack(err)
	}
	if string(headerBytes[0:3]) != "ID3" {
		return nil, errors.New("no id3 tag")
	}
	if err := tag.Header.UnmarshalBinary(headerBytes); err != nil {
		return nil, err
	}
	if tag.Header.MajorVersion != 3 {
		return nil, errors.Errorf("this program only supports v2.3.0; this file is v2.%d.%d", tag.Header.MajorVersion, tag.Header.Revision)
	}
	// TODO: there could be an extended header here
	tagBytes := make([]byte, tag.Header.Size)
	if _, err := f.Read(tagBytes); err != nil {
		return nil, errors.WithStack(err)
	}
	if err := tag.Frames.UnmarshalBinary(tagBytes); err != nil {
		return nil, err
	}
	return tag, nil
}

func (i *ID3v2) MarshalBinary() ([]byte, error) {
	frames := []byte{}
	for _, frame := range i.Frames {
		frameBytes, err := frame.MarshalBinary()
		if err != nil {
			return nil, err
		}
		frames = append(frames, frameBytes...)
	}

	// add padding if necessary
	if len(frames) < i.Header.Size {
		padding := make([]byte, i.Header.Size-len(frames))
		frames = append(frames, padding...)
	}

	// the new frames are bigger than the header; add some extra padding while we rewrite the whole file.
	if len(frames) > i.Header.Size {
		i.Header.Size = len(frames) + AdditionalPaddingSize
		padding := make([]byte, AdditionalPaddingSize)
		frames = append(frames, padding...)
	}

	// marshal header
	header, err := i.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}
	header = append(header, frames...)
	return header, nil
}

// SetTextFrame can be used for any T000-TZZZ frame (excluding TXXX).
// This will overwrite any existing frame with the same id.
// If there is no frame with the same id, it will append a new frame.
// This will remove any other frames with the same id.
func (i *ID3v2) SetTextFrame(id, value string) {
	i.Frames.SetTextInformationFrame(id, value)
}

// marshal to bytes
// oh and then overwrite the header of the file!
// make sure it fits in the header space or else the whole file needs to be rewritten

func (i *ID3v2) String() string {
	var s strings.Builder
	w := tabwriter.NewWriter(&s, 0, 0, 1, '.', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintf(w, "%s\t%s\n", "Header", i.Header.String())
	for _, frame := range i.Frames {
		fmt.Fprintf(w, "%s:\t%v\n", frame.Header, frame.Body.String())
	}
	w.Flush()
	return s.String()
}
