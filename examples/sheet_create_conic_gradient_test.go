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

func TestSheetCreateConicGradient(t *testing.T) {
	width := 240
	height := 240

	canvas := tinyskia.NewCanvas(width, height)
	ctx := canvas.GetContext()

	// Create a conic gradient
	// Start angle is 0
	// Center position is 100, 100
	gradient := ctx.CreateConicGradient(0, 100, 100)

	// Add five color stops
	gradient.AddColorStop(0, color.RGBA{255, 0, 0, 255})      // red
	gradient.AddColorStop(0.25, color.RGBA{255, 165, 0, 255}) // orange
	gradient.AddColorStop(0.5, color.RGBA{255, 255, 0, 255})  // yellow
	gradient.AddColorStop(0.75, color.RGBA{0, 128, 0, 255})   // green
	gradient.AddColorStop(1, color.RGBA{0, 0, 255, 255})      // blue

	// Set fill style and draw rectangle
	ctx.SetFillStyleGradient(gradient)
	ctx.FillRect(20, 20, 200, 200)

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
