(defun swap-elt (source z newvalue)
  (cond
    ((stringp source)
     (string-append
       (subseq source 0 z)
       (create-string 1 newvalue)
       (subseq source (1+ z) (length source))))
    (t
      (let ((s source))
        (while s
          (if (zerop z)
            (set-car newvalue s))
          (setq z (1- z))
          (setq s (cdr s))
          )
        source)
      )
    )
  )
(defun swap-nthcdr (source z newvalue)
  (let ((s source))
    (while s
      (setq z (1- z))
      (if (zerop z)
        (set-cdr newvalue s))
      (setq s (cdr s))
      )
    source)
  )
(defun swap-subseq (seq start end newvalue)
  (cond
    ((stringp seq)
     (string-append (subseq seq 0 start)
                    newvalue
                    (subseq seq end (length seq))))
     (t
       (append (subseq seq 0 start)
               newvalue
               (subseq seq end (length seq))))))
(defun subst (newitem olditem L)
  (let ((result nil))
    (while L
      (if (equal olditem (car L))
        (setq result (cons newitem result))
        (setq result (cons (car L) result)))
      (setq L (cdr L)))
    (nreverse result)))
(defmacro setf (expr newvalue)
  (if (symbolp expr)
    `(setq ,expr ,newvalue)
    (case (car expr)
      (('car)
       (let ((cns (elt expr 1)))
         `(set-car ,newvalue ,cns)))
      (('cdr)
       (let ((cns (elt expr 1)))
         `(set-cdr ,newvalue ,cns)))
      (('elt)
       (let ((seq (elt expr 1)) (z (elt expr 2)))
         `(setf ,seq (swap-elt ,seq ,z ,newvalue))
         ))
      (('nth) ; for pkg/common
       (let ((seq (elt expr 2)) (z (elt expr 1)))
         `(setf ,seq (swap-elt ,seq ,z ,newvalue))
         ))
      (('nthcdr) ; for pkg/common
       (let ((seq (elt expr 2)) (z (elt expr 1)))
         `(setf ,seq (swap-nthcdr ,seq ,z ,newvalue))
       ))
      (('dynamic)
       (let ((name (elt expr 1)))
         `(defdynamic ,name ,newvalue))
       )
      (('subseq)
       (let ((seq (elt expr 1)) (start (elt expr 2)) (end (elt expr 3)))
         `(setf ,seq (swap-subseq ,seq ,start ,end ,newvalue))
         )
       )
      (('setq)
       (let ((name (elt expr 1)) (value (elt expr 2)))
         `(progn (setq ,name ,value) (setf ,name ,newvalue))
         )
       )
      (('assoc)
       (let ((key (elt expr 1)) (m (elt expr 2)))
         `(let ((_m ,m))
            (setf ,m (subst ,newvalue (assoc ,key _m) _m)))
         )
       )
      (t
        (error))
      ) ; case
    ); if
  ) ; defmacro
; vim:set lispwords+=while:
