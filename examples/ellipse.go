// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"embed"
	"log"

	"github.com/arugaz/gg"
)

//go:embed testdata/gopher.png
var img embed.FS

func main() {
	const S = 1000

	dc := gg.NewContext(S, S)
	dc.SetRGBA(0, 0, 0, .1)

	for i := 0; i < 360; i += 15 {
		dc.Push()
		dc.RotateAbout(gg.Radians(float64(i)), S/2, S/2)
		dc.DrawEllipse(S/2, S/2, S*7/16, S/8)
		dc.Fill()
		dc.Pop()
	}

	im, err := gg.LoadImageFromFS(img, "testdata/gopher.png")
	if err != nil {
		panic(err)
	}
	dc.DrawImageAnchored(im, S/2, S/2, .5, .5)

	if err := gg.SavePNG("./testdata/_ellipse.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
