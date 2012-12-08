package csexp

import (
	"fmt"
	"regexp"
	"strconv"
)

type Item struct {
	Type     ItemType
	Position int
	Value    []byte
}

type ItemType int

const (
	ItemError        ItemType = iota
	ItemBracketLeft  // (
	ItemBracketRight // )
	ItemBytes        // []byte
	ItemEOF
)

var (
	reBracketLeft  = regexp.MustCompile(`^\(`)
	reBracketRight = regexp.MustCompile(`^\)`)
	reBytesLength  = regexp.MustCompile(`^(\d+):`)
)

type stateFn func(*lexer) stateFn

type lexer struct {
	input   []byte
	items   chan Item
	start   int
	pos     int
	state   stateFn
	matches [][]byte
}

func (l *lexer) emit(t ItemType) {
	l.items <- Item{t, l.start, l.input[l.start:l.pos]}
}

func (l *lexer) Next() Item {
	item := <-l.items
	return item
}

func (l *lexer) scan(re *regexp.Regexp) bool {
	if l.match(re) {
		l.start = l.pos
		l.pos += len(l.matches[0])
		return true
	}
	return false
}

func (l *lexer) match(re *regexp.Regexp) bool {
	if l.matches = re.FindSubmatch(l.input[l.pos:]); l.matches != nil {
		return true
	}
	return false
}

func (l *lexer) run() {
	for l.state = lexCanonical; l.state != nil; {
		l.state = l.state(l)
	}
	close(l.items)
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- Item{ItemError, l.start, []byte(fmt.Sprintf(format, args...))}
	return nil
}

func lexCanonical(l *lexer) stateFn {
	if l.pos >= len(l.input) {
		l.emit(ItemEOF)
		return nil
	}
	if l.scan(reBracketLeft) {
		l.emit(ItemBracketLeft)
		return lexCanonical
	}
	if l.scan(reBracketRight) {
		l.emit(ItemBracketRight)
		return lexCanonical
	}
	if l.scan(reBytesLength) {
		bytes, _ := strconv.ParseInt(string(l.matches[1]), 10, 64)
		l.start = l.pos
		l.pos += int(bytes)
		l.emit(ItemBytes)

		return lexCanonical
	}
	return l.errorf("Expected expression.") // TODO: Better error.
}

func NewLexer(input []byte) *lexer {
	l := &lexer{input: input, items: make(chan Item)}
	go l.run()
	return l
}
