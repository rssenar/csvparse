package csvparse

import (
	"encoding/csv"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type Record struct {
	ID         int    `json:"id"`
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

func (p *Parser) UnMarshalCSV() (Records []Record, err error) {
	rdr := csv.NewReader(p.file)
	rows, err := rdr.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%v : unable to parse file, csv format required", err)
	}
	for i, v := range rows[:][0] {
		switch {
		case regexp.MustCompile(`(?i)^[Ff]ull[Nn]ame$`).MatchString(v):
			p.header["fullname"] = i
		case regexp.MustCompile(`(?i)^[Ff]irst[Nn]ame|^[Ff]irst [Nn]ame$`).MatchString(v):
			p.header["firstname"] = i
		case regexp.MustCompile(`(?i)^mi$`).MatchString(v):
			p.header["mi"] = i
		case regexp.MustCompile(`(?i)^[Ll]ast[Nn]ame|^[Ll]ast [Nn]ame$`).MatchString(v):
			p.header["lastname"] = i
		case regexp.MustCompile(`(?i)^[Aa]ddress1?$^|[Aa]ddress[ _]1?$`).MatchString(v):
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

	// Check header for required fields
	reqFields := []string{"firstname", "lastname", "address1", "city", "state", "zip"}
	for _, v := range reqFields {
		if _, ok := p.header[v]; ok != true {
			return nil, fmt.Errorf("%v : Missing required fields [firstname, lastname, address1, city, state, zip]", v)
		}
	}

	var r Record
	for _, rec := range rows[:][1:] {
		for hdr := range p.header {
			switch hdr {
			case "fullname":
				r.Fullname = lCase(rec[p.header[hdr]])
			case "firstname":
				r.Firstname = lCase(rec[p.header[hdr]])
			case "mi":
				r.MI = uCase(rec[p.header[hdr]])
			case "lastname":
				r.Lastname = lCase(rec[p.header[hdr]])
			case "address1":
				r.Address1 = lCase(rec[p.header[hdr]])
			case "address2":
				r.Address2 = lCase(rec[p.header[hdr]])
			case "city":
				r.City = lCase(rec[p.header[hdr]])
			case "state":
				r.State = lCase(rec[p.header[hdr]])
			case "zip":
				r.Zip = rec[p.header[hdr]]
			case "zip4":
				r.Zip4 = rec[p.header[hdr]]
			case "hph":
				r.HPH = rec[p.header[hdr]]
			case "bph":
				r.BPH = rec[p.header[hdr]]
			case "cph":
				r.CPH = rec[p.header[hdr]]
			case "email":
				r.Email = lCase(rec[p.header[hdr]])
			case "vin":
				r.VIN = uCase(rec[p.header[hdr]])
			case "year":
				r.Year = rec[p.header[hdr]]
			case "make":
				r.Make = lCase(rec[p.header[hdr]])
			case "model":
				r.Model = lCase(rec[p.header[hdr]])
			case "deldate":
				r.DelDate = rec[p.header[hdr]]
			case "date":
				r.Date = rec[p.header[hdr]]
			case "dsfwalkseq":
				r.DSFwalkseq = uCase(rec[p.header[hdr]])
			case "crrt":
				r.CRRT = rec[p.header[hdr]]
			case "kbb":
				r.KBB = rec[p.header[hdr]]
			}
		}
		Records = append(Records, r)
	}

	// for _, row := range rows[:][1:] {
	// 	w := csv.NewWriter(os.Stdout)
	// 	if err := w.Write(row); err != nil {
	// 		return fmt.Errorf("%v : error writing record to csv", err)
	// 	}
	// 	// Write any buffered data to the underlying writer (standard output).
	// 	w.Flush()

	// 	if err := w.Error(); err != nil {
	// 		return fmt.Errorf("%v : error writing to stdout", err)
	// 	}
	// }

	// for a, b := range rows[:][0] {
	// 	fmt.Println(a, b)
	// }

	// type kv struct {
	// 	Key   string
	// 	Value int
	// }

	// var ss []kv
	// for k, v := range p.header {
	// 	ss = append(ss, kv{k, v})
	// }

	// sort.Slice(ss, func(i, j int) bool {
	// 	return ss[i].Value < ss[j].Value
	// })

	// for _, str := range ss {
	// 	fmt.Printf("%s, %d\n", str.Key, str.Value)
	// }
	// fmt.Printf("%v", x)
	return
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
