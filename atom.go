package gmnlisp

import (
	"context"
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
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

type String string

func (s String) String() string {
	return string(s)
}

func (s String) GoString() string {
	return "\"" + string(s) + "\""
}

func (s String) PrintTo(w io.Writer, m PrintMode) (int, error) {
	if m == PRINC {
		return io.WriteString(w, string(s))
	} else {
		return fmt.Fprintf(w, `"%s"`, unescapeSequenceReplacer.Replace(string(s)))
	}
}

func (s String) Eval(context.Context, *World) (Node, error) {
	return s, nil
}

func (s String) Equals(n Node, m EqlMode) bool {
	ns, ok := n.(String)
	if !ok {
		_ns, ok := n.(String)
		if !ok {
			return false
		}
		if m == STRICT {
			return false
		}
		ns = String(string(_ns))
	}
	if m == EQUALP {
		return strings.EqualFold(string(s), string(ns))
	}
	return string(s) == string(ns)
}

func (s String) firstRuneAndRestString() (Rune, String, bool) {
	if len(s) <= 0 {
		return Rune(utf8.RuneError), "", false
	}
	r, siz := utf8.DecodeRuneInString(string(s))
	return Rune(r), String(s[siz:]), true
}

func (s String) EachRune(f func(Rune) error) error {
	for _, r := range s {
		if err := f(Rune(r)); err != nil {
			return err
		}
	}
	return nil
}

func (s String) FirstAndRest() (Node, Node, bool, func(Node) error) {
	if len(s) <= 0 {
		return nil, Null, false, nil
	}
	r, siz := utf8.DecodeRuneInString(string(s))
	return Rune(r), String(s[siz:]), true, func(value Node) error {
		return ErrNotSupportType
	}
}

func (s String) Add(n Node) (Node, error) {
	if value, ok := n.(String); ok {
		news := make([]byte, 0, len(s)+len(value)+1)
		news = append(news, s...)
		news = append(news, value...)
		return String(news), nil
	}
	return nil, fmt.Errorf("%w: `%s`", ErrNotSupportType, ToString(n, PRINT))
}

func (s String) LessThan(n Node) (bool, error) {
	ns, ok := n.(String)
	if !ok {
		return false, fmt.Errorf("%w: `%s`", ErrNotSupportType, ToString(n, PRINT))
	}
	return string(s) < string(ns), nil
}

var unescapeSequenceReplacer = strings.NewReplacer(
	"\n", "\\n",
	"\r", "\\r",
	"\t", "\\t",
	"\b", "\\b",
	"\"", "\\\"",
)

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
	return nil, fmt.Errorf("%w: `%s`", ErrNotSupportType, ToString(n, PRINT))
}

func (r Rune) Sub(n Node) (Node, error) {
	if value, ok := n.(Integer); ok {
		return r - Rune(value), nil
	}
	if value, ok := n.(Rune); ok {
		return r - value, nil
	}
	return nil, fmt.Errorf("%w: `%s`", ErrNotSupportType, ToString(n, PRINT))
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
