// Copyright 2016 Google Inc.
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package pipeline

import (
	"github.com/lumifloat/tinyskia/blitter"
	"github.com/lumifloat/tinyskia/color"
	"github.com/lumifloat/tinyskia/path"
)

const BYTES_PER_PIXEL = 4

type RasterPipelineBlitter struct {
	Mask          *MaskCtx
	PixmapSrc     *PixmapCtx
	Pixmap        *SubPixmapCtx
	Memset2dColor *color.PremultipliedColorU8
	BlitAntiHRp   RasterPipeline
	BlitRectRp    RasterPipeline
	BlitMaskRp    RasterPipeline
	IsMask        bool
}

func (b *RasterPipelineBlitter) BlitH(x, y uint32, width uint32) {
	r := path.NewScreenIntRectFromXYWHSafe(x, y, width, 1)
	b.BlitRect(r)
}

func (b *RasterPipelineBlitter) BlitAntiH(x, y uint32, aa []uint8, runs []uint16) {
	var maskCtx MaskCtx
	if b.Mask != nil {
		maskCtx = MaskCtx{
			Data:      b.Mask.Data,
			RealWidth: b.Mask.RealWidth,
		}
	}

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
		case color.AlphaU8Transparent:
			// Do nothing
		case color.AlphaU8Opaque:
			b.BlitH(x, y, width)
		default:
			b.BlitAntiHRp.Ctx.CurrentCoverage = float32(alpha) * (1.0 / 255.0)
			rect := path.NewScreenIntRectFromXYWHSafe(x, y, width, 1)
			b.BlitAntiHRp.Run(&rect, &AAMaskCtx{}, &maskCtx, b.PixmapSrc, b.Pixmap)
		}

		x += width
		runOffset += int(run)
		aaOffset += int(run)
	}
}

func (b *RasterPipelineBlitter) BlitV(x, y, height uint32, alpha uint8) {
	bounds := path.NewScreenIntRectFromXYWHSafe(x, y, 1, height)
	mask := blitter.Mask{
		Image:    [2]uint8{alpha, alpha},
		Bounds:   bounds,
		RowBytes: 0,
	}
	b.BlitMask(mask, bounds)
}

func (b *RasterPipelineBlitter) BlitAntiH2(x, y uint32, alpha0, alpha1 uint8) {
	bounds, _ := path.NewScreenIntRectFromXYWH(x, y, 2, 1)
	mask := blitter.Mask{
		Image:    [2]uint8{alpha0, alpha1},
		Bounds:   bounds,
		RowBytes: 2,
	}
	b.BlitMask(mask, bounds)
}

func (b *RasterPipelineBlitter) BlitAntiV2(x, y uint32, alpha0, alpha1 uint8) {
	bounds, _ := path.NewScreenIntRectFromXYWH(x, y, 1, 2)
	mask := blitter.Mask{
		Image:    [2]uint8{alpha0, alpha1},
		Bounds:   bounds,
		RowBytes: 1,
	}
	b.BlitMask(mask, bounds)
}

func (b *RasterPipelineBlitter) BlitRect(rect path.ScreenIntRect) {
	if b.Memset2dColor != nil {
		c := *b.Memset2dColor
		if b.IsMask {
			alpha := c.Alpha()
			for y := uint32(0); y < rect.Height(); y++ {
				start := (int(rect.Y()+y)*b.Pixmap.RealWidth + int(rect.X())) * BYTES_PER_PIXEL
				end := start + int(rect.Width())*BYTES_PER_PIXEL
				data := b.Pixmap.Data[start:end]
				for i := range data {
					data[i] = alpha
				}
			}
		} else {
			for y := uint32(0); y < rect.Height(); y++ {
				start := (int(rect.Y()+y)*b.Pixmap.RealWidth + int(rect.X())) * BYTES_PER_PIXEL
				end := start + int(rect.Width())*BYTES_PER_PIXEL
				data := b.Pixmap.Data[start:end]
				for i := 0; i < len(data); i += BYTES_PER_PIXEL {
					data[i+0] = c.Red()
					data[i+1] = c.Green()
					data[i+2] = c.Blue()
					data[i+3] = c.Alpha()
				}
			}
		}
		return
	}

	var maskCtx MaskCtx
	if b.Mask != nil {
		maskCtx = MaskCtx{
			Data:      b.Mask.Data,
			RealWidth: b.Mask.RealWidth,
		}
	}

	b.BlitRectRp.Run(&rect, &AAMaskCtx{}, &maskCtx, b.PixmapSrc, b.Pixmap)
}

func (b *RasterPipelineBlitter) BlitMask(mask blitter.Mask, clip path.ScreenIntRect) {
	aaMaskCtx := AAMaskCtx{
		Pixels: mask.Image,
		Stride: mask.RowBytes,
		Shift:  int(mask.Bounds.Left() + mask.Bounds.Top()*mask.RowBytes),
	}

	var maskCtx MaskCtx
	if b.Mask != nil {
		maskCtx = MaskCtx{
			Data:      b.Mask.Data,
			RealWidth: b.Mask.RealWidth,
		}
	}

	b.BlitMaskRp.Run(&clip, &aaMaskCtx, &maskCtx, b.PixmapSrc, b.Pixmap)
}
