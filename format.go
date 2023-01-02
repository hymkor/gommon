package gmnlisp

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type runeWriter interface {
	io.Writer
	WriteRune(r rune) (int, error)
}

func printInt(w io.Writer, value Node, base int, args ...int) error {
	width := -1
	padding := -1
	if argc := len(args); argc >= 3 {
		return MakeError(ErrTooManyArguments, "printInt")
	} else if argc == 2 {
		width = args[0]
		padding = args[1]
	} else if argc == 1 {
		width = args[0]
		padding = ' '
	}

	var body string
	if d, ok := value.(Integer); ok {
		body = strconv.FormatInt(int64(d), base)
	} else if f, ok := value.(Float); ok {
		body = strconv.FormatInt(int64(f), base)
	} else {
		return ErrNotSupportType
	}
	if len(body) < width {
		if padding <= 0 {
			padding = ' '
		}
		for i := len(body); i < width; i++ {
			w.Write([]byte{byte(padding)})
		}
	}
	io.WriteString(w, strings.ToUpper(body))
	return nil
}

func funFormatObject(_ context.Context, w *World, list []Node) (Node, error) {
	return tAndNilToWriter(w, list, func(writer runeWriter, list []Node) error {
		var err error
		if IsNull(list[1]) { // ~a (AS-IS)
			_, err = list[0].PrintTo(writer, PRINC)
		} else { // ~s (S expression)
			_, err = list[0].PrintTo(writer, PRINT)
		}
		return err
	})
}

func funFormatChar(_ context.Context, w *World, list []Node) (Node, error) {
	return tAndNilToWriter(w, list, func(writer runeWriter, list []Node) error {
		r, ok := list[0].(Rune)
		if !ok {
			return MakeError(ErrExpectedCharacter, list[0])
		}
		_, err := writer.WriteRune(rune(r))
		return err
	})
}

func funFormatInteger(_ context.Context, w *World, _args []Node) (Node, error) {
	return tAndNilToWriter(w, _args, func(writer runeWriter, args []Node) error {
		radix, ok := args[1].(Integer)
		if !ok {
			return MakeError(ErrExpectedNumber, args[1])
		}
		return printInt(writer, args[0], int(radix))
	})
}

func printFloat(w runeWriter, value Node, mark byte, args ...int) error {
	width := -1
	prec := -1
	if argc := len(args); argc >= 3 {
		return MakeError(ErrTooManyArguments, "printFloat")
	} else if argc == 2 {
		width = args[0]
		prec = args[1]
	} else if argc == 1 {
		width = args[0]
	}
	var body string
	if d, ok := value.(Integer); ok {
		body = strconv.FormatFloat(float64(d), mark, prec, 64)
	} else if f, ok := value.(Float); ok {
		body = strconv.FormatFloat(float64(f), mark, prec, 64)
	} else {
		return ErrNotSupportType
	}
	if len(body) < width {
		for i := len(body); i < width; i++ {
			w.WriteRune(' ')
		}
	}
	io.WriteString(w, body)
	return nil
}

func funFormatFloat(_ context.Context, w *World, args []Node) (Node, error) {
	return tAndNilToWriter(w, args, func(_writer runeWriter, args []Node) error {
		return printFloat(_writer, args[0], 'f')
	})
}

func printSpaces(n int, w io.Writer) {
	for n > 0 {
		w.Write([]byte{' '})
		n--
	}
}

func formatSub(w runeWriter, argv []Node) error {
	format, ok := argv[0].(String)
	if !ok {
		return MakeError(ErrExpectedString, argv[0])
	}
	argv = argv[1:]

	for ok && HasValue(format) {
		var c Rune

		c, format, ok = format.firstRuneAndRestString()
		if !ok {
			break
		}
		if c != '~' {
			c.PrintTo(w, PRINC)
			continue
		}
		if IsNull(format) {
			w.WriteRune('~')
			break
		}
		c, format, ok = format.firstRuneAndRestString()
		if !ok {
			break
		}
		if c == '~' {
			w.Write([]byte{'~'})
			continue
		}
		parameter := []int{}
		for {
			if decimal := strings.IndexByte("0123456789", byte(c)); decimal >= 0 {
				for {
					if IsNull(format) {
						return ErrInvalidFormat
					}
					c, format, ok = format.firstRuneAndRestString()
					if !ok {
						return ErrInvalidFormat
					}
					d := strings.IndexByte("0123456789", byte(c))
					if d < 0 {
						parameter = append(parameter, decimal)
						break
					}
					decimal = decimal*10 + d
				}
			} else if c == '\'' {
				if IsNull(format) {
					return ErrInvalidFormat
				}
				c, format, ok = format.firstRuneAndRestString()
				if !ok {
					return ErrInvalidFormat
				}
				parameter = append(parameter, int(c))
				if IsNull(format) {
					return ErrInvalidFormat
				}
				c, format, ok = format.firstRuneAndRestString()
			} else if c == 'v' || c == 'V' {
				if len(argv) < 1 {
					return ErrTooFewArguments
				}
				decimal, ok := argv[0].(Integer)
				if !ok {
					return ErrExpectedNumber
				}
				parameter = append(parameter, int(decimal))
			} else if c == '#' {
				parameter = append(parameter, int(len(argv)))
			} else {
				break
			}
			if c != ',' {
				break
			}
			if IsNull(format) {
				return ErrInvalidFormat
			}
			c, format, ok = format.firstRuneAndRestString()
			if !ok {
				break
			}
		}

		if c == '%' || c == '&' {
			if len(parameter) >= 1 {
				for n := parameter[0]; n >= 1; n-- {
					w.Write([]byte{'\n'})
				}
			} else {
				w.Write([]byte{'\n'})
			}
			continue
		}
		if len(argv) <= 0 {
			return ErrTooFewArguments
		}

		value := argv[0]
		argv = argv[1:]

		var err error
		switch c {
		case 'd':
			err = printInt(w, value, 10, parameter...)
		case 'x':
			err = printInt(w, value, 16, parameter...)
		case 'o':
			err = printInt(w, value, 8, parameter...)
		case 'b':
			err = printInt(w, value, 2, parameter...)
		case 'f':
			err = printFloat(w, value, 'f', parameter...)
		case 'e':
			err = printFloat(w, value, 'e', parameter...)
		case 'g':
			err = printFloat(w, value, 'g', parameter...)
		case 'a':
			n, err := value.PrintTo(w, PRINC)
			if err != nil {
				return err
			}
			if len(parameter) >= 1 {
				printSpaces(parameter[0]-n, w)
			}
		case 's':
			n, err := value.PrintTo(w, PRINT)
			if err != nil {
				return err
			}
			if len(parameter) >= 1 {
				printSpaces(parameter[0]-n, w)
			}
		default:
			err = fmt.Errorf("Not support code '%c'", c)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func tAndNilToWriter(w *World, argv []Node, f func(runeWriter, []Node) error) (Node, error) {
	if output, ok := argv[0].(io.Writer); ok {
		rw, ok := output.(runeWriter)
		if !ok {
			bw := bufio.NewWriter(output)
			defer bw.Flush()
			rw = bw
		}
		err := f(rw, argv[1:])
		return Null, err
	}
	if IsNull(argv[0]) {
		var buffer strings.Builder
		err := f(&buffer, argv[1:])
		return String(buffer.String()), err
	}
	if True.Equals(argv[0], STRICT) {
		rw, ok := w.shared.stdout._Writer.(runeWriter)
		if !ok {
			bw := bufio.NewWriter(w.shared.stdout._Writer)
			defer bw.Flush()
			rw = bw
		}
		err := f(rw, argv[1:])
		return Null, err
	}
	return nil, MakeError(ErrExpectedWriter, argv[0])
}

func funFormat(ctx context.Context, w *World, argv []Node) (Node, error) {
	return tAndNilToWriter(w, argv, formatSub)
}
