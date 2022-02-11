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
	nextRecordOpCode := int32(0)
	nextRecordReadOnly := false
	nextRecordBitFlags := false
	warnings := []string{}
	for tr.Next() {
		tk := tr.Token()
		switch tk.kind {
		case tokenKindImport:
			toks, err := expectNext(tr, tokenKindStringLiteral)
			if err != nil {
				return f, warnings, err
			}
			imported, err := strconv.Unquote(string(toks[0].concrete))
			if err != nil {
				return f, warnings, err
			}
			f.Imports = append(f.Imports, imported)
			continue
		case tokenKindNewline:
			nextCommentLines = []string{}
			continue
		case tokenKindBlockComment:
			nextCommentLines = append(nextCommentLines, readBlockComment(tr, tk))
			continue
		case tokenKindLineComment:
			nextCommentLines = append(nextCommentLines, sanitizeComment(tk))
			continue
		case tokenKindOpenSquare:
			if err := expectAnyOfNext(tr, tokenKindOpCode, tokenKindFlags); err != nil {
				return f, warnings, err
			}
			switch tr.Token().kind {
			case tokenKindOpCode:
				tr.UnNext()
				var err error
				nextRecordOpCode, err = readOpCode(tr)
				if err != nil {
					return f, warnings, err
				}
			case tokenKindFlags:
				nextRecordBitFlags = true
				if err := expectAnyOfNext(tr, tokenKindCloseSquare); err != nil {
					return f, warnings, err
				}
			}
			continue
		case tokenKindEnum:
			if nextRecordOpCode != 0 {
				return f, warnings, readError(tk, "enums may not have attached op codes")
			}
			en, err := readEnum(tr, nextRecordBitFlags)
			if err != nil {
				return f, warnings, err
			}
			en.Comment = strings.Join(nextCommentLines, "\n")
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
			f.Structs = append(f.Structs, st)
			nextRecordReadOnly = false
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
			f.Consts = append(f.Consts, cons)
		}
		nextCommentLines = []string{}
		nextRecordOpCode = 0
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

func readEnumOptionValue(tr *tokenReader, previousOptions []EnumOption, bitflags bool) (int32, error) {
	if _, err := expectNext(tr, tokenKindEquals); err != nil {
		return 0, err
	}
	if !bitflags {
		toks, err := expectNext(tr, tokenKindIntegerLiteral, tokenKindSemicolon)
		if err != nil {
			return 0, err
		}
		optInteger, err := strconv.ParseInt(string(toks[0].concrete), 0, 32)
		if err != nil {
			return 0, err
		}
		return int32(optInteger), nil
	}
	return readBitflagExpr(tr, previousOptions)
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

func readEnum(tr *tokenReader, bitflags bool) (Enum, error) {
	en := Enum{}
	toks, err := expectNext(tr, tokenKindIdent, tokenKindOpenCurly)
	if err != nil {
		return en, err
	}
	en.Name = string(toks[0].concrete)

	optNewline(tr)

	nextCommentLines := []string{}
	nextDeprecatedMessage := ""
	nextIsDeprecated := false
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

			optValue, err := readEnumOptionValue(tr, en.Options, bitflags)
			if err != nil {
				return en, err
			}
			en.Options = append(en.Options, EnumOption{
				Name:              optName,
				Value:             optValue,
				DeprecatedMessage: nextDeprecatedMessage,
				Deprecated:        nextIsDeprecated,
				Comment:           strings.Join(nextCommentLines, "\n"),
			})
			nextDeprecatedMessage = ""
			nextIsDeprecated = false
			nextCommentLines = []string{}

		case tokenKindOpenSquare:
			if nextIsDeprecated {
				return en, readError(tk, "expected enum option following deprecated annotation")
			}
			msg, err := readDeprecated(tr)
			if err != nil {
				return en, err
			}
			nextIsDeprecated = true
			nextDeprecatedMessage = msg
		case tokenKindBlockComment:
			nextCommentLines = append(nextCommentLines, readBlockComment(tr, tk))
		case tokenKindLineComment:
			nextCommentLines = append(nextCommentLines, sanitizeComment(tk))
		}
	}

	return en, nil
}

func readDeprecated(tr *tokenReader) (string, error) {
	// TODO: can deprecated / op code be followed by a semicolon?
	toks, err := expectNext(tr, tokenKindDeprecated, tokenKindOpenParen, tokenKindStringLiteral,
		tokenKindCloseParen, tokenKindCloseSquare)
	if err != nil {
		return "", err
	}
	msg, err := strconv.Unquote(string(toks[2].concrete))
	if err != nil {
		return "", err
	}
	optNewline(tr)
	return msg, nil
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
	nextCommentTags := []Tag{}
	nextDeprecatedMessage := ""
	nextIsDeprecated := false
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
				Name:              fdName,
				FieldType:         fdType,
				DeprecatedMessage: nextDeprecatedMessage,
				Deprecated:        nextIsDeprecated,
				Comment:           strings.Join(nextCommentLines, "\n"),
				Tags:              nextCommentTags,
			})
			nextDeprecatedMessage = ""
			nextIsDeprecated = false
			nextCommentLines = []string{}
			nextCommentTags = []Tag{}

			skipEndOfLineComments(tr)
		case tokenKindOpenSquare:
			if nextIsDeprecated {
				return st, readError(tk, "expected field following deprecated annotation")
			}
			msg, err := readDeprecated(tr)
			if err != nil {
				return st, err
			}
			nextIsDeprecated = true
			nextDeprecatedMessage = msg
		case tokenKindBlockComment:
			nextCommentLines = append(nextCommentLines, readBlockComment(tr, tk))
		case tokenKindLineComment:
			cmt := sanitizeComment(tk)
			if tag, ok := parseCommentTag(cmt); ok {
				nextCommentTags = append(nextCommentTags, tag)
			}
			nextCommentLines = append(nextCommentLines, cmt)
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
		if nextTk.kind == tokenKindOpenSquare {
			if _, err := expectNext(tr, tokenKindCloseSquare); err != nil {
				return ft, err
			}
			return FieldType{
				Array: &ft,
			}, nil
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
	nextCommentTags := []Tag{}
	nextDeprecatedMessage := ""
	nextIsDeprecated := false
	for tr.Token().kind != tokenKindCloseCurly {
		if err := expectAnyOfNext(tr,
			tokenKindNewline,
			tokenKindIntegerLiteral,
			tokenKindOpenSquare,
			tokenKindBlockComment,
			tokenKindLineComment,
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
				Name:              fdName,
				FieldType:         fdType,
				DeprecatedMessage: nextDeprecatedMessage,
				Deprecated:        nextIsDeprecated,
				Comment:           strings.Join(nextCommentLines, "\n"),
				Tags:              nextCommentTags,
			}
			nextDeprecatedMessage = ""
			nextIsDeprecated = false
			nextCommentLines = []string{}
			nextCommentTags = []Tag{}

			skipEndOfLineComments(tr)
		case tokenKindOpenSquare:
			if nextIsDeprecated {
				return msg, readError(tk, "expected field following deprecated annotation")
			}
			dpMsg, err := readDeprecated(tr)
			if err != nil {
				return msg, err
			}
			nextIsDeprecated = true
			nextDeprecatedMessage = dpMsg
		case tokenKindBlockComment:
			nextCommentLines = append(nextCommentLines, readBlockComment(tr, tk))
		case tokenKindLineComment:
			cmt := sanitizeComment(tk)
			if tag, ok := parseCommentTag(cmt); ok {
				nextCommentTags = append(nextCommentTags, tag)
			}
			nextCommentLines = append(nextCommentLines, cmt)
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
	nextCommentTags := []Tag{}
	nextDeprecatedMessage := ""
	nextIsDeprecated := false
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

			unionFd.Tags = nextCommentTags
			unionFd.Deprecated = nextIsDeprecated
			unionFd.DeprecatedMessage = nextDeprecatedMessage

			union.Fields[uint8(fdInteger)] = unionFd
			nextDeprecatedMessage = ""
			nextIsDeprecated = false
			nextCommentLines = []string{}
			nextCommentTags = []Tag{}

			// This is a close curly-- we must advance past it or the union
			// will read it and believe it is complete
			tr.Next()
			skipEndOfLineComments(tr)
			optNewline(tr)

		case tokenKindOpenSquare:
			if nextIsDeprecated {
				return union, readError(tk, "expected field following deprecated annotation")
			}
			dpMsg, err := readDeprecated(tr)
			if err != nil {
				return union, err
			}
			nextIsDeprecated = true
			nextDeprecatedMessage = dpMsg
		case tokenKindBlockComment:
			nextCommentLines = append(nextCommentLines, readBlockComment(tr, tk))
		case tokenKindLineComment:
			cmt := sanitizeComment(tk)
			if tag, ok := parseCommentTag(cmt); ok {
				nextCommentTags = append(nextCommentTags, tag)
			}
			nextCommentLines = append(nextCommentLines, cmt)
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

func readOpCode(tr *tokenReader) (int32, error) {
	if _, err := expectNext(tr, tokenKindOpCode, tokenKindOpenParen); err != nil {
		return 0, err
	}
	if err := expectAnyOfNext(tr, tokenKindIntegerLiteral, tokenKindStringLiteral); err != nil {
		return 0, err
	}
	var opCode int32
	tk := tr.Token()
	if tk.kind == tokenKindIntegerLiteral {
		content := string(tk.concrete)
		opc, err := strconv.ParseInt(content, 0, 32)
		if err != nil {
			return 0, readError(tk, err.Error())
		}
		opCode = int32(opc)
	} else if tk.kind == tokenKindStringLiteral {
		tk.concrete = bytes.Trim(tk.concrete, "\"")
		if len(tk.concrete) > 4 {
			return 0, readError(tk, "opcode string %q exceeds 4 ascii characters", string(tk.concrete))
		}
		opCode = bytesToOpCode(tk.concrete)
	}
	if _, err := expectNext(tr, tokenKindCloseParen, tokenKindCloseSquare); err != nil {
		return 0, err
	}

	optNewline(tr)

	return opCode, nil
}

func readBlockComment(tr *tokenReader, tk token) string {
	return string(tk.concrete[2 : len(tk.concrete)-2])
}

func sanitizeComment(tk token) string {
	comment := string(tk.concrete[2:])
	comment = strings.Trim(comment, "\r\n")
	return comment
}

func bytesToOpCode(data []byte) int32 {
	opCode := int32(0)
	for _, b := range data {
		opCode <<= 8
		opCode |= int32(b)
	}
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

func parseCommentTag(s string) (Tag, bool) {
	// OK
	//[tag(json:"example,omitempty")]
	//[tag(json:"more colons::")]
	//[tag(boolean)]
	// Not OK
	// [tag(db:unquotedstring)]
	// [tag()]

	if !strings.HasPrefix(s, "[tag(") || !strings.HasSuffix(s, ")]") {
		return Tag{}, false
	}
	s = strings.TrimPrefix(s, "[tag(")
	s = strings.TrimSuffix(s, ")]")
	split := strings.Split(s, ":")
	if len(split) == 0 {
		return Tag{}, false
	}
	if len(split) == 1 {
		return Tag{
			Key:     split[0],
			Boolean: true,
		}, true
	}
	var err error
	value := strings.Join(split[1:], "")
	value, err = strconv.Unquote(value)
	if err != nil {
		return Tag{}, false
	}
	return Tag{
		Key:   split[0],
		Value: value,
	}, true
}
