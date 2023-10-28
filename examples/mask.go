// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/arugaz/gg"
)

func main() {
	const S = 512

	dc := gg.NewContext(S, S)

	im, err := gg.LoadPNG("./testdata/baboon.png")
	if err != nil {
		log.Fatalf("could not load baboon.png: %+v", err)
	}

	dc.DrawRoundedRectangle(0, 0, S, S, 96)
	dc.Clip()
	dc.DrawImage(im, 0, 0)

	if err := gg.SavePNG("./testdata/_mask.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
