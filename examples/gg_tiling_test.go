// Copyright 2016 Michael Fogleman
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package examples

import (
	"os"
	"testing"

	gg "github.com/lumifloat/tinyskia"
)

func TestGGTiling(t *testing.T) {
	const NX = 4
	const NY = 3
	im, err := gg.LoadPNG("gg_gopher.png")
	if err != nil {
		t.Skipf("Image not available: %v", err)
		return
	}
	w := im.Bounds().Size().X
	h := im.Bounds().Size().Y
	c := gg.NewCanvas(w*NX, h*NY)
	dc := c.GetContext()
	for y := 0; y < NY; y++ {
		for x := 0; x < NX; x++ {
			dc.DrawImage(im, float64(x*w), float64(y*h))
		}
	}
	fi, err := os.Create("gg_out.png")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	c.WritePNG(fi, nil)
}
