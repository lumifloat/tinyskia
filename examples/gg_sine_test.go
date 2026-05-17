// Copyright 2016 Michael Fogleman
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package examples

import (
	"image/color"
	"math"
	"os"
	"testing"

	gg "github.com/lumifloat/tinyskia"
)

func TestGGSine(t *testing.T) {
	const W = 1200
	const H = 60
	c := gg.NewCanvas(W, H)
	dc := c.GetContext()

	dc.Translate(-W/2, -H/2)
	dc.Scale(0.95, 0.75)
	dc.Translate(W/2, H/2)

	for i := 0; i < W; i++ {
		a := float64(i) * 2 * math.Pi / W * 8
		x := float64(i)
		y := (math.Sin(a) + 1) / 2 * H
		dc.LineTo(x, y)
	}
	dc.ClosePath()
	dc.SetFillStyleSolidColor(color.RGBA{62, 96, 111, 255}) // #3E606F
	dc.Fill()

	dc.SetStrokeStyleSolidColor(color.RGBA{25, 52, 65, 128}) // #19344180
	dc.SetLineWidth(8)
	dc.Stroke()

	fi, err := os.Create("gg_out.png")
	if err != nil {
		panic(err)
	}
	c.WritePNG(fi, nil)
}
