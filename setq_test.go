package gmnlisp

import (
	"testing"
)

func TestLet(t *testing.T) {
	assertEqual(t, `(let* ((x 2)(y x)) y)`, Integer(2))
	assertEqual(t, `(let ((x 0)) (let ((x 2)(y x)) y))`, Integer(0))
}

func TestDefvar(t *testing.T) {
	assertEqual(t, `(defvar a "ahaha")`, Symbol("a"))
	assertEqual(t, `(defvar a "ahaha")(defvar a "ihihi") a`, String("ahaha"))

	assertEqual(t, `
		(defvar counter 0)
		(defvar a (setq counter (1+ counter)))
		(defvar a (setq counter (1+ counter)))
		counter`, Integer(1))
}

func TestDefparameter(t *testing.T) {
	assertEqual(t, `(defparameter a "ahaha")`, Symbol("a"))
	assertEqual(t, `(defparameter a "ahaha")(defparameter a "ihihi") a`, String("ihihi"))
}

func TestSetf(t *testing.T) {
	assertEqual(t, `(defvar x)
					(setf (car (setq x (cons 1 2))) 3)
					x`, &Cons{Integer(3), Integer(2)})
	assertEqual(t, `(defparameter x (cons 1 2))
					(setf (cdr x) 3)
					x`, &Cons{Integer(1), Integer(3)})
	assertEqual(t, `(defparameter x (list 1 2 3 4))
					(setf (nth 2 x) 0)
					x`, List(Integer(1), Integer(2), Integer(0), Integer(4)))
	assertEqual(t, `(defparameter x (list 1 2 3 4))
					(setf (nthcdr 2 x) (list 7))
					x`, List(Integer(1), Integer(2), Integer(7)))

	assertEqual(t, `(defparameter x (list 1 2 3 4))
					(setf (cadr x) 0)
					x`, List(Integer(1), Integer(0), Integer(3), Integer(4)))

	assertEqual(t, `(defparameter x (list 1 2 3 4))
					(setf (caddr x) 0)
					x`, List(Integer(1), Integer(2), Integer(0), Integer(4)))

	assertEqual(t, `(defparameter x (list 1 2 3 4))
					(setf (cadddr x) 0)
					x`, List(Integer(1), Integer(2), Integer(3), Integer(0)))

	assertEqual(t, `(defparameter x (list 1 2 3 4))
					(setf (cddr x) (list 0))
					x`, List(Integer(1), Integer(2), Integer(0)))

	assertEqual(t, `(defparameter x (list 1 2 3 4))
					(setf (cdddr x) (list 0))
					x`, List(Integer(1), Integer(2), Integer(3), Integer(0)))

	assertEqual(t, `
		(defvar m (list (cons 1 "A") (cons 2 "B") (cons 3 "C")))
		(setf (cdr (assoc 1 m)) "X")
		m`, List(
		&Cons{Integer(1), String("X")},
		&Cons{Integer(2), String("B")},
		&Cons{Integer(3), String("C")}))

	(assertEqual(t, `
		(let ((m '((1 . "A") (2 . "B") (3 . "C"))) pair )
		  (if (setq pair (assoc 1 m))
			  (setf (cdr pair) "X")
		  )
		  m
		)`, List(
		&Cons{Integer(1), String("X")},
		&Cons{Integer(2), String("B")},
		&Cons{Integer(3), String("C")})))
}
