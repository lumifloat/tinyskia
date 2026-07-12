// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

// ClearRect clears all pixels on the bitmap in the given rectangle to transparent black.
func (ctx *Context) ClearRect(x, y, w, h float64) {
	x1 := int(x)
	y1 := int(y)
	x2 := int(x + w)
	y2 := int(y + h)

	bounds := ctx.canvas.im.Bounds()
	if x1 < bounds.Min.X {
		x1 = bounds.Min.X
	}
	if y1 < bounds.Min.Y {
		y1 = bounds.Min.Y
	}
	if x2 > bounds.Max.X {
		x2 = bounds.Max.X
	}
	if y2 > bounds.Max.Y {
		y2 = bounds.Max.Y
	}

	for row := y1; row < y2; row++ {
		offset := row * ctx.canvas.im.Stride
		for col := x1; col < x2; col++ {
			pixelIndex := offset + col*4
			ctx.canvas.im.Pix[pixelIndex+0] = 0
			ctx.canvas.im.Pix[pixelIndex+1] = 0
			ctx.canvas.im.Pix[pixelIndex+2] = 0
			ctx.canvas.im.Pix[pixelIndex+3] = 0
		}
	}
}

// FillRect paints the given rectangle onto the bitmap, using the current fill style.
func (ctx *Context) FillRect(x, y, w, h float64) {
	pd := NewPath2D()
	pd.Rect(x, y, w, h)
	ctx.FillPath(pd)
}

// StrokeRect paints the box that outlines the given rectangle onto the bitmap, using the current stroke style.
func (ctx *Context) StrokeRect(x, y, w, h float64) {
	pd := NewPath2D()
	pd.Rect(x, y, w, h)
	ctx.StrokePath(pd)
}
