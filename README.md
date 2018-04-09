csvparse
=====
A simple package for parsing csv files into structs. The API and techniques inspired from https://github.com/gocarina/gocsv but modified to fit my specific use cases.

### Installation

```go get -u github.com/rssenar/csvparse```

### Sample:

Given this CSV dataset...

```
FIRSTNAME,LASTNAME,ADDRESS_1,CITY,STATE,ZIP,CRRT,DP2,DPC,LOT,LOT_ORD,HPH,CPH,EMAIL,LICENSE,VIN,VYR,VMK,VMD,VML,DIS,ROAMT,DELDATE,IBFLAG,MAIL,TYPE,BPH,CNO,NU,APR,TERM,DATE,SQN,INC
Sherlock,Holmes,1000 BAKER DR,FORT WORTH,TX,76410-3620,C003,0,0,,,,692-250-8078,,,4A3ZM24F67E005557,2007,MITSUBISHI,ECLIPSE,,0,,,,,,,,,,,7/10/15,343820,2
John,Watson,10000 RIDGE DR,FORT WORTH,CA,76240-7530,R040,0,0,,,8205692022,,,,1C3GW46X04N288746,2004,CHRYSLER,SEBRING,,0,,12/31/03,,,,,,,,,10/14/10,343821,2

```
And given this struct...

```
type D struct {
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

data := []*D{}

err = cp.NewDecoder(os.Stdin).DecodeCSV(&data)
if err != nil {
  log.Fatalln(err)
}

```

The output will be...

```
[
  {
   "Full_name": "",
   "First_name": "Sherlock",
   "Middle_name": "",
   "Last_name": "Holmes",
   "Address_1": "1000 Baker Dr",
   "Address_2": "",
   "City": "Fort Worth",
   "State": "TX",
   "Zip": "76410",
   "Zip_4": "3620",
   "Home_phone": "",
   "Email": "",
   "Last_service_date": "2015-07-10T00:00:00Z"
  },
  {
   "Full_name": "",
   "First_name": "John",
   "Middle_name": "",
   "Last_name": "Watson",
   "Address_1": "10000 Ridge Dr",
   "Address_2": "",
   "City": "Fort Worth",
   "State": "CA",
   "Zip": "76240",
   "Zip_4": "7530",
   "Home_phone": "(820) 569-2022",
   "Email": "",
   "Last_service_date": "2010-10-14T00:00:00Z"
  }
 ]
```

### Usage:

Behaviour is dicatated by the struct tags

```
type struct {
	Firstname string    `json:"First_name" csv:"(?i)^first[ _-]?name$" fmt:"tc"`
	Lastname  string    `json:"Last_name" csv:"(?i)^last[ _-]?name$" fmt:"tc"`
}

```
### csv: Tags:

"csv:" tag maps the struct field name to the csv column header. You can use string matching

```

csv:"Firstname"
```

or preferrably regex pattern mathcing "(?i)^last[ _-]?name$" for added flexibility and utility.

```

csv:"(?i)^first[ _-]?name$"

```
if csv struct tag is not provided or tag is "-" (eg. csv:"-"). struct field name will be used for matching.

### fmt: Tags:

"fmt:" tag dictates the formating option for string values.  Options include:

```
fmt: "tc" - Format string to title case, (eg. john smith -> John Smith)
fmt: "uc" - Format string to upper case, (eg. john smith -> JOHN SMITH)
fmt: "lc" - Format string to lower case, (eg. JOHN SMITH -> john smith)
fmt: "fp" - Format phone number, (eg. 9497858798 -> (949) 785-8798)
fmt: "-" - No formating specified
```

By default, leading and trailing spaces are removed. if fmt struct tag is not provided or tag is "-" (eg. fmt:"-"). Original string value will be returned.

Time Values (time.Time):
For struct field that are type time.Time will be formated according to RFC3339 time format (i.e. 2010-10-14T00:00:00Z). Empty time values will return will return Go's zero date value (0001-01-01 00:00:00 +0000 UTC).

### Command Line Tool:

I also proviede a handy tool if you prefer to use this as a CLI.

from csvparse directory, the run the go install command.

```
>> go install ./cmd/csvparse/
```

you can pass the file as a argument (only supports sigle files, will panic if given mutiple files as arguments)

```
>> csvparse testfile.csv
```

or you can pipe input as a data stream from stdin

```
>> cat testfile.csv | csvparse
```

The output will be the parsed and reformated csv representation of the struct field provided. struct field names will be use as the CSV field headers.

```
>> cat test.csv
FIRSTNAME,LASTNAME,ADDRESS_1,CITY,STATE,ZIP,CRRT,DP2,DPC,LOT,LOT_ORD,HPH,CPH,EMAIL,LICENSE,VIN,VYR,VMK,VMD,VML,DIS,ROAMT,DELDATE,IBFLAG,MAIL,TYPE,BPH,CNO,NU,APR,TERM,DATE,SQN,INC
Sherlock,Holmes,1000 BAKER DR,FORT WORTH,TX,76410-3620,C003,0,0,,,,692-250-8078,,,4A3ZM24F67E005557,2007,MITSUBISHI,ECLIPSE,,0,,,,,,,,,,,7/10/15,343820,2
John,Watson,10000 RIDGE DR,FORT WORTH,CA,76240-7530,R040,0,0,,,8205692022,,,,1C3GW46X04N288746,2004,CHRYSLER,SEBRING,,0,,12/31/03,,,,,,,,,10/14/10,343821,2

>> csvparse test.csv
Fullname,Firstname,MI,Lastname,Address1,Address2,City,State,Zip,Zip4,HPH,Email,Date
,Sherlock,,Holmes,1000 Baker Dr,,Fort Worth,TX,76410,3620,,,2015-07-10
,John,,Watson,10000 Ridge Dr,,Fort Worth,CA,76240,7530,(820) 569-2022,,2010-10-14

```

### Performance:

performace is good but looking to improve performance with future concurrency optimizations.

```
>> wc -l $(ls)
   10000 a.csv
   30000 b.csv
   50000 c.csv
  100000 d.csv
  190000 total

>> csvparse a.csv > /dev/null
2018/04/07 20:24:50 CSVParser took 361.195842ms
>> csvparse b.csv > /dev/null
2018/04/07 20:24:55 CSVParser took 866.488151ms
>> csvparse c.csv > /dev/null
2018/04/07 20:25:00 CSVParser took 1.414062947s
>> csvparse d.csv > /dev/null
2018/04/07 20:25:07 CSVParser took 2.810645883s
```

### Caveats & Limitations:

in its current version, the package ONLY supports parsing to []structs.  Passing anything other than []structs results in a panic.

### License:

The MIT License (MIT)