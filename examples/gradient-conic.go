// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"image/color"
	"log"

	"github.com/arugaz/gg"
)

func main() {
	const S = 1000

	dc := gg.NewContext(S, S)

	grad1 := gg.NewConicGradient(S/2, S/2, 0)
	grad1.AddColorStop(0, color.White)
	grad1.AddColorStop(.5, color.RGBA{R: 255, G: 215, A: 255})
	grad1.AddColorStop(1, color.RGBA{R: 255, A: 255})

	dc.SetStrokeStyle(grad1)
	dc.SetLineWidth(20)
	dc.DrawCircle(S/2, S/2, S/2-20)
	dc.Stroke()

	grad2 := gg.NewConicGradient(S/2, S/2, 90)
	grad2.AddColorStop(0, color.RGBA{R: 255, A: 255})
	grad2.AddColorStop(.16, color.RGBA{R: 255, G: 255, A: 255})
	grad2.AddColorStop(.33, color.RGBA{G: 255, A: 255})
	grad2.AddColorStop(.50, color.RGBA{G: 255, B: 255, A: 255})
	grad2.AddColorStop(.66, color.RGBA{B: 255, A: 255})
	grad2.AddColorStop(.83, color.RGBA{R: 255, B: 255, A: 255})
	grad2.AddColorStop(1, color.RGBA{R: 255, A: 255})

	dc.SetFillStyle(grad2)
	dc.DrawCircle(S/2, S/2, S/2-50)
	dc.Fill()

	if err := gg.SavePNG("./testdata/_gradient-conic.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
