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

func TestGGLorem(t *testing.T) {
	var lines = []string{
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod",
		"tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam,",
		"quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo",
		"consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse",
		"cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat",
		"non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
	}

	const W = 800
	const H = 400
	c := gg.NewCanvas(W, H)
	dc := c.GetContext()
	dc.SetFillStyleSolidColor(color.White)
	dc.FillRect(0, 0, float64(W), float64(H))
	dc.SetFillStyleSolidColor(color.Black)

	const h = 24
	for i, line := range lines {
		metrics := dc.MeasureText(line)
		x := W/2 - metrics.Width/2
		y := H/2 - h*len(lines)/2 + i*h
		dc.FillText(line, float64(x), float64(y))
	}
	fi, err := os.Create("gg_out.png")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	c.WritePNG(fi, nil)
}
