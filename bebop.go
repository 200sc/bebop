package bebop

type File struct {
	Structs  []Struct
	Messages []Message
	Enums    []Enum
}

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

type Struct struct {
	Name     string
	Comment  string
	OpCode   int32
	Fields   []Field
	ReadOnly bool
}

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

type Field struct {
	FieldType
	Name    string
	Comment string
	// DeprecatedMessage is only provided if Deprecated is true.
	DeprecatedMessage string
	Deprecated        bool
}

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

type Message struct {
	Name     string
	Comment  string
	OpCode   int32
	Fields   map[uint8]Field
	ReadOnly bool
}

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

type Enum struct {
	Name    string
	Comment string
	Options []EnumOption
}

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

type EnumOption struct {
	Name    string
	Comment string
	Value   int32
	// DeprecatedMessage is only provided if Deprecated is true.
	DeprecatedMessage string
	Deprecated        bool
}

func (eo EnumOption) Equals(eo2 EnumOption) bool {
	return eo == eo2
}

type FieldType struct {
	Simple string
	Map    *MapType
	Array  *FieldType
	// TODO: fight
	// remove the array[t] alias
	// it's added complexity and hurts
}

func (ft FieldType) GoString() string {
	if ft.Map != nil {
		return "map[" + simpleGoString(ft.Map.Key) + "]" + ft.Map.Value.GoString()
	}
	if ft.Array != nil {
		return "[]" + ft.Array.GoString()
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

func (ft FieldType) IsMap() bool {
	return ft.Map != nil
}

func (ft FieldType) IsArray() bool {
	return ft.Array != nil
}

type ArrayType struct {
	Value FieldType
}

func (at ArrayType) Equals(at2 ArrayType) bool {
	return at.Value.Equals(at2.Value)
}

type MapType struct {
	// Keys may only be named types
	Key   string
	Value FieldType
}

func (mt MapType) Equals(mt2 MapType) bool {
	if mt.Key != mt2.Key {
		return false
	}
	return mt.Value.Equals(mt2.Value)
}

func (mt MapType) GoString() string {
	return "map[" + simpleGoString(mt.Key) + "]" + mt.Value.GoString()
}
