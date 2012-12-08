package csexp

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type Kind int

const (
	AtomBytes Kind = iota
	AtomExpression
)

type Atomizer interface {
	Kind() Kind
	Bytes() []byte
	String() string
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

func (bytes *Bytes) Bytes() []byte {
	return []byte(fmt.Sprintf("%d:%s", len(bytes.Value), bytes.Value))
}

func (bytes *Bytes) String() string {
	return strconv.Quote(string(bytes.Value))
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

func (expression *Expression) Bytes() []byte {
	exp := new(bytes.Buffer)
	exp.WriteString("(")
	for _, csexp := range *expression {
		val := csexp.Bytes()
		exp.Write(val)
	}
	exp.WriteString(")")
	return exp.Bytes()
}

func (expression *Expression) String() string {
	exp := new(bytes.Buffer)
	exp.WriteString("(")
	atoms := []string{}
	for _, csexp := range *expression {
		atoms = append(atoms, csexp.String())
	}
	exp.WriteString(strings.Join(atoms, " "))
	exp.WriteString(")")
	return exp.String()
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
		case ItemBytes:
			stack[len(stack)-1].Push(item.Value)
		default:
			panic("unreachable")
		}
	}
	return (*stack[0])[0].(*Expression), nil
}
