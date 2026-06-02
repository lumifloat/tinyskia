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
	"os"
	"testing"

	"github.com/lumifloat/tinyskia"
)

func TestSheetLineDashOffset(t *testing.T) {
	width := 300
	height := 150

	canvas := tinyskia.NewCanvas(width, height)
	ctx := canvas.GetContext()

	// Set dash pattern: 4 units on, 16 units off
	ctx.SetLineDash([]float64{4, 16})

	// Draw line without offset (black)
	ctx.SetStrokeStyleSolidColor(color.RGBA{0, 0, 0, 255}) // black
	ctx.SetLineWidth(2)
	ctx.SetLineDashOffset(0)
	ctx.BeginPath()
	ctx.MoveTo(0, 50)
	ctx.LineTo(300, 50)
	ctx.Stroke()

	// Draw line with offset = 4 (red)
	ctx.SetStrokeStyleSolidColor(color.RGBA{255, 0, 0, 255}) // red
	ctx.SetLineDashOffset(4)
	ctx.BeginPath()
	ctx.MoveTo(0, 100)
	ctx.LineTo(300, 100)
	ctx.Stroke()

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
