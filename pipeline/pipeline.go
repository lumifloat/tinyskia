// Copyright 2016 Google Inc.
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package pipeline

import (
	"github.com/lumifloat/tinyskia/color"
	"github.com/lumifloat/tinyskia/internal/normalized"
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
	Factors []GradientColor
	Biases  []GradientColor
	TValues []normalized.NormalizedF32
}

func (g *GradientCtx) PushConstColor(color GradientColor) {
	g.Factors = append(g.Factors, NewGradientColor(0, 0, 0, 0))
	g.Biases = append(g.Biases, color)
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

func (r RasterPipelineLow) isRasterPipelineKind() {}

type RasterPipeline struct {
	Kind RasterPipelineKind
	Ctx  Context
}

func (b *RasterPipelineBuilder) Compile() *RasterPipeline {
	if len(b.Stages) == 0 {
		return &RasterPipeline{
			Kind: RasterPipelineHigh{},
			Ctx:  Context{},
		}
	}

	isLowpCompatible := true
	for _, s := range b.Stages {
		var table = map[Stage]bool{
			StageMoveSourceToDestination:        true,
			StageMoveDestinationToSource:        true,
			StageClamp0:                         false,
			StageClampA:                         false,
			StagePremultiply:                    true,
			StageUniformColor:                   true,
			StageSeedShader:                     true,
			StageLoadDestination:                true,
			StageStore:                          true,
			StageLoadDestinationU8:              true,
			StageStoreU8:                        true,
			StageGather:                         false,
			StageLoadMaskU8:                     true,
			StageMaskU8:                         true,
			StageScaleU8:                        true,
			StageLerpU8:                         true,
			StageScale1Float:                    true,
			StageLerp1Float:                     true,
			StageDestinationAtop:                true,
			StageDestinationIn:                  true,
			StageDestinationOut:                 true,
			StageDestinationOver:                true,
			StageSourceAtop:                     true,
			StageSourceIn:                       true,
			StageSourceOut:                      true,
			StageSourceOver:                     true,
			StageClear:                          true,
			StageModulate:                       true,
			StageMultiply:                       true,
			StagePlus:                           true,
			StageScreen:                         true,
			StageXor:                            true,
			StageColorBurn:                      false,
			StageColorDodge:                     false,
			StageDarken:                         true,
			StageDifference:                     true,
			StageExclusion:                      true,
			StageHardLight:                      true,
			StageLighten:                        true,
			StageOverlay:                        true,
			StageSoftLight:                      false,
			StageHue:                            false,
			StageSaturation:                     false,
			StageColor:                          false,
			StageLuminosity:                     false,
			StageSourceOverRgba:                 true,
			StageTransform:                      true,
			StageReflect:                        false,
			StageRepeat:                         false,
			StageBilinear:                       false,
			StageBicubic:                        false,
			StagePadX1:                          true,
			StageReflectX1:                      true,
			StageRepeatX1:                       true,
			StageGradient:                       true,
			StageEvenlySpaced2StopGradient:      true,
			StageXYToUnitAngle:                  false,
			StageXYToRadius:                     false,
			StageXYTo2PtConicalFocalOnCircle:    false,
			StageXYTo2PtConicalWellBehaved:      false,
			StageXYTo2PtConicalSmaller:          false,
			StageXYTo2PtConicalGreater:          false,
			StageXYTo2PtConicalStrip:            false,
			StageMask2PtConicalNan:              false,
			StageMask2PtConicalDegenerates:      false,
			StageApplyVectorMask:                false,
			StageAlter2PtConicalCompensateFocal: false,
			StageAlter2PtConicalUnswap:          false,
			StageNegateX:                        false,
			StageApplyConcentricScaleBias:       false,
			StageGammaExpand2:                   false,
			StageGammaExpandDestination2:        false,
			StageGammaCompress2:                 false,
			StageGammaExpand22:                  false,
			StageGammaExpandDestination22:       false,
			StageGammaCompress22:                false,
			StageGammaExpandSrgb:                false,
			StageGammaExpandDestinationSrgb:     false,
			StageGammaCompressSrgb:              false,
		}
		if !table[s] {
			isLowpCompatible = false
			break
		}
	}

	if b.ForceHqPipeline || !isLowpCompatible {
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
		StartHighPipeline(k, rect, aaMaskCtx, maskCtx, &p.Ctx, pixmapSrc, pixmapDst)
	case RasterPipelineLow:
		StartLowPipeline(k, rect, aaMaskCtx, maskCtx, &p.Ctx, pixmapDst)
	}
}
