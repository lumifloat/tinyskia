// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"errors"
	"image/color"
)

var ErrColorStopOffsetOutOfRange = errors.New("color stop offset out of range [0, 1]")

type Gradient interface {
	Style
	AddColorStop(offset float64, color color.Color) error
}

// AddColorStop adds a color stop with the given color to the gradient at the given offset.
// 0.0 is the offset at one end of the gradient, 1.0 is the offset at the other end.
func (g *LinearGradient) AddColorStop(offset float64, color color.Color) error {
	if offset < 0.0 || offset > 1.0 {
		return ErrColorStopOffsetOutOfRange
	}
	g.stops = append(g.stops, stop{pos: offset, color: color})
	return nil
}

// AddColorStop adds a color stop with the given color to the gradient at the given offset.
// 0.0 is the offset at one end of the gradient, 1.0 is the offset at the other end.
func (g *RadialGradient) AddColorStop(offset float64, color color.Color) error {
	if offset < 0.0 || offset > 1.0 {
		return ErrColorStopOffsetOutOfRange
	}
	g.stops = append(g.stops, stop{pos: offset, color: color})
	return nil
}

// AddColorStop adds a color stop with the given color to the gradient at the given offset.
// 0.0 is the offset at one end of the gradient, 1.0 is the offset at the other end.
func (g *ConicGradient) AddColorStop(offset float64, color color.Color) error {
	if offset < 0.0 || offset > 1.0 {
		return ErrColorStopOffsetOutOfRange
	}
	g.stops = append(g.stops, stop{pos: offset, color: color})
	return nil
}
