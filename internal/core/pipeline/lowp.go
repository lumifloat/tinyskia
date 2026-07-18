// Copyright 2016 Google Inc.
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package pipeline

import (
	"image"
)

const LOW_STAGE_WIDTH = 16

// LowPipeline 低精度渲染管线结构体
type LowPipeline struct {
	dst       *image.RGBA
	mask      *image.Alpha
	aaMaskCtx *AAMaskCtx
	ctx       *Context
	r         [LOW_STAGE_WIDTH]uint16
	g         [LOW_STAGE_WIDTH]uint16
	b         [LOW_STAGE_WIDTH]uint16
	a         [LOW_STAGE_WIDTH]uint16
	dr        [LOW_STAGE_WIDTH]uint16
	dg        [LOW_STAGE_WIDTH]uint16
	db        [LOW_STAGE_WIDTH]uint16
	da        [LOW_STAGE_WIDTH]uint16
	tail      int
	dx        int
	dy        int
	stop      bool
}

func StartLowPipeline(stages []Stage, rect image.Rectangle, aaMask *AAMaskCtx, mask *image.Alpha, ctx *Context, dst *image.RGBA) {
	var p LowPipeline
	p.dst = dst
	p.mask = mask
	p.aaMaskCtx = aaMask
	p.ctx = ctx

	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		x := rect.Min.X
		end := rect.Max.X

		for x+LOW_STAGE_WIDTH <= end {
			p.dx = x
			p.dy = y
			p.tail = LOW_STAGE_WIDTH
			p.stop = false

			for _, stage := range stages {
				if p.stop {
					break
				}

				switch stage {
				case StageMoveSourceToDestination:
					// move_source_to_destination
					p.MoveSourceToDestination()

				case StageMoveDestinationToSource:
					// move_destination_to_source
					p.MoveDestinationToSource()

				case StageClamp0:
					// null_fn
					p.Clamp0()

				case StageClampA:
					// null_fn
					p.ClampA()

				case StagePremultiply:
					// premultiply
					p.Premultiply()

				case StageUniformColor:
					// uniform_color
					p.UniformColor()

				case StageSeedShader:
					// seed_shader
					p.SeedShader()

				case StageLoadDestination:
					// load_dst
					p.LoadDestination()

				case StageStore:
					// store
					p.Store()

				case StageLoadDestinationU8:
					// load_dst_u8
					p.LoadDestinationU8()

				case StageStoreU8:
					// store_u8
					p.StoreU8()

				case StageGather:
					// gather
					p.Gather()

				case StageLoadMaskU8:
					// load_mask_u8
					p.LoadMaskU8()

				case StageMaskU8:
					// mask_u8
					p.MaskU8()

				case StageScaleU8:
					// scale_u8
					p.ScaleU8()

				case StageLerpU8:
					// lerp_u8
					p.LerpU8()

				case StageScale1Float:
					// scale_1_float
					p.Scale1Float()

				case StageLerp1Float:
					// lerp_1_float
					p.Lerp1Float()

				case StageDestinationAtop:
					// destination_atop
					p.DestinationAtop()

				case StageDestinationIn:
					// destination_in
					p.DestinationIn()

				case StageDestinationOut:
					// destination_out
					p.DestinationOut()

				case StageDestinationOver:
					// Formula: d + div255(s * inv(da))
					p.DestinationOver()

				case StageSourceAtop:
					// source_atop - sa * da + d * inv(sa)
					p.SourceAtop()

				case StageSourceIn:
					// source_in
					p.SourceIn()

				case StageSourceOut:
					// source_out
					p.SourceOut()

				case StageSourceOver:
					// source_over
					p.SourceOver()

				case StageClear:
					// clear
					p.Clear()

				case StageModulate:
					// modulate
					p.Modulate()

				case StageMultiply:
					// multiply
					p.Multiply()

				case StagePlus:
					// plus
					p.Plus()

				case StageScreen:
					// screen
					p.Screen()

				case StageXor:
					// xor
					p.Xor()

				case StageColorBurn:
					// null_fn
					p.ColorBurn()

				case StageColorDodge:
					// null_fn
					p.ColorDodge()

				case StageDarken:
					// darken
					p.Darken()

				case StageDifference:
					// difference
					p.Difference()

				case StageExclusion:
					// exclusion
					p.Exclusion()

				case StageHardLight:
					// hard_light
					p.HardLight()

				case StageLighten:
					// lighten
					p.Lighten()

				case StageOverlay:
					// overlay
					p.Overlay()

				case StageSoftLight:
					// null_fn
					p.SoftLight()

				case StageHue:
					// null_fn
					p.Hue()

				case StageSaturation:
					// null_fn
					p.Saturation()

				case StageColor:
					// null_fn
					p.Color()

				case StageLuminosity:
					// null_fn
					p.Luminosity()

				case StageSourceOverRgba:
					// source_over_rgba
					p.SourceOverRgba()

				case StageTransform:
					// transform
					p.Transform()

				case StageReflect:
					// null_fn
					p.Reflect()

				case StageRepeat:
					// null_fn
					p.Repeat()

				case StageDecal:
					// decal
					p.Decal()

				case StageCheckDecalMask:
					// check_decal_mask
					p.CheckDecalMask()

				case StageBilinear:
					// null_fn
					p.Bilinear()

				case StageBicubic:
					// null_fn
					p.Bicubic()

				case StagePadX1:
					// pad_x1
					p.PadX1()

				case StageReflectX1:
					// reflect_x1
					p.ReflectX1()

				case StageRepeatX1:
					// repeat_x1
					p.RepeatX1()

				case StageGradient:
					// gradient
					p.Gradient()

				case StageEvenlySpaced2StopGradient:
					// evenly_spaced_2_stop_gradient
					p.EvenlySpaced2StopGradient()

				case StageXYToUnitAngle:
					// null_fn
					p.XYToUnitAngle()

				case StageXYToRadius:
					// xy_to_radius
					p.XYToRadius()

				case StageXYTo2PtConicalFocalOnCircle, StageXYTo2PtConicalWellBehaved,
					StageXYTo2PtConicalSmaller, StageXYTo2PtConicalGreater,
					StageXYTo2PtConicalStrip, StageMask2PtConicalNan,
					StageMask2PtConicalDegenerates, StageApplyVectorMask,
					StageAlter2PtConicalCompensateFocal, StageAlter2PtConicalUnswap,
					StageNegateX, StageApplyConcentricScaleBias,
					StageGammaExpand2, StageGammaExpandDestination2,
					StageGammaCompress2, StageGammaExpand22,
					StageGammaExpandDestination22, StageGammaCompress22,
					StageGammaExpandSrgb, StageGammaExpandDestinationSrgb,
					StageGammaCompressSrgb:
					// null_fn
				}
			}

			x += LOW_STAGE_WIDTH
		}

		if x != end {
			p.dx = x
			p.dy = y
			p.tail = end - x
			p.stop = false

			for _, stage := range stages {
				if p.stop {
					break
				}

				switch stage {
				case StageMoveSourceToDestination:
					// move_source_to_destination
					p.MoveSourceToDestination()

				case StageMoveDestinationToSource:
					// move_destination_to_source
					p.MoveDestinationToSource()

				case StageClamp0:
					// null_fn
					p.Clamp0()

				case StageClampA:
					// null_fn
					p.ClampA()

				case StagePremultiply:
					// premultiply
					p.Premultiply()

				case StageUniformColor:
					// uniform_color
					p.UniformColor()

				case StageSeedShader:
					// seed_shader
					p.SeedShader()

				case StageLoadDestination:
					// load_dst
					p.LoadDestinationTail()

				case StageStore:
					// store
					p.StoreTail()

				case StageLoadDestinationU8:
					// load_dst_u8
					p.LoadDestinationU8()

				case StageStoreU8:
					// store_u8
					p.StoreU8Tail()

				case StageGather:
					// gather
					p.Gather()

				case StageLoadMaskU8:
					// load_mask_u8
					p.LoadMaskU8()

				case StageMaskU8:
					// mask_u8
					p.MaskU8()

				case StageScaleU8:
					// scale_u8
					p.ScaleU8()

				case StageLerpU8:
					// lerp_u8
					p.LerpU8()

				case StageScale1Float:
					// scale_1_float
					p.Scale1Float()

				case StageLerp1Float:
					// lerp_1_float
					p.Lerp1Float()

				case StageDestinationAtop:
					// destination_atop
					p.DestinationAtop()

				case StageDestinationIn:
					// destination_in
					p.DestinationIn()

				case StageDestinationOut:
					// destination_out
					p.DestinationOut()

				case StageDestinationOver:
					// Formula: d + div255(s * inv(da))
					p.DestinationOver()

				case StageSourceAtop:
					// source_atop - sa * da + d * inv(sa)
					p.SourceAtop()

				case StageSourceIn:
					// source_in
					p.SourceIn()

				case StageSourceOut:
					// source_out
					p.SourceOut()

				case StageSourceOver:
					// source_over
					p.SourceOver()

				case StageClear:
					// clear
					p.Clear()

				case StageModulate:
					// modulate
					p.Modulate()

				case StageMultiply:
					// multiply
					p.Multiply()

				case StagePlus:
					// plus
					p.Plus()

				case StageScreen:
					// screen
					p.Screen()

				case StageXor:
					// xor
					p.Xor()

				case StageColorBurn:
					// null_fn
					p.ColorBurn()

				case StageColorDodge:
					// null_fn
					p.ColorDodge()

				case StageDarken:
					// darken
					p.Darken()

				case StageDifference:
					// difference
					p.Difference()

				case StageExclusion:
					// exclusion
					p.Exclusion()

				case StageHardLight:
					// hard_light
					p.HardLight()

				case StageLighten:
					// lighten
					p.Lighten()

				case StageOverlay:
					// overlay
					p.Overlay()

				case StageSoftLight:
					// null_fn
					p.SoftLight()

				case StageHue:
					// null_fn
					p.Hue()

				case StageSaturation:
					// null_fn
					p.Saturation()

				case StageColor:
					// null_fn
					p.Color()

				case StageLuminosity:
					// null_fn
					p.Luminosity()

				case StageSourceOverRgba:
					// source_over_rgba
					p.SourceOverRgbaTail()

				case StageTransform:
					// transform
					p.Transform()

				case StageReflect:
					// null_fn
					p.Reflect()

				case StageRepeat:
					// null_fn
					p.Repeat()

				case StageDecal:
					// decal
					p.Decal()

				case StageCheckDecalMask:
					// check_decal_mask
					p.CheckDecalMask()

				case StageBilinear:
					// null_fn
					p.Bilinear()

				case StageBicubic:
					// null_fn
					p.Bicubic()

				case StagePadX1:
					// pad_x1
					p.PadX1()

				case StageReflectX1:
					// reflect_x1
					p.ReflectX1()

				case StageRepeatX1:
					// repeat_x1
					p.RepeatX1()

				case StageGradient:
					// gradient
					p.Gradient()

				case StageEvenlySpaced2StopGradient:
					// evenly_spaced_2_stop_gradient
					p.EvenlySpaced2StopGradient()

				case StageXYToUnitAngle:
					// null_fn
					p.XYToUnitAngle()

				case StageXYToRadius:
					// xy_to_radius
					p.XYToRadius()

				case StageXYTo2PtConicalFocalOnCircle, StageXYTo2PtConicalWellBehaved,
					StageXYTo2PtConicalSmaller, StageXYTo2PtConicalGreater,
					StageXYTo2PtConicalStrip, StageMask2PtConicalNan,
					StageMask2PtConicalDegenerates, StageApplyVectorMask,
					StageAlter2PtConicalCompensateFocal, StageAlter2PtConicalUnswap,
					StageNegateX, StageApplyConcentricScaleBias,
					StageGammaExpand2, StageGammaExpandDestination2,
					StageGammaCompress2, StageGammaExpand22,
					StageGammaExpandDestination22, StageGammaCompress22,
					StageGammaExpandSrgb, StageGammaExpandDestinationSrgb,
					StageGammaCompressSrgb:
					// null_fn
				}
			}
		}
	}
}
