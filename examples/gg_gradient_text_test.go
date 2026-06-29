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

func TestGGGradientText(t *testing.T) {
	const W = 1024
	const H = 512

	c := gg.NewCanvas(W, H)
	dc := c.GetContext()

	// Set white background
	dc.SetFillStyleSolidColor(color.RGBA{255, 255, 255, 255})
	dc.FillRect(0, 0, W, H)

	dc.SetFont(gg.FontAttr{Family: []string{"impact"}, Size: 128})
	text := "Gradient Text"

	// Calculate center position
	metrics := dc.MeasureText(text)
	textWidth := metrics.Width
	textHeight := metrics.FontBoundingBoxAscent - metrics.FontBoundingBoxDescent

	x := W/2 - textWidth/2
	y := H/2 + textHeight/2

	// Create gradient
	grad := dc.CreateLinearGradient(0, 0, W, H)
	grad.AddColorStop(0, color.RGBA{255, 0, 0, 255})
	grad.AddColorStop(1, color.RGBA{0, 0, 255, 255})

	// Fill text with gradient
	dc.SetFillStyle(grad)
	dc.FillText(text, x, y)

	fi, err := os.Create("gg_out.png")
	if err != nil {
		t.Fatal(err)
	}
	defer fi.Close()
	c.WritePNG(fi, nil)
}
