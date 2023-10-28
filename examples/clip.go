// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/arugaz/gg"
)

func main() {
	const S = 1000

	dc := gg.NewContext(S, S)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	dc.DrawCircle(350, 500, 300)
	dc.Clip()
	dc.DrawCircle(650, 500, 300)
	dc.Clip()
	dc.DrawRectangle(0, 0, S, S)
	dc.SetRGB(.1, .2, .3)
	dc.Fill()

	if err := gg.SavePNG("./testdata/_clip.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
