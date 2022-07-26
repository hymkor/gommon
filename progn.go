package gmnlisp

import (
	"context"
	"errors"
	"fmt"
)

type ErrEarlyReturns struct {
	Value Node
	Name  string
}

func (e *ErrEarlyReturns) Error() string {
	if e.Name == "" {
		return "Unexpected (return)"
	}
	return fmt.Sprintf("Unexpected (return-from %s)", e.Name)
}

func cmdReturn(ctx context.Context, w *World, argv []Node) (Node, error) {
	return nil, &ErrEarlyReturns{Value: argv[0], Name: ""}
}

func cmdReturnFrom(ctx context.Context, w *World, n Node) (Node, error) {
	var argv [2]Node
	if err := listToArray(n, argv[:]); err != nil {
		return nil, err
	}
	symbol, ok := argv[0].(Symbol)
	if !ok {
		return nil, ErrExpectedSymbol
	}
	value, err := argv[1].Eval(ctx, w)
	if err != nil {
		return nil, err
	}
	return nil, &ErrEarlyReturns{Value: value, Name: string(symbol)}
}

func progn(ctx context.Context, w *World, n Node) (value Node, err error) {
	for HasValue(n) {
		value, n, err = w.shiftAndEvalCar(ctx, n)
		if err != nil {
			return nil, err
		}
	}
	return
}

func cmdProgn(ctx context.Context, w *World, c Node) (Node, error) {
	return progn(ctx, w, c)
}

func cmdBlock(ctx context.Context, w *World, node Node) (Node, error) {
	nameNode, statements, err := shift(node)
	if err != nil {
		return nil, err
	}
	nameSymbol, ok := nameNode.(Symbol)
	if !ok {
		return nil, ErrExpectedSymbol
	}
	name := string(nameSymbol)

	var errEarlyReturns *ErrEarlyReturns
	rv, err := progn(ctx, w, statements)
	if errors.As(err, &errEarlyReturns) && errEarlyReturns.Name == name {
		return errEarlyReturns.Value, nil
	}
	return rv, err
}

func cmdCond(ctx context.Context, w *World, list Node) (Node, error) {
	for HasValue(list) {
		var condAndAct Node
		var err error

		condAndAct, list, err = shift(list)
		if err != nil {
			return nil, err
		}
		cond, act, err := w.shiftAndEvalCar(ctx, condAndAct)
		if err != nil {
			return nil, err
		}
		if HasValue(cond) {
			return progn(ctx, w, act)
		}
	}
	return Null, nil
}

func cmdIf(ctx context.Context, w *World, params Node) (Node, error) {
	cond, params, err := w.shiftAndEvalCar(ctx, params)
	if err != nil {
		return nil, err
	}
	thenClause, params, err := shift(params)
	if err != nil {
		return nil, err
	}
	var elseClause Node = Null
	if HasValue(params) {
		elseClause, params, err = shift(params)
		if err != nil {
			return nil, err
		}
		if HasValue(params) {
			return nil, ErrTooManyArguments
		}
	}
	if HasValue(cond) {
		return thenClause.Eval(ctx, w)
	} else if HasValue(elseClause) {
		return elseClause.Eval(ctx, w)
	} else {
		return Null, nil
	}
}

func cmdForeach(ctx context.Context, w *World, args Node) (Node, error) {
	var _symbol Node
	var err error

	_symbol, args, err = shift(args)
	if err != nil {
		return nil, err
	}
	symbol, ok := _symbol.(Symbol)
	if !ok {
		return nil, ErrExpectedSymbol
	}

	var list Node
	var code Node
	list, code, err = w.shiftAndEvalCar(ctx, args)
	if err != nil {
		return nil, err
	}

	var last Node
	for HasValue(list) {
		var value Node

		value, list, err = shift(list)
		if err != nil {
			return nil, err
		}
		w.Set(string(symbol), value)

		last, err = progn(ctx, w, code)
		if err != nil {
			return nil, err
		}
	}
	return last, nil
}

func cmdWhile(ctx context.Context, w *World, n Node) (Node, error) {
	cond, statements, err := shift(n)
	if err != nil {
		return nil, err
	}
	var last Node = Null
	for {
		if err := checkContext(ctx); err != nil {
			return nil, err
		}
		cont, err := cond.Eval(ctx, w)
		if err != nil {
			return nil, err
		}
		if IsNull(cont) {
			return last, nil
		}
		last, err = progn(ctx, w, statements)
		if err != nil {
			return nil, err
		}
	}
}

func cmdQuit(context.Context, *World, Node) (Node, error) {
	return Null, ErrQuit
}
