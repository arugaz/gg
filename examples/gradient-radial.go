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
		W = 400
		H = 200
	)

	dc := gg.NewContext(W, H)

	grad := gg.NewRadialGradient(100, 100, 10, 100, 120, 80)
	grad.AddColorStop(0, color.RGBA{G: 255, A: 255})
	grad.AddColorStop(1, color.RGBA{B: 255, A: 255})

	dc.SetFillStyle(grad)
	dc.DrawRectangle(0, 0, 200, 200)
	dc.Fill()

	dc.SetColor(color.White)
	dc.DrawCircle(100, 100, 10)
	dc.Stroke()
	dc.DrawCircle(100, 120, 80)
	dc.Stroke()

	if err := gg.SavePNG("./testdata/_gradient-radial.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
