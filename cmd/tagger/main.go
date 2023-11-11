package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"tagger"
	"tagger/id3v23/tags"
	"text/template"
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
		fmt.Println("regexpattern", regexpattern)
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
			extracted := map[string]any{}
			for i, name := range re.SubexpNames() {
				if name == "" {
					continue
				}
				extracted[name] = matches[i]
			}
			// apply the overrides (anything not captured the file path but used in the template is required)
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
				panic(fmt.Sprintf("%+v", err))
			}
			//			fmt.Println(tag)
			tag.ApplyConfig(nc)
			//			fmt.Println(tag)
			special["count"] = special["count"].(int) + 1
			var outFile bytes.Buffer
			if err := outFileTmpl.Execute(&outFile, extracted); err != nil {
				panic(fmt.Sprintf("%+v", err))
			}
			if !*dryRun {
				out, err := tag.Output()
				if err != nil {
					panic(fmt.Sprintf("%+v", err))
				}
				if err := os.WriteFile(outFile.String(), out, 0644); err != nil {
					panic(fmt.Sprintf("%+v", err))
				}
				return tag.Close()
			}
			fmt.Printf("[dry run] would have written %q\n", outFile.String())
			return tag.Close()

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
	inputPattern = strings.ReplaceAll(inputPattern, "%ignore%", `(?P<ignore>.+)`)
	inputPattern = strings.Replace(inputPattern, "$disk$", `(?P<disk>\d+)`, 1)
	inputPattern = strings.Replace(inputPattern, "$track$", `(?P<track>\d+)`, 1)
	inputPattern = strings.Replace(inputPattern, "%title%", `(?P<title>.+)`, 1)
	inputPattern = strings.Replace(inputPattern, "%reader%", `(?P<reader>.+)`, 1)
	inputPattern = strings.Replace(inputPattern, "$chapter$", `(?P<chapter>\d+)`, 1)
	inputPattern = strings.Replace(inputPattern, "$part$", `(?P<part>\d+)`, 1)
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
