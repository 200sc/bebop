package bebop

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type tokenReader struct {
	tree               *tokenTree
	r                  *bufio.Reader
	lastToken          token
	nextToken          token
	errs               []locError
	loc                location
	keepNextToken      bool
	optionalSemicolons bool
}

func newTokenReader(r io.Reader) *tokenReader {
	// We buffer the reader to reduce the number
	// of actual read calls we make out to a file,
	// and to ease reading individual bytes
	return &tokenReader{
		r:    bufio.NewReader(r),
		tree: newTokenTree(),
		loc: location{
			line: 1,
		},
	}
}

// UnNext tells the next Next call to not update the returned token
func (tr *tokenReader) UnNext() {
	tr.keepNextToken = true
}

func (tr *tokenReader) setNextToken(tk token) {
	tr.lastToken = tr.nextToken
	tr.nextToken = tk
	tr.nextToken.end = tr.loc
}

func (tr *tokenReader) readByte() (byte, error) {
	b, err := tr.r.ReadByte()
	if err != io.EOF {
		tr.loc.inc(1)
	}
	if b == '\n' {
		tr.loc.incLine()
	}
	return b, err
}

func (tr *tokenReader) unreadByte() {
	err := tr.r.UnreadByte()
	if err != nil {
		panic(fmt.Errorf("unreadByte failed: %w", err))
	}
	tr.loc.inc(-1)
}

func (tr *tokenReader) addError(err error) {
	tr.errs = append(tr.errs, locError{
		loc: tr.loc,
		err: err,
	})
}

func (tr *tokenReader) Err() error {
	if len(tr.errs) == 0 {
		return nil
	}
	errStrs := []string{}
	for _, err := range tr.errs {
		errStrs = append(errStrs, err.Error())
	}
	return fmt.Errorf(strings.Join(errStrs, "\n"))
}

// Token returns the last read token from the underlying reader.
// If no token has been read, it returns an empty token (token{}).
func (tr *tokenReader) Token() token {
	if tr.optionalSemicolons {
		if tr.nextToken.kind == tokenKindNewline || tr.nextToken.kind == tokenKindCloseCurly || tr.nextToken.kind == tokenKindLineComment {
			switch tr.lastToken.kind {
			case tokenKindIdent, tokenKindIntegerLiteral, tokenKindNegativeInf, tokenKindInf,
				tokenKindStringLiteral, tokenKindNaN, tokenKindTrue, tokenKindFalse:
				injectedToken := token{
					kind:     tokenKindSemicolon,
					concrete: []byte{';'},
					start:    tr.loc,
					end:      tr.loc,
				}
				tr.keepNextToken = true
				tr.lastToken = injectedToken
				return injectedToken
			}
		}
	}
	return tr.nextToken
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
	// find all byte-driven tokens
	start := tr.loc
	tk, ok := tr.tree.findFirst(tr)
	if len(tr.errs) != 0 {
		lastErr := tr.errs[len(tr.errs)-1]
		if errors.Is(lastErr, io.EOF) {
			tr.errs = tr.errs[:len(tr.errs)-1]
			return false
		}
		if errors.Is(lastErr, io.ErrUnexpectedEOF) {
			return false
		}
		// other errors should have been corrected
	}
	if ok {
		tk.start = start
		tr.setNextToken(tk)
		return true
	}
	tr.unreadByte()

	// find rune-driven tokens
	// TODO: adapt the finder tree structure to handle runes
	rn, sz, err := tr.r.ReadRune()
	tr.loc.inc(sz)
	if err == io.ErrUnexpectedEOF || err == io.EOF {
		tr.addError(io.ErrUnexpectedEOF)
		return false
	} else if err != nil {
		tr.addError(err)
		return false
	} else if unicode.IsLetter(rn) {
		return tr.nextIdent(rn)
	}
	tr.addError(fmt.Errorf("unexpected token '%v', expected number, letter, or control sequence", string(rn)))
	return false
}

func isHex(b byte) bool {
	return (b >= 'a' && b <= 'f') ||
		(b >= 'A' && b <= 'F')
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
	"flags":      tokenKindFlags,
	"service":    tokenKindService,
	"void":       tokenKindVoid,
}

func (tr *tokenReader) nextIdent(firstRune rune) bool {
	tk := token{
		concrete: []byte(string(firstRune)),
		kind:     tokenKindIdent,
		start:    tr.loc,
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
			tr.addError(err)
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

func (tr *tokenReader) skipFollowingWhitespace() {
	for {
		b, _ := tr.readByte()
		switch b {
		case '\n', ' ', '\r':
			continue
		}
		tr.unreadByte()
		break
	}
}

func simpleToken(k tokenKind) func(*tokenReader, []byte) token {
	return func(_ *tokenReader, concrete []byte) token {
		return token{
			kind:     k,
			concrete: concrete,
		}
	}
}

func numberToken(tr *tokenReader, concrete []byte) token {
	tk := token{
		concrete: concrete,
		kind:     tokenKindIntegerLiteral,
	}
	// second byte is allowed to be 'x' for hex
	secondByte := true
	hex := false
	decimal := false
	invalidLastChar := false
	for {
		b, err := tr.readByte()
		if err == io.EOF {
			if invalidLastChar {
				tr.addError(fmt.Errorf("%w: expected number following %q", io.ErrUnexpectedEOF, (tk.concrete)))
				return token{}
			}
			// stream ended in number
			return tk
		}
		if err != nil {
			tr.addError(err)
			return token{}
		}
		if secondByte && b == 'x' {
			hex = true
			invalidLastChar = true
			tk.concrete = append(tk.concrete, b)
		} else if b == '.' {
			if decimal {
				tr.addError(fmt.Errorf("unexpected second period in float following %q", string(tk.concrete)))
				return token{}
			}
			invalidLastChar = true
			tk.concrete = append(tk.concrete, b)
			tk.kind = tokenKindFloatLiteral
			decimal = true
		} else if isNumeric(b) {
			invalidLastChar = false
			tk.concrete = append(tk.concrete, b)
		} else if hex && isHex(b) {
			invalidLastChar = false
			tk.concrete = append(tk.concrete, b)
		} else if b == 'e' {
			invalidLastChar = true
			// this is allowed if its not the last character, i.e.
			// 1.6e442
			tk.concrete = append(tk.concrete, b)
		} else {
			if invalidLastChar {
				tr.addError(fmt.Errorf("unexpected token '%v', expected number following %q",
					string(b), string(tk.concrete)))
				tk.concrete = append(tk.concrete, 0)
				return tk
			}
			// something else is here
			tr.unreadByte()
			tr.setNextToken(tk)
			return tk
		}
		secondByte = false
	}
}

func lineCommentToken(tr *tokenReader, concrete []byte) token {
	restOfLine, err := tr.r.ReadBytes('\n')
	if err != nil && err != io.EOF {
		tr.addError(err)
		return token{}
	}
	tr.loc.incLine()

	tk := token{
		concrete: concrete,
		kind:     tokenKindLineComment,
	}
	tk.concrete = append(tk.concrete, restOfLine...)
	return tk
}

func blockCommentToken(tr *tokenReader, concrete []byte) token {
	tk := token{
		concrete: concrete,
		kind:     tokenKindBlockComment,
	}
	var lastByte byte
	for {
		b, err := tr.readByte()
		if err == io.EOF {
			tr.addError(fmt.Errorf("%w: block comment missing end token", io.ErrUnexpectedEOF))
			return token{}
		}
		if err != nil {
			tr.addError(err)
			return token{}
		}
		tk.concrete = append(tk.concrete, b)
		if lastByte == '*' && b == '/' {
			tr.skipFollowingWhitespace()
			return tk
		}
		lastByte = b
	}
}

func stringLiteralToken(tr *tokenReader, concrete []byte) token {
	tk := token{
		concrete: concrete,
		kind:     tokenKindStringLiteral,
	}
	var escaping bool
	for {
		b, err := tr.readByte()
		if err == io.EOF {
			tr.addError(fmt.Errorf("%w: waiting for string end quote", io.ErrUnexpectedEOF))
			return token{}
		}
		if err != nil {
			tr.addError(err)
			return token{}
		}
		tk.concrete = append(tk.concrete, b)
		if b == '"' && !escaping {
			return tk
		}
		if b == '\\' && !escaping {
			escaping = true
		} else {
			escaping = false
		}
	}
}

// build a token tree with our set of non-rune-based tokens
func newTokenTree() *tokenTree {
	tt := &tokenTree{}
	// skip whitespace
	tt.skip(' ')
	tt.skip('\t')
	tt.skip('\r')

	// simple terminals
	tt.add([]byte{'='}, simpleToken(tokenKindEquals))
	tt.add([]byte{'['}, simpleToken(tokenKindOpenSquare))
	tt.add([]byte{']'}, simpleToken(tokenKindCloseSquare))
	tt.add([]byte{'{'}, simpleToken(tokenKindOpenCurly))
	tt.add([]byte{'}'}, simpleToken(tokenKindCloseCurly))
	tt.add([]byte{'('}, simpleToken(tokenKindOpenParen))
	tt.add([]byte{')'}, simpleToken(tokenKindCloseParen))
	tt.add([]byte{','}, simpleToken(tokenKindComma))
	tt.add([]byte{';'}, simpleToken(tokenKindSemicolon))
	tt.add([]byte{':'}, simpleToken(tokenKindColon))
	tt.add([]byte{'\n'}, simpleToken(tokenKindNewline))
	tt.add([]byte{'|'}, simpleToken(tokenKindVerticalBar))
	tt.add([]byte{'&'}, simpleToken(tokenKindAmpersand))
	tt.add([]byte{'-', '>'}, simpleToken(tokenKindArrow))
	tt.add([]byte{'>', '>'}, simpleToken(tokenKindDoubleCaretRight))
	tt.add([]byte{'<', '<'}, simpleToken(tokenKindDoubleCaretLeft))
	tt.add([]byte{'-', 'i', 'n', 'f'}, simpleToken(tokenKindNegativeInf))

	// complex terminals
	tt.add([]byte{'"'}, stringLiteralToken)
	tt.add([]byte{'/', '/'}, lineCommentToken)
	tt.add([]byte{'/', '*'}, blockCommentToken)
	for c := byte('0'); c <= '9'; c++ {
		tt.add([]byte{'-', c}, numberToken)
		tt.add([]byte{c}, numberToken)
	}
	return tt
}
