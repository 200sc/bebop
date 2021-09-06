package bebop

func simpleGoString(simple string, settings GenerateSettings) string {
	if simple == typeGUID {
		return "[16]byte"
	}
	if simple == typeDate {
		return "time.Time"
	}
	if alias, ok := settings.importTypeAliases[simple]; ok {
		return alias
	}
	return simple
}

func (ft FieldType) goString(settings GenerateSettings) string {
	if ft.Map != nil {
		return "map[" + simpleGoString(ft.Map.Key, settings) + "]" + ft.Map.Value.goString(settings)
	}
	if ft.Array != nil {
		return "[]" + ft.Array.goString(settings)
	}
	return simpleGoString(ft.Simple, settings)
}

func (mt MapType) goString(settings GenerateSettings) string {
	return "map[" + simpleGoString(mt.Key, settings) + "]" + mt.Value.goString(settings)
}

func (uf UnionField) name() string {
	if uf.Message != nil {
		return uf.Message.Name
	} else {
		return uf.Struct.Name
	}
}
