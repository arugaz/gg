// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gg

import (
	"image/color"
	"math"
	"sort"
)

// stop represents a color stop in a gradient. It defines the position (pos)
// within the gradient and the color to use at that position.
type stop struct {
	pos   float64
	color color.Color
}

// stops is a slice of color stops used in gradients.
// It is a collection of stop points that define the colors and positions in a gradient.
type stops []stop

// Len returns the number of stops in the stops slice.
func (s stops) Len() int {
	return len(s)
}

// Less reports whether the stop at index 'i' should be sorted before the stop at index 'j'.
// Stops are sorted based on their position values in ascending order.
func (s stops) Less(i, j int) bool {
	return s[i].pos < s[j].pos
}

// Swap swaps the positions of two stops in the slice. It is used to re-order stops
// during sorting operations.
func (s stops) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Gradient is an interface representing a gradient pattern that can be used to
// fill shapes with smooth color transitions.
type Gradient interface {
	Pattern

	// AddColorStop appends a color stop to the linear gradient at the specified offset.
	//
	// This method adds a color stop to the gradient with the specified offset and color. Color stops
	// define where colors change within the gradient. The stops are sorted in ascending order based
	// on their offset values.
	AddColorStop(offset float64, color color.Color)
}

// linearGradient represents a linear gradient that can be used as a pattern
// to fill shapes with colors transitioning between specified points.
type linearGradient struct {
	x0, y0, x1, y1 float64
	stops          stops
}

// ColorAt returns the color at the specified (x, y) coordinate within the linear gradient.
//
// This method calculates the color at the given point using the gradient's color stops and
// linear interpolation. If there are no color stops defined in the gradient, it returns
// a transparent color.
func (g *linearGradient) ColorAt(x, y int) color.Color {
	if len(g.stops) == 0 {
		return color.Transparent
	}

	fx, fy := float64(x), float64(y)
	x0, y0, x1, y1 := g.x0, g.y0, g.x1, g.y1
	dx, dy := x1-x0, y1-y0

	// Horizontal
	if dy == 0 && dx != 0 {
		return getColor((fx-x0)/dx, g.stops)
	}

	// Vertical
	if dx == 0 && dy != 0 {
		return getColor((fy-y0)/dy, g.stops)
	}

	// Dot product
	s0 := dx*(fx-x0) + dy*(fy-y0)
	if s0 < 0 {
		return g.stops[0].color
	}

	// Calculate distance to (x0,y0) alone (x0,y0)->(x1,y1)
	mag := math.Hypot(dx, dy)
	u := ((fx-x0)*-dy + (fy-y0)*dx) / (mag * mag)
	x2, y2 := x0+u*-dy, y0+u*dx
	d := math.Hypot(fx-x2, fy-y2) / mag

	return getColor(d, g.stops)
}

// AddColorStop appends a color stop to the linear gradient at the specified offset.
//
// This method adds a color stop to the gradient with the specified offset and color. Color stops
// define where colors change within the gradient. The stops are sorted in ascending order based
// on their offset values.
func (g *linearGradient) AddColorStop(offset float64, color color.Color) {
	g.stops = append(g.stops, stop{pos: offset, color: color})
	sort.Sort(g.stops)
}

// NewLinearGradient creates a new linear gradient pattern that spans from point (x0, y0) to point (x1, y1).
//
// This function creates and returns a linear gradient pattern with the specified start and end points.
// Linear gradients fill an area with colors that transition from one point to another in a straight line.
func NewLinearGradient(x0, y0, x1, y1 float64) Gradient {
	return &linearGradient{
		x0: x0, y0: y0,
		x1: x1, y1: y1,
	}
}

// circle represents a geometric circle with a center point (x, y) and a radius (r).
type circle struct {
	x, y, r float64
}

// radialGradient represents a radial gradient pattern defined by three circles, colors, and stops.
type radialGradient struct {
	c0, c1, cd circle
	a, inva    float64
	mindr      float64
	stops      stops
}

// dot3 calculates the dot product of two 3D vectors (x0, y0, z0) and (x1, y1, z1).
//
// The dot product is a mathematical operation that returns a scalar value by taking
// the sum of the pairwise products of the corresponding elements of the two vectors.
func dot3(x0, y0, z0, x1, y1, z1 float64) float64 {
	return x0*x1 + y0*y1 + z0*z1
}

// ColorAt calculates the color of the radial gradient at the specified point (x, y).
//
// This method implements the color calculation for a radial gradient based on the gradient's
// parameters and color stops. The color is determined by the position of the point within the
// gradient's radial shape and the interpolation of colors defined by the gradient stops.
func (g *radialGradient) ColorAt(x, y int) color.Color {
	if len(g.stops) == 0 {
		return color.Transparent
	}

	// copy from pixman's pixman-radial-gradient.c
	dx, dy := float64(x)+0.5-g.c0.x, float64(y)+0.5-g.c0.y
	b := dot3(dx, dy, g.c0.r, g.cd.x, g.cd.y, g.cd.r)
	c := dot3(dx, dy, -g.c0.r, dx, dy, g.c0.r)

	if g.a == 0 {
		if b == 0 {
			return color.Transparent
		}
		t := 0.5 * c / b
		if t*g.cd.r >= g.mindr {
			return getColor(t, g.stops)
		}
		return color.Transparent
	}

	discr := dot3(b, g.a, 0, b, -c, 0)
	if discr >= 0 {
		sqrtdiscr := math.Sqrt(discr)
		t0 := (b + sqrtdiscr) * g.inva
		t1 := (b - sqrtdiscr) * g.inva

		if t0*g.cd.r >= g.mindr {
			return getColor(t0, g.stops)
		} else if t1*g.cd.r >= g.mindr {
			return getColor(t1, g.stops)
		}
	}

	return color.Transparent
}

// AddColorStop adds a color stop to the gradient at the specified offset.
func (g *radialGradient) AddColorStop(offset float64, color color.Color) {
	g.stops = append(g.stops, stop{pos: offset, color: color})
	sort.Sort(g.stops)
}

// NewRadialGradient creates a new radial gradient object based on the specified parameters.
//
// This function initializes a radial gradient with the given circle coordinates and radii,
// which define the gradient's shape. The radial gradient will interpolate colors smoothly
// from one circle (defined by x0, y0, r0) to another circle (defined by x1, y1, r1).
func NewRadialGradient(x0, y0, r0, x1, y1, r1 float64) Gradient {
	c0 := circle{x0, y0, r0}
	c1 := circle{x1, y1, r1}
	cd := circle{x1 - x0, y1 - y0, r1 - r0}
	a := dot3(cd.x, cd.y, -cd.r, cd.x, cd.y, cd.r)

	var inva float64
	if a != 0 {
		inva = 1.0 / a
	}
	mindr := -c0.r

	return &radialGradient{
		c0:    c0,
		c1:    c1,
		cd:    cd,
		a:     a,
		inva:  inva,
		mindr: mindr,
	}
}

// conicGradient represents a conic (angular) gradient.
type conicGradient struct {
	cx, cy   float64
	rotation float64
	stops    stops
}

// ColorAt returns the color at the specified coordinates (x, y) within the conic gradient.
//
// This method calculates the color at the given (x, y) position based on the angular
// position around the center of the conic gradient. The gradient is defined by its center
// (cx, cy), rotation angle, and color stops. The resulting color represents the gradient
// color at the specified point.
func (g *conicGradient) ColorAt(x, y int) color.Color {
	if len(g.stops) == 0 {
		return color.Transparent
	}

	a := math.Atan2(float64(y)-g.cy, float64(x)-g.cx)
	t := norm(a, -math.Pi, math.Pi) - g.rotation

	if t < 0 {
		t += 1
	}

	return getColor(t, g.stops)
}

// AddColorStop adds a color stop to the conic gradient at the specified offset position.
//
// This method allows you to define color stops for the conic gradient. A color stop
// consists of an offset value and a color. The offset represents the position along
// the angular gradient path where the specified color should take effect.
func (g *conicGradient) AddColorStop(offset float64, color color.Color) {
	g.stops = append(g.stops, stop{pos: offset, color: color})
	sort.Sort(g.stops)
}

// NewConicGradient creates a new conic gradient with the specified center coordinates and rotation angle in degrees.
//
// This function creates a conic gradient, which is a type of gradient that varies in color and
// smoothly rotates around a central point. You can define the center coordinates (cx, cy)
// and the rotation angle in degrees to control the gradient's appearance.
func NewConicGradient(cx, cy, deg float64) Gradient {
	return &conicGradient{
		cx:       cx,
		cy:       cy,
		rotation: normalizeAngle(deg) / 360,
	}
}

// normalizeAngle normalizes an angle value to the range [0, 360) degrees.
//
// This function takes an angle in degrees as input and ensures that it falls within
// the range [0, 360) degrees. If the input angle is negative, it wraps around to
// the positive range.
func normalizeAngle(t float64) float64 {
	t = math.Mod(t, 360)

	if t < 0 {
		t += 360
	}

	return t
}

// norm normalizes a value within a specified range [a, b] to the range [0, 1].
//
// This function takes a value, `value`, and normalizes it to a range of [0, 1]
// based on the given minimum (`a`) and maximum (`b`) values. It's commonly used
// for interpolating or mapping values within a specific range.
func norm(value, a, b float64) float64 {
	return (value - a) * (1.0 / (b - a))
}

// getColor returns the interpolated color at the specified position along a gradient.
//
// This function calculates and returns the interpolated color at a given position
// within a gradient defined by color stops. It linearly interpolates between two
// adjacent stops based on the specified position.
func getColor(pos float64, stops stops) color.Color {
	if pos <= 0.0 || len(stops) == 1 {
		return stops[0].color
	}

	last := stops[len(stops)-1]

	if pos >= last.pos {
		return last.color
	}

	for i, stop := range stops[1:] {
		if pos < stop.pos {
			pos = (pos - stops[i].pos) / (stop.pos - stops[i].pos)
			return colorLerp(stops[i].color, stop.color, pos)
		}
	}

	return last.color
}

// colorLerp linearly interpolates between two colors based on a given position.
//
// This function takes two color.Color values, c0 and c1, and linearly interpolates
// between them based on the position 't' where t is a value between 0 and 1. It returns
// the interpolated color at the specified position.
func colorLerp(c0, c1 color.Color, t float64) color.Color {
	r0, g0, b0, a0 := c0.RGBA()
	r1, g1, b1, a1 := c1.RGBA()

	return color.RGBA{
		R: lerp(r0, r1, t),
		G: lerp(g0, g1, t),
		B: lerp(b0, b1, t),
		A: lerp(a0, a1, t),
	}
}

// lerp performs linear interpolation between two values 'a' and 'b' based on the position 't'.
//
// This helper function is used by colorLerp to perform linear interpolation between two values.
// It takes two uint32 values 'a' and 'b' and returns the interpolated value based on the position 't'.
func lerp(a, b uint32, t float64) uint8 {
	return uint8(int32(float64(a)*(1.0-t)+float64(b)*t) >> 8)
}
