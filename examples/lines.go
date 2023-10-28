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
	const S = 1024

	rnd := rand.New(rand.NewSource(54321))
	dc := gg.NewContext(S, S)

	dc.SetRGB(1, 1, 1)
	dc.Clear()

	for i := 0; i < S; i++ {
		x1 := rnd.Float64() * S
		y1 := rnd.Float64() * S
		x2 := rnd.Float64() * S
		y2 := rnd.Float64() * S
		r := rnd.Float64()
		g := rnd.Float64()
		b := rnd.Float64()
		a := rnd.Float64()*.5 + .5
		w := rnd.Float64()*4 + 1
		dc.SetRGBA(r, g, b, a)
		dc.SetLineWidth(w)
		dc.DrawLine(x1, y1, x2, y2)
		dc.Stroke()
	}

	if err := gg.SavePNG("./testdata/_lines.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
