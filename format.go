package bebop

import (
	"fmt"
	"io"
)

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
			if lastToken.kind == tokenKindIdent {
				fmt.Fprint(w, " ")
			}
			fmt.Fprint(w, string(t.concrete))
		case tokenKindOpenCurly:
			if lastToken.kind == tokenKindIdent {
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
			if lastToken.kind == tokenKindInteger {
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
