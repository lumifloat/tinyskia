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

func TestSheetGetLineDash(t *testing.T) {
	width := 300
	height := 100

	canvas := tinyskia.NewCanvas(width, height)
	ctx := canvas.GetContext()

	// Set line dash pattern
	dashPattern := []float64{10, 20}
	ctx.SetLineDash(dashPattern)

	// Get line dash and verify with assertion
	retrievedDash := ctx.GetLineDash()
	if len(retrievedDash) != len(dashPattern) {
		t.Fatalf("Expected line dash length %d, got %d", len(dashPattern), len(retrievedDash))
	}
	for i, v := range dashPattern {
		if retrievedDash[i] != v {
			t.Fatalf("Expected line dash[%d] = %f, got %f", i, v, retrievedDash[i])
		}
	}

	// Draw a dashed line
	ctx.SetStrokeStyleSolidColor(color.RGBA{0, 0, 0, 255}) // black
	ctx.BeginPath()
	ctx.MoveTo(0, 50)
	ctx.LineTo(300, 50)
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
