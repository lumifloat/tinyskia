package tinyskia

type ImageSmoothingQuality string

const (
	ImageSmoothingQualityLow    ImageSmoothingQuality = "low"
	ImageSmoothingQualityMedium ImageSmoothingQuality = "medium"
	ImageSmoothingQualityHigh   ImageSmoothingQuality = "high"
)

// GetImageSmoothingEnabled eturns whether pattern fills and the drawImage() method will attempt to
// smooth images if their pixels don't line up exactly with the display, when scaling images up.
func (ctx *Context) GetImageSmoothingEnabled() bool {
	return ctx.imageSmoothingEnabled
}

// SetImageSmoothingEnabled changes whether images are smoothed (true) or not (false).
func (ctx *Context) SetImageSmoothingEnabled(enabled bool) {
	// TODO
	ctx.imageSmoothingEnabled = enabled
}

// GetImageSmoothingQuality returns the current image-smoothing-quality preference.
func (ctx *Context) GetImageSmoothingQuality() ImageSmoothingQuality {
	return ctx.imageSmoothingQuality
}

// SetImageSmoothingQuality changes the blur level. Values that are not finite numbers
// greater than or equal to zero are ignored.
func (ctx *Context) SetImageSmoothingQuality(quality ImageSmoothingQuality) {
	// TODO
	ctx.imageSmoothingQuality = quality
}
