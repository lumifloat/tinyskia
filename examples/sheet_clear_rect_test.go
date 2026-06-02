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

func TestSheetClearRect(t *testing.T) {
	width := 200
	height := 150

	canvas := tinyskia.NewCanvas(width, height)
	ctx := canvas.GetContext()

	// Draw yellow background
	ctx.SetFillStyleSolidColor(color.RGBA{255, 255, 102, 255}) // #ff6
	ctx.FillRect(0, 0, float64(width), float64(height))

	// Draw blue triangle
	ctx.BeginPath()
	ctx.SetFillStyleSolidColor(color.RGBA{0, 0, 255, 255}) // blue
	ctx.MoveTo(20, 20)
	ctx.LineTo(180, 20)
	ctx.LineTo(130, 130)
	ctx.ClosePath()
	ctx.Fill()

	// Clear part of the canvas
	ctx.ClearRect(10, 10, 120, 100)

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
