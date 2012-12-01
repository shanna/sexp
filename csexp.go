package csexp

import (
  "fmt"
  "regexp"
  "strconv"
)

type item struct {
  typ itemType
  pos int
  val []byte
}

type itemType int

const (
  itemError itemType = iota
  itemBracketLeft  // (
  itemBracketRight // )
  itemBytes        // []byte
  itemBytesLength  // Integer.
  itemEOF
)

var (
  reBracketLeft  = regexp.MustCompile(`^[(]`)
  reBracketRight = regexp.MustCompile(`^[)]`)
  reBytesLength  = regexp.MustCompile(`^(\d+):`)
)

const eof = -1

type stateFn func(*lexer) stateFn

type lexer struct {
  input   []byte
  items   chan item
  start   int
  pos     int
  state   stateFn
  matches [][]byte
}

func (l *lexer) emit(t itemType) {
  l.items <- item{t, l.start, l.input[l.start:l.pos]}
}

func (l *lexer) next() item {
  item := <-l.items
  return item
}

// Match but advance pointer on match.
func (l *lexer) scan(re *regexp.Regexp) bool {
  if l.matches = re.FindSubmatch(l.input[l.pos:]); l.matches != nil {
    l.start = l.pos
    l.pos  += len(l.matches[0])
    return true
  }
  return false
}

// Match but don't advance pointer.
func (l *lexer) match(re *regexp.Regexp) bool {
  return re.Match(l.input[l.pos:])
}

func (l *lexer) run() {
  for l.state = lexBracketLeft; l.state != nil; {
    l.state = l.state(l)
  }
  close(l.items)
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
  l.items <- item{itemError, l.start, []byte(fmt.Sprintf(format, args...))}
  return nil
}

func lex(input []byte) *lexer {
  l := &lexer{
    input: input,
    items: make(chan item),
  }
  go l.run()
  return l
}

func lexBracketLeft(l *lexer) stateFn {
  if l.pos == len(l.input) {
    l.emit(itemEOF)
    return nil
  }

  if l.scan(reBracketLeft) {
    l.emit(itemBracketLeft)
    return lexExpression
  }

  return l.errorf("Expected '('.")
}

func lexBracketRight(l *lexer) stateFn {
  if l.scan(reBracketRight) {
    l.emit(itemBracketRight)
    return lexBracketLeft
  }
  return l.errorf("Expected ')'.")
}

func lexExpression(l *lexer) stateFn {
  if l.match(reBracketLeft) {
    return lexBracketLeft
  }

  if l.match(reBracketRight) {
    return lexBracketRight
  }

  if l.scan(reBytesLength) {
    bytes, _ := strconv.ParseInt(string(l.matches[1]), 10, 64)
    l.emit(itemBytesLength)

    // Take size of bytes from input.
    // TODO: Complain long and hard if it's out of range.
    l.start = l.pos
    l.pos += int(bytes)
    l.emit(itemBytes)

    return lexExpression
  }
  return l.errorf("Expected expression.")
}
