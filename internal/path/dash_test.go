// Copyright 2014 Google Inc.
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// This module is a mix of SkDashPath, SkDashPathEffect, SkContourMeasure and SkPathMeasure.
package path

import (
	"testing"

	"github.com/chewxy/math32"
	"github.com/lumifloat/tinyskia/internal/helper"
)

func TestStrokeDash(t *testing.T) {
	helper.AssertNil(t, NewStrokeDash([]float32{}, 0.0))
	helper.AssertNil(t, NewStrokeDash([]float32{1.0, -1.0}, 0.0))
	helper.AssertNil(t, NewStrokeDash([]float32{1.0, 2.0, 3.0}, 0.0))
	helper.AssertNil(t, NewStrokeDash([]float32{1.0, -2.0}, 0.0))
	helper.AssertNil(t, NewStrokeDash([]float32{0.0, 0.0}, 0.0))
	helper.AssertNil(t, NewStrokeDash([]float32{1.0, -1.0}, 0.0))
	helper.AssertNil(t, NewStrokeDash([]float32{1.0, 1.0}, math32.Inf(1)))
	helper.AssertNil(t, NewStrokeDash([]float32{1.0, math32.Inf(1)}, 0.0))

}

func TestBug26(t *testing.T) {
	// Create a path with multiple line segments
	pb := NewPathBuilder()
	pb.MoveTo(665.54, 287.3)
	pb.LineTo(675.67, 273.04)
	pb.LineTo(675.52, 271.32)
	pb.LineTo(674.79, 269.61)
	pb.LineTo(674.05, 268.04)
	pb.LineTo(672.88, 266.47)
	pb.LineTo(671.27, 264.9)

	path := pb.Finish()
	if path == nil {
		t.Fatal("Failed to create path")
	}

	// Create a stroke dash pattern
	strokeDash := NewStrokeDash([]float32{6.0, 4.5}, 0.0)
	if strokeDash == nil {
		t.Fatal("Failed to create StrokeDash")
	}

	// Dash the path - this should not return nil
	result := path.Dash(strokeDash, 1.0)
	if result == nil {
		t.Error("path.Dash() returned nil, expected non-nil result")
	}
}
