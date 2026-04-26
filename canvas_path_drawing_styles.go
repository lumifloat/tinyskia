// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

type LineCap int

const (
	LineCapButt LineCap = iota
	LineCapRound
	LineCapSquare
)

type LineJoin int

const (
	LineJoinMiter LineJoin = iota
	LineJoinMiterClip
	LineJoinRound
	LineJoinBevel
)

// SetLineWidth to change the line width.
func (dc *Context) SetLineWidth(width float64) {
	dc.lineWidth = width
}

// GetLineWidth returns the current line width.
func (dc *Context) GetLineWidth() float64 {
	return dc.lineWidth
}

// SetLineCap to change the line cap style.
func (dc *Context) SetLineCap(cap LineCap) {
	dc.lineCap = cap
}

// GetLineCap returns the current line cap style.
func (dc *Context) GetLineCap() LineCap {
	return dc.lineCap
}

// SetLineJoin to change the line join style.
func (dc *Context) SetLineJoin(join LineJoin) {
	dc.lineJoin = join
}

// GetLineJoin returns the current line join style.
func (dc *Context) GetLineJoin() LineJoin {
	return dc.lineJoin
}

// SetMiterLimit to change the miter limit ratio.
func (dc *Context) SetMiterLimit(limit float64) {
	dc.miterLimit = limit
}

// GetMiterLimit returns the current miter limit ratio.
func (dc *Context) GetMiterLimit() float64 {
	return dc.miterLimit
}

// SetLineDash to change the line dash pattern.
func (dc *Context) SetLineDash(segments []float64) {
	dc.lineDash = segments
}

// GetLineDash returns the current line dash pattern.
func (dc *Context) GetLineDash() []float64 {
	return dc.lineDash
}

// SetLineDashOffset to change the line dash offset.
func (dc *Context) SetLineDashOffset(offset float64) {
	dc.lineDashOffset = offset
}

// GetLineDashOffset returns the current line dash offset.
func (dc *Context) GetLineDashOffset() float64 {
	return dc.lineDashOffset
}
