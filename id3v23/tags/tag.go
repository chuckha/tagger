package tags

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/chuckha/tagger"
	"github.com/chuckha/tagger/id3v23/frames"

	"gitlab.com/tozd/go/errors"
)

const (
	MinimalPaddingSize = 1024
)

type ID3v2 struct {
	Header *Header
	Frames *frames.Frames
}

func NewID3v2() *ID3v2 {
	return &ID3v2{
		Header: &Header{
			MajorVersion: 3,
		},
		Frames: &frames.Frames{},
	}
}

// NewID3v2FromFile reads in and unmarshals the entire ID3v2.3 tag.
func NewID3v2FromFile(file string) (*ID3v2, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer f.Close()
	tag := NewID3v2()
	headerBytes := make([]byte, 10)
	if _, err := f.Read(headerBytes); err != nil {
		return nil, errors.WithStack(err)
	}
	if err := tag.Header.UnmarshalBinary(headerBytes); err != nil {
		return nil, err
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

// MarshalBinaryv2 will only marshal the id3v2 tag to binary.
func (i *ID3v2) MarshalBinary() ([]byte, error) {
	// marshal header
	header, err := i.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}
	// TODO: marshal the extra header?
	frames := []byte{}
	for _, frame := range *i.Frames {
		frameBytes, err := frame.MarshalBinary()
		if err != nil {
			return nil, err
		}
		frames = append(frames, frameBytes...)
	}
	tag := append(header, frames...)

	// if the tag fits in the original tag, and the padding is within a threshold, just use the original tag size
	// if the tag fits in the original tag, but there is too much padding; shrink the header
	// if the frames don't fit in the header, make a new header with minimal padding

	if len(tag)+MinimalPaddingSize > i.Header.Size {
		out := make([]byte, len(tag)+MinimalPaddingSize)
		copy(out, tag)
		return out, nil
	}

	// If there's only a little bit of padding left, just keep the same header size
	if i.Header.Size-len(tag) < 3*MinimalPaddingSize {
		out := make([]byte, i.Header.Size)
		copy(out, tag)
		return out, nil
	}

	// otherwise shrink the header
	out := make([]byte, len(tag)+MinimalPaddingSize)
	copy(out, tag)
	return out, nil
}

func (i *ID3v2) ApplyConfig(cfg *tagger.Config) {
	for id, fb := range cfg.Frames {
		i.Frames.ApplyFrame(frames.NewFrame(id, fb))
	}
	// Don't actually want to remove APIC...
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

// An mp3 with an ID3v2 tag contains a header and mp3 bytes.
// Sometimes the entire file needs to be read, but sometimes it doesn't
// what are the inputs and outputs?
// A user says "apply this configuration to the id3v2 tag"
// The resulting id3v2 tag could be larger or smaller than the original.
// if the new header fits in the existing padding, simply rewrite the header with more padding
// if the new header doesn't fit in the existing space, the entire file must be rewritten.

// Write writes the tag to the dst file, using the src file as the original.
// In cases where the src and dst are the same, there is an optimization that can happen: just rewrite the header.
func (t *ID3v2) Write(src, dst string) error {
	file := filepath.Base(dst)
	// get the tag as bytes.
	out, err := t.MarshalBinary()
	if err != nil {
		return errors.WithStack(err)
	}

	// Optimization case:
	// TODO: Consider figuring out some interface to move this a layer up
	// if the output file is the same as the input file AND the header fits, just re-write the header.
	if src == dst && len(out) == t.Header.Size {
		f, err := os.OpenFile(dst, os.O_RDWR, 0644)
		if err != nil {
			return errors.WithStack(err)
		}
		defer f.Close()
		if _, err := f.Write(out); err != nil {
			return errors.WithStack(err)
		}
		return nil
	}

	// normal case, write to a tmp file and then copy it back
	tmpf, err := os.CreateTemp("", file)
	if err != nil {
		return errors.WithStack(err)
	}
	defer os.Remove(tmpf.Name())
	// write the tag to the tmp file
	if _, err := tmpf.Write(out); err != nil {
		return errors.WithStack(err)
	}
	ogf, err := os.Open(src)
	if err != nil {
		return errors.WithStack(err)
	}
	if _, err := ogf.Seek(int64(t.Header.Size), 0); err != nil {
		return errors.WithStack(err)
	}
	// copy the rest of the file
	if _, err := io.Copy(tmpf, ogf); err != nil {
		return errors.WithStack(err)
	}
	if err := ogf.Close(); err != nil {
		return errors.WithStack(err)
	}
	// create the output file if it doesn't exist or open it if it does.
	outfile, err := os.OpenFile(dst, os.O_RDWR, 0644)
	if err != nil {
		return errors.WithStack(err)
	}
	defer outfile.Close()
	if _, err := tmpf.Seek(0, 0); err != nil {
		return errors.WithStack(err)
	}
	// copy the completed temp file to the output file
	n, err := io.Copy(outfile, tmpf)
	if err != nil {
		return errors.WithStack(err)
	}
	fmt.Println("wrote", n, "bytes")
	return nil
}
