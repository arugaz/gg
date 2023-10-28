// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/arugaz/gg"
)

func main() {
	im1, err := gg.LoadPNG("./testdata/baboon.png")
	if err != nil {
		log.Fatalf("could not load baboon.png: %+v", err)
	}

	im2, err := gg.LoadPNG("./testdata/gopher.png")
	if err != nil {
		log.Fatalf("could not load gopher.png: %+v", err)
	}

	s1 := im1.Bounds().Size()
	s2 := im2.Bounds().Size()

	w := max(s1.X, s2.X)
	h := s1.Y + s2.Y

	dc := gg.NewContext(w, h)
	dc.DrawImage(im1, 0, 0)
	dc.DrawImage(im2, 0, s1.Y)

	if err := gg.SavePNG("./testdata/_concat.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
