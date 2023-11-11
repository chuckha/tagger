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
	UnmarshalBinary([]byte) error
	UnmarshalJSON([]byte) error
}

type FrameBodyUnpacker struct {
	FrameBody
}

// UnmarshalJSON will try to unpack the incoming bytes into the correct frame type.
func (f *FrameBodyUnpacker) UnmarshalJSON(data []byte) error {
	fbs := []FrameBody{
		&Comment{}, &TextInformation{}, &AttachedPicture{}, &UserDefinedURL{},
		&PrivateData{}, &UserDefinedTextInformation{}, &MusicCDIdentifier{},
		//		&GeneralEncapsulationObject{}, &TermsOfUse{},
	}
	for _, fb := range fbs {
		if err := fb.UnmarshalJSON(data); err == nil {
			f.FrameBody = fb
			return nil
		}
	}
	return fmt.Errorf("unknown frame type: %q", data)
}

type Frames []*Frame

// implement sort on frames
func (f Frames) Len() int { return len(f) }
func (f Frames) Less(i, j int) bool {
	if string(f[i].Header.ID) == "APIC" {
		return false
	}
	if string(f[j].Header.ID) == "APIC" {
		return true
	}
	return string(f[i].Header.ID) < string(f[j].Header.ID)
}
func (f *Frames) Swap(i, j int) { (*f)[i], (*f)[j] = (*f)[j], (*f)[i] }

func (f *Frames) RemoveFramesWithID(id string) {
	for i := 0; i < len(*f); i++ {
		if string((*f)[i].Header.ID) == id {
			*f = append((*f)[:i], (*f)[i+1:]...)
			i--
		}
	}
}

func (f *Frames) ApplyFrame(frame *Frame) error {
	// TODO: add id3v2.3 rules here for how many of which frame can exist
	switch IDToFrameKind[string(frame.Header.ID)] {
	case TextInformationKind, NonStandardTextInformationKind:
		// remove all of the frames with the same id
		for i := 0; i < len(*f); i++ {
			if (*f)[i].Header.ID == frame.Header.ID {
				*f = append((*f)[:i], (*f)[i+1:]...)
				i--
			}
		}
	default:
		return fmt.Errorf("apply frame cannot handle: %q", frame.Header.ID)
	}

	// append it if it's not handled by any above case
	*f = append(*f, frame)
	return nil
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

func NewFrame(id string, body FrameBody) *Frame {
	return &Frame{
		Header: &FrameHeader{ID: id},
		Body:   body,
	}
}

const (
	TextInformationKind            = "text information"
	NonStandardTextInformationKind = "non-standard text information"
	CommentKind                    = "comment"
	AttachedPictureKind            = "attached picture"
	UserDefinedURLKind             = "user defined url"
	PrivateKind                    = "private"
	UnsynchronizedLyricsKind       = "unsynchronized lyrics"
	UserDefinedTextInformationKind = "user defined text information"
	MusicCDIdentifierKind          = "music cd identifier"
	GeneralEncapsulationObjectKind = "general encapsulation object"
	TermsOfUseKind                 = "terms of use"
)

var IDToFrameKind = map[string]string{
	"TALB": TextInformationKind,
	"TBPM": TextInformationKind,
	"TCOM": TextInformationKind,
	"TCON": TextInformationKind,
	"TCOP": TextInformationKind,
	"TDAT": TextInformationKind,
	"TDLY": TextInformationKind,
	"TENC": TextInformationKind,
	"TEXT": TextInformationKind,
	"TFLT": TextInformationKind,
	"TIME": TextInformationKind,
	"TIT1": TextInformationKind,
	"TIT2": TextInformationKind,
	"TIT3": TextInformationKind,
	"TKEY": TextInformationKind,
	"TLAN": TextInformationKind,
	"TLEN": TextInformationKind,
	"TMED": TextInformationKind,
	"TOAL": TextInformationKind,
	"TOFN": TextInformationKind,
	"TOLY": TextInformationKind,
	"TOPE": TextInformationKind,
	"TORY": TextInformationKind,
	"TOWN": TextInformationKind,
	"TPE1": TextInformationKind,
	"TPE2": TextInformationKind,
	"TPE3": TextInformationKind,
	"TPE4": TextInformationKind,
	"TPOS": TextInformationKind,
	"TPUB": TextInformationKind,
	"TRCK": TextInformationKind,
	"TRDA": TextInformationKind,
	"TRSN": TextInformationKind,
	"TRSO": TextInformationKind,
	"TSIZ": TextInformationKind,
	"TSRC": TextInformationKind,
	"TSSE": TextInformationKind,
	"TYER": TextInformationKind,
	"TCMP": NonStandardTextInformationKind,
	"TDRL": NonStandardTextInformationKind,
	"TDRC": NonStandardTextInformationKind,
	"COMM": CommentKind,
	"APIC": AttachedPictureKind,
	"WXXX": UserDefinedURLKind,
	"PRIV": PrivateKind,
	"USLT": UnsynchronizedLyricsKind,
	"TXXX": UserDefinedTextInformationKind,
	"MCDI": MusicCDIdentifierKind,
	"GEOB": GeneralEncapsulationObjectKind,
	"USER": TermsOfUseKind,
}

func (f *Frame) UnmarshalBinary(data []byte) error {
	switch IDToFrameKind[string(f.Header.ID)] {
	case TextInformationKind, NonStandardTextInformationKind:
		f.Body = &TextInformation{}
	case CommentKind:
		f.Body = &Comment{}
	case AttachedPictureKind:
		f.Body = &AttachedPicture{}
	case UserDefinedURLKind:
		f.Body = &UserDefinedURL{}
	case PrivateKind:
		f.Body = &PrivateData{}
	// case UnsynchronizedLyricsKind:
	// 	f.Body = &UnsynchronizedLyrics{}
	case UserDefinedTextInformationKind:
		f.Body = &UserDefinedTextInformation{}
	case MusicCDIdentifierKind:
		f.Body = &MusicCDIdentifier{}
	// case GeneralEncapsulationObjectKind:
	// 	f.Body = &GeneralEncapsulationObject{}
	// case TermsOfUseKind:
	// 	f.Body = &TermsOfUse{}
	default:
		panic(fmt.Sprintf("unknown frame type: %q", f.Header.ID))
	}
	if err := f.Body.UnmarshalBinary(data); err != nil {
		return err
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
