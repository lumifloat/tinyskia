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

type CanvasFillRule string

const (
	CanvasFillRuleNonzero CanvasFillRule = "nonzero"
	CanvasFillRuleEvenodd CanvasFillRule = "evenodd"
)

// BeginPath resets the current default path.
func (ctx *Context) BeginPath() {
	ctx.path2d.builder = path.NewPathBuilder()
}

// Fills the subpaths of the current default path with the current fill style.
func (ctx *Context) Fill() {
	ctx.FillPathWithFillRule(ctx.path2d, CanvasFillRuleNonzero)
}

// FillWithFillRule fills the subpaths of the current default path with the current fill style, obeying the given fill rule.
func (ctx *Context) FillWithFillRule(fillRule CanvasFillRule) {
	ctx.FillPathWithFillRule(ctx.path2d, fillRule)
}

// FillPath fills the subpaths of the given path with the current fill style.
func (ctx *Context) FillPath(p *Path2D) {
	ctx.FillPathWithFillRule(p, CanvasFillRuleNonzero)
}

// FillPathWithFillRule fills the subpaths of the given path with the current fill style, obeying the given fill rule.
func (ctx *Context) FillPathWithFillRule(p *Path2D, fillRule CanvasFillRule) {
	pp := p.builder.Finish()
	if pp == nil {
		return
	}

	paint := &painter.Paint{
		Shader:          toShader(ctx.fillStyle, ctx.matrix.transform),
		AntiAlias:       ctx.antiAlias,
		BlendMode:       composite(ctx.globalCompositeOperation),
		Colorspace:      ctx.colorspace,
		ForceHQPipeline: ctx.forceHQPipeline,
	}
	paint.FillPath(ctx.canvas.im, ctx.mask, pp, ctx.matrix.transform, rule(fillRule))
}

// Stroke the subpaths of the current default path with the current stroke style.
func (ctx *Context) Stroke() {
	ctx.StrokePath(ctx.path2d)
}

// StrokePath the subpaths of the given path with the current stroke style.
func (ctx *Context) StrokePath(p *Path2D) {
	pp := p.builder.Finish()
	if pp == nil {
		return
	}

	var strokeDash *path.StrokeDash
	if len(ctx.lineDash) > 0 {
		dashArray := make([]float32, len(ctx.lineDash))
		for i, d := range ctx.lineDash {
			dashArray[i] = float32(d)
		}

		if len(dashArray)%2 != 0 {
			doubled := make([]float32, len(dashArray)*2)
			copy(doubled, dashArray)
			copy(doubled[len(dashArray):], dashArray)
			dashArray = doubled
		}

		strokeDash = path.NewStrokeDash(dashArray, float32(ctx.lineDashOffset))
	}

	stroke := path.Stroke{
		Width:      float32(ctx.lineWidth),
		LineCap:    cap(ctx.lineCap),
		LineJoin:   join(ctx.lineJoin),
		MiterLimit: float32(ctx.miterLimit),
	}

	paint := &painter.Paint{
		Shader:          toShader(ctx.strokeStyle, ctx.matrix.transform),
		AntiAlias:       ctx.antiAlias,
		BlendMode:       composite(ctx.globalCompositeOperation),
		Colorspace:      ctx.colorspace,
		ForceHQPipeline: ctx.forceHQPipeline,
	}

	paint.StrokePath(ctx.canvas.im, ctx.mask, pp, ctx.matrix.transform, stroke, strokeDash)
}

// Clip further constrains the clipping region to the current default path, using the given fill rule to determine what points are in the path.
func (ctx *Context) Clip() {
	ctx.ClipPathWithFillRule(ctx.path2d, CanvasFillRuleNonzero)
}

// ClipWithFillRule further constrains the clipping region to the current default path, using the given fill rule to determine what points are in the path.
func (ctx *Context) ClipWithFillRule(fillRule CanvasFillRule) {
	ctx.ClipPathWithFillRule(ctx.path2d, fillRule)
}

// ClipPath further constrains the clipping region to the given path, using the given fill rule to determine what points are in the path.
func (ctx *Context) ClipPath(p *Path2D) {
	ctx.ClipPathWithFillRule(p, CanvasFillRuleNonzero)
}

// ClipPathWithFillRule further constrains the clipping region to the given path, using the given fill rule to determine what points are in the path.
func (ctx *Context) ClipPathWithFillRule(p *Path2D, fillRule CanvasFillRule) {
	pp := p.builder.Finish()
	if pp == nil {
		return
	}

	paint := &painter.Paint{
		Shader:          shader.NewSolidColor(color.NRGBA{255, 255, 255, 255}),
		AntiAlias:       ctx.antiAlias,
		BlendMode:       composite(ctx.globalCompositeOperation),
		Colorspace:      ctx.colorspace,
		ForceHQPipeline: ctx.forceHQPipeline,
	}

	mask := image.NewAlpha(image.Rect(0, 0, ctx.canvas.GetWidth(), ctx.canvas.GetHeight()))

	paint.ClipPath(mask, ctx.mask, pp, ctx.matrix.transform, rule(fillRule))
	ctx.mask = mask
}

func rule(r CanvasFillRule) scan.FillRule {
	switch r {
	case CanvasFillRuleNonzero:
		return scan.FillRuleWinding
	case CanvasFillRuleEvenodd:
		return scan.FillRuleEvenOdd
	default:
		return scan.FillRuleWinding
	}
}
