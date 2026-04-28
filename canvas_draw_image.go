// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"image"

	"github.com/lumifloat/tinyskia/internal/core/scan"
	"github.com/lumifloat/tinyskia/internal/path"
)

// DrawImage draws the specified image at the specified point.
func (dc *Context) DrawImage(im image.Image, x, y float64) {
	// Get image dimensions
	bounds := im.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	if imgWidth <= 0 || imgHeight <= 0 {
		return
	}

	translateTransform := path.NewTransformFromTranslate(float32(x), float32(y))
	patternShader := imageToPatternShader(im, PatternRepeatNone, translateTransform)

	patternShader.Transform(dc.matrix.transform)

	finalTransform := dc.matrix.transform.PreConcat(translateTransform)

	paint := &Paint{
		Shader:          patternShader,
		BlendMode:       dc.composite,
		AntiAlias:       dc.antiAlias,
		Colorspace:      dc.colorspace,
		ForceHQPipeline: dc.forceHQPipeline,
	}

	var maskData []uint8
	if dc.mask != nil {
		maskData = dc.mask.Pix
	}
	blitter := paint.blitter(dc.im.Pix, maskData, dc.Width(), dc.Height())
	if blitter == nil {
		return
	}

	screen, _ := path.NewScreenIntRectFromXYWH(0, 0, uint32(dc.Width()), uint32(dc.Height()))

	rectPath := path.NewPathBuilder()
	rect, _ := path.NewRectFromXYWH(0, 0, float32(imgWidth), float32(imgHeight))
	rectPath.PushRect(rect)
	finalPath := rectPath.Finish()

	transformedPath := finalPath.Transform(finalTransform)

	if dc.antiAlias {
		scan.FillPathAA(transformedPath, int(FillRuleWinding), screen, blitter)
	} else {
		scan.FillPath(transformedPath, int(FillRuleWinding), screen, blitter)
	}
}
