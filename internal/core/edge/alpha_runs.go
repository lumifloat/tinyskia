// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package edge

const AlphaRunEnd uint16 = 0

// AlphaRuns is a sparse array of run-length-encoded alpha (supersampling coverage) values.
//
// Sparseness allows us to independently compose several paths into the
// same AlphaRuns buffer.
type AlphaRuns struct {
	Runs  []uint16 // Run lengths (0 indicates end)
	Alpha []uint8  // Alpha values for each run
}

// NewAlphaRuns creates a new AlphaRuns buffer for the given width.
func NewAlphaRuns(width uint32) *AlphaRuns {
	runs := &AlphaRuns{
		Runs:  make([]uint16, width+1),
		Alpha: make([]uint8, width+1),
	}
	runs.Reset(width)
	return runs
}

// CatchOverflow returns 0-255 given 0-256.
func CatchOverflow(alpha uint16) uint8 {
	return uint8(alpha - (alpha >> 8))
}

// IsEmpty returns true if the scanline contains only a single run, of alpha value 0.
func (r *AlphaRuns) IsEmpty() bool {
	if r.Runs[0] == AlphaRunEnd {
		return true
	}
	run := int(r.Runs[0])
	// 整行都是透明的，跳过渲染
	return r.Alpha[0] == 0 && r.Runs[run] == AlphaRunEnd
}

// Reset reinitializes for a new scanline.
func (r *AlphaRuns) Reset(width uint32) {
	r.Runs[0] = uint16(width)
	r.Runs[width] = AlphaRunEnd
	r.Alpha[0] = 0
}

// Add inserts into the buffer a run starting at (x-offsetX).
//
// if startAlpha > 0
//     one pixel with value += startAlpha,
//         max 255
// if middleCount > 0
//     middleCount pixels with value += maxValue
// if stopAlpha > 0
//     one pixel with value += stopAlpha
//
// Returns the offsetX value that should be passed on the next call,
// assuming we're on the same scanline. If the caller is switching
// scanlines, then offsetX should be 0 when this is called.
func (r *AlphaRuns) Add(x uint32, startAlpha uint8, middleCount int, stopAlpha uint8, maxValue uint8, offsetX int) int {
	xInt := int(x)
	xInt -= offsetX

	runsOffset := offsetX
	alphaOffset := offsetX
	lastAlphaOffset := offsetX

	if startAlpha != 0 {
		BreakRun(r.Runs[runsOffset:], r.Alpha[alphaOffset:], xInt, 1)

		tmp := uint16(r.Alpha[alphaOffset+xInt]) + uint16(startAlpha)
		// was (tmp >> 7), but that seems wrong if we're trying to catch 256
		r.Alpha[alphaOffset+xInt] = uint8(tmp - (tmp >> 8))

		runsOffset += xInt + 1
		alphaOffset += xInt + 1
		xInt = 0
	}

	if middleCount != 0 {
		BreakRun(r.Runs[runsOffset:], r.Alpha[alphaOffset:], xInt, middleCount)
		alphaOffset += xInt
		runsOffset += xInt
		xInt = 0
		for middleCount > 0 {
			a := CatchOverflow(uint16(r.Alpha[alphaOffset]) + uint16(maxValue))
			r.Alpha[alphaOffset] = a

			n := int(r.Runs[runsOffset])

			if n == 0 {
				// Safety check: prevent infinite loop if run length is 0
				break
			}

			alphaOffset += n
			runsOffset += n
			middleCount -= n
		}

		lastAlphaOffset = alphaOffset
	}

	if stopAlpha != 0 {
		BreakRun(r.Runs[runsOffset:], r.Alpha[alphaOffset:], xInt, 1)
		alphaOffset += xInt
		r.Alpha[alphaOffset] += stopAlpha
		lastAlphaOffset = alphaOffset
	}

	return lastAlphaOffset
}

// BreakRun breaks the runs in the buffer at offsets x and x+count, properly
// updating the runs to the right and left.
//
// i.e. from the state AAAABBBB, run-length encoded as A4B4,
// break_run(..., 2, 5) would produce AAAABBBB rle as A2A2B3B1.
// Allows add() to sum another run to some of the new sub-runs.
// i.e. adding ..CCCCC. would produce AADDEEEB, rle as A2D2E3B1.
func BreakRun(runs []uint16, alpha []uint8, x int, count int) {
	origX := x
	runsOffset := 0
	alphaOffset := 0

	for x > 0 {
		n := int(runs[runsOffset])

		if n == 0 {
			// Safety check: prevent infinite loop if run length is 0
			break
		}

		if x < n {
			alpha[alphaOffset+x] = alpha[alphaOffset]
			runs[runsOffset+0] = uint16(x)
			runs[runsOffset+x] = uint16(n - x)
			break
		}
		runsOffset += n
		alphaOffset += n
		x -= n
	}

	runsOffset = origX
	alphaOffset = origX
	x = count

	for {
		n := int(runs[runsOffset])

		if n == 0 {
			// Safety check: prevent infinite loop if run length is 0
			break
		}

		if x < n {
			alpha[alphaOffset+x] = alpha[alphaOffset]
			runs[runsOffset+0] = uint16(x)
			runs[runsOffset+x] = uint16(n - x)
			break
		}

		x -= n
		if x == 0 {
			break
		}

		runsOffset += n
		alphaOffset += n
	}
}

// BreakAt cuts (at offset x in the buffer) a run into two shorter runs with
// matching alpha values.
//
// Used by the RectClipBlitter to trim a RLE encoding to match the
// clipping rectangle.
func BreakAt(alpha []uint8, runs []uint16, x int32) {
	var alphaI int32 = 0
	var runI int32 = 0
	for x > 0 {
		n := int32(runs[runI])
		if x < n {
			alpha[int(alphaI+x)] = alpha[int(alphaI)]
			runs[0] = uint16(x)
			runs[int(x)] = uint16(n - x)
			break
		}

		runI += n
		alphaI += n
		x -= n
	}
}
