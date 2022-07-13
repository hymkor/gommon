package gmnlisp

func cmdCons(w *World, node Node) (Node, error) {
	first, rest, err := w.shiftAndEvalCar(node)
	if err != nil {
		return nil, err
	}
	second, rest, err := w.shiftAndEvalCar(rest)
	if err != nil {
		return nil, err
	}
	if HasValue(rest) {
		return nil, ErrTooFewOrTooManyArguments
	}
	return &Cons{Car: first, Cdr: second}, err
}

func cmdCar(w *World, param Node) (Node, error) {
	first, _, err := w.shiftAndEvalCar(param)
	if err != nil {
		return nil, err
	}
	cons, ok := first.(*Cons)
	if !ok {
		return nil, ErrExpectedCons
	}
	return cons.Car, nil
}

func cmdCdr(w *World, param Node) (Node, error) {
	first, _, err := w.shiftAndEvalCar(param)
	if err != nil {
		return nil, err
	}
	cons, ok := first.(*Cons)
	if !ok {
		return nil, ErrExpectedCons
	}
	return cons.Cdr, nil
}

func cmdQuote(_ *World, param Node) (Node, error) {
	cons, ok := param.(*Cons)
	if !ok {
		return nil, ErrTooFewOrTooManyArguments
	}
	return cons.Car, nil
}

func cmdAtom(_ *World, param Node) (Node, error) {
	cons, ok := param.(*Cons)
	if !ok {
		return nil, ErrExpectedCons
	}
	if _, ok := cons.Car.(*Cons); ok {
		return Null, nil
	}
	return True, nil
}

func cmdEqual(w *World, param Node) (Node, error) {
	first, rest, err := w.shiftAndEvalCar(param)
	if err != nil {
		return nil, err
	}
	for HasValue(rest) {
		var next Node

		next, rest, err = w.shiftAndEvalCar(rest)
		if err != nil {
			return nil, err
		}
		if !first.Equals(next) {
			return Null, nil
		}
	}
	return True, nil
}

func cmdList(w *World, node Node) (Node, error) {
	car, rest, err := w.shiftAndEvalCar(node)
	if err != nil {
		return nil, err
	}
	var cdr Node

	if IsNull(rest) {
		cdr = Null
	} else {
		cdr, err = cmdList(w, rest)
		if err != nil {
			return nil, err
		}
	}
	return &Cons{Car: car, Cdr: cdr}, nil
}
