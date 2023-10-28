// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/arugaz/gg"
)

func main() {
	const (
		S     = 1000
		Minor = 10
		Major = 100
	)

	dc := gg.NewContext(S, S)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	for x := Minor; x < S; x += Minor {
		fx := float64(x) + .5
		dc.DrawLine(fx, 0, fx, S)
	}
	for y := Minor; y < S; y += Minor {
		fy := float64(y) + .5
		dc.DrawLine(0, fy, S, fy)
	}

	dc.SetLineWidth(1)
	dc.SetRGBA(0, 0, 0, .25)
	dc.Stroke()

	for x := Major; x < S; x += Major {
		fx := float64(x) + .5
		dc.DrawLine(fx, 0, fx, S)
	}
	for y := Major; y < S; y += Major {
		fy := float64(y) + .5
		dc.DrawLine(0, fy, S, fy)
	}

	dc.SetLineWidth(1)
	dc.SetRGBA(0, 0, 0, .5)
	dc.Stroke()

	if err := gg.SavePNG("./testdata/_crisp.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
