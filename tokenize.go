package bebop

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

type tokenReader struct {
	r             *bufio.Reader
	lastToken     token
	nextToken     token
	err           error
	loc           location
	keepNextToken bool
}

type location struct {
	absoluteChar int
	lineChar     int
	line         int
}

func (l *location) inc(i int) {
	l.absoluteChar += i
	l.lineChar += i
}

func (l *location) incLine() {
	l.line++
	l.lineChar = 0
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

var singleCharTokens = map[byte]tokenKind{
	'=':  tokenKindEquals,
	'[':  tokenKindOpenSquare,
	']':  tokenKindCloseSquare,
	'{':  tokenKindOpenCurly,
	'}':  tokenKindCloseCurly,
	'(':  tokenKindOpenParen,
	')':  tokenKindCloseParen,
	',':  tokenKindComma,
	';':  tokenKindSemicolon,
	'\n': tokenKindNewline,
}

func (tr *tokenReader) setNextToken(tk token) {
	tr.lastToken = tr.nextToken
	tr.nextToken = tk
	tr.nextToken.loc = tr.loc
}

func (tr *tokenReader) readByte() (byte, error) {
	b, err := tr.r.ReadByte()
	tr.loc.inc(1)
	return b, err
}

func (tr *tokenReader) unreadByte() {
	err := tr.r.UnreadByte()
	if err != nil {
		panic(fmt.Errorf("unreadByte failed: %w", err))
	}
	tr.loc.inc(-1)
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
	for {
		b, err := tr.readByte()
		if err == io.EOF {
			return false
		}
		if kind, ok := singleCharTokens[b]; ok {
			if b == '\n' {
				tr.loc.incLine()
			}
			tr.setNextToken(token{
				concrete: []byte{b},
				kind:     kind,
			})
			return true
		}
		switch b {
		case '"':
			return tr.nextStringLiteral(b)
		case ' ', '\t', '\r':
			continue
		// two token sequences
		case '/':
			b2, err := tr.readByte()
			if err == io.EOF {
				tr.err = fmt.Errorf("eof waiting for '[/, *]' after '/'")
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
				tr.loc.inc(len(restOfLine))
				tr.loc.incLine()

				tr.setNextToken(token{
					concrete: []byte{b, b2},
					kind:     tokenKindLineComment,
				})
				tr.nextToken.concrete = append(tr.nextToken.concrete, restOfLine...)
			} else if b2 == '*' {
				return tr.nextBlockComment(b, b2)
			} else {
				tr.err = fmt.Errorf("unexpected token '%v' waiting for '[/, *]' after '/'", string(b2))
				return false
			}
		case '-':
			b2, err := tr.readByte()
			if err == io.EOF {
				tr.err = fmt.Errorf("eof waiting for (['>', number]) after '-'")
				return false
			}
			if err != nil {
				tr.err = err
				return false
			}
			if b2 == '>' {
				tr.setNextToken(token{
					concrete: []byte{b, b2},
					kind:     tokenKindArrow,
				})
				return true
			}

			// 'inf' case
			if b2 == 'i' {
				b3, err := tr.readByte()
				if err == io.EOF {
					tr.err = fmt.Errorf("eof waiting for 'n' in '-inf'")
					return false
				}
				if err != nil {
					tr.err = err
					return false
				}
				b4, err := tr.readByte()
				if err == io.EOF {
					tr.err = fmt.Errorf("eof waiting for 'f' in '-inf'")
					return false
				}
				if err != nil {
					tr.err = err
					return false
				}
				if b3 != 'n' || b4 != 'f' {
					tr.err = fmt.Errorf("unexpected token %v in '-inf'", string([]byte{b3, b4}))
					return false
				}
				tr.setNextToken(token{
					concrete: []byte{b, b2, b3, b4},
					kind:     tokenKindNegativeInf,
				})
				return true
			}
			// number case
			if isNumeric(b2) {
				return tr.nextNumber([]byte{b, b2})
			}

			tr.err = fmt.Errorf("unexpected token '%v' waiting for (['>', number]) after '-'", string(b2))
			return false
		default:
			if isNumeric(b) {
				return tr.nextNumber([]byte{b})
			}
			tr.unreadByte()
			rn, sz, err := tr.r.ReadRune()
			tr.loc.inc(sz)
			if err == io.ErrUnexpectedEOF || err == io.EOF {
				tr.err = io.ErrUnexpectedEOF
				return false
			} else if err != nil {
				tr.err = err
				return false
			} else if unicode.IsLetter(rn) {
				return tr.nextIdent(rn)
			}
			tr.err = fmt.Errorf("unexpected token '%v', expected number, letter, or control sequence", string(b))
			return false
		}
		return true
	}
}

func (tr *tokenReader) nextNumber(firstBytes []byte) bool {
	tk := token{
		concrete: firstBytes,
		kind:     tokenKindIntegerLiteral,
	}
	// second byte is allowed to be 'x' for hex or '.' for floats
	secondByte := true
	invalidLastChar := false
	for {
		b, err := tr.readByte()
		if err == io.EOF {
			if invalidLastChar {
				tr.err = fmt.Errorf("unexpected eof, expected number following %q", string(tk.concrete))
				return false
			}
			// stream ended in number
			tr.setNextToken(tk)
			return true
		}
		if err != nil {
			tr.err = err
			return false
		}
		if secondByte && b == 'x' {
			invalidLastChar = true
			tk.concrete = append(tk.concrete, b)
		} else if b == '.' {
			invalidLastChar = true
			tk.concrete = append(tk.concrete, b)
			tk.kind = tokenKindFloatLiteral
		} else if isNumeric(b) {
			invalidLastChar = false
			tk.concrete = append(tk.concrete, b)
		} else if b == 'e' {
			invalidLastChar = true
			// this is allowed if its not the last character, i.e.
			// 1.6e442
			tk.concrete = append(tk.concrete, b)
		} else {
			if invalidLastChar {
				tr.err = fmt.Errorf("unexpected token '%v', expected number following %q",
					string(b), string(tk.concrete))
				return false
			}
			// something else is here
			tr.unreadByte()
			tr.setNextToken(tk)
			return true
		}
		secondByte = false
	}
}

func isNumeric(b byte) bool {
	return b >= 0x30 && b <= 0x39
}

var keywords = map[string]tokenKind{
	"readonly":   tokenKindReadOnly,
	"message":    tokenKindMessage,
	"struct":     tokenKindStruct,
	"enum":       tokenKindEnum,
	"deprecated": tokenKindDeprecated,
	"opcode":     tokenKindOpCode,
	"map":        tokenKindMap,
	"array":      tokenKindArray,
	"union":      tokenKindUnion,
	"const":      tokenKindConst,
	"inf":        tokenKindInf,
	"nan":        tokenKindNaN,
	"true":       tokenKindTrue,
	"false":      tokenKindFalse,
	"import":     tokenKindImport,
}

func (tr *tokenReader) nextIdent(firstRune rune) bool {
	tk := token{
		concrete: []byte(string(firstRune)),
		kind:     tokenKindIdent,
	}
	for {
		rn, sz, err := tr.r.ReadRune()
		tr.loc.inc(sz)
		if err == io.EOF {
			// stream ended in ident
			if keywordKind, ok := keywords[string(tk.concrete)]; ok {
				tk.kind = keywordKind
			}
			tr.setNextToken(tk)
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
			// ignore this error, we just called ReadRune
			_ = tr.r.UnreadRune()
			tr.loc.inc(-sz)
			if keywordKind, ok := keywords[string(tk.concrete)]; ok {
				tk.kind = keywordKind
			}
			tr.setNextToken(tk)
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
	var escaping bool
	for {
		b, err := tr.readByte()
		if err == io.EOF {
			tr.err = fmt.Errorf("eof waiting for string end quote")
			return false
		}
		if err != nil {
			tr.err = err
			return false
		}
		tk.concrete = append(tk.concrete, b)
		if b == '"' && !escaping {
			tr.setNextToken(tk)
			return true
		}
		if b == '\\' && !escaping {
			escaping = true
		} else {
			escaping = false
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
		b, err := tr.readByte()
		if err == io.EOF {
			tr.err = fmt.Errorf("block comment missing end token")
			return false
		}
		if err != nil {
			tr.err = err
			return false
		}
		tk.concrete = append(tk.concrete, b)
		if lastByte == '*' && b == '/' {
			tr.skipFollowingWhitespace()
			tr.setNextToken(tk)
			return true
		}
		lastByte = b
	}
}

func (tr *tokenReader) skipFollowingWhitespace() {
	for {
		b, _ := tr.readByte()
		switch b {
		case '\n':
			tr.loc.incLine()
			fallthrough
		case ' ', '\r':
			continue
		}
		tr.unreadByte()
		break
	}
}

func (tr *tokenReader) Err() error {
	if tr.err == nil {
		return nil
	}
	return fmt.Errorf("[%d:%d] %s", tr.loc.line, tr.loc.lineChar, tr.err.Error())
}

var optionalSemicolons = false

// Token returns the last read token from the underlying reader.
// If no token has been read, it returns an empty token (token{}).
func (tr *tokenReader) Token() token {
	if optionalSemicolons {
		if tr.nextToken.kind == tokenKindNewline {
			switch tr.lastToken.kind {
			case tokenKindIdent, tokenKindIntegerLiteral:
				injectedToken := token{
					kind:     tokenKindSemicolon,
					concrete: []byte{';'},
					loc:      tr.loc,
				}
				tr.keepNextToken = true
				tr.lastToken = injectedToken
				return injectedToken
			}
		}
	}
	//fmt.Println("tr:", tr.nextToken.kind.String(), string(tr.nextToken.concrete))
	return tr.nextToken
}
