package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/chuckha/tagger"
	"github.com/chuckha/tagger/id3v23/tags"
)

func main() {
	infofs := flag.NewFlagSet("info", flag.ExitOnError)
	infofs.Usage = func() {
		fmt.Println("tagger info <file>")
	}

	tagfs := flag.NewFlagSet("tag", flag.ExitOnError)
	cfg := tagfs.String("config", "", "path to config file")
	tagfs.Usage = func() {
		fmt.Println("tagger tag -config <cfg.json> <file>")
	}

	templateTagfs := flag.NewFlagSet("template-tag", flag.ExitOnError)
	templateCfg := templateTagfs.String("template-config", "", "path to template config file")
	dryRun := templateTagfs.Bool("dry-run", true, "dry run")
	noisy := templateTagfs.Bool("noisy", false, "noisy")
	templateTagfs.Usage = func() {
		fmt.Println("tagger template-tag -template-config <cfg.json> [-dry-run=false] [-noisy] <dir>")
	}

	stripTagfs := flag.NewFlagSet("strip-tag", flag.ExitOnError)

	if len(os.Args) < 2 {
		flag.Usage = func() {
			fmt.Println("tagger <command> [args]")
			fmt.Println("commands:")
			fmt.Println("  info <file>")
			fmt.Println("  tag --config <cfg.json> <file>")
			fmt.Println("  template-tag --template-config <cfg.json> <dir>")
			fmt.Println("  strip-tag-v1 <file>")
		}
		flag.Usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "info":
		infofs.Parse(os.Args[2:])
		file := infofs.Arg(0)
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
		if err := tag.ApplyFrames(cfg.Frames); err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		fmt.Println(tag)
	case "template-tag":
		templateTagfs.Parse(os.Args[2:])
		dir := templateTagfs.Arg(0)
		tmplfile, err := os.ReadFile(*templateCfg)
		if err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		tmplcfg := tagger.NewTemplateConfig()
		if err := tmplcfg.UnmarshalJSON(tmplfile); err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		if *dryRun {
			tmplcfg.UpdateBehavior(tagger.WriteFile, tagger.Skip)
		}
		if *noisy {
			tmplcfg.UpdateBehavior(tagger.Logging, tagger.Noisy)
		}
		if err := tmplcfg.ProcessDir(dir); err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
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
