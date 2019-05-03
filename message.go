// Copyright (c) 2019 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.package encoding

// Package iso8583 implements encoding and decoding of message as defined in ISO 8583.
// For detail about this standard, see https://en.wikipedia.org/wiki/ISO_8583.
package iso8583

import (
	"errors"
	"fmt"
	"math"
	"sort"

	"github.com/rvflash/iso8583/encoding"
	"github.com/rvflash/iso8583/field"
)

// List of known errors.
var (
	// ErrLen is returned if the length not matches with the expected length.
	ErrLen = errors.New("invalid length")
	// ErrMTI is returned if we failed to fields the data.
	ErrMTI = errors.New("invalid message type identifier")
	// ErrOOR is the data exceeds the bounds.
	ErrOOR = errors.New("out of range")
)

// Fields represents all message's fields.
type Fields map[uint8]*field.Data

// Message represents an ISO 8583 message.
type Message struct {
	MTI    *MTI
	Format encoding.Format
	Header bool
	Fields
}

// Marshal returns the ISO 8683 encoding of Message.
func Marshal(_ *Message) ([]byte, error) {
	return nil, errors.New("not implemented")
}

// Unmarshal parses the ISO 8583-encoded data and stores the result in the Message pointed.
func Unmarshal(data []byte, m *Message) error {
	// Parses the Header.
	data, err := m.header(data)
	if err != nil {
		return err
	}
	// Parses the type indicator.
	data, err = m.mti(data)
	if err != nil {
		return err
	}
	// Parses all bitmaps.
	data, err = m.bitmap(data)
	if err != nil {
		return nil
	}
	return m.fields(data, m.elements())
}

func (m *Message) bitmap(src []byte) (dst []byte, err error) {
	var (
		b, s string
		a, z int
	)
	for {
		z = a + m.Format.LenBitmap()
		if len(src) < z {
			return nil, ErrOOR
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
		f1 := field.New(1)
		f1.Size = len(b)
		err = field.Unmarshal(b, f1)
		if err != nil {
			return nil, err
		}
		m.Fields = Fields{1: f1}
		fmt.Printf("bitmap: %s\n", b)

		return src[z:], nil
	}
}

// elements converts the bitmap (binary) to a list of field positions.
func (m *Message) elements() (list []int) {
	f, ok := m.Fields[1]
	if !ok {
		return
	}
	for k, v := range f.Value {
		if v == '1' && k > 0 {
			list = append(list, k+1)
		}
	}
	sort.Ints(list)

	return
}

func (m *Message) fields(data []byte, list []int) error {
	var (
		a, z int
		n    uint8
		err  error
	)
	for _, v := range list {
		if v > math.MaxInt8 {
			return ErrOOR
		}
		n = uint8(v)
		f := field.New(n)
		z += f.Size
		if len(data) < z {
			err = ErrOOR
		} else {
			err = field.Unmarshal(data[a:z], f)
		}
		if err != nil {
			return fmt.Errorf("#%d: %s (%s[%d:%d])", n, err, data, a, z)
		}
		m.Fields[n] = f
		a += f.Size
	}
	return nil
}

func (m *Message) header(src []byte) (dst []byte, err error) {
	if !m.Header {
		return src, nil
	}
	if len(src) < m.Format.LenHeader() {
		return nil, ErrOOR
	}
	n, err := m.Format.EncodeToDecimal(src[:m.Format.LenHeader()])
	if err != nil {
		return nil, err
	}
	dst = src[m.Format.LenHeader():]
	if len(dst) != int(n) {
		return nil, ErrLen
	}
	return
}

func (m *Message) mti(src []byte) (dst []byte, err error) {
	if len(src) < m.Format.LenMTI() {
		return nil, ErrOOR
	}
	m.MTI, err = ParseMTI(string(src[:m.Format.LenMTI()]))
	if err != nil {
		return nil, err
	}
	return src[m.Format.LenMTI():], nil
}
