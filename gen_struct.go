package bebop

import (
	"io"
	"strconv"
	"strings"

	"github.com/200sc/bebop/iohelp"
)

func (st Struct) generateTypeDefinition(w *iohelp.ErrorWriter, settings GenerateSettings) {
	writeOpCode(w, st.Name, st.OpCode, settings)
	writeRecordAssertion(w, st.Name, settings)
	writeComment(w, 0, st.Comment, settings)
	writeGoStructDef(w, st.Name, settings)
	for _, fd := range st.Fields {
		writeFieldDefinition(fd, w, st.ReadOnly, false, settings)
	}
	writeCloseBlock(w)
}

func (st Struct) generateMarshalBebopTo(w *iohelp.ErrorWriter, settings GenerateSettings) {
	exposedName := exposeName(st.Name, settings)
	if settings.AlwaysUsePointerReceivers {
		writeLine(w, "func (bbp *%s) MarshalBebopTo(buf []byte) int {", exposedName)
	} else {
		writeLine(w, "func (bbp %s) MarshalBebopTo(buf []byte) int {", exposedName)
	}
	startAt := "0"
	if len(st.Fields) == 0 {
		writeLine(w, "\treturn "+startAt)
		writeCloseBlock(w)
		return
	}
	writeLine(w, "\tat := "+startAt)
	for _, fd := range st.Fields {
		name := exposeName(fd.Name, settings)
		if st.ReadOnly {
			name = unexposeName(fd.Name)
		}
		writeFieldByter("bbp."+name, fd.FieldType, w, settings, 1)
	}
	writeLine(w, "\treturn at")
	writeCloseBlock(w)
}

func (st Struct) generateUnmarshalBebop(w *iohelp.ErrorWriter, settings GenerateSettings) {
	exposedName := exposeName(st.Name, settings)
	writeLine(w, "func (bbp *%s) UnmarshalBebop(buf []byte) (err error) {", exposedName)
	if len(st.Fields) > 0 {
		writeLine(w, "\tat := 0")
	}
	for _, fd := range st.Fields {
		name := exposeName(fd.Name, settings)
		if st.ReadOnly {
			name = unexposeName(fd.Name)
		}
		writeFieldReadByter("bbp."+name, fd.FieldType, w, settings, 1, true)
	}
	writeLine(w, "\treturn nil")
	writeCloseBlock(w)
}

func (st Struct) generateMustUnmarshalBebop(w *iohelp.ErrorWriter, settings GenerateSettings) {
	exposedName := exposeName(st.Name, settings)
	writeLine(w, "func (bbp *%s) MustUnmarshalBebop(buf []byte) {", exposedName)
	if len(st.Fields) > 0 {
		writeLine(w, "\tat := 0")
	}
	for _, fd := range st.Fields {
		name := exposeName(fd.Name, settings)
		if st.ReadOnly {
			name = unexposeName(fd.Name)
		}
		writeFieldReadByter("bbp."+name, fd.FieldType, w, settings, 1, false)
	}
	writeCloseBlock(w)
}

func (st Struct) generateEncodeBebop(w *iohelp.ErrorWriter, settings GenerateSettings) {
	exposedName := exposeName(st.Name, settings)
	*settings.isFirstTopLength = true
	if settings.AlwaysUsePointerReceivers {
		writeLine(w, "func (bbp *%s) EncodeBebop(iow io.Writer) (err error) {", exposedName)
	} else {
		writeLine(w, "func (bbp %s) EncodeBebop(iow io.Writer) (err error) {", exposedName)
	}
	if len(st.Fields) == 0 {
		writeLine(w, "\treturn nil")
	} else {
		writeLine(w, "\tw := iohelp.NewErrorWriter(iow)")
		for _, fd := range st.Fields {
			name := exposeName(fd.Name, settings)
			if st.ReadOnly {
				name = unexposeName(fd.Name)
			}
			writeFieldMarshaller("bbp."+name, fd.FieldType, w, settings, 1)
		}
		writeLine(w, "\treturn w.Err")
	}
	writeCloseBlock(w)
}

func (st Struct) generateDecodeBebop(w *iohelp.ErrorWriter, settings GenerateSettings) {
	exposedName := exposeName(st.Name, settings)
	*settings.isFirstTopLength = true
	writeLine(w, "func (bbp *%s) DecodeBebop(ior io.Reader) (err error) {", exposedName)
	if len(st.Fields) == 0 {
		writeLine(w, "\treturn nil")
	} else {
		writeLine(w, "\tr := iohelp.NewErrorReader(ior)")
		for _, fd := range st.Fields {
			name := exposeName(fd.Name, settings)
			if st.ReadOnly {
				name = unexposeName(fd.Name)
			}
			writeStructFieldUnmarshaller("&bbp."+name, fd.FieldType, w, settings, 1)
		}
		writeLine(w, "\treturn r.Err")
	}
	writeCloseBlock(w)
}

func (st Struct) generateSize(w *iohelp.ErrorWriter, settings GenerateSettings) {
	exposedName := exposeName(st.Name, settings)
	if settings.AlwaysUsePointerReceivers {
		writeLine(w, "func (bbp *%s) Size() int {", exposedName)
	} else {
		writeLine(w, "func (bbp %s) Size() int {", exposedName)
	}
	if len(st.Fields) == 0 {
		writeLine(w, "\treturn 0")
	} else {
		writeLine(w, "\tbodyLen := 0")
		for _, fd := range st.Fields {
			name := exposeName(fd.Name, settings)
			if st.ReadOnly {
				name = unexposeName(fd.Name)
			}
			name = "bbp." + name
			writeFieldBodyCount(name, fd.FieldType, w, settings, 1)
		}
		writeLine(w, "\treturn bodyLen")
	}
	writeCloseBlock(w)
}

func (st Struct) generateReadOnlyGetters(w *iohelp.ErrorWriter, settings GenerateSettings) {
	// TODO: slices are not read only, we need to return a copy.
	exposedName := exposeName(st.Name, settings)
	for _, fd := range st.Fields {
		if settings.AlwaysUsePointerReceivers {
			writeLine(w, "func (bbp *%s) Get%s() %s {", exposedName, exposeName(fd.Name, settings), fd.FieldType.goString(settings))
		} else {
			writeLine(w, "func (bbp %s) Get%s() %s {", exposedName, exposeName(fd.Name, settings), fd.FieldType.goString(settings))
		}
		writeLine(w, "\treturn bbp.%s", unexposeName(fd.Name))
		writeCloseBlock(w)
	}
	newFmt := exposeName("New", settings)
	writeLine(w, "func %s%s(", newFmt, exposedName)
	for _, fd := range st.Fields {
		writeLine(w, "\t\t%s %s,", unexposeName(fd.Name), fd.FieldType.goString(settings))
	}
	writeLine(w, ") %s {", exposedName)
	writeLine(w, "\treturn %s{", exposedName)
	for _, fd := range st.Fields {
		writeLine(w, "\t\t%s: %s,", unexposeName(fd.Name), unexposeName(fd.Name))
	}
	writeLine(w, "\t}")
	writeCloseBlock(w)
}

// Generate writes a .go struct definition out to w.
func (st Struct) Generate(w io.Writer, settings GenerateSettings) {
	ew := iohelp.NewErrorWriter(w)
	st.generateTypeDefinition(ew, settings)
	st.generateMarshalBebopTo(ew, settings)
	st.generateUnmarshalBebop(ew, settings)
	if settings.GenerateUnsafeMethods {
		st.generateMustUnmarshalBebop(ew, settings)
	}
	st.generateEncodeBebop(ew, settings)
	st.generateDecodeBebop(ew, settings)
	st.generateSize(ew, settings)

	isEmpty := len(st.Fields) == 0
	writeWrappers(ew, st.Name, isEmpty, settings)
	if st.ReadOnly {
		st.generateReadOnlyGetters(ew, settings)
	}
}

func writeStructFieldUnmarshaller(name string, typ FieldType, w *iohelp.ErrorWriter, settings GenerateSettings, depth int) {
	iName := "i" + strconv.Itoa(depth)
	if typ.Array != nil {
		writeLineWithTabs(w, "%RECV = make([]%TYPE, iohelp.ReadUint32(r))", depth, name, typ.Array.goString(settings))
		if typ.Array.Simple == typeByte {
			writeLineWithTabs(w, "r.Read(%RECV)", depth, name)
		} else {
			writeLineWithTabs(w, "for "+iName+" := range %RECV {", depth, name)
			name = "&(" + name[1:] + "[" + iName + "])"
			writeStructFieldUnmarshaller(name, *typ.Array, w, settings, depth+1)
			writeLineWithTabs(w, "}", depth)
		}
	} else if typ.Map != nil {
		lnName := "ln" + strconv.Itoa(depth)
		if *settings.isFirstTopLength && depth == 1 {
			writeLineWithTabs(w, lnName+" := iohelp.ReadUint32(r)", depth)
			*settings.isFirstTopLength = false
		} else if depth == 1 {
			writeLineWithTabs(w, lnName+" = iohelp.ReadUint32(r)", depth)
		} else {
			writeLineWithTabs(w, lnName+" := iohelp.ReadUint32(r)", depth)
		}
		writeLineWithTabs(w, "%RECV = make(%TYPE, "+lnName+")", depth, name, typ.Map.goString(settings))
		writeLineWithTabs(w, "for "+iName+" := uint32(0); "+iName+" < "+lnName+"; "+iName+"++ {", depth, name)
		ln := getLineWithTabs(settings.typeUnmarshallers[typ.Map.Key], depth+1, "&"+depthName("k", depth))
		w.SafeWrite([]byte(strings.Replace(ln, "=", ":=", 1)))
		name = "&(" + name[1:] + "[" + depthName("k", depth) + "])"
		writeStructFieldUnmarshaller(name, typ.Map.Value, w, settings, depth+1)
		writeLineWithTabs(w, "}", depth)
	} else {
		simpleTyp := typ.Simple
		if alias, ok := settings.importTypeAliases[simpleTyp]; ok {
			simpleTyp = alias
		}
		writeLineWithTabs(w, settings.typeUnmarshallers[simpleTyp], depth, name, typ.goString(settings))
	}
}
