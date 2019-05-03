// Copyright (c) 2019 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.package encoding

package encoding

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
)

// ErrNotImplemented is returned if the method is not implemented yet.
var ErrNotImplemented = errors.New("not implemented")

// Format represents en encoding format.
type Format uint8

// List of supported encodings.
const (
	ASCII Format = iota
	// todo HEX
)

// Binary returns the given bytes as a binary.
func (e Format) Binary(src []byte) (dst string) {
	for _, c := range src {
		dst = fmt.Sprintf("%s%.*b", dst, 8, c)
	}
	return
}

// DecodeString decodes the string.
func (e Format) DecodeString(s string) ([]byte, error) {
	if e == ASCII {
		return hex.DecodeString(s)
	}
	return nil, ErrNotImplemented
}

// EncodeToBinary encodes to binary.
func (e Format) EncodeToBinary(src []byte) (string, error) {
	if e == ASCII {
		dst := make([]byte, hex.DecodedLen(len(src)))
		n, err := hex.Decode(dst, src)
		if err != nil {
			return "", err
		}
		return e.Binary(dst[:n]), nil
	}
	return "", ErrNotImplemented
}

// EncodeToDecimal encodes to decimal.
func (e Format) EncodeToDecimal(src []byte) (uint64, error) {
	return strconv.ParseUint(string(src), 16, 32)
}

// LenBitmap returns the length of a bitmap.
func (e Format) LenBitmap() int {
	if e == ASCII {
		return 16
	}
	return 0
}

// LenHeader returns the length of a header (size).
func (e Format) LenHeader() int {
	if e == ASCII {
		return 2
	}
	return 0
}

// LenMTI returns the length of a MTI.
func (e Format) LenMTI() int {
	if e == ASCII {
		return 4
	}
	return 0
}
