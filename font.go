// Copyright 2016 Michael Fogleman
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"io/ioutil"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"

	color2 "github.com/lumifloat/tinyskia/internal/core/color"
	"github.com/lumifloat/tinyskia/internal/core/scan"
	"github.com/lumifloat/tinyskia/internal/core/shader"
	"github.com/lumifloat/tinyskia/internal/path"
)

// Font holds font-related information
type Font struct {
	face    font.Face
	ttf     *sfnt.Font  // Keep original SFNT font for glyph outlines
	buf     sfnt.Buffer // Buffer for glyph operations
	fsize   float64
	dpi     float64 // dots per inch
	hinting font.Hinting
}

// NewFont creates a new Font with the specified SFNT font, size in points, and DPI.
func NewFont(ttf *sfnt.Font, size, dpi float64) *Font {
	face, err := opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    size,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		panic(err)
	}
	return &Font{
		face:    face,
		ttf:     ttf,
		fsize:   size,
		dpi:     dpi,
		hinting: font.HintingFull,
	}
}

// Face returns the underlying font.Face
func (f *Font) Face() font.Face {
	return f.face
}

// Size returns the font size in points
func (f *Font) Size() float64 {
	return f.fsize
}

// DPI returns the dots per inch
func (f *Font) DPI() float64 {
	return f.dpi
}

// MeasureString calculates the width and height of a string when rendered with the current font.
func (dc *Context) MeasureString(s string) (w, h float64) {
	if dc.font == nil {
		return 0, 0
	}

	d := &font.Drawer{
		Face: dc.font.Face(),
	}

	// Calculate width
	adv := d.MeasureString(s)
	w = float64(adv>>6) + float64(adv&63)/64.0

	// Calculate height from font metrics
	metrics := dc.font.Face().Metrics()
	h = float64(metrics.Ascent+metrics.Descent) / 64.0

	return w, h
}

// MeasureMultilineString calculates the dimensions of a multi-line string.
func (dc *Context) MeasureMultilineString(s string, lineSpacing float64) (width, height float64) {
	lines := splitLines(s)
	maxWidth := 0.0

	for _, line := range lines {
		lineWidth, _ := dc.MeasureString(line)
		if lineWidth > maxWidth {
			maxWidth = lineWidth
		}
	}

	width = maxWidth
	height = float64(len(lines)) * dc.FontHeight() * lineSpacing
	return width, height
}

// WordWrap wraps text to fit within the specified width.
func (dc *Context) WordWrap(text string, wrapWidth float64) []string {
	words := splitWords(text)
	if len(words) == 0 {
		return []string{}
	}

	var lines []string
	currentLine := words[0]

	for i := 1; i < len(words); i++ {
		testLine := currentLine + " " + words[i]
		testWidth, _ := dc.MeasureString(testLine)

		if testWidth <= wrapWidth {
			currentLine = testLine
		} else {
			lines = append(lines, currentLine)
			currentLine = words[i]
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
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

	// Pre-calculate scale once for all glyphs (ppem = pixels per em)
	ppem := fixed.Int26_6(dc.font.fsize * 64)

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

// calculateActualGlyphMetrics calculates the actual ascent and descent of a string
// by examining the glyph outlines. This is more accurate than using font metrics.
func (dc *Context) calculateActualGlyphMetrics(s string) (ascent, descent float64) {
	if dc.font == nil || dc.font.ttf == nil {
		return 0, 0
	}

	ppem := fixed.Int26_6(dc.font.fsize * 64)
	minY, maxY := fixed.Int26_6(1<<30), fixed.Int26_6(-(1 << 30))
	foundGlyph := false

	for _, r := range s {
		glyphIndex, err := dc.font.ttf.GlyphIndex(&dc.font.buf, r)
		if err != nil || glyphIndex == 0 {
			continue // Skip undefined glyphs
		}

		segments, err := dc.font.ttf.LoadGlyph(&dc.font.buf, glyphIndex, ppem, nil)
		if err != nil {
			continue
		}

		if len(segments) == 0 {
			continue
		}

		// Find min and max Y in this glyph's segments
		for _, seg := range segments {
			switch seg.Op {
			case sfnt.SegmentOpMoveTo, sfnt.SegmentOpLineTo:
				if seg.Args[0].Y < minY {
					minY = seg.Args[0].Y
				}
				if seg.Args[0].Y > maxY {
					maxY = seg.Args[0].Y
				}
			case sfnt.SegmentOpQuadTo:
				for i := 0; i < 2; i++ {
					if seg.Args[i].Y < minY {
						minY = seg.Args[i].Y
					}
					if seg.Args[i].Y > maxY {
						maxY = seg.Args[i].Y
					}
				}
			case sfnt.SegmentOpCubeTo:
				for i := 0; i < 3; i++ {
					if seg.Args[i].Y < minY {
						minY = seg.Args[i].Y
					}
					if seg.Args[i].Y > maxY {
						maxY = seg.Args[i].Y
					}
				}
			}
		}
		foundGlyph = true
	}

	if !foundGlyph {
		return 0, 0
	}

	// Convert from fixed point to float
	// In sfnt: Y increases downward (same as screen coordinates)
	// For anchor calculation, we need:
	//   ascent = distance from baseline to the HIGHEST point = -minY (since Y-down)
	//   descent = distance from baseline to the LOWEST point = maxY

	if minY < 0 {
		ascent = float64(-minY) / 64.0
	}
	if maxY > 0 {
		descent = float64(maxY) / 64.0
	}

	return ascent, descent
}

// DrawStringAnchored draws text anchored at a specific point.
// ax and ay should be in the range [0,1].
// (0,0) anchors at top-left, (0.5,0.5) centers, (1,1) anchors at bottom-right.
func (dc *Context) DrawStringAnchored(s string, x, y, ax, ay float64) {
	if s == "" || dc.font == nil {
		return
	}

	w, _ := dc.MeasureString(s)

	// Calculate horizontal anchor offset
	dx := w * ax

	// For vertical positioning, we need to account for how glyphs actually render.
	// After manual Y-flip (screenY = baselineY - glyphY):
	//   - Glyph points with positive Y (above baseline) render HIGHER on screen
	//   - Glyph points with negative Y (below baseline) render LOWER on screen
	//   - The topmost rendered pixel is at: baselineY - max(glyphY)
	//   - The bottommost rendered pixel is at: baselineY - min(glyphY)
	//
	// To achieve proper anchoring:
	//   - When ay=0: top of glyphs should be at y
	//   - When ay=0.5: middle of glyphs should be at y
	//   - When ay=1: bottom of glyphs should be at y
	//
	// We calculate actual glyph extents and use that for positioning.

	actualAscent, actualDescent := dc.calculateActualGlyphMetrics(s)
	visualHeight := actualAscent + actualDescent

	// Calculate dy based on visual height and anchor
	dy := visualHeight * ay

	// The key insight: after Y-flip, the top of the glyph is at (baseline - actualAscent)
	// We want the top to be at (y - dy)
	// So: baseline - actualAscent = y - dy
	// Therefore: baseline = y - dy + actualAscent
	adjustedY := y - dy + actualAscent

	dc.DrawString(s, x-dx, adjustedY)
}

// DrawStringWrapped draws wrapped text with alignment.
func (dc *Context) DrawStringWrapped(s string, x, y, ax, ay, width, lineSpacing float64, align Align) {
	lines := dc.WordWrap(s, width)
	_, lineHeight := dc.MeasureString("X")
	totalHeight := float64(len(lines)) * lineHeight * lineSpacing

	// Calculate starting position based on anchor
	startY := y - totalHeight*ay

	for i, line := range lines {
		lineY := startY + float64(i)*lineHeight*lineSpacing
		lineW, _ := dc.MeasureString(line)

		var lineX float64
		switch align {
		case AlignLeft:
			lineX = x - width*ax
		case AlignCenter:
			lineX = x - lineW/2
		case AlignRight:
			lineX = x - width*ax + (width - lineW)
		default:
			lineX = x - width*ax
		}

		dc.DrawString(line, lineX, lineY)
	}
}

// splitLines splits a string by newline characters
func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	lines = append(lines, s[start:])
	return lines
}

// splitWords splits a string into words
func splitWords(s string) []string {
	var words []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' || s[i] == '\t' || s[i] == '\n' {
			if i > start {
				words = append(words, s[start:i])
			}
			start = i + 1
		}
	}
	if start < len(s) {
		words = append(words, s[start:])
	}
	return words
}

// LoadFontFace loads a TrueType font from a file.
func (dc *Context) LoadFontFace(path string, points float64) error {
	return dc.LoadFontFaceWithDPI(path, points, 72.0)
}

// LoadFontFaceWithDPI loads a TrueType font from a file with custom DPI.
func (dc *Context) LoadFontFaceWithDPI(path string, points, dpi float64) error {
	// Load the font file
	fontBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// Parse the font using sfnt
	ttf, err := sfnt.Parse(fontBytes)
	if err != nil {
		return err
	}

	// Create font object
	dc.font = NewFont(ttf, points, dpi)
	return nil
}

// LoadFontFaceFromData loads a TrueType font from byte data.
func (dc *Context) LoadFontFaceFromData(data []byte, points float64) error {
	return dc.LoadFontFaceFromDataWithDPI(data, points, 72.0)
}

// LoadFontFaceFromDataWithDPI loads a TrueType font from byte data with custom DPI.
func (dc *Context) LoadFontFaceFromDataWithDPI(data []byte, points, dpi float64) error {
	// Parse the font using sfnt
	ttf, err := sfnt.Parse(data)
	if err != nil {
		return err
	}

	// Create font object
	dc.font = NewFont(ttf, points, dpi)
	return nil
}

// SetFontFace sets the current font face.
// Accepts either *Font (for vector rendering with outlines) or font.Face (for bitmap rendering).
func (dc *Context) SetFontFace(f interface{}) {
	switch v := f.(type) {
	case *Font:
		dc.font = v
		dc.fontFace = v.face
		dc.fontHeight = float64(v.face.Metrics().Height) / 64
	case font.Face:
		dc.fontFace = v
		dc.font = nil // Clear font when using plain font.Face
		dc.fontHeight = float64(v.Metrics().Height) / 64
	}
}

// FontHeight returns the approximate height of the current font in pixels.
func (dc *Context) FontHeight() float64 {
	if dc.font == nil {
		return 0
	}

	metrics := dc.font.Face().Metrics()
	return float64(metrics.Ascent+metrics.Descent) / 64.0
}

// GetFont returns the current font
func (dc *Context) GetFont() *Font {
	return dc.font
}
