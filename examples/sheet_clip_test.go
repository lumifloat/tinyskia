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

func TestSheetClip(t *testing.T) {
	width := 180
	height := 150

	canvas := tinyskia.NewCanvas(width, height)
	ctx := canvas.GetContext()

	// Create clipping path
	region := tinyskia.NewPath2D()
	region.Rect(80, 10, 20, 130)
	region.Rect(40, 50, 100, 50)
	ctx.ClipPathWithFillRule(region, tinyskia.FillRuleEvenOdd)

	// Draw clipped content
	ctx.SetFillStyleSolidColor(color.RGBA{0, 0, 255, 255}) // blue
	ctx.FillRect(0, 0, float64(width), float64(height))

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
