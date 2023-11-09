package main

import (
	"flag"
	"fmt"
	"os"
	"tagger/id3v23/tags"
)

// tagger text-info --frame-id TRCK --value abc <file>
// frame id inputs...can they be dynamic?
func main() {
	fmt.Println(os.Args)
	flag.Usage = func() {
		fmt.Println("Usage: tagger <command> [<args>]")
		fmt.Println("Commands:")
		fmt.Println("  info <file>")
		fmt.Println("  trck get <file>")
		fmt.Println("  trck set --value text <file>")
	}

	trckset := flag.NewFlagSet("trck set", flag.ExitOnError)
	trckset.Usage = func() {
		fmt.Println("Usage: tagger trck set <track> <file>")
	}
	trcksetval := trckset.String("value", "", "must be a number or two numbers separated by / (e.g. 1/10)")

	if len(os.Args) < 3 {
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
				trckset.Usage()
				os.Exit(1)
			}
			file := os.Args[5]
			tag, err := tags.NewID3v2FromFile(file)
			if err != nil {
				panic(fmt.Sprintf("%+v", err))
			}
			tag.SetTextFrame("TRCK", *trcksetval)
			if err := tag.Close(); err != nil {
				panic(fmt.Sprintf("%+v", err))
			}
			os.Exit(0)
		default:
			flag.Usage()
			os.Exit(1)
		}
	}
}
