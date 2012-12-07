package csexp

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestBytes(t *testing.T) {
	exp := NewExpression(
		NewBytes([]byte("foo")),
		NewExpression(),
	)
	b := exp.Bytes()
	assert.Equal(t, "(3:foo())", string(b))
}

func TestString(t *testing.T) {
	exp := NewExpression(
		NewBytes([]byte("foo")),
		NewExpression(),
	)
	s := exp.String()
	assert.Equal(t, "(\"foo\" ())", s)
}

func TestStringer(t *testing.T) {
	exp := NewExpression("\u2665s", []byte("bytes"), 12, 0.134)
	assert.Equal(t, "(4:\u2665s5:bytes2:125:0.134)", string(exp.Bytes()))
}

func TestParser(t *testing.T) {
	tree, err := Parse([]byte("(3:foo3:bar(3:baz))"))
	assert.Equal(t, nil, err)
	assert.Equal(t, "(3:foo3:bar(3:baz))", string(tree.Bytes()))
}
