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

func TestSheetClosePath(t *testing.T) {
	width := 300
	height := 150

	canvas := tinyskia.NewCanvas(width, height)
	ctx := canvas.GetContext()

	ctx.BeginPath()
	ctx.Arc(240, 20, 40, 0, math.Pi, false)
	ctx.MoveTo(100, 20)
	ctx.Arc(60, 20, 40, 0, math.Pi, false)
	ctx.MoveTo(215, 80)
	ctx.Arc(150, 80, 65, 0, math.Pi, false)
	ctx.ClosePath()
	ctx.SetLineWidth(6)
	ctx.SetStrokeStyleSolidColor(color.RGBA{0, 0, 0, 255}) // black
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
