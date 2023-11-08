package main

import (
	"flag"
	"fmt"
	"os"
	"tagger/tags"
)

func main() {
	trckset := flag.NewFlagSet("trck set", flag.ExitOnError)
	trckset.Usage = func() {
		fmt.Println("Usage: tagger trck set <track> <file>")
	}
	trcksetval := trckset.String("value", "", "must be a number or two numbers separated by / (e.g. 1/10)")

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "info":
		file := os.Args[2]
		tag, err := tags.NewID3v2FromFile(file)
		if err != nil {
			panic(fmt.Sprintf("%+v", err))
		}
		fmt.Println(tag)
	case "trck":
		switch os.Args[2] {
		case "get":
			fmt.Println("get")
		case "set":
			trckset.Parse(os.Args[3:])
			if trcksetval == nil {
				flag.Usage()
				os.Exit(1)
			}
			file := os.Args[5]
			tag, err := tags.NewID3v2FromFile(file)
			if err != nil {
				panic(fmt.Sprintf("%+v", err))
			}
			tag.SetTextFrame("TRCK", *trcksetval)
		default:
			flag.Usage()
			os.Exit(1)
		}
	}
}
