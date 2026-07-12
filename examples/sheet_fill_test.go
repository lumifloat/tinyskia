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

func TestSheetFill(t *testing.T) {
	width := 300
	height := 150

	canvas := tinyskia.NewCanvas(width, height)
	ctx := canvas.GetContext()

	// Create path
	region := tinyskia.NewPath2D()
	region.MoveTo(30, 90)
	region.LineTo(110, 20)
	region.LineTo(240, 130)
	region.LineTo(60, 130)
	region.LineTo(190, 20)
	region.LineTo(270, 90)
	region.ClosePath()

	// Fill path with evenodd rule
	ctx.SetFillStyleSolidColor(color.RGBA{0, 128, 0, 255}) // green
	ctx.FillPathWithFillRule(region, tinyskia.CanvasFillRuleEvenodd)

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
