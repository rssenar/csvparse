package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/rssenar/csvparse"
)

func main() {
	var j = flag.Bool("j", false, "Enable output to Indented JSON")
	flag.Parse()

	var input io.Reader

	// return the non-flag command-line arguments
	args := flag.Args()

	if len(args) != 0 {
		// verify if file was passed through as a command-line argument
		if len(args) > 1 {
			fmt.Println("Error: Cannot parse multiple files")
			os.Exit(1)
		}
		file, err := os.Open(args[0])
		defer file.Close()
		if err != nil {
			fmt.Printf("%v : No such file or directory\n", args[0])
			os.Exit(1)
		}
		input = file
	} else {
		// verify if file was passed through from os.Stdin
		fi, err := os.Stdin.Stat()
		if err != nil {
			fmt.Printf("%v : Error reading stdin file info\n", err)
			os.Exit(1)
		}
		if fi.Size() == 0 {
			fmt.Println("Input file not specified")
			os.Exit(1)
		}
		input = os.Stdin
	}

	data, err := csvparse.NewDecoder(input).DecodeCSV()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *j {
		// MarshalIndent outputs []*Record to JSON for output to os.Stdout
		// []*Record are stored in the *Parser.Records struct filed
		err := csvparse.NewEncoder(os.Stdout).EncodeJSON(data)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		// MarshaltoCSV marshals []*Record to []string for output to os.Stdout
		// []*Record are stored in the *Parser.Records struct filed
		err = csvparse.NewEncoder(os.Stdout).EncodeCSV(data)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
