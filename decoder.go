package csvparse

import (
	"errors"
	"io"
	"reflect"
	"regexp"
	"time"
)

// CSVDecoder holds the header field map and io.reader interface
type CSVDecoder struct {
	header map[string]int
	file   io.Reader
}

// NewDecoder allocates a new instance of CSVDecoder
func NewDecoder(input io.Reader) *CSVDecoder {
	return &CSVDecoder{
		header: map[string]int{},
		file:   input,
	}
}

// DecodeCSV unmarshalls CSV file to a specified struct type
func (d *CSVDecoder) DecodeCSV(v interface{}) error {
	// Optional timer function for determining function duration
	defer timeTrack(time.Now(), "Unmarshalls CSV file")

	// getCSVRows grabs [][]strings from spcified input
	csvRows, err := GetCSVRows(d.file)
	if err != nil {
		return err
	}

	// returns an error if an empty file was provided
	if len(csvRows) == 0 {
		return errors.New("empty csv file")
	}
	// grab Header row
	headerRow := csvRows[0]
	// grab get remaining rows as body
	body := csvRows[1:]

	// check header row for duplicate fields
	// if duplicate fields found, return error
	if err := CheckForDoubleHeaderNames(headerRow); err != nil {
		return err
	}

	// check interface type (v)
	// if type is not a slice, return error
	slice := CheckInterfaceValue(v)
	if slice.Kind() != reflect.Slice {
		return errors.New("Only slice types are permited")
	}

	// check inner slice type
	// if type is not a struct, return error
	innerType := GetInnerSliceType(v)
	if innerType.Kind() != reflect.Struct {
		return errors.New("Only slice of stucts type permited")
	}

	// allocate new reflect stuct instance from inner type
	innerValueHdr := reflect.New(innerType)

	// range over header row
	for i, csvColumnHdr := range headerRow {
		// range over struct fields
		for j := 0; j < innerValueHdr.Elem().NumField(); j++ {
			var regex string
			// grab regex string from struct tag
			// if not struct tags provided, struct name will be used
			if rgx, ok := reflect.Indirect(innerValueHdr).Type().Field(j).Tag.Lookup("csv"); ok {
				switch rgx {
				case "-":
					regex = reflect.Indirect(innerValueHdr).Type().Field(j).Name
				default:
					regex = rgx
				}
			} else {
				continue
			}
			if regexp.MustCompile(regex).MatchString(csvColumnHdr) {
				d.header[reflect.Indirect(innerValueHdr).Type().Field(j).Name] = i
			}
		}
	}

	// regex = reflect.Indirect(innerValueHdr).Type().Field(j).Name

	for _, csvRow := range body {
		innerValueRow := reflect.New(innerType)
		for j := 0; j < innerValueRow.Elem().NumField(); j++ {
			sFName := reflect.Indirect(innerValueRow).Type().Field(j).Name
			switch innerValueRow.Elem().Type().Field(j).Type {
			case reflect.TypeOf(time.Now()):
				if _, ok := d.header[sFName]; ok {
					val := ParseDate(csvRow[d.header[sFName]])
					innerValueRow.Elem().FieldByName(sFName).Set(reflect.ValueOf(val))
				}
			default:
				if _, ok := d.header[sFName]; ok {
					val := csvRow[d.header[sFName]]
					innerValueRow.Elem().FieldByName(sFName).Set(reflect.ValueOf(val))
				}
			}
		}

		slice.Set(reflect.Append(slice, innerValueRow))
	}
	return nil
}

// func (d *CSVDecoder) DecodeCSV(v interface{}) error {
// 	defer timeTrack(time.Now(), "DecodeCSV Func")

// 	csvRows, err := getCSVRows(d.file)
// 	if err != nil {
// 		return err
// 	}
// 	if len(csvRows) == 0 {
// 		return errors.New("empty csv file")
// 	}
// 	headerRow := csvRows[0]
// 	body := csvRows[1:]

// 	if err := checkForDoubleHeaderNames(headerRow); err != nil {
// 		return err
// 	}

// 	slice := checkifSlice(v)
// 	if slice.Kind() != reflect.Slice {
// 		return errors.New("Only slice type permited")
// 	}
// 	innerType := getInnerSliceType(v)
// 	if innerType.Kind() != reflect.Struct {
// 		return errors.New("Only slice slice of stucts permited")
// 	}

// 	innerValueHdr := reflect.New(innerType)
// 	innerValueHdrLen := innerValueHdr.Elem().NumField()

// 	for i, csvColumnHdr := range headerRow {
// 		for j := 0; j < innerValueHdrLen; j++ {
// 			regex := reflect.Indirect(innerValueHdr).Type().Field(j).Tag.Get("csv")
// 			if regexp.MustCompile(regex).MatchString(csvColumnHdr) {
// 				d.header[reflect.Indirect(innerValueHdr).Type().Field(j).Name] = i
// 			}
// 		}
// 	}

// 	for _, csvRow := range body {
// 		innerValueRow := reflect.New(innerType)
// 		innerValueRowLen := innerValueRow.Elem().NumField()

// 		for j := 0; j < innerValueRowLen; j++ {
// 			sFName := reflect.Indirect(innerValueRow).Type().Field(j).Name
// 			switch innerValueRow.Elem().Type().Field(j).Type {
// 			case reflect.TypeOf(""):
// 				if _, ok := d.header[sFName]; ok {
// 					if format, ok := reflect.Indirect(innerValueRow).Type().Field(j).Tag.Lookup("fmt"); ok {
// 						if format != "-" {
// 							val, err := FormatStringVals(LCase(format), csvRow[d.header[sFName]])
// 							if err != nil {
// 								return err
// 							}
// 							innerValueRow.Elem().FieldByName(sFName).Set(reflect.ValueOf(val))
// 						} else {
// 							val := csvRow[d.header[sFName]]
// 							innerValueRow.Elem().FieldByName(sFName).Set(reflect.ValueOf(val))
// 						}
// 					} else {
// 						val := csvRow[d.header[sFName]]
// 						innerValueRow.Elem().FieldByName(sFName).Set(reflect.ValueOf(val))
// 					}
// 				}
// 			case reflect.TypeOf(time.Now()):
// 				if _, ok := d.header[sFName]; ok {
// 					val := ParseDate(csvRow[d.header[sFName]])
// 					innerValueRow.Elem().FieldByName(sFName).Set(reflect.ValueOf(val))
// 				}
// 			default:
// 				if _, ok := d.header[sFName]; ok {
// 					val := csvRow[d.header[sFName]]
// 					innerValueRow.Elem().FieldByName(sFName).Set(reflect.ValueOf(val))
// 				}
// 			}
// 			if j == innerValueRowLen-1 {
// 				fullname := innerValueRow.Elem().FieldByName("Fullname")
// 				fname := innerValueRow.Elem().FieldByName("Firstname")
// 				lname := innerValueRow.Elem().FieldByName("Lastname")
// 				if fullname.IsValid() && fname.IsValid() && lname.IsValid() {
// 					if fullname.String() != "" && (fname.String() == "" || lname.String() == "") {
// 						name := names.Parse(fullname.String())
// 						innerValueRow.Elem().FieldByName("Firstname").Set(reflect.ValueOf(name.FirstName))
// 						innerValueRow.Elem().FieldByName("MI").Set(reflect.ValueOf(name.MiddleName))
// 						innerValueRow.Elem().FieldByName("Lastname").Set(reflect.ValueOf(name.LastName))
// 					}
// 				}
// 				zip := innerValueRow.Elem().FieldByName("Zip")
// 				if zip.IsValid() {
// 					zip, zip4 := ParseZip(zip.String())
// 					innerValueRow.Elem().FieldByName("Zip").Set(reflect.ValueOf(zip))
// 					if zip4 != "" {
// 						innerValueRow.Elem().FieldByName("Zip4").Set(reflect.ValueOf(zip4))
// 					}
// 				}
// 			}
// 		}
// 		slice.Set(reflect.Append(slice, innerValueRow))
// 	}
// 	return nil
// }
