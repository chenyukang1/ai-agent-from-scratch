// openai 适配层
package openai

import (
	"context"
	"log"
	"time"

	"github.com/chenyukang1/ai-agent-from-scratch/internal/components/chat"
	"github.com/chenyukang1/ai-agent-from-scratch/internal/schema"
	"github.com/sashabaranov/go-openai"
)

type Client struct {
	cli    *openai.Client
	config *ClientConfig
}

type ClientConfig struct {
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

func NewClient(conf *ClientConfig) *Client {
	config := openai.DefaultConfig(conf.APIKey)
	client := openai.NewClientWithConfig(config)
	return &Client{
		cli:    client,
		config: conf,
	}
}

func (c *Client) Generate(ctx context.Context, input []*schema.Message, opts ...chat.Option) (*schema.Message, error) {
	return nil, nil
}

func (c *Client) genRequest(ctx context.Context, input []*schema.Message, opts ...chat.Option) {
	options := chat.GetOptions(&chat.Options{
		Temperature: c.config.Temperature,
		Model:       &c.config.Model,
	}, opts...)
	req := &openai.ChatCompletionRequest{
		Model:                           *options.Model,
		MaxTokens:                       *c.config.MaxCompletionTokens,
		MaxCompletionTokens:             0,
		Temperature:                     0,
		TopP:                            0,
		N:                               0,
		Stream:                          false,
		Stop:                            []string{},
		PresencePenalty:                 0,
		ResponseFormat:                  &openai.ChatCompletionResponseFormat{},
		Seed:                            new(int),
		FrequencyPenalty:                0,
		LogitBias:                       map[string]int{},
		LogProbs:                        false,
		TopLogProbs:                     0,
		User:                            "",
		Functions:                       []openai.FunctionDefinition{},
		FunctionCall:                    c,
		Tools:                           []openai.Tool{},
		ToolChoice:                      c,
		StreamOptions:                   &openai.StreamOptions{},
		ParallelToolCalls:               c,
		Store:                           false,
		ReasoningEffort:                 "",
		Metadata:                        map[string]string{},
		Prediction:                      &openai.Prediction{},
		ChatTemplateKwargs:              map[string]any{},
		ServiceTier:                     "",
		Verbosity:                       "",
		SafetyIdentifier:                "",
		ChatCompletionRequestExtensions: openai.ChatCompletionRequestExtensions{},
	}
	log.Fatal(req)
}
