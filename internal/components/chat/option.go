package chat

import "github.com/chenyukang1/ai-agent-from-scratch/internal/schema"

type Options struct {
	// MaxTokens is the max number of tokens, if reached the max tokens, the model will stop generating, and mostly return an finish reason of "length".
	MaxTokens int
	// MaxCompletionTokens specifies an upper bound for the number of tokens that can be generated for a completion, including visible output tokens and reasoning tokens.
	MaxCompletionTokens int `json:"max_completion_tokens,omitempty"`
	// Temperature is the temperature for the model, which controls the randomness of the model.
	Temperature float32
	// Model is the model name.
	Model string
	// Tools are tool definitions passed to the model for function calling.
	Tools []*schema.ToolInfo
}

type Option struct {
	apply func(opts *Options)
}

func WithOption(apply func(opts *Options)) Option {
	return Option{apply: apply}
}

// WithTools attaches tool definitions to a single Generate call.
func WithTools(tools []*schema.ToolInfo) Option {
	return WithOption(func(opts *Options) {
		opts.Tools = tools
	})
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
