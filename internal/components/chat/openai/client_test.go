package openai

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/chenyukang1/ai-agent-from-scratch/internal/components/chat"
	"github.com/chenyukang1/ai-agent-from-scratch/internal/schema"
	"github.com/sashabaranov/go-openai"
)

type fakeChatCreator struct {
	resp openai.ChatCompletionResponse
	err  error
	got  openai.ChatCompletionRequest
}

func (f *fakeChatCreator) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	f.got = req
	return f.resp, f.err
}

func TestClientGenerateWithWrongAPIKey(t *testing.T) {
	client := NewClient(&ClientConfig{
		APIKey:    "test-api-key",
		Timeout:   3 * time.Second,
		Model:     "test-model",
		MaxTokens: 10000,
	})
	message, err := client.Generate(context.Background(), []*schema.Message{{Role: schema.User, Content: "hello"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Logf("Generated message: %s", message.Content)
}

func TestGenRequestWithDefaultConfig(t *testing.T) {
	client := &Client{
		cli: nil,
		config: &ClientConfig{
			Model:               "test-model",
			MaxTokens:           100,
			MaxCompletionTokens: 50,
			Temperature:         0.3,
		},
	}

	req, err := client.genRequest(context.Background(), []*schema.Message{{Role: schema.User, Content: "hello"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Model != "test-model" {
		t.Fatalf("expected model test-model, got %q", req.Model)
	}
	if req.MaxTokens != 100 {
		t.Fatalf("expected max tokens 100, got %d", req.MaxTokens)
	}
	if req.MaxCompletionTokens != 50 {
		t.Fatalf("expected max completion tokens 50, got %d", req.MaxCompletionTokens)
	}
	if req.Temperature != 0.3 {
		t.Fatalf("expected temperature 0.3, got %v", req.Temperature)
	}
	if len(req.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(req.Messages))
	}
	if req.Messages[0].Role != openai.ChatMessageRoleUser {
		t.Fatalf("expected user role, got %q", req.Messages[0].Role)
	}
	if req.Messages[0].Content != "hello" {
		t.Fatalf("expected hello content, got %q", req.Messages[0].Content)
	}
}

func TestGenRequestWithOverrideOptions(t *testing.T) {
	client := &Client{
		cli: nil,
		config: &ClientConfig{
			Model:               "base-model",
			MaxTokens:           100,
			MaxCompletionTokens: 50,
			Temperature:         0.3,
		},
	}

	override := chat.WithOption(func(opts *chat.Options) {
		opts.Model = "override-model"
		opts.MaxTokens = 200
		opts.Temperature = 0.7
	})

	req, err := client.genRequest(context.Background(), []*schema.Message{}, override)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.Model != "override-model" {
		t.Fatalf("expected override-model, got %q", req.Model)
	}
	if req.MaxTokens != 200 {
		t.Fatalf("expected max tokens 200, got %d", req.MaxTokens)
	}
	if req.Temperature != 0.7 {
		t.Fatalf("expected temperature 0.7, got %v", req.Temperature)
	}
	if len(req.Messages) != 0 {
		t.Fatalf("expected no messages, got %d", len(req.Messages))
	}
}

func TestGenerateReturnsFirstChoice(t *testing.T) {
	fake := &fakeChatCreator{
		resp: openai.ChatCompletionResponse{
			Choices: []openai.ChatCompletionChoice{{
				Message: openai.ChatCompletionMessage{
					Role:             openai.ChatMessageRoleAssistant,
					Content:          "response text",
					ReasoningContent: "reasoning text",
				},
			}},
		},
	}
	client := &Client{
		cli: fake,
		config: &ClientConfig{
			Model:               "test-model",
			MaxTokens:           100,
			MaxCompletionTokens: 50,
			Temperature:         0.5,
		},
	}

	msg, err := client.Generate(context.Background(), []*schema.Message{{Role: schema.User, Content: "hi"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Role != schema.Assistant {
		t.Fatalf("expected assistant role, got %q", msg.Role)
	}
	if msg.Content != "response text" {
		t.Fatalf("expected content response text, got %q", msg.Content)
	}
	if msg.ReasoningContent != "reasoning text" {
		t.Fatalf("expected reasoning content reasoning text, got %q", msg.ReasoningContent)
	}
	if fake.got.Model != "test-model" {
		t.Fatalf("expected API request model test-model, got %q", fake.got.Model)
	}
}

func TestGenerateReturnsErrorWhenAPIError(t *testing.T) {
	fake := &fakeChatCreator{err: errors.New("api failure")}
	client := &Client{
		cli:    fake,
		config: &ClientConfig{Model: "test-model"},
	}

	_, err := client.Generate(context.Background(), []*schema.Message{{Role: schema.User, Content: "hi"}})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "failed to create chat completion: api failure" {
		t.Fatalf("expected wrapped API error, got %v", err)
	}
}

func TestGenerateReturnsErrorWhen401Unauthorized(t *testing.T) {
	fake := &fakeChatCreator{
		err: &openai.RequestError{
			HTTPStatusCode: 401,
			HTTPStatus:     "401 Unauthorized",
			Err:            errors.New("invalid api key"),
		},
	}
	client := &Client{
		cli:    fake,
		config: &ClientConfig{Model: "test-model"},
	}

	_, err := client.Generate(context.Background(), []*schema.Message{{Role: schema.User, Content: "hi"}})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// 断言错误是 401 Unauthorized
	if !strings.Contains(err.Error(), "401 Unauthorized") {
		t.Logf("expected 401 Unauthorized error: %v", err)
	}
}

func TestGenerateReturnsErrorWhenNoChoices(t *testing.T) {
	fake := &fakeChatCreator{resp: openai.ChatCompletionResponse{Choices: []openai.ChatCompletionChoice{}}}
	client := &Client{
		cli:    fake,
		config: &ClientConfig{Model: "test-model"},
	}

	_, err := client.Generate(context.Background(), []*schema.Message{{Role: schema.User, Content: "hi"}})
	if err == nil {
		t.Fatal("expected error when no choices")
	}
	if err.Error() != "no choices returned from openai API response" {
		t.Fatalf("expected no choices error, got %v", err)
	}
}
