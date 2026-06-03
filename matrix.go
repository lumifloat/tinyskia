package tinyskia

import (
	"github.com/lumifloat/tinyskia/internal/path"
)

type Matrix struct {
	transform path.Transform
}

// NewMatrix creates a new matrix.
func NewMatrix(a, b, c, d, e, f float64) *Matrix {
	transform := path.NewTransform(
		float32(a), float32(b), float32(c), float32(d), float32(e), float32(f))
	return &Matrix{transform: transform}
}

// NewMatrixIdentity creates a new identity matrix.
func NewMatrixIdentity() *Matrix {
	return &Matrix{transform: path.NewTransformDefault()}
}

// NewMatrixFromTranslate creates a new translating matrix.
func NewMatrixFromTranslate(x, y float64) *Matrix {
	tf := path.NewTransformFromTranslate(float32(x), float32(y))
	return &Matrix{transform: tf}
}

// NewMatrixFromScale creates a new scaling matrix.
func NewMatrixFromScale(x, y float64) *Matrix {
	tf := path.NewTransformFromScale(float32(x), float32(y))
	return &Matrix{transform: tf}
}

// NewMatrixFromSkew creates a new skewing matrix.
func NewMatrixFromSkew(x, y float64) *Matrix {
	tf := path.NewTransformFromSkew(float32(x), float32(y))
	return &Matrix{transform: tf}
}

// NewMatrixFromRotate creates a new rotating matrix.
func NewMatrixFromRotate(angle float64) *Matrix {
	tf := path.NewTransformFromRotate(float32(angle))
	return &Matrix{transform: tf}
}

// NewTransformFromRotateAt creates a new rotating matrix at the specified position.
func NewTransformFromRotateAt(angle, x, y float64) *Matrix {
	tf := path.NewTransformFromRotateAt(float32(angle), float32(x), float32(y))
	return &Matrix{transform: tf}
}

// Scale scales the current transform.
func (m *Matrix) Scale(x, y float64) *Matrix {
	return &Matrix{transform: m.transform.PostScale(float32(x), float32(y))}
}

// Translate translates the current transform.
func (m *Matrix) Translate(x, y float64) *Matrix {
	return &Matrix{transform: m.transform.PostTranslate(float32(x), float32(y))}
}

// Rotate rotates the current transform by the specified position.
func (m *Matrix) Rotate(angle float64) *Matrix {
	return &Matrix{transform: m.transform.PostRotate(float32(angle))}
}

// RotateAt rotates the current matrix by the specified angle about the specified point.
func (m *Matrix) RotateAt(angle, x, y float64) *Matrix {
	return &Matrix{transform: m.transform.PostRotateAt(float32(angle), float32(x), float32(y))}
}

// Multiply multiplies the current transform.
func (m *Matrix) Multiply(transform *Matrix) *Matrix {
	return &Matrix{transform: m.transform.PostConcat(transform.transform)}
}

// A returns the scale-x value of the matrix.
func (m *Matrix) A() float64 {
	return float64(m.transform.SX)
}

// B returns the skew-y value of the matrix.
func (m *Matrix) B() float64 {
	return float64(m.transform.KY)
}

// C returns the skew-x value of the matrix.
func (m *Matrix) C() float64 {
	return float64(m.transform.KX)
}

// D returns the scale-y value of the matrix.
func (m *Matrix) D() float64 {
	return float64(m.transform.SY)
}

// E returns the translate-x value of the matrix.
func (m *Matrix) E() float64 {
	return float64(m.transform.TX)
}

// F returns the translate-y value of the matrix.
func (m *Matrix) F() float64 {
	return float64(m.transform.TY)
}
