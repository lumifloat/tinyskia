// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Skia uses fixed points pretty chaotically, therefore we cannot use
// strongly typed wrappers. Which is unfortunate.
package fixed

import (
	"math"
)

// FDot6 is a 26.6 fixed point.
type FDot6 int32

// FDot8 is a 24.8 fixed point.
type FDot8 int32

// FDot16 is a 16.16 fixed point.
type FDot16 int32

// FDot6 constants and functions
const (
	FDot6One = 64
)

// NewFDot6FromI32 converts int32 to FDot6.
func NewFDot6FromI32(n int32) FDot6 {
	return FDot6(n << 6)
}

// NewFDot6FromF32 converts float32 to FDot6.
func NewFDot6FromF32(n float32) FDot6 {
	return FDot6(n * 64.0)
}

// Floor returns the floor of FDot6.
func (n FDot6) Floor() FDot6 {
	return FDot6(n >> 6)
}

// Ceil returns the ceiling of FDot6.
func (n FDot6) Ceil() FDot6 {
	return FDot6((n + 63) >> 6)
}

// Round returns the rounded value of FDot6.
func (n FDot6) Round() FDot6 {
	return FDot6((n + 32) >> 6)
}

// ToFDot16 converts FDot6 to FDot16.
func (n FDot6) ToFDot16() FDot16 {
	return FDot16(int32(uint32(n) << 10))
}

// DivToFDot16 divides two FDot6 values and returns FDot16.
func (n FDot6) DivToFDot16(o FDot6) FDot16 {
	if FDot6(int32(int16(n))) == n {
		return FDot16((int32(uint32(n) << 16)) / int32(o))
	}
	v := (int64(uint64(n) << 16)) / int64(o)
	if v < int64(math.MinInt32) {
		return FDot16(math.MinInt32)
	}
	if v > int64(math.MaxInt32) {
		return FDot16(math.MaxInt32)
	}
	return FDot16(v)
}

// CanConvertToFDot16 checks if FDot6 can be converted to FDot16 without overflow.
func (n FDot6) CanConvertToFDot16() bool {
	maxDot6 := int32(math.MaxInt32 >> (16 - 6))
	absN := n
	if absN < 0 {
		absN = -absN
	}
	return int32(absN) <= maxDot6
}

// SmallScale performs small scale multiplication.
func (n FDot6) SmallScale(value uint8) uint8 {
	return uint8((int32(value)*int32(n))>>6) & 0xFF
}

// NewFDot8FromFDot16 converts FDot16 to FDot8.
func NewFDot8FromFDot16(x FDot16) FDot8 {
	return FDot8((x + 0x80) >> 8)
}

// FDot16 constants and functions
const (
	FDot16Half = (1 << 16) / 2
	FDot16One  = 1 << 16
)

// NewFDot16FromF32 converts float32 to FDot16.
func NewFDot16FromF32(x float32) FDot16 {
	if x < float32(math.MinInt32) {
		return FDot16(math.MinInt32)
	}
	if x > float32(math.MaxInt32) {
		return FDot16(math.MaxInt32)
	}
	return FDot16(int32(x * float32(FDot16One)))
}

// FloorToI32 returns the floor of FDot16 as int32.
func (n FDot16) FloorToI32() int32 {
	return int32(n >> 16)
}

// CeilToI32 returns the ceiling of FDot16 as int32.
func (n FDot16) CeilToI32() int32 {
	return int32(n+FDot16One-1) >> 16
}

// RoundToI32 returns the rounded value of FDot16 as int32.
func (n FDot16) RoundToI32() int32 {
	return int32(n+FDot16Half) >> 16
}

// Mul multiplies two FDot16 values.
func (n FDot16) Mul(o FDot16) FDot16 {
	return FDot16((int64(n) * int64(o)) >> 16)
}
