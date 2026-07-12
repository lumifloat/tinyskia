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

func TestSheetLineCap(t *testing.T) {
	width := 150
	height := 150

	lineCaps := []struct {
		name string
		cap  tinyskia.CanvasLineCap
	}{
		{"butt", tinyskia.CanvasLineCapButt},
		{"round", tinyskia.CanvasLineCapRound},
		{"square", tinyskia.CanvasLineCapSquare},
	}

	canvas := tinyskia.NewCanvas(width, height)
	ctx := canvas.GetContext()

	// Draw guides
	ctx.SetStrokeStyleSolidColor(color.RGBA{0, 153, 255, 255}) // #0099ff
	ctx.SetLineWidth(1)
	ctx.BeginPath()
	ctx.MoveTo(10, 10)
	ctx.LineTo(140, 10)
	ctx.MoveTo(10, 140)
	ctx.LineTo(140, 140)
	ctx.Stroke()

	for i, lc := range lineCaps {
		// Draw line with specific line cap
		ctx.SetStrokeStyleSolidColor(color.RGBA{0, 0, 0, 255}) // black
		ctx.SetLineWidth(15)
		ctx.SetLineCap(lc.cap)
		ctx.BeginPath()
		ctx.MoveTo(float64(25+i*50), 10)
		ctx.LineTo(float64(25+i*50), 140)
		ctx.Stroke()

		ctx.Restore()
	}

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
