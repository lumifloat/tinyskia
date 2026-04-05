// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import "github.com/lumifloat/tinyskia/internal/core/pipeline"

// A blending mode.
type BlendMode int

const (
	// Replaces destination with zero: fully transparent.
	BlendModeClear BlendMode = iota
	// Replaces destination.
	BlendModeSource
	// Preserves destination.
	BlendModeDestination
	// Source over destination.
	BlendModeSourceOver
	// Destination over source.
	BlendModeDestinationOver
	// Source trimmed inside destination.
	BlendModeSourceIn
	// Destination trimmed by source.
	BlendModeDestinationIn
	// Source trimmed outside destination.
	BlendModeSourceOut
	// Destination trimmed outside source.
	BlendModeDestinationOut
	// Source inside destination blended with destination.
	BlendModeSourceAtop
	// Destination inside source blended with source.
	BlendModeDestinationAtop
	// Each of source and destination trimmed outside the other.
	BlendModeXor
	// Sum of colors.
	BlendModePlus
	// Product of premultiplied colors; darkens destination.
	BlendModeModulate
	// Multiply inverse of pixels, inverting result; brightens destination.
	BlendModeScreen
	// Multiply or screen, depending on destination.
	BlendModeOverlay
	// Darker of source and destination.
	BlendModeDarken
	// Lighter of source and destination.
	BlendModeLighten
	// Brighten destination to reflect source.
	BlendModeColorDodge
	// Darken destination to reflect source.
	BlendModeColorBurn
	// Multiply or screen, depending on source.
	BlendModeHardLight
	// Lighten or darken, depending on source.
	BlendModeSoftLight
	// Subtract darker from lighter with higher contrast.
	BlendModeDifference
	// Subtract darker from lighter with lower contrast.
	BlendModeExclusion
	// Multiply source with destination, darkening image.
	BlendModeMultiply
	// Hue of source with saturation and luminosity of destination.
	BlendModeHue
	// Saturation of source with hue and luminosity of destination.
	BlendModeSaturation
	// Hue and saturation of source with luminosity of destination.
	BlendModeColor
	// Luminosity of source with hue and saturation of destination.
	BlendModeLuminosity
)

func (dc *Context) SetGlobalCompositeOperation(blendMode BlendMode) {
	dc.blendMode = blendMode
}

func (b BlendMode) ShouldPreScaleCoverage() bool {
	// The most important things we do here are:
	//   1) never pre-scale with rgb coverage if the blend mode involves a source-alpha term;
	//   2) always pre-scale Plus.
	//
	// When we pre-scale with rgb coverage, we scale each of source r,g,b, with a distinct value,
	// and source alpha with one of those three values. This process destructively updates the
	// source-alpha term, so we can't evaluate blend modes that need its original value.
	//
	// Plus always requires pre-scaling as a specific quirk of its implementation in
	// RasterPipeline. This lets us put the clamp inside the blend mode itself rather
	// than as a separate stage that'd come after the lerp.
	//
	// This function is a finer-grained breakdown of SkBlendMode_SupportsCoverageAsAlpha().
	switch b {
	case BlendModeDestination,
		BlendModeDestinationOver,
		BlendModePlus,
		BlendModeDestinationOut,
		BlendModeSourceAtop,
		BlendModeSourceOver,
		BlendModeXor:
		return true
	default:
		return false
	}
}

func (b BlendMode) ToStage() (pipeline.Stage, bool) {
	switch b {
	case BlendModeClear:
		return pipeline.StageClear, true
	case BlendModeSource:
		return 0, false // This stage is a no-op.
	case BlendModeDestination:
		return pipeline.StageMoveDestinationToSource, true
	case BlendModeSourceOver:
		return pipeline.StageSourceOver, true
	case BlendModeDestinationOver:
		return pipeline.StageDestinationOver, true
	case BlendModeSourceIn:
		return pipeline.StageSourceIn, true
	case BlendModeDestinationIn:
		return pipeline.StageDestinationIn, true
	case BlendModeSourceOut:
		return pipeline.StageSourceOut, true
	case BlendModeDestinationOut:
		return pipeline.StageDestinationOut, true
	case BlendModeSourceAtop:
		return pipeline.StageSourceAtop, true
	case BlendModeDestinationAtop:
		return pipeline.StageDestinationAtop, true
	case BlendModeXor:
		return pipeline.StageXor, true
	case BlendModePlus:
		return pipeline.StagePlus, true
	case BlendModeModulate:
		return pipeline.StageModulate, true
	case BlendModeScreen:
		return pipeline.StageScreen, true
	case BlendModeOverlay:
		return pipeline.StageOverlay, true
	case BlendModeDarken:
		return pipeline.StageDarken, true
	case BlendModeLighten:
		return pipeline.StageLighten, true
	case BlendModeColorDodge:
		return pipeline.StageColorDodge, true
	case BlendModeColorBurn:
		return pipeline.StageColorBurn, true
	case BlendModeHardLight:
		return pipeline.StageHardLight, true
	case BlendModeSoftLight:
		return pipeline.StageSoftLight, true
	case BlendModeDifference:
		return pipeline.StageDifference, true
	case BlendModeExclusion:
		return pipeline.StageExclusion, true
	case BlendModeMultiply:
		return pipeline.StageMultiply, true
	case BlendModeHue:
		return pipeline.StageHue, true
	case BlendModeSaturation:
		return pipeline.StageSaturation, true
	case BlendModeColor:
		return pipeline.StageColor, true
	case BlendModeLuminosity:
		return pipeline.StageLuminosity, true
	default:
		return 0, false
	}
}
