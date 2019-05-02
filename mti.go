// Package iso8583 implements encoding and decoding of message as defined in ISO 8583.
// For detail about this standard, see https://en.wikipedia.org/wiki/ISO_8583.
package iso8583

import (
	"fmt"
)

// NewMTI returns a new Message Type Identifier.
// Its variadic behavior allows to give only the first digit with a default behavior for the next.
func NewMTI(digit ...uint8) *MTI {
	var v, c, f, o uint8
	switch len(digit) {
	case 4:
		o = digit[3]
		fallthrough
	case 3:
		f = digit[2]
		fallthrough
	case 2:
		c = digit[1]
		fallthrough
	case 1:
		v = digit[0]
	}
	return &MTI{
		Version:  Version(v),
		Class:    Class(c),
		Function: Function(f),
		Origin:   Origin(o),
	}
}

// ParseMTI parses a string to extract the Message Type Identifier cached inside.
// If we fail, an error is returned instead.
func ParseMTI(s string) (*MTI, error) {
	switch len(s) {
	case 3:
		// deals with missing version (aka GPS style)
		return parse("0" + s)
	case 4:
		return parse(s)
	default:
		return nil, ErrMTI
	}
}

// MTI aka. Message Type Identifier
type MTI struct {
	Version  Version
	Class    Class
	Function Function
	Origin   Origin
}

// Valid returns in success if all parts of the message type identifier are valid.
func (m *MTI) Valid() bool {
	if !m.Version.Valid() {
		return false
	}
	if !m.Class.Valid() {
		return false
	}
	if !m.Function.Valid() {
		return false
	}
	return m.Origin.Valid()
}

// String implements the fmt.Stringer interface.
func (m *MTI) String() string {
	return fmt.Sprintf("%d%d%d%d", m.Version, m.Class, m.Function, m.Origin)
}

func parse(s string) (*MTI, error) {
	var (
		d uint8
		m = new(MTI)
	)
	for i, r := range s {
		if r < '0' || r > '9' {
			return nil, ErrMTI
		}
		d = uint8(r - '0')
		switch i {
		case 0:
			m.Version = Version(d)
		case 1:
			m.Class = Class(d)
		case 2:
			m.Function = Function(d)
		case 3:
			m.Origin = Origin(d)
		}
	}
	if !m.Valid() {
		return nil, ErrMTI
	}
	return m, nil
}