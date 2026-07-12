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

func TestGGInvertMask(t *testing.T) {
	const S = 1024
	c := gg.NewCanvas(S, S)
	dc := c.GetContext()
	dc.SetFillStyleSolidColor(color.RGBA{0, 0, 0, 255})
	dc.FillRect(0, 0, float64(S), float64(S))

	dc.Arc(512, 512, 384, 0, math.Pi*2)
	dc.Clip()

	dc.SetGlobalCompositeOperation(gg.CompositeOperationDestinationOut)
	dc.Fill()

	fi, err := os.Create("gg_out.png")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	c.WritePNG(fi, nil)
}
