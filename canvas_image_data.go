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

	sp := image.Pt(
		sx+ctx.canvas.im.Rect.Min.X,
		sy+ctx.canvas.im.Rect.Min.Y,
	)
	sr := image.Rect(sp.X, sp.Y, sp.X+sw, sp.Y+sh)
	ir := sr.Intersect(ctx.canvas.im.Rect)

	im := image.NewRGBA(image.Rect(0, 0, sr.Dx(), sr.Dy()))
	if ir.Empty() {
		return im, nil
	}

	offset := image.Pt(ir.Min.X-sr.Min.X, ir.Min.Y-sr.Min.Y)

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
	if ir0.Empty() {
		return nil
	}

	// 计算要绘制的区域
	dp := image.Pt(
		dx+ctx.canvas.im.Rect.Min.X,
		dy+ctx.canvas.im.Rect.Min.Y,
	)
	dr := image.Rect(dp.X, dp.Y, dp.X+ir0.Dx(), dp.Y+ir0.Dy())
	ir1 := ctx.canvas.im.Rect.Intersect(dr)
	if ir1.Empty() {
		return nil
	}

	offset := image.Pt(
		ir0.Min.X+(ir1.Min.X-dr.Min.X),
		ir0.Min.Y+(ir1.Min.Y-dr.Min.Y),
	)

	iw := ir1.Dx()
	ih := ir1.Dy()
	switch im := im.(type) {
	case *image.RGBA:
		for y := 0; y < ih; y++ {
			s0 := ctx.canvas.im.PixOffset(ir1.Min.X, ir1.Min.Y+y)
			s1 := im.PixOffset(offset.X, offset.Y+y)
			copy(ctx.canvas.im.Pix[s0:s0+iw*4], im.Pix[s1:s1+iw*4])
		}
	default:
		draw.Draw(ctx.canvas.im, ir1, im, offset, draw.Src)
	}
	return nil
}
