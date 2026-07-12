// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"image"
	"image/color"
)

// CreateImageData returns an image.Image object with the same dimensions and color space as the argument.
// All the pixels in the returned object are transparent black.
func (ctx *Context) CreateImageData(sw, sh int) image.Image {
	return image.NewRGBA(image.Rect(0, 0, sw, sh))
}

// CreateImageDataWithImage returns an image.Image object with the same dimensions.
// All the pixels in the returned object are transparent black.
func (ctx *Context) CreateImageDataWithImage(im image.Image) image.Image {
	return image.NewRGBA(im.Bounds())
}

// GetImageData returns an image containing the pixel data for the specified rectangle
func (ctx *Context) GetImageData(sx, sy, sw, sh int) image.Image {
	canvasWidth := ctx.canvas.GetWidth()
	canvasHeight := ctx.canvas.GetHeight()

	if sw <= 0 || sh <= 0 {
		return image.NewRGBA(image.Rect(0, 0, 0, 0))
	}

	if sx < 0 {
		sx = 0
	}
	if sy < 0 {
		sy = 0
	}
	if sx+sw > canvasWidth {
		sw = canvasWidth - sx
	}
	if sy+sh > canvasHeight {
		sh = canvasHeight - sy
	}

	if sw <= 0 || sh <= 0 {
		return image.NewRGBA(image.Rect(0, 0, 0, 0))
	}

	img := image.NewRGBA(image.Rect(0, 0, sw, sh))

	for y := 0; y < sh; y++ {
		for x := 0; x < sw; x++ {
			canvasX := sx + x
			canvasY := sy + y

			c := ctx.canvas.im.At(canvasX, canvasY)
			img.Set(x, y, c)
		}
	}

	return img
}

// PutImageData paints data from the given image onto the bitmap at the specified position
func (ctx *Context) PutImageData(im image.Image, dx, dy int) {
	ctx.PutImageDataWithDirtyRect(im, dx, dy, 0, 0, im.Bounds().Dx(), im.Bounds().Dy())
}

// PutImageDataWithDirtyRect paints data from the given image onto the bitmap, using the dirty rectangle
func (ctx *Context) PutImageDataWithDirtyRect(im image.Image, dx, dy, dirtyX, dirtyY, dirtyWidth, dirtyHeight int) {
	if im == nil {
		return
	}

	bounds := im.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	if imgWidth <= 0 || imgHeight <= 0 {
		return
	}

	canvasWidth := ctx.canvas.GetWidth()
	canvasHeight := ctx.canvas.GetHeight()

	if dirtyX < 0 {
		dirtyX = 0
	}
	if dirtyY < 0 {
		dirtyY = 0
	}
	if dirtyX+dirtyWidth > imgWidth {
		dirtyWidth = imgWidth - dirtyX
	}
	if dirtyY+dirtyHeight > imgHeight {
		dirtyHeight = imgHeight - dirtyY
	}

	if dirtyWidth <= 0 || dirtyHeight <= 0 {
		return
	}

	startX := dx + dirtyX
	startY := dy + dirtyY
	endX := startX + dirtyWidth
	endY := startY + dirtyHeight

	if startX < 0 {
		dirtyX -= startX
		startX = 0
	}
	if startY < 0 {
		dirtyY -= startY
		startY = 0
	}
	if endX > canvasWidth {
		dirtyWidth -= (endX - canvasWidth)
		endX = canvasWidth
	}
	if endY > canvasHeight {
		dirtyHeight -= (endY - canvasHeight)
		endY = canvasHeight
	}

	if dirtyWidth <= 0 || dirtyHeight <= 0 {
		return
	}

	for y := 0; y < dirtyHeight; y++ {
		for x := 0; x < dirtyWidth; x++ {
			imageX := bounds.Min.X + dirtyX + x
			imageY := bounds.Min.Y + dirtyY + y
			canvasX := startX + x
			canvasY := startY + y

			ctx.canvas.im.SetRGBA(canvasX, canvasY, color.RGBAModel.Convert(im.At(imageX, imageY)).(color.RGBA))
		}
	}
}
