package iso8583

// Version represents the ISO 8583 version used.
type Version uint8

const (
	// V1987 is the first version from 1987, aka 0xxx
	V1987 = iota
	// V1993 is the second version from 1993, aka 1xxx
	V1993
	// V2003 is the third version from 2003, aka 2xxx
	V2003
)

// Valid validates the Version.
func (v Version) Valid() bool {
	return v <= V2003
}