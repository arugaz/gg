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
		S = 1024
		N = S * 2
	)

	dc := gg.NewContext(S, S)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)

	for i := 0; i <= N; i++ {
		t := float64(i) / N
		d := t*S*.4 + 10
		a := t * math.Pi * 2 * 20
		x := S/2 + math.Cos(a)*d
		y := S/2 + math.Sin(a)*d
		r := t * 8
		dc.DrawCircle(x, y, r)
	}

	dc.Fill()

	if err := gg.SavePNG("./testdata/_spiral.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
