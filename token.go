package bebop

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

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

type tokenReader struct {
	r             *bufio.Reader
	nextToken     token
	err           error
	keepNextToken bool
}

func newTokenReader(r io.Reader) *tokenReader {
	// We buffer the reader to reduce the number
	// of actual read calls we make out to a file,
	// and to ease reading individual bytes
	bufferedReader := bufio.NewReader(r)
	return &tokenReader{r: bufferedReader}
}

// UnNext tells the next Next call to not update the returned token
func (tr *tokenReader) UnNext() {
	tr.keepNextToken = true
}

// Next attempts to read the next token in the reader.
// If a token cannot be found, it returns false. If there
// are no tokens because EOF was reached, Err() will return
// nil. Otherwise whatever error encountered will be returned
// by Err().
// The read token, if this returns true, can be obtained via
// Token().
func (tr *tokenReader) Next() bool {
	if tr.keepNextToken {
		tr.keepNextToken = false
		return true
	}
	// read until we get whitespace
	for {
		b, err := tr.r.ReadByte()
		if err == io.EOF {
			return false
		}
		switch b {
		case '"':
			return tr.nextStringLiteral(b)
		case '=':
			tr.nextToken = token{
				concrete: []byte{b},
				kind:     tokenKindEquals,
			}
		case '[':
			tr.nextToken = token{
				concrete: []byte{b},
				kind:     tokenKindOpenSquare,
			}
		case ']':
			tr.nextToken = token{
				concrete: []byte{b},
				kind:     tokenKindCloseSquare,
			}
		case '(':
			tr.nextToken = token{
				concrete: []byte{b},
				kind:     tokenKindOpenParen,
			}
		case ')':
			tr.nextToken = token{
				concrete: []byte{b},
				kind:     tokenKindCloseParen,
			}
		case '{':
			tr.nextToken = token{
				concrete: []byte{b},
				kind:     tokenKindOpenCurly,
			}
		case '}':
			tr.nextToken = token{
				concrete: []byte{b},
				kind:     tokenKindCloseCurly,
			}
		case ',':
			tr.nextToken = token{
				concrete: []byte{b},
				kind:     tokenKindComma,
			}
		case ';':
			tr.nextToken = token{
				concrete: []byte{b},
				kind:     tokenKindSemicolon,
			}
		case '\n':
			tr.nextToken = token{
				concrete: []byte{b},
				kind:     tokenKindNewline,
			}
		case ' ', '\t', '\r':
			continue
		// two token sequences
		case '-':
			b2, err := tr.r.ReadByte()
			if err == io.EOF {
				tr.err = io.ErrUnexpectedEOF
				return false
			}
			if err != nil {
				tr.err = err
				return false
			}
			if b2 != '>' {
				tr.err = fmt.Errorf("invalid token")
				return false
			}
			tr.nextToken = token{
				concrete: []byte{b, b2},
				kind:     tokenKindArrow,
			}
		case '/':
			b2, err := tr.r.ReadByte()
			if err == io.EOF {
				tr.err = io.ErrUnexpectedEOF
				return false
			}
			if err != nil {
				tr.err = err
				return false
			}
			if b2 == '/' {
				restOfLine, err := tr.r.ReadBytes('\n')
				if err != nil && err != io.EOF {
					tr.err = err
					return false
				}

				tr.nextToken = token{
					concrete: []byte{b, b2},
					kind:     tokenKindLineComment,
				}
				tr.nextToken.concrete = append(tr.nextToken.concrete, restOfLine...)
			} else if b2 == '*' {
				return tr.nextBlockComment(b, b2)
			} else {
				tr.err = fmt.Errorf("invalid token")
				return false
			}
		default:
			if isInteger(b) {
				return tr.nextInteger(b)
			} else {
				tr.r.UnreadByte()
				rn, _, err := tr.r.ReadRune()
				if err == io.ErrUnexpectedEOF || err == io.EOF {
					tr.err = io.ErrUnexpectedEOF
					return false
				} else if err != nil {
					tr.err = err
					return false
				} else {
					if unicode.IsLetter(rn) {
						return tr.nextIdent(rn)
					}
				}
			}
			tr.err = fmt.Errorf("invalid token")
			return false
		}
		return true
	}
}

func (tr *tokenReader) nextInteger(firstByte byte) bool {
	tk := token{
		concrete: []byte{firstByte},
		kind:     tokenKindInteger,
	}
	// second byte is allowed to be 'x'
	secondByte := true
	for {
		b, err := tr.r.ReadByte()
		if err == io.EOF {
			// stream ended in number
			tr.nextToken = tk
			return true
		}
		if err != nil {
			tr.err = err
			return false
		}
		if secondByte && b == 'x' {
			tk.concrete = append(tk.concrete, b)
		} else if isInteger(b) {
			tk.concrete = append(tk.concrete, b)
		} else {
			// something else is here
			tr.r.UnreadByte()
			tr.nextToken = tk
			return true
		}
		secondByte = false
	}
}

func isInteger(b byte) bool {
	return b >= 0x30 && b <= 0x39
}

func (tr *tokenReader) nextIdent(firstRune rune) bool {
	tk := token{
		concrete: []byte(string(firstRune)),
		kind:     tokenKindIdent,
	}
	for {
		rn, _, err := tr.r.ReadRune()
		if err == io.EOF {
			// stream ended in ident
			tr.nextToken = tk
			return true
		}
		if err != nil {
			tr.err = err
			return false
		}
		switch {
		case unicode.IsLetter(rn):
		case unicode.IsDigit(rn):
		case rn == '_':
		default:
			tr.r.UnreadRune()
			tr.nextToken = tk
			return true
		}
		tk.concrete = append(tk.concrete, []byte(string(rn))...)
	}
}

func (tr *tokenReader) nextStringLiteral(firstByte byte) bool {
	tk := token{
		concrete: []byte{firstByte},
		kind:     tokenKindStringLiteral,
	}
	for {
		b, err := tr.r.ReadByte()
		if err == io.EOF {
			tr.err = fmt.Errorf("string literal missing end quote")
			return false
		}
		if err != nil {
			tr.err = err
			return false
		}
		// Todo: escaped quotes in strings
		tk.concrete = append(tk.concrete, b)
		if b == '"' {
			tr.nextToken = tk
			return true
		}
	}
}

func (tr *tokenReader) nextBlockComment(b1, b2 byte) bool {
	tk := token{
		concrete: []byte{b1, b2},
		kind:     tokenKindBlockComment,
	}
	var lastByte byte
	for {
		b, err := tr.r.ReadByte()
		if err == io.EOF {
			tr.err = fmt.Errorf("block comment missing end token")
			return false
		}
		if err != nil {
			tr.err = err
			return false
		}
		// Todo: escaped quotes in strings
		tk.concrete = append(tk.concrete, b)
		if lastByte == '*' && b == '/' {
			tr.nextToken = tk
			return true
		}
		lastByte = b
	}
}

func (tr *tokenReader) Err() error {
	return tr.err
}

// Token returns the last read token from the underlying reader.
// If no token has been read, it returns an empty token (token{}).
func (tr *tokenReader) Token() token {
	return tr.nextToken
}

func format(tokens []token, w io.Writer) {
	var lastToken token
	setFirstOnLine := false
	firstOnLine := true
	inRecord := false
	i := 0
	for i < len(tokens) {
		t := tokens[i]
		switch t.kind {
		case tokenKindOpenSquare:
			if firstOnLine {
				if inRecord {
					fmt.Fprint(w, "\t")
				}
				// if this is valid, this is an opcode or deprecated pattern
				// the next tokens are:
				// opcode/deprecated
				// (
				// integer/string literal
				// )
				// ]
				fmt.Fprint(w, string(t.concrete))
				i++
				fmt.Fprint(w, string(tokens[i].concrete))
				i++
				fmt.Fprint(w, string(tokens[i].concrete))
				i++
				fmt.Fprint(w, string(tokens[i].concrete))
				i++
				fmt.Fprint(w, string(tokens[i].concrete))
				i++
				fmt.Fprint(w, string(tokens[i].concrete))
				fmt.Fprint(w, "\n")
				firstOnLine = true
			} else {
				fmt.Fprint(w, string(t.concrete))
			}
		case tokenKindNewline:
			// TODO: respect some provided newlines
			//if lastToken.kind == tokenKindNewline {
			//	fmt.Fprint(w, string(t.concrete))
			//}
			setFirstOnLine = true
		case tokenKindSemicolon:
			fmt.Fprint(w, string(t.concrete)+"\n")
			setFirstOnLine = true
		case tokenKindEquals:
			switch lastToken.kind {
			case tokenKindIdent:
				fmt.Fprint(w, " ")
			}
			fmt.Fprint(w, string(t.concrete))
		case tokenKindOpenCurly:
			switch lastToken.kind {
			case tokenKindIdent:
				fmt.Fprint(w, " ")
			}
			fmt.Fprint(w, string(t.concrete)+"\n")
			setFirstOnLine = true
			inRecord = true
		case tokenKindCloseCurly:
			// Todo: we shouldn't double newline here at the end of the file
			fmt.Fprint(w, string(t.concrete)+"\n\n")
			setFirstOnLine = true
			inRecord = false
		case tokenKindArrow:
			switch lastToken.kind {
			case tokenKindInteger:
				fmt.Fprint(w, " ")
			}
			fmt.Fprint(w, string(t.concrete))
		case tokenKindBlockComment:
			if inRecord && firstOnLine {
				fmt.Fprint(w, "\t")
			}
			fmt.Fprint(w, string(t.concrete)+"\n")
			setFirstOnLine = true
		case tokenKindInteger, tokenKindIdent:
			if inRecord && firstOnLine {
				fmt.Fprint(w, "\t")
			}
			switch lastToken.kind {
			case tokenKindCloseSquare, tokenKindComma, tokenKindArrow,
				tokenKindInteger, tokenKindIdent, tokenKindEquals:
				fmt.Fprint(w, " ")
			}
			fallthrough
		default:
			fmt.Fprint(w, string(t.concrete))
		}
		lastToken = t
		firstOnLine = setFirstOnLine
		setFirstOnLine = false
		i++
	}
}
