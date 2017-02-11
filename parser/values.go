// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/

package parser

import (
	"fmt"
	"io"
	"strconv"
)

type Value interface {
	Encode(w io.Writer, indent int)
}

//

type String struct {
	literal string
}

func (s String) Encode(w io.Writer, indent int) {
	fmt.Fprintf(w, "%s", strconv.Quote(s.literal))
}

//

type Array struct {
	values []Value
}

func (a Array) Encode(w io.Writer, indent int) {
}

func (a *Array) Append(value Value) {
	a.values = append(a.values, value)
}

func (a Array) Count() int {
	return len(a.values)
}

//

type Dictionary struct {
	values map[string]Value
}

func (d Dictionary) Encode(w io.Writer, indent int) {
	// for i := 0; i < indent; i++ {
	// 	fmt.Fprintf(w, "\t")
	// }
	fmt.Fprintln(w, "{")
	for key, value := range d.values {
		fmt.Fprintf(w, "\t%s = ", key)
		value.Encode(w, indent+1)
		fmt.Fprintf(w, ";\n")
	}
	for i := 0; i < indent; i++ {
		fmt.Fprintf(w, "\t")
	}
	fmt.Fprintf(w, "}")
}

func (d *Dictionary) Add(key string, value Value) {
	d.values[key] = value
}

func (d *Dictionary) Count() int {
	return len(d.values)
}
