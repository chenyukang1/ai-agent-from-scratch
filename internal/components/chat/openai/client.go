// openai 适配层
package openai

import (
	"context"
	"fmt"
	"net/http"

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
	config *ChatModelConfig
}

func NewClient(conf *ChatModelConfig) *Client {
	config := openai.DefaultConfig(conf.APIKey)

	if conf.BaseURL != "" {
		config.BaseURL = conf.BaseURL
	}

	if conf.HTTPClient == nil {
		config.HTTPClient = http.DefaultClient
	} else {
		config.HTTPClient = conf.HTTPClient
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

	return toSchemaMessage(resp.Choices[0].Message), nil
}

func toSchemaMessage(message openai.ChatCompletionMessage) *schema.Message {
	outMessage := &schema.Message{
		Role:    toSchemaRole(message.Role),
		Content: message.Content,
	}
	if len(message.ReasoningContent) > 0 {
		outMessage.ReasoningContent = message.ReasoningContent
	}
	if len(message.ToolCalls) > 0 {
		outMessage.ToolCalls = make([]schema.ToolCall, len(message.ToolCalls))
		for i, tc := range message.ToolCalls {
			outMessage.ToolCalls[i] = schema.ToolCall{
				ID:        tc.ID,
				Name:      tc.Function.Name,
				Arguments: tc.Function.Arguments,
			}
		}
	}
	return outMessage
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
	}

	msgs := make([]openai.ChatCompletionMessage, len(input))
	for i, msg := range input {
		msgs[i] = toOpenaiMessage(msg)
	}
	req.Messages = msgs

	if len(options.Tools) > 0 {
		req.Tools = make([]openai.Tool, len(options.Tools))
		for i, info := range options.Tools {
			if info == nil {
				continue
			}
			params := info.Parameters
			if params == nil {
				params = map[string]any{"type": "object", "properties": map[string]any{}}
			}
			req.Tools[i] = openai.Tool{
				Type: openai.ToolTypeFunction,
				Function: &openai.FunctionDefinition{
					Name:        info.Name,
					Description: info.Desc,
					Parameters:  params,
				},
			}
		}
	}

	return req, nil
}

func toOpenaiMessage(msg *schema.Message) openai.ChatCompletionMessage {
	if msg == nil {
		return openai.ChatCompletionMessage{}
	}
	out := openai.ChatCompletionMessage{
		Role:    toOpenaiRole(msg.Role),
		Content: msg.Content,
	}
	if msg.ToolCallID != "" {
		out.ToolCallID = msg.ToolCallID
	}
	if len(msg.ToolCalls) > 0 {
		out.ToolCalls = make([]openai.ToolCall, len(msg.ToolCalls))
		for i, tc := range msg.ToolCalls {
			out.ToolCalls[i] = openai.ToolCall{
				ID:   tc.ID,
				Type: openai.ToolTypeFunction,
				Function: openai.FunctionCall{
					Name:      tc.Name,
					Arguments: tc.Arguments,
				},
			}
		}
	}
	return out
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
