// Copyright (c) 2019 Hervé Gouchet. All rights reserved.
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.package encoding

package iso8583

// Function is the function of the message.
type Function uint8

// List of functions.
const (
	// Request aka xx0x.
	Request = iota
	// RequestResponse aka xx1x.
	RequestResponse
	// Advice aka xx2x.
	Advice
	// AdviceResponse aka xx3x.
	AdviceResponse
	// Notification aka xx4x.
	Notification
	// NotificationAcknowledgement aka xx5x.
	NotificationAcknowledgement
	// Instruction aka xx6x.
	Instruction
	// Instruction aka xx7x.
	InstructionAcknowledgement
)

// Valid validates the Function.
func (f Function) Valid() bool {
	return f <= InstructionAcknowledgement
}
