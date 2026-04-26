// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package tinyskia

import (
	"github.com/lumifloat/tinyskia/internal/core/color"
	"github.com/lumifloat/tinyskia/internal/core/pipeline"
	"github.com/lumifloat/tinyskia/internal/core/shader"
	"github.com/lumifloat/tinyskia/internal/path"
)

// Paint controls how a shape should be painted.
type Paint struct {
	// A paint shader.
	Shader shader.Shader
	// Paint blending mode.
	BlendMode CompositeOperation
	// Enables anti-aliased painting.
	AntiAlias bool
	// Colorspace for blending.
	Colorspace color.ColorSpace
	// Forces the high quality/precision rendering pipeline.
	ForceHQPipeline bool
}

// DefaultPaint returns a Paint with default values.
func DefaultPaint() Paint {
	return Paint{
		Shader:          shader.NewSolidColor(color.ColorBlack),
		BlendMode:       CompositeOperationSourceOver,
		AntiAlias:       true,
		Colorspace:      color.ColorSpaceLinear,
		ForceHQPipeline: false,
	}
}

func (p *Paint) blitter(data, mask []uint8, width, height int) *pipeline.RasterPipelineBlitter {
	size, ok := path.NewIntSize(uint32(width), uint32(height))
	if !ok {
		return nil
	}

	var maskCtx *pipeline.MaskCtx
	if mask != nil {
		maskCtx = &pipeline.MaskCtx{
			Data:      mask,
			RealWidth: uint32(width),
		}
	}

	subPixmapCtx := &pipeline.SubPixmapCtx{
		Data:      data,
		Size:      size,
		RealWidth: width,
	}

	switch p.BlendMode {
	case CompositeOperationDestination:
		return nil
	case CompositeOperationDestinationIn:
		if solid, ok := p.Shader.(*shader.SolidColor); ok && solid.IsOpaque() {
			return nil
		}
	}

	// We can strength-reduce SourceOver into Source when opaque.
	blendMode := p.BlendMode
	if p.Shader.IsOpaque() && blendMode == CompositeOperationSourceOver && maskCtx == nil {
		blendMode = CompositeOperationSource
	}

	// When we're drawing a constant color in Source mode, we can sometimes just memset.
	var memset2dColor color.PremultipliedColorU8
	var useMemset2dColor bool
	if blendMode == CompositeOperationSource && maskCtx == nil {
		if solid, ok := p.Shader.(*shader.SolidColor); ok {
			memset2dColor = solid.Color().Premultiply().ToColorU8()
			useMemset2dColor = true
		}
	}

	// Clear is just a transparent color memset.
	if blendMode == CompositeOperationClear && !p.AntiAlias && maskCtx == nil {
		blendMode = CompositeOperationSource
		memset2dColor = color.PremultipliedColorU8Transparent
		useMemset2dColor = true
	}

	// blit_anti_h_rp
	blitAntiHRpBuilder := pipeline.NewRasterPipelineBuilder()
	blitAntiHRpBuilder.SetForceHqPipeline(p.ForceHQPipeline)
	if !p.Shader.PushStages(p.Colorspace, blitAntiHRpBuilder) {
		return nil
	}

	if maskCtx != nil {
		blitAntiHRpBuilder.Push(pipeline.StageMaskU8)
	}

	if blendMode.ShouldPreScaleCoverage() {
		blitAntiHRpBuilder.Push(pipeline.StageScale1Float)
		blitAntiHRpBuilder.Push(pipeline.StageLoadDestination)
		if stage, ok := expand(p.Colorspace); ok {
			blitAntiHRpBuilder.Push(stage)
		}
		if blendStage, ok := blendMode.stage(); ok {
			blitAntiHRpBuilder.Push(blendStage)
		}
	} else {
		blitAntiHRpBuilder.Push(pipeline.StageLoadDestination)
		if stage, ok := expand(p.Colorspace); ok {
			blitAntiHRpBuilder.Push(stage)
		}
		if blendStage, ok := blendMode.stage(); ok {
			blitAntiHRpBuilder.Push(blendStage)
		}
		blitAntiHRpBuilder.Push(pipeline.StageLerp1Float)
	}
	if stage, ok := compress(p.Colorspace); ok {
		blitAntiHRpBuilder.Push(stage)
	}
	blitAntiHRpBuilder.Push(pipeline.StageStore)
	blitAntiHRp := blitAntiHRpBuilder.Compile()

	// blit_rect_rp
	blitRectRpBuilder := pipeline.NewRasterPipelineBuilder()
	blitRectRpBuilder.SetForceHqPipeline(p.ForceHQPipeline)
	if !p.Shader.PushStages(p.Colorspace, blitRectRpBuilder) {
		return nil
	}

	if maskCtx != nil {
		blitRectRpBuilder.Push(pipeline.StageMaskU8)
	}

	if blendMode == CompositeOperationSourceOver && maskCtx == nil {
		if stage, ok := compress(p.Colorspace); ok {
			blitRectRpBuilder.Push(stage)
		}
		// TODO: ignore when dither_rate is non-zero
		blitRectRpBuilder.Push(pipeline.StageSourceOverRgba)
	} else {
		if blendMode != CompositeOperationSource {
			blitRectRpBuilder.Push(pipeline.StageLoadDestination)
			if blendStage, ok := blendMode.stage(); ok {
				if stage, ok := expand(p.Colorspace); ok {
					blitRectRpBuilder.Push(stage)
				}
				blitRectRpBuilder.Push(blendStage)
			}
		}
		if stage, ok := compress(p.Colorspace); ok {
			blitRectRpBuilder.Push(stage)
		}
		blitRectRpBuilder.Push(pipeline.StageStore)
	}
	blitRectRp := blitRectRpBuilder.Compile()

	// blit_mask_rp
	blitMaskRpBuilder := pipeline.NewRasterPipelineBuilder()
	blitMaskRpBuilder.SetForceHqPipeline(p.ForceHQPipeline)
	if !p.Shader.PushStages(p.Colorspace, blitMaskRpBuilder) {
		return nil
	}

	if maskCtx != nil {
		blitMaskRpBuilder.Push(pipeline.StageMaskU8)
	}

	if blendMode.ShouldPreScaleCoverage() {
		blitMaskRpBuilder.Push(pipeline.StageScaleU8)
		blitMaskRpBuilder.Push(pipeline.StageLoadDestination)
		if stage, ok := expand(p.Colorspace); ok {
			blitMaskRpBuilder.Push(stage)
		}
		if blendStage, ok := blendMode.stage(); ok {
			blitMaskRpBuilder.Push(blendStage)
		}
	} else {
		blitMaskRpBuilder.Push(pipeline.StageLoadDestination)
		if stage, ok := expand(p.Colorspace); ok {
			blitMaskRpBuilder.Push(stage)
		}
		if blendStage, ok := blendMode.stage(); ok {
			blitMaskRpBuilder.Push(blendStage)
		}
		blitMaskRpBuilder.Push(pipeline.StageLerpU8)
	}
	if stage, ok := compress(p.Colorspace); ok {
		blitMaskRpBuilder.Push(stage)
	}
	blitMaskRpBuilder.Push(pipeline.StageStore)
	blitMaskRp := blitMaskRpBuilder.Compile()

	var pixmapCtx *pipeline.PixmapCtx
	if pattern, ok := p.Shader.(*shader.Pattern); ok {
		pixmapCtx = &pipeline.PixmapCtx{
			Data: pattern.Data,
			Size: pattern.Size,
		}
	} else {
		size, _ := path.NewIntSize(1, 1)
		pixmapCtx = &pipeline.PixmapCtx{
			Data: []uint8{0, 0, 0, 0},
			Size: size,
		}
	}

	var maskCtx2 *pipeline.MaskCtx
	if maskCtx != nil {
		maskCtxVal := pipeline.MaskCtx{
			Data:      maskCtx.Data,
			RealWidth: maskCtx.RealWidth,
		}
		maskCtx2 = &maskCtxVal
	}

	var subPixmapCtx2 *pipeline.SubPixmapCtx
	if subPixmapCtx != nil {
		subPixmapCtxVal := pipeline.SubPixmapCtx{
			Data:      subPixmapCtx.Data,
			Size:      subPixmapCtx.Size,
			RealWidth: subPixmapCtx.RealWidth,
		}
		subPixmapCtx2 = &subPixmapCtxVal
	}

	return &pipeline.RasterPipelineBlitter{
		Mask:             maskCtx2,
		PixmapSrc:        pixmapCtx,
		Pixmap:           subPixmapCtx2,
		Memset2dColor:    memset2dColor,
		UseMemset2dColor: useMemset2dColor,
		BlitAntiHRp:      *blitAntiHRp,
		BlitRectRp:       *blitRectRp,
		BlitMaskRp:       *blitMaskRp,
		IsMask:           false,
	}
}

// SetColor sets a paint source to a solid color.
func (p *Paint) SetColor(color color.Color) {
	p.Shader = shader.NewSolidColor(color)
}

// SetColorRGBA8 sets a paint source to a solid color using RGBA8 values.
func (p *Paint) SetColorRGBA8(r, g, b, a uint8) {
	p.SetColor(color.ColorFromRGBA8(r, g, b, a))
}

// IsSolidColor checks that the paint source is a solid color.
func (p *Paint) IsSolidColor() bool {
	_, ok := p.Shader.(*shader.SolidColor)
	return ok
}

// Copy creates a deep copy of the Paint.
func (p *Paint) Copy() Paint {
	return Paint{
		Shader:          p.Shader,
		BlendMode:       p.BlendMode,
		AntiAlias:       p.AntiAlias,
		Colorspace:      p.Colorspace,
		ForceHQPipeline: p.ForceHQPipeline,
	}
}

func expand(self color.ColorSpace) (pipeline.Stage, bool) {
	switch self {
	case color.ColorSpaceLinear:
		return 0, false
	case color.ColorSpaceGamma2:
		return pipeline.StageGammaExpand2, true
	case color.ColorSpaceSimpleSRGB:
		return pipeline.StageGammaExpand22, true
	case color.ColorSpaceFullSRGBGamma:
		return pipeline.StageGammaExpandSrgb, true
	}
	return 0, false
}

func compress(self color.ColorSpace) (pipeline.Stage, bool) {
	switch self {
	case color.ColorSpaceLinear:
		return 0, false
	case color.ColorSpaceGamma2:
		return pipeline.StageGammaCompress2, true
	case color.ColorSpaceSimpleSRGB:
		return pipeline.StageGammaCompress22, true
	case color.ColorSpaceFullSRGBGamma:
		return pipeline.StageGammaCompressSrgb, true
	}
	return 0, false
}

func (b CompositeOperation) stage() (pipeline.Stage, bool) {
	switch b {
	case CompositeOperationClear:
		return pipeline.StageClear, true
	case CompositeOperationSource:
		return 0, false // This stage is a no-op.
	case CompositeOperationDestination:
		return pipeline.StageMoveDestinationToSource, true
	case CompositeOperationSourceOver:
		return pipeline.StageSourceOver, true
	case CompositeOperationDestinationOver:
		return pipeline.StageDestinationOver, true
	case CompositeOperationSourceIn:
		return pipeline.StageSourceIn, true
	case CompositeOperationDestinationIn:
		return pipeline.StageDestinationIn, true
	case CompositeOperationSourceOut:
		return pipeline.StageSourceOut, true
	case CompositeOperationDestinationOut:
		return pipeline.StageDestinationOut, true
	case CompositeOperationSourceAtop:
		return pipeline.StageSourceAtop, true
	case CompositeOperationDestinationAtop:
		return pipeline.StageDestinationAtop, true
	case CompositeOperationXor:
		return pipeline.StageXor, true
	case CompositeOperationPlus:
		return pipeline.StagePlus, true
	case CompositeOperationModulate:
		return pipeline.StageModulate, true
	case CompositeOperationScreen:
		return pipeline.StageScreen, true
	case CompositeOperationOverlay:
		return pipeline.StageOverlay, true
	case CompositeOperationDarken:
		return pipeline.StageDarken, true
	case CompositeOperationLighten:
		return pipeline.StageLighten, true
	case CompositeOperationColorDodge:
		return pipeline.StageColorDodge, true
	case CompositeOperationColorBurn:
		return pipeline.StageColorBurn, true
	case CompositeOperationHardLight:
		return pipeline.StageHardLight, true
	case CompositeOperationSoftLight:
		return pipeline.StageSoftLight, true
	case CompositeOperationDifference:
		return pipeline.StageDifference, true
	case CompositeOperationExclusion:
		return pipeline.StageExclusion, true
	case CompositeOperationMultiply:
		return pipeline.StageMultiply, true
	case CompositeOperationHue:
		return pipeline.StageHue, true
	case CompositeOperationSaturation:
		return pipeline.StageSaturation, true
	case CompositeOperationColor:
		return pipeline.StageColor, true
	case CompositeOperationLuminosity:
		return pipeline.StageLuminosity, true
	default:
		return 0, false
	}
}

func (b CompositeOperation) ShouldPreScaleCoverage() bool {
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
	case CompositeOperationDestination,
		CompositeOperationDestinationOver,
		CompositeOperationPlus,
		CompositeOperationDestinationOut,
		CompositeOperationSourceAtop,
		CompositeOperationSourceOver,
		CompositeOperationXor:
		return true
	default:
		return false
	}
}
