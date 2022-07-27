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

	assertEqual(t, `(reverse '(1 2 3 4))`,
		List(Integer(4), Integer(3), Integer(2), Integer(1)))

	assertEqual(t, `(setq collection '((a . 1) (b . 2) (c . 3)))
			(assoc 'a collection)`,
		&Cons{Car: Symbol("a"), Cdr: Integer(1)})

	assertEqual(t, `(nth 2 '(10 20 30 40))`, Integer(30))
	assertEqual(t, `(nthcdr 2 '(10 20 30 40))`, List(Integer(30), Integer(40)))

	assertEqual(t, `
	(setq collection (list
		(cons "X"  1) (cons "Y"  2) (cons "Z" 4)))
	(subst (cons "X" 7) (assoc "X" collection) collection)`,
		List(
			&Cons{Car: String("X"), Cdr: Integer(7)},
			&Cons{Car: String("Y"), Cdr: Integer(2)},
			&Cons{Car: String("Z"), Cdr: Integer(4)}))
}
