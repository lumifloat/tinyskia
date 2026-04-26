package tinyskia

import (
	"github.com/lumifloat/tinyskia/internal/path"
)

type matrix struct {
	transform path.Transform
}

// NewMatrix creates a new matrix.
func NewMatrix(a, b, c, d, e, f float64) *matrix {
	transform := path.NewTransform(
		float32(a), float32(b), float32(c), float32(d), float32(e), float32(f))
	return &matrix{transform: transform}
}

// NewMatrixIdentity creates a new identity matrix.
func NewMatrixIdentity() *matrix {
	return &matrix{transform: path.NewTransformDefault()}
}

// NewMatrixFromTranslate creates a new translating matrix.
func NewMatrixFromTranslate(x, y float64) *matrix {
	tf := path.NewTransformFromTranslate(float32(x), float32(y))
	return &matrix{transform: tf}
}

// NewMatrixFromScale creates a new scaling matrix.
func NewMatrixFromScale(x, y float64) *matrix {
	tf := path.NewTransformFromScale(float32(x), float32(y))
	return &matrix{transform: tf}
}

// NewMatrixFromSkew creates a new skewing matrix.
func NewMatrixFromSkew(x, y float64) *matrix {
	tf := path.NewTransformFromSkew(float32(x), float32(y))
	return &matrix{transform: tf}
}

// NewMatrixFromRotate creates a new rotating matrix.
func NewMatrixFromRotate(angle float64) *matrix {
	tf := path.NewTransformFromRotate(float32(angle))
	return &matrix{transform: tf}
}

// NewTransformFromRotateAt creates a new rotating matrix at the specified position.
func NewTransformFromRotateAt(angle, x, y float64) *matrix {
	tf := path.NewTransformFromRotateAt(float32(angle), float32(x), float32(y))
	return &matrix{transform: tf}
}

// Scale scales the current transform.
func (m *matrix) Scale(x, y float64) *matrix {
	return &matrix{transform: m.transform.PostScale(float32(x), float32(y))}
}

// Translate translates the current transform.
func (m *matrix) Translate(x, y float64) *matrix {
	return &matrix{transform: m.transform.PostTranslate(float32(x), float32(y))}
}

// Rotate rotates the current transform by the specified position.
func (m *matrix) Rotate(angle float64) *matrix {
	return &matrix{transform: m.transform.PostRotate(float32(angle))}
}

// RotateAt rotates the current matrix by the specified angle about the specified point.
func (m *matrix) RotateAt(angle, x, y float64) *matrix {
	return &matrix{transform: m.transform.PostRotateAt(float32(angle), float32(x), float32(y))}
}

// Multiply multiplies the current transform.
func (m *matrix) Multiply(transform *matrix) *matrix {
	return &matrix{transform: m.transform.PostConcat(transform.transform)}
}
