package csvparse

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
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

// Parser holds the header field map and reader
type Parser struct {
	header  map[string]int
	file    io.Reader
	Records []*Record
}

// New initializes a new parser
func New(input io.Reader) *Parser {
	return &Parser{
		header:  map[string]int{},
		file:    input,
		Records: []*Record{},
	}
}

// UnMarshalCSV unmarshalls CSV file to record struct
func (p *Parser) UnMarshalCSV(vh *bool) error {
	rdr := csv.NewReader(p.file)
	for i := 0; ; i++ {
		row, err := rdr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("%v : unable to parse file, csv format required", err)
		}
		hdrmp := make(map[string]int)
		if i == 0 {
			for i, v := range row {
				if _, ok := hdrmp[v]; ok == true {
					return fmt.Errorf("%v : Duplicate header field", v)
				}
				hdrmp[v]++

				switch {
				case regexp.MustCompile(`(?i)^[Ff]ull[Nn]ame$`).MatchString(v):
					p.header["fullname"] = i
				case regexp.MustCompile(`(?i)^[Ff]irst[Nn]ame$|^[Ff]irst [Nn]ame$`).MatchString(v):
					p.header["firstname"] = i
				case regexp.MustCompile(`(?i)^mi$`).MatchString(v):
					p.header["mi"] = i
				case regexp.MustCompile(`(?i)^[Ll]ast[Nn]ame$|^[Ll]ast [Nn]ame$`).MatchString(v):
					p.header["lastname"] = i
				case regexp.MustCompile(`(?i)^[Aa]ddress1?$|^[Aa]ddress[ _-]1?$`).MatchString(v):
					p.header["address1"] = i
				case regexp.MustCompile(`(?i)^[Aa]ddress2$|^[Aa]ddress[ _-]2$`).MatchString(v):
					p.header["address2"] = i
				case regexp.MustCompile(`(?i)^[Cc]ity$`).MatchString(v):
					p.header["city"] = i
				case regexp.MustCompile(`(?i)^[Ss]tate$|^[Ss][Tt]$`).MatchString(v):
					p.header["state"] = i
				case regexp.MustCompile(`(?i)^[Zz]ip$`).MatchString(v):
					p.header["zip"] = i
				case regexp.MustCompile(`(?i)^[Zz]ip4$|^4zip$`).MatchString(v):
					p.header["zip4"] = i
				case regexp.MustCompile(`(?i)^hph$`).MatchString(v):
					p.header["hph"] = i
				case regexp.MustCompile(`(?i)^bph$`).MatchString(v):
					p.header["bph"] = i
				case regexp.MustCompile(`(?i)^cph$`).MatchString(v):
					p.header["cph"] = i
				case regexp.MustCompile(`(?i)^[Ee]mail$`).MatchString(v):
					p.header["email"] = i
				case regexp.MustCompile(`(?i)^[Vv]in$`).MatchString(v):
					p.header["vin"] = i
				case regexp.MustCompile(`(?i)^[Yy]ear$|^[Vv]yr$`).MatchString(v):
					p.header["year"] = i
				case regexp.MustCompile(`(?i)^[Mm]ake$|^[Vv]mk$`).MatchString(v):
					p.header["make"] = i
				case regexp.MustCompile(`(?i)^[Mm]odel$|^[Vv]md$`).MatchString(v):
					p.header["model"] = i
				case regexp.MustCompile(`(?i)^[Dd]eldate$`).MatchString(v):
					p.header["deldate"] = i
				case regexp.MustCompile(`(?i)^[Dd]ate$`).MatchString(v):
					p.header["date"] = i
				case regexp.MustCompile(`(?i)^DSF_WALK_SEQ$`).MatchString(v):
					p.header["dsfwalkseq"] = i
				case regexp.MustCompile(`(?i)^[Cc]rrt$`).MatchString(v):
					p.header["crrt"] = i
				case regexp.MustCompile(`(?i)^KBB$`).MatchString(v):
					p.header["kbb"] = i
				}
			}

			// Check that all required fields are present
			// you can modify with command line flag -vh=false/true
			if *vh {
				reqFields := []string{"firstname", "lastname", "address1", "city", "state", "zip"}
				for _, v := range reqFields {
					if _, ok := p.header[v]; ok != true {
						return fmt.Errorf("%v : Missing required header field", v)
					}
				}
			}
			continue
		}

		// initialize new record instance then unmarshall records
		r := &Record{}

		for header := range p.header {
			switch header {
			case "fullname":
				r.Fullname = TCase(row[p.header[header]])
			case "firstname":
				r.Firstname = TCase(row[p.header[header]])
			case "mi":
				r.MI = UCase(row[p.header[header]])
			case "lastname":
				r.Lastname = TCase(row[p.header[header]])
			case "address1":
				r.Address1 = TCase(row[p.header[header]])
			case "address2":
				r.Address2 = TCase(row[p.header[header]])
			case "city":
				r.City = TCase(row[p.header[header]])
			case "state":
				r.State = UCase(row[p.header[header]])
			case "zip":
				r.Zip = row[p.header[header]]
			case "zip4":
				r.Zip4 = row[p.header[header]]
			case "hph":
				r.HPH = FormatPhone(row[p.header[header]])
			case "bph":
				r.BPH = FormatPhone(row[p.header[header]])
			case "cph":
				r.CPH = FormatPhone(row[p.header[header]])
			case "email":
				r.Email = LCase(row[p.header[header]])
			case "vin":
				r.VIN = UCase(row[p.header[header]])
			case "year":
				r.Year = row[p.header[header]]
			case "make":
				r.Make = TCase(row[p.header[header]])
			case "model":
				r.Model = TCase(row[p.header[header]])
			case "deldate":
				r.DelDate = ParseDate(row[p.header[header]])
			case "date":
				r.Date = ParseDate(row[p.header[header]])
			case "dsfwalkseq":
				r.DSFwalkseq = StripSep(row[p.header[header]])
			case "crrt":
				r.CRRT = StripSep(row[p.header[header]])
			case "kbb":
				r.KBB = StripSep(row[p.header[header]])
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

		p.Records = append(p.Records, r)
	}
	return nil
}

// MarshaltoCSV marshalls the Record struct then outputs to csv
func (p *Parser) MarshaltoCSV() error {
	hdr := struct {
		PKey, Fullname, Firstname, MI, Lastname, Address1, Address2, City, State, Zip, Zip4,
		HPH, BPH, CPH, Email, VIN, Year, Make, Model, DelDate, Date, DSFwalkseq, CRRT, KBB string
	}{
		"PKey", "Fullname", "Firstname", "MI", "Lastname", "Address1", "Address2", "City", "State", "Zip", "Zip4",
		"HPH", "BPH", "CPH", "Email", "VIN", "Year", "Make", "Model", "DelDate", "Date", "DSFwalkseq", "CRRT", "KBB",
	}
	wtr := csv.NewWriter(os.Stdout)
	for i, r := range p.Records {
		if i == 0 {
			hrow := []string{
				hdr.PKey,
				hdr.Firstname,
				hdr.MI,
				hdr.Lastname,
				hdr.Address1,
				hdr.Address2,
				hdr.City,
				hdr.State,
				hdr.Zip,
				hdr.Zip4,
				hdr.HPH,
				hdr.BPH,
				hdr.CPH,
				hdr.Email,
				hdr.VIN,
				hdr.Year,
				hdr.Make,
				hdr.Model,
				hdr.DelDate,
				hdr.Date,
				hdr.DSFwalkseq,
				hdr.CRRT,
				hdr.KBB,
			}
			if err := wtr.Write(hrow); err != nil {
				return fmt.Errorf("error writing record to csv: %v", err)
			}
			continue
		}

		var DelDate string
		if !r.DelDate.IsZero() {
			DelDate = r.DelDate.Format(time.RFC3339)
		} else {
			DelDate = ""
		}

		var Date string
		if !r.Date.IsZero() {
			Date = r.Date.Format(time.RFC3339)
		} else {
			Date = ""
		}

		row := []string{
			strconv.Itoa(r.PKey),
			r.Firstname,
			r.MI,
			r.Lastname,
			r.Address1,
			r.Address2,
			r.City,
			r.State,
			r.Zip,
			r.Zip4,
			r.HPH,
			r.BPH,
			r.CPH,
			r.Email,
			r.VIN,
			r.Year,
			r.Make,
			r.Model,
			DelDate,
			Date,
			r.DSFwalkseq,
			r.CRRT,
			r.KBB,
		}
		wtr.Write(row)
	}
	wtr.Flush()

	if err := wtr.Error(); err != nil {
		return fmt.Errorf("error writing to output: %v", err)
	}
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
