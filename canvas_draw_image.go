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
	ctx.DrawImageWithSourceRect(im, 0, 0,
		float64(bounds.Dx()), float64(bounds.Dy()), dx, dy, float64(bounds.Dx()), float64(bounds.Dy()))
}

// DrawImageScaled draws the specified image with scaling.
func (ctx *Context) DrawImageScaled(im image.Image, dx, dy, dw, dh float64) {
	bounds := im.Bounds()
	ctx.DrawImageWithSourceRect(im, 0, 0,
		float64(bounds.Dx()), float64(bounds.Dy()), dx, dy, dw, dh)
}

// DrawImageWithSourceRect draws a portion of the specified image with scaling.
func (ctx *Context) DrawImageWithSourceRect(im image.Image, sx, sy, sw, sh, dx, dy, dw, dh float64) {
	// 计算要读取的区域
	sp := image.Pt(
		int(math.Floor(sx))+im.Bounds().Min.X,
		int(math.Floor(sy))+im.Bounds().Min.Y,
	)
	sr := image.Rect(sp.X, sp.Y, sp.X+int(math.Ceil(sw)), sp.Y+int(math.Ceil(sh)))
	ir := sr.Intersect(im.Bounds())
	// 源区域和源图像没有交集直接返回
	if ir.Empty() {
		return
	}

	// 映射到 RGBA 上
	pattern := image.NewRGBA(image.Rect(0, 0, sr.Dx(), sr.Dy()))

	// 目标区域的偏移量
	offset := image.Pt(
		ir.Min.X-sr.Min.X,
		ir.Min.Y-sr.Min.Y,
	)
	switch im := im.(type) {
	case *image.RGBA:
		// 需要复制的范围
		iw := ir.Dx()
		ih := ir.Dy()
		for y := 0; y < ih; y++ {
			s0 := im.PixOffset(ir.Min.X, ir.Min.Y+y)
			s1 := pattern.PixOffset(offset.X, offset.Y+y)
			copy(pattern.Pix[s1:s1+iw*4], im.Pix[s0:s0+iw*4])
		}
	default:
		draw.Draw(pattern, image.Rect(offset.X, offset.Y, offset.X+ir.Dx(), offset.Y+ir.Dy()), im, image.Pt(ir.Min.X, ir.Min.Y), draw.Src)
	}

	// 计算要绘制的区域
	dr, ok := path.NewRectFromXYWHStable(float32(dx), float32(dy), float32(dw), float32(dh))
	if !ok {
		return
	}
	builder := path.NewPathBuilder()
	builder.PushRect(dr)
	rect := builder.Finish()

	// 绘制
	paint := &painter.Paint{
		Shader: shader.NewPattern(
			pattern,
			shader.SpreadModeNoRepeat,
			quality(ctx.imageSmoothingQuality),
			float32(ctx.globalAlpha),
			ctx.matrix.transform.
				PreTranslate(float32(dx), float32(dy)).
				PreScale(float32(dw)/float32(sr.Dx()), float32(dh)/float32(sr.Dy())),
		),
		BlendMode:       composite(ctx.globalCompositeOperation),
		AntiAlias:       ctx.antiAlias,
		Colorspace:      ctx.colorspace,
		ForceHQPipeline: ctx.forceHQPipeline,
	}
	paint.FillPath(ctx.canvas.im, ctx.mask, rect, ctx.matrix.transform, scan.FillRuleWinding)
}
