// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

// MeasureText measures the given text and returns a TextMetrics object with detailed metrics.
func (dc *Context) MeasureText(s string) TextMetrics {
	ft, err := fonts.match(dc.font.Family, dc.font.Weight, dc.font.Style)
	if err != nil {
		return TextMetrics{}
	}

	ppem := fixed.Int26_6(dc.font.Size * 64)

	var advance fixed.Int26_6
	var prev sfnt.GlyphIndex

	// Track actual bounding box across all glyphs
	minX, maxX := fixed.Int26_6(1<<30), fixed.Int26_6(-(1 << 30))
	minY, maxY := fixed.Int26_6(1<<30), fixed.Int26_6(-(1 << 30))

	for _, r := range s {
		curr, err := ft.GlyphIndex(&dc.buf, r)
		if err != nil || curr == 0 {
			prev = 0
			continue
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
			// No visible glyphs
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

	width := float64(advance>>6) + float64(advance&63)/64.0

	metrics, err := ft.Metrics(&dc.buf, ppem, font.HintingFull)
	var fontAscent, fontDescent float64
	if err == nil {
		fontAscent = float64(metrics.Ascent) / 64.0
		fontDescent = float64(metrics.Descent) / 64.0
	} else {
		emSquare := float64(ft.UnitsPerEm())
		scaleFactor := float64(dc.font.Size) / emSquare
		fontAscent = emSquare * scaleFactor * 0.8  // Typical ascent ratio
		fontDescent = emSquare * scaleFactor * 0.2 // Typical descent ratio
	}

	var actualBBoxLeft, actualBBoxRight, actualBBoxAscent, actualBBoxDescent float64
	if minX <= maxX {
		actualBBoxLeft = -float64(minX) / 64.0
		actualBBoxRight = float64(maxX) / 64.0
		actualBBoxAscent = -float64(minY) / 64.0
		actualBBoxDescent = float64(maxY) / 64.0
	} else {
		// No visible glyphs, use font metrics
		actualBBoxLeft = 0
		actualBBoxRight = width
		actualBBoxAscent = fontAscent
		actualBBoxDescent = fontDescent
	}

	return TextMetrics{
		Width:                    width,
		ActualBoundingBoxLeft:    actualBBoxLeft,
		ActualBoundingBoxRight:   actualBBoxRight,
		FontBoundingBoxAscent:    fontAscent,
		FontBoundingBoxDescent:   fontDescent,
		ActualBoundingBoxAscent:  actualBBoxAscent,
		ActualBoundingBoxDescent: actualBBoxDescent,
	}
}

func (dc *Context) FillText(s string, x, y float64) {
	dc.drawText(s, x, y, false)
}

func (dc *Context) StrokeText(s string, x, y float64) {
	dc.drawText(s, x, y, true)
}

func (dc *Context) drawText(s string, x, y float64, stroke bool) {
	ft, err := fonts.match(dc.font.Family, dc.font.Weight, dc.font.Style)
	if err != nil {
		return
	}

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
		curr, err := ft.GlyphIndex(&dc.buf, r)
		if err != nil {
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

		for _, seg := range segments {
			switch seg.Op {
			case sfnt.SegmentOpMoveTo:
				path2d.builder.MoveTo(
					float32(seg.Args[0].X+dot.X)/64.0,
					float32(seg.Args[0].Y+dot.Y)/64.0,
				)
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

		if stroke {
			dc.StrokePath(path2d)
		} else {
			dc.FillPath(path2d)
		}
		dot.X += advance
	}
}
