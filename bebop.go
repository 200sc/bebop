package bebop

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type File struct {
	Structs  []Struct
	Messages []Message
	Enums    []Enum
}

func ReadFile(r io.Reader) (File, error) {
	f := File{}
	tr := newTokenReader(r)
	nextRecordComment := ""
	nextRecordOpCode := int32(0)
	for tr.Next() {
		tk := tr.Token()
		switch tk.kind {
		case tokenKindNewline:
			continue
		// top level:
		case tokenKindIdent:
			switch string(tk.concrete) {
			case "enum":
				if nextRecordOpCode != 0 {
					return f, fmt.Sprintf("enums may not have attached op codes")
				}
				en, err := readEnum(tr)
				if err != nil {
					return f, err
				}
				en.Comment = nextRecordComment
				f.Enums = append(f.Enums, en)

				nextRecordComment = ""
				nextRecordOpCode = 0
			case "struct":
				st, err := readStruct(tr)
				if err != nil {
					return f, err
				}
				st.Comment = nextRecordComment
				st.OpCode = nextRecordOpCode
				f.Structs = append(f.Structs, st)

				nextRecordComment = ""
				nextRecordOpCode = 0
			case "message":
				msg, err := readMessage(tr)
				if err != nil {
					return f, err
				}
				msg.Comment = nextRecordComment
				msg.OpCode = nextRecordOpCode
				f.Messages = append(f.Messages, msg)

				nextRecordComment = ""
				nextRecordOpCode = 0
			}
		case tokenKindBlockComment:
			nextRecordComment = string(tk.concrete[2:len(tk.concrete-2)])
		case tokenKindLineComment:
			nextRecordComment = string(tk.concrete[2:])
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
				base := 10
				if strings.HasPrefix(content, "0x") {
					base = 16
				}
				opCode, err := strconv.ParseInt(content, base, 32)
				if err != nil {
					return f, err
				}
				nextRecordOpCode = int32(opCode)
			} else if tk.kind == tokenKindStringLiteral {
				if len(tk.concrete) > 4 {
					return f, fmt.Errorf("opcode string %s exceeds 4 ascii characters", string(tk.concrete))
				}
				for _, b := range tk.concrete {
					nextRecordOpCode <<= 8
					nextRecordOpCode |= int32(b)
				}

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
	hasNext := !tr.Next()
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
		// transparently skip newlines and comments even if they aren't expected
		if tk.kind == tokenKindNewline || tk.kind == tokenKindBlockComment || tk.kind == tokenKindLineComment {
			return expectNext(tr, kinds...)
		}
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
	nextComment := ""
	nextDeprecatedMessage := ""
	nextIsDeprecated := false
	tr.Next()
	for tr.Token().kind != tokenKindCloseCurly {
		if !tr.Next() {
			return en, fmt.Errorf("enum definition ended early")
		}
		tk := tr.Token()
		switch tk.kind {
		case tokenKindIdent:
			optName := string(tk.concrete)
			if err := expectNext(tr, tokenKindEquals); err != nil {
				return en, err
			}
			if err := expectNext(tr, tokenKindInteger); err != nil {
				return en, err
			}
			optInteger := strconv.ParseInt(string(tr.Token().concrete), 10, 32)
		case tokenKindOpenSquare:
			// deprecated note
			if err := expectNext(tr, tokenKindIdent); err != nil {
				return en, err
			}
			tk = tr.Token()
			if string(tk.concrete) != "deprecated" {
				return fmt.Errorf("invalid ident after leading '[', got %s, wanted %v", string(tk.concrete), "deprecated")
			}
			if err := expectNext(tr, tokenKindOpenParen); err != nil {
				return en, err
			}
			if err := expectNext(tr, tokenKindStringLiteral); err != nil {
				return en, err
			}
			tk = tr.Token()
			nextIsDeprecated = true
			var err error
			nextDeprecatedMessage, err = strconv.Unquote(string(tk.concrete))
			if err != nil {
				return en, err
			}
			if err := expectNext(tr, tokenKindCloseParen); err != nil {
				return en, err
			}
			if err := expectNext(tr, tokenKindCloseSquare); err != nil {
				return en, err
			}
		case tokenKindBlockComment:
			nextComment = string(tk.concrete[2:len(tk.concrete-2)])
		case tokenKindLineComment:
			nextComment = string(tk.concrete[2:])
		}
	}

	return en, nil
}

func readStruct(tr *tokenReader) (Struct, error) {
	return Struct{}, nil
}

func readMessage(tr *tokenReader) (Message, error) {
	return Message{}, nil
}

type Struct struct {
	Name    string
	Comment string
	OpCode  int32
	Fields  []Field
}

type Message struct {
	Name    string
	Comment string
	OpCode  int32
	Fields  map[int32]Field
}

type Field struct {
	// KeyType is only provided if Map is true.
	KeyType   string
	ValueType string
	Name      string
	Comment   string
	// DeprecatedMessage is only provided if Deprecated is true.
	DeprecatedMessage string
	Repeated          bool
	Map               bool
	Deprecated        bool
}

type Enum struct {
	Name    string
	Comment string
	Options []EnumOption
}

type EnumOption struct {
	Name    string
	Comment string
	Value   int32
}
