package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/chuckha/tagger"
	"github.com/chuckha/tagger/id3v23/tags"
)

// tagger info <file>
// tagger tag --config <cfg.json> <file>
// tagger template-tag --template-config <cfg.json> <dir>
// tagger strip-tag-v1 <file>
func main() {

	tagfs := flag.NewFlagSet("tag", flag.ExitOnError)
	cfg := tagfs.String("config", "", "path to config file")

	templateTagfs := flag.NewFlagSet("template-tag", flag.ExitOnError)
	templateCfg := templateTagfs.String("template-config", "", "path to template config file")
	dryRun := templateTagfs.Bool("dry-run", true, "dry run")

	stripTagfs := flag.NewFlagSet("strip-tag", flag.ExitOnError)

	switch os.Args[1] {
	case "info":
		file := os.Args[2]
		tag, err := tags.NewID3v2FromFile(file)
		if err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		fmt.Println(tag)
	case "tag":
		tagfs.Parse(os.Args[2:])
		file := tagfs.Arg(0)
		b, err := os.ReadFile(*cfg)
		if err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		cfg := tagger.NewConfig()
		if err := cfg.UnmarshalJSON(b); err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		fmt.Println(cfg)
		tag, err := tags.NewID3v2FromFile(file)
		if err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		fmt.Println("APPLYING...")
		tag.ApplyConfig(cfg)
		fmt.Println(tag)
	case "template-tag":
		templateTagfs.Parse(os.Args[2:])
		// read in the template config
		tmplfile, err := os.ReadFile(*templateCfg)
		if err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		tmplcfg := tagger.NewTemplateConfig()
		if err := tmplcfg.UnmarshalJSON(tmplfile); err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		fmt.Println(tmplcfg)
		// first, replace the input pattern string with the correct regex and compile it
		regexpattern := subRegex(tmplcfg.FilePattern)
		//		fmt.Println("regexpattern", regexpattern)
		re := regexp.MustCompile(regexpattern)

		// compile the output file template
		outFileTmpl, err := template.New("output").Funcs(tmplFuncs()).Parse(tmplcfg.OutputFilePattern)
		if err != nil {
			panic(fmt.Sprintf("%+v", err))
		}

		// compile the config template
		tmpl, err := template.New("frames").Funcs(tmplFuncs()).Parse(tmplcfg.FramesTemplate)
		if err != nil {
			panic(fmt.Sprintf("%+v", err))
		}

		// read in a directory and walk it
		dir := templateTagfs.Arg(0)
		total := 0
		err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			// don't process dirs, but do kep walking
			if d.IsDir() {
				return nil
			}
			matches := re.FindStringSubmatch(path)
			if len(matches) == 0 {
				// a file to ignore, doesn't match the pattern
				return nil
			}
			total++
			return nil
		})
		if err != nil {
			panic(fmt.Sprintf("%+v", err))
		}

		// create special variables
		special := map[string]any{
			"count": 1,
			"total": total,
		}
		err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			// don't process dirs, but do kep walking
			if d.IsDir() {
				return nil
			}
			matches := re.FindStringSubmatch(path)
			if len(matches) == 0 {
				// a file to ignore, doesn't match the pattern
				return nil
			}
			// then extract the data from the file path
			fmt.Println("working on", path)
			// set up the config template context
			extracted := map[string]any{}
			for i, name := range re.SubexpNames() {
				if name == "" || name == "ignore" {
					continue
				}
				extracted[name] = matches[i]
			}
			// apply the overrides (anything not captured in the file path but used in the template is required)
			// continue setting up the config template context
			for name, override := range tmplcfg.Overrides {
				extracted[name] = override
			}
			extracted["userData"] = tmplcfg.UserData
			extracted["special"] = special
			// get the config for the file
			var b bytes.Buffer
			if err := tmpl.Execute(&b, extracted); err != nil {
				panic(fmt.Sprintf("%+v", err))
			}
			nc := tagger.NewConfig()
			if err := nc.UnmarshalJSON(b.Bytes()); err != nil {
				panic(fmt.Sprintf("%+v", err))
			}
			tag, err := tags.NewID3v2FromFile(path)
			if err != nil {
				var e *tags.NoID3v2IdentifierError
				if !errors.As(err, &e) {
					panic(fmt.Sprintf("%+v", err))
				}
				if !tmplcfg.AddMissingTag() {
					fmt.Printf("skipping %q; no id3 file identifier\n", path)
					return nil
				}
				tag = tags.NewID3v2()
			}
			//			fmt.Println(tag)
			tag.ApplyConfig(nc)
			//			fmt.Println(tag)
			special["count"] = special["count"].(int) + 1
			// generate the outfile name from the outfile pattern
			var outFile bytes.Buffer
			if err := outFileTmpl.Execute(&outFile, extracted); err != nil {
				panic(fmt.Sprintf("%+v", err))
			}
			if !*dryRun {
				return tag.Write(path, outFile.String())
			}
			fmt.Printf("[dry run] would have written %q\n", outFile.String())
			return nil

		})
		if err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
	// apply the template to the file
	// write out the
	// dir := templateTagfs.Arg(0)
	// err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
	// 	if err != nil {
	// 		return err
	// 	}

	// })
	// if err != nil {
	// 	panic(fmt.Sprintf("%+v", err))
	// }
	case "strip-tag":
		stripTagfs.Parse(os.Args[2:])
		file := stripTagfs.Arg(0)
		f, err := os.OpenFile(file, os.O_RDWR, 0644)
		if err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		info, err := f.Stat()
		if err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		n, err := f.Seek(-128, 2)
		if err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		fmt.Println("n", n, info.Size())
		v1Tag := make([]byte, 128)
		if _, err := f.Read(v1Tag); err != nil && err != io.EOF {
			panic(fmt.Sprintf("%+v", err))
		}
		fmt.Println("v1Tag", string(v1Tag))
		if string(v1Tag[0:3]) != "TAG" {
			fmt.Println("no id3v1 tag discovered")
			os.Exit(0)
		}
		fmt.Println("stripping id3v1 tag")
		if err := f.Truncate(info.Size() - 128); err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		f.Close()
	default:
		panic(fmt.Sprintf("unknown command: %q	", os.Args[1]))
	}
}

// subRegex substitutes the easier to read input pattern with the regex pattern
func subRegex(inputPattern string) string {
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

// ID3Tag should not be aware of the file...
// Open the file for Writing if there is no outputFilePattern defined.
// Otherwise, only open the file for reading.

// ID3Tag; it should just be the header data and some metadata like file location.
// Then something else is responsible for writing the updates. The TagWriter.
