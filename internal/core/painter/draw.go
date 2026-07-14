// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package painter

import (
	"image"
	"image/color"

	"github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/font/opentype/tables"
	"github.com/go-text/typesetting/shaping"
	"github.com/lumifloat/tinyskia/internal/core/scan"
	"github.com/lumifloat/tinyskia/internal/core/shader"
	"github.com/lumifloat/tinyskia/internal/path"
	"github.com/lumifloat/tinyskia/internal/text"
	"golang.org/x/image/math/fixed"
)

func (p *Paint) FillPath(dst *image.RGBA, mask *image.Alpha, sp *path.Path, transform path.Transform, fillRule scan.FillRule) {
	if sp == nil || sp.Len() == 0 {
		return
	}
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
	if sp == nil || sp.Len() == 0 {
		return
	}
	var tp *path.Path
	if !transform.IsIdentity() {
		tp = sp.Transform(transform)
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
	if sp == nil || sp.Len() == 0 {
		return
	}
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

func (p *Paint) FillTextShapes(dst *image.RGBA, mask *image.Alpha, shapes []shaping.Output, x, y fixed.Int26_6, transform path.Transform) {
	for _, shape := range shapes {
		scale := float32(shape.Size) / float32(shape.Face.Upem()) / 64
		for _, glyph := range shape.Glyphs {
			data := shape.Face.GlyphData(glyph.GlyphID)
			switch data := data.(type) {
			case font.GlyphBitmap:
				// TODO
			case font.GlyphColor:
				DrawTextPaintTable(p, dst, mask, shape.Face, data.Paint, x, y, scale, transform)
			case font.GlyphOutline:
				tf := transform.PreTranslate(float32(x)/64, float32(y)/64).PreScale(scale, scale)
				p.FillPath(dst, mask, text.Outline(data), tf, scan.FillRuleWinding)
			}
			x += glyph.Advance
		}
	}
}

func (p *Paint) StrokeTextShapes(dst *image.RGBA, mask *image.Alpha, shapes []shaping.Output, x, y fixed.Int26_6, transform path.Transform, stroke path.Stroke, dash *path.StrokeDash) {
	for _, shape := range shapes {
		scale := float32(shape.Size) / float32(shape.Face.Upem()) / 64
		for _, glyph := range shape.Glyphs {
			data := shape.Face.GlyphData(glyph.GlyphID)
			switch data := data.(type) {
			case font.GlyphBitmap:
				// TODO
			case font.GlyphColor:
				p = p.Copy()
				DrawTextPaintTable(p, dst, mask, shape.Face, data.Paint, x, y, scale, transform)
			case font.GlyphOutline:
				tf := transform.PreTranslate(float32(x)/64, float32(y)/64).PreScale(scale, scale)
				p.StrokePath(dst, mask, text.Outline(data), tf, stroke, dash)
			}
			x += glyph.Advance
		}
	}
}

func DrawTextPaintTable(p *Paint, dst *image.RGBA, mask *image.Alpha, face *font.Face, table tables.PaintTable, x, y fixed.Int26_6, scale float32, transform path.Transform) {
	switch t := table.(type) {
	case tables.PaintColrLayers:
		if face.Font.COLR == nil {
			return
		}
		list := face.Font.COLR.LayerList
		layers, err := list.Resolve(t)
		if err != nil {
			return
		}
		for _, layer := range layers {
			DrawTextPaintTable(p, dst, mask, face, layer, x, y, scale, transform)
		}
	case tables.PaintColrLayersResolved:
		for _, layer := range t {
			outline, ok := face.GlyphData(font.GID(layer.GlyphID)).(font.GlyphOutline)
			if !ok {
				continue
			}
			cpal := face.Font.CPAL
			if len(cpal) != 0 && len(cpal[0]) >= int(layer.PaletteIndex) {
				c := cpal[0][layer.PaletteIndex]
				p.Shader = shader.NewSolidColor(color.RGBA{c.Red, c.Green, c.Blue, c.Alpha})
			}
			tf := transform.PreTranslate(float32(x)/64, float32(y)/64).PreScale(scale, scale)
			p.FillPath(dst, mask, text.Outline(outline), tf, scan.FillRuleWinding)
		}
	case tables.PaintGlyph:
		outline, ok := face.GlyphData(font.GID(t.GlyphID)).(font.GlyphOutline)
		if !ok {
			return
		}
		tf := transform.PreTranslate(float32(x)/64, float32(y)/64).PreScale(scale, scale)
		p.FillPath(dst, mask, text.Outline(outline), tf, scan.FillRuleWinding)
	case tables.PaintSolid:
		cpal := face.Font.CPAL
		if len(cpal) != 0 && len(cpal[0]) >= int(t.PaletteIndex) {
			c := cpal[0][t.PaletteIndex]
			p.Shader = shader.NewSolidColor(color.RGBA{c.Red, c.Green, c.Blue, c.Alpha})
		}
	case tables.PaintVarSolid:
		cpal := face.Font.CPAL
		if len(cpal) != 0 && len(cpal[0]) >= int(t.PaletteIndex) {
			c := cpal[0][t.PaletteIndex]
			p.Shader = shader.NewSolidColor(color.RGBA{c.Red, c.Green, c.Blue, c.Alpha})
		}
	case tables.PaintTransform:
		transform = transform.PreConcat(path.NewTransform(
			float32(t.Transform.Xx)/65536.0,
			float32(t.Transform.Yx)/65536.0,
			float32(t.Transform.Xy)/65536.0,
			float32(t.Transform.Yy)/65536.0,
			float32(t.Transform.Dx)/65536.0,
			float32(t.Transform.Dy)/65536.0,
		))
		DrawTextPaintTable(p, dst, mask, face, t.Paint, x, y, scale, transform)
	case tables.PaintTranslate:
		transform = transform.PreConcat(path.NewTransformFromTranslate(
			float32(t.Dx)/65536.0,
			float32(t.Dy)/65536.0,
		))
		DrawTextPaintTable(p, dst, mask, face, t.Paint, x, y, scale, transform)
	default:
		return
	}
}
