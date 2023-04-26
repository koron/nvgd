package embedresource

type Option interface {
	apply(*EmbedResource)
}

type optionFunc func(*EmbedResource)

func (fn optionFunc) apply(res *EmbedResource) {
	fn(res)
}

func WithPrefix(prefix string) Option {
	return optionFunc(func(res *EmbedResource) {
		res.prefix = prefix
	})
}

func WithFallback(fallback string) Option {
	return optionFunc(func(res *EmbedResource) {
		res.fallback = fallback
	})
}
