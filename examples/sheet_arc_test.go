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

func TestSheetArc(t *testing.T) {
	width := 175
	height := 225

	canvas := tinyskia.NewCanvas(width, height)
	ctx := canvas.GetContext()

	ctx.SetStrokeStyleSolidColor(color.RGBA{0, 0, 0, 255}) // black
	ctx.SetFillStyleSolidColor(color.RGBA{0, 0, 0, 255})   // black
	ctx.SetLineWidth(1)

	// Draw shapes
	for i := 0; i <= 3; i++ {
		for j := 0; j <= 2; j++ {
			ctx.BeginPath()
			x := float64(25 + j*50)                        // x coordinate
			y := float64(25 + i*50)                        // y coordinate
			radius := 20.0                                 // arc radius
			startAngle := 0.0                              // starting angle
			endAngle := math.Pi + (math.Pi*float64(j))/2.0 // ending angle
			counterclockwise := i%2 == 1                   // counterclockwise?

			ctx.Arc(x, y, radius, startAngle, endAngle, counterclockwise)

			if i > 1 {
				ctx.Fill() // fill shape
			} else {
				ctx.Stroke() // stroke shape outline
			}
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
