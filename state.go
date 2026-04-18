// Copyright 2016 Michael Fogleman
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"image/color"

	color2 "github.com/lumifloat/tinyskia/internal/core/color"
	"github.com/lumifloat/tinyskia/internal/path"
)

// Save pushes the current state onto the stack.
func (dc *Context) Save() {
	x := *dc
	dc.stack = append(dc.stack, &x)
}

// Restore pops the top state on the stack, restoring the context to that state.
func (dc *Context) Restore() {
	before := *dc
	s := dc.stack
	x, s := s[len(s)-1], s[:len(s)-1]
	*dc = *x
	dc.mask = before.mask
	dc.im = before.im
	// Restore path builder state
	pathData := before.pathBuilder.Finish()
	dc.pathBuilder = path.NewPathBuilder()
	if pathData != nil {
		dc.pathBuilder.PushPath(pathData)
	}
	dc.stack = s
}

// Reset resets the rendering context, which includes the backing buffer, the drawing state stack, path, and styles.
func (dc *Context) Reset() {
	for i := 0; i < len(dc.im.Pix); i++ {
		dc.im.Pix[i] = 0
	}

	dc.stack = nil
	dc.pathBuilder = path.NewPathBuilder()
	dc.dashes = nil
	dc.dashOffset = 0
	dc.lineWidth = 1
	dc.lineCap = LineCapRound
	dc.lineJoin = LineJoinRound
	dc.fillRule = FillRuleWinding
	dc.transform = path.NewTransformDefault()
	dc.blendMode = BlendModeSourceOver
	dc.antiAlias = true
	dc.colorspace = color2.ColorSpaceLinear
	dc.forceHQPipeline = true
	dc.color = color.Transparent
	dc.fillStyle = NewSolidPattern(color.Transparent)
	dc.strokeStyle = NewSolidPattern(color.Transparent)
	dc.mask = nil
	dc.contextLost = false
}

// IsContextLost returns true if the rendering context was lost. Context loss can occur due to driver crashes, running out of memory, etc. In these cases, the canvas loses its backing storage and takes steps to reset the rendering context to its default state.
func (dc *Context) IsContextLost() bool {
	return dc.contextLost
}
