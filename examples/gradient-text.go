// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"image/color"
	"log"

	"github.com/arugaz/gg"
	"golang.org/x/image/font/gofont/gobold"
)

func main() {
	const (
		W = 1024
		H = 512
	)

	dc := gg.NewContext(W, H)
	dc.SetRGB(0, 0, 0)

	if err := dc.LoadFontFaceFromBytes(gobold.TTF, 128); err != nil {
		log.Fatalf("could not load bold font: %+v", err)
	}
	dc.DrawStringAnchored("Gradient Text", W/2, H/2, .5, .5)

	mask := dc.AsMask()

	g := gg.NewLinearGradient(0, 0, W, H)
	g.AddColorStop(0, color.RGBA{R: 255, A: 255})
	g.AddColorStop(1, color.RGBA{B: 255, A: 255})
	dc.SetFillStyle(g)

	if err := dc.SetMask(mask); err != nil {
		log.Fatalf("could not set mask: %+v", err)
	}

	dc.DrawRectangle(0, 0, W, H)
	dc.Fill()

	if err := gg.SavePNG("./testdata/_gradient-text.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
