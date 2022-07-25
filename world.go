package gmnlisp

import (
	"context"
	"fmt"
	"io"
	"os"
)

type World struct {
	parent  *World
	globals map[string]Node
}

type Writer struct {
	_Dummy
	io.Writer
}

func (w *World) each(f func(string, Node) bool) {
	for w != nil {
		for name, value := range w.globals {
			if !f(name, value) {
				return
			}
		}
		w = w.parent
	}
}

func (w *World) Get(name string) (Node, error) {
	for w != nil {
		if value, ok := w.globals[name]; ok {
			return value, nil
		}
		w = w.parent
	}
	return Null, fmt.Errorf("%w `%s`", ErrVariableUnbound, name)
}

func (w *World) Set(name string, value Node) {
	for w != nil {
		if _, ok := w.globals[name]; ok || w.parent == nil {
			w.globals[name] = value
			return
		}
		w = w.parent
	}
}

const standardOutput = "*standard-output*"

func (w *World) Stdout() (io.Writer, error) {
	stdout, err := w.Get(standardOutput)
	if err != nil {
		return nil, err
	}
	_stdout, ok := stdout.(io.Writer)
	if !ok {
		return nil, ErrExpectedWriter
	}
	return _stdout, nil
}

func (w *World) SetStdout(writer io.Writer) {
	w.Set(standardOutput, Writer{Writer: writer})
}

func New() *World {
	return &World{
		globals: map[string]Node{
			standardOutput:        &Writer{Writer: os.Stdout},
			"*":                   Function(cmdMulti),
			"+":                   Function(cmdAdd),
			"-":                   Function(cmdSub),
			"--get-all-symbols--": Function(cmdGetAllSymbols),
			"/":                   Function(cmdDevide),
			"/=":                  Function(cmdNotEqual),
			"<":                   Function(cmdLessThan),
			"<=":                  Function(cmdLessOrEqual),
			"=":                   Function(cmdEqualOp),
			">":                   Function(cmdGreaterThan),
			">=":                  Function(cmdGreaterOrEqual),
			"T":                   True,
			"and":                 Function(cmdAnd),
			"append":              Function(cmdAppend),
			"atom":                Function(cmdAtom),
			"block":               Function(cmdBlock),
			"cadr":                Function(cmdCadr),
			"caddr":               Function(cmdCaddr),
			"cadddr":              Function(cmdCadddr),
			"car":                 Function(cmdCar),
			"cdr":                 Function(cmdCdr),
			"cddr":                Function(cmdCddr),
			"cdddr":               Function(cmdCdddr),
			"close":               Function(cmdClose),
			"command":             Function(cmdCommand),
			"cond":                Function(cmdCond),
			"cons":                Function(cmdCons),
			"defmacro":            Function(cmdDefMacro),
			"defun":               Function(cmdDefun),
			"equal":               Function(cmdEqual),
			"equalp":              Function(cmdEqualOp),
			"exit":                Function(cmdQuit),
			"foreach":             Function(cmdForeach),
			"function":            Function(cmdFunction),
			"funcall":             Function(cmdFunCall),
			"if":                  Function(cmdIf),
			"lambda":              Function(cmdLambda),
			"length":              Function(cmdLength),
			"let":                 Function(cmdLet),
			"list":                Function(cmdList),
			"listp":               Function(cmdListp),
			"load":                Function(cmdLoad),
			"macroexpand":         Function(cmdMacroExpand),
			"mapcar":              Function(cmdMapCar),
			"member":              Function(cmdMember),
			"nil":                 Null,
			"not":                 Function(cmdNot),
			"open":                Function(cmdOpen),
			"or":                  Function(cmdOr),
			"parse-integer":       Function(cmdParseInt),
			"prin1":               Function(cmdPrin1),
			"princ":               Function(cmdPrinc),
			"print":               Function(cmdPrint),
			"progn":               Function(cmdProgn),
			"quit":                Function(cmdQuit),
			"quote":               Function(cmdQuote),
			"read":                Function(cmdRead),
			"read-line":           Function(cmdReadLine),
			"return":              Function(cmdReturn),
			"return-from":         Function(cmdReturnFrom),
			"setq":                Function(cmdSetq),
			"strcase":             Function(cmdStrCase),
			"strcat":              Function(cmdStrCat),
			"strlen":              Function(cmdStrLen),
			"substr":              Function(cmdSubStr),
			"terpri":              Function(cmdTerpri),
			"trace":               Function(cmdTrace),
			"truncate":            Function(cmdTruncate),
			"while":               Function(cmdWhile),
			"write":               Function(cmdWrite),
			"write-line":          Function(cmdWriteLine),
		},
	}
}

func (w *World) evalListAll(ctx context.Context, list Node, result []Node) error {
	if err := listToArray(list, result); err != nil {
		return err
	}
	for i := 0; i < len(result); i++ {
		value, err := result[i].Eval(ctx, w)
		if err != nil {
			return err
		}
		result[i] = value
	}
	return nil
}

func (w *World) evalListToSlice(ctx context.Context, list Node) ([]Node, error) {
	result, err := listToSlice(list)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(result); i++ {
		value, err := result[i].Eval(ctx, w)
		if err != nil {
			return nil, err
		}
		result[i] = value
	}
	return result, nil
}

func (w *World) shiftAndEvalCar(ctx context.Context, list Node) (Node, Node, error) {
	cons, ok := list.(*Cons)
	if !ok {
		return nil, nil, ErrTooFewArguments
	}
	value, err := cons.GetCar().Eval(ctx, w)
	if err != nil {
		return nil, nil, err
	}
	return value, cons.Cdr, nil
}

func (w *World) inject(ctx context.Context, list Node, f func(left, right Node) (Node, error)) (Node, error) {
	result, list, err := w.shiftAndEvalCar(ctx, list)
	if err != nil {
		return nil, err
	}
	for HasValue(list) {
		var next Node
		var err error

		next, list, err = w.shiftAndEvalCar(ctx, list)
		if err != nil {
			return nil, err
		}
		result, err = f(result, next)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (w *World) InterpretNodes(ctx context.Context, ns []Node) (Node, error) {
	var result Node = Null
	var err error

	for _, c := range ns {
		result, err = c.Eval(ctx, w)
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

func (w *World) Interpret(ctx context.Context, code string) (Node, error) {
	compiled, err := ReadString(code)
	if err != nil {
		return nil, err
	}
	return w.InterpretNodes(ctx, compiled)
}

func (w *World) InterpretBytes(ctx context.Context, code []byte) (Node, error) {
	compiled, err := ReadBytes(code)
	if err != nil {
		return nil, err
	}
	return w.InterpretNodes(ctx, compiled)
}

func (w *World) Call(ctx context.Context, f Node, params ...Node) (Node, error) {
	_f, ok := f.(_Callable)
	if !ok {
		return nil, ErrExpectedFunction
	}
	return _f.Call(ctx, w, List(params...))
}
