package bebop

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type FileNamer interface {
	Name() string
}

// ReadFile reads out a bebop file. If r is a FileNamer, like an *os.File,
// the output's FileName will be populated. In addition to fatal errors, string
// warnings may also be output.
func ReadFile(r io.Reader) (File, []string, error) {
	f := File{}
	if fnamer, ok := r.(FileNamer); ok {
		f.FileName = fnamer.Name()
	}
	tr := newTokenReader(r)
	nextCommentLines := []string{}
	nextRecordOpCode := uint32(0)
	nextRecordReadOnly := false
	nextRecordBitFlags := false
	nextDecorations := Decorations{Custom: map[string]string{}}
	warnings := []string{}
	for tr.Next() {
		tk := tr.Token()
		switch tk.kind {
		// file headers
		case tokenKindImport:
			toks, err := expectNext(tr, tokenKindStringLiteral)
			if err != nil {
				return f, warnings, err
			}
			// This cannot fail; string literals are always quoted correctly
			imported, _ := strconv.Unquote(string(toks[0].concrete))
			f.Imports = append(f.Imports, imported)
			continue
			// comments and whitespace
		case tokenKindNewline:
			nextCommentLines = []string{}
			continue
		case tokenKindBlockComment:
			nextCommentLines = append(nextCommentLines, readBlockComment(tr, tk))
			continue
		case tokenKindLineComment:
			nextCommentLines = append(nextCommentLines, sanitizeComment(tk))
			continue
			// annotations
		case tokenKindOpenSquare, tokenKindAtSign:
			k, v, opcodeV, err := readDecoration(tr, tk.kind == tokenKindOpenSquare)
			if err != nil {
				return f, warnings, err
			}
			// TODO: are decorator keys case sensitive
			// TODO: can you put multiple decorators in a row
			// TODO: are semicolons required, allowed, disallowed after decorators
			// TODO: if you put the same decorator twice, what errors
			switch k {
			case "opcode":
				nextRecordOpCode = opcodeV
			case "flags":
				nextRecordBitFlags = true
			case "deprecated":
				if nextDecorations.Deprecated {
					return f, warnings, readError(tk, "deprecated cannot be applied to the same target twice")
				}
				nextDecorations.Deprecated = true
				nextDecorations.DeprecatedMessage = v
			default:
				nextDecorations.Custom[k] = v
			}
			continue
			// top level declarations / statements
		case tokenKindEnum:
			if nextRecordOpCode != 0 {
				return f, warnings, readError(tk, "enums may not have attached op codes")
			}
			en, err := readEnum(tr, nextRecordBitFlags)
			if err != nil {
				return f, warnings, err
			}
			en.Comment = strings.Join(nextCommentLines, "\n")
			en.Decorations = nextDecorations
			nextDecorations = Decorations{Custom: map[string]string{}}

			f.Enums = append(f.Enums, en)
		case tokenKindReadOnly:
			nextRecordReadOnly = true
			if !tr.Next() {
				return f, warnings, readError(tk, "expected (Struct) got no token")
			}
			tk = tr.Token()
			if tk.kind != tokenKindStruct {
				return f, warnings, readError(tk, "expected (Struct) got (%v)", tk.kind)
			}
			fallthrough
		case tokenKindStruct:
			if nextRecordBitFlags {
				return f, warnings, readError(tk, "structs may not use bitflags")
			}
			st, err := readStruct(tr)
			if err != nil {
				return f, warnings, err
			}
			st.Comment = strings.Join(nextCommentLines, "\n")
			st.OpCode = nextRecordOpCode
			st.ReadOnly = nextRecordReadOnly
			nextRecordReadOnly = false

			st.Decorations = nextDecorations
			nextDecorations = Decorations{Custom: map[string]string{}}

			f.Structs = append(f.Structs, st)
		case tokenKindMessage:
			if nextRecordBitFlags {
				return f, warnings, readError(tk, "messages may not use bitflags")
			}
			msg, err := readMessage(tr)
			if err != nil {
				return f, warnings, err
			}
			msg.Comment = strings.Join(nextCommentLines, "\n")
			msg.OpCode = nextRecordOpCode

			msg.Decorations = nextDecorations
			nextDecorations = Decorations{Custom: map[string]string{}}

			f.Messages = append(f.Messages, msg)
		case tokenKindUnion:
			if nextRecordBitFlags {
				return f, warnings, readError(tk, "unions may not use bitflags")
			}
			union, err := readUnion(tr)
			if err != nil {
				return f, warnings, err
			}
			union.Comment = strings.Join(nextCommentLines, "\n")
			union.OpCode = nextRecordOpCode

			union.Decorations = nextDecorations
			nextDecorations = Decorations{Custom: map[string]string{}}

			f.Unions = append(f.Unions, union)
		case tokenKindConst:
			if nextRecordBitFlags {
				return f, warnings, readError(tk, "consts may not use bitflags")
			}
			if nextRecordOpCode != 0 {
				return f, warnings, readError(tk, "consts may not have attached op codes")
			}
			cons, constWarnings, err := readConst(tr)
			warnings = append(warnings, constWarnings...)
			if err != nil {
				return f, warnings, err
			}
			cons.Comment = strings.Join(nextCommentLines, "\n")
			if cons.Name == goPackage && cons.SimpleType == typeString {
				f.GoPackage, _ = strconv.Unquote(cons.Value)
			}

			cons.Decorations = nextDecorations
			nextDecorations = Decorations{Custom: map[string]string{}}

			f.Consts = append(f.Consts, cons)
		}
		nextCommentLines = []string{}
		nextRecordOpCode = 0
	}
	if nextDecorations.Deprecated || len(nextDecorations.Custom) != 0 || nextRecordOpCode != 0 {
		return f, warnings, readError(tr.Token(), "file ended with unattached decoration")
	}
	return f, warnings, nil
}

func expectAnyOfNext(tr *tokenReader, kinds ...tokenKind) error {
	next := tr.Next()
	if tr.Err() != nil {
		return tr.Err()
	}
	if !next {
		kindsStrs := make([]string, len(kinds))
		for i, k := range kinds {
			kindsStrs[i] = k.String()
		}
		return readError(tr.nextToken, "expected (%v), got no token", kindsStr(kinds))
	}
	tk := tr.Token()
	found := false
	for _, k := range kinds {
		if tk.kind == k {
			found = true
		}
	}
	if !found {
		kindsStrs := make([]string, len(kinds))
		for i, k := range kinds {
			kindsStrs[i] = k.String()
		}
		return readError(tk, "expected (%v) got %s", kindsStr(kinds), tk.kind)
	}
	return nil
}

func expectNext(tr *tokenReader, kinds ...tokenKind) ([]token, error) {
	tokens := make([]token, len(kinds))
	for i, k := range kinds {
		next := tr.Next()
		if tr.Err() != nil {
			return tokens, tr.Err()
		}
		if !next {
			return tokens, readError(tr.nextToken, "expected (%v), got no token", kindsStr(kinds))
		}
		tk := tr.Token()
		if tk.kind != k {
			return tokens, readError(tk, "expected (%v) got %s", k, tk.kind)
		}
		tokens[i] = tk
	}
	return tokens, nil
}

func optNewline(tr *tokenReader) {
	tr.Next()
	if tr.Token().kind != tokenKindNewline {
		tr.UnNext()
	}
}

func readUntil(tr *tokenReader, kind tokenKind) ([]token, error) {
	toks := []token{}
	for tr.Next() {
		tk := tr.Token()
		if tk.kind == kind {
			return toks, nil
		}
		toks = append(toks, tk)
	}
	return nil, readError(tr.lastToken, "eof reading until %v", kind)
}

func readEnumOptionValue(tr *tokenReader, previousOptions []EnumOption, bitflags, uinttype bool, bitsize int) (int64, uint64, error) {
	if _, err := expectNext(tr, tokenKindEquals); err != nil {
		return 0, 0, err
	}
	if !bitflags {
		toks, err := expectNext(tr, tokenKindIntegerLiteral, tokenKindSemicolon)
		if err != nil {
			return 0, 0, err
		}
		if uinttype {
			optInteger, err := strconv.ParseUint(string(toks[0].concrete), 0, bitsize)
			if err != nil {
				return 0, 0, err
			}
			return 0, optInteger, nil
		} else {
			optInteger, err := strconv.ParseInt(string(toks[0].concrete), 0, bitsize)
			if err != nil {
				return 0, 0, err
			}
			return optInteger, 0, nil
		}
	}
	return readBitflagExpr(tr, previousOptions, uinttype, bitsize)
}

func readEnum(tr *tokenReader, bitflags bool) (Enum, error) {
	en := Enum{
		SimpleType: "uint32",
	}

	toks, err := expectNext(tr, tokenKindIdent)
	if err != nil {
		return en, err
	}
	en.Name = string(toks[0].concrete)
	err = expectAnyOfNext(tr, tokenKindColon, tokenKindOpenCurly)
	if err != nil {
		return en, err
	}

	switch tr.Token().kind {
	case tokenKindOpenCurly:
		break
	case tokenKindColon:
		enumSizeTokens, err := expectNext(tr, tokenKindIdent, tokenKindOpenCurly)
		if err != nil {
			return en, err
		}
		enumSize := string(enumSizeTokens[0].concrete)
		if !isUintPrimitive(enumSize) && !isIntPrimitive(enumSize) {
			return en, readError(enumSizeTokens[0], "expected an integer enum type")
		}
		en.SimpleType = enumSize
	}
	optNewline(tr)

	bitsize, uinttype := decodeIntegerType(en.SimpleType)
	en.Unsigned = uinttype
	nextCommentLines := []string{}
	nextDecorations := Decorations{Custom: map[string]string{}}

	for tr.Token().kind != tokenKindCloseCurly {
		if !tr.Next() {
			return en, readError(tr.nextToken, "enum definition ended early")
		}
		tk := tr.Token()
		switch tk.kind {
		case tokenKindNewline:
			nextCommentLines = []string{}
		case tokenKindIdent:
			optName := string(tk.concrete)

			signedValue, unsignedValue, err := readEnumOptionValue(tr, en.Options, bitflags, uinttype, bitsize)
			if err != nil {
				return en, err
			}
			en.Options = append(en.Options, EnumOption{
				Name:        optName,
				Value:       signedValue,
				UintValue:   unsignedValue,
				Decorations: nextDecorations,
				Comment:     strings.Join(nextCommentLines, "\n"),
			})
			nextDecorations = Decorations{Custom: map[string]string{}}

			nextCommentLines = []string{}
		case tokenKindOpenSquare, tokenKindAtSign:
			k, v, _, err := readDecoration(tr, tk.kind == tokenKindOpenSquare)
			if err != nil {
				return en, err
			}
			switch k {
			case "opcode":
				return en, readError(tr.nextToken, "opcode annotation not allowed within enum")
			case "flags":
				return en, readError(tr.nextToken, "flags annotation not allowed within enum")
			case "deprecated":
				if nextDecorations.Deprecated {
					return en, readError(tk, "deprecated cannot be applied to the same target twice")
				}
				nextDecorations.Deprecated = true
				nextDecorations.DeprecatedMessage = v
			default:
				nextDecorations.Custom[k] = v
			}
		case tokenKindBlockComment:
			nextCommentLines = append(nextCommentLines, readBlockComment(tr, tk))
		case tokenKindLineComment:
			nextCommentLines = append(nextCommentLines, sanitizeComment(tk))
		}
	}

	return en, nil
}

func skipEndOfLineComments(tr *tokenReader) {
	for tr.Next() {
		nextTk := tr.Token()
		// comments at the end of lines after fields are -not- field comments for the next field
		if nextTk.kind == tokenKindLineComment {
			break
		}
		if nextTk.kind == tokenKindBlockComment {
			// there could be multiple block comments here
			continue
		}
		tr.UnNext()
		break
	}
}
func readStruct(tr *tokenReader) (Struct, error) {
	st := Struct{}
	toks, err := expectNext(tr, tokenKindIdent, tokenKindOpenCurly)
	if err != nil {
		return st, err
	}
	st.Name = string(toks[0].concrete)

	optNewline(tr)

	nextCommentLines := []string{}
	nextDecorations := Decorations{Custom: map[string]string{}}

	for tr.Token().kind != tokenKindCloseCurly {
		if !tr.Next() {
			return st, readError(tr.nextToken, "struct definition ended early")
		}
		tk := tr.Token()
		switch tk.kind {
		case tokenKindNewline:
			nextCommentLines = []string{}
		case tokenKindIdent, tokenKindArray, tokenKindMap:
			tr.UnNext()
			fdType, err := readFieldType(tr)
			if err != nil {
				return st, err
			}
			toks, err := expectNext(tr, tokenKindIdent, tokenKindSemicolon)
			if err != nil {
				return st, err
			}
			fdName := string(toks[0].concrete)
			st.Fields = append(st.Fields, Field{
				Name:        fdName,
				FieldType:   fdType,
				Decorations: nextDecorations,
				Comment:     strings.Join(nextCommentLines, "\n"),
			})
			nextDecorations = Decorations{Custom: map[string]string{}}

			nextCommentLines = []string{}

			skipEndOfLineComments(tr)
		case tokenKindOpenSquare, tokenKindAtSign:
			k, v, _, err := readDecoration(tr, tk.kind == tokenKindOpenSquare)
			if err != nil {
				return st, err
			}
			switch k {
			case "opcode":
				return st, readError(tr.nextToken, "opcode annotation not allowed within struct")
			case "flags":
				return st, readError(tr.nextToken, "flags annotation not allowed within struct")
			case "deprecated":
				if nextDecorations.Deprecated {
					return st, readError(tk, "deprecated cannot be applied to the same target twice")
				}
				nextDecorations.Deprecated = true
				nextDecorations.DeprecatedMessage = v
			default:
				nextDecorations.Custom[k] = v
			}
		case tokenKindBlockComment:
			nextCommentLines = append(nextCommentLines, readBlockComment(tr, tk))
		case tokenKindLineComment:
			nextCommentLines = append(nextCommentLines, sanitizeComment(tk))
		}
	}

	return st, nil
}

func readFieldType(tr *tokenReader) (FieldType, error) {
	ft := FieldType{}
	err := expectAnyOfNext(tr, tokenKindIdent, tokenKindArray, tokenKindMap)
	if err != nil {
		return ft, err
	}
	tk := tr.Token()
	switch tk.kind {
	case tokenKindMap:
		if _, err := expectNext(tr, tokenKindOpenSquare); err != nil {
			return ft, err
		}
		keyType, err := readFieldType(tr)
		if err != nil {
			return ft, err
		}
		if keyType.Map != nil || keyType.Array != nil {
			return ft, readError(tk, "map must begin with simple type")
		}
		if !isPrimitiveType(keyType.Simple) {
			return ft, readError(tk, "map must begin with simple type")
		}
		if _, err := expectNext(tr, tokenKindComma); err != nil {
			return ft, err
		}
		valType, err := readFieldType(tr)
		if err != nil {
			return ft, err
		}
		if _, err := expectNext(tr, tokenKindCloseSquare); err != nil {
			return ft, err
		}
		ft.Map = &MapType{
			Key:   keyType.Simple,
			Value: valType,
		}
	case tokenKindArray:
		if _, err := expectNext(tr, tokenKindOpenSquare); err != nil {
			return ft, err
		}
		arType, err := readFieldType(tr)
		if err != nil {
			return ft, err
		}
		if _, err := expectNext(tr, tokenKindCloseSquare); err != nil {
			return ft, err
		}
		ft.Array = &arType
	case tokenKindIdent:
		ft.Simple = string(tk.concrete)
	}
	if tr.Next() {
		// this might have been followed by []
		nextTk := tr.Token()
		for nextTk.kind == tokenKindOpenSquare {
			if _, err := expectNext(tr, tokenKindCloseSquare); err != nil {
				return ft, err
			}
			ft3 := ft
			ft = FieldType{
				Array: &ft3,
			}
			if !tr.Next() {
				return ft, nil
			}
			nextTk = tr.Token()
		}
		tr.UnNext()
	}
	return ft, nil
}

func readMessage(tr *tokenReader) (Message, error) {
	msg := Message{
		Fields: make(map[uint8]Field),
	}
	toks, err := expectNext(tr, tokenKindIdent, tokenKindOpenCurly)
	if err != nil {
		return msg, err
	}
	msg.Name = string(toks[0].concrete)

	optNewline(tr)

	nextCommentLines := []string{}
	nextDecorations := Decorations{Custom: map[string]string{}}

	for tr.Token().kind != tokenKindCloseCurly {
		if err := expectAnyOfNext(tr,
			tokenKindNewline,
			tokenKindIntegerLiteral,
			tokenKindOpenSquare,
			tokenKindBlockComment,
			tokenKindLineComment,
			tokenKindAtSign,
			tokenKindCloseCurly); err != nil {
			return msg, err
		}
		tk := tr.Token()
		switch tk.kind {
		case tokenKindCloseCurly:
			// break
		case tokenKindNewline:
			nextCommentLines = []string{}
		case tokenKindIntegerLiteral:
			fdInteger, err := strconv.ParseUint(string(tr.Token().concrete), 10, 8)
			if err != nil {
				return msg, readError(tr.nextToken, err.Error())
			}
			if _, ok := msg.Fields[uint8(fdInteger)]; ok {
				return msg, readError(tr.nextToken, "message has duplicate field index %d", fdInteger)
			}
			if _, err := expectNext(tr, tokenKindArrow); err != nil {
				return msg, err
			}
			fdType, err := readFieldType(tr)
			if err != nil {
				return msg, err
			}
			toks, err := expectNext(tr, tokenKindIdent, tokenKindSemicolon)
			if err != nil {
				return msg, err
			}
			fdName := string(toks[0].concrete)

			msg.Fields[uint8(fdInteger)] = Field{
				Name:        fdName,
				FieldType:   fdType,
				Decorations: nextDecorations,
				Comment:     strings.Join(nextCommentLines, "\n"),
			}
			nextDecorations = Decorations{Custom: map[string]string{}}

			nextCommentLines = []string{}

			skipEndOfLineComments(tr)
		case tokenKindOpenSquare, tokenKindAtSign:
			k, v, _, err := readDecoration(tr, tk.kind == tokenKindOpenSquare)
			if err != nil {
				return msg, err
			}
			switch k {
			case "opcode":
				return msg, readError(tr.nextToken, "opcode annotation not allowed within message")
			case "flags":
				return msg, readError(tr.nextToken, "flags annotation not allowed within message")
			case "deprecated":
				if nextDecorations.Deprecated {
					return msg, readError(tk, "deprecated cannot be applied to the same target twice")
				}
				nextDecorations.Deprecated = true
				nextDecorations.DeprecatedMessage = v
			default:
				nextDecorations.Custom[k] = v
			}
		case tokenKindBlockComment:
			nextCommentLines = append(nextCommentLines, readBlockComment(tr, tk))
		case tokenKindLineComment:
			nextCommentLines = append(nextCommentLines, sanitizeComment(tk))
		}
	}

	return msg, nil
}

func readUnion(tr *tokenReader) (Union, error) {
	union := Union{
		Fields: make(map[uint8]UnionField),
	}
	toks, err := expectNext(tr, tokenKindIdent, tokenKindOpenCurly)
	if err != nil {
		return union, err
	}
	union.Name = string(toks[0].concrete)

	optNewline(tr)

	nextCommentLines := []string{}
	nextDecorations := Decorations{Custom: map[string]string{}}

	for tr.Token().kind != tokenKindCloseCurly {
		if !tr.Next() {
			return union, readError(tr.nextToken, "union definition ended early")
		}
		tk := tr.Token()
		switch tk.kind {
		case tokenKindNewline:
			nextCommentLines = []string{}
		case tokenKindIntegerLiteral:
			fdInteger, err := strconv.ParseUint(string(tr.Token().concrete), 10, 8)
			if err != nil {
				return union, readError(tr.nextToken, err.Error())
			}
			if _, ok := union.Fields[uint8(fdInteger)]; ok {
				return union, readError(tr.nextToken, "union has duplicate field index %d", fdInteger)
			}
			if _, err := expectNext(tr, tokenKindArrow); err != nil {
				return union, err
			}
			if err := expectAnyOfNext(tr, tokenKindMessage, tokenKindStruct); err != nil {
				return union, readError(tr.nextToken, "union fields must be messages or structs")
			}
			unionFd := UnionField{}
			tk := tr.Token()
			switch tk.kind {
			case tokenKindMessage:
				msg, err := readMessage(tr)
				if err != nil {
					return union, err
				}
				msg.Comment = strings.Join(nextCommentLines, "\n")
				unionFd.Message = &msg
			case tokenKindStruct:
				st, err := readStruct(tr)
				if err != nil {
					return union, err
				}
				st.Comment = strings.Join(nextCommentLines, "\n")
				unionFd.Struct = &st
			}

			unionFd.Decorations = nextDecorations

			union.Fields[uint8(fdInteger)] = unionFd
			nextDecorations = Decorations{Custom: map[string]string{}}

			nextCommentLines = []string{}

			// This is a close curly-- we must advance past it or the union
			// will read it and believe it is complete
			tr.Next()
			skipEndOfLineComments(tr)
			optNewline(tr)
		case tokenKindOpenSquare, tokenKindAtSign:
			k, v, _, err := readDecoration(tr, tk.kind == tokenKindOpenSquare)
			if err != nil {
				return union, err
			}
			switch k {
			case "opcode":
				return union, readError(tr.nextToken, "opcode annotation not allowed within union")
			case "flags":
				return union, readError(tr.nextToken, "flags annotation not allowed within union")
			case "deprecated":
				if nextDecorations.Deprecated {
					return union, readError(tk, "deprecated cannot be applied to the same target twice")
				}
				nextDecorations.Deprecated = true
				nextDecorations.DeprecatedMessage = v
			default:
				nextDecorations.Custom[k] = v
			}
		case tokenKindBlockComment:
			nextCommentLines = append(nextCommentLines, readBlockComment(tr, tk))
		case tokenKindLineComment:
			nextCommentLines = append(nextCommentLines, sanitizeComment(tk))
		}
	}
	return union, nil
}

func readConst(tr *tokenReader) (Const, []string, error) {
	warnings := []string{}
	cons := Const{}
	toks, err := expectNext(tr, tokenKindIdent, tokenKindIdent, tokenKindEquals)
	if err != nil {
		return cons, warnings, err
	}
	cons.SimpleType = string(toks[0].concrete)
	cons.Name = string(toks[1].concrete)

	if !tr.Next() {
		return cons, warnings, readError(toks[2], "expected value following const type")
	}
	tk := tr.Token()
	cons.Value = string(tk.concrete)
	switch {
	case isUintPrimitive(cons.SimpleType):
		if tk.kind != tokenKindIntegerLiteral {
			return cons, warnings, readError(tk, "%v unassignable to %v", tk.kind, cons.SimpleType)
		}
		_, err := strconv.ParseUint(string(tk.concrete), 0, 64)
		if err != nil {
			// export a warning
			warnings = append(warnings, readError(tk, err.Error()).Error())
		}
	case isIntPrimitive(cons.SimpleType):
		if tk.kind != tokenKindIntegerLiteral {
			return cons, warnings, readError(tk, "%v unassignable to %v", tk.kind, cons.SimpleType)
		}
		_, err := strconv.ParseInt(string(tk.concrete), 0, 64)
		if err != nil {
			// export a warning
			warnings = append(warnings, readError(tk, err.Error()).Error())
		}
	case isFloatPrimitive(cons.SimpleType):
		switch tk.kind {
		case tokenKindInf:
			cons.Value = "math.Inf(1)"
		case tokenKindNegativeInf:
			cons.Value = "math.Inf(-1)"
		case tokenKindNaN:
			cons.Value = "math.NaN()"
		case tokenKindIntegerLiteral:
			_, err := strconv.ParseInt(string(tk.concrete), 0, 64)
			if err != nil {
				// export a warning
				warnings = append(warnings, readError(tk, err.Error()).Error())
			}
		case tokenKindFloatLiteral:
			_, err := strconv.ParseFloat(string(tk.concrete), 64)
			if err != nil {
				// export a warning
				warnings = append(warnings, readError(tk, err.Error()).Error())
			}
		default:
			return cons, warnings, readError(tk, "%v unassignable to %v", tk.kind, cons.SimpleType)
		}
	case cons.SimpleType == typeGUID:
		if tk.kind != tokenKindStringLiteral {
			return cons, warnings, readError(tk, "%v unassignable to %v", tk.kind, cons.SimpleType)
		}
		s := string(bytes.Trim(tk.concrete, "\""))
		// TODO: what guid formats does rainway support?
		if len(strings.ReplaceAll(s, "-", "")) != 32 {
			return cons, warnings, readError(tk, "%q has wrong length for guid", s)
		}
	case cons.SimpleType == typeString:
		if tk.kind != tokenKindStringLiteral {
			return cons, warnings, readError(tk, "%v unassignable to %v", tk.kind, cons.SimpleType)
		}
	case cons.SimpleType == typeBool:
		if tk.kind != tokenKindTrue && tk.kind != tokenKindFalse {
			return cons, warnings, readError(tk, "%v unassignable to %v", tk.kind, cons.SimpleType)
		}
	default:
		return cons, warnings, readError(tk, "invalid type %q for const", cons.SimpleType)
	}
	_, err = expectNext(tr, tokenKindSemicolon)
	if err != nil {
		return cons, warnings, err
	}
	skipEndOfLineComments(tr)
	optNewline(tr)
	return cons, warnings, nil
}

func readDecoration(tr *tokenReader, leadingSquare bool) (k, v string, opcodeV uint32, err error) {
	if err := expectAnyOfNext(tr, tokenKindOpCode, tokenKindDeprecated, tokenKindFlags, tokenKindIdent); err != nil {
		return "", "", 0, err
	}
	tk := tr.Token()
	switch tk.kind {
	case tokenKindOpCode:
		if _, err := expectNext(tr, tokenKindOpenParen); err != nil {
			return "", "", 0, err
		}
		if err := expectAnyOfNext(tr, tokenKindIntegerLiteral, tokenKindStringLiteral); err != nil {
			return "", "", 0, err
		}
		tk := tr.Token()
		if tk.kind == tokenKindIntegerLiteral {
			content := string(tk.concrete)
			opc, err := strconv.ParseUint(content, 0, 32)
			if err != nil {
				return "", "", 0, readError(tk, err.Error())
			}
			opcodeV = uint32(opc)
		} else if tk.kind == tokenKindStringLiteral {
			tk.concrete = bytes.Trim(tk.concrete, "\"")
			if len(tk.concrete) != 4 {
				return "", "", 0, readError(tk, "opcode string %q not 4 ascii characters", string(tk.concrete))
			}
			opcodeV = bytesToOpCode(*(*[4]byte)(tk.concrete))
		}
		if _, err := expectNext(tr, tokenKindCloseParen); err != nil {
			return "", "", 0, err
		}
		k = "opcode"
	case tokenKindDeprecated:
		fallthrough
	case tokenKindIdent:
		if _, err := expectNext(tr, tokenKindOpenParen); err == nil {
			tks, err := expectNext(tr, tokenKindStringLiteral, tokenKindCloseParen)
			if err != nil {
				return "", "", 0, err
			}
			// this cannot error; token readers cannot parse strings
			// with missing terminal quotes.
			v, _ = strconv.Unquote(string(tks[0].concrete))
		}
		k = string(tk.concrete)
	case tokenKindFlags:
		k = "flags"
	}
	if leadingSquare {
		if _, err := expectNext(tr, tokenKindCloseSquare); err != nil {
			return "", "", 0, err
		}
	}

	skipEndOfLineComments(tr)
	optNewline(tr)
	return
}

func readBlockComment(tr *tokenReader, tk token) string {
	return string(tk.concrete[2 : len(tk.concrete)-2])
}

func sanitizeComment(tk token) string {
	comment := string(tk.concrete[2:])
	comment = strings.Trim(comment, "\r\n")
	return comment
}

func bytesToOpCode(data [4]byte) uint32 {
	opCode := uint32(data[0])
	opCode |= (uint32(data[1]) << 8)
	opCode |= (uint32(data[2]) << 16)
	opCode |= (uint32(data[3]) << 24)
	return opCode
}

func readError(tk token, format string, args ...interface{}) error {
	format = fmt.Sprintf("[%d:%d] ", tk.loc.line, tk.loc.lineChar) + format
	return fmt.Errorf(format, args...)
}

func kindsStr(ks []tokenKind) string {
	kindsStrs := make([]string, len(ks))
	for i, k := range ks {
		kindsStrs[i] = k.String()
	}
	return strings.Join(kindsStrs, ", ")
}
