package tagger

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/chuckha/tagger/id3v23/frames"

	"gitlab.com/tozd/go/errors"
)

const (
	MissingID3v2Tag    = "missing-id3v2-tag"
	AddMissingID3v2Tag = "add"
)

type TemplateConfig struct {
	FilePattern       string
	Overrides         map[string]any
	OutputFilePattern string
	// FramesTemplate is a pointer to a frames template json file.
	// This is so the user does not have to do weird escaping within a string.
	FramesTemplate string
	UserData       any
	Behavior       map[string]string
}

func NewTemplateConfig() *TemplateConfig {
	return &TemplateConfig{
		Overrides: make(map[string]any),
		Behavior:  make(map[string]string),
	}
}

func (t *TemplateConfig) UnmarshalJSON(b []byte) error {
	var cfg struct {
		FilePattern       string
		Overrides         map[string]any
		OutputFilePattern string
		FramesTemplate    string
		UserData          any
		Behavior          map[string]string
	}
	if err := json.Unmarshal(b, &cfg); err != nil {
		return errors.WithStack(err)
	}
	t.FilePattern = cfg.FilePattern
	t.Overrides = cfg.Overrides
	t.OutputFilePattern = cfg.OutputFilePattern
	t.UserData = cfg.UserData
	t.Behavior = cfg.Behavior
	b, err := os.ReadFile(cfg.FramesTemplate)
	if err != nil {
		return errors.WithStack(err)
	}
	t.FramesTemplate = string(b)
	return nil
}

func (t *TemplateConfig) AddMissingTag() bool {
	return t.Behavior[MissingID3v2Tag] == AddMissingID3v2Tag
}

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
