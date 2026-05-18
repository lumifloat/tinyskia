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

func TestGGSpiral(t *testing.T) {
	const S = 1024
	const N = 2048
	c := gg.NewCanvas(S, S)
	dc := c.GetContext()
	dc.SetFillStyleSolidColor(color.White)
	dc.FillRect(0, 0, float64(S), float64(S))
	dc.SetFillStyleSolidColor(color.Black)

	for i := 0; i <= N; i++ {
		t := float64(i) / N
		d := t*S*0.4 + 10
		a := t * math.Pi * 2 * 20
		x := S/2 + math.Cos(a)*d
		y := S/2 + math.Sin(a)*d
		r := t * 8
		dc.BeginPath()
		dc.Arc(x, y, r, 0, 2*math.Pi, false)
		dc.Fill()
	}

	fi, err := os.Create("gg_out.png")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	c.WritePNG(fi, nil)
}
