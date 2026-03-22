// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package scalar

import (
	"github.com/chewxy/math32"
)

const (
	// SCALAR_MAX is the maximum value for a scalar.
	SCALAR_MAX float32 = 3.402823466e+38
	// SCALAR_NEARLY_ZERO is a small value used for comparison.
	SCALAR_NEARLY_ZERO float32 = 1.0 / (1 << 12)
	// SCALAR_ROOT_2_OVER_2 is the square root of 2 divided by 2.
	SCALAR_ROOT_2_OVER_2 float32 = 0.707106781
)

// Half returns half of the value.
func Half(s float32) float32 {
	return s * 0.5
}

// Ave returns the average of two values.
func Ave(a, b float32) float32 {
	return (a + b) * 0.5
}

// Sqr returns the square of the value.
func Sqr(s float32) float32 {
	return s * s
}

// Invert returns the inverse of the value.
func Invert(s float32) float32 {
	return 1.0 / s
}

// Works just like SkTPin, returning `max` for NaN/inf.
// A non-panicking clamp.
func Bound(s, min, max float32) float32 {
	return math32.Max(min, math32.Min(s, max))
}

// IsNearlyEqual checks if two values are nearly equal.
func IsNearlyEqual(a, b float32) bool {
	return math32.Abs(a-b) <= SCALAR_NEARLY_ZERO
}

// IsNearlyEqualWithinTolerance checks if two values are nearly equal within a given tolerance.
func IsNearlyEqualWithinTolerance(a, b, tolerance float32) bool {
	return math32.Abs(a-b) <= tolerance
}

// IsNearlyZero checks if the value is nearly zero.
func IsNearlyZero(s float32) bool {
	return IsNearlyZeroWithinTolerance(s, SCALAR_NEARLY_ZERO)
}

// IsNearlyZeroWithinTolerance checks if the value is nearly zero within a given tolerance.
func IsNearlyZeroWithinTolerance(s, tolerance float32) bool {
	return math32.Abs(s) <= tolerance
}

// AlmostEqualUlps compares two floats using Units in the Last Place.
// From SkPathOpsTypes.
func AlmostEqualUlps(a, b float32) bool {
	const ulpsEpsilon int32 = 16
	aBits := f32As2sCompliment(a)
	bBits := f32As2sCompliment(b)

	// Find the difference in ULPs.
	return aBits < bBits+ulpsEpsilon && bBits < aBits+ulpsEpsilon
}

// f32As2sCompliment is a helper for ULPS comparison.
// It converts the float bits to a sign-and-magnitude-like integer for comparison.
// This matches the Rust implementation exactly.
func f32As2sCompliment(s float32) int32 {
	bits := int32(math32.Float32bits(s))
	if bits < 0 {
		bits &= 0x7FFFFFFF
		bits = -bits
	}
	return bits
}
