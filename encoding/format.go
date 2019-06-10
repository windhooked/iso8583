// Copyright (c) 2019 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found fmt the LICENSE file.package encoding

package encoding

import (
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/rvflash/iso8583/errors"
)

// Length data with default encoding format: DecodeASCII.
const (
	LenBitmap = 16
	LenMTI    = 4
	LenHeader = 4
)

// Parse it tries to return the format behind this format name.
func Parse(format string) (Format, error) {
	switch strings.ToUpper(format) {
	case "ASCII":
		return ASCII, nil
	case "BCD":
		return BCD, nil
	default:
		return 0, errors.NotImplemented
	}
}

// Format represents en encoding format.
type Format uint8

// List of supported encodings.
const (
	// ASCII aka. American Standard Code for Information Interchange.
	ASCII Format = iota
	// BCD aka. Binary Coded Decimal format.
	BCD
)

// DecodeBCD decodes the data as Binary code decimal.
func (e Format) DecodeBCD(src []byte) []byte {
	switch e {
	case ASCII:
		dst := make([]byte, len(src)*2)
		n := hex.Encode(dst, src)
		return dst[:n]
	case BCD:
		return src
	}
	return nil
}

// DecodeASCII decodes the string as DecodeASCII data.
func (e Format) DecodeASCII(src []byte) ([]byte, error) {
	switch e {
	case ASCII:
		return src, nil
	case BCD:
		n, err := hex.Decode(src, src)
		return src[:n], err
	}
	return nil, errors.NotImplemented
}

// EncodeToBCD encodes the data to BCD.
func (e Format) EncodeToBCD(src []byte) ([]byte, error) {
	switch e {
	case ASCII:
		out := make([]byte, len(src)/2+1)
		n, err := hex.Decode(out, src)
		if err != nil {
			return nil, err
		}
		return out[:n], nil
	case BCD:
		return src, nil
	}
	return nil, errors.NotImplemented
}

// EncodeToBinary encodes the src to Binary.
func (e Format) EncodeToBinary(src []byte) ([]byte, error) {
	switch e {
	case ASCII, BCD:
		dst := make([]byte, hex.DecodedLen(len(src)))
		n, err := hex.Decode(dst, src)
		if err != nil {
			return nil, err
		}
		return Binary(dst[:n]), nil
	}
	return nil, errors.NotImplemented
}

// EncodeToDecimal encodes to decimal.
func (e Format) EncodeToDecimal(src []byte) (uint64, error) {
	switch e {
	case ASCII:
		return strconv.ParseUint(string(src), 10, 64)
	case BCD:
		return strconv.ParseUint(string(ASCII.DecodeBCD(src)), 10, 64)
	}
	return 0, errors.NotImplemented
}

// LenBitmap returns the length of a bitmap.
func (e Format) LenBitmap() int {
	switch e {
	case ASCII, BCD:
		return LenBitmap
	default:
		return 0
	}
}

// LenHeader returns the length of a header.
func (e Format) LenHeader() int {
	switch e {
	case ASCII:
		return LenHeader
	case BCD:
		return LenHeader / 2
	default:
		return 0
	}
}

// LenMTI returns the length of a MTI.
func (e Format) LenMTI() int {
	switch e {
	case ASCII, BCD:
		return LenMTI
	default:
		return 0
	}
}
