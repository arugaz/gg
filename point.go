// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gg

import (
	"math"

	"golang.org/x/image/math/fixed"
)

// Point represents a two-dimensional point with X and Y coordinates.
type Point struct {
	X, Y float64
}

// Fixed converts a Point to a fixed.Point26_6 representation.
func (a Point) Fixed() fixed.Point26_6 {
	return fixp(a.X, a.Y)
}

// Distance calculates the Euclidean distance between two points 'a' and 'b'.
func (a Point) Distance(b Point) float64 {
	return math.Hypot(a.X-b.X, a.Y-b.Y)
}

// Interpolate performs linear interpolation between two points 'a' and 'b' at a given parameter 't'.
func (a Point) Interpolate(b Point, t float64) Point {
	x := a.X + (b.X-a.X)*t
	y := a.Y + (b.Y-a.Y)*t

	return Point{X: x, Y: y}
}
