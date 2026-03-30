// Copyright 2016 Google Inc.
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package pipeline

import (
	"math"

	"github.com/chewxy/math32"
)

//go:fix inline
func (p *HighPipeline) MoveSourceToDestination() {
	copy(p.dr[:], p.r[:])
	copy(p.dg[:], p.g[:])
	copy(p.db[:], p.b[:])
	copy(p.da[:], p.a[:])
}

//go:fix inline
func (p *HighPipeline) MoveDestinationToSource() {
	copy(p.r[:], p.dr[:])
	copy(p.g[:], p.dg[:])
	copy(p.b[:], p.db[:])
	copy(p.a[:], p.da[:])
}

//go:fix inline
func (p *HighPipeline) Clamp0() {
	p.r[0], p.r[1], p.r[2], p.r[3] = f32max(p.r[0], 0), f32max(p.r[1], 0), f32max(p.r[2], 0), f32max(p.r[3], 0)
	p.r[4], p.r[5], p.r[6], p.r[7] = f32max(p.r[4], 0), f32max(p.r[5], 0), f32max(p.r[6], 0), f32max(p.r[7], 0)
	p.g[0], p.g[1], p.g[2], p.g[3] = f32max(p.g[0], 0), f32max(p.g[1], 0), f32max(p.g[2], 0), f32max(p.g[3], 0)
	p.g[4], p.g[5], p.g[6], p.g[7] = f32max(p.g[4], 0), f32max(p.g[5], 0), f32max(p.g[6], 0), f32max(p.g[7], 0)
	p.b[0], p.b[1], p.b[2], p.b[3] = f32max(p.b[0], 0), f32max(p.b[1], 0), f32max(p.b[2], 0), f32max(p.b[3], 0)
	p.b[4], p.b[5], p.b[6], p.b[7] = f32max(p.b[4], 0), f32max(p.b[5], 0), f32max(p.b[6], 0), f32max(p.b[7], 0)
	p.a[0], p.a[1], p.a[2], p.a[3] = f32max(p.a[0], 0), f32max(p.a[1], 0), f32max(p.a[2], 0), f32max(p.a[3], 0)
	p.a[4], p.a[5], p.a[6], p.a[7] = f32max(p.a[4], 0), f32max(p.a[5], 0), f32max(p.a[6], 0), f32max(p.a[7], 0)
}

//go:fix inline
func (p *HighPipeline) ClampA() {
	p.r[0], p.r[1], p.r[2], p.r[3] = f32min(p.r[0], 1), f32min(p.r[1], 1), f32min(p.r[2], 1), f32min(p.r[3], 1)
	p.r[4], p.r[5], p.r[6], p.r[7] = f32min(p.r[4], 1), f32min(p.r[5], 1), f32min(p.r[6], 1), f32min(p.r[7], 1)
	p.g[0], p.g[1], p.g[2], p.g[3] = f32min(p.g[0], 1), f32min(p.g[1], 1), f32min(p.g[2], 1), f32min(p.g[3], 1)
	p.g[4], p.g[5], p.g[6], p.g[7] = f32min(p.g[4], 1), f32min(p.g[5], 1), f32min(p.g[6], 1), f32min(p.g[7], 1)
	p.b[0], p.b[1], p.b[2], p.b[3] = f32min(p.b[0], 1), f32min(p.b[1], 1), f32min(p.b[2], 1), f32min(p.b[3], 1)
	p.b[4], p.b[5], p.b[6], p.b[7] = f32min(p.b[4], 1), f32min(p.b[5], 1), f32min(p.b[6], 1), f32min(p.b[7], 1)
	p.a[0], p.a[1], p.a[2], p.a[3] = f32min(p.a[0], 1), f32min(p.a[1], 1), f32min(p.a[2], 1), f32min(p.a[3], 1)
	p.a[4], p.a[5], p.a[6], p.a[7] = f32min(p.a[4], 1), f32min(p.a[5], 1), f32min(p.a[6], 1), f32min(p.a[7], 1)
}

//go:fix inline
func (p *HighPipeline) Premultiply() {
	p.r[0], p.r[1], p.r[2], p.r[3] = p.r[0]*p.a[0], p.r[1]*p.a[1], p.r[2]*p.a[2], p.r[3]*p.a[3]
	p.r[4], p.r[5], p.r[6], p.r[7] = p.r[4]*p.a[4], p.r[5]*p.a[5], p.r[6]*p.a[6], p.r[7]*p.a[7]
	p.g[0], p.g[1], p.g[2], p.g[3] = p.g[0]*p.a[0], p.g[1]*p.a[1], p.g[2]*p.a[2], p.g[3]*p.a[3]
	p.g[4], p.g[5], p.g[6], p.g[7] = p.g[4]*p.a[4], p.g[5]*p.a[5], p.g[6]*p.a[6], p.g[7]*p.a[7]
	p.b[0], p.b[1], p.b[2], p.b[3] = p.b[0]*p.a[0], p.b[1]*p.a[1], p.b[2]*p.a[2], p.b[3]*p.a[3]
	p.b[4], p.b[5], p.b[6], p.b[7] = p.b[4]*p.a[4], p.b[5]*p.a[5], p.b[6]*p.a[6], p.b[7]*p.a[7]
}

//go:fix inline
func (p *HighPipeline) UniformColor() {
	uc := p.ctx.UniformColor

	p.r[0], p.r[1], p.r[2], p.r[3] = uc.R, uc.R, uc.R, uc.R
	p.r[4], p.r[5], p.r[6], p.r[7] = uc.R, uc.R, uc.R, uc.R
	p.g[0], p.g[1], p.g[2], p.g[3] = uc.G, uc.G, uc.G, uc.G
	p.g[4], p.g[5], p.g[6], p.g[7] = uc.G, uc.G, uc.G, uc.G
	p.b[0], p.b[1], p.b[2], p.b[3] = uc.B, uc.B, uc.B, uc.B
	p.b[4], p.b[5], p.b[6], p.b[7] = uc.B, uc.B, uc.B, uc.B
	p.a[0], p.a[1], p.a[2], p.a[3] = uc.A, uc.A, uc.A, uc.A
	p.a[4], p.a[5], p.a[6], p.a[7] = uc.A, uc.A, uc.A, uc.A
}

//go:fix inline
func (p *HighPipeline) SeedShader() {
	p.r[0], p.r[1], p.r[2], p.r[3] = float32(p.dx)+0.5, float32(p.dx)+1.5, float32(p.dx)+2.5, float32(p.dx)+3.5
	p.r[4], p.r[5], p.r[6], p.r[7] = float32(p.dx)+4.5, float32(p.dx)+5.5, float32(p.dx)+6.5, float32(p.dx)+7.5

	p.g[0], p.g[1], p.g[2], p.g[3] = float32(p.dy)+0.5, float32(p.dy)+0.5, float32(p.dy)+0.5, float32(p.dy)+0.5
	p.g[4], p.g[5], p.g[6], p.g[7] = float32(p.dy)+0.5, float32(p.dy)+0.5, float32(p.dy)+0.5, float32(p.dy)+0.5
	p.b[0], p.b[1], p.b[2], p.b[3] = 1.0, 1.0, 1.0, 1.0
	p.b[4], p.b[5], p.b[6], p.b[7] = 1.0, 1.0, 1.0, 1.0
	p.a[0], p.a[1], p.a[2], p.a[3] = 0.0, 0.0, 0.0, 0.0
	p.a[4], p.a[5], p.a[6], p.a[7] = 0.0, 0.0, 0.0, 0.0

	p.dr[0], p.dr[1], p.dr[2], p.dr[3] = 0.0, 0.0, 0.0, 0.0
	p.dr[4], p.dr[5], p.dr[6], p.dr[7] = 0.0, 0.0, 0.0, 0.0
	p.dg[0], p.dg[1], p.dg[2], p.dg[3] = 0.0, 0.0, 0.0, 0.0
	p.dg[4], p.dg[5], p.dg[6], p.dg[7] = 0.0, 0.0, 0.0, 0.0
	p.db[0], p.db[1], p.db[2], p.db[3] = 0.0, 0.0, 0.0, 0.0
	p.db[4], p.db[5], p.db[6], p.db[7] = 0.0, 0.0, 0.0, 0.0
	p.da[0], p.da[1], p.da[2], p.da[3] = 0.0, 0.0, 0.0, 0.0
	p.da[4], p.da[5], p.da[6], p.da[7] = 0.0, 0.0, 0.0, 0.0

}

//go:fix inline
func (p *HighPipeline) LoadDestination() {
	const FACTOR = 1.0 / 255.0

	offset := (p.dy*p.pixmapDst.RealWidth + p.dx) * 4
	data := p.pixmapDst.Data[offset : offset+HIGH_STAGE_WIDTH*4]

	p.dr[0], p.dr[1], p.dr[2], p.dr[3] = float32(data[0])*FACTOR, float32(data[4])*FACTOR, float32(data[8])*FACTOR, float32(data[12])*FACTOR
	p.dr[4], p.dr[5], p.dr[6], p.dr[7] = float32(data[16])*FACTOR, float32(data[20])*FACTOR, float32(data[24])*FACTOR, float32(data[28])*FACTOR
	p.dg[0], p.dg[1], p.dg[2], p.dg[3] = float32(data[1])*FACTOR, float32(data[5])*FACTOR, float32(data[9])*FACTOR, float32(data[13])*FACTOR
	p.dg[4], p.dg[5], p.dg[6], p.dg[7] = float32(data[17])*FACTOR, float32(data[21])*FACTOR, float32(data[25])*FACTOR, float32(data[29])*FACTOR
	p.db[0], p.db[1], p.db[2], p.db[3] = float32(data[2])*FACTOR, float32(data[6])*FACTOR, float32(data[10])*FACTOR, float32(data[14])*FACTOR
	p.db[4], p.db[5], p.db[6], p.db[7] = float32(data[18])*FACTOR, float32(data[22])*FACTOR, float32(data[26])*FACTOR, float32(data[30])*FACTOR
	p.da[0], p.da[1], p.da[2], p.da[3] = float32(data[3])*FACTOR, float32(data[7])*FACTOR, float32(data[11])*FACTOR, float32(data[15])*FACTOR
	p.da[4], p.da[5], p.da[6], p.da[7] = float32(data[19])*FACTOR, float32(data[23])*FACTOR, float32(data[27])*FACTOR, float32(data[31])*FACTOR
}

//go:fix inline
func (p *HighPipeline) LoadDestinationTail() {
	const FACTOR = 1.0 / 255.0

	tmp := [HIGH_STAGE_WIDTH * 4]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	offset := (p.dy*p.pixmapDst.RealWidth + p.dx) * 4
	for i := 0; i < p.tail; i++ {
		srcIdx := offset + i*4
		dstIdx := i * 4
		tmp[dstIdx+0] = p.pixmapDst.Data[srcIdx+0]
		tmp[dstIdx+1] = p.pixmapDst.Data[srcIdx+1]
		tmp[dstIdx+2] = p.pixmapDst.Data[srcIdx+2]
		tmp[dstIdx+3] = p.pixmapDst.Data[srcIdx+3]
	}

	// Load from tmp array (rest are already zero/transparent)
	p.dr[0], p.dr[1], p.dr[2], p.dr[3] = float32(tmp[0])*FACTOR, float32(tmp[4])*FACTOR, float32(tmp[8])*FACTOR, float32(tmp[12])*FACTOR
	p.dr[4], p.dr[5], p.dr[6], p.dr[7] = float32(tmp[16])*FACTOR, float32(tmp[20])*FACTOR, float32(tmp[24])*FACTOR, float32(tmp[28])*FACTOR
	p.dg[0], p.dg[1], p.dg[2], p.dg[3] = float32(tmp[1])*FACTOR, float32(tmp[5])*FACTOR, float32(tmp[9])*FACTOR, float32(tmp[13])*FACTOR
	p.dg[4], p.dg[5], p.dg[6], p.dg[7] = float32(tmp[17])*FACTOR, float32(tmp[21])*FACTOR, float32(tmp[25])*FACTOR, float32(tmp[29])*FACTOR
	p.db[0], p.db[1], p.db[2], p.db[3] = float32(tmp[2])*FACTOR, float32(tmp[6])*FACTOR, float32(tmp[10])*FACTOR, float32(tmp[14])*FACTOR
	p.db[4], p.db[5], p.db[6], p.db[7] = float32(tmp[18])*FACTOR, float32(tmp[22])*FACTOR, float32(tmp[26])*FACTOR, float32(tmp[30])*FACTOR
	p.da[0], p.da[1], p.da[2], p.da[3] = float32(tmp[3])*FACTOR, float32(tmp[7])*FACTOR, float32(tmp[11])*FACTOR, float32(tmp[15])*FACTOR
	p.da[4], p.da[5], p.da[6], p.da[7] = float32(tmp[19])*FACTOR, float32(tmp[23])*FACTOR, float32(tmp[27])*FACTOR, float32(tmp[31])*FACTOR
}

//go:fix inline
func (p *HighPipeline) Store() {
	offset := (p.dy*p.pixmapDst.RealWidth + p.dx) * 4
	data := p.pixmapDst.Data[offset : offset+HIGH_STAGE_WIDTH*4]

	data[0], data[4], data[8], data[12] = f32unnorm(p.r[0]), f32unnorm(p.r[1]), f32unnorm(p.r[2]), f32unnorm(p.r[3])
	data[16], data[20], data[24], data[28] = f32unnorm(p.r[4]), f32unnorm(p.r[5]), f32unnorm(p.r[6]), f32unnorm(p.r[7])
	data[1], data[5], data[9], data[13] = f32unnorm(p.g[0]), f32unnorm(p.g[1]), f32unnorm(p.g[2]), f32unnorm(p.g[3])
	data[17], data[21], data[25], data[29] = f32unnorm(p.g[4]), f32unnorm(p.g[5]), f32unnorm(p.g[6]), f32unnorm(p.g[7])
	data[2], data[6], data[10], data[14] = f32unnorm(p.b[0]), f32unnorm(p.b[1]), f32unnorm(p.b[2]), f32unnorm(p.b[3])
	data[18], data[22], data[26], data[30] = f32unnorm(p.b[4]), f32unnorm(p.b[5]), f32unnorm(p.b[6]), f32unnorm(p.b[7])
	data[3], data[7], data[11], data[15] = f32unnorm(p.a[0]), f32unnorm(p.a[1]), f32unnorm(p.a[2]), f32unnorm(p.a[3])
	data[19], data[23], data[27], data[31] = f32unnorm(p.a[4]), f32unnorm(p.a[5]), f32unnorm(p.a[6]), f32unnorm(p.a[7])
}

//go:fix inline
func (p *HighPipeline) StoreTail() {
	tmp := [32]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	tmp[0], tmp[4], tmp[8], tmp[12] = f32unnorm(p.r[0]), f32unnorm(p.r[1]), f32unnorm(p.r[2]), f32unnorm(p.r[3])
	tmp[16], tmp[20], tmp[24], tmp[28] = f32unnorm(p.r[4]), f32unnorm(p.r[5]), f32unnorm(p.r[6]), f32unnorm(p.r[7])
	tmp[1], tmp[5], tmp[9], tmp[13] = f32unnorm(p.g[0]), f32unnorm(p.g[1]), f32unnorm(p.g[2]), f32unnorm(p.g[3])
	tmp[17], tmp[21], tmp[25], tmp[29] = f32unnorm(p.g[4]), f32unnorm(p.g[5]), f32unnorm(p.g[6]), f32unnorm(p.g[7])
	tmp[2], tmp[6], tmp[10], tmp[14] = f32unnorm(p.b[0]), f32unnorm(p.b[1]), f32unnorm(p.b[2]), f32unnorm(p.b[3])
	tmp[18], tmp[22], tmp[26], tmp[30] = f32unnorm(p.b[4]), f32unnorm(p.b[5]), f32unnorm(p.b[6]), f32unnorm(p.b[7])
	tmp[3], tmp[7], tmp[11], tmp[15] = f32unnorm(p.a[0]), f32unnorm(p.a[1]), f32unnorm(p.a[2]), f32unnorm(p.a[3])
	tmp[19], tmp[23], tmp[27], tmp[31] = f32unnorm(p.a[4]), f32unnorm(p.a[5]), f32unnorm(p.a[6]), f32unnorm(p.a[7])

	offset := (p.dy*p.pixmapDst.RealWidth + p.dx) * 4
	for i := 0; i < p.tail; i++ {
		dstIdx := offset + i*4
		p.pixmapDst.Data[dstIdx+0] = tmp[i*4+0]
		p.pixmapDst.Data[dstIdx+1] = tmp[i*4+1]
		p.pixmapDst.Data[dstIdx+2] = tmp[i*4+2]
		p.pixmapDst.Data[dstIdx+3] = tmp[i*4+3]
	}
}

//go:fix inline
func (p *HighPipeline) LoadDestinationU8() {
	// unreachable for highp
}

//go:fix inline
func (p *HighPipeline) LoadDestinationU8Tail() {
	// unreachable for highp
}

//go:fix inline
func (p *HighPipeline) StoreU8() {
	// unreachable for highp
}

//go:fix inline
func (p *HighPipeline) StoreU8Tail() {
	// unreachable for highp
}

//go:fix inline
func (p *HighPipeline) Gather() {
	const FACTOR = 1.0 / 255.0

	ulpsub := func(v float32) float32 {
		bits := math.Float32bits(v)
		return math.Float32frombits(bits - 1)
	}

	w := ulpsub(float32(p.pixmapSrc.Size.Width()))
	h := ulpsub(float32(p.pixmapSrc.Size.Height()))
	iw := int32(p.pixmapSrc.Size.Width())

	for i := 0; i < 8; i++ {
		x := f32min(f32max(p.r[i], 0), w)
		y := f32min(f32max(p.g[i], 0), h)

		offset := (int32(y)*iw + int32(x)) * 4

		p.r[i] = float32(p.pixmapSrc.Data[offset+0]) * FACTOR
		p.g[i] = float32(p.pixmapSrc.Data[offset+1]) * FACTOR
		p.b[i] = float32(p.pixmapSrc.Data[offset+2]) * FACTOR
		p.a[i] = float32(p.pixmapSrc.Data[offset+3]) * FACTOR
	}
}

//go:fix inline
func (p *HighPipeline) LoadMaskU8() {
	// unreachable for highp
}

//go:fix inline
func (p *HighPipeline) MaskU8() {
	const FACTOR = 1.0 / 255.0

	offset := p.maskCtx.Offset(p.dx, p.dy)

	c := [8]float32{0, 0, 0, 0, 0, 0, 0, 0}
	for i := 0; i < p.tail && i < 8; i++ {
		c[i] = float32(p.maskCtx.Data[offset+i]) * FACTOR
	}

	if c[0] == 0 && c[1] == 0 && c[2] == 0 && c[3] == 0 && c[4] == 0 && c[5] == 0 && c[6] == 0 && c[7] == 0 {
		return
	}

	p.r[0], p.r[1], p.r[2], p.r[3] = p.r[0]*c[0], p.r[1]*c[1], p.r[2]*c[2], p.r[3]*c[3]
	p.r[4], p.r[5], p.r[6], p.r[7] = p.r[4]*c[4], p.r[5]*c[5], p.r[6]*c[6], p.r[7]*c[7]
	p.g[0], p.g[1], p.g[2], p.g[3] = p.g[0]*c[0], p.g[1]*c[1], p.g[2]*c[2], p.g[3]*c[3]
	p.g[4], p.g[5], p.g[6], p.g[7] = p.g[4]*c[4], p.g[5]*c[5], p.g[6]*c[6], p.g[7]*c[7]
	p.b[0], p.b[1], p.b[2], p.b[3] = p.b[0]*c[0], p.b[1]*c[1], p.b[2]*c[2], p.b[3]*c[3]
	p.b[4], p.b[5], p.b[6], p.b[7] = p.b[4]*c[4], p.b[5]*c[5], p.b[6]*c[6], p.b[7]*c[7]
	p.a[0], p.a[1], p.a[2], p.a[3] = p.a[0]*c[0], p.a[1]*c[1], p.a[2]*c[2], p.a[3]*c[3]
	p.a[4], p.a[5], p.a[6], p.a[7] = p.a[4]*c[4], p.a[5]*c[5], p.a[6]*c[6], p.a[7]*c[7]
}

//go:fix inline
func (p *HighPipeline) ScaleU8() {
	const FACTOR = 1.0 / 255.0

	data := p.aaMaskCtx.CopyAtXY(p.dx, p.dy, p.tail)

	c0 := float32(data[0]) * FACTOR
	c1 := float32(data[1]) * FACTOR

	p.r[0], p.r[1], p.r[2], p.r[3] = p.r[0]*c0, p.r[1]*c1, 0, 0
	p.r[4], p.r[5], p.r[6], p.r[7] = 0, 0, 0, 0
	p.g[0], p.g[1], p.g[2], p.g[3] = p.g[0]*c0, p.g[1]*c1, 0, 0
	p.g[4], p.g[5], p.g[6], p.g[7] = 0, 0, 0, 0
	p.b[0], p.b[1], p.b[2], p.b[3] = p.b[0]*c0, p.b[1]*c1, 0, 0
	p.b[4], p.b[5], p.b[6], p.b[7] = 0, 0, 0, 0
	p.a[0], p.a[1], p.a[2], p.a[3] = p.a[0]*c0, p.a[1]*c1, 0, 0
	p.a[4], p.a[5], p.a[6], p.a[7] = 0, 0, 0, 0
}

//go:fix inline
func (p *HighPipeline) LerpU8() {
	const FACTOR = 1.0 / 255.0

	data := p.aaMaskCtx.CopyAtXY(p.dx, p.dy, p.tail)

	c0 := float32(data[0]) * FACTOR
	c1 := float32(data[1]) * FACTOR

	p.r[0], p.r[1], p.r[2], p.r[3] = f32lerp(p.dr[0], p.r[0], c0), f32lerp(p.dr[1], p.r[1], c1), 0, 0
	p.r[4], p.r[5], p.r[6], p.r[7] = 0, 0, 0, 0
	p.g[0], p.g[1], p.g[2], p.g[3] = f32lerp(p.dg[0], p.g[0], c0), f32lerp(p.dg[1], p.g[1], c1), 0, 0
	p.g[4], p.g[5], p.g[6], p.g[7] = 0, 0, 0, 0
	p.b[0], p.b[1], p.b[2], p.b[3] = f32lerp(p.db[0], p.b[0], c0), f32lerp(p.db[1], p.b[1], c1), 0, 0
	p.b[4], p.b[5], p.b[6], p.b[7] = 0, 0, 0, 0
	p.a[0], p.a[1], p.a[2], p.a[3] = f32lerp(p.da[0], p.a[0], c0), f32lerp(p.da[1], p.a[1], c1), 0, 0
	p.a[4], p.a[5], p.a[6], p.a[7] = 0, 0, 0, 0
}

//go:fix inline
func (p *HighPipeline) Scale1Float() {
	coverage := p.ctx.CurrentCoverage

	p.r[0], p.r[1], p.r[2], p.r[3] = p.r[0]*coverage, p.r[1]*coverage, p.r[2]*coverage, p.r[3]*coverage
	p.r[4], p.r[5], p.r[6], p.r[7] = p.r[4]*coverage, p.r[5]*coverage, p.r[6]*coverage, p.r[7]*coverage
	p.g[0], p.g[1], p.g[2], p.g[3] = p.g[0]*coverage, p.g[1]*coverage, p.g[2]*coverage, p.g[3]*coverage
	p.g[4], p.g[5], p.g[6], p.g[7] = p.g[4]*coverage, p.g[5]*coverage, p.g[6]*coverage, p.g[7]*coverage
	p.b[0], p.b[1], p.b[2], p.b[3] = p.b[0]*coverage, p.b[1]*coverage, p.b[2]*coverage, p.b[3]*coverage
	p.b[4], p.b[5], p.b[6], p.b[7] = p.b[4]*coverage, p.b[5]*coverage, p.b[6]*coverage, p.b[7]*coverage
	p.a[0], p.a[1], p.a[2], p.a[3] = p.a[0]*coverage, p.a[1]*coverage, p.a[2]*coverage, p.a[3]*coverage
	p.a[4], p.a[5], p.a[6], p.a[7] = p.a[4]*coverage, p.a[5]*coverage, p.a[6]*coverage, p.a[7]*coverage
}

//go:fix inline
func (p *HighPipeline) Lerp1Float() {
	coverage := p.ctx.CurrentCoverage

	p.r[0], p.r[1], p.r[2], p.r[3] = f32lerp(p.dr[0], p.r[0], coverage), f32lerp(p.dr[1], p.r[1], coverage), f32lerp(p.dr[2], p.r[2], coverage), f32lerp(p.dr[3], p.r[3], coverage)
	p.r[4], p.r[5], p.r[6], p.r[7] = f32lerp(p.dr[4], p.r[4], coverage), f32lerp(p.dr[5], p.r[5], coverage), f32lerp(p.dr[6], p.r[6], coverage), f32lerp(p.dr[7], p.r[7], coverage)
	p.g[0], p.g[1], p.g[2], p.g[3] = f32lerp(p.dg[0], p.g[0], coverage), f32lerp(p.dg[1], p.g[1], coverage), f32lerp(p.dg[2], p.g[2], coverage), f32lerp(p.dg[3], p.g[3], coverage)
	p.g[4], p.g[5], p.g[6], p.g[7] = f32lerp(p.dg[4], p.g[4], coverage), f32lerp(p.dg[5], p.g[5], coverage), f32lerp(p.dg[6], p.g[6], coverage), f32lerp(p.dg[7], p.g[7], coverage)
	p.b[0], p.b[1], p.b[2], p.b[3] = f32lerp(p.db[0], p.b[0], coverage), f32lerp(p.db[1], p.b[1], coverage), f32lerp(p.db[2], p.b[2], coverage), f32lerp(p.db[3], p.b[3], coverage)
	p.b[4], p.b[5], p.b[6], p.b[7] = f32lerp(p.db[4], p.b[4], coverage), f32lerp(p.db[5], p.b[5], coverage), f32lerp(p.db[6], p.b[6], coverage), f32lerp(p.db[7], p.b[7], coverage)
	p.a[0], p.a[1], p.a[2], p.a[3] = f32lerp(p.da[0], p.a[0], coverage), f32lerp(p.da[1], p.a[1], coverage), f32lerp(p.da[2], p.a[2], coverage), f32lerp(p.da[3], p.a[3], coverage)
	p.a[4], p.a[5], p.a[6], p.a[7] = f32lerp(p.da[4], p.a[4], coverage), f32lerp(p.da[5], p.a[5], coverage), f32lerp(p.da[6], p.a[6], coverage), f32lerp(p.da[7], p.a[7], coverage)
}

//go:fix inline
func (p *HighPipeline) DestinationAtop() {
	invDa0, invDa1, invDa2, invDa3 := 1.0-p.da[0], 1.0-p.da[1], 1.0-p.da[2], 1.0-p.da[3]
	invDa4, invDa5, invDa6, invDa7 := 1.0-p.da[4], 1.0-p.da[5], 1.0-p.da[6], 1.0-p.da[7]

	p.r[0], p.r[1], p.r[2], p.r[3] = p.dr[0]*p.a[0]+p.r[0]*invDa0, p.dr[1]*p.a[1]+p.r[1]*invDa1, p.dr[2]*p.a[2]+p.r[2]*invDa2, p.dr[3]*p.a[3]+p.r[3]*invDa3
	p.r[4], p.r[5], p.r[6], p.r[7] = p.dr[4]*p.a[4]+p.r[4]*invDa4, p.dr[5]*p.a[5]+p.r[5]*invDa5, p.dr[6]*p.a[6]+p.r[6]*invDa6, p.dr[7]*p.a[7]+p.r[7]*invDa7
	p.g[0], p.g[1], p.g[2], p.g[3] = p.dg[0]*p.a[0]+p.g[0]*invDa0, p.dg[1]*p.a[1]+p.g[1]*invDa1, p.dg[2]*p.a[2]+p.g[2]*invDa2, p.dg[3]*p.a[3]+p.g[3]*invDa3
	p.g[4], p.g[5], p.g[6], p.g[7] = p.dg[4]*p.a[4]+p.g[4]*invDa4, p.dg[5]*p.a[5]+p.g[5]*invDa5, p.dg[6]*p.a[6]+p.g[6]*invDa6, p.dg[7]*p.a[7]+p.g[7]*invDa7
	p.b[0], p.b[1], p.b[2], p.b[3] = p.db[0]*p.a[0]+p.b[0]*invDa0, p.db[1]*p.a[1]+p.b[1]*invDa1, p.db[2]*p.a[2]+p.b[2]*invDa2, p.db[3]*p.a[3]+p.b[3]*invDa3
	p.b[4], p.b[5], p.b[6], p.b[7] = p.db[4]*p.a[4]+p.b[4]*invDa4, p.db[5]*p.a[5]+p.b[5]*invDa5, p.db[6]*p.a[6]+p.b[6]*invDa6, p.db[7]*p.a[7]+p.b[7]*invDa7
	p.a[0], p.a[1], p.a[2], p.a[3] = p.da[0]*p.a[0]+p.a[0]*invDa0, p.a[1]*p.da[1]+p.a[1]*invDa1, p.a[2]*p.da[2]+p.a[2]*invDa2, p.a[3]*p.da[3]+p.a[3]*invDa3
	p.a[4], p.a[5], p.a[6], p.a[7] = p.da[4]*p.a[4]+p.a[4]*invDa4, p.a[5]*p.da[5]+p.a[5]*invDa5, p.a[6]*p.da[6]+p.a[6]*invDa6, p.a[7]*p.da[7]+p.a[7]*invDa7
}

//go:fix inline
func (p *HighPipeline) DestinationIn() {
	p.r[0], p.r[1], p.r[2], p.r[3] = p.dr[0]*p.a[0], p.dr[1]*p.a[1], p.dr[2]*p.a[2], p.dr[3]*p.a[3]
	p.r[4], p.r[5], p.r[6], p.r[7] = p.dr[4]*p.a[4], p.dr[5]*p.a[5], p.dr[6]*p.a[6], p.dr[7]*p.a[7]
	p.g[0], p.g[1], p.g[2], p.g[3] = p.dg[0]*p.a[0], p.dg[1]*p.a[1], p.dg[2]*p.a[2], p.dg[3]*p.a[3]
	p.g[4], p.g[5], p.g[6], p.g[7] = p.dg[4]*p.a[4], p.dg[5]*p.a[5], p.dg[6]*p.a[6], p.dg[7]*p.a[7]
	p.b[0], p.b[1], p.b[2], p.b[3] = p.db[0]*p.a[0], p.db[1]*p.a[1], p.db[2]*p.a[2], p.db[3]*p.a[3]
	p.b[4], p.b[5], p.b[6], p.b[7] = p.db[4]*p.a[4], p.db[5]*p.a[5], p.db[6]*p.a[6], p.db[7]*p.a[7]
	p.a[0], p.a[1], p.a[2], p.a[3] = p.da[0]*p.a[0], p.da[1]*p.a[1], p.da[2]*p.a[2], p.da[3]*p.a[3]
	p.a[4], p.a[5], p.a[6], p.a[7] = p.da[4]*p.a[4], p.da[5]*p.a[5], p.da[6]*p.a[6], p.da[7]*p.a[7]
}

//go:fix inline
func (p *HighPipeline) DestinationOut() {
	invSa0, invSa1, invSa2, invSa3 := 1.0-p.a[0], 1.0-p.a[1], 1.0-p.a[2], 1.0-p.a[3]
	invSa4, invSa5, invSa6, invSa7 := 1.0-p.a[4], 1.0-p.a[5], 1.0-p.a[6], 1.0-p.a[7]

	p.r[0], p.r[1], p.r[2], p.r[3] = p.dr[0]*invSa0, p.dr[1]*invSa1, p.dr[2]*invSa2, p.dr[3]*invSa3
	p.r[4], p.r[5], p.r[6], p.r[7] = p.dr[4]*invSa4, p.dr[5]*invSa5, p.dr[6]*invSa6, p.dr[7]*invSa7
	p.g[0], p.g[1], p.g[2], p.g[3] = p.dg[0]*invSa0, p.dg[1]*invSa1, p.dg[2]*invSa2, p.dg[3]*invSa3
	p.g[4], p.g[5], p.g[6], p.g[7] = p.dg[4]*invSa4, p.dg[5]*invSa5, p.dg[6]*invSa6, p.dg[7]*invSa7
	p.b[0], p.b[1], p.b[2], p.b[3] = p.db[0]*invSa0, p.db[1]*invSa1, p.db[2]*invSa2, p.db[3]*invSa3
	p.b[4], p.b[5], p.b[6], p.b[7] = p.db[4]*invSa4, p.db[5]*invSa5, p.db[6]*invSa6, p.db[7]*invSa7
	p.a[0], p.a[1], p.a[2], p.a[3] = p.da[0]*invSa0, p.da[1]*invSa1, p.da[2]*invSa2, p.da[3]*invSa3
	p.a[4], p.a[5], p.a[6], p.a[7] = p.da[4]*invSa4, p.da[5]*invSa5, p.da[6]*invSa6, p.da[7]*invSa7
}

//go:fix inline
func (p *HighPipeline) DestinationOver() {
	invDa0, invDa1, invDa2, invDa3 := 1.0-p.da[0], 1.0-p.da[1], 1.0-p.da[2], 1.0-p.da[3]
	invDa4, invDa5, invDa6, invDa7 := 1.0-p.da[4], 1.0-p.da[5], 1.0-p.da[6], 1.0-p.da[7]
	p.r[0], p.r[1], p.r[2], p.r[3] = p.r[0]*invDa0+p.dr[0], p.r[1]*invDa1+p.dr[1], p.r[2]*invDa2+p.dr[2], p.r[3]*invDa3+p.dr[3]
	p.r[4], p.r[5], p.r[6], p.r[7] = p.r[4]*invDa4+p.dr[4], p.r[5]*invDa5+p.dr[5], p.r[6]*invDa6+p.dr[6], p.r[7]*invDa7+p.dr[7]
	p.g[0], p.g[1], p.g[2], p.g[3] = p.g[0]*invDa0+p.dg[0], p.g[1]*invDa1+p.dg[1], p.g[2]*invDa2+p.dg[2], p.g[3]*invDa3+p.dg[3]
	p.g[4], p.g[5], p.g[6], p.g[7] = p.g[4]*invDa4+p.dg[4], p.g[5]*invDa5+p.dg[5], p.g[6]*invDa6+p.dg[6], p.g[7]*invDa7+p.dg[7]
	p.b[0], p.b[1], p.b[2], p.b[3] = p.b[0]*invDa0+p.db[0], p.b[1]*invDa1+p.db[1], p.b[2]*invDa2+p.db[2], p.b[3]*invDa3+p.db[3]
	p.b[4], p.b[5], p.b[6], p.b[7] = p.b[4]*invDa4+p.db[4], p.b[5]*invDa5+p.db[5], p.b[6]*invDa6+p.db[6], p.b[7]*invDa7+p.db[7]
	p.a[0], p.a[1], p.a[2], p.a[3] = p.a[0]*invDa0+p.da[0], p.a[1]*invDa1+p.da[1], p.a[2]*invDa2+p.da[2], p.a[3]*invDa3+p.da[3]
	p.a[4], p.a[5], p.a[6], p.a[7] = p.a[4]*invDa4+p.da[4], p.a[5]*invDa5+p.da[5], p.a[6]*invDa6+p.da[6], p.a[7]*invDa7+p.da[7]
}

//go:fix inline
func (p *HighPipeline) SourceAtop() {
	invSa0, invSa1, invSa2, invSa3 := 1.0-p.a[0], 1.0-p.a[1], 1.0-p.a[2], 1.0-p.a[3]
	invSa4, invSa5, invSa6, invSa7 := 1.0-p.a[4], 1.0-p.a[5], 1.0-p.a[6], 1.0-p.a[7]

	p.r[0], p.r[1], p.r[2], p.r[3] = p.r[0]*p.da[0]+p.dr[0]*invSa0, p.r[1]*p.da[1]+p.dr[1]*invSa1, p.r[2]*p.da[2]+p.dr[2]*invSa2, p.r[3]*p.da[3]+p.dr[3]*invSa3
	p.r[4], p.r[5], p.r[6], p.r[7] = p.r[4]*p.da[4]+p.dr[4]*invSa4, p.r[5]*p.da[5]+p.dr[5]*invSa5, p.r[6]*p.da[6]+p.dr[6]*invSa6, p.r[7]*p.da[7]+p.dr[7]*invSa7
	p.g[0], p.g[1], p.g[2], p.g[3] = p.g[0]*p.da[0]+p.dg[0]*invSa0, p.g[1]*p.da[1]+p.dg[1]*invSa1, p.g[2]*p.da[2]+p.dg[2]*invSa2, p.g[3]*p.da[3]+p.dg[3]*invSa3
	p.g[4], p.g[5], p.g[6], p.g[7] = p.g[4]*p.da[4]+p.dg[4]*invSa4, p.g[5]*p.da[5]+p.dg[5]*invSa5, p.g[6]*p.da[6]+p.dg[6]*invSa6, p.g[7]*p.da[7]+p.dg[7]*invSa7
	p.b[0], p.b[1], p.b[2], p.b[3] = p.b[0]*p.da[0]+p.db[0]*invSa0, p.b[1]*p.da[1]+p.db[1]*invSa1, p.b[2]*p.da[2]+p.db[2]*invSa2, p.b[3]*p.da[3]+p.db[3]*invSa3
	p.b[4], p.b[5], p.b[6], p.b[7] = p.b[4]*p.da[4]+p.db[4]*invSa4, p.b[5]*p.da[5]+p.db[5]*invSa5, p.b[6]*p.da[6]+p.db[6]*invSa6, p.b[7]*p.da[7]+p.db[7]*invSa7
	p.a[0], p.a[1], p.a[2], p.a[3] = p.a[0]*p.da[0]+p.da[0]*invSa0, p.a[1]*p.da[1]+p.da[1]*invSa1, p.a[2]*p.da[2]+p.da[2]*invSa2, p.a[3]*p.da[3]+p.da[3]*invSa3
	p.a[4], p.a[5], p.a[6], p.a[7] = p.a[4]*p.da[4]+p.da[4]*invSa4, p.a[5]*p.da[5]+p.da[5]*invSa5, p.a[6]*p.da[6]+p.da[6]*invSa6, p.a[7]*p.da[7]+p.da[7]*invSa7
}

//go:fix inline
func (p *HighPipeline) SourceIn() {
	p.r[0], p.r[1], p.r[2], p.r[3] = p.r[0]*p.da[0], p.r[1]*p.da[1], p.r[2]*p.da[2], p.r[3]*p.da[3]
	p.r[4], p.r[5], p.r[6], p.r[7] = p.r[4]*p.da[4], p.r[5]*p.da[5], p.r[6]*p.da[6], p.r[7]*p.da[7]
	p.g[0], p.g[1], p.g[2], p.g[3] = p.g[0]*p.da[0], p.g[1]*p.da[1], p.g[2]*p.da[2], p.g[3]*p.da[3]
	p.g[4], p.g[5], p.g[6], p.g[7] = p.g[4]*p.da[4], p.g[5]*p.da[5], p.g[6]*p.da[6], p.g[7]*p.da[7]
	p.b[0], p.b[1], p.b[2], p.b[3] = p.b[0]*p.da[0], p.b[1]*p.da[1], p.b[2]*p.da[2], p.b[3]*p.da[3]
	p.b[4], p.b[5], p.b[6], p.b[7] = p.b[4]*p.da[4], p.b[5]*p.da[5], p.b[6]*p.da[6], p.b[7]*p.da[7]
	p.a[0], p.a[1], p.a[2], p.a[3] = p.a[0]*p.da[0], p.a[1]*p.da[1], p.a[2]*p.da[2], p.a[3]*p.da[3]
	p.a[4], p.a[5], p.a[6], p.a[7] = p.a[4]*p.da[4], p.a[5]*p.da[5], p.a[6]*p.da[6], p.a[7]*p.da[7]
}

//go:fix inline
func (p *HighPipeline) SourceOut() {
	invDa0, invDa1, invDa2, invDa3 := 1.0-p.da[0], 1.0-p.da[1], 1.0-p.da[2], 1.0-p.da[3]
	invDa4, invDa5, invDa6, invDa7 := 1.0-p.da[4], 1.0-p.da[5], 1.0-p.da[6], 1.0-p.da[7]

	p.r[0], p.r[1], p.r[2], p.r[3] = p.r[0]*invDa0, p.r[1]*invDa1, p.r[2]*invDa2, p.r[3]*invDa3
	p.r[4], p.r[5], p.r[6], p.r[7] = p.r[4]*invDa4, p.r[5]*invDa5, p.r[6]*invDa6, p.r[7]*invDa7
	p.g[0], p.g[1], p.g[2], p.g[3] = p.g[0]*invDa0, p.g[1]*invDa1, p.g[2]*invDa2, p.g[3]*invDa3
	p.g[4], p.g[5], p.g[6], p.g[7] = p.g[4]*invDa4, p.g[5]*invDa5, p.g[6]*invDa6, p.g[7]*invDa7
	p.b[0], p.b[1], p.b[2], p.b[3] = p.b[0]*invDa0, p.b[1]*invDa1, p.b[2]*invDa2, p.b[3]*invDa3
	p.b[4], p.b[5], p.b[6], p.b[7] = p.b[4]*invDa4, p.b[5]*invDa5, p.b[6]*invDa6, p.b[7]*invDa7
	p.a[0], p.a[1], p.a[2], p.a[3] = p.a[0]*invDa0, p.a[1]*invDa1, p.a[2]*invDa2, p.a[3]*invDa3
	p.a[4], p.a[5], p.a[6], p.a[7] = p.a[4]*invDa4, p.a[5]*invDa5, p.a[6]*invDa6, p.a[7]*invDa7
}

//go:fix inline
func (p *HighPipeline) SourceOver() {
	invSa0, invSa1, invSa2, invSa3 := 1.0-p.a[0], 1.0-p.a[1], 1.0-p.a[2], 1.0-p.a[3]
	invSa4, invSa5, invSa6, invSa7 := 1.0-p.a[4], 1.0-p.a[5], 1.0-p.a[6], 1.0-p.a[7]

	p.r[0], p.r[1], p.r[2], p.r[3] = p.dr[0]*invSa0+p.r[0], p.dr[1]*invSa1+p.r[1], p.dr[2]*invSa2+p.r[2], p.dr[3]*invSa3+p.r[3]
	p.r[4], p.r[5], p.r[6], p.r[7] = p.dr[4]*invSa4+p.r[4], p.dr[5]*invSa5+p.r[5], p.dr[6]*invSa6+p.r[6], p.dr[7]*invSa7+p.r[7]
	p.g[0], p.g[1], p.g[2], p.g[3] = p.dg[0]*invSa0+p.g[0], p.dg[1]*invSa1+p.g[1], p.dg[2]*invSa2+p.g[2], p.dg[3]*invSa3+p.g[3]
	p.g[4], p.g[5], p.g[6], p.g[7] = p.dg[4]*invSa4+p.g[4], p.dg[5]*invSa5+p.g[5], p.dg[6]*invSa6+p.g[6], p.dg[7]*invSa7+p.g[7]
	p.b[0], p.b[1], p.b[2], p.b[3] = p.db[0]*invSa0+p.b[0], p.db[1]*invSa1+p.b[1], p.db[2]*invSa2+p.b[2], p.db[3]*invSa3+p.b[3]
	p.b[4], p.b[5], p.b[6], p.b[7] = p.db[4]*invSa4+p.b[4], p.db[5]*invSa5+p.b[5], p.db[6]*invSa6+p.b[6], p.db[7]*invSa7+p.b[7]
	p.a[0], p.a[1], p.a[2], p.a[3] = p.da[0]*invSa0+p.a[0], p.da[1]*invSa1+p.a[1], p.da[2]*invSa2+p.a[2], p.da[3]*invSa3+p.a[3]
	p.a[4], p.a[5], p.a[6], p.a[7] = p.da[4]*invSa4+p.a[4], p.da[5]*invSa5+p.a[5], p.da[6]*invSa6+p.a[6], p.da[7]*invSa7+p.a[7]
}

//go:fix inline
func (p *HighPipeline) Clear() {
	p.r[0], p.r[1], p.r[2], p.r[3] = 0.0, 0.0, 0.0, 0.0
	p.r[4], p.r[5], p.r[6], p.r[7] = 0.0, 0.0, 0.0, 0.0
	p.g[0], p.g[1], p.g[2], p.g[3] = 0.0, 0.0, 0.0, 0.0
	p.g[4], p.g[5], p.g[6], p.g[7] = 0.0, 0.0, 0.0, 0.0
	p.b[0], p.b[1], p.b[2], p.b[3] = 0.0, 0.0, 0.0, 0.0
	p.b[4], p.b[5], p.b[6], p.b[7] = 0.0, 0.0, 0.0, 0.0
	p.a[0], p.a[1], p.a[2], p.a[3] = 0.0, 0.0, 0.0, 0.0
	p.a[4], p.a[5], p.a[6], p.a[7] = 0.0, 0.0, 0.0, 0.0
}

//go:fix inline
func (p *HighPipeline) Modulate() {
	p.r[0], p.r[1], p.r[2], p.r[3] = p.r[0]*p.dr[0], p.r[1]*p.dr[1], p.r[2]*p.dr[2], p.r[3]*p.dr[3]
	p.r[4], p.r[5], p.r[6], p.r[7] = p.r[4]*p.dr[4], p.r[5]*p.dr[5], p.r[6]*p.dr[6], p.r[7]*p.dr[7]
	p.g[0], p.g[1], p.g[2], p.g[3] = p.g[0]*p.dg[0], p.g[1]*p.dg[1], p.g[2]*p.dg[2], p.g[3]*p.dg[3]
	p.g[4], p.g[5], p.g[6], p.g[7] = p.g[4]*p.dg[4], p.g[5]*p.dg[5], p.g[6]*p.dg[6], p.g[7]*p.dg[7]
	p.b[0], p.b[1], p.b[2], p.b[3] = p.b[0]*p.db[0], p.b[1]*p.db[1], p.b[2]*p.db[2], p.b[3]*p.db[3]
	p.b[4], p.b[5], p.b[6], p.b[7] = p.b[4]*p.db[4], p.b[5]*p.db[5], p.b[6]*p.db[6], p.b[7]*p.db[7]
	p.a[0], p.a[1], p.a[2], p.a[3] = p.a[0]*p.da[0], p.a[1]*p.da[1], p.a[2]*p.da[2], p.a[3]*p.da[3]
	p.a[4], p.a[5], p.a[6], p.a[7] = p.a[4]*p.da[4], p.a[5]*p.da[5], p.a[6]*p.da[6], p.a[7]*p.da[7]
}

//go:fix inline
func (p *HighPipeline) Multiply() {
	invSa0, invSa1, invSa2, invSa3 := 1.0-p.a[0], 1.0-p.a[1], 1.0-p.a[2], 1.0-p.a[3]
	invSa4, invSa5, invSa6, invSa7 := 1.0-p.a[4], 1.0-p.a[5], 1.0-p.a[6], 1.0-p.a[7]
	invDa0, invDa1, invDa2, invDa3 := 1.0-p.da[0], 1.0-p.da[1], 1.0-p.da[2], 1.0-p.da[3]
	invDa4, invDa5, invDa6, invDa7 := 1.0-p.da[4], 1.0-p.da[5], 1.0-p.da[6], 1.0-p.da[7]

	p.r[0], p.r[1], p.r[2], p.r[3] = p.r[0]*invDa0+p.dr[0]*invSa0+p.r[0]*p.dr[0], p.r[1]*invDa1+p.dr[1]*invSa1+p.r[1]*p.dr[1], p.r[2]*invDa2+p.dr[2]*invSa2+p.r[2]*p.dr[2], p.r[3]*invDa3+p.dr[3]*invSa3+p.r[3]*p.dr[3]
	p.r[4], p.r[5], p.r[6], p.r[7] = p.r[4]*invDa4+p.dr[4]*invSa4+p.r[4]*p.dr[4], p.r[5]*invDa5+p.dr[5]*invSa5+p.r[5]*p.dr[5], p.r[6]*invDa6+p.dr[6]*invSa6+p.r[6]*p.dr[6], p.r[7]*invDa7+p.dr[7]*invSa7+p.r[7]*p.dr[7]
	p.g[0], p.g[1], p.g[2], p.g[3] = p.g[0]*invDa0+p.dg[0]*invSa0+p.g[0]*p.dg[0], p.g[1]*invDa1+p.dg[1]*invSa1+p.g[1]*p.dg[1], p.g[2]*invDa2+p.dg[2]*invSa2+p.g[2]*p.dg[2], p.g[3]*invDa3+p.dg[3]*invSa3+p.g[3]*p.dg[3]
	p.g[4], p.g[5], p.g[6], p.g[7] = p.g[4]*invDa4+p.dg[4]*invSa4+p.g[4]*p.dg[4], p.g[5]*invDa5+p.dg[5]*invSa5+p.g[5]*p.dg[5], p.g[6]*invDa6+p.dg[6]*invSa6+p.g[6]*p.dg[6], p.g[7]*invDa7+p.dg[7]*invSa7+p.g[7]*p.dg[7]
	p.b[0], p.b[1], p.b[2], p.b[3] = p.b[0]*invDa0+p.db[0]*invSa0+p.b[0]*p.db[0], p.b[1]*invDa1+p.db[1]*invSa1+p.b[1]*p.db[1], p.b[2]*invDa2+p.db[2]*invSa2+p.b[2]*p.db[2], p.b[3]*invDa3+p.db[3]*invSa3+p.b[3]*p.db[3]
	p.b[4], p.b[5], p.b[6], p.b[7] = p.b[4]*invDa4+p.db[4]*invSa4+p.b[4]*p.db[4], p.b[5]*invDa5+p.db[5]*invSa5+p.b[5]*p.db[5], p.b[6]*invDa6+p.db[6]*invSa6+p.b[6]*p.db[6], p.b[7]*invDa7+p.db[7]*invSa7+p.b[7]*p.db[7]
	p.a[0], p.a[1], p.a[2], p.a[3] = p.a[0]*invDa0+p.da[0]*invSa0+p.a[0]*p.da[0], p.a[1]*invDa1+p.da[1]*invSa1+p.a[1]*p.da[1], p.a[2]*invDa2+p.da[2]*invSa2+p.a[2]*p.da[2], p.a[3]*invDa3+p.da[3]*invSa3+p.a[3]*p.da[3]
	p.a[4], p.a[5], p.a[6], p.a[7] = p.a[4]*invDa4+p.da[4]*invSa4+p.a[4]*p.da[4], p.a[5]*invDa5+p.da[5]*invSa5+p.a[5]*p.da[5], p.a[6]*invDa6+p.da[6]*invSa6+p.a[6]*p.da[6], p.a[7]*invDa7+p.da[7]*invSa7+p.a[7]*p.da[7]
}

//go:fix inline
func (p *HighPipeline) Plus() {
	p.r[0], p.r[1], p.r[2], p.r[3] = f32min(p.r[0]+p.dr[0], 1.0), f32min(p.r[1]+p.dr[1], 1.0), f32min(p.r[2]+p.dr[2], 1.0), f32min(p.r[3]+p.dr[3], 1.0)
	p.r[4], p.r[5], p.r[6], p.r[7] = f32min(p.r[4]+p.dr[4], 1.0), f32min(p.r[5]+p.dr[5], 1.0), f32min(p.r[6]+p.dr[6], 1.0), f32min(p.r[7]+p.dr[7], 1.0)
	p.g[0], p.g[1], p.g[2], p.g[3] = f32min(p.g[0]+p.dg[0], 1.0), f32min(p.g[1]+p.dg[1], 1.0), f32min(p.g[2]+p.dg[2], 1.0), f32min(p.g[3]+p.dg[3], 1.0)
	p.g[4], p.g[5], p.g[6], p.g[7] = f32min(p.g[4]+p.dg[4], 1.0), f32min(p.g[5]+p.dg[5], 1.0), f32min(p.g[6]+p.dg[6], 1.0), f32min(p.g[7]+p.dg[7], 1.0)
	p.b[0], p.b[1], p.b[2], p.b[3] = f32min(p.b[0]+p.db[0], 1.0), f32min(p.b[1]+p.db[1], 1.0), f32min(p.b[2]+p.db[2], 1.0), f32min(p.b[3]+p.db[3], 1.0)
	p.b[4], p.b[5], p.b[6], p.b[7] = f32min(p.b[4]+p.db[4], 1.0), f32min(p.b[5]+p.db[5], 1.0), f32min(p.b[6]+p.db[6], 1.0), f32min(p.b[7]+p.db[7], 1.0)
	p.a[0], p.a[1], p.a[2], p.a[3] = f32min(p.a[0]+p.da[0], 1.0), f32min(p.a[1]+p.da[1], 1.0), f32min(p.a[2]+p.da[2], 1.0), f32min(p.a[3]+p.da[3], 1.0)
	p.a[4], p.a[5], p.a[6], p.a[7] = f32min(p.a[4]+p.da[4], 1.0), f32min(p.a[5]+p.da[5], 1.0), f32min(p.a[6]+p.da[6], 1.0), f32min(p.a[7]+p.da[7], 1.0)
}

//go:fix inline
func (p *HighPipeline) Screen() {
	p.r[0], p.r[1], p.r[2], p.r[3] = p.r[0]+p.dr[0]-p.r[0]*p.dr[0], p.r[1]+p.dr[1]-p.r[1]*p.dr[1], p.r[2]+p.dr[2]-p.r[2]*p.dr[2], p.r[3]+p.dr[3]-p.r[3]*p.dr[3]
	p.r[4], p.r[5], p.r[6], p.r[7] = p.r[4]+p.dr[4]-p.r[4]*p.dr[4], p.r[5]+p.dr[5]-p.r[5]*p.dr[5], p.r[6]+p.dr[6]-p.r[6]*p.dr[6], p.r[7]+p.dr[7]-p.r[7]*p.dr[7]
	p.g[0], p.g[1], p.g[2], p.g[3] = p.g[0]+p.dg[0]-p.g[0]*p.dg[0], p.g[1]+p.dg[1]-p.g[1]*p.dg[1], p.g[2]+p.dg[2]-p.g[2]*p.dg[2], p.g[3]+p.dg[3]-p.g[3]*p.dg[3]
	p.g[4], p.g[5], p.g[6], p.g[7] = p.g[4]+p.dg[4]-p.g[4]*p.dg[4], p.g[5]+p.dg[5]-p.g[5]*p.dg[5], p.g[6]+p.dg[6]-p.g[6]*p.dg[6], p.g[7]+p.dg[7]-p.g[7]*p.dg[7]
	p.b[0], p.b[1], p.b[2], p.b[3] = p.b[0]+p.db[0]-p.b[0]*p.db[0], p.b[1]+p.db[1]-p.b[1]*p.db[1], p.b[2]+p.db[2]-p.b[2]*p.db[2], p.b[3]+p.db[3]-p.b[3]*p.db[3]
	p.b[4], p.b[5], p.b[6], p.b[7] = p.b[4]+p.db[4]-p.b[4]*p.db[4], p.b[5]+p.db[5]-p.b[5]*p.db[5], p.b[6]+p.db[6]-p.b[6]*p.db[6], p.b[7]+p.db[7]-p.b[7]*p.db[7]
	p.a[0], p.a[1], p.a[2], p.a[3] = p.a[0]+p.da[0]-p.a[0]*p.da[0], p.a[1]+p.da[1]-p.a[1]*p.da[1], p.a[2]+p.da[2]-p.a[2]*p.da[2], p.a[3]+p.da[3]-p.a[3]*p.da[3]
	p.a[4], p.a[5], p.a[6], p.a[7] = p.a[4]+p.da[4]-p.a[4]*p.da[4], p.a[5]+p.da[5]-p.a[5]*p.da[5], p.a[6]+p.da[6]-p.a[6]*p.da[6], p.a[7]+p.da[7]-p.a[7]*p.da[7]
}

//go:fix inline
func (p *HighPipeline) Xor() {
	invSa0, invSa1, invSa2, invSa3 := 1.0-p.a[0], 1.0-p.a[1], 1.0-p.a[2], 1.0-p.a[3]
	invSa4, invSa5, invSa6, invSa7 := 1.0-p.a[4], 1.0-p.a[5], 1.0-p.a[6], 1.0-p.a[7]
	invDa0, invDa1, invDa2, invDa3 := 1.0-p.da[0], 1.0-p.da[1], 1.0-p.da[2], 1.0-p.da[3]
	invDa4, invDa5, invDa6, invDa7 := 1.0-p.da[4], 1.0-p.da[5], 1.0-p.da[6], 1.0-p.da[7]

	p.r[0], p.r[1], p.r[2], p.r[3] = p.r[0]*invDa0+p.dr[0]*invSa0, p.r[1]*invDa1+p.dr[1]*invSa1, p.r[2]*invDa2+p.dr[2]*invSa2, p.r[3]*invDa3+p.dr[3]*invSa3
	p.r[4], p.r[5], p.r[6], p.r[7] = p.r[4]*invDa4+p.dr[4]*invSa4, p.r[5]*invDa5+p.dr[5]*invSa5, p.r[6]*invDa6+p.dr[6]*invSa6, p.r[7]*invDa7+p.dr[7]*invSa7
	p.g[0], p.g[1], p.g[2], p.g[3] = p.g[0]*invDa0+p.dg[0]*invSa0, p.g[1]*invDa1+p.dg[1]*invSa1, p.g[2]*invDa2+p.dg[2]*invSa2, p.g[3]*invDa3+p.dg[3]*invSa3
	p.g[4], p.g[5], p.g[6], p.g[7] = p.g[4]*invDa4+p.dg[4]*invSa4, p.g[5]*invDa5+p.dg[5]*invSa5, p.g[6]*invDa6+p.dg[6]*invSa6, p.g[7]*invDa7+p.dg[7]*invSa7
	p.b[0], p.b[1], p.b[2], p.b[3] = p.b[0]*invDa0+p.db[0]*invSa0, p.b[1]*invDa1+p.db[1]*invSa1, p.b[2]*invDa2+p.db[2]*invSa2, p.b[3]*invDa3+p.db[3]*invSa3
	p.b[4], p.b[5], p.b[6], p.b[7] = p.b[4]*invDa4+p.db[4]*invSa4, p.b[5]*invDa5+p.db[5]*invSa5, p.b[6]*invDa6+p.db[6]*invSa6, p.b[7]*invDa7+p.db[7]*invSa7
	p.a[0], p.a[1], p.a[2], p.a[3] = p.a[0]*invDa0+p.da[0]*invSa0, p.a[1]*invDa1+p.da[1]*invSa1, p.a[2]*invDa2+p.da[2]*invSa2, p.a[3]*invDa3+p.da[3]*invSa3
	p.a[4], p.a[5], p.a[6], p.a[7] = p.a[4]*invDa4+p.da[4]*invSa4, p.a[5]*invDa5+p.da[5]*invSa5, p.a[6]*invDa6+p.da[6]*invSa6, p.a[7]*invDa7+p.da[7]*invSa7
}

//go:fix inline
func (p *HighPipeline) ColorBurn() {
	invSa0, invSa1, invSa2, invSa3 := 1.0-p.a[0], 1.0-p.a[1], 1.0-p.a[2], 1.0-p.a[3]
	invSa4, invSa5, invSa6, invSa7 := 1.0-p.a[4], 1.0-p.a[5], 1.0-p.a[6], 1.0-p.a[7]
	invDa0, invDa1, invDa2, invDa3 := 1.0-p.da[0], 1.0-p.da[1], 1.0-p.da[2], 1.0-p.da[3]
	invDa4, invDa5, invDa6, invDa7 := 1.0-p.da[4], 1.0-p.da[5], 1.0-p.da[6], 1.0-p.da[7]

	ch := func(s, d, sa, da, invSa, invDa float32) float32 {
		if d == da {
			return d + s*invDa
		}
		if s == 0 {
			return d * invSa
		}
		return sa*(da-f32min(da, (da-d)*sa*recipFast(s))) + s*invDa + d*invSa
	}

	p.r[0] = ch(p.r[0], p.dr[0], p.a[0], p.da[0], invSa0, invDa0)
	p.g[0] = ch(p.g[0], p.dg[0], p.a[0], p.da[0], invSa0, invDa0)
	p.b[0] = ch(p.b[0], p.db[0], p.a[0], p.da[0], invSa0, invDa0)
	p.r[1] = ch(p.r[1], p.dr[1], p.a[1], p.da[1], invSa1, invDa1)
	p.g[1] = ch(p.g[1], p.dg[1], p.a[1], p.da[1], invSa1, invDa1)
	p.b[1] = ch(p.b[1], p.db[1], p.a[1], p.da[1], invSa1, invDa1)
	p.r[2] = ch(p.r[2], p.dr[2], p.a[2], p.da[2], invSa2, invDa2)
	p.g[2] = ch(p.g[2], p.dg[2], p.a[2], p.da[2], invSa2, invDa2)
	p.b[2] = ch(p.b[2], p.db[2], p.a[2], p.da[2], invSa2, invDa2)
	p.r[3] = ch(p.r[3], p.dr[3], p.a[3], p.da[3], invSa3, invDa3)
	p.g[3] = ch(p.g[3], p.dg[3], p.a[3], p.da[3], invSa3, invDa3)
	p.b[3] = ch(p.b[3], p.db[3], p.a[3], p.da[3], invSa3, invDa3)
	p.r[4] = ch(p.r[4], p.dr[4], p.a[4], p.da[4], invSa4, invDa4)
	p.g[4] = ch(p.g[4], p.dg[4], p.a[4], p.da[4], invSa4, invDa4)
	p.b[4] = ch(p.b[4], p.db[4], p.a[4], p.da[4], invSa4, invDa4)
	p.r[5] = ch(p.r[5], p.dr[5], p.a[5], p.da[5], invSa5, invDa5)
	p.g[5] = ch(p.g[5], p.dg[5], p.a[5], p.da[5], invSa5, invDa5)
	p.b[5] = ch(p.b[5], p.db[5], p.a[5], p.da[5], invSa5, invDa5)
	p.r[6] = ch(p.r[6], p.dr[6], p.a[6], p.da[6], invSa6, invDa6)
	p.g[6] = ch(p.g[6], p.dg[6], p.a[6], p.da[6], invSa6, invDa6)
	p.b[6] = ch(p.b[6], p.db[6], p.a[6], p.da[6], invSa6, invDa6)
	p.r[7] = ch(p.r[7], p.dr[7], p.a[7], p.da[7], invSa7, invDa7)
	p.g[7] = ch(p.g[7], p.dg[7], p.a[7], p.da[7], invSa7, invDa7)
	p.b[7] = ch(p.b[7], p.db[7], p.a[7], p.da[7], invSa7, invDa7)
	p.a[0], p.a[1], p.a[2], p.a[3] = p.a[0]+p.da[0]*invSa0, p.a[1]+p.da[1]*invSa1, p.a[2]+p.da[2]*invSa2, p.a[3]+p.da[3]*invSa3
	p.a[4], p.a[5], p.a[6], p.a[7] = p.a[4]+p.da[4]*invSa4, p.a[5]+p.da[5]*invSa5, p.a[6]+p.da[6]*invSa6, p.a[7]+p.da[7]*invSa7
}

//go:fix inline
func (p *HighPipeline) ColorDodge() {
	invSa0, invSa1, invSa2, invSa3 := 1.0-p.a[0], 1.0-p.a[1], 1.0-p.a[2], 1.0-p.a[3]
	invSa4, invSa5, invSa6, invSa7 := 1.0-p.a[4], 1.0-p.a[5], 1.0-p.a[6], 1.0-p.a[7]
	invDa0, invDa1, invDa2, invDa3 := 1.0-p.da[0], 1.0-p.da[1], 1.0-p.da[2], 1.0-p.da[3]
	invDa4, invDa5, invDa6, invDa7 := 1.0-p.da[4], 1.0-p.da[5], 1.0-p.da[6], 1.0-p.da[7]

	ch := func(s, d, sa, da, invSa, invDa float32) float32 {
		if d == 0 {
			return s * invDa
		}
		if s == sa {
			return s + d*invSa
		}
		return sa*f32min(da, (d*sa)*recipFast(sa-s)) + s*invDa + d*invSa
	}

	p.r[0] = ch(p.r[0], p.dr[0], p.a[0], p.da[0], invSa0, invDa0)
	p.g[0] = ch(p.g[0], p.dg[0], p.a[0], p.da[0], invSa0, invDa0)
	p.b[0] = ch(p.b[0], p.db[0], p.a[0], p.da[0], invSa0, invDa0)
	p.r[1] = ch(p.r[1], p.dr[1], p.a[1], p.da[1], invSa1, invDa1)
	p.g[1] = ch(p.g[1], p.dg[1], p.a[1], p.da[1], invSa1, invDa1)
	p.b[1] = ch(p.b[1], p.db[1], p.a[1], p.da[1], invSa1, invDa1)
	p.r[2] = ch(p.r[2], p.dr[2], p.a[2], p.da[2], invSa2, invDa2)
	p.g[2] = ch(p.g[2], p.dg[2], p.a[2], p.da[2], invSa2, invDa2)
	p.b[2] = ch(p.b[2], p.db[2], p.a[2], p.da[2], invSa2, invDa2)
	p.r[3] = ch(p.r[3], p.dr[3], p.a[3], p.da[3], invSa3, invDa3)
	p.g[3] = ch(p.g[3], p.dg[3], p.a[3], p.da[3], invSa3, invDa3)
	p.b[3] = ch(p.b[3], p.db[3], p.a[3], p.da[3], invSa3, invDa3)
	p.r[4] = ch(p.r[4], p.dr[4], p.a[4], p.da[4], invSa4, invDa4)
	p.g[4] = ch(p.g[4], p.dg[4], p.a[4], p.da[4], invSa4, invDa4)
	p.b[4] = ch(p.b[4], p.db[4], p.a[4], p.da[4], invSa4, invDa4)
	p.r[5] = ch(p.r[5], p.dr[5], p.a[5], p.da[5], invSa5, invDa5)
	p.g[5] = ch(p.g[5], p.dg[5], p.a[5], p.da[5], invSa5, invDa5)
	p.b[5] = ch(p.b[5], p.db[5], p.a[5], p.da[5], invSa5, invDa5)
	p.r[6] = ch(p.r[6], p.dr[6], p.a[6], p.da[6], invSa6, invDa6)
	p.g[6] = ch(p.g[6], p.dg[6], p.a[6], p.da[6], invSa6, invDa6)
	p.b[6] = ch(p.b[6], p.db[6], p.a[6], p.da[6], invSa6, invDa6)
	p.r[7] = ch(p.r[7], p.dr[7], p.a[7], p.da[7], invSa7, invDa7)
	p.g[7] = ch(p.g[7], p.dg[7], p.a[7], p.da[7], invSa7, invDa7)
	p.b[7] = ch(p.b[7], p.db[7], p.a[7], p.da[7], invSa7, invDa7)
	p.a[0], p.a[1], p.a[2], p.a[3] = p.a[0]+p.da[0]*invSa0, p.a[1]+p.da[1]*invSa1, p.a[2]+p.da[2]*invSa2, p.a[3]+p.da[3]*invSa3
	p.a[4], p.a[5], p.a[6], p.a[7] = p.a[4]+p.da[4]*invSa4, p.a[5]+p.da[5]*invSa5, p.a[6]+p.da[6]*invSa6, p.a[7]+p.da[7]*invSa7
}

//go:fix inline
func (p *HighPipeline) Darken() {
	invSa0, invSa1, invSa2, invSa3 := 1.0-p.a[0], 1.0-p.a[1], 1.0-p.a[2], 1.0-p.a[3]
	invSa4, invSa5, invSa6, invSa7 := 1.0-p.a[4], 1.0-p.a[5], 1.0-p.a[6], 1.0-p.a[7]

	p.r[0], p.r[1], p.r[2], p.r[3] = p.r[0]+p.dr[0]-f32max(p.r[0]*p.da[0], p.dr[0]*p.a[0]), p.r[1]+p.dr[1]-f32max(p.r[1]*p.da[1], p.dr[1]*p.a[1]),
		p.r[2]+p.dr[2]-f32max(p.r[2]*p.da[2], p.dr[2]*p.a[2]), p.r[3]+p.dr[3]-f32max(p.r[3]*p.da[3], p.dr[3]*p.a[3])
	p.r[4], p.r[5], p.r[6], p.r[7] = p.r[4]+p.dr[4]-f32max(p.r[4]*p.da[4], p.dr[4]*p.a[4]), p.r[5]+p.dr[5]-f32max(p.r[5]*p.da[5], p.dr[5]*p.a[5]),
		p.r[6]+p.dr[6]-f32max(p.r[6]*p.da[6], p.dr[6]*p.a[6]), p.r[7]+p.dr[7]-f32max(p.r[7]*p.da[7], p.dr[7]*p.a[7])
	p.g[0], p.g[1], p.g[2], p.g[3] = p.g[0]+p.dg[0]-f32max(p.g[0]*p.da[0], p.dg[0]*p.a[0]), p.g[1]+p.dg[1]-f32max(p.g[1]*p.da[1], p.dg[1]*p.a[1]),
		p.g[2]+p.dg[2]-f32max(p.g[2]*p.da[2], p.dg[2]*p.a[2]), p.g[3]+p.dg[3]-f32max(p.g[3]*p.da[3], p.dg[3]*p.a[3])
	p.g[4], p.g[5], p.g[6], p.g[7] = p.g[4]+p.dg[4]-f32max(p.g[4]*p.da[4], p.dg[4]*p.a[4]), p.g[5]+p.dg[5]-f32max(p.g[5]*p.da[5], p.dg[5]*p.a[5]),
		p.g[6]+p.dg[6]-f32max(p.g[6]*p.da[6], p.dg[6]*p.a[6]), p.g[7]+p.dg[7]-f32max(p.g[7]*p.da[7], p.dg[7]*p.a[7])
	p.b[0], p.b[1], p.b[2], p.b[3] = p.b[0]+p.db[0]-f32max(p.b[0]*p.da[0], p.db[0]*p.a[0]), p.b[1]+p.db[1]-f32max(p.b[1]*p.da[1], p.db[1]*p.a[1]),
		p.b[2]+p.db[2]-f32max(p.b[2]*p.da[2], p.db[2]*p.a[2]), p.b[3]+p.db[3]-f32max(p.b[3]*p.da[3], p.db[3]*p.a[3])
	p.b[4], p.b[5], p.b[6], p.b[7] = p.b[4]+p.db[4]-f32max(p.b[4]*p.da[4], p.db[4]*p.a[4]), p.b[5]+p.db[5]-f32max(p.b[5]*p.da[5], p.db[5]*p.a[5]),
		p.b[6]+p.db[6]-f32max(p.b[6]*p.da[6], p.db[6]*p.a[6]), p.b[7]+p.db[7]-f32max(p.b[7]*p.da[7], p.db[7]*p.a[7])
	p.a[0], p.a[1], p.a[2], p.a[3] = p.a[0]+p.da[0]*invSa0, p.a[1]+p.da[1]*invSa1, p.a[2]+p.da[2]*invSa2, p.a[3]+p.da[3]*invSa3
	p.a[4], p.a[5], p.a[6], p.a[7] = p.a[4]+p.da[4]*invSa4, p.a[5]+p.da[5]*invSa5, p.a[6]+p.da[6]*invSa6, p.a[7]+p.da[7]*invSa7

}

//go:fix inline
func (p *HighPipeline) Difference() {
	invSa0, invSa1, invSa2, invSa3 := 1.0-p.a[0], 1.0-p.a[1], 1.0-p.a[2], 1.0-p.a[3]
	invSa4, invSa5, invSa6, invSa7 := 1.0-p.a[4], 1.0-p.a[5], 1.0-p.a[6], 1.0-p.a[7]

	p.r[0], p.r[1], p.r[2], p.r[3] = p.r[0]+p.dr[0]-2*f32min(p.r[0]*p.da[0], p.dr[0]*p.a[0]), p.r[1]+p.dr[1]-2*f32min(p.r[1]*p.da[1], p.dr[1]*p.a[1]),
		p.r[2]+p.dr[2]-2*f32min(p.r[2]*p.da[2], p.dr[2]*p.a[2]), p.r[3]+p.dr[3]-2*f32min(p.r[3]*p.da[3], p.dr[3]*p.a[3])
	p.r[4], p.r[5], p.r[6], p.r[7] = p.r[4]+p.dr[4]-2*f32min(p.r[4]*p.da[4], p.dr[4]*p.a[4]), p.r[5]+p.dr[5]-2*f32min(p.r[5]*p.da[5], p.dr[5]*p.a[5]),
		p.r[6]+p.dr[6]-2*f32min(p.r[6]*p.da[6], p.dr[6]*p.a[6]), p.r[7]+p.dr[7]-2*f32min(p.r[7]*p.da[7], p.dr[7]*p.a[7])
	p.g[0], p.g[1], p.g[2], p.g[3] = p.g[0]+p.dg[0]-2*f32min(p.g[0]*p.da[0], p.dg[0]*p.a[0]), p.g[1]+p.dg[1]-2*f32min(p.g[1]*p.da[1], p.dg[1]*p.a[1]),
		p.g[2]+p.dg[2]-2*f32min(p.g[2]*p.da[2], p.dg[2]*p.a[2]), p.g[3]+p.dg[3]-2*f32min(p.g[3]*p.da[3], p.dg[3]*p.a[3])
	p.g[4], p.g[5], p.g[6], p.g[7] = p.g[4]+p.dg[4]-2*f32min(p.g[4]*p.da[4], p.dg[4]*p.a[4]), p.g[5]+p.dg[5]-2*f32min(p.g[5]*p.da[5], p.dg[5]*p.a[5]),
		p.g[6]+p.dg[6]-2*f32min(p.g[6]*p.da[6], p.dg[6]*p.a[6]), p.g[7]+p.dg[7]-2*f32min(p.g[7]*p.da[7], p.dg[7]*p.a[7])
	p.b[0], p.b[1], p.b[2], p.b[3] = p.b[0]+p.db[0]-2*f32min(p.b[0]*p.da[0], p.db[0]*p.a[0]), p.b[1]+p.db[1]-2*f32min(p.b[1]*p.da[1], p.db[1]*p.a[1]),
		p.b[2]+p.db[2]-2*f32min(p.b[2]*p.da[2], p.db[2]*p.a[2]), p.b[3]+p.db[3]-2*f32min(p.b[3]*p.da[3], p.db[3]*p.a[3])
	p.b[4], p.b[5], p.b[6], p.b[7] = p.b[4]+p.db[4]-2*f32min(p.b[4]*p.da[4], p.db[4]*p.a[4]), p.b[5]+p.db[5]-2*f32min(p.b[5]*p.da[5], p.db[5]*p.a[5]),
		p.b[6]+p.db[6]-2*f32min(p.b[6]*p.da[6], p.db[6]*p.a[6]), p.b[7]+p.db[7]-2*f32min(p.b[7]*p.da[7], p.db[7]*p.a[7])
	p.a[0], p.a[1], p.a[2], p.a[3] = p.a[0]+p.da[0]*invSa0, p.a[1]+p.da[1]*invSa1, p.a[2]+p.da[2]*invSa2, p.a[3]+p.da[3]*invSa3
	p.a[4], p.a[5], p.a[6], p.a[7] = p.a[4]+p.da[4]*invSa4, p.a[5]+p.da[5]*invSa5, p.a[6]+p.da[6]*invSa6, p.a[7]+p.da[7]*invSa7

}

//go:fix inline
func (p *HighPipeline) Exclusion() {
	invSa0, invSa1, invSa2, invSa3 := 1.0-p.a[0], 1.0-p.a[1], 1.0-p.a[2], 1.0-p.a[3]
	invSa4, invSa5, invSa6, invSa7 := 1.0-p.a[4], 1.0-p.a[5], 1.0-p.a[6], 1.0-p.a[7]

	p.r[0], p.r[1], p.r[2], p.r[3] = p.r[0]+p.dr[0]-2*p.r[0]*p.dr[0], p.r[1]+p.dr[1]-2*p.r[1]*p.dr[1], p.r[2]+p.dr[2]-2*p.r[2]*p.dr[2], p.r[3]+p.dr[3]-2*p.r[3]*p.dr[3]
	p.r[4], p.r[5], p.r[6], p.r[7] = p.r[4]+p.dr[4]-2*p.r[4]*p.dr[4], p.r[5]+p.dr[5]-2*p.r[5]*p.dr[5], p.r[6]+p.dr[6]-2*p.r[6]*p.dr[6], p.r[7]+p.dr[7]-2*p.r[7]*p.dr[7]
	p.g[0], p.g[1], p.g[2], p.g[3] = p.g[0]+p.dg[0]-2*p.g[0]*p.dg[0], p.g[1]+p.dg[1]-2*p.g[1]*p.dg[1], p.g[2]+p.dg[2]-2*p.g[2]*p.dg[2], p.g[3]+p.dg[3]-2*p.g[3]*p.dg[3]
	p.g[4], p.g[5], p.g[6], p.g[7] = p.g[4]+p.dg[4]-2*p.g[4]*p.dg[4], p.g[5]+p.dg[5]-2*p.g[5]*p.dg[5], p.g[6]+p.dg[6]-2*p.g[6]*p.dg[6], p.g[7]+p.dg[7]-2*p.g[7]*p.dg[7]
	p.b[0], p.b[1], p.b[2], p.b[3] = p.b[0]+p.db[0]-2*p.b[0]*p.db[0], p.b[1]+p.db[1]-2*p.b[1]*p.db[1], p.b[2]+p.db[2]-2*p.b[2]*p.db[2], p.b[3]+p.db[3]-2*p.b[3]*p.db[3]
	p.b[4], p.b[5], p.b[6], p.b[7] = p.b[4]+p.db[4]-2*p.b[4]*p.db[4], p.b[5]+p.db[5]-2*p.b[5]*p.db[5], p.b[6]+p.db[6]-2*p.b[6]*p.db[6], p.b[7]+p.db[7]-2*p.b[7]*p.db[7]
	p.a[0], p.a[1], p.a[2], p.a[3] = p.a[0]+p.da[0]*invSa0, p.a[1]+p.da[1]*invSa1, p.a[2]+p.da[2]*invSa2, p.a[3]+p.da[3]*invSa3
	p.a[4], p.a[5], p.a[6], p.a[7] = p.a[4]+p.da[4]*invSa4, p.a[5]+p.da[5]*invSa5, p.a[6]+p.da[6]*invSa6, p.a[7]+p.da[7]*invSa7

}

//go:fix inline
func (p *HighPipeline) HardLight() {
	invSa0, invSa1, invSa2, invSa3 := 1.0-p.a[0], 1.0-p.a[1], 1.0-p.a[2], 1.0-p.a[3]
	invSa4, invSa5, invSa6, invSa7 := 1.0-p.a[4], 1.0-p.a[5], 1.0-p.a[6], 1.0-p.a[7]
	invDa0, invDa1, invDa2, invDa3 := 1.0-p.da[0], 1.0-p.da[1], 1.0-p.da[2], 1.0-p.da[3]
	invDa4, invDa5, invDa6, invDa7 := 1.0-p.da[4], 1.0-p.da[5], 1.0-p.da[6], 1.0-p.da[7]

	ch := func(s, d, sa, da, invSa, invDa float32) float32 {
		if 2*s <= sa {
			return s*invDa + d*invSa + 2*s*d
		}
		return s*invDa + d*invSa + sa*da - 2*(da-d)*(sa-s)
	}

	p.r[0] = ch(p.r[0], p.dr[0], p.a[0], p.da[0], invSa0, invDa0)
	p.g[0] = ch(p.g[0], p.dg[0], p.a[0], p.da[0], invSa0, invDa0)
	p.b[0] = ch(p.b[0], p.db[0], p.a[0], p.da[0], invSa0, invDa0)
	p.r[1] = ch(p.r[1], p.dr[1], p.a[1], p.da[1], invSa1, invDa1)
	p.g[1] = ch(p.g[1], p.dg[1], p.a[1], p.da[1], invSa1, invDa1)
	p.b[1] = ch(p.b[1], p.db[1], p.a[1], p.da[1], invSa1, invDa1)
	p.r[2] = ch(p.r[2], p.dr[2], p.a[2], p.da[2], invSa2, invDa2)
	p.g[2] = ch(p.g[2], p.dg[2], p.a[2], p.da[2], invSa2, invDa2)
	p.b[2] = ch(p.b[2], p.db[2], p.a[2], p.da[2], invSa2, invDa2)
	p.r[3] = ch(p.r[3], p.dr[3], p.a[3], p.da[3], invSa3, invDa3)
	p.g[3] = ch(p.g[3], p.dg[3], p.a[3], p.da[3], invSa3, invDa3)
	p.b[3] = ch(p.b[3], p.db[3], p.a[3], p.da[3], invSa3, invDa3)
	p.r[4] = ch(p.r[4], p.dr[4], p.a[4], p.da[4], invSa4, invDa4)
	p.g[4] = ch(p.g[4], p.dg[4], p.a[4], p.da[4], invSa4, invDa4)
	p.b[4] = ch(p.b[4], p.db[4], p.a[4], p.da[4], invSa4, invDa4)
	p.r[5] = ch(p.r[5], p.dr[5], p.a[5], p.da[5], invSa5, invDa5)
	p.g[5] = ch(p.g[5], p.dg[5], p.a[5], p.da[5], invSa5, invDa5)
	p.b[5] = ch(p.b[5], p.db[5], p.a[5], p.da[5], invSa5, invDa5)
	p.r[6] = ch(p.r[6], p.dr[6], p.a[6], p.da[6], invSa6, invDa6)
	p.g[6] = ch(p.g[6], p.dg[6], p.a[6], p.da[6], invSa6, invDa6)
	p.b[6] = ch(p.b[6], p.db[6], p.a[6], p.da[6], invSa6, invDa6)
	p.r[7] = ch(p.r[7], p.dr[7], p.a[7], p.da[7], invSa7, invDa7)
	p.g[7] = ch(p.g[7], p.dg[7], p.a[7], p.da[7], invSa7, invDa7)
	p.b[7] = ch(p.b[7], p.db[7], p.a[7], p.da[7], invSa7, invDa7)
	p.a[0], p.a[1], p.a[2], p.a[3] = p.a[0]+p.da[0]*invSa0, p.a[1]+p.da[1]*invSa1, p.a[2]+p.da[2]*invSa2, p.a[3]+p.da[3]*invSa3
	p.a[4], p.a[5], p.a[6], p.a[7] = p.a[4]+p.da[4]*invSa4, p.a[5]+p.da[5]*invSa5, p.a[6]+p.da[6]*invSa6, p.a[7]+p.da[7]*invSa7

}

//go:fix inline
func (p *HighPipeline) Lighten() {
	invSa0, invSa1, invSa2, invSa3 := 1.0-p.a[0], 1.0-p.a[1], 1.0-p.a[2], 1.0-p.a[3]
	invSa4, invSa5, invSa6, invSa7 := 1.0-p.a[4], 1.0-p.a[5], 1.0-p.a[6], 1.0-p.a[7]

	p.r[0], p.r[1], p.r[2], p.r[3] = p.r[0]+p.dr[0]-f32min(p.r[0]*p.da[0], p.dr[0]*p.a[0]), p.r[1]+p.dr[1]-f32min(p.r[1]*p.da[1], p.dr[1]*p.a[1]),
		p.r[2]+p.dr[2]-f32min(p.r[2]*p.da[2], p.dr[2]*p.a[2]), p.r[3]+p.dr[3]-f32min(p.r[3]*p.da[3], p.dr[3]*p.a[3])
	p.r[4], p.r[5], p.r[6], p.r[7] = p.r[4]+p.dr[4]-f32min(p.r[4]*p.da[4], p.dr[4]*p.a[4]), p.r[5]+p.dr[5]-f32min(p.r[5]*p.da[5], p.dr[5]*p.a[5]),
		p.r[6]+p.dr[6]-f32min(p.r[6]*p.da[6], p.dr[6]*p.a[6]), p.r[7]+p.dr[7]-f32min(p.r[7]*p.da[7], p.dr[7]*p.a[7])
	p.g[0], p.g[1], p.g[2], p.g[3] = p.g[0]+p.dg[0]-f32min(p.g[0]*p.da[0], p.dg[0]*p.a[0]), p.g[1]+p.dg[1]-f32min(p.g[1]*p.da[1], p.dg[1]*p.a[1]),
		p.g[2]+p.dg[2]-f32min(p.g[2]*p.da[2], p.dg[2]*p.a[2]), p.g[3]+p.dg[3]-f32min(p.g[3]*p.da[3], p.dg[3]*p.a[3])
	p.g[4], p.g[5], p.g[6], p.g[7] = p.g[4]+p.dg[4]-f32min(p.g[4]*p.da[4], p.dg[4]*p.a[4]), p.g[5]+p.dg[5]-f32min(p.g[5]*p.da[5], p.dg[5]*p.a[5]),
		p.g[6]+p.dg[6]-f32min(p.g[6]*p.da[6], p.dg[6]*p.a[6]), p.g[7]+p.dg[7]-f32min(p.g[7]*p.da[7], p.dg[7]*p.a[7])
	p.b[0], p.b[1], p.b[2], p.b[3] = p.b[0]+p.db[0]-f32min(p.b[0]*p.da[0], p.db[0]*p.a[0]), p.b[1]+p.db[1]-f32min(p.b[1]*p.da[1], p.db[1]*p.a[1]),
		p.b[2]+p.db[2]-f32min(p.b[2]*p.da[2], p.db[2]*p.a[2]), p.b[3]+p.db[3]-f32min(p.b[3]*p.da[3], p.db[3]*p.a[3])
	p.b[4], p.b[5], p.b[6], p.b[7] = p.b[4]+p.db[4]-f32min(p.b[4]*p.da[4], p.db[4]*p.a[4]), p.b[5]+p.db[5]-f32min(p.b[5]*p.da[5], p.db[5]*p.a[5]),
		p.b[6]+p.db[6]-f32min(p.b[6]*p.da[6], p.db[6]*p.a[6]), p.b[7]+p.db[7]-f32min(p.b[7]*p.da[7], p.db[7]*p.a[7])
	p.a[0], p.a[1], p.a[2], p.a[3] = p.a[0]+p.da[0]*invSa0, p.a[1]+p.da[1]*invSa1, p.a[2]+p.da[2]*invSa2, p.a[3]+p.da[3]*invSa3
	p.a[4], p.a[5], p.a[6], p.a[7] = p.a[4]+p.da[4]*invSa4, p.a[5]+p.da[5]*invSa5, p.a[6]+p.da[6]*invSa6, p.a[7]+p.da[7]*invSa7

}

//go:fix inline
func (p *HighPipeline) Overlay() {
	invSa0, invSa1, invSa2, invSa3 := 1.0-p.a[0], 1.0-p.a[1], 1.0-p.a[2], 1.0-p.a[3]
	invSa4, invSa5, invSa6, invSa7 := 1.0-p.a[4], 1.0-p.a[5], 1.0-p.a[6], 1.0-p.a[7]
	invDa0, invDa1, invDa2, invDa3 := 1.0-p.da[0], 1.0-p.da[1], 1.0-p.da[2], 1.0-p.da[3]
	invDa4, invDa5, invDa6, invDa7 := 1.0-p.da[4], 1.0-p.da[5], 1.0-p.da[6], 1.0-p.da[7]

	ch := func(s, d, sa, da, invSa, invDa float32) float32 {
		if 2*d <= da {
			return s*invDa + d*invSa + 2*s*d
		}
		return s*invDa + d*invSa + sa*da - 2*(da-d)*(sa-s)
	}

	p.r[0] = ch(p.r[0], p.dr[0], p.a[0], p.da[0], invSa0, invDa0)
	p.g[0] = ch(p.g[0], p.dg[0], p.a[0], p.da[0], invSa0, invDa0)
	p.b[0] = ch(p.b[0], p.db[0], p.a[0], p.da[0], invSa0, invDa0)
	p.r[1] = ch(p.r[1], p.dr[1], p.a[1], p.da[1], invSa1, invDa1)
	p.g[1] = ch(p.g[1], p.dg[1], p.a[1], p.da[1], invSa1, invDa1)
	p.b[1] = ch(p.b[1], p.db[1], p.a[1], p.da[1], invSa1, invDa1)
	p.r[2] = ch(p.r[2], p.dr[2], p.a[2], p.da[2], invSa2, invDa2)
	p.g[2] = ch(p.g[2], p.dg[2], p.a[2], p.da[2], invSa2, invDa2)
	p.b[2] = ch(p.b[2], p.db[2], p.a[2], p.da[2], invSa2, invDa2)
	p.r[3] = ch(p.r[3], p.dr[3], p.a[3], p.da[3], invSa3, invDa3)
	p.g[3] = ch(p.g[3], p.dg[3], p.a[3], p.da[3], invSa3, invDa3)
	p.b[3] = ch(p.b[3], p.db[3], p.a[3], p.da[3], invSa3, invDa3)
	p.r[4] = ch(p.r[4], p.dr[4], p.a[4], p.da[4], invSa4, invDa4)
	p.g[4] = ch(p.g[4], p.dg[4], p.a[4], p.da[4], invSa4, invDa4)
	p.b[4] = ch(p.b[4], p.db[4], p.a[4], p.da[4], invSa4, invDa4)
	p.r[5] = ch(p.r[5], p.dr[5], p.a[5], p.da[5], invSa5, invDa5)
	p.g[5] = ch(p.g[5], p.dg[5], p.a[5], p.da[5], invSa5, invDa5)
	p.b[5] = ch(p.b[5], p.db[5], p.a[5], p.da[5], invSa5, invDa5)
	p.r[6] = ch(p.r[6], p.dr[6], p.a[6], p.da[6], invSa6, invDa6)
	p.g[6] = ch(p.g[6], p.dg[6], p.a[6], p.da[6], invSa6, invDa6)
	p.b[6] = ch(p.b[6], p.db[6], p.a[6], p.da[6], invSa6, invDa6)
	p.r[7] = ch(p.r[7], p.dr[7], p.a[7], p.da[7], invSa7, invDa7)
	p.g[7] = ch(p.g[7], p.dg[7], p.a[7], p.da[7], invSa7, invDa7)
	p.b[7] = ch(p.b[7], p.db[7], p.a[7], p.da[7], invSa7, invDa7)
	p.a[0], p.a[1], p.a[2], p.a[3] = p.a[0]+p.da[0]*invSa0, p.a[1]+p.da[1]*invSa1, p.a[2]+p.da[2]*invSa2, p.a[3]+p.da[3]*invSa3
	p.a[4], p.a[5], p.a[6], p.a[7] = p.a[4]+p.da[4]*invSa4, p.a[5]+p.da[5]*invSa5, p.a[6]+p.da[6]*invSa6, p.a[7]+p.da[7]*invSa7

}

//go:fix inline
func (p *HighPipeline) SoftLight() {
	invSa0, invSa1, invSa2, invSa3 := 1.0-p.a[0], 1.0-p.a[1], 1.0-p.a[2], 1.0-p.a[3]
	invSa4, invSa5, invSa6, invSa7 := 1.0-p.a[4], 1.0-p.a[5], 1.0-p.a[6], 1.0-p.a[7]
	invDa0, invDa1, invDa2, invDa3 := 1.0-p.da[0], 1.0-p.da[1], 1.0-p.da[2], 1.0-p.da[3]
	invDa4, invDa5, invDa6, invDa7 := 1.0-p.da[4], 1.0-p.da[5], 1.0-p.da[6], 1.0-p.da[7]

	ch := func(s, d, sa, da, invSa, invDa float32) float32 {
		m := float32(0)
		if da > 0 {
			m = d / da
		}
		s2 := 2 * s
		if s2 <= sa {
			// dark_src: d * (sa + (s2 - sa) * (1 - m))
			return s*invDa + d*invSa + d*(sa+(s2-sa)*(1-m))
		}
		// lite_src: d*sa + da*(s2-sa) * blend(dark_dst, lite_dst)
		m4 := 4 * m
		darkDst := (m4*m4+m4)*(m-1) + 7*m
		liteDst := math32.Sqrt(m) - m
		var selectedDst float32
		if 4*d <= da { // two(two(d)) <= da => 4*d <= da
			selectedDst = darkDst
		} else {
			selectedDst = liteDst
		}
		return s*invDa + d*invSa + d*sa + da*(s2-sa)*selectedDst
	}

	p.r[0] = ch(p.r[0], p.dr[0], p.a[0], p.da[0], invSa0, invDa0)
	p.g[0] = ch(p.g[0], p.dg[0], p.a[0], p.da[0], invSa0, invDa0)
	p.b[0] = ch(p.b[0], p.db[0], p.a[0], p.da[0], invSa0, invDa0)
	p.r[1] = ch(p.r[1], p.dr[1], p.a[1], p.da[1], invSa1, invDa1)
	p.g[1] = ch(p.g[1], p.dg[1], p.a[1], p.da[1], invSa1, invDa1)
	p.b[1] = ch(p.b[1], p.db[1], p.a[1], p.da[1], invSa1, invDa1)
	p.r[2] = ch(p.r[2], p.dr[2], p.a[2], p.da[2], invSa2, invDa2)
	p.g[2] = ch(p.g[2], p.dg[2], p.a[2], p.da[2], invSa2, invDa2)
	p.b[2] = ch(p.b[2], p.db[2], p.a[2], p.da[2], invSa2, invDa2)
	p.r[3] = ch(p.r[3], p.dr[3], p.a[3], p.da[3], invSa3, invDa3)
	p.g[3] = ch(p.g[3], p.dg[3], p.a[3], p.da[3], invSa3, invDa3)
	p.b[3] = ch(p.b[3], p.db[3], p.a[3], p.da[3], invSa3, invDa3)
	p.r[4] = ch(p.r[4], p.dr[4], p.a[4], p.da[4], invSa4, invDa4)
	p.g[4] = ch(p.g[4], p.dg[4], p.a[4], p.da[4], invSa4, invDa4)
	p.b[4] = ch(p.b[4], p.db[4], p.a[4], p.da[4], invSa4, invDa4)
	p.r[5] = ch(p.r[5], p.dr[5], p.a[5], p.da[5], invSa5, invDa5)
	p.g[5] = ch(p.g[5], p.dg[5], p.a[5], p.da[5], invSa5, invDa5)
	p.b[5] = ch(p.b[5], p.db[5], p.a[5], p.da[5], invSa5, invDa5)
	p.r[6] = ch(p.r[6], p.dr[6], p.a[6], p.da[6], invSa6, invDa6)
	p.g[6] = ch(p.g[6], p.dg[6], p.a[6], p.da[6], invSa6, invDa6)
	p.b[6] = ch(p.b[6], p.db[6], p.a[6], p.da[6], invSa6, invDa6)
	p.r[7] = ch(p.r[7], p.dr[7], p.a[7], p.da[7], invSa7, invDa7)
	p.g[7] = ch(p.g[7], p.dg[7], p.a[7], p.da[7], invSa7, invDa7)
	p.b[7] = ch(p.b[7], p.db[7], p.a[7], p.da[7], invSa7, invDa7)
	p.a[0], p.a[1], p.a[2], p.a[3] = p.a[0]+p.da[0]*invSa0, p.a[1]+p.da[1]*invSa1, p.a[2]+p.da[2]*invSa2, p.a[3]+p.da[3]*invSa3
	p.a[4], p.a[5], p.a[6], p.a[7] = p.a[4]+p.da[4]*invSa4, p.a[5]+p.da[5]*invSa5, p.a[6]+p.da[6]*invSa6, p.a[7]+p.da[7]*invSa7

}

//go:fix inline
func (p *HighPipeline) Hue() {
	for i := 0; i < HIGH_STAGE_WIDTH; i++ {
		invSa, invDa := 1.0-p.a[i], 1.0-p.da[i]

		rr, gg, bb := p.r[i]*p.a[i], p.g[i]*p.a[i], p.b[i]*p.a[i]

		mn, mx := f32minmax(rr, gg, bb)
		mnd, mxd := f32minmax(p.dr[i], p.dg[i], p.db[i])
		lum, lumd := rr*0.30+gg*0.59+bb*0.11, p.dr[i]*0.30+p.dg[i]*0.59+p.db[i]*0.11

		rr, gg, bb = f32sat(rr, gg, bb, p.a[i], mn, mx, mnd, mxd)
		rr, gg, bb = f32lum(rr, gg, bb, p.a[i], lum, lumd)
		rr, gg, bb = f32clip(rr, gg, bb, p.a[i]*p.da[i], mn, mx, lum)

		p.r[i] = p.r[i]*invDa + p.dr[i]*invSa + f32max(rr, 0)
		p.g[i] = p.g[i]*invDa + p.dg[i]*invSa + f32max(gg, 0)
		p.b[i] = p.b[i]*invDa + p.db[i]*invSa + f32max(bb, 0)
		p.a[i] = p.a[i] + p.da[i] - p.a[i]*p.da[i]
	}
}

//go:fix inline
func (p *HighPipeline) Saturation() {
	for i := 0; i < HIGH_STAGE_WIDTH; i++ {
		invSa, invDa := 1.0-p.a[i], 1.0-p.da[i]

		rr, gg, bb := p.dr[i]*p.a[i], p.dg[i]*p.a[i], p.db[i]*p.a[i]

		mn, mx := f32minmax(rr, gg, bb)
		mnd, mxd := f32minmax(p.r[i], p.g[i], p.b[i])
		lum, lumd := rr*0.30+gg*0.59+bb*0.11, p.dr[i]*0.30+p.dg[i]*0.59+p.db[i]*0.11

		rr, gg, bb = f32sat(rr, gg, bb, p.a[i], mn, mx, mnd, mxd)
		rr, gg, bb = f32lum(rr, gg, bb, p.a[i], lum, lumd)
		rr, gg, bb = f32clip(rr, gg, bb, p.a[i]*p.da[i], mn, mx, lum)

		p.r[i] = p.r[i]*invDa + p.dr[i]*invSa + f32max(rr, 0)
		p.g[i] = p.g[i]*invDa + p.dg[i]*invSa + f32max(gg, 0)
		p.b[i] = p.b[i]*invDa + p.db[i]*invSa + f32max(bb, 0)
		p.a[i] = p.a[i] + p.da[i] - p.a[i]*p.da[i]
	}
}

//go:fix inline
func (p *HighPipeline) Color() {
	for i := 0; i < HIGH_STAGE_WIDTH; i++ {
		invSa, invDa := 1.0-p.a[i], 1.0-p.da[i]

		rr, gg, bb := p.r[i]*p.da[i], p.g[i]*p.da[i], p.b[i]*p.da[i]

		mn, mx := f32minmax(rr, gg, bb)
		lum, lumd := rr*0.30+gg*0.59+bb*0.11, p.dr[i]*0.30+p.dg[i]*0.59+p.db[i]*0.11

		rr, gg, bb = f32lum(rr, gg, bb, p.da[i], lum, lumd)
		rr, gg, bb = f32clip(rr, gg, bb, p.a[i]*p.da[i], mn, mx, lum)

		p.r[i] = p.r[i]*invDa + p.dr[i]*invSa + f32max(rr, 0)
		p.g[i] = p.g[i]*invDa + p.dg[i]*invSa + f32max(gg, 0)
		p.b[i] = p.b[i]*invDa + p.db[i]*invSa + f32max(bb, 0)
		p.a[i] = p.a[i] + p.da[i] - p.a[i]*p.da[i]
	}
}

//go:fix inline
func (p *HighPipeline) Luminosity() {
	for i := 0; i < HIGH_STAGE_WIDTH; i++ {
		invSa, invDa := 1.0-p.a[i], 1.0-p.da[i]

		rr, gg, bb := p.dr[i]*p.a[i], p.dg[i]*p.a[i], p.db[i]*p.a[i]

		mn, mx := f32minmax(rr, gg, bb)
		lum, lumd := rr*0.30+gg*0.59+bb*0.11, p.r[i]*0.30+p.g[i]*0.59+p.b[i]*0.11

		rr, gg, bb = f32lum(rr, gg, bb, p.da[i], lum, lumd)
		rr, gg, bb = f32clip(rr, gg, bb, p.a[i]*p.da[i], mn, mx, lum)

		p.r[i] = p.r[i]*invDa + p.dr[i]*invSa + f32max(rr, 0)
		p.g[i] = p.g[i]*invDa + p.dg[i]*invSa + f32max(gg, 0)
		p.b[i] = p.b[i]*invDa + p.db[i]*invSa + f32max(bb, 0)
		p.a[i] = p.a[i] + p.da[i]*invSa
	}
}

//go:fix inline
func (p *HighPipeline) SourceOverRgba() {
	const FACTOR = 1.0 / 255.0
	offset := (p.dy*p.pixmapDst.RealWidth + p.dx) * 4
	data := p.pixmapDst.Data[offset : offset+HIGH_STAGE_WIDTH*4]

	p.dr[0], p.dr[1], p.dr[2], p.dr[3] = float32(data[0])*FACTOR, float32(data[4])*FACTOR, float32(data[8])*FACTOR, float32(data[12])*FACTOR
	p.dr[4], p.dr[5], p.dr[6], p.dr[7] = float32(data[16])*FACTOR, float32(data[20])*FACTOR, float32(data[24])*FACTOR, float32(data[28])*FACTOR
	p.dg[0], p.dg[1], p.dg[2], p.dg[3] = float32(data[1])*FACTOR, float32(data[5])*FACTOR, float32(data[9])*FACTOR, float32(data[13])*FACTOR
	p.dg[4], p.dg[5], p.dg[6], p.dg[7] = float32(data[17])*FACTOR, float32(data[21])*FACTOR, float32(data[25])*FACTOR, float32(data[29])*FACTOR
	p.db[0], p.db[1], p.db[2], p.db[3] = float32(data[2])*FACTOR, float32(data[6])*FACTOR, float32(data[10])*FACTOR, float32(data[14])*FACTOR
	p.db[4], p.db[5], p.db[6], p.db[7] = float32(data[18])*FACTOR, float32(data[22])*FACTOR, float32(data[26])*FACTOR, float32(data[30])*FACTOR
	p.da[0], p.da[1], p.da[2], p.da[3] = float32(data[3])*FACTOR, float32(data[7])*FACTOR, float32(data[11])*FACTOR, float32(data[15])*FACTOR
	p.da[4], p.da[5], p.da[6], p.da[7] = float32(data[19])*FACTOR, float32(data[23])*FACTOR, float32(data[27])*FACTOR, float32(data[31])*FACTOR

	p.r[0], p.r[1], p.r[2], p.r[3] = p.dr[0]*(1-p.a[0])+p.r[0], p.dr[1]*(1-p.a[1])+p.r[1], p.dr[2]*(1-p.a[2])+p.r[2], p.dr[3]*(1-p.a[3])+p.r[3]
	p.r[4], p.r[5], p.r[6], p.r[7] = p.dr[4]*(1-p.a[4])+p.r[4], p.dr[5]*(1-p.a[5])+p.r[5], p.dr[6]*(1-p.a[6])+p.r[6], p.dr[7]*(1-p.a[7])+p.r[7]
	p.g[0], p.g[1], p.g[2], p.g[3] = p.dg[0]*(1-p.a[0])+p.g[0], p.dg[1]*(1-p.a[1])+p.g[1], p.dg[2]*(1-p.a[2])+p.g[2], p.dg[3]*(1-p.a[3])+p.g[3]
	p.g[4], p.g[5], p.g[6], p.g[7] = p.dg[4]*(1-p.a[4])+p.g[4], p.dg[5]*(1-p.a[5])+p.g[5], p.dg[6]*(1-p.a[6])+p.g[6], p.dg[7]*(1-p.a[7])+p.g[7]
	p.b[0], p.b[1], p.b[2], p.b[3] = p.db[0]*(1-p.a[0])+p.b[0], p.db[1]*(1-p.a[1])+p.b[1], p.db[2]*(1-p.a[2])+p.b[2], p.db[3]*(1-p.a[3])+p.b[3]
	p.b[4], p.b[5], p.b[6], p.b[7] = p.db[4]*(1-p.a[4])+p.b[4], p.db[5]*(1-p.a[5])+p.b[5], p.db[6]*(1-p.a[6])+p.b[6], p.db[7]*(1-p.a[7])+p.b[7]
	p.a[0], p.a[1], p.a[2], p.a[3] = p.da[0]*(1-p.a[0])+p.a[0], p.da[1]*(1-p.a[1])+p.a[1], p.da[2]*(1-p.a[2])+p.a[2], p.da[3]*(1-p.a[3])+p.a[3]
	p.a[4], p.a[5], p.a[6], p.a[7] = p.da[4]*(1-p.a[4])+p.a[4], p.da[5]*(1-p.a[5])+p.a[5], p.da[6]*(1-p.a[6])+p.a[6], p.da[7]*(1-p.a[7])+p.a[7]

	data[0], data[4], data[8], data[12] = f32unnorm(p.r[0]), f32unnorm(p.r[1]), f32unnorm(p.r[2]), f32unnorm(p.r[3])
	data[16], data[20], data[24], data[28] = f32unnorm(p.r[4]), f32unnorm(p.r[5]), f32unnorm(p.r[6]), f32unnorm(p.r[7])
	data[1], data[5], data[9], data[13] = f32unnorm(p.g[0]), f32unnorm(p.g[1]), f32unnorm(p.g[2]), f32unnorm(p.g[3])
	data[17], data[21], data[25], data[29] = f32unnorm(p.g[4]), f32unnorm(p.g[5]), f32unnorm(p.g[6]), f32unnorm(p.g[7])
	data[2], data[6], data[10], data[14] = f32unnorm(p.b[0]), f32unnorm(p.b[1]), f32unnorm(p.b[2]), f32unnorm(p.b[3])
	data[18], data[22], data[26], data[30] = f32unnorm(p.b[4]), f32unnorm(p.b[5]), f32unnorm(p.b[6]), f32unnorm(p.b[7])
	data[3], data[7], data[11], data[15] = f32unnorm(p.a[0]), f32unnorm(p.a[1]), f32unnorm(p.a[2]), f32unnorm(p.a[3])
	data[19], data[23], data[27], data[31] = f32unnorm(p.a[4]), f32unnorm(p.a[5]), f32unnorm(p.a[6]), f32unnorm(p.a[7])
}

//go:fix inline
func (p *HighPipeline) SourceOverRgbaTail() {
	const FACTOR = 1.0 / 255.0

	dstTmp := [32]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	offset := (p.dy*p.pixmapDst.RealWidth + p.dx) * 4
	for i := 0; i < p.tail; i++ {
		srcIdx := offset + i*4
		dstTmp[i*4+0] = p.pixmapDst.Data[srcIdx+0]
		dstTmp[i*4+1] = p.pixmapDst.Data[srcIdx+1]
		dstTmp[i*4+2] = p.pixmapDst.Data[srcIdx+2]
		dstTmp[i*4+3] = p.pixmapDst.Data[srcIdx+3]
	}

	p.dr[0], p.dr[1], p.dr[2], p.dr[3] = float32(dstTmp[0])*FACTOR, float32(dstTmp[4])*FACTOR, float32(dstTmp[8])*FACTOR, float32(dstTmp[12])*FACTOR
	p.dr[4], p.dr[5], p.dr[6], p.dr[7] = float32(dstTmp[16])*FACTOR, float32(dstTmp[20])*FACTOR, float32(dstTmp[24])*FACTOR, float32(dstTmp[28])*FACTOR
	p.dg[0], p.dg[1], p.dg[2], p.dg[3] = float32(dstTmp[1])*FACTOR, float32(dstTmp[5])*FACTOR, float32(dstTmp[9])*FACTOR, float32(dstTmp[13])*FACTOR
	p.dg[4], p.dg[5], p.dg[6], p.dg[7] = float32(dstTmp[17])*FACTOR, float32(dstTmp[21])*FACTOR, float32(dstTmp[25])*FACTOR, float32(dstTmp[29])*FACTOR
	p.db[0], p.db[1], p.db[2], p.db[3] = float32(dstTmp[2])*FACTOR, float32(dstTmp[6])*FACTOR, float32(dstTmp[10])*FACTOR, float32(dstTmp[14])*FACTOR
	p.db[4], p.db[5], p.db[6], p.db[7] = float32(dstTmp[18])*FACTOR, float32(dstTmp[22])*FACTOR, float32(dstTmp[26])*FACTOR, float32(dstTmp[30])*FACTOR
	p.da[0], p.da[1], p.da[2], p.da[3] = float32(dstTmp[3])*FACTOR, float32(dstTmp[7])*FACTOR, float32(dstTmp[11])*FACTOR, float32(dstTmp[15])*FACTOR
	p.da[4], p.da[5], p.da[6], p.da[7] = float32(dstTmp[19])*FACTOR, float32(dstTmp[23])*FACTOR, float32(dstTmp[27])*FACTOR, float32(dstTmp[31])*FACTOR

	p.r[0], p.r[1], p.r[2], p.r[3] = p.dr[0]*(1-p.a[0])+p.r[0], p.dr[1]*(1-p.a[1])+p.r[1], p.dr[2]*(1-p.a[2])+p.r[2], p.dr[3]*(1-p.a[3])+p.r[3]
	p.r[4], p.r[5], p.r[6], p.r[7] = p.dr[4]*(1-p.a[4])+p.r[4], p.dr[5]*(1-p.a[5])+p.r[5], p.dr[6]*(1-p.a[6])+p.r[6], p.dr[7]*(1-p.a[7])+p.r[7]
	p.g[0], p.g[1], p.g[2], p.g[3] = p.dg[0]*(1-p.a[0])+p.g[0], p.dg[1]*(1-p.a[1])+p.g[1], p.dg[2]*(1-p.a[2])+p.g[2], p.dg[3]*(1-p.a[3])+p.g[3]
	p.g[4], p.g[5], p.g[6], p.g[7] = p.dg[4]*(1-p.a[4])+p.g[4], p.dg[5]*(1-p.a[5])+p.g[5], p.dg[6]*(1-p.a[6])+p.g[6], p.dg[7]*(1-p.a[7])+p.g[7]
	p.b[0], p.b[1], p.b[2], p.b[3] = p.db[0]*(1-p.a[0])+p.b[0], p.db[1]*(1-p.a[1])+p.b[1], p.db[2]*(1-p.a[2])+p.b[2], p.db[3]*(1-p.a[3])+p.b[3]
	p.b[4], p.b[5], p.b[6], p.b[7] = p.db[4]*(1-p.a[4])+p.b[4], p.db[5]*(1-p.a[5])+p.b[5], p.db[6]*(1-p.a[6])+p.b[6], p.db[7]*(1-p.a[7])+p.b[7]
	p.a[0], p.a[1], p.a[2], p.a[3] = p.da[0]*(1-p.a[0])+p.a[0], p.da[1]*(1-p.a[1])+p.a[1], p.da[2]*(1-p.a[2])+p.a[2], p.da[3]*(1-p.a[3])+p.a[3]
	p.a[4], p.a[5], p.a[6], p.a[7] = p.da[4]*(1-p.a[4])+p.a[4], p.da[5]*(1-p.a[5])+p.a[5], p.da[6]*(1-p.a[6])+p.a[6], p.da[7]*(1-p.a[7])+p.a[7]

	srcTmp := [32]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	srcTmp[0], srcTmp[4], srcTmp[8], srcTmp[12] = f32unnorm(p.r[0]), f32unnorm(p.r[1]), f32unnorm(p.r[2]), f32unnorm(p.r[3])
	srcTmp[16], srcTmp[20], srcTmp[24], srcTmp[28] = f32unnorm(p.r[4]), f32unnorm(p.r[5]), f32unnorm(p.r[6]), f32unnorm(p.r[7])
	srcTmp[1], srcTmp[5], srcTmp[9], srcTmp[13] = f32unnorm(p.g[0]), f32unnorm(p.g[1]), f32unnorm(p.g[2]), f32unnorm(p.g[3])
	srcTmp[17], srcTmp[21], srcTmp[25], srcTmp[29] = f32unnorm(p.g[4]), f32unnorm(p.g[5]), f32unnorm(p.g[6]), f32unnorm(p.g[7])
	srcTmp[2], srcTmp[6], srcTmp[10], srcTmp[14] = f32unnorm(p.b[0]), f32unnorm(p.b[1]), f32unnorm(p.b[2]), f32unnorm(p.b[3])
	srcTmp[18], srcTmp[22], srcTmp[26], srcTmp[30] = f32unnorm(p.b[4]), f32unnorm(p.b[5]), f32unnorm(p.b[6]), f32unnorm(p.b[7])
	srcTmp[3], srcTmp[7], srcTmp[11], srcTmp[15] = f32unnorm(p.a[0]), f32unnorm(p.a[1]), f32unnorm(p.a[2]), f32unnorm(p.a[3])
	srcTmp[19], srcTmp[23], srcTmp[27], srcTmp[31] = f32unnorm(p.a[4]), f32unnorm(p.a[5]), f32unnorm(p.a[6]), f32unnorm(p.a[7])

	for i := 0; i < p.tail; i++ {
		dstIdx := offset + i*4
		p.pixmapDst.Data[dstIdx+0] = srcTmp[i*4+0]
		p.pixmapDst.Data[dstIdx+1] = srcTmp[i*4+1]
		p.pixmapDst.Data[dstIdx+2] = srcTmp[i*4+2]
		p.pixmapDst.Data[dstIdx+3] = srcTmp[i*4+3]
	}
}

//go:fix inline
func (p *HighPipeline) Transform() {
	ts := &p.ctx.Transform

	r0, r1, r2, r3 := p.r[0], p.r[1], p.r[2], p.r[3]
	g0, g1, g2, g3 := p.g[0], p.g[1], p.g[2], p.g[3]
	p.r[0], p.g[0] = r0*ts.SX+g0*ts.KX+ts.TX, r0*ts.KY+g0*ts.SY+ts.TY
	p.r[1], p.g[1] = r1*ts.SX+g1*ts.KX+ts.TX, r1*ts.KY+g1*ts.SY+ts.TY
	p.r[2], p.g[2] = r2*ts.SX+g2*ts.KX+ts.TX, r2*ts.KY+g2*ts.SY+ts.TY
	p.r[3], p.g[3] = r3*ts.SX+g3*ts.KX+ts.TX, r3*ts.KY+g3*ts.SY+ts.TY

	r4, r5, r6, r7 := p.r[4], p.r[5], p.r[6], p.r[7]
	g4, g5, g6, g7 := p.g[4], p.g[5], p.g[6], p.g[7]
	p.r[4], p.g[4] = r4*ts.SX+g4*ts.KX+ts.TX, r4*ts.KY+g4*ts.SY+ts.TY
	p.r[5], p.g[5] = r5*ts.SX+g5*ts.KX+ts.TX, r5*ts.KY+g5*ts.SY+ts.TY
	p.r[6], p.g[6] = r6*ts.SX+g6*ts.KX+ts.TX, r6*ts.KY+g6*ts.SY+ts.TY
	p.r[7], p.g[7] = r7*ts.SX+g7*ts.KX+ts.TX, r7*ts.KY+g7*ts.SY+ts.TY
}

//go:fix inline
func (p *HighPipeline) Reflect() {
	reflect := func(v, limit, invLimit float32) float32 {
		return f32abs((v - limit) - (limit+limit)*math32.Floor((v-limit)*(invLimit*0.5)) - limit)
	}
	limitX := &p.ctx.LimitX
	limitY := &p.ctx.LimitY

	p.r[0] = reflect(p.r[0], limitX.Scale, limitX.InvScale)
	p.g[0] = reflect(p.g[0], limitY.Scale, limitY.InvScale)
	p.r[1] = reflect(p.r[1], limitX.Scale, limitX.InvScale)
	p.g[1] = reflect(p.g[1], limitY.Scale, limitY.InvScale)
	p.r[2] = reflect(p.r[2], limitX.Scale, limitX.InvScale)
	p.g[2] = reflect(p.g[2], limitY.Scale, limitY.InvScale)
	p.r[3] = reflect(p.r[3], limitX.Scale, limitX.InvScale)
	p.g[3] = reflect(p.g[3], limitY.Scale, limitY.InvScale)
	p.r[4] = reflect(p.r[4], limitX.Scale, limitX.InvScale)
	p.g[4] = reflect(p.g[4], limitY.Scale, limitY.InvScale)
	p.r[5] = reflect(p.r[5], limitX.Scale, limitX.InvScale)
	p.g[5] = reflect(p.g[5], limitY.Scale, limitY.InvScale)
	p.r[6] = reflect(p.r[6], limitX.Scale, limitX.InvScale)
	p.g[6] = reflect(p.g[6], limitY.Scale, limitY.InvScale)
	p.r[7] = reflect(p.r[7], limitX.Scale, limitX.InvScale)
	p.g[7] = reflect(p.g[7], limitY.Scale, limitY.InvScale)

}

//go:fix inline
func (p *HighPipeline) Repeat() {
	repeat := func(v, limit, invLimit float32) float32 {
		return v - math32.Floor(v*invLimit)*limit
	}
	limitX := &p.ctx.LimitX
	limitY := &p.ctx.LimitY
	p.r[0] = repeat(p.r[0], limitX.Scale, limitX.InvScale)
	p.g[0] = repeat(p.g[0], limitY.Scale, limitY.InvScale)
	p.r[1] = repeat(p.r[1], limitX.Scale, limitX.InvScale)
	p.g[1] = repeat(p.g[1], limitY.Scale, limitY.InvScale)
	p.r[2] = repeat(p.r[2], limitX.Scale, limitX.InvScale)
	p.g[2] = repeat(p.g[2], limitY.Scale, limitY.InvScale)
	p.r[3] = repeat(p.r[3], limitX.Scale, limitX.InvScale)
	p.g[3] = repeat(p.g[3], limitY.Scale, limitY.InvScale)
	p.r[4] = repeat(p.r[4], limitX.Scale, limitX.InvScale)
	p.g[4] = repeat(p.g[4], limitY.Scale, limitY.InvScale)
	p.r[5] = repeat(p.r[5], limitX.Scale, limitX.InvScale)
	p.g[5] = repeat(p.g[5], limitY.Scale, limitY.InvScale)
	p.r[6] = repeat(p.r[6], limitX.Scale, limitX.InvScale)
	p.g[6] = repeat(p.g[6], limitY.Scale, limitY.InvScale)
	p.r[7] = repeat(p.r[7], limitX.Scale, limitX.InvScale)
	p.g[7] = repeat(p.g[7], limitY.Scale, limitY.InvScale)

}

//go:fix inline
func (p *HighPipeline) Bilinear() {
	for i := 0; i < HIGH_STAGE_WIDTH; i++ {
		x, y := p.r[i], p.g[i]
		fx := x + 0.5
		fx = fx - math32.Floor(fx)
		fy := y + 0.5
		fy = fy - math32.Floor(fy)
		wx0, wx1 := 1-fx, fx
		wy0, wy1 := 1-fy, fy
		start := float32(-0.5)
		var r, g, b, a float32
		for j := 0; j < 2; j++ {
			for k := 0; k < 2; k++ {
				sx := x + start + float32(k)
				sy := y + start + float32(j)
				var rr, gg, bb, aa float32
				samplePixel(p.pixmapSrc, &p.ctx.Sampler, sx, sy, &rr, &gg, &bb, &aa)
				w := wx0
				if k == 1 {
					w = wx1
				}
				if j == 1 {
					w *= wy1
				} else {
					w *= wy0
				}
				r += w * rr
				g += w * gg
				b += w * bb
				a += w * aa
			}
		}
		p.r[i], p.g[i], p.b[i], p.a[i] = r, g, b, a
	}
}

//go:fix inline
func (p *HighPipeline) Bicubic() {
	bicubicNear := func(t float32) float32 {
		return t*(t*(-21.0/18.0*t+27.0/18.0)+9.0/18.0) + 1.0/18.0
	}
	bicubicFar := func(t float32) float32 {
		return t * t * (7.0/18.0*t - 6.0/18.0)
	}
	for i := 0; i < HIGH_STAGE_WIDTH; i++ {
		x, y := p.r[i], p.g[i]
		fx := x + 0.5
		fx = fx - math32.Floor(fx)
		fy := y + 0.5
		fy = fy - math32.Floor(fy)
		wx := [4]float32{bicubicFar(1 - fx), bicubicNear(1 - fx), bicubicNear(fx), bicubicFar(fx)}
		wy := [4]float32{bicubicFar(1 - fy), bicubicNear(1 - fy), bicubicNear(fy), bicubicFar(fy)}
		start := float32(-1.5)
		var r, g, b, a float32
		for j := 0; j < 4; j++ {
			for k := 0; k < 4; k++ {
				sx := x + start + float32(k)
				sy := y + start + float32(j)
				var rr, gg, bb, aa float32
				samplePixel(p.pixmapSrc, &p.ctx.Sampler, sx, sy, &rr, &gg, &bb, &aa)
				w := wx[k] * wy[j]
				r += w * rr
				g += w * gg
				b += w * bb
				a += w * aa
			}
		}
		p.r[i], p.g[i], p.b[i], p.a[i] = r, g, b, a
	}
}

//go:fix inline
func (p *HighPipeline) PadX1() {
	p.r[0], p.r[1], p.r[2], p.r[3] = f32normalize(p.r[0]), f32normalize(p.r[1]), f32normalize(p.r[2]), f32normalize(p.r[3])
	p.r[4], p.r[5], p.r[6], p.r[7] = f32normalize(p.r[4]), f32normalize(p.r[5]), f32normalize(p.r[6]), f32normalize(p.r[7])

}

//go:fix inline
func (p *HighPipeline) ReflectX1() {
	v0, v1, v2, v3 := p.r[0]-1, p.r[1]-1, p.r[2]-1, p.r[3]-1
	v4, v5, v6, v7 := p.r[4]-1, p.r[5]-1, p.r[6]-1, p.r[7]-1

	p.r[0], p.r[1], p.r[2], p.r[3] = f32normalize(f32abs(v0-2*math32.Floor(v0*0.5)-1)), f32normalize(f32abs(v1-2*math32.Floor(v1*0.5)-1)), f32normalize(f32abs(v2-2*math32.Floor(v2*0.5)-1)), f32normalize(f32abs(v3-2*math32.Floor(v3*0.5)-1))
	p.r[4], p.r[5], p.r[6], p.r[7] = f32normalize(f32abs(v4-2*math32.Floor(v4*0.5)-1)), f32normalize(f32abs(v5-2*math32.Floor(v5*0.5)-1)), f32normalize(f32abs(v6-2*math32.Floor(v6*0.5)-1)), f32normalize(f32abs(v7-2*math32.Floor(v7*0.5)-1))

}

//go:fix inline
func (p *HighPipeline) RepeatX1() {
	p.r[0], p.r[1], p.r[2], p.r[3] = f32normalize(p.r[0]-math32.Floor(p.r[0])), f32normalize(p.r[1]-math32.Floor(p.r[1])), f32normalize(p.r[2]-math32.Floor(p.r[2])), f32normalize(p.r[3]-math32.Floor(p.r[3]))
	p.r[4], p.r[5], p.r[6], p.r[7] = f32normalize(p.r[4]-math32.Floor(p.r[4])), f32normalize(p.r[5]-math32.Floor(p.r[5])), f32normalize(p.r[6]-math32.Floor(p.r[6])), f32normalize(p.r[7]-math32.Floor(p.r[7]))

}

//go:fix inline
func (p *HighPipeline) Gradient() {
	ctx := &p.ctx.Gradient

	for i := 0; i < 8; i++ {
		t := p.r[i]

		idx := uint32(0)
		for j := 1; j < ctx.Len; j++ {
			if t >= ctx.TValues[j].Get() {
				idx++
			}
		}

		f := ctx.Factors[idx]
		b := ctx.Biases[idx]

		p.r[i] = t*f.R + b.R
		p.g[i] = t*f.G + b.G
		p.b[i] = t*f.B + b.B
		p.a[i] = t*f.A + b.A
	}
}

//go:fix inline
func (p *HighPipeline) EvenlySpaced2StopGradient() {
	esg := &p.ctx.EvenlySpaced2StopGradient

	t0, t1, t2, t3 := p.r[0], p.r[1], p.r[2], p.r[3]
	t4, t5, t6, t7 := p.r[4], p.r[5], p.r[6], p.r[7]

	p.r[0], p.g[0], p.b[0], p.a[0] = t0*esg.Factor.R+esg.Bias.R, t0*esg.Factor.G+esg.Bias.G, t0*esg.Factor.B+esg.Bias.B, t0*esg.Factor.A+esg.Bias.A
	p.r[1], p.g[1], p.b[1], p.a[1] = t1*esg.Factor.R+esg.Bias.R, t1*esg.Factor.G+esg.Bias.G, t1*esg.Factor.B+esg.Bias.B, t1*esg.Factor.A+esg.Bias.A
	p.r[2], p.g[2], p.b[2], p.a[2] = t2*esg.Factor.R+esg.Bias.R, t2*esg.Factor.G+esg.Bias.G, t2*esg.Factor.B+esg.Bias.B, t2*esg.Factor.A+esg.Bias.A
	p.r[3], p.g[3], p.b[3], p.a[3] = t3*esg.Factor.R+esg.Bias.R, t3*esg.Factor.G+esg.Bias.G, t3*esg.Factor.B+esg.Bias.B, t3*esg.Factor.A+esg.Bias.A
	p.r[4], p.g[4], p.b[4], p.a[4] = t4*esg.Factor.R+esg.Bias.R, t4*esg.Factor.G+esg.Bias.G, t4*esg.Factor.B+esg.Bias.B, t4*esg.Factor.A+esg.Bias.A
	p.r[5], p.g[5], p.b[5], p.a[5] = t5*esg.Factor.R+esg.Bias.R, t5*esg.Factor.G+esg.Bias.G, t5*esg.Factor.B+esg.Bias.B, t5*esg.Factor.A+esg.Bias.A
	p.r[6], p.g[6], p.b[6], p.a[6] = t6*esg.Factor.R+esg.Bias.R, t6*esg.Factor.G+esg.Bias.G, t6*esg.Factor.B+esg.Bias.B, t6*esg.Factor.A+esg.Bias.A
	p.r[7], p.g[7], p.b[7], p.a[7] = t7*esg.Factor.R+esg.Bias.R, t7*esg.Factor.G+esg.Bias.G, t7*esg.Factor.B+esg.Bias.B, t7*esg.Factor.A+esg.Bias.A

}

//go:fix inline
func (p *HighPipeline) XYToUnitAngle() {
	unitAngle := func(x, y float32) float32 {
		xAbs, yAbs := math32.Abs(x), math32.Abs(y)
		slope := math32.Min(xAbs, yAbs) / math32.Max(xAbs, yAbs)
		s := slope * slope
		// 7th degree polynomial approximation for atan
		phi := slope * (0.15912117063999176025390625 + s*(-5.185396969318389892578125e-2+s*(2.476101927459239959716796875e-2+s*(-7.0547382347285747528076171875e-3))))
		if xAbs < yAbs {
			phi = 0.25 - phi
		}
		if x < 0 {
			phi = 0.5 - phi
		}
		if y < 0 {
			phi = 1.0 - phi
		}
		if phi != phi { // NaN check
			phi = 0.0
		}
		return phi
	}

	// Manually unrolled for performance (8-way SIMD-like processing)
	p.r[0], p.r[1], p.r[2], p.r[3] = unitAngle(p.r[0], p.g[0]), unitAngle(p.r[1], p.g[1]), unitAngle(p.r[2], p.g[2]), unitAngle(p.r[3], p.g[3])
	p.r[4], p.r[5], p.r[6], p.r[7] = unitAngle(p.r[4], p.g[4]), unitAngle(p.r[5], p.g[5]), unitAngle(p.r[6], p.g[6]), unitAngle(p.r[7], p.g[7])

}

//go:fix inline
func (p *HighPipeline) XYToRadius() {
	radius := func(x, y float32) float32 {
		return math32.Sqrt(x*x + y*y)
	}

	p.r[0], p.r[1], p.r[2], p.r[3] = radius(p.r[0], p.g[0]), radius(p.r[1], p.g[1]), radius(p.r[2], p.g[2]), radius(p.r[3], p.g[3])
	p.r[4], p.r[5], p.r[6], p.r[7] = radius(p.r[4], p.g[4]), radius(p.r[5], p.g[5]), radius(p.r[6], p.g[6]), radius(p.r[7], p.g[7])

}

//go:fix inline
func (p *HighPipeline) XYTo2PtConicalFocalOnCircle() {
	focal := func(x, y float32) float32 {
		return x + y*y/x
	}

	p.r[0], p.r[1], p.r[2], p.r[3] = focal(p.r[0], p.g[0]), focal(p.r[1], p.g[1]), focal(p.r[2], p.g[2]), focal(p.r[3], p.g[3])
	p.r[4], p.r[5], p.r[6], p.r[7] = focal(p.r[4], p.g[4]), focal(p.r[5], p.g[5]), focal(p.r[6], p.g[6]), focal(p.r[7], p.g[7])

}

//go:fix inline
func (p *HighPipeline) XYTo2PtConicalWellBehaved() {
	p0 := p.ctx.TwoPointConicalGradient.P0

	wellBehaved := func(x, y float32) float32 {
		return math32.Sqrt(x*x+y*y) - x*p0
	}

	p.r[0], p.r[1], p.r[2], p.r[3] = wellBehaved(p.r[0], p.g[0]), wellBehaved(p.r[1], p.g[1]), wellBehaved(p.r[2], p.g[2]), wellBehaved(p.r[3], p.g[3])
	p.r[4], p.r[5], p.r[6], p.r[7] = wellBehaved(p.r[4], p.g[4]), wellBehaved(p.r[5], p.g[5]), wellBehaved(p.r[6], p.g[6]), wellBehaved(p.r[7], p.g[7])

}

//go:fix inline
func (p *HighPipeline) XYTo2PtConicalSmaller() {
	p0 := p.ctx.TwoPointConicalGradient.P0

	smaller := func(x, y float32) float32 {
		return -math32.Sqrt(x*x-y*y) - x*p0
	}

	p.r[0], p.r[1], p.r[2], p.r[3] = smaller(p.r[0], p.g[0]), smaller(p.r[1], p.g[1]), smaller(p.r[2], p.g[2]), smaller(p.r[3], p.g[3])
	p.r[4], p.r[5], p.r[6], p.r[7] = smaller(p.r[4], p.g[4]), smaller(p.r[5], p.g[5]), smaller(p.r[6], p.g[6]), smaller(p.r[7], p.g[7])

}

//go:fix inline
func (p *HighPipeline) XYTo2PtConicalGreater() {
	p0 := p.ctx.TwoPointConicalGradient.P0

	greater := func(x, y float32) float32 {
		return math32.Sqrt(x*x+y*y) - x*p0
	}

	p.r[0], p.r[1], p.r[2], p.r[3] = greater(p.r[0], p.g[0]), greater(p.r[1], p.g[1]), greater(p.r[2], p.g[2]), greater(p.r[3], p.g[3])
	p.r[4], p.r[5], p.r[6], p.r[7] = greater(p.r[4], p.g[4]), greater(p.r[5], p.g[5]), greater(p.r[6], p.g[6]), greater(p.r[7], p.g[7])

}

//go:fix inline
func (p *HighPipeline) XYTo2PtConicalStrip() {
	p0 := p.ctx.TwoPointConicalGradient.P0

	strip := func(x, y float32) float32 {
		return x + math32.Sqrt(p0-y*y)
	}

	p.r[0], p.r[1], p.r[2], p.r[3] = strip(p.r[0], p.g[0]), strip(p.r[1], p.g[1]), strip(p.r[2], p.g[2]), strip(p.r[3], p.g[3])
	p.r[4], p.r[5], p.r[6], p.r[7] = strip(p.r[4], p.g[4]), strip(p.r[5], p.g[5]), strip(p.r[6], p.g[6]), strip(p.r[7], p.g[7])

}

//go:fix inline
func (p *HighPipeline) Mask2PtConicalNan() {
	ctx := &p.ctx.TwoPointConicalGradient

	isDegenerate0, isDegenerate1, isDegenerate2, isDegenerate3 := uint32(0), uint32(0), uint32(0), uint32(0)
	isDegenerate4, isDegenerate5, isDegenerate6, isDegenerate7 := uint32(0), uint32(0), uint32(0), uint32(0)

	if p.r[0] != p.r[0] {
		isDegenerate0 = 0xFFFFFFFF
		p.r[0] = 0.0
	}
	if p.r[1] != p.r[1] {
		isDegenerate1 = 0xFFFFFFFF
		p.r[1] = 0.0
	}
	if p.r[2] != p.r[2] {
		isDegenerate2 = 0xFFFFFFFF
		p.r[2] = 0.0
	}
	if p.r[3] != p.r[3] {
		isDegenerate3 = 0xFFFFFFFF
		p.r[3] = 0.0
	}
	if p.r[4] != p.r[4] {
		isDegenerate4 = 0xFFFFFFFF
		p.r[4] = 0.0
	}
	if p.r[5] != p.r[5] {
		isDegenerate5 = 0xFFFFFFFF
		p.r[5] = 0.0
	}
	if p.r[6] != p.r[6] {
		isDegenerate6 = 0xFFFFFFFF
		p.r[6] = 0.0
	}
	if p.r[7] != p.r[7] {
		isDegenerate7 = 0xFFFFFFFF
		p.r[7] = 0.0
	}

	ctx.Mask[0], ctx.Mask[1], ctx.Mask[2], ctx.Mask[3] = ^isDegenerate0, ^isDegenerate1, ^isDegenerate2, ^isDegenerate3
	ctx.Mask[4], ctx.Mask[5], ctx.Mask[6], ctx.Mask[7] = ^isDegenerate4, ^isDegenerate5, ^isDegenerate6, ^isDegenerate7

}

//go:fix inline
func (p *HighPipeline) Mask2PtConicalDegenerates() {
	ctx := &p.ctx.TwoPointConicalGradient

	isDegenerate0, isDegenerate1, isDegenerate2, isDegenerate3 := uint32(0), uint32(0), uint32(0), uint32(0)
	isDegenerate4, isDegenerate5, isDegenerate6, isDegenerate7 := uint32(0), uint32(0), uint32(0), uint32(0)

	if p.r[0] <= 0 || p.r[0] != p.r[0] {
		isDegenerate0 = 0xFFFFFFFF
		p.r[0] = 0.0
	}
	if p.r[1] <= 0 || p.r[1] != p.r[1] {
		isDegenerate1 = 0xFFFFFFFF
		p.r[1] = 0.0
	}
	if p.r[2] <= 0 || p.r[2] != p.r[2] {
		isDegenerate2 = 0xFFFFFFFF
		p.r[2] = 0.0
	}
	if p.r[3] <= 0 || p.r[3] != p.r[3] {
		isDegenerate3 = 0xFFFFFFFF
		p.r[3] = 0.0
	}
	if p.r[4] <= 0 || p.r[4] != p.r[4] {
		isDegenerate4 = 0xFFFFFFFF
		p.r[4] = 0.0
	}
	if p.r[5] <= 0 || p.r[5] != p.r[5] {
		isDegenerate5 = 0xFFFFFFFF
		p.r[5] = 0.0
	}
	if p.r[6] <= 0 || p.r[6] != p.r[6] {
		isDegenerate6 = 0xFFFFFFFF
		p.r[6] = 0.0
	}
	if p.r[7] <= 0 || p.r[7] != p.r[7] {
		isDegenerate7 = 0xFFFFFFFF
		p.r[7] = 0.0
	}

	ctx.Mask[0], ctx.Mask[1], ctx.Mask[2], ctx.Mask[3] = ^isDegenerate0, ^isDegenerate1, ^isDegenerate2, ^isDegenerate3
	ctx.Mask[4], ctx.Mask[5], ctx.Mask[6], ctx.Mask[7] = ^isDegenerate4, ^isDegenerate5, ^isDegenerate6, ^isDegenerate7

}

//go:fix inline
func (p *HighPipeline) ApplyVectorMask() {
	ctx := &p.ctx.TwoPointConicalGradient

	p.r[0] = math32.Float32frombits(math32.Float32bits(p.r[0]) & ctx.Mask[0])
	p.r[1] = math32.Float32frombits(math32.Float32bits(p.r[1]) & ctx.Mask[1])
	p.r[2] = math32.Float32frombits(math32.Float32bits(p.r[2]) & ctx.Mask[2])
	p.r[3] = math32.Float32frombits(math32.Float32bits(p.r[3]) & ctx.Mask[3])
	p.r[4] = math32.Float32frombits(math32.Float32bits(p.r[4]) & ctx.Mask[4])
	p.r[5] = math32.Float32frombits(math32.Float32bits(p.r[5]) & ctx.Mask[5])
	p.r[6] = math32.Float32frombits(math32.Float32bits(p.r[6]) & ctx.Mask[6])
	p.r[7] = math32.Float32frombits(math32.Float32bits(p.r[7]) & ctx.Mask[7])

	p.g[0] = math32.Float32frombits(math32.Float32bits(p.g[0]) & ctx.Mask[0])
	p.g[1] = math32.Float32frombits(math32.Float32bits(p.g[1]) & ctx.Mask[1])
	p.g[2] = math32.Float32frombits(math32.Float32bits(p.g[2]) & ctx.Mask[2])
	p.g[3] = math32.Float32frombits(math32.Float32bits(p.g[3]) & ctx.Mask[3])
	p.g[4] = math32.Float32frombits(math32.Float32bits(p.g[4]) & ctx.Mask[4])
	p.g[5] = math32.Float32frombits(math32.Float32bits(p.g[5]) & ctx.Mask[5])
	p.g[6] = math32.Float32frombits(math32.Float32bits(p.g[6]) & ctx.Mask[6])
	p.g[7] = math32.Float32frombits(math32.Float32bits(p.g[7]) & ctx.Mask[7])

	p.b[0] = math32.Float32frombits(math32.Float32bits(p.b[0]) & ctx.Mask[0])
	p.b[1] = math32.Float32frombits(math32.Float32bits(p.b[1]) & ctx.Mask[1])
	p.b[2] = math32.Float32frombits(math32.Float32bits(p.b[2]) & ctx.Mask[2])
	p.b[3] = math32.Float32frombits(math32.Float32bits(p.b[3]) & ctx.Mask[3])
	p.b[4] = math32.Float32frombits(math32.Float32bits(p.b[4]) & ctx.Mask[4])
	p.b[5] = math32.Float32frombits(math32.Float32bits(p.b[5]) & ctx.Mask[5])
	p.b[6] = math32.Float32frombits(math32.Float32bits(p.b[6]) & ctx.Mask[6])
	p.b[7] = math32.Float32frombits(math32.Float32bits(p.b[7]) & ctx.Mask[7])

	p.a[0] = math32.Float32frombits(math32.Float32bits(p.a[0]) & ctx.Mask[0])
	p.a[1] = math32.Float32frombits(math32.Float32bits(p.a[1]) & ctx.Mask[1])
	p.a[2] = math32.Float32frombits(math32.Float32bits(p.a[2]) & ctx.Mask[2])
	p.a[3] = math32.Float32frombits(math32.Float32bits(p.a[3]) & ctx.Mask[3])
	p.a[4] = math32.Float32frombits(math32.Float32bits(p.a[4]) & ctx.Mask[4])
	p.a[5] = math32.Float32frombits(math32.Float32bits(p.a[5]) & ctx.Mask[5])
	p.a[6] = math32.Float32frombits(math32.Float32bits(p.a[6]) & ctx.Mask[6])
	p.a[7] = math32.Float32frombits(math32.Float32bits(p.a[7]) & ctx.Mask[7])

}

//go:fix inline
func (p *HighPipeline) Alter2PtConicalCompensateFocal() {
	p1 := p.ctx.TwoPointConicalGradient.P1

	p.r[0], p.r[1], p.r[2], p.r[3] = p.r[0]+p1, p.r[1]+p1, p.r[2]+p1, p.r[3]+p1
	p.r[4], p.r[5], p.r[6], p.r[7] = p.r[4]+p1, p.r[5]+p1, p.r[6]+p1, p.r[7]+p1

}

//go:fix inline
func (p *HighPipeline) Alter2PtConicalUnswap() {
	p.r[0], p.r[1], p.r[2], p.r[3] = 1.0-p.r[0], 1.0-p.r[1], 1.0-p.r[2], 1.0-p.r[3]
	p.r[4], p.r[5], p.r[6], p.r[7] = 1.0-p.r[4], 1.0-p.r[5], 1.0-p.r[6], 1.0-p.r[7]

}

//go:fix inline
func (p *HighPipeline) NegateX() {
	p.r[0], p.r[1], p.r[2], p.r[3] = -p.r[0], -p.r[1], -p.r[2], -p.r[3]
	p.r[4], p.r[5], p.r[6], p.r[7] = -p.r[4], -p.r[5], -p.r[6], -p.r[7]

}

//go:fix inline
func (p *HighPipeline) ApplyConcentricScaleBias() {
	tpcg := &p.ctx.TwoPointConicalGradient

	p.r[0], p.r[1], p.r[2], p.r[3] = p.r[0]*tpcg.P0+tpcg.P1, p.r[1]*tpcg.P0+tpcg.P1, p.r[2]*tpcg.P0+tpcg.P1, p.r[3]*tpcg.P0+tpcg.P1
	p.r[4], p.r[5], p.r[6], p.r[7] = p.r[4]*tpcg.P0+tpcg.P1, p.r[5]*tpcg.P0+tpcg.P1, p.r[6]*tpcg.P0+tpcg.P1, p.r[7]*tpcg.P0+tpcg.P1

}

//go:fix inline
func (p *HighPipeline) GammaExpand2() {
	expand := func(x float32) float32 {
		return x * x
	}

	p.r[0], p.r[1], p.r[2], p.r[3] = expand(p.r[0]), expand(p.r[1]), expand(p.r[2]), expand(p.r[3])
	p.r[4], p.r[5], p.r[6], p.r[7] = expand(p.r[4]), expand(p.r[5]), expand(p.r[6]), expand(p.r[7])
	p.g[0], p.g[1], p.g[2], p.g[3] = expand(p.g[0]), expand(p.g[1]), expand(p.g[2]), expand(p.g[3])
	p.g[4], p.g[5], p.g[6], p.g[7] = expand(p.g[4]), expand(p.g[5]), expand(p.g[6]), expand(p.g[7])
	p.b[0], p.b[1], p.b[2], p.b[3] = expand(p.b[0]), expand(p.b[1]), expand(p.b[2]), expand(p.b[3])
	p.b[4], p.b[5], p.b[6], p.b[7] = expand(p.b[4]), expand(p.b[5]), expand(p.b[6]), expand(p.b[7])

}

//go:fix inline
func (p *HighPipeline) GammaExpandDestination2() {
	expand := func(x float32) float32 {
		return x * x
	}

	p.dr[0], p.dr[1], p.dr[2], p.dr[3] = expand(p.dr[0]), expand(p.dr[1]), expand(p.dr[2]), expand(p.dr[3])
	p.dr[4], p.dr[5], p.dr[6], p.dr[7] = expand(p.dr[4]), expand(p.dr[5]), expand(p.dr[6]), expand(p.dr[7])
	p.dg[0], p.dg[1], p.dg[2], p.dg[3] = expand(p.dg[0]), expand(p.dg[1]), expand(p.dg[2]), expand(p.dg[3])
	p.dg[4], p.dg[5], p.dg[6], p.dg[7] = expand(p.dg[4]), expand(p.dg[5]), expand(p.dg[6]), expand(p.dg[7])
	p.db[0], p.db[1], p.db[2], p.db[3] = expand(p.db[0]), expand(p.db[1]), expand(p.db[2]), expand(p.db[3])
	p.db[4], p.db[5], p.db[6], p.db[7] = expand(p.db[4]), expand(p.db[5]), expand(p.db[6]), expand(p.db[7])

}

//go:fix inline
func (p *HighPipeline) GammaCompress2() {
	compress := func(x float32) float32 {
		return math32.Sqrt(x)
	}

	p.r[0], p.r[1], p.r[2], p.r[3] = compress(p.r[0]), compress(p.r[1]), compress(p.r[2]), compress(p.r[3])
	p.r[4], p.r[5], p.r[6], p.r[7] = compress(p.r[4]), compress(p.r[5]), compress(p.r[6]), compress(p.r[7])
	p.g[0], p.g[1], p.g[2], p.g[3] = compress(p.g[0]), compress(p.g[1]), compress(p.g[2]), compress(p.g[3])
	p.g[4], p.g[5], p.g[6], p.g[7] = compress(p.g[4]), compress(p.g[5]), compress(p.g[6]), compress(p.g[7])
	p.b[0], p.b[1], p.b[2], p.b[3] = compress(p.b[0]), compress(p.b[1]), compress(p.b[2]), compress(p.b[3])
	p.b[4], p.b[5], p.b[6], p.b[7] = compress(p.b[4]), compress(p.b[5]), compress(p.b[6]), compress(p.b[7])

}

//go:fix inline
func (p *HighPipeline) GammaExpand22() {
	expand := func(x float32) float32 {
		return math32.Pow(x, 2.2)
	}

	p.r[0], p.r[1], p.r[2], p.r[3] = expand(p.r[0]), expand(p.r[1]), expand(p.r[2]), expand(p.r[3])
	p.r[4], p.r[5], p.r[6], p.r[7] = expand(p.r[4]), expand(p.r[5]), expand(p.r[6]), expand(p.r[7])
	p.g[0], p.g[1], p.g[2], p.g[3] = expand(p.g[0]), expand(p.g[1]), expand(p.g[2]), expand(p.g[3])
	p.g[4], p.g[5], p.g[6], p.g[7] = expand(p.g[4]), expand(p.g[5]), expand(p.g[6]), expand(p.g[7])
	p.b[0], p.b[1], p.b[2], p.b[3] = expand(p.b[0]), expand(p.b[1]), expand(p.b[2]), expand(p.b[3])
	p.b[4], p.b[5], p.b[6], p.b[7] = expand(p.b[4]), expand(p.b[5]), expand(p.b[6]), expand(p.b[7])

}

//go:fix inline
func (p *HighPipeline) GammaExpandDestination22() {
	expand := func(x float32) float32 {
		return math32.Pow(x, 2.2)
	}

	p.dr[0], p.dr[1], p.dr[2], p.dr[3] = expand(p.dr[0]), expand(p.dr[1]), expand(p.dr[2]), expand(p.dr[3])
	p.dr[4], p.dr[5], p.dr[6], p.dr[7] = expand(p.dr[4]), expand(p.dr[5]), expand(p.dr[6]), expand(p.dr[7])
	p.dg[0], p.dg[1], p.dg[2], p.dg[3] = expand(p.dg[0]), expand(p.dg[1]), expand(p.dg[2]), expand(p.dg[3])
	p.dg[4], p.dg[5], p.dg[6], p.dg[7] = expand(p.dg[4]), expand(p.dg[5]), expand(p.dg[6]), expand(p.dg[7])
	p.db[0], p.db[1], p.db[2], p.db[3] = expand(p.db[0]), expand(p.db[1]), expand(p.db[2]), expand(p.db[3])
	p.db[4], p.db[5], p.db[6], p.db[7] = expand(p.db[4]), expand(p.db[5]), expand(p.db[6]), expand(p.db[7])

}

//go:fix inline
func (p *HighPipeline) GammaCompress22() {
	compress := func(x float32) float32 {
		return math32.Pow(x, 0.45454545)
	}

	p.r[0], p.r[1], p.r[2], p.r[3] = compress(p.r[0]), compress(p.r[1]), compress(p.r[2]), compress(p.r[3])
	p.r[4], p.r[5], p.r[6], p.r[7] = compress(p.r[4]), compress(p.r[5]), compress(p.r[6]), compress(p.r[7])
	p.g[0], p.g[1], p.g[2], p.g[3] = compress(p.g[0]), compress(p.g[1]), compress(p.g[2]), compress(p.g[3])
	p.g[4], p.g[5], p.g[6], p.g[7] = compress(p.g[4]), compress(p.g[5]), compress(p.g[6]), compress(p.g[7])
	p.b[0], p.b[1], p.b[2], p.b[3] = compress(p.b[0]), compress(p.b[1]), compress(p.b[2]), compress(p.b[3])
	p.b[4], p.b[5], p.b[6], p.b[7] = compress(p.b[4]), compress(p.b[5]), compress(p.b[6]), compress(p.b[7])

}

//go:fix inline
func (p *HighPipeline) GammaExpandSrgb() {
	expand := func(x float32) float32 {
		if x <= 0.04045 {
			return x / 12.92
		}
		return float32(math32.Pow((x+0.055)/1.055, 2.4))
	}

	p.r[0], p.r[1], p.r[2], p.r[3] = expand(p.r[0]), expand(p.r[1]), expand(p.r[2]), expand(p.r[3])
	p.r[4], p.r[5], p.r[6], p.r[7] = expand(p.r[4]), expand(p.r[5]), expand(p.r[6]), expand(p.r[7])
	p.g[0], p.g[1], p.g[2], p.g[3] = expand(p.g[0]), expand(p.g[1]), expand(p.g[2]), expand(p.g[3])
	p.g[4], p.g[5], p.g[6], p.g[7] = expand(p.g[4]), expand(p.g[5]), expand(p.g[6]), expand(p.g[7])
	p.b[0], p.b[1], p.b[2], p.b[3] = expand(p.b[0]), expand(p.b[1]), expand(p.b[2]), expand(p.b[3])
	p.b[4], p.b[5], p.b[6], p.b[7] = expand(p.b[4]), expand(p.b[5]), expand(p.b[6]), expand(p.b[7])

}

//go:fix inline
func (p *HighPipeline) GammaExpandDestinationSrgb() {
	expand := func(x float32) float32 {
		if x <= 0.04045 {
			return x / 12.92
		}
		return float32(math32.Pow((x+0.055)/1.055, 2.4))
	}

	p.dr[0], p.dr[1], p.dr[2], p.dr[3] = expand(p.dr[0]), expand(p.dr[1]), expand(p.dr[2]), expand(p.dr[3])
	p.dr[4], p.dr[5], p.dr[6], p.dr[7] = expand(p.dr[4]), expand(p.dr[5]), expand(p.dr[6]), expand(p.dr[7])
	p.dg[0], p.dg[1], p.dg[2], p.dg[3] = expand(p.dg[0]), expand(p.dg[1]), expand(p.dg[2]), expand(p.dg[3])
	p.dg[4], p.dg[5], p.dg[6], p.dg[7] = expand(p.dg[4]), expand(p.dg[5]), expand(p.dg[6]), expand(p.dg[7])
	p.db[0], p.db[1], p.db[2], p.db[3] = expand(p.db[0]), expand(p.db[1]), expand(p.db[2]), expand(p.db[3])
	p.db[4], p.db[5], p.db[6], p.db[7] = expand(p.db[4]), expand(p.db[5]), expand(p.db[6]), expand(p.db[7])

}

//go:fix inline
func (p *HighPipeline) GammaCompressSrgb() {
	compress := func(x float32) float32 {
		if x <= 0.0031308 {
			return x * 12.92
		}
		return float32(math32.Pow(x, 0.416666666))*1.055 - 0.055
	}

	p.r[0], p.r[1], p.r[2], p.r[3] = compress(p.r[0]), compress(p.r[1]), compress(p.r[2]), compress(p.r[3])
	p.r[4], p.r[5], p.r[6], p.r[7] = compress(p.r[4]), compress(p.r[5]), compress(p.r[6]), compress(p.r[7])
	p.g[0], p.g[1], p.g[2], p.g[3] = compress(p.g[0]), compress(p.g[1]), compress(p.g[2]), compress(p.g[3])
	p.g[4], p.g[5], p.g[6], p.g[7] = compress(p.g[4]), compress(p.g[5]), compress(p.g[6]), compress(p.g[7])
	p.b[0], p.b[1], p.b[2], p.b[3] = compress(p.b[0]), compress(p.b[1]), compress(p.b[2]), compress(p.b[3])
	p.b[4], p.b[5], p.b[6], p.b[7] = compress(p.b[4]), compress(p.b[5]), compress(p.b[6]), compress(p.b[7])
}

func f32lerp(from, to, t float32) float32 {
	return from + t*(to-from)
}

func f32abs(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}

func f32normalize(v float32) float32 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func recipFast(v float32) float32 {
	if v == 0 {
		return 0
	}
	return 1.0 / v
}

func f32unnorm(x float32) uint8 {
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 255
	}
	return uint8(x*255.0 + 0.5)
}

func f32max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func f32min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func f32minmax(r, g, b float32) (float32, float32) {
	mn, mx := r, r
	if g < mn {
		mn = g
	} else if g > mx {
		mx = g
	}
	if b < mn {
		mn = b
	} else if b > mx {
		mx = b
	}
	return mn, mx
}

func f32sat(r, g, b, a, mna, mxa, mnb, mxb float32) (float32, float32, float32) {
	sat := mxa - mna
	if sat == 0 {
		return 0, 0, 0
	}
	invSat := (mxb - mnb) * a / sat
	return (r - mna) * invSat, (g - mna) * invSat, (b - mna) * invSat
}

func f32lum(r, g, b, a, luma, lumb float32) (float32, float32, float32) {
	diff := lumb*a - luma
	return r + diff, g + diff, b + diff
}

func f32clip(r, g, b, a, mn, mx, lum float32) (float32, float32, float32) {
	s1 := float32(1)
	if mn < 0 && lum != mn {
		s1 = lum / (lum - mn)
	}
	s2 := float32(1)
	if mx > a && mx != lum {
		s2 = (a - lum) / (mx - lum)
	}
	scale := s1 * s2
	r = lum + (r-lum)*scale
	g = lum + (g-lum)*scale
	b = lum + (b-lum)*scale
	return r, g, b
}

func exclusiveReflect(v, limit, invLimit float32) float32 {
	return f32abs((v - limit) - (limit+limit)*math32.Floor((v-limit)*(invLimit*0.5)) - limit)
}

func samplePixel(pixmap *PixmapCtx, ctx *SamplerCtx, x, y float32, r, g, b, a *float32) {
	width := float32(pixmap.Size.Width())
	height := float32(pixmap.Size.Height())
	switch ctx.SpreadMode {
	case 0:
	case 1:
		x = x - float32(math.Floor(float64(x*width)))*ctx.InvWidth
		y = y - float32(math.Floor(float64(y*height)))*ctx.InvHeight
	case 2:
		x = exclusiveReflect(x, width, ctx.InvWidth)
		y = exclusiveReflect(y, height, ctx.InvHeight)
	}
	ix := int(x)
	iy := int(y)
	if ix >= 0 && ix < int(width) && iy >= 0 && iy < int(height) {
		idx := (iy*int(width) + ix) * 4
		*r = float32(pixmap.Data[idx]) / 255
		*g = float32(pixmap.Data[idx+1]) / 255
		*b = float32(pixmap.Data[idx+2]) / 255
		*a = float32(pixmap.Data[idx+3]) / 255
	}
}
