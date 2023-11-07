package tags

import "tagger/frames"

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

func (i *ID3v2) MarshalBinary() ([]byte, error) {
	out := []byte{}
	header, err := i.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}
	out = append(out, header...)
	for _, frame := range i.Frames {
		frameBytes, err := frame.MarshalBinary()
		if err != nil {
			return nil, err
		}
		out = append(out, frameBytes...)
	}
	// add padding

	// marshal header
	// marshal frames
	// marshal padding
	return out, nil
}

// marshal to bytes
// oh and then overwrite the header of the file!
// make sure it fits in the header space or else the whole file needs to be rewritten
