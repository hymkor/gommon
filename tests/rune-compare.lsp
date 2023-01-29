;;; test for (rune)
(test (char< #\a #\b) t)
(test (char> #\a #\b) nil)
(test (char<= #\a #\b) t)
(test (char>= #\a #\b) nil)
(test (char= #\a #\b) nil)
(test (char/= #\a #\b) t)

(test (characterp #\a) t)
(test (characterp 'a) nil)
(test (characterp "2") nil)
