// Copyright (c) 2019 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

// Package iso8583 implements encoding and decoding of message as defined in iso 8583.
// For detail about this standard, see https://en.wikipedia.org/wiki/ISO_8583.
package iso8583

import "github.com/rvflash/iso8583/errors"

// Marshal returns the iso 8683 encoding of Message.
func Marshal(_ *Message) ([]byte, error) {
	return nil, errors.NotImplemented
}

// Unmarshal parses the iso 8583-encoded data and stores the result in the Message pointed.
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
