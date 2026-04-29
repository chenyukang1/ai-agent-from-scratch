package chat

type Options struct {
	// Temperature is the temperature for the model, which controls the randomness of the model.
	Temperature *float32
	// Model is the model name.
	Model *string
}

type Option struct {
	apply func(opts *Options)
}

func GetOptions(base *Options, opts ...Option) *Options {
	if base == nil {
		return &Options{}
	}

	for _, opt := range opts {
		if opt.apply != nil {
			opt.apply(base)
		}
	}

	return base
}
