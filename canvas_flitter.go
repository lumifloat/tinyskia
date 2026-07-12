package tinyskia

type Filter interface {
	// TODO
	filter()
}

// GetFilters
func (ctx *Context) GetFilters() []Filter {
	return ctx.filters
}

// SetFilters
func (ctx *Context) SetFilters(filters []Filter) {
	// TODO
	ctx.filters = filters
}
