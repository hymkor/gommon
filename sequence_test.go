package gmnlisp

import (
	"testing"
)

func TestLength(t *testing.T) {
	assertEqual(t, `(length (list 1 2 3 4))`, Integer(4))
	assertEqual(t, `(length '(list 1 2 3))`, Integer(4))

	assertEqual(t, `(length "12345")`, Integer(5))
}

func TestMapCar(t *testing.T) {
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

	assertEqual(t, `(mapcar #'car '((1 a) (2 b) (3 c)))`,
		List(Integer(1), Integer(2), Integer(3)))
	// assertEqual(t, `(mapcar #'abs '((3 -4 2 -5 -6)))`,
	//	List(Integer(3), Integer(4), Integer(2), Integer(5), Integer(6)))
	assertEqual(t, `(mapcar #'cons '(a b c) '(1 2 3))`,
		List(&Cons{Car: Symbol("a"), Cdr: Integer(1)},
			&Cons{Car: Symbol("b"), Cdr: Integer(2)},
			&Cons{Car: Symbol("c"), Cdr: Integer(3)}))
}

func TestMap(t *testing.T) {
	assertEqual(t, `(map 'string '1+ "123")`, String("234"))
	assertEqual(t, `(map 'list '1+ '(1 2 3))`, List(Integer(2), Integer(3), Integer(4)))
	assertEqual(t, `(length (map 'list #'null '(nil 2 3)))`, Integer(3))
}

func TestMapC(t *testing.T) {
	assertEqual(t, `
		(let ((buffer (create-string-output-stream)) result)
			(setq result (mapc (lambda (c) (format-char buffer (1+ c))) "ABC"))
			(list result (get-output-stream-string buffer))
		)`, List(String("ABC"), String("BCD")))
}

func TestMapCan(t *testing.T) {
	assertEqual(t, `(mapcan (lambda (x) (if (> x 0) (list x))) '(-3 4 0 5 -2 7))`,
		List(Integer(4), Integer(5), Integer(7)))
}

func TestMapList(t *testing.T) {
	assertEqual(t, `(maplist #'append '(1 2 3 4) '(1 2) '(1 2 3))`,
		List(List(Integer(1), Integer(2), Integer(3), Integer(4),
			Integer(1), Integer(2), Integer(1), Integer(2), Integer(3)),
			List(Integer(2), Integer(3), Integer(4), Integer(2), Integer(2), Integer(3))))
}

func TestMapL(t *testing.T) {
	assertEqual(t, `
		(let ((k 0))
			(mapl
				(lambda (x)
					(setq k (+ k (if (member (car x) (cdr x)) 0 1)))
				)
				'(a b a c d b c)
			)
		k)`, Integer(4))
}

func TestMapCon(t *testing.T) {
	assertEqual(t, `
		(mapcon
			(lambda (x)
				(if (member (car x) (cdr x)) (list (car x)))
			)
			'(a b a c d b c b c)
		)`, List(Symbol("a"), Symbol("b"), Symbol("c"), Symbol("b"), Symbol("c")))
}

func TestCoerce(t *testing.T) {
	assertEqual(t, `(coerce '(#\a #\b) 'string)`, String("ab"))
	assertEqual(t, `(coerce '(#\a #\b) 'utf8string)`, UTF8String("ab"))
	assertEqual(t, `(coerce '(#\a #\b) 'utf32string)`, UTF32String("ab"))
	assertEqual(t, `(coerce '(#\a #\b) 'list)`, List(Rune('a'), Rune('b')))
}

func TestConcatenate(t *testing.T) {
	assertEqual(t, `(concatenate 'string "123" "456")`, String("123456"))
	assertEqual(t, `(concatenate 'list '(1 2 3) '(4 5 6))`,
		List(Integer(1), Integer(2), Integer(3), Integer(4), Integer(5), Integer(6)))
}

func TestReverse(t *testing.T) {
	assertEqual(t, `(reverse '(1 2 3 4))`,
		List(Integer(4), Integer(3), Integer(2), Integer(1)))
	assertEqual(t, `(reverse "12345")`, String("54321"))
}

func TestFind(t *testing.T) {
	assertEqual(t, `(find 2 '(1 2 3))`, Integer(2))
	assertEqual(t, `(find 2.0 '(1 2 3))`, Null)
	assertEqual(t, `(find 2.0 '(1 2 3) :test #'(lambda (a b) (equalp a b)))`, Integer(2))
	assertEqual(t, `(find 2.0 '(1 2 3) :test #'(lambda (a b) (eql a b)))`, Null)
}

func TestMember(t *testing.T) {
	assertEqual(t, `(member 'c '(a b c d e))`,
		List(Symbol("c"), Symbol("d"), Symbol("e")))
	assertEqual(t, `(member #\c "abcd")`, String("cd"))
	assertEqual(t, `(member #\C "abcd" :test #'(lambda (a b) (equalp a b)))`, String("cd"))
}

func TestPosition(t *testing.T) {
	assertEqual(t, `(position 'c '(a b c d e))`, Integer(2))
	assertEqual(t, `(position #\c "abcd")`, Integer(2))
	assertEqual(t, `(position #\C "abcd")`, Null)
	assertEqual(t, `(position #\C "abcd" :test #'(lambda (a b) (equalp a b)))`, Integer(2))
}

func TestSubSeq(t *testing.T) {
	assertEqual(t, `(subseq "12345" 2 4)`, String("34"))
	assertEqual(t, `(subseq "12345" 2)`, String("345"))
	assertEqual(t, `(subseq '(1 2 3 4 5) 2 4)`, List(Integer(3), Integer(4)))
	assertEqual(t, `(subseq '(1 2 3 4 5) 2)`, List(Integer(3), Integer(4), Integer(5)))
}

func TestSetfSubSeq(t *testing.T) {
	assertEqual(t, `
		(let ((m (to-utf32 "12345")))
			(setf (subseq m 2 4) (to-utf32 "xx"))
			m)`, UTF32String("12xx5"))

	assertEqual(t, `
		(let ((m (list 1 2 3 4 5)))
			(setf (subseq m 2 4) (list 0 0))
			m)`, List(Integer(1), Integer(2), Integer(0), Integer(0), Integer(5)))
}

func TestElt(t *testing.T) {
	assertEqual(t, `(elt '(a b c) 2)`, Symbol("c"))
	// assertEqual(t, `(elt (vector 'a 'b 'c) 1)`,Symbol("b"))
	assertEqual(t, `(elt "abc" 0)`, Rune('a'))
}
