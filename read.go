package bebop

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

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
		// top level:
		case tokenKindIdent:
			switch string(tk.concrete) {
			case "readonly":
				nextRecordReadOnly = true
				continue
			case "enum":
				if nextRecordOpCode != 0 {
					return f, fmt.Errorf("enums may not have attached op codes")
				}
				en, err := readEnum(tr)
				if err != nil {
					return f, err
				}
				en.Comment = strings.Join(nextCommentLines, "\n")
				f.Enums = append(f.Enums, en)

				nextCommentLines = []string{}
				nextRecordOpCode = 0
			case "struct":
				st, err := readStruct(tr)
				if err != nil {
					return f, err
				}
				st.Comment = strings.Join(nextCommentLines, "\n")
				st.OpCode = nextRecordOpCode
				st.ReadOnly = nextRecordReadOnly
				f.Structs = append(f.Structs, st)

				nextCommentLines = []string{}
				nextRecordOpCode = 0
				nextRecordReadOnly = false
			case "message":
				msg, err := readMessage(tr)
				if err != nil {
					return f, err
				}
				msg.Comment = strings.Join(nextCommentLines, "\n")
				msg.OpCode = nextRecordOpCode
				msg.ReadOnly = nextRecordReadOnly
				f.Messages = append(f.Messages, msg)

				nextCommentLines = []string{}
				nextRecordOpCode = 0
				nextRecordReadOnly = false
			}
		case tokenKindBlockComment:
			nextCommentLines = append(nextCommentLines, string(tk.concrete[2:len(tk.concrete)-2]))
			if tr.Next() {
				if tr.Token().kind != tokenKindNewline {
					tr.UnNext()
				}
			}
		case tokenKindLineComment:
			nextComment := string(tk.concrete[2:])
			nextComment = strings.Trim(nextComment, "\r\n")
			nextCommentLines = append(nextCommentLines, nextComment)
		case tokenKindOpenSquare:
			if err := expectNext(tr, tokenKindIdent); err != nil {
				return f, err
			}
			tk = tr.Token()
			if string(tk.concrete) != "opcode" {
				return f, fmt.Errorf("invalid ident after leading '[', got %s, wanted %v", string(tk.concrete), "opcode")
			}
			if err := expectNext(tr, tokenKindOpenParen); err != nil {
				return f, err
			}
			if err := expectNext(tr, tokenKindInteger, tokenKindStringLiteral); err != nil {
				return f, err
			}
			tk = tr.Token()
			if tk.kind == tokenKindInteger {
				content := string(tk.concrete)
				opCode, err := strconv.ParseInt(content, 0, 32)
				if err != nil {
					return f, err
				}
				nextRecordOpCode = int32(opCode)
			} else if tk.kind == tokenKindStringLiteral {
				tk.concrete = bytes.Trim(tk.concrete, "\"")
				if len(tk.concrete) > 4 {
					return f, fmt.Errorf("opcode string %s exceeds 4 ascii characters", string(tk.concrete))
				}
				nextRecordOpCode = bytesToOpCode(tk.concrete)

			}
			if err := expectNext(tr, tokenKindCloseParen); err != nil {
				return f, err
			}
			if err := expectNext(tr, tokenKindCloseSquare); err != nil {
				return f, err
			}
		}
	}
	return f, nil
}

func expectNext(tr *tokenReader, kinds ...tokenKind) error {
	hasNext := tr.Next()
	if tr.Err() != nil {
		return tr.Err()
	}
	if !hasNext {
		return fmt.Errorf("expected (%v), got no token", kinds)
	}
	tk := tr.Token()
	found := false
	for _, k := range kinds {
		if tk.kind == k {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("expected (%v) got %s", kinds, tk.kind)
	}
	return nil
}

func readEnum(tr *tokenReader) (Enum, error) {
	en := Enum{}
	if err := expectNext(tr, tokenKindIdent); err != nil {
		return en, err
	}
	en.Name = string(tr.Token().concrete)
	if err := expectNext(tr, tokenKindOpenCurly); err != nil {
		return en, err
	}
	tr.Next()
	// optional newline
	if tr.Token().kind != tokenKindNewline {
		tr.UnNext()
	}
	nextCommentLines := []string{}
	nextDeprecatedMessage := ""
	nextIsDeprecated := false
	for tr.Token().kind != tokenKindCloseCurly {
		if !tr.Next() {
			return en, fmt.Errorf("enum definition ended early")
		}
		tk := tr.Token()
		switch tk.kind {
		case tokenKindNewline:
			nextCommentLines = []string{}
		case tokenKindIdent:
			optName := string(tk.concrete)
			if err := expectNext(tr, tokenKindEquals); err != nil {
				return en, err
			}
			if err := expectNext(tr, tokenKindInteger); err != nil {
				return en, err
			}
			optInteger, err := strconv.ParseInt(string(tr.Token().concrete), 10, 32)
			if err != nil {
				return en, err
			}
			if err := expectNext(tr, tokenKindSemicolon); err != nil {
				return en, err
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
				return en, fmt.Errorf("expected enum option following deprecated annotation")
			}
			msg, err := readDeprecated(tr)
			if err != nil {
				return en, err
			}
			nextIsDeprecated = true
			nextDeprecatedMessage = msg
		case tokenKindBlockComment:
			nextCommentLines = append(nextCommentLines, string(tk.concrete[2:len(tk.concrete)-2]))
			if tr.Next() {
				if tr.Token().kind != tokenKindNewline {
					tr.UnNext()
				}
			}
		case tokenKindLineComment:
			nextComment := string(tk.concrete[2:])
			nextComment = strings.Trim(nextComment, "\r\n")
			nextCommentLines = append(nextCommentLines, nextComment)
		}
	}

	return en, nil
}

func readDeprecated(tr *tokenReader) (string, error) {
	if err := expectNext(tr, tokenKindIdent); err != nil {
		return "", err
	}
	tk := tr.Token()
	if string(tk.concrete) != "deprecated" {
		return "", fmt.Errorf("invalid ident after leading '[', got %s, wanted %v", string(tk.concrete), "deprecated")
	}
	if err := expectNext(tr, tokenKindOpenParen); err != nil {
		return "", err
	}
	if err := expectNext(tr, tokenKindStringLiteral); err != nil {
		return "", err
	}
	tk = tr.Token()
	var err error
	msg, err := strconv.Unquote(string(tk.concrete))
	if err != nil {
		return "", err
	}
	if err := expectNext(tr, tokenKindCloseParen); err != nil {
		return "", err
	}
	if err := expectNext(tr, tokenKindCloseSquare); err != nil {
		return "", err
	}
	if err := expectNext(tr, tokenKindNewline); err != nil {
		return "", err
	}
	return msg, nil
}

func readStruct(tr *tokenReader) (Struct, error) {
	st := Struct{}
	if err := expectNext(tr, tokenKindIdent); err != nil {
		return st, err
	}
	st.Name = string(tr.Token().concrete)
	if err := expectNext(tr, tokenKindOpenCurly); err != nil {
		return st, err
	}
	tr.Next()
	// optional newline
	if tr.Token().kind != tokenKindNewline {
		tr.UnNext()
	}

	nextCommentLines := []string{}
	nextDeprecatedMessage := ""
	nextIsDeprecated := false
	for tr.Token().kind != tokenKindCloseCurly {
		if !tr.Next() {
			return st, fmt.Errorf("struct definition ended early")
		}
		tk := tr.Token()
		switch tk.kind {
		case tokenKindNewline:
			nextCommentLines = []string{}
		case tokenKindIdent:
			tr.UnNext()
			fdType, err := readFieldType(tr)
			if err != nil {
				return st, err
			}
			if err := expectNext(tr, tokenKindIdent); err != nil {
				return st, err
			}
			fdName := string(tr.Token().concrete)
			if err := expectNext(tr, tokenKindSemicolon); err != nil {
				return st, err
			}
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

		case tokenKindOpenSquare:
			if nextIsDeprecated {
				return st, fmt.Errorf("expected field following deprecated annotation")
			}
			msg, err := readDeprecated(tr)
			if err != nil {
				return st, err
			}
			nextIsDeprecated = true
			nextDeprecatedMessage = msg
		case tokenKindBlockComment:
			nextCommentLines = append(nextCommentLines, string(tk.concrete[2:len(tk.concrete)-2]))
			if tr.Next() {
				if tr.Token().kind != tokenKindNewline {
					tr.UnNext()
				}
			}
		case tokenKindLineComment:
			nextComment := string(tk.concrete[2:])
			nextComment = strings.Trim(nextComment, "\r\n")
			nextCommentLines = append(nextCommentLines, nextComment)
		}
	}

	return st, nil
}

func readFieldType(tr *tokenReader) (FieldType, error) {
	ft := FieldType{}
	if err := expectNext(tr, tokenKindIdent); err != nil {
		return ft, err
	}
	tk := tr.Token()
	switch string(tk.concrete) {
	case "map":
		if err := expectNext(tr, tokenKindOpenSquare); err != nil {
			return ft, err
		}
		keyType, err := readFieldType(tr)
		if err != nil {
			return ft, err
		}
		if keyType.IsMap() || keyType.IsArray() {
			return ft, fmt.Errorf("map must begin with simple type")
		}
		if !isPrimitiveType(keyType.Simple) {
			return ft, fmt.Errorf("map must being with simple type")
		}
		if err := expectNext(tr, tokenKindComma); err != nil {
			return ft, err
		}
		valType, err := readFieldType(tr)
		if err != nil {
			return ft, err
		}
		if err := expectNext(tr, tokenKindCloseSquare); err != nil {
			return ft, err
		}
		ft.Map = &MapType{
			Key:   keyType.Simple,
			Value: valType,
		}
	case "array":
		if err := expectNext(tr, tokenKindOpenSquare); err != nil {
			return ft, err
		}
		arType, err := readFieldType(tr)
		if err != nil {
			return ft, err
		}
		if err := expectNext(tr, tokenKindCloseSquare); err != nil {
			return ft, err
		}
		ft.Array = &arType
	default:
		ft.Simple = string(tk.concrete)
	}
	if tr.Next() {
		// this might have been followed by []
		nextTk := tr.Token()
		if nextTk.kind == tokenKindOpenSquare {
			if err := expectNext(tr, tokenKindCloseSquare); err != nil {
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
	if err := expectNext(tr, tokenKindIdent); err != nil {
		return msg, err
	}
	msg.Name = string(tr.Token().concrete)
	if err := expectNext(tr, tokenKindOpenCurly); err != nil {
		return msg, err
	}
	tr.Next()
	// optional newline
	if tr.Token().kind != tokenKindNewline {
		tr.UnNext()
	}

	nextCommentLines := []string{}
	nextDeprecatedMessage := ""
	nextIsDeprecated := false
	for tr.Token().kind != tokenKindCloseCurly {
		if !tr.Next() {
			return msg, fmt.Errorf("message definition ended early")
		}
		tk := tr.Token()
		switch tk.kind {
		case tokenKindNewline:
			nextCommentLines = []string{}
		case tokenKindInteger:
			fdInteger, err := strconv.ParseInt(string(tr.Token().concrete), 10, 8)
			if err != nil {
				return msg, err
			}
			if err := expectNext(tr, tokenKindArrow); err != nil {
				return msg, err
			}

			fdType, err := readFieldType(tr)
			if err != nil {
				return msg, err
			}
			if err := expectNext(tr, tokenKindIdent); err != nil {
				return msg, err
			}
			fdName := string(tr.Token().concrete)
			if err := expectNext(tr, tokenKindSemicolon); err != nil {
				return msg, err
			}
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

		case tokenKindOpenSquare:
			if nextIsDeprecated {
				return msg, fmt.Errorf("expected field following deprecated annotation")
			}
			dpMsg, err := readDeprecated(tr)
			if err != nil {
				return msg, err
			}
			nextIsDeprecated = true
			nextDeprecatedMessage = dpMsg
		case tokenKindBlockComment:
			nextCommentLines = append(nextCommentLines, string(tk.concrete[2:len(tk.concrete)-2]))
			if tr.Next() {
				if tr.Token().kind != tokenKindNewline {
					tr.UnNext()
				}
			}
		case tokenKindLineComment:
			nextComment := string(tk.concrete[2:])
			nextComment = strings.Trim(nextComment, "\r\n")
			nextCommentLines = append(nextCommentLines, nextComment)
		}
	}

	return msg, nil
}

func isPrimitiveType(simpleType string) bool {
	_, ok := primitiveTypes[simpleType]
	return ok
}

var primitiveTypes = map[string]struct{}{
	"bool":    {},
	"byte":    {},
	"uint8":   {},
	"uint16":  {},
	"int16":   {},
	"uint32":  {},
	"int32":   {},
	"uint64":  {},
	"int64":   {},
	"float32": {},
	"float64": {},
	"string":  {},
	"guid":    {},
	"date":    {},
}

func bytesToOpCode(data []byte) int32 {
	opCode := int32(0)
	for _, b := range data {
		opCode <<= 8
		opCode |= int32(b)
	}
	return opCode
}
