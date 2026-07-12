// Copyright 2016 Michael Fogleman
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"image"
	"image/color"

	"github.com/lumifloat/tinyskia/internal/core/colorf"
)

// Context is the main drawing context, similar to gg.Context.
// It maintains drawing state and provides a canvas-like API.
type Context struct {
	canvas *Canvas

	mask *image.Alpha

	// CanvasState
	stack       []*Context
	contextLost bool // Tracks if the rendering context was lost

	// CanvasTransform
	matrix *Matrix

	// CanvasCompositing
	globalAlpha              float64
	globalCompositeOperation CompositeOperation

	// CanvasImageSmoothing
	imageSmoothingEnabled bool
	imageSmoothingQuality ImageSmoothingQuality

	// CanvasFillStrokeStyles
	fillStyle   Style
	strokeStyle Style

	// CanvasShadowStyles
	shadowOffsetX float64
	shadowOffsetY float64
	shadowBlur    float64
	shadowColor   color.Color

	// CanvasFilters
	filters []Filter

	// CanvasPathDrawingStyles
	lineWidth      float64
	lineCap        CanvasLineCap
	lineJoin       CanvasLineJoin
	miterLimit     float64
	lineDash       []float64
	lineDashOffset float64

	// CanvasTextDrawingStyles
	lang            string
	font            FontAttr
	textAlign       CanvasTextAlign
	textBaseline    CanvasTextBaseline
	direction       CanvasDirection
	fontKerning     CanvasFontKerning
	fontStretch     CanvasFontStretch
	fontVariantCaps CanvasFontVariantCaps
	textRendering   CanvasTextRendering
	wordSpacing     float64

	// CanvasPath
	path2d *Path2D

	antiAlias       bool
	colorspace      colorf.ColorSpace
	forceHQPipeline bool
}

// SetAntiAlias enables or disables anti-aliasing.
func (ctx *Context) SetAntiAlias(aa bool) {
	ctx.antiAlias = aa
}

func (ctx *Context) SetForceHQPipeline(force bool) {
	ctx.forceHQPipeline = force
}
