package frames

import (
	"fmt"
	"tagger/id3string"
)

// PrivateData are all of the text frames.
// Private Data frames have an ID of PRIV
type PrivateData struct {
	OwnerIdentifier string
	Data            string
}

func (p *PrivateData) UnmarshalBinary(data []byte) error {
	p.OwnerIdentifier = id3string.ExtractNullTerminated(data)
	p.Data = string(data[len(p.OwnerIdentifier)+1:])
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
		p.Data == p2.Data
}
