// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/arugaz/gg"
	"golang.org/x/image/font/gofont/goregular"
)

func main() {
	const S = 1000

	dc := gg.NewContext(S, S)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	dc.Translate(S/2, S/2)
	dc.Scale(40, 40)

	var x0, y0, x1, y1, x2, y2, x3, y3, x4, y4 float64
	x0, y0 = -10, 0
	x1, y1 = -5, -10
	x2, y2 = 0, 0
	x3, y3 = 5, 10
	x4, y4 = 10, 0

	dc.MoveTo(x0, y0)
	dc.LineTo(x1, y1)
	dc.LineTo(x2, y2)
	dc.LineTo(x3, y3)
	dc.LineTo(x4, y4)
	dc.SetHexColor("#FF2D00")
	dc.SetLineWidth(8)
	dc.Stroke()

	dc.MoveTo(x0, y0)
	dc.QuadraticTo(x1, y1, x2, y2)
	dc.QuadraticTo(x3, y3, x4, y4)
	dc.SetHexColor("#3E606F")
	dc.SetLineWidth(16)
	dc.FillPreserve()
	dc.SetRGB(0, 0, 0)
	dc.Stroke()

	dc.DrawCircle(x0, y0, .5)
	dc.DrawCircle(x1, y1, .5)
	dc.DrawCircle(x2, y2, .5)
	dc.DrawCircle(x3, y3, .5)
	dc.DrawCircle(x4, y4, .5)
	dc.SetRGB(1, 1, 1)
	dc.FillPreserve()
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(4)
	dc.Stroke()

	dc.Scale(.05, .05)
	if err := dc.LoadFontFaceFromBytes(goregular.TTF, 100); err != nil {
		log.Fatalf("could not load regular font: %+v", err)
	}
	dc.DrawStringAnchored("G", +100, -100, .5, .5)
	dc.DrawStringAnchored("g", -100, +100, .5, .5)

	if err := gg.SavePNG("./testdata/_quadratic.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
