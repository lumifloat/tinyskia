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

func TestGGInvertMask(t *testing.T) {
	const S = 1024
	dc := gg.NewContext(S, S)
	dc.SetFillStyleSolidColor(color.RGBA{0, 0, 0, 255})
	dc.FillRect(0, 0, float64(S), float64(S))

	dc.Arc(512, 512, 384, 0, math.Pi*2, false)
	dc.Clip()

	dc.SetGlobalCompositeOperation(gg.CompositeOperationDestinationOut)
	dc.Fill()

	dc.SavePNG("gg_out.png")
}
