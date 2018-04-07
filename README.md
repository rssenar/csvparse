csvparse
=====
A simple package for parsing csv files into structs. The API and techniques inspired from https://github.com/gocarina/gocsv but modified to fit my specific use case.

### Installation

```go get -u github.com/rssenar/csvparse```

## Usage:

Give this simple yet common CSV dataset...

```
Company,First Name,Last Name,City,State,Postal Code,Email Address,Insert Date,2 sel
,Michael,Chartan,Great Neck,NY  ,11024,MICHAELCHARTAN@GMAIL.COM,9/16/17 9:26,A
,Ira,kolin,MELVILLE,NY  ,11747,Daddy.kolin@gmail.com,9/16/17 10:36,B

```
Given this struct...

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

The desired output will be...

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