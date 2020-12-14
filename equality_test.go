package bebop

func (f File) equals(f2 File) bool {
	if len(f.Structs) != len(f2.Structs) {
		return false
	}
	for i, st := range f.Structs {
		if !st.equals(f2.Structs[i]) {
			return false
		}
	}
	if len(f.Messages) != len(f2.Messages) {
		return false
	}
	for i, msg := range f.Messages {
		if !msg.equals(f2.Messages[i]) {
			return false
		}
	}
	if len(f.Enums) != len(f2.Enums) {
		return false
	}
	for i, en := range f.Enums {
		if !en.equals(f2.Enums[i]) {
			return false
		}
	}
	return true
}

func (s Struct) equals(s2 Struct) bool {
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
		if !fd.equals(s2.Fields[i]) {
			return false
		}
	}
	return true
}

func (f Field) equals(f2 Field) bool {
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
	return f.FieldType.equals(f2.FieldType)
}

func (m Message) equals(m2 Message) bool {
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
		if !fd.equals(m2.Fields[key]) {
			return false
		}
	}
	for key, fd := range m2.Fields {
		if !fd.equals(m.Fields[key]) {
			return false
		}
	}
	return true
}

func (e Enum) equals(e2 Enum) bool {
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
		if !o.equals(e2.Options[i]) {
			return false
		}
	}
	return true
}

func (eo EnumOption) equals(eo2 EnumOption) bool {
	return eo == eo2
}

func (ft FieldType) equals(ft2 FieldType) bool {
	if ft.Simple != ft2.Simple {
		return false
	}
	if (ft.Map == nil) != (ft2.Map == nil) {
		return false
	}
	if ft.Map != nil && !ft.Map.equals(*ft2.Map) {
		return false
	}
	if (ft.Array == nil) != (ft2.Array == nil) {
		return false
	}
	if ft.Array != nil && !ft.Array.equals(*ft2.Array) {
		return false
	}
	return true
}

func (mt MapType) equals(mt2 MapType) bool {
	if mt.Key != mt2.Key {
		return false
	}
	return mt.Value.equals(mt2.Value)
}
