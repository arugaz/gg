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
	const S = 400
	dc := gg.NewContext(S, S)

	dc.SetRGB(0, 0, 0)

	font, err := gg.FontParse(goregular.TTF)
	if err != nil {
		panic(err)
	}

	face, err := gg.FontNewFace(font, 40)
	if err != nil {
		panic(err)
	}
	defer face.Close()

	dc.SetFontFace(face)
	text := "World, hello!"
	w, h := dc.MeasureString(text)
	dc.Rotate(gg.Radians(10))
	dc.DrawRectangle(100, 180, w, h)
	dc.Stroke()
	dc.DrawStringAnchored(text, 100, 180, 0, 1)

	if err := gg.SavePNG("./testdata/_rotated-text.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
