// Copyright (c) 2019 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.package encoding

package iso8583

import (
	"fmt"
	"math"
	"sort"

	"github.com/rvflash/iso8583/encoding"
	"github.com/rvflash/iso8583/errors"
	"github.com/rvflash/iso8583/field"
)

// Field represents all message's fields.
type Fields map[field.ID]field.Field

// Message represents an iso 8583 message.
type Message struct {
	MTI    *MTI
	Format encoding.Format
	Header bool
	Data   Fields
}

// Fields returns the list of Field Elements.
func (m *Message) Fields() Fields {
	return m.Data
}

// Type returns the Message Type Indicator.
func (m *Message) Type() string {
	if m.MTI == nil || !m.MTI.Valid() {
		return ""
	}
	return m.MTI.String()
}

// todo
// String returns the message as a string.
// Dump string

// bitmap extracts this data and returns the rest of the message.
func (m *Message) bitmap(src []byte) (dst []byte, err error) {
	var (
		b, s string
		a, z int
	)
	for {
		z = a + m.Format.LenBitmap()
		if len(src) < z {
			return nil, errors.OutOfRange
		}
		s, err = m.Format.EncodeToBinary(src[a:z])
		if err != nil {
			return nil, err
		}
		b += s

		// The first byte indicates the presence of an other bitmap.
		if s[0] == '1' {
			a += m.Format.LenBitmap()
			continue
		}
		// Prepares the fields list
		err = m.make([]byte(b))
		if err != nil {
			return nil, err
		}
		return src[z:], nil
	}
}

// elements converts the binary bitmap to a list of field positions.
func (m *Message) elements() (list []int) {
	f1, ok := m.Data[1]
	if !ok {
		return
	}
	for k, v := range f1.String() {
		if v == '1' && k > 0 {
			list = append(list, k+1)
		}
	}
	sort.Ints(list)

	return
}

// fields sets the data elements based on the message and the known fields in the bitmap.
func (m *Message) fields(data []byte, list []int) error {
	var (
		a, s int
		err  error
	)
	for _, v := range list {
		if v > math.MaxInt8 {
			return errors.New(errors.Data, v)
		}
		f := field.New(field.ID(v))
		s, err = f.FixedSize(data[a:])
		if err != nil {
			return errors.New(err, v)
		}
		err = field.Unmarshal(data[a:], f)
		if err != nil {
			fmt.Println(m.Data)
			fmt.Println(string(data))
			fmt.Println(a, s)
			return errors.New(err, v)
		}
		m.Data[f.ID()] = f
		a += s
	}
	return nil
}

// header extracts the header length is needed and returns the rest of the message.
func (m *Message) header(src []byte) (dst []byte, err error) {
	if !m.Header {
		return src, nil
	}
	if len(src) < m.Format.LenHeader() {
		return nil, errors.OutOfRange
	}
	n, err := m.Format.EncodeToDecimal(src[:m.Format.LenHeader()])
	if err != nil {
		return nil, err
	}
	dst = src[m.Format.LenHeader():]
	if len(dst) != int(n) {
		return nil, errors.Length
	}
	return
}

// make sets the bitmap as the first data elements.
func (m *Message) make(bitmap []byte) error {
	f1 := field.New(1)
	f1.Size = len(bitmap)
	err := field.Unmarshal(bitmap, f1)
	if err != nil {
		return err
	}
	m.Data = Fields{1: f1}

	return nil
}

// mti sets the message type indicator and returns the rest of the message.
func (m *Message) mti(src []byte) (dst []byte, err error) {
	if len(src) < m.Format.LenMTI() {
		return nil, errors.OutOfRange
	}
	m.MTI, err = ParseMTI(string(src[:m.Format.LenMTI()]))
	if err != nil {
		return nil, err
	}
	return src[m.Format.LenMTI():], nil
}
