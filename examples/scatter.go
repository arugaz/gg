// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"math/rand"

	"github.com/arugaz/gg"
	"golang.org/x/image/font/gofont/gobold"
)

func main() {
	createPoints := func(n int) []gg.Point {
		var (
			rnd = rand.New(rand.NewSource(54321))
			pts = make([]gg.Point, n)
		)

		for i := 0; i < n; i++ {
			x := 0.5 + rnd.NormFloat64()*0.1
			y := x + rnd.NormFloat64()*0.1
			pts[i] = gg.Point{X: x, Y: y}
		}

		return pts
	}

	const (
		S = 1024
		P = 64
	)

	dc := gg.NewContext(S, S)

	dc.InvertY()
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.Translate(P, P)
	dc.Scale(S-P*2, S-P*2)

	for i := 1; i <= 10; i++ {
		x := float64(i) / 10
		dc.MoveTo(x, 0)
		dc.LineTo(x, 1)
		dc.MoveTo(0, x)
		dc.LineTo(1, x)
	}
	dc.SetRGBA(0, 0, 0, .25)
	dc.SetLineWidth(1)
	dc.Stroke()

	dc.MoveTo(0, 0)
	dc.LineTo(1, 0)
	dc.MoveTo(0, 0)
	dc.LineTo(0, 1)
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(4)
	dc.Stroke()

	points := createPoints(1000)

	dc.SetRGBA(0, 0, 1, .5)
	for _, p := range points {
		dc.DrawCircle(p.X, p.Y, 3.1/S)
		dc.Fill()
	}

	dc.Identity()
	dc.SetRGB(0, 0, 0)

	err := dc.LoadFontFaceFromBytes(gobold.TTF, 32)
	if err != nil {
		log.Fatalf("could not load bold font: %+v", err)
	}
	dc.DrawStringAnchored("Chart Title", S/2, P/2, .5, .5)

	err = dc.LoadFontFaceFromBytes(gobold.TTF, 24)
	if err != nil {
		log.Fatalf("could not load bold font: %+v", err)
	}
	dc.DrawStringAnchored("X Axis Title", S/2, S-P/2, .5, .5)
	dc.Rotate(gg.Radians(-90))
	dc.DrawStringAnchored("Y Axis Title", -S/2, S/P+2, .5, .5)

	if err := gg.SavePNG("./testdata/_scatter.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
