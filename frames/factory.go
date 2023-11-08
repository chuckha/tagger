package frames

import (
	"fmt"
	"strings"
)

var ErrPadding = fmt.Errorf("padding")

const HeaderMinSize = 10

type FrameBody interface {
	String() string
	MarshalBinary() ([]byte, error)
}

type Frames []*Frame

func (f *Frames) SetTextInformationFrame(id, info string) {
	found := false
	for j := 0; j < len(*f); {
		if string((*f)[j].Header.ID) != id {
			j++
			continue
		}
		if !found {
			(*f)[j].Body = &TextInformation{Information: info}
			found = true
			j++
			continue
		}
		*f = append((*f)[:j], (*f)[j+1:]...)
	}

	if !found {
		*f = append(*f, &Frame{
			Header: &FrameHeader{
				ID: id,
			},
			Body: &TextInformation{Information: info},
		})
	}
}

func (f *Frames) UnmarshalBinary(data []byte) error {
	ptr := 0
	for ptr < len(data) {
		if data[ptr] == '\x00' {
			return nil
		}
		header := &FrameHeader{}
		if err := header.UnmarshalBinary(data[ptr : ptr+HeaderMinSize]); err != nil {
			return err
		}
		ptr += HeaderMinSize
		// TODO: read in extended header
		if ptr+header.Size > len(data) {
			fmt.Println("header", header)
		}
		frame := &Frame{Header: header}
		if strings.Contains(header.ID, "\x00") {
			fmt.Println("warning: frame id contains null byte; ignoring the rest of the tag")
			return nil
		}
		if err := frame.UnmarshalBinary(data[ptr : ptr+header.Size]); err != nil {
			return err
		}
		ptr += header.Size
		*f = append(*f, frame)
	}
	return nil
}

type Frame struct {
	Header *FrameHeader
	Body   FrameBody
}

func (f *Frame) UnmarshalBinary(data []byte) error {
	switch string(f.Header.ID) {
	case "TCMP", "TDRL", "TDRC": // non-standard T*** frames
		fallthrough
	case "TALB", "TBPM", "TCOM", "TCON", "TCOP", "TDAT", "TDLY", "TENC", "TEXT", "TFLT",
		"TIME", "TIT1", "TIT2", "TIT3", "TKEY", "TLAN", "TLEN", "TMED", "TOAL", "TOFN",
		"TOLY", "TOPE", "TORY", "TOWN", "TPE1", "TPE2", "TPE3", "TPE4", "TPOS", "TPUB",
		"TRCK", "TRDA", "TRSN", "TRSO", "TSIZ", "TSRC", "TSSE", "TYER":
		ti := &TextInformation{}
		if err := ti.UnmarshalBinary(data); err != nil {
			return err
		}
		f.Body = ti
	case "COMM":
		c := &Comment{}
		if err := c.UnmarshalBinary(data); err != nil {
			return err
		}
		f.Body = c
	case "APIC":
		ap := &AttachedPicture{}
		if err := ap.UnmarshalBinary(data); err != nil {
			return err
		}
		f.Body = ap
	case "WXXX":
		w := &UserDefinedURL{}
		if err := w.UnmarshalBinary(data); err != nil {
			return err
		}
		f.Body = w
	case "PRIV":
		p := &PrivateData{}
		if err := p.UnmarshalBinary(data); err != nil {
			return err
		}
		f.Body = p
	case "USLT":
		u := &UnsynchronizedLyrics{}
		if err := u.UnmarshalBinary(data); err != nil {
			return err
		}
		f.Body = u
	case "TXXX":
		u := &UserDefinedTextInformation{}
		if err := u.UnmarshalBinary(data); err != nil {
			return err
		}
		f.Body = u
	case "MCDI":
		m := &MusicCDIdentifier{}
		if err := m.UnmarshalBinary(data); err != nil {
			return err
		}
		f.Body = m
	case "GEOB":
		g := &GeneralEncapsulationObject{}
		if err := g.UnmarshalBinary(data); err != nil {
			return err
		}
		f.Body = g
	case "USER":
		u := &TermsOfUse{}
		if err := u.UnmarshalBinary(data); err != nil {
			return err
		}
		f.Body = u
	default:
		panic(fmt.Sprintf("unknown frame type: %q", f.Header.ID))
	}
	return nil
}

func (f *Frame) MarshalBinary() ([]byte, error) {
	// marshal the body
	fb, err := f.Body.MarshalBinary()
	if err != nil {
		return nil, err
	}
	// update the header size
	f.Header.Size = len(fb)
	// marshal the header
	fh, err := f.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}
	// concatenate the header and body
	return append(fh, fb...), nil
}
