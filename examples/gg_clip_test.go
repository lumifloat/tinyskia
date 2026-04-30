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

func TestGGClip(t *testing.T) {
	dc := gg.NewContext(1000, 1000)

	dc.BeginPath()
	dc.Arc(350, 500, 300, 0, 2*math.Pi, false)
	dc.Clip()

	dc.BeginPath()
	dc.Arc(650, 500, 300, 0, 2*math.Pi, false)
	dc.Clip()

	dc.BeginPath()
	dc.Rect(0, 0, 1000, 1000)
	dc.SetFillStyleSolidColor(color.RGBA{0, 0, 0, 255})
	dc.Fill()

	dc.SavePNG("gg_out.png")
}
