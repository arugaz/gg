// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/arugaz/gg"
	"golang.org/x/image/font/gofont/gobold"
)

func main() {
	const S = 1024

	dc := gg.NewContext(S, S)

	if err := dc.LoadFontFaceFromBytes(gobold.TTF, 48); err != nil {
		log.Fatalf("could not load bold font: %+v", err)
	}

	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	dc.DrawStringAnchored("Hello, world!", S/2, S/2, .5, .5)

	if err := gg.SavePNG("./testdata/_gofont.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
