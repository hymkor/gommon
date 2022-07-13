package gmnlisp

import (
	"errors"
	"fmt"
	"io"
)

type Lambda struct {
	param []string
	code  Node
	name  string
}

func cmdLambda(_ *World, node Node) (Node, error) {
	return newLambda(node, "")
}

func newLambda(node Node, blockName string) (Node, error) {
	// (lambda (param) code)

	cons, ok := node.(*Cons)
	if !ok {
		return nil, fmt.Errorf("%w for parameter list", ErrExpectedCons)
	}
	params := []string{}
	if err := forEachWithoutEval(cons.Car, func(n Node) error {
		name, ok := n.(Symbol)
		if !ok {
			return ErrExpectedSymbol
		}
		params = append(params, string(name))
		return nil
	}); err != nil {
		return nil, err
	}

	return &Lambda{
		param: params,
		code:  cons.GetCdr(),
		name:  blockName,
	}, nil
}

func (nl *Lambda) PrintTo(w io.Writer) {
	nl.prinX(w, true)
}

func (nl *Lambda) PrincTo(w io.Writer) {
	nl.prinX(w, false)
}

func (nl *Lambda) prinX(w io.Writer, rich bool) {
	io.WriteString(w, "(lambda (")
	dem := ""
	for _, name := range nl.param {
		io.WriteString(w, dem)
		io.WriteString(w, name)
		dem = " "
	}
	io.WriteString(w, ") ")
	if cons, ok := nl.code.(*Cons); ok {
		cons.writeToWithoutKakko(w, rich)
	} else {
		if rich {
			nl.code.PrintTo(w)
		} else {
			nl.code.PrincTo(w)
		}
	}
	io.WriteString(w, ")")
}

func (nl *Lambda) Call(w *World, n Node) (Node, error) {
	globals := map[string]Node{}
	for _, name := range nl.param {
		if cons, ok := n.(*Cons); ok {
			var err error

			globals[name], err = cons.GetCar().Eval(w)
			if err != nil {
				return nil, err
			}
			n = cons.GetCdr()
		} else {
			globals[name] = Null
		}
	}
	w.nameSpace = &_NameSpace{
		globals: globals,
		parent:  w.nameSpace,
	}
	defer func() {
		w.nameSpace = w.nameSpace.parent
	}()

	var errEarlyReturns *ErrEarlyReturns

	result, err := progn(w, nl.code)
	if errors.As(err, &errEarlyReturns) && errEarlyReturns.Name == string(nl.name) {
		return errEarlyReturns.Value, nil
	}
	return result, err
}

func (nl *Lambda) Eval(*World) (Node, error) {
	return nl, nil
}

func (*Lambda) Equals(Node) bool {
	return false
}

func (*Lambda) EqualP(Node) bool {
	return false
}

func cmdDefun(w *World, node Node) (Node, error) {
	cons, ok := node.(*Cons)
	if !ok {
		return nil, ErrExpectedCons
	}
	_name, ok := cons.Car.(Symbol)
	if !ok {
		return nil, ErrExpectedSymbol
	}
	name := string(_name)

	lambda, err := newLambda(cons.Cdr, string(_name))
	if err != nil {
		return nil, err
	}
	w.Set(name, lambda)
	return lambda, nil
}
