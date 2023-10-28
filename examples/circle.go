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

	dc.DrawCircle(S/2, S/2, 400)
	dc.SetRGB(.3, .2, .1)
	dc.Fill()

	if err := gg.SavePNG("./testdata/_circle.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
