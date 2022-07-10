package gmnlisp

import (
	"io"
)

type Cons struct {
	Car Node
	Cdr Node
}

func (c *Cons) GetCar() Node {
	if c.Car == nil {
		return Null
	}
	return c.Car
}

func (c *Cons) GetCdr() Node {
	if c.Cdr == nil {
		return Null
	}
	return c.Cdr
}

func (this *Cons) IsNull() bool {
	return false
}

func HasValue(node Node) bool {
	return node != nil && !node.IsNull()
}

func IsNull(node Node) bool {
	return node == nil || node.IsNull()
}

func (this *Cons) isTailNull() bool {
	if IsNull(this.Cdr) {
		return true
	} else if next, ok := this.Cdr.(*Cons); ok {
		return next.isTailNull()
	} else {
		return false
	}
}

func (this *Cons) writeToWithoutKakko(w io.Writer, rich bool) {
	if IsNull(this.Car) {
		io.WriteString(w, "()")
	} else if rich {
		this.Car.PrintTo(w)
	} else {
		this.Car.PrincTo(w)
	}

	if HasValue(this.Cdr) {
		if this.isTailNull() {
			// output as ( X Y Z ...)

			for p, ok := this.Cdr.(*Cons); ok && HasValue(p); p, ok = p.Cdr.(*Cons) {
				io.WriteString(w, " ")
				if rich {
					p.Car.PrintTo(w)
				} else {
					p.Car.PrincTo(w)
				}
			}
		} else {
			// output as ( X . Y )

			io.WriteString(w, " . ")
			if rich {
				this.GetCdr().PrintTo(w)
			} else {
				this.GetCdr().PrincTo(w)
			}
		}
	}
}

func (this *Cons) PrintTo(w io.Writer) {
	io.WriteString(w, "(")
	this.writeToWithoutKakko(w, true)
	io.WriteString(w, ")")
}

func (this *Cons) PrincTo(w io.Writer) {
	io.WriteString(w, "(")
	this.writeToWithoutKakko(w, false)
	io.WriteString(w, ")")
}

func (this *Cons) Equals(n Node) bool {
	value, ok := n.(*Cons)
	if !ok {
		return false
	}
	return this.GetCar().Equals(value.Car) &&
		this.GetCdr().Equals(value.Cdr)
}
