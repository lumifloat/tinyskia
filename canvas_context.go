// Copyright 2016 Michael Fogleman
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"image"

	color2 "github.com/lumifloat/tinyskia/internal/core/color"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/text/language"
)

// Context is the main drawing context, similar to gg.Context.
// It maintains drawing state and provides a canvas-like API.
type Context struct {
	canvas *Canvas

	mask *image.Alpha

	// FillStrokeStyles
	fillStyle   style
	strokeStyle style

	// Transform
	matrix *Matrix

	// Composite
	composite CompositeOperation

	// PathDrawingStyles
	lineWidth      float64
	lineCap        LineCap
	lineJoin       LineJoin
	miterLimit     float64
	lineDash       []float64
	lineDashOffset float64

	// Path
	path2d *Path2D

	// TextDrawingStyles
	lang        language.Tag
	font        FontAttr
	buf         sfnt.Buffer
	textAlign   TextAlign
	fontKerning FontKerning
	fontBuffer  sfnt.Buffer

	antiAlias       bool
	colorspace      color2.ColorSpace
	forceHQPipeline bool
	stack           []*Context
	contextLost     bool // Tracks if the rendering context was lost
}

// SetAntiAlias enables or disables anti-aliasing.
func (dc *Context) SetAntiAlias(aa bool) {
	dc.antiAlias = aa
}

func (dc *Context) SetForceHQPipeline(force bool) {
	dc.forceHQPipeline = force
}
