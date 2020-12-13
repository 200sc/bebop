package bebop

type token struct {
	kind     tokenKind
	concrete []byte
}

type tokenKind uint8

const (
	tokenKindInvalid tokenKind = iota

	tokenKindIdent         tokenKind = iota
	tokenKindInteger       tokenKind = iota
	tokenKindStringLiteral tokenKind = iota

	tokenKindOpenSquare   tokenKind = iota
	tokenKindCloseSquare  tokenKind = iota
	tokenKindOpenParen    tokenKind = iota
	tokenKindCloseParen   tokenKind = iota
	tokenKindOpenCurly    tokenKind = iota
	tokenKindCloseCurly   tokenKind = iota
	tokenKindSemicolon    tokenKind = iota
	tokenKindNewline      tokenKind = iota
	tokenKindComma        tokenKind = iota
	tokenKindEquals       tokenKind = iota
	tokenKindArrow        tokenKind = iota
	tokenKindLineComment  tokenKind = iota
	tokenKindBlockComment tokenKind = iota
)

var tokenStrings = map[tokenKind]string{
	tokenKindInvalid:       "Invalid",
	tokenKindIdent:         "Ident",
	tokenKindInteger:       "Integer",
	tokenKindStringLiteral: "String Literal",
	tokenKindOpenSquare:    "Open Square",
	tokenKindCloseSquare:   "Close Square",
	tokenKindOpenParen:     "Open Paren",
	tokenKindCloseParen:    "Close Paren",
	tokenKindSemicolon:     "Semicolon",
	tokenKindNewline:       "Newline",
	tokenKindOpenCurly:     "Open Curly",
	tokenKindCloseCurly:    "Close Curly",
	tokenKindComma:         "Comma",
	tokenKindEquals:        "Equals",
	tokenKindArrow:         "Arrow",
	tokenKindLineComment:   "Line Comment",
	tokenKindBlockComment:  "Block Comment",
}

func (tk tokenKind) String() string {
	return tokenStrings[tk]
}
