package sexp

import (
	"bytes"
	"fmt"
	"strconv"
)

type Kind int

const (
	AtomBytes Kind = iota
	AtomExpression
)

type Atomizer interface {
	Kind() Kind
	MarshalSEXP(canonical bool) []byte
}

type Bytes struct {
	// Encoding []byte
	Value []byte
}

func NewBytes(data []byte) *Bytes {
	return &Bytes{data}
}

func (bytes *Bytes) Kind() Kind {
	return AtomBytes
}

func (bytes *Bytes) MarshalSEXP(canonical bool) []byte {
	if canonical {
		return []byte(fmt.Sprintf("%d:%s", len(bytes.Value), bytes.Value))
	}
	return []byte(strconv.Quote(string(bytes.Value)))
}

type Expression []Atomizer

func NewExpression(atoms ...interface{}) *Expression {
	e := &Expression{}
	e.Push(atoms...)
	return e
}

func (a *Expression) Kind() Kind {
	return AtomExpression
}

func (expression *Expression) MarshalSEXP(canonical bool) []byte {
	exp := new(bytes.Buffer)
	exp.WriteString("(")

	// WTB Ternary.
	separator := " "
	if canonical {
		separator = ""
	}

	atoms := [][]byte{}
	for _, sexp := range *expression {
		atoms = append(atoms, sexp.MarshalSEXP(canonical))
	}
	exp.Write(bytes.Join(atoms, []byte(separator)))

	exp.WriteString(")")
	return exp.Bytes()
}

// Convenience method to push Bytes atoms or cast with fmt.Sprint() into Bytes and push.
func (expression *Expression) Push(data ...interface{}) {
	for _, d := range data {
		if atomizer, ok := d.(Atomizer); ok {
			*expression = append(*expression, atomizer)
		} else if bytes, ok := d.([]byte); ok {
			*expression = append(*expression, NewBytes(bytes))
		} else {
			*expression = append(*expression, NewBytes([]byte(fmt.Sprint(d))))
		}
	}
}

// TODO: Stack stuff is ugly.
func Parse(data []byte) (*Expression, error) {
	lexer := NewLexer(data)
	stack := []*Expression{&Expression{}}

	for item := lexer.Next(); item.Type != ItemEOF; item = lexer.Next() {
		switch item.Type {
		case ItemBracketLeft:
			ex := NewExpression()
			stack[len(stack)-1].Push(ex)
			stack = append(stack, ex)
		case ItemBracketRight:
			stack = stack[:len(stack)-1]
		case ItemToken, ItemQuote, ItemVerbatim:
			stack[len(stack)-1].Push(item.Value)
		default:
			panic("unreachable")
		}
	}
	return (*stack[0])[0].(*Expression), nil
}
