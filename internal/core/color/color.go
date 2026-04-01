// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package color

import (
	"github.com/chewxy/math32"
)

const (
	// Represents fully transparent AlphaU8 value.
	AlphaU8Transparent uint8 = 0x00
	// Represents fully opaque AlphaU8 value.
	AlphaU8Opaque uint8 = 0xFF
	// Represents fully transparent Alpha value.
	AlphaTransparent float32 = 0.0
	// Represents fully opaque Alpha value.
	AlphaOpaque float32 = 1.0
)

// A 32-bit RGBA color value.
type ColorU8 struct {
	r, g, b, a uint8
}

// FromRGBA creates a new color.
func ColorU8FromRGBA(r, g, b, a uint8) ColorU8 {
	return ColorU8{r, g, b, a}
}

func (c ColorU8) Red() uint8   { return c.r }
func (c ColorU8) Green() uint8 { return c.g }
func (c ColorU8) Blue() uint8  { return c.b }
func (c ColorU8) Alpha() uint8 { return c.a }

// Check that color is opaque.
func (c ColorU8) IsOpaque() bool {
	return c.Alpha() == AlphaU8Opaque
}

// Converts into a premultiplied color.
func (c ColorU8) Premultiply() PremultipliedColorU8 {
	a := c.Alpha()
	if a != AlphaU8Opaque {
		return PremultipliedColorU8FromRGBAUnchecked(
			PremultiplyU8(c.Red(), a),
			PremultiplyU8(c.Green(), a),
			PremultiplyU8(c.Blue(), a),
			a,
		)
	}
	return PremultipliedColorU8FromRGBAUnchecked(c.Red(), c.Green(), c.Blue(), a)
}

// A 32-bit premultiplied RGBA color value.
type PremultipliedColorU8 struct {
	r, g, b, a uint8
}

var PremultipliedColorU8Transparent = PremultipliedColorU8FromRGBAUnchecked(0, 0, 0, 0)

// FromRGBA creates a new premultiplied color.
func PremultipliedColorU8FromRGBA(r, g, b, a uint8) (PremultipliedColorU8, bool) {
	if r <= a && g <= a && b <= a {
		return PremultipliedColorU8{r, g, b, a}, true
	}
	return PremultipliedColorU8{}, false
}

func PremultipliedColorU8FromRGBAUnchecked(r, g, b, a uint8) PremultipliedColorU8 {
	return PremultipliedColorU8{r, g, b, a}
}

func (c PremultipliedColorU8) Red() uint8   { return c.r }
func (c PremultipliedColorU8) Green() uint8 { return c.g }
func (c PremultipliedColorU8) Blue() uint8  { return c.b }
func (c PremultipliedColorU8) Alpha() uint8 { return c.a }

func (c PremultipliedColorU8) IsOpaque() bool {
	return c.Alpha() == AlphaU8Opaque
}

// Returns a demultiplied color.
func (c PremultipliedColorU8) Demultiply() ColorU8 {
	alpha := c.Alpha()
	if alpha == AlphaU8Opaque {
		return ColorU8(c)
	}
	if alpha == 0 {
		return ColorU8FromRGBA(0, 0, 0, 0)
	}
	a := float64(alpha) / 255.0
	return ColorU8FromRGBA(
		uint8(float64(c.Red())/a+0.5),
		uint8(float64(c.Green())/a+0.5),
		uint8(float64(c.Blue())/a+0.5),
		alpha,
	)
}

// An RGBA color value, holding four floating point components.
type Color struct {
	r, g, b, a float32
}

var (
	ColorTransparent = Color{0, 0, 0, 0}
	ColorBlack       = Color{0, 0, 0, 1}
	ColorWhite       = Color{1, 1, 1, 1}
)

// FromRGBA creates a new color from 4 components.
func ColorFromRGBA(r, g, b, a float32) (Color, bool) {
	if r < 0 || r > 1 || g < 0 || g > 1 || b < 0 || b > 1 || a < 0 || a > 1 {
		return Color{}, false
	}
	return Color{r, g, b, a}, true
}

func ColorFromRGBA8(r, g, b, a uint8) Color {
	return Color{
		float32(r) / 255.0,
		float32(g) / 255.0,
		float32(b) / 255.0,
		float32(a) / 255.0,
	}
}

func (c Color) Red() float32   { return c.r }
func (c Color) Green() float32 { return c.g }
func (c Color) Blue() float32  { return c.b }
func (c Color) Alpha() float32 { return c.a }

func (c *Color) SetRed(val float32)   { c.r = clamp01(val) }
func (c *Color) SetGreen(val float32) { c.g = clamp01(val) }
func (c *Color) SetBlue(val float32)  { c.b = clamp01(val) }
func (c *Color) SetAlpha(val float32) { c.a = clamp01(val) }

func (c *Color) ApplyOpacity(opacity float32) {
	c.a = clamp01(c.a * clamp01(opacity))
}

func (c Color) IsOpaque() bool {
	return c.a == AlphaOpaque
}

func (c Color) Premultiply() PremultipliedColor {
	if c.IsOpaque() {
		return PremultipliedColor{c.r, c.g, c.b, c.a}
	}
	return PremultipliedColor{
		r: c.r * c.a,
		g: c.g * c.a,
		b: c.b * c.a,
		a: c.a,
	}
}

func (c Color) ToColorU8() ColorU8 {
	rgba := colorFloatToU8(c.r, c.g, c.b, c.a)
	return ColorU8FromRGBA(rgba[0], rgba[1], rgba[2], rgba[3])
}

// A premultiplied RGBA color value, holding four floating point components.
type PremultipliedColor struct {
	r, g, b, a float32
}

func (c PremultipliedColor) Red() float32   { return c.r }
func (c PremultipliedColor) Green() float32 { return c.g }
func (c PremultipliedColor) Blue() float32  { return c.b }
func (c PremultipliedColor) Alpha() float32 { return c.a }

func (c PremultipliedColor) Demultiply() Color {
	if c.a == 0 {
		return ColorTransparent
	}
	return Color{
		r: clamp01(c.r / c.a),
		g: clamp01(c.g / c.a),
		b: clamp01(c.b / c.a),
		a: c.a,
	}
}

func (c PremultipliedColor) ToColorU8() PremultipliedColorU8 {
	rgba := colorFloatToU8(c.r, c.g, c.b, c.a)
	return PremultipliedColorU8FromRGBAUnchecked(rgba[0], rgba[1], rgba[2], rgba[3])
}

// Return a*b/255, rounding any fractional bits.
func PremultiplyU8(c, a uint8) uint8 {
	prod := uint64(c)*uint64(a) + 128
	return uint8((prod + (prod >> 8)) >> 8)
}

func colorFloatToU8(r, g, b, a float32) [4]uint8 {
	return [4]uint8{
		uint8(r*255.0 + 0.5),
		uint8(g*255.0 + 0.5),
		uint8(b*255.0 + 0.5),
		uint8(a*255.0 + 0.5),
	}
}

func clamp01(v float32) float32 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

// The colorspace used to interpret pixel values.
type ColorSpace int

const (
	ColorSpaceLinear ColorSpace = iota
	ColorSpaceGamma2
	ColorSpaceSimpleSRGB
	ColorSpaceFullSRGBGamma
)

func (cs ColorSpace) ExpandChannel(x float32) float32 {
	switch cs {
	case ColorSpaceLinear:
		return x
	case ColorSpaceGamma2:
		return x * x
	case ColorSpaceSimpleSRGB:
		return clamp01(math32.Pow(x, 2.2))
	case ColorSpaceFullSRGBGamma:
		if x <= 0.04045 {
			return x / 12.92
		}
		return clamp01(math32.Pow((x+0.055)/1.055, 2.4))
	default:
		return x
	}
}

func (cs ColorSpace) ExpandColor(c Color) Color {
	c.r = cs.ExpandChannel(c.r)
	c.g = cs.ExpandChannel(c.g)
	c.b = cs.ExpandChannel(c.b)
	return c
}

func (cs ColorSpace) CompressChannel(x float32) float32 {
	switch cs {
	case ColorSpaceLinear:
		return x
	case ColorSpaceGamma2:
		return clamp01(math32.Sqrt(x))
	case ColorSpaceSimpleSRGB:
		return clamp01(math32.Pow(x, 1.0/2.2))
	case ColorSpaceFullSRGBGamma:
		if x <= 0.0031308 {
			return x * 12.92
		}
		return clamp01(math32.Pow(x, 1.0/2.4)*1.055 - 0.055)
	default:
		return x
	}
}
