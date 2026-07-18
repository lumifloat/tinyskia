// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"image"
	"image/draw"
	"math"

	"github.com/lumifloat/tinyskia/internal/core/painter"
	"github.com/lumifloat/tinyskia/internal/core/scan"
	"github.com/lumifloat/tinyskia/internal/core/shader"
	"github.com/lumifloat/tinyskia/internal/path"
)

// DrawImage draws the specified image at the specified point.
func (ctx *Context) DrawImage(im image.Image, dx, dy float64) {
	bounds := im.Bounds()
	ctx.DrawImageWithSourceRect(im, float64(bounds.Min.X), float64(bounds.Min.Y),
		float64(bounds.Max.X), float64(bounds.Max.Y), dx, dy, float64(bounds.Dx()), float64(bounds.Dy()))
}

// DrawImageScaled draws the specified image with scaling.
func (ctx *Context) DrawImageScaled(im image.Image, dx, dy, dw, dh float64) {
	bounds := im.Bounds()
	ctx.DrawImageWithSourceRect(im, float64(bounds.Min.X), float64(bounds.Min.Y),
		float64(bounds.Max.X), float64(bounds.Max.Y), dx, dy, dw, dh)
}

// DrawImageWithSourceRect draws a portion of the specified image with scaling.
func (ctx *Context) DrawImageWithSourceRect(im image.Image, sx, sy, sw, sh, dx, dy, dw, dh float64) {
	if sw < 0 {
		sx += sw
		sw = -sw
	}
	if sh < 0 {
		sy += sh
		sh = -sh
	}

	if sw <= 0 || sh <= 0 {
		return
	}

	rect := image.Rect(
		int(0),
		int(0),
		int(math.Ceil(sw)),
		int(math.Ceil(sh)),
	)
	src := image.NewRGBA(rect)
	draw.Draw(src, rect, im, image.Point{int(sx), int(sy)}, draw.Src)

	pattern := shader.NewPattern(
		src,
		shader.SpreadModeNoRepeat,
		shader.FilterQualityBilinear,
		1.0,
		ctx.matrix.transform.
			PreTranslate(float32(dx), float32(dy)).
			PreScale(float32(dw/sw), float32(dh/sh)),
	)

	paint := &painter.Paint{
		Shader:          pattern,
		BlendMode:       composite(ctx.globalCompositeOperation),
		AntiAlias:       ctx.antiAlias,
		Colorspace:      ctx.colorspace,
		ForceHQPipeline: ctx.forceHQPipeline,
	}

	if dw < 0 {
		dx += dw
		dw = -dw
	}
	if dh < 0 {
		dy += dh
		dh = -dh
	}

	if dw <= 0 || dh <= 0 {
		return
	}

	rectPath := path.NewPathBuilder()
	r, _ := path.NewRectFromXYWH(float32(dx), float32(dy), float32(dw), float32(dh))
	rectPath.PushRect(r)
	finalPath := rectPath.Finish()

	paint.FillPath(ctx.canvas.im, ctx.mask, finalPath, ctx.matrix.transform, scan.FillRuleWinding)
}
