package openai

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/chenyukang1/ai-agent-from-scratch/internal/components/chat"
	"github.com/chenyukang1/ai-agent-from-scratch/internal/schema"
)

type ChatModelConfig struct {
	// APIKey is your authentication key
	// Use OpenAI API key or Azure API key depending on the service
	// Required
	APIKey string `json:"api_key"`

	// BaseURL is the OpenAI endpoint URL
	// Format: https://{YOUR_RESOURCE_NAME}.openai.azure.com. YOUR_RESOURCE_NAME is the name of your resource that you have created on Azure.
	// Required for proxy
	BaseURL string `json:"base_url"`

	// Timeout specifies the maximum duration to wait for API responses
	// If HTTPClient is set, Timeout will not be used.
	// Optional. Default: no timeout
	Timeout time.Duration `json:"timeout"`

	// Model specifies the ID of the model to use
	// Required
	Model string `json:"model"`

	// MaxTokens limits the maximum number of tokens that can be generated in the chat completion
	// Optional. Default: model's maximum
	// Deprecated: use MaxCompletionTokens. Not compatible with o1-series models.
	// refs: https://platform.openai.com/docs/api-reference/chat/create#chat-create-max_tokens
	MaxTokens int `json:"max_tokens,omitempty"`

	// MaxCompletionTokens specifies an upper bound for the number of tokens that can be generated for a completion, including visible output tokens and reasoning tokens.
	MaxCompletionTokens int `json:"max_completion_tokens,omitempty"`

	// Temperature specifies what sampling temperature to use
	// Generally recommend altering this or TopP but not both.
	// Range: 0.0 to 2.0. Higher values make output more random
	// Optional. Default: 1.0
	Temperature float32 `json:"temperature,omitempty"`

	// HTTPClient is used to send HTTP requests
	// Optional. Default: http.DefaultClient
	HTTPClient *http.Client `json:"-"`
}

type ChatModel struct {
	client *Client
}

func NewChatModel(conf *ChatModelConfig) *ChatModel {
	client := NewClient(conf)
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
