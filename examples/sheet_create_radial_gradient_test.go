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

func TestSheetCreateRadialGradient(t *testing.T) {
	width := 200
	height := 200

	canvas := tinyskia.NewCanvas(width, height)
	ctx := canvas.GetContext()

	// Create a radial gradient
	// Inner circle at x=110, y=90 with radius 30
	// Outer circle at x=100, y=100 with radius 70
	gradient := ctx.CreateRadialGradient(110, 90, 30, 100, 100, 70)

	// Add three color stops
	gradient.AddColorStop(0, color.RGBA{255, 192, 203, 255})   // pink
	gradient.AddColorStop(0.9, color.RGBA{255, 255, 255, 255}) // white
	gradient.AddColorStop(1, color.RGBA{0, 128, 0, 255})       // green

	// Set fill style and draw rectangle
	ctx.SetFillStyleGradient(gradient)
	ctx.FillRect(20, 20, 160, 160)

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
