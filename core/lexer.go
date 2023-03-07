package core

import (
	"bytes"
	"errors"
	"fmt"
	"unicode"
	"unicode/utf8"
)

const eof = -1

var (
	typeToken         = []byte("type")
	queryToken        = []byte("query")
	mutationToken     = []byte("mutation")
	fragmentToken     = []byte("fragment")
	subscriptionToken = []byte("subscription")
	onToken           = []byte("on")
	trueToken         = []byte("true")
	falseToken        = []byte("false")
	quotesToken       = []byte(`'"`)
	signsToken        = []byte(`+-`)
	spreadToken       = []byte(`...`)
	digitToken        = []byte(`0123456789`)
	dotToken          = []byte(`.`)
	punctuatorToken   = []byte(`!():=[]{|}`)
)

// !()[]{}:=
var punctuator = map[rune]Symbol{
	'!': itemRequired,
	'{': itemObjOpen,
	'}': itemObjClose,
	'[': itemListOpen,
	']': itemListClose,
	'(': itemArgsOpen,
	')': itemArgsClose,
	':': itemColon,
	'=': itemEquals,
}

type stateFn func(*lexer) stateFn
type Symbol int8

const (
	itemError      Symbol = iota // error
	itemEOF                      // end of file
	itemName                     // label
	itemOn                       // "on"
	itemArgsOpen                 // (
	itemArgsClose                // )
	itemListOpen                 // [
	itemListClose                // ]
	itemObjOpen                  // {
	itemObjClose                 // }
	itemColon                    // :
	itemEquals                   // =
	itemRequired                 // !
	itemDirective                // @(directive)
	itemVariable                 // $variable
	itemSpread                   // ...
	itemNumber                   // number
	itemString                   // string
	itemBoolean                  // boolean
	itemUnion                    // |
	itemPunctuator               // punctuation !()[]{}:=
)

type item struct {
	symbol Symbol // The type of this item.
	column int    // The starting column, in bytes, of this item in the input string.
	value  []byte // The value of this item.
	line   int16  // The line number at the start of this item.
}

func (my item) String() string {
	var v string
	switch my.symbol {
	case itemEOF:
		v = "eof"
	case itemError:
		v = "error"
	case itemName:
		v = "name"
	case itemPunctuator:
		v = "punctuator"
	case itemDirective:
		v = "directive"
	case itemVariable:
		v = "variable"
	case itemNumber:
		v = "number"
	case itemString:
		v = "string"
	}
	return v
}

type lexer struct {
	input  []byte // the string being scanned
	cursor int    // current cursor in the input
	start  int    // start cursor of this item
	width  int    // width of last rune read from input
	items  []item // array of scanned items
	line   int16  // 1+number of newlines seen
	err    error
}

func newLexer(input []byte) (lexer, error) {
	var l lexer
	if len(input) == 0 {
		return l, errors.New("query string can't be empty")
	}
	l.input = input
	l.line = 1
	l.run()
	return l, nil
}

// next returns the next byte in the input.
func (my *lexer) next() rune {
	if my.cursor >= len(my.input) {
		my.width = 0
		return eof
	}
	r, w := utf8.DecodeRune(my.input[my.cursor:])
	my.width = w
	my.cursor += my.width
	if r == '\n' {
		my.line++
	}
	return r
}

// peek returns but does not consume the next rune in the input.
func (my *lexer) peek() rune {
	r := my.next()
	my.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (my *lexer) backup() {
	if my.cursor != 0 {
		my.cursor -= my.width
		// Correct newline count.
		if my.width == 1 && my.input[my.cursor] == '\n' {
			my.line--
		}
	}
}

func (my *lexer) current() []byte {
	return my.input[my.start:my.cursor]
}

func (my *lexer) run() {
	for s := lexRoot; s != nil; {
		s = s(my)
	}
}

func (my *lexer) errorf(format string, args ...interface{}) stateFn {
	my.err = fmt.Errorf(format, args...)
	my.items = append(my.items, item{itemError, my.start, my.current(), my.line})
	return nil
}

// ignore skips over the pending input before this point.
func (my *lexer) ignore() {
	my.start = my.cursor
}

// emit passes an item back to the client.
func (my *lexer) emit(t Symbol, options ...func([]byte)) {
	for _, fn := range options {
		fn(my.current())
	}
	my.items = append(my.items, item{t, my.start, my.current(), my.line})
	// Some items contain text internally. If so, count their newlines.
	if t == itemString {
		for i := my.start; i < my.cursor; i++ {
			if my.input[i] == '\n' {
				my.line++
			}
		}
	}
	my.start = my.cursor
}

// accept consumes the next rune if it's from the valid set.
func (my *lexer) accept(valid []byte) (rune, bool) {
	r := my.next()
	if r != eof && bytes.ContainsRune(valid, r) {
		return r, true
	}
	my.backup()
	return r, false
}

// acceptRun consumes a run of runes from the valid set.
func (my *lexer) acceptRun(valid []byte) {
	for bytes.ContainsRune(valid, my.next()) {

	}
	my.backup()
}

// acceptComment consumes a run of runes while till the end of line
func (my *lexer) acceptComment() {
	n := 0
	for r := my.next(); !isEndOfLine(r); r = my.next() {
		n++
	}
}

// acceptAlphaNum consumes a run of runes while they are alpha nums
func (my *lexer) acceptAlphaNum() bool {
	n := 0
	for r := my.next(); isAlphaNumeric(r); r = my.next() {
		n++
	}
	my.backup()
	return n != 0
}

// lexInsideAction scans the elements inside action delimiters.
func lexRoot(l *lexer) stateFn {
	r := l.next()

	switch {
	case r == eof:
		l.emit(itemEOF)
		return nil
	case isEndOfLine(r):
		l.ignore()
	case isSpace(r):
		l.ignore()
	case r == '#':
		l.ignore()
		l.acceptComment()
		l.ignore()
	case r == '@':
		l.ignore()
		l.emit(itemDirective)
	case r == '!':
		l.ignore()
		l.emit(itemRequired)
	case r == '$':
		l.ignore()
		if l.acceptAlphaNum() {
			l.emit(itemVariable)
		}
	case contains(l.current(), punctuatorToken):
		if item, ok := punctuator[r]; ok {
			l.emit(item)
		}
	case r == '"' || r == '\'':
		l.backup()
		return lexString
	case r == '.':
		l.accept(dotToken)
		if equals(l.current(), spreadToken) {
			l.emit(itemSpread)
			return lexRoot
		}
		//fallthrough // '.' can start a number.
	case r == '+' || r == '-' || ('0' <= r && r <= '9'):
		l.backup()
		return lexNumber
	case isAlphaNumeric(r):
		l.backup()
		return lexName
	default:
		return l.errorf("unrecognized character in action: %#U", r)
	}
	return lexRoot
}

func (my *lexer) scanNumber() bool {
	// Optional leading sign.
	my.accept(signsToken)
	my.acceptRun(digitToken)
	if _, ok := my.accept(dotToken); ok {
		my.acceptRun(digitToken)
	}
	// Is it imaginary?
	my.accept([]byte("i"))
	// Next thing mustn't be alphanumeric.
	if isAlphaNumeric(my.peek()) {
		my.next()
		return false
	}
	return true
}

// lexNumber scans a number: decimal and float. This isn't a perfect number scanner
// for instance it accepts "." and "0x0.2" and "089" - but when it's wrong the input
// is invalid and the parser (via strconv) should notice.
func lexNumber(l *lexer) stateFn {
	if !l.scanNumber() {
		return l.errorf("bad number syntax: %q", l.input[l.start:l.cursor])
	}
	l.emit(itemNumber)
	return lexRoot
}

// lexName scans a name.
func lexName(l *lexer) stateFn {
	for {
		r := l.next()
		if r == eof {
			l.emit(itemEOF)
			return nil
		}
		if !isAlphaNumeric(r) {
			l.backup()
			val := l.current()
			switch {
			case equals(val, onToken):
				l.emit(itemOn, lowercase)
			case equals(val, trueToken):
				l.emit(itemBoolean, lowercase)
			case equals(val, falseToken):
				l.emit(itemBoolean, lowercase)
			default:
				l.emit(itemName)
			}
			break
		}
	}
	return lexRoot
}

// lexString scans a string.
func lexString(l *lexer) stateFn {
	if sr, ok := l.accept([]byte(quotesToken)); ok {
		l.ignore()
		var escaped bool
		for {
			r := l.next()
			if r == eof {
				l.emit(itemEOF)
				return nil
			}
			if r == '\\' {
				escaped = true
				continue
			}
			if !escaped && r == sr {
				l.backup()
				l.emit(itemString)
				if _, ok := l.accept(quotesToken); ok {
					l.ignore()
				}
				break
			}
			if escaped {
				escaped = false
			}
		}
	}
	return lexRoot
}

func equals(b, val []byte) bool {
	return bytes.EqualFold(b, val)
}

func contains(b []byte, chars []byte) bool {
	return bytes.ContainsAny(b, string(chars))
}

func lowercase(b []byte) {
	for i := 0; i < len(b); i++ {
		if b[i] >= 'A' && b[i] <= 'Z' {
			b[i] = 'a' + (b[i] - 'A')
		}
	}
}

// isSpace reports whether r is a space character.
func isSpace(r rune) bool {
	return r == ',' || r == ' ' || r == '\t'
}

// isEndOfLine reports whether r is an end-of-line character.
func isEndOfLine(r rune) bool {
	return r == '\r' || r == '\n' || r == eof
}

// isAlphaNumeric reports whether r is an alphabetic, digit, or underscore.
func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}
