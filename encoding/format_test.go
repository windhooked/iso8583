// Copyright (c) 2019 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found fmt the LICENSE file.

package encoding_test

import (
	"strconv"
	"testing"

	"github.com/matryer/is"
	"github.com/rvflash/iso8583/encoding"
	"github.com/rvflash/iso8583/errors"
)

func TestParse(t *testing.T) {
	var (
		are = is.New(t)
		dt  = []struct {
			in  string
			out encoding.Format
			err error
		}{
			{in: "ascii", out: encoding.ASCII},
			{in: "ASCII", out: encoding.ASCII},
			{in: "bcd", out: encoding.BCD},
			{in: "BCD", out: encoding.BCD},
			{in: txt, err: errors.NotImplemented},
		}
	)
	for i, tt := range dt {
		tt := tt
		t.Run("#"+strconv.Itoa(i), func(t *testing.T) {
			out, err := encoding.Parse(tt.in)
			are.Equal(err, tt.err)
			are.Equal(out, tt.out)
		})
	}
}

func TestFormat(t *testing.T) {
	var (
		are = is.New(t)
		dt  = []struct {
			fmt                 encoding.Format
			bitmap, header, mti int
		}{
			{fmt: encoding.ASCII, bitmap: 16, header: 4, mti: 4},
			{fmt: encoding.BCD, bitmap: 16, header: 2, mti: 4},
			{fmt: 255},
		}
	)
	for i, tt := range dt {
		tt := tt
		t.Run("#"+strconv.Itoa(i), func(t *testing.T) {
			are.Equal(tt.fmt.LenBitmap(), tt.bitmap)
			are.Equal(tt.fmt.LenHeader(), tt.header)
			are.Equal(tt.fmt.LenMTI(), tt.mti)
		})
	}
}

func TestFormat_EncodeToDecimal(t *testing.T) {
	var (
		are = is.New(t)
		dt  = []struct {
			fmt encoding.Format
			in  string
			out uint64
			err error
		}{
			{fmt: encoding.ASCII, in: "10", out: 10},
			{fmt: encoding.BCD, in: "10", out: 3130},
			{fmt: 255, err: errors.NotImplemented},
		}
	)
	for i, tt := range dt {
		tt := tt
		t.Run("#"+strconv.Itoa(i), func(t *testing.T) {
			out, err := tt.fmt.EncodeToDecimal([]byte(tt.in))
			are.Equal(err, tt.err)
			are.Equal(out, tt.out)
		})
	}
}
