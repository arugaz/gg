// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"image/color"
	"log"

	"github.com/arugaz/gg"
)

func main() {
	const (
		W = 500
		H = 400
	)

	dc := gg.NewContext(W, H)

	grad := gg.NewLinearGradient(20, 320, 400, 20)
	grad.AddColorStop(0, color.RGBA{G: 255, A: 255})
	grad.AddColorStop(1, color.RGBA{B: 255, A: 255})
	grad.AddColorStop(0.5, color.RGBA{R: 255, A: 255})

	dc.SetColor(color.White)
	dc.DrawRectangle(20, 20, 400-20, 300)
	dc.Stroke()

	dc.SetStrokeStyle(grad)
	dc.SetLineWidth(4)
	dc.MoveTo(10, 10)
	dc.LineTo(410, 10)
	dc.LineTo(410, 100)
	dc.LineTo(10, 100)
	dc.ClosePath()
	dc.Stroke()

	dc.SetFillStyle(grad)
	dc.MoveTo(10, 120)
	dc.LineTo(410, 120)
	dc.LineTo(410, 300)
	dc.LineTo(10, 300)
	dc.ClosePath()
	dc.Fill()

	if err := gg.SavePNG("./testdata/_gradient-linear.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
