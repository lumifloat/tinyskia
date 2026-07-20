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

func TestSheetCreatePattern(t *testing.T) {
	width := 240
	height := 240

	// Create a pattern canvas (offscreen)
	patternCanvas := tinyskia.NewCanvas(50, 50)
	patternCtx := patternCanvas.GetContext()

	// Give the pattern a background color
	patternCtx.SetFillStyleSolidColor(color.RGBA{255, 238, 204, 255}) // #ffeecc
	patternCtx.FillRect(0, 0, 50, 50)

	// Draw an arc
	patternCtx.SetStrokeStyleSolidColor(color.RGBA{0, 0, 0, 255}) // black
	patternCtx.BeginPath()
	patternCtx.Arc(0, 0, 50, 0, 0.5*math.Pi)
	patternCtx.Stroke()

	// Create our primary canvas and fill it with the pattern
	canvas := tinyskia.NewCanvas(width, height)
	ctx := canvas.GetContext()

	// Create pattern from the offscreen canvas
	pattern, err := ctx.CreatePattern(patternCanvas.Image(), tinyskia.RepeatModeNoRepeat)
	if err != nil {
		t.Fatalf("Failed to create pattern: %v", err)
	}

	ctx.SetFillStylePattern(pattern)
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
