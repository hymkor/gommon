package gmnlisp

import (
	"fmt"
	"io"
	"strings"
)

type Node interface {
	Eval(*World) (Node, error)
	Equals(Node) bool
	Equalp(Node) bool
	PrintTo(io.Writer)
	PrincTo(io.Writer)
}

func toString(node Node) string {
	if node == nil {
		return "()"
	}
	var buffer strings.Builder
	node.PrintTo(&buffer)
	return buffer.String()
}

type _TrueType struct{}

func (_TrueType) PrintTo(w io.Writer) {
	io.WriteString(w, "T")
}

func (_TrueType) PrincTo(w io.Writer) {
	io.WriteString(w, "T")
}

func (t _TrueType) Eval(*World) (Node, error) {
	return t, nil
}

var True Node = _TrueType{}

func (_TrueType) Equals(n Node) bool {
	_, ok := n.(_TrueType)
	return ok
}

func (t _TrueType) Equalp(n Node) bool {
	_, ok := n.(_TrueType)
	return ok
}

type _NullType struct{}

func (_NullType) PrintTo(w io.Writer) {
	io.WriteString(w, "nil")
}

func (_NullType) PrincTo(w io.Writer) {
	io.WriteString(w, "nil")
}

func (nt _NullType) Eval(*World) (Node, error) {
	return nt, nil
}

func (nt _NullType) Equals(n Node) bool {
	if n == nil {
		return true
	}
	_, ok := n.(_NullType)
	return ok
}

func (nt _NullType) Equalp(n Node) bool {
	return nt.Equals(n)
}

var Null Node = _NullType{}

type String string

func (s String) PrintTo(w io.Writer) {
	fmt.Fprintf(w, "\"%s\"", string(s))
}

func (s String) PrincTo(w io.Writer) {
	io.WriteString(w, string(s))
}

func (s String) Eval(*World) (Node, error) {
	return s, nil // errors.New("String can not be evaluate.")
}

func (s String) Equals(n Node) bool {
	ns, ok := n.(String)
	return ok && s == ns
}

func (s String) Equalp(n Node) bool {
	if ns, ok := n.(String); ok {
		return strings.EqualFold(string(s), string(ns))
	}
	return false
}

func (s String) Add(n Node) (Node, error) {
	if value, ok := n.(String); ok {
		return s + value, nil
	}
	return nil, fmt.Errorf("%w: `%s`", ErrNotSupportType, toString(n))
}

func (s String) LessThan(n Node) (bool, error) {
	if ns, ok := n.(String); ok {
		return s < ns, nil
	}
	return false, fmt.Errorf("%w: `%s`", ErrNotSupportType, toString(n))
}

type Symbol string

func (s Symbol) PrintTo(w io.Writer) {
	io.WriteString(w, string(s))
}

func (s Symbol) PrincTo(w io.Writer) {
	io.WriteString(w, string(s))
}

func (s Symbol) Eval(w *World) (Node, error) {
	return w.Get(string(s))
}

func (s Symbol) Equals(n Node) bool {
	ns, ok := n.(Symbol)
	return ok && s == ns
}

func (s Symbol) Equalp(n Node) bool {
	ns, ok := n.(Symbol)
	return ok && strings.EqualFold(string(s), string(ns))
}
