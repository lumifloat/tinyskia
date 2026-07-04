// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"image/color"

	"github.com/lumifloat/tinyskia/internal/core/colorf"
)

// Save pushes the current state onto the stack.
func (dc *Context) Save() {
	x := *dc
	dc.stack = append(dc.stack, &x)
}

// Restore pops the top state on the stack, restoring the context to that state.
func (dc *Context) Restore() {
	if len(dc.stack) == 0 {
		return
	}

	currentIm := dc.canvas.im
	currentMask := dc.mask

	savedState := dc.stack[len(dc.stack)-1]
	dc.stack = dc.stack[:len(dc.stack)-1]

	*dc = *savedState

	dc.canvas.im = currentIm
	dc.mask = currentMask
}

// Reset resets the rendering context, which includes the backing buffer, the drawing state stack, path, and styles.
func (dc *Context) Reset() {
	for i := 0; i < len(dc.canvas.im.Pix); i++ {
		dc.canvas.im.Pix[i] = 0
	}

	dc.stack = nil
	dc.path2d = NewPath2D()
	dc.lineDash = nil
	dc.lineDashOffset = 0
	dc.lineWidth = 1
	dc.lineCap = LineCapRound
	dc.lineJoin = LineJoinRound
	dc.matrix = NewMatrixIdentity()
	dc.composite = CompositeOperationSourceOver
	dc.antiAlias = true
	dc.colorspace = colorf.ColorSpaceLinear
	dc.forceHQPipeline = true
	dc.fillStyle = dc.CreateSolidColor(color.Transparent)
	dc.strokeStyle = dc.CreateSolidColor(color.Transparent)
	dc.mask = nil
	dc.contextLost = false
}

// IsContextLost returns true if the rendering context was lost. Context loss can occur due to driver crashes, running out of memory, etc. In these cases, the canvas loses its backing storage and takes steps to reset the rendering context to its default state.
func (dc *Context) IsContextLost() bool {
	return dc.contextLost
}
