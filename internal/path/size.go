// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package path

import (
	"github.com/chewxy/math32"
)

// IntSize is an integer size.
//
// Guarantees:
// - Width and height are positive and non-zero.
type IntSize struct {
	width  uint32
	height uint32
}

// NewIntSize creates a new IntSize from width and height.
func NewIntSize(width, height uint32) (IntSize, bool) {
	if width == 0 || height == 0 || width > math32.MaxInt32 || height > math32.MaxInt32 {
		return IntSize{}, false
	}
	return IntSize{width, height}, true
}

// Width returns width.
func (s IntSize) Width() uint32 {
	return s.width
}

// Height returns height.
func (s IntSize) Height() uint32 {
	return s.height
}

// Dimensions returns width and height.
func (s IntSize) Dimensions() (uint32, uint32) {
	return s.width, s.height
}

// ScaleBy scales current size by the specified factor.
func (s IntSize) ScaleBy(factor float32) (IntSize, bool) {
	return NewIntSize(
		uint32(math32.Round(float32(s.width)*factor)),
		uint32(math32.Round(float32(s.height)*factor)),
	)
}

// ScaleTo scales current size to the specified size.
func (s IntSize) ScaleTo(to IntSize) IntSize {
	return intSizeScale(s, to, false)
}

// ScaleToWidth scales current size to the specified width.
func (s IntSize) ScaleToWidth(newWidth uint32) (IntSize, bool) {
	newHeight := math32.Ceil(float32(newWidth) * float32(s.height) / float32(s.width))
	return NewIntSize(newWidth, uint32(newHeight))
}

// ScaleToHeight scales current size to the specified height.
func (s IntSize) ScaleToHeight(newHeight uint32) (IntSize, bool) {
	newWidth := math32.Ceil(float32(newHeight) * float32(s.width) / float32(s.height))
	return NewIntSize(uint32(newWidth), newHeight)
}

// ToSize converts into Size.
func (s IntSize) ToSize() Size {
	sz, _ := NewSize(float32(s.width), float32(s.height))
	return sz
}

// ToIntRect converts into IntRect at the provided position.
func (s IntSize) ToIntRect(x, y int32) IntRect {
	r, _ := NewIntRectFromXYWH(x, y, s.width, s.height)
	return r
}

func (s IntSize) ToScreenIntRect(x, y uint32) ScreenIntRect {
	return NewScreenIntRectFromXYWHSafe(x, y, uint32(s.width), uint32(s.height))
}

func intSizeScale(s1, s2 IntSize, expand bool) IntSize {
	rw := uint32(math32.Ceil(float32(s2.height) * float32(s1.width) / float32(s1.height)))
	withH := false
	if expand {
		withH = rw <= s2.width
	} else {
		withH = rw >= s2.width
	}

	if !withH {
		sz, _ := NewIntSize(rw, s2.height)
		return sz
	} else {
		h := uint32(math32.Ceil(float32(s2.width) * float32(s1.height) / float32(s1.width)))
		sz, _ := NewIntSize(s2.width, h)
		return sz
	}
}

// Size is a size.
//
// Guarantees:
// - Width and height are positive, non-zero and finite.
type Size struct {
	width  float32
	height float32
}

// NewSize creates a new Size from width and height.
func NewSize(width, height float32) (Size, bool) {
	if width <= 0 || height <= 0 ||
		math32.IsNaN(width) || math32.IsInf(width, 0) ||
		math32.IsNaN(height) || math32.IsInf(height, 0) {
		return Size{}, false
	}
	return Size{width, height}, true
}

// Width returns width.
func (s Size) Width() float32 {
	return s.width
}

// Height returns height.
func (s Size) Height() float32 {
	return s.height
}

// ScaleTo scales current size to specified size.
func (s Size) ScaleTo(to Size) Size {
	return sizeScale(s, to, false)
}

// ExpandTo expands current size to specified size.
func (s Size) ExpandTo(to Size) Size {
	return sizeScale(s, to, true)
}

// ScaleBy scales current size by the specified factor.
func (s Size) ScaleBy(factor float32) (Size, bool) {
	return NewSize(s.width*factor, s.height*factor)
}

// ScaleToWidth scales current size to the specified width.
func (s Size) ScaleToWidth(newWidth float32) (Size, bool) {
	newHeight := newWidth * s.height / s.width
	return NewSize(newWidth, newHeight)
}

// ScaleToHeight scales current size to the specified height.
func (s Size) ScaleToHeight(newHeight float32) (Size, bool) {
	newWidth := newHeight * s.width / s.height
	return NewSize(newWidth, newHeight)
}

// ToIntSize converts into IntSize.
func (s Size) ToIntSize() IntSize {
	w := uint32(math32.Round(s.width))
	if w < 1 {
		w = 1
	}
	h := uint32(math32.Round(s.height))
	if h < 1 {
		h = 1
	}
	sz, _ := NewIntSize(w, h)
	return sz
}

// ToRect converts the current size to Rect at provided position.
func (s Size) ToRect(x, y float32) (Rect, bool) {
	return NewRectFromXYWH(x, y, s.width, s.height)
}

func sizeScale(s1, s2 Size, expand bool) Size {
	rw := s2.height * s1.width / s1.height
	withH := false
	if expand {
		withH = rw <= s2.width
	} else {
		withH = rw >= s2.width
	}

	if !withH {
		sz, _ := NewSize(rw, s2.height)
		return sz
	} else {
		h := s2.width * s1.height / s1.width
		sz, _ := NewSize(s2.width, h)
		return sz
	}
}
