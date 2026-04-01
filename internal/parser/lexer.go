package parser

import (
	"strings"
	"unicode"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) skipWhitespace() {
	// IMPORTANT: we do NOT skip '\n'
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	switch l.ch {
	case '\n':
		l.readChar()
		return Token{Type: NEWLINE, Literal: "\\n"}

	case '@':
		l.readChar()
		return l.readItem()

	case ':':
		l.readChar()
		return l.readIdent()

	case '-':
		l.readChar()
		return Token{Type: MOVE, Literal: "-"}

	case '/':
		l.readChar()
		l.skipWhitespace()
		return l.readIdent()

	case 0:
		return Token{Type: EOF, Literal: ""}

	default:
		if isWordChar(l.ch) {
			return l.readIdent()
		}
	}

	tok := Token{Type: ILLEGAL, Literal: string(l.ch)}
	l.readChar()
	return tok
}

func (l *Lexer) readIdent() Token {
	start := l.position

	for isWordChar(l.ch) || l.ch == ' ' {
		l.readChar()
	}

	lit := l.input[start:l.position]
	lit = normalize(lit)

	res := Token{
		Type:    IDENT,
		Literal: lit,
	}

	if token, ok := literalMap[res.Literal]; ok {
		return token
	}

	return res
}

func (l *Lexer) readItem() Token {
	l.skipWhitespace()

	start := l.position

	for isWordChar(l.ch) || l.ch == ' ' {
		if l.ch == '\n' || l.ch == 0 {
			break
		}
		l.readChar()
	}

	lit := l.input[start:l.position]
	lit = normalize(lit)

	return Token{
		Type:    ITEM,
		Literal: lit,
	}
}

func isWordChar(ch byte) bool {
	return unicode.IsLetter(rune(ch)) ||
		unicode.IsDigit(rune(ch)) ||
		ch == '-' ||
		ch == '\'' ||
		ch == '.'
}

func normalize(s string) string {
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, "'", "")
	s = strings.TrimSpace(s)
	return s
}
