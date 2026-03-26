// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package scan

import (
	"github.com/chewxy/math32"
	"github.com/lumifloat/tinyskia/blitter"
	"github.com/lumifloat/tinyskia/edge"
	"github.com/lumifloat/tinyskia/path"
)

// controls how much we super-sample (when we use that scan conversion)
const supersampleShift uint32 = 2

const shift uint32 = supersampleShift
const scale uint32 = 1 << shift
const mask uint32 = scale - 1

func FillPathAA(
	p *path.Path,
	fillRule int,
	clip path.ScreenIntRect,
	blitter blitter.Blitter,
) {
	// Check for nil or empty path
	if p == nil || p.IsEmpty() {
		return
	}

	// Unlike `path.bounds.to_rect()?.round_out()`,
	// this method rounds out first and then converts into a Rect.
	bounds := p.Bounds()
	ir, ok := path.NewRectFromLTRB(
		math32.Floor(bounds.Left()),
		math32.Floor(bounds.Top()),
		math32.Ceil(bounds.Right()),
		math32.Ceil(bounds.Bottom()),
	)
	if !ok {
		return
	}
	irInt, ok := ir.RoundOut()
	if !ok {
		return
	}

	// If the intersection of the path bounds and the clip bounds
	// will overflow 32767 when << by SHIFT, we can't supersample,
	// so blitter without antialiasing.
	clippedIr, ok := irInt.Intersect(clip.ToIntRect())
	if !ok {
		return
	}

	if rectOverflowsShortShift(clippedIr, int32(shift)) != 0 {
		FillPath(p, fillRule, clip, blitter)
		return
	}

	// Our antialiasing can't handle a clip larger than 32767.
	const maxClipCoord uint32 = 32767
	if clip.Right() > maxClipCoord || clip.Bottom() > maxClipCoord {
		return
	}

	fillPathImpl(p, fillRule, irInt, clip, blitter)
}

func rectOverflowsShortShift(rect path.IntRect, shift int32) int32 {
	return overflowsShortShift(rect.Left(), shift) |
		overflowsShortShift(rect.Top(), shift) |
		overflowsShortShift(rect.Right(), shift) |
		overflowsShortShift(rect.Bottom(), shift)
}

func overflowsShortShift(value int32, shift int32) int32 {
	s := 16 + shift
	return ((value << s) >> s) - value
}

func fillPathImpl(
	p *path.Path,
	fillRule int,
	bounds path.IntRect,
	clip path.ScreenIntRect,
	blitter blitter.Blitter,
) {
	sb, ok := newSuperBlitter(bounds, clip, blitter)
	if !ok {
		return
	}
	// In Go, we must ensure Flush is called since there is no Drop trait.
	defer sb.Flush()

	pathContainedInClip := false
	if sBounds, ok := bounds.ToScreenIntRect(); ok {
		pathContainedInClip = clip.Contains(sBounds)
	}

	FillPathImpl(
		p,
		fillRule,
		clip,
		bounds.Top(),
		bounds.Bottom(),
		int32(shift),
		pathContainedInClip,
		sb,
	)
}

type baseSuperBlitter struct {
	realBlitter blitter.Blitter
	currIy      int32
	width       uint32
	left        uint32
	superLeft   uint32
	currY       int32
	top         int32
}

func newBaseSuperBlitter(
	bounds path.IntRect,
	clipRect path.ScreenIntRect,
	blitter blitter.Blitter,
) (baseSuperBlitter, bool) {
	sectInt, ok := bounds.Intersect(clipRect.ToIntRect())
	if !ok {
		return baseSuperBlitter{}, false
	}
	sect, ok := sectInt.ToScreenIntRect()
	if !ok {
		return baseSuperBlitter{}, false
	}

	return baseSuperBlitter{
		realBlitter: blitter,
		currIy:      int32(sect.Top()) - 1,
		width:       sect.Width(),
		left:        uint32(sect.Left()),
		superLeft:   uint32(sect.Left()) << shift,
		currY:       int32(uint32(sect.Top())<<shift) - 1,
		top:         int32(sect.Top()),
	}, true
}

type superBlitter struct {
	base    baseSuperBlitter
	runs    *edge.AlphaRuns
	offsetX int
}

func newSuperBlitter(
	bounds path.IntRect,
	clipRect path.ScreenIntRect,
	blitter blitter.Blitter,
) (*superBlitter, bool) {
	base, ok := newBaseSuperBlitter(bounds, clipRect, blitter)
	if !ok {
		return nil, false
	}

	return &superBlitter{
		base:    base,
		runs:    edge.NewAlphaRuns(base.width),
		offsetX: 0,
	}, true
}

func (sb *superBlitter) Flush() {
	if sb.base.currIy >= sb.base.top {
		if !sb.runs.IsEmpty() {
			sb.base.realBlitter.BlitAntiH(
				sb.base.left,
				uint32(sb.base.currIy),
				sb.runs.Alpha,
				sb.runs.Runs,
			)
			sb.runs.Reset(sb.base.width)
			sb.offsetX = 0
		}
		sb.base.currIy = sb.base.top - 1
	}
}

func (sb *superBlitter) BlitH(x uint32, y uint32, width uint32) {
	iy := int32(y >> shift)

	if x < sb.base.superLeft {
		width = width - (sb.base.superLeft - x)
		x = 0
	} else {
		x = x - sb.base.superLeft
	}

	if sb.base.currY != int32(y) {
		sb.offsetX = 0
		sb.base.currY = int32(y)
	}

	if iy != sb.base.currIy {
		sb.Flush()
		sb.base.currIy = iy
	}

	start := x
	stop := x + width

	fb := start & mask
	fe := stop & mask
	n := int32(stop>>shift) - int32(start>>shift) - 1

	if n < 0 {
		fb = fe - fb
		n = 0
		fe = 0
	} else {
		if fb == 0 {
			n += 1
		} else {
			fb = scale - fb
		}
	}

	maxValue := uint8((1 << (8 - shift)) - (((y & mask) + 1) >> shift))
	sb.offsetX = sb.runs.Add(
		x>>shift,
		coverageToPartialAlpha(uint32(fb)),
		int(n),
		coverageToPartialAlpha(uint32(fe)),
		maxValue,
		sb.offsetX,
	)
}

// Implement rest of Blitter interface with empty methods or as needed
func (sb *superBlitter) BlitV(x, y, height uint32, alpha uint8)              {}
func (sb *superBlitter) BlitRect(rect path.ScreenIntRect)                    {}
func (sb *superBlitter) BlitAntiH(x, y uint32, alpha []uint8, runs []uint16) {}
func (sb *superBlitter) BlitAntiH2(x, y uint32, alpha0, alpha1 uint8)        {}
func (sb *superBlitter) BlitAntiV2(x, y uint32, alpha0, alpha1 uint8)        {}
func (sb *superBlitter) BlitMask(mask blitter.Mask, clip path.ScreenIntRect) {}

func coverageToPartialAlpha(aa uint32) uint8 {
	aa <<= 8 - 2*shift
	return uint8(aa)
}
