// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gg

import (
	"image"
	"image/color"

	"github.com/golang/freetype/raster"
)

// RepeatOp specifies how an image is repeated or tiled when it is painted onto a larger canvas.
type RepeatOp int

const (
	RepeatBoth RepeatOp = iota // Repeat both indicates that the image is repeated both horizontally (X-axis) and vertically (Y-axis).
	RepeatX                    // Repeat x indicates that the image is repeated horizontally (X-axis) but not vertically (Y-axis).
	RepeatY                    // Repeat y indicates that the image is repeated vertically (Y-axis) but not horizontally (X-axis).
	RepeatNone                 // Repeat none indicates that the image is not repeated and is displayed only once.
)

// Pattern is an interface representing a pattern used for filling areas with a repeating design.
// Implementations of this interface must provide the ability to sample and retrieve color information
// at specific coordinates (x, y) within the pattern.
type Pattern interface {
	ColorAt(x, y int) color.Color
}

// solidPattern is an implementation of the Pattern interface representing a solid color pattern.
// It provides a consistent color regardless of the sampling coordinates (x, y).
type solidPattern struct {
	color color.Color
}

// ColorAt returns the solid color assigned to the solidPattern, regardless of the sampling coordinates (x, y).
func (p *solidPattern) ColorAt(_, _ int) color.Color {
	return p.color
}

// NewSolidPattern creates and returns a solid color pattern using the specified color.
func NewSolidPattern(color color.Color) Pattern {
	return &solidPattern{color: color}
}

// surfacePattern is a pattern based on an image that can be repeated in different ways.
type surfacePattern struct {
	im image.Image
	op RepeatOp
}

// ColorAt returns the color of the pattern at the specified coordinates (x, y).
// The behavior of color retrieval depends on the repetition behavior defined by the RepeatOp.
func (p *surfacePattern) ColorAt(x, y int) color.Color {
	b := p.im.Bounds()

	switch p.op {
	case RepeatX:
		if y >= b.Dy() {
			return color.Transparent
		}
	case RepeatY:
		if x >= b.Dx() {
			return color.Transparent
		}
	case RepeatNone:
		if x >= b.Dx() || y >= b.Dy() {
			return color.Transparent
		}
	}

	x = x%b.Dx() + b.Min.X
	y = y%b.Dy() + b.Min.Y

	return p.im.At(x, y)
}

// NewSurfacePattern creates a new surface pattern from the given image and repetition behavior.
func NewSurfacePattern(im image.Image, op RepeatOp) Pattern {
	return &surfacePattern{im: im, op: op}
}

// patternPainter is responsible for painting patterns onto an image using a mask.
type patternPainter struct {
	im   *image.RGBA
	mask *image.Alpha
	p    Pattern
}

// Paint paints a sequence of spans onto the target image with the specified mask.
// The spans define the area to be painted, and the mask restricts the painting to a specific region.
func (r *patternPainter) Paint(ss []raster.Span, _ bool) {
	b := r.im.Bounds()

	for _, s := range ss {
		if s.Y < b.Min.Y {
			continue
		}
		if s.Y >= b.Max.Y {
			return
		}
		if s.X0 < b.Min.X {
			s.X0 = b.Min.X
		}
		if s.X1 > b.Max.X {
			s.X1 = b.Max.X
		}
		if s.X0 >= s.X1 {
			continue
		}

		const m = 1<<16 - 1
		y := s.Y - r.im.Rect.Min.Y
		x0 := s.X0 - r.im.Rect.Min.X
		// RGBAPainter.Paint() in $GOPATH/src/github.com/golang/freetype/raster/paint.go
		i0 := (s.Y-r.im.Rect.Min.Y)*r.im.Stride + (s.X0-r.im.Rect.Min.X)*4
		i1 := i0 + (s.X1-s.X0)*4

		for i, x := i0, x0; i < i1; i, x = i+4, x+1 {
			ma := s.Alpha

			if r.mask != nil {
				ma = ma * uint32(r.mask.AlphaAt(x, y).A) / 255
				if ma == 0 {
					continue
				}
			}

			c := r.p.ColorAt(x, y)
			cr, cg, cb, ca := c.RGBA()
			dr := uint32(r.im.Pix[i+0])
			dg := uint32(r.im.Pix[i+1])
			db := uint32(r.im.Pix[i+2])
			da := uint32(r.im.Pix[i+3])
			a := (m - (ca * ma / m)) * 0x101

			r.im.Pix[i+0] = uint8((dr*a + cr*ma) / m >> 8)
			r.im.Pix[i+1] = uint8((dg*a + cg*ma) / m >> 8)
			r.im.Pix[i+2] = uint8((db*a + cb*ma) / m >> 8)
			r.im.Pix[i+3] = uint8((da*a + ca*ma) / m >> 8)
		}
	}
}

// newPatternPainter creates a new patternPainter, which is a painter that applies a given pattern to an RGBA image
// while respecting an optional mask. It is used to paint patterns onto an image.
func newPatternPainter(im *image.RGBA, mask *image.Alpha, p Pattern) *patternPainter {
	return &patternPainter{im, mask, p}
}
