package main

import (
	"flag"
	"io"
	"log"
	"os"
	"time"

	"github.com/rssenar/csvparse"
)

type client struct {
	Fullname   string    `json:"Full_name" csv:"(?i)^fullname$" fmt:"tc"`
	Firstname  string    `json:"First_name" csv:"(?i)^first[ _-]?name$" fmt:"tc"`
	MI         string    `json:"Middle_name" csv:"(?i)^mi$" fmt:"uc"`
	Lastname   string    `json:"Last_name" csv:"(?i)^last[ _-]?name$" fmt:"tc"`
	Address1   string    `json:"Address_1" csv:"(?i)^address[ _-]?1?$" fmt:"tc"`
	Address2   string    `json:"Address_2" csv:"(?i)^address[ _-]?2$" fmt:"tc"`
	City       string    `json:"City" csv:"(?i)^city$" fmt:"tc"`
	State      string    `json:"State" csv:"(?i)^state$|^st$" fmt:"uc"`
	Zip        string    `json:"Zip" csv:"(?i)^(zip|postal)[ _]?(code)?$" fmt:"-"`
	Zip4       string    `json:"Zip_4" csv:"(?i)^zip4$|^4zip$" fmt:"-"`
	CRRT       string    `json:"CRRT" csv:"(?i)^crrt$" fmt:"uc"`
	DSFwalkseq string    `json:"DSF_Walk_Sequence" csv:"(?i)^DSF_WALK_SEQ$" fmt:"uc"`
	HPH        string    `json:"Home_phone" csv:"(?i)^hph$|^home[ _]phone$" fmt:"fp"`
	BPH        string    `json:"Business_phone" csv:"(?i)^bph$|^(work|business)[ _]phone$" fmt:"fp"`
	CPH        string    `json:"Mobile_phone" csv:"(?i)^cph$|^mobile[ _]phone$" fmt:"fp"`
	Email      string    `json:"Email" csv:"(?i)^email[ _]?(address)?$" fmt:"lc"`
	VIN        string    `json:"VIN" csv:"(?i)^vin$" fmt:"uc"`
	Year       string    `json:"Veh_Year" csv:"(?i)^year$|^vyr$" fmt:"-"`
	Make       string    `json:"Veh_Make" csv:"(?i)^make$|^vmk$" fmt:"tc"`
	Model      string    `json:"Veh_Model" csv:"(?i)^model$|^vmd$" fmt:"tc"`
	DelDate    time.Time `json:"Delivery_date" csv:"(?i)^del[ ]?date$" fmt:"-"`
	Date       time.Time `json:"Last_service_date" csv:"(?i)^date$" fmt:"-"`
}

func main() {
	// Optional timer function for determining function duration
	defer timeTrack(time.Now(), "Application")
	// Set CLI flags
	var j = flag.Bool("j", false, "Enable output to Indented JSON")
	var concurency = flag.Int("c", 100, "Set number of GoRoutines")
	flag.Parse()

	var input io.Reader
	// return the non-flag command-line arguments
	args := flag.Args()

	if len(args) != 0 {
		// verify if file was passed through as a command-line argument
		if len(args) > 1 {
			log.Fatalln("Error: Cannot parse multiple files")
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
			log.Fatalf("%v : Error reading stdin file info\n", err)
		}
		if fi.Size() == 0 {
			log.Fatalf("Input file not specified")
		}
		input = os.Stdin
	}

	// Pass in black []*Record{} container to be filled
	data := []*client{}

	err := csvparse.NewDecoder(input).DecodeCSV(&data)
	if err != nil {
		log.Fatalln(err)
	}
	switch {
	case *j:
		err = csvparse.NewEncoder(os.Stdout).EncodeJSON(data)
		if err != nil {
			log.Fatalln(err)
		}
	default:
		err = newEncoder(os.Stdout).encodeCSV(data, concurency)
		if err != nil {
			log.Println(err)
		}
	}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
