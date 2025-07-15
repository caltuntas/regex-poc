package main

type TokenType string

type Token struct {
	Type TokenType
	Value string
}

const (
	LITERAL = "LITERAL"
	STAR = "*"
	DOT = "."
	LBRACKET = "["
	RBRACKET = "]"
	ESCAPE = "\\"
	EOF = "EOF"
)

type Lexer struct {
	input string
	position int
	readPosition int
	ch byte
}

func New(pattern string) *Lexer {
	l := &Lexer{input: pattern}
	//l.readChar()
	return l
}

func isLetter(chr byte) bool {
	return chr >= 'a' && chr <= 'z' || chr >= 'A' && chr <= 'Z'
}

func (l *Lexer) NextChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}
func (l *Lexer) NextToken() Token {
	var token Token

	l.readChar()
	if isLetter(l.ch) {
		token.Type = LITERAL
		token.Value = string(l.ch)
	} else if l.ch == '.' {
		token.Type = DOT
		token.Value = string(l.ch)
	} else if l.ch == '*' {
		token.Type = STAR
		token.Value = string(l.ch)
	} else if l.ch == '[' {
		token.Type = LBRACKET
		token.Value = string(l.ch)
	} else if l.ch == ']' {
		token.Type = RBRACKET
		token.Value = string(l.ch)
	} else if l.ch == '\\' {
		token.Type = ESCAPE
		token.Value = string(l.ch)
	}

	return token
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

func (l *Lexer) HasMore() bool {
  return l.position < len(l.input)
}
