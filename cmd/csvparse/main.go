package main

import (
	"io"
	"log"
	"os"

	cp "github.com/rssenar/csvparse"
)

func main() {
	var input io.Reader
	if len(os.Args[1:]) != 0 {
		if len(os.Args[1:]) > 1 {
			log.Fatalln("err: Cannot parse multiple files")
		}
		file, err := os.Open(os.Args[1])
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

	err := p.UnMarshalCSV()
	if err != nil {
		log.Fatalln(err)
	}

	err = p.MarshaltoCSV()
	if err != nil {
		log.Fatalln(err)
	}
	// data, _ := json.MarshalIndent(p.Records, " ", " ")
	// fmt.Println(string(data))
}
