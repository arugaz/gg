// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/arugaz/gg"
)

func main() {
	const (
		W = 400
		H = 500
	)

	dc := gg.NewContext(W, H)

	im, err := gg.LoadPNG("./testdata/gopher.png")
	if err != nil {
		log.Fatalf("could not load gopher.png: %+v", err)
	}

	iw, ih := im.Bounds().Dx(), im.Bounds().Dy()

	dc.SetHexColor("#f00")
	dc.SetLineWidth(2)
	dc.DrawRectangle(1, 1, float64(W)-2., float64(H)-2)
	dc.Stroke()

	dc.SetHexColor("#0000ff")
	dc.SetLineWidth(2)
	dc.DrawRectangle(100, 210, float64(iw), float64(ih))
	dc.Stroke()
	dc.DrawImage(im, 100, 210)

	dc.SetHexColor("#00f")
	dc.SetLineWidth(2)
	dc.Rotate(gg.Radians(11))
	dc.DrawRectangle(100, 0, float64(iw), float64(ih)/2+20)
	dc.StrokePreserve()
	dc.Clip()
	dc.DrawImageAnchored(im, 100, 0, 0, 0)

	if err := gg.SavePNG("./testdata/_rotated-image.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
