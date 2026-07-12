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

func TestGGCircle(t *testing.T) {
	c := gg.NewCanvas(1000, 1000)
	dc := c.GetContext()

	dc.BeginPath()
	dc.Arc(500, 500, 400, 0, 2*math.Pi)
	dc.SetFillStyleSolidColor(color.RGBA{0, 0, 0, 255})
	dc.Fill()

	fi, err := os.Create("gg_out.png")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	c.WritePNG(fi, nil)
}
