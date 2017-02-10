package parser

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func compareTokens(t *testing.T, a []Token, b []Token) {
	if len(a) != len(b) {
		t.Fatalf("Did not get expected tokens len(a) = %d len(b) = %d", len(a), len(b))
	}
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("Did not get expected tokens: a = %v b = %v", a, b)
		}
	}
}

func tokenize(t *testing.T, r io.Reader) []Token {
	tokenizer, err := NewTokenizer(r)
	if err != nil {
		t.Fatal(err)
	}

	var tokens []Token
	for {
		token, err := tokenizer.Next()
		if err != nil {
			t.Fatal(err)
		}

		if token.Type == EndOfFile {
			break
		}

		tokens = append(tokens, token)
	}

	return tokens
}

func tokenizeString(t *testing.T, src string) []Token {
	return tokenize(t, bytes.NewBufferString(src))
}

func tokenizeFile(t *testing.T, path string) []Token {
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	return tokenize(t, file)
}

// func Test_Peek(t *testing.T) {
// 	tokenizer, err := NewTokenizer(bytes.NewBufferString("abc"))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	c1, c2 := tokenizer.peek()
// 	if c1 != 'a' {
// 		t.Fatal("c1 != 'a'")
// 	}
// 	if c2 != 'a' {
// 		t.Fatal("c1 != 'a'")
// 	}
// }

func Test_TokenizeBasic1(t *testing.T) {
	compareTokens(t, tokenizeString(t, "things = { thing = Foo.swift; }"),
		[]Token{
			Token{Type: Identifier, Literal: "things", Line: 1},
			Token{Type: Whitespace, Literal: " ", Line: 1},
			Token{Type: Equals, Line: 1},
			Token{Type: Whitespace, Literal: " ", Line: 1},
			Token{Type: OpenCurly, Line: 1},
			Token{Type: Whitespace, Literal: " ", Line: 1},
			Token{Type: Identifier, Literal: "thing", Line: 1},
			Token{Type: Whitespace, Literal: " ", Line: 1},
			Token{Type: Equals, Line: 1},
			Token{Type: Whitespace, Literal: " ", Line: 1},
			Token{Type: Identifier, Literal: "Foo.swift", Line: 1},
			Token{Type: Semicolon, Line: 1},
			Token{Type: Whitespace, Literal: " ", Line: 1},
			Token{Type: CloseCurly, Line: 1},
			Token{Type: EndOfLine, Line: 1},
		},
	)
}

func Test_TokenizeComment1(t *testing.T) {
	compareTokens(t, tokenizeString(t, "/* Hello */"),
		[]Token{
			Token{Type: Comment, Literal: "/* Hello */", Line: 1},
			Token{Type: EndOfLine, Line: 1},
		},
	)
}

func Test_TokenizeComment2(t *testing.T) {
	compareTokens(t, tokenizeString(t, "{ /* Hello */ }"),
		[]Token{
			Token{Type: OpenCurly, Line: 1},
			Token{Type: Whitespace, Literal: " ", Line: 1},
			Token{Type: Comment, Literal: "/* Hello */", Line: 1},
			Token{Type: Whitespace, Literal: " ", Line: 1},
			Token{Type: CloseCurly, Line: 1},
			Token{Type: EndOfLine, Line: 1},
		},
	)
}

func Test_TokenizeComment3(t *testing.T) {
	compareTokens(t, tokenizeString(t, "/* One */{/* Two */}/* Three */"),
		[]Token{
			Token{Type: Comment, Literal: "/* One */", Line: 1},
			Token{Type: OpenCurly, Line: 1},
			Token{Type: Comment, Literal: "/* Two */", Line: 1},
			Token{Type: CloseCurly, Line: 1},
			Token{Type: Comment, Literal: "/* Three */", Line: 1},
			Token{Type: EndOfLine, Line: 1},
		},
	)
}

func Test_TokenizeComment4(t *testing.T) {
	compareTokens(t, tokenizeString(t, "/*One*/{/*Two*/}/*Three*/"),
		[]Token{
			Token{Type: Comment, Literal: "/*One*/", Line: 1},
			Token{Type: OpenCurly, Line: 1},
			Token{Type: Comment, Literal: "/*Two*/", Line: 1},
			Token{Type: CloseCurly, Line: 1},
			Token{Type: Comment, Literal: "/*Three*/", Line: 1},
			Token{Type: EndOfLine, Line: 1},
		},
	)
}

func Test_TokenizeMinimal(t *testing.T) {
	tokenizeFile(t, "testdata/minimal.pbxproj")
}

func Test_TokenizeFirefox(t *testing.T) {
	tokenizeFile(t, "testdata/firefox.pbxproj")
}

func Test_TokenizeBlockzilla(t *testing.T) {
	tokenizeFile(t, "testdata/blockzilla.pbxproj")
}
