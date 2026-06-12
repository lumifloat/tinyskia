// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"math"

	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

// MeasureText measures the given text and returns a TextMetrics object with detailed metrics.
func (dc *Context) MeasureText(s string) TextMetrics {
	var metrics TextMetrics

	if s == "" {
		return metrics
	}

	ppem := fixed.Int26_6(dc.font.Size * 64)

	var advance fixed.Int26_6
	var prev sfnt.GlyphIndex

	minX, maxX := fixed.Int26_6(1<<30), fixed.Int26_6(-(1 << 30))
	minY, maxY := fixed.Int26_6(1<<30), fixed.Int26_6(-(1 << 30))

	var maxFontAscent, maxFontDescent float64

	for _, r := range s {
		ft, curr := dc.glyph(r)
		if curr == 0 {
			prev = 0
			continue
		}

		hMetrics, err := ft.Metrics(&dc.buf, ppem, font.HintingFull)
		if err == nil {
			fontAscent := math.Abs(float64(hMetrics.Ascent) / 64.0)
			fontDescent := math.Abs(float64(hMetrics.Descent) / 64.0)

			maxFontAscent = math.Max(maxFontAscent, fontAscent)
			maxFontDescent = math.Max(maxFontDescent, fontDescent)
		}

		if dc.fontKerning != FontKerningNone && prev != 0 {
			kern, err := ft.Kern(&dc.buf, prev, curr, ppem, font.HintingFull)
			if err == nil {
				advance += kern
			}
		}
		prev = curr

		b, a, err := ft.GlyphBounds(&dc.buf, curr, ppem, font.HintingFull)
		if err != nil {
			continue
		}
		if b.Empty() {
			advance += a
			continue
		}

		if advance+b.Min.X < minX {
			minX = advance + b.Min.X
		}
		if advance+b.Max.X > maxX {
			maxX = advance + b.Max.X
		}
		if b.Min.Y < minY {
			minY = b.Min.Y
		}
		if b.Max.Y > maxY {
			maxY = b.Max.Y
		}

		advance += a
	}

	metrics.Width = float64(advance) / 64.0

	metrics.FontBoundingBoxAscent = maxFontAscent
	metrics.FontBoundingBoxDescent = maxFontDescent

	if minX < fixed.Int26_6(1<<30) {
		metrics.ActualBoundingBoxLeft = -float64(minX) / 64.0
		metrics.ActualBoundingBoxRight = float64(maxX) / 64.0
		metrics.ActualBoundingBoxAscent = -float64(minY) / 64.0
		metrics.ActualBoundingBoxDescent = float64(maxY) / 64.0
	}

	return metrics
}

func (dc *Context) FillText(s string, x, y float64) {
	dc.drawText(s, x, y, false)
}

func (dc *Context) StrokeText(s string, x, y float64) {
	dc.drawText(s, x, y, true)
}

func (dc *Context) drawText(s string, x, y float64, stroke bool) {
	ppem := fixed.Int26_6(dc.font.Size * 64)

	metrics := dc.MeasureText(s)
	textWidth := metrics.Width

	var offsetX float64
	switch dc.textAlign {
	case TextAlignLeft, TextAlignStart:
		offsetX = 0
	case TextAlignRight, TextAlignEnd:
		offsetX = -textWidth
	case TextAlignCenter:
		offsetX = -textWidth / 2.0
	default:
		offsetX = 0
	}

	x = x + offsetX

	fx := fixed.I(int(x))
	fy := fixed.I(int(y))
	dot := fixed.Point26_6{X: fx, Y: fy}

	var prev sfnt.GlyphIndex

	for _, r := range s {
		path2d := NewPath2D()
		ft, curr := dc.glyph(r)
		if curr == 0 {
			continue
		}

		if dc.fontKerning != FontKerningNone && prev != 0 {
			kern, err := ft.Kern(&dc.buf, prev, curr, ppem, font.HintingFull)
			if err == nil {
				dot.X += kern
			}
		}
		prev = curr

		advance, err := ft.GlyphAdvance(&dc.buf, curr, ppem, font.HintingFull)
		if err != nil {
			continue
		}

		segments, err := ft.LoadGlyph(&dc.buf, curr, ppem, nil)
		if err != nil {
			continue
		}

		var hasPath bool = false
		for _, seg := range segments {
			switch seg.Op {
			case sfnt.SegmentOpMoveTo:
				if hasPath {
					path2d.builder.Close()
				}
				path2d.builder.MoveTo(
					float32(seg.Args[0].X+dot.X)/64.0,
					float32(seg.Args[0].Y+dot.Y)/64.0,
				)
				hasPath = true
			case sfnt.SegmentOpLineTo:
				path2d.builder.LineTo(
					float32(seg.Args[0].X+dot.X)/64.0,
					float32(seg.Args[0].Y+dot.Y)/64.0,
				)
			case sfnt.SegmentOpQuadTo:
				path2d.builder.QuadTo(
					float32(seg.Args[0].X+dot.X)/64.0,
					float32(seg.Args[0].Y+dot.Y)/64.0,
					float32(seg.Args[1].X+dot.X)/64.0,
					float32(seg.Args[1].Y+dot.Y)/64.0,
				)
			case sfnt.SegmentOpCubeTo:
				path2d.builder.CubicTo(
					float32(seg.Args[0].X+dot.X)/64.0,
					float32(seg.Args[0].Y+dot.Y)/64.0,
					float32(seg.Args[1].X+dot.X)/64.0,
					float32(seg.Args[1].Y+dot.Y)/64.0,
					float32(seg.Args[2].X+dot.X)/64.0,
					float32(seg.Args[2].Y+dot.Y)/64.0,
				)
			}
		}
		if hasPath {
			path2d.builder.Close()
		}

		if stroke {
			dc.StrokePath(path2d)
		} else {
			dc.FillPath(path2d)
		}
		dot.X += advance
	}
}

func (dc *Context) glyph(r rune) (*sfnt.Font, sfnt.GlyphIndex) {
	for i := range dc.font.Family {
		chain := dc.fmatch0(dc.font.Family[i])
		for i := range chain {
			ft, err := loadFont(chain[i])
			if err != nil {
				continue
			}
			curr, err := ft.GlyphIndex(&dc.buf, r)
			if err == nil && curr != 0 {
				return ft, curr
			}
		}
	}
	chain := dc.fmatch1(r)
	for i := range chain {
		ft, err := loadFont(chain[i])
		if err != nil {
			continue
		}
		curr, err := ft.GlyphIndex(&dc.buf, r)
		if err == nil && curr != 0 {
			return ft, curr
		}
	}
	return nil, 0
}
