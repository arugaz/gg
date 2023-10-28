// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"math"

	"github.com/arugaz/gg"
)

func main() {
	const (
		W = 1500
		H = 75
	)

	dc := gg.NewContext(W, H)

	dc.ScaleAbout(.95, .75, W/2, H/2)

	for i := 0; i < W; i++ {
		a := float64(i) * 2 * math.Pi / W * 8
		x := float64(i)
		y := (math.Sin(a) + 1) / 2 * H
		dc.LineTo(x, y)
	}

	dc.ClosePath()
	dc.SetHexColor("#3E606F")
	dc.FillPreserve()
	dc.SetHexColor("#19344180")
	dc.SetLineWidth(8)
	dc.Stroke()

	if err := gg.SavePNG("./testdata/_sine.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
