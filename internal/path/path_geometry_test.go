// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// ! A collection of functions to work with Bezier paths.
// !
// ! Mainly for internal use. Do not rely on it!
package path

import (
	"testing"

	"github.com/lumifloat/tinyskia/internal/helper"
	"github.com/lumifloat/tinyskia/internal/numeric/normalized"
)

func TestEvalCubicAt(t *testing.T) {
	src := [4]Point{
		{X: 30.0, Y: 40.0},
		{X: 30.0, Y: 40.0},
		{X: 171.0, Y: 45.0},
		{X: 180.0, Y: 155.0},
	}

	// Test eval_cubic_pos_at
	result1 := evalCubicPosAt(src, normalized.NormalizedF32Zero)
	helper.AssertEqual(t, result1, Point{X: 30.0, Y: 40.0})

	// Test eval_cubic_tangent_at
	result2 := evalCubicTangentAt(src, normalized.NormalizedF32Zero)
	helper.AssertEqual(t, result2, Point{X: 141.0, Y: 5.0})
}

func TestFindCubicMaxCurvature(t *testing.T) {
	src := [4]Point{
		{X: 20.0, Y: 160.0},
		{X: 20.0001, Y: 160.0},
		{X: 160.0, Y: 20.0},
		{X: 160.0001, Y: 20.0},
	}

	var tValues [3]normalized.NormalizedF32
	result := FindCubicMaxCurvature(src, &tValues)

	helper.AssertEqual(t, len(result), 3)
	helper.AssertEqual(t, result[0], normalized.NormalizedF32Zero)
	helper.AssertEqual(t, result[1], normalized.NewNormalizedF32WithClamped(0.5))
	helper.AssertEqual(t, result[2], normalized.NormalizedF32One)
}

func TestChopCubicAtYExtrema(t *testing.T) {
	src := [4]Point{
		{X: 10.0, Y: 20.0},
		{X: 67.0, Y: 437.0},
		{X: 298.0, Y: 213.0},
		{X: 401.0, Y: 214.0},
	}

	var dst [10]Point
	n := ChopCubicAtYExtrema(src, &dst)
	helper.AssertEqual(t, n, 2)

	helper.AssertEqual(t, dst[0], Point{X: 10.0, Y: 20.0})
	helper.AssertEqual(t, dst[1], Point{X: 37.508274, Y: 221.24475})
	helper.AssertEqual(t, dst[2], Point{X: 105.541855, Y: 273.19803})
	helper.AssertEqual(t, dst[3], Point{X: 180.15599, Y: 273.19803})
	helper.AssertEqual(t, dst[4], Point{X: 259.80502, Y: 273.19803})
	helper.AssertEqual(t, dst[5], Point{X: 346.9527, Y: 213.99666})
	helper.AssertEqual(t, dst[6], Point{X: 400.30844, Y: 213.99666})
	helper.AssertEqual(t, dst[7], Point{X: 400.53958, Y: 213.99666})
	helper.AssertEqual(t, dst[8], Point{X: 400.7701, Y: 213.99777})
	helper.AssertEqual(t, dst[9], Point{X: 401.0, Y: 214.0})
}
