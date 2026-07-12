// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import "github.com/lumifloat/tinyskia/internal/path"

type CanvasLineCap string

const (
	CanvasLineCapButt   CanvasLineCap = "butt"
	CanvasLineCapRound  CanvasLineCap = "round"
	CanvasLineCapSquare CanvasLineCap = "square"
)

type CanvasLineJoin string

const (
	CanvasLineJoinRound CanvasLineJoin = "round"
	CanvasLineJoinBevel CanvasLineJoin = "bevel"
	CanvasLineJoinMiter CanvasLineJoin = "miter"
)

// SetLineWidth to change the line width.
func (ctx *Context) SetLineWidth(width float64) {
	ctx.lineWidth = width
}

// GetLineWidth returns the current line width.
func (ctx *Context) GetLineWidth() float64 {
	return ctx.lineWidth
}

// SetLineCap to change the line cap style.
func (ctx *Context) SetLineCap(cap CanvasLineCap) {
	ctx.lineCap = cap
}

// GetLineCap returns the current line cap style.
func (ctx *Context) GetLineCap() CanvasLineCap {
	return ctx.lineCap
}

// SetLineJoin to change the line join style.
func (ctx *Context) SetLineJoin(join CanvasLineJoin) {
	ctx.lineJoin = join
}

// GetLineJoin returns the current line join style.
func (ctx *Context) GetLineJoin() CanvasLineJoin {
	return ctx.lineJoin
}

// SetMiterLimit to change the miter limit ratio.
func (ctx *Context) SetMiterLimit(limit float64) {
	ctx.miterLimit = limit
}

// GetMiterLimit returns the current miter limit ratio.
func (ctx *Context) GetMiterLimit() float64 {
	return ctx.miterLimit
}

// SetLineDash to change the line dash pattern.
func (ctx *Context) SetLineDash(segments []float64) {
	if segments == nil {
		segments = []float64{}
	}
	ctx.lineDash = segments
}

// GetLineDash returns the current line dash pattern.
func (ctx *Context) GetLineDash() []float64 {
	return ctx.lineDash
}

// SetLineDashOffset to change the line dash offset.
func (ctx *Context) SetLineDashOffset(offset float64) {
	ctx.lineDashOffset = offset
}

// GetLineDashOffset returns the current line dash offset.
func (ctx *Context) GetLineDashOffset() float64 {
	return ctx.lineDashOffset
}

func cap(c CanvasLineCap) path.LineCap {
	switch c {
	case CanvasLineCapButt:
		return path.LineCapButt
	case CanvasLineCapRound:
		return path.LineCapRound
	case CanvasLineCapSquare:
		return path.LineCapSquare
	default:
		return path.LineCapButt
	}
}

func join(j CanvasLineJoin) path.LineJoin {
	switch j {
	case CanvasLineJoinRound:
		return path.LineJoinRound
	case CanvasLineJoinBevel:
		return path.LineJoinBevel
	case CanvasLineJoinMiter:
		return path.LineJoinMiter
	default:
		return path.LineJoinRound
	}
}
