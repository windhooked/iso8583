// Copyright (c) 2019 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package errors

import (
	"errors"
	"fmt"
)

// List of known errors.
var (
	// Field is returned if the data is invalid.
	Data = errors.New("invalid data")
	// Length is returned if the length not matches with the expected length.
	Length = errors.New("invalid length")
	// MTI is returned if we failed to fields the data.
	MTI = errors.New("invalid message type identifier")
	// NotImplemented is returned if the method is not implemented yet.
	NotImplemented = errors.New("not implemented")
	// OutOfRange is the data exceeds the bounds.
	OutOfRange = errors.New("out of range")
)

// New returns a new instance of a field error.
func New(err error, num int) error {
	return &Field{err: err, num: num}
}

// Field represents an error about the Field.
type Field struct {
	err error
	num int
}

// Error implements the error interface.
func (e *Field) Error() string {
	return fmt.Sprintf("field #%d: %s", e.num, e.err)
}
