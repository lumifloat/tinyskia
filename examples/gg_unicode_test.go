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

func TestGGUnicode(t *testing.T) {
	const S = 4096 * 2
	const T = 16 * 2
	const F = 28
	c := gg.NewCanvas(S, S)
	dc := c.GetContext()

	switch runtime.GOOS {
	case "windows":
		gg.RegisterFont("C:/Windows/Fonts/arial.ttf", gg.FontFace{Family: "arial"})
	case "darwin":
		gg.RegisterFont("/Library/Fonts/Arial.ttf", gg.FontFace{Family: "arial"})
	default:
		// pass
	}

	dc.SetFillStyleSolidColor(color.RGBA{255, 255, 255, 255})
	dc.FillRect(0, 0, float64(S), float64(S))
	dc.SetFillStyleSolidColor(color.Black)
	dc.SetFont(gg.FontAttr{Family: []string{"arial"}, Size: F})

	for r := 0; r < 256; r++ {
		for c := 0; c < 256; c++ {
			i := r*256 + c
			x := float64(c*T) + T/2
			y := float64(r*T) + T/2

			char := string(rune(i))
			metrics := dc.MeasureText(char)
			charWidth := metrics.Width

			// Center the character in the cell
			drawX := x - charWidth/2
			dc.FillText(char, drawX, y)
		}
	}

	fi, err := os.Create("gg_out.png")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	c.WritePNG(fi, nil)
}
