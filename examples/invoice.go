// Copyright 2023 The gg Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/arugaz/gg"
	"golang.org/x/image/font/gofont/goregular"
)

func main() {
	drawString := func(dc *gg.Context, fontSize float64, s string, x float64, y float64, ax float64, ay float64, width float64, align gg.Align) {
		font, _ := gg.FontParse(goregular.TTF)
		face, _ := gg.FontNewFace(font, fontSize)
		defer face.Close()

		dc.SetFontFace(face)
		dc.DrawStringWrapped(s, x, y, ax, ay, width, fontSize, align)
	}

	const (
		W = 800
		H = 841
		M = 20.
	)

	dc := gg.NewContext(W, H)

	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)

	drawString(dc, 28, "INVOICE", M, M, 0, 0, float64(W)-4*M, gg.AlignRight)
	drawString(dc, 16, "Anonymous", M, 4*M, 0, 0, float64(W)-4*M, gg.AlignRight)
	drawString(dc, 16, "Right Here", M, 5*M, 0, 0, float64(W)-4*M, gg.AlignRight)
	drawString(dc, 16, "69", M, 6*M, 0, 0, float64(W)-4*M, gg.AlignRight)
	drawString(dc, 16, "Bali", M, 7*M, 0, 0, float64(W)-4*M, gg.AlignRight)

	dc.SetRGB(0, 0, 0)
	lineY := M + 20 + 10
	dc.DrawLine(M+20, lineY*4, W-(M+20), lineY*4)
	dc.SetLineWidth(1)
	dc.Stroke()

	drawString(dc, 16, "Anonym", M+20, 11*M, 0, 0, float64(W)-4*M, gg.AlignLeft)
	drawString(dc, 16, "Right There", M+20, 12*M, 0, 0, float64(W)-4*M, gg.AlignLeft)
	drawString(dc, 16, "96", M+20, 13*M, 0, 0, float64(W)-4*M, gg.AlignLeft)
	drawString(dc, 16, "Bojong Gede", M+20, 14*M, 0, 0, float64(W)-4*M, gg.AlignLeft)

	drawString(dc, 16, fmt.Sprintf("Number : %s", "1234"), M+20, 11*M, 0, 0, float64(W)-4*M, gg.AlignRight)
	drawString(dc, 16, fmt.Sprintf("Date  :  %s", time.Now().Format("2006-01-02")), (M+20)-87, 12*M, 0, 0, float64(W)-4*M, gg.AlignRight)
	drawString(dc, 16, fmt.Sprintf("Due Date :  %s", time.Now().Add(time.Hour*24).Format("2006-01-02")), (M+20)-87, 13*M, 0, 0, float64(W)-4*M, gg.AlignRight)

	drawString(dc, 16, "Products", M+20, 18.5*M, 0, 0, float64(W)-4*M, gg.AlignLeft)
	drawString(dc, 16, "Quantity", M-210, 18.5*M, 0, 0, float64(W)-4*M, gg.AlignRight)
	drawString(dc, 16, "Price", M-105, 18.5*M, 0, 0, float64(W)-4*M, gg.AlignRight)
	drawString(dc, 16, "Total", M+8, 18.5*M, 0, 0, float64(W)-4*M, gg.AlignRight)

	dc.DrawLine(M+20, lineY*8, W-(M+20), lineY*8)
	dc.SetLineWidth(1)
	dc.Stroke()

	curM := 21.
	curP := 0
	products := map[string]interface{}{
		"Product 1": map[string]interface{}{
			"quantity": 1,
			"price":    10000,
		},
		"Product 2": map[string]interface{}{
			"quantity": 2,
			"price":    12500,
		},
		"Product 3": map[string]interface{}{
			"quantity": 3,
			"price":    15000,
		},
	}

	for productName, product := range products {
		drawString(dc, 16, productName, M+20, curM*M, 0, 0, float64(W)-4*M, gg.AlignLeft)

		quantity := fmt.Sprintf("%d", product.(map[string]interface{})["quantity"].(int))
		qLen := float64(len(quantity))
		qMargin := qLen * 2.7
		if qLen > 3 {
			qMargin = qLen * 3.5
		}
		drawString(dc, 16, quantity, M-(237.5-qMargin), curM*M, 0, 0, float64(W)-4*M, gg.AlignRight)

		price := fmt.Sprintf("%d", product.(map[string]interface{})["price"].(int))
		pLen := float64(len(price))
		pMargin := pLen * 14
		if pLen > 1 {
			pMargin = (qLen * 10) - (qLen * (qLen * 10))
		}
		total := product.(map[string]interface{})["quantity"].(int) * product.(map[string]interface{})["price"].(int)
		drawString(dc, 16, price, M-(107+pMargin), curM*M, 0, 0, float64(W)-4*M, gg.AlignRight)
		drawString(dc, 16, fmt.Sprintf("%d", total), M+20, curM*M, 0, 0, float64(W)-4*M, gg.AlignRight)
		curP += total
		curM += 1
	}

	lineY = curM + 20 + 10 + float64(len(products)) - 2
	dc.DrawLine(M+20, lineY*9.5, W-(M+20), lineY*9.5)
	dc.SetLineWidth(1)
	dc.Stroke()

	drawString(dc, 16, fmt.Sprintf("Sub Total : %d", curP), M+20, (curM+6.5)*M, 0, 0, float64(W)-4*M, gg.AlignRight)
	drawString(dc, 16, fmt.Sprintf("Vat %.2f%% :     %d", .7, int(float64(curP)*(.7/100))), M+20, (curM+7.5)*M, 0, 0, float64(W)-4*M, gg.AlignRight)

	curM += 1.5
	lineY = curM + 20 + 10 + (float64(len(products)) - 2) + 12.5
	dc.DrawLine(M+440, lineY*9.5, W-(M+20), lineY*9.5)
	dc.SetLineWidth(1)
	dc.Stroke()

	drawString(dc, 16, fmt.Sprintf("Total :   %d", int(float64(curP)*(.7/100))+curP), M+20, (curM+7.5)*M, 0, 0, float64(W)-4*M, gg.AlignRight)

	if err := gg.SavePNG("./testdata/_invoice.png", dc.Image()); err != nil {
		log.Fatalf("could not save to file: %+v", err)
	}
}
