package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rssenar/csvparse"
)

type record struct {
	Fullname   string    `json:"Full_name" csv:"(?i)^fullname$" fmt:"tc"`
	Firstname  string    `json:"First_name" csv:"(?i)^first[ _-]?name$" fmt:"tc"`
	MI         string    `json:"Middle_name" csv:"(?i)^mi$" fmt:"uc"`
	Lastname   string    `json:"Last_name" csv:"(?i)^last[ _-]?name$" fmt:"tc"`
	Address1   string    `json:"Address_1" csv:"(?i)^address[ _-]?1?$" fmt:"tc"`
	Address2   string    `json:"Address_2" csv:"(?i)^address[ _-]?2$" fmt:"tc"`
	City       string    `json:"City" csv:"(?i)^city$" fmt:"tc"`
	State      string    `json:"State" csv:"(?i)^state$|^st$" fmt:"uc"`
	Zip        string    `json:"Zip" csv:"(?i)^(zip|postal)[ _]?(code)?$" fmt:"uc"`
	Zip4       string    `json:"Zip_4" csv:"(?i)^zip4$|^4zip$" fmt:"uc"`
	CRRT       string    `json:"CRRT" csv:"(?i)^crrt$" fmt:"uc"`
	HPH        string    `json:"Home_phone" csv:"(?i)^hph$|^home[ _]phone$" fmt:"fp"`
	BPH        string    `json:"Business_phone" csv:"(?i)^bph$|^(work|business)[ _]phone$" fmt:"fp"`
	CPH        string    `json:"Mobile_phone" csv:"(?i)^cph$|^mobile[ _]phone$" fmt:"fp"`
	Email      string    `json:"Email" csv:"(?i)^email[ _]?(address)?$" fmt:"lc"`
	VIN        string    `json:"VIN" csv:"(?i)^vin$" fmt:"-"`
	Year       string    `json:"Year" csv:"(?i)^year$|^vyr$" fmt:"-"`
	Make       string    `json:"Make" csv:"(?i)^make$|^vmk$" fmt:"tc"`
	Model      string    `json:"Model" csv:"(?i)^model$|^vmd$" fmt:"tc"`
	DelDate    time.Time `json:"Delivery_date" csv:"(?i)^del[ ]?date$" fmt:"-"`
	Date       time.Time `json:"Last_service_date" csv:"(?i)^date$" fmt:"-"`
	DSFwalkseq string    `json:"DSF_Walk_Sequence" csv:"(?i)^DSF_WALK_SEQ$" fmt:"uc"`
	KBB        string    `json:"KBB" csv:"(?i)^kbb$" fmt:"uc"`
}

func main() {
	// var j = flag.Bool("j", false, "Enable output to Indented JSON")
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

	// Pass in black []*Record{} container to be filled
	x := []*record{}

	err := csvparse.NewDecoder(input).DecodeCSV(&x)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// err = csvparse.NewEncoder(os.Stdout).EncodeJSON(x)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	err = csvparse.NewEncoder(os.Stdout).EncodeCSV(x)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// if *j {
	// 	// MarshalIndent outputs []*Record to JSON for output to os.Stdout
	// 	// []*Record are stored in the *Parser.Records struct filed
	// 	err := csvparse.NewEncoder(os.Stdout).EncodeJSON(data)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		os.Exit(1)
	// 	}
	// } else {
	// 	// MarshaltoCSV marshals []*Record to []string for output to os.Stdout
	// 	// []*Record are stored in the *Parser.Records struct filed
	// 	err = csvparse.NewEncoder(os.Stdout).EncodeCSV(data)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		os.Exit(1)
	// 	}
	// }
}
