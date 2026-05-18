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

func TestGGGradientRadial(t *testing.T) {
	c := gg.NewCanvas(200, 200)
	dc := c.GetContext()

	grad := dc.CreateRadialGradient(100, 100, 10, 100, 120, 80)
	grad.AddColorStop(0, color.RGBA{0, 255, 0, 255})
	grad.AddColorStop(1, color.RGBA{0, 0, 255, 255})

	dc.SetFillStyle(grad)
	dc.Rect(0, 0, 200, 200)
	dc.Fill()

	dc.BeginPath()
	dc.SetStrokeStyleSolidColor(color.White)
	dc.Arc(100, 100, 10, 0, 2*math.Pi, false)
	dc.ClosePath()
	dc.Stroke()

	dc.BeginPath()
	dc.Arc(100, 120, 80, 0, 2*math.Pi, false)
	dc.ClosePath()
	dc.Stroke()

	fi, err := os.Create("gg_out.png")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	c.WritePNG(fi, nil)
}
