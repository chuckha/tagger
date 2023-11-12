package tagger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/chuckha/tagger/id3v23/tags"
	"gitlab.com/tozd/go/errors"
)

type Situation string
type Behavior string

const (
	Add   Behavior = "add"
	Noisy Behavior = "noisy"
	Skip  Behavior = "skip"

	MissingTag Situation = "missing-id3v2-tag"
	Logging    Situation = "logging"
	WriteFile  Situation = "write-file"
)

// TemplateConfig is a user defined template config.
// It references a FramesTemplate that is used to generate a Config for file.
// The config, once rendered, will define the tags that should exist on the file.
type TemplateConfig struct {
	FilePattern        *regexp.Regexp
	Overrides          map[string]any
	OutputFileTemplate *template.Template
	// FramesTemplate is a pointer to a frames template json file.
	// This is so the user does not have to do weird escaping within a string.
	FramesTemplate *template.Template
	UserData       any
	Behavior       map[Situation]Behavior

	// special is an internal variable that holds aggregate values across all files.
	// special is available in all templates.
	special map[string]any
}

func NewTemplateConfig() *TemplateConfig {
	return &TemplateConfig{
		Overrides: make(map[string]any),
		Behavior:  make(map[Situation]Behavior),
		special:   make(map[string]any),
	}
}

func (t *TemplateConfig) UnmarshalJSON(b []byte) error {
	var cfg struct {
		FilePattern       string
		Overrides         map[string]any
		OutputFilePattern string
		FramesTemplate    string
		UserData          any
		Behavior          map[Situation]Behavior
	}
	if err := json.Unmarshal(b, &cfg); err != nil {
		return errors.WithStack(err)
	}
	// straight forward conversions
	t.Overrides = cfg.Overrides
	t.UserData = cfg.UserData
	t.Behavior = cfg.Behavior

	// regexp
	t.FilePattern = regexp.MustCompile(subRegex(cfg.FilePattern))

	// templates
	outFileTmpl, err := template.New("output").Funcs(tmplFuncs()).Parse(cfg.OutputFilePattern)
	if err != nil {
		return errors.WithStack(err)
	}
	t.OutputFileTemplate = outFileTmpl

	framesb, err := os.ReadFile(cfg.FramesTemplate)
	if err != nil {
		return errors.WithStack(err)
	}
	tmpl, err := template.New("frames").Funcs(tmplFuncs()).Parse(string(framesb))
	if err != nil {
		return errors.WithStack(err)
	}
	t.FramesTemplate = tmpl
	return nil
}

func (t *TemplateConfig) UpdateBehavior(situation Situation, behavior Behavior) {
	t.Behavior[situation] = behavior
}

func (t *TemplateConfig) SetupSpecial(dir string) error {
	total := 0
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		matches := t.FilePattern.FindStringSubmatch(path)
		if len(matches) == 0 {
			return nil
		}
		total++
		return nil
	})
	if err != nil {
		return err
	}
	t.special["total"] = total

	t.special["count"] = 1
	return nil
}

func (t *TemplateConfig) ProcessDir(dir string) error {
	if err := t.SetupSpecial(dir); err != nil {
		return err
	}
	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		matches := t.FilePattern.FindStringSubmatch(path)
		if len(matches) == 0 {
			return nil
		}
		// TODO: move this into a logging object
		if t.Behavior[Logging] == Noisy {
			fmt.Println("working on", path)
		}
		// set up the config template context
		extracted := map[string]any{}
		for i, name := range t.FilePattern.SubexpNames() {
			if name == "" || name == "ignore" {
				continue
			}
			extracted[name] = matches[i]
		}
		// apply the overrides (anything not captured in the file path but used in the template is required)
		// continue setting up the config template context
		for name, override := range t.Overrides {
			extracted[name] = override
		}
		extracted["userData"] = t.UserData
		extracted["special"] = t.special
		// get the config for the file
		var b bytes.Buffer
		if err := t.FramesTemplate.Execute(&b, extracted); err != nil {
			return errors.WithStack(err)
		}
		nc := NewConfig()
		if err := nc.UnmarshalJSON(b.Bytes()); err != nil {
			return err
		}
		tag, err := tags.NewID3v2FromFile(path)
		if err != nil {
			var e *tags.NoID3v2IdentifierError
			if !errors.As(err, &e) {
				return err
			}
			if !t.AddMissingTag() {
				if t.Noisy() {
					fmt.Printf("skipping %q; no id3 file identifier\n", path)
				}
				return nil
			}
			tag = tags.NewID3v2()
		}
		if err := tag.ApplyFrames(nc.Frames); err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		t.special["count"] = t.special["count"].(int) + 1
		// generate the outfile name from the outfile pattern
		var outFile bytes.Buffer
		if err := t.OutputFileTemplate.Execute(&outFile, extracted); err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		if !t.DryRun() {
			return tag.Write(path, outFile.String())
		}
		fmt.Printf("[dry run] would have written %q\n", outFile.String())
		return nil
	})
}

func (t *TemplateConfig) DryRun() bool {
	return t.Behavior[WriteFile] == Skip
}

func (t *TemplateConfig) Noisy() bool {
	return t.Behavior[Logging] == Noisy
}

func (t *TemplateConfig) AddMissingTag() bool {
	return t.Behavior[MissingTag] == Add
}

// subRegex substitutes the easier to read input pattern with the regex pattern
func subRegex(inputPattern string) string {
	inputPattern = strings.ReplaceAll(inputPattern, "(", `\(`)
	inputPattern = strings.ReplaceAll(inputPattern, ")", `\)`)
	// find any \$(\w+)\$ and replace them with `(?P<$1>\d+))`
	// find any %(\w+)% and replace them with `(?P<$1>.+)
	// input is a user pattern like     "FilePattern": "fables_$volume$_$fable$_aesop_64kb.mp3",
	// output is a regular expression

	digitGrabber := regexp.MustCompile(`\$(\w+)\$`)
	for _, match := range digitGrabber.FindAllStringSubmatch(inputPattern, -1) {
		inputPattern = strings.ReplaceAll(inputPattern, match[0], fmt.Sprintf(`(?P<%s>\d+)`, match[1]))
	}
	wordGrabber := regexp.MustCompile(`%(\w+)%`)
	for _, match := range wordGrabber.FindAllStringSubmatch(inputPattern, -1) {
		inputPattern = strings.ReplaceAll(inputPattern, match[0], fmt.Sprintf(`(?P<%s>.+)`, match[1]))
	}
	return inputPattern
}

func tmplFuncs() template.FuncMap {
	return template.FuncMap{
		"get": get,
	}
}

func get(in any, key string) any {
	switch x := in.(type) {
	case []any:
		d, _ := strconv.Atoi(key)
		return x[d-1]
	default:
		panic(fmt.Sprintf("unknown type %T", x))
	}
}
