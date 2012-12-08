package csexp

import (
	"github.com/bmizerany/assert"
	"testing"
)

func TestLexer(t *testing.T) {
	l := NewLexer([]byte("(3:foo3:bar(3:baz))"))

	item := l.Next()
	assert.Equal(t, []byte("("), item.Value)
	assert.Equal(t, ItemBracketLeft, item.Type)

	item = l.Next()
	assert.Equal(t, "foo", string(item.Value))
	assert.Equal(t, ItemBytes, item.Type)

	item = l.Next()
	assert.Equal(t, "bar", string(item.Value))
	assert.Equal(t, ItemBytes, item.Type)

	item = l.Next()
	assert.Equal(t, "(", string(item.Value))
	assert.Equal(t, ItemBracketLeft, item.Type)

	item = l.Next()
	assert.Equal(t, "baz", string(item.Value))
	assert.Equal(t, ItemBytes, item.Type)

	item = l.Next()
	assert.Equal(t, ")", string(item.Value))
	assert.Equal(t, ItemBracketRight, item.Type)

	item = l.Next()
	assert.Equal(t, ")", string(item.Value))
	assert.Equal(t, ItemBracketRight, item.Type)

	item = l.Next()
	assert.Equal(t, ItemEOF, item.Type)
}
