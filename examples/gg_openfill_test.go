// Copyright 2016 Michael Fogleman
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package examples

import (
	"image/color"
	"math"
	"math/rand"
	"os"
	"testing"

	gg "github.com/lumifloat/tinyskia"
)

func TestGGOpenFill(t *testing.T) {
	c := gg.NewCanvas(1000, 1000)
	dc := c.GetContext()

	for j := 0; j < 10; j++ {
		for i := 0; i < 10; i++ {
			x := float64(i)*100 + 50
			y := float64(j)*100 + 50
			a1 := rand.Float64() * 2 * math.Pi
			a2 := a1 + rand.Float64()*math.Pi + math.Pi/2
			dc.Arc(x, y, 40, a1, a2, false)
		}
	}

	dc.SetFillStyleSolidColor(color.Black)
	dc.Fill()

	dc.SetStrokeStyleSolidColor(color.White)
	dc.SetLineWidth(8)
	dc.Stroke()

	dc.SetStrokeStyleSolidColor(color.RGBA{255, 0, 0, 255})
	dc.SetLineWidth(4)
	dc.Stroke()

	fi, err := os.Create("gg_out.png")
	if err != nil {
		panic(err)
	}
	c.WritePNG(fi, nil)
}
