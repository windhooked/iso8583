// Copyright (c) 2019 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.package encoding

package iso8583_test

import (
	"strconv"
	"testing"

	"github.com/rvflash/iso8583/errors"

	are "github.com/matryer/is"
	"github.com/rvflash/iso8583"
)

func TestParseMTI(t *testing.T) {
	var (
		is = are.New(t)
		dt = []struct {
			in      string
			out     *iso8583.MTI
			literal string
			fail    bool
		}{
			{in: "100", out: iso8583.NewMTI(iso8583.V1987, iso8583.Authorization), literal: "0100"},
			{in: "120", out: iso8583.NewMTI(iso8583.V1987, iso8583.Authorization, iso8583.Advice), literal: "0120"},
			{in: "400", out: iso8583.NewMTI(iso8583.V1987, iso8583.ReversalChargeBack), literal: "0400"},
			{
				in:      "420",
				out:     iso8583.NewMTI(iso8583.V1987, iso8583.ReversalChargeBack, iso8583.Advice),
				literal: "0420",
			},
			{
				in:      "440",
				out:     iso8583.NewMTI(iso8583.V1987, iso8583.ReversalChargeBack, iso8583.Notification),
				literal: "0440",
			},
			{in: "1200", out: iso8583.NewMTI(iso8583.V1993, iso8583.Financial), literal: "1200"},
			{in: "1240", out: iso8583.NewMTI(iso8583.V1993, iso8583.Financial, iso8583.Notification), literal: "1240"},
			{in: "xx0x", fail: true},
		}
	)
	for i, tt := range dt {
		tt := tt
		t.Run("#"+strconv.Itoa(i), func(t *testing.T) {
			out, err := iso8583.ParseMTI(tt.in)
			if tt.fail {
				is.Equal(err, errors.MTI)
			} else {
				is.NoErr(err)
			}
			is.Equal(out, tt.out)
			if !tt.fail {
				is.Equal(out.String(), tt.literal)
			}
		})
	}
}
