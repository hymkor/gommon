package common

import (
	"context"
	"io"
	_ "embed"

	. "github.com/hymkor/gmnlisp"
)

func init() {
	for key, val := range embedFunctions {
		Functions.Set(key, val)
	}
}

var Functions = Variables{
	NewSymbol("member"):       &KWFunction{C: 2, F: funMember},
	NewSymbol("find"):         &KWFunction{C: 2, F: funFind},
	NewSymbol("position"):     &KWFunction{C: 2, F: funPosition},
	NewSymbol("defparameter"): SpecialF(cmdDefparameter),
	NewSymbol("defvar"):       SpecialF(cmdDefvar),
	NewSymbol("prin1"):        &Function{C: 1, F: funPrin1},
	NewSymbol("princ"):        &Function{C: 1, F: funPrinc},
	NewSymbol("print"):        &Function{C: 1, F: funPrint},
	NewSymbol("terpri"):       &Function{C: -1, F: funTerpri},
	NewSymbol("write"):        &KWFunction{C: 1, F: funWrite},
	NewSymbol("write-line"):   SpecialF(cmdWriteLine),
	NewSymbol("when"):         SpecialF(cmdWhen),
	NewSymbol("unless"):       SpecialF(cmdUnless),
	NewSymbol("nth"):          &Function{C: 2, F: funNth},
	NewSymbol("nthcdr"):       &Function{C: 2, F: funNthcdr},
	NewSymbol("cadddr"):       &Function{C: 1, F: funCadddr},
	NewSymbol("caddr"):        &Function{C: 1, F: funCaddr},
	NewSymbol("cadr"):         &Function{C: 1, F: funCadr},
	NewSymbol("cdddr"):        &Function{C: 1, F: funCdddr},
	NewSymbol("cddr"):         &Function{C: 1, F: funCddr},
	//NewSymbol("first"):                       &Function{C: 1, F: funGetCar},
	NewSymbol("second"):         &Function{C: 1, F: funCadr},
	NewSymbol("third"):          &Function{C: 1, F: funCaddr},
	NewSymbol("concatenate"):    defConcatenate,
	NewSymbol("dolist"):         SpecialF(cmdDoList),
	NewSymbol("dotimes"):        SpecialF(cmdDoTimes),
	NewSymbol("with-open-file"): SpecialF(cmdWithOpenFile),
	NewSymbol("coerce"):         &Function{C: 2, F: funCoerce},
	NewSymbol("map"):            &Function{C: -1, F: funMap},
	NewSymbol("typep"):          &Function{C: 2, F: funTypep},
	NewSymbol("subst"):          &Function{C: 3, F: funSubst},
}

var defparameter Callable

func cmdDefparameter(ctx context.Context, w *World, list Node) (Node, error) {
	if IsNull(defparameter) {
		_defparameter, err := w.Get(NewSymbol("defglobal"))
		if err != nil {
			return nil, err
		}
		var ok bool
		defparameter, ok = _defparameter.(Callable)
		if !ok {
			return nil, ErrExpectedFunction
		}
	}
	return defparameter.Call(ctx, w, list)
}

func cmdDefvar(ctx context.Context, w *World, list Node) (Node, error) {
	// from CommonLisp
	var symbolNode Node
	var err error

	symbolNode, list, err = Shift(list)
	if err != nil {
		return nil, err
	}
	symbol, ok := symbolNode.(Symbol)
	if !ok {
		return nil, ErrExpectedSymbol
	}
	w.DefineVariable(symbol, func() Node {
		var value Node = Null
		if HasValue(list) {
			value, _, err = w.ShiftAndEvalCar(ctx, list)
		}
		return value
	})
	return symbol, err
}

func cmdWhen(ctx context.Context, w *World, args Node) (Node, error) {
	cond, args, err := w.ShiftAndEvalCar(ctx, args)
	if err != nil {
		return nil, err
	}
	if IsNull(cond) {
		return Null, nil
	}
	return Progn(ctx, w, args)
}

func cmdUnless(ctx context.Context, w *World, args Node) (Node, error) {
	cond, args, err := w.ShiftAndEvalCar(ctx, args)
	if err != nil {
		return nil, err
	}
	if HasValue(cond) {
		return Null, nil
	}
	return Progn(ctx, w, args)
}

type KWFunction struct {
	C int
	F func(context.Context, *World, []Node, map[Keyword]Node) (Node, error)
}

func (*KWFunction) PrintTo(w io.Writer, m PrintMode) (int, error) {
	return io.WriteString(w, "buildin function")
}

func (f *KWFunction) Eval(context.Context, *World) (Node, error) {
	return f, nil
}

func (f *KWFunction) Equals(n Node, m EqlMode) bool {
	return false
}

func ListToKwargs(ctx context.Context, w *World, list Node) ([]Node, map[Keyword]Node, error) {
	args := []Node{}
	kwargs := map[Keyword]Node{}
	for HasValue(list) {
		var tmp Node
		var err error

		tmp, list, err = w.ShiftAndEvalCar(ctx, list)
		if err != nil {
			return nil, nil, err
		}
		if keyword, ok := tmp.(Keyword); ok {
			tmp, list, err = w.ShiftAndEvalCar(ctx, list)
			if err != nil {
				return nil, nil, err
			}
			kwargs[keyword] = tmp
		} else {
			args = append(args, tmp)
		}
	}
	return args, kwargs, nil
}

func (f *KWFunction) Call(ctx context.Context, w *World, list Node) (Node, error) {
	if err := CheckContext(ctx); err != nil {
		return nil, err
	}
	args, kwargs, err := ListToKwargs(ctx, w, list)
	if err != nil {
		return nil, err
	}
	if f.C >= 0 {
		if len(args) < f.C {
			return nil, ErrTooFewArguments
		}
		if len(args) > f.C {
			return nil, ErrTooManyArguments
		}
	}
	return f.F(ctx, w, args, kwargs)
}
