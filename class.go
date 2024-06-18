package gmnlisp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
)

type _SlotSpec struct {
	identifier Symbol
	reader     []Symbol
	writer     []Symbol
	accessor   []Symbol
	boundp     []Symbol
	initform   func() (Node, error)
	initarg    []Symbol
}

type _Getter struct {
	Symbol
	class map[Symbol]func(*_Receiver) (Node, error)
}

func registerGetter(w *World, getterName, className Symbol, getter func(*_Receiver) (Node, error)) error {
	if _acc, err := w.Get(getterName); err == nil {
		if acc, ok := _acc.(*_Getter); ok {
			acc.class[className] = getter
		} else {
			return fmt.Errorf("%v: already defined as not accessor", getterName)
		}
	} else {
		w.DefineGlobal(getterName, &_Getter{
			Symbol: getterName,
			class: map[Symbol]func(*_Receiver) (Node, error){
				className: getter,
			},
		})
	}
	return nil
}

func (acc *_Getter) findClass(class *_UserClass) func(*_Receiver) (Node, error) {
	if f, ok := acc.class[class.Symbol]; ok {
		return f
	}
	for _, super := range class.Super {
		if f := acc.findClass(super); f != nil {
			return f
		}
	}
	return nil
}

func (acc *_Getter) Call(ctx context.Context, w *World, node Node) (Node, error) {
	_this, _, err := w.ShiftAndEvalCar(ctx, node)
	if err != nil {
		return nil, err
	}
	reciever, ok := _this.(*_Receiver)
	if !ok {
		return nil, errors.New("Expect Class Instance")
	}
	f := acc.findClass(reciever._UserClass)
	if f == nil {
		return nil, fmt.Errorf("reciever %v not found in %v", acc.Symbol.String(), reciever._UserClass.Symbol.String())
	}
	return f(reciever)
}

type _Setter struct {
	Symbol
	class map[Symbol]func(*_Receiver, Node)
}

func registerSetter(w *World, setterName, className Symbol, setter func(*_Receiver, Node)) error {
	if _acc, err := w.Get(setterName); err == nil {
		// println("accessor is found", spec.accessor.String())
		if acc, ok := _acc.(*_Setter); ok {
			acc.class[className] = setter
		} else {
			return fmt.Errorf("%v: already defined as not accessor", setterName.String())
		}
	} else {
		// println("accessor not found", spec.accessor.String())
		w.DefineGlobal(setterName, &_Setter{
			Symbol: setterName,
			class: map[Symbol]func(*_Receiver, Node){
				className: setter,
			},
		})
	}
	return nil
}

func (acc *_Setter) findClass(class *_UserClass) func(*_Receiver, Node) {
	if f := acc.class[class.Symbol]; f != nil {
		return f
	}
	for _, super := range class.Super {
		if f := acc.findClass(super); f != nil {
			return f
		}
	}
	return nil
}

func (acc *_Setter) Call(ctx context.Context, w *World, node Node) (Node, error) {
	value, node, err := w.ShiftAndEvalCar(ctx, node)
	if err != nil {
		return nil, err
	}
	_instance, _, err := w.ShiftAndEvalCar(ctx, node)
	if err != nil {
		return nil, err
	}
	reciever, ok := _instance.(*_Receiver)
	if !ok {
		return nil, errors.New("Expect Class Instance")
	}
	f := acc.findClass(reciever._UserClass)
	if f == nil {
		return nil, errors.New("reciever not found")
	}
	f(reciever, value)
	return acc, nil
}

func readSlotSpec(ctx context.Context, w *World, list Node) (*_SlotSpec, error) {
	cons, ok := list.(*Cons)
	if !ok {
		return nil, fmt.Errorf("[1] %w: %#v", ErrExpectedCons, list)
	}
	identifier, ok := cons.Car.(Symbol)
	if !ok {
		return nil, fmt.Errorf("[1] %w: %#v", ErrExpectedSymbol, cons.Car)
	}
	slotSpec := &_SlotSpec{identifier: identifier}

	list = cons.Cdr
	count := 1
	for IsSome(list) {
		count++
		keywordCons, ok := list.(*Cons)
		if !ok {
			return nil, fmt.Errorf("[%d][1] %w: %#v", count, ErrExpectedCons, list)
		}
		keyword, ok := keywordCons.Car.(Keyword)
		if !ok {
			return nil, fmt.Errorf("[%d][2] %w: %#v", count, ErrExpectedKeyword, keywordCons.Car)
		}
		valueCons, ok := keywordCons.Cdr.(*Cons)
		if !ok {
			return nil, fmt.Errorf("[%d][3] %w: %#v", count, ErrExpectedCons, keywordCons.Cdr)
		}
		var err error
		value := valueCons.Car
		switch keyword {
		case NewKeyword(":reader"):
			if v, ok := value.(Symbol); ok {
				slotSpec.reader = append(slotSpec.reader, v)
			} else {
				return nil, fmt.Errorf(":reader:  %w", ErrExpectedSymbol)
			}
		case NewKeyword(":writer"):
			if v, ok := value.(Symbol); ok {
				slotSpec.writer = append(slotSpec.writer, v)
			} else {
				return nil, fmt.Errorf(":writer: %w", ErrExpectedSymbol)
			}
		case NewKeyword(":accessor"):
			if v, ok := value.(Symbol); ok {
				slotSpec.accessor = append(slotSpec.accessor, v)
			} else {
				return nil, fmt.Errorf(":accessor: %w", ErrExpectedSymbol)
			}
		case NewKeyword(":boundp"):
			if v, ok := value.(Symbol); ok {
				slotSpec.boundp = append(slotSpec.boundp, v)
			} else {
				return nil, fmt.Errorf(":boundp: %w", ErrExpectedSymbol)
			}
		case NewKeyword(":initform"):
			slotSpec.initform = func() (Node, error) { return value.Eval(ctx, w) }
		case NewKeyword(":initarg"):
			if v, ok := value.(Symbol); ok {
				slotSpec.initarg = append(slotSpec.initarg, v)
			} else {
				return nil, fmt.Errorf(":initarg: %w", ErrExpectedSymbol)
			}
		default:
			return nil, fmt.Errorf("invalid keyword %v", keyword)
		}
		if err != nil {
			return nil, fmt.Errorf("[%d][4] %w: %#v", count, err, value)
		}
		if !ok {
			return nil, fmt.Errorf("[%d][5] Domain error: %#v", count, valueCons.Car)
		}
		list = valueCons.Cdr
	}
	return slotSpec, nil
}

type _UserClass struct {
	Symbol
	Super map[Symbol]*_UserClass
	Slot  map[Symbol]*_SlotSpec
}

func (c *_UserClass) Name() Symbol {
	return c.Symbol
}

func (c *_UserClass) InstanceP(n Node) bool {
	_, ok := n.(*_UserClass)
	return ok
}

func (c *_UserClass) Create() Node {
	return &_Receiver{
		_UserClass: c,
		Slot:       make(map[Symbol]Node),
	}
}

func cmdDefClass(ctx context.Context, w *World, args Node) (Node, error) {
	// (defclass class-name (sc-name*) (slot-spec*) class-opt*)

	// class-name
	_className, args, err := Shift(args)
	if err != nil {
		return nil, fmt.Errorf("[1] %w", err)
	}
	className, ok := _className.(Symbol)
	if !ok {
		return nil, fmt.Errorf("[1] %w: %#v", ErrExpectedSymbol, _className)
	}
	class := &_UserClass{
		Symbol: className,
		Super:  make(map[Symbol]*_UserClass),
		Slot:   make(map[Symbol]*_SlotSpec),
	}
	if IsNone(args) {
		w.DefineGlobal(className, class)
		return className, nil
	}
	// (sc-name*) ... super class list
	_scNames, args, err := Shift(args)
	if err != nil {
		return nil, fmt.Errorf("[2] %w", err)
	}
	for IsSome(_scNames) {
		var _super Node
		_super, _scNames, err = w.ShiftAndEvalCar(ctx, _scNames)
		if err != nil {
			return nil, err
		}
		super, ok := _super.(*_UserClass)
		if !ok {
			return nil, errors.New("exepected user class")
		}
		class.Super[super.Symbol] = super
	}
	if IsNone(args) {
		w.DefineGlobal(className, class)
		return className, nil
	}
	// (slot-spec*)
	_slotSpecs, args, err := Shift(args)
	if err != nil {
		return nil, err
	}
	slotSpecs, ok := _slotSpecs.(*Cons)
	if !ok {
		return nil, fmt.Errorf("[3] %w: %#v", ErrExpectedCons, _slotSpecs)
	}
	slotCount := 0
	for p, ok := slotSpecs, true; ok && IsSome(p); p, ok = p.Cdr.(*Cons) {
		slotCount++
		spec, err := readSlotSpec(ctx, w, p.Car)
		if err != nil {
			return nil, fmt.Errorf("[3][%d] %w", slotCount, err)
		}
		class.Slot[spec.identifier] = spec
		getter := func(r *_Receiver) (Node, error) {
			if val, ok := r.Slot[spec.identifier]; ok {
				return val, nil
			}
			return Null, nil
		}
		setter := func(r *_Receiver, value Node) {
			r.Slot[spec.identifier] = value
		}
		if len(spec.reader) > 0 {
			if err := registerGetter(w, spec.reader[0], className, getter); err != nil {
				return nil, err
			}
		}
		if len(spec.writer) > 0 {
			if err := registerSetter(w, spec.writer[0], className, setter); err != nil {
				return nil, err
			}
		}
		if len(spec.accessor) > 0 {
			if err := registerGetter(w, spec.accessor[0], className, getter); err != nil {
				return nil, err
			}
			setterName := NewSymbol("set-" + spec.accessor[0].String())
			if err := registerSetter(w, setterName, className, setter); err != nil {
				return nil, err
			}
		}
	}
	w.DefineGlobal(className, class)
	return className, nil
}

type _Receiver struct {
	*_UserClass
	Slot map[Symbol]Node
}

func (c *_Receiver) PrintTo(w io.Writer, mode PrintMode) (int, error) {
	n, err := c._UserClass.Symbol.PrintTo(w, mode)
	if err != nil {
		return n, err
	}
	delim := []byte{'{'}
	for key, val := range c.Slot {
		n1, err := w.Write(delim)
		n += n1
		if err != nil {
			return n, err
		}
		n1, err = key.PrintTo(w, mode)
		n += n1
		if err != nil {
			return n, err
		}
		n1, err = w.Write([]byte{':'})
		n += n1
		if err != nil {
			return n, err
		}
		n1, err = val.PrintTo(w, mode)
		n += n1
		if err != nil {
			return n, err
		}
		delim = []byte{','}
	}
	n1, err := io.WriteString(w, "}")
	n += n1
	return n, err
}

func (c *_Receiver) String() string {
	var buffer strings.Builder
	c.PrintTo(&buffer, PRINC)
	return buffer.String()
}

func (c *_Receiver) GoString() string {
	var buffer strings.Builder
	c.PrintTo(&buffer, PRINT)
	return buffer.String()
}

func (reciever *_Receiver) callInitForm(classDef *_UserClass) error {
	for _, super := range classDef.Super {
		if err := reciever.callInitForm(super); err != nil {
			return err
		}
	}
	for name, slot1 := range classDef.Slot {
		if _, ok := reciever.Slot[name]; !ok && slot1.initform != nil {
			var err error
			reciever.Slot[name], err = slot1.initform()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (reciever *_Receiver) callInitArg(classDef *_UserClass, initArg Symbol, initVal Node) bool {
	for name, slot := range classDef.Slot {
		for _, slotInitArg := range slot.initarg {
			if slotInitArg == initArg {
				reciever.Slot[name] = initVal
				return true
			}
		}
	}
	for _, super := range classDef.Super {
		if reciever.callInitArg(super, initArg, initVal) {
			return true
		}
	}
	return false
}

func cmdCreate(ctx context.Context, w *World, args Node) (Node, error) {
	_class, args, err := w.ShiftAndEvalCar(ctx, args)
	if err != nil {
		return nil, err
	}
	class, ok := _class.(Class)
	if !ok {
		return nil, errors.New("expect class")
	}
	_this := class.Create()
	this, ok := _this.(*_Receiver)
	if !ok {
		if IsSome(args) {
			return nil, fmt.Errorf("%s does not have slot", class.String())
		}
		return _this, nil
	}
	for IsSome(args) {
		var _initArg Node
		_initArg, args, err = w.ShiftAndEvalCar(ctx, args)
		if err != nil {
			return nil, err
		}
		initArg, ok := _initArg.(Symbol)
		if !ok {
			return nil, ErrExpectedSymbol
		}
		var initVal Node
		initVal, args, err = w.ShiftAndEvalCar(ctx, args)
		if err != nil {
			return nil, err
		}
		if this._UserClass == nil {
			panic("!")
		}
		this.callInitArg(this._UserClass, initArg, initVal)
	}
	return this, this.callInitForm(this._UserClass)
}
