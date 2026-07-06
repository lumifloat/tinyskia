// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package painter

import (
	"image"

	"github.com/go-text/typesetting/font"
	"github.com/lumifloat/tinyskia/internal/core/scan"
	"github.com/lumifloat/tinyskia/internal/path"
)

func (p *Paint) FillPath(dst *image.RGBA, mask *image.Alpha, sp *path.Path, transform path.Transform, fillRule scan.FillRule) {
	var tp *path.Path
	if !transform.IsIdentity() {
		tp = sp.Transform(transform)
	} else {
		tp = sp
	}

	blitter := p.blitter(dst, mask)
	if blitter == nil {
		return
	}
	screen, _ := path.NewScreenIntRectFromXYWH(0, 0, uint32(dst.Rect.Dx()), uint32(dst.Rect.Dy()))
	if p.AntiAlias {
		scan.FillPathAA(tp, int(fillRule), screen, blitter)
	} else {
		scan.FillPath(tp, int(fillRule), screen, blitter)
	}
}

func (p *Paint) StrokePath(dst *image.RGBA, mask *image.Alpha, sp *path.Path, transform path.Transform, stroke path.Stroke, dash *path.StrokeDash) {
	var tp *path.Path
	if !transform.IsIdentity() {
		tp = sp.Transform(transform)
		resScale := path.ComputeResolutionScale(transform)
		stroke.Width = float32(stroke.Width) * resScale
	} else {
		tp = sp
	}

	if dash != nil {
		dashedPath := tp.Dash(dash, path.ComputeResolutionScale(transform))
		if dashedPath != nil {
			tp = dashedPath
		}
	}

	blitter := p.blitter(dst, mask)
	if blitter == nil {
		return
	}
	screen, _ := path.NewScreenIntRectFromXYWH(0, 0, uint32(dst.Rect.Dx()), uint32(dst.Rect.Dy()))
	scale := path.ComputeResolutionScale(transform)
	stroker := path.NewPathStroker()
	ssp := stroker.Stroke(tp, stroke, scale)
	if p.AntiAlias {
		scan.FillPathAA(ssp, int(scan.FillRuleWinding), screen, blitter)
	} else {
		scan.FillPath(ssp, int(scan.FillRuleWinding), screen, blitter)
	}
}

func (p *Paint) ClipPath(dst *image.Alpha, mask *image.Alpha, sp *path.Path, transform path.Transform, fillRule scan.FillRule) {
	var tp *path.Path
	if !transform.IsIdentity() {
		tp = sp.Transform(transform)
	} else {
		tp = sp
	}

	blitter := p.NewMaskBlitter(dst, mask)
	if blitter == nil {
		return
	}
	screen, _ := path.NewScreenIntRectFromXYWH(0, 0, uint32(dst.Rect.Dx()), uint32(dst.Rect.Dy()))
	scan.FillPath(tp, int(fillRule), screen, blitter)
}

// func (p *Paint) FillGlyph(dst *image.RGBA, mask *image.Alpha, glyph shaping.Glyph, scale float32, transform path.Transform, fillRule scan.FillRule) {
// 	text.Outline(glyph, scale, , fillRule)
// }

func (p *Paint) FillTextOutline(dst *image.RGBA, mask *image.Alpha, outline font.GlyphOutline, scale float32, transform path.Transform, fillRule scan.FillRule) {
}
