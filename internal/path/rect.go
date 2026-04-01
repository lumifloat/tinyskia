// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package path

import (
	"github.com/chewxy/math32"
)

// IntRect is an integer rectangle.
//
// Guarantees:
// - Width and height are in 1..=i32.Max range.
// - x+width and y+height does not overflow.
type IntRect struct {
	x      int32
	y      int32
	width  uint32
	height uint32
}

// NewIntRectFromXYWH creates a new IntRect.
func NewIntRectFromXYWH(x, y int32, width, height uint32) (IntRect, bool) {
	// Width and height must be non-zero and fit in i32 range
	if width == 0 || height == 0 || width > math32.MaxInt32 || height > math32.MaxInt32 {
		return IntRect{}, false
	}

	// Check for overflow: x + width and y + height must not overflow i32
	if x > math32.MaxInt32-int32(width) || y > math32.MaxInt32-int32(height) {
		return IntRect{}, false
	}

	return IntRect{x, y, width, height}, true
}

// NewIntRectFromLTRB creates a new IntRect.
func NewIntRectFromLTRB(left, top, right, bottom int32) (IntRect, bool) {
	// right must be > left, bottom must be > top
	if right <= left || bottom <= top {
		return IntRect{}, false
	}

	// Calculate width and height with overflow check
	width := uint32(right - left)
	height := uint32(bottom - top)

	return NewIntRectFromXYWH(left, top, width, height)
}

func (r IntRect) X() int32       { return r.x }
func (r IntRect) Y() int32       { return r.y }
func (r IntRect) Width() uint32  { return r.width }
func (r IntRect) Height() uint32 { return r.height }
func (r IntRect) Left() int32    { return r.x }
func (r IntRect) Top() int32     { return r.y }
func (r IntRect) Right() int32   { return r.x + int32(r.width) }
func (r IntRect) Bottom() int32  { return r.y + int32(r.height) }

// Contains checks that the rect completely includes other Rect.
func (r IntRect) Contains(other IntRect) bool {
	return r.x <= other.x &&
		r.y <= other.y &&
		r.Right() >= other.Right() &&
		r.Bottom() >= other.Bottom()
}

// Intersect returns an intersection of two rectangles.
func (r IntRect) Intersect(other IntRect) (IntRect, bool) {
	left := r.x
	if other.x > left {
		left = other.x
	}
	top := r.y
	if other.y > top {
		top = other.y
	}
	right := r.Right()
	if other.Right() < right {
		right = other.Right()
	}
	bottom := r.Bottom()
	if other.Bottom() < bottom {
		bottom = other.Bottom()
	}

	return NewIntRectFromLTRB(left, top, int32(right), int32(bottom))
}

// Inset insets the rectangle.
func (r IntRect) Inset(dx, dy int32) (IntRect, bool) {
	return NewIntRectFromLTRB(r.Left()+dx, r.Top()+dy, int32(r.Right())-dx, int32(r.Bottom())-dy)
}

// Outset outsets the rectangle.
func (r IntRect) Outset(dx, dy int32) (IntRect, bool) {
	// Simple saturating add/sub without extra utils
	l, t, ri, b := r.Left()-dx, r.Top()-dy, int32(r.Right())+dx, int32(r.Bottom())+dy
	return NewIntRectFromLTRB(l, t, ri, b)
}

// Translate translates the rect by the specified offset.
func (r IntRect) Translate(tx, ty int32) (IntRect, bool) {
	return NewIntRectFromXYWH(r.x+tx, r.y+ty, r.width, r.height)
}

// TranslateTo translates the rect to the specified position.
func (r IntRect) TranslateTo(x, y int32) (IntRect, bool) {
	return NewIntRectFromXYWH(x, y, r.width, r.height)
}

// ToRect converts into Rect.
func (r IntRect) ToRect() Rect {
	rect, _ := NewRectFromLTRB(float32(r.x), float32(r.y), float32(r.Right()), float32(r.Bottom()))
	return rect
}

// ToScreenIntRect converts into ScreenIntRect.
func (r IntRect) ToScreenIntRect() (ScreenIntRect, bool) {
	if r.X() < 0 || r.Y() < 0 {
		return ScreenIntRect{}, false
	}
	return NewScreenIntRectFromXYWH(uint32(r.X()), uint32(r.Y()), uint32(r.Width()), uint32(r.Height()))
}

// Rect is a rectangle defined by left, top, right and bottom edges.
type Rect struct {
	left, top, right, bottom float32
}

// NewRectFromLTRB creates new Rect.
func NewRectFromLTRB(left, top, right, bottom float32) (Rect, bool) {
	if math32.IsNaN(left) || math32.IsInf(left, 0) ||
		math32.IsNaN(top) || math32.IsInf(top, 0) ||
		math32.IsNaN(right) || math32.IsInf(right, 0) ||
		math32.IsNaN(bottom) || math32.IsInf(bottom, 0) {
		return Rect{}, false
	}

	if left <= right && top <= bottom {
		return Rect{left, top, right, bottom}, true
	}
	return Rect{}, false
}

// NewRectFromXYWH creates new Rect.
func NewRectFromXYWH(x, y, w, h float32) (Rect, bool) {
	return NewRectFromLTRB(x, y, x+w, y+h)
}

func (r Rect) Left() float32   { return r.left }
func (r Rect) Top() float32    { return r.top }
func (r Rect) Right() float32  { return r.right }
func (r Rect) Bottom() float32 { return r.bottom }
func (r Rect) X() float32      { return r.left }
func (r Rect) Y() float32      { return r.top }
func (r Rect) Width() float32  { return r.right - r.left }
func (r Rect) Height() float32 { return r.bottom - r.top }
func (r Rect) IsEmpty() bool   { return r.left == r.right || r.top == r.bottom }

// Round converts into an IntRect by adding 0.5 and discarding the fractional portion.
func (r Rect) Round() (IntRect, bool) {
	left := int32(math32.Floor(r.left + 0.5))
	top := int32(math32.Floor(r.top + 0.5))
	right := int32(math32.Floor(r.right + 0.5))
	bottom := int32(math32.Floor(r.bottom + 0.5))

	w, h := right-left, bottom-top
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}

	return NewIntRectFromXYWH(left, top, uint32(w), uint32(h))
}

// RoundOut converts into an IntRect rounding outwards.
func (r Rect) RoundOut() (IntRect, bool) {
	left := int32(math32.Floor(r.left))
	top := int32(math32.Floor(r.top))
	right := int32(math32.Ceil(r.right))
	bottom := int32(math32.Ceil(r.bottom))

	w, h := right-left, bottom-top
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}

	return NewIntRectFromXYWH(left, top, uint32(w), uint32(h))
}

// Intersect returns an intersection of two rectangles.
func (r Rect) Intersect(other Rect) (Rect, bool) {
	return NewRectFromLTRB(
		math32.Max(r.left, other.left),
		math32.Max(r.top, other.top),
		math32.Min(r.right, other.right),
		math32.Min(r.bottom, other.bottom),
	)
}

// Join returns the union of two rectangles.
func (r Rect) Join(other Rect) (Rect, bool) {
	if other.IsEmpty() {
		return r, true
	}
	if r.IsEmpty() {
		return other, true
	}

	return NewRectFromLTRB(
		math32.Min(r.left, other.left),
		math32.Min(r.top, other.top),
		math32.Max(r.right, other.right),
		math32.Max(r.bottom, other.bottom),
	)
}

// NewRectFromPoints creates a Rect from Point array.
func NewRectFromPoints(points []Point) (Rect, bool) {
	if len(points) == 0 {
		return Rect{}, false
	}

	minX, minY := points[0].X, points[0].Y
	maxX, maxY := minX, minY

	for _, p := range points {
		if math32.IsNaN(p.X) || math32.IsInf(p.X, 0) || math32.IsNaN(p.Y) || math32.IsInf(p.Y, 0) {
			return Rect{}, false
		}
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	return NewRectFromLTRB(minX, minY, maxX, maxY)
}

// Inset insets the rectangle by the specified offset.
func (r Rect) Inset(dx, dy float32) (Rect, bool) {
	return NewRectFromLTRB(r.left+dx, r.top+dy, r.right-dx, r.bottom-dy)
}

// Outset outsets the rectangle by the specified offset.
func (r Rect) Outset(dx, dy float32) (Rect, bool) {
	return r.Inset(-dx, -dy)
}

// Transform transforms the rect using the provided Transform.
//
// If the transform is a skew, the result will be a bounding box around the skewed rectangle.
func (r Rect) Transform(transform Transform) (Rect, bool) {
	if transform.IsIdentity() {
		return r, true
	} else if transform.HasSkew() {
		// We need to transform all 4 corners
		lt := Point{X: r.left, Y: r.top}
		rt := Point{X: r.right, Y: r.top}
		lb := Point{X: r.left, Y: r.bottom}
		rb := Point{X: r.right, Y: r.bottom}
		points := []Point{lt, rt, lb, rb}
		transform.MapPoints(points)
		return NewRectFromPoints(points)
	} else {
		// Faster (more common) case - only need to transform 2 points
		lt := Point{X: r.left, Y: r.top}
		rb := Point{X: r.right, Y: r.bottom}
		points := []Point{lt, rb}
		transform.MapPoints(points)
		return NewRectFromPoints(points)
	}
}
