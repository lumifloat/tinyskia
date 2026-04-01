// Copyright 2011 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package scan

import (
	"github.com/lumifloat/tinyskia/internal/core/blitter"
	"github.com/lumifloat/tinyskia/internal/path"
)

const fillRuleWinding int = 0
const fillRuleEvenOdd int = 1

// FillRect fills the specified rectangle with the blitter.
func FillRect(rect path.Rect, clip path.ScreenIntRect, blitter blitter.Blitter) {
	// Check for empty rectangle before attempting to fill
	if rect.Left() >= rect.Right() || rect.Top() >= rect.Bottom() {
		return
	}
	if r, ok := rect.Round(); ok {
		fillIntRect(r, clip, blitter)
	}
}

func fillIntRect(rect path.IntRect, clip path.ScreenIntRect, blitter blitter.Blitter) {
	intersected, ok := rect.Intersect(clip.ToIntRect())
	if !ok {
		return // everything was clipped out
	}

	// Check for empty rectangle before drawing
	if intersected.Width() <= 0 || intersected.Height() <= 0 {
		return
	}

	screenRect, ok := intersected.ToScreenIntRect()
	if !ok {
		return
	}

	blitter.BlitRect(screenRect)
}
