package csvparse

import (
	"encoding/csv"
	"encoding/json"
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
	PKey       int       `json:"Primary_Key"`
	Fullname   string    `json:"full_name"`
	Firstname  string    `json:"first_name"`
	MI         string    `json:"middle_name"`
	Lastname   string    `json:"last_name"`
	Address1   string    `json:"address_1"`
	Address2   string    `json:"address_2"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	Zip        string    `json:"zip"`
	Zip4       string    `json:"zip_4"`
	HPH        string    `json:"home_phone"`
	BPH        string    `json:"business_phone"`
	CPH        string    `json:"mobile_phone"`
	Email      string    `json:"email"`
	VIN        string    `json:"VIN"`
	Year       string    `json:"year"`
	Make       string    `json:"make"`
	Model      string    `json:"model"`
	DelDate    time.Time `json:"delivery_date"`
	Date       time.Time `json:"last_service_date"`
	DSFwalkseq string    `json:"DSF_Walk_Sequence"`
	CRRT       string    `json:"CRRT"`
	KBB        string    `json:"KBB"`
}

// CSVDecoder holds the header field map and reader
type CSVDecoder struct {
	header map[string]int
	file   io.Reader
	vflag  *bool
}

// NewDecoder initializes a new parser
func NewDecoder(input io.Reader, flag *bool) *CSVDecoder {
	return &CSVDecoder{
		header: map[string]int{},
		file:   input,
		vflag:  flag,
	}
}

// DecodeCSV unmarshalls CSV file to record struct
func (d *CSVDecoder) DecodeCSV() ([]*Record, error) {
	var Records []*Record

	rdr := csv.NewReader(d.file)
	for i := 0; ; i++ {
		row, err := rdr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("%v : unable to parse file, csv format required", err)
		}
		hdrmp := make(map[string]int)
		if i == 0 {
			for i, v := range row {
				if _, ok := hdrmp[v]; ok == true {
					return nil, fmt.Errorf("%v : Duplicate header field", v)
				}
				hdrmp[v]++

				switch {
				case regexp.MustCompile(`(?i)^[Ff]ull[Nn]ame$`).MatchString(v):
					d.header["fullname"] = i
				case regexp.MustCompile(`(?i)^[Ff]irst[Nn]ame$|^[Ff]irst [Nn]ame$`).MatchString(v):
					d.header["firstname"] = i
				case regexp.MustCompile(`(?i)^mi$`).MatchString(v):
					d.header["mi"] = i
				case regexp.MustCompile(`(?i)^[Ll]ast[Nn]ame$|^[Ll]ast [Nn]ame$`).MatchString(v):
					d.header["lastname"] = i
				case regexp.MustCompile(`(?i)^[Aa]ddress1?$|^[Aa]ddress[ _-]1?$`).MatchString(v):
					d.header["address1"] = i
				case regexp.MustCompile(`(?i)^[Aa]ddress2$|^[Aa]ddress[ _-]2$`).MatchString(v):
					d.header["address2"] = i
				case regexp.MustCompile(`(?i)^[Cc]ity$`).MatchString(v):
					d.header["city"] = i
				case regexp.MustCompile(`(?i)^[Ss]tate$|^[Ss][Tt]$`).MatchString(v):
					d.header["state"] = i
				case regexp.MustCompile(`(?i)^[Zz]ip$`).MatchString(v):
					d.header["zip"] = i
				case regexp.MustCompile(`(?i)^[Zz]ip4$|^4zip$`).MatchString(v):
					d.header["zip4"] = i
				case regexp.MustCompile(`(?i)^hph$`).MatchString(v):
					d.header["hph"] = i
				case regexp.MustCompile(`(?i)^bph$`).MatchString(v):
					d.header["bph"] = i
				case regexp.MustCompile(`(?i)^cph$`).MatchString(v):
					d.header["cph"] = i
				case regexp.MustCompile(`(?i)^[Ee]mail$`).MatchString(v):
					d.header["email"] = i
				case regexp.MustCompile(`(?i)^[Vv]in$`).MatchString(v):
					d.header["vin"] = i
				case regexp.MustCompile(`(?i)^[Yy]ear$|^[Vv]yr$`).MatchString(v):
					d.header["year"] = i
				case regexp.MustCompile(`(?i)^[Mm]ake$|^[Vv]mk$`).MatchString(v):
					d.header["make"] = i
				case regexp.MustCompile(`(?i)^[Mm]odel$|^[Vv]md$`).MatchString(v):
					d.header["model"] = i
				case regexp.MustCompile(`(?i)^[Dd]eldate$`).MatchString(v):
					d.header["deldate"] = i
				case regexp.MustCompile(`(?i)^[Dd]ate$`).MatchString(v):
					d.header["date"] = i
				case regexp.MustCompile(`(?i)^DSF_WALK_SEQ$`).MatchString(v):
					d.header["dsfwalkseq"] = i
				case regexp.MustCompile(`(?i)^[Cc]rrt$`).MatchString(v):
					d.header["crrt"] = i
				case regexp.MustCompile(`(?i)^KBB$`).MatchString(v):
					d.header["kbb"] = i
				}
			}

			// Check that all required fields are present, this is inactive by defaut
			// you can activate with command line flag -v
			if *d.vflag {
				reqFields := []string{"firstname", "lastname", "address1", "city", "state", "zip"}
				for _, v := range reqFields {
					if _, ok := d.header[v]; ok != true {
						return nil, fmt.Errorf("Error : Missing [ %v ] - required header field", v)
					}
				}
			}
			continue
		}

		// initialize new record instance then unmarshall records
		r := Record{}
		for header := range d.header {
			switch header {
			case "fullname":
				r.Fullname = TCase(row[d.header[header]])
			case "firstname":
				r.Firstname = TCase(row[d.header[header]])
			case "mi":
				r.MI = UCase(row[d.header[header]])
			case "lastname":
				r.Lastname = TCase(row[d.header[header]])
			case "address1":
				r.Address1 = TCase(row[d.header[header]])
			case "address2":
				r.Address2 = TCase(row[d.header[header]])
			case "city":
				r.City = TCase(row[d.header[header]])
			case "state":
				r.State = UCase(row[d.header[header]])
			case "zip":
				r.Zip = row[d.header[header]]
			case "zip4":
				r.Zip4 = row[d.header[header]]
			case "hph":
				r.HPH = FormatPhone(row[d.header[header]])
			case "bph":
				r.BPH = FormatPhone(row[d.header[header]])
			case "cph":
				r.CPH = FormatPhone(row[d.header[header]])
			case "email":
				r.Email = LCase(row[d.header[header]])
			case "vin":
				r.VIN = UCase(row[d.header[header]])
			case "year":
				r.Year = row[d.header[header]]
			case "make":
				r.Make = TCase(row[d.header[header]])
			case "model":
				r.Model = TCase(row[d.header[header]])
			case "deldate":
				r.DelDate = ParseDate(row[d.header[header]])
			case "date":
				r.Date = ParseDate(row[d.header[header]])
			case "dsfwalkseq":
				r.DSFwalkseq = StripSep(row[d.header[header]])
			case "crrt":
				r.CRRT = StripSep(row[d.header[header]])
			case "kbb":
				r.KBB = StripSep(row[d.header[header]])
			}
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

		Records = append(Records, &r)
	}
	return Records, nil
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
func (e *CSVEncoder) EncodeCSV(Records []*Record, outfields []string) error {

	wtr := csv.NewWriter(e.output)

	var header []string
	headerLen := reflect.ValueOf(Records[0]).Elem().NumField()
	for i := 0; i < headerLen; i++ {
		headerName := reflect.Indirect(reflect.ValueOf(Records[i])).Type().Field(i).Name
		if outfields == nil {
			header = append(header, headerName)
		} else {
			for _, o := range outfields {
				if o == headerName {
					header = append(header, headerName)
				}
			}
		}
	}
	if err := wtr.Write(header); err != nil {
		return fmt.Errorf("error writing header to csv: %v", err)
	}

	for _, r := range Records {
		var row []string
		rowLen := reflect.ValueOf(r).Elem().NumField()
		for i := 0; i < rowLen; i++ {
			name := reflect.Indirect(reflect.ValueOf(Records[i])).Type().Field(i).Name
			val := fmt.Sprint(reflect.ValueOf(r).Elem().Field(i))

			if outfields == nil {
				switch name {
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
			} else {
				for _, o := range outfields {
					if o == name {
						switch name {
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
				}
			}
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
func (e *CSVEncoder) EncodeJSON(Records []*Record) error {
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
