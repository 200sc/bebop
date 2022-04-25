package bebop

import (
	"fmt"
	"strconv"
)

func readBitflagExpr(tr *tokenReader, previousOptions []EnumOption) (uint64, error) {
	// to simplify this a little bit, read everything up until ;
	toks, err := readUntil(tr, tokenKindSemicolon)
	if err != nil {
		return 0, err
	}
	// TODO: precedence (although there isn't really a natural precedence for these
	// operations anyway)
	parsed, err := parseBitflagExpr(toks)
	if err != nil {
		return 0, err
	}
	val, err := evaluateBitflagExpr(parsed, previousOptions)
	if err != nil {
		return 0, err
	}
	return val, nil
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

func evaluateBitflagExpr(n bitFlagExprNode, opts []EnumOption) (uint64, error) {
	switch v := n.(type) {
	case identNode:
		name := string(v.tk.concrete)
		for _, o := range opts {
			if o.Name == name {
				return o.Value, nil
			}
		}
		return 0, readError(v.tk, "enum option %v undefined", name)
	case numberNode:
		optInteger, err := strconv.ParseUint(string(v.tk.concrete), 0, 64)
		if err != nil {
			return 0, err
		}
		return optInteger, nil
	case parenNode:
		return evaluateBitflagExpr(v.inner, opts)
	case binOpNode:
		lhs, err := evaluateBitflagExpr(v.lhs, opts)
		if err != nil {
			return 0, err
		}
		rhs, err := evaluateBitflagExpr(v.rhs, opts)
		if err != nil {
			return 0, err
		}
		switch v.op {
		// TODO: confirm that the behavior of these operators in Go
		// matches the expected behavior in bebop
		case tokenKindAmpersand:
			return lhs & rhs, nil
		case tokenKindVerticalBar:
			return lhs | rhs, nil
		case tokenKindDoubleCaretLeft:
			return lhs << rhs, nil
		case tokenKindDoubleCaretRight:
			return lhs >> rhs, nil
		default:
			return 0, fmt.Errorf("undefined binary operator for bitflag %v", v.op)
			// TODO: plus, minus, multiply, the options are endless
		}
	default:
		return 0, fmt.Errorf("undefined bit flag node type %T", n)
	}
}

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
