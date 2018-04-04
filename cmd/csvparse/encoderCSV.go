package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"reflect"
	"sync"
	"time"

	"github.com/blendlabs/go-name-parser"
	cp "github.com/rssenar/csvparse"
)

type csvEncoder struct {
	output io.Writer
}

// NewEncoder initializes a new output
func newEncoder(output io.Writer) *csvEncoder {
	return &csvEncoder{output: output}
}

// EncodeCSV marshalls the Record struct then outputs to csv
func (e *csvEncoder) encodeCSV(data []*client, concurency *int) error {
	defer timeTrack(time.Now(), "encodeCSV")

	tasks := make(chan *client)
	go func() {
		for _, client := range data {
			tasks <- client
		}
		close(tasks)
	}()

	results := make(chan []string)
	var wg sync.WaitGroup
	wg.Add(*concurency)

	go func() {
		wg.Wait()
		close(results)
	}()

	for i := 0; i < *concurency; i++ {
		go func() {
			defer wg.Done()
			for t := range tasks {
				r, err := formatStructFields(t)
				if err != nil {
					log.Println(err)
					continue
				}
				results <- r
			}
		}()
	}
	if err := outputCSV(e.output, results); err != nil {
		log.Printf("could not write to %s: %v", e.output, err)
	}
	return nil
}

func outputCSV(output io.Writer, records <-chan []string) error {
	w := csv.NewWriter(output)
	for r := range records {
		if err := w.Write(r); err != nil {
			log.Fatalf("could not write record to csv: %v", err)
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return fmt.Errorf("writer failed: %v", err)
	}
	return nil
}

func formatStructFields(t *client) ([]string, error) {
	if t.Fullname != "" && (t.Firstname == "" || t.Lastname == "") {
		name := names.Parse(t.Fullname)
		t.Firstname = name.FirstName
		t.MI = name.MiddleName
		t.Lastname = name.LastName
	}
	if t.Zip != "" {
		zip, zip4 := cp.ParseZip(t.Zip)
		t.Zip = zip
		if zip4 != "" {
			t.Zip4 = zip4
		}
	}
	var row []string
	sValue := reflect.ValueOf(t).Elem()
	for i := 0; i < sValue.NumField(); i++ {
		var value string
		name := reflect.Indirect(sValue).Type().Field(i).Name
		switch sValue.Field(i).Type() {
		case reflect.TypeOf(time.Now()):
			time := fmt.Sprint(sValue.FieldByName(name))[:10]
			if time == "0001-01-01" {
				time = ""
			}
			value = time
		default:
			if format, ok := reflect.Indirect(sValue).Type().Field(i).Tag.Lookup("fmt"); ok {
				switch format {
				case "-":
					value = fmt.Sprint(sValue.FieldByName(name))
				default:
					fmtvalue, err := cp.FormatStringVals(format, fmt.Sprint(sValue.FieldByName(name)))
					if err != nil {
						return nil, err
					}
					value = fmtvalue
				}
			} else {
				value = fmt.Sprint(sValue.FieldByName(name))
			}
		}
		row = append(row, value)
	}
	return row, nil
}
