package csvparse

import (
	"encoding/csv"
	"fmt"
	"io"
	"regexp"
	"time"
)

type record struct {
	ID         int       `json:"id"`
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
	Date       time.Time `json:"date"`
	DSFwalkseq string    `json:"DSF_Walk_Sequence"`
	CRRT       string    `json:"CRRT"`
	ErrStat    string    `json:"Status"`
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

func (p *Parser) Columns() error {
	rdr := csv.NewReader(p.file)
	rows, err := rdr.ReadAll()
	if err != nil {
		return fmt.Errorf("%v : unable to parse file, csv format required", err)
	}
	for i, v := range rows[:][0] {
		switch {
		case regexp.MustCompile(`(?i)^[Ff]ul[Nn]ame$`).MatchString(v):
			p.header["fullname"] = i
		case regexp.MustCompile(`(?i)^[Ff]irst[Nn]ame|^[Ff]irst [Nn]ame$`).MatchString(v):
			p.header["firstname"] = i
		case regexp.MustCompile(`(?i)^mi$`).MatchString(v):
			p.header["mi"] = i
		case regexp.MustCompile(`(?i)^[Ll]ast[Nn]ame|^[Ll]ast [Nn]ame$`).MatchString(v):
			p.header["lastname"] = i
		case regexp.MustCompile(`(?i)^[Aa]ddress$|^[Aa]ddress.+1$`).MatchString(v):
			p.header["address1"] = i
		case regexp.MustCompile(`(?i)^[Aa]ddress.+2$`).MatchString(v):
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
	for a, b := range rows[:][0] {
		fmt.Println(a, b)
	}
	fmt.Println(p.header)
	return nil
}
