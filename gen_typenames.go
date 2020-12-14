package bebop

func simpleGoString(simple string) string {
	if simple == "guid" {
		return "[16]byte"
	}
	if simple == "date" {
		return "time.Time"
	}
	return simple
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

func (mt MapType) goString() string {
	return "map[" + simpleGoString(mt.Key) + "]" + mt.Value.goString()
}
