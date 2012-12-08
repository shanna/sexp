package sexp

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestBytes(t *testing.T) {
	exp := NewExpression(
		NewBytes([]byte("foo")),
		NewExpression(),
	)
	b := exp.MarshalSEXP(true)
	assert.Equal(t, "(3:foo())", string(b))
}

func TestAdvanced(t *testing.T) {
	exp := NewExpression(
		NewBytes([]byte("foo")),
		NewExpression(),
	)
	s := exp.MarshalSEXP(false)
	assert.Equal(t, "(\"foo\" ())", string(s))
}

func TestStringer(t *testing.T) {
	exp := NewExpression("\u2665s", []byte("bytes"), 12, 0.134)
	assert.Equal(t, "(4:\u2665s5:bytes2:125:0.134)", string(exp.MarshalSEXP(true)))
}

func TestParser(t *testing.T) {
	tree, err := Parse([]byte("(3:foo3:bar(3:baz))"))
	assert.Equal(t, nil, err)
	assert.Equal(t, "(3:foo3:bar(3:baz))", string(tree.MarshalSEXP(true)))
}
