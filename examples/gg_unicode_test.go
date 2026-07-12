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

func TestGGUnicode(t *testing.T) {
	const S = 4096 * 2
	const T = 16 * 2
	const F = 28
	c := gg.NewCanvas(S, S)
	dc := c.GetContext()

	dc.SetFillStyleSolidColor(color.RGBA{255, 255, 255, 255})
	dc.FillRect(0, 0, float64(S), float64(S))
	dc.SetFillStyleSolidColor(color.Black)
	dc.SetFont(gg.FontAttr{Family: []string{gg.FontGenericSansSerif}, Size: F})

	for r := 0; r < 256; r++ {
		for c := 0; c < 256; c++ {
			i := r*256 + c
			x := float64(c*T) + T/2
			y := float64(r*T) + T/2

			char := string(rune(i))
			dc.SetTextAlign(gg.CanvasTextAlignCenter)
			dc.FillText(char, x, y)
		}
	}

	fi, err := os.Create("gg_out.png")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	c.WritePNG(fi, nil)
}
