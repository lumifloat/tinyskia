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

func TestSheetTextAlign(t *testing.T) {
	width := 350
	height := 150

	canvas := tinyskia.NewCanvas(width, height)
	ctx := canvas.GetContext()

	x := float64(width / 2)

	ctx.BeginPath()
	ctx.MoveTo(x, 0)
	ctx.LineTo(x, float64(height))
	ctx.Stroke()

	ctx.SetFont(tinyskia.FontAttr{Family: []string{"arial"}, Size: 30})

	ctx.SetTextAlign(tinyskia.CanvasTextAlignLeft)
	ctx.FillText("left-aligned", x, 40)

	ctx.SetTextAlign(tinyskia.CanvasTextAlignCenter)
	ctx.FillText("center-aligned", x, 85)

	ctx.SetTextAlign(tinyskia.CanvasTextAlignRight)
	ctx.FillText("right-aligned", x, 130)

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
