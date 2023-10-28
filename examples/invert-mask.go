// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/arugaz/gg"
)

func main() {
	const S = 1024

	dc := gg.NewContext(S, S)
	dc.DrawCircle(S/2, S/2, 384)
	dc.Clip()
	dc.InvertMask()
	dc.DrawRectangle(0, 0, S, S)
	dc.SetRGB(.3, .2, .1)
	dc.Fill()

	if err := gg.SavePNG("./testdata/_invert-mask.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
