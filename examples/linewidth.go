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
	w := .1

	dc.SetRGB(1, 1, 1)
	dc.Clear()

	dc.SetRGB(0, 0, 0)

	for i := 100; i <= S-100; i += 20 {
		x := float64(i)
		dc.DrawLine(x+50, 0, x-50, S)
		dc.SetLineWidth(w)
		dc.Stroke()
		w += .1
	}

	if err := gg.SavePNG("./testdata/_linewidth.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
