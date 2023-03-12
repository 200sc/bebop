package bebop

import (
	"fmt"
)

func readBitflagExpr(tr *tokenReader, previousOptions []EnumOption, uinttype bool, bitsize int) (int64, uint64, error) {
	// to simplify this a little bit, read everything up until ;
	toks, err := readUntil(tr, tokenKindSemicolon)
	if err != nil {
		return 0, 0, err
	}
	// TODO: precedence (although there isn't really a natural precedence for these
	// operations anyway)
	parsed, err := parseBitflagExpr(toks)
	if err != nil {
		return 0, 0, err
	}
	val, uval, err := evaluateBitflagExpr(parsed, previousOptions, uinttype, bitsize)
	if err != nil {
		return 0, 0, err
	}
	return val, uval, nil
}

// productions
// start = expr
// expr = integer | ident
//        expr binop expr |
//        ( expr )

type bitFlagExprNode interface {
	isNode()
}

type parenNode struct {
	inner bitFlagExprNode
}

func (parenNode) isNode() {}

type binOpNode struct {
	lhs, rhs bitFlagExprNode
	op       tokenKind
}

func (binOpNode) isNode() {}

type numberNode struct {
	tk token
}

func (numberNode) isNode() {}

type identNode struct {
	tk token
}

func (identNode) isNode() {}

func parseBitflagExpr(toks []token) (bitFlagExprNode, error) {
	if len(toks) == 0 {
		return nil, fmt.Errorf("empty bitflag expression")
	}
	i := 0
	thisTk := toks[0]
	var lhs bitFlagExprNode
	switch thisTk.kind {
	case tokenKindIdent:
		lhs = identNode{
			tk: thisTk,
		}
	case tokenKindIntegerLiteral:
		lhs = numberNode{
			tk: thisTk,
		}
	case tokenKindOpenParen:
		i++
		var err error
		lhs, i, err = parseParenExpr(i, toks)
		if err != nil {
			return lhs, err
		}
	default:
		return nil, readError(thisTk, "bad token for bitflag expr %v", thisTk.kind)
	}
	i++
	if i >= len(toks) {
		return lhs, nil
	}
	binOpTk := toks[i]
	switch binOpTk.kind {
	case tokenKindAmpersand, tokenKindVerticalBar, tokenKindDoubleCaretLeft, tokenKindDoubleCaretRight:
	default:
		return nil, readError(binOpTk, "undefined binary operator for bitflag %v", binOpTk.kind)
	}
	tree := binOpNode{
		lhs: lhs,
		op:  binOpTk.kind,
	}
	i++
	rhs, err := parseBitflagExpr(toks[i:])
	if err != nil {
		return nil, err
	}
	tree.rhs = rhs
	return tree, nil
}

func parseParenExpr(j int, tokens []token) (node bitFlagExprNode, newI int, err error) {
	startJ := j
	closeNeeded := 1
	for j < len(tokens) {
		if tokens[j].kind == tokenKindOpenParen {
			closeNeeded++
		}
		if tokens[j].kind == tokenKindCloseParen {
			closeNeeded--
			if closeNeeded == 0 {
				break
			}
		}
		j++
	}
	inner, err := parseBitflagExpr(tokens[startJ:j])
	if err != nil {
		return nil, 0, err
	}
	return parenNode{
		inner: inner,
	}, j, nil
}
