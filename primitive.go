package bebop

const (
	typeBool    = "bool"
	typeByte    = "byte"
	typeUint8   = "uint8"
	typeUint16  = "uint16"
	typeInt16   = "int16"
	typeUint32  = "uint32"
	typeInt32   = "int32"
	typeUint64  = "uint64"
	typeInt64   = "int64"
	typeFloat32 = "float32"
	typeFloat64 = "float64"
	typeString  = "string"
	typeGUID    = "guid"
	typeDate    = "date"
)

func isPrimitiveType(simpleType string) bool {
	_, ok := primitiveTypes[simpleType]
	return ok
}

var primitiveTypes = map[string]struct{}{
	typeBool:    {},
	typeByte:    {},
	typeUint8:   {},
	typeUint16:  {},
	typeInt16:   {},
	typeUint32:  {},
	typeInt32:   {},
	typeUint64:  {},
	typeInt64:   {},
	typeFloat32: {},
	typeFloat64: {},
	typeString:  {},
	typeGUID:    {},
	typeDate:    {},
}

func isFloatPrimitive(simpleType string) bool {
	_, ok := floatTypes[simpleType]
	return ok
}

var floatTypes = map[string]struct{}{
	typeFloat32: {},
	typeFloat64: {},
}

func isUintPrimitive(simpleType string) bool {
	_, ok := uintTypes[simpleType]
	return ok
}

var uintTypes = map[string]struct{}{
	typeByte:   {},
	typeUint8:  {},
	typeUint16: {},
	typeUint32: {},
	typeUint64: {},
}

func isIntPrimitive(simpleType string) bool {
	_, ok := intTypes[simpleType]
	return ok
}

var intTypes = map[string]struct{}{
	typeInt16: {},
	typeInt32: {},
	typeInt64: {},
}

func decodeIntegerType(typ string) (bitsize int, uintype bool) {
	switch typ {
	case typeByte, typeUint8:
		return 8, true
	case typeUint16:
		return 16, true
	case typeUint32:
		return 32, true
	case typeUint64:
		return 64, true
	case typeInt16:
		return 16, false
	case typeInt32:
		return 32, false
	case typeInt64:
		return 64, false
	default:
		panic("decodeIntegerType was provided a non-integer")
	}
}
