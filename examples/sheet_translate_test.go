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

func TestSheetTranslate(t *testing.T) {
	width := 200
	height := 150

	canvas := tinyskia.NewCanvas(width, height)
	ctx := canvas.GetContext()

	// Moved square
	ctx.Translate(110, 30)
	ctx.SetFillStyleSolidColor(color.RGBA{255, 0, 0, 255}) // red
	ctx.FillRect(0, 0, 80, 80)

	// Reset current transformation matrix to the identity matrix
	ctx.SetTransform(1, 0, 0, 1, 0, 0)

	// Unmoved square
	ctx.SetFillStyleSolidColor(color.RGBA{128, 128, 128, 255}) // gray
	ctx.FillRect(0, 0, 80, 80)

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
