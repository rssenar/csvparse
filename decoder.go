package csvparse

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/blendlabs/go-name-parser"
)

// CSVDecoder holds the header field map and io.reader interface
type CSVDecoder struct {
	header map[string]int
	file   io.Reader
}

// NewDecoder allocates a new instance of CSVDecoder
func NewDecoder(input io.Reader) *CSVDecoder {
	return &CSVDecoder{
		header: map[string]int{},
		file:   input,
	}
}

// DecodeCSV unmarshalls CSV file to a specified struct type
func (d *CSVDecoder) DecodeCSV(v interface{}) error {
	// Optional timer function for determining function duration
	// defer timeTrack(time.Now(), "DecodeCSVtoStruct")

	// getCSVRows grabs [][]strings from spcified input
	csvRows, err := GetCSVRows(d.file)
	if err != nil {
		return err
	}

	// returns an error if an empty file was provided
	if len(csvRows) == 0 {
		return errors.New("empty csv file")
	}
	// grab Header row
	headerRow := csvRows[0]
	// grab get remaining rows as body
	body := csvRows[1:]

	// check header row for duplicate fields
	// if duplicate fields found, return error
	if err := CheckForDoubleHeaderNames(headerRow); err != nil {
		return err
	}

	// check interface type (v)
	// if type is not a slice, return error
	slice := CheckInterfaceValue(v)
	if slice.Kind() != reflect.Slice {
		return errors.New("Only slice types are permited")
	}

	// check inner slice type
	// if type is not a struct, return error
	innerType := GetInnerSliceType(v)
	if innerType.Kind() != reflect.Struct {
		return errors.New("Only slice of stucts type permited")
	}

	// allocate new reflect stuct instance from inner type
	innerValueHdr := reflect.New(innerType)

	// range over header row
	for i, csvColumnHdr := range headerRow {
		// range over struct fields
		for j := 0; j < innerValueHdr.Elem().NumField(); j++ {
			var regex string
			// grab regex string from struct tag
			// if not struct tags provided, struct name will be used
			if rgx, ok := reflect.Indirect(innerValueHdr).Type().Field(j).Tag.Lookup("csv"); ok {
				switch rgx {
				case "-":
					regex = reflect.Indirect(innerValueHdr).Type().Field(j).Name
				default:
					regex = rgx
				}
			} else {
				continue
			}
			if regexp.MustCompile(regex).MatchString(csvColumnHdr) {
				d.header[reflect.Indirect(innerValueHdr).Type().Field(j).Name] = i
			}
		}
	}

	for _, csvRow := range body {
		innerValueRow := reflect.New(innerType)
		for j := 0; j < innerValueRow.Elem().NumField(); j++ {
			sFName := reflect.Indirect(innerValueRow).Type().Field(j).Name
			switch innerValueRow.Elem().Type().Field(j).Type {
			case reflect.TypeOf(time.Now()):
				if _, ok := d.header[sFName]; ok {
					val := ParseDate(csvRow[d.header[sFName]])
					innerValueRow.Elem().FieldByName(sFName).Set(reflect.ValueOf(val))
				}
			default:
				if _, ok := d.header[sFName]; ok {
					if format, ok := reflect.Indirect(innerValueRow).Type().Field(j).Tag.Lookup("fmt"); ok {
						if format != "-" {
							val := csvRow[d.header[sFName]]
							fmtvalue, err := FormatStringVals(format, val)
							if err != nil {
								log.Fatalln(err)
							}
							innerValueRow.Elem().FieldByName(sFName).SetString(fmtvalue)
						} else {
							val := csvRow[d.header[sFName]]
							innerValueRow.Elem().FieldByName(sFName).SetString(val)
						}
					}
				}
			}
			if j == innerValueRow.Elem().NumField()-1 {
				FullN := innerValueRow.Elem().FieldByName("Fullname")
				FN := innerValueRow.Elem().FieldByName("Firstname")
				LN := innerValueRow.Elem().FieldByName("Lastname")
				if FullN.String() != "" && (FN.String() == "" || LN.String() == "") {
					name := names.Parse(FullN.String())
					innerValueRow.Elem().FieldByName("Firstname").SetString(name.FirstName)
					innerValueRow.Elem().FieldByName("MI").SetString(name.MiddleName)
					innerValueRow.Elem().FieldByName("Lastname").SetString(name.LastName)
				}
				if fmt.Sprint(innerValueRow.Elem().FieldByName("Zip")) != "" {
					zip, zip4 := ParseZip(fmt.Sprint(innerValueRow.Elem().FieldByName("Zip")))
					innerValueRow.Elem().FieldByName("Zip").SetString(zip)
					if innerValueRow.Elem().FieldByName("Zp4").String() != "" {
						innerValueRow.Elem().FieldByName("Zip4").SetString(zip4)
					}
				}
			}
		}
		slice.Set(reflect.Append(slice, innerValueRow))
	}
	return nil
}

// GetCSVRows get [][]strins from io.reader
func GetCSVRows(r io.Reader) ([][]string, error) {
	rdr := csv.NewReader(r)
	rows, err := rdr.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%v : can only parse CSV files", err)
	}
	return rows, nil
}

// GetInnerSliceType gets inner slice tyoe with reflect and return non-pointer value
func GetInnerSliceType(v interface{}) reflect.Type {
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

// CheckInterfaceValue get interface tyoe with reflect and return non-pointer value
func CheckInterfaceValue(v interface{}) reflect.Value {
	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	return value
}

// CheckForDoubleHeaderNames checks for duplicate header fields
// returns error if dupes found
func CheckForDoubleHeaderNames(hdrs []string) error {
	headerMap := make(map[string]bool, len(hdrs))
	for _, v := range hdrs {
		if _, ok := headerMap[v]; ok {
			return fmt.Errorf("Repeated header name: %v", v)
		}
		headerMap[v] = true
	}
	return nil
}

// FormatStringVals applies formating to string
// based on "fmt" struct tag
func FormatStringVals(format, val string) (string, error) {
	switch format {
	case "tc":
		return TCase(val), nil
	case "uc":
		return UCase(val), nil
	case "lc":
		return LCase(val), nil
	case "fp":
		return FormatPhone(val), nil
	case "ss":
		return StripSep(val), nil
	default:
		return "", errors.New("Invalid string format: use [tc, uc, lc, fp, ss]")
	}
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

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
