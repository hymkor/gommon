(defmacro incf (place &rest args)
  (let ((delta (if args (car args) 1)))
    `(setf ,place (+ ,place ,delta))))
(defmacro decf (place &rest args)
  (let ((delta (if args (car args) 1)))
    `(setf ,place (- ,place ,delta))))
(defun swap-elt (source z newvalue)
  (if (stringp source)
    (string-append
      (subseq source 0 z)
      (create-string 1 newvalue)
      (subseq source (1+ z) (length source)))
    (let ((s source))
      (while s
        (if (zerop z)
          (set-car newvalue s))
        (decf z)
        (setq s (cdr s))
        )
      source)))
(defun swap-subseq (seq start end newvalue)
  (if (stringp seq)
    (string-append (subseq seq 0 start)
                   newvalue
                   (subseq seq end (length seq)))
    (let ((orig seq))
      (while seq
        (if (and (<= start 0) (> end 0) newvalue)
          (progn
            (set-car (car newvalue) seq)
            (setq newvalue (cdr newvalue))))
        (decf start)
        (decf end)
        (setq seq (cdr seq)))
      orig)))
(let ((setf-table
        '((car . set-car)
          (cdr . set-cdr)
          (elt . set-elt)
          (dynamic . set-dynamic)
          (subseq . set-subseq)
          (setq . set-setq)
          (assoc . set-assoc))))
  (defmacro setf (expr newvalue)
    (if (symbolp expr)
      `(setq ,expr ,newvalue)
      (let* ((name (car expr))
             (pair (assoc name setf-table))
             (tmp nil)
             (setter
               (if pair
                 (cdr pair)
                 (progn
                   (setq tmp (convert (string-append "set-" (convert name <string>)) <symbol>))
                   (setq setf-table (cons (cons name tmp) setf-table))
                   tmp)))
             (arguments (cdr expr)))
        (cons setter (cons newvalue arguments))))))
(defmacro set-elt (newvalue seq z)
  `(setf ,seq (swap-elt ,seq ,z ,newvalue)))
(defmacro set-dynamic (newvalue name)
  `(defdynamic ,name ,newvalue))
(defmacro set-subseq (newvalue seq start end)
  `(setf ,seq (swap-subseq ,seq ,start ,end ,newvalue)))
(defmacro set-setq (newvalue name avlue)
  (let ((name (elt expr 1)) (value (elt expr 2)))
    `(progn (setq ,name ,value) (setf ,name ,newvalue))))
(defmacro set-assoc (newvalue key m)
  (let ((L (gensym)) (K (gensym)) (tmp (gensym)))
    `(let* ((,L ,m)
            (,K ,key)
            (,tmp nil))
       (while ,L
         (if (and (setq ,tmp (car ,L)) (consp ,tmp) (equal ,K (car ,tmp)))
           (set-car ,newvalue ,L))
         (setq ,L (cdr ,L)))))
  )
(defmacro dolist (vars &rest body)
  (let ((var (car vars))
        (values (elt vars 1))
        (rest (gensym)))
    `(block
       nil
       (let ((,var nil)(,rest ,values))
         (while ,rest
           (setq ,var (car ,rest))
           (setq ,rest (cdr ,rest))
           ,@body)))))
(defmacro dotimes (vars &rest commands)
  (let ((var (car vars))
        (count (elt vars 1))
        (end (gensym)))
    `(let ((,var 0)(,end ,count))
       (while (< ,var ,end)
         (progn ,@commands)
         (setq ,var (+ 1 ,var))))))
; vim:set lispwords+=while,defglobal:
