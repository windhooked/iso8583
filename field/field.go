package field

import (
	"errors"
	"time"
)

// ErrData is returned if the data is invalid.
var ErrData = errors.New("invalid data")

// Format represents the type of data.
type Format uint16

// List of available formats.
const (
	Alpha     Format = 1 << iota // Alpha, including blanks (a).
	Amount                       // First byte are either 'C' to indicate a Credit or 'D' to a negative or Debit value.
	Binary                       // Binary data (b).
	Numeric                      // Numeric values only (n).
	Special                      // Special characters (s).
	Track                        // Track code (z), see ISO 7811 & ISO 7813
	Date                         // Date in format YYMMDD
	YearMonth                    // Date in format YYMM
	MonthDay                     // Date in format MMDD,
	Time                         // Time in format HHMMSS
)

// Type indicates if the data has variable or fixed length.
type Type uint8

// List of length indicator.
const (
	// Fixed length field of six digits (example: n 6).
	Fixed Type = iota
	// LVar numeric field of up to 6 digits in length (example: n.6).
	LVar
	// LLVar alpha field of up to 11 characters in length (example: a..11).
	LLVar
	// LLLVar binary field of up to 999 bits in length (example: b...999).
	LLLVar
)

// Elements is implemented by a Data.
type Elements interface {
	// Int64
	Int64() (int64, error)
	// String
	String() string
	// Int64
	Time() (time.Time, error)
	// Valid
	Valid() bool
}

// Field represents an ISO 8583 data field
type Field struct {
	Format Format
	Type   Type
	Size   int
}

// List all the known fields.
var n = map[uint8]Field{
	1:   {Format: Binary, Size: 64},                                   // Bitmap (128 if secondary or 192 if tertiary).
	2:   {Format: Numeric, Size: 19, Type: LLVar},                     // Primary account number (PAN)
	3:   {Format: Numeric, Size: 6},                                   // Processing code
	4:   {Format: Numeric, Size: 12},                                  // Amount, transaction
	5:   {Format: Numeric, Size: 12},                                  // Amount, settlement
	6:   {Format: Numeric, Size: 12},                                  // Amount, cardholder billing
	7:   {Format: Numeric | MonthDay | Time, Size: 10},                // Transmission date & time
	8:   {Format: Numeric, Size: 8},                                   // Amount, cardholder billing fee
	9:   {Format: Numeric, Size: 8},                                   // Conversion rate, settlement
	10:  {Format: Numeric, Size: 8},                                   // Conversion rate, cardholder billing
	11:  {Format: Numeric, Size: 6},                                   // System trace audit number (STAN)
	12:  {Format: Numeric | Time, Size: 6},                            // Time, local transaction (hhmmss)
	13:  {Format: Numeric | MonthDay, Size: 4},                        // Date, local transaction (MMDD)
	14:  {Format: Numeric | YearMonth, Size: 4},                       // Date, expiration
	15:  {Format: Numeric | MonthDay, Size: 4},                        // Date, settlement
	16:  {Format: Numeric | MonthDay, Size: 4},                        // Date, conversion
	17:  {Format: Numeric | MonthDay, Size: 4},                        // Date, capture
	18:  {Format: Numeric, Size: 4},                                   // Merchant type
	19:  {Format: Numeric, Size: 3},                                   // Acquiring institution country code
	20:  {Format: Numeric, Size: 3},                                   // PAN extended, country code
	21:  {Format: Numeric, Size: 3},                                   // Forwarding institution. country code
	22:  {Format: Numeric, Size: 3},                                   // Point of service entry mode
	23:  {Format: Numeric, Size: 3},                                   // Application PAN sequence number
	24:  {Format: Numeric, Size: 3},                                   // Network International identifier (NII)
	25:  {Format: Numeric, Size: 2},                                   // Point of service condition code
	26:  {Format: Numeric, Size: 2},                                   // Point of service capture code
	27:  {Format: Numeric, Size: 1},                                   // Authorizing identification response length
	28:  {Format: Amount | Numeric, Size: 8},                          // Amount, transaction fee
	29:  {Format: Amount | Numeric, Size: 8},                          // Amount, settlement fee
	30:  {Format: Amount | Numeric, Size: 8},                          // Amount, transaction processing fee
	31:  {Format: Amount | Numeric, Size: 8},                          // Amount, settlement processing fee
	32:  {Format: Numeric, Size: 11, Type: LLVar},                     // Acquiring institution identification code
	33:  {Format: Numeric, Size: 11, Type: LLVar},                     // Forwarding institution identification code
	34:  {Format: Numeric | Special, Size: 28, Type: LLVar},           // Primary account number, extended
	35:  {Format: Track, Size: 37, Type: LLVar},                       // Track 2 data
	36:  {Format: Numeric, Size: 104, Type: LLLVar},                   // Track 3 data
	37:  {Format: Alpha | Numeric, Size: 12},                          // Retrieval reference number
	38:  {Format: Alpha | Numeric, Size: 6},                           // Authorization identification response
	39:  {Format: Alpha | Numeric, Size: 2},                           // Response code
	40:  {Format: Alpha | Numeric, Size: 3},                           // Service restriction code
	41:  {Format: Alpha | Numeric | Special, Size: 8},                 // Card acceptor (CA) terminal identification
	42:  {Format: Alpha | Numeric | Special, Size: 15},                // CA identification code
	43:  {Format: Alpha | Numeric | Special, Size: 40},                // CA address: 1-23 +12: city +2: state +2: country
	44:  {Format: Alpha | Numeric, Size: 25, Type: LLVar},             // Additional response data
	45:  {Format: Alpha | Numeric, Size: 76, Type: LLVar},             // Track 1 data
	46:  {Format: Alpha | Numeric, Size: 999, Type: LLLVar},           // Additional data - ISO
	47:  {Format: Alpha | Numeric, Size: 999, Type: LLLVar},           // Additional data - national
	48:  {Format: Alpha | Numeric, Size: 999, Type: LLLVar},           // Additional data - private
	49:  {Format: Alpha, Size: 3},                                     // Currency code, transaction
	50:  {Format: Alpha, Size: 3},                                     // Currency code, settlement
	51:  {Format: Alpha, Size: 3},                                     // Currency code, cardholder billing
	52:  {Format: Binary, Size: 64},                                   //  data
	53:  {Format: Numeric, Size: 16},                                  // Security related control information
	54:  {Format: Alpha | Numeric, Size: 120, Type: LLLVar},           // Additional amounts
	55:  {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // ICC Elements - EMV having multiple tags
	56:  {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved ISO
	57:  {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved national
	58:  {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved national
	59:  {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved national
	60:  {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved national
	61:  {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved private
	62:  {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved private
	63:  {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved private
	64:  {Format: Binary, Size: 16},                                   // Message authentication code (MAC)
	65:  {Format: Binary, Size: 1},                                    // Bitmap, extended
	66:  {Format: Numeric, Size: 1},                                   // Settlement code
	67:  {Format: Numeric, Size: 2},                                   // Extended payment code
	68:  {Format: Numeric, Size: 3},                                   // Receiving institution country code
	69:  {Format: Numeric, Size: 3},                                   // Settlement institution country code
	70:  {Format: Numeric, Size: 3},                                   // Network management information code
	71:  {Format: Numeric, Size: 4},                                   // Message number
	72:  {Format: Numeric, Size: 4},                                   // Message number, last
	73:  {Format: Numeric | Date, Size: 6},                            // Date, action (YYMMDD)
	74:  {Format: Numeric, Size: 10},                                  // Credits, number
	75:  {Format: Numeric, Size: 10},                                  // Credits, reversal number
	76:  {Format: Numeric, Size: 10},                                  // Debits, number
	77:  {Format: Numeric, Size: 10},                                  // Debits, reversal number
	78:  {Format: Numeric, Size: 10},                                  // Transfer number
	79:  {Format: Numeric, Size: 10},                                  // Transfer, reversal number
	80:  {Format: Numeric, Size: 10},                                  // Inquiries number
	81:  {Format: Numeric, Size: 10},                                  // Authorizations, number
	82:  {Format: Numeric, Size: 12},                                  // Credits, processing fee amount
	83:  {Format: Numeric, Size: 12},                                  // Credits, transaction fee amount
	84:  {Format: Numeric, Size: 12},                                  // Debits, processing fee amount
	85:  {Format: Numeric, Size: 12},                                  // Debits, transaction fee amount
	86:  {Format: Numeric, Size: 16},                                  // Credits, amount
	87:  {Format: Numeric, Size: 16},                                  // Credits, reversal amount
	88:  {Format: Numeric, Size: 16},                                  // Debits, amount
	89:  {Format: Numeric, Size: 16},                                  // Debits, reversal amount
	90:  {Format: Numeric, Size: 42},                                  // Original data elements
	91:  {Format: Alpha | Numeric, Size: 1},                           // File update code
	92:  {Format: Alpha | Numeric, Size: 2},                           // File security code
	93:  {Format: Alpha | Numeric, Size: 5},                           // Response indicator
	94:  {Format: Alpha | Numeric, Size: 7},                           // Service indicator
	95:  {Format: Alpha | Numeric, Size: 42},                          // Replacement amounts
	96:  {Format: Binary, Size: 64},                                   // Message security code
	97:  {Format: Amount | Numeric, Size: 16},                         // Amount, net settlement
	98:  {Format: Alpha | Numeric | Special, Size: 25},                // Payee
	99:  {Format: Numeric, Size: 11, Type: LLVar},                     // Settlement institution identification code
	100: {Format: Numeric, Size: 11, Type: LLVar},                     // Receiving institution identification code
	101: {Format: Alpha | Numeric | Special, Size: 17, Type: LLVar},   // File name
	102: {Format: Alpha | Numeric | Special, Size: 28, Type: LLVar},   // Account identification 1
	103: {Format: Alpha | Numeric | Special, Size: 28, Type: LLVar},   // Account identification 2
	104: {Format: Alpha | Numeric | Special, Size: 100, Type: LLLVar}, // Transaction description
	105: {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved for ISO use
	106: {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved for ISO use
	107: {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved for ISO use
	108: {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved for ISO use
	109: {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved for ISO use
	110: {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved for ISO use
	111: {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved for ISO use
	112: {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved for national use
	113: {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved for national use
	114: {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved for national use
	115: {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved for national use
	116: {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved for national use
	117: {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved for national use
	118: {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved for national use
	119: {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved for national use
	120: {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved for private use
	121: {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved for private use
	122: {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved for private use
	123: {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved for private use
	124: {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved for private use
	125: {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved for private use
	126: {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved for private use
	127: {Format: Alpha | Numeric | Special, Size: 999, Type: LLLVar}, // Reserved for private use
	128: {Format: Binary, Size: 64},                                   // Message authentication code
}
