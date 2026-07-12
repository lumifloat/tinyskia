// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package painter

import (
	"image"
	"image/color"
	"unsafe"

	"github.com/lumifloat/tinyskia/internal/core/colorf"
	"github.com/lumifloat/tinyskia/internal/core/pipeline"
	"github.com/lumifloat/tinyskia/internal/core/shader"
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
	Colorspace colorf.ColorSpace
	// Forces the high quality/precision rendering pipeline.
	ForceHQPipeline bool
}

// DefaultPaint returns a Paint with default values.
func DefaultPaint() Paint {
	return Paint{
		Shader:          shader.NewSolidColor(color.NRGBA{0, 0, 0, 255}),
		BlendMode:       CompositeOperationSourceOver,
		AntiAlias:       true,
		Colorspace:      colorf.ColorSpaceLinear,
		ForceHQPipeline: false,
	}
}

func (p *Paint) blitter(dst *image.RGBA, mask *image.Alpha) *pipeline.RasterPipelineBlitter {
	if dst == nil {
		return nil
	}
	if mask != nil && (mask.Rect.Dx() != dst.Rect.Dx() || mask.Rect.Dy() != dst.Rect.Dy()) {
		return nil
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
	if p.Shader.IsOpaque() && blendMode == CompositeOperationSourceOver && mask == nil {
		blendMode = CompositeOperationCopy
	}

	// When we're drawing a constant color in Source mode, we can sometimes just memset.
	var memset2dColor color.RGBA
	var useMemset2dColor bool
	if blendMode == CompositeOperationCopy && mask == nil {
		if solid, ok := p.Shader.(*shader.SolidColor); ok {
			memset2dColor = color.RGBAModel.Convert(solid.Color()).(color.RGBA)
			useMemset2dColor = true
		}
	}

	// Clear is just a transparent color memset.
	if blendMode == CompositeOperationClear && !p.AntiAlias && mask == nil {
		blendMode = CompositeOperationCopy
		memset2dColor = color.RGBA{0, 0, 0, 0}
		useMemset2dColor = true
	}

	// blit_anti_h_rp
	blitAntiHRpBuilder := pipeline.NewRasterPipelineBuilder()
	blitAntiHRpBuilder.SetForceHqPipeline(p.ForceHQPipeline)
	if !p.Shader.PushStages(p.Colorspace, blitAntiHRpBuilder) {
		return nil
	}

	if mask != nil {
		blitAntiHRpBuilder.Push(pipeline.StageMaskU8)
	}

	if blendMode.ShouldPreScaleCoverage() {
		blitAntiHRpBuilder.Push(pipeline.StageScale1Float)
		blitAntiHRpBuilder.Push(pipeline.StageLoadDestination)
		blitAntiHRpBuilder.PushColorSpaceExpand(p.Colorspace)
		if blendStage, ok := blendMode.stage(); ok {
			blitAntiHRpBuilder.Push(blendStage)
		}
	} else {
		blitAntiHRpBuilder.Push(pipeline.StageLoadDestination)
		blitAntiHRpBuilder.PushColorSpaceExpand(p.Colorspace)
		if blendStage, ok := blendMode.stage(); ok {
			blitAntiHRpBuilder.Push(blendStage)
		}
		blitAntiHRpBuilder.Push(pipeline.StageLerp1Float)
	}
	blitAntiHRpBuilder.PushColorSpaceCompress(p.Colorspace)
	blitAntiHRpBuilder.Push(pipeline.StageStore)
	blitAntiHRp := blitAntiHRpBuilder.Compile()

	// blit_rect_rp
	blitRectRpBuilder := pipeline.NewRasterPipelineBuilder()
	blitRectRpBuilder.SetForceHqPipeline(p.ForceHQPipeline)
	if !p.Shader.PushStages(p.Colorspace, blitRectRpBuilder) {
		return nil
	}

	if mask != nil {
		blitRectRpBuilder.Push(pipeline.StageMaskU8)
	}

	if blendMode == CompositeOperationSourceOver && mask == nil {
		blitRectRpBuilder.PushColorSpaceCompress(p.Colorspace)
		// TODO: ignore when dither_rate is non-zero
		blitRectRpBuilder.Push(pipeline.StageSourceOverRgba)
	} else {
		if blendMode != CompositeOperationCopy {
			blitRectRpBuilder.Push(pipeline.StageLoadDestination)
			if blendStage, ok := blendMode.stage(); ok {
				blitRectRpBuilder.PushColorSpaceExpand(p.Colorspace)
				blitRectRpBuilder.Push(blendStage)
			}
		}
		blitRectRpBuilder.PushColorSpaceCompress(p.Colorspace)
		blitRectRpBuilder.Push(pipeline.StageStore)
	}
	blitRectRp := blitRectRpBuilder.Compile()

	// blit_mask_rp
	blitMaskRpBuilder := pipeline.NewRasterPipelineBuilder()
	blitMaskRpBuilder.SetForceHqPipeline(p.ForceHQPipeline)
	if !p.Shader.PushStages(p.Colorspace, blitMaskRpBuilder) {
		return nil
	}

	if mask != nil {
		blitMaskRpBuilder.Push(pipeline.StageMaskU8)
	}

	if blendMode.ShouldPreScaleCoverage() {
		blitMaskRpBuilder.Push(pipeline.StageScaleU8)
		blitMaskRpBuilder.Push(pipeline.StageLoadDestination)
		blitMaskRpBuilder.PushColorSpaceExpand(p.Colorspace)
		if blendStage, ok := blendMode.stage(); ok {
			blitMaskRpBuilder.Push(blendStage)
		}
	} else {
		blitMaskRpBuilder.Push(pipeline.StageLoadDestination)
		blitMaskRpBuilder.PushColorSpaceExpand(p.Colorspace)
		if blendStage, ok := blendMode.stage(); ok {
			blitMaskRpBuilder.Push(blendStage)
		}
		blitMaskRpBuilder.Push(pipeline.StageLerpU8)
	}
	blitMaskRpBuilder.PushColorSpaceCompress(p.Colorspace)
	blitMaskRpBuilder.Push(pipeline.StageStore)
	blitMaskRp := blitMaskRpBuilder.Compile()

	var src *image.RGBA
	if pattern, ok := p.Shader.(*shader.Pattern); ok {
		src = image.NewRGBA(image.Rect(0, 0, int(pattern.Size.Width()), int(pattern.Size.Height())))
		src.Pix = pattern.Data
	} else {
		src = image.NewRGBA(image.Rect(0, 0, 1, 1))
		src.Pix = []uint8{0, 0, 0, 0}
	}

	return &pipeline.RasterPipelineBlitter{
		Mask:             mask,
		Src:              src,
		Dst:              dst,
		Memset2dColor:    memset2dColor,
		UseMemset2dColor: useMemset2dColor,
		BlitAntiHRp:      *blitAntiHRp,
		BlitRectRp:       *blitRectRp,
		BlitMaskRp:       *blitMaskRp,
		IsMask:           false,
	}
}

func (p *Paint) NewMaskBlitter(dst *image.Alpha, mask *image.Alpha) *pipeline.RasterPipelineBlitter {
	dst0 := (*image.RGBA)(unsafe.Pointer(dst))

	uc := pipeline.UniformColorCtx{
		R0: 0xff, G0: 0xff, B0: 0xff, A0: 0xff,
		R1: 1.0, G1: 1.0, B1: 1.0, A1: 1.0,
	}

	blitAntiHRpBuilder := pipeline.NewRasterPipelineBuilder()
	blitAntiHRpBuilder.PushUniformColor(uc)
	if mask != nil {
		blitAntiHRpBuilder.Push(pipeline.StageMaskU8)
	}
	blitAntiHRpBuilder.Push(pipeline.StageLoadDestinationU8)
	blitAntiHRpBuilder.Push(pipeline.StageLerp1Float)
	blitAntiHRpBuilder.Push(pipeline.StageStoreU8)

	blitRectRpBuilder := pipeline.NewRasterPipelineBuilder()
	blitRectRpBuilder.PushUniformColor(uc)
	if mask != nil {
		blitRectRpBuilder.Push(pipeline.StageMaskU8)
	}
	blitRectRpBuilder.Push(pipeline.StageStoreU8)

	blitMaskRpBuilder := pipeline.NewRasterPipelineBuilder()
	blitMaskRpBuilder.PushUniformColor(uc)
	if mask != nil {
		blitMaskRpBuilder.Push(pipeline.StageMaskU8)
	}
	blitMaskRpBuilder.Push(pipeline.StageLoadDestinationU8)
	blitMaskRpBuilder.Push(pipeline.StageLerpU8)
	blitMaskRpBuilder.Push(pipeline.StageStoreU8)

	return &pipeline.RasterPipelineBlitter{
		Mask:        mask,
		Src:         nil,
		Dst:         dst0,
		BlitAntiHRp: *blitAntiHRpBuilder.Compile(),
		BlitRectRp:  *blitRectRpBuilder.Compile(),
		BlitMaskRp:  *blitMaskRpBuilder.Compile(),
		IsMask:      true,
	}
}

// SetColor sets a paint source to a solid color.
func (p *Paint) SetColor(color color.Color) {
	p.Shader = shader.NewSolidColor(color)
}

// SetColorRGBA8 sets a paint source to a solid color using RGBA8 values.
func (p *Paint) SetColorRGBA8(r, g, b, a uint8) {
	p.SetColor(color.NRGBA{r, g, b, a})
}

// IsSolidColor checks that the paint source is a solid color.
func (p *Paint) IsSolidColor() bool {
	_, ok := p.Shader.(*shader.SolidColor)
	return ok
}

// Copy creates a deep copy of the Paint.
func (p *Paint) Copy() *Paint {
	return &Paint{
		Shader:          p.Shader,
		BlendMode:       p.BlendMode,
		AntiAlias:       p.AntiAlias,
		Colorspace:      p.Colorspace,
		ForceHQPipeline: p.ForceHQPipeline,
	}
}

func (b CompositeOperation) stage() (pipeline.Stage, bool) {
	switch b {
	case CompositeOperationClear:
		return pipeline.StageClear, true
	case CompositeOperationCopy:
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
	case CompositeOperationLighter:
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
		CompositeOperationLighter,
		CompositeOperationDestinationOut,
		CompositeOperationSourceAtop,
		CompositeOperationSourceOver,
		CompositeOperationXor:
		return true
	default:
		return false
	}
}
