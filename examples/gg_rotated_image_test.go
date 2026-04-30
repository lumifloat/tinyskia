// Copyright 2016 Michael Fogleman
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package examples

import (
	"image/color"
	"testing"

	gg "github.com/lumifloat/tinyskia"
)

func TestGGRotatedImage(t *testing.T) {
	const W = 400
	const H = 500
	im, err := gg.LoadPNG("gg_gopher.png")
	if err != nil {
		t.Skipf("Image not available: %v", err)
		return
	}
	iw := im.Bounds().Dx()
	ih := im.Bounds().Dy()

	dc := gg.NewContext(W, H)

	dc.SetStrokeStyleSolidColor(color.RGBA{255, 0, 0, 255}) // #ff0000
	dc.SetLineWidth(1)
	dc.BeginPath()
	dc.Rect(0, 0, float64(W), float64(H))
	dc.Stroke()

	dc.SetStrokeStyleSolidColor(color.RGBA{0, 0, 255, 255}) // #0000ff
	dc.SetLineWidth(2)
	dc.BeginPath()
	dc.Rect(100, 210, float64(iw), float64(ih))
	dc.Stroke()
	dc.DrawImage(im, 100, 210)

	dc.SetStrokeStyleSolidColor(color.RGBA{0, 0, 255, 255})
	dc.SetLineWidth(2)
	dc.Save()
	dc.Rotate(gg.Radians(10))
	dc.BeginPath()
	dc.Rect(100, 0, float64(iw), float64(ih)/2+20.0)
	dc.Stroke()

	dc.Clip()
	dc.DrawImage(im, 100, 0)
	dc.Restore()

	dc.SavePNG("gg_out.png")
}
