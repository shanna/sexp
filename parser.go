package csexp

import (
	"bytes"
	"fmt"
)

type Kind int

const (
	AtomBytes Kind = iota
	AtomExpression
)

type Atomizer interface {
	Kind() Kind
	Bytes() []byte
}

// Bytes Atom
type Bytes struct {
	Value []byte
}

func NewBytes(value []byte) Bytes {
	return Bytes{value}
}

func (a *Bytes) Bytes() []byte {
	return []byte(fmt.Sprintf("%d:%s", len(a.Value), a.Value))
}

func (a *Bytes) Kind() Kind {
	return AtomBytes
}

// (Sub)Expression Atom
type Expression struct {
	Value []Atomizer
}

func NewExpression(values ...Atomizer) Expression {
	return Expression{values}
}

func (a *Expression) Bytes() []byte {
	st := new(bytes.Buffer)
	st.WriteString("(")
	for _, atom := range a.Value {
		st.Write(atom.Bytes())
	}
	st.WriteString(")")
	return st.Bytes()
}

func (a *Expression) Kind() Kind {
	return AtomExpression
}

// Parse a Canonical S-expression from a byte slice.
func Parse(data []byte) *Expression {
	l := newLexer(data)
	s := []*Expression{&Expression{}}

	for item := l.next(); item.typ != itemEOF; item = l.next() {
		switch item.typ {
		case itemBracketLeft:
			e := NewExpression()
			s = append(s, &e)
			s[len(s)-2].Value = append(s[len(s)-2].Value, s[len(s)-1])
		case itemBracketRight:
			s = s[:len(s)-1]
		case itemBytes:
			b := NewBytes(item.val)
			s[len(s)-1].Value = append(s[len(s)-1].Value, &b)
		default:
			panic("unreachable")
		}
	}

	return s[0].Value[0].(*Expression)
}
