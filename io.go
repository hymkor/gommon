package gmnlisp

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type _Dummy struct{}

func (d _Dummy) Eval(context.Context, *World) (Node, error) {
	return d, nil
}

func (d _Dummy) Equals(Node, EqlMode) bool {
	return false
}

func (d _Dummy) PrintTo(w io.Writer, m PrintMode) {
	io.WriteString(w, "(binary)")
}

func openAsRead(fname string) (Node, error) {
	type Reader struct {
		_Dummy
		*bufio.Reader
		io.Closer
	}
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	return &Reader{
		Reader: bufio.NewReader(file),
		Closer: file,
	}, nil
}

func openAsWrite(fname string) (Node, error) {
	type Writer struct {
		_Dummy
		io.WriteCloser
	}
	file, err := os.Create(fname)
	if err != nil {
		return nil, err
	}
	return &Writer{
		WriteCloser: file,
	}, nil
}

func cmdOpen(ctx context.Context, w *World, n Node) (Node, error) {
	fname, n, err := w.shiftAndEvalCar(ctx, n)
	if err != nil {
		return nil, err
	}
	_fname, ok := fname.(String)
	if !ok {
		return nil, fmt.Errorf("%w `%s`", ErrExpectedString, toString(_fname))
	}
	if IsNull(n) {
		return openAsRead(string(_fname))
	}

	mode, n, err := w.shiftAndEvalCar(ctx, n)
	if HasValue(n) {
		return nil, ErrTooManyArguments
	}
	_mode, ok := mode.(String)
	if !ok {
		return nil, fmt.Errorf("%w `%s`", ErrExpectedString, toString(_mode))
	}

	var result Node
	switch string(_mode) {
	case "r":
		result, err = openAsRead(string(_fname))
	case "w":
		result, err = openAsWrite(string(_fname))
	default:
		return nil, fmt.Errorf("no such a option `%s`", string(_mode))
	}
	if errors.Is(err, os.ErrNotExist) {
		return Null, nil
	}
	return result, err
}

func chomp(s string) string {
	L := len(s)
	if L > 0 && s[L-1] == '\n' {
		L--
		s = s[:L]
		if L > 0 && s[L-1] == '\r' {
			s = s[:L-1]
		}
	}
	return s
}

func cmdReadLine(ctx context.Context, w *World, n Node) (Node, error) {
	type ReadStringer interface {
		ReadString(byte) (string, error)
	}

	var argv [1]Node
	if err := w.evalListAll(ctx, n, argv[:]); err != nil {
		return nil, err
	}
	r, ok := argv[0].(ReadStringer)
	if !ok {
		return nil, fmt.Errorf("Expected Reader `%s`", toString(argv[0]))
	}
	s, err := r.ReadString('\n')
	if err == io.EOF {
		return Null, nil
	}
	return String(chomp(s)), err
}

func cmdClose(ctx context.Context, w *World, n Node) (Node, error) {
	var argv [1]Node
	if err := w.evalListAll(ctx, n, argv[:]); err != nil {
		return nil, err
	}
	c, ok := argv[0].(io.Closer)
	if !ok {
		return nil, fmt.Errorf("Expected Closer `%s`", toString(argv[0]))
	}
	return Null, c.Close()
}

func cmdCommand(ctx context.Context, w *World, n Node) (Node, error) {
	argv, err := w.evalListToSlice(ctx, n)
	if err != nil {
		return nil, err
	}
	if len(argv) < 1 {
		return nil, ErrTooFewArguments
	}
	_argv := make([]string, 0, len(argv))
	var buffer strings.Builder
	for _, n := range argv {
		n.PrintTo(&buffer, PRINC)
		_argv = append(_argv, buffer.String())
		buffer.Reset()
	}
	cmd := exec.Command(_argv[0], _argv[1:]...)
	cmd.Stdout, err = w.Stdout()
	if err != nil {
		return nil, err
	}
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return Null, cmd.Run()
}
