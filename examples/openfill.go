// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"math"
	"math/rand"

	"github.com/arugaz/gg"
)

func main() {
	const S = 1000
	var (
		rnd = rand.New(rand.NewSource(54321))
		dc  = gg.NewContext(S, S)
	)

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			x := float64(i)*100 + 50
			y := float64(j)*100 + 50
			a1 := rnd.Float64() * 2 * math.Pi
			a2 := a1 + rnd.Float64()*math.Pi + math.Pi/2
			dc.DrawArc(x, y, 40, a1, a2)
		}
	}
	dc.SetRGB(0, 0, 0)
	dc.FillPreserve()
	dc.SetRGB(1, 1, 1)
	dc.SetLineWidth(8)
	dc.StrokePreserve()
	dc.SetRGB(1, 0, 0)
	dc.SetLineWidth(4)
	dc.StrokePreserve()

	if err := gg.SavePNG("./testdata/_openfill.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
