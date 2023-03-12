package bebop

import (
	"fmt"
	"strconv"
)

func evaluateBitflagExpr(n bitFlagExprNode, opts []EnumOption, uinttype bool, bitsize int) (int64, uint64, error) {
	if uinttype {
		switch bitsize {
		case 8:
			uval, err := evaluateBitflagExprUnsigned[uint8](n, opts)
			return 0, uint64(uval), err
		case 16:
			uval, err := evaluateBitflagExprUnsigned[uint16](n, opts)
			return 0, uint64(uval), err
		case 32:
			uval, err := evaluateBitflagExprUnsigned[uint32](n, opts)
			return 0, uint64(uval), err
		case 64:
			uval, err := evaluateBitflagExprUnsigned[uint64](n, opts)
			return 0, uint64(uval), err
		}
	} else {
		switch bitsize {
		case 16:
			val, err := evaluateBitflagExpSigned[int16](n, opts)
			return int64(val), 0, err
		case 32:
			val, err := evaluateBitflagExpSigned[int32](n, opts)
			return int64(val), 0, err
		case 64:
			val, err := evaluateBitflagExpSigned[int64](n, opts)
			return int64(val), 0, err
		}
	}
	return 0, 0, fmt.Errorf("invalid uinttype / bitsize combination")
}

func evaluateBitflagExpSigned[T signedInteger](n bitFlagExprNode, opts []EnumOption) (T, error) {
	switch v := n.(type) {
	case identNode:
		name := string(v.tk.concrete)
		for _, o := range opts {
			if o.Name == name {
				return T(o.Value), nil
			}
		}
		return 0, readError(v.tk, "enum option %v undefined", name)
	case numberNode:
		optInteger, err := strconv.ParseInt(string(v.tk.concrete), 0, 64)
		if err != nil {
			return 0, err
		}
		return T(optInteger), nil
	case parenNode:
		return evaluateBitflagExpSigned[T](v.inner, opts)
	case binOpNode:
		lhs, err := evaluateBitflagExpSigned[T](v.lhs, opts)
		if err != nil {
			return 0, err
		}
		rhs, err := evaluateBitflagExpSigned[T](v.rhs, opts)
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

func evaluateBitflagExprUnsigned[T unsignedInteger](n bitFlagExprNode, opts []EnumOption) (T, error) {
	switch v := n.(type) {
	case identNode:
		name := string(v.tk.concrete)
		for _, o := range opts {
			if o.Name == name {
				return T(o.UintValue), nil
			}
		}
		return 0, readError(v.tk, "enum option %v undefined", name)
	case numberNode:
		optInteger, err := strconv.ParseUint(string(v.tk.concrete), 0, 64)
		if err != nil {
			return 0, err
		}
		return T(optInteger), nil
	case parenNode:
		return evaluateBitflagExprUnsigned[T](v.inner, opts)
	case binOpNode:
		lhs, err := evaluateBitflagExprUnsigned[T](v.lhs, opts)
		if err != nil {
			return 0, err
		}
		rhs, err := evaluateBitflagExprUnsigned[T](v.rhs, opts)
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

type unsignedInteger interface {
	uint8 | uint16 | uint32 | uint64
}

type signedInteger interface {
	int16 | int32 | int64
}
