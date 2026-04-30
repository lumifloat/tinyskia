// Copyright 2016 Michael Fogleman
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package examples

import (
	"image/color"
	"testing"

	gg "github.com/lumifloat/tinyskia"
)

func TestGGCrisp(t *testing.T) {
	const W = 1000
	const H = 1000
	const Minor = 10
	const Major = 100

	dc := gg.NewContext(W, H)

	// Fill white background
	dc.SetFillStyleSolidColor(color.RGBA{255, 255, 255, 255})
	dc.FillRect(0, 0, float64(W), float64(H))

	// minor grid
	dc.BeginPath()
	for x := Minor; x < W; x += Minor {
		fx := float64(x) + 0.5
		dc.MoveTo(fx, 0)
		dc.LineTo(fx, float64(H))
	}
	for y := Minor; y < H; y += Minor {
		fy := float64(y) + 0.5
		dc.MoveTo(0, fy)
		dc.LineTo(float64(W), fy)
	}
	dc.SetStrokeStyleSolidColor(color.RGBA{0, 0, 0, 64})
	dc.SetLineWidth(1)
	dc.Stroke()

	// major grid
	dc.BeginPath()
	for x := Major; x < W; x += Major {
		fx := float64(x) + 0.5
		dc.MoveTo(fx, 0)
		dc.LineTo(fx, float64(H))
	}
	for y := Major; y < H; y += Major {
		fy := float64(y) + 0.5
		dc.MoveTo(0, fy)
		dc.LineTo(float64(W), fy)
	}
	dc.SetStrokeStyleSolidColor(color.RGBA{0, 0, 0, 128})
	dc.SetLineWidth(1)
	dc.Stroke()

	dc.SavePNG("gg_out.png")
}
