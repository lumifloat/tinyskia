// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package colorf

import (
	"image/color"

	"github.com/chewxy/math32"
	"github.com/lumifloat/tinyskia/internal/numeric/normalized"
)

type RGBAF struct {
	R, G, B, A normalized.NormalizedF32
}

func (c RGBAF) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R*0xffff + 0.5)
	g = uint32(c.G*0xffff + 0.5)
	b = uint32(c.B*0xffff + 0.5)
	a = uint32(c.A*0xffff + 0.5)
	return
}

type NRGBAF struct {
	R, G, B, A normalized.NormalizedF32
}

func (c NRGBAF) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R*c.A*0xffff + 0.5)
	g = uint32(c.G*c.A*0xffff + 0.5)
	b = uint32(c.B*c.A*0xffff + 0.5)
	a = uint32(c.A*0xffff + 0.5)
	return
}

var RGBAFModel = color.ModelFunc(func(c color.Color) color.Color {
	switch c := c.(type) {
	case color.RGBA:
		return RGBAF{
			R: normalized.NewNormalizedF32WithClamped(float32(c.R) / 0xff),
			G: normalized.NewNormalizedF32WithClamped(float32(c.G) / 0xff),
			B: normalized.NewNormalizedF32WithClamped(float32(c.B) / 0xff),
			A: normalized.NewNormalizedF32WithClamped(float32(c.A) / 0xff),
		}
	case color.NRGBA:
		r := float32(c.R) / 0xff
		g := float32(c.G) / 0xff
		b := float32(c.B) / 0xff
		a := float32(c.A) / 0xff
		return RGBAF{
			R: normalized.NewNormalizedF32WithClamped(r * a),
			G: normalized.NewNormalizedF32WithClamped(g * a),
			B: normalized.NewNormalizedF32WithClamped(b * a),
			A: normalized.NewNormalizedF32WithClamped(a),
		}
	case color.RGBA64:
		return RGBAF{
			R: normalized.NewNormalizedF32WithClamped(float32(c.R) / 0xffff),
			G: normalized.NewNormalizedF32WithClamped(float32(c.G) / 0xffff),
			B: normalized.NewNormalizedF32WithClamped(float32(c.B) / 0xffff),
			A: normalized.NewNormalizedF32WithClamped(float32(c.A) / 0xffff),
		}
	case color.NRGBA64:
		r := float32(c.R) / 0xffff
		g := float32(c.G) / 0xffff
		b := float32(c.B) / 0xffff
		a := float32(c.A) / 0xffff
		return RGBAF{
			R: normalized.NewNormalizedF32WithClamped(r * a),
			G: normalized.NewNormalizedF32WithClamped(g * a),
			B: normalized.NewNormalizedF32WithClamped(b * a),
			A: normalized.NewNormalizedF32WithClamped(a),
		}
	case RGBAF:
		return c
	case NRGBAF:
		return RGBAF{
			R: normalized.NewNormalizedF32WithClamped(float32(c.R * c.A)),
			G: normalized.NewNormalizedF32WithClamped(float32(c.G * c.A)),
			B: normalized.NewNormalizedF32WithClamped(float32(c.B * c.A)),
			A: c.A,
		}
	default:
		r, g, b, a := c.RGBA()
		return RGBAF{
			R: normalized.NewNormalizedF32WithClamped(float32(r) / 0xffff),
			G: normalized.NewNormalizedF32WithClamped(float32(g) / 0xffff),
			B: normalized.NewNormalizedF32WithClamped(float32(b) / 0xffff),
			A: normalized.NewNormalizedF32WithClamped(float32(a) / 0xffff),
		}
	}
})

var NRGBAFModel = color.ModelFunc(func(c color.Color) color.Color {
	switch c := c.(type) {
	case color.RGBA:
		a := float32(c.A) / 0xff
		if a == 0 {
			return NRGBAF{R: 0, G: 0, B: 0, A: 0}
		}
		return NRGBAF{
			R: normalized.NewNormalizedF32WithClamped((float32(c.R) / 0xff) / a),
			G: normalized.NewNormalizedF32WithClamped((float32(c.G) / 0xff) / a),
			B: normalized.NewNormalizedF32WithClamped((float32(c.B) / 0xff) / a),
			A: normalized.NewNormalizedF32WithClamped(a),
		}
	case color.NRGBA:
		return NRGBAF{
			R: normalized.NewNormalizedF32WithClamped(float32(c.R) / 0xff),
			G: normalized.NewNormalizedF32WithClamped(float32(c.G) / 0xff),
			B: normalized.NewNormalizedF32WithClamped(float32(c.B) / 0xff),
			A: normalized.NewNormalizedF32WithClamped(float32(c.A) / 0xff),
		}
	case color.RGBA64:
		a := float32(c.A) / 0xffff
		if a == 0 {
			return NRGBAF{R: 0, G: 0, B: 0, A: 0}
		}
		return NRGBAF{
			R: normalized.NewNormalizedF32WithClamped((float32(c.R) / 0xffff) / a),
			G: normalized.NewNormalizedF32WithClamped((float32(c.G) / 0xffff) / a),
			B: normalized.NewNormalizedF32WithClamped((float32(c.B) / 0xffff) / a),
			A: normalized.NewNormalizedF32WithClamped(a),
		}
	case color.NRGBA64:
		return NRGBAF{
			R: normalized.NewNormalizedF32WithClamped(float32(c.R) / 0xffff),
			G: normalized.NewNormalizedF32WithClamped(float32(c.G) / 0xffff),
			B: normalized.NewNormalizedF32WithClamped(float32(c.B) / 0xffff),
			A: normalized.NewNormalizedF32WithClamped(float32(c.A) / 0xffff),
		}
	case NRGBAF:
		return c
	case RGBAF:
		if c.A == 0 {
			return NRGBAF{R: 0, G: 0, B: 0, A: 0}
		}
		return NRGBAF{
			R: normalized.NewNormalizedF32WithClamped(float32(c.R / c.A)),
			G: normalized.NewNormalizedF32WithClamped(float32(c.G / c.A)),
			B: normalized.NewNormalizedF32WithClamped(float32(c.B / c.A)),
			A: c.A,
		}
	default:
		r16, g16, b16, a16 := c.RGBA()
		a := float32(a16) / 0xffff
		if a == 0 {
			return NRGBAF{R: 0, G: 0, B: 0, A: 0}
		}
		return NRGBAF{
			R: normalized.NewNormalizedF32WithClamped((float32(r16) / 0xffff) / a),
			G: normalized.NewNormalizedF32WithClamped((float32(g16) / 0xffff) / a),
			B: normalized.NewNormalizedF32WithClamped((float32(b16) / 0xffff) / a),
			A: normalized.NewNormalizedF32WithClamped(a),
		}
	}
})

func IsOpaque(c color.Color) bool {
	_, _, _, a := c.RGBA()
	return a == 0xffff
}

func ApplyOpacity(c color.Color, opacity float32) color.Color {
	b := NRGBAFModel.Convert(c).(NRGBAF)
	o := normalized.NewNormalizedF32WithClamped(opacity)
	b.A = normalized.NewNormalizedF32WithClamped(float32(b.A * o))
	return b
}

// The colorspace used to interpret pixel values.
type ColorSpace int

const (
	ColorSpaceLinear ColorSpace = iota
	ColorSpaceGamma2
	ColorSpaceSimpleSRGB
	ColorSpaceFullSRGBGamma
)

func (cs ColorSpace) ExpandChannel(x normalized.NormalizedF32) normalized.NormalizedF32 {
	switch cs {
	case ColorSpaceLinear:
		return x
	case ColorSpaceGamma2:
		return normalized.NormalizedF32(x * x)
	case ColorSpaceSimpleSRGB:
		return normalized.NormalizedF32(math32.Pow(float32(x), 2.2))
	case ColorSpaceFullSRGBGamma:
		if x <= 0.04045 {
			return normalized.NormalizedF32(x / 12.92)
		}
		return normalized.NormalizedF32(math32.Pow(float32((x+0.055)/1.055), 2.4))
	default:
		return x
	}
}

func (cs ColorSpace) ExpandColor(c color.Color) color.Color {
	b := NRGBAFModel.Convert(c).(NRGBAF)
	b.R = cs.ExpandChannel(b.R)
	b.G = cs.ExpandChannel(b.G)
	b.B = cs.ExpandChannel(b.B)
	return c
}

func (cs ColorSpace) CompressChannel(x normalized.NormalizedF32) normalized.NormalizedF32 {
	switch cs {
	case ColorSpaceLinear:
		return x
	case ColorSpaceGamma2:
		return normalized.NormalizedF32(math32.Sqrt(float32(x)))
	case ColorSpaceSimpleSRGB:
		return normalized.NormalizedF32(math32.Pow(float32(x), 1.0/2.2))
	case ColorSpaceFullSRGBGamma:
		if x <= 0.0031308 {
			return normalized.NormalizedF32(x * 12.92)
		}
		return normalized.NormalizedF32(math32.Pow(float32(x), 1.0/2.4)*1.055 - 0.055)
	default:
		return x
	}
}
