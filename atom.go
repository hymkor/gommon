package gmnlisp

import (
	"context"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type _TrueType struct{}

func (_TrueType) PrintTo(w io.Writer, m PrintMode) (int, error) {
	return io.WriteString(w, "t")
}

func (t _TrueType) Eval(context.Context, *World) (Node, error) {
	return t, nil
}

var True Node = _TrueType{}

func (_TrueType) Equals(n Node, m EqlMode) bool {
	_, ok := n.(_TrueType)
	return ok
}

type _NullType struct{}

func (_NullType) PrintTo(w io.Writer, m PrintMode) (int, error) {
	return io.WriteString(w, "nil")
}

func (nt _NullType) Eval(context.Context, *World) (Node, error) {
	return nt, nil
}

func (nt _NullType) Equals(n Node, m EqlMode) bool {
	if n == nil {
		return true
	}
	_, ok := n.(_NullType)
	return ok
}

var Null Node = _NullType{}

type Symbol int

var symbols = []string{}

var symbolMap = map[string]Symbol{}

func NewSymbol(s string) Symbol {
	if value, ok := symbolMap[s]; ok {
		return value
	}
	value := Symbol(len(symbolMap))
	symbolMap[s] = value
	symbols = append(symbols, s)
	return value
}

func cmdGensym(ctx context.Context, w *World, node Node) (Node, error) {
	return NewSymbol(fmt.Sprintf("-gensym-%d-", len(symbolMap))), nil
}

func (s Symbol) PrintTo(w io.Writer, m PrintMode) (int, error) {
	return io.WriteString(w, symbols[s])
}

func (s Symbol) Eval(_ context.Context, w *World) (Node, error) {
	return w.Get(s)
}

func (s Symbol) Equals(n Node, m EqlMode) bool {
	ns, ok := n.(Symbol)
	return ok && s == ns
}

func (s Symbol) String() string {
	return symbols[s]
}

func (s Symbol) GoString() string {
	return symbols[s]
}

type Rune rune

func (r Rune) PrintTo(w io.Writer, m PrintMode) (int, error) {
	if m == PRINT {
		switch r {
		case '\t':
			return io.WriteString(w, `#\tab`)
		case '\n':
			return io.WriteString(w, `#\linefeed`)
		case '\r':
			return io.WriteString(w, `#\return`)
		case ' ':
			return io.WriteString(w, `#\space`)
		default:
			if unicode.IsLetter(rune(r)) {
				return fmt.Fprintf(w, `#\%c`, rune(r))
			} else {
				return fmt.Fprintf(w, `#\U%04X`, rune(r))
			}
		}
	} else {
		return fmt.Fprintf(w, "%c", rune(r))
	}
}

func (r Rune) Eval(_ context.Context, w *World) (Node, error) {
	return r, nil
}

func (r Rune) Equals(n Node, m EqlMode) bool {
	if value, ok := n.(Rune); ok {
		if r == value {
			return true
		}
		if m == EQUALP {
			return unicode.ToLower(rune(r)) == unicode.ToLower(rune(value))
		}
		return false
	}
	if m == STRICT {
		return false
	}
	return false
}

func (r Rune) Add(n Node) (Node, error) {
	if value, ok := n.(Integer); ok {
		return r + Rune(value), nil
	}
	if value, ok := n.(Rune); ok {
		return r + value, nil
	}
	return nil, MakeError(ErrNotSupportType, n)
}

func (r Rune) Sub(n Node) (Node, error) {
	if value, ok := n.(Integer); ok {
		return r - Rune(value), nil
	}
	if value, ok := n.(Rune); ok {
		return r - value, nil
	}
	return nil, MakeError(ErrNotSupportType, n)
}

func compareRune(argv []Node, f func(rune) bool) (Node, error) {
	left, ok := argv[0].(Rune)
	if !ok {
		return nil, MakeError(ErrExpectedString, argv[0])
	}
	right, ok := argv[1].(Rune)
	if !ok {
		return nil, MakeError(ErrExpectedString, argv[1])
	}
	cmp := rune(left) - rune(right)
	if f(cmp) {
		return True, nil
	}
	return Null, nil
}

func funRuneLt(ctx context.Context, w *World, argv []Node) (Node, error) {
	return compareRune(argv, func(cmp rune) bool { return cmp < 0 })
}
func funRuneLe(ctx context.Context, w *World, argv []Node) (Node, error) {
	return compareRune(argv, func(cmp rune) bool { return cmp <= 0 })
}
func funRuneEq(ctx context.Context, w *World, argv []Node) (Node, error) {
	return compareRune(argv, func(cmp rune) bool { return cmp == 0 })
}
func funRuneGt(ctx context.Context, w *World, argv []Node) (Node, error) {
	return compareRune(argv, func(cmp rune) bool { return cmp > 0 })
}
func funRuneGe(ctx context.Context, w *World, argv []Node) (Node, error) {
	return compareRune(argv, func(cmp rune) bool { return cmp >= 0 })
}
func funRuneNe(ctx context.Context, w *World, argv []Node) (Node, error) {
	return compareRune(argv, func(cmp rune) bool { return cmp != 0 })
}
func funRuneIndex(ctx context.Context, w *World, argv []Node) (Node, error) {
	_char, ok := argv[0].(Rune)
	if !ok {
		return nil, ErrExpectedCharacter
	}
	char := rune(_char)
	str, ok := argv[1].(String)
	if !ok {
		return nil, ErrExpectedString
	}
	var start int = 0
	if len(argv) >= 3 {
		_start, ok := argv[2].(Integer)
		if !ok {
			return nil, ErrExpectedNumber
		}
		start = int(_start)
	}
	i := 0
	for _, c := range string(str) {
		if i >= start && c == char {
			return Integer(i), nil
		}
		i++
	}
	return Null, nil
}

type Keyword string

func (k Keyword) PrintTo(w io.Writer, m PrintMode) (int, error) {
	return io.WriteString(w, string(k))
}

func (k Keyword) Eval(context.Context, *World) (Node, error) {
	return k, nil
}

func (k Keyword) Equals(n Node, m EqlMode) bool {
	if other, ok := n.(Keyword); ok {
		return strings.EqualFold(string(k), string(other))
	}
	return false
}
