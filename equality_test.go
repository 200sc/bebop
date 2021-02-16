package bebop

import "fmt"

func (f File) equals(f2 File) error {
	if len(f.Structs) != len(f2.Structs) {
		return fmt.Errorf("struct count mismatch: %v vs %v", len(f.Structs), len(f2.Structs))
	}
	for i, st := range f.Structs {
		if err := st.equals(f2.Structs[i]); err != nil {
			return fmt.Errorf("struct %d mismatched: %w", i, err)
		}
	}
	if len(f.Messages) != len(f2.Messages) {
		return fmt.Errorf("message count mismatch: %v vs %v", len(f.Messages), len(f2.Messages))
	}
	for i, msg := range f.Messages {
		if err := msg.equals(f2.Messages[i]); err != nil {
			return fmt.Errorf("message %d mismatched: %w", i, err)
		}
	}
	if len(f.Enums) != len(f2.Enums) {
		return fmt.Errorf("enum count mismatch: %v vs %v", len(f.Enums), len(f2.Enums))
	}
	for i, en := range f.Enums {
		if err := en.equals(f2.Enums[i]); err != nil {
			return fmt.Errorf("enum %d mismatched: %w", i, err)
		}
	}
	if len(f.Unions) != len(f2.Unions) {
		return fmt.Errorf("union count mismatch: %v vs %v", len(f.Unions), len(f2.Unions))
	}
	for i, union := range f.Unions {
		if err := union.equals(f2.Unions[i]); err != nil {
			return fmt.Errorf("union %d mismatched: %w", i, err)
		}
	}
	return nil
}

func (s Struct) equals(s2 Struct) (err error) {
	if s.Name != s2.Name {
		return fmt.Errorf("name mismatch: %v vs %v", s.Name, s2.Name)
	}
	if s.Comment != s2.Comment {
		return fmt.Errorf("comment mismatch: %q vs %q", s.Comment, s2.Comment)
	}
	if s.OpCode != s2.OpCode {
		return fmt.Errorf("opcode mismatch: %v vs %v", s.OpCode, s2.OpCode)
	}
	if len(s.Fields) != len(s2.Fields) {
		return fmt.Errorf("field count mismatch: %v vs %v", len(s.Fields), len(s2.Fields))
	}
	for i, fd := range s.Fields {
		if err := fd.equals(s2.Fields[i]); err != nil {
			return fmt.Errorf("field %d mismatched: %v", i, err)
		}
	}
	return nil
}

func (f Field) equals(f2 Field) error {
	if f.Name != f2.Name {
		return fmt.Errorf("name mismatch: %v vs %v", f.Name, f2.Name)
	}
	if f.Comment != f2.Comment {
		return fmt.Errorf("comment mismatch: %q vs %q", f.Comment, f2.Comment)
	}
	if f.DeprecatedMessage != f2.DeprecatedMessage {
		return fmt.Errorf("deprecated message mismatch: %v vs %v", f.DeprecatedMessage, f2.DeprecatedMessage)
	}
	if f.Deprecated != f2.Deprecated {
		return fmt.Errorf("deprecation mismatch: %v vs %v", f.Deprecated, f2.Deprecated)
	}
	return f.FieldType.equals(f2.FieldType)
}

func (m Message) equals(m2 Message) error {
	if m.Name != m2.Name {
		return fmt.Errorf("name mismatch: %v vs %v", m.Name, m2.Name)
	}
	if m.Comment != m2.Comment {
		return fmt.Errorf("comment mismatch: %q vs %q", m.Comment, m2.Comment)
	}
	if len(m.Fields) != len(m2.Fields) {
		return fmt.Errorf("field count mismatch: %v vs %v", len(m.Fields), len(m2.Fields))
	}
	for key, fd := range m.Fields {
		if err := fd.equals(m2.Fields[key]); err != nil {
			return fmt.Errorf("field %d mismatched: %v", key, err)
		}
	}
	for key, fd := range m2.Fields {
		if err := fd.equals(m.Fields[key]); err != nil {
			return fmt.Errorf("field %d mismatched: %v", key, err)
		}
	}
	return nil
}

func (e Enum) equals(e2 Enum) error {
	if e.Name != e2.Name {
		return fmt.Errorf("name mismatch: %v vs %v", e.Name, e2.Name)
	}
	if e.Comment != e2.Comment {
		return fmt.Errorf("comment mismatch: %q vs %q", e.Comment, e2.Comment)
	}
	if len(e.Options) != len(e2.Options) {
		return fmt.Errorf("option count mismatch: %v vs %v", len(e.Options), len(e2.Options))
	}
	for i, o := range e.Options {
		if err := o.equals(e2.Options[i]); err != nil {
			return fmt.Errorf("option %d mismatched: %v", i, err)
		}
	}
	return nil
}

func (eo EnumOption) equals(eo2 EnumOption) error {
	if eo.Name != eo2.Name {
		return fmt.Errorf("name mismatch: %v vs %v", eo.Name, eo2.Name)
	}
	if eo.Comment != eo2.Comment {
		return fmt.Errorf("comment mismatch: %q vs %q", eo.Comment, eo2.Comment)
	}
	if eo.DeprecatedMessage != eo2.DeprecatedMessage {
		return fmt.Errorf("deprecated message mismatch: %v vs %v", eo.DeprecatedMessage, eo2.DeprecatedMessage)
	}
	if eo.Deprecated != eo2.Deprecated {
		return fmt.Errorf("deprecated mismatch: %v vs %v", eo.Deprecated, eo2.Deprecated)
	}
	if eo.Value != eo2.Value {
		return fmt.Errorf("value mismatch: %v vs %v", eo.Value, eo2.Value)
	}
	return nil
}

func (ft FieldType) equals(ft2 FieldType) error {
	if ft.Simple != ft2.Simple {
		return fmt.Errorf("field is simple mismatch: %v vs %v", ft.Simple, ft2.Simple)
	}
	if (ft.Map == nil) != (ft2.Map == nil) {
		return fmt.Errorf("field is map type mismatch: %v vs %v", ft.Map != nil, ft2.Map != nil)
	}
	if ft.Map != nil {
		if err := ft.Map.equals(*ft2.Map); err != nil {
			return fmt.Errorf("map type mismatch: %v", err)
		}
	}
	if (ft.Array == nil) != (ft2.Array == nil) {
		return fmt.Errorf("field is array type mismatch: %v vs %v", ft.Array != nil, ft2.Array != nil)
	}
	if ft.Array != nil {
		if err := ft.Array.equals(*ft2.Array); err != nil {
			return fmt.Errorf("array type mismatch: %v", err)
		}
	}
	return nil
}

func (mt MapType) equals(mt2 MapType) error {
	if mt.Key != mt2.Key {
		return fmt.Errorf("key mismatch: %v vs %v", mt.Key, mt2.Key)
	}
	if mt.Value != mt2.Value {
		return fmt.Errorf("value mismatch: %v vs %v", mt.Value, mt2.Value)
	}
	return nil
}

func (u Union) equals(u2 Union) error {
	if u.Name != u2.Name {
		return fmt.Errorf("name mismatch: %v vs %v", u.Name, u2.Name)
	}
	if u.Comment != u2.Comment {
		return fmt.Errorf("comment mismatch: %q vs %q", u.Comment, u2.Comment)
	}
	if u.OpCode != u2.OpCode {
		return fmt.Errorf("opcode mismatch: %v vs %v", u.OpCode, u2.OpCode)
	}
	if len(u.Fields) != len(u2.Fields) {
		return fmt.Errorf("field count mismatch: %v vs %v", len(u.Fields), len(u2.Fields))
	}
	for key, fd := range u.Fields {
		if err := fd.equals(u2.Fields[key]); err != nil {
			return fmt.Errorf("field %d mismatch: %w", key, err)
		}
	}
	for key, fd := range u2.Fields {
		if err := fd.equals(u.Fields[key]); err != nil {
			return fmt.Errorf("field %d mismatch: %w", key, err)
		}
	}
	return nil
}

func (uf UnionField) equals(uf2 UnionField) error {
	if uf.Deprecated != uf2.Deprecated {
		return fmt.Errorf("deprecated mismatch: %v vs %v", uf.Deprecated, uf2.Deprecated)
	}
	if uf.DeprecatedMessage != uf2.DeprecatedMessage {
		return fmt.Errorf("deprecated message mismatch: %v vs %v", uf.DeprecatedMessage, uf2.DeprecatedMessage)
	}
	if (uf.Struct == nil) != (uf2.Struct == nil) {
		return fmt.Errorf("field is struct type mismatch: %v vs %v", uf.Struct != nil, uf2.Struct != nil)
	}
	if uf.Struct != nil && uf2.Struct != nil {
		return uf.Struct.equals(*uf2.Struct)
	}
	if (uf.Union == nil) != (uf2.Union == nil) {
		return fmt.Errorf("field is union type mismatch: %v vs %v", uf.Union != nil, uf2.Union != nil)
	}
	if uf.Union != nil && uf2.Union != nil {
		return uf.Union.equals(*uf2.Union)
	}
	if (uf.Message == nil) != (uf2.Message == nil) {
		return fmt.Errorf("field is message type mismatch: %v vs %v", uf.Message != nil, uf2.Message != nil)
	}
	if uf.Message != nil && uf2.Message != nil {
		return uf.Message.equals(*uf2.Message)
	}
	return nil
}
