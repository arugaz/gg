// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gg

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"

	"github.com/arugaz/gg/quantize"
)

// SavePNG saves an image as a PNG file at the specified path.
func SavePNG(path string, im image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	err = EncodePNG(file, im)

	if err := file.Close(); err != nil {
		return err
	}

	return err
}

// EncodePNG encodes an image as a PNG and writes it to the provided io.Writer.
func EncodePNG(w io.Writer, im image.Image) error {
	return png.Encode(w, im)
}

// SaveJPG saves an image as a JPEG file at the specified path with an optional quality setting.
func SaveJPG(path string, im image.Image, quality ...int) error {
	q := jpeg.DefaultQuality
	if len(quality) > 0 {
		q = quality[0]
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	opt := new(jpeg.Options)
	opt.Quality = q

	err = EncodeJPG(file, im, opt)

	if err := file.Close(); err != nil {
		return err
	}

	return err
}

// EncodeJPG encodes an image as a JPEG and writes it to the provided io.Writer with optional encoding options.
func EncodeJPG(w io.Writer, im image.Image, opt *jpeg.Options) error {
	return jpeg.Encode(w, im, opt)
}

// SaveGIF saves an image as a GIF file at the specified path with an optional delay between frames.
func SaveGIF(path string, im []image.Image, fps int) error {
	if im == nil || len(im) == 0 {
		return fmt.Errorf("no frames provided")
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	err = EncodeGIF(file, im, fps)

	if err := file.Close(); err != nil {
		return err
	}

	return err
}

// EncodeGIF encodes an image as a GIF and writes it to the provided io.Writer with an optional delay between frames.
func EncodeGIF(w io.Writer, im []image.Image, fps int) error {
	if im == nil || len(im) == 0 {
		return fmt.Errorf("no frames provided")
	}

	q := new(quantize.MedianCutQuantizer)
	q.Aggregation = quantize.Mode

	p := q.QuantizeMultiple(make(color.Palette, 0, 256), im)

	if q.ReserveTransparent {
		p = append(p, color.RGBA{R: 0, G: 0, B: 0, A: 0})
	}

	animGIF := new(gif.GIF)
	imLen := len(im)
	delay := int(1. / float32(fps) * 100.)

	animGIF.Image = make([]*image.Paletted, imLen)
	animGIF.Delay = make([]int, imLen)
	animGIF.Disposal = make([]byte, imLen)
	animId := make(map[uint32]uint8)

	for i := 0; i < imLen; i++ {
		bounds := im[i].Bounds()
		dst := image.NewPaletted(bounds, p)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				c := im[i].At(x, y)
				cr, cg, cb, ca := c.RGBA()
				o := dst.PixOffset(x, y)
				cid := (cr>>8)<<16 | cg | (cb >> 8)
				if q.ReserveTransparent && ca == 0 {
					dst.Pix[o] = 0
				} else if val, ok := animId[cid]; ok {
					dst.Pix[o] = val
				} else {
					val := uint8(p.Index(c))
					dst.Pix[o] = val
					animId[cid] = val
				}
			}
		}

		animGIF.Image[i] = dst
		animGIF.Delay[i] = delay
		animGIF.Disposal[i] = gif.DisposalBackground
	}

	return gif.EncodeAll(w, animGIF)
}
