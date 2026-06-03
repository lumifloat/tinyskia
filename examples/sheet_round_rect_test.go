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

func TestSheetRoundRect(t *testing.T) {
	width := 700
	height := 300

	canvas := tinyskia.NewCanvas(width, height)
	ctx := canvas.GetContext()

	// Rounded rectangle with zero radius (specified as a number)
	ctx.SetStrokeStyleSolidColor(color.RGBA{255, 0, 0, 255}) // red
	ctx.BeginPath()
	ctx.RoundRect(10, 20, 150, 100, []float64{0})
	ctx.Stroke()

	// Rounded rectangle with 40px radius (single element list)
	ctx.SetStrokeStyleSolidColor(color.RGBA{0, 0, 255, 255}) // blue
	ctx.BeginPath()
	ctx.RoundRect(10, 20, 150, 100, []float64{40})
	ctx.Stroke()

	// Rounded rectangle with 2 different radii
	ctx.SetStrokeStyleSolidColor(color.RGBA{255, 165, 0, 255}) // orange
	ctx.BeginPath()
	ctx.RoundRect(10, 150, 150, 100, []float64{10, 40})
	ctx.Stroke()

	// Rounded rectangle with four different radii
	ctx.SetStrokeStyleSolidColor(color.RGBA{0, 128, 0, 255}) // green
	ctx.BeginPath()
	ctx.RoundRect(400, 20, 200, 100, []float64{0, 30, 50, 60})
	ctx.Stroke()

	// Same rectangle drawn backwards
	ctx.SetStrokeStyleSolidColor(color.RGBA{255, 0, 255, 255}) // magenta
	ctx.BeginPath()
	ctx.RoundRect(400, 150, -200, 100, []float64{0, 30, 50, 60})
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
