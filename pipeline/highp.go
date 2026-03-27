// Copyright 2018 Google Inc.
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package pipeline

import (
	"github.com/lumifloat/tinyskia/path"
)

const HIGH_STAGE_WIDTH = 8

type HighPipeline struct {
	pixmapSrc *PixmapCtx
	pixmapDst *SubPixmapCtx
	aaMaskCtx *AAMaskCtx
	maskCtx   *MaskCtx
	ctx       *Context
	r         [HIGH_STAGE_WIDTH]float32
	g         [HIGH_STAGE_WIDTH]float32
	b         [HIGH_STAGE_WIDTH]float32
	a         [HIGH_STAGE_WIDTH]float32
	dr        [HIGH_STAGE_WIDTH]float32
	dg        [HIGH_STAGE_WIDTH]float32
	db        [HIGH_STAGE_WIDTH]float32
	da        [HIGH_STAGE_WIDTH]float32
	tail      int
	dx        int
	dy        int
}

func StartHighPipeline(
	stages []Stage,
	rect *path.ScreenIntRect,
	aaMaskCtx *AAMaskCtx,
	maskCtx *MaskCtx,
	ctx *Context,
	pixmapSrc *PixmapCtx,
	pixmapDst *SubPixmapCtx,
) {
	var p HighPipeline
	p.pixmapSrc = pixmapSrc
	p.pixmapDst = pixmapDst
	p.maskCtx = maskCtx
	p.aaMaskCtx = aaMaskCtx
	p.ctx = ctx

	for y := rect.Y(); y < rect.Bottom(); y++ {
		x := int(rect.X())
		end := int(rect.Right())

		for x+HIGH_STAGE_WIDTH <= end {
			p.dx = x
			p.dy = int(y)
			p.tail = HIGH_STAGE_WIDTH

			for _, stage := range stages {

				switch stage {
				case StageMoveSourceToDestination:
					// move_source_to_destination
					p.MoveSourceToDestination()

				case StageMoveDestinationToSource:
					// move_destination_to_source
					p.MoveDestinationToSource()

				case StageClamp0:
					// clamp_0
					p.Clamp0()

				case StageClampA:
					// clamp_a
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
					// load_dst_u8 - unreachable for highp
					p.LoadDestinationU8()

				case StageStoreU8:
					// store_u8 - unreachable for highp
					p.StoreU8()

				case StageGather:
					// gather
					p.Gather()

				case StageLoadMaskU8:
					// load_mask_u8 - unreachable for highp
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
					// destination_over
					p.DestinationOver()

				case StageSourceAtop:
					// source_atop
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
					// color burn
					p.ColorBurn()

				case StageColorDodge:
					// color_dodge
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
					// soft_light
					p.SoftLight()

				case StageHue:
					// hue
					p.Hue()

				case StageSaturation:
					// saturation
					p.Saturation()

				case StageColor:
					// color
					p.Color()

				case StageLuminosity:
					// luminosity
					p.Luminosity()

				case StageSourceOverRgba:
					// source_over_rgba
					p.SourceOverRgba()

				case StageTransform:
					// transform
					p.Transform()

				case StageReflect:
					// reflect
					p.Reflect()

				case StageRepeat:
					// repeat
					p.Repeat()

				case StageBilinear:
					// bilinear
					p.Bilinear()

				case StageBicubic:
					// bicubic
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
					// xy_to_unit_angle
					p.XYToUnitAngle()

				case StageXYToRadius:
					// xy_to_radius
					p.XYToRadius()

				case StageXYTo2PtConicalFocalOnCircle:
					// xy_to_2pt_conical_focal_on_circle
					p.XYTo2PtConicalFocalOnCircle()

				case StageXYTo2PtConicalWellBehaved:
					// xy_to_2pt_conical_well_behaved
					p.XYTo2PtConicalWellBehaved()

				case StageXYTo2PtConicalSmaller:
					// xy_to_2pt_conical_smaller
					p.XYTo2PtConicalSmaller()

				case StageXYTo2PtConicalGreater:
					// xy_to_2pt_conical_greater
					p.XYTo2PtConicalGreater()

				case StageXYTo2PtConicalStrip:
					// xy_to_2pt_conical_strip
					p.XYTo2PtConicalStrip()

				case StageMask2PtConicalNan:
					// mask_2pt_conical_nan
					p.Mask2PtConicalNan()

				case StageMask2PtConicalDegenerates:
					// mask_2pt_conical_degenerates
					p.Mask2PtConicalDegenerates()

				case StageApplyVectorMask:
					// apply_vector_mask
					p.ApplyVectorMask()

				case StageAlter2PtConicalCompensateFocal:
					// alter_2pt_conical_compensate_focal
					p.Alter2PtConicalCompensateFocal()

				case StageAlter2PtConicalUnswap:
					// alter_2pt_conical_unswap
					p.Alter2PtConicalUnswap()

				case StageNegateX:
					// negate_x
					p.NegateX()

				case StageApplyConcentricScaleBias:
					// apply_concentric_scale_bias
					p.ApplyConcentricScaleBias()

				case StageGammaExpand2:
					// gamma_expand_2
					p.GammaExpand2()

				case StageGammaExpandDestination2:
					// gamma_expand_dst_2
					p.GammaExpandDestination2()

				case StageGammaCompress2:
					// gamma_compress_2
					p.GammaCompress2()

				case StageGammaExpand22:
					// gamma_expand_22
					p.GammaExpand22()

				case StageGammaExpandDestination22:
					// gamma_expand_dst_22
					p.GammaExpandDestination22()

				case StageGammaCompress22:
					// gamma_compress_22
					p.GammaCompress22()

				case StageGammaExpandSrgb:
					// gamma_expand_srgb
					p.GammaExpandSrgb()

				case StageGammaExpandDestinationSrgb:
					// gamma_expand_dst_srgb
					p.GammaExpandDestinationSrgb()

				case StageGammaCompressSrgb:
					// gamma_compress_srgb
					p.GammaCompressSrgb()

				}
			}

			x += HIGH_STAGE_WIDTH
		}

		if x != end {
			p.dx = x
			p.dy = int(y)
			p.tail = end - x
			for _, stage := range stages {
				switch stage {
				case StageMoveSourceToDestination:
					// move_source_to_destination
					p.MoveSourceToDestination()

				case StageMoveDestinationToSource:
					// move_destination_to_source
					p.MoveDestinationToSource()

				case StageClamp0:
					// clamp_0
					p.Clamp0()

				case StageClampA:
					// clamp_a
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
					// load_dst_u8 - unreachable for highp
					p.LoadDestinationU8Tail()

				case StageStoreU8:
					// store_u8 - unreachable for highp
					p.StoreU8Tail()

				case StageGather:
					// gather
					p.Gather()

				case StageLoadMaskU8:
					// load_mask_u8 - unreachable for highp
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
					// destination_over
					p.DestinationOver()

				case StageSourceAtop:
					// source_atop
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
					// color burn
					p.ColorBurn()

				case StageColorDodge:
					// color_dodge
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
					// soft_light
					p.SoftLight()

				case StageHue:
					// hue
					p.Hue()

				case StageSaturation:
					// saturation
					p.Saturation()

				case StageColor:
					// color
					p.Color()

				case StageLuminosity:
					// luminosity
					p.Luminosity()

				case StageSourceOverRgba:
					// source_over_rgba
					p.SourceOverRgbaTail()

				case StageTransform:
					// transform
					p.Transform()

				case StageReflect:
					// reflect
					p.Reflect()

				case StageRepeat:
					// repeat
					p.Repeat()

				case StageBilinear:
					// bilinear
					p.Bilinear()

				case StageBicubic:
					// bicubic
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
					// xy_to_unit_angle
					p.XYToUnitAngle()

				case StageXYToRadius:
					// xy_to_radius
					p.XYToRadius()

				case StageXYTo2PtConicalFocalOnCircle:
					// xy_to_2pt_conical_focal_on_circle
					p.XYTo2PtConicalFocalOnCircle()

				case StageXYTo2PtConicalWellBehaved:
					// xy_to_2pt_conical_well_behaved
					p.XYTo2PtConicalWellBehaved()

				case StageXYTo2PtConicalSmaller:
					// xy_to_2pt_conical_smaller
					p.XYTo2PtConicalSmaller()

				case StageXYTo2PtConicalGreater:
					// xy_to_2pt_conical_greater
					p.XYTo2PtConicalGreater()

				case StageXYTo2PtConicalStrip:
					// xy_to_2pt_conical_strip
					p.XYTo2PtConicalStrip()

				case StageMask2PtConicalNan:
					// mask_2pt_conical_nan
					p.Mask2PtConicalNan()

				case StageMask2PtConicalDegenerates:
					// mask_2pt_conical_degenerates
					p.Mask2PtConicalDegenerates()

				case StageApplyVectorMask:
					// apply_vector_mask
					p.ApplyVectorMask()

				case StageAlter2PtConicalCompensateFocal:
					// alter_2pt_conical_compensate_focal
					p.Alter2PtConicalCompensateFocal()

				case StageAlter2PtConicalUnswap:
					// alter_2pt_conical_unswap
					p.Alter2PtConicalUnswap()

				case StageNegateX:
					// negate_x
					p.NegateX()

				case StageApplyConcentricScaleBias:
					// apply_concentric_scale_bias
					p.ApplyConcentricScaleBias()

				case StageGammaExpand2:
					// gamma_expand_2
					p.GammaExpand2()

				case StageGammaExpandDestination2:
					// gamma_expand_dst_2
					p.GammaExpandDestination2()

				case StageGammaCompress2:
					// gamma_compress_2
					p.GammaCompress2()

				case StageGammaExpand22:
					// gamma_expand_22
					p.GammaExpand22()

				case StageGammaExpandDestination22:
					// gamma_expand_dst_22
					p.GammaExpandDestination22()

				case StageGammaCompress22:
					// gamma_compress_22
					p.GammaCompress22()

				case StageGammaExpandSrgb:
					// gamma_expand_srgb
					p.GammaExpandSrgb()

				case StageGammaExpandDestinationSrgb:
					// gamma_expand_dst_srgb
					p.GammaExpandDestinationSrgb()

				case StageGammaCompressSrgb:
					// gamma_compress_srgb
					p.GammaCompressSrgb()

				}
			}
		}
	}
}
