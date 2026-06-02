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

func TestSheetBeginPath(t *testing.T) {
	width := 220
	height := 140

	canvas := tinyskia.NewCanvas(width, height)
	ctx := canvas.GetContext()

	// First path (blue)
	ctx.BeginPath()
	ctx.SetStrokeStyleSolidColor(color.RGBA{0, 0, 255, 255}) // blue
	ctx.MoveTo(20, 20)
	ctx.LineTo(200, 20)
	ctx.Stroke()

	// Second path (green)
	ctx.BeginPath()
	ctx.SetStrokeStyleSolidColor(color.RGBA{0, 128, 0, 255}) // green
	ctx.MoveTo(20, 20)
	ctx.LineTo(120, 120)
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
