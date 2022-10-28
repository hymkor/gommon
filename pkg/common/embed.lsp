(defun set-nth (newvalue Z L)
  (set-elt newvalue L Z))
(defun set-nthcdr (newvalue z source)
  (let ((s source))
    (while s
           (setq z (1- z))
           (if (zerop z)
             (set-cdr newvalue s))
           (setq s (cdr s)))
    source))
(defun set-cadr (newvalue L)
  (set-car newvalue (cdr L)))
(defun set-caddr (newvalue L)
  (set-car newvalue (cdr (cdr L))))
(defun set-cadddr (newvalue L)
  (set-car newvalue (cdr (cdr (cdr L)))))
(defun set-cddr (newvalue L)
  (set-cdr newvalue (cdr L)))
