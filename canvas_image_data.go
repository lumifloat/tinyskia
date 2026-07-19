// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"errors"
	"image"
	"image/draw"
)

// CreateImageData returns an image.Image object with the same dimensions and color space as the argument.
// All the pixels in the returned object are transparent black.
func (ctx *Context) CreateImageData(sw, sh int) (image.Image, error) {
	if sw == 0 || sh == 0 {
		return nil, errors.New("sw and sh must not be 0")
	}
	if sw < 0 {
		sw = -sw
	}
	if sh < 0 {
		sh = -sh
	}
	return image.NewRGBA(image.Rect(0, 0, sw, sh)), nil
}

// CreateImageDataWithImage returns an image.Image object with the same dimensions.
// All the pixels in the returned object are transparent black.
func (ctx *Context) CreateImageDataWithImage(im image.Image) (image.Image, error) {
	if im == nil {
		return nil, errors.New("im must not be nil")
	}
	bounds := im.Bounds()
	if bounds.Empty() {
		return nil, errors.New("im must not be empty")
	}
	w := bounds.Dx()
	h := bounds.Dy()
	return image.NewRGBA(image.Rect(0, 0, w, h)), nil
}

// GetImageData returns an image containing the pixel data for the specified rectangle
func (ctx *Context) GetImageData(sx, sy, sw, sh int) (image.Image, error) {
	if sw == 0 || sh == 0 {
		return nil, errors.New("sw and sh must not be 0")
	}

	// 计算要读取的区域
	sp := image.Pt(
		sx+ctx.canvas.im.Rect.Min.X,
		sy+ctx.canvas.im.Rect.Min.Y,
	)
	sr := image.Rect(sp.X, sp.Y, sp.X+sw, sp.Y+sh)
	ir := sr.Intersect(ctx.canvas.im.Rect)

	im := image.NewRGBA(image.Rect(0, 0, sr.Dx(), sr.Dy()))

	// 选取区域和画布没有交集直接返回透明纯黑
	if ir.Empty() {
		return im, nil
	}

	// 目标区域的偏移量
	offset := image.Pt(
		ir.Min.X-sr.Min.X,
		ir.Min.Y-sr.Min.Y,
	)
	// 需要复制的范围
	iw := ir.Dx()
	ih := ir.Dy()
	for y := 0; y < ih; y++ {
		s0 := ctx.canvas.im.PixOffset(ir.Min.X, ir.Min.Y+y)
		s1 := im.PixOffset(offset.X, offset.Y+y)
		copy(im.Pix[s1:s1+iw*4], ctx.canvas.im.Pix[s0:s0+iw*4])
	}
	return im, nil
}

// PutImageData paints data from the given image onto the bitmap at the specified position
func (ctx *Context) PutImageData(im image.Image, dx, dy int) error {
	return ctx.PutImageDataWithDirtyRect(im, dx, dy, 0, 0, im.Bounds().Dx(), im.Bounds().Dy())
}

// PutImageDataWithDirtyRect paints data from the given image onto the bitmap, using the dirty rectangle
func (ctx *Context) PutImageDataWithDirtyRect(im image.Image, dx, dy, dirtyX, dirtyY, dirtyWidth, dirtyHeight int) error {
	if im == nil {
		return errors.New("im must not be nil")
	}

	// 计算要读取的区域
	sp := image.Pt(
		dirtyX+im.Bounds().Min.X,
		dirtyY+im.Bounds().Min.Y,
	)
	sr := image.Rect(sp.X, sp.Y, sp.X+dirtyWidth, sp.Y+dirtyHeight)
	ir0 := im.Bounds().Intersect(sr)
	// 源区域和源图像没有交集直接返回
	if ir0.Empty() {
		return nil
	}

	// 计算要目标的区域
	dp := image.Pt(
		dx+(ir0.Min.X-im.Bounds().Min.X)+ctx.canvas.im.Rect.Min.X,
		dy+(ir0.Min.Y-im.Bounds().Min.Y)+ctx.canvas.im.Rect.Min.Y,
	)
	dr := image.Rect(dp.X, dp.Y, dp.X+ir0.Dx(), dp.Y+ir0.Dy())
	ir1 := ctx.canvas.im.Rect.Intersect(dr)
	// 目标区域和画布没有交集直接返回
	if ir1.Empty() {
		return nil
	}

	// 源区域的偏移量
	offset := image.Pt(
		ir0.Min.X-(dr.Min.X-ir1.Min.X),
		ir0.Min.Y-(dr.Min.Y-ir1.Min.Y),
	)
	iw := ir1.Dx()
	ih := ir1.Dy()
	switch im := im.(type) {
	case *image.RGBA:
		for y := 0; y < ih; y++ {
			s0 := im.PixOffset(offset.X, offset.Y+y)
			s1 := ctx.canvas.im.PixOffset(ir1.Min.X, ir1.Min.Y+y)
			copy(ctx.canvas.im.Pix[s1:s1+iw*4], im.Pix[s0:s0+iw*4])
		}
	default:
		draw.Draw(ctx.canvas.im, ir1, im, offset, draw.Src)
	}
	return nil
}
