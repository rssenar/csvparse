package csvparse_test

import (
	"os"
	"testing"
	"time"

	cp "github.com/rssenar/csvparse"
)

func Test_tCase(t *testing.T) {
	cases := []struct {
		input, expected string
	}{
		{" apPlE ", "Apple"},
		{" sUn ", "Sun"},
		{" nOaH  ", "Noah"},
	}
	for _, c := range cases {
		out := cp.TCase(c.input)
		if out != c.expected {
			t.Errorf("Text should be %v but got %v", c.expected, out)
		}
	}
}

func Test_UCase(t *testing.T) {
	cases := []struct {
		input, expected string
	}{
		{" aPPle  ", "APPLE"},
		{"  sUN ", "SUN"},
		{" noAH ", "NOAH"},
	}
	for _, c := range cases {
		out := cp.UCase(c.input)
		if out != c.expected {
			t.Errorf("Text should be %v but got %v", c.expected, out)
		}
	}
}

func Test_LCase(t *testing.T) {
	cases := []struct {
		input, expected string
	}{
		{"   APPLE  ", "apple"},
		{" SUN ", "sun"},
		{"  NOAH   ", "noah"},
	}
	for _, c := range cases {
		out := cp.LCase(c.input)
		if out != c.expected {
			t.Errorf("Text should be %v but got %v", c.expected, out)
		}
	}
}

func Test_parseZip(t *testing.T) {
	cases := []struct {
		input, zip, zip4 string
	}{
		{"92882-1234", "92882", "1234"},
		{"928821234", "92882", "1234"},
		{"928821234", "92882", "1234"},
		{"9288212", "9288212", ""},
		{"92882123456", "92882123456", ""},
	}
	for _, c := range cases {
		zip, zip4 := cp.ParseZip(c.input)
		if zip != c.zip {
			t.Errorf("Zip should be %v but got %v", c.zip, zip)
		}
		if zip4 != c.zip4 {
			t.Errorf("Zip4 should be %v but got %v", c.zip4, zip4)
		}
	}
}

func Test_FormatPhone(t *testing.T) {
	cases := []struct {
		input, expected string
	}{
		{"9493237895", "(949) 323-7895"},
		{"3237895", "323-7895"},
		{"94932", ""},
		{"94932456748912", ""},
	}
	for _, c := range cases {
		out := cp.FormatPhone(c.input)
		if out != c.expected {
			t.Errorf("Phone should be %v but got %v", c.expected, out)
		}
	}
}

func Test_StripSep(t *testing.T) {
	cases := []struct {
		input, expected string
	}{
		{"#$*string&()&", "string"},
		{"#   $*   string   &()&   ", "string"},
	}
	for _, c := range cases {
		out := cp.StripSep(c.input)
		if out != c.expected {
			t.Errorf("Phone should be %v but got %v", c.expected, out)
		}
	}
}

func Test_TrimZeros(t *testing.T) {
	cases := []struct {
		input, expected string
	}{
		{"00000123", "123"},
		{"000000000000000000123", "123"},
		{"0123", "123"},
	}
	for _, c := range cases {
		out := cp.TrimZeros(c.input)
		if out != c.expected {
			t.Errorf("number should be %v but got %v", c.expected, out)
		}
	}
}

func Test_ParseDate(t *testing.T) {
	cases := []struct {
		input    string
		expected time.Time
	}{
		{"12/31/2003", expDate("12/31/2003")},
		{"12-31-2003", expDate("12-31-2003")},
		{"1-3-03", expDate("1-3-03")},
		{"", expDate("")},
	}
	for _, c := range cases {
		out := cp.ParseDate(c.input)
		if out != c.expected {
			t.Errorf("Date should be %v but got %v", c.expected, out)
		}
	}
}

func expDate(date string) time.Time {
	formats := []string{"1/2/2006", "1-2-2006", "1/2/06", "1-2-06", "2006/1/2", "2006-1-2", time.RFC3339}
	for _, f := range formats {
		if date, err := time.Parse(f, date); err == nil {
			return date
		}
	}
	return time.Time{}
}

func Test_UnMarshalCSV(t *testing.T) {
	file, err := os.Open("testfile.csv")
	defer file.Close()
	if err != nil {
		t.Error("Unable to open test file")
	}

	p := cp.New(file)
	parser, err := p.UnMarshalCSV()
	if err != nil {
		t.Error(err)
	}

	expCase := []cp.Record{
		cp.Record{
			Fullname:   "",
			Firstname:  "Mike",
			MI:         "J.",
			Lastname:   "Smith",
			Address1:   "1000 Kelley Dr",
			Address2:   "",
			City:       "Fort Worth",
			State:      "TX",
			Zip:        "76140",
			Zip4:       "3618",
			HPH:        "",
			BPH:        "",
			CPH:        "(682) 227-5578",
			Email:      "msmith@gmail.com",
			VIN:        "4A3AK24F67E006257",
			Year:       "2007",
			Make:       "Mitsubishi",
			Model:      "Eclipse",
			DSFwalkseq: "B425",
			CRRT:       "C003",
			KBB:        "",
		},
		cp.Record{
			Fullname:   "",
			Firstname:  "Adam",
			MI:         "",
			Lastname:   "Savage",
			Address1:   "10 Anywho St",
			Address2:   "",
			City:       "Corona",
			State:      "CA",
			Zip:        "92882",
			Zip4:       "4588",
			HPH:        "(682) 227-5578",
			BPH:        "",
			CPH:        "",
			Email:      "asavage@yahoo.com",
			VIN:        "1C4RDHDG9DC539254",
			Year:       "2001",
			Make:       "Toyota",
			Model:      "Camry",
			DSFwalkseq: "C312",
			CRRT:       "C001",
			KBB:        "",
		},
		cp.Record{
			Fullname:   "Shepard S. Sam",
			Firstname:  "Shepard",
			MI:         "S.",
			Lastname:   "Sam",
			Address1:   "1 Camino Rd",
			Address2:   "",
			City:       "Anaheim",
			State:      "CA",
			Zip:        "98578",
			Zip4:       "9875",
			HPH:        "(789) 658-1978",
			BPH:        "(684) 578-1234",
			CPH:        "",
			Email:      "ss@gmail.com",
			VIN:        "4A3AK24F67E006257",
			Year:       "2010",
			Make:       "Honda",
			Model:      "Civic",
			DSFwalkseq: "D111",
			CRRT:       "C002",
			KBB:        "",
		},
	}
	for i, r := range parser {
		if r.Fullname != expCase[i].Fullname {
			t.Errorf("Expected %v, got %v", expCase[i].Fullname, r.Fullname)
		}
		if r.Firstname != expCase[i].Firstname {
			t.Errorf("Expected %v, got %v", expCase[i].Firstname, r.Firstname)
		}
		if r.MI != expCase[i].MI {
			t.Errorf("Expected %v, got %v", expCase[i].MI, r.MI)
		}
		if r.Lastname != expCase[i].Lastname {
			t.Errorf("Expected %v, got %v", expCase[i].Lastname, r.Lastname)
		}
		if r.Address1 != expCase[i].Address1 {
			t.Errorf("Expected %v, got %v", expCase[i].Address1, r.Address1)
		}
		if r.Address2 != expCase[i].Address2 {
			t.Errorf("Expected %v, got %v", expCase[i].Address2, r.Address2)
		}
		if r.City != expCase[i].City {
			t.Errorf("Expected %v, got %v", expCase[i].City, r.City)
		}
		if r.State != expCase[i].State {
			t.Errorf("Expected %v, got %v", expCase[i].State, r.State)
		}
		if r.Zip != expCase[i].Zip {
			t.Errorf("Expected %v, got %v", expCase[i].Zip, r.Zip)
		}
		if r.Zip4 != expCase[i].Zip4 {
			t.Errorf("Expected %v, got %v", expCase[i].Zip4, r.Zip4)
		}
		if r.HPH != expCase[i].HPH {
			t.Errorf("Expected %v, got %v", expCase[i].HPH, r.HPH)
		}
		if r.BPH != expCase[i].BPH {
			t.Errorf("Expected %v, got %v", expCase[i].BPH, r.BPH)
		}
		if r.CPH != expCase[i].CPH {
			t.Errorf("Expected %v, got %v", expCase[i].CPH, r.CPH)
		}
		if r.Email != expCase[i].Email {
			t.Errorf("Expected %v, got %v", expCase[i].Email, r.Email)
		}
		if r.VIN != expCase[i].VIN {
			t.Errorf("Expected %v, got %v", expCase[i].VIN, r.VIN)
		}
		if r.Year != expCase[i].Year {
			t.Errorf("Expected %v, got %v", expCase[i].Year, r.Year)
		}
		if r.Make != expCase[i].Make {
			t.Errorf("Expected %v, got %v", expCase[i].Make, r.Make)
		}
		if r.Model != expCase[i].Model {
			t.Errorf("Expected %v, got %v", expCase[i].Model, r.Model)
		}
		if r.DSFwalkseq != expCase[i].DSFwalkseq {
			t.Errorf("Expected %v, got %v", expCase[i].DSFwalkseq, r.DSFwalkseq)
		}
		if r.CRRT != expCase[i].CRRT {
			t.Errorf("Expected %v, got %v", expCase[i].CRRT, r.CRRT)
		}
		if r.KBB != expCase[i].KBB {
			t.Errorf("Expected %v, got %v", expCase[i].KBB, r.KBB)
		}
	}
}
