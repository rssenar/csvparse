package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	cp "github.com/rssenar/csvparse"
)

func main() {
	var (
		vh = flag.Bool("vh", true, "Validate required header fields: (true|false)")
		ot = flag.String("ot", "csv", "Specify output type: (csv|json)")
	)
	flag.Parse()

	var input io.Reader

	// return the non-flag command-line arguments
	args := flag.Args()

	if len(args) != 0 {
		// verify if file was passed through as a command-line argument
		if len(args) > 1 {
			log.Fatalln("err: Cannot parse multiple files")
		}
		file, err := os.Open(args[0])
		defer file.Close()
		if err != nil {
			log.Fatalf("%v : No such file or directory\n", args[0])
		}
		input = file
	} else {
		// verify if file was passed through from os.Stdin
		fi, err := os.Stdin.Stat()
		if err != nil {
			log.Fatalf("%v : Error reading stdin\n", err)
		}
		if fi.Size() == 0 {
			log.Fatalln("err: please pass file to stdin")
		}
		input = os.Stdin
	}

	p := cp.New(input)
	// UnMarshalCSV unmarshalls io.reader into csvparse.Record struct
	// []*Record are stored in the *Parser.Records struct filed
	err := p.UnMarshalCSV(vh)
	if err != nil {
		log.Fatalln(err)
	}

	if *ot == "csv" {
		// MarshaltoCSV marshals []*Record to []string for output to os.Stdout
		// []*Record are stored in the *Parser.Records struct filed
		err = p.MarshaltoCSV()
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		// MarshalIndent outputs []*Record to JSON for output to os.Stdout
		data, err := json.MarshalIndent(p.Records, " ", " ")
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(data))
	}
}
