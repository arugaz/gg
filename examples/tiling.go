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
		NX = 5
		NY = 4
	)

	im, err := gg.LoadPNG("./testdata/gopher.png")
	if err != nil {
		log.Fatalf("could not load gopher.png: %+v", err)
	}
	w := im.Bounds().Dx()
	h := im.Bounds().Dy()

	dc := gg.NewContext(w*NX, h*NY)

	for y := 0; y < NY; y++ {
		for x := 0; x < NX; x++ {
			dc.DrawImage(im, x*w, y*h)
		}
	}

	if err := gg.SavePNG("./testdata/_tiling.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
