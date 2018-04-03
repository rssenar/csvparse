package csvparse

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"time"
)

// CSVEncoder defines output writer.
type CSVEncoder struct {
	output io.Writer
}

// NewEncoder initializes a new output
func NewEncoder(output io.Writer) *CSVEncoder {
	return &CSVEncoder{output: output}
}

// EncodeCSV marshalls the Record struct then outputs to csv
func (e *CSVEncoder) EncodeCSV(v interface{}) error {
	defer timeTrack(time.Now(), "EncodeCSV Func")
	wtr := csv.NewWriter(e.output)

	slice := checkifSlice(v)
	if slice.Kind() != reflect.Slice {
		return errors.New("Need to pass in a slice")
	}
	innerType := getInnerSliceType(v)
	if innerType.Kind() != reflect.Struct {
		return errors.New("Need to pass in a slice of stucts")
	}
	innerValue := slice.Index(0).Elem()

	var header []string
	for i := 0; i < innerValue.NumField(); i++ {
		header = append(header, reflect.Indirect(innerValue).Type().Field(i).Name)
	}
	if err := wtr.Write(header); err != nil {
		return fmt.Errorf("error writing header row: %v", err)
	}

	for i := 0; i < slice.Len(); i++ {
		innerValue := slice.Index(i).Elem()
		var row []string
		for j := 0; j < innerType.NumField(); j++ {
			name := reflect.Indirect(innerValue).Type().Field(j).Name
			row = append(row, innerValue.FieldByName(name).String())
		}
		if err := wtr.Write(row); err != nil {
			return fmt.Errorf("error writing row: %v", err)
		}
	}
	wtr.Flush()
	if err := wtr.Error(); err != nil {
		return fmt.Errorf("error writing to output: %v", err)
	}
	return nil
}

// EncodeJSON marshalls the Record struct then outputs to Indented JSON
func (e *CSVEncoder) EncodeJSON(v interface{}) error {
	defer timeTrack(time.Now(), "EncodeJSON Func")
	data, err := json.MarshalIndent(v, "  ", "  ")
	if err != nil {
		return fmt.Errorf("error encoding to JSON output: %v", err)
	}
	fmt.Fprintln(e.output, string(data))
	return nil
}
