(defclass <point1d> ()
  ((X :initform 0.0 :initarg x :accessor point-x :reader get-point-x :writer let-point-x)))

(defclass <point2d> (<point1d>)
  ((Y :initform 0.0 :initarg y :accessor point-y :reader get-point-y :writer let-point-y)))

(defclass <point3d> (<point2d>)
  ((Z :initform 0.0 :initarg z :accessor point-z :reader get-point-z :writer let-point-z)))

(let
  ((p1 (create <point2d> 'x 3))
   (p2 (create <point3d> 'x 10 'y 20 'z 30)))

  ; test accessor as getter
  (assert-eq (point-x p1) 3)
  (assert-eq (point-y p1) 0.0)

  ; test reader
  (assert-eq (get-point-x p1) 3)
  (assert-eq (get-point-y p1) 0.0)

  ; test accessor as setter
  (setf (point-x p1) 88)
  (assert-eq (point-x p1) 88)

  ; test writer
  (let-point-y 99 p1)
  (assert-eq (point-y p1) 99)

  ; test super class
  (assert-eq (point-x p2) 10)
  (assert-eq (point-y p2) 20)
  (assert-eq (point-z p2) 30)

  (assert-eq (create <integer>) 0)
  (assert-eq (create <float>) 0.0)
  (assert-eq (create <string>) "")

  (assert-eq (instancep p1 <point2d>) t)
  (assert-eq (instancep p2 <point2d>) t)
  (assert-eq (instancep p1 <point3d>) nil)
  (assert-eq (instancep 0 <point2d>) nil)
  (assert-eq (assure <point2d> p1) p1)
  (assert-eq (the <point2d> p1) p1)
  (assert-eq
    (catch
      'c
      (with-handler
        (lambda (c)
          (if (instancep c <domain-error>)
            (throw 'c "OK")
            "NG1"))
        (assure <point3d> p1)
        "NG2"
        )
      )
    "OK")
  )

(assert-eq (subclassp <point1d> <point2d>) nil)
(assert-eq (subclassp <point1d> <point1d>) nil)
(assert-eq (subclassp <point2d> <point1d>) t)

(defclass <foo> ()
  ((X :initarg x :accessor x-of-foo :boundp has-x)))
(let ((f1 (create <foo> 'x 0))
      (f2 (create <foo>)))
  (assert-eq (has-x f1) t)
  (assert-eq (has-x f2) nil)
  (setf (x-of-foo f2) 1)
  (assert-eq (has-x f2) t))

(defmethod initialize-object ((this <foo>) (x <string>))
  (setf (x-of-foo this) x)
  this)
(let ((f1 (create <foo> (string-append "x" "y"))))
  ; (format (standard-output) "~S~%" f1)
  (assert-eq (x-of-foo f1) "xy"))

(assert-eq (subclassp <error> <object>) t)
(assert-eq (subclassp <object> <error>) nil)
(assert-eq (subclassp <error> <error>) nil)

(assert-eq (subclassp <array> <built-in-class>) t)
(assert-eq (subclassp (class-of #(1 2 3)) <built-in-class>) t)
(assert-eq (subclassp (class-of t) <built-in-class>) t)
(assert-eq (subclassp (class-of nil) <built-in-class>) t)
(assert-eq (subclassp <string> <built-in-class>) t)
(assert-eq (subclassp (class-of 's) <built-in-class>) t)
(assert-eq (subclassp <keyword> <built-in-class>) t)
(assert-eq (subclassp (class-of :keyword) <built-in-class>) t)
(assert-eq (subclassp <cons> <built-in-class>) t)
(assert-eq (subclassp (class-of (cons 1 2)) <built-in-class>) t)
(assert-eq (subclassp (class-of (create-string-output-stream)) <built-in-class>) t)
(assert-eq (subclassp <integer> <built-in-class>) t)
(assert-eq (subclassp (class-of 1) <built-in-class>) t)
(assert-eq (subclassp <float> <built-in-class>) t)
(assert-eq (subclassp (class-of 1.1) <built-in-class>) t)
(assert-eq (subclassp <string> <built-in-class>) t)
(assert-eq (subclassp (class-of "a") <built-in-class>) t)
(assert-eq (subclassp <point1d> <built-in-class>) nil)
