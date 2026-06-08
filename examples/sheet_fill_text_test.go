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

func TestSheetFillText(t *testing.T) {
	width := 400
	height := 150

	canvas := tinyskia.NewCanvas(width, height)
	ctx := canvas.GetContext()

	// Set font and draw text
	ctx.SetFont(tinyskia.FontAttr{
		Family: []string{"serif"},
		Size:   50,
	})
	ctx.SetFillStyleSolidColor(color.RGBA{0, 0, 0, 255}) // black
	ctx.FillText("Hello world", 50, 90)

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
