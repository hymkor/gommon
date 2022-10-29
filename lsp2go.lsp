(defun to-safe (s)
  (let ((index nil))
    (while (and s (setq index (string-index "-" s)))
      (setq s (string-append (subseq s 0 index)
                             "_"
                             (subseq s (+ index 1) (length s))))
      )
    s))
(defun defun2lambda (node)
  (if (consp node)
    (case (car node)
      (('defun)
       (set-car 'lambda node)
       (let ((name (elt node 1)))
         (set-cdr (cdr (cdr node)) node)
         name)
       )
      (('defmacro)
       (set-car 'lambda-macro node)
       (let ((name (elt node 1)))
         (set-cdr (cdr (cdr node)) node)
         name)
       )
      (t
        (or (defun2lambda (car node))
            (defun2lambda (cdr node))))
      )
    )
  )

(let ((node nil)
      (name nil)
      (packagename (car *posix-argv*))
      (arguments (cdr *posix-argv*)))
  (format t "package ~a~%" packagename)
  (or (equal packagename "gmnlisp")
      (format t "~%import . \"github.com/hymkor/gmnlisp\""))
  (format t "~%// This code is generated by lsp2go.lsp")
  (while arguments
    (format t " ~a" (car arguments))
    (setq arguments (cdr arguments))
    )
  (format t "~%")
  (format t "var embedFunctions = map[Symbol]Node{~%")
  (while (setq node (read (standard-input) nil nil))
    (if (setq name (defun2lambda node))
      (format t "~aNewSymbol(\"~s\"): &LispString{S: `~s`},~%"
              #\tab
              name
              node)
      )
    )
  (format t "}~%")
  )
; vim:set lispwords+=while:
