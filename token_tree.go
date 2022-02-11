package bebop

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// a tokenTree accepts a series of tokens defined as byte sequences to later
// identify and isolate those tokens given a roughly arbitrary byte reader. It also
// supports ignoring sets of bytes to enable skipping whitespace or other meaningless
// bytes between tokens. Each token string can have its own custom parsing function
// to handle successive bytes once enough bytes have been read to identify the exact
// token parsing strategy that should be taken. The structure is not sufficiently
// generic to accept an io.Reader, and the interface it would need would be arbitrarily
// specific, so it just accepts a tokenReader at this time.
type tokenTree struct {
	successors map[byte]*tokenTree
	skipOver   map[byte]struct{}
	isTerminal bool
	build      func(*tokenReader, []byte) token
}

// skip registers a byte to skip at the root of this tree. It should not be used on
// non-root trees.
func (v *tokenTree) skip(b byte) {
	if v.skipOver == nil {
		v.skipOver = make(map[byte]struct{})
	}
	v.skipOver[b] = struct{}{}
}

// add registers a byte sequence to build a token given the provided build function.
func (v *tokenTree) add(str []byte, build func(*tokenReader, []byte) token) {
	if len(str) == 0 {
		v.isTerminal = true
		v.build = build
		return
	}
	if v.successors == nil {
		v.successors = make(map[byte]*tokenTree)
	}
	if v.successors[str[0]] == nil {
		v.successors[str[0]] = &tokenTree{}
	}
	v.successors[str[0]].add(str[1:], build)
}

// findFirst is a helper to call find with no concrete bytes
func (v *tokenTree) findFirst(tr *tokenReader) (token, bool) {
	return v.find(tr, []byte{})
}

// find looks for a token which matches the bytes read from the token reader,
// or if this tree is a terminal tree (add was called on it with an empty slice),
// builds a token from this tree. If EOF is reached on the token reader,
// false will be returned and either io.EOF or an io.UnexpectedEOF error will be
// added to the token reader, the former if concrete is empty, the latter if concrete
// is non empty (because we failed to finish parsing a token). If the token reader
// provides bytes which are not defined by the tree, the tree will either return that
// there is no valid token (if concrete is empty) or attempt trivial corrective fixes
// to continue reading tokens, adding an error to the token reader at the same time.
func (v *tokenTree) find(tr *tokenReader, concrete []byte) (token, bool) {
	// TODO: allow v.isTerminal to be overridden if more valid bytes exist (for a token set like [=,=>,==,<=,<,<<,>,>>,>=])
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

// nextValidBytes returns all potential valid single bytes to produce a token from this tree, sorted
func (v *tokenTree) nextValidBytes() []string {
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
