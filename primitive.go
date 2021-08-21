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

func isFloatPrimitive(simpleType string) bool {
	_, ok := floatTypes[simpleType]
	return ok
}

var floatTypes = map[string]struct{}{
	"float32": {},
	"float64": {},
}

func isUintPrimitive(simpleType string) bool {
	_, ok := uintTypes[simpleType]
	return ok
}

var uintTypes = map[string]struct{}{
	"byte":   {},
	"uint8":  {},
	"uint16": {},
	"uint32": {},
	"uint64": {},
}

func isIntPrimitive(simpleType string) bool {
	_, ok := intTypes[simpleType]
	return ok
}

var intTypes = map[string]struct{}{
	"int16": {},
	"int32": {},
	"int64": {},
}
