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
)

type LineCap int

const (
	LineCapRound LineCap = iota
	LineCapButt
	LineCapSquare
)

type LineJoin int

const (
	LineJoinMiter LineJoin = iota
	LineJoinMiterClip
	LineJoinRound
	LineJoinBevel
)

type FillRule int

const (
	FillRuleWinding FillRule = iota
	FillRuleEvenOdd
)

type TextAlign int

const (
	// Align to the start edge of the text (left side in left-to-right text, right side in right-to-left text).
	TextAlignStart TextAlign = iota
	// Align to the end edge of the text (right side in left-to-right text, left side in right-to-left text).
	TextAlignEnd
	// Align to the left.
	TextAlignLeft
	// Align to the left.
	TextAlignCenter
	// Align to the center.
	TextAlignRight
)

// Context is the main drawing context, similar to gg.Context.
// It maintains drawing state and provides a canvas-like API.
type Context struct {
	*path2d
	width           int
	height          int
	im              *image.RGBA
	mask            *image.Alpha
	color           color.Color
	fillStyle       Style
	strokeStyle     Style
	dashes          []float64
	dashOffset      float64
	lineWidth       float64
	lineCap         LineCap
	lineJoin        LineJoin
	fillRule        FillRule
	font            *Font
	TextAlign       TextAlign
	transform       path.Transform
	blendMode       BlendMode
	antiAlias       bool
	colorspace      color2.ColorSpace
	forceHQPipeline bool
	stack           []*Context
	contextLost     bool // Tracks if the rendering context was lost
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
		path2d: NewPath2D(),

		width:           bounds.Dx(),
		height:          bounds.Dy(),
		im:              im,
		lineWidth:       1,
		lineCap:         LineCapRound,
		lineJoin:        LineJoinMiter,
		fillRule:        FillRuleWinding,
		transform:       path.NewTransformDefault(),
		blendMode:       BlendModeSourceOver,
		antiAlias:       true,
		colorspace:      color2.ColorSpaceLinear,
		forceHQPipeline: true,
	}
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

// BeginPath resets the current default path.
func (dc *Context) BeginPath() {
	dc.path2d.data = nil
	dc.path2d.builder = path.NewPathBuilder()
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

// SetLineJoinMiter sets the line join to miter.
func (dc *Context) SetLineJoinMiter() {
	dc.lineJoin = LineJoinMiter
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

func (dc *Context) SetForceHQPipeline(force bool) {
	dc.forceHQPipeline = force
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
	dc.SetColor(color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
}

// SetRGB255 sets the current color. r, g, b values should be between 0 and 255,
// inclusive. Alpha will be set to 255 (fully opaque).
func (dc *Context) SetRGB255(r, g, b int) {
	dc.SetRGBA255(r, g, b, 255)
}

// SetRGBA sets the current color. r, g, b, a values should be between 0 and 1,
// inclusive.
func (dc *Context) SetRGBA(r, g, b, a float64) {
	dc.SetColor(color.RGBA{
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

// Fills the subpaths of the current default path with the current fill style, obeying the given fill rule.
func (dc *Context) Fill() {
	dc.FillPath(dc.path2d)
}

// Fills the subpaths of the given path with the current fill style, obeying the given fill rule.
func (dc *Context) FillPath(p *path2d) {
	if p.data == nil {
		p.data = p.builder.Finish()
	}

	var tp *path.Path
	if !dc.transform.IsIdentity() {
		tp = p.data.Transform(dc.transform)
	} else {
		tp = p.data
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
	if dc.antiAlias {
		scan.FillPathAA(tp, int(dc.fillRule), screen, blitter)
	} else {
		scan.FillPath(tp, int(dc.fillRule), screen, blitter)
	}
}

// Stroke the subpaths of the current default path with the current stroke style.
func (dc *Context) Stroke() {
	dc.StrokePath(dc.path2d)
}

// StrokePath the subpaths of the given path with the current stroke style.
func (dc *Context) StrokePath(p *path2d) {
	if p.data == nil {
		p.data = p.builder.Finish()
	}

	var tp *path.Path
	if !dc.transform.IsIdentity() {
		tp = p.data.Transform(dc.transform)
	} else {
		tp = p.data
	}

	stroke := &path.Stroke{
		Width:      float32(dc.lineWidth),
		LineCap:    path.LineCap(dc.lineCap),
		LineJoin:   path.LineJoin(dc.lineJoin),
		MiterLimit: 4.0, // Default miter limit
	}

	if len(dc.dashes) > 0 {
		dashArray := make([]float32, len(dc.dashes))
		for i, d := range dc.dashes {
			dashArray[i] = float32(d)
		}
		stroke.Dash = path.NewStrokeDash(dashArray, float32(dc.dashOffset))
	}

	paint := &Paint{
		Shader:          toShader(dc.strokeStyle, dc.transform),
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

	resScale := path.ComputeResolutionScale(dc.transform)

	stroker := path.NewPathStroker()
	strokedPath := stroker.Stroke(tp, *stroke, resScale)
	if strokedPath == nil {
		return
	}

	if dc.antiAlias {
		scan.FillPathAA(strokedPath, int(FillRuleWinding), screen, blitter)
	} else {
		scan.FillPath(strokedPath, int(FillRuleWinding), screen, blitter)
	}
}

// Further constrains the clipping region to the current default path, using the given fill rule to determine what points are in the path.
func (dc *Context) Clip() {
	dc.ClipPath(dc.path2d)
}

// Further constrains the clipping region to the given path, using the given fill rule to determine what points are in the path.
func (dc *Context) ClipPath(p *path2d) {
	if p.data == nil {
		p.data = p.builder.Finish()
	}

	var tp *path.Path
	if !dc.transform.IsIdentity() {
		tp = p.data.Transform(dc.transform)
	} else {
		tp = p.data
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
	scan.FillPathAA(tp, int(dc.fillRule), screen, blitter)

	// Extract alpha channel from RGBA to Alpha mask
	for i := 0; i < width*height; i++ {
		clipMask.Pix[i] = tempRGBA.Pix[i*4+3] // Alpha channel
	}

	// Intersect with existing mask (take minimum of alpha values)
	if dc.mask == nil {
		dc.mask = clipMask
	} else {
		// Create new mask by intersecting old mask and clip mask
		// Intersection = min(old_alpha, clip_alpha) for each pixel
		mask := image.NewAlpha(image.Rect(0, 0, width, height))
		for i := range mask.Pix {
			clipAlpha := clipMask.Pix[i]
			oldAlpha := dc.mask.Pix[i]
			if clipAlpha < oldAlpha {
				mask.Pix[i] = clipAlpha
			} else {
				mask.Pix[i] = oldAlpha
			}
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

// DrawImage draws the specified image at the specified point.
func (dc *Context) DrawImage(im image.Image, x, y int) {
	// Get image dimensions
	bounds := im.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	if imgWidth <= 0 || imgHeight <= 0 {
		return
	}

	// Match Rust tiny-skia's draw_pixmap approach:
	// 1. Create pattern with translation only
	// 2. Transform both path and shader by dc.transform

	// Step 1: Pattern starts with translation only
	translateTransform := path.NewTransformFromTranslate(float32(x), float32(y))
	patternShader := imageToPatternShader(im, RepeatNone, translateTransform)

	// Step 2: Apply dc.transform to the pattern shader
	patternShader.Transform(dc.transform)

	// For the path, we need to apply BOTH translation and dc.transform
	// Following Rust's approach: transform the rect by the combined transform
	finalTransform := dc.transform.PreConcat(translateTransform)

	// Create paint with the pattern shader
	paint := &Paint{
		Shader:          patternShader,
		BlendMode:       dc.blendMode,
		AntiAlias:       dc.antiAlias,
		Colorspace:      dc.colorspace,
		ForceHQPipeline: dc.forceHQPipeline,
	}

	// Create blitter
	var maskData []uint8
	if dc.mask != nil {
		maskData = dc.mask.Pix
	}
	blitter := paint.blitter(dc.im.Pix, maskData, dc.Width(), dc.Height())
	if blitter == nil {
		return
	}

	// Create a rectangle path for the image bounds
	screen, _ := path.NewScreenIntRectFromXYWH(0, 0, uint32(dc.Width()), uint32(dc.Height()))

	rectPath := path.NewPathBuilder()
	rect, _ := path.NewRectFromXYWH(0, 0, float32(imgWidth), float32(imgHeight))
	rectPath.PushRect(rect)
	finalPath := rectPath.Finish()

	// Transform the path by the final transform
	transformedPath := finalPath.Transform(finalTransform)

	if dc.antiAlias {
		scan.FillPathAA(transformedPath, int(dc.fillRule), screen, blitter)
	} else {
		scan.FillPath(transformedPath, int(dc.fillRule), screen, blitter)
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
	// NewTransformFromRotate expects degrees, so convert from radians
	angleDegrees := angle * 180.0 / math.Pi
	rotateTransform := path.NewTransformFromRotate(float32(angleDegrees))
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
