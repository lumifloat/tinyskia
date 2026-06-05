// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
//
// -----------------------------------------------------------------------
// Portions of this code are derived from MDN Web Docs (by Mozilla and individual contributors),
// used under the terms of the Creative Commons Attribution-ShareAlike License (CC-BY-SA) v2.5
// or later, or the MIT License where applicable.
// Source: https://developer.mozilla.org/
package examples

import (
	"image/color"
	"math"
	"os"
	"testing"

	"github.com/lumifloat/tinyskia"
)

func TestSheetFillStyle(t *testing.T) {
	width := 150
	height := 150

	canvas := tinyskia.NewCanvas(width, height)
	ctx := canvas.GetContext()

	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			r := uint8(math.Floor(255 - 42.5*float64(i)))
			g := uint8(math.Floor(255 - 42.5*float64(j)))
			b := uint8(0)
			ctx.SetFillStyleSolidColor(color.RGBA{r, g, b, 255})
			ctx.FillRect(float64(j*25), float64(i*25), 25, 25)
		}
	}

	outputPath := "sheet_out.png"
	fi, err := os.Create(outputPath)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer fi.Close()
	if err := canvas.WritePNG(fi, nil); err != nil {
		t.Fatalf("Failed to save PNG: %v", err)
	}
}
