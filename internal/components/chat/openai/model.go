package openai

import (
	"context"
	"log"
	"time"

	"github.com/chenyukang1/ai-agent-from-scratch/internal/components/chat"
	"github.com/chenyukang1/ai-agent-from-scratch/internal/schema"
)

type ChatModelConfig struct {
	// APIKey is your authentication key
	// Use OpenAI API key or Azure API key depending on the service
	// Required
	APIKey string `json:"api_key"`

	// Timeout specifies the maximum duration to wait for API responses
	// If HTTPClient is set, Timeout will not be used.
	// Optional. Default: no timeout
	Timeout time.Duration `json:"timeout"`

	// Model specifies the ID of the model to use
	// Required
	Model string `json:"model"`

	// MaxCompletionTokens specifies an upper bound for the number of tokens that can be generated for a completion, including visible output tokens and reasoning tokens.
	MaxCompletionTokens *int `json:"max_completion_tokens,omitempty"`

	// Temperature specifies what sampling temperature to use
	// Generally recommend altering this or TopP but not both.
	// Range: 0.0 to 2.0. Higher values make output more random
	// Optional. Default: 1.0
	Temperature *float32 `json:"temperature,omitempty"`
}

type ChatModel struct {
	client *Client
}

func NewChatModel(conf *ChatModelConfig) *ChatModel {
	client := NewClient(&ClientConfig{
		APIKey:              conf.APIKey,
		Timeout:             conf.Timeout,
		Model:               conf.Model,
		MaxCompletionTokens: conf.MaxCompletionTokens,
		Temperature:         conf.Temperature,
	})
	return &ChatModel{client: client}
}

func (cm *ChatModel) Generate(ctx context.Context, input []*schema.Message, opts ...chat.Option) (*schema.Message, error) {
	output, err := cm.client.Generate(ctx, input, opts...)
	if err != nil {
		return nil, err
	}
	log.Fatal(output)
	return nil, nil
}
