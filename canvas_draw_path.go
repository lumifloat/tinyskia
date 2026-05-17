// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"image"

	color2 "github.com/lumifloat/tinyskia/internal/core/color"
	"github.com/lumifloat/tinyskia/internal/core/scan"
	"github.com/lumifloat/tinyskia/internal/core/shader"
	"github.com/lumifloat/tinyskia/internal/path"
)

// Fills the subpaths of the current default path with the current fill style.
func (dc *Context) Fill() {
	dc.FillPathWithFillRule(dc.path2d, FillRuleWinding)
}

// FillWithFillRule fills the subpaths of the current default path with the current fill style, obeying the given fill rule.
func (dc *Context) FillWithFillRule(fillRule FillRule) {
	dc.FillPathWithFillRule(dc.path2d, fillRule)
}

// FillPath fills the subpaths of the given path with the current fill style.
func (dc *Context) FillPath(p *Path2D) {
	dc.FillPathWithFillRule(p, FillRuleWinding)
}

// FillPathWithFillRule fills the subpaths of the given path with the current fill style, obeying the given fill rule.
func (dc *Context) FillPathWithFillRule(p *Path2D, fillRule FillRule) {
	pp := p.builder.Finish()
	if pp == nil {
		return
	}

	var tp *path.Path
	if !dc.matrix.transform.IsIdentity() {
		tp = pp.Transform(dc.matrix.transform)
	} else {
		tp = pp
	}

	paint := &Paint{
		Shader:          toShader(dc.fillStyle, dc.matrix.transform),
		AntiAlias:       dc.antiAlias,
		BlendMode:       dc.composite,
		Colorspace:      dc.colorspace,
		ForceHQPipeline: dc.forceHQPipeline,
	}
	var maskData []uint8
	if dc.mask != nil {
		maskData = dc.mask.Pix
	}
	blitter := paint.blitter(dc.canvas.im.Pix, maskData, dc.canvas.GetWidth(), dc.canvas.GetHeight())
	screen, _ := path.NewScreenIntRectFromXYWH(0, 0, uint32(dc.canvas.GetWidth()), uint32(dc.canvas.GetHeight()))
	if dc.antiAlias {
		scan.FillPathAA(tp, int(fillRule), screen, blitter)
	} else {
		scan.FillPath(tp, int(fillRule), screen, blitter)
	}
}

// Stroke the subpaths of the current default path with the current stroke style.
func (dc *Context) Stroke() {
	dc.StrokePath(dc.path2d)
}

// StrokePath the subpaths of the given path with the current stroke style.
func (dc *Context) StrokePath(p *Path2D) {
	pp := p.builder.Finish()
	if pp == nil {
		return
	}

	var tp *path.Path
	var lineWidth float32
	if !dc.matrix.transform.IsIdentity() {
		tp = pp.Transform(dc.matrix.transform)
		resScale := path.ComputeResolutionScale(dc.matrix.transform)
		lineWidth = float32(dc.lineWidth) * resScale
	} else {
		tp = pp
		lineWidth = float32(dc.lineWidth)
	}

	if len(dc.lineDash) > 0 {
		dashArray := make([]float32, len(dc.lineDash))
		for i, d := range dc.lineDash {
			dashArray[i] = float32(d)
		}
		strokeDash := path.NewStrokeDash(dashArray, float32(dc.lineDashOffset))
		if strokeDash != nil {
			dashedPath := tp.Dash(strokeDash, path.ComputeResolutionScale(dc.matrix.transform))
			if dashedPath != nil {
				tp = dashedPath
			}
		}
	}

	stroke := path.Stroke{
		Width:      lineWidth,
		LineCap:    path.LineCap(dc.lineCap),
		LineJoin:   path.LineJoin(dc.lineJoin),
		MiterLimit: 4.0, // Default miter limit
	}

	paint := &Paint{
		Shader:          toShader(dc.strokeStyle, dc.matrix.transform),
		AntiAlias:       dc.antiAlias,
		BlendMode:       dc.composite,
		Colorspace:      dc.colorspace,
		ForceHQPipeline: dc.forceHQPipeline,
	}
	var maskData []uint8
	if dc.mask != nil {
		maskData = dc.mask.Pix
	}
	blitter := paint.blitter(dc.canvas.im.Pix, maskData, dc.canvas.GetWidth(), dc.canvas.GetHeight())
	screen, _ := path.NewScreenIntRectFromXYWH(0, 0, uint32(dc.canvas.GetWidth()), uint32(dc.canvas.GetHeight()))

	resScale := path.ComputeResolutionScale(dc.matrix.transform)

	stroker := path.NewPathStroker()
	strokedPath := stroker.Stroke(tp, stroke, resScale)
	if strokedPath == nil {
		return
	}

	if dc.antiAlias {
		scan.FillPathAA(strokedPath, int(FillRuleWinding), screen, blitter)
	} else {
		scan.FillPath(strokedPath, int(FillRuleWinding), screen, blitter)
	}
}

// Further constrains the clipping region to the current default path, using the given fill rule to determine what points are in the path.
func (dc *Context) Clip() {
	dc.ClipPath(dc.path2d)
}

// Further constrains the clipping region to the given path, using the given fill rule to determine what points are in the path.
func (dc *Context) ClipPath(p *Path2D) {
	pp := p.builder.Finish()
	if pp == nil {
		return
	}

	var tp *path.Path
	if !dc.matrix.transform.IsIdentity() {
		tp = pp.Transform(dc.matrix.transform)
	} else {
		tp = pp
	}

	// Create a temporary alpha mask for the clip path
	width := dc.canvas.GetWidth()
	height := dc.canvas.GetHeight()
	clipMask := image.NewAlpha(image.Rect(0, 0, width, height))

	// Render path to a temporary RGBA image first (blitter expects RGBA)
	tempRGBA := image.NewRGBA(image.Rect(0, 0, width, height))
	paint := &Paint{
		Shader:          shader.NewSolidColor(color2.ColorFromRGBA8(255, 255, 255, 255)),
		AntiAlias:       dc.antiAlias,
		BlendMode:       CompositeOperationSource,
		Colorspace:      dc.colorspace,
		ForceHQPipeline: dc.forceHQPipeline,
	}
	blitter := paint.blitter(tempRGBA.Pix, nil, width, height)
	screen, _ := path.NewScreenIntRectFromXYWH(0, 0, uint32(width), uint32(height))
	scan.FillPathAA(tp, int(FillRuleWinding), screen, blitter)

	// Extract alpha channel from RGBA to Alpha mask
	for i := 0; i < width*height; i++ {
		clipMask.Pix[i] = tempRGBA.Pix[i*4+3] // Alpha channel
	}

	// Intersect with existing mask (take minimum of alpha values)
	if dc.mask == nil {
		dc.mask = clipMask
	} else {
		// Create new mask by intersecting old mask and clip mask
		// Intersection = min(old_alpha, clip_alpha) for each pixel
		mask := image.NewAlpha(image.Rect(0, 0, width, height))
		for i := range mask.Pix {
			clipAlpha := clipMask.Pix[i]
			oldAlpha := dc.mask.Pix[i]
			if clipAlpha < oldAlpha {
				mask.Pix[i] = clipAlpha
			} else {
				mask.Pix[i] = oldAlpha
			}
		}
		dc.mask = mask
	}
}
