package csexp

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestParser(t *testing.T) {
	tree := Parse([]byte("(3:foo3:bar(3:baz))"))
	assert.Equal(t, "(3:foo3:bar(3:baz))", string(tree.Bytes()))
}
