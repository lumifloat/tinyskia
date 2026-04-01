// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"github.com/lumifloat/tinyskia/internal/path"
)

// MaskType is a mask type.
type MaskType int

const (
	// MaskTypeAlpha transfers only the Alpha channel from Pixmap to Mask.
	MaskTypeAlpha MaskType = iota
	// MaskTypeLuminance transfers RGB channels as luminance from Pixmap to Mask.
	MaskTypeLuminance
)

// Mask is a mask.
type Mask struct {
	data []uint8
	size path.IntSize
}

// SubMaskRef is a reference to a sub-portion of a mask.
type SubMask struct {
	Data      []uint8
	Size      path.IntSize
	RealWidth uint32
}

// NewMask creates a new mask, allocating a buffer of the given size.
func NewMask(width, height uint32) *Mask {
	size, ok := path.NewIntSize(width, height)
	if !ok {
		return nil
	}
	return &Mask{
		data: make([]uint8, width*height),
		size: size,
	}
}

// FromVec creates a new mask by taking ownership over a mask buffer.
func FromVec(data []uint8, size path.IntSize) *Mask {
	dataLen := uint64(size.Width()) * uint64(size.Height())
	if uint64(len(data)) != dataLen {
		return nil
	}

	return &Mask{data: data, size: size}
}

// Width returns mask's width.
func (m *Mask) Width() uint32 {
	return m.size.Width()
}

// Height returns mask's height.
func (m *Mask) Height() uint32 {
	return m.size.Height()
}

// Data returns the internal data.
func (m *Mask) Data() []uint8 {
	return m.data
}

// Take consumes the mask and returns its owned internal data.
func (m *Mask) Take() []uint8 {
	d := m.data
	m.data = nil
	return d
}

func (m *Mask) AsSubMask() *SubMask {
	return &SubMask{
		Data:      m.data,
		Size:      m.size,
		RealWidth: uint32(m.size.Width()),
	}
}

func (m *Mask) SubMask(rect path.IntRect) *SubMask {
	r, ok := m.size.ToIntRect(0, 0).Intersect(rect)
	if !ok {
		return nil
	}
	rowBytes := m.Width()
	offset := uint64(r.Top())*uint64(rowBytes) + uint64(r.Left())

	size, _ := path.NewIntSize(r.Width(), r.Height())
	return &SubMask{
		Data:      m.data[offset:],
		Size:      size,
		RealWidth: uint32(m.size.Width()),
	}
}

// Invert inverts the mask.
func (m *Mask) Invert() {
	for i := range m.data {
		m.data[i] = 255 - m.data[i]
	}
}

// Clear clears the mask.
func (m *Mask) Clear() {
	for i := range m.data {
		m.data[i] = 0
	}
}
