// Copyright 2016 Michael Fogleman
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package examples

import (
	"log"
	"testing"

	gg "github.com/lumifloat/tinyskia"
)

func TestGGMask(t *testing.T) {
	im, err := gg.LoadImage("gg_baboon.png")
	if err != nil {
		log.Fatal(err)
	}

	dc := gg.NewContext(512, 512)
	dc.RoundRect(0, 0, 512, 512, []float64{64})
	dc.Clip()
	dc.DrawImage(im, 0, 0)
	dc.SavePNG("gg_out.png")
}
