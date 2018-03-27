package csvparse

import (
	"encoding/csv"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type Record struct {
	Fullname   string `json:"full_name"`
	Firstname  string `json:"first_name"`
	MI         string `json:"middle_name"`
	Lastname   string `json:"last_name"`
	Address1   string `json:"address_1"`
	Address2   string `json:"address_2"`
	City       string `json:"city"`
	State      string `json:"state"`
	Zip        string `json:"zip"`
	Zip4       string `json:"zip_4"`
	HPH        string `json:"home_phone"`
	BPH        string `json:"business_phone"`
	CPH        string `json:"mobile_phone"`
	Email      string `json:"email"`
	VIN        string `json:"VIN"`
	Year       string `json:"year"`
	Make       string `json:"make"`
	Model      string `json:"model"`
	DelDate    string `json:"delivery_date"`
	Date       string `json:"date"`
	DSFwalkseq string `json:"DSF_Walk_Sequence"`
	CRRT       string `json:"CRRT"`
	KBB        string `json:"KBB"`
}

type Parser struct {
	header map[string]int
	file   io.ReadCloser
}

func New(input io.ReadCloser) *Parser {
	return &Parser{
		header: map[string]int{},
		file:   input,
	}
}

func (p *Parser) UnMarshalCSV() ([]Record, error) {

	var records []Record
	rdr := csv.NewReader(p.file)

	for i := 0; ; i++ {
		row, err := rdr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("%v : unable to parse file, csv format required", err)
		}
		if i == 0 {
			for i, v := range row {
				switch {
				case regexp.MustCompile(`(?i)^[Ff]ull[Nn]ame$`).MatchString(v):
					p.header["fullname"] = i
				case regexp.MustCompile(`(?i)^[Ff]irst[Nn]ame|^[Ff]irst [Nn]ame$`).MatchString(v):
					p.header["firstname"] = i
				case regexp.MustCompile(`(?i)^mi$`).MatchString(v):
					p.header["mi"] = i
				case regexp.MustCompile(`(?i)^[Ll]ast[Nn]ame|^[Ll]ast [Nn]ame$`).MatchString(v):
					p.header["lastname"] = i
				case regexp.MustCompile(`(?i)^[Aa]ddress1?$^|[Aa]ddress[ _-]1?$`).MatchString(v):
					p.header["address1"] = i
				case regexp.MustCompile(`(?i)^[Aa]ddress2$|^[Aa]ddress 2$`).MatchString(v):
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
		}

		var r Record
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
				r.HPH = row[p.header[header]]
			case "bph":
				r.BPH = row[p.header[header]]
			case "cph":
				r.CPH = row[p.header[header]]
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
				r.DelDate = row[p.header[header]]
			case "date":
				r.Date = row[p.header[header]]
			case "dsfwalkseq":
				r.DSFwalkseq = uCase(row[p.header[header]])
			case "crrt":
				r.CRRT = row[p.header[header]]
			case "kbb":
				r.KBB = row[p.header[header]]
			}
		}
		records = append(records, r)
	}

	// Check header for required fields
	reqFields := []string{"firstname", "lastname", "address1", "city", "state", "zip"}
	for _, v := range reqFields {
		if _, ok := p.header[v]; ok != true {
			return nil, fmt.Errorf("%v : Missing required header fields [firstname, lastname, address1, city, state, zip]", v)
		}
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
