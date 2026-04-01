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

func TestIntSize(t *testing.T) {
	_, f1 := NewIntSize(0, 0)
	helper.AssertEqual(t, f1, false)
	_, f2 := NewIntSize(1, 0)
	helper.AssertEqual(t, f2, false)
	_, f3 := NewIntSize(0, 1)
	helper.AssertEqual(t, f3, false)

	size1, f4 := NewIntSize(3, 4)
	helper.AssertEqual(t, f4, true)
	size2, f5 := NewIntRectFromXYWH(1, 2, 3, 4)
	helper.AssertEqual(t, f5, true)
	helper.AssertEqual(t, size1.ToIntRect(1, 2), size2)
}
