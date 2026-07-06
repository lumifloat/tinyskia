// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"image"
	"image/color"

	"github.com/lumifloat/tinyskia/internal/core/painter"
	"github.com/lumifloat/tinyskia/internal/core/scan"
	"github.com/lumifloat/tinyskia/internal/core/shader"
	"github.com/lumifloat/tinyskia/internal/path"
)

// BeginPath resets the current default path.
func (dc *Context) BeginPath() {
	dc.path2d.builder = path.NewPathBuilder()
}

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

	paint := &painter.Paint{
		Shader:          toShader(dc.fillStyle, dc.matrix.transform),
		AntiAlias:       dc.antiAlias,
		BlendMode:       dc.composite,
		Colorspace:      dc.colorspace,
		ForceHQPipeline: dc.forceHQPipeline,
	}
	paint.FillPath(dc.canvas.im, dc.mask, pp, dc.matrix.transform, scan.FillRule(fillRule))
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

	var strokeDash *path.StrokeDash
	if len(dc.lineDash) > 0 {
		dashArray := make([]float32, len(dc.lineDash))
		for i, d := range dc.lineDash {
			dashArray[i] = float32(d)
		}

		if len(dashArray)%2 != 0 {
			doubled := make([]float32, len(dashArray)*2)
			copy(doubled, dashArray)
			copy(doubled[len(dashArray):], dashArray)
			dashArray = doubled
		}

		strokeDash = path.NewStrokeDash(dashArray, float32(dc.lineDashOffset))
	}

	stroke := path.Stroke{
		Width:      float32(dc.lineWidth),
		LineCap:    path.LineCap(dc.lineCap),
		LineJoin:   path.LineJoin(dc.lineJoin),
		MiterLimit: float32(dc.miterLimit),
	}

	paint := &painter.Paint{
		Shader:          toShader(dc.strokeStyle, dc.matrix.transform),
		AntiAlias:       dc.antiAlias,
		BlendMode:       dc.composite,
		Colorspace:      dc.colorspace,
		ForceHQPipeline: dc.forceHQPipeline,
	}

	paint.StrokePath(dc.canvas.im, dc.mask, pp, dc.matrix.transform, stroke, strokeDash)
}

// Clip further constrains the clipping region to the current default path, using the given fill rule to determine what points are in the path.
func (dc *Context) Clip() {
	dc.ClipPathWithFillRule(dc.path2d, FillRuleWinding)
}

// ClipWithFillRule further constrains the clipping region to the current default path, using the given fill rule to determine what points are in the path.
func (dc *Context) ClipWithFillRule(fillRule FillRule) {
	dc.ClipPathWithFillRule(dc.path2d, fillRule)
}

// ClipPath further constrains the clipping region to the given path, using the given fill rule to determine what points are in the path.
func (dc *Context) ClipPath(p *Path2D) {
	dc.ClipPathWithFillRule(p, FillRuleWinding)
}

// ClipPathWithFillRule further constrains the clipping region to the given path, using the given fill rule to determine what points are in the path.
func (dc *Context) ClipPathWithFillRule(p *Path2D, fillRule FillRule) {
	pp := p.builder.Finish()
	if pp == nil {
		return
	}

	paint := &painter.Paint{
		Shader:          shader.NewSolidColor(color.NRGBA{255, 255, 255, 255}),
		AntiAlias:       dc.antiAlias,
		BlendMode:       CompositeOperationSource,
		Colorspace:      dc.colorspace,
		ForceHQPipeline: dc.forceHQPipeline,
	}

	mask := image.NewAlpha(image.Rect(0, 0, dc.canvas.GetWidth(), dc.canvas.GetHeight()))

	paint.ClipPath(mask, dc.mask, pp, dc.matrix.transform, scan.FillRule(fillRule))
	dc.mask = mask
}
