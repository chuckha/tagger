package frames

import (
	"encoding/json"
	"fmt"

	"github.com/chuckha/tagger/id3string"

	"gitlab.com/tozd/go/errors"
)

// Private Data frames have an ID of PRIV
type PrivateData struct {
	OwnerIdentifier string
	Data            []byte
}

func (p *PrivateData) UnmarshalBinary(data []byte) error {
	p.OwnerIdentifier = id3string.ExtractNullTerminatedASCII(data)
	p.Data = data[len(p.OwnerIdentifier)+1:]
	return nil
}

func (b *PrivateData) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, b); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (p *PrivateData) String() string {
	return fmt.Sprintf("ownerid: %q; data: %q", p.OwnerIdentifier, p.Data)
}

func (p *PrivateData) MarshalBinary() ([]byte, error) {
	return append(append([]byte(p.OwnerIdentifier), '\x00'), []byte(p.Data)...), nil
}

func (p *PrivateData) Equal(p2 *PrivateData) bool {
	return p.OwnerIdentifier == p2.OwnerIdentifier &&
		id3string.EqualBytes(p.Data, p2.Data)
}
