// Copyright 2016 Google Inc.
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package pipeline

import (
	"math"
)

//go:fix inline
func (p *LowPipeline) MoveSourceToDestination() {
	copy(p.dr[:], p.r[:])
	copy(p.dg[:], p.g[:])
	copy(p.db[:], p.b[:])
	copy(p.da[:], p.a[:])
}

//go:fix inline
func (p *LowPipeline) MoveDestinationToSource() {
	copy(p.r[:], p.dr[:])
	copy(p.g[:], p.dg[:])
	copy(p.b[:], p.db[:])
	copy(p.a[:], p.da[:])
}

//go:fix inline
func (p *LowPipeline) Clamp0() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) ClampA() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) Premultiply() {
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		p.r[i] = uint16(((uint32(p.r[i]) * uint32(p.a[i])) + 255) >> 8)
		p.r[i+1] = uint16(((uint32(p.r[i+1]) * uint32(p.a[i+1])) + 255) >> 8)
		p.r[i+2] = uint16(((uint32(p.r[i+2]) * uint32(p.a[i+2])) + 255) >> 8)
		p.r[i+3] = uint16(((uint32(p.r[i+3]) * uint32(p.a[i+3])) + 255) >> 8)
		p.r[i+4] = uint16(((uint32(p.r[i+4]) * uint32(p.a[i+4])) + 255) >> 8)
		p.r[i+5] = uint16(((uint32(p.r[i+5]) * uint32(p.a[i+5])) + 255) >> 8)
		p.r[i+6] = uint16(((uint32(p.r[i+6]) * uint32(p.a[i+6])) + 255) >> 8)
		p.r[i+7] = uint16(((uint32(p.r[i+7]) * uint32(p.a[i+7])) + 255) >> 8)

		p.g[i] = uint16(((uint32(p.g[i]) * uint32(p.a[i])) + 255) >> 8)
		p.g[i+1] = uint16(((uint32(p.g[i+1]) * uint32(p.a[i+1])) + 255) >> 8)
		p.g[i+2] = uint16(((uint32(p.g[i+2]) * uint32(p.a[i+2])) + 255) >> 8)
		p.g[i+3] = uint16(((uint32(p.g[i+3]) * uint32(p.a[i+3])) + 255) >> 8)
		p.g[i+4] = uint16(((uint32(p.g[i+4]) * uint32(p.a[i+4])) + 255) >> 8)
		p.g[i+5] = uint16(((uint32(p.g[i+5]) * uint32(p.a[i+5])) + 255) >> 8)
		p.g[i+6] = uint16(((uint32(p.g[i+6]) * uint32(p.a[i+6])) + 255) >> 8)
		p.g[i+7] = uint16(((uint32(p.g[i+7]) * uint32(p.a[i+7])) + 255) >> 8)

		p.b[i] = uint16(((uint32(p.b[i]) * uint32(p.a[i])) + 255) >> 8)
		p.b[i+1] = uint16(((uint32(p.b[i+1]) * uint32(p.a[i+1])) + 255) >> 8)
		p.b[i+2] = uint16(((uint32(p.b[i+2]) * uint32(p.a[i+2])) + 255) >> 8)
		p.b[i+3] = uint16(((uint32(p.b[i+3]) * uint32(p.a[i+3])) + 255) >> 8)
		p.b[i+4] = uint16(((uint32(p.b[i+4]) * uint32(p.a[i+4])) + 255) >> 8)
		p.b[i+5] = uint16(((uint32(p.b[i+5]) * uint32(p.a[i+5])) + 255) >> 8)
		p.b[i+6] = uint16(((uint32(p.b[i+6]) * uint32(p.a[i+6])) + 255) >> 8)
		p.b[i+7] = uint16(((uint32(p.b[i+7]) * uint32(p.a[i+7])) + 255) >> 8)
	}
}

//go:fix inline
func (p *LowPipeline) UniformColor() {
	uniformColor := p.ctx.UniformColor
	r := uint16(uniformColor.RGBA[0])
	g := uint16(uniformColor.RGBA[1])
	b := uint16(uniformColor.RGBA[2])
	a := uint16(uniformColor.RGBA[3])

	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		p.r[i], p.r[i+1], p.r[i+2], p.r[i+3] = r, r, r, r
		p.r[i+4], p.r[i+5], p.r[i+6], p.r[i+7] = r, r, r, r

		p.g[i], p.g[i+1], p.g[i+2], p.g[i+3] = g, g, g, g
		p.g[i+4], p.g[i+5], p.g[i+6], p.g[i+7] = g, g, g, g

		p.b[i], p.b[i+1], p.b[i+2], p.b[i+3] = b, b, b, b
		p.b[i+4], p.b[i+5], p.b[i+6], p.b[i+7] = b, b, b, b

		p.a[i], p.a[i+1], p.a[i+2], p.a[i+3] = a, a, a, a
		p.a[i+4], p.a[i+5], p.a[i+6], p.a[i+7] = a, a, a, a
	}
}

//go:fix inline
func (p *LowPipeline) SeedShader() {
	// Sets up pixel coordinates for shader processing
	// x = dx + [0.5, 1.5, 2.5, ..., 15.5]
	// y = dy + 0.5 (constant for all pixels)
	iota := [16]float32{0.5, 1.5, 2.5, 3.5, 4.5, 5.5, 6.5, 7.5, 8.5, 9.5, 10.5, 11.5, 12.5, 13.5, 14.5, 15.5}
	dxFloat := float32(p.dx)
	dyFloat := float32(p.dy) + 0.5

	// Calculate x coordinates: dx + iota
	for i := 0; i < 16; i++ {
		x := dxFloat + iota[i]
		xBits := math.Float32bits(x)
		p.r[i] = uint16(xBits & 0xFFFF)
		p.g[i] = uint16(xBits >> 16)
	}

	// Calculate y coordinates: dy + 0.5 (constant)
	yBits := math.Float32bits(dyFloat)
	for i := 0; i < 16; i++ {
		p.b[i] = uint16(yBits & 0xFFFF)
		p.a[i] = uint16(yBits >> 16)
	}

}

//go:fix inline
func (p *LowPipeline) LoadDestination() {
	baseIdx := (p.dy*p.pixmap.RealWidth + p.dx) * 4
	data := p.pixmap.Data[baseIdx : baseIdx+LOW_STAGE_WIDTH*4]

	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		off0 := i * 4
		off1 := off0 + 4
		off2 := off0 + 8
		off3 := off0 + 12
		off4 := off0 + 16
		off5 := off0 + 20
		off6 := off0 + 24
		off7 := off0 + 28

		p.dr[i] = uint16(data[off0])
		p.dr[i+1] = uint16(data[off1])
		p.dr[i+2] = uint16(data[off2])
		p.dr[i+3] = uint16(data[off3])
		p.dr[i+4] = uint16(data[off4])
		p.dr[i+5] = uint16(data[off5])
		p.dr[i+6] = uint16(data[off6])
		p.dr[i+7] = uint16(data[off7])

		p.dg[i] = uint16(data[off0+1])
		p.dg[i+1] = uint16(data[off1+1])
		p.dg[i+2] = uint16(data[off2+1])
		p.dg[i+3] = uint16(data[off3+1])
		p.dg[i+4] = uint16(data[off4+1])
		p.dg[i+5] = uint16(data[off5+1])
		p.dg[i+6] = uint16(data[off6+1])
		p.dg[i+7] = uint16(data[off7+1])

		p.db[i] = uint16(data[off0+2])
		p.db[i+1] = uint16(data[off1+2])
		p.db[i+2] = uint16(data[off2+2])
		p.db[i+3] = uint16(data[off3+2])
		p.db[i+4] = uint16(data[off4+2])
		p.db[i+5] = uint16(data[off5+2])
		p.db[i+6] = uint16(data[off6+2])
		p.db[i+7] = uint16(data[off7+2])

		p.da[i] = uint16(data[off0+3])
		p.da[i+1] = uint16(data[off1+3])
		p.da[i+2] = uint16(data[off2+3])
		p.da[i+3] = uint16(data[off3+3])
		p.da[i+4] = uint16(data[off4+3])
		p.da[i+5] = uint16(data[off5+3])
		p.da[i+6] = uint16(data[off6+3])
		p.da[i+7] = uint16(data[off7+3])
	}
}

//go:fix inline
func (p *LowPipeline) LoadDestinationTail() {
	baseIdx := (p.dy*p.pixmap.RealWidth + p.dx) * 4
	data := p.pixmap.Data[baseIdx : baseIdx+p.tail*4]
	for i := 0; i < p.tail; i++ {
		off := i * 4
		p.dr[i] = uint16(data[off])
		p.dg[i] = uint16(data[off+1])
		p.db[i] = uint16(data[off+2])
		p.da[i] = uint16(data[off+3])
	}
}

//go:fix inline
func (p *LowPipeline) Store() {
	baseIdx := (p.dy*p.pixmap.RealWidth + p.dx) * 4
	data := p.pixmap.Data[baseIdx : baseIdx+p.tail*4]

	for i := 0; i < p.tail; i += 8 {
		off0 := i * 4
		off1 := off0 + 4
		off2 := off0 + 8
		off3 := off0 + 12
		off4 := off0 + 16
		off5 := off0 + 20
		off6 := off0 + 24
		off7 := off0 + 28

		data[off0] = uint8(p.r[i])
		data[off1] = uint8(p.r[i+1])
		data[off2] = uint8(p.r[i+2])
		data[off3] = uint8(p.r[i+3])
		data[off4] = uint8(p.r[i+4])
		data[off5] = uint8(p.r[i+5])
		data[off6] = uint8(p.r[i+6])
		data[off7] = uint8(p.r[i+7])

		data[off0+1] = uint8(p.g[i])
		data[off1+1] = uint8(p.g[i+1])
		data[off2+1] = uint8(p.g[i+2])
		data[off3+1] = uint8(p.g[i+3])
		data[off4+1] = uint8(p.g[i+4])
		data[off5+1] = uint8(p.g[i+5])
		data[off6+1] = uint8(p.g[i+6])
		data[off7+1] = uint8(p.g[i+7])

		data[off0+2] = uint8(p.b[i])
		data[off1+2] = uint8(p.b[i+1])
		data[off2+2] = uint8(p.b[i+2])
		data[off3+2] = uint8(p.b[i+3])
		data[off4+2] = uint8(p.b[i+4])
		data[off5+2] = uint8(p.b[i+5])
		data[off6+2] = uint8(p.b[i+6])
		data[off7+2] = uint8(p.b[i+7])

		data[off0+3] = uint8(p.a[i])
		data[off1+3] = uint8(p.a[i+1])
		data[off2+3] = uint8(p.a[i+2])
		data[off3+3] = uint8(p.a[i+3])
		data[off4+3] = uint8(p.a[i+4])
		data[off5+3] = uint8(p.a[i+5])
		data[off6+3] = uint8(p.a[i+6])
		data[off7+3] = uint8(p.a[i+7])
	}
}

//go:fix inline
func (p *LowPipeline) StoreTail() {
	baseIdx := (p.dy*p.pixmap.RealWidth + p.dx) * 4
	data := p.pixmap.Data[baseIdx : baseIdx+p.tail*4]
	for i := 0; i < p.tail; i++ {
		off := i * 4
		data[off] = uint8(p.r[i])
		data[off+1] = uint8(p.g[i])
		data[off+2] = uint8(p.b[i])
		data[off+3] = uint8(p.a[i])
	}
}

//go:fix inline
func (p *LowPipeline) LoadDestinationU8() {
	baseIdx := (p.dy*p.pixmap.RealWidth + p.dx) * 4
	data := p.pixmap.Data[baseIdx : baseIdx+LOW_STAGE_WIDTH*4]
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		off0 := i * 4
		p.dr[i] = uint16(data[off0])
		p.dr[i+1] = uint16(data[off0+4])
		p.dr[i+2] = uint16(data[off0+8])
		p.dr[i+3] = uint16(data[off0+12])
		p.dr[i+4] = uint16(data[off0+16])
		p.dr[i+5] = uint16(data[off0+20])
		p.dr[i+6] = uint16(data[off0+24])
		p.dr[i+7] = uint16(data[off0+28])

		p.dg[i] = uint16(data[off0+1])
		p.dg[i+1] = uint16(data[off0+5])
		p.dg[i+2] = uint16(data[off0+9])
		p.dg[i+3] = uint16(data[off0+13])
		p.dg[i+4] = uint16(data[off0+17])
		p.dg[i+5] = uint16(data[off0+21])
		p.dg[i+6] = uint16(data[off0+25])
		p.dg[i+7] = uint16(data[off0+29])

		p.db[i] = uint16(data[off0+2])
		p.db[i+1] = uint16(data[off0+6])
		p.db[i+2] = uint16(data[off0+10])
		p.db[i+3] = uint16(data[off0+14])
		p.db[i+4] = uint16(data[off0+18])
		p.db[i+5] = uint16(data[off0+22])
		p.db[i+6] = uint16(data[off0+26])
		p.db[i+7] = uint16(data[off0+30])

		p.da[i] = uint16(data[off0+3])
		p.da[i+1] = uint16(data[off0+7])
		p.da[i+2] = uint16(data[off0+11])
		p.da[i+3] = uint16(data[off0+15])
		p.da[i+4] = uint16(data[off0+19])
		p.da[i+5] = uint16(data[off0+23])
		p.da[i+6] = uint16(data[off0+27])
		p.da[i+7] = uint16(data[off0+31])
	}
}

//go:fix inline
func (p *LowPipeline) LoadDestinationU8Tail() {
	// TODO
}

//go:fix inline
func (p *LowPipeline) StoreU8() {
	baseIdx := (p.dy*p.pixmap.RealWidth + p.dx) * 4
	data := p.pixmap.Data[baseIdx : baseIdx+p.tail*4]
	for i := 0; i < p.tail; i++ {
		off := i * 4
		data[off] = uint8(p.r[i])
		data[off+1] = uint8(p.g[i])
		data[off+2] = uint8(p.b[i])
		data[off+3] = uint8(p.a[i])
	}
}

//go:fix inline
func (p *LowPipeline) StoreU8Tail() {
	// TODO
}

//go:fix inline
func (p *LowPipeline) Gather() {
	// This stage samples colors from a texture based on UV coordinates
	// stored in p.r/p.g and uses the shader context
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		// Gather would sample from texture using UV coordinates
		// For lowp, this is typically handled by highp pipeline
	}
}

//go:fix inline
func (p *LowPipeline) LoadMaskU8() {
	if p.maskCtx != nil && len(p.maskCtx.Data) > 0 {
		baseIdx := int(p.maskCtx.RealWidth)*p.dy + p.dx
		maskData := p.maskCtx.Data
		for i := 0; i < LOW_STAGE_WIDTH && baseIdx+i < len(maskData); i += 8 {
			p.a[i] = uint16(maskData[baseIdx+i])
			p.a[i+1] = uint16(maskData[baseIdx+i+1])
			p.a[i+2] = uint16(maskData[baseIdx+i+2])
			p.a[i+3] = uint16(maskData[baseIdx+i+3])
			p.a[i+4] = uint16(maskData[baseIdx+i+4])
			p.a[i+5] = uint16(maskData[baseIdx+i+5])
			p.a[i+6] = uint16(maskData[baseIdx+i+6])
			p.a[i+7] = uint16(maskData[baseIdx+i+7])
		}
	}
}

//go:fix inline
func (p *LowPipeline) MaskU8() {
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		mask := p.a[i]
		p.r[i] = uint16((uint32(p.r[i])*uint32(mask) + 255) >> 8)
		p.g[i] = uint16((uint32(p.g[i])*uint32(mask) + 255) >> 8)
		p.b[i] = uint16((uint32(p.b[i])*uint32(mask) + 255) >> 8)
		p.a[i] = uint16((uint32(p.a[i])*uint32(mask) + 255) >> 8)
	}
}

//go:fix inline
func (p *LowPipeline) ScaleU8() {
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		p.r[i] = uint16((uint32(p.r[i])*uint32(p.dr[i]) + 255) >> 8)
		p.g[i] = uint16((uint32(p.g[i])*uint32(p.dg[i]) + 255) >> 8)
		p.b[i] = uint16((uint32(p.b[i])*uint32(p.db[i]) + 255) >> 8)
	}
}

//go:fix inline
func (p *LowPipeline) LerpU8() {
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		invA0 := 255 - p.a[i]
		p.r[i] = uint16((uint32(p.r[i])*uint32(p.a[i]) + uint32(p.dr[i])*uint32(invA0) + 255) >> 8)
	}
}

//go:fix inline
func (p *LowPipeline) Scale1Float() {
	c := p.ctx.CurrentCoverage
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		p.r[i] = uint16((uint32(p.r[i])*uint32(c) + 255) >> 8)
		p.r[i+1] = uint16((uint32(p.r[i+1])*uint32(c) + 255) >> 8)
		p.r[i+2] = uint16((uint32(p.r[i+2])*uint32(c) + 255) >> 8)
		p.r[i+3] = uint16((uint32(p.r[i+3])*uint32(c) + 255) >> 8)
		p.r[i+4] = uint16((uint32(p.r[i+4])*uint32(c) + 255) >> 8)
		p.r[i+5] = uint16((uint32(p.r[i+5])*uint32(c) + 255) >> 8)
		p.r[i+6] = uint16((uint32(p.r[i+6])*uint32(c) + 255) >> 8)
		p.r[i+7] = uint16((uint32(p.r[i+7])*uint32(c) + 255) >> 8)

		p.g[i] = uint16((uint32(p.g[i])*uint32(c) + 255) >> 8)
		p.g[i+1] = uint16((uint32(p.g[i+1])*uint32(c) + 255) >> 8)
		p.g[i+2] = uint16((uint32(p.g[i+2])*uint32(c) + 255) >> 8)
		p.g[i+3] = uint16((uint32(p.g[i+3])*uint32(c) + 255) >> 8)
		p.g[i+4] = uint16((uint32(p.g[i+4])*uint32(c) + 255) >> 8)
		p.g[i+5] = uint16((uint32(p.g[i+5])*uint32(c) + 255) >> 8)
		p.g[i+6] = uint16((uint32(p.g[i+6])*uint32(c) + 255) >> 8)
		p.g[i+7] = uint16((uint32(p.g[i+7])*uint32(c) + 255) >> 8)

		p.b[i] = uint16((uint32(p.b[i])*uint32(c) + 255) >> 8)
		p.b[i+1] = uint16((uint32(p.b[i+1])*uint32(c) + 255) >> 8)
		p.b[i+2] = uint16((uint32(p.b[i+2])*uint32(c) + 255) >> 8)
		p.b[i+3] = uint16((uint32(p.b[i+3])*uint32(c) + 255) >> 8)
		p.b[i+4] = uint16((uint32(p.b[i+4])*uint32(c) + 255) >> 8)
		p.b[i+5] = uint16((uint32(p.b[i+5])*uint32(c) + 255) >> 8)
		p.b[i+6] = uint16((uint32(p.b[i+6])*uint32(c) + 255) >> 8)
		p.b[i+7] = uint16((uint32(p.b[i+7])*uint32(c) + 255) >> 8)

		p.a[i] = uint16((uint32(p.a[i])*uint32(c) + 255) >> 8)
		p.a[i+1] = uint16((uint32(p.a[i+1])*uint32(c) + 255) >> 8)
		p.a[i+2] = uint16((uint32(p.a[i+2])*uint32(c) + 255) >> 8)
		p.a[i+3] = uint16((uint32(p.a[i+3])*uint32(c) + 255) >> 8)
		p.a[i+4] = uint16((uint32(p.a[i+4])*uint32(c) + 255) >> 8)
		p.a[i+5] = uint16((uint32(p.a[i+5])*uint32(c) + 255) >> 8)
		p.a[i+6] = uint16((uint32(p.a[i+6])*uint32(c) + 255) >> 8)
		p.a[i+7] = uint16((uint32(p.a[i+7])*uint32(c) + 255) >> 8)
	}
}

//go:fix inline
func (p *LowPipeline) Lerp1Float() {
	// where c = current_coverage
	c := p.ctx.CurrentCoverage
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		// lerp(dr, r, c) = dr + (r - dr) * c / 255
		p.r[i] = uint16(int32(p.dr[i]) + ((int32(p.r[i])-int32(p.dr[i]))*int32(c)+128)>>8)
		p.r[i+1] = uint16(int32(p.dr[i+1]) + ((int32(p.r[i+1])-int32(p.dr[i+1]))*int32(c)+128)>>8)
		p.r[i+2] = uint16(int32(p.dr[i+2]) + ((int32(p.r[i+2])-int32(p.dr[i+2]))*int32(c)+128)>>8)
		p.r[i+3] = uint16(int32(p.dr[i+3]) + ((int32(p.r[i+3])-int32(p.dr[i+3]))*int32(c)+128)>>8)
		p.r[i+4] = uint16(int32(p.dr[i+4]) + ((int32(p.r[i+4])-int32(p.dr[i+4]))*int32(c)+128)>>8)
		p.r[i+5] = uint16(int32(p.dr[i+5]) + ((int32(p.r[i+5])-int32(p.dr[i+5]))*int32(c)+128)>>8)
		p.r[i+6] = uint16(int32(p.dr[i+6]) + ((int32(p.r[i+6])-int32(p.dr[i+6]))*int32(c)+128)>>8)
		p.r[i+7] = uint16(int32(p.dr[i+7]) + ((int32(p.r[i+7])-int32(p.dr[i+7]))*int32(c)+128)>>8)

		p.g[i] = uint16(int32(p.dg[i]) + ((int32(p.g[i])-int32(p.dg[i]))*int32(c)+128)>>8)
		p.g[i+1] = uint16(int32(p.dg[i+1]) + ((int32(p.g[i+1])-int32(p.dg[i+1]))*int32(c)+128)>>8)
		p.g[i+2] = uint16(int32(p.dg[i+2]) + ((int32(p.g[i+2])-int32(p.dg[i+2]))*int32(c)+128)>>8)
		p.g[i+3] = uint16(int32(p.dg[i+3]) + ((int32(p.g[i+3])-int32(p.dg[i+3]))*int32(c)+128)>>8)
		p.g[i+4] = uint16(int32(p.dg[i+4]) + ((int32(p.g[i+4])-int32(p.dg[i+4]))*int32(c)+128)>>8)
		p.g[i+5] = uint16(int32(p.dg[i+5]) + ((int32(p.g[i+5])-int32(p.dg[i+5]))*int32(c)+128)>>8)
		p.g[i+6] = uint16(int32(p.dg[i+6]) + ((int32(p.g[i+6])-int32(p.dg[i+6]))*int32(c)+128)>>8)
		p.g[i+7] = uint16(int32(p.dg[i+7]) + ((int32(p.g[i+7])-int32(p.dg[i+7]))*int32(c)+128)>>8)

		p.b[i] = uint16(int32(p.db[i]) + ((int32(p.b[i])-int32(p.db[i]))*int32(c)+128)>>8)
		p.b[i+1] = uint16(int32(p.db[i+1]) + ((int32(p.b[i+1])-int32(p.db[i+1]))*int32(c)+128)>>8)
		p.b[i+2] = uint16(int32(p.db[i+2]) + ((int32(p.b[i+2])-int32(p.db[i+2]))*int32(c)+128)>>8)
		p.b[i+3] = uint16(int32(p.db[i+3]) + ((int32(p.b[i+3])-int32(p.db[i+3]))*int32(c)+128)>>8)
		p.b[i+4] = uint16(int32(p.db[i+4]) + ((int32(p.b[i+4])-int32(p.db[i+4]))*int32(c)+128)>>8)
		p.b[i+5] = uint16(int32(p.db[i+5]) + ((int32(p.b[i+5])-int32(p.db[i+5]))*int32(c)+128)>>8)
		p.b[i+6] = uint16(int32(p.db[i+6]) + ((int32(p.b[i+6])-int32(p.db[i+6]))*int32(c)+128)>>8)
		p.b[i+7] = uint16(int32(p.db[i+7]) + ((int32(p.b[i+7])-int32(p.db[i+7]))*int32(c)+128)>>8)

		p.a[i] = uint16(int32(p.da[i]) + ((int32(p.a[i])-int32(p.da[i]))*int32(c)+128)>>8)
		p.a[i+1] = uint16(int32(p.da[i+1]) + ((int32(p.a[i+1])-int32(p.da[i+1]))*int32(c)+128)>>8)
		p.a[i+2] = uint16(int32(p.da[i+2]) + ((int32(p.a[i+2])-int32(p.da[i+2]))*int32(c)+128)>>8)
		p.a[i+3] = uint16(int32(p.da[i+3]) + ((int32(p.a[i+3])-int32(p.da[i+3]))*int32(c)+128)>>8)
		p.a[i+4] = uint16(int32(p.da[i+4]) + ((int32(p.a[i+4])-int32(p.da[i+4]))*int32(c)+128)>>8)
		p.a[i+5] = uint16(int32(p.da[i+5]) + ((int32(p.a[i+5])-int32(p.da[i+5]))*int32(c)+128)>>8)
		p.a[i+6] = uint16(int32(p.da[i+6]) + ((int32(p.a[i+6])-int32(p.da[i+6]))*int32(c)+128)>>8)
		p.a[i+7] = uint16(int32(p.da[i+7]) + ((int32(p.a[i+7])-int32(p.da[i+7]))*int32(c)+128)>>8)
	}
}

//go:fix inline
func (p *LowPipeline) DestinationAtop() {
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		invDa0, invDa1 := 255-p.da[i], 255-p.da[i+1]
		invDa2, invDa3 := 255-p.da[i+2], 255-p.da[i+3]
		invDa4, invDa5 := 255-p.da[i+4], 255-p.da[i+5]
		invDa6, invDa7 := 255-p.da[i+6], 255-p.da[i+7]

		p.r[i] = uint16((uint32(p.dr[i])*uint32(p.a[i]) + uint32(p.r[i])*uint32(invDa0) + 255) >> 8)
		p.r[i+1] = uint16((uint32(p.dr[i+1])*uint32(p.a[i+1]) + uint32(p.r[i+1])*uint32(invDa1) + 255) >> 8)
		p.r[i+2] = uint16((uint32(p.dr[i+2])*uint32(p.a[i+2]) + uint32(p.r[i+2])*uint32(invDa2) + 255) >> 8)
		p.r[i+3] = uint16((uint32(p.dr[i+3])*uint32(p.a[i+3]) + uint32(p.r[i+3])*uint32(invDa3) + 255) >> 8)
		p.r[i+4] = uint16((uint32(p.dr[i+4])*uint32(p.a[i+4]) + uint32(p.r[i+4])*uint32(invDa4) + 255) >> 8)
		p.r[i+5] = uint16((uint32(p.dr[i+5])*uint32(p.a[i+5]) + uint32(p.r[i+5])*uint32(invDa5) + 255) >> 8)
		p.r[i+6] = uint16((uint32(p.dr[i+6])*uint32(p.a[i+6]) + uint32(p.r[i+6])*uint32(invDa6) + 255) >> 8)
		p.r[i+7] = uint16((uint32(p.dr[i+7])*uint32(p.a[i+7]) + uint32(p.r[i+7])*uint32(invDa7) + 255) >> 8)

		p.g[i] = uint16((uint32(p.dg[i])*uint32(p.a[i]) + uint32(p.g[i])*uint32(invDa0) + 255) >> 8)
		p.g[i+1] = uint16((uint32(p.dg[i+1])*uint32(p.a[i+1]) + uint32(p.g[i+1])*uint32(invDa1) + 255) >> 8)
		p.g[i+2] = uint16((uint32(p.dg[i+2])*uint32(p.a[i+2]) + uint32(p.g[i+2])*uint32(invDa2) + 255) >> 8)
		p.g[i+3] = uint16((uint32(p.dg[i+3])*uint32(p.a[i+3]) + uint32(p.g[i+3])*uint32(invDa3) + 255) >> 8)
		p.g[i+4] = uint16((uint32(p.dg[i+4])*uint32(p.a[i+4]) + uint32(p.g[i+4])*uint32(invDa4) + 255) >> 8)
		p.g[i+5] = uint16((uint32(p.dg[i+5])*uint32(p.a[i+5]) + uint32(p.g[i+5])*uint32(invDa5) + 255) >> 8)
		p.g[i+6] = uint16((uint32(p.dg[i+6])*uint32(p.a[i+6]) + uint32(p.g[i+6])*uint32(invDa6) + 255) >> 8)
		p.g[i+7] = uint16((uint32(p.dg[i+7])*uint32(p.a[i+7]) + uint32(p.g[i+7])*uint32(invDa7) + 255) >> 8)

		p.b[i] = uint16((uint32(p.db[i])*uint32(p.a[i]) + uint32(p.b[i])*uint32(invDa0) + 255) >> 8)
		p.b[i+1] = uint16((uint32(p.db[i+1])*uint32(p.a[i+1]) + uint32(p.b[i+1])*uint32(invDa1) + 255) >> 8)
		p.b[i+2] = uint16((uint32(p.db[i+2])*uint32(p.a[i+2]) + uint32(p.b[i+2])*uint32(invDa2) + 255) >> 8)
		p.b[i+3] = uint16((uint32(p.db[i+3])*uint32(p.a[i+3]) + uint32(p.b[i+3])*uint32(invDa3) + 255) >> 8)
		p.b[i+4] = uint16((uint32(p.db[i+4])*uint32(p.a[i+4]) + uint32(p.b[i+4])*uint32(invDa4) + 255) >> 8)
		p.b[i+5] = uint16((uint32(p.db[i+5])*uint32(p.a[i+5]) + uint32(p.b[i+5])*uint32(invDa5) + 255) >> 8)
		p.b[i+6] = uint16((uint32(p.db[i+6])*uint32(p.a[i+6]) + uint32(p.b[i+6])*uint32(invDa6) + 255) >> 8)
		p.b[i+7] = uint16((uint32(p.db[i+7])*uint32(p.a[i+7]) + uint32(p.b[i+7])*uint32(invDa7) + 255) >> 8)

		// Alpha channel: same formula as RGB - div255(d * sa + s * inv(da))
		p.a[i] = uint16((uint32(p.da[i])*uint32(p.a[i]) + uint32(p.a[i])*uint32(invDa0) + 255) >> 8)
		p.a[i+1] = uint16((uint32(p.da[i+1])*uint32(p.a[i+1]) + uint32(p.a[i+1])*uint32(invDa1) + 255) >> 8)
		p.a[i+2] = uint16((uint32(p.da[i+2])*uint32(p.a[i+2]) + uint32(p.a[i+2])*uint32(invDa2) + 255) >> 8)
		p.a[i+3] = uint16((uint32(p.da[i+3])*uint32(p.a[i+3]) + uint32(p.a[i+3])*uint32(invDa3) + 255) >> 8)
		p.a[i+4] = uint16((uint32(p.da[i+4])*uint32(p.a[i+4]) + uint32(p.a[i+4])*uint32(invDa4) + 255) >> 8)
		p.a[i+5] = uint16((uint32(p.da[i+5])*uint32(p.a[i+5]) + uint32(p.a[i+5])*uint32(invDa5) + 255) >> 8)
		p.a[i+6] = uint16((uint32(p.da[i+6])*uint32(p.a[i+6]) + uint32(p.a[i+6])*uint32(invDa6) + 255) >> 8)
		p.a[i+7] = uint16((uint32(p.da[i+7])*uint32(p.a[i+7]) + uint32(p.a[i+7])*uint32(invDa7) + 255) >> 8)
	}
}

//go:fix inline
func (p *LowPipeline) DestinationIn() {
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		p.r[i] = uint16((uint32(p.dr[i])*uint32(p.a[i]) + 255) >> 8)
		p.r[i+1] = uint16((uint32(p.dr[i+1])*uint32(p.a[i+1]) + 255) >> 8)
		p.r[i+2] = uint16((uint32(p.dr[i+2])*uint32(p.a[i+2]) + 255) >> 8)
		p.r[i+3] = uint16((uint32(p.dr[i+3])*uint32(p.a[i+3]) + 255) >> 8)
		p.r[i+4] = uint16((uint32(p.dr[i+4])*uint32(p.a[i+4]) + 255) >> 8)
		p.r[i+5] = uint16((uint32(p.dr[i+5])*uint32(p.a[i+5]) + 255) >> 8)
		p.r[i+6] = uint16((uint32(p.dr[i+6])*uint32(p.a[i+6]) + 255) >> 8)
		p.r[i+7] = uint16((uint32(p.dr[i+7])*uint32(p.a[i+7]) + 255) >> 8)

		p.g[i] = uint16((uint32(p.dg[i])*uint32(p.a[i]) + 255) >> 8)
		p.g[i+1] = uint16((uint32(p.dg[i+1])*uint32(p.a[i+1]) + 255) >> 8)
		p.g[i+2] = uint16((uint32(p.dg[i+2])*uint32(p.a[i+2]) + 255) >> 8)
		p.g[i+3] = uint16((uint32(p.dg[i+3])*uint32(p.a[i+3]) + 255) >> 8)
		p.g[i+4] = uint16((uint32(p.dg[i+4])*uint32(p.a[i+4]) + 255) >> 8)
		p.g[i+5] = uint16((uint32(p.dg[i+5])*uint32(p.a[i+5]) + 255) >> 8)
		p.g[i+6] = uint16((uint32(p.dg[i+6])*uint32(p.a[i+6]) + 255) >> 8)
		p.g[i+7] = uint16((uint32(p.dg[i+7])*uint32(p.a[i+7]) + 255) >> 8)

		p.b[i] = uint16((uint32(p.db[i])*uint32(p.a[i]) + 255) >> 8)
		p.b[i+1] = uint16((uint32(p.db[i+1])*uint32(p.a[i+1]) + 255) >> 8)
		p.b[i+2] = uint16((uint32(p.db[i+2])*uint32(p.a[i+2]) + 255) >> 8)
		p.b[i+3] = uint16((uint32(p.db[i+3])*uint32(p.a[i+3]) + 255) >> 8)
		p.b[i+4] = uint16((uint32(p.db[i+4])*uint32(p.a[i+4]) + 255) >> 8)
		p.b[i+5] = uint16((uint32(p.db[i+5])*uint32(p.a[i+5]) + 255) >> 8)
		p.b[i+6] = uint16((uint32(p.db[i+6])*uint32(p.a[i+6]) + 255) >> 8)
		p.b[i+7] = uint16((uint32(p.db[i+7])*uint32(p.a[i+7]) + 255) >> 8)

		p.a[i] = uint16((uint32(p.da[i])*uint32(p.a[i]) + 255) >> 8)
		p.a[i+1] = uint16((uint32(p.da[i+1])*uint32(p.a[i+1]) + 255) >> 8)
		p.a[i+2] = uint16((uint32(p.da[i+2])*uint32(p.a[i+2]) + 255) >> 8)
		p.a[i+3] = uint16((uint32(p.da[i+3])*uint32(p.a[i+3]) + 255) >> 8)
		p.a[i+4] = uint16((uint32(p.da[i+4])*uint32(p.a[i+4]) + 255) >> 8)
		p.a[i+5] = uint16((uint32(p.da[i+5])*uint32(p.a[i+5]) + 255) >> 8)
		p.a[i+6] = uint16((uint32(p.da[i+6])*uint32(p.a[i+6]) + 255) >> 8)
		p.a[i+7] = uint16((uint32(p.da[i+7])*uint32(p.a[i+7]) + 255) >> 8)
	}
}

//go:fix inline
func (p *LowPipeline) DestinationOut() {
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		invSa0, invSa1 := 255-p.a[i], 255-p.a[i+1]
		invSa2, invSa3 := 255-p.a[i+2], 255-p.a[i+3]
		invSa4, invSa5 := 255-p.a[i+4], 255-p.a[i+5]
		invSa6, invSa7 := 255-p.a[i+6], 255-p.a[i+7]

		p.r[i] = uint16((uint32(p.dr[i])*uint32(invSa0) + 255) >> 8)
		p.r[i+1] = uint16((uint32(p.dr[i+1])*uint32(invSa1) + 255) >> 8)
		p.r[i+2] = uint16((uint32(p.dr[i+2])*uint32(invSa2) + 255) >> 8)
		p.r[i+3] = uint16((uint32(p.dr[i+3])*uint32(invSa3) + 255) >> 8)
		p.r[i+4] = uint16((uint32(p.dr[i+4])*uint32(invSa4) + 255) >> 8)
		p.r[i+5] = uint16((uint32(p.dr[i+5])*uint32(invSa5) + 255) >> 8)
		p.r[i+6] = uint16((uint32(p.dr[i+6])*uint32(invSa6) + 255) >> 8)
		p.r[i+7] = uint16((uint32(p.dr[i+7])*uint32(invSa7) + 255) >> 8)

		p.g[i] = uint16((uint32(p.dg[i])*uint32(invSa0) + 255) >> 8)
		p.g[i+1] = uint16((uint32(p.dg[i+1])*uint32(invSa1) + 255) >> 8)
		p.g[i+2] = uint16((uint32(p.dg[i+2])*uint32(invSa2) + 255) >> 8)
		p.g[i+3] = uint16((uint32(p.dg[i+3])*uint32(invSa3) + 255) >> 8)
		p.g[i+4] = uint16((uint32(p.dg[i+4])*uint32(invSa4) + 255) >> 8)
		p.g[i+5] = uint16((uint32(p.dg[i+5])*uint32(invSa5) + 255) >> 8)
		p.g[i+6] = uint16((uint32(p.dg[i+6])*uint32(invSa6) + 255) >> 8)
		p.g[i+7] = uint16((uint32(p.dg[i+7])*uint32(invSa7) + 255) >> 8)

		p.b[i] = uint16((uint32(p.db[i])*uint32(invSa0) + 255) >> 8)
		p.b[i+1] = uint16((uint32(p.db[i+1])*uint32(invSa1) + 255) >> 8)
		p.b[i+2] = uint16((uint32(p.db[i+2])*uint32(invSa2) + 255) >> 8)
		p.b[i+3] = uint16((uint32(p.db[i+3])*uint32(invSa3) + 255) >> 8)
		p.b[i+4] = uint16((uint32(p.db[i+4])*uint32(invSa4) + 255) >> 8)
		p.b[i+5] = uint16((uint32(p.db[i+5])*uint32(invSa5) + 255) >> 8)
		p.b[i+6] = uint16((uint32(p.db[i+6])*uint32(invSa6) + 255) >> 8)
		p.b[i+7] = uint16((uint32(p.db[i+7])*uint32(invSa7) + 255) >> 8)

		p.a[i] = uint16((uint32(p.da[i])*uint32(invSa0) + 255) >> 8)
		p.a[i+1] = uint16((uint32(p.da[i+1])*uint32(invSa1) + 255) >> 8)
		p.a[i+2] = uint16((uint32(p.da[i+2])*uint32(invSa2) + 255) >> 8)
		p.a[i+3] = uint16((uint32(p.da[i+3])*uint32(invSa3) + 255) >> 8)
		p.a[i+4] = uint16((uint32(p.da[i+4])*uint32(invSa4) + 255) >> 8)
		p.a[i+5] = uint16((uint32(p.da[i+5])*uint32(invSa5) + 255) >> 8)
		p.a[i+6] = uint16((uint32(p.da[i+6])*uint32(invSa6) + 255) >> 8)
		p.a[i+7] = uint16((uint32(p.da[i+7])*uint32(invSa7) + 255) >> 8)
	}
}

//go:fix inline
func (p *LowPipeline) DestinationOver() {
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		invDa0, invDa1 := 255-p.da[i], 255-p.da[i+1]
		invDa2, invDa3 := 255-p.da[i+2], 255-p.da[i+3]
		invDa4, invDa5 := 255-p.da[i+4], 255-p.da[i+5]
		invDa6, invDa7 := 255-p.da[i+6], 255-p.da[i+7]

		p.r[i] = uint16(uint32(p.dr[i]) + (uint32(p.r[i])*uint32(invDa0)+255)>>8)
		p.r[i+1] = uint16(uint32(p.dr[i+1]) + (uint32(p.r[i+1])*uint32(invDa1)+255)>>8)
		p.r[i+2] = uint16(uint32(p.dr[i+2]) + (uint32(p.r[i+2])*uint32(invDa2)+255)>>8)
		p.r[i+3] = uint16(uint32(p.dr[i+3]) + (uint32(p.r[i+3])*uint32(invDa3)+255)>>8)
		p.r[i+4] = uint16(uint32(p.dr[i+4]) + (uint32(p.r[i+4])*uint32(invDa4)+255)>>8)
		p.r[i+5] = uint16(uint32(p.dr[i+5]) + (uint32(p.r[i+5])*uint32(invDa5)+255)>>8)
		p.r[i+6] = uint16(uint32(p.dr[i+6]) + (uint32(p.r[i+6])*uint32(invDa6)+255)>>8)
		p.r[i+7] = uint16(uint32(p.dr[i+7]) + (uint32(p.r[i+7])*uint32(invDa7)+255)>>8)

		p.g[i] = uint16(uint32(p.dg[i]) + (uint32(p.g[i])*uint32(invDa0)+255)>>8)
		p.g[i+1] = uint16(uint32(p.dg[i+1]) + (uint32(p.g[i+1])*uint32(invDa1)+255)>>8)
		p.g[i+2] = uint16(uint32(p.dg[i+2]) + (uint32(p.g[i+2])*uint32(invDa2)+255)>>8)
		p.g[i+3] = uint16(uint32(p.dg[i+3]) + (uint32(p.g[i+3])*uint32(invDa3)+255)>>8)
		p.g[i+4] = uint16(uint32(p.dg[i+4]) + (uint32(p.g[i+4])*uint32(invDa4)+255)>>8)
		p.g[i+5] = uint16(uint32(p.dg[i+5]) + (uint32(p.g[i+5])*uint32(invDa5)+255)>>8)
		p.g[i+6] = uint16(uint32(p.dg[i+6]) + (uint32(p.g[i+6])*uint32(invDa6)+255)>>8)
		p.g[i+7] = uint16(uint32(p.dg[i+7]) + (uint32(p.g[i+7])*uint32(invDa7)+255)>>8)

		p.b[i] = uint16(uint32(p.db[i]) + (uint32(p.b[i])*uint32(invDa0)+255)>>8)
		p.b[i+1] = uint16(uint32(p.db[i+1]) + (uint32(p.b[i+1])*uint32(invDa1)+255)>>8)
		p.b[i+2] = uint16(uint32(p.db[i+2]) + (uint32(p.b[i+2])*uint32(invDa2)+255)>>8)
		p.b[i+3] = uint16(uint32(p.db[i+3]) + (uint32(p.b[i+3])*uint32(invDa3)+255)>>8)
		p.b[i+4] = uint16(uint32(p.db[i+4]) + (uint32(p.b[i+4])*uint32(invDa4)+255)>>8)
		p.b[i+5] = uint16(uint32(p.db[i+5]) + (uint32(p.b[i+5])*uint32(invDa5)+255)>>8)
		p.b[i+6] = uint16(uint32(p.db[i+6]) + (uint32(p.b[i+6])*uint32(invDa6)+255)>>8)
		p.b[i+7] = uint16(uint32(p.db[i+7]) + (uint32(p.b[i+7])*uint32(invDa7)+255)>>8)

		p.a[i] = uint16(uint32(p.da[i]) + (uint32(p.a[i])*uint32(invDa0)+255)>>8)
		p.a[i+1] = uint16(uint32(p.da[i+1]) + (uint32(p.a[i+1])*uint32(invDa1)+255)>>8)
		p.a[i+2] = uint16(uint32(p.da[i+2]) + (uint32(p.a[i+2])*uint32(invDa2)+255)>>8)
		p.a[i+3] = uint16(uint32(p.da[i+3]) + (uint32(p.a[i+3])*uint32(invDa3)+255)>>8)
		p.a[i+4] = uint16(uint32(p.da[i+4]) + (uint32(p.a[i+4])*uint32(invDa4)+255)>>8)
		p.a[i+5] = uint16(uint32(p.da[i+5]) + (uint32(p.a[i+5])*uint32(invDa5)+255)>>8)
		p.a[i+6] = uint16(uint32(p.da[i+6]) + (uint32(p.a[i+6])*uint32(invDa6)+255)>>8)
		p.a[i+7] = uint16(uint32(p.da[i+7]) + (uint32(p.a[i+7])*uint32(invDa7)+255)>>8)
	}
}

//go:fix inline
func (p *LowPipeline) SourceAtop() {
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		invSa0, invSa1 := 255-p.a[i], 255-p.a[i+1]
		invSa2, invSa3 := 255-p.a[i+2], 255-p.a[i+3]
		invSa4, invSa5 := 255-p.a[i+4], 255-p.a[i+5]
		invSa6, invSa7 := 255-p.a[i+6], 255-p.a[i+7]

		p.r[i] = uint16((uint32(p.r[i])*uint32(p.da[i]) + uint32(p.dr[i])*uint32(invSa0) + 255) >> 8)
		p.r[i+1] = uint16((uint32(p.r[i+1])*uint32(p.da[i+1]) + uint32(p.dr[i+1])*uint32(invSa1) + 255) >> 8)
		p.r[i+2] = uint16((uint32(p.r[i+2])*uint32(p.da[i+2]) + uint32(p.dr[i+2])*uint32(invSa2) + 255) >> 8)
		p.r[i+3] = uint16((uint32(p.r[i+3])*uint32(p.da[i+3]) + uint32(p.dr[i+3])*uint32(invSa3) + 255) >> 8)
		p.r[i+4] = uint16((uint32(p.r[i+4])*uint32(p.da[i+4]) + uint32(p.dr[i+4])*uint32(invSa4) + 255) >> 8)
		p.r[i+5] = uint16((uint32(p.r[i+5])*uint32(p.da[i+5]) + uint32(p.dr[i+5])*uint32(invSa5) + 255) >> 8)
		p.r[i+6] = uint16((uint32(p.r[i+6])*uint32(p.da[i+6]) + uint32(p.dr[i+6])*uint32(invSa6) + 255) >> 8)
		p.r[i+7] = uint16((uint32(p.r[i+7])*uint32(p.da[i+7]) + uint32(p.dr[i+7])*uint32(invSa7) + 255) >> 8)

		p.g[i] = uint16((uint32(p.g[i])*uint32(p.da[i]) + uint32(p.dg[i])*uint32(invSa0) + 255) >> 8)
		p.g[i+1] = uint16((uint32(p.g[i+1])*uint32(p.da[i+1]) + uint32(p.dg[i+1])*uint32(invSa1) + 255) >> 8)
		p.g[i+2] = uint16((uint32(p.g[i+2])*uint32(p.da[i+2]) + uint32(p.dg[i+2])*uint32(invSa2) + 255) >> 8)
		p.g[i+3] = uint16((uint32(p.g[i+3])*uint32(p.da[i+3]) + uint32(p.dg[i+3])*uint32(invSa3) + 255) >> 8)
		p.g[i+4] = uint16((uint32(p.g[i+4])*uint32(p.da[i+4]) + uint32(p.dg[i+4])*uint32(invSa4) + 255) >> 8)
		p.g[i+5] = uint16((uint32(p.g[i+5])*uint32(p.da[i+5]) + uint32(p.dg[i+5])*uint32(invSa5) + 255) >> 8)
		p.g[i+6] = uint16((uint32(p.g[i+6])*uint32(p.da[i+6]) + uint32(p.dg[i+6])*uint32(invSa6) + 255) >> 8)
		p.g[i+7] = uint16((uint32(p.g[i+7])*uint32(p.da[i+7]) + uint32(p.dg[i+7])*uint32(invSa7) + 255) >> 8)

		p.b[i] = uint16((uint32(p.b[i])*uint32(p.da[i]) + uint32(p.db[i])*uint32(invSa0) + 255) >> 8)
		p.b[i+1] = uint16((uint32(p.b[i+1])*uint32(p.da[i+1]) + uint32(p.db[i+1])*uint32(invSa1) + 255) >> 8)
		p.b[i+2] = uint16((uint32(p.b[i+2])*uint32(p.da[i+2]) + uint32(p.db[i+2])*uint32(invSa2) + 255) >> 8)
		p.b[i+3] = uint16((uint32(p.b[i+3])*uint32(p.da[i+3]) + uint32(p.db[i+3])*uint32(invSa3) + 255) >> 8)
		p.b[i+4] = uint16((uint32(p.b[i+4])*uint32(p.da[i+4]) + uint32(p.db[i+4])*uint32(invSa4) + 255) >> 8)
		p.b[i+5] = uint16((uint32(p.b[i+5])*uint32(p.da[i+5]) + uint32(p.db[i+5])*uint32(invSa5) + 255) >> 8)
		p.b[i+6] = uint16((uint32(p.b[i+6])*uint32(p.da[i+6]) + uint32(p.db[i+6])*uint32(invSa6) + 255) >> 8)
		p.b[i+7] = uint16((uint32(p.b[i+7])*uint32(p.da[i+7]) + uint32(p.db[i+7])*uint32(invSa7) + 255) >> 8)

		p.a[i] = uint16((uint32(p.a[i])*uint32(p.da[i]) + uint32(p.da[i])*uint32(invSa0) + 255) >> 8)
		p.a[i+1] = uint16((uint32(p.a[i+1])*uint32(p.da[i+1]) + uint32(p.da[i+1])*uint32(invSa1) + 255) >> 8)
		p.a[i+2] = uint16((uint32(p.a[i+2])*uint32(p.da[i+2]) + uint32(p.da[i+2])*uint32(invSa2) + 255) >> 8)
		p.a[i+3] = uint16((uint32(p.a[i+3])*uint32(p.da[i+3]) + uint32(p.da[i+3])*uint32(invSa3) + 255) >> 8)
		p.a[i+4] = uint16((uint32(p.a[i+4])*uint32(p.da[i+4]) + uint32(p.da[i+4])*uint32(invSa4) + 255) >> 8)
		p.a[i+5] = uint16((uint32(p.a[i+5])*uint32(p.da[i+5]) + uint32(p.da[i+5])*uint32(invSa5) + 255) >> 8)
		p.a[i+6] = uint16((uint32(p.a[i+6])*uint32(p.da[i+6]) + uint32(p.da[i+6])*uint32(invSa6) + 255) >> 8)
		p.a[i+7] = uint16((uint32(p.a[i+7])*uint32(p.da[i+7]) + uint32(p.da[i+7])*uint32(invSa7) + 255) >> 8)
	}
}

//go:fix inline
func (p *LowPipeline) SourceIn() {
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		p.r[i] = uint16((uint32(p.r[i])*uint32(p.da[i]) + 255) >> 8)
		p.r[i+1] = uint16((uint32(p.r[i+1])*uint32(p.da[i+1]) + 255) >> 8)
		p.r[i+2] = uint16((uint32(p.r[i+2])*uint32(p.da[i+2]) + 255) >> 8)
		p.r[i+3] = uint16((uint32(p.r[i+3])*uint32(p.da[i+3]) + 255) >> 8)
		p.r[i+4] = uint16((uint32(p.r[i+4])*uint32(p.da[i+4]) + 255) >> 8)
		p.r[i+5] = uint16((uint32(p.r[i+5])*uint32(p.da[i+5]) + 255) >> 8)
		p.r[i+6] = uint16((uint32(p.r[i+6])*uint32(p.da[i+6]) + 255) >> 8)
		p.r[i+7] = uint16((uint32(p.r[i+7])*uint32(p.da[i+7]) + 255) >> 8)

		p.g[i] = uint16((uint32(p.g[i])*uint32(p.da[i]) + 255) >> 8)
		p.g[i+1] = uint16((uint32(p.g[i+1])*uint32(p.da[i+1]) + 255) >> 8)
		p.g[i+2] = uint16((uint32(p.g[i+2])*uint32(p.da[i+2]) + 255) >> 8)
		p.g[i+3] = uint16((uint32(p.g[i+3])*uint32(p.da[i+3]) + 255) >> 8)
		p.g[i+4] = uint16((uint32(p.g[i+4])*uint32(p.da[i+4]) + 255) >> 8)
		p.g[i+5] = uint16((uint32(p.g[i+5])*uint32(p.da[i+5]) + 255) >> 8)
		p.g[i+6] = uint16((uint32(p.g[i+6])*uint32(p.da[i+6]) + 255) >> 8)
		p.g[i+7] = uint16((uint32(p.g[i+7])*uint32(p.da[i+7]) + 255) >> 8)

		p.b[i] = uint16((uint32(p.b[i])*uint32(p.da[i]) + 255) >> 8)
		p.b[i+1] = uint16((uint32(p.b[i+1])*uint32(p.da[i+1]) + 255) >> 8)
		p.b[i+2] = uint16((uint32(p.b[i+2])*uint32(p.da[i+2]) + 255) >> 8)
		p.b[i+3] = uint16((uint32(p.b[i+3])*uint32(p.da[i+3]) + 255) >> 8)
		p.b[i+4] = uint16((uint32(p.b[i+4])*uint32(p.da[i+4]) + 255) >> 8)
		p.b[i+5] = uint16((uint32(p.b[i+5])*uint32(p.da[i+5]) + 255) >> 8)
		p.b[i+6] = uint16((uint32(p.b[i+6])*uint32(p.da[i+6]) + 255) >> 8)
		p.b[i+7] = uint16((uint32(p.b[i+7])*uint32(p.da[i+7]) + 255) >> 8)

		p.a[i] = uint16((uint32(p.a[i])*uint32(p.da[i]) + 255) >> 8)
		p.a[i+1] = uint16((uint32(p.a[i+1])*uint32(p.da[i+1]) + 255) >> 8)
		p.a[i+2] = uint16((uint32(p.a[i+2])*uint32(p.da[i+2]) + 255) >> 8)
		p.a[i+3] = uint16((uint32(p.a[i+3])*uint32(p.da[i+3]) + 255) >> 8)
		p.a[i+4] = uint16((uint32(p.a[i+4])*uint32(p.da[i+4]) + 255) >> 8)
		p.a[i+5] = uint16((uint32(p.a[i+5])*uint32(p.da[i+5]) + 255) >> 8)
		p.a[i+6] = uint16((uint32(p.a[i+6])*uint32(p.da[i+6]) + 255) >> 8)
		p.a[i+7] = uint16((uint32(p.a[i+7])*uint32(p.da[i+7]) + 255) >> 8)
	}
}

//go:fix inline
func (p *LowPipeline) SourceOut() {
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		invDa0, invDa1 := 255-p.da[i], 255-p.da[i+1]
		invDa2, invDa3 := 255-p.da[i+2], 255-p.da[i+3]
		invDa4, invDa5 := 255-p.da[i+4], 255-p.da[i+5]
		invDa6, invDa7 := 255-p.da[i+6], 255-p.da[i+7]

		p.r[i] = uint16((uint32(p.r[i])*uint32(invDa0) + 255) >> 8)
		p.r[i+1] = uint16((uint32(p.r[i+1])*uint32(invDa1) + 255) >> 8)
		p.r[i+2] = uint16((uint32(p.r[i+2])*uint32(invDa2) + 255) >> 8)
		p.r[i+3] = uint16((uint32(p.r[i+3])*uint32(invDa3) + 255) >> 8)
		p.r[i+4] = uint16((uint32(p.r[i+4])*uint32(invDa4) + 255) >> 8)
		p.r[i+5] = uint16((uint32(p.r[i+5])*uint32(invDa5) + 255) >> 8)
		p.r[i+6] = uint16((uint32(p.r[i+6])*uint32(invDa6) + 255) >> 8)
		p.r[i+7] = uint16((uint32(p.r[i+7])*uint32(invDa7) + 255) >> 8)

		p.g[i] = uint16((uint32(p.g[i])*uint32(invDa0) + 255) >> 8)
		p.g[i+1] = uint16((uint32(p.g[i+1])*uint32(invDa1) + 255) >> 8)
		p.g[i+2] = uint16((uint32(p.g[i+2])*uint32(invDa2) + 255) >> 8)
		p.g[i+3] = uint16((uint32(p.g[i+3])*uint32(invDa3) + 255) >> 8)
		p.g[i+4] = uint16((uint32(p.g[i+4])*uint32(invDa4) + 255) >> 8)
		p.g[i+5] = uint16((uint32(p.g[i+5])*uint32(invDa5) + 255) >> 8)
		p.g[i+6] = uint16((uint32(p.g[i+6])*uint32(invDa6) + 255) >> 8)
		p.g[i+7] = uint16((uint32(p.g[i+7])*uint32(invDa7) + 255) >> 8)

		p.b[i] = uint16((uint32(p.b[i])*uint32(invDa0) + 255) >> 8)
		p.b[i+1] = uint16((uint32(p.b[i+1])*uint32(invDa1) + 255) >> 8)
		p.b[i+2] = uint16((uint32(p.b[i+2])*uint32(invDa2) + 255) >> 8)
		p.b[i+3] = uint16((uint32(p.b[i+3])*uint32(invDa3) + 255) >> 8)
		p.b[i+4] = uint16((uint32(p.b[i+4])*uint32(invDa4) + 255) >> 8)
		p.b[i+5] = uint16((uint32(p.b[i+5])*uint32(invDa5) + 255) >> 8)
		p.b[i+6] = uint16((uint32(p.b[i+6])*uint32(invDa6) + 255) >> 8)
		p.b[i+7] = uint16((uint32(p.b[i+7])*uint32(invDa7) + 255) >> 8)

		p.a[i] = uint16((uint32(p.a[i])*uint32(invDa0) + 255) >> 8)
		p.a[i+1] = uint16((uint32(p.a[i+1])*uint32(invDa1) + 255) >> 8)
		p.a[i+2] = uint16((uint32(p.a[i+2])*uint32(invDa2) + 255) >> 8)
		p.a[i+3] = uint16((uint32(p.a[i+3])*uint32(invDa3) + 255) >> 8)
		p.a[i+4] = uint16((uint32(p.a[i+4])*uint32(invDa4) + 255) >> 8)
		p.a[i+5] = uint16((uint32(p.a[i+5])*uint32(invDa5) + 255) >> 8)
		p.a[i+6] = uint16((uint32(p.a[i+6])*uint32(invDa6) + 255) >> 8)
		p.a[i+7] = uint16((uint32(p.a[i+7])*uint32(invDa7) + 255) >> 8)
	}
}

//go:fix inline
func (p *LowPipeline) SourceOver() {
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		invSa0, invSa1, invSa2, invSa3 := 255-p.a[i], 255-p.a[i+1], 255-p.a[i+2], 255-p.a[i+3]
		invSa4, invSa5, invSa6, invSa7 := 255-p.a[i+4], 255-p.a[i+5], 255-p.a[i+6], 255-p.a[i+7]

		p.r[i] = uint16(uint32(p.r[i]) + (uint32(p.dr[i])*uint32(invSa0)+255)>>8)
		p.r[i+1] = uint16(uint32(p.r[i+1]) + (uint32(p.dr[i+1])*uint32(invSa1)+255)>>8)
		p.r[i+2] = uint16(uint32(p.r[i+2]) + (uint32(p.dr[i+2])*uint32(invSa2)+255)>>8)
		p.r[i+3] = uint16(uint32(p.r[i+3]) + (uint32(p.dr[i+3])*uint32(invSa3)+255)>>8)
		p.r[i+4] = uint16(uint32(p.r[i+4]) + (uint32(p.dr[i+4])*uint32(invSa4)+255)>>8)
		p.r[i+5] = uint16(uint32(p.r[i+5]) + (uint32(p.dr[i+5])*uint32(invSa5)+255)>>8)
		p.r[i+6] = uint16(uint32(p.r[i+6]) + (uint32(p.dr[i+6])*uint32(invSa6)+255)>>8)
		p.r[i+7] = uint16(uint32(p.r[i+7]) + (uint32(p.dr[i+7])*uint32(invSa7)+255)>>8)

		p.g[i] = uint16(uint32(p.g[i]) + (uint32(p.dg[i])*uint32(invSa0)+255)>>8)
		p.g[i+1] = uint16(uint32(p.g[i+1]) + (uint32(p.dg[i+1])*uint32(invSa1)+255)>>8)
		p.g[i+2] = uint16(uint32(p.g[i+2]) + (uint32(p.dg[i+2])*uint32(invSa2)+255)>>8)
		p.g[i+3] = uint16(uint32(p.g[i+3]) + (uint32(p.dg[i+3])*uint32(invSa3)+255)>>8)
		p.g[i+4] = uint16(uint32(p.g[i+4]) + (uint32(p.dg[i+4])*uint32(invSa4)+255)>>8)
		p.g[i+5] = uint16(uint32(p.g[i+5]) + (uint32(p.dg[i+5])*uint32(invSa5)+255)>>8)
		p.g[i+6] = uint16(uint32(p.g[i+6]) + (uint32(p.dg[i+6])*uint32(invSa6)+255)>>8)
		p.g[i+7] = uint16(uint32(p.g[i+7]) + (uint32(p.dg[i+7])*uint32(invSa7)+255)>>8)

		p.b[i] = uint16(uint32(p.b[i]) + (uint32(p.db[i])*uint32(invSa0)+255)>>8)
		p.b[i+1] = uint16(uint32(p.b[i+1]) + (uint32(p.db[i+1])*uint32(invSa1)+255)>>8)
		p.b[i+2] = uint16(uint32(p.b[i+2]) + (uint32(p.db[i+2])*uint32(invSa2)+255)>>8)
		p.b[i+3] = uint16(uint32(p.b[i+3]) + (uint32(p.db[i+3])*uint32(invSa3)+255)>>8)
		p.b[i+4] = uint16(uint32(p.b[i+4]) + (uint32(p.db[i+4])*uint32(invSa4)+255)>>8)
		p.b[i+5] = uint16(uint32(p.b[i+5]) + (uint32(p.db[i+5])*uint32(invSa5)+255)>>8)
		p.b[i+6] = uint16(uint32(p.b[i+6]) + (uint32(p.db[i+6])*uint32(invSa6)+255)>>8)
		p.b[i+7] = uint16(uint32(p.b[i+7]) + (uint32(p.db[i+7])*uint32(invSa7)+255)>>8)

		p.a[i] = uint16(uint32(p.a[i]) + (uint32(p.da[i])*uint32(invSa0)+255)>>8)
		p.a[i+1] = uint16(uint32(p.a[i+1]) + (uint32(p.da[i+1])*uint32(invSa1)+255)>>8)
		p.a[i+2] = uint16(uint32(p.a[i+2]) + (uint32(p.da[i+2])*uint32(invSa2)+255)>>8)
		p.a[i+3] = uint16(uint32(p.a[i+3]) + (uint32(p.da[i+3])*uint32(invSa3)+255)>>8)
		p.a[i+4] = uint16(uint32(p.a[i+4]) + (uint32(p.da[i+4])*uint32(invSa4)+255)>>8)
		p.a[i+5] = uint16(uint32(p.a[i+5]) + (uint32(p.da[i+5])*uint32(invSa5)+255)>>8)
		p.a[i+6] = uint16(uint32(p.a[i+6]) + (uint32(p.da[i+6])*uint32(invSa6)+255)>>8)
		p.a[i+7] = uint16(uint32(p.a[i+7]) + (uint32(p.da[i+7])*uint32(invSa7)+255)>>8)
	}
}

//go:fix inline
func (p *LowPipeline) Clear() {
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		p.r[i], p.r[i+1], p.r[i+2], p.r[i+3] = 0, 0, 0, 0
		p.r[i+4], p.r[i+5], p.r[i+6], p.r[i+7] = 0, 0, 0, 0

		p.g[i], p.g[i+1], p.g[i+2], p.g[i+3] = 0, 0, 0, 0
		p.g[i+4], p.g[i+5], p.g[i+6], p.g[i+7] = 0, 0, 0, 0

		p.b[i], p.b[i+1], p.b[i+2], p.b[i+3] = 0, 0, 0, 0
		p.b[i+4], p.b[i+5], p.b[i+6], p.b[i+7] = 0, 0, 0, 0

		p.a[i], p.a[i+1], p.a[i+2], p.a[i+3] = 0, 0, 0, 0
		p.a[i+4], p.a[i+5], p.a[i+6], p.a[i+7] = 0, 0, 0, 0
	}
}

//go:fix inline
func (p *LowPipeline) Modulate() {
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		p.r[i] = uint16((uint32(p.r[i])*uint32(p.dr[i]) + 255) >> 8)
		p.r[i+1] = uint16((uint32(p.r[i+1])*uint32(p.dr[i+1]) + 255) >> 8)
		p.r[i+2] = uint16((uint32(p.r[i+2])*uint32(p.dr[i+2]) + 255) >> 8)
		p.r[i+3] = uint16((uint32(p.r[i+3])*uint32(p.dr[i+3]) + 255) >> 8)
		p.r[i+4] = uint16((uint32(p.r[i+4])*uint32(p.dr[i+4]) + 255) >> 8)
		p.r[i+5] = uint16((uint32(p.r[i+5])*uint32(p.dr[i+5]) + 255) >> 8)
		p.r[i+6] = uint16((uint32(p.r[i+6])*uint32(p.dr[i+6]) + 255) >> 8)
		p.r[i+7] = uint16((uint32(p.r[i+7])*uint32(p.dr[i+7]) + 255) >> 8)

		p.g[i] = uint16((uint32(p.g[i])*uint32(p.dg[i]) + 255) >> 8)
		p.g[i+1] = uint16((uint32(p.g[i+1])*uint32(p.dg[i+1]) + 255) >> 8)
		p.g[i+2] = uint16((uint32(p.g[i+2])*uint32(p.dg[i+2]) + 255) >> 8)
		p.g[i+3] = uint16((uint32(p.g[i+3])*uint32(p.dg[i+3]) + 255) >> 8)
		p.g[i+4] = uint16((uint32(p.g[i+4])*uint32(p.dg[i+4]) + 255) >> 8)
		p.g[i+5] = uint16((uint32(p.g[i+5])*uint32(p.dg[i+5]) + 255) >> 8)
		p.g[i+6] = uint16((uint32(p.g[i+6])*uint32(p.dg[i+6]) + 255) >> 8)
		p.g[i+7] = uint16((uint32(p.g[i+7])*uint32(p.dg[i+7]) + 255) >> 8)

		p.b[i] = uint16((uint32(p.b[i])*uint32(p.db[i]) + 255) >> 8)
		p.b[i+1] = uint16((uint32(p.b[i+1])*uint32(p.db[i+1]) + 255) >> 8)
		p.b[i+2] = uint16((uint32(p.b[i+2])*uint32(p.db[i+2]) + 255) >> 8)
		p.b[i+3] = uint16((uint32(p.b[i+3])*uint32(p.db[i+3]) + 255) >> 8)
		p.b[i+4] = uint16((uint32(p.b[i+4])*uint32(p.db[i+4]) + 255) >> 8)
		p.b[i+5] = uint16((uint32(p.b[i+5])*uint32(p.db[i+5]) + 255) >> 8)
		p.b[i+6] = uint16((uint32(p.b[i+6])*uint32(p.db[i+6]) + 255) >> 8)
		p.b[i+7] = uint16((uint32(p.b[i+7])*uint32(p.db[i+7]) + 255) >> 8)

		p.a[i] = uint16((uint32(p.a[i])*uint32(p.da[i]) + 255) >> 8)
		p.a[i+1] = uint16((uint32(p.a[i+1])*uint32(p.da[i+1]) + 255) >> 8)
		p.a[i+2] = uint16((uint32(p.a[i+2])*uint32(p.da[i+2]) + 255) >> 8)
		p.a[i+3] = uint16((uint32(p.a[i+3])*uint32(p.da[i+3]) + 255) >> 8)
		p.a[i+4] = uint16((uint32(p.a[i+4])*uint32(p.da[i+4]) + 255) >> 8)
		p.a[i+5] = uint16((uint32(p.a[i+5])*uint32(p.da[i+5]) + 255) >> 8)
		p.a[i+6] = uint16((uint32(p.a[i+6])*uint32(p.da[i+6]) + 255) >> 8)
		p.a[i+7] = uint16((uint32(p.a[i+7])*uint32(p.da[i+7]) + 255) >> 8)
	}
}

//go:fix inline
func (p *LowPipeline) Multiply() {
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		invDa0, invDa1, invDa2, invDa3 := 255-p.da[i], 255-p.da[i+1], 255-p.da[i+2], 255-p.da[i+3]
		invDa4, invDa5, invDa6, invDa7 := 255-p.da[i+4], 255-p.da[i+5], 255-p.da[i+6], 255-p.da[i+7]
		invSa0, invSa1, invSa2, invSa3 := 255-p.a[i], 255-p.a[i+1], 255-p.a[i+2], 255-p.a[i+3]
		invSa4, invSa5, invSa6, invSa7 := 255-p.a[i+4], 255-p.a[i+5], 255-p.a[i+6], 255-p.a[i+7]

		p.r[i] = uint16((uint32(p.r[i])*uint32(invDa0) + uint32(p.dr[i])*uint32(invSa0) + uint32(p.r[i])*uint32(p.dr[i]) + 255) >> 8)
		p.r[i+1] = uint16((uint32(p.r[i+1])*uint32(invDa1) + uint32(p.dr[i+1])*uint32(invSa1) + uint32(p.r[i+1])*uint32(p.dr[i+1]) + 255) >> 8)
		p.r[i+2] = uint16((uint32(p.r[i+2])*uint32(invDa2) + uint32(p.dr[i+2])*uint32(invSa2) + uint32(p.r[i+2])*uint32(p.dr[i+2]) + 255) >> 8)
		p.r[i+3] = uint16((uint32(p.r[i+3])*uint32(invDa3) + uint32(p.dr[i+3])*uint32(invSa3) + uint32(p.r[i+3])*uint32(p.dr[i+3]) + 255) >> 8)
		p.r[i+4] = uint16((uint32(p.r[i+4])*uint32(invDa4) + uint32(p.dr[i+4])*uint32(invSa4) + uint32(p.r[i+4])*uint32(p.dr[i+4]) + 255) >> 8)
		p.r[i+5] = uint16((uint32(p.r[i+5])*uint32(invDa5) + uint32(p.dr[i+5])*uint32(invSa5) + uint32(p.r[i+5])*uint32(p.dr[i+5]) + 255) >> 8)
		p.r[i+6] = uint16((uint32(p.r[i+6])*uint32(invDa6) + uint32(p.dr[i+6])*uint32(invSa6) + uint32(p.r[i+6])*uint32(p.dr[i+6]) + 255) >> 8)
		p.r[i+7] = uint16((uint32(p.r[i+7])*uint32(invDa7) + uint32(p.dr[i+7])*uint32(invSa7) + uint32(p.r[i+7])*uint32(p.dr[i+7]) + 255) >> 8)

		p.g[i] = uint16((uint32(p.g[i])*uint32(invDa0) + uint32(p.dg[i])*uint32(invSa0) + uint32(p.g[i])*uint32(p.dg[i]) + 255) >> 8)
		p.g[i+1] = uint16((uint32(p.g[i+1])*uint32(invDa1) + uint32(p.dg[i+1])*uint32(invSa1) + uint32(p.g[i+1])*uint32(p.dg[i+1]) + 255) >> 8)
		p.g[i+2] = uint16((uint32(p.g[i+2])*uint32(invDa2) + uint32(p.dg[i+2])*uint32(invSa2) + uint32(p.g[i+2])*uint32(p.dg[i+2]) + 255) >> 8)
		p.g[i+3] = uint16((uint32(p.g[i+3])*uint32(invDa3) + uint32(p.dg[i+3])*uint32(invSa3) + uint32(p.g[i+3])*uint32(p.dg[i+3]) + 255) >> 8)
		p.g[i+4] = uint16((uint32(p.g[i+4])*uint32(invDa4) + uint32(p.dg[i+4])*uint32(invSa4) + uint32(p.g[i+4])*uint32(p.dg[i+4]) + 255) >> 8)
		p.g[i+5] = uint16((uint32(p.g[i+5])*uint32(invDa5) + uint32(p.dg[i+5])*uint32(invSa5) + uint32(p.g[i+5])*uint32(p.dg[i+5]) + 255) >> 8)
		p.g[i+6] = uint16((uint32(p.g[i+6])*uint32(invDa6) + uint32(p.dg[i+6])*uint32(invSa6) + uint32(p.g[i+6])*uint32(p.dg[i+6]) + 255) >> 8)
		p.g[i+7] = uint16((uint32(p.g[i+7])*uint32(invDa7) + uint32(p.dg[i+7])*uint32(invSa7) + uint32(p.g[i+7])*uint32(p.dg[i+7]) + 255) >> 8)

		p.b[i] = uint16((uint32(p.b[i])*uint32(invDa0) + uint32(p.db[i])*uint32(invSa0) + uint32(p.b[i])*uint32(p.db[i]) + 255) >> 8)
		p.b[i+1] = uint16((uint32(p.b[i+1])*uint32(invDa1) + uint32(p.db[i+1])*uint32(invSa1) + uint32(p.b[i+1])*uint32(p.db[i+1]) + 255) >> 8)
		p.b[i+2] = uint16((uint32(p.b[i+2])*uint32(invDa2) + uint32(p.db[i+2])*uint32(invSa2) + uint32(p.b[i+2])*uint32(p.db[i+2]) + 255) >> 8)
		p.b[i+3] = uint16((uint32(p.b[i+3])*uint32(invDa3) + uint32(p.db[i+3])*uint32(invSa3) + uint32(p.b[i+3])*uint32(p.db[i+3]) + 255) >> 8)
		p.b[i+4] = uint16((uint32(p.b[i+4])*uint32(invDa4) + uint32(p.db[i+4])*uint32(invSa4) + uint32(p.b[i+4])*uint32(p.db[i+4]) + 255) >> 8)
		p.b[i+5] = uint16((uint32(p.b[i+5])*uint32(invDa5) + uint32(p.db[i+5])*uint32(invSa5) + uint32(p.b[i+5])*uint32(p.db[i+5]) + 255) >> 8)
		p.b[i+6] = uint16((uint32(p.b[i+6])*uint32(invDa6) + uint32(p.db[i+6])*uint32(invSa6) + uint32(p.b[i+6])*uint32(p.db[i+6]) + 255) >> 8)
		p.b[i+7] = uint16((uint32(p.b[i+7])*uint32(invDa7) + uint32(p.db[i+7])*uint32(invSa7) + uint32(p.b[i+7])*uint32(p.db[i+7]) + 255) >> 8)

		p.a[i] = uint16((uint32(p.a[i])*uint32(invDa0) + uint32(p.da[i])*uint32(invSa0) + uint32(p.a[i])*uint32(p.da[i]) + 255) >> 8)
		p.a[i+1] = uint16((uint32(p.a[i+1])*uint32(invDa1) + uint32(p.da[i+1])*uint32(invSa1) + uint32(p.a[i+1])*uint32(p.da[i+1]) + 255) >> 8)
		p.a[i+2] = uint16((uint32(p.a[i+2])*uint32(invDa2) + uint32(p.da[i+2])*uint32(invSa2) + uint32(p.a[i+2])*uint32(p.da[i+2]) + 255) >> 8)
		p.a[i+3] = uint16((uint32(p.a[i+3])*uint32(invDa3) + uint32(p.da[i+3])*uint32(invSa3) + uint32(p.a[i+3])*uint32(p.da[i+3]) + 255) >> 8)
		p.a[i+4] = uint16((uint32(p.a[i+4])*uint32(invDa4) + uint32(p.da[i+4])*uint32(invSa4) + uint32(p.a[i+4])*uint32(p.da[i+4]) + 255) >> 8)
		p.a[i+5] = uint16((uint32(p.a[i+5])*uint32(invDa5) + uint32(p.da[i+5])*uint32(invSa5) + uint32(p.a[i+5])*uint32(p.da[i+5]) + 255) >> 8)
		p.a[i+6] = uint16((uint32(p.a[i+6])*uint32(invDa6) + uint32(p.da[i+6])*uint32(invSa6) + uint32(p.a[i+6])*uint32(p.da[i+6]) + 255) >> 8)
		p.a[i+7] = uint16((uint32(p.a[i+7])*uint32(invDa7) + uint32(p.da[i+7])*uint32(invSa7) + uint32(p.a[i+7])*uint32(p.da[i+7]) + 255) >> 8)
	}
}

//go:fix inline
func (p *LowPipeline) Plus() {
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		sum := uint32(p.r[i]) + uint32(p.dr[i])
		if sum > 255 {
			p.r[i] = 255
		} else {
			p.r[i] = uint16(sum)
		}

		sum = uint32(p.r[i+1]) + uint32(p.dr[i+1])
		if sum > 255 {
			p.r[i+1] = 255
		} else {
			p.r[i+1] = uint16(sum)
		}

		sum = uint32(p.r[i+2]) + uint32(p.dr[i+2])
		if sum > 255 {
			p.r[i+2] = 255
		} else {
			p.r[i+2] = uint16(sum)
		}

		sum = uint32(p.r[i+3]) + uint32(p.dr[i+3])
		if sum > 255 {
			p.r[i+3] = 255
		} else {
			p.r[i+3] = uint16(sum)
		}

		sum = uint32(p.r[i+4]) + uint32(p.dr[i+4])
		if sum > 255 {
			p.r[i+4] = 255
		} else {
			p.r[i+4] = uint16(sum)
		}

		sum = uint32(p.r[i+5]) + uint32(p.dr[i+5])
		if sum > 255 {
			p.r[i+5] = 255
		} else {
			p.r[i+5] = uint16(sum)
		}

		sum = uint32(p.r[i+6]) + uint32(p.dr[i+6])
		if sum > 255 {
			p.r[i+6] = 255
		} else {
			p.r[i+6] = uint16(sum)
		}

		sum = uint32(p.r[i+7]) + uint32(p.dr[i+7])
		if sum > 255 {
			p.r[i+7] = 255
		} else {
			p.r[i+7] = uint16(sum)
		}

		// G channel
		sum = uint32(p.g[i]) + uint32(p.dg[i])
		if sum > 255 {
			p.g[i] = 255
		} else {
			p.g[i] = uint16(sum)
		}
		sum = uint32(p.g[i+1]) + uint32(p.dg[i+1])
		if sum > 255 {
			p.g[i+1] = 255
		} else {
			p.g[i+1] = uint16(sum)
		}
		sum = uint32(p.g[i+2]) + uint32(p.dg[i+2])
		if sum > 255 {
			p.g[i+2] = 255
		} else {
			p.g[i+2] = uint16(sum)
		}
		sum = uint32(p.g[i+3]) + uint32(p.dg[i+3])
		if sum > 255 {
			p.g[i+3] = 255
		} else {
			p.g[i+3] = uint16(sum)
		}
		sum = uint32(p.g[i+4]) + uint32(p.dg[i+4])
		if sum > 255 {
			p.g[i+4] = 255
		} else {
			p.g[i+4] = uint16(sum)
		}
		sum = uint32(p.g[i+5]) + uint32(p.dg[i+5])
		if sum > 255 {
			p.g[i+5] = 255
		} else {
			p.g[i+5] = uint16(sum)
		}
		sum = uint32(p.g[i+6]) + uint32(p.dg[i+6])
		if sum > 255 {
			p.g[i+6] = 255
		} else {
			p.g[i+6] = uint16(sum)
		}
		sum = uint32(p.g[i+7]) + uint32(p.dg[i+7])
		if sum > 255 {
			p.g[i+7] = 255
		} else {
			p.g[i+7] = uint16(sum)
		}

		// B channel
		sum = uint32(p.b[i]) + uint32(p.db[i])
		if sum > 255 {
			p.b[i] = 255
		} else {
			p.b[i] = uint16(sum)
		}
		sum = uint32(p.b[i+1]) + uint32(p.db[i+1])
		if sum > 255 {
			p.b[i+1] = 255
		} else {
			p.b[i+1] = uint16(sum)
		}
		sum = uint32(p.b[i+2]) + uint32(p.db[i+2])
		if sum > 255 {
			p.b[i+2] = 255
		} else {
			p.b[i+2] = uint16(sum)
		}
		sum = uint32(p.b[i+3]) + uint32(p.db[i+3])
		if sum > 255 {
			p.b[i+3] = 255
		} else {
			p.b[i+3] = uint16(sum)
		}
		sum = uint32(p.b[i+4]) + uint32(p.db[i+4])
		if sum > 255 {
			p.b[i+4] = 255
		} else {
			p.b[i+4] = uint16(sum)
		}
		sum = uint32(p.b[i+5]) + uint32(p.db[i+5])
		if sum > 255 {
			p.b[i+5] = 255
		} else {
			p.b[i+5] = uint16(sum)
		}
		sum = uint32(p.b[i+6]) + uint32(p.db[i+6])
		if sum > 255 {
			p.b[i+6] = 255
		} else {
			p.b[i+6] = uint16(sum)
		}
		sum = uint32(p.b[i+7]) + uint32(p.db[i+7])
		if sum > 255 {
			p.b[i+7] = 255
		} else {
			p.b[i+7] = uint16(sum)
		}

		// A channel
		sum = uint32(p.a[i]) + uint32(p.da[i])
		if sum > 255 {
			p.a[i] = 255
		} else {
			p.a[i] = uint16(sum)
		}
		sum = uint32(p.a[i+1]) + uint32(p.da[i+1])
		if sum > 255 {
			p.a[i+1] = 255
		} else {
			p.a[i+1] = uint16(sum)
		}
		sum = uint32(p.a[i+2]) + uint32(p.da[i+2])
		if sum > 255 {
			p.a[i+2] = 255
		} else {
			p.a[i+2] = uint16(sum)
		}
		sum = uint32(p.a[i+3]) + uint32(p.da[i+3])
		if sum > 255 {
			p.a[i+3] = 255
		} else {
			p.a[i+3] = uint16(sum)
		}
		sum = uint32(p.a[i+4]) + uint32(p.da[i+4])
		if sum > 255 {
			p.a[i+4] = 255
		} else {
			p.a[i+4] = uint16(sum)
		}
		sum = uint32(p.a[i+5]) + uint32(p.da[i+5])
		if sum > 255 {
			p.a[i+5] = 255
		} else {
			p.a[i+5] = uint16(sum)
		}
		sum = uint32(p.a[i+6]) + uint32(p.da[i+6])
		if sum > 255 {
			p.a[i+6] = 255
		} else {
			p.a[i+6] = uint16(sum)
		}
		sum = uint32(p.a[i+7]) + uint32(p.da[i+7])
		if sum > 255 {
			p.a[i+7] = 255
		} else {
			p.a[i+7] = uint16(sum)
		}
	}

}

//go:fix inline
func (p *LowPipeline) Screen() {
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		p.r[i] = uint16(uint32(p.r[i]) + uint32(p.dr[i]) - (uint32(p.r[i])*uint32(p.dr[i])+255)>>8)
		p.r[i+1] = uint16(uint32(p.r[i+1]) + uint32(p.dr[i+1]) - (uint32(p.r[i+1])*uint32(p.dr[i+1])+255)>>8)
		p.r[i+2] = uint16(uint32(p.r[i+2]) + uint32(p.dr[i+2]) - (uint32(p.r[i+2])*uint32(p.dr[i+2])+255)>>8)
		p.r[i+3] = uint16(uint32(p.r[i+3]) + uint32(p.dr[i+3]) - (uint32(p.r[i+3])*uint32(p.dr[i+3])+255)>>8)
		p.r[i+4] = uint16(uint32(p.r[i+4]) + uint32(p.dr[i+4]) - (uint32(p.r[i+4])*uint32(p.dr[i+4])+255)>>8)
		p.r[i+5] = uint16(uint32(p.r[i+5]) + uint32(p.dr[i+5]) - (uint32(p.r[i+5])*uint32(p.dr[i+5])+255)>>8)
		p.r[i+6] = uint16(uint32(p.r[i+6]) + uint32(p.dr[i+6]) - (uint32(p.r[i+6])*uint32(p.dr[i+6])+255)>>8)
		p.r[i+7] = uint16(uint32(p.r[i+7]) + uint32(p.dr[i+7]) - (uint32(p.r[i+7])*uint32(p.dr[i+7])+255)>>8)

		p.g[i] = uint16(uint32(p.g[i]) + uint32(p.dg[i]) - (uint32(p.g[i])*uint32(p.dg[i])+255)>>8)
		p.g[i+1] = uint16(uint32(p.g[i+1]) + uint32(p.dg[i+1]) - (uint32(p.g[i+1])*uint32(p.dg[i+1])+255)>>8)
		p.g[i+2] = uint16(uint32(p.g[i+2]) + uint32(p.dg[i+2]) - (uint32(p.g[i+2])*uint32(p.dg[i+2])+255)>>8)
		p.g[i+3] = uint16(uint32(p.g[i+3]) + uint32(p.dg[i+3]) - (uint32(p.g[i+3])*uint32(p.dg[i+3])+255)>>8)
		p.g[i+4] = uint16(uint32(p.g[i+4]) + uint32(p.dg[i+4]) - (uint32(p.g[i+4])*uint32(p.dg[i+4])+255)>>8)
		p.g[i+5] = uint16(uint32(p.g[i+5]) + uint32(p.dg[i+5]) - (uint32(p.g[i+5])*uint32(p.dg[i+5])+255)>>8)
		p.g[i+6] = uint16(uint32(p.g[i+6]) + uint32(p.dg[i+6]) - (uint32(p.g[i+6])*uint32(p.dg[i+6])+255)>>8)
		p.g[i+7] = uint16(uint32(p.g[i+7]) + uint32(p.dg[i+7]) - (uint32(p.g[i+7])*uint32(p.dg[i+7])+255)>>8)

		p.b[i] = uint16(uint32(p.b[i]) + uint32(p.db[i]) - (uint32(p.b[i])*uint32(p.db[i])+255)>>8)
		p.b[i+1] = uint16(uint32(p.b[i+1]) + uint32(p.db[i+1]) - (uint32(p.b[i+1])*uint32(p.db[i+1])+255)>>8)
		p.b[i+2] = uint16(uint32(p.b[i+2]) + uint32(p.db[i+2]) - (uint32(p.b[i+2])*uint32(p.db[i+2])+255)>>8)
		p.b[i+3] = uint16(uint32(p.b[i+3]) + uint32(p.db[i+3]) - (uint32(p.b[i+3])*uint32(p.db[i+3])+255)>>8)
		p.b[i+4] = uint16(uint32(p.b[i+4]) + uint32(p.db[i+4]) - (uint32(p.b[i+4])*uint32(p.db[i+4])+255)>>8)
		p.b[i+5] = uint16(uint32(p.b[i+5]) + uint32(p.db[i+5]) - (uint32(p.b[i+5])*uint32(p.db[i+5])+255)>>8)
		p.b[i+6] = uint16(uint32(p.b[i+6]) + uint32(p.db[i+6]) - (uint32(p.b[i+6])*uint32(p.db[i+6])+255)>>8)
		p.b[i+7] = uint16(uint32(p.b[i+7]) + uint32(p.db[i+7]) - (uint32(p.b[i+7])*uint32(p.db[i+7])+255)>>8)

		p.a[i] = uint16(uint32(p.a[i]) + uint32(p.da[i]) - (uint32(p.a[i])*uint32(p.da[i])+255)>>8)
		p.a[i+1] = uint16(uint32(p.a[i+1]) + uint32(p.da[i+1]) - (uint32(p.a[i+1])*uint32(p.da[i+1])+255)>>8)
		p.a[i+2] = uint16(uint32(p.a[i+2]) + uint32(p.da[i+2]) - (uint32(p.a[i+2])*uint32(p.da[i+2])+255)>>8)
		p.a[i+3] = uint16(uint32(p.a[i+3]) + uint32(p.da[i+3]) - (uint32(p.a[i+3])*uint32(p.da[i+3])+255)>>8)
		p.a[i+4] = uint16(uint32(p.a[i+4]) + uint32(p.da[i+4]) - (uint32(p.a[i+4])*uint32(p.da[i+4])+255)>>8)
		p.a[i+5] = uint16(uint32(p.a[i+5]) + uint32(p.da[i+5]) - (uint32(p.a[i+5])*uint32(p.da[i+5])+255)>>8)
		p.a[i+6] = uint16(uint32(p.a[i+6]) + uint32(p.da[i+6]) - (uint32(p.a[i+6])*uint32(p.da[i+6])+255)>>8)
		p.a[i+7] = uint16(uint32(p.a[i+7]) + uint32(p.da[i+7]) - (uint32(p.a[i+7])*uint32(p.da[i+7])+255)>>8)
	}
}

//go:fix inline
func (p *LowPipeline) Xor() {
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		invDa0, invDa1, invDa2, invDa3 := 255-p.da[i], 255-p.da[i+1], 255-p.da[i+2], 255-p.da[i+3]
		invDa4, invDa5, invDa6, invDa7 := 255-p.da[i+4], 255-p.da[i+5], 255-p.da[i+6], 255-p.da[i+7]
		invSa0, invSa1, invSa2, invSa3 := 255-p.a[i], 255-p.a[i+1], 255-p.a[i+2], 255-p.a[i+3]
		invSa4, invSa5, invSa6, invSa7 := 255-p.a[i+4], 255-p.a[i+5], 255-p.a[i+6], 255-p.a[i+7]

		p.r[i] = uint16((uint32(p.r[i])*uint32(invDa0) + uint32(p.dr[i])*uint32(invSa0) + 255) >> 8)
		p.r[i+1] = uint16((uint32(p.r[i+1])*uint32(invDa1) + uint32(p.dr[i+1])*uint32(invSa1) + 255) >> 8)
		p.r[i+2] = uint16((uint32(p.r[i+2])*uint32(invDa2) + uint32(p.dr[i+2])*uint32(invSa2) + 255) >> 8)
		p.r[i+3] = uint16((uint32(p.r[i+3])*uint32(invDa3) + uint32(p.dr[i+3])*uint32(invSa3) + 255) >> 8)
		p.r[i+4] = uint16((uint32(p.r[i+4])*uint32(invDa4) + uint32(p.dr[i+4])*uint32(invSa4) + 255) >> 8)
		p.r[i+5] = uint16((uint32(p.r[i+5])*uint32(invDa5) + uint32(p.dr[i+5])*uint32(invSa5) + 255) >> 8)
		p.r[i+6] = uint16((uint32(p.r[i+6])*uint32(invDa6) + uint32(p.dr[i+6])*uint32(invSa6) + 255) >> 8)
		p.r[i+7] = uint16((uint32(p.r[i+7])*uint32(invDa7) + uint32(p.dr[i+7])*uint32(invSa7) + 255) >> 8)

		p.g[i] = uint16((uint32(p.g[i])*uint32(invDa0) + uint32(p.dg[i])*uint32(invSa0) + 255) >> 8)
		p.g[i+1] = uint16((uint32(p.g[i+1])*uint32(invDa1) + uint32(p.dg[i+1])*uint32(invSa1) + 255) >> 8)
		p.g[i+2] = uint16((uint32(p.g[i+2])*uint32(invDa2) + uint32(p.dg[i+2])*uint32(invSa2) + 255) >> 8)
		p.g[i+3] = uint16((uint32(p.g[i+3])*uint32(invDa3) + uint32(p.dg[i+3])*uint32(invSa3) + 255) >> 8)
		p.g[i+4] = uint16((uint32(p.g[i+4])*uint32(invDa4) + uint32(p.dg[i+4])*uint32(invSa4) + 255) >> 8)
		p.g[i+5] = uint16((uint32(p.g[i+5])*uint32(invDa5) + uint32(p.dg[i+5])*uint32(invSa5) + 255) >> 8)
		p.g[i+6] = uint16((uint32(p.g[i+6])*uint32(invDa6) + uint32(p.dg[i+6])*uint32(invSa6) + 255) >> 8)
		p.g[i+7] = uint16((uint32(p.g[i+7])*uint32(invDa7) + uint32(p.dg[i+7])*uint32(invSa7) + 255) >> 8)

		p.b[i] = uint16((uint32(p.b[i])*uint32(invDa0) + uint32(p.db[i])*uint32(invSa0) + 255) >> 8)
		p.b[i+1] = uint16((uint32(p.b[i+1])*uint32(invDa1) + uint32(p.db[i+1])*uint32(invSa1) + 255) >> 8)
		p.b[i+2] = uint16((uint32(p.b[i+2])*uint32(invDa2) + uint32(p.db[i+2])*uint32(invSa2) + 255) >> 8)
		p.b[i+3] = uint16((uint32(p.b[i+3])*uint32(invDa3) + uint32(p.db[i+3])*uint32(invSa3) + 255) >> 8)
		p.b[i+4] = uint16((uint32(p.b[i+4])*uint32(invDa4) + uint32(p.db[i+4])*uint32(invSa4) + 255) >> 8)
		p.b[i+5] = uint16((uint32(p.b[i+5])*uint32(invDa5) + uint32(p.db[i+5])*uint32(invSa5) + 255) >> 8)
		p.b[i+6] = uint16((uint32(p.b[i+6])*uint32(invDa6) + uint32(p.db[i+6])*uint32(invSa6) + 255) >> 8)
		p.b[i+7] = uint16((uint32(p.b[i+7])*uint32(invDa7) + uint32(p.db[i+7])*uint32(invSa7) + 255) >> 8)

		p.a[i] = uint16((uint32(p.a[i])*uint32(invDa0) + uint32(p.da[i])*uint32(invSa0) + 255) >> 8)
		p.a[i+1] = uint16((uint32(p.a[i+1])*uint32(invDa1) + uint32(p.da[i+1])*uint32(invSa1) + 255) >> 8)
		p.a[i+2] = uint16((uint32(p.a[i+2])*uint32(invDa2) + uint32(p.da[i+2])*uint32(invSa2) + 255) >> 8)
		p.a[i+3] = uint16((uint32(p.a[i+3])*uint32(invDa3) + uint32(p.da[i+3])*uint32(invSa3) + 255) >> 8)
		p.a[i+4] = uint16((uint32(p.a[i+4])*uint32(invDa4) + uint32(p.da[i+4])*uint32(invSa4) + 255) >> 8)
		p.a[i+5] = uint16((uint32(p.a[i+5])*uint32(invDa5) + uint32(p.da[i+5])*uint32(invSa5) + 255) >> 8)
		p.a[i+6] = uint16((uint32(p.a[i+6])*uint32(invDa6) + uint32(p.da[i+6])*uint32(invSa6) + 255) >> 8)
		p.a[i+7] = uint16((uint32(p.a[i+7])*uint32(invDa7) + uint32(p.da[i+7])*uint32(invSa7) + 255) >> 8)
	}
}

//go:fix inline
func (p *LowPipeline) ColorBurn() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) ColorDodge() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) Darken() {
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		for j := 0; j < 8; j++ {
			// Formula: s + d - div255(max(s * da, d * sa))
			prod1 := uint32(p.r[i+j]) * uint32(p.da[i+j])
			prod2 := uint32(p.dr[i+j]) * uint32(p.a[i+j])
			maxProd := u16max(uint16(prod2), uint16(prod1))
			p.r[i+j] = uint16(uint32(p.r[i+j]) + uint32(p.dr[i+j]) - ((uint32(maxProd) + 255) >> 8))

			prod1 = uint32(p.g[i+j]) * uint32(p.da[i+j])
			prod2 = uint32(p.dg[i+j]) * uint32(p.a[i+j])
			maxProd = u16max(uint16(prod2), uint16(prod1))
			p.g[i+j] = uint16(uint32(p.g[i+j]) + uint32(p.dg[i+j]) - ((uint32(maxProd) + 255) >> 8))

			prod1 = uint32(p.b[i+j]) * uint32(p.da[i+j])
			prod2 = uint32(p.db[i+j]) * uint32(p.a[i+j])
			maxProd = u16max(uint16(prod2), uint16(prod1))
			p.b[i+j] = uint16(uint32(p.b[i+j]) + uint32(p.db[i+j]) - ((uint32(maxProd) + 255) >> 8))

			// Alpha channel: source_over formula
			invSa := 255 - p.a[i+j]
			p.a[i+j] = uint16(uint32(p.a[i+j]) + (uint32(p.da[i+j])*uint32(invSa)+255)>>8)
		}
	}

}

//go:fix inline
func (p *LowPipeline) Difference() {
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		for j := 0; j < 8; j++ {
			// Formula: s + d - 2 * div255(min(s * da, d * sa))
			prod1 := uint32(p.r[i+j]) * uint32(p.da[i+j])
			prod2 := uint32(p.dr[i+j]) * uint32(p.a[i+j])
			minProd := u16min(uint16(prod2), uint16(prod1))
			p.r[i+j] = uint16(uint32(p.r[i+j]) + uint32(p.dr[i+j]) - 2*((uint32(minProd)+255)>>8))

			prod1 = uint32(p.g[i+j]) * uint32(p.da[i+j])
			prod2 = uint32(p.dg[i+j]) * uint32(p.a[i+j])
			minProd = u16min(uint16(prod2), uint16(prod1))
			p.g[i+j] = uint16(uint32(p.g[i+j]) + uint32(p.dg[i+j]) - 2*((uint32(minProd)+255)>>8))

			prod1 = uint32(p.b[i+j]) * uint32(p.da[i+j])
			prod2 = uint32(p.db[i+j]) * uint32(p.a[i+j])
			minProd = u16min(uint16(prod2), uint16(prod1))
			p.b[i+j] = uint16(uint32(p.b[i+j]) + uint32(p.db[i+j]) - 2*((uint32(minProd)+255)>>8))

			// Alpha channel: source_over formula
			invSa := 255 - p.a[i+j]
			p.a[i+j] = uint16(uint32(p.a[i+j]) + (uint32(p.da[i+j])*uint32(invSa)+255)>>8)
		}
	}
}

//go:fix inline
func (p *LowPipeline) Exclusion() {
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		for j := 0; j < 8; j++ {
			// Formula: s + d - 2 * div255(s * d)
			prod := uint32(p.r[i+j]) * uint32(p.dr[i+j])
			p.r[i+j] = uint16(uint32(p.r[i+j]) + uint32(p.dr[i+j]) - 2*((prod+255)>>8))

			prod = uint32(p.g[i+j]) * uint32(p.dg[i+j])
			p.g[i+j] = uint16(uint32(p.g[i+j]) + uint32(p.dg[i+j]) - 2*((prod+255)>>8))

			prod = uint32(p.b[i+j]) * uint32(p.db[i+j])
			p.b[i+j] = uint16(uint32(p.b[i+j]) + uint32(p.db[i+j]) - 2*((prod+255)>>8))

			// Alpha channel: source_over formula
			invSa := 255 - p.a[i+j]
			p.a[i+j] = uint16(uint32(p.a[i+j]) + (uint32(p.da[i+j])*uint32(invSa)+255)>>8)
		}
	}

}

//go:fix inline
func (p *LowPipeline) HardLight() {
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		for j := 0; j < 8; j++ {
			// Formula: div255(s * inv(da) + d * inv(sa) + (2*s <= sa ? 2*s*d : sa*da - 2*(sa-s)*(da-d)))
			invDa := uint32(255 - p.da[i+j])
			invSa := uint32(255 - p.a[i+j])
			sa := uint32(p.a[i+j])
			da := uint32(p.da[i+j])

			// R channel
			s := uint32(p.r[i+j])
			d := uint32(p.dr[i+j])
			term := s*invDa + d*invSa
			if 2*s <= sa {
				term += 2 * s * d
			} else {
				term += sa*da - 2*(sa-s)*(da-d)
			}
			p.r[i+j] = uint16((term + 255) >> 8)

			// G channel
			s = uint32(p.g[i+j])
			d = uint32(p.dg[i+j])
			term = s*invDa + d*invSa
			if 2*s <= sa {
				term += 2 * s * d
			} else {
				term += sa*da - 2*(sa-s)*(da-d)
			}
			p.g[i+j] = uint16((term + 255) >> 8)

			// B channel
			s = uint32(p.b[i+j])
			d = uint32(p.db[i+j])
			term = s*invDa + d*invSa
			if 2*s <= sa {
				term += 2 * s * d
			} else {
				term += sa*da - 2*(sa-s)*(da-d)
			}
			p.b[i+j] = uint16((term + 255) >> 8)

			// Alpha channel: source_over formula
			p.a[i+j] = uint16(uint32(p.a[i+j]) + (uint32(p.da[i+j])*invSa+255)>>8)
		}
	}
}

//go:fix inline
func (p *LowPipeline) Lighten() {
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		for j := 0; j < 8; j++ {
			// Formula: s + d - div255(min(s * da, d * sa))
			prod1 := uint32(p.r[i+j]) * uint32(p.da[i+j])
			prod2 := uint32(p.dr[i+j]) * uint32(p.a[i+j])
			minProd := u16min(uint16(prod2), uint16(prod1))
			p.r[i+j] = uint16(uint32(p.r[i+j]) + uint32(p.dr[i+j]) - ((uint32(minProd) + 255) >> 8))

			prod1 = uint32(p.g[i+j]) * uint32(p.da[i+j])
			prod2 = uint32(p.dg[i+j]) * uint32(p.a[i+j])
			minProd = u16min(uint16(prod2), uint16(prod1))
			p.g[i+j] = uint16(uint32(p.g[i+j]) + uint32(p.dg[i+j]) - ((uint32(minProd) + 255) >> 8))

			prod1 = uint32(p.b[i+j]) * uint32(p.da[i+j])
			prod2 = uint32(p.db[i+j]) * uint32(p.a[i+j])
			minProd = u16min(uint16(prod2), uint16(prod1))
			p.b[i+j] = uint16(uint32(p.b[i+j]) + uint32(p.db[i+j]) - ((uint32(minProd) + 255) >> 8))

			// Alpha channel: source_over formula
			invSa := 255 - p.a[i+j]
			p.a[i+j] = uint16(uint32(p.a[i+j]) + (uint32(p.da[i+j])*uint32(invSa)+255)>>8)
		}
	}
}

//go:fix inline
func (p *LowPipeline) Overlay() {
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		for j := 0; j < 8; j++ {
			sa := uint32(p.a[i+j])
			da := uint32(p.da[i+j])
			sr, sg, sb := uint32(p.r[i+j]), uint32(p.g[i+j]), uint32(p.b[i+j])
			dr, dg, db := uint32(p.dr[i+j]), uint32(p.dg[i+j]), uint32(p.db[i+j])
			invSaVal, invDaVal := uint32(255-p.a[i+j]), uint32(255-p.da[i+j])

			// R channel
			term := sr*invDaVal + dr*invSaVal
			if 2*dr <= da {
				term += 2 * sr * dr
			} else {
				term += sa*da - 2*(sa-sr)*(da-dr)
			}
			p.r[i+j] = uint16((term + 255) >> 8)

			// G channel
			term = sg*invDaVal + dg*invSaVal
			if 2*dg <= da {
				term += 2 * sg * dg
			} else {
				term += sa*da - 2*(sa-sg)*(da-dg)
			}
			p.g[i+j] = uint16((term + 255) >> 8)

			// B channel
			term = sb*invDaVal + db*invSaVal
			if 2*db <= da {
				term += 2 * sb * db
			} else {
				term += sa*da - 2*(sa-sb)*(da-db)
			}
			p.b[i+j] = uint16((term + 255) >> 8)

			// Alpha channel: source_over formula
			p.a[i+j] = uint16(sa + (da*invSaVal+255)>>8)
		}
	}
}

//go:fix inline
func (p *LowPipeline) SoftLight() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) Hue() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) Saturation() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) Color() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) Luminosity() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) SourceOverRgba() {
	// TODO
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		invSa0 := 255 - p.a[i]
		p.r[i] = uint16(uint32(p.r[i]) + (uint32(p.dr[i])*uint32(invSa0)+255)>>8)
		p.g[i] = uint16(uint32(p.g[i]) + (uint32(p.dg[i])*uint32(invSa0)+255)>>8)
		p.b[i] = uint16(uint32(p.b[i]) + (uint32(p.db[i])*uint32(invSa0)+255)>>8)
		p.a[i] = uint16(uint32(p.a[i]) + (uint32(p.da[i])*uint32(invSa0)+255)>>8)
	}
}

//go:fix inline
func (p *LowPipeline) SourceOverRgbaTail() {
	// TODO
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		invSa0 := 255 - p.a[i]
		p.r[i] = uint16(uint32(p.r[i]) + (uint32(p.dr[i])*uint32(invSa0)+255)>>8)
		p.g[i] = uint16(uint32(p.g[i]) + (uint32(p.dg[i])*uint32(invSa0)+255)>>8)
		p.b[i] = uint16(uint32(p.b[i]) + (uint32(p.db[i])*uint32(invSa0)+255)>>8)
		p.a[i] = uint16(uint32(p.a[i]) + (uint32(p.da[i])*uint32(invSa0)+255)>>8)
	}
}

//go:fix inline
func (p *LowPipeline) Transform() {
	// nx = r * sx + b * kx + tx
	// ny = r * ky + b * sy + ty
	ts := &p.ctx.Transform

	// Join r,g into x and b,a into y as float32 arrays
	var x, y [16]float32
	for i := 0; i < 16; i++ {
		// Convert u16 pair to float32 using bit casting
		x[i] = math.Float32frombits(uint32(p.r[i]) | uint32(p.g[i])<<16)
		y[i] = math.Float32frombits(uint32(p.b[i]) | uint32(p.a[i])<<16)
	}

	// Apply transform: nx = x*sx + y*kx + tx
	// ny = x*ky + y*sy + ty
	var nx, ny [16]float32
	for i := 0; i < 16; i++ {
		nx[i] = x[i]*ts.SX + y[i]*ts.KX + ts.TX
		ny[i] = x[i]*ts.KY + y[i]*ts.SY + ts.TY
	}

	// Split back to u16
	for i := 0; i < 16; i++ {
		nxBits := math.Float32bits(nx[i])
		nyBits := math.Float32bits(ny[i])
		p.r[i] = uint16(nxBits & 0xFFFF)
		p.g[i] = uint16(nxBits >> 16)
		p.b[i] = uint16(nyBits & 0xFFFF)
		p.a[i] = uint16(nyBits >> 16)
	}
}

//go:fix inline
func (p *LowPipeline) Reflect() {
	// null_fn

}

//go:fix inline
func (p *LowPipeline) Repeat() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) Bilinear() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) Bicubic() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) PadX1() {
	// Clamps x to [0, 1] range using normalize()
	for i := 0; i < LOW_STAGE_WIDTH; i++ {
		// Convert to float32
		x := math.Float32frombits(uint32(p.r[i]) | uint32(p.g[i])<<16)
		// Normalize/clamp to [0, 1] range
		if x < 0.0 {
			x = 0.0
		} else if x > 1.0 {
			x = 1.0
		}
		// Convert back to u16 representation
		xBits := math.Float32bits(x)
		p.r[i] = uint16(xBits & 0xFFFF)
		p.g[i] = uint16(xBits >> 16)
	}
}

//go:fix inline
func (p *LowPipeline) ReflectX1() {
	// Mirrors x at integer boundaries: x = |x - 1 - 2*floor((x-1)*0.5)| - 1
	for i := 0; i < LOW_STAGE_WIDTH; i++ {
		// Convert to float32
		x := math.Float32frombits(uint32(p.r[i]) | uint32(p.g[i])<<16)
		// Apply reflect formula: x = |x - 1 - 2*floor((x-1)*0.5)| - 1
		xMinus1 := x - 1.0
		floored := float32(int(xMinus1 * 0.5))
		reflected := xMinus1 - 2.0*floored - 1.0
		if reflected < 0 {
			reflected = -reflected
		}
		// Convert back to u16 representation
		xBits := math.Float32bits(reflected)
		p.r[i] = uint16(xBits & 0xFFFF)
		p.g[i] = uint16(xBits >> 16)
	}
}

//go:fix inline
func (p *LowPipeline) RepeatX1() {
	// Repeats pattern every integer: x = x - floor(x)
	for i := 0; i < LOW_STAGE_WIDTH; i++ {
		// Convert to float32
		x := math.Float32frombits(uint32(p.r[i]) | uint32(p.g[i])<<16)
		// Apply repeat formula: x = x - floor(x)
		floored := float32(int(x))
		if floored > x {
			floored--
		}
		repeated := x - floored
		// Convert back to u16 representation
		xBits := math.Float32bits(repeated)
		p.r[i] = uint16(xBits & 0xFFFF)
		p.g[i] = uint16(xBits >> 16)
	}
}

//go:fix inline
func (p *LowPipeline) Gradient() {
	// Interpolates gradient colors based on t value in p.r/p.g
	ctx := p.ctx.Gradient
	if ctx.Len > 0 {
		// Join r,g into t values as float32
		var t [16]float32
		for i := 0; i < 16; i++ {
			t[i] = math.Float32frombits(uint32(p.r[i]) | uint32(p.g[i])<<16)
		}

		// For each pixel, find which stop interval it falls into
		for i := 0; i < LOW_STAGE_WIDTH; i++ {
			ti := t[i]
			idx := uint16(0)
			// Find the stop index where t >= stop.t
			for j := 1; j < ctx.Len; j++ {
				if ti >= ctx.TValues[j].Get() {
					idx = uint16(j)
				}
			}

			// Use bias color from gradient context (already interpolated)
			if int(idx) < 16 {
				color := ctx.Biases[idx]
				p.r[i] = uint16(color.R)
				p.g[i] = uint16(color.G)
				p.b[i] = uint16(color.B)
				p.a[i] = uint16(color.A)
			}
		}
	}
}

//go:fix inline
func (p *LowPipeline) EvenlySpaced2StopGradient() {
	factor := p.ctx.EvenlySpaced2StopGradient.Factor
	bias := p.ctx.EvenlySpaced2StopGradient.Bias
	for i := 0; i < LOW_STAGE_WIDTH; i += 8 {
		p.r[i] = uint16(factor.R*float32(p.r[i]) + bias.R)
		p.g[i] = uint16(factor.G*float32(p.g[i]) + bias.G)
		p.b[i] = uint16(factor.B*float32(p.b[i]) + bias.B)
		p.a[i] = uint16(factor.A*float32(p.a[i]) + bias.A)
	}
}

//go:fix inline
func (p *LowPipeline) XYToUnitAngle() {
	for i := 0; i < LOW_STAGE_WIDTH; i++ {
		// Convert x,y to float32
		x := math.Float32frombits(uint32(p.r[i]) | uint32(p.g[i])<<16)
		y := math.Float32frombits(uint32(p.b[i]) | uint32(p.a[i])<<16)
		// Calculate radius: r = sqrt(x^2 + y^2)
		r := float32(math.Sqrt(float64(x*x + y*y)))
		// Convert back to u16 representation
		rBits := math.Float32bits(r)
		p.r[i] = uint16(rBits & 0xFFFF)
		p.g[i] = uint16(rBits >> 16)
		p.b[i] = uint16(0)
		p.a[i] = uint16(0)
	}
}

//go:fix inline
func (p *LowPipeline) XYToRadius() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) XYTo2PtConicalFocalOnCircle() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) XYTo2PtConicalWellBehaved() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) XYTo2PtConicalSmaller() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) XYTo2PtConicalGreater() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) XYTo2PtConicalStrip() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) Mask2PtConicalNan() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) Mask2PtConicalDegenerates() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) ApplyVectorMask() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) Alter2PtConicalCompensateFocal() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) Alter2PtConicalUnswap() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) NegateX() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) ApplyConcentricScaleBias() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) GammaExpand2() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) GammaExpandDestination2() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) GammaCompress2() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) GammaExpand22() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) GammaExpandDestination22() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) GammaCompress22() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) GammaExpandSrgb() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) GammaExpandDestinationSrgb() {
	// null_fn
}

//go:fix inline
func (p *LowPipeline) GammaCompressSrgb() {
	// null_fn
}

func u16min(a, b uint16) uint16 {
	if a < b {
		return a
	}
	return b
}

func u16max(a, b uint16) uint16 {
	if a > b {
		return a
	}
	return b
}
