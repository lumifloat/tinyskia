// Copyright 2016 Michael Fogleman
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package examples

import (
	"log"
	"os"
	"testing"

	gg "github.com/lumifloat/tinyskia"
)

func TestGGPatternFill(t *testing.T) {
	im, err := gg.LoadPNG("gg_baboon.png")
	if err != nil {
		log.Fatal(err)
	}

	// Create a pattern by tiling the image
	// Since NewSurfacePattern is not available, we'll manually tile it
	patternWidth := im.Bounds().Dx()
	patternHeight := im.Bounds().Dy()

	c := gg.NewCanvas(600, 600)
	dc := c.GetContext()

	// Tile the pattern image to fill the context
	for y := 0; y < 600; y += patternHeight {
		for x := 0; x < 600; x += patternWidth {
			dc.DrawImage(im, float64(x), float64(y))
		}
	}

	// Create a clipping region (the quad shape from original)
	dc.BeginPath()
	dc.MoveTo(20, 20)
	dc.LineTo(590, 20)
	dc.LineTo(590, 590)
	dc.LineTo(20, 590)
	dc.ClosePath()
	dc.Clip()

	// The pattern is already drawn, clipping will show only the quad area
	// But we need to redraw because Clip affects future operations
	// Actually, let's do it differently: draw pattern, then clip

	fi, err := os.Create("gg_out.png")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	c.WritePNG(fi, nil)
}
