// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/font/opentype"
	"github.com/go-text/typesetting/fontscan"
	"github.com/go-text/typesetting/shaping"
	"golang.org/x/image/math/fixed"
)

// MeasureText measures the given text and returns a TextMetrics object with detailed metrics.
func (ctx *Context) MeasureText(s string) (metrics TextMetrics) {
	runes := []rune(s)
	if len(runes) == 0 {
		return metrics
	}

	ppem := fixed.Int26_6(ctx.font.Size * 64)
	input := shaping.Input{
		Text:     runes,
		RunStart: 0,
		RunEnd:   len(runes),
		Size:     ppem,
		FontFeatures: []shaping.FontFeature{
			{
				Tag:   opentype.NewTag('k', 'e', 'r', 'n'),
				Value: uint32(ctx.fontKerning),
			},
		},
	}

	var segmenter shaping.Segmenter
	flock.Lock()
	fonts.SetQuery(fontscan.Query{
		Families: ctx.font.Family,
		Aspect: font.Aspect{
			Weight: font.Weight(ctx.font.Weight),
			Style:  font.StyleNormal,
		},
	})
	outputs := segmenter.Split(input, fonts)
	flock.Unlock()

	shaper := shaping.HarfbuzzShaper{}

	var width fixed.Int26_6
	var fascent, fdescent fixed.Int26_6
	var ascent, descent, left, right fixed.Int26_6

	for _, output := range outputs {
		shape := shaper.Shape(output)

		for _, glyph := range shape.Glyphs {
			l := width + glyph.XOffset + glyph.XBearing
			r := l + glyph.Width

			if l < left {
				left = l
			}
			if r > right {
				right = r
			}
			width += glyph.Advance
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

func (dc *Context) FillText(s string, x, y float64) {
	dc.drawText(s, x, y, false)
}

func (dc *Context) StrokeText(s string, x, y float64) {
	dc.drawText(s, x, y, true)
}

func (dc *Context) drawText(s string, x, y float64, stroke bool) {
	runes := []rune(s)
	if len(runes) == 0 {
		return
	}

	ppem := fixed.Int26_6(dc.font.Size * 64)
	input := shaping.Input{
		Text:     runes,
		RunStart: 0,
		RunEnd:   len(runes),
		Size:     ppem,
		FontFeatures: []shaping.FontFeature{
			{
				Tag:   opentype.NewTag('k', 'e', 'r', 'n'),
				Value: uint32(dc.fontKerning),
			},
		},
	}

	var segmenter shaping.Segmenter
	flock.Lock()
	fonts.SetQuery(fontscan.Query{
		Families: dc.font.Family,
		Aspect: font.Aspect{
			Weight: font.Weight(dc.font.Weight),
			Style:  font.Style(dc.font.Style),
		},
	})
	outputs := segmenter.Split(input, fonts)
	flock.Unlock()

	fx := float32(x)
	fy := float32(y)

	shaper := shaping.HarfbuzzShaper{}
	shapes := []shaping.Output{}
	width := float32(0)
	for _, output := range outputs {
		shape := shaper.Shape(output)
		shapes = append(shapes, shape)
		width += float32(shape.Advance) / 64
	}
	switch dc.textAlign {
	case TextAlignLeft, TextAlignStart:
		break
	case TextAlignRight, TextAlignEnd:
		fx -= width
	case TextAlignCenter:
		fx -= width / 2.0
	}

	for _, shape := range shapes {
		upem := shape.Face.Upem()
		scale := float32(ppem) / float32(upem) / 64.0
		for _, glyph := range shape.Glyphs {
			data := shape.Face.GlyphData(glyph.GlyphID)
			switch d := data.(type) {
			case font.GlyphOutline:
				path2d := outline(d, scale, fx, fy)
				if stroke {
					dc.StrokePath(path2d)
				} else {
					dc.FillPath(path2d)
				}
			}

			fx += float32(glyph.Advance) / 64
		}
	}

}

func outline(outline font.GlyphOutline, scale float32, x, y float32) *Path2D {
	var path2d = NewPath2D()
	var hasPath = false
	for _, s := range outline.Segments {
		switch s.Op {
		case opentype.SegmentOpMoveTo:
			if hasPath {
				path2d.builder.Close()
			}
			path2d.builder.MoveTo(
				s.Args[0].X*scale+x,
				-s.Args[0].Y*scale+y,
			)
			hasPath = true
		case opentype.SegmentOpLineTo:
			path2d.builder.LineTo(
				s.Args[0].X*scale+x,
				-s.Args[0].Y*scale+y,
			)
		case opentype.SegmentOpQuadTo:
			path2d.builder.QuadTo(
				s.Args[0].X*scale+x,
				-s.Args[0].Y*scale+y,
				s.Args[1].X*scale+x,
				-s.Args[1].Y*scale+y,
			)
		case opentype.SegmentOpCubeTo:
			path2d.builder.CubicTo(
				s.Args[0].X*scale+x,
				-s.Args[0].Y*scale+y,
				s.Args[1].X*scale+x,
				-s.Args[1].Y*scale+y,
				s.Args[2].X*scale+x,
				-s.Args[2].Y*scale+y,
			)
		}
	}
	if hasPath {
		path2d.builder.Close()
	}
	return path2d
}
