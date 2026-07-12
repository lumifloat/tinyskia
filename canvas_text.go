// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/fontscan"
	"github.com/go-text/typesetting/shaping"
	"github.com/lumifloat/tinyskia/internal/core/painter"
	"github.com/lumifloat/tinyskia/internal/path"
	"github.com/lumifloat/tinyskia/internal/text"
	"golang.org/x/image/math/fixed"
)

// MeasureText measures the given text and returns a TextMetrics object with detailed metrics.
func (ctx *Context) MeasureText(s string) (metrics TextMetrics) {
	runes := []rune(s)
	if len(runes) == 0 {
		return metrics
	}

	ppem := fixed.Int26_6(ctx.font.Size * 64)

	features := []shaping.FontFeature{}
	kerning(ctx.fontKerning, features)
	variant(ctx.fontVariantCaps, features)

	input := shaping.Input{
		Text:         runes,
		RunStart:     0,
		RunEnd:       len(runes),
		Size:         ppem,
		Direction:    direction(ctx.direction),
		FontFeatures: features,
	}

	// TODO ADD PREFERENCE
	query := fontscan.Query{
		Families: ctx.font.Family,
		Aspect: font.Aspect{
			Weight:  font.Weight(ctx.font.Weight),
			Style:   font.Style(ctx.font.Style),
			Stretch: stretch(ctx.fontStretch),
		},
	}

	shapes, width := text.Shape(input, query)

	var fascent, fdescent fixed.Int26_6
	var ascent, descent, left, right fixed.Int26_6

	for _, shape := range shapes {
		for _, glyph := range shape.Glyphs {
			l := width + glyph.XOffset + glyph.XBearing
			r := l + glyph.Width

			if l < left {
				left = l
			}
			if r > right {
				right = r
			}
		}

		if ascent < shape.GlyphBounds.Ascent {
			ascent = shape.GlyphBounds.Ascent
		}
		if descent < shape.GlyphBounds.Descent {
			descent = shape.GlyphBounds.Descent
		}
		if fascent < shape.LineBounds.Ascent {
			fascent = shape.LineBounds.Ascent
		}
		if fdescent > shape.LineBounds.Descent {
			fdescent = shape.LineBounds.Descent
		}
	}

	metrics.Width = float64(width) / 64

	metrics.ActualBoundingBoxLeft = float64(left) / 64
	metrics.ActualBoundingBoxRight = float64(right) / 64

	metrics.ActualBoundingBoxAscent = float64(ascent) / 64
	metrics.ActualBoundingBoxDescent = float64(descent) / 64

	metrics.FontBoundingBoxAscent = float64(fascent) / 64
	metrics.FontBoundingBoxDescent = float64(fdescent) / 64

	return metrics
}

func (ctx *Context) FillText(s string, x, y float64) {
	runes := []rune(s)
	if len(runes) == 0 {
		return
	}

	ppem := fixed.Int26_6(ctx.font.Size * 64)

	features := []shaping.FontFeature{}
	kerning(ctx.fontKerning, features)
	variant(ctx.fontVariantCaps, features)

	input := shaping.Input{
		Text:         runes,
		RunStart:     0,
		RunEnd:       len(runes),
		Size:         ppem,
		Direction:    direction(ctx.direction),
		FontFeatures: features,
	}

	// TODO ADD PREFERENCE
	query := fontscan.Query{
		Families: ctx.font.Family,
		Aspect: font.Aspect{
			Weight:  font.Weight(ctx.font.Weight),
			Style:   font.Style(ctx.font.Style),
			Stretch: stretch(ctx.fontStretch),
		},
	}

	shapes, width := text.Shape(input, query)

	fx := fixed.Int26_6(x*64 + 0.5)
	fy := fixed.Int26_6(y*64 + 0.5)

	switch ctx.textAlign {
	case CanvasTextAlignLeft, CanvasTextAlignStart:
		break
	case CanvasTextAlignRight, CanvasTextAlignEnd:
		fx -= width
	case CanvasTextAlignCenter:
		fx -= width / 2.0
	}

	paint := &painter.Paint{
		Shader:          toShader(ctx.fillStyle, ctx.matrix.transform),
		AntiAlias:       ctx.antiAlias,
		BlendMode:       composite(ctx.globalCompositeOperation),
		Colorspace:      ctx.colorspace,
		ForceHQPipeline: ctx.forceHQPipeline,
	}

	paint.FillTextShapes(ctx.canvas.im, ctx.mask, shapes, fx, fy, ctx.matrix.transform)
}

func (ctx *Context) StrokeText(s string, x, y float64) {
	runes := []rune(s)
	if len(runes) == 0 {
		return
	}

	ppem := fixed.Int26_6(ctx.font.Size * 64)

	features := []shaping.FontFeature{}
	kerning(ctx.fontKerning, features)
	variant(ctx.fontVariantCaps, features)

	input := shaping.Input{
		Text:         runes,
		RunStart:     0,
		RunEnd:       len(runes),
		Size:         ppem,
		Direction:    direction(ctx.direction),
		FontFeatures: features,
	}

	// TODO ADD PREFERENCE
	query := fontscan.Query{
		Families: ctx.font.Family,
		Aspect: font.Aspect{
			Weight:  font.Weight(ctx.font.Weight),
			Style:   font.Style(ctx.font.Style),
			Stretch: stretch(ctx.fontStretch),
		},
	}

	shapes, width := text.Shape(input, query)

	fx := fixed.Int26_6(x*64 + 0.5)
	fy := fixed.Int26_6(y*64 + 0.5)

	switch ctx.textAlign {
	case CanvasTextAlignLeft, CanvasTextAlignStart:
		break
	case CanvasTextAlignRight, CanvasTextAlignEnd:
		fx -= width
	case CanvasTextAlignCenter:
		fx -= width / 2.0
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

	paint.StrokeTextShapes(ctx.canvas.im, ctx.mask, shapes, fx, fy, ctx.matrix.transform, stroke, strokeDash)
}
