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
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)

	if err := dc.LoadFontFace("/usr/share/fonts/noto-cjk/NotoSansCJK-Regular.ttc", 96); err != nil {
		log.Fatalf("could not load font: %+v", err)
	}
	dc.DrawStringAnchored("こんにちは世界！", S/2, S/2, 0.5, 0.5)

	if err := gg.SavePNG("./testdata/_text.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
