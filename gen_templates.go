package bebop

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

func writeLineWithTabs(w io.Writer, format string, depth int, args ...string) {
	var assigner string
	var receiver string
	var typename string
	var namespace string
	var baretype string
	if len(args) > 0 {
		assigner = args[0]
		if assigner[0] == '&' || assigner[0] == '*' {
			receiver = assigner[1:]
		} else {
			if assigner[0] == '(' {
				receiver = "(*" + assigner[1:]
			} else {
				receiver = "*" + assigner
			}
		}
	}
	if len(args) > 1 {
		typename = args[1]
		if strings.Contains(typename, ".") {
			splitType := strings.Split(typename, ".")
			if len(splitType) == 2 {
				namespace = splitType[0]
				baretype = splitType[1]
			}
		}
	}
	// add tabs
	tbs := strings.Repeat("\t", depth)
	format = tbs + format
	format = strings.Replace(format, "\n", "\n"+tbs, -1)

	format = strings.Replace(format, fillReciever, receiver, -1)
	format = strings.Replace(format, fillAssigner, assigner, -1)
	format = strings.Replace(format, fillTypename, typename, -1)
	format = strings.Replace(format, fillBareType, baretype, -1)
	format = strings.Replace(format, fillNamespace, namespace, -1)
	format = strings.Replace(format, fillKey, depthName("k", depth), -1)
	format = strings.Replace(format, fillValue, depthName("v", depth), -1)

	fmt.Fprint(w, format+"\n")
}

const (
	fillAssigner = "%ASGN"
	fillReciever = "%RECV"
	fillTypename = "%TYPE"
	// type name without namepsace
	fillBareType  = "%BARETYPE"
	fillNamespace = "%NAMESPACE"
	fillKey       = "%KNAME"
	fillValue     = "%VNAME"

	fmtErrReturn            = "if err != nil {\n\treturn err\n}"
	fmtAddSizeToAt          = "at += (%ASGN).Size()"
	fmtAdd4PlusLenToAt      = "at += 4 + len(%ASGN)"
	fmtAdd4ToAt             = "at += 4"
	fmtAddSizeToBodyLen     = "bodyLen += (%ASGN).Size()"
	fmtAdd4PlusLenToBodyLen = "bodyLen += 4 + len(%ASGN)"
	fmtAdd4ToBodyLen        = "bodyLen += 4"

	fmtMakeType           = "(%RECV), err = Make%TYPE(r)\n" + fmtErrReturn
	fmtMakeNamespacedType = "(%RECV), err = %NAMESPACE.Make%BARETYPE(r)\n" + fmtErrReturn
	fmtMakePrivateType    = "(%RECV), err = make%TYPE(r)\n" + fmtErrReturn

	fmtMake           = "%ASGN, err = Make%TYPEFromBytes(buf[at:])\n"
	fmtMakeNamespaced = "%ASGN, err = %NAMESPACE.Make%BARETYPEFromBytes(buf[at:])\n"
	fmtMakePrivate    = "%ASGN, err = make%TYPEFromBytes(buf[at:])\n"

	fmtMustMake           = "%ASGN = MustMake%TYPEFromBytes(buf[at:])\n"
	fmtMustMakeNamespaced = "%ASGN = %NAMESPACE.MustMake%BARETYPEFromBytes(buf[at:])\n"
	fmtMustMakePrivate    = "%ASGN = mustMake%TYPEFromBytes(buf[at:])\n"

	fmtMarshal = "(%ASGN).MarshalBebopTo(buf[at:])\n"
	fmtEncode  = "err = (%ASGN).EncodeBebop(w)\n"
)

var fixedSizeTypes = map[string]uint8{
	typeBool:    1,
	typeByte:    1,
	typeUint8:   1,
	typeUint16:  2,
	typeInt16:   2,
	typeUint32:  4,
	typeInt32:   4,
	typeUint64:  8,
	typeInt64:   8,
	typeFloat32: 4,
	typeFloat64: 8,
	typeGUID:    16,
	typeDate:    8,
}

func fixedTitleString(typ string) string {
	if typ == typeGUID {
		return "GUID"
	}
	return strings.Title(typ)
}

func makeFormatType(namespace string, settings GenerateSettings) string {
	// namespacing takes precedence over privacy-- its assumed that
	// a namespaced import will not be compiled as private.
	if namespace != "" {
		return fmtMakeNamespacedType
	}
	if settings.PrivateDefinitions {
		return fmtMakePrivateType
	}
	return fmtMakeType
}

func makeFormat(namespace string, settings GenerateSettings) string {
	if namespace != "" {
		return fmtMakeNamespaced
	}
	if settings.PrivateDefinitions {
		return fmtMakePrivate
	}
	return fmtMake
}

func mustMakeFormat(namespace string, settings GenerateSettings) string {
	if namespace != "" {
		return fmtMustMakeNamespaced
	}
	if settings.PrivateDefinitions {
		return fmtMustMakePrivate
	}
	return fmtMustMake
}

func (f File) typeUnmarshallers(settings GenerateSettings) map[string]string {
	out := make(map[string]string)
	for typ := range fixedSizeTypes {
		out[typ] = "%RECV = iohelp.Read" + fixedTitleString(typ) + "(r)"
	}
	out[typeString] = "%RECV = iohelp.ReadString(r)"
	for _, en := range f.Enums {
		out[en.Name] = "%RECV = %TYPE(iohelp.ReadUint32(r))"
	}
	for _, st := range f.Structs {
		out[st.Name] = makeFormatType(st.Namespace, settings)
	}
	for _, msg := range f.Messages {
		out[msg.Name] = makeFormatType(msg.Namespace, settings)
	}
	for _, union := range f.Unions {
		uout := union.typeUnmarshallers(settings)
		for k, v := range uout {
			out[k] = v
		}
	}
	return out
}

func (u Union) typeUnmarshallers(settings GenerateSettings) map[string]string {
	out := make(map[string]string, 1+len(u.Fields))
	out[u.Name] = makeFormatType(u.Namespace, settings)
	for _, ufd := range u.Fields {
		if ufd.Struct != nil {
			out[ufd.Struct.Name] = makeFormatType(ufd.Struct.Namespace, settings)
		}
		if ufd.Message != nil {
			out[ufd.Message.Name] = makeFormatType(ufd.Message.Namespace, settings)
		}
	}
	return out
}

func (f File) typeMarshallers() map[string]string {
	out := make(map[string]string)
	for typ := range fixedSizeTypes {
		out[typ] = "iohelp.Write" + fixedTitleString(typ) + "(w, %ASGN)"
	}
	out[typeString] = "iohelp.WriteUint32(w, uint32(len(%ASGN)))\n" +
		"w.Write([]byte(%ASGN))"
	out[typeDate] = "if %ASGN != (time.Time{}) {\n" +
		"\tiohelp.WriteInt64(w, ((%ASGN).UnixNano() / 100))\n" +
		"} else {\n" +
		"\tiohelp.WriteInt64(w, 0)\n" +
		"}"
	for _, en := range f.Enums {
		out[en.Name] = "iohelp.WriteUint32(w, uint32(%ASGN))"
	}
	for _, st := range f.Structs {
		out[st.Name] = fmtEncode + fmtErrReturn
	}
	for _, msg := range f.Messages {
		out[msg.Name] = fmtEncode + fmtErrReturn
	}
	for _, union := range f.Unions {
		uout := union.typeMarshallers()
		for k, v := range uout {
			out[k] = v
		}
	}
	return out
}

func (u Union) typeMarshallers() map[string]string {
	out := make(map[string]string)
	out[u.Name] = fmtEncode + fmtErrReturn
	for _, ufd := range u.Fields {
		if ufd.Struct != nil {
			out[ufd.Struct.Name] = fmtEncode + fmtErrReturn
		}
		if ufd.Message != nil {
			out[ufd.Message.Name] = fmtEncode + fmtErrReturn
		}
	}
	return out
}

func (f File) typeLengthers() map[string]string {
	out := make(map[string]string)
	for typ, sz := range fixedSizeTypes {
		out[typ] = "bodyLen += " + strconv.Itoa(int(sz))
	}
	out[typeString] = fmtAdd4PlusLenToBodyLen
	for _, en := range f.Enums {
		out[en.Name] = fmtAdd4ToBodyLen
	}
	for _, st := range f.Structs {
		out[st.Name] = fmtAddSizeToBodyLen
	}
	for _, msg := range f.Messages {
		out[msg.Name] = fmtAddSizeToBodyLen
	}
	for _, union := range f.Unions {
		uout := union.typeLengthers()
		for k, v := range uout {
			out[k] = v
		}
	}
	return out
}

func (u Union) typeLengthers() map[string]string {
	out := make(map[string]string)
	out[u.Name] = fmtAddSizeToBodyLen
	for _, ufd := range u.Fields {
		if ufd.Struct != nil {
			out[ufd.Struct.Name] = fmtAddSizeToBodyLen
		}
		if ufd.Message != nil {
			out[ufd.Message.Name] = fmtAddSizeToBodyLen
		}
	}
	return out
}

func (f File) typeByters() map[string]string {
	out := make(map[string]string)
	for typ, sz := range fixedSizeTypes {
		out[typ] = "iohelp.Write" + fixedTitleString(typ) + "Bytes(buf[at:], %ASGN)\n" +
			"at += " + strconv.Itoa(int(sz))
	}
	out[typeString] = "iohelp.WriteUint32Bytes(buf[at:], uint32(len(%ASGN)))\n" +
		"copy(buf[at+4:at+4+len(%ASGN)], []byte(%ASGN))\n" + fmtAdd4PlusLenToAt

	out[typeDate] = "if %ASGN != (time.Time{}) {\n" +
		"\tiohelp.WriteInt64Bytes(buf[at:], ((%ASGN).UnixNano() / 100))\n" +
		"} else {\n" +
		"\tiohelp.WriteInt64Bytes(buf[at:], 0)\n" +
		"}\n" +
		"at += 8"
	for _, en := range f.Enums {
		out[en.Name] = "iohelp.WriteUint32Bytes(buf[at:], uint32(%ASGN))\n" + fmtAdd4ToAt
	}
	for _, st := range f.Structs {
		out[st.Name] = fmtMarshal + fmtAddSizeToAt
	}
	for _, msg := range f.Messages {
		out[msg.Name] = fmtMarshal + fmtAddSizeToAt
	}
	for _, union := range f.Unions {
		uout := union.typeByters()
		for k, v := range uout {
			out[k] = v
		}
	}
	return out
}

func (u Union) typeByters() map[string]string {
	out := map[string]string{}
	out[u.Name] = fmtMarshal + fmtAddSizeToAt
	for _, ufd := range u.Fields {
		if ufd.Struct != nil {
			out[ufd.Struct.Name] = fmtMarshal + fmtAddSizeToAt
		}
		if ufd.Message != nil {
			out[ufd.Message.Name] = fmtMarshal + fmtAddSizeToAt
		}
	}
	return out
}

func (f File) typeByteReaders(gs GenerateSettings) map[string]string {
	out := make(map[string]string)
	for typ, sz := range fixedSizeTypes {
		out[typ] = "%ASGN = iohelp.Read" + fixedTitleString(typ) + "Bytes(buf[at:])\n" +
			"at += " + strconv.Itoa(int(sz))
	}
	stringRead := "ReadStringBytes(buf[at:])"
	if gs.SharedMemoryStrings {
		stringRead = "ReadStringBytesSharedMemory(buf[at:])"
	}

	out[typeString] = "%ASGN = iohelp.Must" + stringRead + "\n" + fmtAdd4PlusLenToAt
	out["string&safe"] = "%ASGN, err = iohelp." + stringRead + "\n" + fmtErrReturn + "\n" + fmtAdd4PlusLenToAt

	for _, en := range f.Enums {
		out[en.Name] = "%ASGN = %TYPE(iohelp.ReadUint32Bytes(buf[at:]))\n" + fmtAdd4ToAt
	}
	for _, st := range f.Structs {
		out[st.Name] = mustMakeFormat(st.Namespace, gs) + fmtAddSizeToAt
		out[st.Name+"&safe"] = makeFormat(st.Namespace, gs) + fmtErrReturn + "\n" + fmtAddSizeToAt
	}
	for _, msg := range f.Messages {
		out[msg.Name] = mustMakeFormat(msg.Namespace, gs) + fmtAddSizeToAt
		out[msg.Name+"&safe"] = makeFormat(msg.Namespace, gs) + fmtErrReturn + "\n" + fmtAddSizeToAt
	}
	for _, union := range f.Unions {
		uout := union.typeByteReaders(gs)
		for k, v := range uout {
			out[k] = v
		}
	}
	return out
}

func (u Union) typeByteReaders(settings GenerateSettings) map[string]string {
	out := map[string]string{}
	out[u.Name] = mustMakeFormat(u.Namespace, settings) + fmtAddSizeToAt
	out[u.Name+"&safe"] = makeFormat(u.Namespace, settings) + fmtAddSizeToAt + fmtErrReturn + "\n" + fmtAddSizeToAt
	for _, ufd := range u.Fields {
		if ufd.Struct != nil {
			st := ufd.Struct
			out[st.Name] = mustMakeFormat(st.Namespace, settings) + fmtAddSizeToAt
			out[st.Name+"&safe"] = makeFormat(st.Namespace, settings) + fmtErrReturn + "\n" + fmtAddSizeToAt
		}
		if ufd.Message != nil {
			msg := ufd.Message
			out[msg.Name] = mustMakeFormat(msg.Namespace, settings) + fmtAddSizeToAt
			out[msg.Name+"&safe"] = makeFormat(msg.Namespace, settings) + fmtErrReturn + "\n" + fmtAddSizeToAt
		}
	}
	return out
}
