package wide

import (
	"github.com/chewxy/math32"
)

type F32x2 [2]float32
type I32x2 [2]int32

func Splat(v float32) F32x2 {
	return F32x2{v, v}
}

func (v F32x2) Floor() F32x2 {
	trunc := v.TruncInt()
	roundtrip := F32x2{
		float32(trunc[0]),
		float32(trunc[1]),
	}

	mask := roundtrip.CmpGt(v)
	return roundtrip.Sub(mask.Blend(
		F32x2{1, 1},
		F32x2{0, 0},
	))
}

func (v F32x2) Abs() F32x2 {
	return F32x2{
		math32.Abs(v[0]),
		math32.Abs(v[1]),
	}
}

func (v F32x2) Max(rhs F32x2) F32x2 {
	return F32x2{
		math32.Max(v[0], rhs[0]),
		math32.Max(v[1], rhs[1]),
	}
}

func (v F32x2) Min(rhs F32x2) F32x2 {
	return F32x2{
		math32.Min(v[0], rhs[0]),
		math32.Min(v[1], rhs[1]),
	}
}

func (v F32x2) CmpEq(rhs F32x2) F32x2 {
	var res F32x2
	if v[0] == rhs[0] {
		res[0] = math32.MaxUint32
	} else {
		res[0] = 0.0
	}
	if v[1] == rhs[1] {
		res[1] = math32.MaxUint32
	} else {
		res[1] = 0.0
	}
	return res
}

func (v F32x2) CmpNe(rhs F32x2) F32x2 {
	var res F32x2
	if v[0] != rhs[0] {
		res[0] = math32.MaxUint32
	} else {
		res[0] = 0.0
	}
	if v[1] != rhs[1] {
		res[1] = math32.MaxUint32
	} else {
		res[1] = 0.0
	}
	return res
}

func (v F32x2) CmpGe(rhs F32x2) F32x2 {
	var res F32x2
	if v[0] >= rhs[0] {
		res[0] = math32.MaxUint32
	} else {
		res[0] = 0.0
	}
	if v[1] >= rhs[1] {
		res[1] = math32.MaxUint32
	} else {
		res[1] = 0.0
	}
	return res
}

func (v F32x2) CmpGt(rhs F32x2) F32x2 {
	var res F32x2
	if v[0] > rhs[0] {
		res[0] = math32.MaxUint32
	} else {
		res[0] = 0.0
	}
	if v[1] > rhs[1] {
		res[1] = math32.MaxUint32
	} else {
		res[1] = 0.0
	}
	return res
}

func (v F32x2) CmpLe(rhs F32x2) F32x2 {
	var res F32x2
	if v[0] <= rhs[0] {
		res[0] = math32.MaxUint32
	} else {
		res[0] = 0.0
	}
	if v[1] <= rhs[1] {
		res[1] = math32.MaxUint32
	} else {
		res[1] = 0.0
	}
	return res
}

func (v F32x2) CmpLt(rhs F32x2) F32x2 {
	var res F32x2
	if v[0] < rhs[0] {
		res[0] = math32.MaxUint32
	} else {
		res[0] = 0.0
	}
	if v[1] < rhs[1] {
		res[1] = math32.MaxUint32
	} else {
		res[1] = 0.0
	}
	return res
}

func (v F32x2) Blend(t, f F32x2) F32x2 {
	return F32x2{
		F32GenericBitBlend(v[0], t[0], f[0]),
		F32GenericBitBlend(v[1], t[1], f[1]),
	}
}

func (v F32x2) Round() F32x2 {
	return F32x2{
		math32.Round(v[0]),
		math32.Round(v[1]),
	}
}

func (v F32x2) RoundInt() I32x2 {
	rounded := v.Round()
	return I32x2{
		int32(rounded[0]),
		int32(rounded[1]),
	}
}

// TODO
func (v F32x2) TruncInt() I32x2 {
	return I32x2{
		int32(v[0]),
		int32(v[1]),
	}
}

func (v F32x2) RecipFast() F32x2 {
	return F32x2{
		1.0 / v[0],
		1.0 / v[1],
	}
}

func (v F32x2) RecipSqrt() F32x2 {
	return F32x2{
		1.0 / math32.Sqrt(v[0]),
		1.0 / math32.Sqrt(v[1]),
	}
}

func (v F32x2) Sqrt() F32x2 {
	return F32x2{
		math32.Sqrt(v[0]),
		math32.Sqrt(v[1]),
	}
}

func (v F32x2) Add(rhs F32x2) F32x2 {
	return F32x2{v[0] + rhs[0], v[1] + rhs[1]}
}

func (v F32x2) Sub(rhs F32x2) F32x2 {
	return F32x2{v[0] - rhs[0], v[1] - rhs[1]}
}

func (v F32x2) Mul(rhs F32x2) F32x2 {
	return F32x2{v[0] * rhs[0], v[1] * rhs[1]}
}

func (v F32x2) Div(rhs F32x2) F32x2 {
	return F32x2{v[0] / rhs[0], v[1] / rhs[1]}
}

func (v F32x2) BitAnd(rhs F32x2) F32x2 {
	return F32x2{
		math32.Float32frombits(math32.Float32bits(v[0]) & math32.Float32bits(rhs[0])),
		math32.Float32frombits(math32.Float32bits(v[1]) & math32.Float32bits(rhs[1])),
	}
}

func (v F32x2) BitOr(rhs F32x2) F32x2 {
	return F32x2{
		math32.Float32frombits(math32.Float32bits(v[0]) | math32.Float32bits(rhs[0])),
		math32.Float32frombits(math32.Float32bits(v[1]) | math32.Float32bits(rhs[1])),
	}
}

func (v F32x2) BitXor(rhs F32x2) F32x2 {
	return F32x2{
		math32.Float32frombits(math32.Float32bits(v[0]) ^ math32.Float32bits(rhs[0])),
		math32.Float32frombits(math32.Float32bits(v[1]) ^ math32.Float32bits(rhs[1])),
	}
}

func (v F32x2) Neg() F32x2 {
	return F32x2{
		-v[0],
		-v[1],
	}
}

func (v F32x2) Not() F32x2 {
	return F32x2{
		math32.Float32frombits(math32.Float32bits(v[0]) ^ math32.MaxUint8),
		math32.Float32frombits(math32.Float32bits(v[1]) ^ math32.MaxUint8),
	}
}

func (v F32x2) Eq(rhs F32x2) bool {
	return v[0] == rhs[0] && v[1] == rhs[1]
}
