// Copyright (c) 2019 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.package encoding

package field

import (
	"fmt"
	"strconv"
	"time"
	"unicode"

	"github.com/rvflash/iso8583/errors"
)

// Marshal returns the encoded value of v.
// todo
func Marshal(v *Data) ([]byte, error) {
	if !v.Valid() {
		return nil, errors.Data
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

// Unmarshal parses the gives data and stores the result into the Field pointed.
func Unmarshal(data []byte, d *Data) (err error) {
	d.Size, err = d.FixedSize(data)
	if err != nil {
		return err
	}
	if len(data) < d.Size {
		return errors.OutOfRange
	}
	d.Value = data[d.prefixSize():d.Size]
	d.Size -= d.prefixSize()

	if !d.Valid() {
		return errors.Data
	}
	return nil
}

// New returns a new instance of Field.
func New(num ID) *Data {
	return &Data{
		Element: n[num],
		Pos:     num,
	}
}

// Field represents a data element.
type Data struct {
	Element
	Pos   ID
	Value []byte
}

// FixedSize implements the Field interface.
func (d *Data) FixedSize(raw []byte) (int, error) {
	prefix := d.prefixSize()
	switch {
	case prefix == 0:
		return d.Size, nil
	case len(raw) < prefix:
		return 0, errors.OutOfRange
	default:
		// Extracts the prefix header showing the length of the field.
		i, err := strconv.Atoi(string(raw[:prefix]))
		if err != nil {
			return 0, err
		}
		// Adds to it the prefix length.
		i = i + prefix
		if i == 0 {
			return 0, errors.Length
		}
		return i, nil
	}
}

func (d *Data) prefixSize() int {
	var prefix int
	switch d.Type {
	case LLLVar:
		prefix++
		fallthrough
	case LLVar:
		prefix++
		fallthrough
	case LVar:
		prefix++
	}
	return prefix
}

// ID implements the Field interface.
func (d *Data) ID() ID {
	return d.Pos
}

// Int64 implements the Field interface.
func (d *Data) Int64() (int64, error) {
	if d.Value == nil || d.Size == 0 {
		return 0, errors.Data
	}
	s := d.String()
	switch {
	case d.Format&Amount != 0:
		if len(s) < 1 {
			return 0, errors.Data
		}
		s = s[1:]
		fallthrough
	case d.Format&Numeric != 0:
		return strconv.ParseInt(s, 10, 64)
	default:
		return 0, errors.Data
	}
}

// String implements the Field interface.
func (d *Data) String() string {
	if d.Value == nil || d.Size == 0 {
		return ""
	}
	return string(d.Value)
}

// List of supported date formats
const (
	timeFmt      = "150405"
	monthDayFmt  = "0102"
	yearFmt      = "06"
	yearMonthFmt = yearFmt + "01"
	dateFmt      = yearFmt + monthDayFmt
)

// Time implements the Field interface.
func (d *Data) Time() (time.Time, error) {
	if d.Value == nil || d.Size == 0 {
		return time.Time{}, errors.Data
	}
	switch {
	case d.Format&MonthDay != 0:
		var (
			layout = yearFmt + monthDayFmt
			value  = time.Now().UTC().Format(yearFmt) + d.String()
		)
		if d.Format&Time != 0 {
			return time.Parse(layout+timeFmt, value)
		}
		return time.Parse(layout, value)
	case d.Format&Date != 0:
		return time.Parse(dateFmt, d.String())
	case d.Format&Time != 0:
		return time.Parse(timeFmt, d.String())
	case d.Format&YearMonth != 0:
		return time.Parse(yearMonthFmt, d.String())
	default:
		return time.Time{}, errors.Data
	}
}

// Valid implements the Field interface.
func (d *Data) Valid() bool {
	if d.Value == nil || d.Size == 0 {
		return false
	}
	switch d.Format {
	case Alpha:
		return are(string(d.Value), unicode.IsLetter, unicode.IsSpace)
	case Alpha | Numeric:
		return are(string(d.Value), unicode.IsLetter, unicode.IsSpace, unicode.IsNumber)
	case Alpha | Numeric | Special, Track:
		return are(string(d.Value), unicode.IsLetter, unicode.IsSpace, unicode.IsNumber, unicode.IsSymbol)
	case Numeric | Amount:
		if len(d.Value) < 2 {
			return false
		}
		return are(string(d.Value[:1]), isAmount) && are(string(d.Value[1:]), unicode.IsNumber)
	case Binary:
		return are(string(d.Value), isBinary)
	default:
		return are(string(d.Value), unicode.IsNumber)
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
