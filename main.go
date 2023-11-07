package main

import (
	"fmt"
	"os"

	"tagger/tags"
)

func perr(err error) {
	panic(fmt.Sprintf("%+v", err))
}

// cleanup
// if an mp3 has multiple TRCK tags, ask which one to keep
// Configuration: keep first, keep last, keep all, keep none?
// Duplicate tag configuration: keep first; ask; keep all
// override for specific tags

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: tagger <file>")
		os.Exit(1)
	}

	file := os.Args[1]
	fmt.Printf("file: %q\n", file)
	f, err := os.Open(file)
	if err != nil {
		perr(err)
	}
	tag := tags.NewID3v2()
	headerBytes := make([]byte, 10)
	if _, err := f.Read(headerBytes); err != nil {
		perr(err)
	}
	if string(headerBytes[0:3]) != "ID3" {
		fmt.Printf("%q has no id3 tag\n", file)
		os.Exit(0)
	}
	if err := tag.Header.UnmarshalBinary(headerBytes); err != nil {
		perr(err)
	}
	if tag.Header.MajorVersion != 3 {
		fmt.Printf("%q is version v2.%d.%d tag\n", file, tag.Header.MajorVersion, tag.Header.Revision)
		os.Exit(0)
	}
	//	fmt.Println(tag.Header)
	// there could be an extended header here
	tagBytes := make([]byte, tag.Header.Size)
	if _, err := f.Read(tagBytes); err != nil {
		perr(err)
	}
	if err := tag.Frames.UnmarshalBinary(tagBytes); err != nil {
		perr(err)
	}

	// w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.AlignRight|tabwriter.Debug)
	// fmt.Fprintf(w, "%s\t%s\n", "Header", tag.Header.String())
	// for _, frame := range tag.Frames {
	// 	fmt.Fprintf(w, "%s:\t%v\n", frame.Header, frame.Body.String())
	// }
	// w.Flush()
}
