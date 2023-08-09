package gmnlisp

import (
	"io"

	"github.com/hymkor/gmnlisp/pkg/parser"
)

var (
	ampRest         = NewSymbol("&rest")
	commaSymbol     = NewSymbol(",")
	nulSymbol       = NewSymbol("")
	quoteSymbol     = NewSymbol("quote")
	backQuoteSymbol = NewSymbol("backquote")
	slashSymbol     = NewSymbol("/")
	colonRest       = NewKeyword(":rest")
)

type stdFactory struct{}

func (stdFactory) Cons(car, cdr Node) Node           { return &Cons{Car: car, Cdr: cdr} }
func (stdFactory) Int(n int64) Node                  { return Integer(n) }
func (stdFactory) Float(f float64) Node              { return Float(f) }
func (stdFactory) String(s string) Node              { return String(s) }
func (stdFactory) Array(list []Node, dim []int) Node { return &Array{list: list, dim: dim} }
func (stdFactory) Keyword(s string) Node             { return NewKeyword(s) }
func (stdFactory) Rune(r rune) Node                  { return Rune(r) }
func (stdFactory) Symbol(s string) Node              { return NewSymbol(s) }
func (stdFactory) Null() Node                        { return Null }
func (stdFactory) True() Node                        { return True }

func ReadNode(rs io.RuneScanner) (Node, error) {
	return parser.Read[Node](stdFactory{}, rs)
}

func tryParseAsNumber(token string) (Node, bool, error) {
	return parser.TryParseAsNumber[Node](stdFactory{}, token)
}

func newQuote(value Node) Node {
	return parser.NewQuote[Node](stdFactory{}, value)
}

func ReadAll(rs io.RuneScanner) ([]Node, error) {
	result := []Node{}
	for {
		token, err := ReadNode(rs)
		if err != nil {
			if err == io.EOF {
				return result, nil
			}
			return nil, err
		}
		result = append(result, token)
	}
}
