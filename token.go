package bebop

import "fmt"

type token struct {
	kind     tokenKind
	concrete []byte
	start    location
	end      location
}

func (t token) String() string {
	return fmt.Sprintf("%s-%s %s", t.start.String(), t.end.String(), t.kind.String())
}

type tokenKind uint8

const (
	tokenKindInvalid tokenKind = iota

	tokenKindIdent
	tokenKindIntegerLiteral
	tokenKindFloatLiteral
	tokenKindStringLiteral

	tokenKindReadOnly
	tokenKindStruct
	tokenKindMessage
	tokenKindEnum
	tokenKindDeprecated
	tokenKindOpCode
	tokenKindMap
	tokenKindArray
	tokenKindUnion
	tokenKindConst
	tokenKindInf
	tokenKindNegativeInf
	tokenKindNaN
	tokenKindTrue
	tokenKindFalse
	tokenKindImport
	tokenKindFlags
	tokenKindService

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
	tokenKindVerticalBar
	tokenKindAmpersand
	tokenKindDoubleCaretLeft
	tokenKindDoubleCaretRight
	tokenKindColon

	tokenKindNewline

	tokenKindFinal
)

var tokenStrings = map[tokenKind]string{
	tokenKindInvalid:          "Invalid",
	tokenKindIdent:            "Ident",
	tokenKindIntegerLiteral:   "Integer Literal",
	tokenKindReadOnly:         "Readonly",
	tokenKindStruct:           "Struct",
	tokenKindMessage:          "Message",
	tokenKindEnum:             "Enum",
	tokenKindUnion:            "Union",
	tokenKindDeprecated:       "Deprecated",
	tokenKindOpCode:           "OpCode",
	tokenKindMap:              "Map",
	tokenKindArray:            "Array",
	tokenKindStringLiteral:    "String Literal",
	tokenKindOpenSquare:       "Open Square",
	tokenKindCloseSquare:      "Close Square",
	tokenKindOpenParen:        "Open Paren",
	tokenKindCloseParen:       "Close Paren",
	tokenKindSemicolon:        "Semicolon",
	tokenKindOpenCurly:        "Open Curly",
	tokenKindCloseCurly:       "Close Curly",
	tokenKindComma:            "Comma",
	tokenKindEquals:           "Equals",
	tokenKindArrow:            "Arrow",
	tokenKindLineComment:      "Line Comment",
	tokenKindBlockComment:     "Block Comment",
	tokenKindNewline:          "Newline",
	tokenKindConst:            "Const",
	tokenKindFloatLiteral:     "Floating Point Literal",
	tokenKindInf:              "Infinity",
	tokenKindNegativeInf:      "Negative Infinity",
	tokenKindNaN:              "NaN",
	tokenKindTrue:             "True",
	tokenKindFalse:            "False",
	tokenKindImport:           "Import",
	tokenKindFlags:            "Flags",
	tokenKindVerticalBar:      "Vertical Bar",
	tokenKindAmpersand:        "Ampersand",
	tokenKindDoubleCaretLeft:  "Double Caret Left",
	tokenKindDoubleCaretRight: "Double Caret Right",
	tokenKindColon:            "Colon",
	tokenKindService:          "Service",
}

func (tk tokenKind) String() string {
	return tokenStrings[tk]
}
