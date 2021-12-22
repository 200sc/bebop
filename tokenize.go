package bebop

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"unicode"
)

type locError struct {
	loc location
	err error
}

func (l locError) Unwrap() error {
	return l.err
}

func (l locError) Error() string {
	return fmt.Sprintf("[%d:%d] %s", l.loc.line, l.loc.lineChar, l.err)
}

type tokenReader struct {
	finder             *tokenFinder
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
		r:      bufio.NewReader(r),
		finder: newTokenFinder(),
	}
}

// UnNext tells the next Next call to not update the returned token
func (tr *tokenReader) UnNext() {
	tr.keepNextToken = true
}

func (tr *tokenReader) setNextToken(tk token) {
	tr.lastToken = tr.nextToken
	tr.nextToken = tk
	tr.nextToken.loc = tr.loc
}

func (tr *tokenReader) readByte() (byte, error) {
	b, err := tr.r.ReadByte()
	if err != io.EOF {
		tr.loc.inc(1)
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
	tk, ok := tr.finder.findFirst(tr)
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
		if tk.kind == tokenKindNewline {
			tr.loc.incLine()
		}
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
					loc:      tr.loc,
				}
				tr.keepNextToken = true
				tr.lastToken = injectedToken
				return injectedToken
			}
		}
	}
	return tr.nextToken
}

type location struct {
	lineChar int
	line     int
}

func (l *location) inc(i int) {
	l.lineChar += i
}

func (l *location) incLine() {
	l.line++
	l.lineChar = 0
}

type tokenFinder struct {
	successors map[byte]*tokenFinder
	skipOver   map[byte]struct{}
	isTerminal bool
	build      func(*tokenReader, []byte) token
}

func (v *tokenFinder) skip(b byte) {
	if v.skipOver == nil {
		v.skipOver = make(map[byte]struct{})
	}
	v.skipOver[b] = struct{}{}
}

func (v *tokenFinder) add(str []byte, build func(*tokenReader, []byte) token) {
	if len(str) == 0 {
		v.isTerminal = true
		v.build = build
		return
	}
	if v.successors == nil {
		v.successors = make(map[byte]*tokenFinder)
	}
	if v.successors[str[0]] == nil {
		v.successors[str[0]] = &tokenFinder{}
	}
	v.successors[str[0]].add(str[1:], build)
}

func (v *tokenFinder) findFirst(tr *tokenReader) (token, bool) {
	return v.find(tr, []byte{})
}

func (v *tokenFinder) find(tr *tokenReader, concrete []byte) (token, bool) {
	// after find is executed, if false is returned, the last read token
	// (which should be reevaluated for other token cases) will be tr.lastToken.
	if v.isTerminal {
		return v.build(tr, concrete), true
	}
	var b byte
	var err error
	for {
		b, err = tr.readByte()
		if err == io.EOF {
			if len(concrete) != 0 {
				nextValid := v.nextValidBytes()
				tr.addError(fmt.Errorf("%w: waiting for [%v] after '%v'", io.ErrUnexpectedEOF, strings.Join(nextValid, ", "), string(concrete)))
			} else {
				tr.addError(io.EOF)
			}
			return token{}, false
		}
		if err != nil {
			tr.addError(err)
			return token{}, false
		}
		if v.skipOver != nil {
			if _, ok := v.skipOver[b]; ok {
				continue
			}
		}
		break
	}
	t, ok := v.successors[b]
	if !ok {
		if len(concrete) != 0 {
			nextValid := v.nextValidBytes()

			tr.addError(fmt.Errorf("unexpected token '%v' waiting for [%v] after '%v'", string(b), strings.Join(nextValid, ", "), string(concrete)))
			// greedily choose the first option
			b = byte(nextValid[0][0])
			t = v.successors[b]
		} else {
			return token{}, false
		}
	}
	concrete = append(concrete, b)
	return t.find(tr, concrete)
}

func (v *tokenFinder) nextValidBytes() []string {
	opts := make(map[string]struct{}, len(v.successors))
	for k := range v.successors {
		if isNumeric(k) {
			// Assumption: if this accepts a digit, it accepts any number
			opts["number"] = struct{}{}
			continue
		}
		opts[string(k)] = struct{}{}
	}
	strs := make([]string, len(opts))
	i := 0
	for k := range opts {
		strs[i] = k
		i++
	}
	sort.Strings(strs)

	return strs
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

func newTokenFinder() *tokenFinder {
	tf := &tokenFinder{}
	tf.skip(' ')
	tf.skip('\t')
	tf.skip('\r')
	// TODO: support terminals that are prefixed by other smaller terminals
	// e.g. right now we can't have == and = as terminals
	tf.add([]byte{'='}, simpleToken(tokenKindEquals))
	tf.add([]byte{'['}, simpleToken(tokenKindOpenSquare))
	tf.add([]byte{']'}, simpleToken(tokenKindCloseSquare))
	tf.add([]byte{'{'}, simpleToken(tokenKindOpenCurly))
	tf.add([]byte{'}'}, simpleToken(tokenKindCloseCurly))
	tf.add([]byte{'('}, simpleToken(tokenKindOpenParen))
	tf.add([]byte{')'}, simpleToken(tokenKindCloseParen))
	tf.add([]byte{','}, simpleToken(tokenKindComma))
	tf.add([]byte{';'}, simpleToken(tokenKindSemicolon))
	tf.add([]byte{'\n'}, simpleToken(tokenKindNewline))
	tf.add([]byte{'|'}, simpleToken(tokenKindVerticalBar))
	tf.add([]byte{'&'}, simpleToken(tokenKindAmpersand))
	tf.add([]byte{'-', '>'}, simpleToken(tokenKindArrow))
	tf.add([]byte{'>', '>'}, simpleToken(tokenKindDoubleCaretRight))
	tf.add([]byte{'<', '<'}, simpleToken(tokenKindDoubleCaretLeft))
	tf.add([]byte{'-', 'i', 'n', 'f'}, simpleToken(tokenKindNegativeInf))
	tf.add([]byte{'"'}, stringLiteralToken)
	for c := byte('0'); c <= '9'; c++ {
		tf.add([]byte{'-', c}, numberToken)
		tf.add([]byte{c}, numberToken)
	}
	tf.add([]byte{'/', '/'}, lineCommentToken)
	tf.add([]byte{'/', '*'}, blockCommentToken)
	// */ is also technically a token, but it is parsed as a part of parsing /*
	return tf
}
