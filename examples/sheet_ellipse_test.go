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

func TestSheetEllipse(t *testing.T) {
	width := 300
	height := 150

	canvas := tinyskia.NewCanvas(width, height)
	ctx := canvas.GetContext()

	// Draw red ellipse
	ctx.SetFillStyleSolidColor(color.RGBA{255, 0, 0, 255}) // red
	ctx.BeginPath()
	ctx.Ellipse(60, 75, 50, 30, math.Pi*0.25, 0, math.Pi*1.5, false)
	ctx.Fill()

	// Draw blue ellipse
	ctx.SetFillStyleSolidColor(color.RGBA{0, 0, 255, 255}) // blue
	ctx.BeginPath()
	ctx.Ellipse(150, 75, 50, 30, math.Pi*0.25, 0, math.Pi, false)
	ctx.Fill()

	// Draw green ellipse (counterclockwise)
	ctx.SetFillStyleSolidColor(color.RGBA{0, 128, 0, 255}) // green
	ctx.BeginPath()
	ctx.Ellipse(240, 75, 50, 30, math.Pi*0.25, 0, math.Pi, true)
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
