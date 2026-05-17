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

func TestGGLineWidth(t *testing.T) {
	c := gg.NewCanvas(1000, 1000)
	dc := c.GetContext()
	dc.SetFillStyleSolidColor(color.White)
	dc.FillRect(0, 0, 1000, 1000)
	dc.SetStrokeStyleSolidColor(color.Black)
	w := 0.1
	for i := 100; i <= 900; i += 20 {
		x := float64(i)
		dc.BeginPath()
		dc.MoveTo(x+50, 0)
		dc.LineTo(x-50, 1000)
		dc.SetLineWidth(w)
		dc.Stroke()
		w += 0.1
	}
	fi, err := os.Create("gg_out.png")
	if err != nil {
		panic(err)
	}
	c.WritePNG(fi, nil)
}
