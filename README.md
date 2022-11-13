[![GoDoc](https://godoc.org/github.com/hymkor/gmnlisp?status.svg)](https://godoc.org/github.com/hymkor/gmnlisp)

gmnlisp
=======

Gmnlisp is a small Lisp implementation in Go.
( Now under constructing. Experimental implementing )

![Example image](factorial.png)

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/hymkor/gmnlisp"
)

func main() {
    lisp := gmnlisp.New()
    lisp = lisp.Let(gmnlisp.Variables{
        gmnlisp.NewSymbol("a"): gmnlisp.Integer(1),
        gmnlisp.NewSymbol("b"): gmnlisp.Integer(2),
    })
    value, err := lisp.Interpret(context.TODO(), "(+ a b)")
    if err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        return
    }
    value.PrintTo(os.Stdout, gmnlisp.PRINT)
    fmt.Println()
}
```

```
$ go run examples/example1.go
3
```

gmnlpp - Text preprocessor by gmnlisp
-------------------------------------

This page was generated by a preprocessor with built-in gmnlisp.
The text before proprocessed is [here](https://github.com/hymkor/gmnlisp/blob/master/_README.md)

Support Types
-------------

integer , float , string , symbol , cons , list , array , character , t/nil

Support functions
-----------------

#### List and Sequence

- (cons OBJ1 OBJ2)
- (car CONS)
- (cdr CONS)
- (quote OBJ)
- 'OBJ
- (list OBJ...)
- (rest LIST)
- (length SEQUENCE)
- (last LIST)
- (reverse LIST)
- (nreverse LIST)
- (append LIST...)
- (assoc OBJ LIST)
- (subseq SEQUENCE Z1 Z2)
- (elt SEQUENCE INDEX)
- (member ATOM LIST)

#### array

- #(...) , #2a((...) (...)) , #3a(((.. ..))) ...
- (create-array '(DIM...) INITIAL-VALUE)
- (array-dimensions ARRAY)
- (aref ARRAY INDEX...)
- (setf (aref ARAY INDEX) NEWVALUE)
- (set-aref NEWVALUE ARRAY INDEX...)
- (arrayp OBJ)

#### Variables

- (defdyncamic NAME FORM)
- (defglobal NAME FORM)
- (dynamic VAR)
- (dynamic-let ((VAR FORM)...) BODY-FORM...)
- (let ((VAR FORM)... ) BODY-FORM...)
- (let\* ((VAR FORM)...) BODY-FORM...)
- (set-car NEW-CAR CONS)
- (set-cdr NEW-CDR CONS)
- (setf PLACE FORM)
- (setq VAR FORM)

#### Operator

- (= EXP1 EXP2...)
- (/= EXP1 EXP2...)
- (\< EXP1 EXP2...)
- (\<= EXP1 EXP2...)
- (\> EXP1 EXP2...)
- (\>= EXP1 EXP2...)
- (eq EXP1 EXP2...)
- (eql EXP1 EXP2...)
- (equal EXP1 EXP2...)
- (equalp EXP1 EXP2...)
- (+ EXP1 EXP2...)
- (- EXP1 EXP2...)
- (\* EXP1 EXP2...)
- (/ EXP1 EXP2...)
- (mod EXP1 EXP2)
- (rem EXP1 EXP2)
- (1+ EXP)
- (1- EXP)
- (incf VAR [VALUE]) [MACRO]
- (decf VAR [VALUE]) [MACRO]
- (and EXP1 EXP2..)
- (or EXP1 EXP2..)
- (not EXP)
- (string= STRING1 STRING2)
- (string/= STRING1 STRING2)
- (string&lt; STRING1 STRING2)
- (string&gt; STRING1 STRING2)
- (string&gt;= STRING1 STRING2)
- (string&lt;= STRING1 STRING2)
- (string-index SUBSTRING STRING [START])
- (string-append STRING...)
- (create-String I [INITIAL-CHARACTER])

#### test

- (atom OBJ)
- (consp OBJ)
- (evenp OBJ)
- (floatp OBJ)
- (functionp OBJ)
- (integerp OBJ)
- (listp OBJ)
- (minusp OBJ)
- (null OBJ)
- (numberp OBJ)
- (oddp OBJ)
- (plusp OBJ)
- (stringp OBJ)
- (symbolp OBJ)
- (zerop OBJ)

#### Convert

- (convert OBJ &lt;float&gt;)
- (convert OBJ &lt;integer&gt;)
- (convert OBJ &lt;list&gt;)
- (convert OBJ &lt;string&gt;)
- (convert OBJ &lt;symbol&gt;)
- (parse-number STRING)
- (truncate X)
- (floor X)
- (ceiling X)
- (round X)

#### Branch and Loop

- (case KEYFORM ((KEY...) FORM...)... [(t FORM...)])
- (cond (TEST FORM...)...)
- (for ((VAR INIT [STEP])...) (END-TEST RESULT...) FORM... )
- (if TEST-FORM THEN-FORM ELSE-FORM)
- (progn FORM...)
- (prog1 FORM...)
- (prog2 FORM...)
- (while TEST-FORM BODY-FORM...)
- (tagbody {TAG|FORM}...)
    - (go TAG)
- (dolist (VAR '(VALUES..)) FORM...) [MACRO]
- (dotimes (VAR N) FORM...) [MACRO]
- (when TEST-FORM THEN-FORM...)
- (unless TEST-FORM ELSE-FORM...)

#### Functions

- (lambda (IDENTIFIER... [&amp;rest IDENTIFIER]) FORM...)
- (defun FUNCTION-NAME (IDENTIFIER... [&amp;rest IDENTIFIER]) FORM...)
- (labels ((FUNCTION-NAME LAMBDA-LIST FORM...)...) BODY-FORM...)
- (flet ((FUNCTION-NAME LAMBDA-LIST FORM...)...) BODY-FORM...)

#### Constant

- most-postive-fixnum
- most-negative-fixnum
- pi

#### Function Reference

- (function FUNCTION)
- #'FUNCTION

#### Macro

- (defmacro NAME (IDENTIFIER... [&amp;rest IDENTIFIER]) FORM...)

#### Mapping

- (mapcar #'FUNCTION LIST)
- (mapc #'FUNCTION LIST)
- (mapcan #'FUNCTION LIST)
- (maplist #'FUNCTION LIST)
- (mapl #'FUNCTION LIST)
- (mapcon #'FUNCTION LIST)
- (apply #'FUNCTION [PARAMS...] LIST)
- (funcall #'FUNCTION EXP1...)

#### I/O

- (close STREAM)
- (create-string-input-stream STRING)
- (create-string-output-stream)
    - (get-output-stream-string STRSTREAM)
- (error-output)
- (file-length FILENAME ELEMENT-CLASS)
- (format {OUTPUT-STREAM|t|nil} FORMAT-STRING OBJ..)
    - (format t "..") is same as (format (standard-output) "..")
    - (format nil "..") is same as (let ((B create-string-output-stream)) (format B "..") (get-output-stream-string B))
- (format-char {OUTPUT-STREAM|t|nil} CHAR)
- (format-float {OUTPUT-STREAM|t|nil} FLOAT)
- (format-integer {OUTPUT-STREAM|t|nil} INTEGER RADIX)
- (format-object {OUTPUT-STREAM|t|nil} OBJ ESCAPE-P)
- (open-input-file FILENAME)
- (open-output-file FILENAME)
- (probe-file FILENAME)
- (read [STREAM [EOF-FLAG [EOF-VALUE]]])
- (read-line [STREAM [EOF-FLAG [EOF-VALUE]]])
- (standard-input)
- (standard-output)
- (with-open-input-file (NAME FILENAME) FORM...)
- (with-open-output-file (NAME FILENAME) FORM...)

#### Exceptions

- (block {SYMBOL|nil} FORM...)
    - (return RESULT-FORM)
    - (return-from SYMBOL RESULT-FORM)
- (catch TAG-FORM FORM...)
    - (throw TAG-FORM RESULT-FORM)
- (unwind-protect FORM CLEANUP-FORM...)
- (with-handler HANDLER FORM...)

#### Quit

- (exit)
- (quit)
- (abort)
