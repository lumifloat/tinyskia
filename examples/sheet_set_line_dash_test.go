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
	"os"
	"testing"

	"github.com/lumifloat/tinyskia"
)

func TestSheetSetLineDash(t *testing.T) {
	width := 300
	height := 180

	canvas := tinyskia.NewCanvas(width, height)
	ctx := canvas.GetContext()

	y := 15.0

	drawDashedLine := func(pattern []float64) {
		ctx.BeginPath()
		ctx.SetLineDash(pattern)
		ctx.MoveTo(0, y)
		ctx.LineTo(300, y)
		ctx.Stroke()
		y += 20
	}

	drawDashedLine([]float64{})
	drawDashedLine([]float64{1, 1})
	drawDashedLine([]float64{10, 10})
	drawDashedLine([]float64{20, 5})
	drawDashedLine([]float64{15, 3, 3, 3})
	drawDashedLine([]float64{20, 3, 3, 3, 3, 3, 3, 3})
	drawDashedLine([]float64{12, 3, 3}) // Equals [12, 3, 3, 12, 3, 3]

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
