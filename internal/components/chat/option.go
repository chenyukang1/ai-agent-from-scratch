package chat

type Options struct {
	// Temperature is the temperature for the model, which controls the randomness of the model.
	Temperature float32
	// MaxTokens is the max number of tokens, if reached the max tokens, the model will stop generating, and mostly return an finish reason of "length".
	MaxTokens int
	// MaxCompletionTokens specifies an upper bound for the number of tokens that can be generated for a completion, including visible output tokens and reasoning tokens.
	MaxCompletionTokens int `json:"max_completion_tokens,omitempty"`
	// Model is the model name.
	Model string
}

type Option struct {
	apply func(opts *Options)
}

func WithOption(apply func(opts *Options)) Option {
	return Option{apply: apply}
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
