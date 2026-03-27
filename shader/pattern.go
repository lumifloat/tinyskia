// Copyright 2006 The Android Open Source Project
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package shader

import (
	"github.com/chewxy/math32"
	"github.com/lumifloat/tinyskia/color"
	"github.com/lumifloat/tinyskia/internal/normalized"
	"github.com/lumifloat/tinyskia/path"
	"github.com/lumifloat/tinyskia/pipeline"
)

// Controls how much filtering to be done when transforming images.
type FilterQuality int

const (
	// Nearest-neighbor. Low quality, but fastest.
	FilterQualityNearest FilterQuality = iota
	// Bilinear.
	FilterQualityBilinear
	// Bicubic. High quality, but slow.
	FilterQualityBicubic
)

// A pattern shader.
//
// Essentially a SkImageShader.
type Pattern struct {
	size      path.IntSize
	quality   FilterQuality
	spread    SpreadMode
	opacity   normalized.NormalizedF32
	transform path.Transform
}

func NewPattern(size path.IntSize, spread SpreadMode, quality FilterQuality, opacity float32, transform path.Transform) Shader {
	return &Pattern{
		size:      size,
		quality:   quality,
		spread:    spread,
		opacity:   normalized.NewNormalizedF32WithClamped(opacity),
		transform: transform,
	}
}

func (p *Pattern) IsOpaque() bool {
	return false
}

func (p *Pattern) PushStages(cs color.ColorSpace, builder *pipeline.RasterPipelineBuilder) bool {
	ts, ok := p.transform.Invert()
	if !ok {
		// failed to invert a pattern transform. Nothing will be rendered
		return false
	}

	builder.Push(pipeline.StageSeedShader)
	builder.PushTransform(ts)

	quality := p.quality

	if ts.IsIdentity() || ts.IsTranslate() {
		quality = FilterQualityNearest
	}

	if quality == FilterQualityBilinear {
		if ts.IsTranslate() {
			if ts.TX == math32.Trunc(ts.TX) && ts.TY == math32.Trunc(ts.TY) {
				// When the matrix is just an integer translate, bilerp == nearest neighbor.
				quality = FilterQualityNearest
			}
		}
	}

	switch quality {
	case FilterQualityNearest:
		builder.Ctx.LimitX = pipeline.TileCtx{
			Scale:    float32(p.size.Width()),
			InvScale: 1.0 / float32(p.size.Width()),
		}

		builder.Ctx.LimitY = pipeline.TileCtx{
			Scale:    float32(p.size.Height()),
			InvScale: 1.0 / float32(p.size.Height()),
		}

		switch p.spread {
		case SpreadModePad:
			// The gather() stage will clamp for us.
		case SpreadModeRepeat:
			builder.Push(pipeline.StageRepeat)
		case SpreadModeReflect:
			builder.Push(pipeline.StageReflect)
		}

		builder.Push(pipeline.StageGather)

	case FilterQualityBilinear:
		builder.Ctx.Sampler = pipeline.SamplerCtx{
			SpreadMode: pipeline.SpreadMode(p.spread),
			InvWidth:   1.0 / float32(p.size.Width()),
			InvHeight:  1.0 / float32(p.size.Height()),
		}
		builder.Push(pipeline.StageBilinear)

	case FilterQualityBicubic:
		builder.Ctx.Sampler = pipeline.SamplerCtx{
			SpreadMode: pipeline.SpreadMode(p.spread),
			InvWidth:   1.0 / float32(p.size.Width()),
			InvHeight:  1.0 / float32(p.size.Height()),
		}
		builder.Push(pipeline.StageBicubic)

		// Bicubic filtering naturally produces out of range values on both sides of [0,1].
		builder.Push(pipeline.StageClamp0)
		builder.Push(pipeline.StageClampA)
	}

	// Unlike Skia, we do not support global opacity and only Pattern allows it.
	if p.opacity != 1.0 {
		builder.Ctx.CurrentCoverage = p.opacity.Get()
		builder.Push(pipeline.StageScale1Float)
	}

	if cs != color.ColorSpaceLinear {
		switch cs {
		case color.ColorSpaceGamma2:
			builder.Push(pipeline.StageGammaExpand2)
		case color.ColorSpaceSimpleSRGB:
			builder.Push(pipeline.StageGammaExpand22)
		case color.ColorSpaceFullSRGBGamma:
			builder.Push(pipeline.StageGammaExpandSrgb)
		}
	}

	return true
}

func (p *Pattern) Transform(ts path.Transform) {
	p.transform = p.transform.PreConcat(ts)
}

func (p *Pattern) ApplyOpacity(opacity float32) {
	p.opacity = normalized.NewNormalizedF32WithClamped(p.opacity.Get() * opacity)
}
