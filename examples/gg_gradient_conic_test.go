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

func TestGGGradientConic(t *testing.T) {
	dc := gg.NewContext(400, 400)

	grad1 := dc.CreateConicGradient(0, 200, 200)
	grad1.AddColorStop(0.0, color.Black)
	grad1.AddColorStop(0.5, color.RGBA{255, 215, 0, 255})
	grad1.AddColorStop(1.0, color.RGBA{255, 0, 0, 255})

	grad2 := dc.CreateConicGradient(math.Pi/2, 200, 200)
	grad2.AddColorStop(0.00, color.RGBA{255, 0, 0, 255})
	grad2.AddColorStop(0.16, color.RGBA{255, 255, 0, 255})
	grad2.AddColorStop(0.33, color.RGBA{0, 255, 0, 255})
	grad2.AddColorStop(0.50, color.RGBA{0, 255, 255, 255})
	grad2.AddColorStop(0.66, color.RGBA{0, 0, 255, 255})
	grad2.AddColorStop(0.83, color.RGBA{255, 0, 255, 255})
	grad2.AddColorStop(1.00, color.RGBA{255, 0, 0, 255})

	dc.SetStrokeStyle(grad1)
	dc.SetLineWidth(20)
	dc.Arc(200, 200, 180, 0, 2*math.Pi, false)
	dc.Stroke()

	dc.BeginPath()
	dc.SetFillStyle(grad2)
	dc.Arc(200, 200, 150, 0, 2*math.Pi, false)
	dc.Fill()

	dc.SavePNG("gg_out.png")
}
