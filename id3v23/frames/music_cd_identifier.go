package frames

import (
	"encoding/json"

	"gitlab.com/tozd/go/errors"
)

// MusicCDIdentifier have the ID MCDI
type MusicCDIdentifier struct {
	TableOfContents []byte
}

func (m *MusicCDIdentifier) UnmarshalBinary(data []byte) error {
	m.TableOfContents = data
	return nil
}

func (m *MusicCDIdentifier) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, m); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (m *MusicCDIdentifier) String() string {
	return "toc: <MCDI TOC parsing is not implemented>"
}

func (m *MusicCDIdentifier) MarshalBinary() ([]byte, error) {
	return m.TableOfContents, nil
}

func (m *MusicCDIdentifier) Equal(m2 *MusicCDIdentifier) bool {
	return string(m.TableOfContents) == string(m2.TableOfContents)
}
