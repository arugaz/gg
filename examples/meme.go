// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/arugaz/gg"
	"golang.org/x/image/font/gofont/goregular"
)

func main() {
	const S = 1024

	dc := gg.NewContext(S, S)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	if err := dc.LoadFontFaceFromBytes(goregular.TTF, 72); err != nil {
		log.Fatalf("could not load regular font: %+v", err)
	}

	dc.SetRGB(0, 0, 0)
	s := "ONE DOES NOT SIMPLY"
	n := 5

	for dy := -n; dy <= n; dy++ {
		for dx := -n; dx <= n; dx++ {
			if dx*dx+dy*dy >= n*n {
				continue
			}
			x := S/2 + float64(dx)
			y := S/2 + float64(dy)
			dc.DrawStringAnchored(s, x, y, .5, .5)
		}
	}

	dc.SetRGB(1, 1, 1)
	dc.DrawStringAnchored(s, S/2, S/2, 0.5, 0.5)

	if err := gg.SavePNG("./testdata/_meme.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
