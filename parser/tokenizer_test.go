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
			Token{Type: Identifier, Literal: "things"},
			Token{Type: Whitespace, Literal: " "},
			Token{Type: Equals},
			Token{Type: Whitespace, Literal: " "},
			Token{Type: OpenCurly},
			Token{Type: Whitespace, Literal: " "},
			Token{Type: Identifier, Literal: "thing"},
			Token{Type: Whitespace, Literal: " "},
			Token{Type: Equals},
			Token{Type: Whitespace, Literal: " "},
			Token{Type: Identifier, Literal: "Foo.swift"},
			Token{Type: Semicolon},
			Token{Type: Whitespace, Literal: " "},
			Token{Type: CloseCurly},
			Token{Type: EndOfLine},
		},
	)
}

func Test_TokenizeComment1(t *testing.T) {
	compareTokens(t, tokenizeString(t, "/* Hello */"),
		[]Token{
			Token{Type: Comment, Literal: "/* Hello */"},
			Token{Type: EndOfLine},
		},
	)
}

func Test_TokenizeComment2(t *testing.T) {
	compareTokens(t, tokenizeString(t, "{ /* Hello */ }"),
		[]Token{
			Token{Type: OpenCurly},
			Token{Type: Whitespace, Literal: " "},
			Token{Type: Comment, Literal: "/* Hello */"},
			Token{Type: Whitespace, Literal: " "},
			Token{Type: CloseCurly},
			Token{Type: EndOfLine},
		},
	)
}

func Test_TokenizeComment3(t *testing.T) {
	compareTokens(t, tokenizeString(t, "/* One */{/* Two */}/* Three */"),
		[]Token{
			Token{Type: Comment, Literal: "/* One */"},
			Token{Type: OpenCurly},
			Token{Type: Comment, Literal: "/* Two */"},
			Token{Type: CloseCurly},
			Token{Type: Comment, Literal: "/* Three */"},
			Token{Type: EndOfLine},
		},
	)
}

func Test_TokenizeComment4(t *testing.T) {
	compareTokens(t, tokenizeString(t, "/*One*/{/*Two*/}/*Three*/"),
		[]Token{
			Token{Type: Comment, Literal: "/*One*/"},
			Token{Type: OpenCurly},
			Token{Type: Comment, Literal: "/*Two*/"},
			Token{Type: CloseCurly},
			Token{Type: Comment, Literal: "/*Three*/"},
			Token{Type: EndOfLine},
		},
	)
}

func Test_TokenizeFirefox(t *testing.T) {
	tokenizeFile(t, "testdata/firefox.pbxproj")
}

func Test_TokenizeBlockzilla(t *testing.T) {
	tokenizeFile(t, "testdata/blockzilla.pbxproj")
}
