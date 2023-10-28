// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gg

import "math"

// quadratic calculates points on a quadratic Bézier curve at a given parameter 't'.
//
// The quadratic Bézier curve is defined by three control points: (x0, y0), (x1, y1), and (x2, y2).
// The parameter 't' varies between 0 and 1, where t=0 corresponds to the start point (x0, y0) and
// t=1 corresponds to the end point (x2, y2) of the Bézier curve.
func quadratic(x0, y0, x1, y1, x2, y2, t float64) (x, y float64) {
	u := 1 - t
	a := u * u
	b := 2 * u * t
	c := t * t
	x = a*x0 + b*x1 + c*x2
	y = a*y0 + b*y1 + c*y2

	return x, y
}

// QuadraticBezier computes a series of points on a quadratic Bézier curve defined by three control points.
//
// The quadratic Bézier curve is determined by the starting point (x0, y0), a control point (x1, y1),
// and the ending point (x2, y2). The function calculates 'n' equidistant points along the curve, where 'n'
// is determined by the curve's total length. If the calculated number of points is less than 4, it defaults
// to 4 points.
func QuadraticBezier(x0, y0, x1, y1, x2, y2 float64) []Point {
	l := math.Hypot(x1-x0, y1-y0) + math.Hypot(x2-x1, y2-y1)
	n := int(l + 0.5)

	if n < 4 {
		n = 4
	}

	d := float64(n) - 1
	result := make([]Point, n)

	for i := 0; i < n; i++ {
		t := float64(i) / d
		x, y := quadratic(x0, y0, x1, y1, x2, y2, t)
		result[i] = Point{x, y}
	}

	return result
}

// cubic calculates points on a cubic Bézier curve at a given parameter 't'.
//
// The cubic Bézier curve is defined by four control points: (x0, y0), (x1, y1), (x2, y2), and (x3, y3).
// The parameter 't' varies between 0 and 1, where t=0 corresponds to the start point (x0, y0) and t=1
// corresponds to the end point (x3, y3) of the Bézier curve.
func cubic(x0, y0, x1, y1, x2, y2, x3, y3, t float64) (x, y float64) {
	u := 1 - t
	a := u * u * u
	b := 3 * u * u * t
	c := 3 * u * t * t
	d := t * t * t
	x = a*x0 + b*x1 + c*x2 + d*x3
	y = a*y0 + b*y1 + c*y2 + d*y3

	return x, y
}

// CubicBezier computes a series of points on a cubic Bézier curve defined by four control points.
//
// The cubic Bézier curve is determined by the starting point (x0, y0), two control points (x1, y1) and (x2, y2),
// and the ending point (x3, y3). The function calculates 'n' equidistant points along the curve, where 'n' is
// determined by the cumulative length of the curve. If the calculated number of points is less than 4, it defaults
// to 4 points.
func CubicBezier(x0, y0, x1, y1, x2, y2, x3, y3 float64) []Point {
	l := math.Hypot(x1-x0, y1-y0) + math.Hypot(x2-x1, y2-y1) + math.Hypot(x3-x2, y3-y2)
	n := int(l + 0.5)

	if n < 4 {
		n = 4
	}

	d := float64(n) - 1
	result := make([]Point, n)

	for i := 0; i < n; i++ {
		t := float64(i) / d
		x, y := cubic(x0, y0, x1, y1, x2, y2, x3, y3, t)
		result[i] = Point{x, y}
	}

	return result
}
