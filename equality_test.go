package bebop

import (
	"fmt"
)

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
	if len(f.Consts) != len(f2.Consts) {
		return fmt.Errorf("const count mismatch: %v vs %v", len(f.Consts), len(f2.Consts))
	}
	for i, cons := range f.Consts {
		if err := cons.equals(f2.Consts[i]); err != nil {
			return fmt.Errorf("const %d mismatched: %w", i, err)
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
	if err := s.Decorations.equals(s2.Decorations); err != nil {
		return fmt.Errorf("decorations mismatched: %v", err)
	}
	return nil
}

func (d Decorations) equals(d2 Decorations) (err error) {
	if d.Deprecated != d2.Deprecated {
		return fmt.Errorf("deprecation mismatch: %v vs %v", d.Deprecated, d2.Deprecated)
	}
	if d.DeprecatedMessage != d2.DeprecatedMessage {
		return fmt.Errorf("deprecated message mismatch: %v vs %v", d.DeprecatedMessage, d2.DeprecatedMessage)
	}
	if len(d.Custom) != len(d2.Custom) {
		return fmt.Errorf("custom decorations mismatch: %v vs %v", len(d.Custom), len(d2.Custom))
	}
	for k, v := range d.Custom {
		v2 := d2.Custom[k]
		if v != v2 {
			return fmt.Errorf("decoration %v mismatch: %v vs %v", k, v, v2)
		}
	}
	for k, v := range d2.Custom {
		v2 := d.Custom[k]
		if v != v2 {
			return fmt.Errorf("decoration %v mismatch: %v vs %v", k, v, v2)
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
	if err := f.Decorations.equals(f2.Decorations); err != nil {
		return fmt.Errorf("decorations mismatched: %v", err)
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
	if m.OpCode != m2.OpCode {
		return fmt.Errorf("opcode mismatch: %v vs %v", m.OpCode, m2.OpCode)
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
	if err := m.Decorations.equals(m2.Decorations); err != nil {
		return fmt.Errorf("decorations mismatched: %v", err)
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
	if e.SimpleType != e2.SimpleType {
		return fmt.Errorf("simple type mismatch: %q vs %q", e.SimpleType, e2.SimpleType)
	}
	if e.Unsigned != e2.Unsigned {
		return fmt.Errorf("unsigned mismatch: %v vs %v", e.Unsigned, e2.Unsigned)
	}
	if err := e.Decorations.equals(e2.Decorations); err != nil {
		return fmt.Errorf("decorations mismatched: %v", err)
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
	if eo.Value != eo2.Value {
		return fmt.Errorf("value mismatch: %v vs %v", eo.Value, eo2.Value)
	}
	if err := eo.Decorations.equals(eo2.Decorations); err != nil {
		return fmt.Errorf("decorations mismatched: %v", err)
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
	if err := mt.Value.equals(mt2.Value); err != nil {
		return fmt.Errorf("value mismatch: %v", err)
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
	if err := u.Decorations.equals(u2.Decorations); err != nil {
		return fmt.Errorf("decorations mismatched: %v", err)
	}
	return nil
}

func (uf UnionField) equals(uf2 UnionField) error {
	if err := uf.Decorations.equals(uf2.Decorations); err != nil {
		return fmt.Errorf("decorations mismatched: %v", err)
	}
	if (uf.Struct == nil) != (uf2.Struct == nil) {
		return fmt.Errorf("field is struct type mismatch: %v vs %v", uf.Struct != nil, uf2.Struct != nil)
	}
	if uf.Struct != nil && uf2.Struct != nil {
		return uf.Struct.equals(*uf2.Struct)
	}
	if (uf.Message == nil) != (uf2.Message == nil) {
		return fmt.Errorf("field is message type mismatch: %v vs %v", uf.Message != nil, uf2.Message != nil)
	}
	if uf.Message != nil && uf2.Message != nil {
		return uf.Message.equals(*uf2.Message)
	}
	return nil
}

func (c Const) equals(c2 Const) error {
	if c.SimpleType != c2.SimpleType {
		return fmt.Errorf("simple type mismatch: %v vs %v", c.SimpleType, c2.SimpleType)
	}
	if c.Name != c2.Name {
		return fmt.Errorf("name mismatch: %v vs %v", c.Name, c2.Name)
	}
	if c.Value != c2.Value {
		return fmt.Errorf("value mismatch: %v vs %v", c.Value, c2.Value)
	}
	if err := c.Decorations.equals(c2.Decorations); err != nil {
		return fmt.Errorf("decorations mismatched: %v", err)
	}
	return nil
}
