// Copyright 2016 Michael Fogleman
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package examples

import (
	"testing"

	gg "github.com/lumifloat/tinyskia"
)

func TestGGConcat(t *testing.T) {
	im1, err := gg.LoadPNG("gg_baboon.png")
	if err != nil {
		t.Skipf("Skipping test - baboon.png not found: %v", err)
		return
	}

	im2, err := gg.LoadPNG("gg_gopher.png")
	if err != nil {
		t.Skipf("Skipping test - gopher.png not found: %v", err)
		return
	}

	s1 := im1.Bounds().Size()
	s2 := im2.Bounds().Size()

	width := s1.X
	if s2.X > width {
		width = s2.X
	}
	height := s1.Y + s2.Y

	dc := gg.NewContext(width, height)

	dc.DrawImage(im1, 0, 0)
	dc.DrawImage(im2, 0, float64(s1.Y))

	dc.SavePNG("gg_out.png")
}
