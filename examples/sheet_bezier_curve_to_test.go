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

func TestSheetBezierCurveTo(t *testing.T) {
	width := 300
	height := 150

	canvas := tinyskia.NewCanvas(width, height)
	ctx := canvas.GetContext()

	// Define point coordinates
	start := struct{ x, y float64 }{50, 20}
	cp1 := struct{ x, y float64 }{230, 30}
	cp2 := struct{ x, y float64 }{150, 80}
	end := struct{ x, y float64 }{250, 100}

	// Cubic Bezier curve
	ctx.BeginPath()
	ctx.SetStrokeStyleSolidColor(color.RGBA{0, 0, 0, 255}) // black
	ctx.MoveTo(start.x, start.y)
	ctx.BezierCurveTo(cp1.x, cp1.y, cp2.x, cp2.y, end.x, end.y)
	ctx.Stroke()

	// Start and end points (blue)
	ctx.SetFillStyleSolidColor(color.RGBA{0, 0, 255, 255}) // blue
	ctx.BeginPath()
	ctx.Arc(start.x, start.y, 5, 0, 2*math.Pi) // start point
	ctx.Arc(end.x, end.y, 5, 0, 2*math.Pi)     // end point
	ctx.Fill()

	// Control points (red)
	ctx.SetFillStyleSolidColor(color.RGBA{255, 0, 0, 255}) // red
	ctx.BeginPath()
	ctx.Arc(cp1.x, cp1.y, 5, 0, 2*math.Pi) // control point 1
	ctx.Arc(cp2.x, cp2.y, 5, 0, 2*math.Pi) // control point 2
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
