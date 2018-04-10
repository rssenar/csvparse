csvparse
=====
A simple package for parsing csv files into structs. The API and techniques inspired from  [github.com/gocarina/gocsv](https://github.com/gocarina/gocsv) but modified to fit my specific use cases.

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

Behaviour is dicatated through the struct tags

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

Included is a command line tool if you prefer to use this as a CLI.

from csvparse directory, the run the go install command.

```
>> go install ./cmd/csvparse/
```

you can pass a single file or mutliple files as a argument. Output files will be labeled with the current file name plus _parsed.csv.

```
>> ls
testfile.csv

>> csvparse testfile.csv
2018/04/09 16:27:01 testfile.csv was parser in 801.853858ms

>> ls
testfile.csv		testfile_parsed.csv

```

```
>> ls
testfile1.csv	testfile2.csv	testfile3.csv

>> wc -l $(ls)
   10000 testfile1.csv
   30000 testfile2.csv
   50000 testfile3.csv
   90000 total

>> csvparse $(ls)
2018/04/09 16:29:00 testfile1.csv was parser in 792.097509ms
2018/04/09 16:29:02 testfile2.csv was parser in 2.313945665s
2018/04/09 16:29:06 testfile3.csv was parser in 3.939320765s

>> ls
testfile1.csv		testfile1_parsed.csv	testfile2.csv		testfile2_parsed.csv	testfile3.csv		testfile3_parsed.csv


```

You can also pipe in input from stdin to be printed to stdout or redirected to another file

```
>> cat testfile.csv | csvparse

>> cat testfile.csv | csvparse > out.csv
```

The output file will be the parsed and reformated as a csv representation of the struct provided. Struct field names will be use as the CSV field headers.

```
struct {
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

>> cat test.csv
FIRSTNAME,LASTNAME,ADDRESS_1,CITY,STATE,ZIP,CRRT,DP2,DPC,LOT,LOT_ORD,HPH,CPH,EMAIL,LICENSE,VIN,VYR,VMK,VMD,VML,DIS,ROAMT,DELDATE,IBFLAG,MAIL,TYPE,BPH,CNO,NU,APR,TERM,DATE,SQN,INC
Sherlock,Holmes,1000 BAKER DR,FORT WORTH,TX,76410-3620,C003,0,0,,,,692-250-8078,,,4A3ZM24F67E005557,2007,MITSUBISHI,ECLIPSE,,0,,,,,,,,,,,7/10/15,343820,2
John,Watson,10000 RIDGE DR,FORT WORTH,CA,76240-7530,R040,0,0,,,8205692022,,,,1C3GW46X04N288746,2004,CHRYSLER,SEBRING,,0,,12/31/03,,,,,,,,,10/14/10,343821,2

>> csvparse test.csv
Fullname,Firstname,MI,Lastname,Address1,Address2,City,State,Zip,Zip4,HPH,Email,Date
,Sherlock,,Holmes,1000 Baker Dr,,Fort Worth,TX,76410,3620,,,2015-07-10
,John,,Watson,10000 Ridge Dr,,Fort Worth,CA,76240,7530,(820) 569-2022,,2010-10-14

```

### Specific Use Cases:

For my specific use case, I typically require parsing of Zip codes (e.g 92882-4578) to Zip & ZIp4 components.  By default, if ZIP fields matches the zip regex patern, Zip & ZIp4 components will be parsed to their corresponding fields.

```
  {
   "Zip": "76410-3620",
   "Zip_4": "",
  }

  will be parsed to...

  {
   "Zip": "76410",
   "Zip_4": "3620",
  }

```

Another specialized us case is Fullname parsing.  If a fullname field is provided amd there is a missing First or Last Name field, Fullname will be parsed to First. Middle and Last Name fields respectively.

Fullname parsing is handled by [github.com/blendlabs/go-name-parser](https://github.com/blendlabs/go-name-parser) package.

```
  {
   "Full_name": "Sherlock H. Holmes",
   "First_name": "Ed",
   "Middle_name": "",
   "Last_name": "",
  }

  will be parsed to...

  {
   "Full_name": "Sherlock H. Holmes",
   "First_name": "Sherlock",
   "Middle_name": "H.",
   "Last_name": "Holmes",
  }

```

### Caveats & Limitations:

in its current version, the package ONLY supports parsing to []structs.  Passing anything other than []structs results in a panic.

### License:

The MIT License (MIT)