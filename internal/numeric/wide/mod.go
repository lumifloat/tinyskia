package wide

import (
	"github.com/chewxy/math32"
)

func F32GenericBitBlend(mask, y, n float32) float32 {
	m := math32.Float32bits(mask)
	yBits := math32.Float32bits(y)
	nBits := math32.Float32bits(n)

	resBits := nBits ^ ((nBits ^ yBits) & m)

	return math32.Float32frombits(resBits)
}

func F64GenericBitBlend(mask, y, n float64) float64 {
	m := math32.Float64bits(mask)
	yBits := math32.Float64bits(y)
	nBits := math32.Float64bits(n)

	resBits := nBits ^ ((nBits ^ yBits) & m)

	return math32.Float64frombits(resBits)
}
