package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rssenar/csvparse"
)

func main() {
	var (
		i = flag.Bool("i", false, "Enable output to Indented JSON")
		v = flag.Bool("v", false, "Enable header field validation")
		h = flag.String("h", "", "Specify Output fields [PKey,Fullname,Firstname,MI,Lastname,Address1,Address2,City,State,Zip,Zip4,HPH,BPH,CPH,Email,VIN,Year,Make,Model,DelDate,Date,DSFwalkse,CRRT,KBB]")
	)
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

	data, err := csvparse.NewDecoder(input, v).DecodeCSV()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Set output field options if -h flags are specified
	var fields []string
	if *h == "" {
		fields = nil
	} else {
		fields = strings.Split(*h, ",")
	}

	if *i {
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
		err = csvparse.NewEncoder(os.Stdout).EncodeCSV(data, fields)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
