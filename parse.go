package csvparse

import (
	"encoding/csv"
	"fmt"
	"io"
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

// Parser holds the header field map and reader
type Parser struct {
	header map[string]int
	file   io.ReadCloser
}

// New initializes a new parser
func New(input io.ReadCloser) *Parser {
	return &Parser{
		header: map[string]int{},
		file:   input,
	}
}

// UnMarshalCSV unmarshalls CSV file to record struct
func (p *Parser) UnMarshalCSV() ([]*Record, error) {
	records := []*Record{}
	rdr := csv.NewReader(p.file)

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
			reqFields := []string{"firstname", "lastname", "address1", "city", "state", "zip"}
			for _, v := range reqFields {
				if _, ok := p.header[v]; ok != true {
					return nil, fmt.Errorf("%v : Missing required header fields [firstname, lastname, address1, city, state, zip]", v)
				}
			}
			// continue
		}

		// initialize new record instance then unmarshall records
		r := &Record{}
		for header := range p.header {
			switch header {
			case "fullname":
				r.Fullname = tCase(row[p.header[header]])
			case "firstname":
				r.Firstname = tCase(row[p.header[header]])
			case "mi":
				r.MI = uCase(row[p.header[header]])
			case "lastname":
				r.Lastname = tCase(row[p.header[header]])
			case "address1":
				r.Address1 = tCase(row[p.header[header]])
			case "address2":
				r.Address2 = tCase(row[p.header[header]])
			case "city":
				r.City = tCase(row[p.header[header]])
			case "state":
				r.State = uCase(row[p.header[header]])
			case "zip":
				r.Zip = row[p.header[header]]
			case "zip4":
				r.Zip4 = row[p.header[header]]
			case "hph":
				r.HPH = formatPhone(row[p.header[header]])
			case "bph":
				r.BPH = formatPhone(row[p.header[header]])
			case "cph":
				r.CPH = formatPhone(row[p.header[header]])
			case "email":
				r.Email = lCase(row[p.header[header]])
			case "vin":
				r.VIN = uCase(row[p.header[header]])
			case "year":
				r.Year = row[p.header[header]]
			case "make":
				r.Make = tCase(row[p.header[header]])
			case "model":
				r.Model = tCase(row[p.header[header]])
			case "deldate":
				r.DelDate = parseDate(row[p.header[header]])
			case "date":
				r.Date = parseDate(row[p.header[header]])
			case "dsfwalkseq":
				r.DSFwalkseq = stripSep(row[p.header[header]])
			case "crrt":
				r.CRRT = stripSep(row[p.header[header]])
			case "kbb":
				r.KBB = stripSep(row[p.header[header]])
			}
		}

		// validate Fullname, First Name, Last Name & MI
		if r.Fullname != "" && (r.Firstname == "" || r.Lastname == "") {
			name := names.Parse(r.Fullname)
			r.Firstname = tCase(name.FirstName)
			r.MI = uCase(name.MiddleName)
			r.Lastname = tCase(name.LastName)
		} else {
			r.Firstname = tCase(r.Firstname)
			r.MI = uCase(r.MI)
			r.Lastname = tCase(r.Lastname)
		}

		// parse Zip code field into Zip & Zip4
		zip, zip4 := parseZip(r.Zip)
		r.Zip = zip
		if zip4 != "" {
			r.Zip4 = zip4
		}

		records = append(records, r)
	}
	return records, nil
}

func tCase(f string) string {
	return strings.TrimSpace(strings.Title(strings.ToLower(f)))
}

func uCase(f string) string {
	return strings.TrimSpace(strings.ToUpper(f))
}

func lCase(f string) string {
	return strings.TrimSpace(strings.ToLower(f))
}

func parseZip(zip string) (string, string) {
	switch {
	case regexp.MustCompile(`(?i)^[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]$`).MatchString(zip):
		return trimZeros(zip[:5]), trimZeros(zip[5:])
	case regexp.MustCompile(`(?i)^[0-9][0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9]$`).MatchString(zip):
		zsplit := strings.Split(zip, "-")
		return trimZeros(zsplit[0]), trimZeros(zsplit[1])
	case regexp.MustCompile(`(?i)^[0-9][0-9][0-9][0-9][0-9] [0-9][0-9][0-9][0-9]$`).MatchString(zip):
		zsplit := strings.Split(zip, " ")
		return trimZeros(zsplit[0]), trimZeros(zsplit[1])
	default:
		return zip, ""
	}
}

func trimZeros(s string) string {
	for i := 0; i < len(s); i++ {
		s = strings.TrimPrefix(s, "0")
	}
	return s
}

func formatPhone(p string) string {
	p = stripSep(p)
	switch len(p) {
	case 10:
		return fmt.Sprintf("(%v) %v-%v", p[0:3], p[3:6], p[6:10])
	case 7:
		return fmt.Sprintf("%v-%v", p[0:3], p[3:7])
	default:
		return ""
	}
}

func stripSep(p string) string {
	sep := []string{"'", "#", "%", "$", "-", ".", "*", "(", ")", ":", ";", "{", "}", "|", " "}
	for _, v := range sep {
		p = strings.Replace(p, v, "", -1)
	}
	return p
}

func parseDate(d string) time.Time {
	if d != "" {
		formats := []string{"1/2/2006", "1-2-2006", "1/2/06", "1-2-06",
			"2006/1/2", "2006-1-2", time.RFC3339}
		for _, f := range formats {
			if date, err := time.Parse(f, d); err == nil {
				return date
			}
		}
	}
	return time.Time{}
}
