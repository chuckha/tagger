package tags

import (
	"fmt"
	"io"
	"os"
	"strings"
	"tagger/id3v23/frames"
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

	// f is file pointer to the mp3. This should be opened for writing.
	f           *os.File
	updated     bool
	fullRewrite bool
}

func NewID3v2() *ID3v2 {
	return &ID3v2{
		Header:  &Header{},
		Frames:  make(frames.Frames, 0),
		Padding: make([]byte, 0),
	}
}

func NewID3v2FromFile(file string) (*ID3v2, error) {
	f, err := os.OpenFile(file, os.O_RDWR, 0644)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	tag := NewID3v2()
	tag.f = f
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

func (i *ID3v2) Close() error {
	defer i.f.Close()
	if !i.updated {
		return nil
	}
	b, err := i.MarshalBinary()
	if err != nil {
		return err
	}
	if !i.fullRewrite {
		if _, err := i.f.Seek(0, 0); err != nil {
			return errors.WithStack(err)
		}
		_, err := i.f.Write(b)
		return errors.WithStack(err)
	}
	// read in the rest of the file
	mp3Bytes, err := io.ReadAll(i.f)
	if err != nil {
		return errors.WithStack(err)
	}
	// write the tag
	if _, err := i.f.Seek(0, 0); err != nil {
		return errors.WithStack(err)
	}
	if _, err := i.f.Write(b); err != nil {
		return errors.WithStack(err)
	}
	// write the rest of the file
	if _, err := i.f.Write(mp3Bytes); err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(i.f.Close())
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

	var padding []byte
	switch len(frames) <= i.Header.Size {
	case true:
		// Pad the remaining space in the tag.
		padding = make([]byte, i.Header.Size-len(frames))
	case false:
		// Tag size has grown too much. requires an entire file rewrite. Add additional padding.
		i.fullRewrite = true
		i.Header.Size = len(frames) + AdditionalPaddingSize
		padding = make([]byte, AdditionalPaddingSize)
	}
	frames = append(frames, padding...)

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
	i.updated = true
	i.Frames.SetTextInformationFrame(id, value)
}

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
