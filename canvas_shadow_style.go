package tinyskia

import "image/color"

// GetShadowOffsetX returns the shadow offset x.
func (ctx *Context) GetShadowOffsetX() float64 {
	return ctx.shadowOffsetX
}

// SetShadowOffsetX sets the shadow offset x.
func (ctx *Context) SetShadowOffsetX(x float64) {
	// TODO
	ctx.shadowOffsetX = x
}

// GetShadowOffsetY returns the shadow offset y.
func (ctx *Context) GetShadowOffsetY() float64 {
	return ctx.shadowOffsetY
}

// SetShadowOffsetY sets the shadow offset y.
func (ctx *Context) SetShadowOffsetY(y float64) {
	// TODO
	ctx.shadowOffsetY = y
}

// GetShadowBlur returns the shadow blur.
func (ctx *Context) GetShadowBlur() float64 {
	return ctx.shadowBlur
}

// SetShadowBlur sets the shadow blur.
func (ctx *Context) SetShadowBlur(blur float64) {
	// TODO
	ctx.shadowBlur = blur
}

// GetShadowColor returns the shadow color.
func (ctx *Context) GetShadowColor() color.Color {
	return ctx.shadowColor
}

// SetShadowColor sets the shadow color.
func (ctx *Context) SetShadowColor(color color.Color) {
	// TODO
	ctx.shadowColor = color
}
