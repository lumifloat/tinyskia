// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package path

import (
	"testing"

	"github.com/chewxy/math32"
	"github.com/lumifloat/tinyskia/internal/helper"
)

func TestIntRect(t *testing.T) {
	// Test basic validation
	_, f1 := NewIntRectFromXYWH(0, 0, 0, 0)
	helper.AssertEqual(t, f1, false)
	_, f2 := NewIntRectFromXYWH(0, 0, 1, 0)
	helper.AssertEqual(t, f2, false)
	_, f3 := NewIntRectFromXYWH(0, 0, 0, 1)
	helper.AssertEqual(t, f3, false)

	// Test overflow cases
	_, f4 := NewIntRectFromXYWH(0, 0, math32.MaxUint32, math32.MaxUint32)
	helper.AssertEqual(t, f4, false)
	_, f5 := NewIntRectFromXYWH(0, 0, 1, math32.MaxUint32)
	helper.AssertEqual(t, f5, false)
	_, f6 := NewIntRectFromXYWH(0, 0, math32.MaxUint32, 1)
	helper.AssertEqual(t, f6, false)

	_, f7 := NewIntRectFromXYWH(math32.MaxInt32, 0, 1, 1)
	helper.AssertEqual(t, f7, false)
	_, f8 := NewIntRectFromXYWH(0, math32.MaxInt32, 1, 1)
	helper.AssertEqual(t, f8, false)

	// Test intersection - no intersection
	r1, ok1 := NewIntRectFromXYWH(1, 2, 3, 4)
	helper.AssertEqual(t, ok1, true)
	r2, ok2 := NewIntRectFromXYWH(11, 12, 13, 14)
	helper.AssertEqual(t, ok2, true)
	_, intersectOk1 := r1.Intersect(r2)
	helper.AssertEqual(t, intersectOk1, false)

	// Test intersection - second inside first
	r3, ok3 := NewIntRectFromXYWH(1, 2, 30, 40)
	helper.AssertEqual(t, ok3, true)
	r4, ok4 := NewIntRectFromXYWH(11, 12, 13, 14)
	helper.AssertEqual(t, ok4, true)
	intersect1, intersectOk2 := r3.Intersect(r4)
	helper.AssertEqual(t, intersectOk2, true)
	expected1, _ := NewIntRectFromXYWH(11, 12, 13, 14)
	helper.AssertEqual(t, intersect1, expected1)

	// Test intersection - partial overlap
	r5, ok5 := NewIntRectFromXYWH(1, 2, 30, 40)
	helper.AssertEqual(t, ok5, true)
	r6, ok6 := NewIntRectFromXYWH(11, 12, 50, 60)
	helper.AssertEqual(t, ok6, true)
	intersect2, intersectOk3 := r5.Intersect(r6)
	helper.AssertEqual(t, intersectOk3, true)
	expected2, _ := NewIntRectFromXYWH(11, 12, 20, 30)
	helper.AssertEqual(t, intersect2, expected2)
}

func TestRect(t *testing.T) {
	// Test basic validation
	_, f1 := NewRectFromLTRB(10.0, 10.0, 5.0, 10.0)
	helper.AssertEqual(t, f1, false)
	_, f2 := NewRectFromLTRB(10.0, 10.0, 10.0, 5.0)
	helper.AssertEqual(t, f2, false)
	_, f3 := NewRectFromLTRB(math32.NaN(), 10.0, 10.0, 10.0)
	helper.AssertEqual(t, f3, false)
	_, f4 := NewRectFromLTRB(10.0, math32.NaN(), 10.0, 10.0)
	helper.AssertEqual(t, f4, false)
	_, f5 := NewRectFromLTRB(10.0, 10.0, math32.NaN(), 10.0)
	helper.AssertEqual(t, f5, false)
	_, f6 := NewRectFromLTRB(10.0, 10.0, 10.0, math32.NaN())
	helper.AssertEqual(t, f6, false)
	_, f7 := NewRectFromLTRB(10.0, 10.0, 10.0, math32.Inf(1))
	helper.AssertEqual(t, f7, false)

	// Test valid rect
	rect, ok := NewRectFromLTRB(10.0, 20.0, 30.0, 40.0)
	helper.AssertEqual(t, ok, true)
	helper.AssertEqual(t, rect.Left(), float32(10.0))
	helper.AssertEqual(t, rect.Top(), float32(20.0))
	helper.AssertEqual(t, rect.Right(), float32(30.0))
	helper.AssertEqual(t, rect.Bottom(), float32(40.0))
	helper.AssertEqual(t, rect.Width(), float32(20.0))
	helper.AssertEqual(t, rect.Height(), float32(20.0))

	// Test negative coordinates
	rect2, ok2 := NewRectFromLTRB(-30.0, 20.0, -10.0, 40.0)
	helper.AssertEqual(t, ok2, true)
	helper.AssertEqual(t, rect2.Width(), float32(20.0))
	helper.AssertEqual(t, rect2.Height(), float32(20.0))
}
