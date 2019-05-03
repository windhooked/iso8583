// Copyright (c) 2019 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.package encoding

package field

import (
	"fmt"
	"strconv"
	"time"
	"unicode"
)

// Marshal returns the encoded value of v.
func Marshal(v *Data) ([]byte, error) {
	if !v.Valid() {
		return nil, ErrData
	}
	switch v.Type {
	case LVar, LLVar, LLLVar:
		switch {
		case v.Format == Numeric:
			return []byte(fmt.Sprintf("%0*d", v.Size, v.Value)), nil
		case v.Format&Amount != 0:
			return []byte(fmt.Sprintf("%s%0*d", v.Value[1:], v.Size-1, v.Value[1:])), nil
		default:
			return []byte(fmt.Sprintf("%*v", v.Size, v.Value)), nil
		}
	}
	return v.Value, nil
}

// Unmarshal parses the gives data and stores the result into the Data pointed.
func Unmarshal(data interface{}, v *Data) error {
	switch d := data.(type) {
	case []byte:
		v.Value = d
	case string:
		v.Value = []byte(d)
	case int:
		v.Value = []byte(strconv.Itoa(d))
	case int64:
		v.Value = []byte(strconv.FormatInt(d, 10))
	}
	if !v.Valid() {
		return ErrData
	}
	return nil
}

// New returns a new instance of Data.
func New(num uint8) *Data {
	return &Data{
		Format: n[num].Format,
		Size:   n[num].Size,
		Type:   n[num].Type,
	}
}

// Data represents a data element.
type Data struct {
	Format Format
	Type   Type
	Size   int
	Value  []byte
}

// Int64 tries to returns the data value as an int64.
func (v *Data) Int64() (int64, error) {
	if v.Value == nil || v.Size == 0 {
		return 0, ErrData
	}
	switch {
	case v.Format&Numeric != 0:
		return strconv.ParseInt(v.String(), 16, 64)
	case v.Format&Amount != 0:
		return strconv.ParseInt(v.String()[1:], 16, 64)
	}
	return 0, ErrData
}

// String implements the fmt.Stringer interface.
func (v *Data) String() string {
	if v.Value == nil || v.Size == 0 {
		return ""
	}
	return string(v.Value)
}

// List of supported date formats
const (
	timeFmt      = "150405"
	monthDayFmt  = "0102"
	yearFmt      = "06"
	yearMonthFmt = yearFmt + "01"
	dateFmt      = yearFmt + monthDayFmt
)

// Time tries to returns the data value as a time.Time.
func (v *Data) Time() (time.Time, error) {
	if v.Value == nil || v.Size == 0 {
		return time.Time{}, ErrData
	}
	switch {
	case v.Format&MonthDay != 0:
		var (
			layout = yearFmt + monthDayFmt
			value  = time.Now().Format(yearFmt) + v.String()
		)
		if v.Format&Time != 0 {
			return time.Parse(layout+timeFmt, value)
		}
		return time.Parse(layout, value)
	case v.Format&Date != 0:
		return time.Parse(dateFmt, v.String())
	case v.Format&Time != 0:
		return time.Parse(timeFmt, v.String())
	case v.Format&YearMonth != 0:
		return time.Parse(yearMonthFmt, v.String())
	default:
		return time.Time{}, ErrData
	}
}

// Valid returns in success if the value of the data is valid.
func (v *Data) Valid() bool {
	if v.Value == nil || v.Size == 0 {
		return false
	}
	switch v.Format {
	case Alpha:
		return are(string(v.Value), unicode.IsLetter, unicode.IsSpace)
	case Alpha | Numeric:
		return are(string(v.Value), unicode.IsLetter, unicode.IsSpace, unicode.IsNumber)
	case Alpha | Numeric | Special, Track:
		return are(string(v.Value), unicode.IsLetter, unicode.IsSpace, unicode.IsNumber, unicode.IsSymbol)
	case Numeric | Amount:
		if len(v.Value) < 2 {
			return false
		}
		return are(string(v.Value[:1]), isAmount) && are(string(v.Value[1:]), unicode.IsNumber)
	case Binary:
		return are(string(v.Value), isBinary)
	default:
		return are(string(v.Value), unicode.IsNumber)
	}
}

func are(s string, fn ...func(r rune) bool) bool {
	var is = func(r rune) bool {
		for _, f := range fn {
			if f(r) {
				return true
			}
		}
		return false
	}
	for _, r := range s {
		if !is(r) {
			return false
		}
	}
	return true
}

// List of types of amount.
const (
	credit = 'C'
	debit  = 'D'
)

func isAmount(r rune) bool {
	return r == credit || r == debit
}

func isBinary(r rune) bool {
	return r == '0' || r == '1'
}
