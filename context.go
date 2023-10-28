// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package gg provides a simple API for rendering 2D graphics in pure Go.
package gg

import (
	"errors"
	"image"
	"image/color"
	"io"
	"io/fs"
	"math"
	"strings"

	"github.com/golang/freetype/raster"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/f64"
)

// LineCap defines the possible line cap styles for drawing paths in a rendering context.
type LineCap int

const (
	LineCapRound  LineCap = iota // Line cap style with round ends.
	LineCapButt                  // Line cap style with flat ends.
	LineCapSquare                // Line cap style with square ends.
)

// LineJoin defines the possible line join styles for connecting path segments in a rendering context.
type LineJoin int

const (
	LineJoinRound LineJoin = iota // Line join style with round connections.
	LineJoinBevel                 // Line join style with beveled connections.
)

// FillRule specifies the fill rule for determining the interior of complex paths.
type FillRule int

const (
	FillRuleWinding FillRule = iota // Winding fill rule for complex paths.
	FillRuleEvenOdd                 // Even-Odd fill rule for complex paths.
)

// Align defines text alignment options for text rendering.
type Align int

const (
	AlignLeft   Align = iota // Left alignment for text rendering.
	AlignCenter              // Center alignment for text rendering.
	AlignRight               // Right alignment for text rendering.
)

// defaultFillStyle represents the default fill style used in the rendering context.
var defaultFillStyle = NewSolidPattern(color.White)

// defaultStrokeStyle represents the default stroke style used in the rendering context.
var defaultStrokeStyle = NewSolidPattern(color.Black)

// Context represents a 2D rendering context used for drawing operations.
type Context struct {
	width         int
	height        int
	rasterizer    *raster.Rasterizer
	im            *image.RGBA
	mask          *image.Alpha
	color         color.Color
	fillPattern   Pattern
	strokePattern Pattern
	strokePath    raster.Path
	fillPath      raster.Path
	start         Point
	current       Point
	hasCurrent    bool
	dashes        []float64
	dashOffset    float64
	lineWidth     float64
	lineCap       LineCap
	lineJoin      LineJoin
	fillRule      FillRule
	fontFace      font.Face
	fontHeight    float64
	matrix        Matrix
	stack         []*Context
	interp        draw.Interpolator
	frames        []image.Image
}

// NewContext creates a new rendering context with the specified width and height.
//
// This function initializes a new rendering context with the provided 'width' and 'height' dimensions.
// The context is used for performing 2D drawing operations, including path rendering, text rendering,
// and manipulation of the rendering state.
func NewContext(width, height int) *Context {
	return NewContextForRGBA(image.NewRGBA(image.Rect(0, 0, width, height)))
}

// NewContextForImage creates a new rendering context based on an existing image.Image.
//
// This function initializes a new rendering context using the dimensions and content of the provided 'im' image.
// The context is used for performing 2D drawing operations, including path rendering, text rendering,
// and manipulation of the rendering state.
func NewContextForImage(im image.Image) *Context {
	return NewContextForRGBA(ImageToRGBA(im))
}

// NewContextForRGBA prepares a context for rendering onto the specified image.
// No copy is made.
func NewContextForRGBA(im *image.RGBA) *Context {
	w := im.Bounds().Size().X
	h := im.Bounds().Size().Y

	return &Context{
		width:         w,
		height:        h,
		rasterizer:    raster.NewRasterizer(w, h),
		im:            im,
		color:         color.Transparent,
		fillPattern:   defaultFillStyle,
		strokePattern: defaultStrokeStyle,
		lineWidth:     1,
		fillRule:      FillRuleWinding,
		fontFace:      basicfont.Face7x13,
		fontHeight:    13,
		matrix:        Identity(),
		interp:        draw.BiLinear,
	}
}

// GetCurrentPoint returns the current point (cursor position) in the rendering context.
//
// This method retrieves the current cursor position in the rendering context. The current point represents
// the location where the most recent drawing operation ended or where the cursor was explicitly positioned.
func (dc *Context) GetCurrentPoint() (Point, bool) {
	if dc.hasCurrent {
		return dc.current, true
	}

	return Point{}, false
}

// Image returns the image.RGBA associated with the rendering context.
//
// This method retrieves the image.RGBA associated with the rendering context. The image.RGBA represents
// the canvas or target image where the 2D drawing operations are performed and rendered.
func (dc *Context) Image() image.Image {
	return dc.im
}

// Frames returns the frames slice.
//
// This method returns the frames slice. It is used for creating animated image.
func (dc *Context) Frames() []image.Image {
	return dc.frames
}

// Width returns the width of the rendering context.
//
// This method retrieves the width of the rendering context in pixels.
func (dc *Context) Width() int {
	return dc.width
}

// Height returns the height of the rendering context.
//
// This method retrieves the height of the rendering context in pixels.
func (dc *Context) Height() int {
	return dc.height
}

// SetDash sets the dash pattern for stroking lines.
//
// This method allows you to set a custom dash pattern for stroking lines. You can provide a sequence of 'dashes' values
// to specify the lengths of dashes and gaps between them. For example, to create a dashed line with a 10-pixel dash
// followed by a 5-pixel gap, you can call SetDash(10, 5).
func (dc *Context) SetDash(dashes ...float64) {
	dc.dashes = dashes
}

// SetDashOffset sets the offset for the dash pattern.
//
// This method allows you to set the offset (phase) for the dash pattern when stroking lines. The 'offset' value determines
// where the dash pattern starts along the path. You can specify the 'offset' in units of the current line width.
func (dc *Context) SetDashOffset(offset float64) {
	dc.dashOffset = offset
}

// SetLineWidth sets the line width for drawing operations.
//
// This method allows you to set the line width for stroking lines and drawing shapes. The 'lineWidth' parameter specifies
// the width of lines in user space units (typically pixels).
func (dc *Context) SetLineWidth(lineWidth float64) {
	dc.lineWidth = lineWidth
}

// SetLineCap sets the line cap style for the end of stroked lines.
//
// This method allows you to set the line cap style for the ends of stroked lines. You can choose from three predefined
// styles: LineCapRound, LineCapButt, and LineCapSquare.
func (dc *Context) SetLineCap(lineCap LineCap) {
	dc.lineCap = lineCap
}

// SetLineCapRound sets the line cap style to Round for the end of stroked lines.
//
// This method sets the line cap style to Round for the ends of stroked lines, which results in rounded line endings.
func (dc *Context) SetLineCapRound() {
	dc.lineCap = LineCapRound
}

// SetLineCapButt sets the line cap style to Butt for the end of stroked lines.
//
// This method sets the line cap style to Butt for the ends of stroked lines, which results in flat line endings.
func (dc *Context) SetLineCapButt() {
	dc.lineCap = LineCapButt
}

// SetLineCapSquare sets the line cap style to Square for the end of stroked lines.
//
// This method sets the line cap style to Square for the ends of stroked lines, which results in square line endings.
func (dc *Context) SetLineCapSquare() {
	dc.lineCap = LineCapSquare
}

// SetLineJoin sets the line join style for the intersection of stroked lines.
//
// This method allows you to set the line join style for the intersections of stroked lines. You can choose from two predefined
// styles: LineJoinRound and LineJoinBevel.
func (dc *Context) SetLineJoin(lineJoin LineJoin) {
	dc.lineJoin = lineJoin
}

// SetLineJoinRound sets the line join style to Round for the intersection of stroked lines.
//
// This method sets the line join style to Round for the intersections of stroked lines, which results in rounded intersections.
func (dc *Context) SetLineJoinRound() {
	dc.lineJoin = LineJoinRound
}

// SetLineJoinBevel sets the line join style to Bevel for the intersection of stroked lines.
//
// This method sets the line join style to Bevel for the intersections of stroked lines, which results in beveled intersections.
func (dc *Context) SetLineJoinBevel() {
	dc.lineJoin = LineJoinBevel
}

// SetFillRule sets the fill rule for determining the interior of filled shapes.
//
// This method allows you to set the fill rule for determining the interior of filled shapes. You can choose from two predefined
// fill rules: FillRuleWinding and FillRuleEvenOdd.
func (dc *Context) SetFillRule(fillRule FillRule) {
	dc.fillRule = fillRule
}

// SetFillRuleWinding sets the fill rule to Winding for determining the interior of filled shapes.
//
// This method sets the fill rule to Winding for determining the interior of filled shapes. In this rule, a point is considered
// inside a shape if a ray from the point in any direction crosses an odd number of times with the shape's boundary.
func (dc *Context) SetFillRuleWinding() {
	dc.fillRule = FillRuleWinding
}

// SetFillRuleEvenOdd sets the fill rule to Even-Odd for determining the interior of filled shapes.
//
// This method sets the fill rule to Even-Odd for determining the interior of filled shapes. In this rule, a point is considered
// inside a shape if a ray from the point in any direction crosses an odd number of times with the shape's boundary.
func (dc *Context) SetFillRuleEvenOdd() {
	dc.fillRule = FillRuleEvenOdd
}

// setFillAndStrokeColor sets the fill and stroke color of the rendering context.
//
// This method allows you to set both the fill and stroke colors of the rendering context to the specified color 'c'. The 'c'
// parameter should be a color.Color value representing the desired color.
func (dc *Context) setFillAndStrokeColor(c color.Color) {
	dc.color = c
	dc.fillPattern = NewSolidPattern(c)
	dc.strokePattern = NewSolidPattern(c)
}

// SetFillStyle sets the fill style of the rendering context.
//
// This method allows you to set the fill style of the rendering context using the specified 'pattern'. The 'pattern' can
// be any implementation of the Pattern interface, such as SolidPattern, Gradient, or TexturePattern. If a SolidPattern is
// used, the method also updates the current 'color' of the context.
func (dc *Context) SetFillStyle(pattern Pattern) {
	// if pattern is SolidPattern, also change dc.color(for dc.Clear, dc.drawString)
	if fillStyle, ok := pattern.(*solidPattern); ok {
		dc.color = fillStyle.color
	}
	dc.fillPattern = pattern
}

// SetStrokeStyle sets the stroke style of the rendering context.
//
// This method allows you to set the stroke style of the rendering context using the specified 'pattern'. The 'pattern' can
// be any implementation of the Pattern interface, such as SolidPattern, Gradient, or TexturePattern.
func (dc *Context) SetStrokeStyle(pattern Pattern) {
	dc.strokePattern = pattern
}

// SetColor sets the fill and stroke color of the rendering context to the specified color.
//
// This method sets both the fill and stroke colors of the rendering context to the specified color 'c'. The 'c' parameter
// should be a color.Color value representing the desired color.
func (dc *Context) SetColor(c color.Color) {
	dc.setFillAndStrokeColor(c)
}

// SetHexColor sets the fill and stroke color of the rendering context using a hexadecimal color string.
//
// This method sets both the fill and stroke colors of the rendering context to the colors specified in the hexadecimal
// color string 'x'. The 'x' parameter should be a string in the format "#RGB" or "#RRGGBB" or "#RRGGBBAA".
func (dc *Context) SetHexColor(x string) {
	r, g, b, a := ParseHexColor(x)
	dc.SetRGBA255(r, g, b, a)
}

// SetRGBA255 sets the fill and stroke color of the rendering context to the specified RGBA color using 8-bit values.
//
// This method sets both the fill and stroke colors of the rendering context to the specified RGBA color using 8-bit values
// for red, green, blue, and alpha. The 'r', 'g', 'b', and 'a' parameters should be integers in the range [0, 255].
func (dc *Context) SetRGBA255(r, g, b, a int) {
	dc.color = color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	dc.setFillAndStrokeColor(dc.color)
}

// SetRGB255 sets the fill and stroke color of the rendering context to the specified RGB color using 8-bit values.
//
// This method sets both the fill and stroke colors of the rendering context to the specified RGB color using 8-bit values
// for red, green, and blue. The 'r', 'g', and 'b' parameters should be integers in the range [0, 255], and an alpha value
// of 255 (fully opaque) is used.
func (dc *Context) SetRGB255(r, g, b int) {
	dc.SetRGBA255(r, g, b, 255)
}

// SetRGBA sets the fill and stroke color of the rendering context to the specified RGBA color using floating-point values.
//
// This method sets both the fill and stroke colors of the rendering context to the specified RGBA color using floating-point
// values in the range [0.0, 1.0]. The 'r', 'g', 'b', and 'a' parameters should represent the red, green, blue, and alpha
// components, respectively.
func (dc *Context) SetRGBA(r, g, b, a float64) {
	dc.color = color.NRGBA{
		R: uint8(r * 255),
		G: uint8(g * 255),
		B: uint8(b * 255),
		A: uint8(a * 255),
	}
	dc.setFillAndStrokeColor(dc.color)
}

// SetRGB sets the fill and stroke color of the rendering context to the specified RGB color using floating-point values.
//
// This method sets both the fill and stroke colors of the rendering context to the specified RGB color using floating-point
// values in the range [0.0, 1.0]. The 'r', 'g', and 'b' parameters should represent the red, green, and blue components, and
// an alpha value of 1.0 (fully opaque) is used.
func (dc *Context) SetRGB(r, g, b float64) {
	dc.SetRGBA(r, g, b, 1)
}

// MoveTo starts a new subpath at the specified point (x, y).
//
// This method begins a new subpath in the rendering context with the current point set to the specified coordinates (x, y).
// If there was a previous subpath and the current point was set, this method adds that point to the fill path and starts the
// stroke path from the new point.
func (dc *Context) MoveTo(x, y float64) {
	if dc.hasCurrent {
		dc.fillPath.Add1(dc.start.Fixed())
	}
	x, y = dc.TransformPoint(x, y)
	p := Point{x, y}
	dc.strokePath.Start(p.Fixed())
	dc.fillPath.Start(p.Fixed())
	dc.start = p
	dc.current = p
	dc.hasCurrent = true
}

// LineTo adds a straight line segment to the current subpath.
//
// This method appends a straight line segment from the current point to the specified coordinates (x, y) to both the fill and
// stroke paths. If there is no current subpath, this method behaves as if MoveTo() was called with the same coordinates.
func (dc *Context) LineTo(x, y float64) {
	if !dc.hasCurrent {
		dc.MoveTo(x, y)
	} else {
		x, y = dc.TransformPoint(x, y)
		p := Point{x, y}
		dc.strokePath.Add1(p.Fixed())
		dc.fillPath.Add1(p.Fixed())
		dc.current = p
	}
}

// QuadraticTo adds a quadratic Bézier curve segment to the current subpath.
//
// This method appends a quadratic Bézier curve segment defined by the current point (start point) and two control points
// (x1, y1 and x2, y2) to both the fill and stroke paths. If there is no current subpath, this method behaves as if MoveTo()
// was called with the starting point coordinates (x1, y1).
func (dc *Context) QuadraticTo(x1, y1, x2, y2 float64) {
	if !dc.hasCurrent {
		dc.MoveTo(x1, y1)
	}
	x1, y1 = dc.TransformPoint(x1, y1)
	x2, y2 = dc.TransformPoint(x2, y2)
	p1 := Point{x1, y1}
	p2 := Point{x2, y2}
	dc.strokePath.Add2(p1.Fixed(), p2.Fixed())
	dc.fillPath.Add2(p1.Fixed(), p2.Fixed())
	dc.current = p2
}

// CubicTo adds a cubic Bézier curve segment to the current subpath.
//
// This method appends a cubic Bézier curve segment defined by the current point (start point) and three control points
// (x1, y1, x2, y2, x3, y3) to both the fill and stroke paths. If there is no current subpath, this method behaves as if
// MoveTo() was called with the starting point coordinates (x1, y1).
func (dc *Context) CubicTo(x1, y1, x2, y2, x3, y3 float64) {
	if !dc.hasCurrent {
		dc.MoveTo(x1, y1)
	}

	x0, y0 := dc.current.X, dc.current.Y
	x1, y1 = dc.TransformPoint(x1, y1)
	x2, y2 = dc.TransformPoint(x2, y2)
	x3, y3 = dc.TransformPoint(x3, y3)
	points := CubicBezier(x0, y0, x1, y1, x2, y2, x3, y3)
	previous := dc.current.Fixed()

	for _, p := range points[1:] {
		f := p.Fixed()
		if f == previous {
			// TODO: this fixes some rendering issues but not all
			continue
		}
		previous = f
		dc.strokePath.Add1(f)
		dc.fillPath.Add1(f)
		dc.current = p
	}
}

// ClosePath closes the current subpath by adding a straight line to the starting point.
//
// This method closes the current subpath by appending a straight line segment from the current point to the starting point
// of the subpath. If there is no current subpath, this method has no effect.
func (dc *Context) ClosePath() {
	if dc.hasCurrent {
		dc.strokePath.Add1(dc.start.Fixed())
		dc.fillPath.Add1(dc.start.Fixed())
		dc.current = dc.start
	}
}

// ClearPath clears both the fill and stroke paths, and resets the current point.
//
// This method clears both the fill and stroke paths and resets the current point, effectively starting a new subpath.
func (dc *Context) ClearPath() {
	dc.strokePath.Clear()
	dc.fillPath.Clear()
	dc.hasCurrent = false
}

// NewSubPath starts a new subpath without adding any line segments.
//
// This method begins a new subpath without adding any line segments to it. If there was a previous subpath and the current
// point was set, this method adds that point to the fill path.
func (dc *Context) NewSubPath() {
	if dc.hasCurrent {
		dc.fillPath.Add1(dc.start.Fixed())
	}
	dc.hasCurrent = false
}

// capper returns the raster.Capper based on the current line cap style.
//
// This method returns a raster.Capper based on the current line cap style (LineCapButt, LineCapRound, or LineCapSquare).
// It is used in the stroke operation to determine how line ends should be drawn.
func (dc *Context) capper() raster.Capper {
	switch dc.lineCap {
	case LineCapButt:
		return raster.ButtCapper
	case LineCapRound:
		return raster.RoundCapper
	case LineCapSquare:
		return raster.SquareCapper
	}

	return nil
}

// joiner returns the raster.Joiner based on the current line join style.
//
// This method returns a raster.Joiner based on the current line join style (LineJoinBevel or LineJoinRound).
// It is used in the stroke operation to determine how line segments are joined at corners.
func (dc *Context) joiner() raster.Joiner {
	switch dc.lineJoin {
	case LineJoinBevel:
		return raster.BevelJoiner
	case LineJoinRound:
		return raster.RoundJoiner
	}

	return nil
}

// stroke applies stroke painting to the current path.
//
// This method applies stroke painting to the current path using the specified painter. If a dash pattern is set, it applies
// the dashed pattern to the path. It also uses the current line width, line cap, and line join styles for stroke rendering.
// The resulting stroke is rendered by the painter.
func (dc *Context) stroke(painter raster.Painter) {
	path := dc.strokePath
	if len(dc.dashes) > 0 {
		path = dashed(path, dc.dashes, dc.dashOffset)
	} else {
		// TODO: this is a temporary workaround to remove tiny segments
		// that result in rendering issues
		path = rasterPath(flattenPath(path))
	}
	r := dc.rasterizer
	r.UseNonZeroWinding = true
	r.Clear()
	r.AddStroke(path, fix(dc.lineWidth), dc.capper(), dc.joiner())
	r.Rasterize(painter)
}

// fill applies fill painting to the current path.
//
// This method applies fill painting to the current path using the specified painter. It takes into account the current fill
// rule to determine how the path should be filled. If there is a current point and a subpath, it appends the starting point
// to the path before filling. The resulting fill is rendered by the painter.
func (dc *Context) fill(painter raster.Painter) {
	path := dc.fillPath
	if dc.hasCurrent {
		path = make(raster.Path, len(dc.fillPath))
		copy(path, dc.fillPath)
		path.Add1(dc.start.Fixed())
	}
	r := dc.rasterizer
	r.UseNonZeroWinding = dc.fillRule == FillRuleWinding
	r.Clear()
	r.AddPath(path)
	r.Rasterize(painter)
}

// StrokePreserve applies the stroke operation to the current path, preserving the path for future rendering.
//
// This method applies the stroke operation to the current path and preserves the path for future rendering. It selects
// the appropriate painter based on the stroke pattern and the presence of a mask. It then calls the stroke method to
// render the stroke. After the stroke is applied, the current path remains intact.
func (dc *Context) StrokePreserve() {
	var painter raster.Painter
	if dc.mask == nil {
		if pattern, ok := dc.strokePattern.(*solidPattern); ok {
			// with a nil mask and a solid color pattern, we can be more efficient
			// TODO: refactor so we don't have to do this type assertion stuff?
			p := raster.NewRGBAPainter(dc.im)
			p.SetColor(pattern.color)
			painter = p
		}
	}
	if painter == nil {
		painter = newPatternPainter(dc.im, dc.mask, dc.strokePattern)
	}
	dc.stroke(painter)
}

// Stroke applies the stroke operation to the current path and clears the path.
//
// This method applies the stroke operation to the current path and clears the path. It uses the StrokePreserve method to
// render the stroke and then clears the path, making it ready for further path construction.
func (dc *Context) Stroke() {
	dc.StrokePreserve()
	dc.ClearPath()
}

// FillPreserve applies the fill operation to the current path, preserving the path for future rendering.
//
// This method applies the fill operation to the current path and preserves the path for future rendering. It selects
// the appropriate painter based on the fill pattern and the presence of a mask. It then calls the fill method to render
// the fill. After the fill is applied, the current path remains intact.
func (dc *Context) FillPreserve() {
	var painter raster.Painter
	if dc.mask == nil {
		if pattern, ok := dc.fillPattern.(*solidPattern); ok {
			// with a nil mask and a solid color pattern, we can be more efficient
			// TODO: refactor so we don't have to do this type assertion stuff?
			p := raster.NewRGBAPainter(dc.im)
			p.SetColor(pattern.color)
			painter = p
		}
	}
	if painter == nil {
		painter = newPatternPainter(dc.im, dc.mask, dc.fillPattern)
	}
	dc.fill(painter)
}

// Fill applies the fill operation to the current path and clears the path.
//
// This method applies the fill operation to the current path and clears the path. It uses the FillPreserve method to
// render the fill and then clears the path, making it ready for further path construction.
func (dc *Context) Fill() {
	dc.FillPreserve()
	dc.ClearPath()
}

// ClipPreserve applies the clip operation to the current path, preserving the path for future rendering.
//
// This method applies the clip operation to the current path and preserves the path for future rendering. It creates a
// mask using the current path, and if a mask already exists, it combines the new mask with the existing one. The mask
// is used to define the clipping area for future rendering.
func (dc *Context) ClipPreserve() {
	clip := image.NewAlpha(image.Rect(0, 0, dc.width, dc.height))
	painter := raster.NewAlphaOverPainter(clip)
	dc.fill(painter)
	if dc.mask == nil {
		dc.mask = clip
	} else {
		mask := image.NewAlpha(image.Rect(0, 0, dc.width, dc.height))
		draw.DrawMask(mask, mask.Bounds(), clip, image.Point{}, dc.mask, image.Point{}, draw.Over)
		dc.mask = mask
	}
}

// SetMask sets the mask of the rendering context.
//
// This method sets the mask of the rendering context to the specified image.Alpha mask. The mask size must match the
// size of the context's image. If the mask size does not match, an error is returned.
func (dc *Context) SetMask(mask *image.Alpha) error {
	if mask.Bounds().Size() != dc.im.Bounds().Size() {
		return errors.New("mask size must match context size")
	}

	dc.mask = mask

	return nil
}

// AsMask converts the context's image to a mask.
//
// This method converts the context's image to a mask, which can be used for masking in future rendering operations. The
// original image is copied to create the mask, and the mask is returned.
func (dc *Context) AsMask() *image.Alpha {
	mask := image.NewAlpha(dc.im.Bounds())

	draw.Draw(mask, dc.im.Bounds(), dc.im, image.Point{}, draw.Src)

	return mask
}

// InvertMask inverts the current mask or creates a new one if none exists.
//
// This method inverts the current mask if one exists. If no mask is set, it creates a new alpha mask and inverts it.
// The inverted mask can be used for masking in future rendering operations.
func (dc *Context) InvertMask() {
	if dc.mask == nil {
		dc.mask = image.NewAlpha(dc.im.Bounds())
	} else {
		for i, a := range dc.mask.Pix {
			dc.mask.Pix[i] = 255 - a
		}
	}
}

// Clip applies the clip operation to the current path and clears the path.
//
// This method applies the clip operation to the current path and clears the path. It uses the ClipPreserve method to
// define the clipping area for future rendering. After applying the clip, the current path is cleared.
func (dc *Context) Clip() {
	dc.ClipPreserve()
	dc.ClearPath()
}

// ResetClip clears the current clipping mask.
//
// This method clears the current clipping mask, allowing for unrestricted rendering. It effectively removes any
// clipping effects applied by previous Clip or ClipPreserve operations.
func (dc *Context) ResetClip() {
	dc.mask = nil
}

// Clear sets the entire context's image to a uniform color, effectively clearing the rendering area.
//
// This method fills the entire context's image with the specified uniform color, effectively clearing the rendering area.
// The uniform color is set to the context's current color. The Clear operation overwrites any existing content in the
// image with the specified color, making the entire image uniform.
func (dc *Context) Clear() {
	src := image.NewUniform(dc.color)
	draw.Draw(dc.im, dc.im.Bounds(), src, image.Point{}, draw.Src)
}

// SetPixel sets the color of a single pixel at the specified coordinates.
//
// This method sets the color of a single pixel at the given (x, y) coordinates in the context's image to the current color.
// It effectively paints a single pixel with the specified color.
func (dc *Context) SetPixel(x, y int) {
	dc.im.Set(x, y, dc.color)
}

// DrawPoint draws a filled circle at the specified point.
//
// This method draws a filled circle with the specified radius (r) centered at the coordinates (x, y).
func (dc *Context) DrawPoint(x, y, r float64) {
	dc.Push()
	tx, ty := dc.TransformPoint(x, y)
	dc.Identity()
	dc.DrawCircle(tx, ty, r)
	dc.Pop()
}

// DrawLine draws a straight line segment between two points.
//
// This method draws a straight line segment between the points (x1, y1) and (x2, y2).
func (dc *Context) DrawLine(x1, y1, x2, y2 float64) {
	dc.MoveTo(x1, y1)
	dc.LineTo(x2, y2)
}

// DrawRectangle draws a filled rectangle at the specified coordinates with the given width and height.
//
// This method draws a filled rectangle with its top-left corner at coordinates (x, y), a width of 'w', and a height of 'h'.
func (dc *Context) DrawRectangle(x, y, w, h float64) {
	dc.NewSubPath()
	dc.MoveTo(x, y)
	dc.LineTo(x+w, y)
	dc.LineTo(x+w, y+h)
	dc.LineTo(x, y+h)
	dc.ClosePath()
}

// DrawRoundedRectangle draws a filled rounded rectangle with the specified coordinates, dimensions, and corner radius.
//
// This method draws a filled rounded rectangle with its top-left corner at coordinates (x, y), a width of 'w', a height of 'h', and rounded corners with a radius of 'r'.
func (dc *Context) DrawRoundedRectangle(x, y, w, h, r float64) {
	x0, x1, x2, x3 := x, x+r, x+w-r, x+w
	y0, y1, y2, y3 := y, y+r, y+h-r, y+h
	dc.NewSubPath()
	dc.MoveTo(x1, y0)
	dc.LineTo(x2, y0)
	dc.DrawArc(x2, y1, r, Radians(270), Radians(360))
	dc.LineTo(x3, y2)
	dc.DrawArc(x2, y2, r, Radians(0), Radians(90))
	dc.LineTo(x1, y3)
	dc.DrawArc(x1, y2, r, Radians(90), Radians(180))
	dc.LineTo(x0, y1)
	dc.DrawArc(x1, y1, r, Radians(180), Radians(270))
	dc.ClosePath()
}

// DrawEllipticalArc draws a series of connected line segments to approximate an elliptical arc.
//
// This method approximates an elliptical arc within the specified bounding box defined by (x, y), width 'rx', and height 'ry'. The arc is drawn between the angles 'angle1' and 'angle2'.
func (dc *Context) DrawEllipticalArc(x, y, rx, ry, angle1, angle2 float64) {
	const n = 16
	for i := 0; i < n; i++ {
		p1 := float64(i+0) / n
		p2 := float64(i+1) / n
		a1 := angle1 + (angle2-angle1)*p1
		a2 := angle1 + (angle2-angle1)*p2
		x0 := x + rx*math.Cos(a1)
		y0 := y + ry*math.Sin(a1)
		x1 := x + rx*math.Cos((a1+a2)/2)
		y1 := y + ry*math.Sin((a1+a2)/2)
		x2 := x + rx*math.Cos(a2)
		y2 := y + ry*math.Sin(a2)
		cx := 2*x1 - x0/2 - x2/2
		cy := 2*y1 - y0/2 - y2/2
		if i == 0 {
			if dc.hasCurrent {
				dc.LineTo(x0, y0)
			} else {
				dc.MoveTo(x0, y0)
			}
		}
		dc.QuadraticTo(cx, cy, x2, y2)
	}
}

// DrawEllipse draws an ellipse within the specified bounding box.
//
// This method draws an ellipse within the bounding box defined by the center point (x, y) and the radii 'rx' and 'ry'. The ellipse is drawn as a series of connected line segments that approximate the ellipse shape.
func (dc *Context) DrawEllipse(x, y, rx, ry float64) {
	dc.NewSubPath()
	dc.DrawEllipticalArc(x, y, rx, ry, 0, 2*math.Pi)
	dc.ClosePath()
}

// DrawArc draws an arc within the specified parameters.
//
// This method draws an arc within a circle with center (x, y) and radius 'r' starting at 'angle1' and ending at 'angle2' (in radians). The arc is drawn as a series of connected line segments.
func (dc *Context) DrawArc(x, y, r, angle1, angle2 float64) {
	dc.DrawEllipticalArc(x, y, r, r, angle1, angle2)
}

// DrawCircle draws a circle centered at (x, y) with the specified radius 'r'.
//
// This method draws a circle with center at (x, y) and radius 'r'. The circle is drawn as a series of connected line segments.
func (dc *Context) DrawCircle(x, y, r float64) {
	dc.NewSubPath()
	dc.DrawEllipticalArc(x, y, r, r, 0, 2*math.Pi)
	dc.ClosePath()
}

// DrawRegularPolygon draws a regular polygon with 'n' sides.
//
// This method draws a regular polygon with 'n' sides, centered at (x, y) with a given radius 'r'. The 'rotation' parameter specifies the initial rotation angle. The polygon is drawn as a series of connected line segments.
func (dc *Context) DrawRegularPolygon(n int, x, y, r, rotation float64) {
	angle := 2 * math.Pi / float64(n)
	rotation -= math.Pi / 2
	if n%2 == 0 {
		rotation += angle / 2
	}
	dc.NewSubPath()
	for i := 0; i < n; i++ {
		a := rotation + angle*float64(i)
		dc.LineTo(x+r*math.Cos(a), y+r*math.Sin(a))
	}
	dc.ClosePath()
}

// SetInterpolator sets the drawing interpolator for the context.
//
// This method allows you to set a custom drawing interpolator for the context. The interpolator is responsible for determining how paths and shapes are transformed and rendered, affecting the quality and smoothness of the output. Providing a nil interpolator will result in a panic.
func (dc *Context) SetInterpolator(interp draw.Interpolator) {
	if interp == nil {
		panic(errors.New("gg: invalid interpolator"))
	}
	dc.interp = interp
}

// DrawImage draws an image at the specified (x, y) coordinates.
//
// This method draws the given image at the specified (x, y) coordinates in the context's image. The top-left corner of the image will be positioned at the (x, y) coordinates.
func (dc *Context) DrawImage(im image.Image, x, y int) {
	dc.DrawImageAnchored(im, x, y, 0, 0)
}

// DrawImageAnchored draws an image anchored at the specified (x, y) coordinates.
// The anchor point is x - w * ax, y - h * ay, where w, h is the size of the
// image.
//
// This method draws the given image anchored at the specified (x, y) coordinates in the context's image. The anchor point (ax, ay) determines the relative position within the image that aligns with the (x, y) coordinates. For example, (0, 0) represents the top-left corner of the image, (0.5, 0.5) represents the center, and (1, 1) represents the bottom-right corner.
func (dc *Context) DrawImageAnchored(im image.Image, x, y int, ax, ay float64) {
	s := im.Bounds().Size()
	x -= int(ax * float64(s.X))
	y -= int(ay * float64(s.Y))
	var (
		fx  = float64(x)
		fy  = float64(y)
		m   = dc.matrix.Translate(fx, fy)
		s2d = f64.Aff3{m.XX, m.XY, m.X0, m.YX, m.YY, m.Y0}
		opt *draw.Options
	)
	if dc.mask != nil {
		opt = &draw.Options{
			DstMask:  dc.mask,
			DstMaskP: image.Point{},
		}
	}
	dc.interp.Transform(dc.im, s2d, im, im.Bounds(), draw.Over, opt)
}

// SetFontFace sets the font face for text rendering.
//
// This method sets the font face for rendering text. The provided `fontFace` is used for subsequent text operations. You can obtain a font face using the `LoadFontFace` or related methods.
func (dc *Context) SetFontFace(fontFace font.Face) {
	dc.fontFace = fontFace
	dc.fontHeight = (float64(fontFace.Metrics().Height) / 64) * 72 / 96
}

// LoadFontFace loads a font face from a file and sets it for text rendering.
//
// This method loads a font face from the specified file path and sets it for rendering text. The `points` parameter determines the font size in points.
func (dc *Context) LoadFontFace(path string, points float64) error {
	face, err := LoadFontFace(path, points)
	if err != nil {
		return err
	}

	dc.fontFace = face
	dc.fontHeight = points * 72 / 96

	return nil
}

// LoadFontFaceFromReader loads a font face from an io.Reader and sets it for text rendering.
//
// This method loads a font face from the specified io.Reader and sets it for rendering text. The `points` parameter determines the font size in points.
func (dc *Context) LoadFontFaceFromReader(r io.Reader, points float64) error {
	face, err := LoadFontFaceFromReader(r, points)
	if err != nil {
		return err
	}

	dc.fontFace = face
	dc.fontHeight = points * 72 / 96

	return nil
}

// LoadFontFaceFromFS loads a font face from a file system (fs.FS) and sets it for text rendering.
//
// This method loads a font face from the file system (fs.FS) and sets it for rendering text. The `points` parameter determines the font size in points.
func (dc *Context) LoadFontFaceFromFS(fsys fs.FS, path string, points float64) error {
	face, err := LoadFontFaceFromFS(fsys, path, points)
	if err != nil {
		return err
	}

	dc.fontFace = face
	dc.fontHeight = points * 72 / 96

	return nil
}

// LoadFontFaceFromBytes loads a font face from raw font data and sets it for text rendering.
//
// This method loads a font face from raw font data provided as a byte slice and sets it for rendering text. The `points` parameter determines the font size in points.
func (dc *Context) LoadFontFaceFromBytes(raw []byte, points float64) error {
	face, err := LoadFontFaceFromBytes(raw, points)
	if err != nil {
		return err
	}

	dc.fontFace = face
	dc.fontHeight = points * 72 / 96

	return nil
}

// FontHeight returns the height of the currently set font face.
//
// This method returns the height of the font face currently set for text rendering. The height is measured in points.
func (dc *Context) FontHeight() float64 {
	return dc.fontHeight
}

// drawString renders a text string onto an RGBA image at the specified coordinates.
//
// This method renders the given text string `s` onto the provided RGBA image `im` at the specified (x, y) coordinates. The text is rendered using the current font face, color, and other text rendering settings of the context.
func (dc *Context) drawString(im *image.RGBA, s string, x, y float64) {
	d := &font.Drawer{
		Dst:  im,
		Src:  image.NewUniform(dc.color),
		Face: dc.fontFace,
		Dot:  fixp(x, y),
	}
	// based on Drawer.DrawString() in golang.org/x/image/font/font.go
	prevC := rune(-1)
	for _, c := range s {
		if prevC >= 0 {
			d.Dot.X += d.Face.Kern(prevC, c)
		}
		dr, mask, maskp, advance, ok := d.Face.Glyph(d.Dot, c)
		if !ok {
			// TODO: is falling back on the U+FFFD glyph the responsibility of
			// the Drawer or the Face?
			// TODO: set prevC = '\ufffd'?
			continue
		}
		sr := dr.Sub(dr.Min)
		fx, fy := float64(dr.Min.X), float64(dr.Min.Y)
		m := dc.matrix.Translate(fx, fy)
		s2d := f64.Aff3{m.XX, m.XY, m.X0, m.YX, m.YY, m.Y0}
		dc.interp.Transform(d.Dst, s2d, d.Src, sr, draw.Over, &draw.Options{
			SrcMask:  mask,
			SrcMaskP: maskp,
		})
		d.Dot.X += advance
		prevC = c
	}
}

// DrawString renders a text string at the specified coordinates.
//
// This method renders the given text string `s` at the specified (x, y) coordinates on the context's image. The text is rendered using the current font face, color, and other text rendering settings of the context. The (x, y) coordinates represent the baseline position for rendering the text.
func (dc *Context) DrawString(s string, x, y float64) {
	dc.DrawStringAnchored(s, x, y, 0, 0)
}

// DrawStringAnchored renders a text string anchored at the specified coordinates.
//
// This method renders the given text string `s` anchored at the specified (x, y) coordinates on the context's image. The anchor point is determined by the `ax` (X-axis) and `ay` (Y-axis) values, which represent the relative position within the text bounding box. The text is rendered using the current font face, color, and other text rendering settings of the context.
func (dc *Context) DrawStringAnchored(s string, x, y, ax, ay float64) {
	w, h := dc.MeasureString(s)
	x -= ax * w
	y += ay * h
	if dc.mask == nil {
		dc.drawString(dc.im, s, x, y)
	} else {
		im := image.NewRGBA(image.Rect(0, 0, dc.width, dc.height))
		dc.drawString(im, s, x, y)
		draw.DrawMask(dc.im, dc.im.Bounds(), im, image.Point{}, dc.mask, image.Point{}, draw.Over)
	}
}

// DrawStringWrapped renders a text string wrapped within a specified width.
//
// This method renders the given text string `s` wrapped within the specified `width` while anchored at the (x, y) coordinates. The text is wrapped into multiple lines to fit the given width. The `ax` (X-axis) and `ay` (Y-axis) values determine the anchor point's relative position within the text bounding box, and the `lineSpacing` controls the vertical spacing between lines. The `align` parameter specifies the horizontal alignment of the text.
func (dc *Context) DrawStringWrapped(s string, x, y, ax, ay, width, lineSpacing float64, align Align) {
	lines := dc.WordWrap(s, width)

	// sync h formula with MeasureMultilineString
	h := float64(len(lines)) * dc.fontHeight * lineSpacing
	h -= (lineSpacing - 1) * dc.fontHeight

	x -= ax * width
	y -= ay * h
	switch align {
	case AlignLeft:
		ax = 0
	case AlignCenter:
		ax = 0.5
		x += width / 2
	case AlignRight:
		ax = 1
		x += width
	}
	ay = 1
	for _, line := range lines {
		dc.DrawStringAnchored(line, x, y, ax, ay)
		y += dc.fontHeight * lineSpacing
	}
}

// MeasureMultilineString measures the width and height of a multiline text string.
//
// This method calculates the dimensions of a multiline text string `s`, taking into account the specified `lineSpacing` factor for vertical line spacing. It returns the width and height of the multiline text in pixels.
func (dc *Context) MeasureMultilineString(s string, lineSpacing float64) (width, height float64) {
	lines := strings.Split(s, "\n")

	// sync h formula with DrawStringWrapped
	height = float64(len(lines)) * dc.fontHeight * lineSpacing
	height -= (lineSpacing - 1) * dc.fontHeight

	d := &font.Drawer{
		Face: dc.fontFace,
	}

	// max width from lines
	for _, line := range lines {
		adv := d.MeasureString(line)
		currentWidth := float64(adv >> 6) // from gg.Context.MeasureString
		if currentWidth > width {
			width = currentWidth
		}
	}

	return width, height
}

// MeasureString measures the width and height of a single-line text string.
//
// This method calculates the dimensions of a single-line text string `s` and returns its width and the standard line height (font height).
func (dc *Context) MeasureString(s string) (w, h float64) {
	d := &font.Drawer{
		Face: dc.fontFace,
	}
	a := d.MeasureString(s)

	return float64(a >> 6), dc.fontHeight
}

// WordWrap wraps a text string to fit within a specified width.
//
// This method takes a text string `s` and wraps it to fit within a given width `w`, breaking it into multiple lines as necessary to prevent text from exceeding the specified width. The result is returned as a slice of strings, each representing a wrapped line of text.
func (dc *Context) WordWrap(s string, w float64) []string {
	return wordWrap(dc, s, w)
}

// Identity resets the current transformation matrix to the identity matrix.
//
// This method sets the transformation matrix of the drawing context to the identity matrix, effectively removing any prior transformations.
func (dc *Context) Identity() {
	dc.matrix = Identity()
}

// Translate applies a translation to the current transformation matrix.
//
// This method translates the drawing context by the specified horizontal (x) and vertical (y) amounts. It modifies the current transformation matrix to reflect the translation.
func (dc *Context) Translate(x, y float64) {
	dc.matrix = dc.matrix.Translate(x, y)
}

// Scale applies a scaling transformation to the current matrix.
//
// This method scales the drawing context by the specified horizontal (x) and vertical (y) scaling factors. It modifies the current transformation matrix to reflect the scaling.
func (dc *Context) Scale(x, y float64) {
	dc.matrix = dc.matrix.Scale(x, y)
}

// ScaleAbout applies a scaling transformation about a specified point.
//
// This method scales the drawing context by the specified horizontal (sx) and vertical (sy) scaling factors around the point (x, y). It is equivalent to translating the context to the center of scaling, scaling it, and then translating it back to its original position.
func (dc *Context) ScaleAbout(sx, sy, x, y float64) {
	dc.Translate(x, y)
	dc.Scale(sx, sy)
	dc.Translate(-x, -y)
}

// Rotate applies a rotation transformation to the current transformation matrix.
//
// This method rotates the drawing context by the specified angle in radians. It modifies the current transformation matrix to reflect the rotation.
func (dc *Context) Rotate(angle float64) {
	dc.matrix = dc.matrix.Rotate(angle)
}

// RotateAbout applies a rotation transformation about a specified point.
//
// This method rotates the drawing context by the specified angle in radians around the point (x, y). It is equivalent to translating the context to the center of rotation, rotating it, and then translating it back to its original position.
func (dc *Context) RotateAbout(angle, x, y float64) {
	dc.Translate(x, y)
	dc.Rotate(angle)
	dc.Translate(-x, -y)
}

// Shear applies a shear transformation to the current transformation matrix.
//
// This method shears the drawing context by the specified horizontal (x) and vertical (y) shear factors. It modifies the current transformation matrix to reflect the shear.
func (dc *Context) Shear(x, y float64) {
	dc.matrix = dc.matrix.Shear(x, y)
}

// ShearAbout applies a shear transformation about a specified point.
//
// This method shears the drawing context by the specified horizontal (sx) and vertical (sy) shear factors around the point (x, y). It is equivalent to translating the context to the center of shearing, shearing it, and then translating it back to its original position.
func (dc *Context) ShearAbout(sx, sy, x, y float64) {
	dc.Translate(x, y)
	dc.Shear(sx, sy)
	dc.Translate(-x, -y)
}

// TransformPoint applies the current transformation matrix to the given (x, y) coordinates.
//
// This method transforms the (x, y) coordinates using the current transformation matrix. It returns the transformed coordinates (tx, ty).
func (dc *Context) TransformPoint(x, y float64) (tx, ty float64) {
	return dc.matrix.TransformPoint(x, y)
}

// InvertY inverts the Y-axis of the current transformation.
//
// This method inverts the Y-axis of the current drawing context. It effectively flips the vertical orientation of subsequent drawings.
func (dc *Context) InvertY() {
	dc.Translate(0, float64(dc.height))
	dc.Scale(1, -1)
}

// Push saves the current drawing context by creating a copy of it and pushing it onto the context stack.
//
// This method saves the current state of the drawing context by creating a copy of it and pushing it onto the context stack.
// You can later restore this state using the Pop method.
func (dc *Context) Push() {
	x := *dc
	dc.stack = append(dc.stack, &x)
}

// Pop restores the drawing context to a previously saved state from the context stack.
//
// This method pops the topmost state from the context stack and restores the drawing context to this saved state.
// It effectively reverts the current drawing state to a previous state that was saved using the Push method.
func (dc *Context) Pop() {
	var (
		before = *dc
		s      = dc.stack
		ctx    *Context
	)
	ctx, dc.stack = s[len(s)-1], s[:len(s)-1]
	*dc = *ctx
	dc.mask = before.mask
	dc.strokePath = before.strokePath
	dc.fillPath = before.fillPath
	dc.start = before.start
	dc.current = before.current
	dc.hasCurrent = before.hasCurrent
}

// AppendFrame appends the current context image to the frames slice.
//
// This method copies the current context image and appends it to the frames slice. This method is used for creating animated image.
func (dc *Context) AppendFrame() {
	if dc.frames == nil {
		dc.frames = make([]image.Image, 0)
	}
	dc.frames = append(dc.frames, RGBAToImage(dc.im))
}
