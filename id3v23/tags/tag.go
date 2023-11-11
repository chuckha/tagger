package tags

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"tagger"
	"tagger/id3v23/frames"
	"text/tabwriter"

	"gitlab.com/tozd/go/errors"
)

const (
	AdditionalPaddingSize = 2048
)

type ID3v2 struct {
	Header  *Header
	Frames  *frames.Frames
	Padding []byte

	// f is file pointer to the mp3. This should be opened for writing.
	f           *os.File
	fullRewrite bool
}

func NewID3v2() *ID3v2 {
	return &ID3v2{
		Header:  &Header{},
		Frames:  &frames.Frames{},
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

// Output returns what to write in the file starting at byte 0 (always implied).
func (i *ID3v2) Output() ([]byte, error) {
	b, err := i.MarshalBinary()
	if err != nil {
		return nil, err
	}
	if !i.fullRewrite {
		return b, nil
	}
	mp3Bytes, err := io.ReadAll(i.f)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return append(b, mp3Bytes...), nil
}

func (i *ID3v2) Close() error {
	return i.f.Close()
}

func (i *ID3v2) MarshalBinary() ([]byte, error) {
	frames := []byte{}
	for _, frame := range *i.Frames {
		frameBytes, err := frame.MarshalBinary()
		if err != nil {
			return nil, err
		}
		frames = append(frames, frameBytes...)
	}

	var padding []byte
	switch len(frames) <= i.Header.Size {
	case true:
		// if it shrunk a LOT, reduce the padding and do a whole rewrite
		if i.Header.Size-len(frames) > AdditionalPaddingSize {
			i.fullRewrite = true
			i.Header.Size = len(frames) + AdditionalPaddingSize
			padding = make([]byte, AdditionalPaddingSize)
		} else {
			// Pad the remaining space in the tag.
			padding = make([]byte, i.Header.Size-len(frames))
		}
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

func (i *ID3v2) ApplyConfig(cfg *tagger.Config) {
	for id, fb := range cfg.Frames {
		i.Frames.ApplyFrame(frames.NewFrame(id, fb))
	}
	i.Frames.RemoveFramesWithID("APIC")
	sort.Sort(i.Frames)
}

func (i *ID3v2) String() string {
	var s strings.Builder
	w := tabwriter.NewWriter(&s, 0, 0, 1, '.', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintf(w, "%s\t%s\n", "Header", i.Header.String())
	for _, frame := range *i.Frames {
		fmt.Fprintf(w, "%s:\t%v\n", frame.Header, frame.Body.String())
	}
	w.Flush()
	return s.String()
}
