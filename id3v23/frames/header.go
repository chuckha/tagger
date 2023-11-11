package frames

import (
	"fmt"

	"github.com/chuckha/tagger/id3math"
)

const (
	FlagTagAlterPreservation  = 0b10000000
	FlagFileAlterPreservation = 0b01000000
	FlagReadOnly              = 0b00100000
	FlagCompression           = 0b00010000
	FlagEncryption            = 0b00001000
	FlagGroupingIdentity      = 0b00000100
)

type FrameHeader struct {
	ID   string
	Size int

	// PreserveTagOnAlteration: this flag tells the software what to do with this frame if it is
	//   unknown and the tag is altered in any way. This applies to all
	//   kinds of alterations, including adding more padding and reordering
	//   the frames.
	PreserveTagOnAlteration bool

	// PreserveFileOnAlteration: this flag tells the software what to do with this frame if it is
	//   unknown and the file, excluding the tag, is altered. This does not
	//   apply when the audio is completely replaced with other audio data.
	PreserveFileOnAlteration bool

	// ReadOnly: this flag, if set, tells the software that the contents of this
	//   frame is intended to be read only. Changing the contents might
	//   break something, e.g. a signature. If the contents are changed,
	//   without knowledge in why the frame was flagged read only and
	//   without taking the proper means to compensate, e.g. recalculating
	//   the signature, the bit should be cleared.
	ReadOnly bool

	// Compression: this flag indicates whether or not the frame is compressed.
	Compressed bool

	// Encryption: this flag indicates wether or not the frame is encrypted. If set
	// one byte indicating with which method it was encrypted will be
	// appended to the frame header. See section 4.26. for more
	// information about encryption method registration.
	Encrypted bool

	// ContainsGroupingIdentity: this flag indicates whether or not this frame belongs in a group
	// with other frames. If set a group identifier byte is added to the
	// frame header. Every frame with the same group identifier belongs
	// to the same group.
	ContainsGroupingIdentity bool
}

func (f *FrameHeader) UnmarshalBinary(data []byte) error {
	// TODO: this is not true, it could have a longer than 10 byte header if the flags are set
	if len(data) != 10 {
		panic("expected a 10 byte header")
	}

	f.ID = string(data[0:4])
	f.Size = id3math.BytesToInt(data[4:8])
	f.PreserveTagOnAlteration = data[8]&FlagTagAlterPreservation == FlagTagAlterPreservation
	f.PreserveFileOnAlteration = data[8]&FlagFileAlterPreservation == FlagFileAlterPreservation
	f.ReadOnly = data[8]&FlagReadOnly == FlagReadOnly
	f.Compressed = data[9]&FlagCompression == FlagCompression
	f.Encrypted = data[9]&FlagEncryption == FlagEncryption
	f.ContainsGroupingIdentity = data[9]&FlagGroupingIdentity == FlagGroupingIdentity
	return nil
}

func (f *FrameHeader) String() string {
	flags := f.FlagsAsBytes()
	return fmt.Sprintf("%s (%s) size: %d, flags: %08b %08b", f.ID, Descriptions[string(f.ID)], f.Size, flags[0], flags[1])
}

func (f *FrameHeader) MarshalBinary() ([]byte, error) {
	size := id3math.IntToBytes(f.Size)
	flags := f.FlagsAsBytes()
	return append([]byte(f.ID), size[0], size[1], size[2], size[3], flags[0], flags[1]), nil
}

func (f *FrameHeader) FlagsAsBytes() [2]byte {
	var flags [2]byte
	if f.PreserveTagOnAlteration {
		flags[0] |= FlagTagAlterPreservation
	}
	if f.PreserveFileOnAlteration {
		flags[0] |= FlagFileAlterPreservation
	}
	if f.ReadOnly {
		flags[0] |= FlagReadOnly
	}
	if f.Compressed {
		flags[1] |= FlagCompression
	}
	if f.Encrypted {
		flags[1] |= FlagEncryption
	}
	if f.ContainsGroupingIdentity {
		flags[1] |= FlagGroupingIdentity
	}
	return flags

}
