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
func (ctx *Context) Save() {
	x := *ctx
	ctx.stack = append(ctx.stack, &x)
}

// Restore pops the top state on the stack, restoring the context to that state.
func (ctx *Context) Restore() {
	if len(ctx.stack) == 0 {
		return
	}

	currentIm := ctx.canvas.im
	currentMask := ctx.mask

	savedState := ctx.stack[len(ctx.stack)-1]
	ctx.stack = ctx.stack[:len(ctx.stack)-1]

	*ctx = *savedState

	ctx.canvas.im = currentIm
	ctx.mask = currentMask
}

// Reset resets the rendering context, which includes the backing buffer, the drawing state stack, path, and styles.
func (ctx *Context) Reset() {
	for i := 0; i < len(ctx.canvas.im.Pix); i++ {
		ctx.canvas.im.Pix[i] = 0
	}
	ctx.mask = nil

	// CanvasState
	ctx.stack = nil
	ctx.contextLost = false

	// CanvasTransform
	ctx.matrix = NewMatrixIdentity()

	// CanvasCompositing
	ctx.globalAlpha = 1.0
	ctx.globalCompositeOperation = CompositeOperationSourceOver

	// CanvasImageSmoothing
	ctx.imageSmoothingEnabled = true
	ctx.imageSmoothingQuality = ImageSmoothingQualityLow

	// CanvasFillStrokeStyles
	ctx.fillStyle = &SolidColor{color: color.RGBA{0, 0, 0, 255}}
	ctx.strokeStyle = &SolidColor{color: color.RGBA{0, 0, 0, 255}}

	// CanvasShadowStyles
	ctx.shadowOffsetX = 0
	ctx.shadowOffsetY = 0
	ctx.shadowBlur = 0
	ctx.shadowColor = color.RGBA{0, 0, 0, 255}

	// CanvasFilters
	ctx.filters = nil

	// CanvasPathDrawingStyles
	ctx.lineWidth = 1
	ctx.lineCap = CanvasLineCapButt
	ctx.lineJoin = CanvasLineJoinMiter
	ctx.miterLimit = 10
	ctx.lineDash = []float64{}
	ctx.lineDashOffset = 0

	// CanvasTextDrawingStyles
	ctx.lang = "inherit"
	ctx.font = FontAttr{
		Family: []string{string(FontGenericSansSerif)},
		Weight: FontWeightNormal,
		Style:  FontStyleNormal,
		Size:   10,
	}
	ctx.textAlign = CanvasTextAlignStart
	ctx.direction = CanvasDirectionLTR
	ctx.fontKerning = CanvasFontKerningAuto
	ctx.fontStretch = CanvasFontStretchNormal
	ctx.fontVariantCaps = CanvasFontVariantCapsNormal
	ctx.textRendering = CanvasTextRenderingAuto
	ctx.wordSpacing = 0

	// CanvasPath
	ctx.path2d = NewPath2D()

	ctx.antiAlias = true
	ctx.colorspace = colorf.ColorSpaceLinear
	ctx.forceHQPipeline = true
}

// IsContextLost returns true if the rendering context was lost. Context loss can occur due to driver crashes, running out of memory, etc. In these cases, the canvas loses its backing storage and takes steps to reset the rendering context to its default state.
func (ctx *Context) IsContextLost() bool {
	return ctx.contextLost
}
