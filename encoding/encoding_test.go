// Copyright (c) 2019 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found fmt the LICENSE file.

package encoding_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/matryer/is"
	"github.com/rvflash/iso8583/encoding"
)

const (
	txt    = `Rv`
	txtHex = `5276`
	dec    = `45`
	decBin = `0011010000110101` //  with leading 0
	hex    = `A020000000800010`
	hexBin = `0100000100110000001100100011000000110000001100000011000000110000` +
		`0011000000110000001110000011000000110000001100000011000100110000`
)

func TestBinary(t *testing.T) {
	var (
		are = is.New(t)
		dt  = []struct {
			in, out []byte
		}{
			{in: []byte(hex), out: []byte(hexBin)},
			{in: []byte(dec), out: []byte(decBin)},
		}
	)
	var out []byte
	for i, tt := range dt {
		tt := tt
		t.Run("#"+strconv.Itoa(i), func(t *testing.T) {
			out = encoding.Binary(tt.in)
			are.Equal(out, tt.out)
		})
	}
}

func TestLeftBCD(t *testing.T) {
	var (
		are = is.New(t)
		dt  = []struct {
			in, left, right string
		}{
			{in: "45", left: "45", right: "45"},
			{in: "456", left: "4560", right: "0456"},
			{in: "0", left: "00", right: "00"},
			{in: "456ab8", left: "456AB8", right: "456AB8"},
		}
	)
	for i, tt := range dt {
		tt := tt
		t.Run("#"+strconv.Itoa(i), func(t *testing.T) {
			are.Equal(fmt.Sprintf("%X", encoding.LeftBCD([]byte(tt.in))), tt.left)
			are.Equal(fmt.Sprintf("%X", encoding.RightBCD([]byte(tt.in))), tt.right)
		})
	}
}

func TestLeftBCD2(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	_ = encoding.LeftBCD([]byte(txt))
}

func TestX(t *testing.T) {
	is.New(t).Equal(encoding.X([]byte(txt)), []byte(txtHex))
}
