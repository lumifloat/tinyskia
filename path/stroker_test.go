// Copyright 2008 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Based on SkStroke.cpp
package path

import (
	"math"
	"testing"

	"github.com/lumifloat/tinyskia/internal/helper"
)

func TestAutoClose(t *testing.T) {
	// A triangle.
	pb := NewPathBuilder()
	pb.MoveTo(10.0, 10.0)
	pb.LineTo(20.0, 50.0)
	pb.LineTo(30.0, 10.0)
	pb.Close()
	path := pb.Finish()

	stroke := defaultStroke()
	strokePath := newPathStroker().Stroke(path, stroke, 1.0)
	if strokePath == nil {
		t.Fatal("Expected stroked path, got nil")
	}

	iter := strokePath.Segments()
	iter.SetAutoClose(true)

	helper.AssertEqual(t, iter.Next(), PathSegmentMoveTo(Point{10.485071, 9.878732}))
	helper.AssertEqual(t, iter.Next(), PathSegmentLineTo(Point{20.485071, 49.878731}))
	helper.AssertEqual(t, iter.Next(), PathSegmentLineTo(Point{20.0, 50.0}))
	helper.AssertEqual(t, iter.Next(), PathSegmentLineTo(Point{19.514929, 49.878731}))
	helper.AssertEqual(t, iter.Next(), PathSegmentLineTo(Point{29.514929, 9.878732}))
	helper.AssertEqual(t, iter.Next(), PathSegmentLineTo(Point{30.0, 10.0}))
	helper.AssertEqual(t, iter.Next(), PathSegmentLineTo(Point{30.0, 10.5}))
	helper.AssertEqual(t, iter.Next(), PathSegmentLineTo(Point{10.0, 10.5}))
	helper.AssertEqual(t, iter.Next(), PathSegmentLineTo(Point{10.0, 10.0}))
	helper.AssertEqual(t, iter.Next(), PathSegmentLineTo(Point{10.485071, 9.878732}))
	helper.AssertEqual(t, iter.Next(), PathSegmentClose(struct{}{}))
	helper.AssertEqual(t, iter.Next(), PathSegmentMoveTo(Point{9.3596115, 9.5}))
	helper.AssertEqual(t, iter.Next(), PathSegmentLineTo(Point{30.640388, 9.5}))
	helper.AssertEqual(t, iter.Next(), PathSegmentLineTo(Point{20.485071, 50.121269}))
	helper.AssertEqual(t, iter.Next(), PathSegmentLineTo(Point{19.514929, 50.121269}))
	helper.AssertEqual(t, iter.Next(), PathSegmentLineTo(Point{9.514929, 10.121268}))
	helper.AssertEqual(t, iter.Next(), PathSegmentLineTo(Point{9.3596115, 9.5}))
	helper.AssertEqual(t, iter.Next(), PathSegmentClose(struct{}{}))
}

func TestCubic1(t *testing.T) {
	pb := NewPathBuilder()
	pb.MoveTo(51.0161362, 1511.52478)
	pb.CubicTo(
		51.0161362, 1511.52478,
		51.0161362, 1511.52478,
		51.0161362, 1511.52478,
	)
	path := pb.Finish()

	stroke := defaultStroke()
	stroke.Width = 0.394537568

	// Should return nil for degenerate cubic
	result := newPathStroker().Stroke(path, stroke, 1.0)
	if result != nil {
		t.Error("Expected nil result for degenerate cubic")
	}
}

func TestCubic2(t *testing.T) {
	pb := NewPathBuilder()
	pb.MoveTo(math.Float32frombits(0x424c1086), math.Float32frombits(0x44bcf0cb)) // 51.0161362, 1511.52478
	pb.CubicTo(
		math.Float32frombits(0x424c107c), math.Float32frombits(0x44bcf0cb), // 51.0160980, 1511.52478
		math.Float32frombits(0x424c10c2), math.Float32frombits(0x44bcf0cb), // 51.0163651, 1511.52478
		math.Float32frombits(0x424c1119), math.Float32frombits(0x44bcf0ca), // 51.0166969, 1511.52466
	)
	path := pb.Finish()

	stroke := defaultStroke()
	stroke.Width = 0.394537568

	// Should produce a valid result
	result := newPathStroker().Stroke(path, stroke, 1.0)
	if result == nil {
		t.Error("Expected valid result for non-degenerate cubic")
	}
}

func TestBig(t *testing.T) {
	// Skia uses `kStrokeAndFill_Style` here, but we do not support it.

	pb := NewPathBuilder()
	pb.MoveTo(math.Float32frombits(0x46380000), math.Float32frombits(0xc6380000)) // 11776, -11776
	pb.LineTo(math.Float32frombits(0x46a00000), math.Float32frombits(0xc6a00000)) // 20480, -20480
	pb.LineTo(math.Float32frombits(0x468c0000), math.Float32frombits(0xc68c0000)) // 17920, -17920
	pb.LineTo(math.Float32frombits(0x46100000), math.Float32frombits(0xc6100000)) // 9216, -9216
	pb.LineTo(math.Float32frombits(0x46380000), math.Float32frombits(0xc6380000)) // 11776, -11776
	pb.Close()
	path := pb.Finish()

	stroke := defaultStroke()
	stroke.Width = 1.49679073e+10

	// Should handle large stroke widths without crashing
	result := newPathStroker().Stroke(path, stroke, 1.0)
	if result == nil {
		t.Error("Expected valid result for large stroke width")
	}
}

func TestQuadStrokerOneOff(t *testing.T) {
	pb := NewPathBuilder()
	pb.MoveTo(math.Float32frombits(0x43c99223), math.Float32frombits(0x42b7417e))
	pb.QuadTo(
		math.Float32frombits(0x4285d839), math.Float32frombits(0x43ed6645),
		math.Float32frombits(0x43c941c8), math.Float32frombits(0x42b3ace3),
	)
	path := pb.Finish()

	stroke := defaultStroke()
	stroke.Width = 164.683548

	// Should handle this edge case
	result := newPathStroker().Stroke(path, stroke, 1.0)
	if result == nil {
		t.Error("Expected valid result for quad stroker edge case")
	}
}

func TestCubicStrokerOneOff(t *testing.T) {
	pb := NewPathBuilder()
	pb.MoveTo(math.Float32frombits(0x433f5370), math.Float32frombits(0x43d1f4b3))
	pb.CubicTo(
		math.Float32frombits(0x4331cb76), math.Float32frombits(0x43ea3340),
		math.Float32frombits(0x4388f498), math.Float32frombits(0x42f7f08d),
		math.Float32frombits(0x43f1cd32), math.Float32frombits(0x42802ec1),
	)
	path := pb.Finish()

	stroke := defaultStroke()
	stroke.Width = 42.835968

	// Should handle this edge case
	result := newPathStroker().Stroke(path, stroke, 1.0)
	if result == nil {
		t.Error("Expected valid result for cubic stroker edge case")
	}
}
