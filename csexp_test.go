package csexp

import (
  "github.com/bmizerany/assert"
  "testing"
)

func TestLexer(t *testing.T) {
  l := lex([]byte("(3:foo3:bar(3:baz))"))

  item := l.next()
  assert.Equal(t, []byte("("), item.val)
  assert.Equal(t, itemBracketLeft, item.typ)

  item = l.next()
  assert.Equal(t, "3:", string(item.val))
  assert.Equal(t, itemBytesLength, item.typ)

  item = l.next()
  assert.Equal(t, "foo", string(item.val))
  assert.Equal(t, itemBytes, item.typ)

  item = l.next()
  assert.Equal(t, "3:", string(item.val))
  assert.Equal(t, itemBytesLength, item.typ)

  item = l.next()
  assert.Equal(t, "bar", string(item.val))
  assert.Equal(t, itemBytes, item.typ)

  item = l.next()
  assert.Equal(t, "(", string(item.val))
  assert.Equal(t, itemBracketLeft, item.typ)

  item = l.next()
  assert.Equal(t, "3:", string(item.val))
  assert.Equal(t, itemBytesLength, item.typ)

  item = l.next()
  assert.Equal(t, "baz", string(item.val))
  assert.Equal(t, itemBytes, item.typ)

  item = l.next()
  assert.Equal(t, ")", string(item.val))
  assert.Equal(t, itemBracketRight, item.typ)

  item = l.next()
  assert.Equal(t, ")", string(item.val))
  assert.Equal(t, itemBracketRight, item.typ)

  item = l.next()
  assert.Equal(t, itemEOF, item.typ)
}

