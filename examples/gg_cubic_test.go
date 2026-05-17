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

func TestGGCubic(t *testing.T) {
	const S = 1000
	c := gg.NewCanvas(S, S)
	dc := c.GetContext()
	dc.SetFillStyleSolidColor(color.RGBA{255, 255, 255, 255})
	dc.FillRect(0, 0, float64(S), float64(S))

	dc.Translate(float64(S/2), float64(S/2))

	var x0, y0, x1, y1, x2, y2, x3, y3 float64
	x0, y0 = -400, 0
	x1, y1 = -320, -320
	x2, y2 = 320, 320
	x3, y3 = 400, 0

	dc.BeginPath()
	dc.MoveTo(x0, y0)
	dc.BezierCurveTo(x1, y1, x2, y2, x3, y3)
	dc.SetFillStyleSolidColor(color.RGBA{0, 0, 0, 51})
	dc.Fill()

	dc.SetStrokeStyleSolidColor(color.RGBA{0, 0, 0, 255})
	dc.SetLineWidth(8)
	dc.SetLineDash([]float64{16, 24})
	dc.Stroke()

	dc.BeginPath()
	dc.MoveTo(x0, y0)
	dc.LineTo(x1, y1)
	dc.LineTo(x2, y2)
	dc.LineTo(x3, y3)
	dc.SetStrokeStyleSolidColor(color.RGBA{255, 0, 0, 102}) // 0.4 alpha
	dc.SetLineWidth(2)
	dc.SetLineDash([]float64{4, 8, 1, 8})
	dc.Stroke()

	fi, err := os.Create("gg_out.png")
	if err != nil {
		panic(err)
	}
	c.WritePNG(fi, nil)
}
