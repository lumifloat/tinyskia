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

func TestGGMask(t *testing.T) {
	im, err := gg.LoadImage("gg_baboon.png")
	if err != nil {
		log.Fatal(err)
	}

	c := gg.NewCanvas(512, 512)
	dc := c.GetContext()
	dc.RoundRect(0, 0, 512, 512, []float64{64})
	dc.Clip()
	dc.DrawImage(im, 0, 0)
	fi, err := os.Create("gg_out.png")
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	c.WritePNG(fi, nil)
}
