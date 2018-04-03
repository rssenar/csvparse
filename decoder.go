package csvparse

import (
	"errors"
	"io"
	"reflect"
	"regexp"
	"time"

	"github.com/blendlabs/go-name-parser"
)

// CSVDecoder holds the header field map and reader
type CSVDecoder struct {
	header map[string]int
	file   io.Reader
}

// NewDecoder initializes a new parser
func NewDecoder(input io.Reader) *CSVDecoder {
	return &CSVDecoder{
		header: map[string]int{},
		file:   input,
	}
}

// DecodeCSV unmarshalls CSV file to record struct
func (d *CSVDecoder) DecodeCSV(v interface{}) error {
	defer timeTrack(time.Now(), "DecodeCSV Func")
	csvRows, err := getCSVRows(d.file)
	if err != nil {
		return err
	}
	if len(csvRows) == 0 {
		return errors.New("empty csv file given")
	}
	headerRow := csvRows[0]
	body := csvRows[1:]

	if err := checkForDoubleHeaderNames(headerRow); err != nil {
		return err
	}

	slice := checkifSlice(v)
	if slice.Kind() != reflect.Slice {
		return errors.New("Need to pass in a slice")
	}
	innerType := getInnerSliceType(v)
	if innerType.Kind() != reflect.Struct {
		return errors.New("Need to pass in a slice of stucts")
	}

	innerValueHdr := reflect.New(innerType)
	innerValueHdrLen := innerValueHdr.Elem().NumField()

	for i, csvColumnHdr := range headerRow {
		for j := 0; j < innerValueHdrLen; j++ {
			regex := reflect.Indirect(innerValueHdr).Type().Field(j).Tag.Get("csv")
			if regexp.MustCompile(regex).MatchString(csvColumnHdr) {
				d.header[reflect.Indirect(innerValueHdr).Type().Field(j).Name] = i
			}
		}
	}

	for _, csvRow := range body {
		innerValueRow := reflect.New(innerType)
		innerValueRowLen := innerValueRow.Elem().NumField()

		for j := 0; j < innerValueRowLen; j++ {
			sFName := reflect.Indirect(innerValueRow).Type().Field(j).Name
			switch innerValueRow.Elem().Type().Field(j).Type {
			case reflect.TypeOf(""):
				if _, ok := d.header[sFName]; ok {
					format := reflect.Indirect(innerValueRow).Type().Field(j).Tag.Get("fmt")
					val := formatStringVals(sFName, format, csvRow[d.header[sFName]])
					innerValueRow.Elem().FieldByName(sFName).Set(reflect.ValueOf(val))
				}
			case reflect.TypeOf(time.Now()):
				if _, ok := d.header[sFName]; ok {
					val := ParseDate(csvRow[d.header[sFName]])
					innerValueRow.Elem().FieldByName(sFName).Set(reflect.ValueOf(val))
				}
			default:
				if _, ok := d.header[sFName]; ok {
					val := csvRow[d.header[sFName]]
					innerValueRow.Elem().FieldByName(sFName).Set(reflect.ValueOf(val))
				}
			}
			if j == innerValueRowLen-1 {
				fullname := innerValueRow.Elem().FieldByName("Fullname")
				fname := innerValueRow.Elem().FieldByName("Firstname")
				lname := innerValueRow.Elem().FieldByName("Lastname")
				if fullname.IsValid() && fname.IsValid() && lname.IsValid() {
					if fullname.String() != "" && (fname.String() == "" || lname.String() == "") {
						name := names.Parse(fullname.String())
						innerValueRow.Elem().FieldByName("Firstname").Set(reflect.ValueOf(name.FirstName))
						innerValueRow.Elem().FieldByName("MI").Set(reflect.ValueOf(name.MiddleName))
						innerValueRow.Elem().FieldByName("Lastname").Set(reflect.ValueOf(name.LastName))
					}
				}
				zip := innerValueRow.Elem().FieldByName("Zip")
				if zip.IsValid() {
					zip, zip4 := ParseZip(zip.String())
					innerValueRow.Elem().FieldByName("Zip").Set(reflect.ValueOf(zip))
					if zip4 != "" {
						innerValueRow.Elem().FieldByName("Zip4").Set(reflect.ValueOf(zip4))
					}
				}
			}
		}
		slice.Set(reflect.Append(slice, innerValueRow))
	}
	return nil
}
