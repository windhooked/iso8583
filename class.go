// Copyright (c) 2019 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.package encoding

package iso8583

// Class is the purpose of the message.
type Class uint8

// List of message class.
const (
	// Authorization determine if funds are available aka x1xx.
	// It get an approval but do not post to account for reconciliation.
	Authorization = iota + 1
	// Financial determine if funds are available, get an approval and post directly to the account, aka x2xx.
	Financial
	// FileActions is used for hot-card, TMS and other exchanges, aka x3xx.
	FileActions
	// ReversalChargeBack reverses the action of a previous authorization, aka x4xx.
	ReversalChargeBack
	// Reconciliation transmits settlement information message  aka x5xx.
	Reconciliation
	// Administrative transmits administrative advice  aka x6xx.
	Administrative
	// FeeCollection ... (no doc)  aka x7xx.
	FeeCollection
	// NetworkManagement is used for secure key exchange, logon, echo test and other network functions aka x8xx.
	NetworkManagement
)

// Valid validates the Class.
func (c Class) Valid() bool {
	return c >= Authorization && c <= NetworkManagement
}
