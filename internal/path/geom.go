// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package path

import (
	"github.com/chewxy/math32"
)

// ScreenIntRect is a screen IntRect.
//
// Guarantees:
// - X and Y are in 0..=math.MaxInt32 range.
// - Width and height are in 1..=math.MaxInt32 range.
// - x+width and y+height does not overflow.
type ScreenIntRect struct {
	x      uint32
	y      uint32
	width  uint32
	height uint32
}

// NewScreenIntRectFromXYWH creates a new ScreenIntRect.
func NewScreenIntRectFromXYWH(x, y, width, height uint32) (ScreenIntRect, bool) {
	// Check if values fit in int32 range
	if x > math32.MaxInt32 || y > math32.MaxInt32 || width > math32.MaxInt32 || height > math32.MaxInt32 {
		return ScreenIntRect{}, false
	}

	// Checked add to prevent overflow
	result := x + width
	if result < x {
		return ScreenIntRect{}, false
	}
	result = y + height
	if result < y {
		return ScreenIntRect{}, false
	}

	// Width and height must be at least 1 (LengthU32 equivalent)
	if width == 0 || height == 0 {
		return ScreenIntRect{}, false
	}

	return ScreenIntRect{
		x:      x,
		y:      y,
		width:  width,
		height: height,
	}, true
}

// NewScreenIntRectFromXYWHSafe creates a new ScreenIntRect without checks.
// Only use when you're certain the values are valid.
func NewScreenIntRectFromXYWHSafe(x, y uint32, width, height uint32) ScreenIntRect {
	return ScreenIntRect{
		x:      x,
		y:      y,
		width:  width,
		height: height,
	}
}

// X returns rect's X position.
func (r ScreenIntRect) X() uint32 {
	return r.x
}

// Y returns rect's Y position.
func (r ScreenIntRect) Y() uint32 {
	return r.y
}

// Width returns rect's width.
func (r ScreenIntRect) Width() uint32 {
	return r.width
}

// Height returns rect's height.
func (r ScreenIntRect) Height() uint32 {
	return r.height
}

// Left returns rect's left edge.
func (r ScreenIntRect) Left() uint32 {
	return r.x
}

// Top returns rect's top edge.
func (r ScreenIntRect) Top() uint32 {
	return r.y
}

// Right returns rect's right edge.
// The right edge is at least 1.
func (r ScreenIntRect) Right() uint32 {
	// No overflow is guaranteed by constructors.
	return r.x + r.width
}

// Bottom returns rect's bottom edge.
// The bottom edge is at least 1.
func (r ScreenIntRect) Bottom() uint32 {
	// No overflow is guaranteed by constructors.
	return r.y + r.height
}

// Size returns rect's size.
func (r ScreenIntRect) Size() IntSize {
	// Can't fail because ScreenIntRect is always valid.
	size, _ := NewIntSize(r.width, r.height)
	return size
}

// Contains checks that the rect completely includes other Rect.
func (r ScreenIntRect) Contains(other ScreenIntRect) bool {
	return r.x <= other.x &&
		r.y <= other.y &&
		r.Right() >= other.Right() &&
		r.Bottom() >= other.Bottom()
}

// ToIntRect converts into a IntRect.
func (r ScreenIntRect) ToIntRect() IntRect {
	// Everything is already checked by constructors.
	rect, _ := NewIntRectFromXYWH(
		int32(r.x),
		int32(r.y),
		r.Width(),
		r.Height(),
	)
	return rect
}

// ToRect converts into a Rect.
func (r ScreenIntRect) ToRect() Rect {
	// Can't fail because ScreenIntRect is always valid.
	// And uint32 always fits into float32.
	rect, _ := NewRectFromLTRB(
		float32(r.x),
		float32(r.y),
		float32(r.x)+float32(r.width),
		float32(r.y)+float32(r.height),
	)
	return rect
}
