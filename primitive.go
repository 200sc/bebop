package bebop

func isPrimitiveType(simpleType string) bool {
	_, ok := primitiveTypes[simpleType]
	return ok
}

var primitiveTypes = map[string]struct{}{
	"bool":    {},
	"byte":    {},
	"uint8":   {},
	"uint16":  {},
	"int16":   {},
	"uint32":  {},
	"int32":   {},
	"uint64":  {},
	"int64":   {},
	"float32": {},
	"float64": {},
	"string":  {},
	"guid":    {},
	"date":    {},
}
