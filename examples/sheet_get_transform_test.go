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

func TestSheetGetTransform(t *testing.T) {
	width := 300
	height := 150

	canvas := tinyskia.NewCanvas(width, height)
	ctx := canvas.GetContext()

	// Set transform on first context
	ctx.SetTransformValues(1, 0.2, 0.8, 1, 0, 0)

	// Draw a rectangle with the transform
	ctx.SetFillStyleSolidColor(color.RGBA{255, 0, 0, 255}) // red
	ctx.FillRect(25, 25, 50, 50)

	// Get the stored transform and verify with assertions
	storedTransform := ctx.GetTransform()

	// Assert transform values with tolerance for float32->float64 conversion
	expectedA := 1.0
	expectedB := 0.2
	expectedC := 0.8
	expectedD := 1.0
	expectedE := 0.0
	expectedF := 0.0
	tolerance := 1e-5

	if math.Abs(storedTransform.A()-expectedA) > tolerance {
		t.Fatalf("Expected transform.A = %f, got %f", expectedA, storedTransform.A())
	}
	if math.Abs(storedTransform.B()-expectedB) > tolerance {
		t.Fatalf("Expected transform.B = %f, got %f", expectedB, storedTransform.B())
	}
	if math.Abs(storedTransform.C()-expectedC) > tolerance {
		t.Fatalf("Expected transform.C = %f, got %f", expectedC, storedTransform.C())
	}
	if math.Abs(storedTransform.D()-expectedD) > tolerance {
		t.Fatalf("Expected transform.D = %f, got %f", expectedD, storedTransform.D())
	}
	if math.Abs(storedTransform.E()-expectedE) > tolerance {
		t.Fatalf("Expected transform.E = %f, got %f", expectedE, storedTransform.E())
	}
	if math.Abs(storedTransform.F()-expectedF) > tolerance {
		t.Fatalf("Expected transform.F = %f, got %f", expectedF, storedTransform.F())
	}

	// Reset transform and apply the stored transform
	ctx.SetTransformValues(1, 0, 0, 1, 0, 0) // identity
	ctx.SetTransformValues(storedTransform.A(), storedTransform.B(), storedTransform.C(), storedTransform.D(), storedTransform.E(), storedTransform.F())

	// Draw an arc with the same transform
	ctx.SetFillStyleSolidColor(color.RGBA{0, 0, 255, 255}) // blue
	ctx.BeginPath()
	ctx.Arc(50, 50, 50, 0, 2*math.Pi)
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
