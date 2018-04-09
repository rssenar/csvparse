package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	cp "github.com/rssenar/csvparse"
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
	// defer timeTrack(time.Now(), "CSVParser")
	flag.Parse()
	args := flag.Args()

	var input io.Reader
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

	// Pass in blank []*Record{} container to be filled
	data := []*client{}
	var err error
	err = cp.NewDecoder(input).DecodeCSV(&data)
	// err = newEncoder(os.Stdout).encodeCSV(data, 1000)
	jdata, err := json.MarshalIndent(data, " ", " ")
	fmt.Println(string(jdata))
	if err != nil {
		log.Fatalln(err)
	}
}

// type csvEncoder struct{ output io.Writer }

// // NewEncoder initializes a new output
// func newEncoder(output io.Writer) *csvEncoder { return &csvEncoder{output: output} }

// // EncodeCSV marshalls the Record struct then outputs to csv
// func (e *csvEncoder) encodeCSV(data []*client, concurency int) error {
// 	// defer timeTrack(time.Now(), "EncodeStructtoCSV")

// 	var Client *client

// 	tasks := make(chan *client)
// 	go func() {
// 		for _, Client = range data {
// 			tasks <- Client
// 		}
// 		close(tasks)
// 	}()

// 	results := make(chan []string)
// 	var wg sync.WaitGroup

// 	wg.Add(concurency)

// 	go func() {
// 		wg.Wait()
// 		close(results)
// 	}()

// 	for i := 0; i < concurency; i++ {
// 		go func() {
// 			defer wg.Done()
// 			for t := range tasks {
// 				r, err := process(t)
// 				if err != nil {
// 					log.Println(err)
// 					continue
// 				}
// 				results <- r
// 			}
// 		}()
// 	}
// 	if err := print(e.output, results, Client); err != nil {
// 		log.Printf("could not write to %s: %v", e.output, err)
// 	}
// 	return nil
// }

// func process(t *client) ([]string, error) {
// 	if t.Fullname != "" && (t.Firstname == "" || t.Lastname == "") {
// 		name := names.Parse(t.Fullname)
// 		t.Firstname = name.FirstName
// 		t.MI = name.MiddleName
// 		t.Lastname = name.LastName
// 	}
// 	if t.Zip != "" {
// 		zip, zip4 := cp.ParseZip(t.Zip)
// 		t.Zip = zip
// 		if zip4 != "" {
// 			t.Zip4 = zip4
// 		}
// 	}

// 	var row []string
// 	sValue := reflect.ValueOf(t).Elem()
// 	for i := 0; i < sValue.NumField(); i++ {
// 		var value string
// 		name := reflect.Indirect(sValue).Type().Field(i).Name
// 		switch sValue.Field(i).Type() {
// 		case reflect.TypeOf(time.Now()):
// 			time := fmt.Sprint(sValue.FieldByName(name))[:10]
// 			if time == "0001-01-01" {
// 				time = ""
// 			}
// 			value = time
// 		default:
// 			if format, ok := reflect.Indirect(sValue).Type().Field(i).Tag.Lookup("fmt"); ok {
// 				switch format {
// 				case "-":
// 					value = fmt.Sprint(sValue.FieldByName(name))
// 				default:
// 					fmtvalue, err := cp.FormatStringVals(format, fmt.Sprint(sValue.FieldByName(name)))
// 					if err != nil {
// 						return nil, err
// 					}
// 					value = fmtvalue
// 				}
// 			} else {
// 				value = fmt.Sprint(sValue.FieldByName(name))
// 			}
// 		}
// 		row = append(row, value)
// 	}
// 	return row, nil
// }

// func print(output io.Writer, records <-chan []string, c *client) error {
// 	w := csv.NewWriter(output)

// 	var header []string
// 	sValue := reflect.ValueOf(c).Elem()
// 	for i := 0; i < sValue.NumField(); i++ {
// 		name := reflect.Indirect(sValue).Type().Field(i).Name
// 		header = append(header, name)
// 	}
// 	if err := w.Write(header); err != nil {
// 		log.Fatalf("could not write header to csv: %v", err)
// 	}

// 	for r := range records {
// 		if err := w.Write(r); err != nil {
// 			log.Fatalf("could not write record to csv: %v", err)
// 		}
// 	}
// 	w.Flush()
// 	if err := w.Error(); err != nil {
// 		return fmt.Errorf("writer failed: %v", err)
// 	}
// 	return nil
// }

// func timeTrack(start time.Time, name string) {
// 	elapsed := time.Since(start)
// 	log.Printf("%s took %s", name, elapsed)
// }
