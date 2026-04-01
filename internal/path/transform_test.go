// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package path

import (
	"testing"

	"github.com/lumifloat/tinyskia/internal/helper"
)

func TestTransform(t *testing.T) {
	// Test identity transform
	identity := Identity()
	expected1 := NewTransform(1.0, 0.0, 0.0, 1.0, 0.0, 0.0)
	helper.AssertEqual(t, identity, expected1)

	// Test scale transform
	scale := NewTransformFromScale(1.0, 2.0)
	expected2 := NewTransform(1.0, 0.0, 0.0, 2.0, 0.0, 0.0)
	helper.AssertEqual(t, scale, expected2)

	// Test skew transform
	skew := NewTransformFromSkew(2.0, 3.0)
	expected3 := NewTransform(1.0, 3.0, 2.0, 1.0, 0.0, 0.0)
	helper.AssertEqual(t, skew, expected3)

	// Test translate transform
	translate := NewTransformFromTranslate(5.0, 6.0)
	expected4 := NewTransform(1.0, 0.0, 0.0, 1.0, 5.0, 6.0)
	helper.AssertEqual(t, translate, expected4)

	// Test identity properties
	ts := Identity()
	helper.AssertEqual(t, ts.IsIdentity(), true)
	helper.AssertEqual(t, ts.IsScale(), false)
	helper.AssertEqual(t, ts.IsSkew(), false)
	helper.AssertEqual(t, ts.IsTranslate(), false)
	helper.AssertEqual(t, ts.IsScaleTranslate(), false)
	helper.AssertEqual(t, ts.HasScale(), false)
	helper.AssertEqual(t, ts.HasSkew(), false)
	helper.AssertEqual(t, ts.HasTranslate(), false)

	// Test scale properties
	ts2 := NewTransformFromScale(2.0, 3.0)
	helper.AssertEqual(t, ts2.IsIdentity(), false)
	helper.AssertEqual(t, ts2.IsScale(), true)
	helper.AssertEqual(t, ts2.IsSkew(), false)
	helper.AssertEqual(t, ts2.IsTranslate(), false)
	helper.AssertEqual(t, ts2.IsScaleTranslate(), true)
	helper.AssertEqual(t, ts2.HasScale(), true)
	helper.AssertEqual(t, ts2.HasSkew(), false)
	helper.AssertEqual(t, ts2.HasTranslate(), false)

	// Test skew properties
	ts3 := NewTransformFromSkew(2.0, 3.0)
	helper.AssertEqual(t, ts3.IsIdentity(), false)
	helper.AssertEqual(t, ts3.IsScale(), false)
	helper.AssertEqual(t, ts3.IsSkew(), true)
	helper.AssertEqual(t, ts3.IsTranslate(), false)
	helper.AssertEqual(t, ts3.IsScaleTranslate(), false)
	helper.AssertEqual(t, ts3.HasScale(), false)
	helper.AssertEqual(t, ts3.HasSkew(), true)
	helper.AssertEqual(t, ts3.HasTranslate(), false)

	// Test translate properties
	ts4 := NewTransformFromTranslate(2.0, 3.0)
	helper.AssertEqual(t, ts4.IsIdentity(), false)
	helper.AssertEqual(t, ts4.IsScale(), false)
	helper.AssertEqual(t, ts4.IsSkew(), false)
	helper.AssertEqual(t, ts4.IsTranslate(), true)
	helper.AssertEqual(t, ts4.IsScaleTranslate(), true)
	helper.AssertEqual(t, ts4.HasScale(), false)
	helper.AssertEqual(t, ts4.HasSkew(), false)
	helper.AssertEqual(t, ts4.HasTranslate(), true)

	// Test general transform properties
	ts5 := NewTransform(1.0, 2.0, 3.0, 4.0, 5.0, 6.0)
	helper.AssertEqual(t, ts5.IsIdentity(), false)
	helper.AssertEqual(t, ts5.IsScale(), false)
	helper.AssertEqual(t, ts5.IsSkew(), false)
	helper.AssertEqual(t, ts5.IsTranslate(), false)
	helper.AssertEqual(t, ts5.IsScaleTranslate(), false)
	helper.AssertEqual(t, ts5.HasScale(), true)
	helper.AssertEqual(t, ts5.HasSkew(), true)
	helper.AssertEqual(t, ts5.HasTranslate(), true)

	// Test edge cases for has_* methods
	ts6 := NewTransformFromScale(1.0, 1.0)
	helper.AssertEqual(t, ts6.HasScale(), false)

	ts7 := NewTransformFromSkew(0.0, 0.0)
	helper.AssertEqual(t, ts7.HasSkew(), false)

	ts8 := NewTransformFromTranslate(0.0, 0.0)
	helper.AssertEqual(t, ts8.HasTranslate(), false)
}

func TestTransformConcat(t *testing.T) {
	// Test pre-scale
	ts := NewTransform(1.2, 3.4, -5.6, -7.8, 1.2, 3.4)
	ts = ts.PreScale(2.0, -4.0)
	expected1 := NewTransform(2.4, 6.8, 22.4, 31.2, 1.2, 3.4)
	helper.AssertEqual(t, ts, expected1)

	// Test post-scale
	ts2 := NewTransform(1.2, 3.4, -5.6, -7.8, 1.2, 3.4)
	ts2 = ts2.PostScale(2.0, -4.0)
	expected2 := NewTransform(2.4, -13.6, -11.2, 31.2, 2.4, -13.6)
	helper.AssertEqual(t, ts2, expected2)
}
