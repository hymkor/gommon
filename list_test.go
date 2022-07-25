package gmnlisp

import (
	"testing"
)

func TestList(t *testing.T) {
	assertEqual(t, "(car '(1 2))", Integer(1))
	assertEqual(t, "(car '(1 . 2))", Integer(1))
	assertEqual(t, "(cdr '(1 . 2))", Integer(2))
	assertEqual(t, "(cdr '(1 2))", List(Integer(2)))

	assertEqual(t, `(list 1 2 3 4)`,
		List(Integer(1), Integer(2), Integer(3), Integer(4)))
	assertEqual(t, `(append '(1 2) '(3 4))`, List(Integer(1), Integer(2), Integer(3), Integer(4)))
	assertEqual(t, `(append '(1 2) '(3 4) '(5 6))`,
		List(Integer(1), Integer(2), Integer(3), Integer(4), Integer(5), Integer(6)))
	assertEqual(t, `(member 'c '(a b c d e))`,
		List(Symbol("c"), Symbol("d"), Symbol("e")))

	assertEqual(t, `(cadr '(1 2 3))`, Integer(2))
	assertEqual(t, `(caddr '(1 2 3 4 5 ))`, Integer(3))
	assertEqual(t, `(cadddr '(1 2 3 4 5))`, Integer(4))
	assertEqual(t, `(cddr '(1 2 3 4 5))`,
		List(Integer(3), Integer(4), Integer(5)))
	assertEqual(t, `(cdddr '(1 2 3 4 5))`,
		List(Integer(4), Integer(5)))

	assertEqual(t, "(cons 1 2)", &Cons{Car: Integer(1), Cdr: Integer(2)})

	assertEqual(t, `(mapcar (function +) '(1 2 3) '(4 5 6))`,
		List(Integer(5), Integer(7), Integer(9)))
	assertEqual(t, `(mapcar #'+ '(1 2 3) '(4 5 6))`,
		List(Integer(5), Integer(7), Integer(9)))
	assertEqual(t, `(mapcar '+ '(1 2 3) '(4 5 6))`,
		List(Integer(5), Integer(7), Integer(9)))
	assertEqual(t, `(mapcar (lambda (a b) (+ a b)) '(1 2 3) '(4 5 6))`,
		List(Integer(5), Integer(7), Integer(9)))
	assertEqual(t, `(mapcar #'(lambda (a b) (+ a b)) '(1 2 3) '(4 5 6))`,
		List(Integer(5), Integer(7), Integer(9)))
	assertEqual(t, `(listp ())`, True)
	assertEqual(t, `(listp 1)`, Null)
	assertEqual(t, `(listp '(1 2 3))`, True)

	assertEqual(t, `(length (list 1 2 3 4))`, Integer(4))
	assertEqual(t, `(length '(list 1 2 3))`, Integer(4))
}
