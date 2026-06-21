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
	"runtime"
	"testing"

	"github.com/lumifloat/tinyskia"
)

func TestSheetFontKerning(t *testing.T) {
	width := 400
	height := 150

	canvas := tinyskia.NewCanvas(width, height)
	ctx := canvas.GetContext()

	switch runtime.GOOS {
	case "windows":
		tinyskia.RegisterFontWithFile("C:/Windows/Fonts/arial.ttf", tinyskia.FontFace{Family: "arial"})
	case "darwin":
		tinyskia.RegisterFontWithFile("/Library/Fonts/Arial.ttf", tinyskia.FontFace{Family: "arial"})
	default:
		// pass
	}

	ctx.SetFont(tinyskia.FontAttr{Family: []string{"arial"}, Size: 30})

	// Default (auto)
	ctx.FillText("AVA Ta We (default)", 5, 30)

	// Font kerning: normal
	ctx.SetFontKerning(tinyskia.FontKerningNormal)
	ctx.FillText("AVA Ta We (normal)", 5, 70)

	// Font kerning: none
	ctx.SetFontKerning(tinyskia.FontKerningNone)
	ctx.FillText("AVA Ta We (none)", 5, 110)

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
