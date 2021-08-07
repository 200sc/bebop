package bebop

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// ReadFile reads out a bebop file.
func ReadFile(r io.Reader) (File, error) {
	f := File{}
	tr := newTokenReader(r)
	nextCommentLines := []string{}
	nextRecordOpCode := int32(0)
	nextRecordReadOnly := false
	for tr.Next() {
		tk := tr.Token()
		switch tk.kind {
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
			var err error
			nextRecordOpCode, err = readOpCode(tr)
			if err != nil {
				return f, err
			}
			continue
		case tokenKindEnum:
			if nextRecordOpCode != 0 {
				return f, readError(tk, "enums may not have attached op codes")
			}
			if nextRecordReadOnly {
				return f, readError(tk, "enums cannot be readonly")
			}
			en, err := readEnum(tr)
			if err != nil {
				return f, err
			}
			en.Comment = strings.Join(nextCommentLines, "\n")
			f.Enums = append(f.Enums, en)
		case tokenKindReadOnly:
			nextRecordReadOnly = true
			if !tr.Next() {
				return f, readError(tk, "expected (Struct) got no token")
			}
			tk = tr.Token()
			if tk.kind != tokenKindStruct {
				return f, readError(tk, "expected (Struct) got (%v)", tk.kind)
			}
			fallthrough
		case tokenKindStruct:
			st, err := readStruct(tr)
			if err != nil {
				return f, err
			}
			st.Comment = strings.Join(nextCommentLines, "\n")
			st.OpCode = nextRecordOpCode
			st.ReadOnly = nextRecordReadOnly
			f.Structs = append(f.Structs, st)
			nextRecordReadOnly = false
		case tokenKindMessage:
			if nextRecordReadOnly {
				return f, readError(tk, "messages cannot be readonly")
			}
			msg, err := readMessage(tr)
			if err != nil {
				return f, err
			}
			msg.Comment = strings.Join(nextCommentLines, "\n")
			msg.OpCode = nextRecordOpCode
			f.Messages = append(f.Messages, msg)
		case tokenKindUnion:
			if nextRecordReadOnly {
				return f, readError(tk, "unions cannot be readonly")
			}
			union, err := readUnion(tr)
			if err != nil {
				return f, err
			}
			union.Comment = strings.Join(nextCommentLines, "\n")
			union.OpCode = nextRecordOpCode
			f.Unions = append(f.Unions, union)
		}
		nextCommentLines = []string{}
		nextRecordOpCode = 0
	}
	return f, nil
}

func expectAnyOfNext(tr *tokenReader, kinds ...tokenKind) error {
	next := tr.Next()
	if tr.Err() != nil {
		return tr.Err()
	}
	if !next {
		return readError(tr.nextToken, "expected (%v), got no token", kinds)
	}
	tk := tr.Token()
	found := false
	for _, k := range kinds {
		if tk.kind == k {
			found = true
		}
	}
	if !found {
		return readError(tk, "expected (%v) got %s", kinds, tk.kind)
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
			return tokens, readError(tr.nextToken, "expected (%v), got no token", kinds)
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

func readEnum(tr *tokenReader) (Enum, error) {
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
			toks, err = expectNext(tr, tokenKindEquals, tokenKindInteger, tokenKindSemicolon)
			if err != nil {
				return en, err
			}
			optInteger, err := strconv.ParseInt(string(toks[1].concrete), 10, 32)
			if err != nil {
				return en, readError(toks[1], err.Error())
			}
			en.Options = append(en.Options, EnumOption{
				Name:              optName,
				Value:             int32(optInteger),
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
			})
			nextDeprecatedMessage = ""
			nextIsDeprecated = false
			nextCommentLines = []string{}

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
	nextDeprecatedMessage := ""
	nextIsDeprecated := false
	for tr.Token().kind != tokenKindCloseCurly {
		if !tr.Next() {
			return msg, readError(tr.nextToken, "message definition ended early")
		}
		tk := tr.Token()
		switch tk.kind {
		case tokenKindNewline:
			nextCommentLines = []string{}
		case tokenKindInteger:
			fdInteger, err := strconv.ParseInt(string(tr.Token().concrete), 10, 8)
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
			}
			nextDeprecatedMessage = ""
			nextIsDeprecated = false
			nextCommentLines = []string{}

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
		case tokenKindInteger:
			fdInteger, err := strconv.ParseInt(string(tr.Token().concrete), 10, 8)
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

			unionFd.Deprecated = nextIsDeprecated
			unionFd.DeprecatedMessage = nextDeprecatedMessage

			union.Fields[uint8(fdInteger)] = unionFd
			nextDeprecatedMessage = ""
			nextIsDeprecated = false
			nextCommentLines = []string{}

			skipEndOfLineComments(tr)
			tr.Next()
			tr.Next()

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
			nextCommentLines = append(nextCommentLines, sanitizeComment(tk))
		}
	}

	return union, nil
}

func readOpCode(tr *tokenReader) (int32, error) {
	if _, err := expectNext(tr, tokenKindOpCode, tokenKindOpenParen); err != nil {
		return 0, err
	}
	if err := expectAnyOfNext(tr, tokenKindInteger, tokenKindStringLiteral); err != nil {
		return 0, err
	}
	var opCode int32
	tk := tr.Token()
	if tk.kind == tokenKindInteger {
		content := string(tk.concrete)
		opc, err := strconv.ParseInt(content, 0, 32)
		if err != nil {
			return 0, readError(tk, err.Error())
		}
		opCode = int32(opc)
	} else if tk.kind == tokenKindStringLiteral {
		tk.concrete = bytes.Trim(tk.concrete, "\"")
		if len(tk.concrete) > 4 {
			return 0, readError(tk, "opcode string %s exceeds 4 ascii characters", string(tk.concrete))
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
