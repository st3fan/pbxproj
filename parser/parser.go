package parser

import (
	"errors"
	"fmt"
	"io"
)

type Parser struct {
	tokenizer *Tokenizer
	lookahead *Token
}

type Value interface {
	Something() string
}

//

type String struct {
	String string
}

func (s String) Something() string {
	return s.String
}

//

type Array struct {
	Values []Value
}

func (a Array) Something() string {
	return "TODO"
}

func (a *Array) Append(value Value) {
	a.Values = append(a.Values, value)
}

func (a Array) Count() int {
	return len(a.Values)
}

//

type Dictionary struct {
	Values map[string]*Value
}

func (d Dictionary) Something() string {
	return "TODO"
}

func (d *Dictionary) Add(key string, value Value) {
	d.Values[key] = &value
}

func (d *Dictionary) Count() int {
	return len(d.Values)
}

func NewParser(r io.Reader) (*Parser, error) {
	tokenizer, err := NewTokenizer(r)
	if err != nil {
		return nil, err
	}
	return &Parser{tokenizer: tokenizer}, nil
}

func (p *Parser) peek() (Token, error) {
	if p.lookahead != nil {
		return *p.lookahead, nil
	}
	token, err := p.token()
	if err != nil {
		return Token{}, err
	}
	p.lookahead = &token
	return token, err
}

func (p *Parser) token() (Token, error) {
	if p.lookahead != nil {
		token := *p.lookahead
		p.lookahead = nil
		return token, nil
	}

	for {
		token, err := p.tokenizer.Next()
		if err != nil {
			return token, err
		}

		if token.Type != Whitespace && token.Type != Comment && token.Type != EndOfLine {
			return token, err
		}
	}
}

func (p *Parser) expect(tokenType int) (Token, error) {
	token, err := p.token()
	if err != nil {
		return Token{}, err
	}
	if token.Type != tokenType {
		return Token{}, fmt.Errorf("Expected %d but got %d", tokenType, token.Type)
	}
	return token, nil
}

func (p *Parser) parseHeader() error {
	token, err := p.token()
	if err != nil {
		return err
	}

	if token.Type != Header {
		return fmt.Errorf("Expected Header but got: %s", token)
	}

	return nil
}

func (p *Parser) parseDictionary() (*Dictionary, error) {
	dictionary := &Dictionary{Values: make(map[string]*Value)}

	if _, err := p.expect(OpenCurly); err != nil {
		return nil, err
	}

	for {
		// TODO Can we combine this in one operation? Also used in parseArray.
		token, err := p.peek()
		if err != nil {
			return nil, err
		}
		if token.Type == CloseCurly {
			p.token()
			return dictionary, nil
		}

		key, err := p.parseString()
		if err != nil {
			return nil, err
		}

		if _, err := p.expect(Equals); err != nil {
			return nil, err
		}

		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		if _, err := p.expect(Semicolon); err != nil {
			return nil, err
		}

		dictionary.Add(key.String, value)
	}

	return dictionary, nil
}

func (p *Parser) parseString() (*String, error) {
	token, err := p.token()
	if err != nil {
		return nil, err
	}
	if token.Type == Identifier {
		return &String{String: token.Literal}, nil
	}
	return nil, fmt.Errorf("Expected Identifier or String but got %d", token.Type)
}

func (p *Parser) parseArray() (*Array, error) {
	array := &Array{}

	if _, err := p.expect(OpenParen); err != nil {
		return nil, err
	}

	for {
		token, err := p.peek()
		if err != nil {
			return nil, err
		}
		if token.Type == CloseParen {
			p.token()
			return array, nil
		}

		value, err := p.parseValue()
		if err != nil {
			return nil, fmt.Errorf("Expected value: %s", err)
		}
		if _, err := p.expect(Comma); err != nil {
			return nil, err
		}
		array.Append(value)
	}

	return nil, errors.New("Expected CloseParen or Value")
}

func (p *Parser) parseDictionaryOrArray() (Value, error) {
	token, err := p.peek()
	if err != nil {
		return nil, err
	}

	if token.Type == OpenCurly {
		return p.parseDictionary()
	}

	if token.Type == OpenParen {
		return p.parseArray()
	}

	return nil, fmt.Errorf("Expected OpenCurly or OpenParen but got: %s", token)
}

func (p *Parser) parseValue() (Value, error) {
	token, err := p.peek()
	if err != nil {
		return nil, err
	}

	if token.Type == OpenCurly {
		return p.parseDictionary()
	}

	if token.Type == OpenParen {
		return p.parseArray()
	}

	if token.Type == Identifier {
		return p.parseString()
	}

	return nil, fmt.Errorf("Expected OpenCurly, OpenParen or Identifier but got %v at %d", token, token.Line)
}

func (p *Parser) Parse() (Value, error) {
	if err := p.parseHeader(); err != nil {
		return nil, err
	}
	return p.parseDictionaryOrArray()
}
