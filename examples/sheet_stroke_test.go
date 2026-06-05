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

func TestSheetStroke(t *testing.T) {
	width := 300
	height := 180

	canvas := tinyskia.NewCanvas(width, height)
	ctx := canvas.GetContext()

	// First sub-path
	ctx.SetLineWidth(26)
	ctx.SetStrokeStyleSolidColor(color.RGBA{255, 165, 0, 255}) // orange
	ctx.MoveTo(20, 20)
	ctx.LineTo(160, 20)
	ctx.Stroke()

	// Second sub-path
	ctx.SetLineWidth(14)
	ctx.SetStrokeStyleSolidColor(color.RGBA{0, 128, 0, 255}) // green
	ctx.MoveTo(20, 80)
	ctx.LineTo(220, 80)
	ctx.Stroke()

	// Third sub-path
	ctx.SetLineWidth(4)
	ctx.SetStrokeStyleSolidColor(color.RGBA{255, 192, 203, 255}) // pink
	ctx.MoveTo(20, 140)
	ctx.LineTo(280, 140)
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
