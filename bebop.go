// Package bebop provides structures, tokenizing, parsing, and code generation for the bebop file type
package bebop

// Version is the library version. Should be used by CLI tools when passed a '--version' flag.
const Version = "v0.3.3"

// A File is a structured representation of a .bop file.
type File struct {
	// FileName is an optional argument defining where this
	// bebop file came from. This argument is only used to
	// determine where relative import files lie. If relative
	// imports are not used, this argument is not read. If
	// FileName is a relative path, it will be treated as
	// relative to os.Getwd().
	FileName string

	// GoPackage is the value of this file's go_package const,
	// should it be defined and string-typed.
	GoPackage string

	Structs  []Struct
	Messages []Message
	Enums    []Enum
	Unions   []Union
	Consts   []Const
	Imports  []string
}

// goPackage defines the constant used in bebop files as a hint to
// our compiler for which package a file should belong to. E.g.
// defining 'const go_package = github.com/user/repo/schema' will
// cause the file to define itself under the "schema" package and
// other bebop files will import it as github.com/user/repo/schema.
const goPackage = "go_package"

// A Struct is a record type where all fields are required.
type Struct struct {
	Name    string
	Comment string
	Fields  []Field
	// If OpCode is defined, wire encodings of the struct can be
	// preceded by the OpCode.
	OpCode uint32
	// Namespace is only provided for imported types, and only
	// used in code generation.
	Namespace string
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
	// Tags are not written by default, and must be enabled via a compiler flag.
	Tags []Tag
	// DeprecatedMessage is only provided if Deprecated is true.
	DeprecatedMessage string
	Deprecated        bool
}

// A Message is a record type where all fields are optional and keyed to indices.
type Message struct {
	Name    string
	Comment string
	Fields  map[uint8]Field
	OpCode  uint32
	// Namespace is only provided for imported types, and only
	// used in code generation.
	Namespace string
}

// A Union is like a message where explicitly one field will be provided.
type Union struct {
	Name    string
	Comment string
	Fields  map[uint8]UnionField
	OpCode  uint32
	// Namespace is only provided for imported types, and only
	// used in code generation.
	Namespace string
}

// A UnionField is either a Message, Struct, or Union, defined inline.
type UnionField struct {
	Message *Message
	Struct  *Struct
	// Tags are not written by default, ard must be enabled via a compiler flag.
	Tags []Tag
	// DeprecatedMessage is only provided if Deprecated is true.
	DeprecatedMessage string
	Deprecated        bool
}

// An Enum is a definition that will generate typed enumerable options.
type Enum struct {
	Name    string
	Comment string
	Options []EnumOption
	// Namespace is only provided for imported types, and only
	// used in code generation.
	Namespace string
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

// A const is a simple type - value pair that is compiled as
// a constant into generated code.
type Const struct {
	// Consts do not support map or array (or record) types
	SimpleType string
	Comment    string
	Name       string
	Value      string
}

// A Tag is a Go struct field tag, e.g. `json:"userId,omitempty"`
type Tag struct {
	Key   string
	Value string
	// Boolean is set if Value is empty, in the form `key`, not `key:""`.
	Boolean bool
}
