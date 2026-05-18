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

func TestGGGradientLinear(t *testing.T) {
	c := gg.NewCanvas(420, 330)
	dc := c.GetContext()

	grad := dc.CreateLinearGradient(20, 320, 400, 20)
	grad.AddColorStop(0, color.RGBA{0, 255, 0, 255})
	grad.AddColorStop(1, color.RGBA{0, 0, 255, 255})
	grad.AddColorStop(0.5, color.RGBA{255, 0, 0, 255})

	dc.SetStrokeStyleSolidColor(color.White)
	dc.Rect(20, 20, 400-20, 300)
	dc.Stroke()

	dc.BeginPath()
	dc.SetStrokeStyle(grad)
	dc.SetLineWidth(4)
	dc.MoveTo(10, 10)
	dc.LineTo(410, 10)
	dc.LineTo(410, 100)
	dc.LineTo(10, 100)
	dc.ClosePath()
	dc.Stroke()

	dc.BeginPath()
	dc.SetFillStyle(grad)
	dc.MoveTo(10, 120)
	dc.LineTo(410, 120)
	dc.LineTo(410, 300)
	dc.LineTo(10, 300)
	dc.ClosePath()
	dc.Fill()

	fi, err := os.Create("gg_out.png")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	c.WritePNG(fi, nil)
}
