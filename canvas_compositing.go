// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import "github.com/lumifloat/tinyskia/internal/core/painter"

// A Composite mode.
type CompositeOperation string

const (
	// Replaces destination with zero: fully transparent.
	CompositeOperationClear CompositeOperation = "clear"
	// Replaces destination.
	CompositeOperationCopy CompositeOperation = "copy"
	// Preserves destination.
	CompositeOperationDestination CompositeOperation = "destination"
	// Source over destination.
	CompositeOperationSourceOver CompositeOperation = "source-over"
	// Destination over source.
	CompositeOperationDestinationOver CompositeOperation = "destination-over"
	// Source trimmed inside destination.
	CompositeOperationSourceIn CompositeOperation = "source-in"
	// Destination trimmed by source.
	CompositeOperationDestinationIn CompositeOperation = "destination-in"
	// Source trimmed outside destination.
	CompositeOperationSourceOut CompositeOperation = "source-out"
	// Destination trimmed outside source.
	CompositeOperationDestinationOut CompositeOperation = "destination-out"
	// Source inside destination blended with destination.
	CompositeOperationSourceAtop CompositeOperation = "source-atop"
	// Destination inside source blended with source.
	CompositeOperationDestinationAtop CompositeOperation = "destination-atop"
	// Each of source and destination trimmed outside the other.
	CompositeOperationXor CompositeOperation = "xor"
	// Sum of colors.
	CompositeOperationLighter CompositeOperation = "lighter"
	// Multiply source with destination, darkening image.
	CompositeOperationMultiply CompositeOperation = "multiply"
	// Multiply inverse of pixels, inverting result; brightens destination.
	CompositeOperationScreen CompositeOperation = "screen"
	// Multiply or screen, depending on destination.
	CompositeOperationOverlay CompositeOperation = "overlay"
	// Darker of source and destination.
	CompositeOperationDarken CompositeOperation = "darken"
	// Lighter of source and destination.
	CompositeOperationLighten CompositeOperation = "lighten"
	// Brighten destination to reflect source.
	CompositeOperationColorDodge CompositeOperation = "color-dodge"
	// Darken destination to reflect source.
	CompositeOperationColorBurn CompositeOperation = "color-burn"
	// Multiply or screen, depending on source.
	CompositeOperationHardLight CompositeOperation = "hard-light"
	// Lighten or darken, depending on source.
	CompositeOperationSoftLight CompositeOperation = "soft-light"
	// Subtract darker from lighter with higher contrast.
	CompositeOperationDifference CompositeOperation = "difference"
	// Subtract darker from lighter with lower contrast.
	CompositeOperationExclusion CompositeOperation = "exclusion"
	// Hue of source with saturation and luminosity of destination.
	CompositeOperationHue CompositeOperation = "hue"
	// Saturation of source with hue and luminosity of destination.
	CompositeOperationSaturation CompositeOperation = "saturation"
	// Hue and saturation of source with luminosity of destination.
	CompositeOperationColor CompositeOperation = "color"
	// Luminosity of source with hue and saturation of destination.
	CompositeOperationLuminosity CompositeOperation = "luminosity"
)

// GetGlobalAlpha return this's global alpha.
func (ctx *Context) GetGlobalAlpha() float64 {
	return ctx.globalAlpha
}

// SetGlobalAlpha set this's global alpha to the given value.
func (ctx *Context) SetGlobalAlpha(alpha float64) {
	// TODO
	ctx.globalAlpha = alpha
}

// GetGlobalCompositeOperation return this's current compositing and blending operator.
func (ctx *Context) GetGlobalCompositeOperation() CompositeOperation {
	return ctx.globalCompositeOperation
}

// SetGlobalCompositeOperation set this's current compositing and blending operator to the given value.
func (ctx *Context) SetGlobalCompositeOperation(op CompositeOperation) {
	ctx.globalCompositeOperation = op
}

func composite(op CompositeOperation) painter.CompositeOperation {
	switch op {
	case CompositeOperationClear:
		return painter.CompositeOperationClear
	case CompositeOperationCopy:
		return painter.CompositeOperationCopy
	case CompositeOperationDestination:
		return painter.CompositeOperationDestination
	case CompositeOperationSourceOver:
		return painter.CompositeOperationSourceOver
	case CompositeOperationDestinationOver:
		return painter.CompositeOperationDestinationOver
	case CompositeOperationSourceIn:
		return painter.CompositeOperationSourceIn
	case CompositeOperationDestinationIn:
		return painter.CompositeOperationDestinationIn
	case CompositeOperationSourceOut:
		return painter.CompositeOperationSourceOut
	case CompositeOperationDestinationOut:
		return painter.CompositeOperationDestinationOut
	case CompositeOperationSourceAtop:
		return painter.CompositeOperationSourceAtop
	case CompositeOperationDestinationAtop:
		return painter.CompositeOperationDestinationAtop
	case CompositeOperationXor:
		return painter.CompositeOperationXor
	case CompositeOperationLighter:
		return painter.CompositeOperationLighter
	case CompositeOperationMultiply:
		return painter.CompositeOperationMultiply
	case CompositeOperationScreen:
		return painter.CompositeOperationScreen
	case CompositeOperationOverlay:
		return painter.CompositeOperationOverlay
	case CompositeOperationDarken:
		return painter.CompositeOperationDarken
	case CompositeOperationLighten:
		return painter.CompositeOperationLighten
	case CompositeOperationColorDodge:
		return painter.CompositeOperationColorDodge
	case CompositeOperationColorBurn:
		return painter.CompositeOperationColorBurn
	case CompositeOperationHardLight:
		return painter.CompositeOperationHardLight
	case CompositeOperationSoftLight:
		return painter.CompositeOperationSoftLight
	case CompositeOperationDifference:
		return painter.CompositeOperationDifference
	case CompositeOperationExclusion:
		return painter.CompositeOperationExclusion
	case CompositeOperationHue:
		return painter.CompositeOperationHue
	case CompositeOperationSaturation:
		return painter.CompositeOperationSaturation
	case CompositeOperationColor:
		return painter.CompositeOperationColor
	case CompositeOperationLuminosity:
		return painter.CompositeOperationLuminosity
	default:
		return painter.CompositeOperationSourceOver
	}
}
