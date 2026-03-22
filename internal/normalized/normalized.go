// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package normalized

import "github.com/chewxy/math32"

const (
	FLOAT32_EPSILON float32 = 1.19209290e-07
)

type NormalizedF32Exclusive float32

func NewNormalizedF32Exclusive(v float32) (NormalizedF32Exclusive, bool) {
	if v > 0.0 && v < 1.0 {
		return NormalizedF32Exclusive(v), true
	}
	return 0.0, false
}

func NewNormalizedF32ExclusiveWithBounded(v float32) NormalizedF32Exclusive {
	const epsilon float32 = 1.1920929e-7

	lower := epsilon
	upper := 1.0 - epsilon

	if v < lower {
		v = lower
	} else if v > upper {
		v = upper
	}
	return NormalizedF32Exclusive(v)
}

func (v NormalizedF32Exclusive) Get() float32 {
	return float32(v)
}

func (v NormalizedF32Exclusive) ToNormalized() NormalizedF32 {
	return NormalizedF32(v)
}

type NormalizedF32 float32

const (
	NormalizedF32Zero = NormalizedF32(0.0)
	NormalizedF32One  = NormalizedF32(1.0)
	NormalizedF32Half = NormalizedF32(0.5)
)

func NewNormalizedF32(v float32) (NormalizedF32, bool) {
	if v < 0.0 {
		return 0.0, false
	}
	if v > 1.0 {
		return 0.0, false
	}
	return NormalizedF32(v), true
}

func NewNormalizedF32WithClamped(v float32) NormalizedF32 {
	if v <= 0.0 {
		return 0.0
	}
	if v >= 1.0 {
		return 1.0
	}
	return NormalizedF32(v)
}

func (v NormalizedF32) Get() float32 {
	return float32(v)
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
