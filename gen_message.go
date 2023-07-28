package bebop

import (
	"io"
	"sort"
	"strconv"
	"strings"
)

type fieldWithNumber struct {
	UnionField UnionField
	Field
	num uint8
}

func (msg Message) generateMarshalBebopTo(w io.Writer, settings GenerateSettings, fields []fieldWithNumber) {
	exposedName := exposeName(msg.Name, settings)
	writeLine(w, "func (bbp %s) MarshalBebopTo(buf []byte) int {", exposedName)
	writeLine(w, "\tat := 0")
	writeLine(w, "\tiohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-4))")
	writeLine(w, "\tat += 4")
	for _, fd := range fields {
		if fd.Deprecated {
			continue
		}
		name := exposeName(fd.Name, settings)
		num := strconv.Itoa(int(fd.num))
		name = "*bbp." + name
		writeLineWithTabs(w, "if %RECV != nil {", 1, name)
		writeLineWithTabs(w, "buf[at] = %ASGN", 2, num)
		writeLineWithTabs(w, "at++", 2)
		writeFieldByter(name, fd.FieldType, w, settings, 2)
		writeLineWithTabs(w, "}", 1)
	}
	writeLine(w, "\treturn at")
	writeCloseBlock(w)
}

func (msg Message) generateUnmarshalBebop(w io.Writer, settings GenerateSettings, fields []fieldWithNumber) {
	exposedName := exposeName(msg.Name, settings)
	writeLine(w, "func (bbp *%s) UnmarshalBebop(buf []byte) (err error) {", exposedName)
	writeLine(w, "\tat := 0")
	writeLine(w, "\t_ = iohelp.ReadUint32Bytes(buf[at:])")
	writeLine(w, "\tbuf = buf[4:]")
	writeLine(w, "\tfor {")
	writeLine(w, "\t\tswitch buf[at] {")
	for _, fd := range fields {
		name := exposeName(fd.Name, settings)
		writeLine(w, "\t\tcase %d:", fd.num)
		writeLine(w, "\t\t\tat += 1")
		writeLine(w, "\t\t\tbbp.%[1]s = new(%[2]s)", name, fd.FieldType.goString(settings))
		writeFieldReadByter("(*bbp."+name+")", fd.FieldType, w, settings, 3, true)
	}
	writeLine(w, "\t\tdefault:")
	writeLine(w, "\t\t\treturn nil")
	writeLine(w, "\t\t}")
	writeLine(w, "\t}")
	writeCloseBlock(w)
}

func (msg Message) generateMustUnmarshalBebop(w io.Writer, settings GenerateSettings, fields []fieldWithNumber) {
	exposedName := exposeName(msg.Name, settings)
	writeLine(w, "func (bbp *%s) MustUnmarshalBebop(buf []byte) {", exposedName)
	writeLine(w, "\tat := 0")
	writeLine(w, "\t_ = iohelp.ReadUint32Bytes(buf[at:])")
	writeLine(w, "\tbuf = buf[4:]")
	writeLine(w, "\tfor {")
	writeLine(w, "\t\tswitch buf[at] {")
	for _, fd := range fields {
		name := exposeName(fd.Name, settings)
		writeLine(w, "\t\tcase %d:", fd.num)
		writeLine(w, "\t\t\tat += 1")
		writeLine(w, "\t\t\tbbp.%[1]s = new(%[2]s)", name, fd.FieldType.goString(settings))
		writeFieldReadByter("(*bbp."+name+")", fd.FieldType, w, settings, 3, false)
	}
	writeLine(w, "\t\tdefault:")
	writeLine(w, "\t\t\treturn")
	writeLine(w, "\t\t}")
	writeLine(w, "\t}")
	writeCloseBlock(w)
}

func (msg Message) generateEncodeBebop(w io.Writer, settings GenerateSettings, fields []fieldWithNumber) {
	exposedName := exposeName(msg.Name, settings)
	*settings.isFirstTopLength = true
	writeLine(w, "func (bbp %s) EncodeBebop(iow io.Writer) (err error) {", exposedName)
	writeLine(w, "\tw := iohelp.NewErrorWriter(iow)")
	writeLine(w, "\tiohelp.WriteUint32(w, uint32(bbp.Size()-4))")
	for _, fd := range fields {
		if fd.Deprecated {
			continue
		}
		name := exposeName(fd.Name, settings)
		num := strconv.Itoa(int(fd.num))
		name = "*bbp." + name
		writeLineWithTabs(w, "if %RECV != nil {", 1, name)
		writeLineWithTabs(w, "w.Write([]byte{%ASGN})", 2, num)
		writeFieldMarshaller(name, fd.FieldType, w, settings, 2)
		writeLineWithTabs(w, "}", 1)
	}
	writeLine(w, "\tw.Write([]byte{0})")
	writeLine(w, "\treturn w.Err")
	writeCloseBlock(w)
}

func (msg Message) generateDecodeBebop(w io.Writer, settings GenerateSettings, fields []fieldWithNumber) {
	exposedName := exposeName(msg.Name, settings)
	*settings.isFirstTopLength = true
	writeLine(w, "func (bbp *%s) DecodeBebop(ior io.Reader) (err error) {", exposedName)
	writeLine(w, "\tr := iohelp.NewErrorReader(ior)")
	writeLine(w, "\tbodyLen := iohelp.ReadUint32(r)")
	writeLine(w, "\tr.Reader = &io.LimitedReader{R:r.Reader, N:int64(bodyLen)}")
	writeLine(w, "\tfor {")
	writeLine(w, "\t\tswitch iohelp.ReadByte(r) {")
	for _, fd := range fields {
		writeLine(w, "\t\tcase %d:", fd.num)
		name := exposeName(fd.Name, settings)
		writeLine(w, "\t\t\tbbp.%[1]s = new(%[2]s)", name, fd.FieldType.goString(settings))
		writeMessageFieldUnmarshaller("bbp."+name, fd.FieldType, w, settings, 3)
	}
	// ref: https://github.com/RainwayApp/bebop/wiki/Wire-format#messages, final paragraph
	// we're allowed to skip parsing all remaining fields if we see one that we don't know about.
	writeLine(w, "\t\tdefault:")
	writeLine(w, "\t\t\tio.ReadAll(r)")
	writeLine(w, "\t\t\treturn r.Err")
	writeLine(w, "\t\t}")
	writeLine(w, "\t}")
	writeCloseBlock(w)
}

func (msg Message) generateSize(w io.Writer, settings GenerateSettings, fields []fieldWithNumber) {
	exposedName := exposeName(msg.Name, settings)
	writeLine(w, "func (bbp %s) Size() int {", exposedName)
	// size at front (4) + 0 byte (1)
	// q: why do messages end in a 0 byte?
	// a: (I think) because then we can loop reading a single byte for each field, and if we read 0
	// we know we're done and don't have to unread the byte
	writeLine(w, "\tbodyLen := 5")
	for _, fd := range fields {
		if fd.Deprecated {
			continue
		}
		name := exposeName(fd.Name, settings)
		name = "*bbp." + name
		writeLineWithTabs(w, "if %RECV != nil {", 1, name)
		writeLineWithTabs(w, "bodyLen += 1", 2)
		writeFieldBodyCount(name, fd.FieldType, w, settings, 2)
		writeLineWithTabs(w, "}", 1)
	}
	writeLine(w, "\treturn bodyLen")
	writeCloseBlock(w)
}

// Generate writes a .go message definition out to w.
func (msg Message) Generate(w io.Writer, settings GenerateSettings) {
	fields := make([]fieldWithNumber, 0, len(msg.Fields))
	for i, fd := range msg.Fields {
		fields = append(fields, fieldWithNumber{
			Field: fd,
			num:   i,
		})
	}
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].num < fields[j].num
	})
	writeRecordTypeDefinition(w, msg.Name, msg.OpCode, msg.Comment, settings, fields)
	msg.generateMarshalBebopTo(w, settings, fields)
	msg.generateUnmarshalBebop(w, settings, fields)
	if settings.GenerateUnsafeMethods {
		msg.generateMustUnmarshalBebop(w, settings, fields)
	}
	msg.generateEncodeBebop(w, settings, fields)
	msg.generateDecodeBebop(w, settings, fields)
	msg.generateSize(w, settings, fields)
	isEmpty := len(msg.Fields) == 0
	writeWrappers(w, msg.Name, isEmpty, settings)
}

func writeMessageFieldUnmarshaller(name string, typ FieldType, w io.Writer, settings GenerateSettings, depth int) {
	if typ.Array != nil {
		writeLineWithTabs(w, "%RECV = make([]%TYPE, iohelp.ReadUint32(r))", depth, name, typ.Array.goString(settings))
		if typ.Array.Simple == typeByte {
			writeLineWithTabs(w, "r.Read(%RECV)", depth, name)
		} else {
			writeLineWithTabs(w, "for i := range %RECV {", depth, name)
			writeMessageFieldUnmarshaller("("+name+")[i]", *typ.Array, w, settings, depth+1)
			writeLineWithTabs(w, "}", depth)
		}
	} else if typ.Map != nil {
		lnName := depthName("ln", depth)
		writeLineWithTabs(w, lnName+" := iohelp.ReadUint32(r)", depth)
		writeLineWithTabs(w, "%RECV = make("+typ.Map.goString(settings)+")", depth, name)
		writeLineWithTabs(w, "for i := uint32(0); i < "+lnName+"; i++ {", depth, name)
		ln := getLineWithTabs(settings.typeUnmarshallers[typ.Map.Key], depth+1, "&"+depthName("k", depth))
		w.Write([]byte(strings.Replace(ln, "=", ":=", 1)))
		writeMessageFieldUnmarshaller("("+name+")["+depthName("k", depth)+"]", typ.Map.Value, w, settings, depth+1)
		writeLineWithTabs(w, "}", depth)
	} else {
		simpleTyp := typ.Simple
		if alias, ok := settings.importTypeAliases[simpleTyp]; ok {
			simpleTyp = alias
		}
		writeLineWithTabs(w, settings.typeUnmarshallers[simpleTyp], depth, name, typ.goString(settings))
	}
}
