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
	"fmt"
	"image/color"
	"math"
	"os"
	"testing"

	"github.com/lumifloat/tinyskia"
)

func TestSheetMiterLimit(t *testing.T) {
	width := 350
	height := 650

	canvas := tinyskia.NewCanvas(width, height)
	ctx := canvas.GetContext()

	// Clear canvas
	ctx.ClearRect(0, 0, float64(width), float64(height))

	// Draw reference lines (blue)
	ctx.SetStrokeStyleSolidColor(color.RGBA{9, 153, 255, 255}) // #09f
	ctx.SetLineWidth(2)
	ctx.StrokeRect(-5, 50, 160, 50)

	// Set line style for zigzag
	ctx.SetStrokeStyleSolidColor(color.RGBA{0, 0, 0, 255}) // black
	ctx.SetLineWidth(10)
	ctx.SetLineCap(tinyskia.LineCapButt)
	ctx.SetLineJoin(tinyskia.LineJoinMiter)

	// Test different miterLimit values
	miterLimits := []float64{2.0, 4.2, 10.0, 20.0}
	yPositions := []float64{100, 250, 400, 550}

	for i, miterLimit := range miterLimits {
		ctx.SetMiterLimit(miterLimit)

		// Draw zigzag pattern
		ctx.BeginPath()
		ctx.MoveTo(0, yPositions[i])
		for j := 0; j < 24; j++ {
			dy := 25.0
			if j%2 != 0 {
				dy = -25.0
			}
			x := math.Pow(float64(j), 1.5) * 2
			y := yPositions[i] - 25 + dy
			ctx.LineTo(x, y)
		}
		ctx.Stroke()

		// Draw label
		ctx.SetFillStyleSolidColor(color.RGBA{0, 0, 0, 255})
		ctx.SetFont(tinyskia.FontAttr{Family: []string{"arial"}, Size: 12})
		ctx.FillText(fmt.Sprintf("MiterLimit=%0.1f", miterLimit), 250, yPositions[i]+5)
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
