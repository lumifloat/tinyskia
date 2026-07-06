// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"image"

	"github.com/lumifloat/tinyskia/internal/core/painter"
	"github.com/lumifloat/tinyskia/internal/core/scan"
	"github.com/lumifloat/tinyskia/internal/path"
)

// DrawImage draws the specified image at the specified point.
func (dc *Context) DrawImage(im image.Image, dx, dy float64) {
	bounds := im.Bounds()
	dc.DrawImageWithSourceRect(im, 0, 0, float64(bounds.Dx()), float64(bounds.Dy()), dx, dy, 0, 0)
}

// DrawImageScaled draws the specified image with scaling.
func (dc *Context) DrawImageScaled(im image.Image, dx, dy, dw, dh float64) {
	bounds := im.Bounds()
	dc.DrawImageWithSourceRect(im, 0, 0, float64(bounds.Dx()), float64(bounds.Dy()), dx, dy, dw, dh)
}

// DrawImageWithSourceRect draws a portion of the specified image with scaling.
func (dc *Context) DrawImageWithSourceRect(im image.Image, sx, sy, sw, sh, dx, dy, dw, dh float64) {
	bounds := im.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	if imgWidth <= 0 || imgHeight <= 0 {
		return
	}

	if sw < 0 {
		sx += sw
		sw = -sw
	}
	if sh < 0 {
		sy += sh
		sh = -sh
	}

	if sx < 0 {
		sx = 0
	}
	if sy < 0 {
		sy = 0
	}
	if sx+sw > float64(imgWidth) {
		sw = float64(imgWidth) - sx
	}
	if sy+sh > float64(imgHeight) {
		sh = float64(imgHeight) - sy
	}

	if sw <= 0 || sh <= 0 {
		return
	}

	sourceX := int(sx)
	sourceY := int(sy)
	sourceWidth := int(sw)
	sourceHeight := int(sh)

	if sourceX < bounds.Min.X {
		sourceX = bounds.Min.X
	}
	if sourceY < bounds.Min.Y {
		sourceY = bounds.Min.Y
	}
	if sourceX+sourceWidth > bounds.Max.X {
		sourceWidth = bounds.Max.X - sourceX
	}
	if sourceY+sourceHeight > bounds.Max.Y {
		sourceHeight = bounds.Max.Y - sourceY
	}

	if sourceWidth <= 0 || sourceHeight <= 0 {
		return
	}

	subImg := image.NewRGBA(image.Rect(0, 0, sourceWidth, sourceHeight))
	for y := 0; y < sourceHeight; y++ {
		for x := 0; x < sourceWidth; x++ {
			subImg.Set(x, y, im.At(sourceX+x, sourceY+y))
		}
	}

	subBounds := subImg.Bounds()
	subWidth := subBounds.Dx()
	subHeight := subBounds.Dy()

	if subWidth <= 0 || subHeight <= 0 {
		return
	}

	if dw <= 0 {
		dw = float64(subWidth)
	}
	if dh <= 0 {
		dh = float64(subHeight)
	}

	translateTransform := path.NewTransformFromTranslate(float32(dx), float32(dy))
	patternShader := imageToPatternShader(subImg, PatternRepeatNone, translateTransform)

	patternShader.Transform(dc.matrix.transform)

	finalTransform := dc.matrix.transform.PreConcat(translateTransform)

	paint := &painter.Paint{
		Shader:          patternShader,
		BlendMode:       dc.composite,
		AntiAlias:       dc.antiAlias,
		Colorspace:      dc.colorspace,
		ForceHQPipeline: dc.forceHQPipeline,
	}
	rectPath := path.NewPathBuilder()
	rect, _ := path.NewRectFromXYWH(0, 0, float32(dw), float32(dh))
	rectPath.PushRect(rect)
	finalPath := rectPath.Finish()

	paint.FillPath(dc.canvas.im, dc.mask, finalPath, finalTransform, scan.FillRuleWinding)
}
