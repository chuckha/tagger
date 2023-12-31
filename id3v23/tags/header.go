package tags

import (
	"fmt"

	"github.com/chuckha/tagger/id3math"

	"gitlab.com/tozd/go/errors"
)

type NoID3v2IdentifierError struct {
	actual []byte
}

func NewNoID3v2IdentifierError(actual []byte) *NoID3v2IdentifierError {
	return &NoID3v2IdentifierError{actual: actual}
}

func (n *NoID3v2IdentifierError) Error() string {
	return fmt.Sprintf("expected ID3v2 identifier, got %q", n.actual)
}

type Header struct {
	FileIdentifier    []byte
	MajorVersion      byte
	Revision          byte
	Unsynchronisation bool
	ExtendedHeader    bool
	Experimental      bool
	// If the tag changes in anyway, this may be out dated and should be updated upon writing.
	Size int
}

func (h *Header) UnmarshalBinary(data []byte) error {
	// header is always 10 bytes, but there could be an extended header
	if len(data) != 10 {
		return errors.Errorf("expected 10 bytes, got %d", len(data))
	}
	h.FileIdentifier = data[0:3]
	if string(h.FileIdentifier) != "ID3" {
		return errors.WithStack(NewNoID3v2IdentifierError(h.FileIdentifier))
	}
	h.MajorVersion = data[3]
	h.Revision = data[4]
	if h.MajorVersion != 3 {
		return errors.Errorf("this program only supports v2.3.0; this file is v2.%d.%d", h.MajorVersion, h.Revision)
	}
	h.Unsynchronisation = data[5]&FlagUnsynchronisation == FlagUnsynchronisation
	h.ExtendedHeader = data[5]&FlagExtendedHeader == FlagExtendedHeader
	h.Experimental = data[5]&FlagExperimental == FlagExperimental
	h.Size = id3math.SyncSafeToInt(data[6:10])
	return nil
}

func (h *Header) MarshalBinary() ([]byte, error) {
	flags := 0
	if h.Unsynchronisation {
		flags = flags | FlagUnsynchronisation
	}
	if h.ExtendedHeader {
		flags = flags | FlagExtendedHeader
	}
	if h.Experimental {
		flags = flags | FlagExperimental
	}
	size := id3math.IntToSyncSafe(h.Size)
	return []byte{'I', 'D', '3', h.MajorVersion, h.Revision, byte(flags), size[0], size[1], size[2], size[3]}, nil
}

const (
	FlagUnsynchronisation = 0b10000000
	FlagExtendedHeader    = 0b01000000
	FlagExperimental      = 0b00100000
)

func (i Header) String() string {
	return fmt.Sprintf("%sv2.%d.%d; size: %d bytes", i.FileIdentifier, i.MajorVersion, i.Revision, i.Size)
}

func NewHeader(data []byte) Header {
	return Header{
		FileIdentifier:    data[0:3],
		MajorVersion:      data[3],
		Revision:          data[4],
		Unsynchronisation: data[5]&FlagUnsynchronisation == FlagUnsynchronisation,
		ExtendedHeader:    data[5]&FlagExtendedHeader == FlagExtendedHeader,
		Experimental:      data[5]&FlagExperimental == FlagExperimental,
		Size:              id3math.SyncSafeToInt(data[6:10]),
	}
}

func (h *Header) Equal(h2 *Header) bool {
	return h.FileIdentifier[0] == h2.FileIdentifier[0] &&
		h.FileIdentifier[1] == h2.FileIdentifier[1] &&
		h.FileIdentifier[2] == h2.FileIdentifier[2] &&
		h.MajorVersion == h2.MajorVersion &&
		h.Revision == h2.Revision &&
		h.Unsynchronisation == h2.Unsynchronisation &&
		h.ExtendedHeader == h2.ExtendedHeader &&
		h.Experimental == h2.Experimental &&
		h.Size == h2.Size
}
