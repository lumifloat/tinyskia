// Copyright 2016 Michael Fogleman
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"fmt"
	"io/ioutil"

	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"

	color2 "github.com/lumifloat/tinyskia/internal/core/color"
	"github.com/lumifloat/tinyskia/internal/core/scan"
	"github.com/lumifloat/tinyskia/internal/core/shader"
	"github.com/lumifloat/tinyskia/internal/path"
)

// Font holds font-related information
type Font struct {
	ttf     *sfnt.Font  // Keep original SFNT font for glyph outlines
	buf     sfnt.Buffer // Buffer for glyph operations
	size    float64
	dpi     float64 // dots per inch
	hinting font.Hinting
}

// NewFont creates a new Font with the specified SFNT font, size in points, and DPI.
func NewFont(ttf *sfnt.Font, size, dpi float64) *Font {
	return &Font{
		ttf:     ttf,
		size:    size,
		dpi:     dpi,
		hinting: font.HintingFull,
	}
}

// Size returns the font size in points
func (f *Font) Size() float64 {
	return f.size
}

// DPI returns the dots per inch
func (f *Font) DPI() float64 {
	return f.dpi
}

// TextMetrics holds text measurement information, similar to HTML5 Canvas TextMetrics.
type TextMetrics struct {
	// Width is the advance width of the text (inline box width).
	Width float64

	// ActualBoundingBoxLeft is the distance from the alignment point to the left side
	// of the bounding rectangle. Positive values indicate left direction.
	ActualBoundingBoxLeft float64

	// ActualBoundingBoxRight is the distance from the alignment point to the right side
	// of the bounding rectangle. Positive values indicate right direction.
	ActualBoundingBoxRight float64

	// FontBoundingBoxAscent is the distance from the baseline to the font's ascent metric,
	// in CSS pixels. Positive values indicate upward direction.
	FontBoundingBoxAscent float64

	// FontBoundingBoxDescent is the distance from the baseline to the font's descent metric,
	// in CSS pixels. Positive values indicate downward direction.
	FontBoundingBoxDescent float64

	// ActualBoundingBoxAscent is the distance from the baseline to the top of the
	// bounding rectangle of the given text. Positive values indicate upward.
	ActualBoundingBoxAscent float64

	// ActualBoundingBoxDescent is the distance from the baseline to the bottom of the
	// bounding rectangle of the given text. Positive values indicate downward.
	ActualBoundingBoxDescent float64
}

// MeasureText measures the given text and returns a TextMetrics object with detailed metrics.
func (dc *Context) MeasureText(s string) TextMetrics {
	if dc.font == nil || dc.font.ttf == nil || s == "" {
		return TextMetrics{}
	}

	ppem := fixed.Int26_6(dc.font.size * 64)
	var advance fixed.Int26_6
	var prev sfnt.GlyphIndex

	// Track actual bounding box across all glyphs
	minX, maxX := fixed.Int26_6(1<<30), fixed.Int26_6(-(1 << 30))
	minY, maxY := fixed.Int26_6(1<<30), fixed.Int26_6(-(1 << 30))

	for _, r := range s {
		curr, err := dc.font.ttf.GlyphIndex(&dc.font.buf, r)
		if err != nil || curr == 0 {
			prev = 0
			continue
		}

		if prev != 0 {
			kern, err := dc.font.ttf.Kern(&dc.font.buf, prev, curr, ppem, dc.font.hinting)
			if err == nil {
				advance += kern
			}
		}
		prev = curr

		b, a, err := dc.font.ttf.GlyphBounds(&dc.font.buf, curr, ppem, dc.font.hinting)
		if err != nil {
			continue
		}
		advance += a
		if b.Empty() {
			// No visible glyphs
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
	}

	width := float64(advance>>6) + float64(advance&63)/64.0

	metrics, err := dc.font.ttf.Metrics(&dc.font.buf, ppem, dc.font.hinting)
	var fontAscent, fontDescent float64
	if err == nil {
		fontAscent = float64(metrics.Ascent) / 64.0
		fontDescent = float64(metrics.Descent) / 64.0
	} else {
		emSquare := float64(dc.font.ttf.UnitsPerEm())
		scaleFactor := dc.font.size / emSquare
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

// DrawString draws text at the specified position using TinySkia's rasterization.
// It converts each glyph outline to a path and fills it using FillPath.
func (dc *Context) DrawString(s string, x, y float64) {
	if dc.font == nil || s == "" || dc.font.ttf == nil {
		return
	}

	// Get shader from fillStyle - default to black if not set
	fillShader := toShader(dc.fillStyle, dc.transform)
	if fillShader == nil {
		fillShader = shader.NewSolidColor(color2.ColorBlack)
	}

	// Convert coordinates to fixed point (26.6 format)
	fx := fixed.I(int(x))
	fy := fixed.I(int(y))
	dot := fixed.Point26_6{X: fx, Y: fy}

	ppem := fixed.Int26_6(dc.font.size * 64)

	// Cache space glyph advance width
	var spaceAdvance fixed.Int26_6
	spaceGlyph, _ := dc.font.ttf.GlyphIndex(&dc.font.buf, ' ')
	if spaceGlyph != 0 {
		advance, err := dc.font.ttf.GlyphAdvance(&dc.font.buf, spaceGlyph, ppem, font.HintingNone)
		if err == nil {
			spaceAdvance = advance
		}
	}

	for _, r := range s {
		// Get glyph index
		glyphIndex, err := dc.font.ttf.GlyphIndex(&dc.font.buf, r)
		if err != nil || glyphIndex == 0 {
			// Glyph not found, advance by cached space width
			dot.X += spaceAdvance
			continue
		}

		// Load glyph segments
		segments, err := dc.font.ttf.LoadGlyph(&dc.font.buf, glyphIndex, ppem, nil)
		if err != nil {
			dot.X += spaceAdvance
			continue
		}

		// Skip empty glyphs but still advance
		if len(segments) == 0 {
			advance, _ := dc.font.ttf.GlyphAdvance(&dc.font.buf, glyphIndex, ppem, font.HintingNone)
			dot.X += advance
			continue
		}

		// Convert glyph outline to path using a fresh path builder
		glyphPathBuilder := path.NewPathBuilder()
		// Temporarily replace the context's path builder
		oldPathBuilder := dc.pathBuilder
		dc.pathBuilder = glyphPathBuilder

		// Convert segments to path
		dc.segmentsToPath(segments, dot)

		// Render this glyph
		pathData := dc.pathBuilder.Finish()
		if pathData != nil {
			// Transform the path
			var transformedPath *path.Path
			if !dc.transform.IsIdentity() {
				transformedPath = pathData.Transform(dc.transform)
			} else {
				transformedPath = pathData
			}

			// Create paint - disable AA for text rendering to avoid performance issues
			paint := &Paint{
				Shader:          fillShader,
				AntiAlias:       false, // Disable AA for text to avoid O(n²) complexity
				BlendMode:       dc.blendMode,
				Colorspace:      dc.colorspace,
				ForceHQPipeline: false,
			}

			var maskData []uint8
			if dc.mask != nil {
				maskData = dc.mask.Pix
			}
			blitter := paint.blitter(dc.im.Pix, maskData, dc.Width(), dc.Height())
			screen, _ := path.NewScreenIntRectFromXYWH(0, 0, uint32(dc.Width()), uint32(dc.Height()))

			// Use even-odd fill rule for text to correctly handle holes
			scan.FillPathAA(transformedPath, int(FillRuleEvenOdd), screen, blitter)
		}

		// Restore the original path builder
		dc.pathBuilder = oldPathBuilder

		// Advance to next glyph position
		advance, _ := dc.font.ttf.GlyphAdvance(&dc.font.buf, glyphIndex, ppem, font.HintingNone)
		dot.X += advance
	}
}

// segmentsToPath converts sfnt Segments to a path
// sfnt segments use Y-down coordinate system (same as screen)
func (dc *Context) segmentsToPath(segments sfnt.Segments, dot fixed.Point26_6) {
	if len(segments) == 0 {
		return
	}

	for _, seg := range segments {
		switch seg.Op {
		case sfnt.SegmentOpMoveTo:
			// Convert from 26.6 fixed point to float32 and add offset
			x := float32(seg.Args[0].X)/64.0 + float32(dot.X)/64.0
			y := float32(seg.Args[0].Y)/64.0 + float32(dot.Y)/64.0
			dc.pathBuilder.MoveTo(x, y)
		case sfnt.SegmentOpLineTo:
			x := float32(seg.Args[0].X)/64.0 + float32(dot.X)/64.0
			y := float32(seg.Args[0].Y)/64.0 + float32(dot.Y)/64.0
			dc.pathBuilder.LineTo(x, y)
		case sfnt.SegmentOpQuadTo:
			// Quadratic Bezier: control point, end point
			ctrlX := float32(seg.Args[0].X)/64.0 + float32(dot.X)/64.0
			ctrlY := float32(seg.Args[0].Y)/64.0 + float32(dot.Y)/64.0
			endX := float32(seg.Args[1].X)/64.0 + float32(dot.X)/64.0
			endY := float32(seg.Args[1].Y)/64.0 + float32(dot.Y)/64.0
			dc.pathBuilder.QuadTo(ctrlX, ctrlY, endX, endY)
		case sfnt.SegmentOpCubeTo:
			// Cubic Bezier: two control points, end point
			ctrl1X := float32(seg.Args[0].X)/64.0 + float32(dot.X)/64.0
			ctrl1Y := float32(seg.Args[0].Y)/64.0 + float32(dot.Y)/64.0
			ctrl2X := float32(seg.Args[1].X)/64.0 + float32(dot.X)/64.0
			ctrl2Y := float32(seg.Args[1].Y)/64.0 + float32(dot.Y)/64.0
			endX := float32(seg.Args[2].X)/64.0 + float32(dot.X)/64.0
			endY := float32(seg.Args[2].Y)/64.0 + float32(dot.Y)/64.0
			dc.pathBuilder.CubicTo(ctrl1X, ctrl1Y, ctrl2X, ctrl2Y, endX, endY)
		}
	}
}

// LoadFontFace loads a TrueType font from a file.
func (dc *Context) LoadFontFace(path string, points float64) error {
	return dc.LoadFontFaceWithDPI(path, points, 72.0)
}

// LoadFontFaceWithDPI loads a TrueType font from a file with custom DPI.
// Supports both single fonts (TTF/OTF) and font collections (TTC/OTC).
// For collections, it loads the first font (index 0).
func (dc *Context) LoadFontFaceWithDPI(path string, points, dpi float64) error {
	// Load the font file
	fontBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// Try to parse as a collection first (TTC/OTC)
	collection, err := sfnt.ParseCollection(fontBytes)
	if err == nil && collection.NumFonts() > 0 {
		// Successfully parsed as collection, use the first font
		ttf, err := collection.Font(0)
		if err != nil {
			return fmt.Errorf("failed to get first font from collection: %v", err)
		}
		// Create font object with raw data
		dc.font = NewFont(ttf, points, dpi)
		return nil
	}

	// Not a collection or failed to parse as collection, try as single font
	ttf, err := sfnt.Parse(fontBytes)
	if err != nil {
		return err
	}

	// Create font object with raw data
	dc.font = NewFont(ttf, points, dpi)
	return nil
}

// LoadFontFaceFromData loads a TrueType font from byte data.
func (dc *Context) LoadFontFaceFromData(data []byte, points float64) error {
	return dc.LoadFontFaceFromDataWithDPI(data, points, 72.0)
}

// LoadFontFaceFromDataWithDPI loads a TrueType font from byte data with custom DPI.
// Supports both single fonts (TTF/OTF) and font collections (TTC/OTC).
// For collections, it loads the first font (index 0).
func (dc *Context) LoadFontFaceFromDataWithDPI(data []byte, points, dpi float64) error {
	// Try to parse as a collection first (TTC/OTC)
	collection, err := sfnt.ParseCollection(data)
	if err == nil && collection.NumFonts() > 0 {
		// Successfully parsed as collection, use the first font
		ttf, err := collection.Font(0)
		if err != nil {
			return fmt.Errorf("failed to get first font from collection: %v", err)
		}
		dc.font = NewFont(ttf, points, dpi)
		return nil
	}

	// Not a collection or failed to parse as collection, try as single font
	ttf, err := sfnt.Parse(data)
	if err != nil {
		return err
	}

	dc.font = NewFont(ttf, points, dpi)
	return nil
}

// SetFontFace sets the current font face.
// Accepts either *Font (for vector rendering with outlines) or font.Face (for bitmap rendering).
func (dc *Context) SetFontFace(f *Font) {
	dc.font = f
}

// GetFont returns the current font
func (dc *Context) GetFont() *Font {
	return dc.font
}
