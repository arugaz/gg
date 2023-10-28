// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gg

import (
	"math"

	"github.com/golang/freetype/raster"
	"golang.org/x/image/math/fixed"
)

// flattenPath converts a raster.Path into a slice of slices of Point, representing flattened path segments.
//
// This function processes a raster.Path, which is typically a series of fixed-point commands and coordinates, and flattens it into a list of connected Point segments.
func flattenPath(p raster.Path) [][]Point {
	var result [][]Point
	var path []Point
	var cx, cy float64

	for i := 0; i < len(p); {
		switch p[i] {
		case 0:
			if len(path) > 0 {
				result = append(result, path)
				path = nil
			}
			x := unfix(p[i+1])
			y := unfix(p[i+2])
			path = append(path, Point{x, y})
			cx, cy = x, y
			i += 4
		case 1:
			x := unfix(p[i+1])
			y := unfix(p[i+2])
			path = append(path, Point{x, y})
			cx, cy = x, y
			i += 4
		case 2:
			x1 := unfix(p[i+1])
			y1 := unfix(p[i+2])
			x2 := unfix(p[i+3])
			y2 := unfix(p[i+4])
			points := QuadraticBezier(cx, cy, x1, y1, x2, y2)
			path = append(path, points...)
			cx, cy = x2, y2
			i += 6
		case 3:
			x1 := unfix(p[i+1])
			y1 := unfix(p[i+2])
			x2 := unfix(p[i+3])
			y2 := unfix(p[i+4])
			x3 := unfix(p[i+5])
			y3 := unfix(p[i+6])
			points := CubicBezier(cx, cy, x1, y1, x2, y2, x3, y3)
			path = append(path, points...)
			cx, cy = x3, y3
			i += 8
		default:
			panic("bad path")
		}
	}

	if len(path) > 0 {
		result = append(result, path)
	}

	return result
}

// dashPath creates a dashed representation of a path.
//
// This function takes a list of connected path segments represented as a slice of slices of Point and converts them into a dashed representation based on the specified dash pattern and offset.
func dashPath(paths [][]Point, dashes []float64, offset float64) [][]Point {
	var result [][]Point

	if len(dashes) == 0 {
		return paths
	}

	if len(dashes) == 1 {
		dashes = append(dashes, dashes[0])
	}

	for _, path := range paths {
		if len(path) < 2 {
			continue
		}
		previous := path[0]
		pathIndex := 1
		dashIndex := 0
		segmentLength := 0.0

		if offset != 0 {
			var totalLength float64
			for _, dashLength := range dashes {
				totalLength += dashLength
			}
			offset = math.Mod(offset, totalLength)
			if offset < 0 {
				offset += totalLength
			}
			for i, dashLength := range dashes {
				offset -= dashLength
				if offset < 0 {
					dashIndex = i
					segmentLength = dashLength + offset
					break
				}
			}
		}

		var segment []Point
		segment = append(segment, previous)

		for pathIndex < len(path) {
			dashLength := dashes[dashIndex]
			point := path[pathIndex]
			d := previous.Distance(point)
			maxd := dashLength - segmentLength
			if d > maxd {
				t := maxd / d
				p := previous.Interpolate(point, t)
				segment = append(segment, p)
				if dashIndex%2 == 0 && len(segment) > 1 {
					result = append(result, segment)
				}
				segment = nil
				segment = append(segment, p)
				segmentLength = 0
				previous = p
				dashIndex = (dashIndex + 1) % len(dashes)
			} else {
				segment = append(segment, point)
				previous = point
				segmentLength += d
				pathIndex++
			}
		}

		if dashIndex%2 == 0 && len(segment) > 1 {
			result = append(result, segment)
		}
	}

	return result
}

// rasterPath converts a path into a raster representation.
//
// This function takes a slice of slices of Point, where each inner slice represents a connected path segment. It converts these path segments into a raster representation for rendering.
func rasterPath(paths [][]Point) raster.Path {
	var result raster.Path

	for _, path := range paths {
		var previous fixed.Point26_6

		for i, point := range path {
			f := point.Fixed()

			if i == 0 {
				result.Start(f)
			} else {
				dx := f.X - previous.X
				dy := f.Y - previous.Y
				if dx < 0 {
					dx = -dx
				}
				if dy < 0 {
					dy = -dy
				}
				if dx+dy > 4 {
					// TODO: this is a hack for cases where two points are
					// too close - causes rendering issues with joins / caps
					result.Add1(f)
				}
			}

			previous = f
		}
	}

	return result
}

// dashed creates a dashed version of a raster.Path.
//
// This function takes an input raster.Path, a slice of dash lengths, and an offset, and creates a dashed version of the input path. The resulting raster.Path represents the dashed line.
func dashed(path raster.Path, dashes []float64, offset float64) raster.Path {
	return rasterPath(dashPath(flattenPath(path), dashes, offset))
}
