package bebop

func (f File) customRecordTypes() map[string]struct{} {
	out := make(map[string]struct{})
	for _, st := range f.Structs {
		out[st.Name] = struct{}{}
	}
	for _, msg := range f.Messages {
		out[msg.Name] = struct{}{}
	}
	for _, union := range f.Unions {
		out[union.Name] = struct{}{}
		for _, ufd := range union.Fields {
			if ufd.Struct != nil {
				out[ufd.Struct.Name] = struct{}{}
			}
			if ufd.Message != nil {
				out[ufd.Message.Name] = struct{}{}
			}
		}
	}
	return out
}

func (f File) usedTypes() map[string]bool {
	out := make(map[string]bool)
	for _, st := range f.Structs {
		stOut := st.usedTypes()
		for k, v := range stOut {
			out[k] = v
		}
	}
	for _, msg := range f.Messages {
		msgOut := msg.usedTypes()
		for k, v := range msgOut {
			out[k] = v
		}
	}
	for _, union := range f.Unions {
		unionOut := union.usedTypes()
		for k, v := range unionOut {
			out[k] = v
		}
	}
	return out
}

func (st Struct) usedTypes() map[string]bool {
	out := make(map[string]bool)
	for _, fd := range st.Fields {
		fdTypes := fd.usedTypes()
		for k, v := range fdTypes {
			out[k] = v
		}
	}
	return out
}

func (msg Message) usedTypes() map[string]bool {
	out := make(map[string]bool)
	for _, fd := range msg.Fields {
		fdTypes := fd.usedTypes()
		for k, v := range fdTypes {
			out[k] = v
		}
	}
	return out
}

func (u Union) usedTypes() map[string]bool {
	out := make(map[string]bool)
	for _, ufd := range u.Fields {
		fdTypes := ufd.usedTypes()
		for k, v := range fdTypes {
			out[k] = v
		}
	}
	return out
}

func (ft FieldType) usedTypes() map[string]bool {
	if ft.Array != nil {
		return ft.Array.usedTypes()
	}
	if ft.Map != nil {
		valTypes := ft.Map.Value.usedTypes()
		valTypes[ft.Map.Key] = true
		return valTypes
	}
	return map[string]bool{ft.Simple: true}
}

func (ufd UnionField) usedTypes() map[string]bool {
	if ufd.Struct != nil {
		return ufd.Struct.usedTypes()
	}
	if ufd.Message != nil {
		return ufd.Message.usedTypes()
	}
	return map[string]bool{}
}
