package gmnlisp

import (
	"context"
	"fmt"
)

// funGetCar implements (car X) and (setf (car X) Y)
func funGetCar(_ context.Context, _ *World, argv []Node) (Node, func(Node) error, error) {
	cons, ok := argv[0].(*Cons)
	if !ok {
		return nil, nil, fmt.Errorf("%w: %s", ErrExpectedCons, ToString(argv[0], PRINT))
	}
	return cons.Car, func(value Node) error { cons.Car = value; return nil }, nil
}

// funGetCdr implements (cdr X) and (setf (cdr X) Y)
func funGetCdr(_ context.Context, _ *World, argv []Node) (Node, func(Node) error, error) {
	cons, ok := argv[0].(*Cons)
	if !ok {
		return nil, nil, ErrExpectedCons
	}
	return cons.Cdr, func(value Node) error { cons.Cdr = value; return nil }, nil
}

// funList implements (list A B ...)
func funList(_ context.Context, _ *World, list []Node) (Node, error) {
	var cons Node = Null
	for i := len(list) - 1; i >= 0; i-- {
		cons = &Cons{
			Car: list[i],
			Cdr: cons,
		}
	}
	return cons, nil
}

func lastOfList(node Node) (*Cons, error) {
	for {
		cons, ok := node.(*Cons)
		if !ok {
			return nil, fmt.Errorf("%w `%s`", ErrExpectedCons, node)
		}
		if IsNull(cons.Cdr) {
			return cons, nil
		}
		node = cons.Cdr
	}
}

// funLast implements (last LIST)
func funLast(_ context.Context, _ *World, list []Node) (Node, error) {
	if IsNull(list[0]) {
		return Null, nil
	}
	tail, err := lastOfList(list[0])
	if err != nil {
		return nil, err
	}
	return tail.Car, nil
}

// funAppend implements (append LIST1 LIST2 ...)
func funAppend(_ context.Context, _ *World, list []Node) (Node, error) {
	var value Node
	var err error
	var buffer ListBuilder

	if len(list) <= 0 {
		return Null, nil
	}
	if len(list) == 1 {
		return list[0], nil
	}
	for _, list1 := range list[:len(list)-1] {
		for HasValue(list1) {
			value, list1, err = Shift(list1)
			if err != nil {
				return nil, err
			}
			buffer.Add(value)
		}
	}
	head := buffer.Sequence()
	tail, err := lastOfList(head)
	if err != nil {
		return nil, err
	}
	tail.Cdr = list[len(list)-1]
	return head, nil
}

// funCons implements (cons CAR CDR)
func funCons(_ context.Context, _ *World, argv []Node) (Node, error) {
	return &Cons{Car: argv[0], Cdr: argv[1]}, nil
}

// funLispp implements (listp VALUE)
func funListp(_ context.Context, _ *World, argv []Node) (Node, error) {
	if IsNull(argv[0]) {
		return True, nil
	}
	if _, ok := argv[0].(*Cons); ok {
		return True, nil
	}
	return Null, nil
}

// funAssoc implements (assoc KEY LIST)
func funAssoc(_ context.Context, _ *World, argv []Node) (Node, error) {
	key := argv[0]
	list := argv[1]

	for HasValue(list) {
		var element Node
		var err error

		element, list, err = Shift(list)
		if err != nil {
			return nil, err
		}
		cons, ok := element.(*Cons)
		if !ok {
			return nil, ErrExpectedCons
		}
		if key.Equals(cons.Car, STRICT) {
			return cons, nil
		}
	}
	return Null, nil
}

func subst(newItem, oldItem, list Node) Node {
	if list.Equals(oldItem, STRICT) {
		return newItem
	}
	cons, ok := list.(*Cons)
	if ok {
		return &Cons{
			Car: subst(newItem, oldItem, cons.Car),
			Cdr: subst(newItem, oldItem, cons.Cdr),
		}
	}
	return list
}

// funSubst implements (subst NEWITEM OLDITEM LIST)
func funSubst(_ context.Context, _ *World, argv []Node) (Node, error) {
	return subst(argv[0], argv[1], argv[2]), nil
}

// funReplaca implements (replaca X Y) == (setf (car X) Y)
func funReplaca(_ context.Context, _ *World, argv []Node) (Node, error) {
	cons, ok := argv[0].(*Cons)
	if !ok {
		return nil, ErrExpectedCons
	}
	cons.Car = argv[1]
	return cons, nil
}

// funReplacd implements (replacd X Y) == (setf (cdr X) Y)
func funReplacd(_ context.Context, _ *World, argv []Node) (Node, error) {
	cons, ok := argv[0].(*Cons)
	if !ok {
		return nil, ErrExpectedCons
	}
	cons.Cdr = argv[1]
	return cons, nil
}
