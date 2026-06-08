// Copyright 2016 Michael Fogleman
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package examples

import (
	"image/color"
	"os"
	"testing"

	gg "github.com/lumifloat/tinyskia"
	"golang.org/x/image/font/gofont/goregular"
)

func TestGoFont(t *testing.T) {
	const S = 1024
	c := gg.NewCanvas(S, S)
	dc := c.GetContext()
	dc.SetFillStyleSolidColor(color.RGBA{255, 255, 255, 255})
	dc.FillRect(0, 0, float64(S), float64(S))

	if err := gg.RegisterFontWithData(goregular.TTF, gg.FontFace{Family: "go"}); err != nil {
		t.Errorf("Font not available: %v", err)
		return
	}
	dc.SetFont(gg.FontAttr{Family: []string{"go"}, Size: 48})
	text := "Hello, world!"
	metrics := dc.MeasureText(text)
	textWidth := metrics.Width
	textHeight := metrics.FontBoundingBoxAscent - metrics.FontBoundingBoxDescent

	x := S/2 - textWidth/2
	y := S/2 + textHeight/2

	dc.SetFillStyleSolidColor(color.Black)
	dc.FillText(text, x, y)

	fi, err := os.Create("gg_out.png")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	c.WritePNG(fi, nil)
}
