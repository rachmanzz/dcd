package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/rachmanzz/dcd/data"
	"github.com/rachmanzz/dcd/parse"
	"github.com/rachmanzz/dcd/render"
)

var version = "0.3.0"

func main() {
	dataFile := flag.String("data", "", "JSON file with variables")
	format := flag.String("format", "docx", "Output format: docx or pdf")
	showVersion := flag.Bool("version", false, "Show version")
	flag.StringVar(dataFile, "d", "", "JSON file with variables (shorthand)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: dcd [OPTIONS] <input.dcd> [output]\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *showVersion {
		fmt.Printf("dcd version %s\n", version)
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	input := flag.Arg(0)
	output := flag.Arg(1)
	if output == "" {
		output = "output.docx"
		if *format == "pdf" {
			output = "output.pdf"
		}
	}

	doc, err := parse.Parse(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var src any
	if *dataFile != "" {
		b, err := os.ReadFile(*dataFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading --data: %v\n", err)
			os.Exit(1)
		}
		if err := json.Unmarshal(b, &src); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing --data: %v\n", err)
			os.Exit(1)
		}
	}
	ds := data.NewDataSet(src)

	var r render.Renderer
	switch *format {
	case "pdf":
		fmt.Fprintf(os.Stderr, "Error: PDF output was removed. Use DOCX format instead.\n")
		os.Exit(1)
	default:
		r = render.NewDocxRenderer()
	}

	if err := render.New(doc, ds, r).Run(output); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %s\n", output)
}
