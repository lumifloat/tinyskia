// Copyright 2016 Google Inc.
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package pipeline

import (
	"github.com/lumifloat/tinyskia/color"
	"github.com/lumifloat/tinyskia/path"
)

// SpreadMode is an alias for shaders.SpreadMode to avoid circular imports.
// Defined here as int to match the shaders package type.
type SpreadMode int

const (
	SpreadModePad SpreadMode = iota
	SpreadModeReflect
	SpreadModeRepeat
)

const MAX_STAGES = 32

type Stage int

const (
	StageMoveSourceToDestination Stage = iota
	StageMoveDestinationToSource
	StageClamp0
	StageClampA
	StagePremultiply
	StageUniformColor
	StageSeedShader
	StageLoadDestination
	StageStore
	StageLoadDestinationU8
	StageStoreU8
	StageGather
	StageLoadMaskU8
	StageMaskU8
	StageScaleU8
	StageLerpU8
	StageScale1Float
	StageLerp1Float
	StageDestinationAtop
	StageDestinationIn
	StageDestinationOut
	StageDestinationOver
	StageSourceAtop
	StageSourceIn
	StageSourceOut
	StageSourceOver
	StageClear
	StageModulate
	StageMultiply
	StagePlus
	StageScreen
	StageXor
	StageColorBurn
	StageColorDodge
	StageDarken
	StageDifference
	StageExclusion
	StageHardLight
	StageLighten
	StageOverlay
	StageSoftLight
	StageHue
	StageSaturation
	StageColor
	StageLuminosity
	StageSourceOverRgba
	StageTransform
	StageReflect
	StageRepeat
	StageBilinear
	StageBicubic
	StagePadX1
	StageReflectX1
	StageRepeatX1
	StageGradient
	StageEvenlySpaced2StopGradient
	StageXYToUnitAngle
	StageXYToRadius
	StageXYTo2PtConicalFocalOnCircle
	StageXYTo2PtConicalWellBehaved
	StageXYTo2PtConicalSmaller
	StageXYTo2PtConicalGreater
	StageXYTo2PtConicalStrip
	StageMask2PtConicalNan
	StageMask2PtConicalDegenerates
	StageApplyVectorMask
	StageAlter2PtConicalCompensateFocal
	StageAlter2PtConicalUnswap
	StageNegateX
	StageApplyConcentricScaleBias
	StageGammaExpand2
	StageGammaExpandDestination2
	StageGammaCompress2
	StageGammaExpand22
	StageGammaExpandDestination22
	StageGammaCompress22
	StageGammaExpandSrgb
	StageGammaExpandDestinationSrgb
	StageGammaCompressSrgb
)

const STAGES_COUNT = int(StageGammaCompressSrgb) + 1

type PixmapCtx struct {
	Data []uint8
	Size path.IntSize
}

type SubPixmapCtx struct {
	Data      []uint8
	Size      path.IntSize
	RealWidth int
}

func (c *SubPixmapCtx) Offset(dx, dy int) int {
	return c.RealWidth*dy + dx
}

type AAMaskCtx struct {
	Pixels [2]uint8
	Stride uint32
	Shift  int
}

func (c *AAMaskCtx) CopyAtXY(dx, dy, tail int) [2]uint8 {
	offset := int(c.Stride)*dy + dx - c.Shift
	switch {
	case offset == 0 && tail == 1:
		return [2]byte{c.Pixels[0], 0}
	case offset == 0 && tail == 2:
		return [2]byte{c.Pixels[0], c.Pixels[1]}
	case offset == 1 && tail == 1:
		return [2]byte{c.Pixels[1], 0}
	default:
		return [2]byte{0, 0}
	}
}

type MaskCtx struct {
	Data      []uint8
	RealWidth uint32
}

func (c *MaskCtx) Offset(dx, dy int) int {
	return int(c.RealWidth)*dy + dx
}

type Context struct {
	CurrentCoverage           float32
	Sampler                   SamplerCtx
	UniformColor              UniformColorCtx
	EvenlySpaced2StopGradient EvenlySpaced2StopGradientCtx
	Gradient                  GradientCtx
	TwoPointConicalGradient   TwoPointConicalGradientCtx
	LimitX                    TileCtx
	LimitY                    TileCtx
	Transform                 path.Transform
}

type SamplerCtx struct {
	SpreadMode SpreadMode
	InvWidth   float32
	InvHeight  float32
}

type UniformColorCtx struct {
	R, G, B, A float32
	RGBA       [4]uint16
}

type GradientColor struct {
	R, G, B, A float32
}

func NewGradientColor(r, g, b, a float32) GradientColor {
	return GradientColor{r, g, b, a}
}

type EvenlySpaced2StopGradientCtx struct {
	Factor GradientColor
	Bias   GradientColor
}

type GradientCtx struct {
	Len     int
	Factors [16]GradientColor
	Biases  [16]GradientColor
	TValues [16]float32 // NormalizedF32
}

func (g *GradientCtx) PushConstColor(color GradientColor) {
	if g.Len < 16 {
		g.Factors[g.Len] = NewGradientColor(0, 0, 0, 0)
		g.Biases[g.Len] = color
		g.TValues[g.Len] = 0
		g.Len++
	}
}

type TwoPointConicalGradientCtx struct {
	Mask [8]uint32 // u32x8
	P0   float32
	P1   float32
}

type TileCtx struct {
	Scale    float32
	InvScale float32
}

type RasterPipelineBuilder struct {
	Stages          []Stage
	ForceHqPipeline bool
	Ctx             Context
}

func NewRasterPipelineBuilder() *RasterPipelineBuilder {
	return &RasterPipelineBuilder{
		Stages: make([]Stage, 0, MAX_STAGES),
		Ctx:    Context{},
	}
}

func (b *RasterPipelineBuilder) SetForceHqPipeline(hq bool) {
	b.ForceHqPipeline = hq
}

func (b *RasterPipelineBuilder) Push(stage Stage) {
	if len(b.Stages) < MAX_STAGES {
		b.Stages = append(b.Stages, stage)
	}
}

func (b *RasterPipelineBuilder) PushTransform(ts path.Transform) {
	if ts.IsFinite() && !ts.IsIdentity() {
		b.Stages = append(b.Stages, StageTransform)
		b.Ctx.Transform = ts
	}
}

func (b *RasterPipelineBuilder) PushUniformColor(c color.PremultipliedColor) {
	r := c.Red()
	g := c.Green()
	bl := c.Blue()
	a := c.Alpha()
	rgba := [4]uint16{
		uint16(r*255.0 + 0.5),
		uint16(g*255.0 + 0.5),
		uint16(bl*255.0 + 0.5),
		uint16(a*255.0 + 0.5),
	}

	b.Ctx.UniformColor = UniformColorCtx{R: r, G: g, B: bl, A: a, RGBA: rgba}
	b.Stages = append(b.Stages, StageUniformColor)
}

type RasterPipelineKind interface {
	isRasterPipelineKind()
}

type RasterPipelineHigh []Stage

func (r RasterPipelineHigh) isRasterPipelineKind() {}

type RasterPipelineLow []Stage

type RasterBatchLow struct {
	Functions []LowpStageFn
}

func (r RasterPipelineLow) isRasterPipelineKind() {}

type RasterPipeline struct {
	Kind RasterPipelineKind
	Ctx  Context
}

// HighpStageFn and LowpStageFn would be defined in highp and lowp packages
type HighpStageFn func()
type LowpStageFn func(p *LowPipeline)

// GetHighpStage returns the highp stage function for a given stage.
// This is a placeholder - actual implementation would be in highp package.
func GetHighpStage(s Stage) HighpStageFn {
	// TODO: Implement stage function lookup
	return nil
}

// GetHighpTailVariant returns the tail variant of a highp stage function.
// This is a placeholder - actual implementation would be in highp package.
func GetHighpTailVariant(fn HighpStageFn) HighpStageFn {
	// TODO: Implement tail variant lookup
	return fn
}

// GetLowpStage returns the lowp stage function for a given stage.
// This is a placeholder - actual implementation would be in lowp package.
func GetLowpStage(s Stage) LowpStageFn {
	// TODO: Implement stage function lookup
	return nil
}

// GetLowpTailVariant returns the tail variant of a lowp stage function.
// This is a placeholder - actual implementation would be in lowp package.
func GetLowpTailVariant(fn LowpStageFn) LowpStageFn {
	// TODO: Implement tail variant lookup
	return fn
}

func (b *RasterPipelineBuilder) Compile() *RasterPipeline {
	if len(b.Stages) == 0 {
		return &RasterPipeline{
			Kind: RasterPipelineHigh{},
			Ctx:  Context{},
		}
	}

	// In Go, we'd typically look up function implementations from a registry
	// This mirrors the logic of checking lowp compatibility and cloning for tail functions
	isLowpCompatible := true
	// for _, s := range b.Stages {
	// 	if GetLowpStage(s) == nil {
	// 		isLowpCompatible = false
	// 		break
	// 	}
	// }

	if b.ForceHqPipeline || !isLowpCompatible {
		fns := make([]HighpStageFn, len(b.Stages))
		for i, s := range b.Stages {
			fns[i] = GetHighpStage(s)
		}
		// Tail logic would involve replacing specific stages with their _tail variants
		tailFns := make([]HighpStageFn, len(fns))
		copy(tailFns, fns)
		for i, fn := range tailFns {
			tailFns[i] = GetHighpTailVariant(fn)
		}

		return &RasterPipeline{
			Kind: RasterPipelineHigh(b.Stages),
			Ctx:  b.Ctx,
		}
	} else {
		return &RasterPipeline{
			Kind: RasterPipelineLow(b.Stages),
			Ctx:  b.Ctx,
		}
	}
}

func (p *RasterPipeline) Run(
	rect *path.ScreenIntRect,
	aaMaskCtx *AAMaskCtx,
	maskCtx *MaskCtx,
	pixmapSrc *PixmapCtx,
	pixmapDst *SubPixmapCtx,
) {
	switch k := p.Kind.(type) {
	case RasterPipelineHigh:
		// TODO: Implement highp pipeline execution
		StartHighPipeline(k, rect, aaMaskCtx, maskCtx, &p.Ctx, pixmapSrc, pixmapDst)
	case RasterPipelineLow:
		StartLowPipeline(k, rect, aaMaskCtx, maskCtx, &p.Ctx, pixmapDst)
	}
}
