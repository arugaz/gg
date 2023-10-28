// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gg

import (
	"fmt"
	"image"
	"image/draw"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/math/fixed"
)

// Radians converts an angle in degrees to radians.
func Radians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// Degrees converts an angle in radians to degrees.
func Degrees(radians float64) float64 {
	return radians * 180 / math.Pi
}

// ImageToRGBA converts an image.Image to an *image.RGBA by creating a new *image.RGBA with the same
// bounds as the source image and copying the image content into it.
func ImageToRGBA(src image.Image) *image.RGBA {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)

	return dst
}

// RGBAToImage converts an *image.RGBA to an image.Image by creating a new image.Image with the same
// bounds as the source image and copying the image content into it.
func RGBAToImage(src *image.RGBA) image.Image {
	dst := image.NewRGBA(src.Bounds())
	copy(dst.Pix, src.Pix)

	return dst
}

// ParseHexColor parses a hexadecimal color string (e.g., "#RGB" or "#RRGGBB" or "#RRGGBBAA") and returns the
// corresponding red (r), green (g), blue (b), and alpha (a) values. If the alpha component is not
// provided, it defaults to 255 (fully opaque).
func ParseHexColor(x string) (r, g, b, a int) {
	x = strings.TrimPrefix(x, "#")
	a = 255

	if len(x) == 3 {
		fmt.Sscanf(x, "%1x%1x%1x", &r, &g, &b)
		r |= r << 4
		g |= g << 4
		b |= b << 4
	}

	if len(x) == 6 {
		fmt.Sscanf(x, "%02x%02x%02x", &r, &g, &b)
	}

	if len(x) == 8 {
		fmt.Sscanf(x, "%02x%02x%02x%02x", &r, &g, &b, &a)
	}

	return r, g, b, a
}

// fixp creates a fixed.Point26_6 value from the provided floating-point coordinates.
func fixp(x, y float64) fixed.Point26_6 {
	return fixed.Point26_6{X: fix(x), Y: fix(y)}
}

// fix converts a floating-point value to a fixed-point representation in fixed.Int26_6 format.
// The function multiplies the input value by 64 and rounds the result to the nearest integer.
func fix(x float64) fixed.Int26_6 {
	return fixed.Int26_6(math.Round(x * 64))
}

// unfix converts a fixed-point value in fixed.Int26_6 format to a floating-point representation.
func unfix(x fixed.Int26_6) float64 {
	const shift, mask = 6, 1<<6 - 1

	if x >= 0 {
		return float64(x>>shift) + float64(x&mask)/64
	}

	x = -x
	if x >= 0 {
		return -(float64(x>>shift) + float64(x&mask)/64)
	}

	return 0
}

// checkfsys checks if a provided file system (fs.FS) is valid. If the file system is nil, it uses
// the default file system and resolves the path relative to the current working directory if the
// path is an absolute path.
func checkfsys(fsys fs.FS, path string) (fs.FS, string, error) {
	switch {
	case filepath.IsAbs(path):
		var (
			err  error
			orig = path
			root = filepath.FromSlash("/")
		)

		path, err = filepath.Rel(root, path)
		if err != nil {
			return nil, "", fmt.Errorf("could not find relative path for %q from %q: %w", orig, root, err)
		}

		fsys = os.DirFS(root)
	default:
		fsys = os.DirFS(".")
	}

	return fsys, path, nil
}
