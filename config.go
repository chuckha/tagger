package tagger

import (
	"encoding/json"
	"fmt"

	"github.com/chuckha/tagger/id3v23/frames"

	"gitlab.com/tozd/go/errors"
)

// Config specifies the ID3 Frames to be added to every mp3 file that goes through the pipeline.
// There is an extraction language available to extract values from the path and put them into tags.
// This also supports file-based additions like lyrics or pictures as well as things like compression.
type Config struct {
	Frames map[string]frames.FrameBody
}

func NewConfig() *Config {
	return &Config{
		Frames: make(map[string]frames.FrameBody),
	}
}

func (c *Config) UnmarshalJSON(data []byte) error {
	var cfg struct {
		Frames map[string]json.RawMessage
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return errors.WithStack(err)
	}
	for k, data := range cfg.Frames {
		switch frames.IDToFrameKind[k] {
		case frames.TextInformationKind:
			ti := &frames.TextInformation{}
			if err := ti.UnmarshalJSON(data); err != nil {
				return errors.WithStack(err)
			}
			c.Frames[k] = ti
		case frames.AttachedPictureKind:
			ap := &frames.AttachedPicture{}
			if err := ap.UnmarshalJSON(data); err != nil {
				return errors.WithStack(err)
			}
			c.Frames[k] = ap
		default:
			panic(fmt.Sprintf("config does not support frame %q", k))
		}
	}
	return nil
}

/*

// read in a directory and walk it
// first, replace the input pattern string with the correct regex and compile it
// then extract the data from the file path
// apply the overrides (anything not captured the file path but used in the template is required)
// execute the template with the input data + overrides
// apply the template to the file
// write out the

// reads in data from the file name
// then populates the config with a template
// then runs the program.

{
	InputPattern: "Harry Potter and the Half Blood Prince Disk %disk%/%reader% - Chapter $chapter$-$part$.mp3",
	Overrides: {
		Title: "Harry Potter and the Half-Blood Prince",
		Reader: "Stephen Fry"
	},

	RenamePattern: "%track%-%title%.%ext%"
	FramesTemplate: `{
		"TALB": "Harry Potter and the Half-Blood Prince",
		"TEXT": "J.K. Rowling",
		"TPE1": "Stephen Fry",
		"TIT2": "{{.ChapterTitle}}",
	}`
}

If you want to use variables

1. Extract data from the file title
2. get specific data from custom input (?).
2. populate the frames template.
3. apply the frame config to the file.

*/
