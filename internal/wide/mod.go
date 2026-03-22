package wide

import "math"

func F32GenericBitBlend(mask, y, n float32) float32 {
	m := math.Float32bits(mask)
	yBits := math.Float32bits(y)
	nBits := math.Float32bits(n)

	resBits := nBits ^ ((nBits ^ yBits) & m)

	return math.Float32frombits(resBits)
}

func F64GenericBitBlend(mask, y, n float64) float64 {
	m := math.Float64bits(mask)
	yBits := math.Float64bits(y)
	nBits := math.Float64bits(n)

	resBits := nBits ^ ((nBits ^ yBits) & m)

	return math.Float64frombits(resBits)
}
