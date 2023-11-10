package bebop

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/200sc/bebop/internal/importgraph"
)

// GenerateSettings holds customization options for what Generate should do.
type GenerateSettings struct {
	// PackageName is optional if the target bebop file defines a go_package
	// constant. If both are provided, PackageName will take precedence.
	PackageName string

	typeByters        map[string]string
	typeByteReaders   map[string]string
	typeMarshallers   map[string]string
	typeUnmarshallers map[string]string
	typeLengthers     map[string]string
	customRecordTypes map[string]struct{}

	ImportGenerationMode
	imported          []File
	importTypeAliases map[string]string

	nextLength       *int
	isFirstTopLength *bool

	GenerateUnsafeMethods     bool
	SharedMemoryStrings       bool
	GenerateFieldTags         bool
	PrivateDefinitions        bool
	AlwaysUsePointerReceivers bool
}

type ImportGenerationMode uint8

const (
	// ImportGenerationModeSeparate will generate separate go files for
	// every bebop file, and will assume that imported files are
	// already generated. If imported file types are used and their containing
	// files do not contain a go_package constant, this mode will fail.
	ImportGenerationModeSeparate ImportGenerationMode = iota

	// ImportGenerationModeCombined will generate one go file including
	// all definitions from all imports. This does not require go_package
	// is defined anywhere, and maintains compatibility with files compilable
	// by the original bebopc compiler.
	ImportGenerationModeCombined ImportGenerationMode = iota
)

var allImportModes = map[ImportGenerationMode]struct{}{
	ImportGenerationModeSeparate: {},
	ImportGenerationModeCombined: {},
}

func (gs GenerateSettings) Validate() error {
	if _, ok := allImportModes[gs.ImportGenerationMode]; !ok {
		return fmt.Errorf("unknown import mode: %d", gs.ImportGenerationMode)
	}
	return nil
}

// Validate verifies a File can be successfully generated.
func (f File) Validate() error {
	allOpCodes := map[uint32]string{}
	allConsts := map[string]struct{}{}
	for _, c := range f.Consts {
		if _, ok := allConsts[c.Name]; ok {
			return fmt.Errorf("const has duplicated name %s", c.Name)
		}
		allConsts[c.Name] = struct{}{}
	}
	customTypes := map[string]struct{}{}
	structTypeUsage := map[string]map[string]bool{}
	for _, en := range f.Enums {
		if _, ok := primitiveTypes[en.Name]; ok {
			return fmt.Errorf("enum shares primitive type name %s", en.Name)
		}
		if _, ok := customTypes[en.Name]; ok {
			return fmt.Errorf("enum has duplicated name %s", en.Name)
		}
		if en.Namespace != "" {
			nameWithoutPrefix := strings.TrimPrefix(en.Name, en.Namespace+".")
			customTypes[nameWithoutPrefix] = struct{}{}
		}
		customTypes[en.Name] = struct{}{}
		optionNames := map[string]struct{}{}
		optionValues := map[int64]struct{}{}
		unsignedOptionValues := map[uint64]struct{}{}
		for _, opt := range en.Options {
			if _, ok := optionNames[opt.Name]; ok {
				return fmt.Errorf("enum %s has duplicate option name %s", en.Name, opt.Name)
			}
			optionNames[opt.Name] = struct{}{}
			if en.Unsigned {
				if _, ok := unsignedOptionValues[opt.UintValue]; ok {
					return fmt.Errorf("enum %s has duplicate option value %d", en.Name, opt.UintValue)
				}
				unsignedOptionValues[opt.UintValue] = struct{}{}
			} else {
				if _, ok := optionValues[opt.Value]; ok {
					return fmt.Errorf("enum %s has duplicate option value %d", en.Name, opt.Value)
				}
				optionValues[opt.Value] = struct{}{}
			}
		}
	}
	for _, st := range f.Structs {
		if _, ok := primitiveTypes[st.Name]; ok {
			return fmt.Errorf("struct shares primitive type name %s", st.Name)
		}
		if _, ok := customTypes[st.Name]; ok {
			return fmt.Errorf("struct has duplicated name %s", st.Name)
		}
		if st.Namespace != "" {
			nameWithoutPrefix := strings.TrimPrefix(st.Name, st.Namespace+".")
			customTypes[nameWithoutPrefix] = struct{}{}
			structTypeUsage[nameWithoutPrefix] = st.usedTypes()
		}
		customTypes[st.Name] = struct{}{}
		structTypeUsage[st.Name] = st.usedTypes()
		stNames := map[string]struct{}{}
		for _, fd := range st.Fields {
			if _, ok := stNames[fd.Name]; ok {
				return fmt.Errorf("struct %s has duplicate field name %s", st.Name, fd.Name)
			}
			stNames[fd.Name] = struct{}{}
		}
		if st.OpCode != 0 {
			if conflict, ok := allOpCodes[st.OpCode]; ok {
				return fmt.Errorf("struct %s has duplicate opcode %02x (duplicated in %s)", st.Name, st.OpCode, conflict)
			}
			allOpCodes[st.OpCode] = st.Name
		}
	}
	for _, msg := range f.Messages {
		if _, ok := primitiveTypes[msg.Name]; ok {
			return fmt.Errorf("message shares primitive type name %s", msg.Name)
		}
		if _, ok := customTypes[msg.Name]; ok {
			return fmt.Errorf("message has duplicated name %s", msg.Name)
		}
		if msg.Namespace != "" {
			nameWithoutPrefix := strings.TrimPrefix(msg.Name, msg.Namespace+".")
			customTypes[nameWithoutPrefix] = struct{}{}
		}
		customTypes[msg.Name] = struct{}{}
		msgNames := map[string]struct{}{}
		for _, fd := range msg.Fields {
			if _, ok := msgNames[fd.Name]; ok {
				return fmt.Errorf("message %s has duplicate field name %s", msg.Name, fd.Name)
			}

			msgNames[fd.Name] = struct{}{}
		}
		if msg.OpCode != 0 {
			if conflict, ok := allOpCodes[msg.OpCode]; ok {
				return fmt.Errorf("message %s has duplicate opcode %02x (duplicated in %s)", msg.Name, msg.OpCode, conflict)
			}
			allOpCodes[msg.OpCode] = msg.Name
		}
	}
	for _, un := range f.Unions {
		if _, ok := primitiveTypes[un.Name]; ok {
			return fmt.Errorf("union shares primitive type name %s", un.Name)
		}
		if _, ok := customTypes[un.Name]; ok {
			return fmt.Errorf("union has duplicated name %s", un.Name)
		}
		if un.Namespace != "" {
			nameWithoutPrefix := strings.TrimPrefix(un.Name, un.Namespace+".")
			customTypes[nameWithoutPrefix] = struct{}{}
		}
		customTypes[un.Name] = struct{}{}
		unionNames := map[string]struct{}{}
		for _, fd := range un.Fields {
			if _, ok := unionNames[fd.name()]; ok {
				return fmt.Errorf("union %s has duplicate field name %s", un.Name, fd.name())
			}
			unionNames[fd.name()] = struct{}{}
		}
		if un.OpCode != 0 {
			if conflict, ok := allOpCodes[un.OpCode]; ok {
				return fmt.Errorf("union %s has duplicate opcode %02x (duplicated in %s)", un.Name, un.OpCode, conflict)
			}
			allOpCodes[un.OpCode] = un.Name
		}
	}
	allTypes := customTypes
	for typ := range primitiveTypes {
		allTypes[typ] = struct{}{}
	}
	for _, st := range f.Structs {
		for _, fd := range st.Fields {
			if err := typeDefined(fd.FieldType, allTypes); err != nil {
				return err
			}
		}
	}
	for _, msg := range f.Messages {
		for _, fd := range msg.Fields {
			if err := typeDefined(fd.FieldType, allTypes); err != nil {
				return err
			}
		}
	}

	// Todo: within a given struct, enum, or message, a field / option cannot
	// have a duplicate name

	// Determine which structs include themselves as required fields (which would lead
	// to the struct taking up infinite size)
	delta := true
	// This is an unbounded loop because for this case:
	// fooStruct:{barStruct}, barStruct{bizzStruct}, bizzStruct{bazStruct}, bazStruct{fooStruct}
	// After each iteration we have these updated usages:
	// 1. fooStruct:{barStruct, bizzStruct}, barStruct{bizzStruct, bazStruct}, bizzStruct{bazStruct, fooStruct}, bazStruct{fooStruct, barStruct}
	// 2. fooStruct:{barStruct, bizzStruct, bazStruct, fooStruct}, barStruct{bizzStruct, bazStruct, fooStruct, barStruct}, bizzStruct{bazStruct, fooStruct, barStruct, bizzStruct}, bazStruct{fooStruct, barStruct, bizzStruct, bazStruct}
	// ... and as the chain of structs gets longer the required iterations also increases.
	for delta {
		delta = false
		for stName, usage := range structTypeUsage {
			for stName2, usage2 := range structTypeUsage {
				if stName == stName2 {
					continue
				}
				// If struct1 includes struct2, it also includes
				// all fields that struct2 includes
				if usage[stName2] {
					for k, v := range usage2 {
						if !usage[k] {
							delta = true
						}
						usage[k] = v
					}
				}
			}
			structTypeUsage[stName] = usage
		}
	}
	for stName, usage := range structTypeUsage {
		if usage[stName] {
			return fmt.Errorf("struct %s recursively includes itself as a required field", stName)
		}
	}

	return nil
}

func typeDefined(ft FieldType, allTypes map[string]struct{}) error {
	if ft.Array != nil {
		return typeDefined(*ft.Array, allTypes)
	}
	if ft.Map != nil {
		if _, ok := allTypes[ft.Map.Key]; !ok {
			return fmt.Errorf("map key type %s undefined", ft.Map.Key)
		}
		return typeDefined(ft.Map.Value, allTypes)
	}
	if _, ok := allTypes[ft.Simple]; !ok {
		return fmt.Errorf("type %s undefined", ft.Simple)
	}
	return nil
}

// Generate writes a .go file out to w. If f has imports, it will
// parse the files at those imports and generate them according to
// settings.
func (f File) Generate(w io.Writer, settings GenerateSettings) error {
	if err := settings.Validate(); err != nil {
		return fmt.Errorf("invalid generation settings: %w", err)
	}
	settings.nextLength = new(int)
	settings.isFirstTopLength = new(bool)

	if len(f.Imports) != 0 {
		thisFilePath := f.FileName
		if !path.IsAbs(f.FileName) {
			// assume relative to our current directory
			wd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("cannot determine working directory for import lookup: %v", err)
			}
			thisFilePath = filepath.Join(wd, thisFilePath)
		}
		thisDir := filepath.Dir(thisFilePath)
		type bebopImport struct {
			from string
			to   string
		}
		imports := make([]bebopImport, len(f.Imports))
		importGraph := importgraph.NewDgraph()
		for i, imp := range f.Imports {
			imports[i] = bebopImport{
				from: f.GoPackage,
				to:   imp,
			}
		}
		// TODO: why are imports not scoped to a namespace?
		imported := map[string]struct{}{}
		for i := 0; i < len(imports); i++ {
			imp := imports[i]
			impPath := filepath.Join(thisDir, imp.to)

			impF, err := os.Open(impPath)
			if err != nil {
				return fmt.Errorf("failed to open imported file %s: %w", imp.to, err)
			}
			impFile, _, err := ReadFile(impF)
			if err != nil {
				impF.Close()
				return fmt.Errorf("failed to parse imported file %s: %w", imp.to, err)
			}
			impF.Close()
			importGraph.AddEdge(imp.from, impFile.GoPackage)
			if _, ok := imported[impPath]; ok {
				continue
			}
			settings.imported = append(settings.imported, impFile)
			for _, subImp := range impFile.Imports {
				imports = append(imports, bebopImport{
					from: impFile.GoPackage,
					to:   subImp,
				})
			}
			imported[impPath] = struct{}{}
		}
		if settings.ImportGenerationMode == ImportGenerationModeSeparate {
			if err := importGraph.FindCycle(); err != nil {
				return err
			}
		}
	}
	imports := []string{}
	potentialImports := []string{}
	settings.importTypeAliases = make(map[string]string)
	switch settings.ImportGenerationMode {
	case ImportGenerationModeSeparate:
		for _, imp := range settings.imported {
			if imp.GoPackage == "" {
				return fmt.Errorf("cannot import %s: file has no %s const", filepath.Base(imp.FileName), goPackage)
			}
			packageName := imp.GoPackage
			packageNamespace := path.Base(packageName)
			for _, st := range imp.Structs {
				st.Namespace = packageNamespace
				namespacedName := packageNamespace + "." + st.Name
				settings.importTypeAliases[st.Name] = namespacedName
				st.Name = namespacedName
				f.Structs = append(f.Structs, st)
			}
			for _, un := range imp.Unions {
				un.Namespace = packageNamespace
				namespacedName := packageNamespace + "." + un.Name
				settings.importTypeAliases[un.Name] = namespacedName
				un.Name = namespacedName
				f.Unions = append(f.Unions, un)
			}
			for _, msg := range imp.Messages {
				msg.Namespace = packageNamespace
				namespacedName := packageNamespace + "." + msg.Name
				settings.importTypeAliases[msg.Name] = namespacedName
				msg.Name = namespacedName
				f.Messages = append(f.Messages, msg)
			}
			for _, enm := range imp.Enums {
				enm.Namespace = packageNamespace
				namespacedName := packageNamespace + "." + enm.Name
				settings.importTypeAliases[enm.Name] = namespacedName
				enm.Name = namespacedName
				f.Enums = append(f.Enums, enm)
			}
			potentialImports = append(potentialImports, packageName)
		}
	case ImportGenerationModeCombined:
		// treat all imported files as a part of this file. Do not observe GoPackage.
		for _, imp := range settings.imported {
			f.Consts = append(f.Consts, imp.Consts...)
			f.Structs = append(f.Structs, imp.Structs...)
			f.Unions = append(f.Unions, imp.Unions...)
			f.Messages = append(f.Messages, imp.Messages...)
			f.Enums = append(f.Enums, imp.Enums...)
		}
	}

	if err := f.Validate(); err != nil {
		return fmt.Errorf("cannot generate file: %w", err)
	}
	settings.typeMarshallers = f.typeMarshallers()
	settings.typeByters = f.typeByters()
	settings.typeByteReaders = f.typeByteReaders(settings)
	settings.typeUnmarshallers = f.typeUnmarshallers(settings)
	settings.typeLengthers = f.typeLengthers()
	settings.customRecordTypes = f.customRecordTypes()

	usedTypes := f.usedTypes()
	if settings.PackageName == "" && f.GoPackage != "" {
		settings.PackageName = path.Base(f.GoPackage)
	} else if settings.PackageName == "" {
		return fmt.Errorf("no package name is defined, provide a %s const or an explicit package name setting", goPackage)
	}
	for _, imp := range potentialImports {
		for typ := range usedTypes {
			if alias, ok := settings.importTypeAliases[typ]; ok {
				if strings.HasPrefix(alias, path.Base(imp)+".") {
					imports = append(imports, imp)
					break
				}
			}
		}
	}

	writeLine(w, "// Code generated by bebopc-go; DO NOT EDIT.")
	writeLine(w, "")
	writeLine(w, "package %s", settings.PackageName)
	writeLine(w, "")

	if len(f.Messages)+len(f.Structs)+len(f.Unions) != 0 {
		imports = append(imports, "github.com/200sc/bebop")
		imports = append(imports, "github.com/200sc/bebop/iohelp")
	}

	if len(f.Messages)+len(f.Structs)+len(f.Unions) != 0 {
		imports = append(imports, "io")
	}
	for _, c := range f.Consts {
		if c.impossibleGoConst() {
			imports = append(imports, "math")
			break
		}
	}
	if usedTypes[typeDate] {
		imports = append(imports, "time")
	}

	if len(imports) != 0 {
		writeLine(w, "import (")
		for _, i := range imports {
			writeLine(w, "\t%q", i)
		}
		writeLine(w, ")")
		writeLine(w, "")
	}

	impossibleGoConsts := []Const{}

	if len(f.Consts) != 0 {
		writeLine(w, "const (")
		for _, con := range f.Consts {
			if con.impossibleGoConst() {
				impossibleGoConsts = append(impossibleGoConsts, con)
			} else {
				con.Generate(w, settings)
			}
		}
		writeLine(w, ")")
		writeLine(w, "")
	}
	if len(impossibleGoConsts) != 0 {
		writeLine(w, "var (")
		for _, con := range impossibleGoConsts {
			con.Generate(w, settings)
		}
		writeLine(w, ")")
		writeLine(w, "")
	}
	// Namespaced types are imported from another package, and must not be generated.
	// They must, however, be defined up til this point so we know how to create them
	// as components of records in this package.
	for _, en := range f.Enums {
		if en.Namespace != "" {
			continue
		}
		en.Generate(w, settings)
	}
	for _, st := range f.Structs {
		if st.Namespace != "" {
			continue
		}
		st.Generate(w, settings)
	}
	for _, msg := range f.Messages {
		if msg.Namespace != "" {
			continue
		}
		msg.Generate(w, settings)
	}
	for _, union := range f.Unions {
		if union.Namespace != "" {
			continue
		}
		union.Generate(w, settings)
	}
	return nil
}

// Generate writes a .go enum definition out to w.
func (en Enum) Generate(w io.Writer, settings GenerateSettings) {
	exposedName := exposeName(en.Name, settings)
	writeComment(w, 0, en.Comment, settings)
	writeLine(w, "type %s %s", exposedName, en.SimpleType)
	writeLine(w, "")
	if len(en.Options) != 0 {
		writeLine(w, "const (")
		for _, opt := range en.Options {
			writeComment(w, 1, opt.Comment, settings)
			if opt.Deprecated {
				writeLine(w, "\t// Deprecated: %s", opt.DeprecatedMessage)
			}
			if en.Unsigned {
				writeLine(w, "\t%s_%s %s = %d", exposedName, opt.Name, exposedName, opt.UintValue)
			} else {
				writeLine(w, "\t%s_%s %s = %d", exposedName, opt.Name, exposedName, opt.Value)
			}
		}
		writeLine(w, ")")
		writeLine(w, "")
	}
}

func (con Const) impossibleGoConst() bool {
	// unique floating point values (inf, nan) cannot be represented as consts in go,
	// at least not intuitively. This probably doesn't matter-- its rare to
	// rely on inf or nan floating points anyway, and even if you could get these as
	// consts you would need to use math.IsInf or math.IsNaN for many use cases.
	if con.SimpleType == typeFloat32 || con.SimpleType == typeFloat64 {
		switch con.Value {
		case "math.Inf(-1)":
			return true
		case "math.Inf(1)":
			return true
		case "math.NaN()":
			return true
		}
	}
	return false
}

func (con Const) Generate(w io.Writer, settings GenerateSettings) {
	writeComment(w, 0, con.Comment, settings)
	writeLine(w, "\t%s = %v", exposeName(con.Name, settings), con.Value)
}

func writeFieldDefinition(fd Field, w io.Writer, readOnly bool, message bool, settings GenerateSettings) {
	writeComment(w, 1, fd.Comment, settings)
	if fd.Deprecated {
		writeLine(w, "\t// Deprecated: %s", fd.DeprecatedMessage)
	}

	name := exposeName(fd.Name, settings)
	if readOnly {
		name = unexposeName(fd.Name)
	}
	typ := fd.FieldType.goString(settings)
	if message {
		typ = "*" + typ
	}
	if settings.GenerateFieldTags && len(fd.Tags) != 0 {
		formattedTags := []string{}
		for _, tag := range fd.Tags {
			if tag.Boolean {
				formattedTags = append(formattedTags, tag.Key)
			} else {
				formattedTags = append(formattedTags, fmt.Sprintf("%s:%q", tag.Key, tag.Value))
			}
		}
		writeLine(w, "\t%s %s `%s`", name, typ, strings.Join(formattedTags, " "))
	} else {
		writeLine(w, "\t%s %s", name, typ)
	}
}

func depthName(name string, depth int) string {
	return name + strconv.Itoa(depth)
}

func lengthName(settings GenerateSettings) string {
	(*settings.nextLength)++
	return "ln" + strconv.Itoa(*settings.nextLength)
}

func writeFieldByter(name string, typ FieldType, w io.Writer, settings GenerateSettings, depth int) {
	if typ.Array != nil {
		writeLineWithTabs(w, "iohelp.WriteUint32Bytes(buf[at:], uint32(len(%ASGN)))", depth, name)
		writeLineWithTabs(w, "at += 4", depth)
		if typ.Array.Simple == typeByte || typ.Array.Simple == typeUint8 {
			writeLineWithTabs(w, "copy(buf[at:at+len(%ASGN)], %ASGN)", depth, name)
			writeLineWithTabs(w, "at += len(%ASGN)", depth, name)
			return
		}
		writeLineWithTabs(w, "for _, %VNAME := range %ASGN {", depth, name)
		writeFieldByter(depthName("v", depth), *typ.Array, w, settings, depth+1)
		writeLineWithTabs(w, "}", depth)
	} else if typ.Map != nil {
		writeLineWithTabs(w, "iohelp.WriteUint32Bytes(buf[at:], uint32(len(%ASGN)))", depth, name)
		writeLineWithTabs(w, "at += 4", depth)
		writeLineWithTabs(w, "for %KNAME, %VNAME := range %ASGN {", depth, name)
		writeLineWithTabs(w, settings.typeByters[typ.Map.Key], depth+1, depthName("k", depth))
		writeFieldByter(depthName("v", depth), typ.Map.Value, w, settings, depth+1)
		writeLineWithTabs(w, "}", depth)
	} else {
		simpleTyp := typ.Simple
		if alias, ok := settings.importTypeAliases[simpleTyp]; ok {
			simpleTyp = alias
		}
		writeLineWithTabs(w, settings.typeByters[simpleTyp], depth, name, typ.goString(settings))
	}
}

func writeLengthCheck(w io.Writer, ln string, depth int, args ...string) {
	writeLineWithTabs(w, "if len(buf[at:]) < "+ln+" {", depth, args...)
	writeLineWithTabs(w, "\treturn io.ErrUnexpectedEOF", depth, args...)
	writeLineWithTabs(w, "}", depth, args...)
}

func writeFieldReadByter(name string, typ FieldType, w io.Writer, settings GenerateSettings, depth int, safe bool) {
	if typ.Array != nil {
		if safe {
			writeLengthCheck(w, "4", depth)
		}

		writeLineWithTabs(w, "%ASGN = make([]%TYPE, iohelp.ReadUint32Bytes(buf[at:]))", depth, name, typ.Array.goString(settings))
		writeLineWithTabs(w, "at += 4", depth)
		if safe {
			if sz, ok := fixedSizeTypes[typ.Array.Simple]; ok {
				writeLengthCheck(w, "len(%ASGN)*"+strconv.Itoa(int(sz)), depth, name)
				safe = false
			}
		}
		if typ.Array.Simple == typeByte || typ.Array.Simple == typeUint8 {
			writeLineWithTabs(w, "copy(%ASGN, buf[at:at+len(%ASGN)])", depth, name)
			writeLineWithTabs(w, "at += len(%ASGN)", depth, name)
			return
		}
		iName := depthName("i", depth)
		writeLineWithTabs(w, "for "+iName+" := range %ASGN {", depth, name)
		writeFieldReadByter("("+name+")["+iName+"]", *typ.Array, w, settings, depth+1, safe)
		writeLineWithTabs(w, "}", depth)
	} else if typ.Map != nil {
		lnName := lengthName(settings)
		writeLineWithTabs(w, lnName+" := iohelp.ReadUint32Bytes(buf[at:])", depth)
		writeLineWithTabs(w, "at += 4", depth)
		writeLineWithTabs(w, "%ASGN = make(%TYPE,"+lnName+")", depth, name, typ.Map.goString(settings))
		writeLineWithTabs(w, "for i := uint32(0); i < "+lnName+"; i++ {", depth, name)
		var ln string
		if format, ok := settings.typeByteReaders[typ.Map.Key+"&safe"]; ok && safe {
			ln = getLineWithTabs(format, depth+1, depthName("k", depth), simpleGoString(typ.Map.Key, settings))
		} else {
			if sz, ok := fixedSizeTypes[typ.Map.Key]; ok && safe {
				writeLengthCheck(w, strconv.Itoa(int(sz)), depth+1, depthName("k", depth))
			}
			ln = getLineWithTabs(settings.typeByteReaders[typ.Map.Key], depth+1, depthName("k", depth), typ.goString(settings))
		}
		w.Write([]byte(strings.Replace(ln, "=", ":=", 1)))
		writeFieldReadByter("("+name+")["+depthName("k", depth)+"]", typ.Map.Value, w, settings, depth+1, safe)
		writeLineWithTabs(w, "}", depth)
	} else {
		simpleTyp := typ.Simple
		if alias, ok := settings.importTypeAliases[simpleTyp]; ok {
			simpleTyp = alias
		}
		if format, ok := settings.typeByteReaders[simpleTyp+"&safe"]; ok && safe {
			writeLineWithTabs(w, format, depth, name, typ.goString(settings))
		} else {
			if sz, ok := fixedSizeTypes[simpleTyp]; ok && safe {
				writeLengthCheck(w, strconv.Itoa(int(sz)), depth, name)
			}
			writeLineWithTabs(w, settings.typeByteReaders[simpleTyp], depth, name, typ.goString(settings))
		}
	}
}

func writeFieldMarshaller(name string, typ FieldType, w io.Writer, settings GenerateSettings, depth int) {
	if typ.Array != nil {
		writeLineWithTabs(w, "iohelp.WriteUint32(w, uint32(len(%ASGN)))", depth, name)
		if typ.Array.Simple == typeByte {
			writeLineWithTabs(w, "w.Write(%ASGN)", depth, name)
		} else {
			writeLineWithTabs(w, "for _, elem := range %ASGN {", depth, name)
			writeFieldMarshaller("elem", *typ.Array, w, settings, depth+1)
			writeLineWithTabs(w, "}", depth)
		}
	} else if typ.Map != nil {
		writeLineWithTabs(w, "iohelp.WriteUint32(w, uint32(len(%ASGN)))", depth, name)
		writeLineWithTabs(w, "for %KNAME, %VNAME := range %ASGN {", depth, name)
		writeLineWithTabs(w, settings.typeMarshallers[typ.Map.Key], depth+1, depthName("k", depth))
		writeFieldMarshaller(depthName("v", depth), typ.Map.Value, w, settings, depth+1)
		writeLineWithTabs(w, "}", depth)
	} else {
		simpleTyp := typ.Simple
		if alias, ok := settings.importTypeAliases[simpleTyp]; ok {
			simpleTyp = alias
		}
		writeLineWithTabs(w, settings.typeMarshallers[simpleTyp], depth, name, typ.goString(settings))
	}
}

func typeNeedsElem(typ string, settings GenerateSettings) bool {
	switch typ {
	case "":
		return true
	case typeString:
		return true
	}
	if _, ok := primitiveTypes[typ]; ok {
		return false
	}
	_, ok := settings.customRecordTypes[typ]
	return ok
}

func writeFieldBodyCount(name string, typ FieldType, w io.Writer, settings GenerateSettings, depth int) {
	if typ.Array != nil {
		writeLineWithTabs(w, "bodyLen += 4", depth)
		if sz, ok := fixedSizeTypes[typ.Array.Simple]; ok {
			// short circuit-- write length times elem size
			writeLineWithTabs(w, "bodyLen += len(%ASGN) * "+strconv.Itoa(int(sz)), depth, name)
			return
		}
		writeLineWithTabs(w, "for _, elem := range %ASGN {", depth, name)
		writeFieldBodyCount("elem", *typ.Array, w, settings, depth+1)
		writeLineWithTabs(w, "}", depth)
	} else if typ.Map != nil {
		writeLineWithTabs(w, "bodyLen += 4", depth, name)
		useV := typeNeedsElem(typ.Map.Value.Simple, settings)
		useK := typ.Map.Key == typeString
		if useV && useK {
			writeLineWithTabs(w, "for %KNAME, %VNAME := range %ASGN {", depth, name)
		} else if useV {
			writeLineWithTabs(w, "for _, %VNAME := range %ASGN {", depth, name)
		} else if useK {
			writeLineWithTabs(w, "for %KNAME := range %ASGN {", depth, name)
		} else {
			writeLineWithTabs(w, "for range %ASGN {", depth, name)
		}
		writeLineWithTabs(w, settings.typeLengthers[typ.Map.Key], depth+1, depthName("k", depth))
		writeFieldBodyCount(depthName("v", depth), typ.Map.Value, w, settings, depth+1)
		writeLineWithTabs(w, "}", depth)
	} else {
		simpleTyp := typ.Simple
		if alias, ok := settings.importTypeAliases[simpleTyp]; ok {
			simpleTyp = alias
		}
		writeLineWithTabs(w, settings.typeLengthers[simpleTyp], depth, name, typ.goString(settings))
	}
}

func exposeName(name string, settings GenerateSettings) string {
	if settings.PrivateDefinitions {
		return unexposeName(name)
	}
	if name == "" {
		return ""
	}
	return strings.ToUpper(string(name[0])) + name[1:]
}

func unexposeName(name string) string {
	if name == "" {
		return ""
	}
	return strings.ToLower(string(name[0])) + name[1:]
}

func writeWrappers(w io.Writer, name string, isEmpty bool, settings GenerateSettings) {
	writeMarshalBebop(w, name, isEmpty, settings)
	writeMake(w, name, isEmpty, settings)
	writeMakeFromBytes(w, name, isEmpty, settings)
	if settings.GenerateUnsafeMethods {
		writeMustMakeFromBytes(w, name, isEmpty, settings)
	}
}

func writeMarshalBebop(w io.Writer, name string, isEmpty bool, settings GenerateSettings) {
	exposedName := exposeName(name, settings)
	if settings.AlwaysUsePointerReceivers {
		writeLine(w, "func (bbp *%s) MarshalBebop() []byte {", exposedName)
	} else {
		writeLine(w, "func (bbp %s) MarshalBebop() []byte {", exposedName)
	}
	if isEmpty {
		writeLine(w, "\treturn []byte{}")
	} else {
		writeLine(w, "\tbuf := make([]byte, bbp.Size())")
		writeLine(w, "\tbbp.MarshalBebopTo(buf)")
		writeLine(w, "\treturn buf")
	}
	writeCloseBlock(w)
}

func writeMake(w io.Writer, name string, isEmpty bool, settings GenerateSettings) {
	exposedName := exposeName(name, settings)
	makeName := exposeName("Make", settings)
	writeLine(w, "func %[2]s%[1]s(r *iohelp.ErrorReader) (%[1]s, error) {", exposedName, makeName)
	if isEmpty {
		writeLine(w, "\treturn %s{}, nil", exposedName)
	} else {
		writeLine(w, "\tv := %s{}", exposedName)
		writeLine(w, "\terr := v.DecodeBebop(r)")
		writeLine(w, "\treturn v, err")
	}
	writeCloseBlock(w)
}

func writeMakeFromBytes(w io.Writer, name string, isEmpty bool, settings GenerateSettings) {
	exposedName := exposeName(name, settings)
	makeName := exposeName("Make", settings)
	writeLine(w, "func %[2]s%[1]sFromBytes(buf []byte) (%[1]s, error) {", exposedName, makeName)
	if isEmpty {
		writeLine(w, "\treturn %s{}, nil", exposedName)
	} else {
		writeLine(w, "\tv := %s{}", exposedName)
		writeLine(w, "\terr := v.UnmarshalBebop(buf)")
		writeLine(w, "\treturn v, err")
	}
	writeCloseBlock(w)
}

func writeMustMakeFromBytes(w io.Writer, name string, isEmpty bool, settings GenerateSettings) {
	exposedName := exposeName(name, settings)
	makeName := exposeName("MustMake", settings)
	writeLine(w, "func %[2]s%[1]sFromBytes(buf []byte) %[1]s {", exposedName, makeName)
	if isEmpty {
		writeLine(w, "\treturn %s{}", exposedName)
	} else {
		writeLine(w, "\tv := %s{}", exposedName)
		writeLine(w, "\tv.MustUnmarshalBebop(buf)")
		writeLine(w, "\treturn v")
	}
	writeCloseBlock(w)
}

func writeLine(w io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(w, format+"\n", args...)
}

func getLineWithTabs(format string, depth int, args ...string) string {
	b := new(bytes.Buffer)
	writeLineWithTabs(b, format, depth, args...)
	return b.String()
}

func writeComment(w io.Writer, depth int, comment string, settings GenerateSettings) {
	if comment == "" {
		return
	}
	tbs := strings.Repeat("\t", depth)

	commentLines := strings.Split(comment, "\n")
	for _, cm := range commentLines {
		// If you have tag comments and are generating them as tags,
		// you probably don't want them showing up in your code as comments too.
		if settings.GenerateFieldTags {
			if _, ok := parseCommentTag(cm); ok {
				continue
			}
		}
		writeLine(w, tbs+"//%s", cm)
	}
}

func writeCloseBlock(w io.Writer) {
	writeLine(w, "}")
	writeLine(w, "")
}

func writeOpCode(w io.Writer, name string, opCode uint32, settings GenerateSettings) {
	if opCode != 0 {
		writeLine(w, "const %sOpCode = 0x%x", exposeName(name, settings), opCode)
		writeLine(w, "")
	}
}

func writeRecordAssertion(w io.Writer, name string, settings GenerateSettings) {
	writeLine(w, "var _ bebop.Record = &%s{}", exposeName(name, settings))
	writeLine(w, "")
}

func writeGoStructDef(w io.Writer, name string, settings GenerateSettings) {
	writeLine(w, "type %s struct {", exposeName(name, settings))
}

func writeRecordTypeDefinition(w io.Writer, name string, opCode uint32, comment string, settings GenerateSettings, fields []fieldWithNumber) {
	writeOpCode(w, name, opCode, settings)
	writeRecordAssertion(w, name, settings)
	writeComment(w, 0, comment, settings)
	writeGoStructDef(w, name, settings)
	for _, fd := range fields {
		writeFieldDefinition(fd.Field, w, false, true, settings)
	}
	writeCloseBlock(w)
}
