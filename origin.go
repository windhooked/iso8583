// Copyright (c) 2019 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.package encoding

package iso8583

// Origin defines the message source within the payment chain.
type Origin uint8

const (
	// Acquirer aka xxx0.
	Acquirer = iota
	// AcquirerRepeat aka xxx1.
	AcquirerRepeat
	// Issuer aka xxx2.
	Issuer
	// IssuerRepeat aka xxx3.
	IssuerRepeat
	// Other aka xxx4.
	Other
	// OtherRepeat aka xxx5.
	OtherRepeat
)

// Valid validates the Origin.
func (o Origin) Valid() bool {
	return o <= OtherRepeat
}
