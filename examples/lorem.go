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
	lines := []string{
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod",
		"tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam,",
		"quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo",
		"consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse",
		"cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat",
		"non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
	}

	const (
		W = 800
		H = 400
	)

	dc := gg.NewContext(W, H)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)

	if err := dc.LoadFontFaceFromBytes(gobold.TTF, 14); err != nil {
		log.Fatalf("could not load bold font: %+v", err)
	}

	const h = 24
	for i, line := range lines {
		y := H/2 - h*len(lines)/2 + i*h
		dc.DrawStringAnchored(line, 400, float64(y), .5, .5)
	}

	if err := gg.SavePNG("./testdata/_lorem.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
