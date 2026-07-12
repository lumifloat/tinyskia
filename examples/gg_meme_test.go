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
)

func TestGGMeme(t *testing.T) {
	const S = 1024
	c := gg.NewCanvas(S, S)
	dc := c.GetContext()
	dc.SetFillStyleSolidColor(color.White)
	dc.FillRect(0, 0, float64(S), float64(S))

	s := "ONE DOES NOT SIMPLY"
	dc.SetFont(gg.FontAttr{Family: []string{"impact"}, Size: 96})
	metrics := dc.MeasureText(s)
	textWidth := metrics.Width
	textHeight := metrics.FontBoundingBoxAscent - metrics.FontBoundingBoxDescent

	x := S/2 - textWidth/2
	y := S/2 + textHeight/2

	dc.SetStrokeStyleSolidColor(color.Black)
	dc.SetLineWidth(12)
	dc.SetLineJoin(gg.CanvasLineJoinRound)
	dc.StrokeText(s, x, y)

	dc.SetFillStyleSolidColor(color.White)
	dc.FillText(s, x, y)

	fi, err := os.Create("gg_out.png")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	c.WritePNG(fi, nil)
}
