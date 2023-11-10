package bebop

import (
	"io"
	"sort"
	"strconv"
)

func (u Union) generateMarshalBebopTo(w io.Writer, settings GenerateSettings, fields []fieldWithNumber) {
	exposedName := exposeName(u.Name, settings)
	if settings.AlwaysUsePointerReceivers {
		writeLine(w, "func (bbp *%s) MarshalBebopTo(buf []byte) int {", exposedName)
	} else {
		writeLine(w, "func (bbp %s) MarshalBebopTo(buf []byte) int {", exposedName)
	}
	writeLine(w, "\tat := 0")
	// 5 = 4 bytes of size + 1 byte discriminator
	writeLine(w, "\tiohelp.WriteUint32Bytes(buf[at:], uint32(bbp.Size()-5))")
	writeLine(w, "\tat += 4")
	for _, fd := range fields {
		name := exposeName(fd.Name, settings)
		num := strconv.Itoa(int(fd.num))
		name = "*bbp." + name
		writeLineWithTabs(w, "if %RECV != nil {", 1, name)
		writeLineWithTabs(w, "buf[at] = %ASGN", 2, num)
		writeLineWithTabs(w, "at++", 2)
		writeFieldByter(name, fd.FieldType, w, settings, 2)
		writeLineWithTabs(w, "return at", 2)
		writeLineWithTabs(w, "}", 1)
	}
	writeLine(w, "\treturn at")
	writeCloseBlock(w)
}

func (u Union) generateUnmarshalBebop(w io.Writer, settings GenerateSettings, fields []fieldWithNumber) {
	exposedName := exposeName(u.Name, settings)
	writeLine(w, "func (bbp *%s) UnmarshalBebop(buf []byte) (err error) {", exposedName)
	writeLine(w, "\tat := 0")
	writeLine(w, "\t_ = iohelp.ReadUint32Bytes(buf[at:])")
	writeLine(w, "\tbuf = buf[4:]")
	writeLine(w, "\tif len(buf) == 0 {")
	writeLine(w, "\t\treturn iohelp.ErrUnpopulatedUnion")
	writeLine(w, "\t}")
	writeLine(w, "\tfor {")
	writeLine(w, "\t\tswitch buf[at] {")
	for _, fd := range fields {
		name := exposeName(fd.Name, settings)
		writeLine(w, "\t\tcase %d:", fd.num)
		writeLine(w, "\t\t\tat += 1")
		writeLine(w, "\t\t\tbbp.%[1]s = new(%[2]s)", name, fd.FieldType.goString(settings))
		writeFieldReadByter("(*bbp."+name+")", fd.FieldType, w, settings, 3, true)
		writeLine(w, "\t\t\treturn nil")
	}
	writeLine(w, "\t\tdefault:")
	writeLine(w, "\t\t\treturn nil")
	writeLine(w, "\t\t}")
	writeLine(w, "\t}")
	writeCloseBlock(w)
}

func (u Union) generateMustUnmarshalBebop(w io.Writer, settings GenerateSettings, fields []fieldWithNumber) {
	exposedName := exposeName(u.Name, settings)
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
		writeLine(w, "\t\t\treturn")
	}
	writeLine(w, "\t\tdefault:")
	writeLine(w, "\t\t\treturn")
	writeLine(w, "\t\t}")
	writeLine(w, "\t}")
	writeCloseBlock(w)
}

func (u Union) generateEncodeBebop(w io.Writer, settings GenerateSettings, fields []fieldWithNumber) {
	exposedName := exposeName(u.Name, settings)
	*settings.isFirstTopLength = true
	if settings.AlwaysUsePointerReceivers {
		writeLine(w, "func (bbp *%s) EncodeBebop(iow io.Writer) (err error) {", exposedName)
	} else {
		writeLine(w, "func (bbp %s) EncodeBebop(iow io.Writer) (err error) {", exposedName)
	}
	writeLine(w, "\tw := iohelp.NewErrorWriter(iow)")
	writeLine(w, "\tiohelp.WriteUint32(w, uint32(bbp.Size()-5))")
	for _, fd := range fields {
		name := exposeName(fd.Name, settings)
		num := strconv.Itoa(int(fd.num))
		name = "*bbp." + name
		writeLineWithTabs(w, "if %RECV != nil {", 1, name)
		writeLineWithTabs(w, "w.Write([]byte{%ASGN})", 2, num)
		writeFieldMarshaller(name, fd.FieldType, w, settings, 2)
		writeLineWithTabs(w, "return w.Err", 2)
		writeLineWithTabs(w, "}", 1)
	}
	writeLine(w, "\treturn w.Err")
	writeCloseBlock(w)
}

func (u Union) generateDecodeBebop(w io.Writer, settings GenerateSettings, fields []fieldWithNumber) {
	exposedName := exposeName(u.Name, settings)
	*settings.isFirstTopLength = true
	writeLine(w, "func (bbp *%s) DecodeBebop(ior io.Reader) (err error) {", exposedName)
	writeLine(w, "\tr := iohelp.NewErrorReader(ior)")
	writeLine(w, "\tbodyLen := iohelp.ReadUint32(r)")
	writeLine(w, "\tr.Reader = &io.LimitedReader{R: r.Reader, N: int64(bodyLen) + 1}")
	writeLine(w, "\tfor {")
	writeLine(w, "\t\tswitch iohelp.ReadByte(r) {")
	for _, fd := range fields {
		writeLine(w, "\t\tcase %d:", fd.num)
		name := exposeName(fd.Name, settings)
		writeLine(w, "\t\t\tbbp.%[1]s = new(%[2]s)", name, fd.FieldType.goString(settings))
		writeMessageFieldUnmarshaller("bbp."+name, fd.FieldType, w, settings, 3)
		writeLine(w, "\t\t\tio.ReadAll(r)")
		writeLine(w, "\t\t\treturn r.Err")
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

func (u Union) generateSize(w io.Writer, settings GenerateSettings, fields []fieldWithNumber) {
	exposedName := exposeName(u.Name, settings)
	if settings.AlwaysUsePointerReceivers {
		writeLine(w, "func (bbp *%s) Size() int {", exposedName)
	} else {
		writeLine(w, "func (bbp %s) Size() int {", exposedName)
	}
	// size at front (4)
	writeLine(w, "\tbodyLen := 4")
	for _, fd := range fields {
		name := exposeName(fd.Name, settings)
		name = "*bbp." + name
		writeLineWithTabs(w, "if %RECV != nil {", 1, name)
		writeLineWithTabs(w, "bodyLen += 1", 2)
		writeFieldBodyCount(name, fd.FieldType, w, settings, 2)
		writeLineWithTabs(w, "return bodyLen", 2)
		writeLineWithTabs(w, "}", 1)
	}
	writeLine(w, "\treturn bodyLen")
	writeCloseBlock(w)
}

// Generate writes a .go union definition out to w.
func (u Union) Generate(w io.Writer, settings GenerateSettings) {
	fields := make([]fieldWithNumber, 0, len(u.Fields))
	for i, ufd := range u.Fields {
		var fd Field
		if ufd.Struct != nil {
			fd.FieldType.Simple = ufd.Struct.Name
		}
		if ufd.Message != nil {
			fd.FieldType.Simple = ufd.Message.Name
		}
		fd.Name = fd.FieldType.Simple
		fd.Tags = ufd.Tags
		fields = append(fields, fieldWithNumber{
			UnionField: ufd,
			Field:      fd,
			num:        i,
		})
	}
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].num < fields[j].num
	})
	for _, field := range fields {
		if field.UnionField.Struct != nil {
			field.UnionField.Struct.Generate(w, settings)
		}
		if field.UnionField.Message != nil {
			field.UnionField.Message.Generate(w, settings)
		}
	}
	writeRecordTypeDefinition(w, u.Name, u.OpCode, u.Comment, settings, fields)
	u.generateMarshalBebopTo(w, settings, fields)
	u.generateUnmarshalBebop(w, settings, fields)
	if settings.GenerateUnsafeMethods {
		u.generateMustUnmarshalBebop(w, settings, fields)
	}
	u.generateEncodeBebop(w, settings, fields)
	u.generateDecodeBebop(w, settings, fields)
	u.generateSize(w, settings, fields)
	isEmpty := len(u.Fields) == 0
	writeWrappers(w, u.Name, isEmpty, settings)
}
