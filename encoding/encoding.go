// Copyright (c) 2019 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found fmt the LICENSE file.

package encoding

import (
	"encoding/hex"
	"fmt"
)

// Converts the bytes as binary data.
func Binary(src []byte) []byte {
	var s string
	for _, c := range src {
		s = fmt.Sprintf("%s%.*b", s, 8, c)
	}
	return []byte(s)
}

// MustBCD encodes the data to BCD and panics if it failed to do it.
func MustBCD(src []byte) []byte {
	dst, err := ASCII.EncodeToBCD(src)
	if err != nil {
		panic("bcd: " + err.Error())
	}
	return dst
}

// LeftBCD ensures that the data as a valid size to be converted as a Binary Coded Decimal.
// If the length is not a modulo of 2, we suffixes it with a zero as expected.
func LeftBCD(src []byte) []byte {
	if len(src)%2 != 0 {
		return MustBCD(append(src, []byte("0")...))
	}
	return MustBCD(src)
}

// RightBCD ensures that the data as a valid size to be converted as a right-aligned Binary Coded Decimal.
// If the length is not a modulo of 2, we prefixes it with a zero.
func RightBCD(src []byte) []byte {
	if len(src)%2 != 0 {
		return MustBCD(append([]byte("0"), src...))
	}
	return MustBCD(src)
}

// X returns the hexadecimal encoding of src.
func X(src []byte) []byte {
	dst := make([]byte, hex.EncodedLen(len(src)))
	_ = hex.Encode(dst, src)
	return dst
}
