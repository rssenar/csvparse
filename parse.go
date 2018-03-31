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

// Record represents a customer
type Record struct {
	Fullname   string    `json:"Full_name" csv:"(?i)^fullname$"`
	Firstname  string    `json:"First_name" csv:"(?i)^first[ _-]?name$"`
	MI         string    `json:"Middle_name" csv:"(?i)^mi$"`
	Lastname   string    `json:"Last_name" csv:"(?i)^last[ _-]?name$"`
	Address1   string    `json:"Address_1" csv:"(?i)^address[ _-]?1?$"`
	Address2   string    `json:"Address_2" csv:"(?i)^address[ _-]?2$"`
	City       string    `json:"City" csv:"(?i)^[Cc]ity$"`
	State      string    `json:"State" csv:"(?i)^state$|^st$"`
	Zip        string    `json:"Zip" csv:"(?i)^zip$"`
	Zip4       string    `json:"Zip_4" csv:"(?i)^zip4$|^4zip$"`
	HPH        string    `json:"Home_phone" csv:"(?i)^hph$"`
	BPH        string    `json:"Business_phone" csv:"(?i)^bph$"`
	CPH        string    `json:"Mobile_phone" csv:"(?i)^cph$"`
	Email      string    `json:"Email" csv:"(?i)^email$"`
	VIN        string    `json:"VIN" csv:"(?i)^vin$"`
	Year       string    `json:"Year" csv:"(?i)^year$|^vyr$"`
	Make       string    `json:"Make" csv:"(?i)^make$|^vmk$"`
	Model      string    `json:"Model" csv:"(?i)^model$|^vmd$"`
	DelDate    time.Time `json:"Delivery_date" csv:"(?i)^del[ ]?date[s]?$"`
	Date       time.Time `json:"Last_service_date" csv:"(?i)^date[s]?$"`
	DSFwalkseq string    `json:"DSF_Walk_Sequence" csv:"(?i)^DSF_WALK_SEQ$"`
	CRRT       string    `json:"CRRT" csv:"(?i)^crrt$"`
	KBB        string    `json:"KBB" csv:"(?i)^kbb$"`
}

// CSVDecoder holds the header field map and reader
type CSVDecoder struct {
	header  map[string]int
	file    io.Reader
	records []Record
}

// NewDecoder initializes a new parser
func NewDecoder(input io.Reader) *CSVDecoder {
	return &CSVDecoder{
		header:  map[string]int{},
		file:    input,
		records: []Record{},
	}
}

// DecodeCSV unmarshalls CSV file to record struct
func (d *CSVDecoder) DecodeCSV() ([]Record, error) {
	csvRows, err := getCSVRows(d.file)
	if err != nil {
		return nil, err
	}
	if len(csvRows) == 0 {
		return nil, errors.New("empty csv file given")
	}
	headerRow := csvRows[0]
	body := csvRows[1:]

	if err := checkForDoubleHeaderNames(headerRow); err != nil {
		return nil, err
	}

	r := Record{}
	sValue := reflect.ValueOf(&r)
	sLen := sValue.Elem().NumField()

	for i, csvColumnHdr := range headerRow {
		for j := 0; j < sLen; j++ {
			if regexp.MustCompile(reflect.Indirect(sValue).Type().Field(j).Tag.Get("csv")).MatchString(csvColumnHdr) {
				d.header[reflect.Indirect(sValue).Type().Field(j).Name] = i
			}
		}
	}

	for _, csvRow := range body {
		for j := 0; j < sLen; j++ {
			sFName := reflect.Indirect(sValue).Type().Field(j).Name
			switch sValue.Elem().Field(j).Type() {
			case reflect.TypeOf(""):
				if _, ok := d.header[sFName]; ok {
					val := reformatStringVals(sFName, csvRow[d.header[sFName]])
					sValue.Elem().FieldByName(sFName).Set(reflect.ValueOf(val))
				}
			case reflect.TypeOf(time.Now()):
				if _, ok := d.header[sFName]; ok {
					val := ParseDate(csvRow[d.header[sFName]])
					sValue.Elem().FieldByName(sFName).Set(reflect.ValueOf(val))
				}
			}
		}
		d.records = append(d.records, r)
	}

	// validate Fullname, First Name, Last Name & MI
	if r.Fullname != "" && (r.Firstname == "" || r.Lastname == "") {
		name := names.Parse(r.Fullname)
		r.Firstname = TCase(name.FirstName)
		r.MI = UCase(name.MiddleName)
		r.Lastname = TCase(name.LastName)
	} else {
		r.Firstname = TCase(r.Firstname)
		r.MI = UCase(r.MI)
		r.Lastname = TCase(r.Lastname)
	}

	// parse Zip code field into Zip & Zip4
	zip, zip4 := ParseZip(r.Zip)
	r.Zip = zip
	if zip4 != "" {
		r.Zip4 = zip4
	}
	return d.records, nil
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

func reformatStringVals(sname, val string) string {
	switch sname {
	case "Fullname":
		return TCase(val)
	case "Firstname":
		return TCase(val)
	case "MI":
		return UCase(val)
	case "Lastname":
		return TCase(val)
	case "Address1":
		return TCase(val)
	case "Address2":
		return TCase(val)
	case "City":
		return TCase(val)
	case "State":
		return UCase(val)
	case "Zip":
		return UCase(val)
	case "Zip4":
		return UCase(val)
	case "HPH":
		return FormatPhone(val)
	case "BPH":
		return FormatPhone(val)
	case "CPH":
		return FormatPhone(val)
	case "Email":
		return LCase(val)
	case "VIN":
		return UCase(val)
	case "Year":
		return UCase(val)
	case "Make":
		return TCase(val)
	case "Model":
		return TCase(val)
	case "DSFwalkseq":
		return StripSep(val)
	case "CRRT":
		return StripSep(val)
	case "KBB":
		return UCase(val)
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

// EncodeCSV marshalls the Record struct then outputs to csv
func (e *CSVEncoder) EncodeCSV(Records []Record) error {
	wtr := csv.NewWriter(e.output)

	var header []string
	headerLen := reflect.ValueOf(&Records[0]).Elem().NumField()

	for i := 0; i < headerLen; i++ {
		headerName := reflect.Indirect(reflect.ValueOf(&Records[0])).Type().Field(i).Name
		header = append(header, headerName)
	}

	if err := wtr.Write(header); err != nil {
		return fmt.Errorf("error writing header to csv: %v", err)
	}

	for _, r := range Records {
		var row []string

		rowLen := reflect.ValueOf(r).NumField()

		for i := 0; i < rowLen; i++ {
			fName := reflect.Indirect(reflect.ValueOf(Records[0])).Type().Field(i).Name
			val := fmt.Sprint(reflect.ValueOf(r).Field(i))

			switch fName {
			case "DelDate":
				if !r.DelDate.IsZero() {
					val = fmt.Sprintf("%v/%v/%v", int(r.DelDate.Month()), r.DelDate.Day(), r.DelDate.Year())
				} else {
					val = ""
				}
			case "Date":
				if !r.Date.IsZero() {
					val = fmt.Sprintf("%v/%v/%v", int(r.DelDate.Month()), r.DelDate.Day(), r.DelDate.Year())
				} else {
					val = ""
				}
			}
			row = append(row, val)
		}
		if err := wtr.Write(row); err != nil {
			return fmt.Errorf("error writing row to csv: %v", err)
		}
	}
	wtr.Flush()
	if err := wtr.Error(); err != nil {
		return fmt.Errorf("error writing to output: %v", err)
	}
	return nil
}

// EncodeJSON marshalls the Record struct then outputs to Indented JSON
func (e *CSVEncoder) EncodeJSON(Records []Record) error {
	data, err := json.MarshalIndent(Records, " ", " ")
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
