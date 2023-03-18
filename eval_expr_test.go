package bebop

import "testing"

func Test_evalauteBitflagExpr(t *testing.T) {
	t.Parallel()
	t.Run("invalid combination", func(t *testing.T) {
		t.Parallel()
		_, _, err := evaluateBitflagExpr(nil, nil, false, 128)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
	t.Run("invalid node type", func(t *testing.T) {
		t.Parallel()
		t.Run("signed", func(t *testing.T) {
			t.Parallel()
			_, err := evaluateBitflagExpSigned[int16](nil, nil)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
		t.Run("unsigned", func(t *testing.T) {
			t.Parallel()
			_, err := evaluateBitflagExprUnsigned[uint16](nil, nil)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	})
	t.Run("invalid binary operator", func(t *testing.T) {
		t.Parallel()
		t.Run("signed", func(t *testing.T) {
			t.Parallel()
			n := binOpNode{
				lhs: numberNode{
					tk: token{
						concrete: []byte("1"),
					},
				},
				rhs: numberNode{
					tk: token{
						concrete: []byte("1"),
					},
				},
				op: tokenKindCloseSquare,
			}
			_, err := evaluateBitflagExpSigned[int16](n, nil)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
		t.Run("unsigned", func(t *testing.T) {
			t.Parallel()
			n := binOpNode{
				lhs: numberNode{
					tk: token{
						concrete: []byte("1"),
					},
				},
				rhs: numberNode{
					tk: token{
						concrete: []byte("1"),
					},
				},
				op: tokenKindCloseSquare,
			}
			_, err := evaluateBitflagExprUnsigned[uint16](n, nil)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	})
}
