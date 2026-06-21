// Copyright 2016 Michael Fogleman
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package examples

import (
	"image/color"
	"os"
	"runtime"
	"testing"

	gg "github.com/lumifloat/tinyskia"
)

func TestGGRotatedText(t *testing.T) {
	const S = 400
	c := gg.NewCanvas(S, S)
	dc := c.GetContext()

	switch runtime.GOOS {
	case "windows":
		gg.RegisterFontWithFile("C:/Windows/Fonts/arial.ttf", gg.FontFace{Family: "arial"})
	case "darwin":
		gg.RegisterFontWithFile("/Library/Fonts/Arial.ttf", gg.FontFace{Family: "arial"})
	default:
		// pass
	}

	dc.SetFillStyleSolidColor(color.RGBA{255, 255, 255, 255})
	dc.FillRect(0, 0, float64(S), float64(S))
	dc.SetFillStyleSolidColor(color.Black)

	dc.SetFont(gg.FontAttr{Family: []string{"arial"}, Size: 40})
	text := "Hello, world!"
	metrics := dc.MeasureText(text)
	textWidth := metrics.Width
	textHeight := metrics.FontBoundingBoxAscent - metrics.FontBoundingBoxDescent

	x := 100.0
	y := 180.0 + metrics.FontBoundingBoxAscent // Adjust for baseline

	dc.Save()
	dc.Rotate(gg.Radians(10))

	dc.SetStrokeStyleSolidColor(color.Black)
	dc.SetLineWidth(1)
	dc.BeginPath()
	dc.MoveTo(x, y-textHeight)
	dc.LineTo(x+textWidth, y-textHeight)
	dc.LineTo(x+textWidth, y)
	dc.LineTo(x, y)
	dc.ClosePath()
	dc.Stroke()

	dc.FillText(text, x, y)
	dc.Restore()

	fi, err := os.Create("gg_out.png")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	c.WritePNG(fi, nil)
}
