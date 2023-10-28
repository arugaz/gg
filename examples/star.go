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
	makePoints := func(n int, x, y, r float64) []gg.Point {
		result := make([]gg.Point, n)
		for i := 0; i < n; i++ {
			a := float64(i)*2*math.Pi/float64(n) - math.Pi/2
			result[i] = gg.Point{X: x + r*math.Cos(a), Y: y + r*math.Sin(a)}
		}
		return result
	}

	const (
		S = 512
		N = 5
	)

	var (
		points = makePoints(N, S, S, 400)
		dc     = gg.NewContext(S*2, S*2)
	)

	dc.SetHexColor("fff")
	dc.Clear()

	for i := 0; i < N+1; i++ {
		index := (i * 2) % N
		p := points[index]
		dc.LineTo(p.X, p.Y)
	}

	dc.SetRGBA(.5, 0, 0, 1)
	dc.SetFillRule(gg.FillRuleEvenOdd)
	dc.FillPreserve()
	dc.SetRGBA(1, 0, 0, .5)
	dc.SetLineWidth(16)
	dc.Stroke()

	if err := gg.SavePNG("./testdata/_star.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
