// Copyright 2016 Michael Fogleman
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package examples

import (
	"image/color"
	"math"
	"testing"

	gg "github.com/lumifloat/tinyskia"
)

func Polygon(n int, x, y, r float64) [][2]float64 {
	result := make([][2]float64, n)
	for i := 0; i < n; i++ {
		a := float64(i)*2*math.Pi/float64(n) - math.Pi/2
		result[i] = [2]float64{x + r*math.Cos(a), y + r*math.Sin(a)}
	}
	return result
}

func TestGGStar(t *testing.T) {
	n := 5
	points := Polygon(n, 512, 512, 400)

	dc := gg.NewContext(1024, 1024)
	dc.SetFillStyleSolidColor(color.RGBA{255, 255, 255, 255})
	dc.FillRect(0, 0, 1024, 1024)

	for i := 0; i < n+1; i++ {
		index := (i * 2) % n
		p := points[index]
		dc.LineTo(p[0], p[1])
	}

	dc.SetFillStyleSolidColor(color.RGBA{0, 128, 0, 255})
	dc.FillWithFillRule(gg.FillRuleEvenOdd)

	dc.SetStrokeStyleSolidColor(color.RGBA{0, 255, 0, 128})
	dc.SetLineWidth(16)
	dc.Stroke()

	dc.SavePNG("gg_out.png")
}
