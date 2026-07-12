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

func TestSheetRotate(t *testing.T) {
	width := 200
	height := 150

	canvas := tinyskia.NewCanvas(width, height)
	ctx := canvas.GetContext()

	// Point of transform origin
	ctx.SetFillStyleSolidColor(color.RGBA{0, 0, 255, 255}) // blue
	ctx.BeginPath()
	ctx.Arc(0, 0, 5, 0, 2*math.Pi)
	ctx.Fill()

	// Non-rotated rectangle
	ctx.SetFillStyleSolidColor(color.RGBA{128, 128, 128, 255}) // gray
	ctx.FillRect(100, 0, 80, 20)

	// Rotated rectangle
	ctx.Rotate((45 * math.Pi) / 180)
	ctx.SetFillStyleSolidColor(color.RGBA{255, 0, 0, 255}) // red
	ctx.FillRect(100, 0, 80, 20)

	// Reset transformation matrix to the identity matrix
	ctx.SetTransformValues(1, 0, 0, 1, 0, 0)

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
