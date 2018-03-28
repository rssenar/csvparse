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
	vhdr := flag.Bool("vhdr", true, "Validate required header fields")
	otype := flag.String("otype", "csv", "Specify output type: csv | json")
	flag.Parse()

	var input io.Reader

	args := flag.Args()

	if len(args) != 0 {
		if len(args[1:]) > 1 {
			log.Fatalln("err: Cannot parse multiple files")
		}
		file, err := os.Open(args[1])
		defer file.Close()
		if err != nil {
			log.Fatalf("%v : No such file or directory\n", os.Args[1])
		}
		input = file
	} else {
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

	err := p.UnMarshalCSV(*vhdr)
	if err != nil {
		log.Fatalln(err)
	}

	if *otype == "csv" {
		err = p.MarshaltoCSV()
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		data, err := json.MarshalIndent(p.Records, " ", " ")
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(string(data))
	}
}
