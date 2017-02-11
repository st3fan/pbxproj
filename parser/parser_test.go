// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package parser

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func parse(t *testing.T, src io.Reader) Value {
	parser, err := NewParser(src)
	if err != nil {
		t.Fatal(err)
	}
	value, err := parser.Parse()
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func parseString(t *testing.T, src string) Value {
	return parse(t, bytes.NewBufferString(src))
}

func parseFile(t *testing.T, path string) Value {
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	return parse(t, file)
}

func failToParseString(t *testing.T, src string) {
	parser, err := NewParser(bytes.NewBufferString(src))
	if err != nil {
		t.Fatal(err)
	}
	_, err = parser.Parse()
	if err == nil {
		t.Fatal("Did not fail to parse: " + src)
	}
}

func Test_ParseEmptyArray(t *testing.T) {
	value := parseString(t, "// !$*UTF8*$!\n()\n")
	array, ok := value.(*Array)
	if !ok {
		t.Fatal("Expected an Array but got:", value)
	}
	if array.Count() != 0 {
		t.Fatal("Expected an Array with Count 0 but got: ", array.Count())
	}
}

func Test_ParseStringArray(t *testing.T) {
	value := parseString(t, "// !$*UTF8*$!\n(\"Foo\", \"Bar\", \"Baz\", )\n")
	array, ok := value.(*Array)
	if !ok {
		t.Fatal("Expected an Array but got:", value)
	}
	if array.Count() != 3 {
		t.Fatal("Expected an Array with Count 3 but got: ", array.Count())
	}
}

func Test_ParseStringArrayWithoutCommas(t *testing.T) {
	failToParseString(t, "// !$*UTF8*$!\n(\"Foo\" \"Bar\" \"Baz\")\n")
}

// Dictionary

func Test_ParseEmptyDictionary(t *testing.T) {
	value := parseString(t, "// !$*UTF8*$!\n{})\n")
	dictionary, ok := value.(*Dictionary)
	if !ok {
		t.Fatal("Expected a Dictionary but got:", value)
	}
	if dictionary.Count() != 0 {
		t.Fatal("Expected a Dictionary with Count 0 but got: ", dictionary.Count())
	}
}

func Test_ParseDictionary(t *testing.T) {
	value := parseString(t, "// !$*UTF8*$!\n{foo=1; bar=2;})\n")
	dictionary, ok := value.(*Dictionary)
	if !ok {
		t.Fatal("Expected a Dictionary but got:", value)
	}
	if dictionary.Count() != 2 {
		t.Fatal("Expected a Dictionary with Count 2 but got: ", dictionary.Count())
	}
}

// Test files

func Test_ParseMinimal(t *testing.T) {
	parseFile(t, "testdata/minimal.pbxproj")
}

func Test_ParseFirefox(t *testing.T) {
	parseFile(t, "testdata/firefox.pbxproj")
}

func Test_ParseBlockzilla(t *testing.T) {
	parseFile(t, "testdata/blockzilla.pbxproj")
}
