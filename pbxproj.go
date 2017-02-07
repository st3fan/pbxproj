package main

import (
	"bytes"
	"fmt"

	"github.com/st3fan/pbxproj/parser"
)

func main() {
	// src, err := os.Open("project.pbxproj")
	// if err != nil {
	// 	panic(err)
	// }

	src := bytes.NewBufferString(`foo = "bar";`)

	tokenizer, err := parser.NewTokenizer(src)
	if err != nil {
		panic(err)
	}

	var token parser.Token
	for token.Type != parser.EndOfFile {
		token, err := tokenizer.Next()
		if err != nil {
			panic(err)
		}

		if token.Type == parser.EndOfFile {
			break
		}

		fmt.Printf("%+v\n", token)
	}
}
