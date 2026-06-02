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

func TestSheetLineJoin(t *testing.T) {
	width := 150
	height := 150

	lineJoins := []struct {
		name string
		join tinyskia.LineJoin
	}{
		{"round", tinyskia.LineJoinRound},
		{"bevel", tinyskia.LineJoinBevel},
		{"miter", tinyskia.LineJoinMiter},
	}

	canvas := tinyskia.NewCanvas(width, height)
	ctx := canvas.GetContext()

	ctx.SetLineWidth(10)
	ctx.SetStrokeStyleSolidColor(color.RGBA{0, 0, 0, 255}) // black

	for i, lj := range lineJoins {
		ctx.SetLineJoin(lj.join)
		ctx.BeginPath()
		ctx.MoveTo(-5, float64(5+i*40))
		ctx.LineTo(35, float64(45+i*40))
		ctx.LineTo(75, float64(5+i*40))
		ctx.LineTo(115, float64(45+i*40))
		ctx.LineTo(155, float64(5+i*40))
		ctx.Stroke()
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
