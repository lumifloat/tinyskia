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
	"image/png"
	"os"
	"testing"

	"github.com/lumifloat/tinyskia"
)

func TestSheetGetImageData(t *testing.T) {
	width := 700
	height := 350

	canvas := tinyskia.NewCanvas(width, height)
	ctx := canvas.GetContext()

	// Load the gopher image
	imageFile, err := os.Open("gg_gopher.png")
	if err != nil {
		t.Fatalf("Failed to open image file: %v", err)
	}
	defer imageFile.Close()

	img, err := png.Decode(imageFile)
	if err != nil {
		t.Fatalf("Failed to decode image: %v", err)
	}

	// Draw the original image
	ctx.DrawImageScaled(img, 0, 0, 233, 320)

	// Get image data from a region
	imageData, err := ctx.GetImageData(10, 20, 80, 230)
	if err != nil {
		t.Fatalf("Failed to get image data: %v", err)
	}

	// Put the image data at different positions
	ctx.PutImageData(imageData, 260, 0)
	ctx.PutImageData(imageData, 380, 50)
	ctx.PutImageData(imageData, 500, 100)

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
