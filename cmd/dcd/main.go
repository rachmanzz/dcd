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

func main() {
	dataFile := flag.String("data", "", "JSON file with variables")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: dcd [--data file.json] <input.dcd> [output]\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	input := flag.Arg(0)
	output := flag.Arg(1)
	if output == "" {
		output = "output.docx"
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

	r := render.NewDocxRenderer()

	if err := render.New(doc, ds, r).Run(output); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %s\n", output)
}
