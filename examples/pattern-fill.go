// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/arugaz/gg"
)

func main() {
	const S = 600

	dc := gg.NewContext(S, S)

	im, err := gg.LoadPNG("./testdata/baboon.png")
	if err != nil {
		log.Fatalf("could not load baboon.png: %+v", err)
	}

	pattern := gg.NewSurfacePattern(im, gg.RepeatBoth)

	dc.MoveTo(0, 0)
	dc.LineTo(S, 0)
	dc.LineTo(S, S)
	dc.LineTo(0, S)
	dc.ClosePath()
	dc.SetFillStyle(pattern)
	dc.Fill()

	if err := gg.SavePNG("./testdata/_pattern-fill.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
