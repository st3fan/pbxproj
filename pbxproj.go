package main

import (
	"fmt"
	"io"
	"os"

	"github.com/st3fan/pbxproj/parser"
)

type Project struct {
	Root parser.Value
}

func Open(path string) (Project, error) {
	src, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	return Read(src)
}

func Read(r io.Reader) (Project, error) {
	parser, err := parser.NewParser(r)
	if err != nil {
		return Project{}, nil
	}
	root, err := parser.Parse()
	if err != nil {
		return Project{}, nil
	}
	return Project{Root: root}, nil
}

func (p Project) Encode(w io.Writer) {
	p.Root.Encode(w, 0)
	fmt.Fprintln(w)
}

func main() {
	project, err := Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	project.Encode(os.Stdout)
}
