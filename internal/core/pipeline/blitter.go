// Copyright 2016 Google Inc.
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package pipeline

import (
	"image"
	"image/color"
	"unsafe"
)

const BYTES_PER_PIXEL = 4

type RasterPipelineBlitter struct {
	Mask             *image.Alpha
	Src              *image.RGBA
	Dst              *image.RGBA
	Memset2dColor    color.RGBA
	UseMemset2dColor bool
	BlitAntiHRp      RasterPipeline
	BlitRectRp       RasterPipeline
	BlitMaskRp       RasterPipeline
	IsMask           bool
}

func (b *RasterPipelineBlitter) BlitH(x, y uint32, width uint32) {
	r := image.Rect(int(x), int(y), int(x+width), int(y+1))
	b.BlitRect(r)
}

func (b *RasterPipelineBlitter) BlitAntiH(x, y uint32, aa []uint8, runs []uint16) {
	aaOffset := 0
	runOffset := 0
	for runOffset < len(runs) {
		run := runs[runOffset]
		if run == 0 {
			break
		}
		width := uint32(run)

		alpha := aa[aaOffset]
		switch alpha {
		case 0x00: // Transparent
			// Do nothing
		case 0xFF: // Opaque
			b.BlitH(x, y, width)
		default:
			b.BlitAntiHRp.Ctx.CurrentCoverage = float32(alpha) * (1.0 / 255.0)
			rect := image.Rect(int(x), int(y), int(x+width), int(y+1))
			b.BlitAntiHRp.Run(rect, &AAMaskCtx{}, b.Mask, b.Src, b.Dst)
		}

		x += width
		runOffset += int(run)
		aaOffset += int(run)
	}
}

func (b *RasterPipelineBlitter) BlitV(x, y, height uint32, alpha uint8) {
	rect := image.Rect(int(x), int(y), int(x)+1, int(y+height))
	mask := image.NewAlpha(rect)
	mask.Pix[0] = alpha
	mask.Pix[1] = alpha
	b.BlitMask(mask, rect)
}

func (b *RasterPipelineBlitter) BlitAntiH2(x, y uint32, alpha0, alpha1 uint8) {
	rect := image.Rect(int(x), int(y), int(x)+2, int(y)+1)
	mask := image.NewAlpha(rect)
	mask.Pix[0] = alpha0
	mask.Pix[1] = alpha1
	b.BlitMask(mask, rect)
}

func (b *RasterPipelineBlitter) BlitAntiV2(x, y uint32, alpha0, alpha1 uint8) {
	rect := image.Rect(int(x), int(y), int(x)+1, int(y)+2)
	mask := image.NewAlpha(rect)
	mask.Pix[0] = alpha0
	mask.Pix[1] = alpha1
	b.BlitMask(mask, rect)
}

func (b *RasterPipelineBlitter) BlitRect(rect image.Rectangle) {
	if b.UseMemset2dColor {
		if b.IsMask {
			mask := (*image.Alpha)(unsafe.Pointer(b.Dst))
			width := rect.Dx()
			for y := rect.Min.Y; y < rect.Max.Y; y++ {
				idx := mask.PixOffset(rect.Min.X, y)
				data := mask.Pix[idx : idx+width]
				for i := range data {
					data[i] = b.Memset2dColor.A
				}
			}
		} else {
			width := rect.Dx() * BYTES_PER_PIXEL
			for y := rect.Min.Y; y < rect.Max.Y; y++ {
				idx := b.Dst.PixOffset(rect.Min.X, y)
				data := b.Dst.Pix[idx : idx+width]
				for i := 0; i < len(data); i += BYTES_PER_PIXEL {
					data[i+0] = b.Memset2dColor.R
					data[i+1] = b.Memset2dColor.G
					data[i+2] = b.Memset2dColor.B
					data[i+3] = b.Memset2dColor.A
				}
			}
		}
		return
	}

	b.BlitRectRp.Run(rect, &AAMaskCtx{}, b.Mask, b.Src, b.Dst)
}

func (b *RasterPipelineBlitter) BlitMask(mask *image.Alpha, clip image.Rectangle) {
	aaMaskCtx := AAMaskCtx{
		Pixels: [2]uint8{mask.Pix[0], mask.Pix[1]},
		Stride: uint32(mask.Stride),
		Shift:  int(mask.Rect.Min.X + mask.Rect.Min.Y*mask.Stride),
	}

	b.BlitMaskRp.Run(clip, &aaMaskCtx, b.Mask, b.Src, b.Dst)
}
