// Package bebop provides structures, tokenizing, parsing, and code generation for the bebop file type
package bebop

// A File is a structured representation of a .bop file.
type File struct {
	Structs  []Struct
	Messages []Message
	Enums    []Enum
}

// Equals reports whether two Files are equivalent.
func (f File) Equals(f2 File) bool {
	if len(f.Structs) != len(f2.Structs) {
		return false
	}
	for i, st := range f.Structs {
		if !st.Equals(f2.Structs[i]) {
			return false
		}
	}
	if len(f.Messages) != len(f2.Messages) {
		return false
	}
	for i, msg := range f.Messages {
		if !msg.Equals(f2.Messages[i]) {
			return false
		}
	}
	if len(f.Enums) != len(f2.Enums) {
		return false
	}
	for i, en := range f.Enums {
		if !en.Equals(f2.Enums[i]) {
			return false
		}
	}
	return true
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

// Equals reports whether two Structs are equivalent.
func (s Struct) Equals(s2 Struct) bool {
	if s.Name != s2.Name {
		return false
	}
	if s.Comment != s2.Comment {
		return false
	}
	if s.OpCode != s2.OpCode {
		return false
	}
	if len(s.Fields) != len(s2.Fields) {
		return false
	}
	for i, fd := range s.Fields {
		if !fd.Equals(s2.Fields[i]) {
			return false
		}
	}
	return true
}

// A Field is an individual, typed data component making up
// a Struct or Message.
type Field struct {
	FieldType
	Name    string
	Comment string
	// DeprecatedMessage should only be non-empty if Deprecated is true.
	DeprecatedMessage string
	Deprecated        bool
}

// Equals reports whether two Fields are equivalent.
func (f Field) Equals(f2 Field) bool {
	if f.Name != f2.Name {
		return false
	}
	if f.Comment != f2.Comment {
		return false
	}
	if f.DeprecatedMessage != f2.DeprecatedMessage {
		return false
	}
	if f.Deprecated != f2.Deprecated {
		return false
	}
	return f.FieldType.Equals(f2.FieldType)
}

// A Message is a record type where all fields are optional and keyed to indices.
type Message struct {
	Name     string
	Comment  string
	Fields   map[uint8]Field
	OpCode   int32
	ReadOnly bool
}

// Equals reports whether two Messages are equivalent.
func (m Message) Equals(m2 Message) bool {
	if m.Name != m2.Name {
		return false
	}
	if m.Comment != m2.Comment {
		return false
	}
	if len(m.Fields) != len(m2.Fields) {
		return false
	}
	for key, fd := range m.Fields {
		if !fd.Equals(m2.Fields[key]) {
			return false
		}
	}
	for key, fd := range m2.Fields {
		if !fd.Equals(m.Fields[key]) {
			return false
		}
	}
	return true
}

// An Enum is a definition that will generate typed enumerable options.
type Enum struct {
	Name    string
	Comment string
	Options []EnumOption
}

// Equals reports whether two Enums are equivalent.
func (e Enum) Equals(e2 Enum) bool {
	if e.Name != e2.Name {
		return false
	}
	if e.Comment != e2.Comment {
		return false
	}
	if len(e.Options) != len(e2.Options) {
		return false
	}
	for i, o := range e.Options {
		if !o.Equals(e2.Options[i]) {
			return false
		}
	}
	return true
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

// Equals reports whether two EnumOptions are equivalent.
func (eo EnumOption) Equals(eo2 EnumOption) bool {
	return eo == eo2
}

// A FieldType is a union of three choices: Simple types, array types, and map types.
// Only one of the three should be provided for a given FieldType.
type FieldType struct {
	Simple string
	Map    *MapType
	Array  *FieldType
}

func (ft FieldType) goString() string {
	if ft.Map != nil {
		return "map[" + simpleGoString(ft.Map.Key) + "]" + ft.Map.Value.goString()
	}
	if ft.Array != nil {
		return "[]" + ft.Array.goString()
	}
	return simpleGoString(ft.Simple)
}

func simpleGoString(simple string) string {
	if simple == "guid" {
		return "[16]byte"
	}
	if simple == "date" {
		return "time.Time"
	}
	return simple
}

// Equals reports whether two FieldTypes are equivalent.
func (ft FieldType) Equals(ft2 FieldType) bool {
	if ft.Simple != ft2.Simple {
		return false
	}
	if (ft.Map == nil) != (ft2.Map == nil) {
		return false
	}
	if ft.Map != nil && !ft.Map.Equals(*ft2.Map) {
		return false
	}
	if (ft.Array == nil) != (ft2.Array == nil) {
		return false
	}
	if ft.Array != nil && !ft.Array.Equals(*ft2.Array) {
		return false
	}
	return true
}

// A MapType is a key-value type pair, where the key must be
// a primitive builtin type.
type MapType struct {
	// Keys may only be named types
	Key   string
	Value FieldType
}

// Equals reports whether two MapTypes are equivalent.
func (mt MapType) Equals(mt2 MapType) bool {
	if mt.Key != mt2.Key {
		return false
	}
	return mt.Value.Equals(mt2.Value)
}

func (mt MapType) goString() string {
	return "map[" + simpleGoString(mt.Key) + "]" + mt.Value.goString()
}
