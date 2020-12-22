// Package bebop provides structures, tokenizing, parsing, and code generation for the bebop file type
package bebop

// A File is a structured representation of a .bop file.
type File struct {
	Structs  []Struct
	Messages []Message
	Enums    []Enum
}

// A Struct is a record type where all fields are required.
type Struct struct {
	Name    string
	Comment string
	Fields  []Field
	// If OpCode is defined, wire encodings of the struct will be
	// preceded by the OpCode.
	OpCode int32
	// If ReadOnly is true, generated code for the struct will
	// provide field getters instead of exporting fields.
	ReadOnly bool
}

// A Field is an individual, typed data component making up
// a Struct or Message.
type Field struct {
	FieldType
	Name    string
	Comment string
	// DeprecatedMessage is only provided if Deprecated is true.
	DeprecatedMessage string
	Deprecated        bool
}

// A Message is a record type where all fields are optional and keyed to indices.
type Message struct {
	Name    string
	Comment string
	Fields  map[uint8]Field
	OpCode  int32
}

// An Enum is a definition that will generate typed enumerable options.
type Enum struct {
	Name    string
	Comment string
	Options []EnumOption
}

// An EnumOption is one possible value for a field typed as a specific Enum.
type EnumOption struct {
	Name    string
	Comment string
	// DeprecatedMessage is only provided if Deprecated is true.
	DeprecatedMessage string
	Value             int32
	Deprecated        bool
}

// A FieldType is a union of three choices: Simple types, array types, and map types.
// Only one of the three should be provided for a given FieldType.
type FieldType struct {
	Simple string
	Map    *MapType
	Array  *FieldType
}

// A MapType is a key-value type pair, where the key must be
// a primitive builtin type.
type MapType struct {
	// Keys may only be named types
	Key   string
	Value FieldType
}
