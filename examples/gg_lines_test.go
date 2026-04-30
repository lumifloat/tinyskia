// Copyright 2016 Michael Fogleman
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package examples

import (
	"image/color"
	"math/rand"
	"testing"

	gg "github.com/lumifloat/tinyskia"
)

func TestGGLines(t *testing.T) {
	const W = 1024
	const H = 1024
	dc := gg.NewContext(W, H)
	dc.SetFillStyleSolidColor(color.Black)
	dc.FillRect(0, 0, float64(W), float64(H))
	source := rand.NewSource(1)
	rnd := rand.New(source)
	for i := 0; i < 1000; i++ {
		x1 := rnd.Float64() * W
		y1 := rnd.Float64() * H
		x2 := rnd.Float64() * W
		y2 := rnd.Float64() * H
		r := rnd.Float64()
		g := rnd.Float64()
		b := rnd.Float64()
		a := rnd.Float64()*0.5 + 0.5
		w := rnd.Float64()*4 + 1
		dc.SetStrokeStyleSolidColor(color.RGBA{
			uint8(r * 255),
			uint8(g * 255),
			uint8(b * 255),
			uint8(a * 255),
		})
		dc.SetLineWidth(w)
		dc.BeginPath()
		dc.MoveTo(x1, y1)
		dc.LineTo(x2, y2)
		dc.Stroke()
	}
	dc.SavePNG("gg_out.png")
}
