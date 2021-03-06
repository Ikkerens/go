// errorcheck -0 -m -l

// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test escape analysis with respect to field assignments.

package escape

var sink interface{}

type X struct {
	p1 *int
	p2 *int
	a  [2]*int
}

type Y struct {
	x X
}

func field0() {
	i := 0 // ERROR "moved to heap: i$"
	var x X
	x.p1 = &i
	sink = x.p1 // ERROR "x\.p1 escapes to heap"
}

func field1() {
	i := 0 // ERROR "moved to heap: i$"
	var x X
	// BAD: &i should not escape
	x.p1 = &i
	sink = x.p2 // ERROR "x\.p2 escapes to heap"
}

func field3() {
	i := 0 // ERROR "moved to heap: i$"
	var x X
	x.p1 = &i
	sink = x  // ERROR "x escapes to heap"
}

func field4() {
	i := 0 // ERROR "moved to heap: i$"
	var y Y
	y.x.p1 = &i
	x := y.x
	sink = x // ERROR "x escapes to heap"
}

func field5() {
	i := 0 // ERROR "moved to heap: i$"
	var x X
	// BAD: &i should not escape here
	x.a[0] = &i
	sink = x.a[1] // ERROR "x\.a\[1\] escapes to heap"
}

// BAD: we are not leaking param x, only x.p2
func field6(x *X) { // ERROR "leaking param content: x$"
	sink = x.p2 // ERROR "x\.p2 escapes to heap"
}

func field6a() {
	i := 0 // ERROR "moved to heap: i$"
	var x X
	// BAD: &i should not escape
	x.p1 = &i
	field6(&x)
}

func field7() {
	i := 0
	var y Y
	y.x.p1 = &i
	x := y.x
	var y1 Y
	y1.x = x
	_ = y1.x.p1
}

func field8() {
	i := 0 // ERROR "moved to heap: i$"
	var y Y
	y.x.p1 = &i
	x := y.x
	var y1 Y
	y1.x = x
	sink = y1.x.p1 // ERROR "y1\.x\.p1 escapes to heap"
}

func field9() {
	i := 0 // ERROR "moved to heap: i$"
	var y Y
	y.x.p1 = &i
	x := y.x
	var y1 Y
	y1.x = x
	sink = y1.x // ERROR "y1\.x escapes to heap"
}

func field10() {
	i := 0 // ERROR "moved to heap: i$"
	var y Y
	// BAD: &i should not escape
	y.x.p1 = &i
	x := y.x
	var y1 Y
	y1.x = x
	sink = y1.x.p2 // ERROR "y1\.x\.p2 escapes to heap"
}

func field11() {
	i := 0         // ERROR "moved to heap: i$"
	x := X{p1: &i}
	sink = x.p1    // ERROR "x\.p1 escapes to heap"
}

func field12() {
	i := 0 // ERROR "moved to heap: i$"
	// BAD: &i should not escape
	x := X{p1: &i}
	sink = x.p2    // ERROR "x\.p2 escapes to heap"
}

func field13() {
	i := 0          // ERROR "moved to heap: i$"
	x := &X{p1: &i} // ERROR "field13 &X literal does not escape$"
	sink = x.p1     // ERROR "x\.p1 escapes to heap"
}

func field14() {
	i := 0 // ERROR "moved to heap: i$"
	// BAD: &i should not escape
	x := &X{p1: &i} // ERROR "field14 &X literal does not escape$"
	sink = x.p2     // ERROR "x\.p2 escapes to heap"
}

func field15() {
	i := 0          // ERROR "moved to heap: i$"
	x := &X{p1: &i} // ERROR "&X literal escapes to heap$"
	sink = x        // ERROR "x escapes to heap"
}

func field16() {
	i := 0 // ERROR "moved to heap: i$"
	var x X
	// BAD: &i should not escape
	x.p1 = &i
	var iface interface{} = x // ERROR "x escapes to heap"
	x1 := iface.(X)
	sink = x1.p2 // ERROR "x1\.p2 escapes to heap"
}

func field17() {
	i := 0 // ERROR "moved to heap: i$"
	var x X
	x.p1 = &i
	var iface interface{} = x // ERROR "x escapes to heap"
	x1 := iface.(X)
	sink = x1.p1 // ERROR "x1\.p1 escapes to heap"
}

func field18() {
	i := 0 // ERROR "moved to heap: i$"
	var x X
	// BAD: &i should not escape
	x.p1 = &i
	var iface interface{} = x // ERROR "x escapes to heap"
	y, _ := iface.(Y)         // Put X, but extracted Y. The cast will fail, so y is zero initialized.
	sink = y                  // ERROR "y escapes to heap"
}
