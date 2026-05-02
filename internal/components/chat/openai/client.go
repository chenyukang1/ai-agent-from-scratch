// openai 适配层
package openai

import (
	"context"
	"fmt"
	"time"

	"github.com/chenyukang1/ai-agent-from-scratch/internal/components/chat"
	"github.com/chenyukang1/ai-agent-from-scratch/internal/schema"
	"github.com/sashabaranov/go-openai"
)

// For testing purpose, we can use a fake chat creator to mock the API response.
type chatCompletionCreator interface {
	CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error)
}

type Client struct {
	cli    chatCompletionCreator
	config *ClientConfig
}

type ClientConfig struct {
	// APIKey is your authentication key
	// Use OpenAI API key or Azure API key depending on the service
	// Required
	APIKey string `json:"api_key"`

	// BseURL is the base URL for API requests
	// For OpenAI API, use https://api.openai.com/v1
	// For Azure OpenAI, use the endpoint URL of your Azure OpenAI resource, e.g. https://your-resource-name.openai.azure.com/
	// Optional. Default: https://api.openai.com/v1
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
}

func NewClient(conf *ClientConfig) *Client {
	config := openai.DefaultConfig(conf.APIKey)
	if conf.BaseURL != "" {
		config.BaseURL = conf.BaseURL
	}
	client := openai.NewClientWithConfig(config)
	return &Client{
		cli:    client,
		config: conf,
	}
}

func (c *Client) Generate(ctx context.Context, input []*schema.Message, opts ...chat.Option) (*schema.Message, error) {
	req, err := c.genRequest(ctx, input, opts...)
	if err != nil {
		return nil, err
	}
	resp, err := c.cli.CreateChatCompletion(ctx, *req)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat completion: %v", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned from openai API response")
	}

	message := resp.Choices[0].Message
	outMessage := &schema.Message{
		Role:    toSchemaRole(message.Role),
		Content: message.Content,
	}
	if len(message.ReasoningContent) > 0 {
		outMessage.ReasoningContent = message.ReasoningContent
	}

	return outMessage, nil
}

func (c *Client) genRequest(ctx context.Context, input []*schema.Message, opts ...chat.Option) (*openai.ChatCompletionRequest, error) {
	options := chat.GetOptions(&chat.Options{
		Temperature:         c.config.Temperature,
		MaxTokens:           c.config.MaxTokens,
		MaxCompletionTokens: c.config.MaxCompletionTokens,
		Model:               c.config.Model,
	}, opts...)
	req := &openai.ChatCompletionRequest{
		Model:               options.Model,
		MaxTokens:           options.MaxTokens,
		MaxCompletionTokens: options.MaxCompletionTokens,
		Temperature:         options.Temperature,
		Stream:              false,
		ResponseFormat:      &openai.ChatCompletionResponseFormat{},
	}

	msgs := make([]openai.ChatCompletionMessage, len(input))
	for i, msg := range input {
		msgs[i] = openai.ChatCompletionMessage{
			Role:    toOpenaiRole(msg.Role),
			Content: msg.Content,
		}
	}
	req.Messages = msgs
	return req, nil
}

func toOpenaiRole(role schema.RoleType) string {
	switch role {
	case schema.User:
		return openai.ChatMessageRoleUser
	case schema.Assistant:
		return openai.ChatMessageRoleAssistant
	case schema.System:
		return openai.ChatMessageRoleSystem
	case schema.Tool:
		return openai.ChatMessageRoleTool
	default:
		return string(role)
	}
}

func toSchemaRole(s string) schema.RoleType {
	switch s {
	case openai.ChatMessageRoleUser:
		return schema.User
	case openai.ChatMessageRoleAssistant:
		return schema.Assistant
	case openai.ChatMessageRoleSystem:
		return schema.System
	case openai.ChatMessageRoleTool:
		return schema.Tool
	default:
		panic(fmt.Sprintf("unimplemented role: %s", s))
	}
}
