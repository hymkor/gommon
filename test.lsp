(defmacro assert (func expect)
  (let ((result (gensym)))
    `(let ((,result ,func))
       (if (not (equal ,result ,expect))
         (progn
           (format t "Failed: ~s~%" (quote ,func))
           (format t "  expect: ~s but ~s~%" ,expect ,result)
           (quit))))))

(let ((A (create-array '(3 2) 1)))
  (set-array-elt 2 A 1 1)
  (set-array-elt (create-array '(2) 4) A 0)
  (assert (elt A 0 0) 4)
  (assert (elt A 0 1) 4)
  (assert (elt A 0) (create-array '(2) 4))
  (assert (elt A 1 1) 2)
  (assert (elt A 1 0) 1)
  )

