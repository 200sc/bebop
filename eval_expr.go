package bebop

import (
	"fmt"
	"strconv"
)

func evaluateBitflagExpr(n bitFlagExprNode, opts []EnumOption, uinttype bool, bitsize int) (int64, uint64, error) {
	if uinttype {
		switch bitsize {
		case 8:
			uval, err := evaluateBitflagExpr_uint8(n, opts)
			return 0, uint64(uval), err
		case 16:
			uval, err := evaluateBitflagExpr_uint16(n, opts)
			return 0, uint64(uval), err
		case 32:
			uval, err := evaluateBitflagExpr_uint32(n, opts)
			return 0, uint64(uval), err
		case 64:
			uval, err := evaluateBitflagExpr_uint64(n, opts)
			return 0, uint64(uval), err
		}
	} else {
		switch bitsize {
		case 16:
			val, err := evaluateBitflagExpr_int16(n, opts)
			return int64(val), 0, err
		case 32:
			val, err := evaluateBitflagExpr_int32(n, opts)
			return int64(val), 0, err
		case 64:
			val, err := evaluateBitflagExpr_int64(n, opts)
			return int64(val), 0, err
		}
	}
	panic("invalid uinttype / bitsize combination")
}

func evaluateBitflagExpr_int64(n bitFlagExprNode, opts []EnumOption) (int64, error) {
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
		optInteger, err := strconv.ParseInt(string(v.tk.concrete), 0, 64)
		if err != nil {
			return 0, err
		}
		return optInteger, nil
	case parenNode:
		return evaluateBitflagExpr_int64(v.inner, opts)
	case binOpNode:
		lhs, err := evaluateBitflagExpr_int64(v.lhs, opts)
		if err != nil {
			return 0, err
		}
		rhs, err := evaluateBitflagExpr_int64(v.rhs, opts)
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

func evaluateBitflagExpr_int32(n bitFlagExprNode, opts []EnumOption) (int32, error) {
	switch v := n.(type) {
	case identNode:
		name := string(v.tk.concrete)
		for _, o := range opts {
			if o.Name == name {
				return int32(o.Value), nil
			}
		}
		return 0, readError(v.tk, "enum option %v undefined", name)
	case numberNode:
		optInteger, err := strconv.ParseInt(string(v.tk.concrete), 0, 32)
		if err != nil {
			return 0, err
		}
		return int32(optInteger), nil
	case parenNode:
		return evaluateBitflagExpr_int32(v.inner, opts)
	case binOpNode:
		lhs, err := evaluateBitflagExpr_int32(v.lhs, opts)
		if err != nil {
			return 0, err
		}
		rhs, err := evaluateBitflagExpr_int32(v.rhs, opts)
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

func evaluateBitflagExpr_int16(n bitFlagExprNode, opts []EnumOption) (int16, error) {
	switch v := n.(type) {
	case identNode:
		name := string(v.tk.concrete)
		for _, o := range opts {
			if o.Name == name {
				return int16(o.Value), nil
			}
		}
		return 0, readError(v.tk, "enum option %v undefined", name)
	case numberNode:
		optInteger, err := strconv.ParseInt(string(v.tk.concrete), 0, 16)
		if err != nil {
			return 0, err
		}
		return int16(optInteger), nil
	case parenNode:
		return evaluateBitflagExpr_int16(v.inner, opts)
	case binOpNode:
		lhs, err := evaluateBitflagExpr_int16(v.lhs, opts)
		if err != nil {
			return 0, err
		}
		rhs, err := evaluateBitflagExpr_int16(v.rhs, opts)
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

func evaluateBitflagExpr_uint64(n bitFlagExprNode, opts []EnumOption) (uint64, error) {
	switch v := n.(type) {
	case identNode:
		name := string(v.tk.concrete)
		for _, o := range opts {
			if o.Name == name {
				return o.UintValue, nil
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
		return evaluateBitflagExpr_uint64(v.inner, opts)
	case binOpNode:
		lhs, err := evaluateBitflagExpr_uint64(v.lhs, opts)
		if err != nil {
			return 0, err
		}
		rhs, err := evaluateBitflagExpr_uint64(v.rhs, opts)
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

func evaluateBitflagExpr_uint32(n bitFlagExprNode, opts []EnumOption) (uint32, error) {
	switch v := n.(type) {
	case identNode:
		name := string(v.tk.concrete)
		for _, o := range opts {
			if o.Name == name {
				return uint32(o.UintValue), nil
			}
		}
		return 0, readError(v.tk, "enum option %v undefined", name)
	case numberNode:
		optInteger, err := strconv.ParseUint(string(v.tk.concrete), 0, 32)
		if err != nil {
			return 0, err
		}
		return uint32(optInteger), nil
	case parenNode:
		return evaluateBitflagExpr_uint32(v.inner, opts)
	case binOpNode:
		lhs, err := evaluateBitflagExpr_uint32(v.lhs, opts)
		if err != nil {
			return 0, err
		}
		rhs, err := evaluateBitflagExpr_uint32(v.rhs, opts)
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

func evaluateBitflagExpr_uint16(n bitFlagExprNode, opts []EnumOption) (uint16, error) {
	switch v := n.(type) {
	case identNode:
		name := string(v.tk.concrete)
		for _, o := range opts {
			if o.Name == name {
				return uint16(o.UintValue), nil
			}
		}
		return 0, readError(v.tk, "enum option %v undefined", name)
	case numberNode:
		optInteger, err := strconv.ParseUint(string(v.tk.concrete), 0, 16)
		if err != nil {
			return 0, err
		}
		return uint16(optInteger), nil
	case parenNode:
		return evaluateBitflagExpr_uint16(v.inner, opts)
	case binOpNode:
		lhs, err := evaluateBitflagExpr_uint16(v.lhs, opts)
		if err != nil {
			return 0, err
		}
		rhs, err := evaluateBitflagExpr_uint16(v.rhs, opts)
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

func evaluateBitflagExpr_uint8(n bitFlagExprNode, opts []EnumOption) (uint8, error) {
	switch v := n.(type) {
	case identNode:
		name := string(v.tk.concrete)
		for _, o := range opts {
			if o.Name == name {
				return uint8(o.UintValue), nil
			}
		}
		return 0, readError(v.tk, "enum option %v undefined", name)
	case numberNode:
		optInteger, err := strconv.ParseUint(string(v.tk.concrete), 0, 8)
		if err != nil {
			return 0, err
		}
		return uint8(optInteger), nil
	case parenNode:
		return evaluateBitflagExpr_uint8(v.inner, opts)
	case binOpNode:
		lhs, err := evaluateBitflagExpr_uint8(v.lhs, opts)
		if err != nil {
			return 0, err
		}
		rhs, err := evaluateBitflagExpr_uint8(v.rhs, opts)
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
