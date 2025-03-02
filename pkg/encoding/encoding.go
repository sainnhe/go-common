// Package encoding defines interfaces that convert data to and from byte-level and textual representations.
package encoding

// Type is the type of the encoding format.
type Type int

const (
	// TypeNil indicates no type.
	TypeNil Type = 0
	// TypeJSON is the JSON type.
	TypeJSON Type = 1
	// TypeYAML is the YAML type.
	TypeYAML Type = 2
)
