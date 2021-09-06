package bebop

import (
	"io"
	"strconv"
	"strings"
)

func (st Struct) generateTypeDefinition(w io.Writer, settings GenerateSettings) {
	writeOpCode(w, st.Name, st.OpCode)
	writeRecordAssertion(w, st.Name)
	writeComment(w, 0, st.Comment)
	writeGoStructDef(w, st.Name)
	for _, fd := range st.Fields {
		writeFieldDefinition(fd, w, st.ReadOnly, false, settings)
	}
	writeCloseBlock(w)
}

func (st Struct) generateMarshalBebopTo(w io.Writer, settings GenerateSettings) {
	exposedName := exposeName(st.Name)
	writeLine(w, "func (bbp %s) MarshalBebopTo(buf []byte) int {", exposedName)
	startAt := "0"
	if st.OpCode != 0 {
		writeLine(w, "\tiohelp.WriteUint32Bytes(buf, uint32(%sOpCode))", exposedName)
		startAt = "4"
	}
	if len(st.Fields) == 0 {
		writeLine(w, "\treturn "+startAt)
		writeCloseBlock(w)
		return
	}
	writeLine(w, "\tat := "+startAt)
	for _, fd := range st.Fields {
		name := exposeName(fd.Name)
		if st.ReadOnly {
			name = unexposeName(fd.Name)
		}
		writeFieldByter("bbp."+name, fd.FieldType, w, settings, 1)
	}
	writeLine(w, "\treturn at")
	writeCloseBlock(w)
}

func (st Struct) generateUnmarshalBebop(w io.Writer, settings GenerateSettings) {
	exposedName := exposeName(st.Name)
	writeLine(w, "func (bbp *%s) UnmarshalBebop(buf []byte) (err error) {", exposedName)
	if st.OpCode != 0 {
		if len(st.Fields) > 0 {
			writeLine(w, "\tat := 4")
		}
	} else {
		if len(st.Fields) > 0 {
			writeLine(w, "\tat := 0")
		}
	}
	for _, fd := range st.Fields {
		name := exposeName(fd.Name)
		if st.ReadOnly {
			name = unexposeName(fd.Name)
		}
		writeFieldReadByter("bbp."+name, fd.FieldType, w, settings, 1, true)
	}
	writeLine(w, "\treturn nil")
	writeCloseBlock(w)
}

func (st Struct) generateMustUnmarshalBebop(w io.Writer, settings GenerateSettings) {
	exposedName := exposeName(st.Name)
	writeLine(w, "func (bbp *%s) MustUnmarshalBebop(buf []byte) {", exposedName)
	if st.OpCode != 0 {
		if len(st.Fields) > 0 {
			writeLine(w, "\tat := 4")
		}
	} else {
		if len(st.Fields) > 0 {
			writeLine(w, "\tat := 0")
		}
	}
	for _, fd := range st.Fields {
		name := exposeName(fd.Name)
		if st.ReadOnly {
			name = unexposeName(fd.Name)
		}
		writeFieldReadByter("bbp."+name, fd.FieldType, w, settings, 1, false)
	}
	writeCloseBlock(w)
}

func (st Struct) generateEncodeBebop(w io.Writer, settings GenerateSettings) {
	exposedName := exposeName(st.Name)
	*settings.isFirstTopLength = true
	writeLine(w, "func (bbp %s) EncodeBebop(iow io.Writer) (err error) {", exposedName)
	if len(st.Fields) == 0 && st.OpCode == 0 {
		writeLine(w, "\treturn nil")
	} else {
		writeLine(w, "\tw := iohelp.NewErrorWriter(iow)")
		if st.OpCode != 0 {
			writeLine(w, "\tiohelp.WriteUint32(w, uint32(%sOpCode))", exposedName)
		}
		for _, fd := range st.Fields {
			name := exposeName(fd.Name)
			if st.ReadOnly {
				name = unexposeName(fd.Name)
			}
			writeFieldMarshaller("bbp."+name, fd.FieldType, w, settings, 1)
		}
		writeLine(w, "\treturn w.Err")
	}
	writeCloseBlock(w)
}

func (st Struct) generateDecodeBebop(w io.Writer, settings GenerateSettings) {
	exposedName := exposeName(st.Name)
	*settings.isFirstTopLength = true
	writeLine(w, "func (bbp *%s) DecodeBebop(ior io.Reader) (err error) {", exposedName)
	if len(st.Fields) == 0 && st.OpCode == 0 {
		writeLine(w, "\treturn nil")
	} else {
		writeLine(w, "\tr := iohelp.NewErrorReader(ior)")
		if st.OpCode != 0 {
			writeLine(w, "\tr.Read(make([]byte, 4))")
			writeLine(w, "\tif r.Err != nil {\n\t\treturn r.Err\n\t}")
		}
		for _, fd := range st.Fields {
			name := exposeName(fd.Name)
			if st.ReadOnly {
				name = unexposeName(fd.Name)
			}
			writeStructFieldUnmarshaller("&bbp."+name, fd.FieldType, w, settings, 1)
		}
		writeLine(w, "\treturn r.Err")
	}
	writeCloseBlock(w)
}

func (st Struct) generateSize(w io.Writer, settings GenerateSettings) {
	exposedName := exposeName(st.Name)
	writeLine(w, "func (bbp %s) Size() int {", exposedName)
	if len(st.Fields) == 0 && st.OpCode == 0 {
		writeLine(w, "\treturn 0")
	} else {
		writeLine(w, "\tbodyLen := 0")
		if st.OpCode != 0 {
			writeLine(w, "\tbodyLen += 4")
		}
		for _, fd := range st.Fields {
			name := exposeName(fd.Name)
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

func (st Struct) generateReadOnlyGetters(w io.Writer, settings GenerateSettings) {
	// TODO: slices are not read only, we need to return a copy.
	exposedName := exposeName(st.Name)
	for _, fd := range st.Fields {
		writeLine(w, "func (bbp %s) Get%s() %s {", exposedName, exposeName(fd.Name), fd.FieldType.goString(settings))
		writeLine(w, "\treturn bbp.%s", unexposeName(fd.Name))
		writeCloseBlock(w)
	}
	writeLine(w, "func New%s(", exposedName)
	for _, fd := range st.Fields {
		writeLine(w, "\t\t%s %s,", unexposeName(fd.Name), fd.FieldType.goString(settings))
	}
	writeLine(w, "\t) %s {", exposedName)
	writeLine(w, "\treturn %s{", exposedName)
	for _, fd := range st.Fields {
		writeLine(w, "\t\t%s: %s,", unexposeName(fd.Name), unexposeName(fd.Name))
	}
	writeLine(w, "\t}")
	writeCloseBlock(w)
}

// Generate writes a .go struct definition out to w.
func (st Struct) Generate(w io.Writer, settings GenerateSettings) {
	st.generateTypeDefinition(w, settings)
	st.generateMarshalBebopTo(w, settings)
	st.generateUnmarshalBebop(w, settings)
	if settings.GenerateUnsafeMethods {
		st.generateMustUnmarshalBebop(w, settings)
	}
	st.generateEncodeBebop(w, settings)
	st.generateDecodeBebop(w, settings)
	st.generateSize(w, settings)

	isEmpty := len(st.Fields) == 0 && st.OpCode == 0
	writeWrappers(w, st.Name, isEmpty, settings)
	if st.ReadOnly {
		st.generateReadOnlyGetters(w, settings)
	}
}

func writeStructFieldUnmarshaller(name string, typ FieldType, w io.Writer, settings GenerateSettings, depth int) {
	iName := "i" + strconv.Itoa(depth)
	if typ.Array != nil {
		writeLineWithTabs(w, "%RECV = make([]%TYPE, iohelp.ReadUint32(r))", depth, name, typ.Array.goString(settings))
		writeLineWithTabs(w, "for "+iName+" := range %RECV {", depth, name)
		name = "&(" + name[1:] + "[" + iName + "])"
		writeStructFieldUnmarshaller(name, *typ.Array, w, settings, depth+1)
		writeLineWithTabs(w, "}", depth)
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
		w.Write([]byte(strings.Replace(ln, "=", ":=", 1)))
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
