// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"math/rand"

	"github.com/arugaz/gg"
)

func main() {
	rnd := rand.New(rand.NewSource(54321))

	random := func() float64 {
		return rnd.Float64()*2 - 1
	}
	point := func() (x, y float64) {
		return random(), random()
	}
	drawCurve := func(dc *gg.Context) {
		dc.SetRGBA(0, 0, 0, .1)
		dc.FillPreserve()
		dc.SetRGB(0, 0, 0)
		dc.SetLineWidth(12)
		dc.Stroke()
	}
	drawPoints := func(dc *gg.Context) {
		dc.SetRGBA(1, 0, 0, .5)
		dc.SetLineWidth(2)
		dc.Stroke()
	}
	randomQuadratic := func(dc *gg.Context) {
		x0, y0 := point()
		x1, y1 := point()
		x2, y2 := point()
		dc.MoveTo(x0, y0)
		dc.QuadraticTo(x1, y1, x2, y2)
		drawCurve(dc)
		dc.MoveTo(x0, y0)
		dc.LineTo(x1, y1)
		dc.LineTo(x2, y2)
		drawPoints(dc)
	}
	randomCubic := func(dc *gg.Context) {
		x0, y0 := point()
		x1, y1 := point()
		x2, y2 := point()
		x3, y3 := point()
		dc.MoveTo(x0, y0)
		dc.CubicTo(x1, y1, x2, y2, x3, y3)
		drawCurve(dc)
		dc.MoveTo(x0, y0)
		dc.LineTo(x1, y1)
		dc.LineTo(x2, y2)
		dc.LineTo(x3, y3)
		drawPoints(dc)
	}

	const (
		S = 256
		W = 8
		H = 8
	)

	dc := gg.NewContext(S*W, S*H)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	for j := 0; j < H; j++ {
		for i := 0; i < W; i++ {
			x := float64(i)*S + S/2
			y := float64(j)*S + S/2
			dc.Push()
			dc.Translate(x, y)
			dc.Scale(S/2, S/2)
			if j%2 == 0 {
				randomCubic(dc)
			} else {
				randomQuadratic(dc)
			}
			dc.Pop()
			dc.AppendFrame()
		}
	}

	if err := gg.SavePNG("./testdata/_beziers.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
	if err := gg.SaveJPG("./testdata/_beziers.jpeg", dc.Image(), 80); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
	if err := gg.SaveGIF("./testdata/_beziers.gif", dc.Frames(), 30); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
