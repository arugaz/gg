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
	makePoints := func(n int) []gg.Point {
		result := make([]gg.Point, n)
		for i := 0; i < n; i++ {
			a := float64(i)*2*math.Pi/float64(n) - math.Pi/2
			result[i] = gg.Point{X: math.Cos(a), Y: math.Sin(a)}
		}
		return result
	}

	const (
		W = 1200
		H = W / 10
		S = 100
	)

	dc := gg.NewContext(W, H)
	dc.SetHexColor("#FFFFFF")
	dc.Clear()

	var (
		n      = 5
		points = makePoints(n)
		rnd    = rand.New(rand.NewSource(54321))
	)

	for x := S / 2; x < W; x += S {
		dc.Push()

		s := rnd.Float64()*S/4 + S/4
		dc.Translate(float64(x), H/2)
		dc.Rotate(rnd.Float64() * 2 * math.Pi)
		dc.Scale(s, s)

		for i := 0; i < n+1; i++ {
			index := (i * 2) % n
			p := points[index]
			dc.LineTo(p.X, p.Y)
		}

		dc.SetLineWidth(10)
		dc.SetHexColor("#FC0")
		dc.StrokePreserve()
		dc.SetHexColor("#FFE43A")
		dc.Fill()
		dc.Pop()
	}

	if err := gg.SavePNG("./testdata/_stars.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
