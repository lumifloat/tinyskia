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

func TestGGEllipse(t *testing.T) {
	const S = 1024
	dc := gg.NewContext(S, S)

	dc.SetFillStyleSolidColor(color.RGBA{R: 0, G: 0, B: 0, A: 26})
	for i := 0; i < 360; i += 15 {
		dc.Save()
		dc.Translate(-float64(S/2), -float64(S/2))
		dc.Rotate(gg.Radians(float64(i)))
		dc.Translate(float64(S/2), float64(S/2))
		dc.BeginPath()
		dc.Ellipse(float64(S/2), float64(S/2), float64(S)*7/16, float64(S)/8, 0, 0, 2*math.Pi, false)
		dc.Fill()
		dc.Restore()
	}

	if im, err := gg.LoadImage("gg_gopher.png"); err == nil {
		dc.DrawImage(im, S/2-float64(im.Bounds().Dx()/2), (S/2 - float64(im.Bounds().Dy()/2)))
	}

	dc.SavePNG("gg_out.png")
}
