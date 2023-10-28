// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"net/http"

	"github.com/arugaz/gg"
)

func main() {
	const (
		S = 4096 * 2
		T = 16 * 2
	)

	dc := gg.NewContext(S, S)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)

	resp, err := http.Get("https://www.1001fonts.com/download/font/xolonium.regular.otf")
	if err != nil {
		log.Fatalf("could not download font: %+v", err)
	}
	defer resp.Body.Close()

	if err := dc.LoadFontFaceFromReader(resp.Body, 28); err != nil {
		log.Fatalf("could not load xolonium-regular font: %+v", err)
	}

	for r := 0; r < 256; r++ {
		for c := 0; c < 256; c++ {
			i := r*256 + c
			x := float64(c*T) + T/2
			y := float64(r*T) + T/2
			dc.DrawStringAnchored(string(rune(i)), x, y, .5, .5)
		}
	}

	if err := gg.SavePNG("./testdata/_unicode.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
