#lang racket
(require racket/struct)

(struct monkey (items operation test inspected)
  #:mutable
  #:methods gen:custom-write
  [(define write-proc
     (make-constructor-style-printer
      (lambda (obj) 'monkey:)
      (lambda (obj) (list (monkey-items obj) (monkey-operation obj) (monkey-test obj) (monkey-inspected obj)))))])

(struct operation (operator operand)
  #:methods gen:custom-write
  [(define write-proc
     (make-constructor-style-printer
      (lambda (obj) 'operation:)
      (lambda (obj) (list (operation-operator obj) (operation-operand obj)))))])

(struct test (divisor if-true if-false)
  #:methods gen:custom-write
  [(define write-proc
     (make-constructor-style-printer
      (lambda (obj) 'test:)
      (lambda (obj) (list (test-divisor obj) (test-if-true obj) (test-if-false obj)))))])

(define (parse-starting-items line)
  (if (not (string-prefix? line "  Starting items: "))
      (raise (string-append-immutable "error parsing starting items: '" line "'"))
      (map
       (lambda(num-str) (string->number num-str))
       (string-split (substring line 18) ", "))))

(define (parse-operation line)
  (if (not (string-prefix? line "  Operation: new = old "))
      (raise "error parsing operation" line)
      ((lambda()
         (define op (string-ref line 23))
         (define operand (if (equal? "old" (substring line 25)) "old" (string->number (substring line 25))))
         (operation op operand)))))

(define (parse-num-end-str prefix)
  (lambda(line)
    (if (not (string-prefix? line prefix))
        (raise "error parsing" line)
        (string->number (substring line (string-length prefix))))))

(define parse-test-condition (parse-num-end-str "  Test: divisible by "))
(define parse-test-if-true (parse-num-end-str "    If true: throw to monkey "))
(define parse-test-if-false (parse-num-end-str "    If false: throw to monkey "))

(define (parse-test line1 line2 line3)
  (test
   (parse-test-condition line1)
   (parse-test-if-true line2)
   (parse-test-if-false line3)))

(define (read-record)
  (monkey (parse-starting-items (read-line))
          (parse-operation (read-line))
          (parse-test (read-line) (read-line) (read-line))
          0))

(define input-file-path "input.txt")

(define (read-record-loop)
  (define (loop records)
    (if (not (eq? (peek-char) eof))
        ((lambda()
           (read-line)
           (define record (read-record))
           (read-line)
           (loop (append records (list record)))))
        records))
  (loop '()))

(define (or-old operand worry) (if (equal? "old" operand) worry operand))

(define (apply-op oper worry)
  (define op (operation-operator oper))
  (define operand (operation-operand oper))
  (case op
    [(#\+)
     (+ (or-old operand worry) worry)]
    [(#\*)
     (* (or-old operand worry) worry)]
    [else (raise "unknown operator" (op))]))

(define monkeys '())
(define common-multiplier 0)

(define (load-monkeys)
  (set! monkeys
        (with-input-from-file input-file-path
          (lambda ()
            (read-record-loop))))
  (set! common-multiplier
        (for/product ([m monkeys])
          (test-divisor (monkey-test m)))))

(define (do-round monkeys divisor)
  (for ([m monkeys])
    (monkey-inspect m divisor)
    (monkey-test-throw m)));

(define (monkey-inspect m divisor)
  (set-monkey-items!
   m
   (map
    (lambda (worry)
      (set-monkey-inspected! m (add1 (monkey-inspected m)))
      (quotient (apply-op (monkey-operation m) worry) divisor)
      )
    (monkey-items m))));

(define (monkey-test-throw m)
  ; (println m)
  (define test (monkey-test m))
  (define divisor (test-divisor test))
  (define if-true-monkey (list-ref monkeys (test-if-true test)))
  (define if-false-monkey (list-ref monkeys (test-if-false test)))
  (define items (monkey-items m))
  (for ([item items])
    (if (zero? (modulo item divisor))
        (set-monkey-items! if-true-monkey (append (monkey-items if-true-monkey) (list item )))
        (set-monkey-items! if-false-monkey (append (monkey-items if-false-monkey) (list item)))
        ))
  (set-monkey-items! m '())
  )

(define (stress-relief)
  (define all-items '())
  (for-each (lambda (m) (set! all-items (append all-items (monkey-items m)))) monkeys)
  (for-each
   (lambda (m)
     (set-monkey-items! m (map (lambda(item) (modulo item common-multiplier)) (monkey-items m))))
   monkeys)
  )

(define (collect-inspected) (map (lambda (m) (monkey-inspected m)) monkeys))

(define (print-ans)
  (define inspect-counts (collect-inspected))
  (println inspect-counts)
  (define sorted-list (sort inspect-counts >))
  (println sorted-list)
  (println (* (list-ref sorted-list 0) (list-ref sorted-list 1))))

(load-monkeys)

(for ([i 20])
  (do-round monkeys 3))

(print-ans)

(load-monkeys)

(for ([i 10000])
  (do-round monkeys 1)
  (stress-relief))

(print-ans)