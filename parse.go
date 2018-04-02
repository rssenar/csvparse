package csvparse

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strings"
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

func getInnerSliceType(v interface{}) reflect.Type {
	outType := reflect.TypeOf(v)
	if outType.Kind() == reflect.Ptr {
		outType = outType.Elem()
	}
	inType := outType.Elem()
	if inType.Kind() == reflect.Ptr {
		inType = inType.Elem()
	}
	return inType
}

func checkifSlice(v interface{}) reflect.Value {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return value
}

func getCSVRows(r io.Reader) ([][]string, error) {
	rdr := csv.NewReader(r)
	rows, err := rdr.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%v : unable to read file", err)
	}
	return rows, nil
}

func checkForDoubleHeaderNames(hdrs []string) error {
	headerMap := make(map[string]bool, len(hdrs))
	for _, v := range hdrs {
		if _, ok := headerMap[v]; ok {
			return fmt.Errorf("Repeated header name: %v", v)
		}
		headerMap[v] = true
	}
	return nil
}

func formatStringVals(name, format, val string) string {
	switch format {
	case "tc":
		return TCase(val)
	case "uc":
		return UCase(val)
	case "lc":
		return LCase(val)
	case "fp":
		return FormatPhone(val)
	case "ss":
		return StripSep(val)
	}
	return val
}

// CSVEncoder outputs JSON values to an output stream.
type CSVEncoder struct {
	output io.Writer
}

// NewEncoder initializes a new parser
func NewEncoder(output io.Writer) *CSVEncoder {
	return &CSVEncoder{output: output}
}

// // EncodeCSV marshalls the Record struct then outputs to csv
// func (e *CSVEncoder) EncodeCSV(Records []Record) error {
// 	wtr := csv.NewWriter(e.output)

// 	var header []string
// 	headerLen := reflect.ValueOf(&Records[0]).Elem().NumField()

// 	for i := 0; i < headerLen; i++ {
// 		headerName := reflect.Indirect(reflect.ValueOf(&Records[0])).Type().Field(i).Name
// 		header = append(header, headerName)
// 	}

// 	if err := wtr.Write(header); err != nil {
// 		return fmt.Errorf("error writing header to csv: %v", err)
// 	}

// 	for _, r := range Records {
// 		var row []string
// 		rowLen := reflect.ValueOf(r).NumField()
// 		for i := 0; i < rowLen; i++ {
// 			val := fmt.Sprint(reflect.ValueOf(r).Field(i))
// 			switch reflect.Indirect(reflect.ValueOf(Records[0])).Type().Field(i).Name {
// 			case "DelDate", "Date":
// 				if !r.DelDate.IsZero() {
// 					val = fmt.Sprintf("%v/%v/%v", int(r.DelDate.Month()), r.DelDate.Day(), r.DelDate.Year())
// 				} else {
// 					val = ""
// 				}
// 			}
// 			row = append(row, val)
// 		}
// 		if err := wtr.Write(row); err != nil {
// 			return fmt.Errorf("error writing row to csv: %v", err)
// 		}
// 	}
// 	wtr.Flush()
// 	if err := wtr.Error(); err != nil {
// 		return fmt.Errorf("error writing to output: %v", err)
// 	}
// 	return nil
// }

// EncodeJSON marshalls the Record struct then outputs to Indented JSON
func (e *CSVEncoder) EncodeJSON(v interface{}) error {
	data, err := json.MarshalIndent(v, " ", " ")
	if err != nil {
		return fmt.Errorf("error encoding to JSON output: %v", err)
	}
	fmt.Fprintln(e.output, string(data))
	return nil
}

// TCase transforms string to title case and trims leading & trailing white space
func TCase(f string) string {
	return strings.TrimSpace(strings.Title(strings.ToLower(f)))
}

// UCase transforms string to upper case and trims leading & trailing white space
func UCase(f string) string {
	return strings.TrimSpace(strings.ToUpper(f))
}

// LCase transforms string to lower case and trims leading & trailing white space
func LCase(f string) string {
	return strings.TrimSpace(strings.ToLower(f))
}

// ParseZip perses ZIP code to Zip & Zip4
func ParseZip(zip string) (string, string) {
	if zip == "" {
		return "", ""
	}
	switch {
	case regexp.MustCompile(`(?i)^[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]$`).MatchString(zip):
		return TrimZeros(zip[:5]), TrimZeros(zip[5:])
	case regexp.MustCompile(`(?i)^[0-9][0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9]$`).MatchString(zip):
		zsplit := strings.Split(zip, "-")
		return TrimZeros(zsplit[0]), TrimZeros(zsplit[1])
	case regexp.MustCompile(`(?i)^[0-9][0-9][0-9][0-9][0-9] [0-9][0-9][0-9][0-9]$`).MatchString(zip):
		zsplit := strings.Split(zip, " ")
		return TrimZeros(zsplit[0]), TrimZeros(zsplit[1])
	default:
		return zip, ""
	}
}

// TrimZeros removed leading Zeros
func TrimZeros(s string) string {
	l := len(s)
	for i := 1; i <= l; i++ {
		s = strings.TrimPrefix(s, "0")
	}
	return s
}

// FormatPhone re-formats phone field
func FormatPhone(p string) string {
	p = StripSep(p)
	switch len(p) {
	case 10:
		return fmt.Sprintf("(%v) %v-%v", p[0:3], p[3:6], p[6:10])
	case 7:
		return fmt.Sprintf("%v-%v", p[0:3], p[3:7])
	default:
		return ""
	}
}

// StripSep removes irrelevant characters
func StripSep(p string) string {
	sep := []string{"'", "#", "%", "$", "-", "+", ".", "*", "(", ")", ":", ";", "{", "}", "|", "&", " "}
	for _, v := range sep {
		p = strings.Replace(p, v, "", -1)
	}
	return p
}

// ParseDate converts date string to time.Time
func ParseDate(d string) time.Time {
	if d != "" {
		formats := []string{"1/2/2006", "1-2-2006", "1/2/06", "1-2-06", "2006/1/2", "2006-1-2", time.RFC3339}
		for _, f := range formats {
			if date, err := time.Parse(f, d); err == nil {
				return date
			}
		}
	}
	return time.Time{}
}
