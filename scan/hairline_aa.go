// Copyright 2011 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package scan

import (
	"github.com/lumifloat/tinyskia/blitter"
	"github.com/lumifloat/tinyskia/edge"
	"github.com/lumifloat/tinyskia/internal/fixed"
	"github.com/lumifloat/tinyskia/path"
)

type fixedRect struct {
	left   fixed.FDot16
	top    fixed.FDot16
	right  fixed.FDot16
	bottom fixed.FDot16
}

func newFixedRectFromRect(src path.Rect) fixedRect {
	return fixedRect{
		left:   fixed.NewFDot16FromF32(src.Left()),
		top:    fixed.NewFDot16FromF32(src.Top()),
		right:  fixed.NewFDot16FromF32(src.Right()),
		bottom: fixed.NewFDot16FromF32(src.Bottom()),
	}
}

// alphaMul multiplies value by 0..256, and shifts the result down 8
// (i.e. return (value * alpha256) >> 8)
func alphaMul(value uint8, alpha256 int32) uint8 {
	a := (int32(value) * alpha256) >> 8
	return uint8(a)
}

// FillRectAA fills a rectangle with anti-aliasing
func FillRectAA(rect path.Rect, clip path.ScreenIntRect, blitter blitter.Blitter) {
	r, ok := rect.Intersect(clip.ToRect())
	if !ok {
		return // everything was clipped out
	}

	fr := newFixedRectFromRect(r)
	fillFixedRect(fr, blitter)
}

func fillFixedRect(rect fixedRect, blitter blitter.Blitter) {
	fillDot8(
		fixed.NewFDot8FromFDot16(rect.left),
		fixed.NewFDot8FromFDot16(rect.top),
		fixed.NewFDot8FromFDot16(rect.right),
		fixed.NewFDot8FromFDot16(rect.bottom),
		true,
		blitter,
	)
}

func fillDot8(l, t, r, b fixed.FDot8, fillInner bool, blitter blitter.Blitter) {
	toAlpha := func(a int32) uint8 {
		return uint8(a)
	}

	// check for empty now that we're in our reduced precision space
	if l >= r || t >= b {
		return
	}

	top := int32(t >> 8)
	if top == ((b - 1) >> 8) {
		// just one scanline high
		doScanline(l, top, r, toAlpha(int32(b-t-1)), blitter)
		return
	}

	if t&0xFF != 0 {
		doScanline(l, top, r, toAlpha(int32(256-(t&0xFF))), blitter)
		top += 1
	}

	bottom := int32(b >> 8)
	height := bottom - top
	if height > 0 {
		h64 := uint32(height)
		left := int32(l >> 8)
		if left == ((r - 1) >> 8) {
			// just 1-pixel wide
			blitter.BlitV(uint32(left), uint32(top), h64, toAlpha(int32(r-l-1)))
		} else {
			if l&0xFF != 0 {
				blitter.BlitV(uint32(left), uint32(top), h64, toAlpha(int32(256-(l&0xFF))))
				left += 1
			}

			right := int32(r >> 8)
			width := right - left
			if fillInner {
				if width > 0 {
					rect, ok := path.NewScreenIntRectFromXYWH(uint32(left), uint32(top), uint32(width), h64)
					if ok {
						blitter.BlitRect(rect)
					}
				}
			}

			if r&0xFF != 0 {
				blitter.BlitV(uint32(right), uint32(top), h64, toAlpha(int32(r&0xFF)))
			}
		}
	}

	if b&0xFF != 0 {
		doScanline(l, bottom, r, toAlpha(int32(b&0xFF)), blitter)
	}
}

func doScanline(l fixed.FDot8, top int32, r fixed.FDot8, alpha uint8, blitter blitter.Blitter) {
	if l >= r {
		return
	}

	// Convert top to uint32 like Rust's match u32::try_from(top)
	if top < 0 {
		return
	}
	t64 := uint32(top)

	const oneLen = 1

	if (l >> 8) == ((r - 1) >> 8) {
		// 1x1 pixel
		leftVal := int32(l >> 8)
		if leftVal >= 0 {
			blitter.BlitV(uint32(leftVal), t64, oneLen, alphaMul(alpha, int32(r-l)))
		}
		return
	}

	left := int32(l >> 8)

	if l&0xFF != 0 {
		if left >= 0 {
			blitter.BlitV(uint32(left), t64, oneLen, alphaMul(alpha, int32(256-(l&0xFF))))
		}
		left += 1
	}

	right := int32(r >> 8)
	width := right - left
	if width > 0 {
		if left >= 0 {
			callHLineBlitter(uint32(left), true, t64, uint32(width), alpha, blitter)
		}
	}

	if r&0xFF != 0 {
		if right >= 0 {
			blitter.BlitV(uint32(right), t64, oneLen, alphaMul(alpha, int32(r&0xFF)))
		}
	}
}

func callHLineBlitter(x uint32, hasY bool, y uint32, count uint32, alpha uint8, blitter blitter.Blitter) {
	const hlineStackBuffer = 100

	runs := make([]uint16, hlineStackBuffer+1)
	aa := make([]uint8, hlineStackBuffer)

	for count > 0 {
		aa[0] = alpha

		n := count
		if n > hlineStackBuffer {
			n = hlineStackBuffer
		}

		// Ensure n fits in uint16 for runs array
		if n > 65535 {
			n = 65535
		}

		runs[0] = uint16(n)
		runs[n] = 0 // terminator
		if hasY {
			blitter.BlitAntiH(x, y, aa, runs)
		}
		x += n

		// Match Rust's break condition: if n >= count || count == 0
		if n >= count {
			break
		}

		count -= n
	}
}

// strokePathAA strokes a path with anti-aliasing
func StrokePathAA(p *path.Path, lineCap path.LineCap, clip path.ScreenIntRect, hasClip bool, blitter blitter.Blitter) {
	strokePathImpl(p, lineCap, clip, hasClip, antiHairLineRgn, blitter)
}

// antiHairLineRgn renders antialiased lines in a region
func antiHairLineRgn(points []path.Point, clip path.ScreenIntRect, hasClip bool, blitter blitter.Blitter) {
	const max = 32767.0
	fixedBounds, _ := path.NewRectFromLTRB(-max, -max, max, max)

	var clipBounds path.Rect
	if hasClip {
		// Antialiased hairlines can draw up to 1/2 of a pixel outside of
		// their bounds, so we need to outset the clip before calling the
		// clipper. To make the numerics safer, we outset by a whole pixel.
		clipBounds, _ = clip.ToRect().Outset(1.0, 1.0)
	}

	for i := 0; i < len(points)-1; i++ {
		pts := [2]path.Point{points[i], points[i+1]}

		// We have to pre-clip the line to fit in a Fixed, so we just chop the line.
		if !path.Intersect(pts, fixedBounds, &pts) {
			continue
		}

		if hasClip {
			tmp := pts
			if !path.Intersect(tmp, clipBounds, &pts) {
				continue
			}
		}

		x0 := fixed.NewFDot6FromF32(pts[0].X)
		y0 := fixed.NewFDot6FromF32(pts[0].Y)
		x1 := fixed.NewFDot6FromF32(pts[1].X)
		y1 := fixed.NewFDot6FromF32(pts[1].Y)

		if hasClip {
			// Compute line bounds in fixed point
			left := x0
			if x1 < left {
				left = x1
			}
			top := y0
			if y1 < top {
				top = y1
			}
			rightVal := x0
			if x1 > rightVal {
				rightVal = x1
			}
			bottom := y0
			if y1 > bottom {
				bottom = y1
			}

			// Convert to integer rect like Rust does
			ir, ok := path.NewIntRectFromLTRB(
				int32(fixed.FDot6Floor(left))-1,
				int32(fixed.FDot6Floor(top))-1,
				int32(fixed.FDot6Ceil(rightVal))+1,
				int32(fixed.FDot6Ceil(bottom))+1,
			)
			if !ok {
				// Match Rust: return early if IntRect creation fails
				return
			}

			// Check intersection with clip
			intersected, ok := ir.Intersect(clip.ToIntRect())
			if !ok {
				continue
			}

			// If not fully contained, use subclip
			if !clip.ToIntRect().Contains(intersected) {
				subclip, ok := intersected.ToScreenIntRect()
				if ok {
					doAntiHairline(x0, y0, x1, y1, subclip.ToIntRect(), true, blitter)
				}
				continue
			}

			// fall through to no-clip case (line is fully inside clip)
		}

		doAntiHairline(x0, y0, x1, y1, path.IntRect{}, false, blitter)
	}
}

type blitterKind fixed.FDot6

const (
	blitterKindHLine blitterKind = iota
	blitterKindHorish
	blitterKindVLine
	blitterKindVertish
)

// doAntiHairline draws an antialiased hairline
func doAntiHairline(x0, y0, x1, y1 fixed.FDot6, clipOpt path.IntRect, hasClip bool, blitter blitter.Blitter) {
	if anyBadInts(x0, y0, x1, y1) {
		return
	}

	if fixed.FDot6Abs(x1-x0) > fixed.NewFDot6FromI32(511) ||
		fixed.FDot6Abs(y1-y0) > fixed.NewFDot6FromI32(511) {
		hx := (x0 >> 1) + (x1 >> 1)
		hy := (y0 >> 1) + (y1 >> 1)
		doAntiHairline(x0, y0, hx, hy, clipOpt, hasClip, blitter)
		doAntiHairline(hx, hy, x1, y1, clipOpt, hasClip, blitter)
		return
	}

	var scaleStart, scaleStop int32
	var istart, istop, fstart, slope int32
	var kind blitterKind

	if fixed.FDot6Abs(x1-x0) > fixed.FDot6Abs(y1-y0) {
		if x0 > x1 {
			x0, x1 = x1, x0
			y0, y1 = y1, y0
		}

		istart = fixed.FDot6Floor(x0)
		istop = fixed.FDot6Ceil(x1)
		fstart = fixed.FDot6ToFDot16(y0)
		if y0 == y1 {
			// completely horizontal, take fast case
			slope = 0
			kind = blitterKindHLine
		} else {
			slope = fixed.FDot6DivToFDot16(y1-y0, x1-x0)
			fstart += (slope*int32(32-(x0&63)) + 32) >> 6
			kind = blitterKindHorish
		}

		if istop-istart == 1 {
			scaleStart = x1 - x0
			scaleStop = 0
		} else {
			scaleStart = 64 - (x0 & 63)
			scaleStop = x1 & 63
		}

		if hasClip {
			clip := clipOpt
			if istart >= clip.Right() || istop <= clip.Left() {
				return
			}

			if istart < clip.Left() {
				fstart += slope * (clip.Left() - istart)
				istart = clip.Left()
				scaleStart = 64
				if istop-istart == 1 {
					scaleStart = contribution64(fixed.FDot6(x1))
					scaleStop = 0
				}
			}

			if istop > clip.Right() {
				istop = clip.Right()
				scaleStop = 0
			}

			if istart == istop {
				return
			}

			var top, bottom int32
			if slope >= 0 {
				// T2B
				top = fixed.FDot16FloorToI32(fstart - 32768)
				bottom = fixed.FDot16CeilToI32(fstart + (istop-istart-1)*slope + 32768)
			} else {
				// B2T
				bottom = fixed.FDot16CeilToI32(fstart + 32768)
				top = fixed.FDot16FloorToI32(fstart + (istop-istart-1)*slope - 32768)
			}

			top -= 1
			bottom += 1

			if top >= clip.Bottom() || bottom <= clip.Top() {
				return
			}

			if clip.Top() <= top && clip.Bottom() >= bottom {
				hasClip = false
			}
		}
	} else {
		if y0 > y1 {
			x0, x1 = x1, x0
			y0, y1 = y1, y0
		}

		istart = fixed.FDot6Floor(y0)
		istop = fixed.FDot6Ceil(y1)
		fstart = fixed.FDot6ToFDot16(x0)
		if x0 == x1 {
			if y0 == y1 {
				// are we zero length? nothing to do
				return
			}
			// completely vertical, take fast case
			slope = 0
			kind = blitterKindVLine
		} else {
			slope = fixed.FDot6DivToFDot16(x1-x0, y1-y0)
			fstart += (slope*int32(32-(y0&63)) + 32) >> 6
			kind = blitterKindVertish
		}

		if istop-istart == 1 {
			scaleStart = y1 - y0
			scaleStop = 0
		} else {
			scaleStart = 64 - (y0 & 63)
			scaleStop = y1 & 63
		}

		if hasClip {
			clip := clipOpt
			if istart >= clip.Bottom() || istop <= clip.Top() {
				return
			}

			if istart < clip.Top() {
				fstart += slope * (clip.Top() - istart)
				istart = clip.Top()
				scaleStart = 64
				if istop-istart == 1 {
					scaleStart = contribution64(fixed.FDot6(y1))
					scaleStop = 0
				}
			}
			if istop > clip.Bottom() {
				istop = clip.Bottom()
				scaleStop = 0
			}

			if istart == istop {
				return
			}

			var left, right int32
			if slope >= 0 {
				// L2R
				left = fixed.FDot16FloorToI32(fstart - 32768)
				right = fixed.FDot16CeilToI32(fstart + (istop-istart-1)*slope + 32768)
			} else {
				// R2L
				right = fixed.FDot16CeilToI32(fstart + 32768)
				left = fixed.FDot16FloorToI32(fstart + (istop-istart-1)*slope - 32768)
			}

			left -= 1
			right += 1

			if left >= clip.Right() || right <= clip.Left() {
				return
			}

			if clip.Left() <= left && clip.Right() >= right {
				hasClip = false
			}
		}
	}

	activeBlitter := blitter
	if hasClip {
		// Convert IntRect to ScreenIntRect for rectClipBlitter
		screenClip, ok := clipOpt.ToScreenIntRect()
		if !ok {
			return
		}
		activeBlitter = &rectClipBlitter{blitter: blitter, clip: screenClip}
	}

	var hair antiHairBlitter
	switch kind {
	case blitterKindHLine:
		hair = &hLineAntiHairBlitter{activeBlitter}
	case blitterKindHorish:
		hair = &horishAntiHairBlitter{activeBlitter}
	case blitterKindVLine:
		hair = &vLineAntiHairBlitter{activeBlitter}
	case blitterKindVertish:
		hair = &vertishAntiHairBlitter{activeBlitter}
	}

	if istart < 0 || istop < 0 {
		return
	}

	uStart := uint32(istart)
	uStop := uint32(istop)

	fstart = hair.drawCap(uStart, fixed.FDot16(fstart), fixed.FDot16(slope), scaleStart)
	uStart += 1

	stopOffset := uint32(0)
	if scaleStop > 0 {
		stopOffset = 1
	}

	fullSpans := uStop - uStart - stopOffset
	if fullSpans > 0 {
		fstart = hair.drawLine(uStart, uStart+fullSpans, fixed.FDot16(fstart), fixed.FDot16(slope))
	}

	if scaleStop > 0 {
		hair.drawCap(uStop-1, fixed.FDot16(fstart), fixed.FDot16(slope), scaleStop)
	}
}

// anyBadInts checks for invalid integer values
func anyBadInts(a, b, c, d fixed.FDot6) bool {
	bad := func(x int64) int64 { return x & -x }
	return ((bad(int64(a)) | bad(int64(b)) | bad(int64(c)) | bad(int64(d))) >> 63) != 0
}

// contribution64 computes the fractional part of ordinate
func contribution64(ordinate int32) int32 {
	return ((ordinate - 1) & 63) + 1
}

// antiHairBlitter interface for different blitter types
type antiHairBlitter interface {
	drawCap(x uint32, fy fixed.FDot16, slope fixed.FDot16, mod64 int32) fixed.FDot16
	drawLine(x uint32, stopx uint32, fy fixed.FDot16, slope fixed.FDot16) fixed.FDot16
}

// hLineAntiHairBlitter for horizontal lines
type hLineAntiHairBlitter struct{ blitter blitter.Blitter }

func (h *hLineAntiHairBlitter) drawCap(x uint32, fy fixed.FDot16, _ fixed.FDot16, mod64 int32) fixed.FDot16 {
	fy += 32768 // HALF
	if fy < 0 {
		fy = 0
	}
	y := fy >> 16
	a := uint8(fy >> 8)

	ma := fixed.FDot6SmallScale(a, mod64)
	if ma != 0 {
		callHLineBlitter(x, true, uint32(y), 1, ma, h.blitter)
	}

	ma = fixed.FDot6SmallScale(255-a, mod64)
	if ma != 0 {
		if int32(y) > 0 {
			callHLineBlitter(x, true, uint32(int32(y)-1), 1, ma, h.blitter)
		}
	}
	return fy - 32768
}

func (h *hLineAntiHairBlitter) drawLine(x uint32, stopX uint32, fy fixed.FDot16, _ fixed.FDot16) fixed.FDot16 {
	count := stopX - x
	if count <= 0 {
		return fy
	}
	fy += 32768
	if fy < 0 {
		fy = 0
	}
	y := fy >> 16
	a := uint8(fy >> 8)

	if a != 0 {
		callHLineBlitter(x, true, uint32(y), count, a, h.blitter)
	}
	if 255-a != 0 {
		if int32(y) > 0 {
			callHLineBlitter(x, true, uint32(int32(y)-1), count, 255-a, h.blitter)
		}
	}
	return fy - 32768
}

// horishAntiHairBlitter for mostly horizontal lines
type horishAntiHairBlitter struct{ blitter blitter.Blitter }

func (h *horishAntiHairBlitter) drawCap(x uint32, fy fixed.FDot16, slope fixed.FDot16, mod64 int32) fixed.FDot16 {
	fy += 32768 // HALF
	if fy < 0 {
		fy = 0
	}

	lowerY := fy >> 16
	a := i32ToAlpha(fy >> 8)
	a0 := fixed.FDot6SmallScale(255-a, mod64)
	a1 := fixed.FDot6SmallScale(a, mod64)
	vY := uint32(lowerY)
	if vY < 1 {
		vY = 1
	}
	h.blitter.BlitAntiV2(x, vY-1, a0, a1)

	return fy + slope - 32768
}

func (h *horishAntiHairBlitter) drawLine(x uint32, stopX uint32, fy fixed.FDot16, slope fixed.FDot16) fixed.FDot16 {
	count := stopX - x
	if count <= 0 {
		return fy
	}

	fy += 32768 // HALF
	loopX := x
	for {
		if fy < 0 {
			fy = 0
		}
		lowerY := fy >> 16
		a := i32ToAlpha(fy >> 8)
		vY := uint32(lowerY)
		if vY < 1 {
			vY = 1
		}
		h.blitter.BlitAntiV2(loopX, vY-1, 255-a, a)
		fy += slope
		loopX++
		if loopX >= stopX {
			break
		}
	}

	return fy - 32768
}

// vLineAntiHairBlitter for vertical lines
type vLineAntiHairBlitter struct{ blitter blitter.Blitter }

func (v *vLineAntiHairBlitter) drawCap(y uint32, fx fixed.FDot16, _ fixed.FDot16, mod64 int32) fixed.FDot16 {
	fx += 32768 // HALF
	if fx < 0 {
		fx = 0
	}
	xCoord := fx >> 16
	a := uint8(fx >> 8)

	ma := fixed.FDot6SmallScale(a, mod64)
	if ma != 0 {
		v.blitter.BlitV(uint32(xCoord), y, 1, ma)
	}

	ma = fixed.FDot6SmallScale(255-a, mod64)
	if ma != 0 {
		xVal := uint32(0)
		if int32(xCoord) > 0 {
			xVal = uint32(int32(xCoord) - 1)
		}
		v.blitter.BlitV(xVal, y, 1, ma)
	}

	return fx - 32768
}

func (v *vLineAntiHairBlitter) drawLine(y uint32, stopY uint32, fx fixed.FDot16, _ fixed.FDot16) fixed.FDot16 {
	count := stopY - y
	if count <= 0 {
		return fx
	}

	fx += 32768 // HALF
	if fx < 0 {
		fx = 0
	}
	xCoord := fx >> 16
	a := uint8(fx >> 8)

	if a != 0 {
		v.blitter.BlitV(uint32(xCoord), y, count, a)
	}

	a = 255 - a
	if a != 0 {
		xVal := uint32(0)
		if int32(xCoord) > 0 {
			xVal = uint32(int32(xCoord) - 1)
		}
		v.blitter.BlitV(xVal, y, count, a)
	}

	return fx - 32768
}

// vertishAntiHairBlitter for mostly vertical lines
type vertishAntiHairBlitter struct{ blitter blitter.Blitter }

func (v *vertishAntiHairBlitter) drawCap(y uint32, fx fixed.FDot16, dx fixed.FDot16, mod64 int32) fixed.FDot16 {
	fx += 32768 // HALF
	if fx < 0 {
		fx = 0
	}
	xCoord := fx >> 16
	a := uint8(fx >> 8)

	ma := fixed.FDot6SmallScale(a, mod64)
	maa := fixed.FDot6SmallScale(255-a, mod64)

	// Match Rust's x.max(1) - 1 logic
	xVal := uint32(0)
	if int32(xCoord) > 0 {
		xVal = uint32(int32(xCoord) - 1)
	}

	if ma != 0 || maa != 0 {
		v.blitter.BlitAntiH2(xVal, y, maa, ma)
	}

	return fx + dx - 32768
}

func (v *vertishAntiHairBlitter) drawLine(y uint32, stopY uint32, fx fixed.FDot16, dx fixed.FDot16) fixed.FDot16 {
	count := stopY - y
	if count <= 0 {
		return fx
	}

	fx += 32768 // HALF
	loopY := y
	for {
		if fx < 0 {
			fx = 0
		}
		xCoord := fx >> 16
		a := uint8(fx >> 8)

		// Match Rust's x.max(1) - 1 logic
		xVal := uint32(0)
		if int32(xCoord) > 0 {
			xVal = uint32(int32(xCoord) - 1)
		}

		// Call BlitAntiH2 with the two alpha values
		v.blitter.BlitAntiH2(xVal, loopY, 255-a, a)

		fx += dx
		loopY++
		if loopY >= stopY {
			break
		}
	}

	return fx - 32768
}

func i32ToAlpha(a int32) uint8 {
	return uint8(a & 0xFF)
}

// rectClipBlitter clips drawing to a rectangle
type rectClipBlitter struct {
	blitter blitter.Blitter
	clip    path.ScreenIntRect
}

func (r *rectClipBlitter) BlitH(x, y, count uint32) {
	// TODO: implement clipping for BlitH
	r.blitter.BlitH(x, y, count)
}

func (r *rectClipBlitter) BlitMask(mask blitter.Mask, rect path.ScreenIntRect) {
	// TODO: implement masking with clipping
	r.blitter.BlitMask(mask, rect)
}

func (r *rectClipBlitter) BlitRect(rect path.ScreenIntRect) {
	// TODO: implement rect clipping
	r.blitter.BlitRect(rect)
}

func (r *rectClipBlitter) BlitAntiH(x, y uint32, aa []uint8, runs []uint16) {
	yInRect := func(y uint32, rect path.ScreenIntRect) bool {
		return (y - uint32(rect.Top())) < rect.Height()
	}

	if !yInRect(y, r.clip) || x >= uint32(r.clip.Right()) {
		return
	}

	var x0 = x
	var x1 = x + computeAntiWidth(runs)

	if x1 <= uint32(r.clip.Left()) {
		return
	}

	if x0 < uint32(r.clip.Left()) {
		dx := uint32(r.clip.Left()) - x0
		edge.BreakAt(aa, runs, int32(dx))
		aa = aa[dx:]
		runs = runs[dx:]
		x0 = uint32(r.clip.Left())
	}

	if x1 > uint32(r.clip.Right()) {
		x1 = uint32(r.clip.Right())
		edge.BreakAt(aa, runs, int32(x1-x0))
		runs[x1-x0] = 0
	}

	r.blitter.BlitAntiH(x0, y, aa, runs)
}

func (r *rectClipBlitter) BlitV(x, y uint32, height uint32, alpha uint8) {
	xInRect := func(x uint32, rect path.ScreenIntRect) bool {
		return (x - uint32(rect.Left())) < rect.Width()
	}

	if !xInRect(x, r.clip) {
		return
	}

	y0 := y
	y1 := y + height

	if y0 < uint32(r.clip.Top()) {
		y0 = uint32(r.clip.Top())
	}
	if y1 > uint32(r.clip.Bottom()) {
		y1 = uint32(r.clip.Bottom())
	}

	if y0 < y1 {
		r.blitter.BlitV(x, y0, y1-y0, alpha)
	}
}

func (r *rectClipBlitter) BlitAntiH2(x, y uint32, a0, a1 uint8) {
	r.BlitAntiH(x, y, []uint8{a0, a1}, []uint16{1, 1, 0})
}

func (r *rectClipBlitter) BlitAntiV2(x, y uint32, a0, a1 uint8) {
	r.BlitAntiH(x, y, []uint8{a0}, []uint16{1, 0})
	r.BlitAntiH(x, y+1, []uint8{a1}, []uint16{1, 0})
}

// computeAntiWidth computes the width of anti-aliased runs
func computeAntiWidth(runs []uint16) uint32 {
	var width uint32
	for i := 0; i < len(runs) && runs[i] != 0; {
		count := uint32(runs[i])
		width += count
		i += int(count)
	}
	return width
}
