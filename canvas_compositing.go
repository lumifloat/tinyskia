// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

// A CompositeOperation mode.
type CompositeOperation int

const (
	// Replaces destination with zero: fully transparent.
	CompositeOperationClear CompositeOperation = iota
	// Replaces destination.
	CompositeOperationSource
	// Preserves destination.
	CompositeOperationDestination
	// Source over destination.
	CompositeOperationSourceOver
	// Destination over source.
	CompositeOperationDestinationOver
	// Source trimmed inside destination.
	CompositeOperationSourceIn
	// Destination trimmed by source.
	CompositeOperationDestinationIn
	// Source trimmed outside destination.
	CompositeOperationSourceOut
	// Destination trimmed outside source.
	CompositeOperationDestinationOut
	// Source inside destination blended with destination.
	CompositeOperationSourceAtop
	// Destination inside source blended with source.
	CompositeOperationDestinationAtop
	// Each of source and destination trimmed outside the other.
	CompositeOperationXor
	// Sum of colors.
	CompositeOperationPlus
	// Product of premultiplied colors; darkens destination.
	CompositeOperationModulate
	// Multiply inverse of pixels, inverting result; brightens destination.
	CompositeOperationScreen
	// Multiply or screen, depending on destination.
	CompositeOperationOverlay
	// Darker of source and destination.
	CompositeOperationDarken
	// Lighter of source and destination.
	CompositeOperationLighten
	// Brighten destination to reflect source.
	CompositeOperationColorDodge
	// Darken destination to reflect source.
	CompositeOperationColorBurn
	// Multiply or screen, depending on source.
	CompositeOperationHardLight
	// Lighten or darken, depending on source.
	CompositeOperationSoftLight
	// Subtract darker from lighter with higher contrast.
	CompositeOperationDifference
	// Subtract darker from lighter with lower contrast.
	CompositeOperationExclusion
	// Multiply source with destination, darkening image.
	CompositeOperationMultiply
	// Hue of source with saturation and luminosity of destination.
	CompositeOperationHue
	// Saturation of source with hue and luminosity of destination.
	CompositeOperationSaturation
	// Hue and saturation of source with luminosity of destination.
	CompositeOperationColor
	// Luminosity of source with hue and saturation of destination.
	CompositeOperationLuminosity
)

// SetGlobalCompositeOperation change the current compositing and blending operator.
func (dc *Context) SetGlobalCompositeOperation(op CompositeOperation) {
	dc.composite = op
}

// GetGlobalCompositeOperation returns the current compositing and blending operator, from the values defined in Compositing and Blending.
func (dc *Context) GetGlobalCompositeOperation() CompositeOperation {
	return dc.composite
}
