// Copyright 2016 Michael Fogleman
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"

	color2 "github.com/lumifloat/tinyskia/internal/core/color"
	"github.com/lumifloat/tinyskia/internal/core/scan"
	"github.com/lumifloat/tinyskia/internal/core/shader"
	"github.com/lumifloat/tinyskia/internal/path"
	"golang.org/x/image/font"
)

type LineCap int

const (
	LineCapRound LineCap = iota
	LineCapButt
	LineCapSquare
)

type LineJoin int

const (
	LineJoinRound LineJoin = iota
	LineJoinBevel
)

type FillRule int

const (
	FillRuleWinding FillRule = iota
	FillRuleEvenOdd
)

type Align int

const (
	AlignLeft Align = iota
	AlignCenter
	AlignRight
)

// Context is the main drawing context, similar to gg.Context.
// It maintains drawing state and provides a canvas-like API.
type Context struct {
	width           int
	height          int
	im              *image.RGBA
	mask            *image.Alpha
	color           color.Color
	fillStyle       Style
	strokeStyle     Style
	pathBuilder     *path.PathBuilder
	start           path.Point
	current         path.Point
	hasCurrent      bool
	dashes          []float64
	dashOffset      float64
	lineWidth       float64
	lineCap         LineCap
	lineJoin        LineJoin
	fillRule        FillRule
	fontFace        font.Face
	fontHeight      float64
	transform       path.Transform
	blendMode       BlendMode
	antiAlias       bool
	colorspace      color2.ColorSpace
	forceHQPipeline bool
	stack           []*Context
}

func NewContext(width, height int) *Context {
	return NewContextForRGBA(image.NewRGBA(image.Rect(0, 0, width, height)))
}

// NewContextForImage creates a context from an existing image.Image.
// No copy is made.
func NewContextForImage(im image.Image) *Context {
	return NewContextForRGBA(imageToRGBA(im))
}

func imageToRGBA(im image.Image) *image.RGBA {
	bounds := im.Bounds()
	rgba := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba.Set(x, y, im.At(x, y))
		}
	}
	return rgba
}

func NewContextForRGBA(im *image.RGBA) *Context {
	bounds := im.Bounds()
	return &Context{
		width:           bounds.Dx(),
		height:          bounds.Dy(),
		im:              im,
		pathBuilder:     path.NewPathBuilder(),
		lineWidth:       1,
		lineCap:         LineCapRound,
		lineJoin:        LineJoinRound,
		fillRule:        FillRuleWinding,
		fontFace:        nil, // Will be set when font support is added
		fontHeight:      0,
		transform:       path.NewTransformDefault(),
		blendMode:       BlendModeSourceOver,
		antiAlias:       true,
		colorspace:      color2.ColorSpaceLinear,
		forceHQPipeline: false,
	}
}

// GetCurrentPoint will return the current point and if there is a current point.
// The point will have been transformed by the context's transformation matrix.
func (dc *Context) GetCurrentPoint() (path.Point, bool) {
	if dc.hasCurrent {
		return dc.current, true
	}
	return path.Point{}, false
}

// Image returns the image that has been drawn by this context.
func (dc *Context) Image() image.Image {
	return dc.im
}

// Width returns the width of the image in pixels.
func (dc *Context) Width() int {
	return dc.width
}

// Height returns the height of the image in pixels.
func (dc *Context) Height() int {
	return dc.height
}

// SavePNG encodes the image as a PNG and writes it to disk.
func (dc *Context) SavePNG(path string) error {
	return SavePNG(path, dc.Image())
}

// EncodePNG encodes the image as a PNG and writes it to the provided io.Writer.
func (dc *Context) EncodePNG(w io.Writer) error {
	return png.Encode(w, dc.Image())
}

// State Management (Push/Pop)

// Push saves the current state of the context for later retrieval. These
// can be nested.
func (dc *Context) Push() {
	x := *dc
	dc.stack = append(dc.stack, &x)
}

// Pop restores the last saved context state from the stack.
func (dc *Context) Pop() {
	before := *dc
	s := dc.stack
	x, s := s[len(s)-1], s[:len(s)-1]
	*dc = *x
	dc.mask = before.mask
	dc.im = before.im
	pathData := before.pathBuilder.Finish()
	dc.pathBuilder = path.NewPathBuilder()
	if pathData != nil {
		// Restore path data if needed
	}
	dc.start = before.start
	dc.current = before.current
	dc.hasCurrent = before.hasCurrent
	dc.stack = s
}

// SetDash sets the current dash pattern to use. Call with zero arguments to
// disable dashes. The values specify the lengths of each dash, with
// alternating on and off lengths.
func (dc *Context) SetDash(dashes ...float64) {
	dc.dashes = dashes
}

// SetDashOffset sets the initial offset into the dash pattern to use when
// stroking dashed paths.
func (dc *Context) SetDashOffset(offset float64) {
	dc.dashOffset = offset
}

// SetLineWidth sets the line width for stroking paths.
func (dc *Context) SetLineWidth(lineWidth float64) {
	dc.lineWidth = lineWidth
}

// SetLineCap sets the line cap style (Butt, Round, Square).
func (dc *Context) SetLineCap(lineCap LineCap) {
	dc.lineCap = lineCap
}

// SetLineCapRound sets the line cap to round.
func (dc *Context) SetLineCapRound() {
	dc.lineCap = LineCapRound
}

// SetLineCapButt sets the line cap to butt.
func (dc *Context) SetLineCapButt() {
	dc.lineCap = LineCapButt
}

// SetLineCapSquare sets the line cap to square.
func (dc *Context) SetLineCapSquare() {
	dc.lineCap = LineCapSquare
}

// SetLineJoin sets the line join style (Bevel, Round, Miter).
func (dc *Context) SetLineJoin(lineJoin LineJoin) {
	dc.lineJoin = lineJoin
}

// SetLineJoinRound sets the line join to round.
func (dc *Context) SetLineJoinRound() {
	dc.lineJoin = LineJoinRound
}

// SetLineJoinBevel sets the line join to bevel.
func (dc *Context) SetLineJoinBevel() {
	dc.lineJoin = LineJoinBevel
}

func (dc *Context) SetFillRule(fillRule FillRule) {
	dc.fillRule = fillRule
}

func (dc *Context) SetFillRuleWinding() {
	dc.fillRule = FillRuleWinding
}

func (dc *Context) SetFillRuleEvenOdd() {
	dc.fillRule = FillRuleEvenOdd
}

// SetAntiAlias enables or disables anti-aliasing.
func (dc *Context) SetAntiAlias(aa bool) {
	dc.antiAlias = aa
}

// SetFillStyle sets current fill style.
// Accepts Gradient, SolidPattern, SurfacePattern or other Style implementations.
func (dc *Context) SetFillStyle(style Style) {
	dc.fillStyle = style
	// If it's a solid pattern, also update color
	if solid, ok := style.(*solidPattern); ok {
		dc.color = solid.color
	}
}

// SetStrokeStyle sets current stroke style.
// Accepts Gradient, SolidPattern, SurfacePattern or other Style implementations.
func (dc *Context) SetStrokeStyle(style Style) {
	dc.strokeStyle = style
}

// SetColor sets the current color(for both fill and stroke).
func (dc *Context) SetColor(c color.Color) {
	dc.color = c
	dc.fillStyle = NewSolidPattern(c)
	dc.strokeStyle = NewSolidPattern(c)
}

// SetHexColor sets the current color using a hex string. The leading pound
// sign (#) is optional. Both 3- and 6-digit variations are supported. 8 digits
// may be provided to set the alpha value as well.
func (dc *Context) SetHexColor(hexStr string) {
	// Remove leading # if present
	if len(hexStr) > 0 && hexStr[0] == '#' {
		hexStr = hexStr[1:]
	}

	var r, g, b, a uint8
	switch len(hexStr) {
	case 3:
		// 3-digit hex (RGB)
		r = parseHexChar(hexStr[0]) * 17
		g = parseHexChar(hexStr[1]) * 17
		b = parseHexChar(hexStr[2]) * 17
		a = 255
	case 6:
		// 6-digit hex (RGB)
		r = parseHexByte(hexStr[0:2])
		g = parseHexByte(hexStr[2:4])
		b = parseHexByte(hexStr[4:6])
		a = 255
	case 8:
		// 8-digit hex (RGBA)
		r = parseHexByte(hexStr[0:2])
		g = parseHexByte(hexStr[2:4])
		b = parseHexByte(hexStr[4:6])
		a = parseHexByte(hexStr[6:8])
	default:
		// Invalid format, use black
		r, g, b, a = 0, 0, 0, 255
	}

	dc.SetRGBA255(int(r), int(g), int(b), int(a))
}

// parseHexChar converts a single hex character to its numeric value
func parseHexChar(c byte) uint8 {
	switch {
	case c >= '0' && c <= '9':
		return c - '0'
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10
	default:
		return 0
	}
}

// parseHexByte converts two hex characters to a byte value
func parseHexByte(s string) uint8 {
	if len(s) != 2 {
		return 0
	}
	return parseHexChar(s[0])*16 + parseHexChar(s[1])
}

// SetRGBA255 sets the current color. r, g, b, a values should be between 0 and
// 255, inclusive.
func (dc *Context) SetRGBA255(r, g, b, a int) {
	dc.SetColor(color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
}

// SetRGB255 sets the current color. r, g, b values should be between 0 and 255,
// inclusive. Alpha will be set to 255 (fully opaque).
func (dc *Context) SetRGB255(r, g, b int) {
	dc.SetRGBA255(r, g, b, 255)
}

// SetRGBA sets the current color. r, g, b, a values should be between 0 and 1,
// inclusive.
func (dc *Context) SetRGBA(r, g, b, a float64) {
	dc.SetColor(color.NRGBA{
		uint8(r * 255),
		uint8(g * 255),
		uint8(b * 255),
		uint8(a * 255),
	})
}

// SetRGB sets the current color. r, g, b values should be between 0 and 1,
// inclusive. Alpha will be set to 1 (fully opaque).
func (dc *Context) SetRGB(r, g, b float64) {
	dc.SetRGBA(r, g, b, 1)
}

// MoveTo starts a new subpath within the current path starting at the
// specified point.
func (dc *Context) MoveTo(x, y float64) {
	p := path.Point{X: float32(x), Y: float32(y)}
	dc.pathBuilder.MoveTo(p.X, p.Y)
	dc.start = p
	dc.current = p
	dc.hasCurrent = true
}

// LineTo adds a line segment to the current path starting at the current
// point. If there is no current point, it is equivalent to MoveTo(x, y)
func (dc *Context) LineTo(x, y float64) {
	if !dc.hasCurrent {
		dc.MoveTo(x, y)
	} else {
		p := path.Point{X: float32(x), Y: float32(y)}
		dc.pathBuilder.LineTo(p.X, p.Y)
		dc.current = p
	}
}

// QuadraticTo adds a quadratic bezier curve to the current path starting at
// the current point. If there is no current point, it first performs
// MoveTo(x1, y1)
func (dc *Context) QuadraticTo(x1, y1, x2, y2 float64) {
	if !dc.hasCurrent {
		dc.MoveTo(x1, y1)
	}
	p1 := path.Point{X: float32(x1), Y: float32(y1)}
	p2 := path.Point{X: float32(x2), Y: float32(y2)}
	dc.pathBuilder.QuadTo(p1.X, p1.Y, p2.X, p2.Y)
	dc.current = p2
}

// CubicTo adds a cubic bezier curve to the current path starting at the
// current point. If there is no current point, it first performs
// MoveTo(x1, y1).
func (dc *Context) CubicTo(x1, y1, x2, y2, x3, y3 float64) {
	if !dc.hasCurrent {
		dc.MoveTo(x1, y1)
	}
	p1 := path.Point{X: float32(x1), Y: float32(y1)}
	p2 := path.Point{X: float32(x2), Y: float32(y2)}
	p3 := path.Point{X: float32(x3), Y: float32(y3)}
	dc.pathBuilder.CubicTo(p1.X, p1.Y, p2.X, p2.Y, p3.X, p3.Y)
	dc.current = p3
}

// ClosePath adds a line segment from the current point to the beginning
// of the current subpath. If there is no current point, this is a no-op.
func (dc *Context) ClosePath() {
	if dc.hasCurrent {
		dc.pathBuilder.Close()
	}
}

// ClearPath clears the current path. There is no current point after this
// operation.
func (dc *Context) ClearPath() {
	dc.pathBuilder.Clear()
	dc.hasCurrent = false
}

// NewSubPath starts a new subpath within the current path. There is no current
// point after this operation.
func (dc *Context) NewSubPath() {
	dc.hasCurrent = false
}

// Fill fills the current path with the current color. Open subpaths
// are implicity closed. The path is cleared after this operation.
func (dc *Context) Fill() {
	dc.FillPreserve()
	dc.ClearPath()
}

// FillPreserve fills the current path with the current color. Open subpaths
// are implicity closed. The path is preserved after this operation.
func (dc *Context) FillPreserve() {
	p := dc.pathBuilder.Finish()
	if p == nil {
		return
	}

	// Both path and shader need to be transformed
	// This matches tiny-skia's behavior where transform applies to both
	var transformedPath *path.Path
	if !dc.transform.IsIdentity() {
		transformedPath = p.Transform(dc.transform)
	} else {
		transformedPath = p
	}

	paint := &Paint{
		Shader:          toShader(dc.fillStyle, dc.transform),
		AntiAlias:       dc.antiAlias,
		BlendMode:       dc.blendMode,
		Colorspace:      dc.colorspace,
		ForceHQPipeline: dc.forceHQPipeline,
	}
	var maskData []uint8
	if dc.mask != nil {
		maskData = dc.mask.Pix
	}
	blitter := paint.blitter(dc.im.Pix, maskData, dc.Width(), dc.Height())
	screen, _ := path.NewScreenIntRectFromXYWH(0, 0, uint32(dc.Width()), uint32(dc.Height()))
	scan.FillPathAA(transformedPath, int(dc.fillRule), screen, blitter)
}

// Stroke strokes the current path with the current color, line width,
// line cap, line join and dash settings. The path is cleared after this
// operation.
func (dc *Context) Stroke() {
	dc.StrokePreserve()
	dc.ClearPath()
}

// StrokePreserve strokes the current path with the current color, line width,
// line cap, line join and dash settings. The path is preserved after this
// operation.
func (dc *Context) StrokePreserve() {
	pathData := dc.pathBuilder.Finish()
	if pathData == nil {
		return
	}

	// Apply transform to path
	var transformedPath *path.Path
	if !dc.transform.IsIdentity() {
		transformedPath = pathData.Transform(dc.transform)
	} else {
		transformedPath = pathData
	}

	// Convert LineCap from context enum to path stroker enum
	var lineCap path.LineCap
	switch dc.lineCap {
	case LineCapRound:
		lineCap = path.LineCapRound
	case LineCapButt:
		lineCap = path.LineCapButt
	case LineCapSquare:
		lineCap = path.LineCapSquare
	}

	// Convert LineJoin from context enum to path stroker enum
	var lineJoin path.LineJoin
	switch dc.lineJoin {
	case LineJoinRound:
		lineJoin = path.LineJoinRound
	case LineJoinBevel:
		lineJoin = path.LineJoinBevel
	default:
		lineJoin = path.LineJoinMiter // Default to Miter
	}

	// Build stroke options
	stroke := &path.Stroke{
		Width:      float32(dc.lineWidth),
		LineCap:    lineCap,
		LineJoin:   lineJoin,
		MiterLimit: 4.0, // Default miter limit
	}

	// Add dashing if specified
	if len(dc.dashes) > 0 {
		// Convert dashes from float64 to float32
		dashArray := make([]float32, len(dc.dashes))
		for i, d := range dc.dashes {
			dashArray[i] = float32(d)
		}
		stroke.Dash = path.NewStrokeDash(dashArray, float32(dc.dashOffset))
	}

	// Compute resolution scale based on transform
	resScale := path.ComputeResolutionScale(dc.transform)

	// Stroke the path to get a filled outline
	stroker := path.NewPathStroker()
	strokedPath := stroker.Stroke(transformedPath, *stroke, resScale)
	if strokedPath == nil {
		return
	}

	// For lineWidth <= 1, disable AA to ensure visibility
	// AA makes very thin lines nearly invisible due to low coverage
	useAA := dc.antiAlias && dc.lineWidth > 1.0

	// Fill the stroked path
	paint := &Paint{
		Shader:          toShader(dc.strokeStyle, path.NewTransformDefault()),
		AntiAlias:       useAA,
		BlendMode:       dc.blendMode,
		Colorspace:      dc.colorspace,
		ForceHQPipeline: dc.forceHQPipeline,
	}
	var maskData []uint8
	if dc.mask != nil {
		maskData = dc.mask.Pix
	}
	blitter := paint.blitter(dc.im.Pix, maskData, dc.Width(), dc.Height())
	screen, _ := path.NewScreenIntRectFromXYWH(0, 0, uint32(dc.Width()), uint32(dc.Height()))

	// Use AA or non-AA fill based on line width
	if useAA {
		scan.FillPathAA(strokedPath, int(dc.fillRule), screen, blitter)
	} else {
		scan.FillPath(strokedPath, int(dc.fillRule), screen, blitter)
	}
}

// Clip updates the clipping region by intersecting the current
// clipping region with the current path as it would be filled by dc.Fill().
// The path is cleared after this operation.
func (dc *Context) Clip() {
	dc.ClipPreserve()
	dc.ClearPath()
}

// ClipPreserve updates the clipping region by intersecting the current
// clipping region with the current path as it would be filled by dc.Fill().
// The path is preserved after this operation.
func (dc *Context) ClipPreserve() {
	pathData := dc.pathBuilder.Finish()
	if pathData == nil {
		return
	}

	// Apply transform to path for clipping
	var transformedPath *path.Path
	if !dc.transform.IsIdentity() {
		transformedPath = pathData.Transform(dc.transform)
	} else {
		transformedPath = pathData
	}

	// Create a temporary alpha mask for the clip path
	width := dc.Width()
	height := dc.Height()
	clipMask := image.NewAlpha(image.Rect(0, 0, width, height))

	// Render path to a temporary RGBA image first (blitter expects RGBA)
	tempRGBA := image.NewRGBA(image.Rect(0, 0, width, height))
	paint := &Paint{
		Shader:          shader.NewSolidColor(color2.ColorFromRGBA8(255, 255, 255, 255)),
		AntiAlias:       dc.antiAlias,
		BlendMode:       BlendModeSource,
		Colorspace:      dc.colorspace,
		ForceHQPipeline: dc.forceHQPipeline,
	}
	blitter := paint.blitter(tempRGBA.Pix, nil, width, height)
	screen, _ := path.NewScreenIntRectFromXYWH(0, 0, uint32(width), uint32(height))
	scan.FillPathAA(transformedPath, int(dc.fillRule), screen, blitter)

	// Extract alpha channel from RGBA to Alpha mask
	for i := 0; i < width*height; i++ {
		clipMask.Pix[i] = tempRGBA.Pix[i*4+3] // Alpha channel
	}

	// Intersect with existing mask (same as gg's draw.DrawMask with draw.Over)
	if dc.mask == nil {
		dc.mask = clipMask
	} else {
		// Create new mask by combining old mask and clip mask
		mask := image.NewAlpha(image.Rect(0, 0, width, height))
		for i := range mask.Pix {
			// draw.Over compositing: result = clip + old * (1 - clip_alpha/255)
			clipAlpha := uint32(clipMask.Pix[i])
			oldAlpha := uint32(dc.mask.Pix[i])
			mask.Pix[i] = uint8(clipAlpha + oldAlpha*(255-clipAlpha)/255)
		}
		dc.mask = mask
	}
}

// ResetClip clears the clipping region.
func (dc *Context) ResetClip() {
	dc.mask = nil
}

// SetMask allows you to directly set the *image.Alpha to be used as a clipping
// mask. It must be the same size as the context, else an error is returned
// and the mask is unchanged.
func (dc *Context) SetMask(maskImg *image.Alpha) error {
	if maskImg == nil {
		return nil
	}
	// Check size
	if maskImg.Rect.Dx() != dc.Width() || maskImg.Rect.Dy() != dc.Height() {
		return fmt.Errorf("mask size %dx%d does not match context size %dx%d",
			maskImg.Rect.Dx(), maskImg.Rect.Dy(), dc.Width(), dc.Height())
	}
	dc.mask = maskImg
	return nil
}

// AsMask returns an *image.Alpha representing the alpha channel of this
// context. This can be useful for advanced clipping operations where you first
// render the mask geometry and then use it as a mask.
func (dc *Context) AsMask() *image.Alpha {
	// Create alpha mask from im
	width := dc.Width()
	height := dc.Height()
	alpha := image.NewAlpha(image.Rect(0, 0, width, height))

	// Extract alpha channel from im
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			offset := (y*width + x) * 4                // RGBA
			alpha.Pix[y*width+x] = dc.im.Pix[offset+3] // Alpha channel
		}
	}

	return alpha
}

// InvertMask inverts the alpha values in the current clipping mask such that
// a fully transparent region becomes fully opaque and vice versa.
func (dc *Context) InvertMask() {
	if dc.mask == nil {
		return
	}
	// Invert mask data
	for i := range dc.mask.Pix {
		dc.mask.Pix[i] = 255 - dc.mask.Pix[i]
	}
}

// Font Functions - implemented in font.go
// See font.go for: SetFontFace, LoadFontFace, FontHeight, MeasureString,
// MeasureMultilineString, WordWrap, DrawString, DrawStringAnchored, DrawStringWrapped

// Clear fills the entire image with the current color.
func (dc *Context) Clear() {
	// Fill with current color
	r, g, b, a := dc.color.RGBA()
	r8 := uint8(r >> 8)
	g8 := uint8(g >> 8)
	b8 := uint8(b >> 8)
	a8 := uint8(a >> 8)

	for i := 0; i < len(dc.im.Pix); i += 4 {
		dc.im.Pix[i] = r8
		dc.im.Pix[i+1] = g8
		dc.im.Pix[i+2] = b8
		dc.im.Pix[i+3] = a8
	}
}

// SetPixel sets the color of the specified pixel using the current color.
func (dc *Context) SetPixel(x, y int) {
	if x < 0 || y < 0 || x >= dc.Width() || y >= dc.Height() {
		return
	}
	// Set to current color
	r, g, b, a := dc.color.RGBA()
	offset := (y*dc.Width() + x) * 4
	dc.im.Pix[offset] = uint8(r >> 8)
	dc.im.Pix[offset+1] = uint8(g >> 8)
	dc.im.Pix[offset+2] = uint8(b >> 8)
	dc.im.Pix[offset+3] = uint8(a >> 8)
}

// DrawPoint is like DrawCircle but ensures that a circle of the specified
// size is drawn regardless of the current transformation matrix. The position
// is still transformed, but not the shape of the point.
func (dc *Context) DrawPoint(x, y, r float64) {
	// Save current transform
	savedTransform := dc.transform
	// Reset to identity to draw point without transformation
	dc.Identity()
	// Draw circle at the transformed position
	dc.DrawCircle(x, y, r)
	// Restore transform
	dc.transform = savedTransform
}

// DrawLine draws a line segment from (x1,y1) to (x2,y2).
func (dc *Context) DrawLine(x1, y1, x2, y2 float64) {
	dc.MoveTo(x1, y1)
	dc.LineTo(x2, y2)
	dc.Stroke()
}

// DrawRectangle draws a rectangle at (x,y) with the specified width and height.
func (dc *Context) DrawRectangle(x, y, w, h float64) {
	dc.NewSubPath()
	dc.MoveTo(x, y)
	dc.LineTo(x+w, y)
	dc.LineTo(x+w, y+h)
	dc.LineTo(x, y+h)
	dc.ClosePath()
}

// DrawRoundedRectangle draws a rounded rectangle at (x,y) with the specified width, height and corner radius.
func (dc *Context) DrawRoundedRectangle(x, y, w, h, r float64) {
	dc.NewSubPath()
	if w <= 0 || h <= 0 {
		return
	}
	if r <= 0 {
		// Degenerate case - just a rectangle
		dc.DrawRectangle(x, y, w, h)
		return
	}

	// Clamp radius to half the dimensions
	if r > w/2 {
		r = w / 2
	}
	if r > h/2 {
		r = h / 2
	}

	// Move to starting point (right edge of top-left corner)
	dc.MoveTo(x+r, y)

	// Top edge
	dc.LineTo(x+w-r, y)

	// Top-right corner - draw as arc
	dc.drawCornerArc(x+w-r, y+r, r, -math.Pi/2, 0)

	// Right edge
	dc.LineTo(x+w, y+h-r)

	// Bottom-right corner
	dc.drawCornerArc(x+w-r, y+h-r, r, 0, math.Pi/2)

	// Bottom edge
	dc.LineTo(x+r, y+h)

	// Bottom-left corner
	dc.drawCornerArc(x+r, y+h-r, r, math.Pi/2, math.Pi)

	// Left edge
	dc.LineTo(x, y+r)

	// Top-left corner
	dc.drawCornerArc(x+r, y+r, r, math.Pi, 3*math.Pi/2)

	dc.ClosePath()
}

// drawCornerArc draws a quarter-circle arc for rounded rectangle corners
func (dc *Context) drawCornerArc(cx, cy, r, startAngle, endAngle float64) {
	const kappa = 0.5522847498307935

	sx := cx + r*math.Cos(startAngle)
	sy := cy + r*math.Sin(startAngle)
	ex := cx + r*math.Cos(endAngle)
	ey := cy + r*math.Sin(endAngle)

	cp1x := sx - kappa*r*math.Sin(startAngle)
	cp1y := sy + kappa*r*math.Cos(startAngle)
	cp2x := ex + kappa*r*math.Sin(endAngle)
	cp2y := ey - kappa*r*math.Cos(endAngle)

	dc.CubicTo(cp1x, cp1y, cp2x, cp2y, ex, ey)
}

// DrawEllipse draws an ellipse centered at (x,y) with the specified x and y radii.
func (dc *Context) DrawEllipse(x, y, rx, ry float64) {
	dc.NewSubPath()
	if rx <= 0 || ry <= 0 {
		return
	}

	// Approximate ellipse with 4 cubic bezier curves
	// Using kappa constant for circle approximation
	const kappa = 0.5522847498307935
	cx := kappa * rx
	cy := kappa * ry

	x0 := x - rx
	y0 := y - ry
	x1 := x + rx
	y1 := y + ry

	dc.MoveTo(x, y0)
	dc.CubicTo(x+cx, y0, x1, y-cy, x1, y)
	dc.CubicTo(x1, y+cy, x+cx, y1, x, y1)
	dc.CubicTo(x-cx, y1, x0, y+cy, x0, y)
	dc.CubicTo(x0, y-cy, x-cx, y0, x, y0)
	dc.ClosePath()
}

// DrawEllipticalArc draws an elliptical arc centered at (x,y) with the specified radii and angle range.
func (dc *Context) DrawEllipticalArc(x, y, rx, ry, angle1, angle2 float64) {
	dc.NewSubPath()
	if rx <= 0 || ry <= 0 {
		return
	}

	// Normalize angles to [0, 2π)
	angle1 = math.Mod(angle1, 2*math.Pi)
	if angle1 < 0 {
		angle1 += 2 * math.Pi
	}
	angle2 = math.Mod(angle2, 2*math.Pi)
	if angle2 < 0 {
		angle2 += 2 * math.Pi
	}

	// Calculate sweep angle
	sweep := angle2 - angle1
	if sweep < 0 {
		sweep += 2 * math.Pi
	}

	// Approximate with bezier curves
	// Split into segments of at most 90 degrees
	numSegments := int(math.Ceil(sweep / (math.Pi / 2)))
	if numSegments < 1 {
		numSegments = 1
	}
	segmentAngle := sweep / float64(numSegments)

	// Kappa for elliptical arc approximation
	kappa := 4.0 / 3.0 * math.Tan(segmentAngle/4.0)

	for i := 0; i < numSegments; i++ {
		startAngle := angle1 + float64(i)*segmentAngle
		endAngle := startAngle + segmentAngle

		// Start point
		sx := x + rx*math.Cos(startAngle)
		sy := y + ry*math.Sin(startAngle)

		// End point
		ex := x + rx*math.Cos(endAngle)
		ey := y + ry*math.Sin(endAngle)

		// Control points
		cp1x := sx - kappa*rx*math.Sin(startAngle)
		cp1y := sy + kappa*ry*math.Cos(startAngle)
		cp2x := ex + kappa*rx*math.Sin(endAngle)
		cp2y := ey - kappa*ry*math.Cos(endAngle)

		if i == 0 {
			dc.MoveTo(sx, sy)
		}
		dc.CubicTo(cp1x, cp1y, cp2x, cp2y, ex, ey)
	}
}

// DrawArc draws a circular arc centered at (x,y) with the specified radius and angle range.
func (dc *Context) DrawArc(x, y, r, angle1, angle2 float64) {
	dc.DrawEllipticalArc(x, y, r, r, angle1, angle2)
}

// DrawCircle draws a circle centered at (x,y) with the specified radius.
func (dc *Context) DrawCircle(x, y, r float64) {
	dc.DrawEllipse(x, y, r, r)
}

// DrawRegularPolygon draws a regular polygon with n sides, centered at (x,y) with the specified radius and rotation.
func (dc *Context) DrawRegularPolygon(n int, x, y, r, rotation float64) {
	if n < 3 {
		return
	}
	dc.NewSubPath()

	angleStep := 2.0 * math.Pi / float64(n)

	for i := 0; i < n; i++ {
		angle := float64(i)*angleStep + rotation
		pX := x + r*math.Cos(angle)
		pY := y + r*math.Sin(angle)

		if i == 0 {
			dc.MoveTo(pX, pY)
		} else {
			dc.LineTo(pX, pY)
		}
	}
	dc.ClosePath()
}

// DrawImage draws the specified image at the specified point.
func (dc *Context) DrawImage(im image.Image, x, y int) {
	// TODO: Implement using internal packages
	bounds := im.Bounds()
	for py := 0; py < bounds.Dy(); py++ {
		for px := 0; px < bounds.Dx(); px++ {
			r, g, b, a := im.At(px+bounds.Min.X, py+bounds.Min.Y).RGBA()
			dstX := x + px
			dstY := y + py
			if dstX >= 0 && dstX < dc.Width() && dstY >= 0 && dstY < dc.Height() {
				offset := (dstY*dc.Width() + dstX) * 4
				dc.im.Pix[offset] = uint8(r >> 8)
				dc.im.Pix[offset+1] = uint8(g >> 8)
				dc.im.Pix[offset+2] = uint8(b >> 8)
				dc.im.Pix[offset+3] = uint8(a >> 8)
			}
		}
	}
}

// DrawImageAnchored draws the specified image at the specified anchor point.
// The anchor point is x - w * ax, y - h * ay, where w, h is the size of the
// image. Use ax=0.5, ay=0.5 to center the image at the specified point.
func (dc *Context) DrawImageAnchored(im image.Image, x, y int, ax, ay float64) {
	bounds := im.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	anchorX := int(float64(x) - float64(w)*ax)
	anchorY := int(float64(y) - float64(h)*ay)
	dc.DrawImage(im, anchorX, anchorY)
}

// Identity resets the current transformation matrix to the identity matrix.
// This results in no translating, scaling, rotating, or shearing.
func (dc *Context) Identity() {
	dc.transform = path.NewTransformDefault()
}

// Translate updates the current matrix with a translation.
func (dc *Context) Translate(x, y float64) {
	translateTransform := path.NewTransformFromTranslate(float32(x), float32(y))
	dc.transform = dc.transform.PreConcat(translateTransform)
}

// Scale updates the current matrix with a scaling factor.
// Scaling occurs about the origin.
func (dc *Context) Scale(x, y float64) {
	scaleTransform := path.NewTransformFromScale(float32(x), float32(y))
	dc.transform = dc.transform.PreConcat(scaleTransform)
}

// ScaleAbout updates the current matrix with a scaling factor.
// Scaling occurs about the specified point.
func (dc *Context) ScaleAbout(sx, sy, x, y float64) {
	dc.Translate(x, y)
	dc.Scale(sx, sy)
	dc.Translate(-x, -y)
}

// Rotate updates the current matrix with a anticlockwise rotation.
// Rotation occurs about the origin. Angle is specified in radians.
func (dc *Context) Rotate(angle float64) {
	rotateTransform := path.NewTransformFromRotate(float32(angle))
	dc.transform = dc.transform.PreConcat(rotateTransform)
}

// RotateAbout updates the current matrix with a anticlockwise rotation.
// Rotation occurs about the specified point. Angle is specified in radians.
func (dc *Context) RotateAbout(angle, x, y float64) {
	dc.Translate(x, y)
	dc.Rotate(angle)
	dc.Translate(-x, -y)
}

// Shear updates the current matrix with a shearing angle.
// Shearing occurs about the origin.
func (dc *Context) Shear(x, y float64) {
	sx := float32(x)
	sy := float32(y)
	// Use PreConcat to apply shear transform
	shearTransform := path.NewTransformFromSkew(sx, sy)
	dc.transform = dc.transform.PreConcat(shearTransform)
}

// ShearAbout updates the current matrix with a shearing angle.
// Shearing occurs about the specified point.
func (dc *Context) ShearAbout(sx, sy, x, y float64) {
	dc.Translate(x, y)
	dc.Shear(sx, sy)
	dc.Translate(-x, -y)
}

// TransformPoint multiplies the specified point by the current matrix,
// returning a transformed position.
func (dc *Context) TransformPoint(x, y float64) (tx, ty float64) {
	pts := []path.Point{{X: float32(x), Y: float32(y)}}
	dc.transform.MapPoints(pts)
	return float64(pts[0].X), float64(pts[0].Y)
}

// InvertY flips the Y axis so that Y grows from bottom to top and Y=0 is at
// the bottom of the image.
func (dc *Context) InvertY() {
	dc.Translate(0, float64(dc.Height()))
	dc.Scale(1, -1)
}
