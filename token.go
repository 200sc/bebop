package bebop

type token struct {
	kind           tokenKind
	concrete       []byte
	newlineFollows bool
}

type tokenKind uint8

const (
	tokenKindInvalid tokenKind = iota

	tokenKindIdent
	tokenKindInteger
	tokenKindStringLiteral

	tokenKindReadOnly
	tokenKindStruct
	tokenKindMessage
	tokenKindEnum
	tokenKindDeprecated
	tokenKindOpCode
	tokenKindMap
	tokenKindArray

	tokenKindOpenSquare
	tokenKindCloseSquare
	tokenKindOpenParen
	tokenKindCloseParen
	tokenKindOpenCurly
	tokenKindCloseCurly
	tokenKindSemicolon
	tokenKindComma
	tokenKindEquals
	tokenKindArrow
	tokenKindLineComment
	tokenKindBlockComment

	tokenKindNewline
)

var tokenStrings = map[tokenKind]string{
	tokenKindInvalid:       "Invalid",
	tokenKindIdent:         "Ident",
	tokenKindInteger:       "Integer",
	tokenKindReadOnly:      "Readonly",
	tokenKindStruct:        "Struct",
	tokenKindMessage:       "Message",
	tokenKindEnum:          "Enum",
	tokenKindDeprecated:    "Deprecated",
	tokenKindOpCode:        "OpCode",
	tokenKindMap:           "Map",
	tokenKindArray:         "Array",
	tokenKindStringLiteral: "String Literal",
	tokenKindOpenSquare:    "Open Square",
	tokenKindCloseSquare:   "Close Square",
	tokenKindOpenParen:     "Open Paren",
	tokenKindCloseParen:    "Close Paren",
	tokenKindSemicolon:     "Semicolon",
	tokenKindOpenCurly:     "Open Curly",
	tokenKindCloseCurly:    "Close Curly",
	tokenKindComma:         "Comma",
	tokenKindEquals:        "Equals",
	tokenKindArrow:         "Arrow",
	tokenKindLineComment:   "Line Comment",
	tokenKindBlockComment:  "Block Comment",
	tokenKindNewline:       "Newline",
}

func (tk tokenKind) String() string {
	return tokenStrings[tk]
}
