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

func TestSheetArcTo(t *testing.T) {
	width := 250
	height := 150

	canvas := tinyskia.NewCanvas(width, height)
	ctx := canvas.GetContext()

	// Draw tangent lines (gray)
	ctx.BeginPath()
	ctx.SetStrokeStyleSolidColor(color.RGBA{128, 128, 128, 255}) // gray
	ctx.SetLineWidth(1)
	ctx.MoveTo(200, 20)
	ctx.LineTo(200, 130)
	ctx.LineTo(50, 20)
	ctx.Stroke()

	// Draw arc (black)
	ctx.BeginPath()
	ctx.SetStrokeStyleSolidColor(color.RGBA{0, 0, 0, 255}) // black
	ctx.SetLineWidth(5)
	ctx.MoveTo(200, 20)
	ctx.ArcTo(200, 130, 50, 20, 40)
	ctx.Stroke()

	// Draw start point (blue)
	ctx.BeginPath()
	ctx.SetFillStyleSolidColor(color.RGBA{0, 0, 255, 255}) // blue
	ctx.Arc(200, 20, 5, 0, 2*math.Pi, false)
	ctx.Fill()

	// Draw control points (red)
	ctx.BeginPath()
	ctx.SetFillStyleSolidColor(color.RGBA{255, 0, 0, 255}) // red
	ctx.Arc(200, 130, 5, 0, 2*math.Pi, false)              // control point 1
	ctx.Arc(50, 20, 5, 0, 2*math.Pi, false)                // control point 2
	ctx.Fill()

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
