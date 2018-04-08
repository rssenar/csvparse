csvparse
=====
A simple package for parsing csv files into structs. The API and techniques inspired from https://github.com/gocarina/gocsv but modified to fit my specific use case.

### Installation

```go get -u github.com/rssenar/csvparse```

### Sample:

Give this simple yet common CSV dataset...

```
FIRSTNAME,LASTNAME,ADDRESS_1,CITY,STATE,ZIP,CRRT,DP2,DPC,LOT,LOT_ORD,HPH,CPH,EMAIL,LICENSE,VIN,VYR,VMK,VMD,VML,DIS,ROAMT,DELDATE,IBFLAG,MAIL,TYPE,BPH,CNO,NU,APR,TERM,DATE,SQN,INC
Toby,Avdeef,1000 KELLEY DR,FORT WORTH,TX,76140-3618,C003,0,0,,,,682-227-5578,,,4A3AK24F67E006257,2007,MITSUBISHI,ECLIPSE,,0,,,,,,,,,,,7/10/15,343820,2
Natalie,Jackson,10000 IRON RIDGE DR,FORT WORTH,TX,76140-7527,R040,0,0,,,8175681080,,,,1C3EL46X04N213746,2004,CHRYSLER,SEBRING,,0,,12/31/03,,,,,,,,,10/14/10,343821,2

```
And given this struct...

```
type client struct {
	Fullname  string    `json:"Full_name" csv:"(?i)^fullname$" fmt:"tc"`
	Firstname string    `json:"First_name" csv:"(?i)^first[ _-]?name$" fmt:"tc"`
	MI        string    `json:"Middle_name" csv:"(?i)^mi$" fmt:"uc"`
	Lastname  string    `json:"Last_name" csv:"(?i)^last[ _-]?name$" fmt:"tc"`
	Address1  string    `json:"Address_1" csv:"(?i)^address[ _-]?1?$" fmt:"tc"`
	Address2  string    `json:"Address_2" csv:"(?i)^address[ _-]?2$" fmt:"tc"`
	City      string    `json:"City" csv:"(?i)^city$" fmt:"tc"`
	State     string    `json:"State" csv:"(?i)^state$|^st$" fmt:"uc"`
	Zip       string    `json:"Zip" csv:"(?i)^(zip|postal)[ _]?(code)?$" fmt:"-"`
	Zip4      string    `json:"Zip_4" csv:"(?i)^zip4$|^4zip$" fmt:"-"`
	HPH       string    `json:"Home_phone" csv:"(?i)^hph$|^home[ _]phone$" fmt:"fp"`
	Email     string    `json:"Email" csv:"(?i)^email[ _]?(address)?$" fmt:"lc"`
	Date      time.Time `json:"Last_service_date" csv:"(?i)^date$" fmt:"-"`
}
```

The output will be...

```
[
  {
   "Full_name": "",
   "First_name": "Toby",
   "Middle_name": "",
   "Last_name": "Avdeef",
   "Address_1": "1000 KELLEY DR",
   "Address_2": "",
   "City": "FORT WORTH",
   "State": "TX",
   "Zip": "76140-3618",
   "Zip_4": "",
   "Home_phone": "",
   "Email": "",
   "Last_service_date": "2015-07-10T00:00:00Z"
  },
  {
   "Full_name": "",
   "First_name": "Natalie",
   "Middle_name": "",
   "Last_name": "Jackson",
   "Address_1": "10000 IRON RIDGE DR",
   "Address_2": "",
   "City": "FORT WORTH",
   "State": "TX",
   "Zip": "76140-7527",
   "Zip_4": "",
   "Home_phone": "8175681080",
   "Email": "",
   "Last_service_date": "2010-10-14T00:00:00Z"
  }
 ]
```

### Usage:

Behaviour is dicatated by the struct tags

```

csv:"(?i)^first[ _-]?name$" fmt:"tc"`
csv:"(?i)^last[ _-]?name$" fmt:"tc"`
csv:"(?i)^date$" fmt:"-"`

```
#### csv: Tags:
"csv:" tag maps the struct field name to the csv column header. You can use string matching

```

csv:"Firstname"
```

or preferrably regex pattern mathcing "(?i)^last[ _-]?name$" for added flexibility and utility.

```

csv:"(?i)^first[ _-]?name$"

```
#### fmt: Tags:
"fmt:" tag dictates the formating option for string values.  Options include:

fmt: "tc" - Format string to title case, (eg. john smith -> John Smith)
fmt: "uc" - Format string to upper case, (eg. john smith -> JOHN SMITH)
fmt: "lc" - Format string to lower case, (eg. JOHN SMITH -> john smith)
fmt: "fp" - Format phone number, (eg. 9497858798 -> (949) 785-8798)

By default, leading and trailing spaces are removed
