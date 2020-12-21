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

type nextResp struct {
	hasNext    bool
	hasNewline bool
}

func (tr *tokenReader) setNextToken(tk token) {
	tr.lastToken = tr.nextToken
	tr.nextToken = tk
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
		b, err := tr.r.ReadByte()
		if err == io.EOF {
			return false
		}
		if kind, ok := singleCharTokens[b]; ok {
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
			tr.setNextToken(token{
				concrete: []byte{b, b2},
				kind:     tokenKindArrow,
			})
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

				tr.setNextToken(token{
					concrete: []byte{b, b2},
					kind:     tokenKindLineComment,
				})
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
			}
			tr.r.UnreadByte()
			rn, _, err := tr.r.ReadRune()
			if err == io.ErrUnexpectedEOF || err == io.EOF {
				tr.err = io.ErrUnexpectedEOF
				return false
			} else if err != nil {
				tr.err = err
				return false
			} else if unicode.IsLetter(rn) {
				return tr.nextIdent(rn)
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
			tr.setNextToken(tk)
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
			tr.setNextToken(tk)
			return true
		}
		secondByte = false
	}
}

func isInteger(b byte) bool {
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
			tr.r.UnreadRune()
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
		b, err := tr.r.ReadByte()
		if err == io.EOF {
			tr.err = fmt.Errorf("string literal missing end quote")
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
		b, err := tr.r.ReadByte()
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
		b, _ := tr.r.ReadByte()
		switch b {
		case ' ', '\r', '\n':
			continue
		}
		tr.r.UnreadByte()
		break
	}
}

func (tr *tokenReader) Err() error {
	return tr.err
}

var optionalSemicolons = false

// Token returns the last read token from the underlying reader.
// If no token has been read, it returns an empty token (token{}).
func (tr *tokenReader) Token() token {
	if optionalSemicolons {
		if tr.nextToken.kind == tokenKindNewline {
			switch tr.lastToken.kind {
			case tokenKindIdent, tokenKindInteger:
				injectedToken := token{
					kind:     tokenKindSemicolon,
					concrete: []byte{';'},
				}
				tr.keepNextToken = true
				tr.lastToken = injectedToken
				return injectedToken
			}
		}
	}
	return tr.nextToken
}
