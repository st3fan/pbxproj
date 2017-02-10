package parser

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
)

const (
	EndOfFile = iota + 1
	EndOfLine

	Whitespace
	Comment

	Header
	Identifier

	OpenCurly
	CloseCurly
	Equals
	Semicolon
	Comma
	OpenParen
	CloseParen

	Illegal
)

type Token struct {
	Type    int
	Literal string
	Line    int
	Pos     int
}

func (t Token) String() string {
	switch t.Type {
	case EndOfFile:
		return "Token(EndOfFile)"
	case EndOfLine:
		return "Token(EndOfLine)"
	case Whitespace:
		return "Token(Whitespace)"
	case Comment:
		return fmt.Sprintf("Token(Comment, %s)", strconv.Quote(t.Literal))
	case Header:
		return "Token(Header)"
	case Identifier:
		return fmt.Sprintf("Token(Identifier, %s)", strconv.Quote(t.Literal))
	case OpenCurly:
		return "Token(OpenCurly)"
	case CloseCurly:
		return "Token(CloseCurly)"
	case OpenParen:
		return "Token(OpenParen)"
	case CloseParen:
		return "Token(CloseParen)"
	case Equals:
		return "Token(Equals)"
	case Semicolon:
		return "Token(Semicolon)"
	case Comma:
		return "Token(Comma)"
	}

	return fmt.Sprintf("Token(UNKNOWN, %d, %s)", t.Type, strconv.Quote(t.Literal))
}

type Tokenizer struct {
	scanner *bufio.Scanner
	line    int
	buffer  string
	index   int
}

var eof = rune(0)

func (t *Tokenizer) read() rune {
	if t.index == len(t.buffer) {
		return eof
	}
	c := t.buffer[t.index]
	t.index += 1
	return rune(c)
}

func (t *Tokenizer) unread() {
	if t.index > 0 {
		t.index -= 1
	}
}

func (t *Tokenizer) peek() (rune, rune) {
	if len(t.buffer)-t.index >= 2 {
		return rune(t.buffer[t.index]), rune(t.buffer[t.index+1])
	} else if len(t.buffer)-t.index >= 1 {
		return rune(t.buffer[t.index]), eof
	}
	return eof, eof
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func (t *Tokenizer) peekWhitespace() bool {
	c, _ := t.peek()
	return isWhitespace(c)
}

func (t *Tokenizer) scanWhitespace() (Token, error) {
	var buf bytes.Buffer
	buf.WriteRune(t.read())

	for {
		if ch := t.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			t.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return Token{Type: Whitespace, Literal: buf.String(), Line: t.line}, nil
}

func isIdentifierStartChar(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') ||
		c == '$' || c == '/' || c == '.' || c == '_' || c == '-'
}

func isIdentifierChar(c rune) bool {
	return isIdentifierStartChar(c) || c == '_' || c == '-'
}

func (t *Tokenizer) peekIdentifier() bool {
	c, _ := t.peek()
	return isIdentifierStartChar(c)
}

func (t *Tokenizer) scanIdentifier() (Token, error) {
	var buf bytes.Buffer
	buf.WriteRune(t.read())

	for {
		if ch := t.read(); ch == eof {
			break
		} else if !isIdentifierChar(ch) {
			t.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return Token{Type: Identifier, Literal: buf.String(), Line: t.line}, nil
}

func (t *Tokenizer) peekComment() bool {
	c1, c2 := t.peek()
	return c1 == '/' && c2 == '*'
}

func (t *Tokenizer) scanComment() (Token, error) {
	var buf bytes.Buffer
	buf.WriteRune(t.read())

	c1 := t.read()
	if c1 == eof {
		return Token{Line: t.line}, errors.New("Unexpected")
	}
	buf.WriteRune(c1)

	c2 := t.read()
	if c2 == eof {
		return Token{Line: t.line}, errors.New("Unexpected")
	}
	buf.WriteRune(c2)

	for {
		ch := t.read()
		if ch == eof {
			break // Error
		}
		buf.WriteRune(ch)

		if buf.Len() >= 4 {
			bytes := buf.Bytes()
			if bytes[len(bytes)-2] == '*' && bytes[len(bytes)-1] == '/' {
				return Token{Type: Comment, Literal: buf.String(), Line: t.line}, nil
			}
		}
	}

	return Token{Line: t.line}, fmt.Errorf("Non closed comment at %d:%d", t.line, t.index)
}

func (t *Tokenizer) peekString() bool {
	c, _ := t.peek()
	return c == '"'
}

func (t *Tokenizer) scanString() (Token, error) {
	var buf bytes.Buffer
	buf.WriteRune(t.read())

	for {
		ch := t.read()
		if ch == eof {
			break // Error
		}
		buf.WriteRune(ch)

		if ch == '"' {
			bytes := buf.Bytes()
			if bytes[len(bytes)-2] != '\\' {
				return Token{Type: Identifier, Literal: buf.String(), Line: t.line}, nil
			}
		}
	}

	return Token{Line: t.line}, fmt.Errorf("Non closed string at %d:%d", t.line, t.index)
}

func NewTokenizer(r io.Reader) (*Tokenizer, error) {
	return &Tokenizer{scanner: bufio.NewScanner(r), index: -1}, nil
}

func (t *Tokenizer) Next() (Token, error) {
	if t.index == -1 {
		if !t.scanner.Scan() {
			if err := t.scanner.Err(); err != nil {
				return Token{Line: t.line}, err
			} else {
				return Token{Type: EndOfFile, Line: t.line}, nil
			}
		}

		t.line += 1
		t.buffer = t.scanner.Text()
		t.index = 0
	}

	if t.buffer == "// !$*UTF8*$!" {
		t.index = -1
		return Token{Type: Header, Literal: t.buffer, Line: t.line}, nil
	}

	if t.index == len(t.buffer) {
		t.index = -1
		return Token{Type: EndOfLine, Literal: "", Line: t.line}, nil
	}

	if t.peekWhitespace() {
		return t.scanWhitespace()
	}

	if t.peekComment() {
		return t.scanComment()
	}

	if t.peekIdentifier() {
		return t.scanIdentifier()
	}

	if t.peekString() {
		return t.scanString()
	}

	c := t.read()

	switch c {
	case '{':
		return Token{Type: OpenCurly, Line: t.line}, nil
	case '}':
		return Token{Type: CloseCurly, Line: t.line}, nil
	case '(':
		return Token{Type: OpenParen, Line: t.line}, nil
	case ')':
		return Token{Type: CloseParen, Line: t.line}, nil
	case '=':
		return Token{Type: Equals, Line: t.line}, nil
	case ';':
		return Token{Type: Semicolon, Line: t.line}, nil
	case ',':
		return Token{Type: Comma, Line: t.line}, nil
	}

	return Token{Type: Illegal, Literal: "", Line: t.line}, fmt.Errorf("Unknown token at %d:%d", t.line, t.index)
}
