// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gg

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/fs"
	"os"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// LoadImage loads an image from the specified file path and returns it as an image.Image.
// It opens the file, decodes the image, and returns the decoded image.
func LoadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	im, err := LoadImageFromReader(file)

	if err := file.Close(); err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return im, nil
}

// LoadImageFromFS loads an image from the file system provided by a file system (fs.FS)
// and returns it as an image.Image. It opens the specified file within the file system,
// decodes the image, and returns the decoded image.
func LoadImageFromFS(fsys fs.FS, path string) (image.Image, error) {
	var err error
	if fsys == nil {
		fsys, path, err = checkfsys(fsys, path)
		if err != nil {
			return nil, err
		}
	}

	file, err := fsys.Open(path)
	if err != nil {
		return nil, err
	}

	im, err := LoadImageFromReader(file)

	if err := file.Close(); err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return im, nil
}

// LoadImageFromBytes decodes an image from a byte slice and returns it as an image.Image.
func LoadImageFromBytes(raw []byte) (image.Image, error) {
	im, err := LoadImageFromReader(bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}

	return im, nil
}

// LoadImageFromReader loads an image from an io.Reader and returns it as an image.Image.
// It decodes the image, and returns the decoded image.
func LoadImageFromReader(r io.Reader) (image.Image, error) {
	im, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	return im, nil
}

// LoadPNG loads a PNG image from the specified file path and returns it as an image.Image.
// It opens the file, decodes the PNG image, and returns the decoded image.
func LoadPNG(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	im, err := LoadPNGFromReader(file)

	if err := file.Close(); err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return im, nil
}

// LoadPNGFromFS loads a PNG image from the specified file path within a file system (fs.FS)
// and returns it as an image.Image. It opens the specified file within the file system,
// decodes the PNG image, and returns the decoded image.
func LoadPNGFromFS(fsys fs.FS, path string) (image.Image, error) {
	var err error
	if fsys == nil {
		fsys, path, err = checkfsys(fsys, path)
		if err != nil {
			return nil, err
		}
	}

	file, err := fsys.Open(path)
	if err != nil {
		return nil, err
	}

	im, err := LoadPNGFromReader(file)

	if err := file.Close(); err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return im, nil
}

// LoadPNGFromBytes decodes a PNG image from a byte slice and returns it as an image.Image.
func LoadPNGFromBytes(raw []byte) (image.Image, error) {
	im, err := LoadPNGFromReader(bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}

	return im, nil
}

// LoadPNGFromReader loads a PNG image from an io.Reader and returns it as an image.Image.
// It decodes the PNG image, and returns the decoded image.
func LoadPNGFromReader(r io.Reader) (image.Image, error) {
	im, err := png.Decode(r)
	if err != nil {
		return nil, err
	}

	return im, nil
}

// LoadJPG loads a JPEG image from the specified file path and returns it as an image.Image.
// It opens the file, decodes the JPEG image, and returns the decoded image.
func LoadJPG(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	im, err := LoadJPGFromReader(file)

	if err := file.Close(); err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return im, nil
}

// LoadJPGFromFS loads a JPEG image from the specified file path within a file system (fs.FS)
// and returns it as an image.Image. It opens the specified file within the file system,
// decodes the JPEG image, and returns the decoded image.
func LoadJPGFromFS(fsys fs.FS, path string) (image.Image, error) {
	var err error
	if fsys == nil {
		fsys, path, err = checkfsys(fsys, path)
		if err != nil {
			return nil, err
		}
	}

	file, err := fsys.Open(path)
	if err != nil {
		return nil, err
	}

	im, err := LoadJPGFromReader(file)

	if err := file.Close(); err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return im, nil
}

// LoadJPGFromBytes decodes a JPEG image from a byte slice and returns it as an image.Image.
func LoadJPGFromBytes(raw []byte) (image.Image, error) {
	im, err := LoadJPGFromReader(bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}

	return im, nil
}

// LoadJPGFromReader loads a JPEG image from an io.Reader and returns it as an image.Image.
// It decodes the JPEG image, and returns the decoded image.
func LoadJPGFromReader(r io.Reader) (image.Image, error) {
	im, err := jpeg.Decode(r)
	if err != nil {
		return nil, err
	}

	return im, nil
}

// LoadFontFace loads a font face from a TrueType or OpenType font file at the specified file path
// and returns it as a font.Face with the specified point size.
func LoadFontFace(path string, points float64) (font.Face, error) {
	fontBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return LoadFontFaceFromBytes(fontBytes, points)
}

// LoadFontFaceFromReader loads a font face from a TrueType or OpenType font file at the specified io.Reader
// and returns it as a font.Face with the specified point size.
func LoadFontFaceFromReader(r io.Reader, points float64) (font.Face, error) {
	fontBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return LoadFontFaceFromBytes(fontBytes, points)
}

// LoadFontFaceFromFS loads a font face from a TrueType or OpenType font file at the specified file path within a file system (fs.FS)
// and returns it as a font.Face with the specified point size.
func LoadFontFaceFromFS(fsys fs.FS, path string, points float64) (font.Face, error) {
	var err error
	if fsys == nil {
		fsys, path, err = checkfsys(fsys, path)
		if err != nil {
			return nil, err
		}
	}

	fontBytes, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, err
	}

	return LoadFontFaceFromBytes(fontBytes, points)
}

// LoadFontFaceFromBytes creates a font face from a byte slice containing TrueType or OpenType font data,
// and sets the specified point size.
func LoadFontFaceFromBytes(raw []byte, points float64) (font.Face, error) {
	f, err := FontParse(raw)
	if err != nil {
		return nil, err
	}

	return FontNewFace(f, points)
}

// FontParse parses TrueType or OpenType font data from a byte slice and returns an *opentype.Font.
func FontParse(raw []byte) (*opentype.Font, error) {
	f, err := opentype.Parse(raw)

	if err != nil && strings.Contains(err.Error(), "font collection") {
		return FontParseCollection(raw)
	}

	return f, err
}

// FontParseCollection parses a TrueType or OpenType font collection from a byte slice and returns an *opentype.Font.
func FontParseCollection(raw []byte, i ...int) (*opentype.Font, error) {
	in := 0
	if len(i) > 0 {
		in = i[0]
	}

	c, err := opentype.ParseCollection(raw)
	if err != nil {
		return nil, err
	}

	return c.Font(in)
}

// FontNewFace creates a font face from a parsed *opentype.Font with the specified point size and hinting.
func FontNewFace(f *opentype.Font, points float64, hinting ...font.Hinting) (font.Face, error) {
	hint := font.HintingNone
	if len(hinting) > 0 {
		hint = hinting[0]
	}

	return opentype.NewFace(f, &opentype.FaceOptions{
		Size:    points,
		DPI:     72,
		Hinting: hint,
	})
}
