// Copyright (c) 2019 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.package encoding

package iso8583_test

import (
	"fmt"
	"testing"

	"github.com/matryer/is"
	"github.com/rvflash/iso8583"
)

func TestUnmarshal(t *testing.T) {
	are := is.New(t)
	iso := new(iso8583.Message)
	data := []byte("0800823A0000000000000400000000000000042009061390000109061304200420001")
	err := iso8583.Unmarshal(data, iso)
	are.NoErr(err)
	fmt.Println(iso.Fields)
}
