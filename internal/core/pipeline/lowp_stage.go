// Copyright 2016 Google Inc.
// Copyright 2020 Yevhenii Reizner
// Copyright 2026 LumiFloat
//
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
package pipeline

import (
	"unsafe"

	"github.com/chewxy/math32"
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
	p.r[0], p.r[1], p.r[2], p.r[3] = (p.r[0]*p.a[0]+255)>>8, (p.r[1]*p.a[1]+255)>>8, (p.r[2]*p.a[2]+255)>>8, (p.r[3]*p.a[3]+255)>>8
	p.r[4], p.r[5], p.r[6], p.r[7] = (p.r[4]*p.a[4]+255)>>8, (p.r[5]*p.a[5]+255)>>8, (p.r[6]*p.a[6]+255)>>8, (p.r[7]*p.a[7]+255)>>8
	p.r[8], p.r[9], p.r[10], p.r[11] = (p.r[8]*p.a[8]+255)>>8, (p.r[9]*p.a[9]+255)>>8, (p.r[10]*p.a[10]+255)>>8, (p.r[11]*p.a[11]+255)>>8
	p.r[12], p.r[13], p.r[14], p.r[15] = (p.r[12]*p.a[12]+255)>>8, (p.r[13]*p.a[13]+255)>>8, (p.r[14]*p.a[14]+255)>>8, (p.r[15]*p.a[15]+255)>>8

	p.g[0], p.g[1], p.g[2], p.g[3] = (p.g[0]*p.a[0]+255)>>8, (p.g[1]*p.a[1]+255)>>8, (p.g[2]*p.a[2]+255)>>8, (p.g[3]*p.a[3]+255)>>8
	p.g[4], p.g[5], p.g[6], p.g[7] = (p.g[4]*p.a[4]+255)>>8, (p.g[5]*p.a[5]+255)>>8, (p.g[6]*p.a[6]+255)>>8, (p.g[7]*p.a[7]+255)>>8
	p.g[8], p.g[9], p.g[10], p.g[11] = (p.g[8]*p.a[8]+255)>>8, (p.g[9]*p.a[9]+255)>>8, (p.g[10]*p.a[10]+255)>>8, (p.g[11]*p.a[11]+255)>>8
	p.g[12], p.g[13], p.g[14], p.g[15] = (p.g[12]*p.a[12]+255)>>8, (p.g[13]*p.a[13]+255)>>8, (p.g[14]*p.a[14]+255)>>8, (p.g[15]*p.a[15]+255)>>8

	p.b[0], p.b[1], p.b[2], p.b[3] = (p.b[0]*p.a[0]+255)>>8, (p.b[1]*p.a[1]+255)>>8, (p.b[2]*p.a[2]+255)>>8, (p.b[3]*p.a[3]+255)>>8
	p.b[4], p.b[5], p.b[6], p.b[7] = (p.b[4]*p.a[4]+255)>>8, (p.b[5]*p.a[5]+255)>>8, (p.b[6]*p.a[6]+255)>>8, (p.b[7]*p.a[7]+255)>>8
	p.b[8], p.b[9], p.b[10], p.b[11] = (p.b[8]*p.a[8]+255)>>8, (p.b[9]*p.a[9]+255)>>8, (p.b[10]*p.a[10]+255)>>8, (p.b[11]*p.a[11]+255)>>8
	p.b[12], p.b[13], p.b[14], p.b[15] = (p.b[12]*p.a[12]+255)>>8, (p.b[13]*p.a[13]+255)>>8, (p.b[14]*p.a[14]+255)>>8, (p.b[15]*p.a[15]+255)>>8
}

//go:fix inline
func (p *LowPipeline) UniformColor() {
	uniformColor := p.ctx.UniformColor
	r := uniformColor.RGBA[0]
	g := uniformColor.RGBA[1]
	b := uniformColor.RGBA[2]
	a := uniformColor.RGBA[3]

	p.r[0], p.r[1], p.r[2], p.r[3] = r, r, r, r
	p.r[4], p.r[5], p.r[6], p.r[7] = r, r, r, r
	p.r[8], p.r[9], p.r[10], p.r[11] = r, r, r, r
	p.r[12], p.r[13], p.r[14], p.r[15] = r, r, r, r
	p.g[0], p.g[1], p.g[2], p.g[3] = g, g, g, g
	p.g[4], p.g[5], p.g[6], p.g[7] = g, g, g, g
	p.g[8], p.g[9], p.g[10], p.g[11] = g, g, g, g
	p.g[12], p.g[13], p.g[14], p.g[15] = g, g, g, g
	p.b[0], p.b[1], p.b[2], p.b[3] = b, b, b, b
	p.b[4], p.b[5], p.b[6], p.b[7] = b, b, b, b
	p.b[8], p.b[9], p.b[10], p.b[11] = b, b, b, b
	p.b[12], p.b[13], p.b[14], p.b[15] = b, b, b, b
	p.a[0], p.a[1], p.a[2], p.a[3] = a, a, a, a
	p.a[4], p.a[5], p.a[6], p.a[7] = a, a, a, a
	p.a[8], p.a[9], p.a[10], p.a[11] = a, a, a, a
	p.a[12], p.a[13], p.a[14], p.a[15] = a, a, a, a
}

//go:fix inline
func (p *LowPipeline) SeedShader() {
	dxFloat := float32(p.dx)
	dyFloat := float32(p.dy) + 0.5

	// Calculate x coordinates: dx + [0.5, 1.5, ..., 15.5]
	x := [16]float32{
		dxFloat + 0.5, dxFloat + 1.5, dxFloat + 2.5, dxFloat + 3.5,
		dxFloat + 4.5, dxFloat + 5.5, dxFloat + 6.5, dxFloat + 7.5,
		dxFloat + 8.5, dxFloat + 9.5, dxFloat + 10.5, dxFloat + 11.5,
		dxFloat + 12.5, dxFloat + 13.5, dxFloat + 14.5, dxFloat + 15.5,
	}

	// Convert float32 array to uint16 arrays using unsafe pointer casting
	// This matches Rust's split function: each f32 is split into two u16
	// The memory layout is preserved exactly (platform-independent)
	xPtr := (*[16][2]uint16)(unsafe.Pointer(&x[0]))
	for i := 0; i < 16; i++ {
		p.r[i] = (*xPtr)[i][0] // Low 16 bits of float32
		p.g[i] = (*xPtr)[i][1] // High 16 bits of float32
	}

	// Calculate y coordinates: dy + 0.5 (constant for all pixels)
	y := [16]float32{
		dyFloat, dyFloat, dyFloat, dyFloat,
		dyFloat, dyFloat, dyFloat, dyFloat,
		dyFloat, dyFloat, dyFloat, dyFloat,
		dyFloat, dyFloat, dyFloat, dyFloat,
	}

	// Convert float32 array to uint16 arrays using same approach
	yPtr := (*[16][2]uint16)(unsafe.Pointer(&y[0]))
	for i := 0; i < 16; i++ {
		p.b[i] = (*yPtr)[i][0] // Low 16 bits of float32
		p.a[i] = (*yPtr)[i][1] // High 16 bits of float32
	}
}

//go:fix inline
func (p *LowPipeline) LoadDestination() {
	baseIdx := (p.dy*p.pixmap.RealWidth + p.dx) * 4
	data := p.pixmap.Data[baseIdx : baseIdx+LOW_STAGE_WIDTH*4]

	p.dr[0], p.dr[1], p.dr[2], p.dr[3] = uint16(data[0]), uint16(data[4]), uint16(data[8]), uint16(data[12])
	p.dr[4], p.dr[5], p.dr[6], p.dr[7] = uint16(data[16]), uint16(data[20]), uint16(data[24]), uint16(data[28])
	p.dg[0], p.dg[1], p.dg[2], p.dg[3] = uint16(data[1]), uint16(data[5]), uint16(data[9]), uint16(data[13])
	p.dg[4], p.dg[5], p.dg[6], p.dg[7] = uint16(data[17]), uint16(data[21]), uint16(data[25]), uint16(data[29])
	p.db[0], p.db[1], p.db[2], p.db[3] = uint16(data[2]), uint16(data[6]), uint16(data[10]), uint16(data[14])
	p.db[4], p.db[5], p.db[6], p.db[7] = uint16(data[18]), uint16(data[22]), uint16(data[26]), uint16(data[30])
	p.da[0], p.da[1], p.da[2], p.da[3] = uint16(data[3]), uint16(data[7]), uint16(data[11]), uint16(data[15])
	p.da[4], p.da[5], p.da[6], p.da[7] = uint16(data[19]), uint16(data[23]), uint16(data[27]), uint16(data[31])

	p.dr[8], p.dr[9], p.dr[10], p.dr[11] = uint16(data[32]), uint16(data[36]), uint16(data[40]), uint16(data[44])
	p.dr[12], p.dr[13], p.dr[14], p.dr[15] = uint16(data[48]), uint16(data[52]), uint16(data[56]), uint16(data[60])
	p.dg[8], p.dg[9], p.dg[10], p.dg[11] = uint16(data[33]), uint16(data[37]), uint16(data[41]), uint16(data[45])
	p.dg[12], p.dg[13], p.dg[14], p.dg[15] = uint16(data[49]), uint16(data[53]), uint16(data[57]), uint16(data[61])
	p.db[8], p.db[9], p.db[10], p.db[11] = uint16(data[34]), uint16(data[38]), uint16(data[42]), uint16(data[46])
	p.db[12], p.db[13], p.db[14], p.db[15] = uint16(data[50]), uint16(data[54]), uint16(data[58]), uint16(data[62])
	p.da[8], p.da[9], p.da[10], p.da[11] = uint16(data[35]), uint16(data[39]), uint16(data[43]), uint16(data[47])
	p.da[12], p.da[13], p.da[14], p.da[15] = uint16(data[51]), uint16(data[55]), uint16(data[59]), uint16(data[63])
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
	data := p.pixmap.Data[baseIdx : baseIdx+LOW_STAGE_WIDTH*4]

	data[0], data[4], data[8], data[12] = uint8(p.r[0]), uint8(p.r[1]), uint8(p.r[2]), uint8(p.r[3])
	data[16], data[20], data[24], data[28] = uint8(p.r[4]), uint8(p.r[5]), uint8(p.r[6]), uint8(p.r[7])
	data[1], data[5], data[9], data[13] = uint8(p.g[0]), uint8(p.g[1]), uint8(p.g[2]), uint8(p.g[3])
	data[17], data[21], data[25], data[29] = uint8(p.g[4]), uint8(p.g[5]), uint8(p.g[6]), uint8(p.g[7])
	data[2], data[6], data[10], data[14] = uint8(p.b[0]), uint8(p.b[1]), uint8(p.b[2]), uint8(p.b[3])
	data[18], data[22], data[26], data[30] = uint8(p.b[4]), uint8(p.b[5]), uint8(p.b[6]), uint8(p.b[7])
	data[3], data[7], data[11], data[15] = uint8(p.a[0]), uint8(p.a[1]), uint8(p.a[2]), uint8(p.a[3])
	data[19], data[23], data[27], data[31] = uint8(p.a[4]), uint8(p.a[5]), uint8(p.a[6]), uint8(p.a[7])

	data[32], data[36], data[40], data[44] = uint8(p.r[8]), uint8(p.r[9]), uint8(p.r[10]), uint8(p.r[11])
	data[48], data[52], data[56], data[60] = uint8(p.r[12]), uint8(p.r[13]), uint8(p.r[14]), uint8(p.r[15])
	data[33], data[37], data[41], data[45] = uint8(p.g[8]), uint8(p.g[9]), uint8(p.g[10]), uint8(p.g[11])
	data[49], data[53], data[57], data[61] = uint8(p.g[12]), uint8(p.g[13]), uint8(p.g[14]), uint8(p.g[15])
	data[34], data[38], data[42], data[46] = uint8(p.b[8]), uint8(p.b[9]), uint8(p.b[10]), uint8(p.b[11])
	data[50], data[54], data[58], data[62] = uint8(p.b[12]), uint8(p.b[13]), uint8(p.b[14]), uint8(p.b[15])
	data[35], data[39], data[43], data[47] = uint8(p.a[8]), uint8(p.a[9]), uint8(p.a[10]), uint8(p.a[11])
	data[51], data[55], data[59], data[63] = uint8(p.a[12]), uint8(p.a[13]), uint8(p.a[14]), uint8(p.a[15])
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

	p.da[0], p.da[1], p.da[2], p.da[3] = uint16(data[3]), uint16(data[7]), uint16(data[11]), uint16(data[15])
	p.da[4], p.da[5], p.da[6], p.da[7] = uint16(data[19]), uint16(data[23]), uint16(data[27]), uint16(data[31])
	p.da[8], p.da[9], p.da[10], p.da[11] = uint16(data[35]), uint16(data[39]), uint16(data[43]), uint16(data[47])
	p.da[12], p.da[13], p.da[14], p.da[15] = uint16(data[51]), uint16(data[55]), uint16(data[59]), uint16(data[63])
}

//go:fix inline
func (p *LowPipeline) LoadDestinationU8Tail() {
	baseIdx := (p.dy*p.pixmap.RealWidth + p.dx) * 4
	data := p.pixmap.Data[baseIdx : baseIdx+p.tail*4]

	for i := 0; i < p.tail; i++ {
		off := i * 4
		p.da[i] = uint16(data[off+3])
	}
}

//go:fix inline
func (p *LowPipeline) StoreU8() {
	baseIdx := (p.dy*p.pixmap.RealWidth + p.dx) * 4
	data := p.pixmap.Data[baseIdx : baseIdx+LOW_STAGE_WIDTH]

	data[0], data[1], data[2], data[3] = uint8(p.a[0]), uint8(p.a[1]), uint8(p.a[2]), uint8(p.a[3])
	data[4], data[5], data[6], data[7] = uint8(p.a[4]), uint8(p.a[5]), uint8(p.a[6]), uint8(p.a[7])
	data[8], data[9], data[10], data[11] = uint8(p.a[8]), uint8(p.a[9]), uint8(p.a[10]), uint8(p.a[11])
	data[12], data[13], data[14], data[15] = uint8(p.a[12]), uint8(p.a[13]), uint8(p.a[14]), uint8(p.a[15])
}

//go:fix inline
func (p *LowPipeline) StoreU8Tail() {
	baseIdx := (p.dy*p.pixmap.RealWidth + p.dx) * 4
	data := p.pixmap.Data[baseIdx : baseIdx+p.tail]

	for i := 0; i < p.tail; i++ {
		data[i] = uint8(p.a[i])
	}
}

//go:fix inline
func (p *LowPipeline) Gather() {
}

//go:fix inline
func (p *LowPipeline) LoadMaskU8() {
	baseIdx := int(p.maskCtx.RealWidth)*p.dy + p.dx
	maskData := p.maskCtx.Data

	var c [LOW_STAGE_WIDTH]uint16
	for i := 0; i < p.tail; i++ {
		c[i] = uint16(maskData[baseIdx+i])
	}

	p.r = [LOW_STAGE_WIDTH]uint16{}
	p.g = [LOW_STAGE_WIDTH]uint16{}
	p.b = [LOW_STAGE_WIDTH]uint16{}
	p.a = c
}

//go:fix inline
func (p *LowPipeline) MaskU8() {
	baseIdx := int(p.maskCtx.RealWidth)*p.dy + p.dx
	maskData := p.maskCtx.Data

	var c [LOW_STAGE_WIDTH]uint16
	for i := 0; i < p.tail; i++ {
		c[i] = uint16(maskData[baseIdx+i])
	}

	p.r[0], p.r[1], p.r[2], p.r[3] = (p.r[0]*c[0]+255)>>8, (p.r[1]*c[1]+255)>>8, (p.r[2]*c[2]+255)>>8, (p.r[3]*c[3]+255)>>8
	p.r[4], p.r[5], p.r[6], p.r[7] = (p.r[4]*c[4]+255)>>8, (p.r[5]*c[5]+255)>>8, (p.r[6]*c[6]+255)>>8, (p.r[7]*c[7]+255)>>8
	p.r[8], p.r[9], p.r[10], p.r[11] = (p.r[8]*c[8]+255)>>8, (p.r[9]*c[9]+255)>>8, (p.r[10]*c[10]+255)>>8, (p.r[11]*c[11]+255)>>8
	p.r[12], p.r[13], p.r[14], p.r[15] = (p.r[12]*c[12]+255)>>8, (p.r[13]*c[13]+255)>>8, (p.r[14]*c[14]+255)>>8, (p.r[15]*c[15]+255)>>8

	p.g[0], p.g[1], p.g[2], p.g[3] = (p.g[0]*c[0]+255)>>8, (p.g[1]*c[1]+255)>>8, (p.g[2]*c[2]+255)>>8, (p.g[3]*c[3]+255)>>8
	p.g[4], p.g[5], p.g[6], p.g[7] = (p.g[4]*c[4]+255)>>8, (p.g[5]*c[5]+255)>>8, (p.g[6]*c[6]+255)>>8, (p.g[7]*c[7]+255)>>8
	p.g[8], p.g[9], p.g[10], p.g[11] = (p.g[8]*c[8]+255)>>8, (p.g[9]*c[9]+255)>>8, (p.g[10]*c[10]+255)>>8, (p.g[11]*c[11]+255)>>8
	p.g[12], p.g[13], p.g[14], p.g[15] = (p.g[12]*c[12]+255)>>8, (p.g[13]*c[13]+255)>>8, (p.g[14]*c[14]+255)>>8, (p.g[15]*c[15]+255)>>8

	p.b[0], p.b[1], p.b[2], p.b[3] = (p.b[0]*c[0]+255)>>8, (p.b[1]*c[1]+255)>>8, (p.b[2]*c[2]+255)>>8, (p.b[3]*c[3]+255)>>8
	p.b[4], p.b[5], p.b[6], p.b[7] = (p.b[4]*c[4]+255)>>8, (p.b[5]*c[5]+255)>>8, (p.b[6]*c[6]+255)>>8, (p.b[7]*c[7]+255)>>8
	p.b[8], p.b[9], p.b[10], p.b[11] = (p.b[8]*c[8]+255)>>8, (p.b[9]*c[9]+255)>>8, (p.b[10]*c[10]+255)>>8, (p.b[11]*c[11]+255)>>8
	p.b[12], p.b[13], p.b[14], p.b[15] = (p.b[12]*c[12]+255)>>8, (p.b[13]*c[13]+255)>>8, (p.b[14]*c[14]+255)>>8, (p.b[15]*c[15]+255)>>8

	p.a[0], p.a[1], p.a[2], p.a[3] = (p.a[0]*c[0]+255)>>8, (p.a[1]*c[1]+255)>>8, (p.a[2]*c[2]+255)>>8, (p.a[3]*c[3]+255)>>8
	p.a[4], p.a[5], p.a[6], p.a[7] = (p.a[4]*c[4]+255)>>8, (p.a[5]*c[5]+255)>>8, (p.a[6]*c[6]+255)>>8, (p.a[7]*c[7]+255)>>8
	p.a[8], p.a[9], p.a[10], p.a[11] = (p.a[8]*c[8]+255)>>8, (p.a[9]*c[9]+255)>>8, (p.a[10]*c[10]+255)>>8, (p.a[11]*c[11]+255)>>8
	p.a[12], p.a[13], p.a[14], p.a[15] = (p.a[12]*c[12]+255)>>8, (p.a[13]*c[13]+255)>>8, (p.a[14]*c[14]+255)>>8, (p.a[15]*c[15]+255)>>8
}

//go:fix inline
func (p *LowPipeline) ScaleU8() {
	baseIdx := int(p.maskCtx.RealWidth)*p.dy + p.dx
	maskData := p.maskCtx.Data

	c0 := uint16(maskData[baseIdx])
	c1 := uint16(maskData[baseIdx+1])

	p.r[0], p.r[1], p.r[2], p.r[3] = (p.r[0]*c0+255)>>8, (p.r[1]*c1+255)>>8, 0, 0
	p.r[4], p.r[5], p.r[6], p.r[7] = 0, 0, 0, 0
	p.r[8], p.r[9], p.r[10], p.r[11] = 0, 0, 0, 0
	p.r[12], p.r[13], p.r[14], p.r[15] = 0, 0, 0, 0

	p.g[0], p.g[1], p.g[2], p.g[3] = (p.g[0]*c0+255)>>8, (p.g[1]*c1+255)>>8, 0, 0
	p.g[4], p.g[5], p.g[6], p.g[7] = 0, 0, 0, 0
	p.g[8], p.g[9], p.g[10], p.g[11] = 0, 0, 0, 0
	p.g[12], p.g[13], p.g[14], p.g[15] = 0, 0, 0, 0

	p.b[0], p.b[1], p.b[2], p.b[3] = (p.b[0]*c0+255)>>8, (p.b[1]*c1+255)>>8, 0, 0
	p.b[4], p.b[5], p.b[6], p.b[7] = 0, 0, 0, 0
	p.b[8], p.b[9], p.b[10], p.b[11] = 0, 0, 0, 0
	p.b[12], p.b[13], p.b[14], p.b[15] = 0, 0, 0, 0

	p.a[0], p.a[1], p.a[2], p.a[3] = (p.a[0]*c0+255)>>8, (p.a[1]*c1+255)>>8, 0, 0
	p.a[4], p.a[5], p.a[6], p.a[7] = 0, 0, 0, 0
	p.a[8], p.a[9], p.a[10], p.a[11] = 0, 0, 0, 0
	p.a[12], p.a[13], p.a[14], p.a[15] = 0, 0, 0, 0
}

//go:fix inline
func (p *LowPipeline) LerpU8() {
	baseIdx := int(p.maskCtx.RealWidth)*p.dy + p.dx
	maskData := p.maskCtx.Data

	c0 := uint16(maskData[baseIdx])
	c1 := uint16(maskData[baseIdx+1])
	invC0, invC1 := 255-c0, 255-c1

	p.r[0], p.r[1], p.r[2], p.r[3] = (p.dr[0]*invC0+p.r[0]*c0+255)>>8, (p.dr[1]*invC1+p.r[1]*c1+255)>>8, 0, 0
	p.r[4], p.r[5], p.r[6], p.r[7] = 0, 0, 0, 0
	p.r[8], p.r[9], p.r[10], p.r[11] = 0, 0, 0, 0
	p.r[12], p.r[13], p.r[14], p.r[15] = 0, 0, 0, 0

	p.g[0], p.g[1], p.g[2], p.g[3] = (p.dg[0]*invC0+p.g[0]*c0+255)>>8, (p.dg[1]*invC1+p.g[1]*c1+255)>>8, 0, 0
	p.g[4], p.g[5], p.g[6], p.g[7] = 0, 0, 0, 0
	p.g[8], p.g[9], p.g[10], p.g[11] = 0, 0, 0, 0
	p.g[12], p.g[13], p.g[14], p.g[15] = 0, 0, 0, 0

	p.b[0], p.b[1], p.b[2], p.b[3] = (p.db[0]*invC0+p.b[0]*c0+255)>>8, (p.db[1]*invC1+p.b[1]*c1+255)>>8, 0, 0
	p.b[4], p.b[5], p.b[6], p.b[7] = 0, 0, 0, 0
	p.b[8], p.b[9], p.b[10], p.b[11] = 0, 0, 0, 0
	p.b[12], p.b[13], p.b[14], p.b[15] = 0, 0, 0, 0

	p.a[0], p.a[1], p.a[2], p.a[3] = (p.da[0]*invC0+p.a[0]*c0+255)>>8, (p.da[1]*invC1+p.a[1]*c1+255)>>8, 0, 0
	p.a[4], p.a[5], p.a[6], p.a[7] = 0, 0, 0, 0
	p.a[8], p.a[9], p.a[10], p.a[11] = 0, 0, 0, 0
	p.a[12], p.a[13], p.a[14], p.a[15] = 0, 0, 0, 0
}

//go:fix inline
func (p *LowPipeline) Scale1Float() {
	c := uint16(p.ctx.CurrentCoverage)

	p.r[0], p.r[1], p.r[2], p.r[3] = (p.r[0]*c+255)>>8, (p.r[1]*c+255)>>8, (p.r[2]*c+255)>>8, (p.r[3]*c+255)>>8
	p.r[4], p.r[5], p.r[6], p.r[7] = (p.r[4]*c+255)>>8, (p.r[5]*c+255)>>8, (p.r[6]*c+255)>>8, (p.r[7]*c+255)>>8
	p.r[8], p.r[9], p.r[10], p.r[11] = (p.r[8]*c+255)>>8, (p.r[9]*c+255)>>8, (p.r[10]*c+255)>>8, (p.r[11]*c+255)>>8
	p.r[12], p.r[13], p.r[14], p.r[15] = (p.r[12]*c+255)>>8, (p.r[13]*c+255)>>8, (p.r[14]*c+255)>>8, (p.r[15]*c+255)>>8

	p.g[0], p.g[1], p.g[2], p.g[3] = (p.g[0]*c+255)>>8, (p.g[1]*c+255)>>8, (p.g[2]*c+255)>>8, (p.g[3]*c+255)>>8
	p.g[4], p.g[5], p.g[6], p.g[7] = (p.g[4]*c+255)>>8, (p.g[5]*c+255)>>8, (p.g[6]*c+255)>>8, (p.g[7]*c+255)>>8
	p.g[8], p.g[9], p.g[10], p.g[11] = (p.g[8]*c+255)>>8, (p.g[9]*c+255)>>8, (p.g[10]*c+255)>>8, (p.g[11]*c+255)>>8
	p.g[12], p.g[13], p.g[14], p.g[15] = (p.g[12]*c+255)>>8, (p.g[13]*c+255)>>8, (p.g[14]*c+255)>>8, (p.g[15]*c+255)>>8

	p.b[0], p.b[1], p.b[2], p.b[3] = (p.b[0]*c+255)>>8, (p.b[1]*c+255)>>8, (p.b[2]*c+255)>>8, (p.b[3]*c+255)>>8
	p.b[4], p.b[5], p.b[6], p.b[7] = (p.b[4]*c+255)>>8, (p.b[5]*c+255)>>8, (p.b[6]*c+255)>>8, (p.b[7]*c+255)>>8
	p.b[8], p.b[9], p.b[10], p.b[11] = (p.b[8]*c+255)>>8, (p.b[9]*c+255)>>8, (p.b[10]*c+255)>>8, (p.b[11]*c+255)>>8
	p.b[12], p.b[13], p.b[14], p.b[15] = (p.b[12]*c+255)>>8, (p.b[13]*c+255)>>8, (p.b[14]*c+255)>>8, (p.b[15]*c+255)>>8

	p.a[0], p.a[1], p.a[2], p.a[3] = (p.a[0]*c+255)>>8, (p.a[1]*c+255)>>8, (p.a[2]*c+255)>>8, (p.a[3]*c+255)>>8
	p.a[4], p.a[5], p.a[6], p.a[7] = (p.a[4]*c+255)>>8, (p.a[5]*c+255)>>8, (p.a[6]*c+255)>>8, (p.a[7]*c+255)>>8
	p.a[8], p.a[9], p.a[10], p.a[11] = (p.a[8]*c+255)>>8, (p.a[9]*c+255)>>8, (p.a[10]*c+255)>>8, (p.a[11]*c+255)>>8
	p.a[12], p.a[13], p.a[14], p.a[15] = (p.a[12]*c+255)>>8, (p.a[13]*c+255)>>8, (p.a[14]*c+255)>>8, (p.a[15]*c+255)>>8
}

//go:fix inline
func (p *LowPipeline) Lerp1Float() {
	c := int32(p.ctx.CurrentCoverage)
	invC := 255 - c

	p.r[0], p.r[1], p.r[2], p.r[3] = uint16((int32(p.dr[0])*invC+int32(p.r[0])*c+128)>>8), uint16((int32(p.dr[1])*invC+int32(p.r[1])*c+128)>>8), uint16((int32(p.dr[2])*invC+int32(p.r[2])*c+128)>>8), uint16((int32(p.dr[3])*invC+int32(p.r[3])*c+128)>>8)
	p.r[4], p.r[5], p.r[6], p.r[7] = uint16((int32(p.dr[4])*invC+int32(p.r[4])*c+128)>>8), uint16((int32(p.dr[5])*invC+int32(p.r[5])*c+128)>>8), uint16((int32(p.dr[6])*invC+int32(p.r[6])*c+128)>>8), uint16((int32(p.dr[7])*invC+int32(p.r[7])*c+128)>>8)
	p.r[8], p.r[9], p.r[10], p.r[11] = uint16((int32(p.dr[8])*invC+int32(p.r[8])*c+128)>>8), uint16((int32(p.dr[9])*invC+int32(p.r[9])*c+128)>>8), uint16((int32(p.dr[10])*invC+int32(p.r[10])*c+128)>>8), uint16((int32(p.dr[11])*invC+int32(p.r[11])*c+128)>>8)
	p.r[12], p.r[13], p.r[14], p.r[15] = uint16((int32(p.dr[12])*invC+int32(p.r[12])*c+128)>>8), uint16((int32(p.dr[13])*invC+int32(p.r[13])*c+128)>>8), uint16((int32(p.dr[14])*invC+int32(p.r[14])*c+128)>>8), uint16((int32(p.dr[15])*invC+int32(p.r[15])*c+128)>>8)

	p.g[0], p.g[1], p.g[2], p.g[3] = uint16((int32(p.dg[0])*invC+int32(p.g[0])*c+128)>>8), uint16((int32(p.dg[1])*invC+int32(p.g[1])*c+128)>>8), uint16((int32(p.dg[2])*invC+int32(p.g[2])*c+128)>>8), uint16((int32(p.dg[3])*invC+int32(p.g[3])*c+128)>>8)
	p.g[4], p.g[5], p.g[6], p.g[7] = uint16((int32(p.dg[4])*invC+int32(p.g[4])*c+128)>>8), uint16((int32(p.dg[5])*invC+int32(p.g[5])*c+128)>>8), uint16((int32(p.dg[6])*invC+int32(p.g[6])*c+128)>>8), uint16((int32(p.dg[7])*invC+int32(p.g[7])*c+128)>>8)
	p.g[8], p.g[9], p.g[10], p.g[11] = uint16((int32(p.dg[8])*invC+int32(p.g[8])*c+128)>>8), uint16((int32(p.dg[9])*invC+int32(p.g[9])*c+128)>>8), uint16((int32(p.dg[10])*invC+int32(p.g[10])*c+128)>>8), uint16((int32(p.dg[11])*invC+int32(p.g[11])*c+128)>>8)
	p.g[12], p.g[13], p.g[14], p.g[15] = uint16((int32(p.dg[12])*invC+int32(p.g[12])*c+128)>>8), uint16((int32(p.dg[13])*invC+int32(p.g[13])*c+128)>>8), uint16((int32(p.dg[14])*invC+int32(p.g[14])*c+128)>>8), uint16((int32(p.dg[15])*invC+int32(p.g[15])*c+128)>>8)

	p.b[0], p.b[1], p.b[2], p.b[3] = uint16((int32(p.db[0])*invC+int32(p.b[0])*c+128)>>8), uint16((int32(p.db[1])*invC+int32(p.b[1])*c+128)>>8), uint16((int32(p.db[2])*invC+int32(p.b[2])*c+128)>>8), uint16((int32(p.db[3])*invC+int32(p.b[3])*c+128)>>8)
	p.b[4], p.b[5], p.b[6], p.b[7] = uint16((int32(p.db[4])*invC+int32(p.b[4])*c+128)>>8), uint16((int32(p.db[5])*invC+int32(p.b[5])*c+128)>>8), uint16((int32(p.db[6])*invC+int32(p.b[6])*c+128)>>8), uint16((int32(p.db[7])*invC+int32(p.b[7])*c+128)>>8)
	p.b[8], p.b[9], p.b[10], p.b[11] = uint16((int32(p.db[8])*invC+int32(p.b[8])*c+128)>>8), uint16((int32(p.db[9])*invC+int32(p.b[9])*c+128)>>8), uint16((int32(p.db[10])*invC+int32(p.b[10])*c+128)>>8), uint16((int32(p.db[11])*invC+int32(p.b[11])*c+128)>>8)
	p.b[12], p.b[13], p.b[14], p.b[15] = uint16((int32(p.db[12])*invC+int32(p.b[12])*c+128)>>8), uint16((int32(p.db[13])*invC+int32(p.b[13])*c+128)>>8), uint16((int32(p.db[14])*invC+int32(p.b[14])*c+128)>>8), uint16((int32(p.db[15])*invC+int32(p.b[15])*c+128)>>8)

	p.a[0], p.a[1], p.a[2], p.a[3] = uint16((int32(p.da[0])*invC+int32(p.a[0])*c+128)>>8), uint16((int32(p.da[1])*invC+int32(p.a[1])*c+128)>>8), uint16((int32(p.da[2])*invC+int32(p.a[2])*c+128)>>8), uint16((int32(p.da[3])*invC+int32(p.a[3])*c+128)>>8)
	p.a[4], p.a[5], p.a[6], p.a[7] = uint16((int32(p.da[4])*invC+int32(p.a[4])*c+128)>>8), uint16((int32(p.da[5])*invC+int32(p.a[5])*c+128)>>8), uint16((int32(p.da[6])*invC+int32(p.a[6])*c+128)>>8), uint16((int32(p.da[7])*invC+int32(p.a[7])*c+128)>>8)
	p.a[8], p.a[9], p.a[10], p.a[11] = uint16((int32(p.da[8])*invC+int32(p.a[8])*c+128)>>8), uint16((int32(p.da[9])*invC+int32(p.a[9])*c+128)>>8), uint16((int32(p.da[10])*invC+int32(p.a[10])*c+128)>>8), uint16((int32(p.da[11])*invC+int32(p.a[11])*c+128)>>8)
	p.a[12], p.a[13], p.a[14], p.a[15] = uint16((int32(p.da[12])*invC+int32(p.a[12])*c+128)>>8), uint16((int32(p.da[13])*invC+int32(p.a[13])*c+128)>>8), uint16((int32(p.da[14])*invC+int32(p.a[14])*c+128)>>8), uint16((int32(p.da[15])*invC+int32(p.a[15])*c+128)>>8)
}

//go:fix inline
func (p *LowPipeline) DestinationAtop() {
	invDa0, invDa1, invDa2, invDa3 := 255-p.da[0], 255-p.da[1], 255-p.da[2], 255-p.da[3]
	invDa4, invDa5, invDa6, invDa7 := 255-p.da[4], 255-p.da[5], 255-p.da[6], 255-p.da[7]
	invDa8, invDa9, invDa10, invDa11 := 255-p.da[8], 255-p.da[9], 255-p.da[10], 255-p.da[11]
	invDa12, invDa13, invDa14, invDa15 := 255-p.da[12], 255-p.da[13], 255-p.da[14], 255-p.da[15]

	p.r[0], p.r[1], p.r[2], p.r[3] = uint16((uint32(p.dr[0])*uint32(p.a[0])+uint32(p.r[0])*uint32(invDa0)+255)>>8), uint16((uint32(p.dr[1])*uint32(p.a[1])+uint32(p.r[1])*uint32(invDa1)+255)>>8), uint16((uint32(p.dr[2])*uint32(p.a[2])+uint32(p.r[2])*uint32(invDa2)+255)>>8), uint16((uint32(p.dr[3])*uint32(p.a[3])+uint32(p.r[3])*uint32(invDa3)+255)>>8)
	p.r[4], p.r[5], p.r[6], p.r[7] = uint16((uint32(p.dr[4])*uint32(p.a[4])+uint32(p.r[4])*uint32(invDa4)+255)>>8), uint16((uint32(p.dr[5])*uint32(p.a[5])+uint32(p.r[5])*uint32(invDa5)+255)>>8), uint16((uint32(p.dr[6])*uint32(p.a[6])+uint32(p.r[6])*uint32(invDa6)+255)>>8), uint16((uint32(p.dr[7])*uint32(p.a[7])+uint32(p.r[7])*uint32(invDa7)+255)>>8)
	p.r[8], p.r[9], p.r[10], p.r[11] = uint16((uint32(p.dr[8])*uint32(p.a[8])+uint32(p.r[8])*uint32(invDa8)+255)>>8), uint16((uint32(p.dr[9])*uint32(p.a[9])+uint32(p.r[9])*uint32(invDa9)+255)>>8), uint16((uint32(p.dr[10])*uint32(p.a[10])+uint32(p.r[10])*uint32(invDa10)+255)>>8), uint16((uint32(p.dr[11])*uint32(p.a[11])+uint32(p.r[11])*uint32(invDa11)+255)>>8)
	p.r[12], p.r[13], p.r[14], p.r[15] = uint16((uint32(p.dr[12])*uint32(p.a[12])+uint32(p.r[12])*uint32(invDa12)+255)>>8), uint16((uint32(p.dr[13])*uint32(p.a[13])+uint32(p.r[13])*uint32(invDa13)+255)>>8), uint16((uint32(p.dr[14])*uint32(p.a[14])+uint32(p.r[14])*uint32(invDa14)+255)>>8), uint16((uint32(p.dr[15])*uint32(p.a[15])+uint32(p.r[15])*uint32(invDa15)+255)>>8)

	p.g[0], p.g[1], p.g[2], p.g[3] = uint16((uint32(p.dg[0])*uint32(p.a[0])+uint32(p.g[0])*uint32(invDa0)+255)>>8), uint16((uint32(p.dg[1])*uint32(p.a[1])+uint32(p.g[1])*uint32(invDa1)+255)>>8), uint16((uint32(p.dg[2])*uint32(p.a[2])+uint32(p.g[2])*uint32(invDa2)+255)>>8), uint16((uint32(p.dg[3])*uint32(p.a[3])+uint32(p.g[3])*uint32(invDa3)+255)>>8)
	p.g[4], p.g[5], p.g[6], p.g[7] = uint16((uint32(p.dg[4])*uint32(p.a[4])+uint32(p.g[4])*uint32(invDa4)+255)>>8), uint16((uint32(p.dg[5])*uint32(p.a[5])+uint32(p.g[5])*uint32(invDa5)+255)>>8), uint16((uint32(p.dg[6])*uint32(p.a[6])+uint32(p.g[6])*uint32(invDa6)+255)>>8), uint16((uint32(p.dg[7])*uint32(p.a[7])+uint32(p.g[7])*uint32(invDa7)+255)>>8)
	p.g[8], p.g[9], p.g[10], p.g[11] = uint16((uint32(p.dg[8])*uint32(p.a[8])+uint32(p.g[8])*uint32(invDa8)+255)>>8), uint16((uint32(p.dg[9])*uint32(p.a[9])+uint32(p.g[9])*uint32(invDa9)+255)>>8), uint16((uint32(p.dg[10])*uint32(p.a[10])+uint32(p.g[10])*uint32(invDa10)+255)>>8), uint16((uint32(p.dg[11])*uint32(p.a[11])+uint32(p.g[11])*uint32(invDa11)+255)>>8)
	p.g[12], p.g[13], p.g[14], p.g[15] = uint16((uint32(p.dg[12])*uint32(p.a[12])+uint32(p.g[12])*uint32(invDa12)+255)>>8), uint16((uint32(p.dg[13])*uint32(p.a[13])+uint32(p.g[13])*uint32(invDa13)+255)>>8), uint16((uint32(p.dg[14])*uint32(p.a[14])+uint32(p.g[14])*uint32(invDa14)+255)>>8), uint16((uint32(p.dg[15])*uint32(p.a[15])+uint32(p.g[15])*uint32(invDa15)+255)>>8)

	p.b[0], p.b[1], p.b[2], p.b[3] = uint16((uint32(p.db[0])*uint32(p.a[0])+uint32(p.b[0])*uint32(invDa0)+255)>>8), uint16((uint32(p.db[1])*uint32(p.a[1])+uint32(p.b[1])*uint32(invDa1)+255)>>8), uint16((uint32(p.db[2])*uint32(p.a[2])+uint32(p.b[2])*uint32(invDa2)+255)>>8), uint16((uint32(p.db[3])*uint32(p.a[3])+uint32(p.b[3])*uint32(invDa3)+255)>>8)
	p.b[4], p.b[5], p.b[6], p.b[7] = uint16((uint32(p.db[4])*uint32(p.a[4])+uint32(p.b[4])*uint32(invDa4)+255)>>8), uint16((uint32(p.db[5])*uint32(p.a[5])+uint32(p.b[5])*uint32(invDa5)+255)>>8), uint16((uint32(p.db[6])*uint32(p.a[6])+uint32(p.b[6])*uint32(invDa6)+255)>>8), uint16((uint32(p.db[7])*uint32(p.a[7])+uint32(p.b[7])*uint32(invDa7)+255)>>8)
	p.b[8], p.b[9], p.b[10], p.b[11] = uint16((uint32(p.db[8])*uint32(p.a[8])+uint32(p.b[8])*uint32(invDa8)+255)>>8), uint16((uint32(p.db[9])*uint32(p.a[9])+uint32(p.b[9])*uint32(invDa9)+255)>>8), uint16((uint32(p.db[10])*uint32(p.a[10])+uint32(p.b[10])*uint32(invDa10)+255)>>8), uint16((uint32(p.db[11])*uint32(p.a[11])+uint32(p.b[11])*uint32(invDa11)+255)>>8)
	p.b[12], p.b[13], p.b[14], p.b[15] = uint16((uint32(p.db[12])*uint32(p.a[12])+uint32(p.b[12])*uint32(invDa12)+255)>>8), uint16((uint32(p.db[13])*uint32(p.a[13])+uint32(p.b[13])*uint32(invDa13)+255)>>8), uint16((uint32(p.db[14])*uint32(p.a[14])+uint32(p.b[14])*uint32(invDa14)+255)>>8), uint16((uint32(p.db[15])*uint32(p.a[15])+uint32(p.b[15])*uint32(invDa15)+255)>>8)

	p.a[0], p.a[1], p.a[2], p.a[3] = uint16((uint32(p.da[0])*uint32(p.a[0])+uint32(p.a[0])*uint32(invDa0)+255)>>8), uint16((uint32(p.da[1])*uint32(p.a[1])+uint32(p.a[1])*uint32(invDa1)+255)>>8), uint16((uint32(p.da[2])*uint32(p.a[2])+uint32(p.a[2])*uint32(invDa2)+255)>>8), uint16((uint32(p.da[3])*uint32(p.a[3])+uint32(p.a[3])*uint32(invDa3)+255)>>8)
	p.a[4], p.a[5], p.a[6], p.a[7] = uint16((uint32(p.da[4])*uint32(p.a[4])+uint32(p.a[4])*uint32(invDa4)+255)>>8), uint16((uint32(p.da[5])*uint32(p.a[5])+uint32(p.a[5])*uint32(invDa5)+255)>>8), uint16((uint32(p.da[6])*uint32(p.a[6])+uint32(p.a[6])*uint32(invDa6)+255)>>8), uint16((uint32(p.da[7])*uint32(p.a[7])+uint32(p.a[7])*uint32(invDa7)+255)>>8)
	p.a[8], p.a[9], p.a[10], p.a[11] = uint16((uint32(p.da[8])*uint32(p.a[8])+uint32(p.a[8])*uint32(invDa8)+255)>>8), uint16((uint32(p.da[9])*uint32(p.a[9])+uint32(p.a[9])*uint32(invDa9)+255)>>8), uint16((uint32(p.da[10])*uint32(p.a[10])+uint32(p.a[10])*uint32(invDa10)+255)>>8), uint16((uint32(p.da[11])*uint32(p.a[11])+uint32(p.a[11])*uint32(invDa11)+255)>>8)
	p.a[12], p.a[13], p.a[14], p.a[15] = uint16((uint32(p.da[12])*uint32(p.a[12])+uint32(p.a[12])*uint32(invDa12)+255)>>8), uint16((uint32(p.da[13])*uint32(p.a[13])+uint32(p.a[13])*uint32(invDa13)+255)>>8), uint16((uint32(p.da[14])*uint32(p.a[14])+uint32(p.a[14])*uint32(invDa14)+255)>>8), uint16((uint32(p.da[15])*uint32(p.a[15])+uint32(p.a[15])*uint32(invDa15)+255)>>8)
}

//go:fix inline
func (p *LowPipeline) DestinationIn() {
	p.r[0], p.r[1], p.r[2], p.r[3] = uint16((uint32(p.dr[0])*uint32(p.a[0])+255)>>8), uint16((uint32(p.dr[1])*uint32(p.a[1])+255)>>8), uint16((uint32(p.dr[2])*uint32(p.a[2])+255)>>8), uint16((uint32(p.dr[3])*uint32(p.a[3])+255)>>8)
	p.r[4], p.r[5], p.r[6], p.r[7] = uint16((uint32(p.dr[4])*uint32(p.a[4])+255)>>8), uint16((uint32(p.dr[5])*uint32(p.a[5])+255)>>8), uint16((uint32(p.dr[6])*uint32(p.a[6])+255)>>8), uint16((uint32(p.dr[7])*uint32(p.a[7])+255)>>8)
	p.r[8], p.r[9], p.r[10], p.r[11] = uint16((uint32(p.dr[8])*uint32(p.a[8])+255)>>8), uint16((uint32(p.dr[9])*uint32(p.a[9])+255)>>8), uint16((uint32(p.dr[10])*uint32(p.a[10])+255)>>8), uint16((uint32(p.dr[11])*uint32(p.a[11])+255)>>8)
	p.r[12], p.r[13], p.r[14], p.r[15] = uint16((uint32(p.dr[12])*uint32(p.a[12])+255)>>8), uint16((uint32(p.dr[13])*uint32(p.a[13])+255)>>8), uint16((uint32(p.dr[14])*uint32(p.a[14])+255)>>8), uint16((uint32(p.dr[15])*uint32(p.a[15])+255)>>8)

	p.g[0], p.g[1], p.g[2], p.g[3] = uint16((uint32(p.dg[0])*uint32(p.a[0])+255)>>8), uint16((uint32(p.dg[1])*uint32(p.a[1])+255)>>8), uint16((uint32(p.dg[2])*uint32(p.a[2])+255)>>8), uint16((uint32(p.dg[3])*uint32(p.a[3])+255)>>8)
	p.g[4], p.g[5], p.g[6], p.g[7] = uint16((uint32(p.dg[4])*uint32(p.a[4])+255)>>8), uint16((uint32(p.dg[5])*uint32(p.a[5])+255)>>8), uint16((uint32(p.dg[6])*uint32(p.a[6])+255)>>8), uint16((uint32(p.dg[7])*uint32(p.a[7])+255)>>8)
	p.g[8], p.g[9], p.g[10], p.g[11] = uint16((uint32(p.dg[8])*uint32(p.a[8])+255)>>8), uint16((uint32(p.dg[9])*uint32(p.a[9])+255)>>8), uint16((uint32(p.dg[10])*uint32(p.a[10])+255)>>8), uint16((uint32(p.dg[11])*uint32(p.a[11])+255)>>8)
	p.g[12], p.g[13], p.g[14], p.g[15] = uint16((uint32(p.dg[12])*uint32(p.a[12])+255)>>8), uint16((uint32(p.dg[13])*uint32(p.a[13])+255)>>8), uint16((uint32(p.dg[14])*uint32(p.a[14])+255)>>8), uint16((uint32(p.dg[15])*uint32(p.a[15])+255)>>8)

	p.b[0], p.b[1], p.b[2], p.b[3] = uint16((uint32(p.db[0])*uint32(p.a[0])+255)>>8), uint16((uint32(p.db[1])*uint32(p.a[1])+255)>>8), uint16((uint32(p.db[2])*uint32(p.a[2])+255)>>8), uint16((uint32(p.db[3])*uint32(p.a[3])+255)>>8)
	p.b[4], p.b[5], p.b[6], p.b[7] = uint16((uint32(p.db[4])*uint32(p.a[4])+255)>>8), uint16((uint32(p.db[5])*uint32(p.a[5])+255)>>8), uint16((uint32(p.db[6])*uint32(p.a[6])+255)>>8), uint16((uint32(p.db[7])*uint32(p.a[7])+255)>>8)
	p.b[8], p.b[9], p.b[10], p.b[11] = uint16((uint32(p.db[8])*uint32(p.a[8])+255)>>8), uint16((uint32(p.db[9])*uint32(p.a[9])+255)>>8), uint16((uint32(p.db[10])*uint32(p.a[10])+255)>>8), uint16((uint32(p.db[11])*uint32(p.a[11])+255)>>8)
	p.b[12], p.b[13], p.b[14], p.b[15] = uint16((uint32(p.db[12])*uint32(p.a[12])+255)>>8), uint16((uint32(p.db[13])*uint32(p.a[13])+255)>>8), uint16((uint32(p.db[14])*uint32(p.a[14])+255)>>8), uint16((uint32(p.db[15])*uint32(p.a[15])+255)>>8)

	p.a[0], p.a[1], p.a[2], p.a[3] = uint16((uint32(p.da[0])*uint32(p.a[0])+255)>>8), uint16((uint32(p.da[1])*uint32(p.a[1])+255)>>8), uint16((uint32(p.da[2])*uint32(p.a[2])+255)>>8), uint16((uint32(p.da[3])*uint32(p.a[3])+255)>>8)
	p.a[4], p.a[5], p.a[6], p.a[7] = uint16((uint32(p.da[4])*uint32(p.a[4])+255)>>8), uint16((uint32(p.da[5])*uint32(p.a[5])+255)>>8), uint16((uint32(p.da[6])*uint32(p.a[6])+255)>>8), uint16((uint32(p.da[7])*uint32(p.a[7])+255)>>8)
	p.a[8], p.a[9], p.a[10], p.a[11] = uint16((uint32(p.da[8])*uint32(p.a[8])+255)>>8), uint16((uint32(p.da[9])*uint32(p.a[9])+255)>>8), uint16((uint32(p.da[10])*uint32(p.a[10])+255)>>8), uint16((uint32(p.da[11])*uint32(p.a[11])+255)>>8)
	p.a[12], p.a[13], p.a[14], p.a[15] = uint16((uint32(p.da[12])*uint32(p.a[12])+255)>>8), uint16((uint32(p.da[13])*uint32(p.a[13])+255)>>8), uint16((uint32(p.da[14])*uint32(p.a[14])+255)>>8), uint16((uint32(p.da[15])*uint32(p.a[15])+255)>>8)
}

//go:fix inline
func (p *LowPipeline) DestinationOut() {
	invSa0, invSa1, invSa2, invSa3 := 255-p.a[0], 255-p.a[1], 255-p.a[2], 255-p.a[3]
	invSa4, invSa5, invSa6, invSa7 := 255-p.a[4], 255-p.a[5], 255-p.a[6], 255-p.a[7]
	invSa8, invSa9, invSa10, invSa11 := 255-p.a[8], 255-p.a[9], 255-p.a[10], 255-p.a[11]
	invSa12, invSa13, invSa14, invSa15 := 255-p.a[12], 255-p.a[13], 255-p.a[14], 255-p.a[15]

	p.r[0], p.r[1], p.r[2], p.r[3] = uint16((uint32(p.dr[0])*uint32(invSa0)+255)>>8), uint16((uint32(p.dr[1])*uint32(invSa1)+255)>>8), uint16((uint32(p.dr[2])*uint32(invSa2)+255)>>8), uint16((uint32(p.dr[3])*uint32(invSa3)+255)>>8)
	p.r[4], p.r[5], p.r[6], p.r[7] = uint16((uint32(p.dr[4])*uint32(invSa4)+255)>>8), uint16((uint32(p.dr[5])*uint32(invSa5)+255)>>8), uint16((uint32(p.dr[6])*uint32(invSa6)+255)>>8), uint16((uint32(p.dr[7])*uint32(invSa7)+255)>>8)
	p.r[8], p.r[9], p.r[10], p.r[11] = uint16((uint32(p.dr[8])*uint32(invSa8)+255)>>8), uint16((uint32(p.dr[9])*uint32(invSa9)+255)>>8), uint16((uint32(p.dr[10])*uint32(invSa10)+255)>>8), uint16((uint32(p.dr[11])*uint32(invSa11)+255)>>8)
	p.r[12], p.r[13], p.r[14], p.r[15] = uint16((uint32(p.dr[12])*uint32(invSa12)+255)>>8), uint16((uint32(p.dr[13])*uint32(invSa13)+255)>>8), uint16((uint32(p.dr[14])*uint32(invSa14)+255)>>8), uint16((uint32(p.dr[15])*uint32(invSa15)+255)>>8)

	p.g[0], p.g[1], p.g[2], p.g[3] = uint16((uint32(p.dg[0])*uint32(invSa0)+255)>>8), uint16((uint32(p.dg[1])*uint32(invSa1)+255)>>8), uint16((uint32(p.dg[2])*uint32(invSa2)+255)>>8), uint16((uint32(p.dg[3])*uint32(invSa3)+255)>>8)
	p.g[4], p.g[5], p.g[6], p.g[7] = uint16((uint32(p.dg[4])*uint32(invSa4)+255)>>8), uint16((uint32(p.dg[5])*uint32(invSa5)+255)>>8), uint16((uint32(p.dg[6])*uint32(invSa6)+255)>>8), uint16((uint32(p.dg[7])*uint32(invSa7)+255)>>8)
	p.g[8], p.g[9], p.g[10], p.g[11] = uint16((uint32(p.dg[8])*uint32(invSa8)+255)>>8), uint16((uint32(p.dg[9])*uint32(invSa9)+255)>>8), uint16((uint32(p.dg[10])*uint32(invSa10)+255)>>8), uint16((uint32(p.dg[11])*uint32(invSa11)+255)>>8)
	p.g[12], p.g[13], p.g[14], p.g[15] = uint16((uint32(p.dg[12])*uint32(invSa12)+255)>>8), uint16((uint32(p.dg[13])*uint32(invSa13)+255)>>8), uint16((uint32(p.dg[14])*uint32(invSa14)+255)>>8), uint16((uint32(p.dg[15])*uint32(invSa15)+255)>>8)

	p.b[0], p.b[1], p.b[2], p.b[3] = uint16((uint32(p.db[0])*uint32(invSa0)+255)>>8), uint16((uint32(p.db[1])*uint32(invSa1)+255)>>8), uint16((uint32(p.db[2])*uint32(invSa2)+255)>>8), uint16((uint32(p.db[3])*uint32(invSa3)+255)>>8)
	p.b[4], p.b[5], p.b[6], p.b[7] = uint16((uint32(p.db[4])*uint32(invSa4)+255)>>8), uint16((uint32(p.db[5])*uint32(invSa5)+255)>>8), uint16((uint32(p.db[6])*uint32(invSa6)+255)>>8), uint16((uint32(p.db[7])*uint32(invSa7)+255)>>8)
	p.b[8], p.b[9], p.b[10], p.b[11] = uint16((uint32(p.db[8])*uint32(invSa8)+255)>>8), uint16((uint32(p.db[9])*uint32(invSa9)+255)>>8), uint16((uint32(p.db[10])*uint32(invSa10)+255)>>8), uint16((uint32(p.db[11])*uint32(invSa11)+255)>>8)
	p.b[12], p.b[13], p.b[14], p.b[15] = uint16((uint32(p.db[12])*uint32(invSa12)+255)>>8), uint16((uint32(p.db[13])*uint32(invSa13)+255)>>8), uint16((uint32(p.db[14])*uint32(invSa14)+255)>>8), uint16((uint32(p.db[15])*uint32(invSa15)+255)>>8)

	p.a[0], p.a[1], p.a[2], p.a[3] = uint16((uint32(p.da[0])*uint32(invSa0)+255)>>8), uint16((uint32(p.da[1])*uint32(invSa1)+255)>>8), uint16((uint32(p.da[2])*uint32(invSa2)+255)>>8), uint16((uint32(p.da[3])*uint32(invSa3)+255)>>8)
	p.a[4], p.a[5], p.a[6], p.a[7] = uint16((uint32(p.da[4])*uint32(invSa4)+255)>>8), uint16((uint32(p.da[5])*uint32(invSa5)+255)>>8), uint16((uint32(p.da[6])*uint32(invSa6)+255)>>8), uint16((uint32(p.da[7])*uint32(invSa7)+255)>>8)
	p.a[8], p.a[9], p.a[10], p.a[11] = uint16((uint32(p.da[8])*uint32(invSa8)+255)>>8), uint16((uint32(p.da[9])*uint32(invSa9)+255)>>8), uint16((uint32(p.da[10])*uint32(invSa10)+255)>>8), uint16((uint32(p.da[11])*uint32(invSa11)+255)>>8)
	p.a[12], p.a[13], p.a[14], p.a[15] = uint16((uint32(p.da[12])*uint32(invSa12)+255)>>8), uint16((uint32(p.da[13])*uint32(invSa13)+255)>>8), uint16((uint32(p.da[14])*uint32(invSa14)+255)>>8), uint16((uint32(p.da[15])*uint32(invSa15)+255)>>8)
}

//go:fix inline
func (p *LowPipeline) DestinationOver() {
	invDa0, invDa1, invDa2, invDa3 := 255-p.da[0], 255-p.da[1], 255-p.da[2], 255-p.da[3]
	invDa4, invDa5, invDa6, invDa7 := 255-p.da[4], 255-p.da[5], 255-p.da[6], 255-p.da[7]
	invDa8, invDa9, invDa10, invDa11 := 255-p.da[8], 255-p.da[9], 255-p.da[10], 255-p.da[11]
	invDa12, invDa13, invDa14, invDa15 := 255-p.da[12], 255-p.da[13], 255-p.da[14], 255-p.da[15]

	p.r[0], p.r[1], p.r[2], p.r[3] = uint16(uint32(p.dr[0])+(uint32(p.r[0])*uint32(invDa0)+255)>>8), uint16(uint32(p.dr[1])+(uint32(p.r[1])*uint32(invDa1)+255)>>8), uint16(uint32(p.dr[2])+(uint32(p.r[2])*uint32(invDa2)+255)>>8), uint16(uint32(p.dr[3])+(uint32(p.r[3])*uint32(invDa3)+255)>>8)
	p.r[4], p.r[5], p.r[6], p.r[7] = uint16(uint32(p.dr[4])+(uint32(p.r[4])*uint32(invDa4)+255)>>8), uint16(uint32(p.dr[5])+(uint32(p.r[5])*uint32(invDa5)+255)>>8), uint16(uint32(p.dr[6])+(uint32(p.r[6])*uint32(invDa6)+255)>>8), uint16(uint32(p.dr[7])+(uint32(p.r[7])*uint32(invDa7)+255)>>8)
	p.r[8], p.r[9], p.r[10], p.r[11] = uint16(uint32(p.dr[8])+(uint32(p.r[8])*uint32(invDa8)+255)>>8), uint16(uint32(p.dr[9])+(uint32(p.r[9])*uint32(invDa9)+255)>>8), uint16(uint32(p.dr[10])+(uint32(p.r[10])*uint32(invDa10)+255)>>8), uint16(uint32(p.dr[11])+(uint32(p.r[11])*uint32(invDa11)+255)>>8)
	p.r[12], p.r[13], p.r[14], p.r[15] = uint16(uint32(p.dr[12])+(uint32(p.r[12])*uint32(invDa12)+255)>>8), uint16(uint32(p.dr[13])+(uint32(p.r[13])*uint32(invDa13)+255)>>8), uint16(uint32(p.dr[14])+(uint32(p.r[14])*uint32(invDa14)+255)>>8), uint16(uint32(p.dr[15])+(uint32(p.r[15])*uint32(invDa15)+255)>>8)

	p.g[0], p.g[1], p.g[2], p.g[3] = uint16(uint32(p.dg[0])+(uint32(p.g[0])*uint32(invDa0)+255)>>8), uint16(uint32(p.dg[1])+(uint32(p.g[1])*uint32(invDa1)+255)>>8), uint16(uint32(p.dg[2])+(uint32(p.g[2])*uint32(invDa2)+255)>>8), uint16(uint32(p.dg[3])+(uint32(p.g[3])*uint32(invDa3)+255)>>8)
	p.g[4], p.g[5], p.g[6], p.g[7] = uint16(uint32(p.dg[4])+(uint32(p.g[4])*uint32(invDa4)+255)>>8), uint16(uint32(p.dg[5])+(uint32(p.g[5])*uint32(invDa5)+255)>>8), uint16(uint32(p.dg[6])+(uint32(p.g[6])*uint32(invDa6)+255)>>8), uint16(uint32(p.dg[7])+(uint32(p.g[7])*uint32(invDa7)+255)>>8)
	p.g[8], p.g[9], p.g[10], p.g[11] = uint16(uint32(p.dg[8])+(uint32(p.g[8])*uint32(invDa8)+255)>>8), uint16(uint32(p.dg[9])+(uint32(p.g[9])*uint32(invDa9)+255)>>8), uint16(uint32(p.dg[10])+(uint32(p.g[10])*uint32(invDa10)+255)>>8), uint16(uint32(p.dg[11])+(uint32(p.g[11])*uint32(invDa11)+255)>>8)
	p.g[12], p.g[13], p.g[14], p.g[15] = uint16(uint32(p.dg[12])+(uint32(p.g[12])*uint32(invDa12)+255)>>8), uint16(uint32(p.dg[13])+(uint32(p.g[13])*uint32(invDa13)+255)>>8), uint16(uint32(p.dg[14])+(uint32(p.g[14])*uint32(invDa14)+255)>>8), uint16(uint32(p.dg[15])+(uint32(p.g[15])*uint32(invDa15)+255)>>8)

	p.b[0], p.b[1], p.b[2], p.b[3] = uint16(uint32(p.db[0])+(uint32(p.b[0])*uint32(invDa0)+255)>>8), uint16(uint32(p.db[1])+(uint32(p.b[1])*uint32(invDa1)+255)>>8), uint16(uint32(p.db[2])+(uint32(p.b[2])*uint32(invDa2)+255)>>8), uint16(uint32(p.db[3])+(uint32(p.b[3])*uint32(invDa3)+255)>>8)
	p.b[4], p.b[5], p.b[6], p.b[7] = uint16(uint32(p.db[4])+(uint32(p.b[4])*uint32(invDa4)+255)>>8), uint16(uint32(p.db[5])+(uint32(p.b[5])*uint32(invDa5)+255)>>8), uint16(uint32(p.db[6])+(uint32(p.b[6])*uint32(invDa6)+255)>>8), uint16(uint32(p.db[7])+(uint32(p.b[7])*uint32(invDa7)+255)>>8)
	p.b[8], p.b[9], p.b[10], p.b[11] = uint16(uint32(p.db[8])+(uint32(p.b[8])*uint32(invDa8)+255)>>8), uint16(uint32(p.db[9])+(uint32(p.b[9])*uint32(invDa9)+255)>>8), uint16(uint32(p.db[10])+(uint32(p.b[10])*uint32(invDa10)+255)>>8), uint16(uint32(p.db[11])+(uint32(p.b[11])*uint32(invDa11)+255)>>8)
	p.b[12], p.b[13], p.b[14], p.b[15] = uint16(uint32(p.db[12])+(uint32(p.b[12])*uint32(invDa12)+255)>>8), uint16(uint32(p.db[13])+(uint32(p.b[13])*uint32(invDa13)+255)>>8), uint16(uint32(p.db[14])+(uint32(p.b[14])*uint32(invDa14)+255)>>8), uint16(uint32(p.db[15])+(uint32(p.b[15])*uint32(invDa15)+255)>>8)

	p.a[0], p.a[1], p.a[2], p.a[3] = uint16(uint32(p.da[0])+(uint32(p.a[0])*uint32(invDa0)+255)>>8), uint16(uint32(p.da[1])+(uint32(p.a[1])*uint32(invDa1)+255)>>8), uint16(uint32(p.da[2])+(uint32(p.a[2])*uint32(invDa2)+255)>>8), uint16(uint32(p.da[3])+(uint32(p.a[3])*uint32(invDa3)+255)>>8)
	p.a[4], p.a[5], p.a[6], p.a[7] = uint16(uint32(p.da[4])+(uint32(p.a[4])*uint32(invDa4)+255)>>8), uint16(uint32(p.da[5])+(uint32(p.a[5])*uint32(invDa5)+255)>>8), uint16(uint32(p.da[6])+(uint32(p.a[6])*uint32(invDa6)+255)>>8), uint16(uint32(p.da[7])+(uint32(p.a[7])*uint32(invDa7)+255)>>8)
	p.a[8], p.a[9], p.a[10], p.a[11] = uint16(uint32(p.da[8])+(uint32(p.a[8])*uint32(invDa8)+255)>>8), uint16(uint32(p.da[9])+(uint32(p.a[9])*uint32(invDa9)+255)>>8), uint16(uint32(p.da[10])+(uint32(p.a[10])*uint32(invDa10)+255)>>8), uint16(uint32(p.da[11])+(uint32(p.a[11])*uint32(invDa11)+255)>>8)
	p.a[12], p.a[13], p.a[14], p.a[15] = uint16(uint32(p.da[12])+(uint32(p.a[12])*uint32(invDa12)+255)>>8), uint16(uint32(p.da[13])+(uint32(p.a[13])*uint32(invDa13)+255)>>8), uint16(uint32(p.da[14])+(uint32(p.a[14])*uint32(invDa14)+255)>>8), uint16(uint32(p.da[15])+(uint32(p.a[15])*uint32(invDa15)+255)>>8)
}

//go:fix inline
func (p *LowPipeline) SourceAtop() {
	// Formula: div255(s * da + d * inv(sa))
	ch := func(s, d, da, sa uint16) uint16 {
		invSa := uint16(255 - sa)
		return uint16((uint32(s)*uint32(da) + uint32(d)*uint32(invSa) + 255) >> 8)
	}

	// Blend R channel (4 at a time)
	p.r[0], p.r[1], p.r[2], p.r[3] = ch(p.r[0], p.dr[0], p.da[0], p.a[0]), ch(p.r[1], p.dr[1], p.da[1], p.a[1]), ch(p.r[2], p.dr[2], p.da[2], p.a[2]), ch(p.r[3], p.dr[3], p.da[3], p.a[3])
	p.r[4], p.r[5], p.r[6], p.r[7] = ch(p.r[4], p.dr[4], p.da[4], p.a[4]), ch(p.r[5], p.dr[5], p.da[5], p.a[5]), ch(p.r[6], p.dr[6], p.da[6], p.a[6]), ch(p.r[7], p.dr[7], p.da[7], p.a[7])
	p.r[8], p.r[9], p.r[10], p.r[11] = ch(p.r[8], p.dr[8], p.da[8], p.a[8]), ch(p.r[9], p.dr[9], p.da[9], p.a[9]), ch(p.r[10], p.dr[10], p.da[10], p.a[10]), ch(p.r[11], p.dr[11], p.da[11], p.a[11])
	p.r[12], p.r[13], p.r[14], p.r[15] = ch(p.r[12], p.dr[12], p.da[12], p.a[12]), ch(p.r[13], p.dr[13], p.da[13], p.a[13]), ch(p.r[14], p.dr[14], p.da[14], p.a[14]), ch(p.r[15], p.dr[15], p.da[15], p.a[15])

	// Blend G channel (4 at a time)
	p.g[0], p.g[1], p.g[2], p.g[3] = ch(p.g[0], p.dg[0], p.da[0], p.a[0]), ch(p.g[1], p.dg[1], p.da[1], p.a[1]), ch(p.g[2], p.dg[2], p.da[2], p.a[2]), ch(p.g[3], p.dg[3], p.da[3], p.a[3])
	p.g[4], p.g[5], p.g[6], p.g[7] = ch(p.g[4], p.dg[4], p.da[4], p.a[4]), ch(p.g[5], p.dg[5], p.da[5], p.a[5]), ch(p.g[6], p.dg[6], p.da[6], p.a[6]), ch(p.g[7], p.dg[7], p.da[7], p.a[7])
	p.g[8], p.g[9], p.g[10], p.g[11] = ch(p.g[8], p.dg[8], p.da[8], p.a[8]), ch(p.g[9], p.dg[9], p.da[9], p.a[9]), ch(p.g[10], p.dg[10], p.da[10], p.a[10]), ch(p.g[11], p.dg[11], p.da[11], p.a[11])
	p.g[12], p.g[13], p.g[14], p.g[15] = ch(p.g[12], p.dg[12], p.da[12], p.a[12]), ch(p.g[13], p.dg[13], p.da[13], p.a[13]), ch(p.g[14], p.dg[14], p.da[14], p.a[14]), ch(p.g[15], p.dg[15], p.da[15], p.a[15])

	// Blend B channel (4 at a time)
	p.b[0], p.b[1], p.b[2], p.b[3] = ch(p.b[0], p.db[0], p.da[0], p.a[0]), ch(p.b[1], p.db[1], p.da[1], p.a[1]), ch(p.b[2], p.db[2], p.da[2], p.a[2]), ch(p.b[3], p.db[3], p.da[3], p.a[3])
	p.b[4], p.b[5], p.b[6], p.b[7] = ch(p.b[4], p.db[4], p.da[4], p.a[4]), ch(p.b[5], p.db[5], p.da[5], p.a[5]), ch(p.b[6], p.db[6], p.da[6], p.a[6]), ch(p.b[7], p.db[7], p.da[7], p.a[7])
	p.b[8], p.b[9], p.b[10], p.b[11] = ch(p.b[8], p.db[8], p.da[8], p.a[8]), ch(p.b[9], p.db[9], p.da[9], p.a[9]), ch(p.b[10], p.db[10], p.da[10], p.a[10]), ch(p.b[11], p.db[11], p.da[11], p.a[11])
	p.b[12], p.b[13], p.b[14], p.b[15] = ch(p.b[12], p.db[12], p.da[12], p.a[12]), ch(p.b[13], p.db[13], p.da[13], p.a[13]), ch(p.b[14], p.db[14], p.da[14], p.a[14]), ch(p.b[15], p.db[15], p.da[15], p.a[15])

	// Alpha channel (4 at a time)
	p.a[0], p.a[1], p.a[2], p.a[3] = ch(p.a[0], p.da[0], p.da[0], p.a[0]), ch(p.a[1], p.da[1], p.da[1], p.a[1]), ch(p.a[2], p.da[2], p.da[2], p.a[2]), ch(p.a[3], p.da[3], p.da[3], p.a[3])
	p.a[4], p.a[5], p.a[6], p.a[7] = ch(p.a[4], p.da[4], p.da[4], p.a[4]), ch(p.a[5], p.da[5], p.da[5], p.a[5]), ch(p.a[6], p.da[6], p.da[6], p.a[6]), ch(p.a[7], p.da[7], p.da[7], p.a[7])
	p.a[8], p.a[9], p.a[10], p.a[11] = ch(p.a[8], p.da[8], p.da[8], p.a[8]), ch(p.a[9], p.da[9], p.da[9], p.a[9]), ch(p.a[10], p.da[10], p.da[10], p.a[10]), ch(p.a[11], p.da[11], p.da[11], p.a[11])
	p.a[12], p.a[13], p.a[14], p.a[15] = ch(p.a[12], p.da[12], p.da[12], p.a[12]), ch(p.a[13], p.da[13], p.da[13], p.a[13]), ch(p.a[14], p.da[14], p.da[14], p.a[14]), ch(p.a[15], p.da[15], p.da[15], p.a[15])
}

//go:fix inline
func (p *LowPipeline) SourceIn() {
	// Formula: div255(s * da)
	ch := func(s, da uint16) uint16 {
		return uint16((uint32(s)*uint32(da) + 255) >> 8)
	}

	// Blend R channel (4 at a time)
	p.r[0], p.r[1], p.r[2], p.r[3] = ch(p.r[0], p.da[0]), ch(p.r[1], p.da[1]), ch(p.r[2], p.da[2]), ch(p.r[3], p.da[3])
	p.r[4], p.r[5], p.r[6], p.r[7] = ch(p.r[4], p.da[4]), ch(p.r[5], p.da[5]), ch(p.r[6], p.da[6]), ch(p.r[7], p.da[7])
	p.r[8], p.r[9], p.r[10], p.r[11] = ch(p.r[8], p.da[8]), ch(p.r[9], p.da[9]), ch(p.r[10], p.da[10]), ch(p.r[11], p.da[11])
	p.r[12], p.r[13], p.r[14], p.r[15] = ch(p.r[12], p.da[12]), ch(p.r[13], p.da[13]), ch(p.r[14], p.da[14]), ch(p.r[15], p.da[15])

	// Blend G channel (4 at a time)
	p.g[0], p.g[1], p.g[2], p.g[3] = ch(p.g[0], p.da[0]), ch(p.g[1], p.da[1]), ch(p.g[2], p.da[2]), ch(p.g[3], p.da[3])
	p.g[4], p.g[5], p.g[6], p.g[7] = ch(p.g[4], p.da[4]), ch(p.g[5], p.da[5]), ch(p.g[6], p.da[6]), ch(p.g[7], p.da[7])
	p.g[8], p.g[9], p.g[10], p.g[11] = ch(p.g[8], p.da[8]), ch(p.g[9], p.da[9]), ch(p.g[10], p.da[10]), ch(p.g[11], p.da[11])
	p.g[12], p.g[13], p.g[14], p.g[15] = ch(p.g[12], p.da[12]), ch(p.g[13], p.da[13]), ch(p.g[14], p.da[14]), ch(p.g[15], p.da[15])

	// Blend B channel (4 at a time)
	p.b[0], p.b[1], p.b[2], p.b[3] = ch(p.b[0], p.da[0]), ch(p.b[1], p.da[1]), ch(p.b[2], p.da[2]), ch(p.b[3], p.da[3])
	p.b[4], p.b[5], p.b[6], p.b[7] = ch(p.b[4], p.da[4]), ch(p.b[5], p.da[5]), ch(p.b[6], p.da[6]), ch(p.b[7], p.da[7])
	p.b[8], p.b[9], p.b[10], p.b[11] = ch(p.b[8], p.da[8]), ch(p.b[9], p.da[9]), ch(p.b[10], p.da[10]), ch(p.b[11], p.da[11])
	p.b[12], p.b[13], p.b[14], p.b[15] = ch(p.b[12], p.da[12]), ch(p.b[13], p.da[13]), ch(p.b[14], p.da[14]), ch(p.b[15], p.da[15])

	// Alpha channel (4 at a time)
	p.a[0], p.a[1], p.a[2], p.a[3] = ch(p.a[0], p.da[0]), ch(p.a[1], p.da[1]), ch(p.a[2], p.da[2]), ch(p.a[3], p.da[3])
	p.a[4], p.a[5], p.a[6], p.a[7] = ch(p.a[4], p.da[4]), ch(p.a[5], p.da[5]), ch(p.a[6], p.da[6]), ch(p.a[7], p.da[7])
	p.a[8], p.a[9], p.a[10], p.a[11] = ch(p.a[8], p.da[8]), ch(p.a[9], p.da[9]), ch(p.a[10], p.da[10]), ch(p.a[11], p.da[11])
	p.a[12], p.a[13], p.a[14], p.a[15] = ch(p.a[12], p.da[12]), ch(p.a[13], p.da[13]), ch(p.a[14], p.da[14]), ch(p.a[15], p.da[15])
}

//go:fix inline
func (p *LowPipeline) SourceOut() {
	// Formula: div255(s * inv(da))
	ch := func(s, invDa uint16) uint16 {
		return uint16((uint32(s)*uint32(invDa) + 255) >> 8)
	}

	// Precompute inverse alpha values (4 at a time)
	invDa0, invDa1, invDa2, invDa3 := uint16(255-p.da[0]), uint16(255-p.da[1]), uint16(255-p.da[2]), uint16(255-p.da[3])
	invDa4, invDa5, invDa6, invDa7 := uint16(255-p.da[4]), uint16(255-p.da[5]), uint16(255-p.da[6]), uint16(255-p.da[7])
	invDa8, invDa9, invDa10, invDa11 := uint16(255-p.da[8]), uint16(255-p.da[9]), uint16(255-p.da[10]), uint16(255-p.da[11])
	invDa12, invDa13, invDa14, invDa15 := uint16(255-p.da[12]), uint16(255-p.da[13]), uint16(255-p.da[14]), uint16(255-p.da[15])

	// Blend R channel (4 at a time)
	p.r[0], p.r[1], p.r[2], p.r[3] = ch(p.r[0], invDa0), ch(p.r[1], invDa1), ch(p.r[2], invDa2), ch(p.r[3], invDa3)
	p.r[4], p.r[5], p.r[6], p.r[7] = ch(p.r[4], invDa4), ch(p.r[5], invDa5), ch(p.r[6], invDa6), ch(p.r[7], invDa7)
	p.r[8], p.r[9], p.r[10], p.r[11] = ch(p.r[8], invDa8), ch(p.r[9], invDa9), ch(p.r[10], invDa10), ch(p.r[11], invDa11)
	p.r[12], p.r[13], p.r[14], p.r[15] = ch(p.r[12], invDa12), ch(p.r[13], invDa13), ch(p.r[14], invDa14), ch(p.r[15], invDa15)

	// Blend G channel (4 at a time)
	p.g[0], p.g[1], p.g[2], p.g[3] = ch(p.g[0], invDa0), ch(p.g[1], invDa1), ch(p.g[2], invDa2), ch(p.g[3], invDa3)
	p.g[4], p.g[5], p.g[6], p.g[7] = ch(p.g[4], invDa4), ch(p.g[5], invDa5), ch(p.g[6], invDa6), ch(p.g[7], invDa7)
	p.g[8], p.g[9], p.g[10], p.g[11] = ch(p.g[8], invDa8), ch(p.g[9], invDa9), ch(p.g[10], invDa10), ch(p.g[11], invDa11)
	p.g[12], p.g[13], p.g[14], p.g[15] = ch(p.g[12], invDa12), ch(p.g[13], invDa13), ch(p.g[14], invDa14), ch(p.g[15], invDa15)

	// Blend B channel (4 at a time)
	p.b[0], p.b[1], p.b[2], p.b[3] = ch(p.b[0], invDa0), ch(p.b[1], invDa1), ch(p.b[2], invDa2), ch(p.b[3], invDa3)
	p.b[4], p.b[5], p.b[6], p.b[7] = ch(p.b[4], invDa4), ch(p.b[5], invDa5), ch(p.b[6], invDa6), ch(p.b[7], invDa7)
	p.b[8], p.b[9], p.b[10], p.b[11] = ch(p.b[8], invDa8), ch(p.b[9], invDa9), ch(p.b[10], invDa10), ch(p.b[11], invDa11)
	p.b[12], p.b[13], p.b[14], p.b[15] = ch(p.b[12], invDa12), ch(p.b[13], invDa13), ch(p.b[14], invDa14), ch(p.b[15], invDa15)

	// Alpha channel (4 at a time)
	p.a[0], p.a[1], p.a[2], p.a[3] = ch(p.a[0], invDa0), ch(p.a[1], invDa1), ch(p.a[2], invDa2), ch(p.a[3], invDa3)
	p.a[4], p.a[5], p.a[6], p.a[7] = ch(p.a[4], invDa4), ch(p.a[5], invDa5), ch(p.a[6], invDa6), ch(p.a[7], invDa7)
	p.a[8], p.a[9], p.a[10], p.a[11] = ch(p.a[8], invDa8), ch(p.a[9], invDa9), ch(p.a[10], invDa10), ch(p.a[11], invDa11)
	p.a[12], p.a[13], p.a[14], p.a[15] = ch(p.a[12], invDa12), ch(p.a[13], invDa13), ch(p.a[14], invDa14), ch(p.a[15], invDa15)
}

//go:fix inline
func (p *LowPipeline) SourceOver() {
	// Formula: s + div255(d * inv(sa))
	ch := func(s, d, invSa uint16) uint16 {
		return uint16(uint32(s) + ((uint32(d)*uint32(invSa) + 255) >> 8))
	}

	// Precompute inverse alpha values (4 at a time)
	invSa0, invSa1, invSa2, invSa3 := uint16(255-p.a[0]), uint16(255-p.a[1]), uint16(255-p.a[2]), uint16(255-p.a[3])
	invSa4, invSa5, invSa6, invSa7 := uint16(255-p.a[4]), uint16(255-p.a[5]), uint16(255-p.a[6]), uint16(255-p.a[7])
	invSa8, invSa9, invSa10, invSa11 := uint16(255-p.a[8]), uint16(255-p.a[9]), uint16(255-p.a[10]), uint16(255-p.a[11])
	invSa12, invSa13, invSa14, invSa15 := uint16(255-p.a[12]), uint16(255-p.a[13]), uint16(255-p.a[14]), uint16(255-p.a[15])

	// Blend R channel (4 at a time)
	p.r[0], p.r[1], p.r[2], p.r[3] = ch(p.r[0], p.dr[0], invSa0), ch(p.r[1], p.dr[1], invSa1), ch(p.r[2], p.dr[2], invSa2), ch(p.r[3], p.dr[3], invSa3)
	p.r[4], p.r[5], p.r[6], p.r[7] = ch(p.r[4], p.dr[4], invSa4), ch(p.r[5], p.dr[5], invSa5), ch(p.r[6], p.dr[6], invSa6), ch(p.r[7], p.dr[7], invSa7)
	p.r[8], p.r[9], p.r[10], p.r[11] = ch(p.r[8], p.dr[8], invSa8), ch(p.r[9], p.dr[9], invSa9), ch(p.r[10], p.dr[10], invSa10), ch(p.r[11], p.dr[11], invSa11)
	p.r[12], p.r[13], p.r[14], p.r[15] = ch(p.r[12], p.dr[12], invSa12), ch(p.r[13], p.dr[13], invSa13), ch(p.r[14], p.dr[14], invSa14), ch(p.r[15], p.dr[15], invSa15)

	// Blend G channel (4 at a time)
	p.g[0], p.g[1], p.g[2], p.g[3] = ch(p.g[0], p.dg[0], invSa0), ch(p.g[1], p.dg[1], invSa1), ch(p.g[2], p.dg[2], invSa2), ch(p.g[3], p.dg[3], invSa3)
	p.g[4], p.g[5], p.g[6], p.g[7] = ch(p.g[4], p.dg[4], invSa4), ch(p.g[5], p.dg[5], invSa5), ch(p.g[6], p.dg[6], invSa6), ch(p.g[7], p.dg[7], invSa7)
	p.g[8], p.g[9], p.g[10], p.g[11] = ch(p.g[8], p.dg[8], invSa8), ch(p.g[9], p.dg[9], invSa9), ch(p.g[10], p.dg[10], invSa10), ch(p.g[11], p.dg[11], invSa11)
	p.g[12], p.g[13], p.g[14], p.g[15] = ch(p.g[12], p.dg[12], invSa12), ch(p.g[13], p.dg[13], invSa13), ch(p.g[14], p.dg[14], invSa14), ch(p.g[15], p.dg[15], invSa15)

	// Blend B channel (4 at a time)
	p.b[0], p.b[1], p.b[2], p.b[3] = ch(p.b[0], p.db[0], invSa0), ch(p.b[1], p.db[1], invSa1), ch(p.b[2], p.db[2], invSa2), ch(p.b[3], p.db[3], invSa3)
	p.b[4], p.b[5], p.b[6], p.b[7] = ch(p.b[4], p.db[4], invSa4), ch(p.b[5], p.db[5], invSa5), ch(p.b[6], p.db[6], invSa6), ch(p.b[7], p.db[7], invSa7)
	p.b[8], p.b[9], p.b[10], p.b[11] = ch(p.b[8], p.db[8], invSa8), ch(p.b[9], p.db[9], invSa9), ch(p.b[10], p.db[10], invSa10), ch(p.b[11], p.db[11], invSa11)
	p.b[12], p.b[13], p.b[14], p.b[15] = ch(p.b[12], p.db[12], invSa12), ch(p.b[13], p.db[13], invSa13), ch(p.b[14], p.db[14], invSa14), ch(p.b[15], p.db[15], invSa15)

	// Alpha channel (4 at a time)
	p.a[0], p.a[1], p.a[2], p.a[3] = ch(p.a[0], p.da[0], invSa0), ch(p.a[1], p.da[1], invSa1), ch(p.a[2], p.da[2], invSa2), ch(p.a[3], p.da[3], invSa3)
	p.a[4], p.a[5], p.a[6], p.a[7] = ch(p.a[4], p.da[4], invSa4), ch(p.a[5], p.da[5], invSa5), ch(p.a[6], p.da[6], invSa6), ch(p.a[7], p.da[7], invSa7)
	p.a[8], p.a[9], p.a[10], p.a[11] = ch(p.a[8], p.da[8], invSa8), ch(p.a[9], p.da[9], invSa9), ch(p.a[10], p.da[10], invSa10), ch(p.a[11], p.da[11], invSa11)
	p.a[12], p.a[13], p.a[14], p.a[15] = ch(p.a[12], p.da[12], invSa12), ch(p.a[13], p.da[13], invSa13), ch(p.a[14], p.da[14], invSa14), ch(p.a[15], p.da[15], invSa15)
}

//go:fix inline
func (p *LowPipeline) Clear() {
	p.r = [LOW_STAGE_WIDTH]uint16{}
	p.g = [LOW_STAGE_WIDTH]uint16{}
	p.b = [LOW_STAGE_WIDTH]uint16{}
	p.a = [LOW_STAGE_WIDTH]uint16{}
}

//go:fix inline
func (p *LowPipeline) Modulate() {
	// Formula: div255(s * d)
	ch := func(s, d uint16) uint16 {
		return uint16((uint32(s)*uint32(d) + 255) >> 8)
	}

	// Blend R channel (4 at a time)
	p.r[0], p.r[1], p.r[2], p.r[3] = ch(p.r[0], p.dr[0]), ch(p.r[1], p.dr[1]), ch(p.r[2], p.dr[2]), ch(p.r[3], p.dr[3])
	p.r[4], p.r[5], p.r[6], p.r[7] = ch(p.r[4], p.dr[4]), ch(p.r[5], p.dr[5]), ch(p.r[6], p.dr[6]), ch(p.r[7], p.dr[7])
	p.r[8], p.r[9], p.r[10], p.r[11] = ch(p.r[8], p.dr[8]), ch(p.r[9], p.dr[9]), ch(p.r[10], p.dr[10]), ch(p.r[11], p.dr[11])
	p.r[12], p.r[13], p.r[14], p.r[15] = ch(p.r[12], p.dr[12]), ch(p.r[13], p.dr[13]), ch(p.r[14], p.dr[14]), ch(p.r[15], p.dr[15])

	// Blend G channel (4 at a time)
	p.g[0], p.g[1], p.g[2], p.g[3] = ch(p.g[0], p.dg[0]), ch(p.g[1], p.dg[1]), ch(p.g[2], p.dg[2]), ch(p.g[3], p.dg[3])
	p.g[4], p.g[5], p.g[6], p.g[7] = ch(p.g[4], p.dg[4]), ch(p.g[5], p.dg[5]), ch(p.g[6], p.dg[6]), ch(p.g[7], p.dg[7])
	p.g[8], p.g[9], p.g[10], p.g[11] = ch(p.g[8], p.dg[8]), ch(p.g[9], p.dg[9]), ch(p.g[10], p.dg[10]), ch(p.g[11], p.dg[11])
	p.g[12], p.g[13], p.g[14], p.g[15] = ch(p.g[12], p.dg[12]), ch(p.g[13], p.dg[13]), ch(p.g[14], p.dg[14]), ch(p.g[15], p.dg[15])

	// Blend B channel (4 at a time)
	p.b[0], p.b[1], p.b[2], p.b[3] = ch(p.b[0], p.db[0]), ch(p.b[1], p.db[1]), ch(p.b[2], p.db[2]), ch(p.b[3], p.db[3])
	p.b[4], p.b[5], p.b[6], p.b[7] = ch(p.b[4], p.db[4]), ch(p.b[5], p.db[5]), ch(p.b[6], p.db[6]), ch(p.b[7], p.db[7])
	p.b[8], p.b[9], p.b[10], p.b[11] = ch(p.b[8], p.db[8]), ch(p.b[9], p.db[9]), ch(p.b[10], p.db[10]), ch(p.b[11], p.db[11])
	p.b[12], p.b[13], p.b[14], p.b[15] = ch(p.b[12], p.db[12]), ch(p.b[13], p.db[13]), ch(p.b[14], p.db[14]), ch(p.b[15], p.db[15])

	// Alpha channel (4 at a time)
	p.a[0], p.a[1], p.a[2], p.a[3] = ch(p.a[0], p.da[0]), ch(p.a[1], p.da[1]), ch(p.a[2], p.da[2]), ch(p.a[3], p.da[3])
	p.a[4], p.a[5], p.a[6], p.a[7] = ch(p.a[4], p.da[4]), ch(p.a[5], p.da[5]), ch(p.a[6], p.da[6]), ch(p.a[7], p.da[7])
	p.a[8], p.a[9], p.a[10], p.a[11] = ch(p.a[8], p.da[8]), ch(p.a[9], p.da[9]), ch(p.a[10], p.da[10]), ch(p.a[11], p.da[11])
	p.a[12], p.a[13], p.a[14], p.a[15] = ch(p.a[12], p.da[12]), ch(p.a[13], p.da[13]), ch(p.a[14], p.da[14]), ch(p.a[15], p.da[15])
}

//go:fix inline
func (p *LowPipeline) Multiply() {
	// Formula: div255(s * inv(da) + d * inv(sa) + s * d)
	ch := func(s, d, invSa, invDa uint16) uint16 {
		return uint16((uint32(s)*uint32(invDa) + uint32(d)*uint32(invSa) + uint32(s)*uint32(d) + 255) >> 8)
	}

	// Precompute inverse alpha values (4 at a time)
	invSa0, invSa1, invSa2, invSa3 := uint16(255-p.a[0]), uint16(255-p.a[1]), uint16(255-p.a[2]), uint16(255-p.a[3])
	invSa4, invSa5, invSa6, invSa7 := uint16(255-p.a[4]), uint16(255-p.a[5]), uint16(255-p.a[6]), uint16(255-p.a[7])
	invSa8, invSa9, invSa10, invSa11 := uint16(255-p.a[8]), uint16(255-p.a[9]), uint16(255-p.a[10]), uint16(255-p.a[11])
	invSa12, invSa13, invSa14, invSa15 := uint16(255-p.a[12]), uint16(255-p.a[13]), uint16(255-p.a[14]), uint16(255-p.a[15])

	invDa0, invDa1, invDa2, invDa3 := uint16(255-p.da[0]), uint16(255-p.da[1]), uint16(255-p.da[2]), uint16(255-p.da[3])
	invDa4, invDa5, invDa6, invDa7 := uint16(255-p.da[4]), uint16(255-p.da[5]), uint16(255-p.da[6]), uint16(255-p.da[7])
	invDa8, invDa9, invDa10, invDa11 := uint16(255-p.da[8]), uint16(255-p.da[9]), uint16(255-p.da[10]), uint16(255-p.da[11])
	invDa12, invDa13, invDa14, invDa15 := uint16(255-p.da[12]), uint16(255-p.da[13]), uint16(255-p.da[14]), uint16(255-p.da[15])

	// Blend R channel (4 at a time)
	p.r[0], p.r[1], p.r[2], p.r[3] = ch(p.r[0], p.dr[0], invSa0, invDa0), ch(p.r[1], p.dr[1], invSa1, invDa1), ch(p.r[2], p.dr[2], invSa2, invDa2), ch(p.r[3], p.dr[3], invSa3, invDa3)
	p.r[4], p.r[5], p.r[6], p.r[7] = ch(p.r[4], p.dr[4], invSa4, invDa4), ch(p.r[5], p.dr[5], invSa5, invDa5), ch(p.r[6], p.dr[6], invSa6, invDa6), ch(p.r[7], p.dr[7], invSa7, invDa7)
	p.r[8], p.r[9], p.r[10], p.r[11] = ch(p.r[8], p.dr[8], invSa8, invDa8), ch(p.r[9], p.dr[9], invSa9, invDa9), ch(p.r[10], p.dr[10], invSa10, invDa10), ch(p.r[11], p.dr[11], invSa11, invDa11)
	p.r[12], p.r[13], p.r[14], p.r[15] = ch(p.r[12], p.dr[12], invSa12, invDa12), ch(p.r[13], p.dr[13], invSa13, invDa13), ch(p.r[14], p.dr[14], invSa14, invDa14), ch(p.r[15], p.dr[15], invSa15, invDa15)

	// Blend G channel (4 at a time)
	p.g[0], p.g[1], p.g[2], p.g[3] = ch(p.g[0], p.dg[0], invSa0, invDa0), ch(p.g[1], p.dg[1], invSa1, invDa1), ch(p.g[2], p.dg[2], invSa2, invDa2), ch(p.g[3], p.dg[3], invSa3, invDa3)
	p.g[4], p.g[5], p.g[6], p.g[7] = ch(p.g[4], p.dg[4], invSa4, invDa4), ch(p.g[5], p.dg[5], invSa5, invDa5), ch(p.g[6], p.dg[6], invSa6, invDa6), ch(p.g[7], p.dg[7], invSa7, invDa7)
	p.g[8], p.g[9], p.g[10], p.g[11] = ch(p.g[8], p.dg[8], invSa8, invDa8), ch(p.g[9], p.dg[9], invSa9, invDa9), ch(p.g[10], p.dg[10], invSa10, invDa10), ch(p.g[11], p.dg[11], invSa11, invDa11)
	p.g[12], p.g[13], p.g[14], p.g[15] = ch(p.g[12], p.dg[12], invSa12, invDa12), ch(p.g[13], p.dg[13], invSa13, invDa13), ch(p.g[14], p.dg[14], invSa14, invDa14), ch(p.g[15], p.dg[15], invSa15, invDa15)

	// Blend B channel (4 at a time)
	p.b[0], p.b[1], p.b[2], p.b[3] = ch(p.b[0], p.db[0], invSa0, invDa0), ch(p.b[1], p.db[1], invSa1, invDa1), ch(p.b[2], p.db[2], invSa2, invDa2), ch(p.b[3], p.db[3], invSa3, invDa3)
	p.b[4], p.b[5], p.b[6], p.b[7] = ch(p.b[4], p.db[4], invSa4, invDa4), ch(p.b[5], p.db[5], invSa5, invDa5), ch(p.b[6], p.db[6], invSa6, invDa6), ch(p.b[7], p.db[7], invSa7, invDa7)
	p.b[8], p.b[9], p.b[10], p.b[11] = ch(p.b[8], p.db[8], invSa8, invDa8), ch(p.b[9], p.db[9], invSa9, invDa9), ch(p.b[10], p.db[10], invSa10, invDa10), ch(p.b[11], p.db[11], invSa11, invDa11)
	p.b[12], p.b[13], p.b[14], p.b[15] = ch(p.b[12], p.db[12], invSa12, invDa12), ch(p.b[13], p.db[13], invSa13, invDa13), ch(p.b[14], p.db[14], invSa14, invDa14), ch(p.b[15], p.db[15], invSa15, invDa15)

	// Alpha channel (4 at a time)
	p.a[0], p.a[1], p.a[2], p.a[3] = ch(p.a[0], p.da[0], invSa0, invDa0), ch(p.a[1], p.da[1], invSa1, invDa1), ch(p.a[2], p.da[2], invSa2, invDa2), ch(p.a[3], p.da[3], invSa3, invDa3)
	p.a[4], p.a[5], p.a[6], p.a[7] = ch(p.a[4], p.da[4], invSa4, invDa4), ch(p.a[5], p.da[5], invSa5, invDa5), ch(p.a[6], p.da[6], invSa6, invDa6), ch(p.a[7], p.da[7], invSa7, invDa7)
	p.a[8], p.a[9], p.a[10], p.a[11] = ch(p.a[8], p.da[8], invSa8, invDa8), ch(p.a[9], p.da[9], invSa9, invDa9), ch(p.a[10], p.da[10], invSa10, invDa10), ch(p.a[11], p.da[11], invSa11, invDa11)
	p.a[12], p.a[13], p.a[14], p.a[15] = ch(p.a[12], p.da[12], invSa12, invDa12), ch(p.a[13], p.da[13], invSa13, invDa13), ch(p.a[14], p.da[14], invSa14, invDa14), ch(p.a[15], p.da[15], invSa15, invDa15)
}

//go:fix inline
func (p *LowPipeline) Plus() {
	// Formula: (s + d).min(255)
	ch := func(s, d uint16) uint16 {
		sum := uint32(s) + uint32(d)
		if sum > 255 {
			return 255
		}
		return uint16(sum)
	}

	// Blend R channel (4 at a time)
	p.r[0], p.r[1], p.r[2], p.r[3] = ch(p.r[0], p.dr[0]), ch(p.r[1], p.dr[1]), ch(p.r[2], p.dr[2]), ch(p.r[3], p.dr[3])
	p.r[4], p.r[5], p.r[6], p.r[7] = ch(p.r[4], p.dr[4]), ch(p.r[5], p.dr[5]), ch(p.r[6], p.dr[6]), ch(p.r[7], p.dr[7])
	p.r[8], p.r[9], p.r[10], p.r[11] = ch(p.r[8], p.dr[8]), ch(p.r[9], p.dr[9]), ch(p.r[10], p.dr[10]), ch(p.r[11], p.dr[11])
	p.r[12], p.r[13], p.r[14], p.r[15] = ch(p.r[12], p.dr[12]), ch(p.r[13], p.dr[13]), ch(p.r[14], p.dr[14]), ch(p.r[15], p.dr[15])

	// Blend G channel (4 at a time)
	p.g[0], p.g[1], p.g[2], p.g[3] = ch(p.g[0], p.dg[0]), ch(p.g[1], p.dg[1]), ch(p.g[2], p.dg[2]), ch(p.g[3], p.dg[3])
	p.g[4], p.g[5], p.g[6], p.g[7] = ch(p.g[4], p.dg[4]), ch(p.g[5], p.dg[5]), ch(p.g[6], p.dg[6]), ch(p.g[7], p.dg[7])
	p.g[8], p.g[9], p.g[10], p.g[11] = ch(p.g[8], p.dg[8]), ch(p.g[9], p.dg[9]), ch(p.g[10], p.dg[10]), ch(p.g[11], p.dg[11])
	p.g[12], p.g[13], p.g[14], p.g[15] = ch(p.g[12], p.dg[12]), ch(p.g[13], p.dg[13]), ch(p.g[14], p.dg[14]), ch(p.g[15], p.dg[15])

	// Blend B channel (4 at a time)
	p.b[0], p.b[1], p.b[2], p.b[3] = ch(p.b[0], p.db[0]), ch(p.b[1], p.db[1]), ch(p.b[2], p.db[2]), ch(p.b[3], p.db[3])
	p.b[4], p.b[5], p.b[6], p.b[7] = ch(p.b[4], p.db[4]), ch(p.b[5], p.db[5]), ch(p.b[6], p.db[6]), ch(p.b[7], p.db[7])
	p.b[8], p.b[9], p.b[10], p.b[11] = ch(p.b[8], p.db[8]), ch(p.b[9], p.db[9]), ch(p.b[10], p.db[10]), ch(p.b[11], p.db[11])
	p.b[12], p.b[13], p.b[14], p.b[15] = ch(p.b[12], p.db[12]), ch(p.b[13], p.db[13]), ch(p.b[14], p.db[14]), ch(p.b[15], p.db[15])

	// Alpha channel (4 at a time)
	p.a[0], p.a[1], p.a[2], p.a[3] = ch(p.a[0], p.da[0]), ch(p.a[1], p.da[1]), ch(p.a[2], p.da[2]), ch(p.a[3], p.da[3])
	p.a[4], p.a[5], p.a[6], p.a[7] = ch(p.a[4], p.da[4]), ch(p.a[5], p.da[5]), ch(p.a[6], p.da[6]), ch(p.a[7], p.da[7])
	p.a[8], p.a[9], p.a[10], p.a[11] = ch(p.a[8], p.da[8]), ch(p.a[9], p.da[9]), ch(p.a[10], p.da[10]), ch(p.a[11], p.da[11])
	p.a[12], p.a[13], p.a[14], p.a[15] = ch(p.a[12], p.da[12]), ch(p.a[13], p.da[13]), ch(p.a[14], p.da[14]), ch(p.a[15], p.da[15])
}

//go:fix inline
func (p *LowPipeline) Screen() {
	// Formula: s + d - div255(s * d)
	ch := func(s, d uint16) uint16 {
		return uint16(uint32(s) + uint32(d) - ((uint32(s)*uint32(d) + 255) >> 8))
	}

	// Blend R channel (4 at a time)
	p.r[0], p.r[1], p.r[2], p.r[3] = ch(p.r[0], p.dr[0]), ch(p.r[1], p.dr[1]), ch(p.r[2], p.dr[2]), ch(p.r[3], p.dr[3])
	p.r[4], p.r[5], p.r[6], p.r[7] = ch(p.r[4], p.dr[4]), ch(p.r[5], p.dr[5]), ch(p.r[6], p.dr[6]), ch(p.r[7], p.dr[7])
	p.r[8], p.r[9], p.r[10], p.r[11] = ch(p.r[8], p.dr[8]), ch(p.r[9], p.dr[9]), ch(p.r[10], p.dr[10]), ch(p.r[11], p.dr[11])
	p.r[12], p.r[13], p.r[14], p.r[15] = ch(p.r[12], p.dr[12]), ch(p.r[13], p.dr[13]), ch(p.r[14], p.dr[14]), ch(p.r[15], p.dr[15])

	// Blend G channel (4 at a time)
	p.g[0], p.g[1], p.g[2], p.g[3] = ch(p.g[0], p.dg[0]), ch(p.g[1], p.dg[1]), ch(p.g[2], p.dg[2]), ch(p.g[3], p.dg[3])
	p.g[4], p.g[5], p.g[6], p.g[7] = ch(p.g[4], p.dg[4]), ch(p.g[5], p.dg[5]), ch(p.g[6], p.dg[6]), ch(p.g[7], p.dg[7])
	p.g[8], p.g[9], p.g[10], p.g[11] = ch(p.g[8], p.dg[8]), ch(p.g[9], p.dg[9]), ch(p.g[10], p.dg[10]), ch(p.g[11], p.dg[11])
	p.g[12], p.g[13], p.g[14], p.g[15] = ch(p.g[12], p.dg[12]), ch(p.g[13], p.dg[13]), ch(p.g[14], p.dg[14]), ch(p.g[15], p.dg[15])

	// Blend B channel (4 at a time)
	p.b[0], p.b[1], p.b[2], p.b[3] = ch(p.b[0], p.db[0]), ch(p.b[1], p.db[1]), ch(p.b[2], p.db[2]), ch(p.b[3], p.db[3])
	p.b[4], p.b[5], p.b[6], p.b[7] = ch(p.b[4], p.db[4]), ch(p.b[5], p.db[5]), ch(p.b[6], p.db[6]), ch(p.b[7], p.db[7])
	p.b[8], p.b[9], p.b[10], p.b[11] = ch(p.b[8], p.db[8]), ch(p.b[9], p.db[9]), ch(p.b[10], p.db[10]), ch(p.b[11], p.db[11])
	p.b[12], p.b[13], p.b[14], p.b[15] = ch(p.b[12], p.db[12]), ch(p.b[13], p.db[13]), ch(p.b[14], p.db[14]), ch(p.b[15], p.db[15])

	// Alpha channel (4 at a time)
	p.a[0], p.a[1], p.a[2], p.a[3] = ch(p.a[0], p.da[0]), ch(p.a[1], p.da[1]), ch(p.a[2], p.da[2]), ch(p.a[3], p.da[3])
	p.a[4], p.a[5], p.a[6], p.a[7] = ch(p.a[4], p.da[4]), ch(p.a[5], p.da[5]), ch(p.a[6], p.da[6]), ch(p.a[7], p.da[7])
	p.a[8], p.a[9], p.a[10], p.a[11] = ch(p.a[8], p.da[8]), ch(p.a[9], p.da[9]), ch(p.a[10], p.da[10]), ch(p.a[11], p.da[11])
	p.a[12], p.a[13], p.a[14], p.a[15] = ch(p.a[12], p.da[12]), ch(p.a[13], p.da[13]), ch(p.a[14], p.da[14]), ch(p.a[15], p.da[15])
}

//go:fix inline
func (p *LowPipeline) Xor() {
	// Formula: div255(s * inv(da) + d * inv(sa))
	ch := func(s, d, invSa, invDa uint16) uint16 {
		return uint16((uint32(s)*uint32(invDa) + uint32(d)*uint32(invSa) + 255) >> 8)
	}

	// Precompute inverse alpha values (4 at a time)
	invSa0, invSa1, invSa2, invSa3 := uint16(255-p.a[0]), uint16(255-p.a[1]), uint16(255-p.a[2]), uint16(255-p.a[3])
	invSa4, invSa5, invSa6, invSa7 := uint16(255-p.a[4]), uint16(255-p.a[5]), uint16(255-p.a[6]), uint16(255-p.a[7])
	invSa8, invSa9, invSa10, invSa11 := uint16(255-p.a[8]), uint16(255-p.a[9]), uint16(255-p.a[10]), uint16(255-p.a[11])
	invSa12, invSa13, invSa14, invSa15 := uint16(255-p.a[12]), uint16(255-p.a[13]), uint16(255-p.a[14]), uint16(255-p.a[15])

	invDa0, invDa1, invDa2, invDa3 := uint16(255-p.da[0]), uint16(255-p.da[1]), uint16(255-p.da[2]), uint16(255-p.da[3])
	invDa4, invDa5, invDa6, invDa7 := uint16(255-p.da[4]), uint16(255-p.da[5]), uint16(255-p.da[6]), uint16(255-p.da[7])
	invDa8, invDa9, invDa10, invDa11 := uint16(255-p.da[8]), uint16(255-p.da[9]), uint16(255-p.da[10]), uint16(255-p.da[11])
	invDa12, invDa13, invDa14, invDa15 := uint16(255-p.da[12]), uint16(255-p.da[13]), uint16(255-p.da[14]), uint16(255-p.da[15])

	// Blend R channel (4 at a time)
	p.r[0], p.r[1], p.r[2], p.r[3] = ch(p.r[0], p.dr[0], invSa0, invDa0), ch(p.r[1], p.dr[1], invSa1, invDa1), ch(p.r[2], p.dr[2], invSa2, invDa2), ch(p.r[3], p.dr[3], invSa3, invDa3)
	p.r[4], p.r[5], p.r[6], p.r[7] = ch(p.r[4], p.dr[4], invSa4, invDa4), ch(p.r[5], p.dr[5], invSa5, invDa5), ch(p.r[6], p.dr[6], invSa6, invDa6), ch(p.r[7], p.dr[7], invSa7, invDa7)
	p.r[8], p.r[9], p.r[10], p.r[11] = ch(p.r[8], p.dr[8], invSa8, invDa8), ch(p.r[9], p.dr[9], invSa9, invDa9), ch(p.r[10], p.dr[10], invSa10, invDa10), ch(p.r[11], p.dr[11], invSa11, invDa11)
	p.r[12], p.r[13], p.r[14], p.r[15] = ch(p.r[12], p.dr[12], invSa12, invDa12), ch(p.r[13], p.dr[13], invSa13, invDa13), ch(p.r[14], p.dr[14], invSa14, invDa14), ch(p.r[15], p.dr[15], invSa15, invDa15)

	// Blend G channel (4 at a time)
	p.g[0], p.g[1], p.g[2], p.g[3] = ch(p.g[0], p.dg[0], invSa0, invDa0), ch(p.g[1], p.dg[1], invSa1, invDa1), ch(p.g[2], p.dg[2], invSa2, invDa2), ch(p.g[3], p.dg[3], invSa3, invDa3)
	p.g[4], p.g[5], p.g[6], p.g[7] = ch(p.g[4], p.dg[4], invSa4, invDa4), ch(p.g[5], p.dg[5], invSa5, invDa5), ch(p.g[6], p.dg[6], invSa6, invDa6), ch(p.g[7], p.dg[7], invSa7, invDa7)
	p.g[8], p.g[9], p.g[10], p.g[11] = ch(p.g[8], p.dg[8], invSa8, invDa8), ch(p.g[9], p.dg[9], invSa9, invDa9), ch(p.g[10], p.dg[10], invSa10, invDa10), ch(p.g[11], p.dg[11], invSa11, invDa11)
	p.g[12], p.g[13], p.g[14], p.g[15] = ch(p.g[12], p.dg[12], invSa12, invDa12), ch(p.g[13], p.dg[13], invSa13, invDa13), ch(p.g[14], p.dg[14], invSa14, invDa14), ch(p.g[15], p.dg[15], invSa15, invDa15)

	// Blend B channel (4 at a time)
	p.b[0], p.b[1], p.b[2], p.b[3] = ch(p.b[0], p.db[0], invSa0, invDa0), ch(p.b[1], p.db[1], invSa1, invDa1), ch(p.b[2], p.db[2], invSa2, invDa2), ch(p.b[3], p.db[3], invSa3, invDa3)
	p.b[4], p.b[5], p.b[6], p.b[7] = ch(p.b[4], p.db[4], invSa4, invDa4), ch(p.b[5], p.db[5], invSa5, invDa5), ch(p.b[6], p.db[6], invSa6, invDa6), ch(p.b[7], p.db[7], invSa7, invDa7)
	p.b[8], p.b[9], p.b[10], p.b[11] = ch(p.b[8], p.db[8], invSa8, invDa8), ch(p.b[9], p.db[9], invSa9, invDa9), ch(p.b[10], p.db[10], invSa10, invDa10), ch(p.b[11], p.db[11], invSa11, invDa11)
	p.b[12], p.b[13], p.b[14], p.b[15] = ch(p.b[12], p.db[12], invSa12, invDa12), ch(p.b[13], p.db[13], invSa13, invDa13), ch(p.b[14], p.db[14], invSa14, invDa14), ch(p.b[15], p.db[15], invSa15, invDa15)

	// Alpha channel (4 at a time)
	p.a[0], p.a[1], p.a[2], p.a[3] = ch(p.a[0], p.da[0], invSa0, invDa0), ch(p.a[1], p.da[1], invSa1, invDa1), ch(p.a[2], p.da[2], invSa2, invDa2), ch(p.a[3], p.da[3], invSa3, invDa3)
	p.a[4], p.a[5], p.a[6], p.a[7] = ch(p.a[4], p.da[4], invSa4, invDa4), ch(p.a[5], p.da[5], invSa5, invDa5), ch(p.a[6], p.da[6], invSa6, invDa6), ch(p.a[7], p.da[7], invSa7, invDa7)
	p.a[8], p.a[9], p.a[10], p.a[11] = ch(p.a[8], p.da[8], invSa8, invDa8), ch(p.a[9], p.da[9], invSa9, invDa9), ch(p.a[10], p.da[10], invSa10, invDa10), ch(p.a[11], p.da[11], invSa11, invDa11)
	p.a[12], p.a[13], p.a[14], p.a[15] = ch(p.a[12], p.da[12], invSa12, invDa12), ch(p.a[13], p.da[13], invSa13, invDa13), ch(p.a[14], p.da[14], invSa14, invDa14), ch(p.a[15], p.da[15], invSa15, invDa15)
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
	// Formula: s + d - div255(max(s * da, d * sa))
	ch := func(s, d, sa, da uint16) uint16 {
		prod1 := uint32(s) * uint32(da)
		prod2 := uint32(d) * uint32(sa)
		maxProd := u16max(uint16(prod2), uint16(prod1))
		return uint16(uint32(s) + uint32(d) - ((uint32(maxProd) + 255) >> 8))
	}

	// Blend R channel (4 at a time)
	p.r[0], p.r[1], p.r[2], p.r[3] = ch(p.r[0], p.dr[0], p.a[0], p.da[0]), ch(p.r[1], p.dr[1], p.a[1], p.da[1]), ch(p.r[2], p.dr[2], p.a[2], p.da[2]), ch(p.r[3], p.dr[3], p.a[3], p.da[3])
	p.r[4], p.r[5], p.r[6], p.r[7] = ch(p.r[4], p.dr[4], p.a[4], p.da[4]), ch(p.r[5], p.dr[5], p.a[5], p.da[5]), ch(p.r[6], p.dr[6], p.a[6], p.da[6]), ch(p.r[7], p.dr[7], p.a[7], p.da[7])
	p.r[8], p.r[9], p.r[10], p.r[11] = ch(p.r[8], p.dr[8], p.a[8], p.da[8]), ch(p.r[9], p.dr[9], p.a[9], p.da[9]), ch(p.r[10], p.dr[10], p.a[10], p.da[10]), ch(p.r[11], p.dr[11], p.a[11], p.da[11])
	p.r[12], p.r[13], p.r[14], p.r[15] = ch(p.r[12], p.dr[12], p.a[12], p.da[12]), ch(p.r[13], p.dr[13], p.a[13], p.da[13]), ch(p.r[14], p.dr[14], p.a[14], p.da[14]), ch(p.r[15], p.dr[15], p.a[15], p.da[15])

	// Blend G channel (4 at a time)
	p.g[0], p.g[1], p.g[2], p.g[3] = ch(p.g[0], p.dg[0], p.a[0], p.da[0]), ch(p.g[1], p.dg[1], p.a[1], p.da[1]), ch(p.g[2], p.dg[2], p.a[2], p.da[2]), ch(p.g[3], p.dg[3], p.a[3], p.da[3])
	p.g[4], p.g[5], p.g[6], p.g[7] = ch(p.g[4], p.dg[4], p.a[4], p.da[4]), ch(p.g[5], p.dg[5], p.a[5], p.da[5]), ch(p.g[6], p.dg[6], p.a[6], p.da[6]), ch(p.g[7], p.dg[7], p.a[7], p.da[7])
	p.g[8], p.g[9], p.g[10], p.g[11] = ch(p.g[8], p.dg[8], p.a[8], p.da[8]), ch(p.g[9], p.dg[9], p.a[9], p.da[9]), ch(p.g[10], p.dg[10], p.a[10], p.da[10]), ch(p.g[11], p.dg[11], p.a[11], p.da[11])
	p.g[12], p.g[13], p.g[14], p.g[15] = ch(p.g[12], p.dg[12], p.a[12], p.da[12]), ch(p.g[13], p.dg[13], p.a[13], p.da[13]), ch(p.g[14], p.dg[14], p.a[14], p.da[14]), ch(p.g[15], p.dg[15], p.a[15], p.da[15])

	// Blend B channel (4 at a time)
	p.b[0], p.b[1], p.b[2], p.b[3] = ch(p.b[0], p.db[0], p.a[0], p.da[0]), ch(p.b[1], p.db[1], p.a[1], p.da[1]), ch(p.b[2], p.db[2], p.a[2], p.da[2]), ch(p.b[3], p.db[3], p.a[3], p.da[3])
	p.b[4], p.b[5], p.b[6], p.b[7] = ch(p.b[4], p.db[4], p.a[4], p.da[4]), ch(p.b[5], p.db[5], p.a[5], p.da[5]), ch(p.b[6], p.db[6], p.a[6], p.da[6]), ch(p.b[7], p.db[7], p.a[7], p.da[7])
	p.b[8], p.b[9], p.b[10], p.b[11] = ch(p.b[8], p.db[8], p.a[8], p.da[8]), ch(p.b[9], p.db[9], p.a[9], p.da[9]), ch(p.b[10], p.db[10], p.a[10], p.da[10]), ch(p.b[11], p.db[11], p.a[11], p.da[11])
	p.b[12], p.b[13], p.b[14], p.b[15] = ch(p.b[12], p.db[12], p.a[12], p.da[12]), ch(p.b[13], p.db[13], p.a[13], p.da[13]), ch(p.b[14], p.db[14], p.a[14], p.da[14]), ch(p.b[15], p.db[15], p.a[15], p.da[15])

	// Alpha channel: source_over formula (4 at a time)
	invSa0, invSa1, invSa2, invSa3 := uint16(255-p.a[0]), uint16(255-p.a[1]), uint16(255-p.a[2]), uint16(255-p.a[3])
	invSa4, invSa5, invSa6, invSa7 := uint16(255-p.a[4]), uint16(255-p.a[5]), uint16(255-p.a[6]), uint16(255-p.a[7])
	invSa8, invSa9, invSa10, invSa11 := uint16(255-p.a[8]), uint16(255-p.a[9]), uint16(255-p.a[10]), uint16(255-p.a[11])
	invSa12, invSa13, invSa14, invSa15 := uint16(255-p.a[12]), uint16(255-p.a[13]), uint16(255-p.a[14]), uint16(255-p.a[15])

	p.a[0], p.a[1], p.a[2], p.a[3] = p.a[0]+uint16((uint32(p.da[0])*uint32(invSa0)+255)>>8), p.a[1]+uint16((uint32(p.da[1])*uint32(invSa1)+255)>>8), p.a[2]+uint16((uint32(p.da[2])*uint32(invSa2)+255)>>8), p.a[3]+uint16((uint32(p.da[3])*uint32(invSa3)+255)>>8)
	p.a[4], p.a[5], p.a[6], p.a[7] = p.a[4]+uint16((uint32(p.da[4])*uint32(invSa4)+255)>>8), p.a[5]+uint16((uint32(p.da[5])*uint32(invSa5)+255)>>8), p.a[6]+uint16((uint32(p.da[6])*uint32(invSa6)+255)>>8), p.a[7]+uint16((uint32(p.da[7])*uint32(invSa7)+255)>>8)
	p.a[8], p.a[9], p.a[10], p.a[11] = p.a[8]+uint16((uint32(p.da[8])*uint32(invSa8)+255)>>8), p.a[9]+uint16((uint32(p.da[9])*uint32(invSa9)+255)>>8), p.a[10]+uint16((uint32(p.da[10])*uint32(invSa10)+255)>>8), p.a[11]+uint16((uint32(p.da[11])*uint32(invSa11)+255)>>8)
	p.a[12], p.a[13], p.a[14], p.a[15] = p.a[12]+uint16((uint32(p.da[12])*uint32(invSa12)+255)>>8), p.a[13]+uint16((uint32(p.da[13])*uint32(invSa13)+255)>>8), p.a[14]+uint16((uint32(p.da[14])*uint32(invSa14)+255)>>8), p.a[15]+uint16((uint32(p.da[15])*uint32(invSa15)+255)>>8)
}

//go:fix inline
func (p *LowPipeline) Difference() {
	// Formula: s + d - 2 * div255(min(s * da, d * sa))
	ch := func(s, d, sa, da uint16) uint16 {
		prod1 := uint32(s) * uint32(da)
		prod2 := uint32(d) * uint32(sa)
		minProd := u16min(uint16(prod2), uint16(prod1))
		return uint16(uint32(s) + uint32(d) - 2*((uint32(minProd)+255)>>8))
	}

	// Blend R channel (4 at a time)
	p.r[0], p.r[1], p.r[2], p.r[3] = ch(p.r[0], p.dr[0], p.a[0], p.da[0]), ch(p.r[1], p.dr[1], p.a[1], p.da[1]), ch(p.r[2], p.dr[2], p.a[2], p.da[2]), ch(p.r[3], p.dr[3], p.a[3], p.da[3])
	p.r[4], p.r[5], p.r[6], p.r[7] = ch(p.r[4], p.dr[4], p.a[4], p.da[4]), ch(p.r[5], p.dr[5], p.a[5], p.da[5]), ch(p.r[6], p.dr[6], p.a[6], p.da[6]), ch(p.r[7], p.dr[7], p.a[7], p.da[7])
	p.r[8], p.r[9], p.r[10], p.r[11] = ch(p.r[8], p.dr[8], p.a[8], p.da[8]), ch(p.r[9], p.dr[9], p.a[9], p.da[9]), ch(p.r[10], p.dr[10], p.a[10], p.da[10]), ch(p.r[11], p.dr[11], p.a[11], p.da[11])
	p.r[12], p.r[13], p.r[14], p.r[15] = ch(p.r[12], p.dr[12], p.a[12], p.da[12]), ch(p.r[13], p.dr[13], p.a[13], p.da[13]), ch(p.r[14], p.dr[14], p.a[14], p.da[14]), ch(p.r[15], p.dr[15], p.a[15], p.da[15])

	// Blend G channel (4 at a time)
	p.g[0], p.g[1], p.g[2], p.g[3] = ch(p.g[0], p.dg[0], p.a[0], p.da[0]), ch(p.g[1], p.dg[1], p.a[1], p.da[1]), ch(p.g[2], p.dg[2], p.a[2], p.da[2]), ch(p.g[3], p.dg[3], p.a[3], p.da[3])
	p.g[4], p.g[5], p.g[6], p.g[7] = ch(p.g[4], p.dg[4], p.a[4], p.da[4]), ch(p.g[5], p.dg[5], p.a[5], p.da[5]), ch(p.g[6], p.dg[6], p.a[6], p.da[6]), ch(p.g[7], p.dg[7], p.a[7], p.da[7])
	p.g[8], p.g[9], p.g[10], p.g[11] = ch(p.g[8], p.dg[8], p.a[8], p.da[8]), ch(p.g[9], p.dg[9], p.a[9], p.da[9]), ch(p.g[10], p.dg[10], p.a[10], p.da[10]), ch(p.g[11], p.dg[11], p.a[11], p.da[11])
	p.g[12], p.g[13], p.g[14], p.g[15] = ch(p.g[12], p.dg[12], p.a[12], p.da[12]), ch(p.g[13], p.dg[13], p.a[13], p.da[13]), ch(p.g[14], p.dg[14], p.a[14], p.da[14]), ch(p.g[15], p.dg[15], p.a[15], p.da[15])

	// Blend B channel (4 at a time)
	p.b[0], p.b[1], p.b[2], p.b[3] = ch(p.b[0], p.db[0], p.a[0], p.da[0]), ch(p.b[1], p.db[1], p.a[1], p.da[1]), ch(p.b[2], p.db[2], p.a[2], p.da[2]), ch(p.b[3], p.db[3], p.a[3], p.da[3])
	p.b[4], p.b[5], p.b[6], p.b[7] = ch(p.b[4], p.db[4], p.a[4], p.da[4]), ch(p.b[5], p.db[5], p.a[5], p.da[5]), ch(p.b[6], p.db[6], p.a[6], p.da[6]), ch(p.b[7], p.db[7], p.a[7], p.da[7])
	p.b[8], p.b[9], p.b[10], p.b[11] = ch(p.b[8], p.db[8], p.a[8], p.da[8]), ch(p.b[9], p.db[9], p.a[9], p.da[9]), ch(p.b[10], p.db[10], p.a[10], p.da[10]), ch(p.b[11], p.db[11], p.a[11], p.da[11])
	p.b[12], p.b[13], p.b[14], p.b[15] = ch(p.b[12], p.db[12], p.a[12], p.da[12]), ch(p.b[13], p.db[13], p.a[13], p.da[13]), ch(p.b[14], p.db[14], p.a[14], p.da[14]), ch(p.b[15], p.db[15], p.a[15], p.da[15])

	// Alpha channel: source_over formula (4 at a time)
	invSa0, invSa1, invSa2, invSa3 := uint16(255-p.a[0]), uint16(255-p.a[1]), uint16(255-p.a[2]), uint16(255-p.a[3])
	invSa4, invSa5, invSa6, invSa7 := uint16(255-p.a[4]), uint16(255-p.a[5]), uint16(255-p.a[6]), uint16(255-p.a[7])
	invSa8, invSa9, invSa10, invSa11 := uint16(255-p.a[8]), uint16(255-p.a[9]), uint16(255-p.a[10]), uint16(255-p.a[11])
	invSa12, invSa13, invSa14, invSa15 := uint16(255-p.a[12]), uint16(255-p.a[13]), uint16(255-p.a[14]), uint16(255-p.a[15])

	p.a[0], p.a[1], p.a[2], p.a[3] = p.a[0]+uint16((uint32(p.da[0])*uint32(invSa0)+255)>>8), p.a[1]+uint16((uint32(p.da[1])*uint32(invSa1)+255)>>8), p.a[2]+uint16((uint32(p.da[2])*uint32(invSa2)+255)>>8), p.a[3]+uint16((uint32(p.da[3])*uint32(invSa3)+255)>>8)
	p.a[4], p.a[5], p.a[6], p.a[7] = p.a[4]+uint16((uint32(p.da[4])*uint32(invSa4)+255)>>8), p.a[5]+uint16((uint32(p.da[5])*uint32(invSa5)+255)>>8), p.a[6]+uint16((uint32(p.da[6])*uint32(invSa6)+255)>>8), p.a[7]+uint16((uint32(p.da[7])*uint32(invSa7)+255)>>8)
	p.a[8], p.a[9], p.a[10], p.a[11] = p.a[8]+uint16((uint32(p.da[8])*uint32(invSa8)+255)>>8), p.a[9]+uint16((uint32(p.da[9])*uint32(invSa9)+255)>>8), p.a[10]+uint16((uint32(p.da[10])*uint32(invSa10)+255)>>8), p.a[11]+uint16((uint32(p.da[11])*uint32(invSa11)+255)>>8)
	p.a[12], p.a[13], p.a[14], p.a[15] = p.a[12]+uint16((uint32(p.da[12])*uint32(invSa12)+255)>>8), p.a[13]+uint16((uint32(p.da[13])*uint32(invSa13)+255)>>8), p.a[14]+uint16((uint32(p.da[14])*uint32(invSa14)+255)>>8), p.a[15]+uint16((uint32(p.da[15])*uint32(invSa15)+255)>>8)
}

//go:fix inline
func (p *LowPipeline) Exclusion() {
	// Formula: s + d - 2 * div255(s * d)
	ch := func(s, d uint16) uint16 {
		prod := uint32(s) * uint32(d)
		return uint16(uint32(s) + uint32(d) - 2*((prod+255)>>8))
	}

	// Blend R channel (4 at a time)
	p.r[0], p.r[1], p.r[2], p.r[3] = ch(p.r[0], p.dr[0]), ch(p.r[1], p.dr[1]), ch(p.r[2], p.dr[2]), ch(p.r[3], p.dr[3])
	p.r[4], p.r[5], p.r[6], p.r[7] = ch(p.r[4], p.dr[4]), ch(p.r[5], p.dr[5]), ch(p.r[6], p.dr[6]), ch(p.r[7], p.dr[7])
	p.r[8], p.r[9], p.r[10], p.r[11] = ch(p.r[8], p.dr[8]), ch(p.r[9], p.dr[9]), ch(p.r[10], p.dr[10]), ch(p.r[11], p.dr[11])
	p.r[12], p.r[13], p.r[14], p.r[15] = ch(p.r[12], p.dr[12]), ch(p.r[13], p.dr[13]), ch(p.r[14], p.dr[14]), ch(p.r[15], p.dr[15])

	// Blend G channel (4 at a time)
	p.g[0], p.g[1], p.g[2], p.g[3] = ch(p.g[0], p.dg[0]), ch(p.g[1], p.dg[1]), ch(p.g[2], p.dg[2]), ch(p.g[3], p.dg[3])
	p.g[4], p.g[5], p.g[6], p.g[7] = ch(p.g[4], p.dg[4]), ch(p.g[5], p.dg[5]), ch(p.g[6], p.dg[6]), ch(p.g[7], p.dg[7])
	p.g[8], p.g[9], p.g[10], p.g[11] = ch(p.g[8], p.dg[8]), ch(p.g[9], p.dg[9]), ch(p.g[10], p.dg[10]), ch(p.g[11], p.dg[11])
	p.g[12], p.g[13], p.g[14], p.g[15] = ch(p.g[12], p.dg[12]), ch(p.g[13], p.dg[13]), ch(p.g[14], p.dg[14]), ch(p.g[15], p.dg[15])

	// Blend B channel (4 at a time)
	p.b[0], p.b[1], p.b[2], p.b[3] = ch(p.b[0], p.db[0]), ch(p.b[1], p.db[1]), ch(p.b[2], p.db[2]), ch(p.b[3], p.db[3])
	p.b[4], p.b[5], p.b[6], p.b[7] = ch(p.b[4], p.db[4]), ch(p.b[5], p.db[5]), ch(p.b[6], p.db[6]), ch(p.b[7], p.db[7])
	p.b[8], p.b[9], p.b[10], p.b[11] = ch(p.b[8], p.db[8]), ch(p.b[9], p.db[9]), ch(p.b[10], p.db[10]), ch(p.b[11], p.db[11])
	p.b[12], p.b[13], p.b[14], p.b[15] = ch(p.b[12], p.db[12]), ch(p.b[13], p.db[13]), ch(p.b[14], p.db[14]), ch(p.b[15], p.db[15])

	// Alpha channel: source_over formula (4 at a time)
	invSa0, invSa1, invSa2, invSa3 := uint16(255-p.a[0]), uint16(255-p.a[1]), uint16(255-p.a[2]), uint16(255-p.a[3])
	invSa4, invSa5, invSa6, invSa7 := uint16(255-p.a[4]), uint16(255-p.a[5]), uint16(255-p.a[6]), uint16(255-p.a[7])
	invSa8, invSa9, invSa10, invSa11 := uint16(255-p.a[8]), uint16(255-p.a[9]), uint16(255-p.a[10]), uint16(255-p.a[11])
	invSa12, invSa13, invSa14, invSa15 := uint16(255-p.a[12]), uint16(255-p.a[13]), uint16(255-p.a[14]), uint16(255-p.a[15])

	p.a[0], p.a[1], p.a[2], p.a[3] = p.a[0]+uint16((uint32(p.da[0])*uint32(invSa0)+255)>>8), p.a[1]+uint16((uint32(p.da[1])*uint32(invSa1)+255)>>8), p.a[2]+uint16((uint32(p.da[2])*uint32(invSa2)+255)>>8), p.a[3]+uint16((uint32(p.da[3])*uint32(invSa3)+255)>>8)
	p.a[4], p.a[5], p.a[6], p.a[7] = p.a[4]+uint16((uint32(p.da[4])*uint32(invSa4)+255)>>8), p.a[5]+uint16((uint32(p.da[5])*uint32(invSa5)+255)>>8), p.a[6]+uint16((uint32(p.da[6])*uint32(invSa6)+255)>>8), p.a[7]+uint16((uint32(p.da[7])*uint32(invSa7)+255)>>8)
	p.a[8], p.a[9], p.a[10], p.a[11] = p.a[8]+uint16((uint32(p.da[8])*uint32(invSa8)+255)>>8), p.a[9]+uint16((uint32(p.da[9])*uint32(invSa9)+255)>>8), p.a[10]+uint16((uint32(p.da[10])*uint32(invSa10)+255)>>8), p.a[11]+uint16((uint32(p.da[11])*uint32(invSa11)+255)>>8)
	p.a[12], p.a[13], p.a[14], p.a[15] = p.a[12]+uint16((uint32(p.da[12])*uint32(invSa12)+255)>>8), p.a[13]+uint16((uint32(p.da[13])*uint32(invSa13)+255)>>8), p.a[14]+uint16((uint32(p.da[14])*uint32(invSa14)+255)>>8), p.a[15]+uint16((uint32(p.da[15])*uint32(invSa15)+255)>>8)
}

//go:fix inline
func (p *LowPipeline) HardLight() {
	// Formula: div255(s * inv(da) + d * inv(sa) + (2*s <= sa ? 2*s*d : sa*da - 2*(sa-s)*(da-d)))
	ch := func(s, d, sa, da, invSa, invDa uint16) uint16 {
		s32, d32 := uint32(s), uint32(d)
		sa32, da32 := uint32(sa), uint32(da)
		invSa32, invDa32 := uint32(invSa), uint32(invDa)

		// Base term: s * inv(da) + d * inv(sa)
		term := s32*invDa32 + d32*invSa32

		// Branch term
		if 2*s32 <= sa32 {
			term += 2 * s32 * d32
		} else {
			term += sa32*da32 - 2*(sa32-s32)*(da32-d32)
		}

		return uint16((term + 255) >> 8)
	}

	// Precompute inverse alpha values (4 at a time)
	invSa0, invSa1, invSa2, invSa3 := uint16(255-p.a[0]), uint16(255-p.a[1]), uint16(255-p.a[2]), uint16(255-p.a[3])
	invSa4, invSa5, invSa6, invSa7 := uint16(255-p.a[4]), uint16(255-p.a[5]), uint16(255-p.a[6]), uint16(255-p.a[7])
	invSa8, invSa9, invSa10, invSa11 := uint16(255-p.a[8]), uint16(255-p.a[9]), uint16(255-p.a[10]), uint16(255-p.a[11])
	invSa12, invSa13, invSa14, invSa15 := uint16(255-p.a[12]), uint16(255-p.a[13]), uint16(255-p.a[14]), uint16(255-p.a[15])

	invDa0, invDa1, invDa2, invDa3 := uint16(255-p.da[0]), uint16(255-p.da[1]), uint16(255-p.da[2]), uint16(255-p.da[3])
	invDa4, invDa5, invDa6, invDa7 := uint16(255-p.da[4]), uint16(255-p.da[5]), uint16(255-p.da[6]), uint16(255-p.da[7])
	invDa8, invDa9, invDa10, invDa11 := uint16(255-p.da[8]), uint16(255-p.da[9]), uint16(255-p.da[10]), uint16(255-p.da[11])
	invDa12, invDa13, invDa14, invDa15 := uint16(255-p.da[12]), uint16(255-p.da[13]), uint16(255-p.da[14]), uint16(255-p.da[15])

	// Blend R channel (4 at a time)
	p.r[0], p.r[1], p.r[2], p.r[3] = ch(p.r[0], p.dr[0], p.a[0], p.da[0], invSa0, invDa0), ch(p.r[1], p.dr[1], p.a[1], p.da[1], invSa1, invDa1), ch(p.r[2], p.dr[2], p.a[2], p.da[2], invSa2, invDa2), ch(p.r[3], p.dr[3], p.a[3], p.da[3], invSa3, invDa3)
	p.r[4], p.r[5], p.r[6], p.r[7] = ch(p.r[4], p.dr[4], p.a[4], p.da[4], invSa4, invDa4), ch(p.r[5], p.dr[5], p.a[5], p.da[5], invSa5, invDa5), ch(p.r[6], p.dr[6], p.a[6], p.da[6], invSa6, invDa6), ch(p.r[7], p.dr[7], p.a[7], p.da[7], invSa7, invDa7)
	p.r[8], p.r[9], p.r[10], p.r[11] = ch(p.r[8], p.dr[8], p.a[8], p.da[8], invSa8, invDa8), ch(p.r[9], p.dr[9], p.a[9], p.da[9], invSa9, invDa9), ch(p.r[10], p.dr[10], p.a[10], p.da[10], invSa10, invDa10), ch(p.r[11], p.dr[11], p.a[11], p.da[11], invSa11, invDa11)
	p.r[12], p.r[13], p.r[14], p.r[15] = ch(p.r[12], p.dr[12], p.a[12], p.da[12], invSa12, invDa12), ch(p.r[13], p.dr[13], p.a[13], p.da[13], invSa13, invDa13), ch(p.r[14], p.dr[14], p.a[14], p.da[14], invSa14, invDa14), ch(p.r[15], p.dr[15], p.a[15], p.da[15], invSa15, invDa15)

	// Blend G channel (4 at a time)
	p.g[0], p.g[1], p.g[2], p.g[3] = ch(p.g[0], p.dg[0], p.a[0], p.da[0], invSa0, invDa0), ch(p.g[1], p.dg[1], p.a[1], p.da[1], invSa1, invDa1), ch(p.g[2], p.dg[2], p.a[2], p.da[2], invSa2, invDa2), ch(p.g[3], p.dg[3], p.a[3], p.da[3], invSa3, invDa3)
	p.g[4], p.g[5], p.g[6], p.g[7] = ch(p.g[4], p.dg[4], p.a[4], p.da[4], invSa4, invDa4), ch(p.g[5], p.dg[5], p.a[5], p.da[5], invSa5, invDa5), ch(p.g[6], p.dg[6], p.a[6], p.da[6], invSa6, invDa6), ch(p.g[7], p.dg[7], p.a[7], p.da[7], invSa7, invDa7)
	p.g[8], p.g[9], p.g[10], p.g[11] = ch(p.g[8], p.dg[8], p.a[8], p.da[8], invSa8, invDa8), ch(p.g[9], p.dg[9], p.a[9], p.da[9], invSa9, invDa9), ch(p.g[10], p.dg[10], p.a[10], p.da[10], invSa10, invDa10), ch(p.g[11], p.dg[11], p.a[11], p.da[11], invSa11, invDa11)
	p.g[12], p.g[13], p.g[14], p.g[15] = ch(p.g[12], p.dg[12], p.a[12], p.da[12], invSa12, invDa12), ch(p.g[13], p.dg[13], p.a[13], p.da[13], invSa13, invDa13), ch(p.g[14], p.dg[14], p.a[14], p.da[14], invSa14, invDa14), ch(p.g[15], p.dg[15], p.a[15], p.da[15], invSa15, invDa15)

	// Blend B channel (4 at a time)
	p.b[0], p.b[1], p.b[2], p.b[3] = ch(p.b[0], p.db[0], p.a[0], p.da[0], invSa0, invDa0), ch(p.b[1], p.db[1], p.a[1], p.da[1], invSa1, invDa1), ch(p.b[2], p.db[2], p.a[2], p.da[2], invSa2, invDa2), ch(p.b[3], p.db[3], p.a[3], p.da[3], invSa3, invDa3)
	p.b[4], p.b[5], p.b[6], p.b[7] = ch(p.b[4], p.db[4], p.a[4], p.da[4], invSa4, invDa4), ch(p.b[5], p.db[5], p.a[5], p.da[5], invSa5, invDa5), ch(p.b[6], p.db[6], p.a[6], p.da[6], invSa6, invDa6), ch(p.b[7], p.db[7], p.a[7], p.da[7], invSa7, invDa7)
	p.b[8], p.b[9], p.b[10], p.b[11] = ch(p.b[8], p.db[8], p.a[8], p.da[8], invSa8, invDa8), ch(p.b[9], p.db[9], p.a[9], p.da[9], invSa9, invDa9), ch(p.b[10], p.db[10], p.a[10], p.da[10], invSa10, invDa10), ch(p.b[11], p.db[11], p.a[11], p.da[11], invSa11, invDa11)
	p.b[12], p.b[13], p.b[14], p.b[15] = ch(p.b[12], p.db[12], p.a[12], p.da[12], invSa12, invDa12), ch(p.b[13], p.db[13], p.a[13], p.da[13], invSa13, invDa13), ch(p.b[14], p.db[14], p.a[14], p.da[14], invSa14, invDa14), ch(p.b[15], p.db[15], p.a[15], p.da[15], invSa15, invDa15)

	// Alpha channel: source_over formula (4 at a time)
	p.a[0], p.a[1], p.a[2], p.a[3] = p.a[0]+uint16((uint32(p.da[0])*uint32(invSa0)+255)>>8), p.a[1]+uint16((uint32(p.da[1])*uint32(invSa1)+255)>>8), p.a[2]+uint16((uint32(p.da[2])*uint32(invSa2)+255)>>8), p.a[3]+uint16((uint32(p.da[3])*uint32(invSa3)+255)>>8)
	p.a[4], p.a[5], p.a[6], p.a[7] = p.a[4]+uint16((uint32(p.da[4])*uint32(invSa4)+255)>>8), p.a[5]+uint16((uint32(p.da[5])*uint32(invSa5)+255)>>8), p.a[6]+uint16((uint32(p.da[6])*uint32(invSa6)+255)>>8), p.a[7]+uint16((uint32(p.da[7])*uint32(invSa7)+255)>>8)
	p.a[8], p.a[9], p.a[10], p.a[11] = p.a[8]+uint16((uint32(p.da[8])*uint32(invSa8)+255)>>8), p.a[9]+uint16((uint32(p.da[9])*uint32(invSa9)+255)>>8), p.a[10]+uint16((uint32(p.da[10])*uint32(invSa10)+255)>>8), p.a[11]+uint16((uint32(p.da[11])*uint32(invSa11)+255)>>8)
	p.a[12], p.a[13], p.a[14], p.a[15] = p.a[12]+uint16((uint32(p.da[12])*uint32(invSa12)+255)>>8), p.a[13]+uint16((uint32(p.da[13])*uint32(invSa13)+255)>>8), p.a[14]+uint16((uint32(p.da[14])*uint32(invSa14)+255)>>8), p.a[15]+uint16((uint32(p.da[15])*uint32(invSa15)+255)>>8)
}

//go:fix inline
func (p *LowPipeline) Lighten() {
	// Formula: s + d - div255(min(s * da, d * sa))
	ch := func(s, d, sa, da uint16) uint16 {
		prod1 := uint32(s) * uint32(da)
		prod2 := uint32(d) * uint32(sa)
		minProd := u16min(uint16(prod2), uint16(prod1))
		return uint16(uint32(s) + uint32(d) - ((uint32(minProd) + 255) >> 8))
	}

	// Blend R channel (4 at a time)
	p.r[0], p.r[1], p.r[2], p.r[3] = ch(p.r[0], p.dr[0], p.a[0], p.da[0]), ch(p.r[1], p.dr[1], p.a[1], p.da[1]), ch(p.r[2], p.dr[2], p.a[2], p.da[2]), ch(p.r[3], p.dr[3], p.a[3], p.da[3])
	p.r[4], p.r[5], p.r[6], p.r[7] = ch(p.r[4], p.dr[4], p.a[4], p.da[4]), ch(p.r[5], p.dr[5], p.a[5], p.da[5]), ch(p.r[6], p.dr[6], p.a[6], p.da[6]), ch(p.r[7], p.dr[7], p.a[7], p.da[7])
	p.r[8], p.r[9], p.r[10], p.r[11] = ch(p.r[8], p.dr[8], p.a[8], p.da[8]), ch(p.r[9], p.dr[9], p.a[9], p.da[9]), ch(p.r[10], p.dr[10], p.a[10], p.da[10]), ch(p.r[11], p.dr[11], p.a[11], p.da[11])
	p.r[12], p.r[13], p.r[14], p.r[15] = ch(p.r[12], p.dr[12], p.a[12], p.da[12]), ch(p.r[13], p.dr[13], p.a[13], p.da[13]), ch(p.r[14], p.dr[14], p.a[14], p.da[14]), ch(p.r[15], p.dr[15], p.a[15], p.da[15])

	// Blend G channel (4 at a time)
	p.g[0], p.g[1], p.g[2], p.g[3] = ch(p.g[0], p.dg[0], p.a[0], p.da[0]), ch(p.g[1], p.dg[1], p.a[1], p.da[1]), ch(p.g[2], p.dg[2], p.a[2], p.da[2]), ch(p.g[3], p.dg[3], p.a[3], p.da[3])
	p.g[4], p.g[5], p.g[6], p.g[7] = ch(p.g[4], p.dg[4], p.a[4], p.da[4]), ch(p.g[5], p.dg[5], p.a[5], p.da[5]), ch(p.g[6], p.dg[6], p.a[6], p.da[6]), ch(p.g[7], p.dg[7], p.a[7], p.da[7])
	p.g[8], p.g[9], p.g[10], p.g[11] = ch(p.g[8], p.dg[8], p.a[8], p.da[8]), ch(p.g[9], p.dg[9], p.a[9], p.da[9]), ch(p.g[10], p.dg[10], p.a[10], p.da[10]), ch(p.g[11], p.dg[11], p.a[11], p.da[11])
	p.g[12], p.g[13], p.g[14], p.g[15] = ch(p.g[12], p.dg[12], p.a[12], p.da[12]), ch(p.g[13], p.dg[13], p.a[13], p.da[13]), ch(p.g[14], p.dg[14], p.a[14], p.da[14]), ch(p.g[15], p.dg[15], p.a[15], p.da[15])

	// Blend B channel (4 at a time)
	p.b[0], p.b[1], p.b[2], p.b[3] = ch(p.b[0], p.db[0], p.a[0], p.da[0]), ch(p.b[1], p.db[1], p.a[1], p.da[1]), ch(p.b[2], p.db[2], p.a[2], p.da[2]), ch(p.b[3], p.db[3], p.a[3], p.da[3])
	p.b[4], p.b[5], p.b[6], p.b[7] = ch(p.b[4], p.db[4], p.a[4], p.da[4]), ch(p.b[5], p.db[5], p.a[5], p.da[5]), ch(p.b[6], p.db[6], p.a[6], p.da[6]), ch(p.b[7], p.db[7], p.a[7], p.da[7])
	p.b[8], p.b[9], p.b[10], p.b[11] = ch(p.b[8], p.db[8], p.a[8], p.da[8]), ch(p.b[9], p.db[9], p.a[9], p.da[9]), ch(p.b[10], p.db[10], p.a[10], p.da[10]), ch(p.b[11], p.db[11], p.a[11], p.da[11])
	p.b[12], p.b[13], p.b[14], p.b[15] = ch(p.b[12], p.db[12], p.a[12], p.da[12]), ch(p.b[13], p.db[13], p.a[13], p.da[13]), ch(p.b[14], p.db[14], p.a[14], p.da[14]), ch(p.b[15], p.db[15], p.a[15], p.da[15])

	// Alpha channel: source_over formula (4 at a time)
	invSa0, invSa1, invSa2, invSa3 := uint16(255-p.a[0]), uint16(255-p.a[1]), uint16(255-p.a[2]), uint16(255-p.a[3])
	invSa4, invSa5, invSa6, invSa7 := uint16(255-p.a[4]), uint16(255-p.a[5]), uint16(255-p.a[6]), uint16(255-p.a[7])
	invSa8, invSa9, invSa10, invSa11 := uint16(255-p.a[8]), uint16(255-p.a[9]), uint16(255-p.a[10]), uint16(255-p.a[11])
	invSa12, invSa13, invSa14, invSa15 := uint16(255-p.a[12]), uint16(255-p.a[13]), uint16(255-p.a[14]), uint16(255-p.a[15])

	p.a[0], p.a[1], p.a[2], p.a[3] = p.a[0]+uint16((uint32(p.da[0])*uint32(invSa0)+255)>>8), p.a[1]+uint16((uint32(p.da[1])*uint32(invSa1)+255)>>8), p.a[2]+uint16((uint32(p.da[2])*uint32(invSa2)+255)>>8), p.a[3]+uint16((uint32(p.da[3])*uint32(invSa3)+255)>>8)
	p.a[4], p.a[5], p.a[6], p.a[7] = p.a[4]+uint16((uint32(p.da[4])*uint32(invSa4)+255)>>8), p.a[5]+uint16((uint32(p.da[5])*uint32(invSa5)+255)>>8), p.a[6]+uint16((uint32(p.da[6])*uint32(invSa6)+255)>>8), p.a[7]+uint16((uint32(p.da[7])*uint32(invSa7)+255)>>8)
	p.a[8], p.a[9], p.a[10], p.a[11] = p.a[8]+uint16((uint32(p.da[8])*uint32(invSa8)+255)>>8), p.a[9]+uint16((uint32(p.da[9])*uint32(invSa9)+255)>>8), p.a[10]+uint16((uint32(p.da[10])*uint32(invSa10)+255)>>8), p.a[11]+uint16((uint32(p.da[11])*uint32(invSa11)+255)>>8)
	p.a[12], p.a[13], p.a[14], p.a[15] = p.a[12]+uint16((uint32(p.da[12])*uint32(invSa12)+255)>>8), p.a[13]+uint16((uint32(p.da[13])*uint32(invSa13)+255)>>8), p.a[14]+uint16((uint32(p.da[14])*uint32(invSa14)+255)>>8), p.a[15]+uint16((uint32(p.da[15])*uint32(invSa15)+255)>>8)
}

//go:fix inline
func (p *LowPipeline) Overlay() {
	// Precompute inverse alpha values (4 at a time)
	invSa0, invSa1, invSa2, invSa3 := uint16(255-p.a[0]), uint16(255-p.a[1]), uint16(255-p.a[2]), uint16(255-p.a[3])
	invSa4, invSa5, invSa6, invSa7 := uint16(255-p.a[4]), uint16(255-p.a[5]), uint16(255-p.a[6]), uint16(255-p.a[7])
	invSa8, invSa9, invSa10, invSa11 := uint16(255-p.a[8]), uint16(255-p.a[9]), uint16(255-p.a[10]), uint16(255-p.a[11])
	invSa12, invSa13, invSa14, invSa15 := uint16(255-p.a[12]), uint16(255-p.a[13]), uint16(255-p.a[14]), uint16(255-p.a[15])

	invDa0, invDa1, invDa2, invDa3 := uint16(255-p.da[0]), uint16(255-p.da[1]), uint16(255-p.da[2]), uint16(255-p.da[3])
	invDa4, invDa5, invDa6, invDa7 := uint16(255-p.da[4]), uint16(255-p.da[5]), uint16(255-p.da[6]), uint16(255-p.da[7])
	invDa8, invDa9, invDa10, invDa11 := uint16(255-p.da[8]), uint16(255-p.da[9]), uint16(255-p.da[10]), uint16(255-p.da[11])
	invDa12, invDa13, invDa14, invDa15 := uint16(255-p.da[12]), uint16(255-p.da[13]), uint16(255-p.da[14]), uint16(255-p.da[15])

	ch := func(s, d, sa, da, invSa, invDa uint16) uint16 {
		s32, d32 := uint32(s), uint32(d)
		sa32, da32 := uint32(sa), uint32(da)
		invSa32, invDa32 := uint32(invSa), uint32(invDa)

		// Base term: s * inv(da) + d * inv(sa)
		term := s32*invDa32 + d32*invSa32

		// Branch term
		if 2*d32 <= da32 {
			term += 2 * s32 * d32
		} else {
			term += sa32*da32 - 2*(sa32-s32)*(da32-d32)
		}

		return uint16((term + 255) >> 8)
	}

	p.r[0], p.r[1], p.r[2], p.r[3] = ch(p.r[0], p.dr[0], p.a[0], p.da[0], invSa0, invDa0), ch(p.r[1], p.dr[1], p.a[1], p.da[1], invSa1, invDa1), ch(p.r[2], p.dr[2], p.a[2], p.da[2], invSa2, invDa2), ch(p.r[3], p.dr[3], p.a[3], p.da[3], invSa3, invDa3)
	p.r[4], p.r[5], p.r[6], p.r[7] = ch(p.r[4], p.dr[4], p.a[4], p.da[4], invSa4, invDa4), ch(p.r[5], p.dr[5], p.a[5], p.da[5], invSa5, invDa5), ch(p.r[6], p.dr[6], p.a[6], p.da[6], invSa6, invDa6), ch(p.r[7], p.dr[7], p.a[7], p.da[7], invSa7, invDa7)
	p.r[8], p.r[9], p.r[10], p.r[11] = ch(p.r[8], p.dr[8], p.a[8], p.da[8], invSa8, invDa8), ch(p.r[9], p.dr[9], p.a[9], p.da[9], invSa9, invDa9), ch(p.r[10], p.dr[10], p.a[10], p.da[10], invSa10, invDa10), ch(p.r[11], p.dr[11], p.a[11], p.da[11], invSa11, invDa11)
	p.r[12], p.r[13], p.r[14], p.r[15] = ch(p.r[12], p.dr[12], p.a[12], p.da[12], invSa12, invDa12), ch(p.r[13], p.dr[13], p.a[13], p.da[13], invSa13, invDa13), ch(p.r[14], p.dr[14], p.a[14], p.da[14], invSa14, invDa14), ch(p.r[15], p.dr[15], p.a[15], p.da[15], invSa15, invDa15)

	p.g[0], p.g[1], p.g[2], p.g[3] = ch(p.g[0], p.dg[0], p.a[0], p.da[0], invSa0, invDa0), ch(p.g[1], p.dg[1], p.a[1], p.da[1], invSa1, invDa1), ch(p.g[2], p.dg[2], p.a[2], p.da[2], invSa2, invDa2), ch(p.g[3], p.dg[3], p.a[3], p.da[3], invSa3, invDa3)
	p.g[4], p.g[5], p.g[6], p.g[7] = ch(p.g[4], p.dg[4], p.a[4], p.da[4], invSa4, invDa4), ch(p.g[5], p.dg[5], p.a[5], p.da[5], invSa5, invDa5), ch(p.g[6], p.dg[6], p.a[6], p.da[6], invSa6, invDa6), ch(p.g[7], p.dg[7], p.a[7], p.da[7], invSa7, invDa7)
	p.g[8], p.g[9], p.g[10], p.g[11] = ch(p.g[8], p.dg[8], p.a[8], p.da[8], invSa8, invDa8), ch(p.g[9], p.dg[9], p.a[9], p.da[9], invSa9, invDa9), ch(p.g[10], p.dg[10], p.a[10], p.da[10], invSa10, invDa10), ch(p.g[11], p.dg[11], p.a[11], p.da[11], invSa11, invDa11)
	p.g[12], p.g[13], p.g[14], p.g[15] = ch(p.g[12], p.dg[12], p.a[12], p.da[12], invSa12, invDa12), ch(p.g[13], p.dg[13], p.a[13], p.da[13], invSa13, invDa13), ch(p.g[14], p.dg[14], p.a[14], p.da[14], invSa14, invDa14), ch(p.g[15], p.dg[15], p.a[15], p.da[15], invSa15, invDa15)

	p.b[0], p.b[1], p.b[2], p.b[3] = ch(p.b[0], p.db[0], p.a[0], p.da[0], invSa0, invDa0), ch(p.b[1], p.db[1], p.a[1], p.da[1], invSa1, invDa1), ch(p.b[2], p.db[2], p.a[2], p.da[2], invSa2, invDa2), ch(p.b[3], p.db[3], p.a[3], p.da[3], invSa3, invDa3)
	p.b[4], p.b[5], p.b[6], p.b[7] = ch(p.b[4], p.db[4], p.a[4], p.da[4], invSa4, invDa4), ch(p.b[5], p.db[5], p.a[5], p.da[5], invSa5, invDa5), ch(p.b[6], p.db[6], p.a[6], p.da[6], invSa6, invDa6), ch(p.b[7], p.db[7], p.a[7], p.da[7], invSa7, invDa7)
	p.b[8], p.b[9], p.b[10], p.b[11] = ch(p.b[8], p.db[8], p.a[8], p.da[8], invSa8, invDa8), ch(p.b[9], p.db[9], p.a[9], p.da[9], invSa9, invDa9), ch(p.b[10], p.db[10], p.a[10], p.da[10], invSa10, invDa10), ch(p.b[11], p.db[11], p.a[11], p.da[11], invSa11, invDa11)
	p.b[12], p.b[13], p.b[14], p.b[15] = ch(p.b[12], p.db[12], p.a[12], p.da[12], invSa12, invDa12), ch(p.b[13], p.db[13], p.a[13], p.da[13], invSa13, invDa13), ch(p.b[14], p.db[14], p.a[14], p.da[14], invSa14, invDa14), ch(p.b[15], p.db[15], p.a[15], p.da[15], invSa15, invDa15)

	p.a[0], p.a[1], p.a[2], p.a[3] = p.a[0]+uint16((uint32(p.da[0])*uint32(invSa0)+255)>>8), p.a[1]+uint16((uint32(p.da[1])*uint32(invSa1)+255)>>8), p.a[2]+uint16((uint32(p.da[2])*uint32(invSa2)+255)>>8), p.a[3]+uint16((uint32(p.da[3])*uint32(invSa3)+255)>>8)
	p.a[4], p.a[5], p.a[6], p.a[7] = p.a[4]+uint16((uint32(p.da[4])*uint32(invSa4)+255)>>8), p.a[5]+uint16((uint32(p.da[5])*uint32(invSa5)+255)>>8), p.a[6]+uint16((uint32(p.da[6])*uint32(invSa6)+255)>>8), p.a[7]+uint16((uint32(p.da[7])*uint32(invSa7)+255)>>8)
	p.a[8], p.a[9], p.a[10], p.a[11] = p.a[8]+uint16((uint32(p.da[8])*uint32(invSa8)+255)>>8), p.a[9]+uint16((uint32(p.da[9])*uint32(invSa9)+255)>>8), p.a[10]+uint16((uint32(p.da[10])*uint32(invSa10)+255)>>8), p.a[11]+uint16((uint32(p.da[11])*uint32(invSa11)+255)>>8)
	p.a[12], p.a[13], p.a[14], p.a[15] = p.a[12]+uint16((uint32(p.da[12])*uint32(invSa12)+255)>>8), p.a[13]+uint16((uint32(p.da[13])*uint32(invSa13)+255)>>8), p.a[14]+uint16((uint32(p.da[14])*uint32(invSa14)+255)>>8), p.a[15]+uint16((uint32(p.da[15])*uint32(invSa15)+255)>>8)
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
	// Load destination RGBA (4 pixels at a time, unrolled)
	baseIdx := (p.dy*p.pixmap.RealWidth + p.dx) * 4
	data := p.pixmap.Data[baseIdx : baseIdx+LOW_STAGE_WIDTH*4]

	// Load dr (red channel)
	p.dr[0], p.dr[1], p.dr[2], p.dr[3] = uint16(data[0]), uint16(data[4]), uint16(data[8]), uint16(data[12])
	p.dr[4], p.dr[5], p.dr[6], p.dr[7] = uint16(data[16]), uint16(data[20]), uint16(data[24]), uint16(data[28])
	p.dr[8], p.dr[9], p.dr[10], p.dr[11] = uint16(data[32]), uint16(data[36]), uint16(data[40]), uint16(data[44])
	p.dr[12], p.dr[13], p.dr[14], p.dr[15] = uint16(data[48]), uint16(data[52]), uint16(data[56]), uint16(data[60])

	// Load dg (green channel)
	p.dg[0], p.dg[1], p.dg[2], p.dg[3] = uint16(data[1]), uint16(data[5]), uint16(data[9]), uint16(data[13])
	p.dg[4], p.dg[5], p.dg[6], p.dg[7] = uint16(data[17]), uint16(data[21]), uint16(data[25]), uint16(data[29])
	p.dg[8], p.dg[9], p.dg[10], p.dg[11] = uint16(data[33]), uint16(data[37]), uint16(data[41]), uint16(data[45])
	p.dg[12], p.dg[13], p.dg[14], p.dg[15] = uint16(data[49]), uint16(data[53]), uint16(data[57]), uint16(data[61])

	// Load db (blue channel)
	p.db[0], p.db[1], p.db[2], p.db[3] = uint16(data[2]), uint16(data[6]), uint16(data[10]), uint16(data[14])
	p.db[4], p.db[5], p.db[6], p.db[7] = uint16(data[18]), uint16(data[22]), uint16(data[26]), uint16(data[30])
	p.db[8], p.db[9], p.db[10], p.db[11] = uint16(data[34]), uint16(data[38]), uint16(data[42]), uint16(data[46])
	p.db[12], p.db[13], p.db[14], p.db[15] = uint16(data[50]), uint16(data[54]), uint16(data[58]), uint16(data[62])

	// Load da (alpha channel)
	p.da[0], p.da[1], p.da[2], p.da[3] = uint16(data[3]), uint16(data[7]), uint16(data[11]), uint16(data[15])
	p.da[4], p.da[5], p.da[6], p.da[7] = uint16(data[19]), uint16(data[23]), uint16(data[27]), uint16(data[31])
	p.da[8], p.da[9], p.da[10], p.da[11] = uint16(data[35]), uint16(data[39]), uint16(data[43]), uint16(data[47])
	p.da[12], p.da[13], p.da[14], p.da[15] = uint16(data[51]), uint16(data[55]), uint16(data[59]), uint16(data[63])

	// source_over blend: src + dst * inv(src_alpha)
	// inv(a) = 255 - a
	// div255(v) = (v + 255) >> 8
	invSa0, invSa1, invSa2, invSa3 := uint16(255-p.a[0]), uint16(255-p.a[1]), uint16(255-p.a[2]), uint16(255-p.a[3])
	invSa4, invSa5, invSa6, invSa7 := uint16(255-p.a[4]), uint16(255-p.a[5]), uint16(255-p.a[6]), uint16(255-p.a[7])
	invSa8, invSa9, invSa10, invSa11 := uint16(255-p.a[8]), uint16(255-p.a[9]), uint16(255-p.a[10]), uint16(255-p.a[11])
	invSa12, invSa13, invSa14, invSa15 := uint16(255-p.a[12]), uint16(255-p.a[13]), uint16(255-p.a[14]), uint16(255-p.a[15])

	// Blend all 16 pixels (4 at a time, unrolled)
	p.r[0], p.r[1], p.r[2], p.r[3] = p.r[0]+uint16((uint32(p.dr[0])*uint32(invSa0)+255)>>8), p.r[1]+uint16((uint32(p.dr[1])*uint32(invSa1)+255)>>8), p.r[2]+uint16((uint32(p.dr[2])*uint32(invSa2)+255)>>8), p.r[3]+uint16((uint32(p.dr[3])*uint32(invSa3)+255)>>8)
	p.r[4], p.r[5], p.r[6], p.r[7] = p.r[4]+uint16((uint32(p.dr[4])*uint32(invSa4)+255)>>8), p.r[5]+uint16((uint32(p.dr[5])*uint32(invSa5)+255)>>8), p.r[6]+uint16((uint32(p.dr[6])*uint32(invSa6)+255)>>8), p.r[7]+uint16((uint32(p.dr[7])*uint32(invSa7)+255)>>8)
	p.r[8], p.r[9], p.r[10], p.r[11] = p.r[8]+uint16((uint32(p.dr[8])*uint32(invSa8)+255)>>8), p.r[9]+uint16((uint32(p.dr[9])*uint32(invSa9)+255)>>8), p.r[10]+uint16((uint32(p.dr[10])*uint32(invSa10)+255)>>8), p.r[11]+uint16((uint32(p.dr[11])*uint32(invSa11)+255)>>8)
	p.r[12], p.r[13], p.r[14], p.r[15] = p.r[12]+uint16((uint32(p.dr[12])*uint32(invSa12)+255)>>8), p.r[13]+uint16((uint32(p.dr[13])*uint32(invSa13)+255)>>8), p.r[14]+uint16((uint32(p.dr[14])*uint32(invSa14)+255)>>8), p.r[15]+uint16((uint32(p.dr[15])*uint32(invSa15)+255)>>8)

	p.g[0], p.g[1], p.g[2], p.g[3] = p.g[0]+uint16((uint32(p.dg[0])*uint32(invSa0)+255)>>8), p.g[1]+uint16((uint32(p.dg[1])*uint32(invSa1)+255)>>8), p.g[2]+uint16((uint32(p.dg[2])*uint32(invSa2)+255)>>8), p.g[3]+uint16((uint32(p.dg[3])*uint32(invSa3)+255)>>8)
	p.g[4], p.g[5], p.g[6], p.g[7] = p.g[4]+uint16((uint32(p.dg[4])*uint32(invSa4)+255)>>8), p.g[5]+uint16((uint32(p.dg[5])*uint32(invSa5)+255)>>8), p.g[6]+uint16((uint32(p.dg[6])*uint32(invSa6)+255)>>8), p.g[7]+uint16((uint32(p.dg[7])*uint32(invSa7)+255)>>8)
	p.g[8], p.g[9], p.g[10], p.g[11] = p.g[8]+uint16((uint32(p.dg[8])*uint32(invSa8)+255)>>8), p.g[9]+uint16((uint32(p.dg[9])*uint32(invSa9)+255)>>8), p.g[10]+uint16((uint32(p.dg[10])*uint32(invSa10)+255)>>8), p.g[11]+uint16((uint32(p.dg[11])*uint32(invSa11)+255)>>8)
	p.g[12], p.g[13], p.g[14], p.g[15] = p.g[12]+uint16((uint32(p.dg[12])*uint32(invSa12)+255)>>8), p.g[13]+uint16((uint32(p.dg[13])*uint32(invSa13)+255)>>8), p.g[14]+uint16((uint32(p.dg[14])*uint32(invSa14)+255)>>8), p.g[15]+uint16((uint32(p.dg[15])*uint32(invSa15)+255)>>8)

	p.b[0], p.b[1], p.b[2], p.b[3] = p.b[0]+uint16((uint32(p.db[0])*uint32(invSa0)+255)>>8), p.b[1]+uint16((uint32(p.db[1])*uint32(invSa1)+255)>>8), p.b[2]+uint16((uint32(p.db[2])*uint32(invSa2)+255)>>8), p.b[3]+uint16((uint32(p.db[3])*uint32(invSa3)+255)>>8)
	p.b[4], p.b[5], p.b[6], p.b[7] = p.b[4]+uint16((uint32(p.db[4])*uint32(invSa4)+255)>>8), p.b[5]+uint16((uint32(p.db[5])*uint32(invSa5)+255)>>8), p.b[6]+uint16((uint32(p.db[6])*uint32(invSa6)+255)>>8), p.b[7]+uint16((uint32(p.db[7])*uint32(invSa7)+255)>>8)
	p.b[8], p.b[9], p.b[10], p.b[11] = p.b[8]+uint16((uint32(p.db[8])*uint32(invSa8)+255)>>8), p.b[9]+uint16((uint32(p.db[9])*uint32(invSa9)+255)>>8), p.b[10]+uint16((uint32(p.db[10])*uint32(invSa10)+255)>>8), p.b[11]+uint16((uint32(p.db[11])*uint32(invSa11)+255)>>8)
	p.b[12], p.b[13], p.b[14], p.b[15] = p.b[12]+uint16((uint32(p.db[12])*uint32(invSa12)+255)>>8), p.b[13]+uint16((uint32(p.db[13])*uint32(invSa13)+255)>>8), p.b[14]+uint16((uint32(p.db[14])*uint32(invSa14)+255)>>8), p.b[15]+uint16((uint32(p.db[15])*uint32(invSa15)+255)>>8)

	p.a[0], p.a[1], p.a[2], p.a[3] = p.a[0]+uint16((uint32(p.da[0])*uint32(invSa0)+255)>>8), p.a[1]+uint16((uint32(p.da[1])*uint32(invSa1)+255)>>8), p.a[2]+uint16((uint32(p.da[2])*uint32(invSa2)+255)>>8), p.a[3]+uint16((uint32(p.da[3])*uint32(invSa3)+255)>>8)
	p.a[4], p.a[5], p.a[6], p.a[7] = p.a[4]+uint16((uint32(p.da[4])*uint32(invSa4)+255)>>8), p.a[5]+uint16((uint32(p.da[5])*uint32(invSa5)+255)>>8), p.a[6]+uint16((uint32(p.da[6])*uint32(invSa6)+255)>>8), p.a[7]+uint16((uint32(p.da[7])*uint32(invSa7)+255)>>8)
	p.a[8], p.a[9], p.a[10], p.a[11] = p.a[8]+uint16((uint32(p.da[8])*uint32(invSa8)+255)>>8), p.a[9]+uint16((uint32(p.da[9])*uint32(invSa9)+255)>>8), p.a[10]+uint16((uint32(p.da[10])*uint32(invSa10)+255)>>8), p.a[11]+uint16((uint32(p.da[11])*uint32(invSa11)+255)>>8)
	p.a[12], p.a[13], p.a[14], p.a[15] = p.a[12]+uint16((uint32(p.da[12])*uint32(invSa12)+255)>>8), p.a[13]+uint16((uint32(p.da[13])*uint32(invSa13)+255)>>8), p.a[14]+uint16((uint32(p.da[14])*uint32(invSa14)+255)>>8), p.a[15]+uint16((uint32(p.da[15])*uint32(invSa15)+255)>>8)

	// Store result back to destination (4 pixels at a time, unrolled)
	data[0], data[4], data[8], data[12] = uint8(p.r[0]), uint8(p.r[1]), uint8(p.r[2]), uint8(p.r[3])
	data[16], data[20], data[24], data[28] = uint8(p.r[4]), uint8(p.r[5]), uint8(p.r[6]), uint8(p.r[7])
	data[32], data[36], data[40], data[44] = uint8(p.r[8]), uint8(p.r[9]), uint8(p.r[10]), uint8(p.r[11])
	data[48], data[52], data[56], data[60] = uint8(p.r[12]), uint8(p.r[13]), uint8(p.r[14]), uint8(p.r[15])

	data[1], data[5], data[9], data[13] = uint8(p.g[0]), uint8(p.g[1]), uint8(p.g[2]), uint8(p.g[3])
	data[17], data[21], data[25], data[29] = uint8(p.g[4]), uint8(p.g[5]), uint8(p.g[6]), uint8(p.g[7])
	data[33], data[37], data[41], data[45] = uint8(p.g[8]), uint8(p.g[9]), uint8(p.g[10]), uint8(p.g[11])
	data[49], data[53], data[57], data[61] = uint8(p.g[12]), uint8(p.g[13]), uint8(p.g[14]), uint8(p.g[15])

	data[2], data[6], data[10], data[14] = uint8(p.b[0]), uint8(p.b[1]), uint8(p.b[2]), uint8(p.b[3])
	data[18], data[22], data[26], data[30] = uint8(p.b[4]), uint8(p.b[5]), uint8(p.b[6]), uint8(p.b[7])
	data[34], data[38], data[42], data[46] = uint8(p.b[8]), uint8(p.b[9]), uint8(p.b[10]), uint8(p.b[11])
	data[50], data[54], data[58], data[62] = uint8(p.b[12]), uint8(p.b[13]), uint8(p.b[14]), uint8(p.b[15])

	data[3], data[7], data[11], data[15] = uint8(p.a[0]), uint8(p.a[1]), uint8(p.a[2]), uint8(p.a[3])
	data[19], data[23], data[27], data[31] = uint8(p.a[4]), uint8(p.a[5]), uint8(p.a[6]), uint8(p.a[7])
	data[35], data[39], data[43], data[47] = uint8(p.a[8]), uint8(p.a[9]), uint8(p.a[10]), uint8(p.a[11])
	data[51], data[55], data[59], data[63] = uint8(p.a[12]), uint8(p.a[13]), uint8(p.a[14]), uint8(p.a[15])
}

//go:fix inline
func (p *LowPipeline) SourceOverRgbaTail() {
	// Load destination RGBA for tail pixels
	baseIdx := (p.dy*p.pixmap.RealWidth + p.dx) * 4
	data := p.pixmap.Data[baseIdx : baseIdx+p.tail*4]

	for i := 0; i < p.tail; i++ {
		off := i * 4
		p.dr[i] = uint16(data[off])
		p.dg[i] = uint16(data[off+1])
		p.db[i] = uint16(data[off+2])
		p.da[i] = uint16(data[off+3])
	}

	// source_over blend: src + dst * inv(src_alpha)
	for i := 0; i < p.tail; i++ {
		invSa := uint16(255 - p.a[i])
		p.r[i] = p.r[i] + uint16((uint32(p.dr[i])*uint32(invSa)+255)>>8)
		p.g[i] = p.g[i] + uint16((uint32(p.dg[i])*uint32(invSa)+255)>>8)
		p.b[i] = p.b[i] + uint16((uint32(p.db[i])*uint32(invSa)+255)>>8)
		p.a[i] = p.a[i] + uint16((uint32(p.da[i])*uint32(invSa)+255)>>8)
	}

	// Store result back
	for i := 0; i < p.tail; i++ {
		off := i * 4
		data[off] = uint8(p.r[i])
		data[off+1] = uint8(p.g[i])
		data[off+2] = uint8(p.b[i])
		data[off+3] = uint8(p.a[i])
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
		x[i] = math32.Float32frombits(uint32(p.r[i]) | uint32(p.g[i])<<16)
		y[i] = math32.Float32frombits(uint32(p.b[i]) | uint32(p.a[i])<<16)
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
		nxBits := math32.Float32bits(nx[i])
		nyBits := math32.Float32bits(ny[i])
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
		x := math32.Float32frombits(uint32(p.r[i]) | uint32(p.g[i])<<16)
		// Normalize/clamp to [0, 1] range
		if x < 0.0 {
			x = 0.0
		} else if x > 1.0 {
			x = 1.0
		}
		// Convert back to u16 representation
		xBits := math32.Float32bits(x)
		p.r[i] = uint16(xBits & 0xFFFF)
		p.g[i] = uint16(xBits >> 16)
	}
}

//go:fix inline
func (p *LowPipeline) ReflectX1() {
	// Mirrors x at integer boundaries: x = |x - 1 - 2*floor((x-1)*0.5)| - 1
	for i := 0; i < LOW_STAGE_WIDTH; i++ {
		// Convert to float32
		x := math32.Float32frombits(uint32(p.r[i]) | uint32(p.g[i])<<16)
		// Apply reflect formula: x = |x - 1 - 2*floor((x-1)*0.5)| - 1
		xMinus1 := x - 1.0
		floored := float32(int(xMinus1 * 0.5))
		reflected := xMinus1 - 2.0*floored - 1.0
		if reflected < 0 {
			reflected = -reflected
		}
		// Convert back to u16 representation
		xBits := math32.Float32bits(reflected)
		p.r[i] = uint16(xBits & 0xFFFF)
		p.g[i] = uint16(xBits >> 16)
	}
}

//go:fix inline
func (p *LowPipeline) RepeatX1() {
	// Repeats pattern every integer: x = x - floor(x)
	for i := 0; i < LOW_STAGE_WIDTH; i++ {
		// Convert to float32
		x := math32.Float32frombits(uint32(p.r[i]) | uint32(p.g[i])<<16)
		// Apply repeat formula: x = x - floor(x)
		floored := float32(int(x))
		if floored > x {
			floored--
		}
		repeated := x - floored
		// Convert back to u16 representation
		xBits := math32.Float32bits(repeated)
		p.r[i] = uint16(xBits & 0xFFFF)
		p.g[i] = uint16(xBits >> 16)
	}
}

//go:fix inline
func (p *LowPipeline) Gradient() {
	ctx := p.ctx.Gradient

	// Join r,g into t values as float32
	var t [16]float32
	for i := 0; i < 16; i++ {
		t[i] = math32.Float32frombits(uint32(p.r[i]) | uint32(p.g[i])<<16)
	}

	// Find stop indices for all pixels
	var idx [16]uint16
	for j := 1; j < ctx.Len; j++ {
		tt := ctx.TValues[j].Get()
		if t[0] >= tt {
			idx[0]++
		}
		if t[1] >= tt {
			idx[1]++
		}
		if t[2] >= tt {
			idx[2]++
		}
		if t[3] >= tt {
			idx[3]++
		}
		if t[4] >= tt {
			idx[4]++
		}
		if t[5] >= tt {
			idx[5]++
		}
		if t[6] >= tt {
			idx[6]++
		}
		if t[7] >= tt {
			idx[7]++
		}
		if t[8] >= tt {
			idx[8]++
		}
		if t[9] >= tt {
			idx[9]++
		}
		if t[10] >= tt {
			idx[10]++
		}
		if t[11] >= tt {
			idx[11]++
		}
		if t[12] >= tt {
			idx[12]++
		}
		if t[13] >= tt {
			idx[13]++
		}
		if t[14] >= tt {
			idx[14]++
		}
		if t[15] >= tt {
			idx[15]++
		}
	}

	f0, f1, f2, f3 := ctx.Factors[idx[0]], ctx.Factors[idx[1]], ctx.Factors[idx[2]], ctx.Factors[idx[3]]
	f4, f5, f6, f7 := ctx.Factors[idx[4]], ctx.Factors[idx[5]], ctx.Factors[idx[6]], ctx.Factors[idx[7]]
	f8, f9, f10, f11 := ctx.Factors[idx[8]], ctx.Factors[idx[9]], ctx.Factors[idx[10]], ctx.Factors[idx[11]]
	f12, f13, f14, f15 := ctx.Factors[idx[12]], ctx.Factors[idx[13]], ctx.Factors[idx[14]], ctx.Factors[idx[15]]

	b0, b1, b2, b3 := ctx.Biases[idx[0]], ctx.Biases[idx[1]], ctx.Biases[idx[2]], ctx.Biases[idx[3]]
	b4, b5, b6, b7 := ctx.Biases[idx[4]], ctx.Biases[idx[5]], ctx.Biases[idx[6]], ctx.Biases[idx[7]]
	b8, b9, b10, b11 := ctx.Biases[idx[8]], ctx.Biases[idx[9]], ctx.Biases[idx[10]], ctx.Biases[idx[11]]
	b12, b13, b14, b15 := ctx.Biases[idx[12]], ctx.Biases[idx[13]], ctx.Biases[idx[14]], ctx.Biases[idx[15]]

	p.r[0], p.r[1], p.r[2], p.r[3] = uint16((t[0]*f0.R+b0.R)*255.0+0.5), uint16((t[1]*f1.R+b1.R)*255.0+0.5), uint16((t[2]*f2.R+b2.R)*255.0+0.5), uint16((t[3]*f3.R+b3.R)*255.0+0.5)
	p.r[4], p.r[5], p.r[6], p.r[7] = uint16((t[4]*f4.R+b4.R)*255.0+0.5), uint16((t[5]*f5.R+b5.R)*255.0+0.5), uint16((t[6]*f6.R+b6.R)*255.0+0.5), uint16((t[7]*f7.R+b7.R)*255.0+0.5)
	p.r[8], p.r[9], p.r[10], p.r[11] = uint16((t[8]*f8.R+b8.R)*255.0+0.5), uint16((t[9]*f9.R+b9.R)*255.0+0.5), uint16((t[10]*f10.R+b10.R)*255.0+0.5), uint16((t[11]*f11.R+b11.R)*255.0+0.5)
	p.r[12], p.r[13], p.r[14], p.r[15] = uint16((t[12]*f12.R+b12.R)*255.0+0.5), uint16((t[13]*f13.R+b13.R)*255.0+0.5), uint16((t[14]*f14.R+b14.R)*255.0+0.5), uint16((t[15]*f15.R+b15.R)*255.0+0.5)

	p.g[0], p.g[1], p.g[2], p.g[3] = uint16((t[0]*f0.G+b0.G)*255.0+0.5), uint16((t[1]*f1.G+b1.G)*255.0+0.5), uint16((t[2]*f2.G+b2.G)*255.0+0.5), uint16((t[3]*f3.G+b3.G)*255.0+0.5)
	p.g[4], p.g[5], p.g[6], p.g[7] = uint16((t[4]*f4.G+b4.G)*255.0+0.5), uint16((t[5]*f5.G+b5.G)*255.0+0.5), uint16((t[6]*f6.G+b6.G)*255.0+0.5), uint16((t[7]*f7.G+b7.G)*255.0+0.5)
	p.g[8], p.g[9], p.g[10], p.g[11] = uint16((t[8]*f8.G+b8.G)*255.0+0.5), uint16((t[9]*f9.G+b9.G)*255.0+0.5), uint16((t[10]*f10.G+b10.G)*255.0+0.5), uint16((t[11]*f11.G+b11.G)*255.0+0.5)
	p.g[12], p.g[13], p.g[14], p.g[15] = uint16((t[12]*f12.G+b12.G)*255.0+0.5), uint16((t[13]*f13.G+b13.G)*255.0+0.5), uint16((t[14]*f14.G+b14.G)*255.0+0.5), uint16((t[15]*f15.G+b15.G)*255.0+0.5)

	p.b[0], p.b[1], p.b[2], p.b[3] = uint16((t[0]*f0.B+b0.B)*255.0+0.5), uint16((t[1]*f1.B+b1.B)*255.0+0.5), uint16((t[2]*f2.B+b2.B)*255.0+0.5), uint16((t[3]*f3.B+b3.B)*255.0+0.5)
	p.b[4], p.b[5], p.b[6], p.b[7] = uint16((t[4]*f4.B+b4.B)*255.0+0.5), uint16((t[5]*f5.B+b5.B)*255.0+0.5), uint16((t[6]*f6.B+b6.B)*255.0+0.5), uint16((t[7]*f7.B+b7.B)*255.0+0.5)
	p.b[8], p.b[9], p.b[10], p.b[11] = uint16((t[8]*f8.B+b8.B)*255.0+0.5), uint16((t[9]*f9.B+b9.B)*255.0+0.5), uint16((t[10]*f10.B+b10.B)*255.0+0.5), uint16((t[11]*f11.B+b11.B)*255.0+0.5)
	p.b[12], p.b[13], p.b[14], p.b[15] = uint16((t[12]*f12.B+b12.B)*255.0+0.5), uint16((t[13]*f13.B+b13.B)*255.0+0.5), uint16((t[14]*f14.B+b14.B)*255.0+0.5), uint16((t[15]*f15.B+b15.B)*255.0+0.5)

	p.a[0], p.a[1], p.a[2], p.a[3] = uint16((t[0]*f0.A+b0.A)*255.0+0.5), uint16((t[1]*f1.A+b1.A)*255.0+0.5), uint16((t[2]*f2.A+b2.A)*255.0+0.5), uint16((t[3]*f3.A+b3.A)*255.0+0.5)
	p.a[4], p.a[5], p.a[6], p.a[7] = uint16((t[4]*f4.A+b4.A)*255.0+0.5), uint16((t[5]*f5.A+b5.A)*255.0+0.5), uint16((t[6]*f6.A+b6.A)*255.0+0.5), uint16((t[7]*f7.A+b7.A)*255.0+0.5)
	p.a[8], p.a[9], p.a[10], p.a[11] = uint16((t[8]*f8.A+b8.A)*255.0+0.5), uint16((t[9]*f9.A+b9.A)*255.0+0.5), uint16((t[10]*f10.A+b10.A)*255.0+0.5), uint16((t[11]*f11.A+b11.A)*255.0+0.5)
	p.a[12], p.a[13], p.a[14], p.a[15] = uint16((t[12]*f12.A+b12.A)*255.0+0.5), uint16((t[13]*f13.A+b13.A)*255.0+0.5), uint16((t[14]*f14.A+b14.A)*255.0+0.5), uint16((t[15]*f15.A+b15.A)*255.0+0.5)
}

//go:fix inline
func (p *LowPipeline) EvenlySpaced2StopGradient() {
	factor := p.ctx.EvenlySpaced2StopGradient.Factor
	bias := p.ctx.EvenlySpaced2StopGradient.Bias
	for i := 0; i < LOW_STAGE_WIDTH; i++ {
		t := math32.Float32frombits(uint32(p.r[i]) | uint32(p.g[i])<<16)

		rf := t*factor.R + bias.R
		gf := t*factor.G + bias.G
		bf := t*factor.B + bias.B
		af := t*factor.A + bias.A

		p.r[i] = uint16(rf*255.0 + 0.5)
		p.g[i] = uint16(gf*255.0 + 0.5)
		p.b[i] = uint16(bf*255.0 + 0.5)
		p.a[i] = uint16(af*255.0 + 0.5)
	}
}

//go:fix inline
func (p *LowPipeline) XYToUnitAngle() {
}

//go:fix inline
func (p *LowPipeline) XYToRadius() {
	// Join r and g into x coordinates
	var x [LOW_STAGE_WIDTH]float32
	for i := 0; i < LOW_STAGE_WIDTH; i++ {
		x[i] = math32.Float32frombits(uint32(p.r[i]) | uint32(p.g[i])<<16)
	}

	// Join b and a into y coordinates
	var y [LOW_STAGE_WIDTH]float32
	for i := 0; i < LOW_STAGE_WIDTH; i++ {
		y[i] = math32.Float32frombits(uint32(p.b[i]) | uint32(p.a[i])<<16)
	}

	// Calculate radius: r = sqrt(x^2 + y^2)
	for i := 0; i < LOW_STAGE_WIDTH; i++ {
		x[i] = math32.Sqrt(x[i]*x[i] + y[i]*y[i])
	}

	// Split x back to r,g and y back to b,a (matching Rust implementation)
	for i := 0; i < LOW_STAGE_WIDTH; i++ {
		// Split x radius back to r, g
		xBits := math32.Float32bits(x[i])
		p.r[i] = uint16(xBits & 0xFFFF)
		p.g[i] = uint16(xBits >> 16)

		// Split y radius back to b, a
		yBits := math32.Float32bits(y[i])
		p.b[i] = uint16(yBits & 0xFFFF)
		p.a[i] = uint16(yBits >> 16)
	}
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
