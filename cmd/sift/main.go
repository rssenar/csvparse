package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/rssenar/sift"
)

// D is the stuct the CSV file will be unmarshalled to
type D struct {
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
	DSFwalkseq string    `json:"DSF_Walk_Sequence" csv:"(?i)^DSF_WALK_SEQ$" fmt:"-"`
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
	flag.Parse()
	args := flag.Args()
	if len(args) != 0 {
		parseArgsInput(args)
	} else {
		fi, err := os.Stdin.Stat()
		if err != nil {
			log.Fatalf("%v : Error reading stdin\n", err)
		}
		if fi.Size() == 0 {
			log.Fatalln("Unspecified Input file")
		}
		parseStdinInput(os.Stdin)
	}
}

func parseStdinInput(input io.Reader) {
	start := time.Now()
	data := []*D{}
	err := sift.NewDecoder(input).DecodeCSV(&data)
	if err != nil {
		log.Fatalln(err)
	}
	w := csv.NewWriter(os.Stdout)
	sValue := reflect.ValueOf(data)

	var hrow []string
	for j := 0; j < sValue.Index(0).Elem().NumField(); j++ {
		hrow = append(hrow, fmt.Sprint(reflect.Indirect(sValue.Index(0).Elem()).Type().Field(j).Name))
	}
	if err := w.Write(hrow); err != nil {
		log.Fatalf("could not write header to csv: %v", err)
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatalf("writer failed: %v", err)
	}

	for i := 0; i < sValue.Len(); i++ {
		var row []string
		for j := 0; j < sValue.Index(i).Elem().NumField(); j++ {
			vType := sValue.Index(i).Elem().Field(j).Type()

			switch vType {
			case reflect.TypeOf(time.Now()):
				time := fmt.Sprint(sValue.Index(i).Elem().Field(j))[:10]
				if time == "0001-01-01" {
					time = ""
				}
				row = append(row, time)
			default:
				row = append(row, fmt.Sprint(sValue.Index(i).Elem().Field(j)))
			}
		}
		if err := w.Write(row); err != nil {
			log.Fatalf("could not write header to csv: %v", err)
		}
		w.Flush()
		if err := w.Error(); err != nil {
			log.Fatalf("writer failed: %v", err)
		}
	}
	Elapsed := time.Since(start)
	log.Printf("csvparse took %v", Elapsed)
}

func parseArgsInput(args []string) {
	for _, filename := range args {
		start := time.Now()
		data := []*D{}
		inputfile, err := os.Open(filename)
		defer inputfile.Close()
		if err != nil {
			log.Fatalln(err)
		}
		err = sift.NewDecoder(inputfile).DecodeCSV(&data)
		if err != nil {
			log.Fatalln(err)
		}
		outputfile, err := os.Create(fmt.Sprintf("%v_parsed.csv", inputfile.Name()[:len(inputfile.Name())-4]))
		defer outputfile.Close()
		if err != nil {
			log.Fatalln(err)
		}
		w := csv.NewWriter(outputfile)
		sValue := reflect.ValueOf(data)

		var hrow []string
		for j := 0; j < sValue.Index(0).Elem().NumField(); j++ {
			hrow = append(hrow, fmt.Sprint(reflect.Indirect(sValue.Index(0).Elem()).Type().Field(j).Name))
		}
		if err := w.Write(hrow); err != nil {
			log.Fatalf("could not write header to csv: %v", err)
		}
		w.Flush()
		if err := w.Error(); err != nil {
			log.Fatalf("writer failed: %v", err)
		}

		for i := 0; i < sValue.Len(); i++ {
			var row []string
			for j := 0; j < sValue.Index(i).Elem().NumField(); j++ {
				vType := sValue.Index(i).Elem().Field(j).Type()
				switch vType {
				case reflect.TypeOf(time.Now()):
					time := fmt.Sprint(sValue.Index(i).Elem().Field(j))[:10]
					if time == "0001-01-01" {
						time = ""
					}
					row = append(row, time)
				default:
					row = append(row, fmt.Sprint(sValue.Index(i).Elem().Field(j)))
				}
			}
			if err := w.Write(row); err != nil {
				log.Fatalf("could not write header to csv: %v", err)
			}
			w.Flush()
			if err := w.Error(); err != nil {
				log.Fatalf("writer failed: %v", err)
			}
		}
		Elapsed := time.Since(start)
		log.Printf("%v was parser in %v", inputfile.Name(), Elapsed)
	}
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

// val := sValue.Index(i).Elem().Field(j)
// name := reflect.Indirect(sValue.Index(i).Elem()).Type().Field(j).Name
