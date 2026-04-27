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

func random() float64 {
	return rand.Float64()*2 - 1
}

func point(size float64) (x, y float64) {
	return random() * size, random() * size
}

func drawCurve(dc *gg.Context) {
	dc.SetFillStyleSolidColor(color.RGBA{0, 0, 0, 26})
	dc.Fill()
	dc.SetStrokeStyleSolidColor(color.RGBA{0, 0, 0, 255})
	dc.SetLineWidth(12)
	dc.Stroke()
}

func drawPoints(dc *gg.Context) {
	dc.SetStrokeStyleSolidColor(color.RGBA{255, 0, 0, 128})
	dc.SetLineWidth(2)
	dc.Stroke()
}

func randomQuadratic(dc *gg.Context, size float64) {
	x0, y0 := point(size)
	x1, y1 := point(size)
	x2, y2 := point(size)
	dc.BeginPath()
	dc.MoveTo(x0, y0)
	dc.QuadraticCurveTo(x1, y1, x2, y2)
	drawCurve(dc)
	dc.BeginPath()
	dc.MoveTo(x0, y0)
	dc.LineTo(x1, y1)
	dc.LineTo(x2, y2)
	drawPoints(dc)
}

func randomCubic(dc *gg.Context, size float64) {
	x0, y0 := point(size)
	x1, y1 := point(size)
	x2, y2 := point(size)
	x3, y3 := point(size)
	dc.BeginPath()
	dc.MoveTo(x0, y0)
	dc.BezierCurveTo(x1, y1, x2, y2, x3, y3)
	drawCurve(dc)
	dc.BeginPath()
	dc.MoveTo(x0, y0)
	dc.LineTo(x1, y1)
	dc.LineTo(x2, y2)
	dc.LineTo(x3, y3)
	drawPoints(dc)
}

func TestGGBezier(t *testing.T) {
	const (
		S = 256
		W = 8
		H = 8
	)
	dc := gg.NewContext(S*W, S*H)
	dc.SetFillStyleSolidColor(color.RGBA{255, 255, 255, 255})
	dc.FillRect(0, 0, float64(S*W), float64(S*H))

	for j := 0; j < H; j++ {
		for i := 0; i < W; i++ {
			x := float64(i)*S + S/2
			y := float64(j)*S + S/2
			dc.Save()
			dc.Translate(x, y)
			if j%2 == 0 {
				randomCubic(dc, S/2)
			} else {
				randomQuadratic(dc, S/2)
			}
			dc.Restore()
		}
	}
	dc.SavePNG("out.png")
}
